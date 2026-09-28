/*
Copyright (c) 2026 Dell Inc. or its subsidiaries. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package apiv27

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	v2 "github.com/dell/gopowerscale/api/v2"
	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExportCreateWithZoneAndXprtsec(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	tests := []struct {
		name        string
		export      *v2.Export
		zone        string
		xprtsec     string
		setupMock   func()
		expectError bool
		errorMsg    string
	}{
		{
			name:    "Success - mTLS export",
			export:  &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:    "System",
			xprtsec: "mtls",
			setupMock: func() {
				client.On("Post", anyArgs...).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "Success - TLS export",
			export:  &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:    "System",
			xprtsec: "tls",
			setupMock: func() {
				client.On("Post", anyArgs...).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "Success - TLS:mTLS export",
			export:  &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:    "System",
			xprtsec: "tls:mtls",
			setupMock: func() {
				client.On("Post", anyArgs...).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:    "Success - empty xprtsec omits the field",
			export:  &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:    "System",
			xprtsec: "",
			setupMock: func() {
				client.On("Post", anyArgs...).Return(nil).Once()
			},
			expectError: false,
		},
		{
			name:        "Error - empty zone",
			export:      &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:        "",
			xprtsec:     "mtls",
			setupMock:   func() {},
			expectError: true,
			errorMsg:    "zone cannot be empty",
		},
		{
			name:        "Error - empty paths",
			export:      &v2.Export{Paths: &[]string{}},
			zone:        "System",
			xprtsec:     "mtls",
			setupMock:   func() {},
			expectError: true,
			errorMsg:    "no path set",
		},
		{
			name:        "Error - nil paths",
			export:      &v2.Export{Paths: nil},
			zone:        "System",
			xprtsec:     "mtls",
			setupMock:   func() {},
			expectError: true,
			errorMsg:    "no path set",
		},
		{
			name:        "Error - nil export",
			export:      nil,
			zone:        "System",
			xprtsec:     "mtls",
			setupMock:   func() {},
			expectError: true,
			errorMsg:    "export cannot be nil",
		},
		{
			name:    "Error - API failure",
			export:  &v2.Export{Paths: &[]string{"/ifs/data/vol1"}},
			zone:    "System",
			xprtsec: "mtls",
			setupMock: func() {
				client.On("Post", anyArgs...).Return(errors.New("API error")).Once()
			},
			expectError: true,
			errorMsg:    "API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.ExpectedCalls = nil
			tt.setupMock()

			exportID, err := ExportCreateWithZoneAndXprtsec(ctx, client, tt.export, tt.zone, tt.xprtsec)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
				assert.Equal(t, 0, exportID)
			} else {
				assert.NoError(t, err)
			}

			client.AssertExpectations(t)
		})
	}
}

func TestExportCreateWithZoneAndXprtsec_DualPathArchitecture(t *testing.T) {
	// This test documents the dual-path architecture:
	// - v27 API is used when xprtsec is specified
	// - v2 API is used when xprtsec is empty (handled by caller)

	t.Run("v27 API endpoint verification", func(t *testing.T) {
		ctx := context.Background()
		client := &mocks.Client{}

		export := &v2.Export{Paths: &[]string{"/ifs/data/vol1"}}

		// Verify that v27 API endpoint is used and the request carries xprtsec
		client.On("Post", ctx, "platform/27/protocols/nfs/exports", mock.Anything, mock.Anything,
			mock.Anything, mock.Anything, mock.Anything,
		).Return(nil).Run(func(args mock.Arguments) {
			req := args.Get(5).(*ExportReq)
			assert.NotNil(t, req.Xprtsec)
			assert.Equal(t, "mtls", *req.Xprtsec)
		}).Once()

		_, err := ExportCreateWithZoneAndXprtsec(ctx, client, export, "System", "mtls")
		assert.NoError(t, err)

		client.AssertExpectations(t)
	})
}

// TestExportReqJSON verifies that the v27 request marshals to a single flat
// object: the embedded v2 fields are promoted alongside xprtsec, and xprtsec is
// omitted entirely when empty.
func TestExportReqJSON(t *testing.T) {
	export := &v2.Export{
		Paths:       &[]string{"/ifs/data/vol1"},
		Description: "k8s volume",
	}

	// With xprtsec set
	body, err := json.Marshal(newExportReq(export, "mtls"))
	assert.NoError(t, err)

	var decoded map[string]interface{}
	assert.NoError(t, json.Unmarshal(body, &decoded))
	assert.Equal(t, "mtls", decoded["xprtsec"])
	assert.Equal(t, "k8s volume", decoded["description"])
	assert.Equal(t, []interface{}{"/ifs/data/vol1"}, decoded["paths"])

	// Without xprtsec the key must be absent, so older endpoints never see it
	body, err = json.Marshal(newExportReq(export, ""))
	assert.NoError(t, err)

	decoded = map[string]interface{}{}
	assert.NoError(t, json.Unmarshal(body, &decoded))
	_, present := decoded["xprtsec"]
	assert.False(t, present, "xprtsec must be omitted when empty")
	assert.Equal(t, "k8s volume", decoded["description"])
}

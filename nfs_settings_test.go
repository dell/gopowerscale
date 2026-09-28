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

package gopowerscale

import (
	"errors"
	"testing"

	"github.com/dell/gopowerscale/mocks"
	"github.com/dell/gopowerscale/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetNfsSettingsGlobal(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Successfully retrieve cluster-wide NFS TLS settings
	expected := &openapi.V27NfsSettingsGlobal{
		NfsTLSMode:        openapi.PtrString("none:tls:mtls"),
		NfsTLSMinVersion:  openapi.PtrString("1.3"),
		NfsTLSCiphers:     openapi.PtrString(""),
		NfsTLSOcspMode:    openapi.PtrString("none"),
		NfsTLSOcspAddress: openapi.PtrString(""),
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V27NfsSettingsGlobalResponse)
		*resp = openapi.V27NfsSettingsGlobalResponse{Settings: expected}
	}).Once()
	settings, err := client.GetNfsSettingsGlobal(defaultCtx)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, "none:tls:mtls", *settings.NfsTLSMode)
	assert.Equal(t, "1.3", *settings.NfsTLSMinVersion)

	// Test case: mTLS-only cluster configuration
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V27NfsSettingsGlobalResponse)
		*resp = openapi.V27NfsSettingsGlobalResponse{
			Settings: &openapi.V27NfsSettingsGlobal{NfsTLSMode: openapi.PtrString("mtls")},
		}
	}).Once()
	settings, err = client.GetNfsSettingsGlobal(defaultCtx)
	assert.NoError(t, err)
	assert.Equal(t, "mtls", *settings.NfsTLSMode)

	// Test case: API error (e.g. OneFS < 9.16.0 returns 404)
	testErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(testErr).Once()
	settings, err = client.GetNfsSettingsGlobal(defaultCtx)
	assert.ErrorIs(t, err, testErr)
	assert.Nil(t, settings)

	// Test case: Empty settings in an otherwise successful response
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V27NfsSettingsGlobalResponse)
		*resp = openapi.V27NfsSettingsGlobalResponse{Settings: nil}
	}).Once()
	settings, err = client.GetNfsSettingsGlobal(defaultCtx)
	assert.Error(t, err)
	assert.Nil(t, settings)
}

func TestUpdateNfsSettingsGlobal(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Successfully update cluster-wide TLS mode
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.UpdateNfsSettingsGlobal(defaultCtx, &openapi.V27NfsSettingsGlobal{
		NfsTLSMode:       openapi.PtrString("tls:mtls"),
		NfsTLSMinVersion: openapi.PtrString("1.3"),
	})
	assert.NoError(t, err)

	// Test case: Nil settings are rejected before any API call is made
	err = client.UpdateNfsSettingsGlobal(defaultCtx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "settings cannot be nil")

	// Test case: API error is propagated
	testErr := errors.New("test error")
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(testErr).Once()
	err = client.UpdateNfsSettingsGlobal(defaultCtx, &openapi.V27NfsSettingsGlobal{
		NfsTLSMode: openapi.PtrString("mtls"),
	})
	assert.ErrorIs(t, err, testErr)
}

func TestGetNfsSettingsExport(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Successfully retrieve site-wide default xprtsec
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V27NfsSettingsExportResponse)
		*resp = openapi.V27NfsSettingsExportResponse{
			Settings: &openapi.V27NfsSettingsExport{Xprtsec: openapi.PtrString("none:tls:mtls")},
		}
	}).Once()
	settings, err := client.GetNfsSettingsExport(defaultCtx)
	assert.NoError(t, err)
	assert.NotNil(t, settings)
	assert.Equal(t, "none:tls:mtls", *settings.Xprtsec)

	// Test case: API error is propagated
	testErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(testErr).Once()
	settings, err = client.GetNfsSettingsExport(defaultCtx)
	assert.ErrorIs(t, err, testErr)
	assert.Nil(t, settings)

	// Test case: Empty settings in an otherwise successful response
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V27NfsSettingsExportResponse)
		*resp = openapi.V27NfsSettingsExportResponse{Settings: nil}
	}).Once()
	settings, err = client.GetNfsSettingsExport(defaultCtx)
	assert.Error(t, err)
	assert.Nil(t, settings)
}

func TestUpdateNfsSettingsExport(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Successfully update the site-wide default xprtsec
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.UpdateNfsSettingsExport(defaultCtx, &openapi.V27NfsSettingsExport{
		Xprtsec: openapi.PtrString("tls:mtls"),
	})
	assert.NoError(t, err)

	// Test case: Nil settings are rejected before any API call is made
	err = client.UpdateNfsSettingsExport(defaultCtx, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "settings cannot be nil")

	// Test case: API error is propagated (e.g. empty intersection with cluster mode)
	testErr := errors.New("test error")
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(testErr).Once()
	err = client.UpdateNfsSettingsExport(defaultCtx, &openapi.V27NfsSettingsExport{
		Xprtsec: openapi.PtrString("mtls"),
	})
	assert.ErrorIs(t, err, testErr)
}

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
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dell/gopowerscale/api"
	"github.com/dell/gopowerscale/mocks"
	"github.com/dell/gopowerscale/openapi"
	"github.com/stretchr/testify/assert"
)

func TestGetNfsSettingsGlobal(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   interface{}
		responseStatus int
		wantErr        bool
		wantNilResult  bool
	}{
		{
			name: "Success - Get global NFS TLS settings",
			responseBody: openapi.V27NfsSettingsGlobalResponse{
				Settings: &openapi.V27NfsSettingsGlobal{
					NfsTLSMode:        stringPtr("none:tls:mtls"),
					NfsTLSMinVersion:  stringPtr("1.3"),
					NfsTLSCiphers:     stringPtr(""),
					NfsTLSOcspMode:    stringPtr("none"),
					NfsTLSOcspAddress: stringPtr(""),
				},
			},
			responseStatus: http.StatusOK,
			wantErr:        false,
			wantNilResult:  false,
		},
		{
			name: "Success - mTLS only configuration",
			responseBody: openapi.V27NfsSettingsGlobalResponse{
				Settings: &openapi.V27NfsSettingsGlobal{
					NfsTLSMode:       stringPtr("mtls"),
					NfsTLSMinVersion: stringPtr("1.3"),
				},
			},
			responseStatus: http.StatusOK,
			wantErr:        false,
			wantNilResult:  false,
		},
		{
			name:           "Error - 404 Not Found (OneFS < 9.16.0)",
			responseBody:   map[string]string{"error": "Not Found"},
			responseStatus: http.StatusNotFound,
			wantErr:        true,
			wantNilResult:  true,
		},
		{
			name:           "Error - 403 Forbidden (insufficient privileges)",
			responseBody:   map[string]string{"error": "Forbidden"},
			responseStatus: http.StatusForbidden,
			wantErr:        true,
			wantNilResult:  true,
		},
		{
			name: "Error - Empty settings response",
			responseBody: openapi.V27NfsSettingsGlobalResponse{
				Settings: nil,
			},
			responseStatus: http.StatusOK,
			wantErr:        true,
			wantNilResult:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/platform/latest/" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"latest": "9.16"})
					return
				}
				assert.Equal(t, "/"+nfsSettingsGlobalPath+"/", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create API client
			client, _ := api.New(context.Background(), server.URL, "user", "pass", "", 0, 0, nil)

			// Call function
			result, err := GetNfsSettingsGlobal(context.Background(), client)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNilResult {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUpdateNfsSettingsGlobal(t *testing.T) {
	tests := []struct {
		name           string
		settings       *openapi.V27NfsSettingsGlobal
		responseStatus int
		wantErr        bool
	}{
		{
			name: "Success - Update TLS mode to mtls only",
			settings: &openapi.V27NfsSettingsGlobal{
				NfsTLSMode:       stringPtr("mtls"),
				NfsTLSMinVersion: stringPtr("1.3"),
			},
			responseStatus: http.StatusNoContent,
			wantErr:        false,
		},
		{
			name: "Success - Update with cipher list",
			settings: &openapi.V27NfsSettingsGlobal{
				NfsTLSMode:       stringPtr("tls:mtls"),
				NfsTLSMinVersion: stringPtr("1.3"),
				NfsTLSCiphers:    stringPtr("HIGH:!aNULL:!MD5"),
			},
			responseStatus: http.StatusNoContent,
			wantErr:        false,
		},
		{
			name:           "Error - Nil settings",
			settings:       nil,
			responseStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Error - 400 Bad Request (invalid TLS mode)",
			settings: &openapi.V27NfsSettingsGlobal{
				NfsTLSMode: stringPtr("invalid"),
			},
			responseStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Error - 403 Forbidden (insufficient privileges)",
			settings: &openapi.V27NfsSettingsGlobal{
				NfsTLSMode: stringPtr("mtls"),
			},
			responseStatus: http.StatusForbidden,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.settings == nil {
				// The nil-settings guard must reject the request before the
				// client is ever used. Use a real mock with no expectations
				// so the test fails loudly if a call is attempted.
				mockClient := &mocks.Client{}
				err := UpdateNfsSettingsGlobal(context.Background(), mockClient, tt.settings)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "settings cannot be nil")
				mockClient.AssertNotCalled(t, "Put", anyArgs...)
				mockClient.AssertExpectations(t)
				return
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/platform/latest/" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"latest": "9.16"})
					return
				}
				assert.Equal(t, "/"+nfsSettingsGlobalPath+"/", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)

				w.WriteHeader(tt.responseStatus)
			}))
			defer server.Close()

			// Create API client
			client, _ := api.New(context.Background(), server.URL, "user", "pass", "", 0, 0, nil)

			// Call function
			err := UpdateNfsSettingsGlobal(context.Background(), client, tt.settings)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetNfsSettingsExport(t *testing.T) {
	tests := []struct {
		name           string
		responseBody   interface{}
		responseStatus int
		wantErr        bool
		wantNilResult  bool
	}{
		{
			name: "Success - Get export default settings",
			responseBody: openapi.V27NfsSettingsExportResponse{
				Settings: &openapi.V27NfsSettingsExport{
					Xprtsec: stringPtr("none:tls:mtls"),
				},
			},
			responseStatus: http.StatusOK,
			wantErr:        false,
			wantNilResult:  false,
		},
		{
			name: "Success - mTLS only default",
			responseBody: openapi.V27NfsSettingsExportResponse{
				Settings: &openapi.V27NfsSettingsExport{
					Xprtsec: stringPtr("mtls"),
				},
			},
			responseStatus: http.StatusOK,
			wantErr:        false,
			wantNilResult:  false,
		},
		{
			name:           "Error - 404 Not Found",
			responseBody:   map[string]string{"error": "Not Found"},
			responseStatus: http.StatusNotFound,
			wantErr:        true,
			wantNilResult:  true,
		},
		{
			name: "Error - Empty settings response",
			responseBody: openapi.V27NfsSettingsExportResponse{
				Settings: nil,
			},
			responseStatus: http.StatusOK,
			wantErr:        true,
			wantNilResult:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/platform/latest/" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"latest": "9.16"})
					return
				}
				assert.Equal(t, "/"+nfsSettingsExportPath+"/", r.URL.Path)
				assert.Equal(t, http.MethodGet, r.Method)

				w.WriteHeader(tt.responseStatus)
				json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create API client
			client, _ := api.New(context.Background(), server.URL, "user", "pass", "", 0, 0, nil)

			// Call function
			result, err := GetNfsSettingsExport(context.Background(), client)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.wantNilResult {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}
		})
	}
}

func TestUpdateNfsSettingsExport(t *testing.T) {
	tests := []struct {
		name           string
		settings       *openapi.V27NfsSettingsExport
		responseStatus int
		wantErr        bool
	}{
		{
			name: "Success - Update default xprtsec to mtls",
			settings: &openapi.V27NfsSettingsExport{
				Xprtsec: stringPtr("mtls"),
			},
			responseStatus: http.StatusNoContent,
			wantErr:        false,
		},
		{
			name: "Success - Update default xprtsec to tls:mtls",
			settings: &openapi.V27NfsSettingsExport{
				Xprtsec: stringPtr("tls:mtls"),
			},
			responseStatus: http.StatusNoContent,
			wantErr:        false,
		},
		{
			name:           "Error - Nil settings",
			settings:       nil,
			responseStatus: http.StatusBadRequest,
			wantErr:        true,
		},
		{
			name: "Error - 400 Bad Request (empty intersection)",
			settings: &openapi.V27NfsSettingsExport{
				Xprtsec: stringPtr("mtls"),
			},
			responseStatus: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.settings == nil {
				// The nil-settings guard must reject the request before the
				// client is ever used. Use a real mock with no expectations
				// so the test fails loudly if a call is attempted.
				mockClient := &mocks.Client{}
				err := UpdateNfsSettingsExport(context.Background(), mockClient, tt.settings)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "settings cannot be nil")
				mockClient.AssertNotCalled(t, "Put", anyArgs...)
				mockClient.AssertExpectations(t)
				return
			}

			// Create test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/platform/latest/" {
					w.WriteHeader(http.StatusOK)
					json.NewEncoder(w).Encode(map[string]string{"latest": "9.16"})
					return
				}
				assert.Equal(t, "/"+nfsSettingsExportPath+"/", r.URL.Path)
				assert.Equal(t, http.MethodPut, r.Method)

				w.WriteHeader(tt.responseStatus)
			}))
			defer server.Close()

			// Create API client
			client, _ := api.New(context.Background(), server.URL, "user", "pass", "", 0, 0, nil)

			// Call function
			err := UpdateNfsSettingsExport(context.Background(), client, tt.settings)

			// Assertions
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

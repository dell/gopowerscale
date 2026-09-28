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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dell/gopowerscale/api"
	"github.com/stretchr/testify/assert"
)

// capturedRequest records what actually went out on the wire.
type capturedRequest struct {
	method string
	path   string
	body   string
}

// newCompatTestClient starts a fake OneFS that reports the given PAPI version
// from /platform/latest and records every subsequent request.
func newCompatTestClient(t *testing.T, onefsVersion string) (*Client, *[]capturedRequest, func()) {
	t.Helper()

	captured := &[]capturedRequest{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/platform/latest") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]string{"latest": onefsVersion})
			return
		}

		raw, _ := io.ReadAll(r.Body)
		*captured = append(*captured, capturedRequest{
			method: r.Method,
			path:   r.URL.Path,
			body:   string(raw),
		})

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]int{"id": 42})
	}))

	apiClient, err := api.New(context.Background(), server.URL, "user", "pass", "", 0, 0, nil)
	assert.NoError(t, err)

	return &Client{apiClient}, captured, server.Close
}

// Scenario 1: OneFS older than 9.16.0 (no NFS-over-TLS support at all).
// Existing callers must keep hitting the v2 endpoint, and the request must not
// contain an xprtsec key that an older OneFS would reject as an unknown property.
func TestCompat_OneFSOlderThan9_16_NoMTLS(t *testing.T) {
	client, captured, closeFn := newCompatTestClient(t, "8.2")
	defer closeFn()

	// The pre-existing API that all current consumers call.
	_, err := client.ExportWithZoneAndPath(defaultCtx, "/ifs/data/vol1", "System", "legacy volume")
	assert.NoError(t, err)

	// The new API, invoked without opting in to transport security.
	_, err = client.ExportWithZoneAndPathAndXprtsec(defaultCtx, "/ifs/data/vol2", "System", "legacy volume", "")
	assert.NoError(t, err)

	assert.Len(t, *captured, 2)
	for _, req := range *captured {
		assert.Equal(t, http.MethodPost, req.method)
		assert.Contains(t, req.path, "/platform/2/protocols/nfs/exports",
			"pre-9.16 OneFS must only ever see the v2 endpoint")
		assert.NotContains(t, req.path, "/platform/27/",
			"the v27 endpoint does not exist on this OneFS and must not be called")
		assert.NotContains(t, req.body, "xprtsec",
			"xprtsec must never be sent to an OneFS that does not understand it")
	}
}

// Scenario 2: OneFS 9.16.0+ but the user is not using mTLS. The library must
// still take the legacy v2 path, so behaviour is unchanged for these users.
func TestCompat_OneFS9_16_WithoutMTLS_UsesLegacyPath(t *testing.T) {
	client, captured, closeFn := newCompatTestClient(t, "9.16")
	defer closeFn()

	_, err := client.ExportWithZoneAndPath(defaultCtx, "/ifs/data/vol1", "System", "non-mtls volume")
	assert.NoError(t, err)

	_, err = client.ExportWithZoneAndPathAndXprtsec(defaultCtx, "/ifs/data/vol2", "System", "non-mtls volume", "")
	assert.NoError(t, err)

	_, err = client.ExportVolumeWithZoneAndPathAndXprtsec(defaultCtx, "/ifs/data/vol3", "System", "non-mtls volume", "")
	assert.NoError(t, err)

	assert.Len(t, *captured, 3)
	for _, req := range *captured {
		assert.Contains(t, req.path, "/platform/2/protocols/nfs/exports",
			"opting out of xprtsec must keep using the v2 endpoint even on 9.16+")
		assert.NotContains(t, req.body, "xprtsec",
			"no xprtsec key should be sent when the feature is not used")
	}
}

// Scenario 3: mTLS/TLS explicitly requested on OneFS 9.16.0+. The library must
// switch to the v27 endpoint and send the xprtsec field verbatim.
func TestCompat_OneFS9_16_WithMTLS_UsesV27(t *testing.T) {
	for _, xprtsec := range []string{"mtls", "tls", "tls:mtls", "none:tls:mtls"} {
		t.Run(xprtsec, func(t *testing.T) {
			client, captured, closeFn := newCompatTestClient(t, "9.16")
			defer closeFn()

			id, err := client.ExportWithZoneAndPathAndXprtsec(
				defaultCtx, "/ifs/data/vol1", "System", "mtls volume", xprtsec)
			assert.NoError(t, err)
			assert.Equal(t, 42, id, "export id must be parsed from the v27 response")

			assert.Len(t, *captured, 1)
			req := (*captured)[0]
			assert.Equal(t, http.MethodPost, req.method)
			assert.Contains(t, req.path, "/platform/27/protocols/nfs/exports",
				"xprtsec requests must go to the v27 endpoint")

			var body map[string]interface{}
			assert.NoError(t, json.Unmarshal([]byte(req.body), &body))
			assert.Equal(t, xprtsec, body["xprtsec"], "xprtsec must be sent verbatim")
			// The embedded v2 fields must still be present and flat.
			assert.Equal(t, "mtls volume", body["description"])
			assert.Equal(t, []interface{}{"/ifs/data/vol1"}, body["paths"])
		})
	}
}

// The legacy v2 request body must be byte-for-byte what it was before xprtsec
// support existed. This guards against any silent drift in the legacy path.
func TestCompat_LegacyV2BodyUnchanged(t *testing.T) {
	client, captured, closeFn := newCompatTestClient(t, "9.16")
	defer closeFn()

	_, err := client.ExportWithZoneAndPath(defaultCtx, "/ifs/data/vol1", "System", "k8s volume")
	assert.NoError(t, err)

	assert.Len(t, *captured, 1)
	// Exactly the payload the pre-xprtsec implementation produced: paths and
	// description only, no xprtsec and no other new keys.
	assert.JSONEq(t,
		`{"paths":["/ifs/data/vol1"],"description":"k8s volume"}`,
		(*captured)[0].body,
	)
}

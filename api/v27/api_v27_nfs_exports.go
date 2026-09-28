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
	"errors"

	"github.com/dell/gopowerscale/api"
	apiv2 "github.com/dell/gopowerscale/api/v2"
)

const (
	nfsExportsPathV27 = "platform/27/protocols/nfs/exports"
)

// ExportReq is the request body for creating an NFS export via the v27 API.
// It embeds the v2 export request and adds the xprtsec transport-security
// field introduced in OneFS 9.16.0 (PAPI v27).
//
// xprtsec deliberately lives here rather than on the v2 request type so that
// it is structurally impossible to send the field to an older endpoint that
// does not recognise it. The embedded fields are promoted by encoding/json,
// so the wire format is a single flat object.
type ExportReq struct {
	apiv2.ExportReq
	// Xprtsec is the colon-delimited allowlist of transport security modes
	// (none, tls, mtls) this export permits. Omitted when empty, in which
	// case the cluster default applies.
	Xprtsec *string `json:"xprtsec,omitempty"`
}

// newExportReq builds a v27 export request from a v2 export definition plus an
// explicit xprtsec value. An empty xprtsec leaves the field unset.
func newExportReq(export *apiv2.Export, xprtsec string) *ExportReq {
	req := &ExportReq{ExportReq: *apiv2.NewExportReq(export)}
	if xprtsec != "" {
		req.Xprtsec = &xprtsec
	}
	return req
}

// ExportCreateWithZoneAndXprtsec creates an NFS export using the v27 API with xprtsec support.
// This function should be used when xprtsec parameter is specified (mTLS/TLS).
// Requires OneFS 9.16.0+ (PAPI v27).
//
// Parameters:
//   - ctx: Context for cancellation and timeouts
//   - client: API client
//   - export: Export configuration
//   - zone: Access zone name (e.g., "System")
//   - xprtsec: Transport security policy (e.g., "mtls", "tls:mtls"). Empty
//     omits the field so the cluster default applies.
//
// Returns:
//   - Export ID on success
//   - Error if export is nil, no path is set, zone is empty, or export creation fails
//
// Example:
//
//	paths := []string{"/ifs/data/vol1"}
//	export := &apiv2.Export{
//	    Paths:       &paths,
//	    Description: "k8s volume",
//	}
//	// mTLS-only export
//	exportID, err := ExportCreateWithZoneAndXprtsec(ctx, client, export, "System", "mtls")
func ExportCreateWithZoneAndXprtsec(
	ctx context.Context,
	client api.Client,
	export *apiv2.Export,
	zone string,
	xprtsec string,
) (int, error) {
	if export == nil {
		return 0, errors.New("export cannot be nil")
	}
	if export.Paths == nil || len(*export.Paths) == 0 {
		return 0, errors.New("no path set")
	}
	if zone == "" {
		return 0, errors.New("zone cannot be empty")
	}

	var resp apiv2.Export

	// Use v27 API endpoint for export creation with xprtsec support
	if err := client.Post(
		ctx,
		nfsExportsPathV27,
		"",
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		newExportReq(export, xprtsec),
		&resp); err != nil {
		return 0, err
	}

	return resp.ID, nil
}

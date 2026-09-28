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

package openapi

// V27NfsSettingsExport represents site-wide default values for NFS exports.
// Added in OneFS 9.16.0 (PAPI v27) for NFS over TLS (RFC 9289) support.
type V27NfsSettingsExport struct {
	// Xprtsec specifies the default transport security (in-flight encryption) policy for new exports.
	// Colon-delimited allowlist of none, tls and mtls (e.g., "none:tls:mtls" or "tls:mtls").
	// Default: "none:tls:mtls" (allows all modes).
	// The effective policy is the intersection with cluster-wide nfs_tls_mode.
	Xprtsec *string `json:"xprtsec,omitempty"`
}

// V27NfsSettingsExportResponse wraps the export settings response.
type V27NfsSettingsExportResponse struct {
	Settings *V27NfsSettingsExport `json:"settings,omitempty"`
}

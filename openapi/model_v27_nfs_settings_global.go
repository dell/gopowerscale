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

// V27NfsSettingsGlobal represents cluster-wide NFS settings including TLS configuration.
// Added in OneFS 9.16.0 (PAPI v27) for NFS over TLS (RFC 9289) support.
type V27NfsSettingsGlobal struct {
	// NfsTLSMode specifies cluster-wide NFS-over-TLS transport security policy.
	// Colon-separated tokens from: none, tls, mtls.
	// Example: "none:tls:mtls" (default) allows plaintext, TLS, and mTLS.
	// Tokens may appear in any order; duplicate tokens are accepted.
	NfsTLSMode *string `json:"nfs_tls_mode,omitempty"`

	// NfsTLSMinVersion specifies minimum TLS protocol version accepted for NFS-over-TLS connections.
	// Valid values: "1.2", "1.3" (default).
	// TLS 1.1 and earlier are never permitted.
	NfsTLSMinVersion *string `json:"nfs_tls_min_version,omitempty"`

	// NfsTLSCiphers specifies colon-separated OpenSSL cipher list used for NFS-over-TLS connections.
	// Applies to both TLS 1.2 and TLS 1.3.
	// Empty string (default) selects the built-in default cipher list.
	// Maximum length: 4096 characters.
	NfsTLSCiphers *string `json:"nfs_tls_ciphers,omitempty"`

	// NfsTLSOcspMode specifies OCSP revocation-checking mode for client certificates in mutual-TLS NFS connections.
	// Valid values:
	//   - "none" (default): disables OCSP
	//   - "strict": requires a reachable responder confirming the cert is valid
	//   - "no_responder_ok": allows the handshake when the responder is unreachable
	//   - "no_response_ok": allows the handshake when the responder returns no definitive answer
	NfsTLSOcspMode *string `json:"nfs_tls_ocsp_mode,omitempty"`

	// NfsTLSOcspAddress specifies OCSP responder URL (http:// only; OCSP does not use https)
	// used to check client certificate revocation for mutual-TLS NFS connections.
	// Empty string (default) disables OCSP.
	NfsTLSOcspAddress *string `json:"nfs_tls_ocsp_address,omitempty"`
}

// V27NfsSettingsGlobalResponse wraps the global NFS settings response.
type V27NfsSettingsGlobalResponse struct {
	Settings *V27NfsSettingsGlobal `json:"settings,omitempty"`
}

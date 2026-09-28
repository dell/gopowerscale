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

	apiv27 "github.com/dell/gopowerscale/api/v27"
	"github.com/dell/gopowerscale/openapi"
)

// NfsSettingsGlobal is the cluster-wide NFS TLS settings.
type NfsSettingsGlobal *openapi.V27NfsSettingsGlobal

// NfsSettingsExport is the site-wide default export settings.
type NfsSettingsExport *openapi.V27NfsSettingsExport

// GetNfsSettingsGlobal retrieves cluster-wide NFS TLS settings.
// Requires OneFS 9.16.0+ (PAPI v27) and ISI_PRIV_NFS_SETTINGS_GLOBAL privilege.
//
// Returns error if:
//   - OneFS version < 9.16.0 (PAPI v27 not available)
//   - Insufficient privileges
//   - Network or API errors
//
// Example:
//
//	settings, err := client.GetNfsSettingsGlobal(ctx)
//	if err != nil {
//	    log.Error("Failed to get NFS TLS settings", err)
//	}
//	log.Info("Cluster TLS mode:", *settings.NfsTLSMode)
func (c *Client) GetNfsSettingsGlobal(ctx context.Context) (NfsSettingsGlobal, error) {
	return apiv27.GetNfsSettingsGlobal(ctx, c.API)
}

// UpdateNfsSettingsGlobal modifies cluster-wide NFS TLS settings.
// Requires OneFS 9.16.0+ (PAPI v27), ISI_PRIV_NFS_SETTINGS_GLOBAL and ISI_PRIV_NFS_TLS privileges.
//
// Validation rules:
//   - nfs_tls_mode: colon-separated list containing only tokens: none, tls, mtls
//   - nfs_tls_min_version: must be "1.2" or "1.3"
//   - nfs_tls_ciphers: max 4096 characters, valid OpenSSL cipher list
//   - nfs_tls_ocsp_mode: one of: none, strict, no_responder_ok, no_response_ok
//   - nfs_tls_ocsp_address: valid HTTP URL or empty string
//
// Example:
//
//	tlsMode := "tls:mtls"  // Allow TLS and mTLS, but not plaintext
//	minVersion := "1.3"
//	settings := &openapi.V27NfsSettingsGlobal{
//	    NfsTLSMode:       &tlsMode,
//	    NfsTLSMinVersion: &minVersion,
//	}
//	err := client.UpdateNfsSettingsGlobal(ctx, settings)
func (c *Client) UpdateNfsSettingsGlobal(ctx context.Context, settings NfsSettingsGlobal) error {
	return apiv27.UpdateNfsSettingsGlobal(ctx, c.API, settings)
}

// GetNfsSettingsExport retrieves site-wide default export settings including default xprtsec.
// Requires OneFS 9.16.0+ (PAPI v27) and ISI_PRIV_NFS_SETTINGS_EXPORT privilege.
//
// Returns error if:
//   - OneFS version < 9.16.0 (PAPI v27 not available)
//   - Insufficient privileges
//   - Network or API errors
//
// Example:
//
//	settings, err := client.GetNfsSettingsExport(ctx)
//	if err != nil {
//	    log.Error("Failed to get export default settings", err)
//	}
//	log.Info("Default xprtsec:", *settings.Xprtsec)
func (c *Client) GetNfsSettingsExport(ctx context.Context) (NfsSettingsExport, error) {
	return apiv27.GetNfsSettingsExport(ctx, c.API)
}

// UpdateNfsSettingsExport modifies site-wide default export settings including default xprtsec.
// Requires OneFS 9.16.0+ (PAPI v27), ISI_PRIV_NFS_SETTINGS_EXPORT and ISI_PRIV_NFS_TLS privileges.
//
// Validation rules:
//   - xprtsec: colon-separated list containing only tokens: none, tls, mtls
//   - intersection with cluster nfs_tls_mode must not be empty
//
// Example:
//
//	xprtsec := "tls:mtls"  // New exports default to TLS/mTLS only
//	settings := &openapi.V27NfsSettingsExport{
//	    Xprtsec: &xprtsec,
//	}
//	err := client.UpdateNfsSettingsExport(ctx, settings)
func (c *Client) UpdateNfsSettingsExport(ctx context.Context, settings NfsSettingsExport) error {
	return apiv27.UpdateNfsSettingsExport(ctx, c.API, settings)
}

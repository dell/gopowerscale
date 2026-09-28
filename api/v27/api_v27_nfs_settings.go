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
	"fmt"

	"github.com/dell/gopowerscale/api"
	"github.com/dell/gopowerscale/openapi"
)

const (
	nfsSettingsGlobalPath = "platform/27/protocols/nfs/settings/global"
	nfsSettingsExportPath = "platform/27/protocols/nfs/settings/export"
)

// GetNfsSettingsGlobal retrieves cluster-wide NFS TLS settings.
// Requires ISI_PRIV_NFS_SETTINGS_GLOBAL privilege.
// Returns error if OneFS version < 9.16.0 (PAPI v27 not available).
func GetNfsSettingsGlobal(ctx context.Context, client api.Client) (*openapi.V27NfsSettingsGlobal, error) {
	var resp openapi.V27NfsSettingsGlobalResponse

	if err := client.Get(ctx, nfsSettingsGlobalPath, "", nil, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get NFS global settings: %w", err)
	}

	if resp.Settings == nil {
		return nil, fmt.Errorf("empty settings response from %s", nfsSettingsGlobalPath)
	}

	return resp.Settings, nil
}

// UpdateNfsSettingsGlobal modifies cluster-wide NFS TLS settings.
// Requires ISI_PRIV_NFS_SETTINGS_GLOBAL and ISI_PRIV_NFS_TLS privileges.
// Returns error if OneFS version < 9.16.0 (PAPI v27 not available).
//
// Validation rules:
//   - nfs_tls_mode: colon-separated list containing only tokens: none, tls, mtls
//   - nfs_tls_min_version: must be "1.2" or "1.3"
//   - nfs_tls_ciphers: max 4096 characters, valid OpenSSL cipher list
//   - nfs_tls_ocsp_mode: one of: none, strict, no_responder_ok, no_response_ok
//   - nfs_tls_ocsp_address: valid HTTP URL or empty string
func UpdateNfsSettingsGlobal(ctx context.Context, client api.Client, settings *openapi.V27NfsSettingsGlobal) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	if err := client.Put(ctx, nfsSettingsGlobalPath, "", nil, nil, settings, nil); err != nil {
		return fmt.Errorf("failed to update NFS global settings: %w", err)
	}

	return nil
}

// GetNfsSettingsExport retrieves site-wide default export settings including default xprtsec.
// Requires ISI_PRIV_NFS_SETTINGS_EXPORT privilege.
// Returns error if OneFS version < 9.16.0 (PAPI v27 not available).
func GetNfsSettingsExport(ctx context.Context, client api.Client) (*openapi.V27NfsSettingsExport, error) {
	var resp openapi.V27NfsSettingsExportResponse

	if err := client.Get(ctx, nfsSettingsExportPath, "", nil, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get NFS export settings: %w", err)
	}

	if resp.Settings == nil {
		return nil, fmt.Errorf("empty settings response from %s", nfsSettingsExportPath)
	}

	return resp.Settings, nil
}

// UpdateNfsSettingsExport modifies site-wide default export settings including default xprtsec.
// Requires ISI_PRIV_NFS_SETTINGS_EXPORT and ISI_PRIV_NFS_TLS privileges.
// Returns error if OneFS version < 9.16.0 (PAPI v27 not available).
//
// Validation rules:
//   - xprtsec: colon-separated list containing only tokens: none, tls, mtls
//   - intersection with cluster nfs_tls_mode must not be empty
func UpdateNfsSettingsExport(ctx context.Context, client api.Client, settings *openapi.V27NfsSettingsExport) error {
	if settings == nil {
		return fmt.Errorf("settings cannot be nil")
	}

	if err := client.Put(ctx, nfsSettingsExportPath, "", nil, nil, settings, nil); err != nil {
		return fmt.Errorf("failed to update NFS export settings: %w", err)
	}

	return nil
}

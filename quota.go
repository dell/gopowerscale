/*
Copyright (c) 2022 Dell Inc, or its subsidiaries.

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
package goisilon

import (
	"context"

	api "github.com/dell/goisilon/api/v1"
	apiV5 "github.com/dell/goisilon/api/v5"
)

// Quota maps to an Isilon filesystem quota.
type Quota *api.IsiQuota

// QuotaList maps to a set of quotas
type QuotaList []*api.IsiQuota

// QuotaResp returns a quota response with resume field
type QuotaResp *api.IsiQuotaListRespResume

// GetQuota returns a specific quota by volume name
func (c *Client) GetQuota(ctx context.Context, name string) (Quota, error) {
	quota, err := api.GetIsiQuota(ctx, c.API, c.API.VolumePath(name))
	if err != nil {
		return nil, err
	}

	return quota, nil
}

// GetAllQuotas returns all quotas on the cluster
func (c *Client) GetAllQuotas(ctx context.Context) (QuotaList, error) {
	quotas, err := api.GetAllIsiQuota(ctx, c.API)
	if err != nil {
		return nil, err
	}

	return quotas, nil
}

// GetQuotasWithResume returns a list of quota with resume field
func (c *Client) GetQuotasWithResume(ctx context.Context, resume string) (QuotaResp, error) {
	quotas, err := api.GetIsiQuotaWithResume(ctx, c.API, resume)
	if err != nil {
		return nil, err
	}

	return quotas, nil
}

// GetQuotaByID returns a specific quota by ID
func (c *Client) GetQuotaByID(ctx context.Context, ID string) (Quota, error) {
	quota, err := api.GetIsiQuotaByID(ctx, c.API, ID)
	if err != nil {
		return nil, err
	}

	return quota, nil
}

// GetQuotaWithPath returns a specific quota by path
func (c *Client) GetQuotaWithPath(ctx context.Context, path string) (Quota, error) {
	quota, err := api.GetIsiQuota(ctx, c.API, path)
	if err != nil {
		return nil, err
	}

	return quota, nil
}

// TODO: Add a means to set/update more fields of the quota

// CreateQuota creates a new hard directory quota with the specified size
// and container option
func (c *Client) CreateQuota(
	ctx context.Context, name string, container bool, size, softLimit, advisoryLimit, softGracePrd int64,
) (string, error) {
	return api.CreateIsiQuota(
		ctx, c.API, c.API.VolumePath(name), container, size, softLimit, advisoryLimit, softGracePrd)
}

// CreateQuotaWithPath creates a new hard directory quota with the specified size
// and container option
func (c *Client) CreateQuotaWithPath(
	ctx context.Context, path string, container bool, size, softLimit, advisoryLimit, softGracePrd int64,
) (string, error) {
	return api.CreateIsiQuota(
		ctx, c.API, path, container, size, softLimit, advisoryLimit, softGracePrd)
}

// SetQuotaSize sets the max size (hard threshold) of a quota for a volume
func (c *Client) SetQuotaSize(
	ctx context.Context, name string, size, softLimit, advisoryLimit, softGracePrd int64,
) (string, error) {
	return api.SetIsiQuotaHardThreshold(
		ctx, c.API, c.API.VolumePath(name), size, softLimit, advisoryLimit, softGracePrd)
}

// UpdateQuotaSize modifies the max size (hard threshold) of a quota for a volume
func (c *Client) UpdateQuotaSize(
	ctx context.Context, name string, size, softLimit, advisoryLimit, softGracePrd int64,
) error {
	return api.UpdateIsiQuotaHardThreshold(
		ctx, c.API, c.API.VolumePath(name), size, softLimit, advisoryLimit, softGracePrd)
}

// UpdateQuotaSizeByID modifies the max size (hard threshold) of a quota for a volume
func (c *Client) UpdateQuotaSizeByID(
	ctx context.Context, ID string, size, softLimit, advisoryLimit, softGracePrd int64,
) error {
	return api.UpdateIsiQuotaHardThresholdByID(
		ctx, c.API, ID, size, softLimit, advisoryLimit, softGracePrd)
}

// ClearQuota removes the quota from a volume
func (c *Client) ClearQuota(ctx context.Context, name string) error {
	return api.DeleteIsiQuota(ctx, c.API, c.API.VolumePath(name))
}

// ClearQuotaWithPath removes the quota from a volume with IsiPath as a parameter
func (c *Client) ClearQuotaWithPath(ctx context.Context, path string) error {
	return api.DeleteIsiQuota(ctx, c.API, path)
}

// ClearQuotaByID removes the quota from a volume by quota id
func (c *Client) ClearQuotaByID(ctx context.Context, id string) error {
	return api.DeleteIsiQuotaByID(ctx, c.API, id)
}

// ClearQuotaByIDWithZone removes the quota from a volume by quota id with access zone
func (c *Client) ClearQuotaByIDWithZone(ctx context.Context, id, zone string) error {
	return api.DeleteIsiQuotaByIDWithZone(ctx, c.API, id, zone)
}

// IsQuotaLicenseActivated checks if SmartQuotas has been activated (either licensed or in evaluation)
func (c *Client) IsQuotaLicenseActivated(ctx context.Context) (bool, error) {
	return apiV5.IsQuotaLicenseActivated(ctx, c.API)
}

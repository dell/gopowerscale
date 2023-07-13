/*
	Copyright (c) 2023 Dell Inc, or its subsidiaries.

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
	apiv12 "github.com/dell/goisilon/api/v12"
	"github.com/dell/goisilon/openapi"
)

// ListALlSmbSharesWithStructParams returns all the smb shares with params
func (c *Client) ListALlSmbSharesWithStructParams(ctx context.Context, params apiv12.ListV12SmbSharesParams) ([]openapi.V12SmbShareExtended, error) {
	var result []openapi.V12SmbShareExtended
	shares, err := apiv12.ListSmbShares(ctx, params, c.API)
	if err != nil {
		return nil, err
	}
	result = shares.Shares
	for shares.Resume != nil {
		resumeParam := apiv12.ListV12SmbSharesParams{Resume: shares.Resume}
		shares, err = apiv12.ListSmbShares(ctx, resumeParam, c.API)
		if err != nil {
			return nil, err
		}
		result = append(result, shares.Shares...)
	}
	return result, nil
}

// ListSmbSharesWithStructParams returns the smb shares with params
func (c *Client) ListSmbSharesWithStructParams(ctx context.Context, params apiv12.ListV12SmbSharesParams) (*openapi.V12SmbShares, error) {
	return apiv12.ListSmbShares(ctx, params, c.API)
}

// GetSmbShareWithStructParams return the specific smb share with params
func (c *Client) GetSmbShareWithStructParams(ctx context.Context, params apiv12.GetV12SmbShareParams) (*openapi.V12SmbSharesExtended, error) {
	return apiv12.GetSmbShare(ctx, params, c.API)
}

// CreateSmbShareWithStructParams creates a smb share with params
func (c *Client) CreateSmbShareWithStructParams(ctx context.Context, params apiv12.CreateV12SmbShareRequest) (*openapi.Createv12SmbShareResponse, error) {
	return apiv12.CreateSmbShare(ctx, params, c.API)
}

// DeleteSmbShareWithStructParams delete a specific smb share with params6570
func (c *Client) DeleteSmbShareWithStructParams(ctx context.Context, params apiv12.DeleteV12SmbShareRequest) error {
	return apiv12.DeleteSmbShare(ctx, params, c.API)
}

// UpdateSmbShareWithStructParams updates a smb share with params
func (c *Client) UpdateSmbShareWithStructParams(ctx context.Context, params apiv12.UpdateV12SmbShareRequest) error {
	return apiv12.UpdateSmbShare(ctx, params, c.API)
}

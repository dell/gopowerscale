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
package v12

import (
	"context"
	"github.com/dell/goisilon/api"
	"github.com/dell/goisilon/openapi"
)

const sharesPath = "/platform/12/protocols/smb/shares"

// ListV12SmbSharesParams contains smb shares params
type ListV12SmbSharesParams struct {
	Sort         *string `json:"sort,omitempty"`
	Zone         *string `json:"zone,omitempty"`
	Resume       *string `json:"resume,omitempty"`
	ResolveNames *bool   `json:"resolve_names,omitempty"`
	Limit        *int32  `json:"limit,omitempty"`
	Offset       *int32  `json:"offset,omitempty"`
	Scope        *string `json:"scope,omitempty"`
	Dir          *string `json:"dir,omitempty"`
}

// ListSmbShares GETs all smb shares.
func ListSmbShares(
	ctx context.Context,
	params ListV12SmbSharesParams,
	client api.Client) (*openapi.V12SmbShares, error) {

	var resp openapi.V12SmbShares
	if err := client.Get(
		ctx,
		sharesPath,
		"",
		api.StructToOrderedValues(params),
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// GetV12SmbShareParams contains smb share params
type GetV12SmbShareParams struct {
	V12SmbShareId string
	Scope         *string `json:"scope,omitempty"`
	ResolveNames  *bool   `json:"resolve_names,omitempty"`
	Zone          *string `json:"zone,omitempty"`
}

// GetSmbShare GET smb share.
func GetSmbShare(
	ctx context.Context,
	params GetV12SmbShareParams,
	client api.Client) (*openapi.V12SmbSharesExtended, error) {

	var resp openapi.V12SmbSharesExtended
	if err := client.Get(
		ctx,
		sharesPath,
		params.V12SmbShareId,
		api.StructToOrderedValues(params),
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// CreateV12SmbShareRequest contains request body
type CreateV12SmbShareRequest struct {
	V12SmbShare *openapi.V12SmbShare
	Zone        *string `json:"zone,omitempty"`
}

// CreateSmbShare POST smb share.
func CreateSmbShare(
	ctx context.Context,
	r CreateV12SmbShareRequest,
	client api.Client) (*openapi.Createv12SmbShareResponse, error) {

	var resp openapi.Createv12SmbShareResponse
	if err := client.Post(
		ctx,
		sharesPath,
		"",
		nil,
		nil, r.V12SmbShare,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// UpdateV12SmbShareRequest contains request body
type UpdateV12SmbShareRequest struct {
	V12SmbShareId string
	V12SmbShare   *openapi.V12SmbShareExtendedExtended
	Zone          *string `json:"zone,omitempty"`
}

// UpdateSmbShare UPDATE smb share.
func UpdateSmbShare(
	ctx context.Context,
	r UpdateV12SmbShareRequest,
	client api.Client) error {
	err := client.Put(
		ctx,
		sharesPath,
		r.V12SmbShareId,
		api.StructToOrderedValues(r),
		nil, r.V12SmbShare,
		nil)

	return err
}

// DeleteV12SmbShareRequest contains request params
type DeleteV12SmbShareRequest struct {
	V12SmbShareId string
	Zone          *string `json:"zone,omitempty"`
}

// DeleteSmbShare Delete one export.
func DeleteSmbShare(
	ctx context.Context, r DeleteV12SmbShareRequest,
	client api.Client) error {
	err := client.Delete(
		ctx,
		sharesPath,
		r.V12SmbShareId,
		api.StructToOrderedValues(r),
		nil, nil)

	return err
}

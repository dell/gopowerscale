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
package v4

import (
	"context"
	"github.com/dell/goisilon/api"
	"github.com/dell/goisilon/openapi"
)

const exportsPath = "/platform/4/protocols/nfs/exports"

type ListV4NfsExportsParams struct {
	Sort   *string `json:"sort,omitempty"`
	Zone   *string `json:"zone,omitempty"`
	Resume *string `json:"resume,omitempty"`
	Scope  *string `json:"scope,omitempty"`
	Limit  *int32  `json:"limit,omitempty"`
	Offset *int32  `json:"offset,omitempty"`
	Path   *string `json:"path,omitempty"`
	Check  *bool   `json:"check,omitempty"`
	Dir    *string `json:"dir,omitempty"`
}

// ListNfsExports GETs all exports.
func ListNfsExports(
	ctx context.Context,
	params ListV4NfsExportsParams,
	client api.Client) (*openapi.V2NfsExports, error) {

	var resp openapi.V2NfsExports
	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.StructToOrderedValues(params),
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

type GetV2NfsExportRequest struct {
	V2NfsExportId string
	Scope         *string `json:"scope,omitempty"`
	Zone          *string `json:"zone,omitempty"`
}

// GetNfsExport GET export.
func GetNfsExport(
	ctx context.Context,
	params GetV2NfsExportRequest,
	client api.Client) (*openapi.V2NfsExportsExtended, error) {

	var resp openapi.V2NfsExportsExtended
	if err := client.Get(
		ctx,
		exportsPath,
		params.V2NfsExportId,
		api.StructToOrderedValues(params),
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

type CreateV4NfsExportRequest struct {
	V4NfsExport             *openapi.V2NfsExport
	Force                   *bool   `json:"force,omitempty"`
	IgnoreUnresolvableHosts *bool   `json:"ignore_unresolvable_hosts,omitempty"`
	Zone                    *string `json:"zone,omitempty"`
	IgnoreConflicts         *bool   `json:"ignore_conflicts,omitempty"`
	IgnoreBadPaths          *bool   `json:"ignore_bad_paths,omitempty"`
	IgnoreBadAuth           *bool   `json:"ignore_bad_auth,omitempty"`
}

// CreateNfsExport Create one export.
func CreateNfsExport(
	ctx context.Context, r CreateV4NfsExportRequest,
	client api.Client) (*openapi.Createv3EventEventResponse, error) {
	var resp openapi.Createv3EventEventResponse
	if err := client.Post(
		ctx,
		exportsPath,
		"",
		api.StructToOrderedValues(r),
		nil, r.V4NfsExport,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

type UpdateV4NfsExportRequest struct {
	V2NfsExportId           string
	V2NfsExport             *openapi.V2NfsExportExtendedExtended
	Force                   *bool   `json:"force,omitempty"`
	IgnoreUnresolvableHosts *bool   `json:"ignore_unresolvable_hosts,omitempty"`
	Zone                    *string `json:"zone,omitempty"`
	IgnoreConflicts         *bool   `json:"ignore_conflicts,omitempty"`
	IgnoreBadPaths          *bool   `json:"ignore_bad_paths,omitempty"`
	IgnoreBadAuth           *bool   `json:"ignore_bad_auth,omitempty"`
}

// UpdateNfsExport Update one export.
func UpdateNfsExport(
	ctx context.Context, r UpdateV4NfsExportRequest,
	client api.Client) error {
	err := client.Put(
		ctx,
		exportsPath,
		r.V2NfsExportId,
		api.StructToOrderedValues(r),
		nil, r.V2NfsExport,
		nil)

	return err
}

type DeleteV4NfsExportRequest struct {
	V2NfsExportId string
	Zone          *string `json:"zone,omitempty"`
}

// DeleteNfsExport Delete one export.
func DeleteNfsExport(
	ctx context.Context, r DeleteV4NfsExportRequest,
	client api.Client) error {
	err := client.Delete(
		ctx,
		exportsPath,
		r.V2NfsExportId,
		api.StructToOrderedValues(r),
		nil, nil)

	return err
}

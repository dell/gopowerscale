/* 
 Copyright (c) 2019 Dell Inc, or its subsidiaries.

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
package v2

import (
	"context"
	"errors"
	"strconv"

	"github.com/dell/goisilon/api"
	"github.com/dell/goisilon/api/json"
)

// Export is an Isilon Export.
type Export struct {
	ID               int          `json:"id,omitmarshal"`
	Paths            *[]string    `json:"paths,omitempty"`
	Clients          *[]string    `json:"clients,omitempty"`
	RootClients      *[]string    `json:"root_clients,omitempty"`
	ReadWriteClients *[]string    `json:"read_write_clients,omitempty"`
	ReadOnlyClients  *[]string    `json:"read_only_clients,omitempty"`
	MapAll           *UserMapping `json:"map_all,omitempty"`
	MapRoot          *UserMapping `json:"map_root,omitempty"`
	MapNonRoot       *UserMapping `json:"map_non_root,omitempty"`
	MapFailure       *UserMapping `json:"map_failure,omitempty"`
	Description      string       `json:"description,omitempty"`
	Zone             string       `json:"zone,omitempty"`
}

// ExportList is a list of Isilon Exports.
type ExportList []*Export

// Exports is the whole struct of return JSON from the exports REST API
type Exports struct {
	Digest  string    `json:"digest,omitempty"`
	Exports []*Export `json:"exports,omitempty"`
	Resume  string    `json:"resume,omitempty"`
	Total   int       `json:"total,omitempty"`
}

// MarshalJSON marshals an ExportList to JSON.
func (l ExportList) MarshalJSON() ([]byte, error) {
	exports := struct {
		Exports []*Export `json:"exports,omitempty"`
	}{l}
	return json.Marshal(exports)
}

// UnmarshalJSON unmarshals an ExportList from JSON.
func (l *ExportList) UnmarshalJSON(text []byte) error {
	exports := struct {
		Exports []*Export `json:"exports,omitempty"`
	}{}
	if err := json.Unmarshal(text, &exports); err != nil {
		return err
	}
	*l = exports.Exports
	return nil
}

// ExportsList GETs all exports.
func ExportsList(
	ctx context.Context,
	client api.Client) ([]*Export, error) {

	var resp ExportList

	if err := client.Get(
		ctx,
		exportsPath,
		"",
		nil,
		nil,
		&resp); err != nil {

		return nil, err
	}

	return resp, nil
}

// ExportsListWithZone GETs all exports in the specified zone.
func ExportsListWithZone(
	ctx context.Context,
	client api.Client, zone string) ([]*Export, error) {

	var resp ExportList

	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		&resp); err != nil {

		return nil, err
	}

	return resp, nil
}

// ExportInspect GETs an export.
func ExportInspect(
	ctx context.Context,
	client api.Client,
	id int) (*Export, error) {

	var resp ExportList

	if err := client.Get(
		ctx,
		exportsPath,
		strconv.Itoa(id),
		nil,
		nil,
		&resp); err != nil {

		return nil, err
	}

	if len(resp) == 0 {
		return nil, nil
	}

	return resp[0], nil
}

// ExportCreate POSTs an Export object to the Isilon server.
func ExportCreate(
	ctx context.Context,
	client api.Client,
	export *Export) (int, error) {

	if export.Paths != nil && len(*export.Paths) == 0 {
		return 0, errors.New("no path set")
	}

	var resp Export

	if err := client.Post(
		ctx,
		exportsPath,
		"",
		nil,
		nil,
		export,
		&resp); err != nil {

		return 0, err
	}

	return resp.ID, nil
}

// ExportCreateWithZone POSTs an Export object with zone to the Isilon server.
func ExportCreateWithZone(
	ctx context.Context,
	client api.Client,
	export *Export, zone string) (int, error) {

	if export.Paths != nil && len(*export.Paths) == 0 {
		return 0, errors.New("no path set")
	}
	if zone == "" {
		return 0, errors.New("zone cannot be empty")
	}

	var resp Export

	if err := client.Post(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		export,
		&resp); err != nil {

		return 0, err
	}

	return resp.ID, nil
}

// ExportUpdate PUTs an Export object to the Isilon server.
func ExportUpdate(
	ctx context.Context,
	client api.Client,
	export *Export) error {

	return client.Put(
		ctx,
		exportsPath,
		strconv.Itoa(export.ID),
		nil,
		nil,
		export,
		nil)
}

// ExportUpdateWithZone PUTs an Export object in a specified zone to the Isilon server.
func ExportUpdateWithZone(
	ctx context.Context,
	client api.Client,
	export *Export,
	zone string) error {

	return client.Put(
		ctx,
		exportsPath,
		strconv.Itoa(export.ID),
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		export,
		nil)
}

// ExportDelete DELETEs an Export object on the Isilon server.
func ExportDelete(
	ctx context.Context,
	client api.Client,
	id int) error {

	return client.Delete(
		ctx,
		exportsPath,
		strconv.Itoa(id),
		nil,
		nil,
		nil)
}

// ExportDeleteWithZone DELETEs an Export object in the specified zone on the Isilon server.
func ExportDeleteWithZone(
	ctx context.Context,
	client api.Client,
	id int, zone string) error {

	return client.Delete(
		ctx,
		exportsPath,
		strconv.Itoa(id),
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		nil)
}

// SetExportClients sets an Export's clients property.
func SetExportClients(
	ctx context.Context,
	client api.Client,
	id int,
	addrs ...string) error {

	return ExportUpdate(ctx, client, &Export{ID: id, Clients: &addrs})
}

// SetExportRootClients sets an Export's root_clients property.
func SetExportRootClients(
	ctx context.Context,
	client api.Client,
	id int,
	addrs ...string) error {

	return ExportUpdate(ctx, client, &Export{ID: id, RootClients: &addrs})
}

// Unexport is an alias for ExportDelete.
func Unexport(
	ctx context.Context,
	client api.Client,
	id int) error {

	return ExportDelete(ctx, client, id)
}

// UnexportWithZone is an alias for ExportDeleteWithZone.
func UnexportWithZone(
	ctx context.Context,
	client api.Client,
	id int, zone string) error {

	return ExportDeleteWithZone(ctx, client, id, zone)
}

// ExportsListWithResume GETs the next page of exports based on the resume token from the previous call.
func ExportsListWithResume(
	ctx context.Context,
	client api.Client, resume string) (*Exports, error) {
	var resp Exports

	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("resume"), []byte(resume)},
		},
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// ExportsListWithLimit GETs a number of exports in the default sequence and the number is the parameter limit.
func ExportsListWithLimit(
	ctx context.Context,
	client api.Client, limit string) (*Exports, error) {
	var resp Exports

	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("limit"), []byte(limit)},
		},
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// ExportsListWithParams GETs exports based on the parmapeters.
func ExportsListWithParams(
	ctx context.Context,
	client api.Client, params api.OrderedValues) (*Exports, error) {
	var resp Exports

	if err := client.Get(
		ctx,
		exportsPath,
		"",
		params,
		nil,
		&resp); err != nil {

		return nil, err
	}

	return &resp, nil
}

// GetExportWithPath GETs an export with the specified target path
func GetExportWithPath(
	ctx context.Context,
	client api.Client,
	path string) (*Export, error) {
	var resp ExportList
	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("path"), []byte(path)},
		},
		nil,
		&resp); err != nil {

		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}

	return resp[0], nil
}

// GetExportWithPathAndZone GETs an export with the specified target path and access zone
func GetExportWithPathAndZone(
	ctx context.Context,
	client api.Client,
	path, zone string) (*Export, error) {
	var resp ExportList
	if zone == "" {
		zone = "System"
	}
	if err := client.Get(
		ctx,
		exportsPath,
		"",
		api.OrderedValues{
			{[]byte("path"), []byte(path)},
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		&resp); err != nil {

		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}

	return resp[0], nil
}

// GetExportByIDWithZone get the export by export id and access zone
func GetExportByIDWithZone(
	ctx context.Context,
	client api.Client,
	id int,
	zone string) (*Export, error) {
	var resp ExportList
	if err := client.Get(
		ctx,
		exportsPath,
		strconv.Itoa(id),
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		&resp); err != nil {
		return nil, err
	}
	if len(resp) == 0 {
		return nil, nil
	}

	return resp[0], nil
}

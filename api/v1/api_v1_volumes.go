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
package v1

import (
	"context"
	"github.com/dell/goisilon/api"
	"path"
)

// const defaultACL = "public_read_write"
const defaultACL = "0777"

var (
	aclQS           = api.OrderedValues{{[]byte("acl")}}
	metadataQS      = api.OrderedValues{{[]byte("metadata")}}
	recursiveTrueQS = api.OrderedValues{
		{[]byte("recursive"), []byte("true")},
	}
	sizeQS = api.OrderedValues{
		{[]byte("detail"), []byte("size")},
		{[]byte("max-depth"), []byte("-1")},
	}
	mergeQS = api.OrderedValues{
		{[]byte("merge"), []byte("True")},
	}
)

// GetIsiVolumes queries a list of all volumes on the cluster
func GetIsiVolumes(
	ctx context.Context,
	client api.Client) (resp *getIsiVolumesResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volumes/
	err = client.Get(ctx, realNamespacePath(client), "", nil, nil, &resp)
	return resp, err
}

// CreateIsiVolume makes a new volume on the cluster
func CreateIsiVolume(
	ctx context.Context,
	client api.Client,
	name string) (resp *getIsiVolumesResp, err error) {

	return CreateIsiVolumeWithACL(ctx, client, name, defaultACL)
}

// CreateIsiVolumeWithIsiPath makes a new volume with isiPath on the cluster
func CreateIsiVolumeWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, name, isiVolumePathPermissions string) (resp *getIsiVolumesResp, err error) {
	return CreateIsiVolumeWithACLAndIsiPath(ctx, client, isiPath, name, isiVolumePathPermissions)
}

// CreateIsiVolumeWithIsiPathMetaData makes a new volume with isiPath on the cluster
func CreateIsiVolumeWithIsiPathMetaData(
	ctx context.Context,
	client api.Client,
	isiPath, name, isiVolumePathPermissions string, metadata map[string]string) (resp *getIsiVolumesResp, err error) {
	return CreateIsiVolumeWithACLAndIsiPathMetaData(ctx, client, isiPath, name, isiVolumePathPermissions, metadata)
}

// CreateIsiVolumeWithACL makes a new volume on the cluster with the specified permissions
func CreateIsiVolumeWithACL(
	ctx context.Context,
	client api.Client,
	name, ACL string) (resp *getIsiVolumesResp, err error) {

	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/volume_name
	//             x-isi-ifs-target-type: container
	//             x-isi-ifs-access-control: ACL
	//
	//             PUT https://1.2.3.4:8080/namespace/path/to/volumes/volume_name?acl
	//             {authoritative: "acl",
	//              action: "update",
	//              owner: {name: "username", type: "user"},
	//              group: {name: "groupname", type: "group"}
	//             }

	createVolumeHeaders := map[string]string{
		"x-isi-ifs-target-type":    "container",
		"x-isi-ifs-access-control": ACL,
	}

	// create the volume
	err = client.Put(
		ctx,
		realNamespacePath(client),
		name,
		nil,
		createVolumeHeaders,
		nil,
		&resp)

	// The following code is completely pointless and also counterproductive -
	// The folder is already owned by client.User() because that is the user
	// that we authenticated to the API with. It's useless additional work that
	// also fails if the parent doesn't have an ACL granting us std_write_owner
	/*
		if err != nil {
			return resp, err
		}
		var data = &AclRequest{
			"acl",
			"update",
			&Ownership{client.User(), "user"},
			nil,
		}

		if group := client.Group(); group != "" {
			data.Group = &Ownership{group, "group"}
		}

		// set the ownership of the volume
		err = client.Put(
			ctx,
			realNamespacePath(client),
			name,
			aclQS,
			nil,
			data,
			&resp)
	*/

	return resp, err
}

// CreateIsiVolumeWithACLAndIsiPath makes a new volume on the cluster with the specified permissions and isiPath
func CreateIsiVolumeWithACLAndIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, name, ACL string) (resp *getIsiVolumesResp, err error) {

	createVolumeHeaders := map[string]string{
		"x-isi-ifs-target-type":    "container",
		"x-isi-ifs-access-control": ACL,
	}
	// create the volume
	err = client.Put(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		nil,
		createVolumeHeaders,
		nil,
		&resp)
	return resp, err
}

// CreateIsiVolumeWithACLAndIsiPathMetadata makes a new volume on the cluster with the specified permissions and isiPath
func CreateIsiVolumeWithACLAndIsiPathMetaData(
	ctx context.Context,
	client api.Client,
	isiPath, name, ACL string, metadata map[string]string) (resp *getIsiVolumesResp, err error) {
	var createVolumeHeaders = make(map[string]string)

	createVolumeHeaders["x-isi-ifs-target-type"] = "container"
	createVolumeHeaders["x-isi-ifs-access-control"] = ACL

	if len(metadata) != 0 {
		for key, value := range metadata {
			createVolumeHeaders[key] = value
		}
	}

	// create the volume
	err = client.Put(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		nil,
		createVolumeHeaders,
		nil,
		&resp)
	return resp, err
}

// GetIsiVolume queries the attributes of a volume on the cluster
func GetIsiVolume(
	ctx context.Context,
	client api.Client,
	name string) (resp *getIsiVolumeAttributesResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/?metadata
	err = client.Get(
		ctx,
		realNamespacePath(client),
		name,
		metadataQS,
		nil,
		&resp)
	return resp, err
}

// GetIsiVolumeWithIsiPath queries the attributes of a volume with isiPath on the cluster
func GetIsiVolumeWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, name string) (resp *getIsiVolumeAttributesResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/?metadata
	err = client.Get(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		metadataQS,
		nil,
		&resp)
	return resp, err
}

// GetIsiVolumeWithoutMetadata is used to check whether a volume exists thus the url does not append the metadata parameter.
func GetIsiVolumeWithoutMetadata(
	ctx context.Context,
	client api.Client,
	name string) (err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/
	err = client.Get(
		ctx,
		realNamespacePath(client),
		name,
		nil,
		nil,
		&getIsiVolumeResp{})

	return err
}

// GetIsiVolumeWithoutMetadataWithIsiPath is used to check whether a volume exists with isiPath thus the url does not append the metadata parameter.
func GetIsiVolumeWithoutMetadataWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, name string) (err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume/
	err = client.Get(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		nil,
		nil,
		&getIsiVolumeResp{})

	return err
}

// DeleteIsiVolume removes a volume from the cluster
func DeleteIsiVolume(
	ctx context.Context,
	client api.Client,
	name string) (resp *getIsiVolumesResp, err error) {

	err = client.Delete(
		ctx,
		realNamespacePath(client),
		name,
		recursiveTrueQS,
		nil,
		&resp)
	return resp, err
}

// DeleteIsiVolumeWithIsiPath removes a volume from the cluster with isiPath
func DeleteIsiVolumeWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, name string) (resp *getIsiVolumesResp, err error) {

	err = client.Delete(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		recursiveTrueQS,
		nil,
		&resp)
	return resp, err
}

// CopyIsiVolume creates a new volume on the cluster based on an existing volume
func CopyIsiVolume(
	ctx context.Context,
	client api.Client,
	sourceName, destinationName string) (resp *getIsiVolumesResp, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/destination_volume_name?merge=True
	//             x-isi-ifs-copy-source: /path/to/volumes/source_volume_name

	// copy the volume
	err = client.Put(
		ctx,
		realNamespacePath(client),
		destinationName,
		mergeQS,
		map[string]string{
			"x-isi-ifs-copy-source": path.Join(
				"/",
				realNamespacePath(client),
				sourceName),
		},
		nil,
		&resp)
	return resp, err
}

// CopyIsiVolumeWithIsiPath creates a new volume with isiPath on the cluster based on an existing volume
func CopyIsiVolumeWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, sourceName, destinationName string) (resp *getIsiVolumesResp, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/destination_volume_name?merge=True
	//             x-isi-ifs-copy-source: /path/to/volumes/source_volume_name
	//             x-isi-ifs-mode-mask: preserve

	// copy the volume
	err = client.Put(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		destinationName,
		mergeQS,
		map[string]string{
			"x-isi-ifs-copy-source": path.Join(
				"/",
				GetRealNamespacePathWithIsiPath(isiPath),
				sourceName),
			"x-isi-ifs-mode-mask": "preserve",
		},
		nil,
		&resp)
	return resp, err
}

// GetIsiVolumeWithSize lists size of all the children files and subfolders in a directory
func GetIsiVolumeWithSize(
	ctx context.Context,
	client api.Client,
	isiPath, name string) (resp *getIsiVolumeSizeResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/volume?detail=size&max-depth=-1
	err = client.Get(
		ctx,
		GetRealNamespacePathWithIsiPath(isiPath),
		name,
		sizeQS,
		nil,
		&resp)

	return resp, err
}

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
package v1

import (
	"context"
	"errors"
	"fmt"
	"github.com/dell/goisilon/api"
	"path"
)

// GetIsiSnapshots queries a list of all snapshots on the cluster
func GetIsiSnapshots(
	ctx context.Context,
	client api.Client) (resp *getIsiSnapshotsResp, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/snapshot/snapshots
	err = client.Get(ctx, snapshotsPath, "", nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetIsiSnapshot queries an individual snapshot on the cluster
func GetIsiSnapshot(
	ctx context.Context,
	client api.Client,
	id int64) (*IsiSnapshot, error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/snapshot/snapshots/123
	snapshotURL := fmt.Sprintf("%s/%d", snapshotsPath, id)
	var resp *getIsiSnapshotsResp
	err := client.Get(ctx, snapshotURL, "", nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	// PAPI returns the snapshot data in a JSON list with the same structure as
	// when querying all snapshots.  Since this is for a single Id, we just
	// want the first (and should be only) entry in the list.
	return resp.SnapshotList[0], nil
}

// GetIsiSnapshotByIdentity queries an individual snapshot on the cluster
// parm identity string: snapshot id or name
func GetIsiSnapshotByIdentity(
	ctx context.Context,
	client api.Client,
	identity string) (*IsiSnapshot, error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/snapshot/snapshots/id|name
	snapshotURL := fmt.Sprintf("%s/%s", snapshotsPath, identity)
	var resp *getIsiSnapshotsResp
	err := client.Get(ctx, snapshotURL, "", nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	// PAPI returns the snapshot data in a JSON list with the same structure as
	// when querying all snapshots.  Since this is for a single Id, we just
	// want the first (and should be only) entry in the list.
	return resp.SnapshotList[0], nil
}

// CreateIsiSnapshot makes a new snapshot on the cluster
func CreateIsiSnapshot(
	ctx context.Context,
	client api.Client,
	path, name string) (resp *IsiSnapshot, err error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/snapshot/snapshots
	//            Content-Type: application/json
	//            {path: "/path/to/volume"
	//             name: "snapshot_name"  <--- optional
	//            }
	if path == "" {
		return nil, errors.New("no path set")
	}

	data := &SnapshotPath{Path: path}
	if name != "" {
		data.Name = name
	}

	err = client.Post(ctx, snapshotsPath, "", nil, nil, data, &resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CopyIsiSnapshot copies all files/directories in a snapshot to a new directory
func CopyIsiSnapshot(
	ctx context.Context,
	client api.Client,
	sourceSnapshotName, sourceVolume, destinationName string, zonePath, accessZone string) (resp *IsiVolume, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/destination_volume_name?merge=True
	//             x-isi-ifs-copy-source: /path/to/snapshot/volumes/source_volume_name

	headers := map[string]string{
		"x-isi-ifs-copy-source": path.Join(
			"/",
			realVolumeSnapshotPath(client, sourceSnapshotName, zonePath, accessZone),
			sourceVolume),
	}

	// copy the volume
	err = client.Put(ctx, realNamespacePath(client), destinationName, mergeQS, headers, nil, &resp)
	return resp, err
}

// CopyIsiSnapshotWithIsiPath copies all files/directories in a snapshot in under the defined isiPath to a new directory
func CopyIsiSnapshotWithIsiPath(
	ctx context.Context,
	client api.Client,
	isiPath, snapshotSourceVolumeIsiPath, sourceSnapshotName, sourceVolume, destinationName string, accessZone string) (resp *IsiCopySnapshotResp, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080/namespace/path/to/volumes/destination_volume_name?merge=True
	//             x-isi-ifs-copy-source: /path/to/snapshot/volumes/source_volume_name
	//             x-isi-ifs-mode-mask: preserve
	zone, err := GetZoneByName(ctx, client, accessZone)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"x-isi-ifs-copy-source": path.Join(
			"/",
			GetRealVolumeSnapshotPathWithIsiPath(snapshotSourceVolumeIsiPath, zone.Path, sourceSnapshotName, accessZone),
			sourceVolume),
		"x-isi-ifs-mode-mask": "preserve",
	}
	// copy the volume
	err = client.Put(ctx, GetRealNamespacePathWithIsiPath(isiPath), destinationName, mergeQS, headers, nil, &resp)
	return resp, err
}

// RemoveIsiSnapshot deletes a snapshot from the cluster
func RemoveIsiSnapshot(
	ctx context.Context,
	client api.Client,
	id int64) error {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/snapshot/snapshots/123
	snapshotURL := fmt.Sprintf("%s/%d", snapshotsPath, id)
	err := client.Delete(ctx, snapshotURL, "", nil, nil, nil)

	return err
}

// GetIsiSnapshotFolderWithSize lists size of all the children files and subfolders in a sanpshot directory
func GetIsiSnapshotFolderWithSize(
	ctx context.Context,
	client api.Client,
	isiPath, name, volume string, accessZone string) (resp *getIsiVolumeSizeResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/namespace/path/to/snapshot?detail=size&max-depth=-1
	zone, err := GetZoneByName(ctx, client, accessZone)
	if err != nil {
		return nil, err
	}

	err = client.Get(
		ctx,
		GetRealVolumeSnapshotPathWithIsiPath(isiPath, zone.Path, name, accessZone),
		volume,
		sizeQS,
		nil,
		&resp)

	return resp, err
}

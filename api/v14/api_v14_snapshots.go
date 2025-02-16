/*
Copyright (c) 2025 Dell Inc, or its subsidiaries.

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

package v14

import (
	"context"
	"fmt"

	"github.com/dell/goisilon/api"
)

// GetIsiWritableSnapshots retrieves a list of writable snapshots.
//
// ctx: the context.
// client: the API client.
// Returns a list of writable snapshots and an error in case of failure.
func GetIsiWritableSnapshots(
	ctx context.Context,
	client api.Client,
) ([]*IsiWritableSnapshot, error) {
	var resp *IsiWritableSnapshotQueryResponse
	err := client.Get(ctx, writableSnapshotPath, "", nil, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get writable snapshots from array: %v", err)
	}

	result := make([]*IsiWritableSnapshot, 0, resp.Total)
	result = append(result, resp.Writable...)

	for {
		if resp.Resume == "" {
			break
		}

		query := api.OrderedValues{
			{[]byte("resume"), []byte(resp.Resume)},
		}

		resp = nil
		err = client.Get(ctx, writableSnapshotPath, "", query, nil, &resp)
		if err != nil {
			return nil, fmt.Errorf("failed to get writable snapshots (query mode) from array: %v", err)
		}

		result = append(result, resp.Writable...)
	}

	return result, err
}

// GetIsiWritableSnapshot retrieves a writable snapshot.
//
// ctx: the context.
// client: the API client.
// snapshotPath: the path of the snapshot.
//
// Returns the snapshot on success and error in case of failure.
func GetIsiWritableSnapshot(
	ctx context.Context,
	client api.Client,
	snapshotPath string,
) (*IsiWritableSnapshot, error) {
	var resp *IsiWritableSnapshotQueryResponse
	err := client.Get(ctx, writableSnapshotPath+snapshotPath, "", nil, nil, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to get writable snapshot: %v", err)
	}

	return resp.Writable[0], nil
}

// CreateWritableSnapshot creates a writable snapshot.
//
// ctx: the context.
// client: the API client.
// sourceSnapshot: the source snapshot name or ID.
// destination: the destination path, must not be nested under the source snapshot.
//
// Returns the response and error.
func CreateWritableSnapshot(
	ctx context.Context,
	client api.Client,
	sourceSnapshot string,
	destination string,
) (resp *IsiWritableSnapshot, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080//platform/14/snapshot/writable
	//             Body: {"src_snap": sourceSnapshot, "dst_path": destination}

	body := map[string]string{
		"src_snap": sourceSnapshot,
		"dst_path": destination,
	}

	err = client.Post(ctx, writableSnapshotPath, "", nil, nil, body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create writable snapshot: %v", err)
	}

	return resp, err
}

func RemoveWritableSnapshot(
	ctx context.Context,
	client api.Client,
	snapshotPath string,
) error {
	return client.Delete(ctx, writableSnapshotPath+snapshotPath, "", nil, nil, nil)
}

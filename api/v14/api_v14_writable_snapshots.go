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

package v14

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dell/gopowerscale/api"
)

const (
	writableSnapshotsPath = "platform/14/snapshot/writable"
)

// CreateIsiWritableSnapshot creates a new writable snapshot on the cluster.
//
// Example:
//
//	request := &IsiWritableSnapshotRequest{
//	    SourceSnapshot:      "12345",
//	    DestinationPath: "/ifs/data/writable-vol-1",
//	}
//	ws, err := CreateIsiWritableSnapshot(ctx, client, request)
//	// Returns writable snapshot with fields: ID, DstPath, SrcID, SrcPath, SrcSnap, Created, State, LogSize, PhysSize
func CreateIsiWritableSnapshot(
	ctx context.Context,
	client api.Client,
	request *IsiWritableSnapshotRequest,
) (*IsiWritableSnapshotResponse, error) {
	// Validate input
	if request == nil {
		return nil, errors.New("request cannot be nil")
	}
	if request.DestinationPath == "" {
		return nil, errors.New("destination_path cannot be empty")
	}

	// PAPI call: POST https://1.2.3.4:8080/platform/14/snapshot/writable
	var resp *IsiWritableSnapshotResponse
	err := client.Post(ctx, writableSnapshotsPath, "", nil, nil, request, &resp)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// GetIsiWritableSnapshot retrieves details of a writable snapshot by path.
func GetIsiWritableSnapshot(
	ctx context.Context,
	client api.Client,
	path string,
) (*IsiWritableSnapshotResponse, error) {
	// Validate input
	if path == "" {
		return nil, errors.New("path cannot be empty")
	}

	// PAPI call: GET https://1.2.3.4:8080/platform/14/snapshot/writable/{path}
	endpoint := fmt.Sprintf("%s/%s", writableSnapshotsPath, strings.Trim(path, "/"))
	var listResp *IsiWritableSnapshotListResponse
	err := client.Get(ctx, endpoint, "", nil, nil, &listResp)
	if err != nil {
		return nil, err
	}

	if listResp == nil {
		return nil, errors.New("unexpected nil response from API")
	}

	// Extract first element from Writable array
	if len(listResp.Writable) == 0 {
		return nil, errors.New("writable snapshot not found in response")
	}

	return &listResp.Writable[0], nil
}

// DeleteIsiWritableSnapshot deletes a writable snapshot from the cluster.
func DeleteIsiWritableSnapshot(
	ctx context.Context,
	client api.Client,
	path string,
) error {
	// Validate input
	if path == "" {
		return errors.New("path cannot be empty")
	}

	// PAPI call: DELETE https://1.2.3.4:8080/platform/14/snapshot/writable/{path}
	endpoint := fmt.Sprintf("%s/%s", writableSnapshotsPath, strings.Trim(path, "/"))
	err := client.Delete(ctx, endpoint, "", nil, nil, nil)

	return err
}

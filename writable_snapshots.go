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

package gopowerscale

import (
	"context"

	apiv14 "github.com/dell/gopowerscale/api/v14"
)

// WritableSnapshot represents a writable snapshot.
type WritableSnapshot *apiv14.IsiWritableSnapshotResponse

// CreateWritableSnapshot creates a new writable snapshot from an existing snapshot.
//
// Example:
//
//	ws, err := client.CreateWritableSnapshot(ctx, "/ifs/data/writable-vol-1", "3754")
//	// Returns writable snapshot with fields: ID, DstPath, SrcID, SrcPath, SrcSnap, Created, State, LogSize, PhysSize
func (c *Client) CreateWritableSnapshot(
	ctx context.Context,
	destinationPath string,
	sourceSnapshot string,
) (WritableSnapshot, error) {
	request := &apiv14.IsiWritableSnapshotRequest{
		DestinationPath: destinationPath,
		SourceSnapshot:  sourceSnapshot,
	}

	return apiv14.CreateIsiWritableSnapshot(ctx, c.API, request)
}

// GetWritableSnapshot retrieves details of a writable snapshot by path.
func (c *Client) GetWritableSnapshot(
	ctx context.Context,
	path string,
) (WritableSnapshot, error) {
	return apiv14.GetIsiWritableSnapshot(ctx, c.API, path)
}

// DeleteWritableSnapshot deletes a writable snapshot from the cluster.
func (c *Client) DeleteWritableSnapshot(
	ctx context.Context,
	path string,
) error {
	return apiv14.DeleteIsiWritableSnapshot(ctx, c.API, path)
}

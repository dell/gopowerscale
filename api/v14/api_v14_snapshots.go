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

// CreateWriteableSnapshot creates a writeable snapshot.
//
// ctx: the context.
// client: the API client.
// sourceSnapshot: the source snapshot name or ID.
// destination: the destination path, must not be nested under the source snapshot.
//
// Returns the response and error.
func CreateWriteableSnapshot(
	ctx context.Context,
	client api.Client,
	sourceSnapshot string,
	destination string,
) (resp *IsiWriteableSnapshot, err error) {
	// PAPI calls: PUT https://1.2.3.4:8080//platform/14/snapshot/writable
	//             Body: {"src_snap": sourceSnapshot, "dst_path": destination}

	body := map[string]string{
		"src_snap": sourceSnapshot,
		"dst_path": destination,
	}

	err = client.Post(ctx, writeableSnapshotPath, "", nil, nil, body, &resp)
	if err != nil {
		return nil, fmt.Errorf("failed to create writeable snapshot: %v", err)
	}

	return resp, err
}

func RemoveWriteableSnapshot(
	ctx context.Context,
	client api.Client,
	snapshotPath string,
) error {
	return client.Delete(ctx, writeableSnapshotPath+snapshotPath, "", nil, nil, nil)
}

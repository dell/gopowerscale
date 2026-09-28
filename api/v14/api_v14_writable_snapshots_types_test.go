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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsiWritableSnapshotRequest_JSONMarshaling(t *testing.T) {
	req := &IsiWritableSnapshotRequest{
		SourceSnapshot:  "12345",
		DestinationPath: "/ifs/data/writable-vol-1",
	}

	// Marshal to JSON
	data, err := json.Marshal(req)
	assert.NoError(t, err)
	assert.Contains(t, string(data), "src_snap")
	assert.Contains(t, string(data), "dst_path")

	// Unmarshal from JSON
	var decoded IsiWritableSnapshotRequest
	err = json.Unmarshal(data, &decoded)
	assert.NoError(t, err)
	assert.Equal(t, req.SourceSnapshot, decoded.SourceSnapshot)
	assert.Equal(t, req.DestinationPath, decoded.DestinationPath)
}

func TestIsiWritableSnapshotResponse_JSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"id": 24315648,
		"dst_path": "/ifs/writeable",
		"src_id": 3754,
		"src_path": "/ifs/sky",
		"src_snap": "Snapshot: 2026Jul22, 5:05 PM",
		"created": 1783420362,
		"state": "active",
		"log_size": 0,
		"phys_size": 2048
	}`

	var resp IsiWritableSnapshotResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	assert.NoError(t, err)
	assert.Equal(t, int64(24315648), resp.ID)
	assert.Equal(t, "/ifs/writeable", resp.DstPath)
	assert.Equal(t, int64(3754), resp.SrcID)
	assert.Equal(t, "/ifs/sky", resp.SrcPath)
	assert.Equal(t, "Snapshot: 2026Jul22, 5:05 PM", resp.SrcSnap)
	assert.Equal(t, int64(1783420362), resp.Created)
	assert.Equal(t, "active", resp.State)
	assert.Equal(t, int64(0), resp.LogSize)
	assert.Equal(t, int64(2048), resp.PhysSize)
}

func TestIsiWritableSnapshotListResponse_JSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"writable": [
			{
				"id": 24315648,
				"dst_path": "/ifs/writeable",
				"src_id": 3754,
				"src_path": "/ifs/sky",
				"src_snap": "Snapshot: 2026Jul22, 5:05 PM",
				"created": 1783420362,
				"state": "active",
				"log_size": 0,
				"phys_size": 2048
			}
		]
	}`

	var listResp IsiWritableSnapshotListResponse
	err := json.Unmarshal([]byte(jsonData), &listResp)
	assert.NoError(t, err)
	assert.Len(t, listResp.Writable, 1)
	assert.Equal(t, int64(24315648), listResp.Writable[0].ID)
	assert.Equal(t, "/ifs/writeable", listResp.Writable[0].DstPath)
	assert.Equal(t, "active", listResp.Writable[0].State)
}

func TestIsiWritableSnapshotListResponse_EmptyArray(t *testing.T) {
	jsonData := `{"writable": []}`

	var listResp IsiWritableSnapshotListResponse
	err := json.Unmarshal([]byte(jsonData), &listResp)
	assert.NoError(t, err)
	assert.Len(t, listResp.Writable, 0)
}

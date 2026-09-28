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

// IsiWritableSnapshotRequest represents the request body for creating a writable snapshot.
type IsiWritableSnapshotRequest struct {
	DestinationPath string `json:"dst_path"`
	SourceSnapshot  string `json:"src_snap"`
}

// IsiWritableSnapshotResponse represents an individual writable snapshot returned by the OneFS API.
type IsiWritableSnapshotResponse struct {
	ID       int64  `json:"id"`
	DstPath  string `json:"dst_path"`
	SrcID    int64  `json:"src_id"`
	SrcPath  string `json:"src_path"`
	SrcSnap  string `json:"src_snap"`
	Created  int64  `json:"created"`
	State    string `json:"state"`
	LogSize  int64  `json:"log_size"`
	PhysSize int64  `json:"phys_size"`
}

// IsiWritableSnapshotListResponse represents the wrapper response from the GET endpoint
// containing an array of writable snapshots.
type IsiWritableSnapshotListResponse struct {
	Writable []IsiWritableSnapshotResponse `json:"writable"`
}

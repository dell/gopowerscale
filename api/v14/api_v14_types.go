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

package v14

// IsiClusterAcs Cluster ACS status.
type IsiClusterAcs struct {
	// list of failed nodes serial number.
	FailedNodesSn []string `json:"failed_nodes_sn,omitempty"`
	// the number of joined nodes.
	JoinedNodes *int32 `json:"joined_nodes,omitempty"`
	// the status of license activation.
	LicenseStatus *string `json:"license_status,omitempty"`
	// the status of SRS enablement.
	SrsStatus *string `json:"srs_status,omitempty"`
	// total nodes number of the cluster.
	TotalNodes *int32 `json:"total_nodes,omitempty"`
	// list of unresponsive nodes serial number.
	UnresponsiveSn []string `json:"unresponsive_sn,omitempty"`
}

// An IsiWriteableSnapshot is a writable snapshot.
type IsiWriteableSnapshot struct {
	// The Unix Epoch time the writable snapshot was created.
	Created int64 `json:"created"`
	// The /ifs path of user supplied source snapshot. This will be null for writable snapshots pending delete.
	SrcPath string `json:"src_path"`
	// The user supplied /ifs path of writable snapshot.
	DstPath string `json:"dst_path"`
	// The system ID given to the writable snapshot.
	ID int64 `json:"id"`
	// The system ID of the user supplied source snapshot.
	SrcID int64 `json:"src_id"`
	// The user supplied source snapshot name or ID. This will be null for writable snapshots pending delete.
	SrcSnap string `json:"src_snap"`
	// The sum in bytes of logical size of files in this writable snapshot.
	LogSize int64 `json:"log_size"`
	// The amount of storage in bytes used to store this writable snapshot.
	PhysicalSize int64 `json:"phys_size"`
	// Writable Snapshot state.
	State string `json:"state"`
}

// IsiWriteableSnapshotQueryResponse is the response to a writable snapshot query.
type IsiWriteableSnapshotQueryResponse struct {
	// Total number of items available.
	Total int64 `json:"total,omitempty"`
	// Used to continue a query. This is null for the last page.
	Resume string `json:"resume,omitempty"`
	// List of writable snapshots.
	Writeable []*IsiWriteableSnapshot `json:"writable"`
}

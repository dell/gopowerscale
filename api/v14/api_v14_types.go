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

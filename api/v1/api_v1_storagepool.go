/*
Copyright (c) 2022-2026 Dell Inc, or its subsidiaries.

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

	"github.com/dell/gopowerscale/api"
)

const (
	nodePoolsPath = "platform/1/storagepool/nodepools"
)

// GetIsiNodePools queries the node pools on the cluster
func GetIsiNodePools(
	ctx context.Context,
	client api.Client,
) (*IsiNodePoolsResp, error) {
	var resp IsiNodePoolsResp
	err := client.Get(ctx, nodePoolsPath, "", nil, nil, &resp)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

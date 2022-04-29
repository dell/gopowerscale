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
package goisilon

import (
	"context"

	apiv3 "github.com/dell/goisilon/api/v3"
)

// Stats is Isilon statistics data structure .
type Stats *apiv3.IsiStatsResp
type Clients *apiv3.ExportClientList

//GetStatistics returns statistics from Isilon. Keys indicate type of statistics expected
func (c *Client) GetStatistics(
	ctx context.Context,
	keys []string) (Stats, error) {

	stats, err := apiv3.GetIsiStats(ctx, c.API, keys)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// IsIOinProgress checks whether a volume on a node has IO in progress
func (c *Client) IsIOinProgress(
	ctx context.Context) (Clients, error) {

	// query the volume without using the metadata parameter, use whether an error (typically, JSONError instance with "404 Not Found" status code) is returned to indicate whether the volume already exists.
	stats, err := apiv3.IsIOinProgress(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return stats, nil
}

// ClusterConfig represents the configuration of cluster in k8s (namespace API).
type ClusterConfig *apiv3.IsiClusterConfig

// GetClusterConfig returns information about the configuration of cluster
func (c *Client) GetClusterConfig(ctx context.Context) (ClusterConfig, error) {
	clusterConfig, err := apiv3.GetIsiClusterConfig(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

// GetLocalSerial returns the local serial which is the serial number of cluster
func (c *Client) GetLocalSerial(ctx context.Context) (string, error) {
	clusterConfig, err := c.GetClusterConfig(ctx)
	if err != nil {
		return "", err
	}
	return clusterConfig.LocalSerial, nil
}

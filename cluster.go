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

	apiv14 "github.com/dell/goisilon/api/v14"
	apiv3 "github.com/dell/goisilon/api/v3"
	apiv7 "github.com/dell/goisilon/api/v7"
)

// Stats is Isilon statistics data structure .
type Stats *apiv3.IsiStatsResp
type FloatStats *apiv3.IsiFloatStatsResp
type Clients *apiv3.ExportClientList

// ClusterConfig represents the configuration of cluster in k8s (namespace API).
type ClusterConfig *apiv3.IsiClusterConfig
type ClusterIdentity *apiv3.IsiClusterIdentity
type ClusterNodes *apiv3.IsiClusterNodes
type ClusterAcs *apiv14.IsiClusterAcs
type ClusterInternalNetworks *apiv7.IsiClusterInternalNetworks

// GetStatistics returns statistics from Isilon. Keys indicate type of statistics expected
func (c *Client) GetStatistics(
	ctx context.Context,
	keys []string) (Stats, error) {

	stats, err := apiv3.GetIsiStats(ctx, c.API, keys)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// GetFloatStatistics returns float statistics from Isilon. Keys indicate type of statistics expected
func (c *Client) GetFloatStatistics(
	ctx context.Context,
	keys []string) (FloatStats, error) {

	stats, err := apiv3.GetIsiFloatStats(ctx, c.API, keys)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

// IsIOinProgress checks whether a volume on a node has IO in progress
func (c *Client) IsIOInProgress(
	ctx context.Context) (Clients, error) {

	// query the volume without using the metadata parameter, use whether an error (typically, JSONError instance with "404 Not Found" status code) is returned to indicate whether the volume already exists.
	clientList, err := apiv3.IsIOInProgress(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clientList, nil
}

// GetClusterConfig returns information about the configuration of cluster
func (c *Client) GetClusterConfig(ctx context.Context) (ClusterConfig, error) {
	clusterConfig, err := apiv3.GetIsiClusterConfig(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterConfig, nil
}

// GetClusterIdentity returns the login information
func (c *Client) GetClusterIdentity(ctx context.Context) (ClusterIdentity, error) {
	clusterIdentity, err := apiv3.GetIsiClusterIdentity(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterIdentity, nil
}

// GetClusterAcs returns the ACS status
func (c *Client) GetClusterAcs(ctx context.Context) (ClusterAcs, error) {
	clusterAcs, err := apiv14.GetIsiClusterAcs(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterAcs, nil
}

// GetClusterInternalNetworks internal networks settings
func (c *Client) GetClusterInternalNetworks(ctx context.Context) (ClusterInternalNetworks, error) {
	clusterInternalNetworks, err := apiv7.GetIsiClusterInternalNetworks(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterInternalNetworks, nil
}

// GetClusterNodes list the nodes on this cluster
func (c *Client) GetClusterNodes(ctx context.Context) (ClusterNodes, error) {
	clusterNodes, err := apiv3.GetIsiClusterNodes(ctx, c.API)
	if err != nil {
		return nil, err
	}
	return clusterNodes, nil
}

// GetClusterNode retrieves one node on this cluster
func (c *Client) GetClusterNode(ctx context.Context, nodeID int) (ClusterNodes, error) {
	clusterNodes, err := apiv3.GetIsiClusterNode(ctx, c.API, nodeID)
	if err != nil {
		return nil, err
	}
	return clusterNodes, nil
}

// GetLocalSerial returns the local serial which is the serial number of cluster
func (c *Client) GetLocalSerial(ctx context.Context) (string, error) {
	clusterConfig, err := c.GetClusterConfig(ctx)
	if err != nil {
		return "", err
	}
	return clusterConfig.LocalSerial, nil
}

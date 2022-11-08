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
package v3

import (
	"context"
	"strings"

	"github.com/dell/goisilon/api"
)

// constants
const (
	platfromStatsPath = "platform/3/statistics"
)

// GetIsiStats queries the attributes of a volume on the cluster
func GetIsiStats(
	ctx context.Context,
	client api.Client,
	keys []string) (resp *IsiStatsResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/3/statistics/current?keys=ifs.bytes.avail

	keysStr := strings.Join(keys, ",")
	statsOv := api.OrderedValues{{[]byte("keys"), []byte(keysStr)}}

	err = client.Get(
		ctx,
		string(platfromStatsPath),
		"current",
		statsOv,
		nil,
		&resp)

	return resp, err
}

// GetIsiFloatStats queries the float attributes of a cluster
func GetIsiFloatStats(
	ctx context.Context,
	client api.Client,
	keys []string) (resp *IsiFloatStatsResp, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/3/statistics/current?keys=ifs.bytes.avail

	keysStr := strings.Join(keys, ",")
	statsOv := api.OrderedValues{{[]byte("keys"), []byte(keysStr)}}

	err = client.Get(
		ctx,
		string(platfromStatsPath),
		"current",
		statsOv,
		nil,
		&resp)

	return resp, err
}

// IsIOInProgress returns the list of clients currently performing IO on the particular array
func IsIOInProgress(ctx context.Context,
	client api.Client) (resp *ExportClientList, err error) {
	err = client.Get(
		ctx, string(platfromStatsPath), "summary/client", api.OrderedValues{
			{[]byte("numeric"), []byte("true")}, // numeric=true returns the response faster since it does not performs reverse lookup and returns IP addresses
		},
		nil,
		&resp)
	return resp, err
}

// GetIsiClusterConfig queries the config information of OneFS cluster
func GetIsiClusterConfig(
	ctx context.Context,
	client api.Client) (clusterConfig *IsiClusterConfig, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/3/cluster/config
	// This will return configuration information of the cluster
	var clusterConfigResp IsiClusterConfig
	err = client.Get(ctx, clusterConfigPath, "", nil, nil, &clusterConfigResp)
	if err != nil {
		return nil, err
	}

	return &clusterConfigResp, nil
}

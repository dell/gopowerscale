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
	"testing"
)

func TestGetStatistics(*testing.T) {
	keyArray := []string{"ifs.bytes.avail", "ifs.bytes.total"}
	stats, err := client.GetStatistics(defaultCtx, keyArray)
	if err != nil || len(stats.StatsList) != 2 {
		panic("Couldn't get statistics.")
	}
	if stats.StatsList[0].Value <= 0 {
		panic("Statistics returned bad value.")
	}
}

func TestGetFloatStatistics(*testing.T) {
	floatStatsKeyArray := []string{"cluster.disk.bytes.in.rate", "ifs.bytes.total", "cluster.disk.xfers.in.rate"}
	stats, err := client.GetFloatStatistics(defaultCtx, floatStatsKeyArray)
	if err != nil || len(stats.StatsList) != 3 {
		panic("Couldn't get float statistics.")
	}
	if stats.StatsList[0].Value <= 0 {
		panic("Statistics returned bad value.")
	}
}

// Test if the local serial can be returned normally
func TestGetLocalSerial(t *testing.T) {
	// Get local serial
	localSerial, err := client.GetLocalSerial(defaultCtx)
	if err != nil {
		panic(err)
	}
	println(localSerial)
}

func TestGetClusterConfig(t *testing.T) {
	config, err := client.GetClusterConfig(defaultCtx)
	assertNoError(t, err)
	assertNotNil(t, config.OnefsVersion)
	assertNotNil(t, config.Timezone)
}

func TestGetClusterIdentity(t *testing.T) {
	identity, err := client.GetClusterIdentity(defaultCtx)
	assertNoError(t, err)
	assertNotNil(t, identity)
	assertNotNil(t, identity.Name)
}

func TestGetClusterAcs(t *testing.T) {
	acs, err := client.GetClusterAcs(defaultCtx)
	assertNoError(t, err)
	assertNotNil(t, acs)
}

func TestGetClusterInternalNetworks(t *testing.T) {
	networks, err := client.GetClusterInternalNetworks(defaultCtx)
	assertNoError(t, err)
	assertNotNil(t, networks)
}

func TestGetClusterNodes(t *testing.T) {
	nodes, err := client.GetClusterNodes(defaultCtx)
	assertNoError(t, err)
	assertNotNil(t, nodes)
	assertNotEqual(t, int(*nodes.Total), 0)
}

func TestGetClusterNode(t *testing.T) {
	nodeID := 1
	nodes, err := client.GetClusterNode(defaultCtx, nodeID)
	assertNoError(t, err)
	assertNotNil(t, nodes)
	assertEqual(t, int(*nodes.Total), 1)
	assertLen(t, nodes.Nodes, 1)
}

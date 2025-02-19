/*
Copyright (c) 2022-2025 Dell Inc, or its subsidiaries.

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
	"errors"
	"testing"

	apiv3 "github.com/dell/goisilon/api/v3"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetStatistics(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	keyArray := []string{"ifs.bytes.avail", "ifs.bytes.total"}
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv3.IsiStatsResp)
		stats := &apiv3.IsiStats{
			Key: "ifs.bytes.avail",
		}
		*resp = &apiv3.IsiStatsResp{
			StatsList: []*apiv3.IsiStats{
				stats,
			},
		}
	}).Once()
	_, err := client.GetStatistics(defaultCtx, keyArray)
	if err == nil {
		assert.Equal(t, nil, err)
	}

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get statistics")).Once()
	_, err = client.GetStatistics(defaultCtx, keyArray)
	if err != nil {
		assert.Equal(t, errors.New("failed to get statistics"), err)
	}
}

func TestGetFloatStatistics(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	floatStatsKeyArray := []string{"cluster.disk.bytes.in.rate", "ifs.bytes.total", "cluster.disk.xfers.in.rate"}
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv3.IsiFloatStatsResp)
		*resp = &apiv3.IsiFloatStatsResp{}
	}).Once()
	_, err := client.GetFloatStatistics(defaultCtx, floatStatsKeyArray)
	if err == nil {
		assert.Equal(t, nil, err)
	}

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get float statistics")).Once()
	_, err = client.GetFloatStatistics(defaultCtx, floatStatsKeyArray)
	if err != nil {
		assert.Equal(t, errors.New("failed to get float statistics"), err)
	}
}

func TestIsIOInProgress(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv3.ExportClientList)
		*resp = &apiv3.ExportClientList{}
	}).Once()
	_, err := client.IsIOInProgress(defaultCtx)
	if err == nil {
		assert.Equal(t, nil, err)
	}

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("io in progress failed")).Once()
	_, err = client.IsIOInProgress(defaultCtx)
	if err != nil {
		assert.Equal(t, errors.New("io in progress failed"), err)
	}
}

// Test if the local serial can be returned normally
func TestGetLocalSerial(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Get local serial
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetLocalSerial(defaultCtx)
	if err != nil {
		assert.Equal(t, nil, err)
	}

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get local serial")).Once()
	_, err = client.GetLocalSerial(defaultCtx)
	if err != nil {
		assert.Equal(t, errors.New("failed to get local serial"), err)
	}
}

func TestGetClusterConfig(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterConfig(defaultCtx)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster config")).Once()
	_, err = client.GetClusterConfig(defaultCtx)
	assert.Equal(t, errors.New("failed to get cluster config"), err)
}

func TestGetClusterIdentity(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterIdentity(defaultCtx)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster indentity")).Once()
	_, err = client.GetClusterIdentity(defaultCtx)
	assert.Equal(t, errors.New("failed to get cluster indentity"), err)
}

func TestGetClusterAcs(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterAcs(defaultCtx)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster acls")).Once()
	_, err = client.GetClusterAcs(defaultCtx)
	assert.Equal(t, errors.New("failed to get cluster acls"), err)
}

func TestGetClusterInternalNetworks(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterInternalNetworks(defaultCtx)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster internal networks")).Once()
	_, err = client.GetClusterInternalNetworks(defaultCtx)
	assert.Equal(t, errors.New("failed to get cluster internal networks"), err)
}

func TestGetClusterNodes(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterNodes(defaultCtx)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster nodes")).Once()
	_, err = client.GetClusterNodes(defaultCtx)
	assert.Equal(t, errors.New("failed to get cluster nodes"), err)
}

func TestGetClusterNode(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	nodeID := 1
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetClusterNode(defaultCtx, nodeID)
	assert.Equal(t, nil, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("failed to get cluster node")).Once()
	_, err = client.GetClusterNode(defaultCtx, nodeID)
	assert.Equal(t, errors.New("failed to get cluster node"), err)
}

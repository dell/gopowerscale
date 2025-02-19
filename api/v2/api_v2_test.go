/*
Copyright (c) 2025 Dell Inc, or its subsidiaries.

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

package v2

import (
	"testing"

	"github.com/dell/goisilon/api"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

const (
	testVolumePath = "/ifs/test/volume"
)

type Client struct {
	// API is the underlying OneFS API client.
	API api.Client
}

func TestRealNamespacePath(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Once()
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()

	expected := "namespace" + testVolumePath
	actual := realNamespacePath(client)

	assert.Equal(t, expected, actual)
}

func TestRealExportsPath(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Once()
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()

	expected := "platform/2/protocols/nfs/exports" + testVolumePath
	actual := realExportsPath(client)

	assert.Equal(t, expected, actual)
}

func TestRealVolumeSnapshotPath(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Once()
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()

	expected := "namespace/ifs/.snapshot/snapshotName/test/volume"
	actual := realVolumeSnapshotPath(client, "snapshotName")

	assert.Equal(t, expected, actual)
}

func TestGetAbsoluteSnapshotPath(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Once()
	client.On("VolumePath", anyArgs...).Return(testVolumePath).Once()

	expected := "/ifs/.snapshot/snapshotName/test/volume"
	actual := GetAbsoluteSnapshotPath(client, "snapshotName", "volumeName")

	assert.Equal(t, expected, actual)
}

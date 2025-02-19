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
package v1

import (
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetRealVolumeSnapshotPathWithIsiPath(t *testing.T) {
	value := GetRealVolumeSnapshotPathWithIsiPath("abc/ifs/xyz", "zonepath/abc", "name", "System")
	assert.Equal(t, "namespace/abc/zonepath/abc/.snapshot/name/xyz", value)

	value = GetRealVolumeSnapshotPathWithIsiPath("abc/ifs/xyz", "zonepath/abc", "name", "")
	assert.Equal(t, "namespace/zonepath/abc/.snapshot/name", value)
}

func TestGetAbsoluteSnapshotPath(t *testing.T) {
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	client.On("VolumePath", anyArgs...).Return("").Twice()
	value := GetAbsoluteSnapshotPath(client, "snapshotName", "volumeName", "zonepath/abc")
	assert.Equal(t, "zonepath/abc/.snapshot/snapshotName", value)
}

func TestRealVolumeSnapshotPath(t *testing.T) {
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	client.On("VolumesPath", anyArgs...).Return("xyz/ifs/abc").Twice()
	value := realVolumeSnapshotPath(client, "snapshotName", "volumeName", "zonepath/abc")
	assert.Equal(t, "namespace/xyz/volumeName/.snapshot/snapshotName/abc", value)
}

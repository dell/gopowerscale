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
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetVolumeACL(t *testing.T) {
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	_, err := client.GetVolumeACL(defaultCtx, "test_get_volume_acl")
	assert.Nil(t, err)
}

func TestSetVolumeOwnerToCurrentUser(t *testing.T) {
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("User", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetVolumeOwnerToCurrentUser(defaultCtx, "test_set_volume_owner")
	assert.Nil(t, err)
}

func TestSetVolumeOwner(t *testing.T) {
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetVolumeOwner(defaultCtx, "test_set_volume_owner", "rexray")
	assert.Nil(t, err)
}

func TestSetVolumeMode(t *testing.T) {
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetVolumeMode(defaultCtx, "test_set_volume_owner", 777)
	assert.Nil(t, err)
}

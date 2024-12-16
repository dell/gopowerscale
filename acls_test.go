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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	api "github.com/dell/goisilon/api/v2"
)

func TestGetVolumeACL(t *testing.T) {
	volumeName := "test_get_volume_acl"

	// make sure the volume exists
	client.CreateVolume(defaultCtx, volumeName)
	volume, err := client.GetVolume(defaultCtx, volumeName, volumeName)
	assertNoError(t, err)
	assertNotNil(t, volume)

	defer client.DeleteVolume(defaultCtx, volume.Name)

	username := client.API.User()
	user, err := client.GetUserByNameOrUID(defaultCtx, &username, nil)
	assertNoError(t, err)
	assertNotNil(t, user)

	acl, err := client.GetVolumeACL(defaultCtx, volume.Name)
	assertNoError(t, err)
	assertNotNil(t, acl)

	assertNotNil(t, acl.Owner)
	assertNotNil(t, acl.Owner.Name)
	assert.Equal(t, user.Name, *acl.Owner.Name)
	assertNotNil(t, acl.Owner.Type)
	assert.Equal(t, api.PersonaTypeUser, *acl.Owner.Type)
	assertNotNil(t, acl.Owner.ID)
	assert.Equal(t, user.OnDiskUserIdentity.ID, fmt.Sprintf("UID:%s", acl.Owner.ID.ID))
	assert.Equal(t, api.PersonaIDTypeUID, acl.Owner.ID.Type)
}

func TestSetVolumeOwnerToCurrentUser(t *testing.T) {
	volumeName := "test_set_volume_owner"

	// make sure the volume exists
	client.CreateVolume(defaultCtx, volumeName)
	volume, err := client.GetVolume(defaultCtx, volumeName, volumeName)
	assertNoError(t, err)
	assertNotNil(t, volume)

	defer client.DeleteVolume(defaultCtx, volume.Name)

	username := client.API.User()
	user, err := client.GetUserByNameOrUID(defaultCtx, &username, nil)
	assertNoError(t, err)
	assertNotNil(t, user)

	acl, err := client.GetVolumeACL(defaultCtx, volume.Name)
	assertNoError(t, err)
	assertNotNil(t, acl)

	assertNotNil(t, acl.Owner)
	assertNotNil(t, acl.Owner.Name)
	assert.Equal(t, user.Name, *acl.Owner.Name)
	assertNotNil(t, acl.Owner.Type)
	assert.Equal(t, api.PersonaTypeUser, *acl.Owner.Type)
	assertNotNil(t, acl.Owner.ID)
	assert.Equal(t, user.OnDiskUserIdentity.ID, fmt.Sprintf("UID:%s", acl.Owner.ID.ID))
	assert.Equal(t, api.PersonaIDTypeUID, acl.Owner.ID.Type)

	err = client.SetVolumeOwner(defaultCtx, volume.Name, "rexray")
	if err != nil {
		t.Skipf("Unable to change volume owner: %s - is efs.bam.chown_unrestricted set?", err)
	}
	assertNoError(t, err)

	acl, err = client.GetVolumeACL(defaultCtx, volume.Name)
	assertNoError(t, err)
	assertNotNil(t, acl)

	assertNotNil(t, acl.Owner)
	assertNotNil(t, acl.Owner.Name)
	assert.Equal(t, "rexray", *acl.Owner.Name)
	assertNotNil(t, acl.Owner.Type)
	assert.Equal(t, api.PersonaTypeUser, *acl.Owner.Type)
	assertNotNil(t, acl.Owner.ID)
	assert.Equal(t, "2000", acl.Owner.ID.ID)
	assert.Equal(t, api.PersonaIDTypeUID, acl.Owner.ID.Type)

	err = client.SetVolumeOwnerToCurrentUser(defaultCtx, volume.Name)
	assertNoError(t, err)

	acl, err = client.GetVolumeACL(defaultCtx, volume.Name)
	assertNoError(t, err)
	assertNotNil(t, acl)

	assertNotNil(t, acl.Owner)
	assertNotNil(t, acl.Owner.Name)
	assert.Equal(t, client.API.User(), *acl.Owner.Name)
	assertNotNil(t, acl.Owner.Type)
	assert.Equal(t, api.PersonaTypeUser, *acl.Owner.Type)
	assertNotNil(t, acl.Owner.ID)
	assert.Equal(t, "10", acl.Owner.ID.ID)
	assert.Equal(t, api.PersonaIDTypeUID, acl.Owner.ID.Type)
}

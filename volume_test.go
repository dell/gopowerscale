/*
Copyright (c) 2019-2025 Dell Inc, or its subsidiaries.

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
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"testing"

	apiv1 "github.com/dell/goisilon/api/v1"
	apiv2 "github.com/dell/goisilon/api/v2"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	isiPath                     = "/ifs/data/csi"
	isiVolumePathPermissions    = "077"
	volumeName                  = "test_get_create_volume_name"
	sourceVolumeName            = "test_copy_source_volume_name"
	destinationVolumeName       = "test_copy_destination_volume_name"
	subDirectoryName            = "test_sub_directory"
	sourceSubDirectoryPath      = fmt.Sprintf("%s/%s", sourceVolumeName, subDirectoryName)
	destinationSubDirectoryPath = fmt.Sprintf("%s/%s", destinationVolumeName, subDirectoryName)
	dirPath                     = "dA/dAA/dAAA"
)

// TODO - As part of PR job runs, observing GetVolumes, is not returning updated number of volumes
// hence the reason commented, this changed seems to be mostly related to PowerScale upgrade from 8.1.2.0
// 8.1.3.0
/*func TestVolumeList(*testing.T) {
	volumeName1 := "test_get_volumes_name1"
	volumeName2 := "test_get_volumes_name2"

	// identify all volumes on the cluster
	volumeMap := make(map[string]bool)
	volumes, err := client.GetVolumes(defaultCtx)
	if err != nil {
		panic(err)
	}
	for _, volume := range volumes {
		volumeMap[volume.Name] = true
	}
	initialVolumeCount := len(volumes)

	// Add the test volumes
	testVolume1, err := client.CreateVolume(defaultCtx, volumeName1)
	if err != nil {
		panic(err)
	}
	testVolume2, err := client.CreateVolume(defaultCtx, volumeName2)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName1)
	defer client.DeleteVolume(defaultCtx, volumeName2)

	// get the updated volume list
	volumes, err = client.GetVolumes(defaultCtx)
	if err != nil {
		panic(err)
	}

	// verify that the new volumes are there as well as all the old volumes.
	if len(volumes) != initialVolumeCount+2 {
		panic(fmt.Sprintf("Incorrect number of volumes.  Expected: %d Actual: %d\n", initialVolumeCount+2, len(volumes)))
	}
	// remove the original volumes and add the new ones.  in the end, we
	// should only have the volumes we just created and nothing more.
	for _, volume := range volumes {
		if _, found := volumeMap[volume.Name]; found == true {
			// this volume existed prior to the test start
			delete(volumeMap, volume.Name)
		} else {
			// this volume is new
			volumeMap[volume.Name] = true
		}
	}
	if len(volumeMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new volumes.  Expected: 2 Actual: %d\n", len(volumeMap)))
	}
	if _, found := volumeMap[testVolume1.Name]; found == false {
		panic(fmt.Sprintf("testVolume1 was not in the volume list\n"))
	}
	if _, found := volumeMap[testVolume2.Name]; found == false {
		panic(fmt.Sprintf("testVolume2 was not in the volume list\n"))
	}

}*/

func TestGetVolume(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Volume exists - with id
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err := client.GetVolume(defaultCtx, "testVolumeId", "testVolumeName")
	assert.Nil(t, err)
	assert.Equal(t, "testVolumeId", testVolume.Name)

	// Test case: Volume exists - without id
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err = client.GetVolume(defaultCtx, "", "testVolumeName")
	assert.Nil(t, err)
	assert.Equal(t, "testVolumeName", testVolume.Name)

	// Test case: Volume does not exist
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err = client.GetVolume(defaultCtx, "nonExistentVolumeID", "nonExistentVolumeName")
	assert.NotNil(t, err)
}

func TestGetVolumeWithIsiPath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Volume exists - with id
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err := client.GetVolumeWithIsiPath(defaultCtx, isiPath, "testVolumeId", "testVolumeName")
	assert.Nil(t, err)
	assert.Equal(t, "testVolumeId", testVolume.Name)

	// Test case: Volume exists - without id
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err = client.GetVolumeWithIsiPath(defaultCtx, isiPath, "", "testVolumeName")
	assert.Nil(t, err)
	assert.Equal(t, "testVolumeName", testVolume.Name)

	// Test case: Volume does not exist
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err = client.GetVolumeWithIsiPath(defaultCtx, isiPath, "nonExistentVolumeID", "nonExistentVolumeName")
	assert.NotNil(t, err)
}

func TestIsVolumeExistent(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Volume exists
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	isExistent := client.IsVolumeExistent(defaultCtx, "volumeId", "volumeName")
	assert.True(t, isExistent)

	// Test case: Volume does not exist
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	isExistent = client.IsVolumeExistent(defaultCtx, "volumeId", "volumeName")
	assert.False(t, isExistent)
}

func TestIsVolumeExistentWithIsiPath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Volume exists
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	isExistent := client.IsVolumeExistentWithIsiPath(defaultCtx, isiPath, "volumeId", "volumeName")
	assert.True(t, isExistent)

	// Test case: Volume does not exist
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	isExistent = client.IsVolumeExistentWithIsiPath(defaultCtx, isiPath, "volumeId", "volumeName")
	assert.False(t, isExistent)
}

func TestGetVolumes(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test case: Volume exists - with id
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumesResp)
		*resp = &apiv1.GetIsiVolumesResp{}
	}).Once()
	testVolumes, err := client.GetVolumes(defaultCtx)
	assert.Nil(t, err)
	assert.Empty(t, testVolumes)

	// Test case: Volume exists - without id
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumesResp)
		*resp = &apiv1.GetIsiVolumesResp{
			Children: []*apiv1.VolumeName{{Name: "testVolume"}},
		}
	}).Once()
	testVolumes, err = client.GetVolumes(defaultCtx)
	assert.Nil(t, err)
	assert.Equal(t, "testVolume", testVolumes[0].Name)

	// Test case: Volume does not exist
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	testVolumes, err = client.GetVolumes(defaultCtx)
	assert.NotNil(t, err)
	assert.Nil(t, testVolumes)
}

func TestCreateVolume(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	volumeName := "test_get_create_volume_name"
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	testVolume, err := client.CreateVolume(defaultCtx, volumeName)
	assert.Nil(t, err)
	assert.Equal(t, volumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume creation failed")).Once()
	testVolume, err = client.CreateVolume(defaultCtx, volumeName)
	assert.ErrorContains(t, err, "volume creation failed")
	assert.Nil(t, testVolume)
}

func TestCreateVolumeWithIsipath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	volumeName := "test_get_create_volume_name"
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	testVolume, err := client.CreateVolumeWithIsipath(defaultCtx, isiPath, volumeName, isiVolumePathPermissions)
	assert.Nil(t, err)
	assert.Equal(t, volumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume creation failed")).Once()
	testVolume, err = client.CreateVolumeWithIsipath(defaultCtx, isiPath, volumeName, isiVolumePathPermissions)
	assert.ErrorContains(t, err, "volume creation failed")
	assert.Nil(t, testVolume)
}

func TestCreateVolumeWithIsipathMetaData(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	volumeName := "test_get_create_volume_name"
	testHeader := map[string]string{
		"x-csi-pv-name":      "pv-name",
		"x-csi-pv-claimname": "pv-claimname",
		"x-csi-pv-namespace": "pv-namesace",
	}

	// success
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	testVolume, err := client.CreateVolumeWithIsipathMetaData(defaultCtx, isiPath, volumeName, isiVolumePathPermissions, testHeader)
	assert.Nil(t, err)
	assert.Equal(t, volumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume creation failed")).Once()
	testVolume, err = client.CreateVolumeWithIsipathMetaData(defaultCtx, isiPath, volumeName, isiVolumePathPermissions, testHeader)
	assert.ErrorContains(t, err, "volume creation failed")
	assert.Nil(t, testVolume)
}

func TestCreateVolumeNoACL(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	volumeName := "test_get_create_volume_name"
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	testVolume, err := client.CreateVolumeNoACL(defaultCtx, volumeName)
	assert.Nil(t, err)
	assert.Equal(t, volumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume creation failed")).Once()
	testVolume, err = client.CreateVolumeNoACL(defaultCtx, volumeName)
	assert.ErrorContains(t, err, "volume creation failed")
	assert.Nil(t, testVolume)
}

func TestDeleteVolume(t *testing.T) {
	volumeName := "test_remove_volume_name"
	// remove the volume
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err = client.DeleteVolume(defaultCtx, volumeName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	err = client.DeleteVolume(defaultCtx, volumeName)
	assert.ErrorContains(t, err, "not found")
}

func TestDeleteIsiVolumeWithIsiPath(t *testing.T) {
	volumeName := "test_remove_volume_name"
	// remove the volume
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err = client.DeleteVolumeWithIsiPath(defaultCtx, isiPath, volumeName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	err = client.DeleteVolumeWithIsiPath(defaultCtx, isiPath, volumeName)
	assert.ErrorContains(t, err, "not found")
}

func TestCopyVolume(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err := client.CopyVolume(defaultCtx, sourceVolumeName, destinationVolumeName)
	assert.Nil(t, err)
	assert.Equal(t, destinationVolumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume copy failed")).Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	testVolume, err = client.CopyVolume(defaultCtx, sourceVolumeName, destinationVolumeName)
	assert.ErrorContains(t, err, "volume copy failed")
	assert.Nil(t, testVolume)
}

func TestCopyVolumeWithIsiPath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()
	testVolume, err := client.CopyVolumeWithIsiPath(defaultCtx, isiPath, sourceVolumeName, destinationVolumeName)
	assert.Nil(t, err)
	assert.Equal(t, destinationVolumeName, testVolume.Name)

	// negative
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume copy failed")).Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	testVolume, err = client.CopyVolumeWithIsiPath(defaultCtx, isiPath, sourceVolumeName, destinationVolumeName)
	assert.ErrorContains(t, err, "volume copy failed")
	assert.Nil(t, testVolume)
}

func TestExportVolume(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Times(3)
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err := client.ExportVolume(defaultCtx, "testing")
	assert.Nil(t, err)
}

func TestExportVolumeWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err := client.ExportVolumeWithZone(defaultCtx, sourceVolumeName, "zone", "description")
	assert.Nil(t, err)
}

func TestExportVolumeWithZoneAndPath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err := client.ExportVolumeWithZoneAndPath(defaultCtx, isiPath, "zone", "description")
	assert.Nil(t, err)
}

func TestUnexportVolume(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Times(3)
	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(nil).Once()
	err := client.UnexportVolume(defaultCtx, "testing")
	assert.Nil(t, err)
}

func TestQueryVolumeChildren(t *testing.T) {
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	_, err := client.QueryVolumeChildren(defaultCtx, "testing")
	assert.Nil(t, err)
}

func TestCreateVolumeDir(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// success
	newDirMode := apiv2.FileMode(0o755)
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.CreateVolumeDir(defaultCtx, volumeName, dirPath, os.FileMode(newDirMode), false, false)
	assert.Nil(t, err)

	// negative
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("volume creation failed")).Once()
	err = client.CreateVolumeDir(defaultCtx, volumeName, dirPath, os.FileMode(newDirMode), false, false)
	assert.ErrorContains(t, err, "volume creation failed")
}

func TestGetVolumeExportMap(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumesResp)
		*resp = &apiv1.GetIsiVolumesResp{}
	}).Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	testVolumes, err := client.GetVolumeExportMap(defaultCtx, false)
	assert.Nil(t, err)
	assert.Empty(t, testVolumes)
}

type bufReadCloser struct {
	b *bytes.Buffer
}

func (b *bufReadCloser) Read(p []byte) (n int, err error) {
	return b.b.Read(p)
}

func (b *bufReadCloser) Close() error {
	return nil
}

func VolumeQueryChildrenTest(t *testing.T) {
	// TODO: Need to fix this as it is failing with Isilon 8.1
	skipTest(t)

	var (
		ctx = defaultCtx
		// context.WithValue(defaultCtx, log.LevelKey(), log.InfoLevel)

		err      error
		volume   Volume
		children VolumeChildrenMap

		volumeName   = "test_volume_query_children"
		dirPath0     = "dA"
		dirPath0a    = "dA/dAA"
		dirPath1     = "dA/dAA/dAAA"
		dirPath2     = "dB/dBB"
		dirPath3     = "dC"
		dirPath4     = "dC/dCC"
		fileName0    = "an_empty_file"
		fileName1    = fileName0
		fileName2    = fileName0
		fileName3    = fileName0
		volDirPath0  = path.Join(volumeName, dirPath0)
		volDirPath0a = path.Join(volumeName, dirPath0a)
		volDirPath1  = path.Join(volumeName, dirPath1)
		volDirPath2  = path.Join(volumeName, dirPath2)
		volDirPath3  = path.Join(volumeName, dirPath3)
		volDirPath4  = path.Join(volumeName, dirPath4)
		volFilePath0 = path.Join(volumeName, fileName0)
		volFilePath1 = path.Join(volDirPath2, fileName1)
		volFilePath2 = path.Join(volDirPath3, fileName2)
		volFilePath3 = path.Join(volDirPath4, fileName3)
		dirPath0Key  = path.Join(client.API.VolumePath(volumeName), dirPath0)
		dirPath0aKey = path.Join(client.API.VolumePath(volumeName), dirPath0a)
		dirPath1Key  = path.Join(client.API.VolumePath(volumeName), dirPath1)
		dirPath2Key  = path.Join(client.API.VolumePath(volumeName), dirPath2)
		dirPath3Key  = path.Join(client.API.VolumePath(volumeName), dirPath3)
		dirPath4Key  = path.Join(client.API.VolumePath(volumeName), dirPath4)
		filePath0Key = path.Join(client.API.VolumePath(volumeName), fileName0)
		filePath1Key = path.Join(dirPath2Key, fileName1)
		filePath2Key = path.Join(dirPath3Key, fileName2)
		filePath3Key = path.Join(dirPath4Key, fileName3)

		newUserName  = client.API.User()
		newGroupName = newUserName
		badUserID    = "999"
		badGroupID   = "999"
		badUserName  = "Unknown User"
		badGroupName = "Unknown Group"

		volChildCount = 9

		newDirMode  = apiv2.FileMode(0o755)
		newFileMode = apiv2.FileMode(0o644)
		badDirMode  = apiv2.FileMode(0o700)
		badFileMode = apiv2.FileMode(0o400)

		newDirACL = &apiv2.ACL{
			Action:        &apiv2.PActionTypeReplace,
			Authoritative: &apiv2.PAuthoritativeTypeMode,
			Owner: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   newUserName,
					Type: apiv2.PersonaIDTypeUser,
				},
			},
			Group: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   newGroupName,
					Type: apiv2.PersonaIDTypeGroup,
				},
			},
			Mode: &newDirMode,
		}

		newFileACL = &apiv2.ACL{
			Action:        &apiv2.PActionTypeReplace,
			Authoritative: &apiv2.PAuthoritativeTypeMode,
			Owner: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   newUserName,
					Type: apiv2.PersonaIDTypeUser,
				},
			},
			Group: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   newGroupName,
					Type: apiv2.PersonaIDTypeGroup,
				},
			},
			Mode: &newFileMode,
		}

		badDirACL = &apiv2.ACL{
			Action:        &apiv2.PActionTypeReplace,
			Authoritative: &apiv2.PAuthoritativeTypeMode,
			Owner: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   badUserID,
					Type: apiv2.PersonaIDTypeUID,
				},
			},
			Group: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   badGroupID,
					Type: apiv2.PersonaIDTypeGID,
				},
			},
			Mode: &badDirMode,
		}

		badFileACL = &apiv2.ACL{
			Action:        &apiv2.PActionTypeReplace,
			Authoritative: &apiv2.PAuthoritativeTypeMode,
			Owner: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   badUserID,
					Type: apiv2.PersonaIDTypeUID,
				},
			},
			Group: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   badGroupID,
					Type: apiv2.PersonaIDTypeGID,
				},
			},
			Mode: &badFileMode,
		}
	)

	defer client.ForceDeleteVolume(ctx, volumeName)

	setACLsWithPaths := func(
		ctx context.Context, acl *apiv2.ACL, paths ...string,
	) {
		for _, p := range paths {
			assertNoError(
				t,
				apiv2.ACLUpdate(ctx, client.API, p, acl))
		}
	}

	assertNewFileACL := func(cm map[string]*apiv2.ContainerChild, k string) {
		if !assert.Equal(t, newFileMode, *cm[k].Mode) ||
			!assert.Equal(t, newUserName, *cm[k].Owner) ||
			!assert.Equal(t, newGroupName, *cm[k].Group) {
			t.FailNow()
		}
	}

	assertBadFileACL := func(cm map[string]*apiv2.ContainerChild, k string) {
		if !assert.Equal(t, badFileMode, *cm[k].Mode) ||
			!assert.Equal(t, badUserName, *cm[k].Owner) ||
			!assert.Equal(t, badGroupName, *cm[k].Group) {
			t.FailNow()
		}
	}

	assertNewDirACL := func(cm map[string]*apiv2.ContainerChild, k string) {
		if !assert.Equal(t, newDirMode, *cm[k].Mode) ||
			!assert.Equal(t, newUserName, *cm[k].Owner) ||
			!assert.Equal(t, newGroupName, *cm[k].Group) {
			t.FailNow()
		}
	}

	assertBadDirACL := func(cm map[string]*apiv2.ContainerChild, k string) {
		if !assert.Equal(t, badDirMode, *cm[k].Mode) ||
			!assert.Equal(t, badUserName, *cm[k].Owner) ||
			!assert.Equal(t, badGroupName, *cm[k].Group) {
			t.FailNow()
		}
	}

	createObjs := func(ctx context.Context, createType int) {
		// make sure the volume exists
		client.CreateVolume(ctx, volumeName)
		volume, err = client.GetVolume(ctx, volumeName, volumeName)
		assertNoError(t, err)
		assertNotNil(t, volume)

		// assert the volume has no children
		children, err = client.QueryVolumeChildren(ctx, volumeName)
		assertNoError(t, err)
		assertLen(t, children, 0)

		switch createType {
		case 0:
			// create dirPath1
			assertError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath1,
				os.FileMode(newDirMode),
				false,
				false))

			// create dirPath1 again, recursively
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath1,
				os.FileMode(newDirMode),
				false,
				true))

			// create the second directory path
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath2,
				os.FileMode(newDirMode),
				true,
				true))

			// create file0
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volumeName,
				fileName0,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				false))

			// create file1
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath2,
				fileName1,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				false))

			// try and create file1 again; should fail
			assertError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath2,
				fileName1,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				false))

			// try and create file1 again; should work
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath2,
				fileName1,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				true))

			// create the third directory path
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath4,
				os.FileMode(newDirMode),
				true,
				true))

			setACLsWithPaths(
				ctx, newDirACL,
				volDirPath0, volDirPath1, volDirPath2, volDirPath3, volDirPath4)
			setACLsWithPaths(ctx, newFileACL, volFilePath0, volFilePath1)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, volChildCount)
			assertNewDirACL(children, dirPath0Key)
			assertNewDirACL(children, dirPath1Key)
			assertNewDirACL(children, dirPath2Key)
			assertNewDirACL(children, dirPath3Key)
			assertNewDirACL(children, dirPath4Key)
			assertNewFileACL(children, filePath0Key)
			assertNewFileACL(children, filePath1Key)
		case 1:
			// test a single root dir with good perms that has an empty sub-dir
			// with bad perms
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath4,
				os.FileMode(newDirMode),
				true,
				true))
			setACLsWithPaths(ctx, newDirACL, volDirPath3, volDirPath4)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, 2)
			assertNewDirACL(children, dirPath3Key)
			assertNewDirACL(children, dirPath4Key)
		case 2:
			// test a single, empty root dir with bad perms
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath3,
				os.FileMode(newDirMode),
				true,
				true))
			setACLsWithPaths(ctx, newDirACL, volDirPath3)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, 1)
			assertNewDirACL(children, dirPath3Key)
		case 3:
			// test a single, root dir with bad perms that has a single file in
			// it where the file has good perms
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath3,
				os.FileMode(newDirMode),
				true,
				true))
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath3,
				fileName2,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				true))
			setACLsWithPaths(ctx, newFileACL, volFilePath2)
			setACLsWithPaths(ctx, newDirACL, volDirPath3)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, 2)
			assertNewDirACL(children, dirPath3Key)
			assertNewFileACL(children, filePath2Key)
		case 4:
			// test a single, root dir with bad perms that has a single sub-dir
			// with good perms, and the sub-dir contains a file with good perms
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath4,
				os.FileMode(newDirMode),
				true,
				true))
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath4,
				fileName3,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				true))
			setACLsWithPaths(ctx, newFileACL, volFilePath3)
			setACLsWithPaths(ctx, newDirACL, volDirPath3, volDirPath4)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, 3)
			assertNewDirACL(children, dirPath3Key)
			assertNewDirACL(children, dirPath4Key)
			assertNewFileACL(children, filePath3Key)
		case 5:
			// test a single, root dir with good perms that has a single sub-dir
			// with bad perms, and the sub-dir contains a file with good perms
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath4,
				os.FileMode(newDirMode),
				true,
				true))
			assertNoError(t, apiv2.ContainerCreateFile(
				ctx,
				client.API,
				volDirPath4,
				fileName3,
				0,
				newFileMode,
				&bufReadCloser{&bytes.Buffer{}},
				true))
			setACLsWithPaths(ctx, newFileACL, volFilePath3)
			setACLsWithPaths(ctx, newDirACL, volDirPath3, volDirPath4)
			children, err = client.QueryVolumeChildren(ctx, volumeName)
			assertNoError(t, err)
			assertLen(t, children, 3)
			assertNewDirACL(children, dirPath3Key)
			assertNewDirACL(children, dirPath4Key)
			assertNewFileACL(children, filePath3Key)
		case 6:
			// test /dA/dAA/dAAA where dA has bad perms; the volume delete
			// should fail
			assertNoError(t, client.CreateVolumeDir(
				ctx,
				volumeName,
				dirPath1,
				os.FileMode(newDirMode),
				true,
				true))
		}
	}

	// assert that ForceDeleteVolume works
	createObjs(ctx, 0)
	setACLsWithPaths(ctx, badFileACL, volFilePath0, volFilePath1)
	setACLsWithPaths(ctx,
		badDirACL, volDirPath4, volDirPath3, volDirPath1, volDirPath0)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, volChildCount)
	assertBadDirACL(children, dirPath0Key)
	assertBadDirACL(children, dirPath1Key)
	assertBadDirACL(children, dirPath3Key)
	assertBadDirACL(children, dirPath4Key)
	assertBadFileACL(children, filePath0Key)
	assertBadFileACL(children, filePath1Key)
	// force delete the volume
	assertNoError(t, client.ForceDeleteVolume(ctx, volumeName))

	// assert that an initial delete will result in the removal of files
	// and directories not in conflict
	createObjs(ctx, 0)
	setACLsWithPaths(ctx, badFileACL, volFilePath0, volFilePath1)
	setACLsWithPaths(ctx,
		badDirACL, volDirPath4, volDirPath3, volDirPath1, volDirPath0)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, volChildCount)
	assertBadDirACL(children, dirPath0Key)
	assertBadDirACL(children, dirPath1Key)
	assertBadDirACL(children, dirPath3Key)
	assertBadDirACL(children, dirPath4Key)
	assertBadFileACL(children, filePath0Key)
	assertBadFileACL(children, filePath1Key)
	// attempt to delete the volume; should fail, but the following paths
	// will have been removed:
	//
	// - /dB
	// - /dB/dBB
	// - /dB/dBB/an_empty_file
	// - /an_empty_file
	assertError(t, client.DeleteVolume(ctx, volumeName))
	setACLsWithPaths(
		ctx, newDirACL,
		volDirPath0, volDirPath1, volDirPath3, volDirPath4)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, volChildCount-4)
	assertNewDirACL(children, dirPath0Key)
	assertNewDirACL(children, dirPath1Key)
	assertNewDirACL(children, dirPath3Key)
	assertNewDirACL(children, dirPath4Key)
	// attempt to delete the volume; should succeed
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert that a file with bad perms deep in the hierarchy won't prevent
	// a delete
	createObjs(ctx, 0)
	setACLsWithPaths(ctx, badFileACL, volFilePath1)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, volChildCount)
	assertBadFileACL(children, filePath1Key)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert that a root file with bad perms will not prevent a delete
	createObjs(ctx, 0)
	setACLsWithPaths(ctx, badFileACL, volFilePath0)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, volChildCount)
	assertBadFileACL(children, filePath0Key)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert that a root-directory with an empty sub-dir with bad perms will
	// not prevent a delete
	createObjs(ctx, 1)
	setACLsWithPaths(ctx, badDirACL, volDirPath4)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 2)
	assertBadDirACL(children, dirPath4Key)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert that an empty root-directory with bad perms will not prevent a
	// delete
	createObjs(ctx, 2)
	setACLsWithPaths(ctx, badDirACL, volDirPath3)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 1)
	assertBadDirACL(children, dirPath3Key)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert that a root-directory with bad perms that contains a file with
	// good perms will prevent a delete
	createObjs(ctx, 3)
	setACLsWithPaths(ctx, badDirACL, volDirPath3)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 2)
	assertBadDirACL(children, dirPath3Key)
	assertNewFileACL(children, filePath2Key)
	assertError(t, client.DeleteVolume(ctx, volumeName))
	setACLsWithPaths(ctx, newDirACL, volDirPath3)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert the previous scenario will be handled by a force delete
	createObjs(ctx, 3)
	setACLsWithPaths(ctx, badDirACL, volDirPath3)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 2)
	assertBadDirACL(children, dirPath3Key)
	assertNewFileACL(children, filePath2Key)
	assertNoError(t, client.ForceDeleteVolume(ctx, volumeName))

	// assert that a root-directory with bad perms that contains a sub-dir
	// with good perms that contains a file with good perms prevents a
	// delete
	createObjs(ctx, 4)
	setACLsWithPaths(ctx, badDirACL, volDirPath3)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertBadDirACL(children, dirPath3Key)
	assertNewDirACL(children, dirPath4Key)
	assertNewFileACL(children, filePath3Key)
	assertError(t, client.DeleteVolume(ctx, volumeName))
	setACLsWithPaths(ctx, newDirACL, volDirPath3)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert the previous scenario will be handled by a force delete
	createObjs(ctx, 4)
	setACLsWithPaths(ctx, badDirACL, volDirPath3)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertBadDirACL(children, dirPath3Key)
	assertNewDirACL(children, dirPath4Key)
	assertNewFileACL(children, filePath3Key)
	assertNoError(t, client.ForceDeleteVolume(ctx, volumeName))

	// assert a single, root dir with good perms that has a single sub-dir
	// with bad perms, and the sub-dir contains a file with good perms prevents
	// a delete
	createObjs(ctx, 5)
	setACLsWithPaths(ctx, badDirACL, volDirPath4)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertNewDirACL(children, dirPath3Key)
	assertBadDirACL(children, dirPath4Key)
	assertNewFileACL(children, filePath3Key)
	assertError(t, client.DeleteVolume(ctx, volumeName))
	setACLsWithPaths(ctx, newDirACL, volDirPath4)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert the previous scenario will be handled by a force delete
	createObjs(ctx, 5)
	setACLsWithPaths(ctx, badDirACL, volDirPath4)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertNewDirACL(children, dirPath3Key)
	assertBadDirACL(children, dirPath4Key)
	assertNewFileACL(children, filePath3Key)
	assertNoError(t, client.ForceDeleteVolume(ctx, volumeName))

	// test /dA/dAA/dAAA where dA has bad perms; the volume delete
	// should fail
	createObjs(ctx, 6)
	setACLsWithPaths(ctx, badDirACL, volDirPath0a)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertNewDirACL(children, dirPath0Key)
	assertBadDirACL(children, dirPath0aKey)
	assertNewDirACL(children, dirPath1Key)
	assertError(t, client.DeleteVolume(ctx, volumeName))
	setACLsWithPaths(ctx, newDirACL, volDirPath0a)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertNoError(t, client.DeleteVolume(ctx, volumeName))

	// assert the previous scenario will be handled by a force delete
	createObjs(ctx, 6)
	setACLsWithPaths(ctx, badDirACL, volDirPath0a)
	children, err = client.QueryVolumeChildren(ctx, volumeName)
	assertNoError(t, err)
	assertLen(t, children, 3)
	assertNewDirACL(children, dirPath0Key)
	assertBadDirACL(children, dirPath0aKey)
	assertNewDirACL(children, dirPath1Key)
	assertNoError(t, client.ForceDeleteVolume(ctx, volumeName))
}

func TestGetVolumeSize(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeSizeResp)
		*resp = &apiv1.GetIsiVolumeSizeResp{}
	}).Once()
	size, err := client.GetVolumeSize(defaultCtx, isiPath, volumeName)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), size)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeSizeResp)
		*resp = &apiv1.GetIsiVolumeSizeResp{AttributeMap: []struct {
			Name string `json:"name"`
			Size int64  `json:"size"`
		}{
			{
				Name: "vol1",
				Size: 1024,
			},
			{
				Name: "vol2",
				Size: 512,
			},
		}}
	}).Once()
	size, err = client.GetVolumeSize(defaultCtx, isiPath, volumeName)
	assert.Nil(t, err)
	assert.Equal(t, int64(1536), size)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	size, err = client.GetVolumeSize(defaultCtx, isiPath, volumeName)
	assert.ErrorContains(t, err, "not found")
	assert.Equal(t, int64(0), size)
}

func TestForceDeleteVolume(t *testing.T) {
	// Test case: successful deletion
	client.API.(*mocks.Client).On("User", anyArgs[0:6]...).Return("user").Run(nil).Once()
	client.API.(*mocks.Client).On("VolumesPath", anyArgs[0:6]...).Return("").Run(nil).Times(3)
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(nil).Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err := client.ForceDeleteVolume(context.Background(), "testvol")
	assert.NoError(t, err)
}

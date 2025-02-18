/*
Copyright (c) 2019 Dell Inc, or its subsidiaries.

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
	"errors"
	"fmt"
	"os"
	"testing"

	apiv1 "github.com/dell/goisilon/api/v1"
	"github.com/dell/gounity/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSnapshotsGet(_ *testing.T) {
	snapshotPath := "test_snapshots_get_volume"
	snapshotName1 := "test_snapshots_get_snapshot_0"
	snapshotName2 := "test_snapshots_get_snapshot_1"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// identify all snapshots on the cluster
	snapshotMap := make(map[int64]string)
	snapshots, err := client.GetSnapshots(defaultCtx)
	if err != nil {
		panic(err)
	}
	for _, snapshot := range snapshots {
		snapshotMap[snapshot.ID] = snapshot.Name
	}
	initialSnapshotCount := len(snapshots)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}
	testSnapshot2, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName2)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)
	defer client.RemoveSnapshot(defaultCtx, testSnapshot2.ID, snapshotName2)

	// get the updated snapshot list
	snapshots, err = client.GetSnapshots(defaultCtx)
	if err != nil {
		panic(err)
	}

	// verify that the new snapshots are there as well as all the old snapshots.
	if len(snapshots) != initialSnapshotCount+2 {
		panic(fmt.Sprintf("Incorrect number of snapshots.  Expected: %d Actual: %d\n", initialSnapshotCount+2, len(snapshots)))
	}
	// remove the original snapshots and add the new ones.  in the end, we
	// should only have the snapshots we just created and nothing more.
	for _, snapshot := range snapshots {
		if _, found := snapshotMap[snapshot.ID]; found == true {
			// this snapshot existed prior to the test start
			delete(snapshotMap, snapshot.ID)
		} else {
			// this snapshot is new
			snapshotMap[snapshot.ID] = snapshot.Name
		}
	}
	if len(snapshotMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new snapshots.  Expected: 2 Actual: %d\n", len(snapshotMap)))
	}
	if _, found := snapshotMap[testSnapshot1.ID]; found == false {
		panic("testSnapshot1 was not in the snapshot list\n")
	}
	if _, found := snapshotMap[testSnapshot2.ID]; found == false {
		panic("testSnapshot2 was not in the snapshot list\n")
	}
}

func TestSnapshotsGetByPath(_ *testing.T) {
	snapshotPath1 := "test_snapshots_get_by_path_volume_1"
	snapshotPath2 := "test_snapshots_get_by_path_volume_2"
	snapshotName1 := "test_snapshots_get_by_path_snapshot_1"
	snapshotName2 := "test_snapshots_get_by_path_snapshot_2"
	snapshotName3 := "test_snapshots_get_by_path_snapshot_3"

	// create the two test volumes
	_, err := client.CreateVolume(defaultCtx, snapshotPath1)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath1)
	_, err = client.CreateVolume(defaultCtx, snapshotPath2)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath2)

	// identify all snapshots on the cluster
	snapshotMap := make(map[int64]string)
	snapshots, err := client.GetSnapshotsByPath(defaultCtx, snapshotPath1)
	if err != nil {
		panic(err)
	}
	for _, snapshot := range snapshots {
		snapshotMap[snapshot.ID] = snapshot.Name
	}
	initialSnapshotCount := len(snapshots)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath1, snapshotName1)
	if err != nil {
		panic(err)
	}
	testSnapshot2, err := client.CreateSnapshot(
		defaultCtx, snapshotPath2, snapshotName2)
	if err != nil {
		panic(err)
	}
	testSnapshot3, err := client.CreateSnapshot(
		defaultCtx, snapshotPath1, snapshotName3)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)
	defer client.RemoveSnapshot(defaultCtx, testSnapshot2.ID, snapshotName2)
	defer client.RemoveSnapshot(defaultCtx, testSnapshot3.ID, snapshotName3)

	// get the updated snapshot list
	snapshots, err = client.GetSnapshotsByPath(defaultCtx, snapshotPath1)
	if err != nil {
		panic(err)
	}

	// verify that the new snapshots in the given path are there as well
	// as all the old snapshots in that path.
	if len(snapshots) != initialSnapshotCount+2 {
		panic(fmt.Sprintf("Incorrect number of snapshots for path (%s).  Expected: %d Actual: %d\n", snapshotPath1, initialSnapshotCount+2, len(snapshots)))
	}
	// remove the original snapshots and add the new ones.  in the end, we
	// should only have the snapshots we just created and nothing more.
	for _, snapshot := range snapshots {
		if _, found := snapshotMap[snapshot.ID]; found == true {
			// this snapshot existed prior to the test start
			delete(snapshotMap, snapshot.ID)
		} else {
			// this snapshot is new
			snapshotMap[snapshot.ID] = snapshot.Name
		}
	}
	if len(snapshotMap) != 2 {
		panic(fmt.Sprintf("Incorrect number of new snapshots.  Expected: 2 Actual: %d\n", len(snapshotMap)))
	}
	if _, found := snapshotMap[testSnapshot1.ID]; found == false {
		panic("testSnapshot1 was not in the snapshot list\n")
	}
	if _, found := snapshotMap[testSnapshot3.ID]; found == false {
		panic("testSnapshot3 was not in the snapshot list\n")
	}
}

func TestSnapshotCreate(_ *testing.T) {
	snapshotPath := "test_snapshot_create_volume"
	snapshotName := "test_snapshot_create_snapshot"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// make sure the snapshot doesn't exist yet
	snapshot, err := client.GetSnapshot(defaultCtx, -1, snapshotName)
	if err == nil && snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) already exists.\n", snapshotName))
	}

	// Add the test snapshot
	testSnapshot, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot.ID, snapshotName)

	// get the updated snapshot list
	snapshot, err = client.GetSnapshot(
		defaultCtx, testSnapshot.ID, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot == nil {
		panic(fmt.Sprintf("Snapshot (%s) was not created.\n", snapshotName))
	}
	if snapshot.Name != snapshotName {
		panic(fmt.Sprintf("Snapshot name not set properly.  Expected: (%s) Actual: (%s)\n", snapshotName, snapshot.Name))
	}
	if snapshot.Path != client.API.VolumePath(snapshotPath) {
		panic(fmt.Sprintf("Snapshot path not set properly.  Expected: (%s) Actual: (%s)\n", snapshotPath, snapshot.Path))
	}
}

func TestSnapshotRemove(_ *testing.T) {
	snapshotPath := "test_snapshot_remove_volume"
	snapshotName := "test_snapshot_remove_snapshot"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// make sure the snapshot exists
	client.CreateSnapshot(defaultCtx, snapshotPath, snapshotName)
	snapshot, err := client.GetSnapshot(defaultCtx, -1, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot == nil {
		panic(fmt.Sprintf("Test not setup properly.  No test snapshot (%s).", snapshotName))
	}

	// remove the snapshot
	err = client.RemoveSnapshot(defaultCtx, snapshot.ID, snapshotName)
	if err != nil {
		panic(err)
	}

	// make sure the snapshot was removed
	snapshot, err = client.GetSnapshot(defaultCtx, snapshot.ID, snapshotName)
	if err != nil {
		panic(err)
	}
	if snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) was not removed.\n%+v\n", snapshotName, snapshot))
	}
}

func TestSnapshotCopy(_ *testing.T) {
	accessZone := "System"
	sourceSnapshotPath := "test_snapshot_copy_src_volume"
	sourceSnapshotName := "test_snapshot_copy_src_snapshot"
	destinationVolume := "test_snapshot_copy_dst_volume"
	subdirectoryName := "test_snapshot_copy_sub_dir"
	sourceSubDirectory := fmt.Sprintf("%s/%s", sourceSnapshotPath, subdirectoryName)
	destinationSubDirectory := fmt.Sprintf("%s/%s", destinationVolume, subdirectoryName)

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, sourceSnapshotPath)
	if err != nil {
		panic(err)
	}
	//	defer client.DeleteVolume(snapshotPath)
	// create a subdirectory in the test tvolume
	_, err = client.CreateVolume(defaultCtx, sourceSubDirectory)
	if err != nil {
		panic(err)
	}

	// make sure the snapshot doesn't exist yet
	snapshot, err := client.GetSnapshot(defaultCtx, -1, sourceSnapshotName)
	if err == nil && snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) already exists.\n", sourceSnapshotName))
	}

	// Add the test snapshot
	testSnapshot, err := client.CreateSnapshot(
		defaultCtx, sourceSnapshotPath, sourceSnapshotName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot.ID, sourceSnapshotName)
	// remove the sub directory
	err = client.DeleteVolume(defaultCtx, sourceSubDirectory)
	if err != nil {
		panic(err)
	}

	// copy the snapshot to the destination volume
	_, err = client.CreateVolume(defaultCtx, destinationVolume)
	if err != nil {
		panic(err)
	}
	copiedVolume, err := client.CopySnapshot(
		defaultCtx, testSnapshot.ID, testSnapshot.Name, accessZone, destinationVolume)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, destinationVolume)

	if copiedVolume.Name != destinationVolume {
		panic(fmt.Sprintf("Copied volume has incorrect name.  Expected: (%s) Acutal: (%s)", destinationVolume, copiedVolume.Name))
	}

	// make sure the destination volume was created
	volume, err := client.GetVolume(defaultCtx, "", destinationVolume)
	if err != nil || volume == nil {
		panic(fmt.Sprintf("Destination volume: (%s) was not created.\n", destinationVolume))
	}
	// make sure the sub directory was also created
	subDirectory, err := client.GetVolume(defaultCtx, "", destinationSubDirectory)
	if err != nil {
		panic(fmt.Sprintf("Destination sub directory: (%s) was not created.\n", subdirectoryName))
	}
	if subDirectory.Name != destinationSubDirectory {
		panic(fmt.Sprintf("Sub directory has incorrect name.  Expected: (%s) Acutal: (%s)", destinationSubDirectory, subDirectory.Name))
	}
}

func TestSnapshotCopyWithIsiPath(_ *testing.T) {
	sourceSnapshotPath := "test_snapshot_copy_src_volume"
	sourceSnapshotName := "test_snapshot_copy_src_snapshot"
	destinationVolume := "test_snapshot_copy_dst_volume"
	subdirectoryName := "test_snapshot_copy_sub_dir"
	defaultAccessZone := "System"
	sourceSubDirectory := fmt.Sprintf("%s/%s", sourceSnapshotPath, subdirectoryName)
	destinationSubDirectory := fmt.Sprintf("%s/%s", destinationVolume, subdirectoryName)

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, sourceSnapshotPath)
	if err != nil {
		panic(err)
	}
	//	defer client.DeleteVolume(snapshotPath)
	// create a subdirectory in the test tvolume
	_, err = client.CreateVolume(defaultCtx, sourceSubDirectory)
	if err != nil {
		panic(err)
	}

	// make sure the snapshot doesn't exist yet
	snapshot, err := client.GetSnapshot(defaultCtx, -1, sourceSnapshotName)
	if err == nil && snapshot != nil {
		panic(fmt.Sprintf("Snapshot (%s) already exists.\n", sourceSnapshotName))
	}

	// Add the test snapshot
	testSnapshot, err := client.CreateSnapshot(
		defaultCtx, sourceSnapshotPath, sourceSnapshotName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot.ID, sourceSnapshotName)
	// remove the sub directory
	err = client.DeleteVolume(defaultCtx, sourceSubDirectory)
	if err != nil {
		panic(err)
	}

	// copy the snapshot to the destination volume
	_, err = client.CreateVolume(defaultCtx, destinationVolume)
	if err != nil {
		panic(err)
	}
	newIsiPath := os.Getenv("GOISILON_VOLUMEPATH")
	copiedVolume, err := client.CopySnapshotWithIsiPath(
		defaultCtx, newIsiPath, newIsiPath, testSnapshot.ID, testSnapshot.Name, destinationVolume, defaultAccessZone)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, destinationVolume)

	if copiedVolume.Name != destinationVolume {
		panic(fmt.Sprintf("Copied volume has incorrect name.  Expected: (%s) Acutal: (%s)", destinationVolume, copiedVolume.Name))
	}

	// make sure the destination volume was created
	volume, err := client.GetVolume(defaultCtx, "", destinationVolume)
	if err != nil || volume == nil {
		panic(fmt.Sprintf("Destination volume: (%s) was not created.\n", destinationVolume))
	}
	// make sure the sub directory was also created
	subDirectory, err := client.GetVolume(defaultCtx, "", destinationSubDirectory)
	if err != nil {
		panic(fmt.Sprintf("Destination sub directory: (%s) was not created.\n", subdirectoryName))
	}
	if subDirectory.Name != destinationSubDirectory {
		panic(fmt.Sprintf("Sub directory has incorrect name.  Expected: (%s) Acutal: (%s)", destinationSubDirectory, subDirectory.Name))
	}
}

func TestSnapshotGetByIdentity(_ *testing.T) {
	snapshotPath := "test_snapshots_get_volume"
	snapshotName1 := "test_snapshots_get_snapshot_0"
	snapshotName2 := "test_snapshots_get_snapshot_1"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}
	testSnapshot2, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName2)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)
	defer client.RemoveSnapshot(defaultCtx, testSnapshot2.ID, snapshotName2)

	snapshot1, err := client.GetIsiSnapshotByIdentity(defaultCtx, snapshotName1)
	if err != nil {
		panic("failed to get testSnapshot1\n")
	}
	if snapshot1.ID != testSnapshot1.ID {
		panic(fmt.Sprintf("testSnapshot1: id %d is incorrect\n", snapshot1.ID))
	}
	snapshot2, err := client.GetIsiSnapshotByIdentity(defaultCtx, snapshotName2)
	if err != nil {
		panic("failed to get testSnapshot2\n")
	}
	if snapshot2.ID != testSnapshot2.ID {
		panic(fmt.Sprintf("testSnapshot2: id %d is incorrect\n", snapshot2.ID))
	}
}

func TestSnapshotIsExistent(_ *testing.T) {
	snapshotPath := "test_snapshots_exist_volume"
	snapshotName1 := "test_snapshots_exist_snapshot_0"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// check if snapshotName1 exists
	if client.IsSnapshotExistent(defaultCtx, snapshotName1) {
		panic(fmt.Sprintf("found snapshot %s, expected not found\n", snapshotName1))
	}

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}

	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)

	// check if snapshotName1 exists
	if !client.IsSnapshotExistent(defaultCtx, snapshotName1) {
		panic(fmt.Sprintf("not found snapshot %s, expected found\n", snapshotName1))
	}
}

func TestSnapshotExportWithZone(_ *testing.T) {
	snapshotPath := "test_snapshots_export_volume"
	snapshotName1 := "test_snapshots_export_snapshot_0"
	defaultAccessZone := "System"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}

	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)

	// export snapshot
	id, err := client.ExportSnapshotWithZone(defaultCtx, snapshotName1, snapshotPath, defaultAccessZone, "")
	if err != nil {
		panic(fmt.Sprintf("failed to export snapshot, name %s path %s \n", snapshotName1, snapshotPath))
	}

	// unexport snapshot
	defer client.UnexportByIDWithZone(defaultCtx, id, defaultAccessZone)
}

func TestGetRealVolumeSnapshotPathWithIsiPath(_ *testing.T) {
	volName := "volFromSnap0"
	newIsiPath := os.Getenv("GOISILON_VOLUMEPATH")
	accessZone := "System"
	name := "snapshottest"
	fmt.Println(apiv1.GetRealVolumeSnapshotPathWithIsiPath(newIsiPath, volName, name, accessZone))
}

func TestSnapshotSizeGet(_ *testing.T) {
	snapshotPath := "test_snapshots_get_volume"
	snapshotName1 := "test_snapshots_get_snapshot_0"
	accessZone := "System"

	// create the test volume
	_, err := client.CreateVolume(defaultCtx, snapshotPath)
	if err != nil {
		panic(err)
	}
	defer client.DeleteVolume(defaultCtx, snapshotPath)

	// Add the test snapshots
	testSnapshot1, err := client.CreateSnapshot(
		defaultCtx, snapshotPath, snapshotName1)
	if err != nil {
		panic(err)
	}

	// make sure we clean up when we're done
	defer client.RemoveSnapshot(defaultCtx, testSnapshot1.ID, snapshotName1)

	newIsiPath := os.Getenv("GOISILON_VOLUMEPATH")
	// get the updated snapshot list
	totalSize, err := client.GetSnapshotFolderSize(defaultCtx, newIsiPath, snapshotName1, accessZone)
	if err != nil {
		panic(err)
	}
	if totalSize < 0 {
		panic(fmt.Sprintf("Snapshot folder size %d is not correct.\n", totalSize))
	}
}

func TestGetWritableSnapshots(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv14.IsiWritableSnapshotQueryResponse)
		*resp = &apiv14.IsiWritableSnapshotQueryResponse{
			Writable: []*apiv14.IsiWritableSnapshot{
				{
					ID: 100,
				},
			},
		}
	}).Once()

	result, err := client.GetWritableSnapshots(context.Background())
	assert.Nil(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(100), result[0].ID)
}

func TestGetWritableSnapshot(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv14.IsiWritableSnapshotQueryResponse)
		*resp = &apiv14.IsiWritableSnapshotQueryResponse{
			Writable: []*apiv14.IsiWritableSnapshot{
				{
					ID: 100,
				},
			},
		}
	}).Once()

	result, err := client.GetWritableSnapshot(context.Background(), "/ifs/data1")
	assert.Nil(t, err)
	assert.Equal(t, int64(100), result.ID)

	result, err = client.GetWritableSnapshot(context.Background(), "/data1")
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "invalid snapshot path, must start with /ifs/: /data1")
}

func TestCreateWritableSnapshot(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshotPath := "/ifs/data1"
	sourceSnapshot := "snapshot_source_1"
	destinationPath := "/ifs/data2"

	client.API.(*mocks.Client).On("Post", anyArgs[0:7]...).Return(nil).Run(func(args mock.Arguments) {
		body := args.Get(5).(map[string]string)
		assert.Equal(t, sourceSnapshot, body["src_snap"])
		assert.Equal(t, destinationPath, body["dst_path"])

		resp := args.Get(6).(**apiv14.IsiWritableSnapshot)
		*resp = &apiv14.IsiWritableSnapshot{
			ID:           100,
			SrcPath:      snapshotPath,
			DstPath:      destinationPath,
			SrcSnap:      sourceSnapshot,
			State:        apiv14.WritableSnapshotStateActive,
			LogSize:      100,
			PhysicalSize: 200,
		}
	}).Once()

	result, err := client.CreateWritableSnapshot(context.Background(), sourceSnapshot, destinationPath)
	assert.Nil(t, err)
	assert.Equal(t, int64(100), result.ID)
	assert.Equal(t, snapshotPath, result.SrcPath)
	assert.Equal(t, destinationPath, result.DstPath)
	assert.Equal(t, sourceSnapshot, result.SrcSnap)
	assert.Equal(t, apiv14.WritableSnapshotStateActive, result.State)
	assert.Equal(t, int64(100), result.LogSize)
	assert.Equal(t, int64(200), result.PhysicalSize)

	// Test case: error in API call
	client.API.(*mocks.Client).On("Post", anyArgs[0:7]...).Return(errors.New("API call failed")).Once()

	result, err = client.CreateWritableSnapshot(context.Background(), sourceSnapshot, destinationPath)
	assert.ErrorContains(t, err, "API call failed")
	assert.Nil(t, result)
}

func TestRemoveWritableSnapshot(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		tgt := args.Get(1).(string)
		assert.Equal(t, "platform/14/snapshot/writable/ifs/data1", tgt)
	}).Once()

	err := client.RemoveWritableSnapshot(context.Background(), "/ifs/data1")
	assert.Nil(t, err)

	// Test case: error in API call
	err = client.RemoveWritableSnapshot(context.Background(), "/data1")
	assert.ErrorContains(t, err, "invalid snapshot path, must start with /ifs/: /data1")
}

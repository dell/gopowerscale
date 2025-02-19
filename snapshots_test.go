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
	"context"
	"errors"
	"fmt"
	"testing"

	apiv1 "github.com/dell/goisilon/api/v1"
	apiv14 "github.com/dell/goisilon/api/v14"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSnapshots(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{}
	}).Once()
	_, err = client.GetSnapshots(context.Background())
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err = client.GetSnapshots(context.Background())
	assert.NotNil(t, err)
}

func TestCreateSnapshot(t *testing.T) {
	volName := "testVol"
	snapshotName := "testSnapshot"
	volumePath := "/path/to/volume"

	// Successful snapshot creation
	client.API.(*mocks.Client).On("VolumePath", volName).Return(volumePath).Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	// Call the function
	_, err := client.CreateSnapshot(context.Background(), volName, snapshotName)
	assert.Nil(t, err)
}

func TestGetSnapshotsByPath(t *testing.T) {
	path := "testPath"
	volumePath := "/path/to/volume"

	// Mock snapshots
	mockSnapshots := SnapshotList{
		&apiv1.IsiSnapshot{Path: volumePath},
		&apiv1.IsiSnapshot{Path: "/another/path"},
	}

	// Successful retrieval of snapshots
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: mockSnapshots,
		}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", path).Return(volumePath).Twice()

	// Call the function
	result, err := client.GetSnapshotsByPath(context.Background(), path)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, volumePath, result[0].Path)

	// Retrieval failure
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("retrieval failed")).Once()

	// Call the function
	result, err = client.GetSnapshotsByPath(context.Background(), path)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(result))
}

func TestGetIsiSnapshotByIdentity(t *testing.T) {
	// Test case 1: Successful retrieval of snapshot by identity
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshot := &apiv1.IsiSnapshot{
		ID:    1,
		Name:  "test_snapshot",
		Path:  "/path/to/snapshot",
		State: "available",
	}

	// Mock the Get method to simulate a successful response
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{snapshot},
		}
	}).Once()

	// Call the GetIsiSnapshotByIdentity function
	result, err := client.GetIsiSnapshotByIdentity(context.Background(), "test_identity")

	// Assertions
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, snapshot.ID, result.ID)
	assert.Equal(t, snapshot.Name, result.Name)
	assert.Equal(t, snapshot.Path, result.Path)
	assert.Equal(t, snapshot.State, result.State)

	// Test case 2: Error when snapshot is not found
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("snapshot not found")).Once()

	// Call the GetIsiSnapshotByIdentity function
	result, err = client.GetIsiSnapshotByIdentity(context.Background(), "test_identity")

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)

	// Test case 3: Error handling in Get method (e.g., network error)
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("network error")).Once()

	// Call the GetIsiSnapshotByIdentity function
	result, err = client.GetIsiSnapshotByIdentity(context.Background(), "test_identity")

	// Assertions
	assert.NotNil(t, err)
	assert.Nil(t, result)

	// Assert expectations on mocks
	client.API.(*mocks.Client).AssertExpectations(t)
}

func TestIsSnapshotExistent(t *testing.T) {
	// Test case 1: Snapshot exists
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshot := &apiv1.IsiSnapshot{
		ID:    1,
		Name:  "test_snapshot",
		Path:  "/path/to/snapshot",
		State: "available",
	}

	// Mock the GetIsiSnapshotByIdentity to simulate a successful response
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{snapshot},
		}
	}).Once()

	// Call the IsSnapshotExistent function
	result := client.IsSnapshotExistent(context.Background(), "test_identity")

	// Assertions
	assert.True(t, result) // Snapshot exists

	// Test case 2: Snapshot does not exist
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("snapshot not found")).Once()

	// Call the IsSnapshotExistent function
	result = client.IsSnapshotExistent(context.Background(), "test_identity")

	// Assertions
	assert.False(t, result) // Snapshot does not exist

	// Test case 3: Error when fetching snapshot (e.g., network error)
	client.API.(*mocks.Client).On("Get", mock.Anything, "platform/1/snapshot/snapshots/test_identity", "", mock.Anything, mock.Anything, mock.Anything).
		Return(fmt.Errorf("network error")).Once()

	// Call the IsSnapshotExistent function
	result = client.IsSnapshotExistent(context.Background(), "test_identity")

	// Assertions
	assert.False(t, result) // Snapshot fetch failed, so it does not exist

	// Assert expectations on mocks
	client.API.(*mocks.Client).AssertExpectations(t)
}

func TestRemoveSnapshot(t *testing.T) {
	// Clear previous expectations
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Define snapshot ID and name for the test.
	snapshotID := int64(123)
	snapshotName := "test-snapshot"

	// Mock GetSnapshot to return a snapshot from GetIsiSnapshotsResp
	client.API.(*mocks.Client).On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()

	// Mock Delete to succeed (no error)
	client.API.(*mocks.Client).On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Once()

	// Call RemoveSnapshot
	err := client.RemoveSnapshot(context.Background(), snapshotID, snapshotName)

	// Assert that no error occurred
	assert.Nil(t, err)

	// Assert that the mock expectations were met
	client.API.(*mocks.Client).AssertExpectations(t)

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Twice()
	// client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	err = client.RemoveSnapshot(context.Background(), snapshotID, snapshotName)
	assert.NotNil(t, err)
}

func TestCreateSnapshotWithPath(t *testing.T) {
	path := "/path/to/snapshot"
	var snapshotName string

	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err := client.CreateSnapshotWithPath(context.Background(), path, snapshotName)
	assert.Nil(t, err)
}

func TestGetSnapshotFolderSize(t *testing.T) {
	client := &Client{
		API: new(mocks.Client),
	}

	ctx := context.Background()
	var isiPath, accessZone, name string

	// Mock GetSnapshot to return a snapshot from GetIsiSnapshotsResp
	client.API.(*mocks.Client).On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   1,
					Name: "test_snapshot",
					Path: "/path/to/snapshot",
				},
			},
			Total:  1,
			Resume: "",
		}
	}).Once()

	// Mock GetZoneByName to return a zone
	client.API.(*mocks.Client).On("Get", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: []*apiv1.IsiZone{
				{
					Name: "zone1",
					Path: "/ifs/data",
				},
			},
		}
	}).Once()

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiVolumeSizeResp)
		*resp = &apiv1.GetIsiVolumeSizeResp{
			AttributeMap: []struct {
				Name string `json:"name"`
				Size int64  `json:"size"`
			}{
				{Name: "test", Size: 12345},
			},
		}
	}).Once()

	// Call GetSnapshotFolderSize
	_, err := client.GetSnapshotFolderSize(ctx, isiPath, name, accessZone)

	// Assert that no error occurred
	assert.Nil(t, err)
}

func TestCopySnapshot(t *testing.T) {
	// Clear previous expectations
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshotID := int64(123)
	snapshotName := "test-snapshot"

	// Mock GetSnapshot to return a snapshot from GetIsiSnapshotsResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()

	// Mock GetZoneByName to return a zones from GetIsiZonesResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiZonesResp
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: []*apiv1.IsiZone{
				{
					Name: snapshotName,
					Path: "/ifs/data",
				},
			},
		}
	}).Once()

	// Mock CopyIsiSnapshot to return a volumes from IsiVolume
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Run(nil).Once()

	// Mock GetVolume to return a volumes from GetIsiVolumeAttributesResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiVolumeAttributesResp
		resp := args.Get(5).(**apiv1.GetIsiVolumeAttributesResp)
		*resp = &apiv1.GetIsiVolumeAttributesResp{}
	}).Once()

	_, err := client.CopySnapshot(context.Background(), snapshotID, snapshotName, "test-zone", "test-snapshot-copy")
	assert.Nil(t, err)

	// Negative Scenarios
	// Case 1 - Unable to get snapshot
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("unable to retrieve snapshot")).Run(nil).Once()
	_, err = client.CopySnapshot(context.Background(), snapshotID, "", "test-zone", "test-snapshot-copy")
	assert.Error(t, err)

	// Case 2 - Unable to get zone
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("unable to retrieve zone")).Run(nil).Once()
	_, err = client.CopySnapshot(context.Background(), snapshotID, snapshotName, "", "test-snapshot-copy")
	assert.Error(t, err)
}

func TestGetSnapshotIsiPath(t *testing.T) {
	// Clear previous expectations
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshotID := int64(123)
	snapshotName := "test-snapshot"
	isiPath := "/ifs/data"
	accessZone := "test-zone"

	// Mock GetIsiSnapshotByIdentity to return a snapshot from GetIsiSnapshotsResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()

	// Mock GetZoneByName to return a zones from GetIsiZonesResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiZonesResp
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: []*apiv1.IsiZone{
				{
					Name: snapshotName,
					Path: "/ifs/data",
				},
			},
		}
	}).Once()

	_, err := client.GetSnapshotIsiPath(context.Background(), isiPath, "test-snapshot", accessZone)
	assert.Nil(t, err)

	// Negative Scenarios
	// Case 1 - Unable to get Isi Snapshot By Identity
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("unable to retrieve snapshot")).Run(nil).Once()
	_, err = client.GetSnapshotIsiPath(context.Background(), isiPath, "", accessZone)
	assert.Error(t, err)

	// Case 2 - Unable to get zone
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("unable to retrieve zone")).Run(nil).Once()
	_, err = client.GetSnapshotIsiPath(context.Background(), isiPath, "test-snapshot", "")
	assert.Error(t, err)

	// Case 3 - when zone.Path and isiPath mismatch
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiSnapshotsResp
		resp := args.Get(5).(**apiv1.GetIsiSnapshotsResp)
		*resp = &apiv1.GetIsiSnapshotsResp{
			SnapshotList: []*apiv1.IsiSnapshot{
				{
					ID:   snapshotID,
					Name: snapshotName,
				},
			},
			Total: 1,
		}
	}).Once()

	// Mock GetZoneByName to return a zones from GetIsiZonesResp
	client.API.(*mocks.Client).On("VolumesPath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for GetIsiZonesResp
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: []*apiv1.IsiZone{
				{
					Name: snapshotName,
					Path: "/ifs/data1",
				},
			},
		}
	}).Once()

	_, err = client.GetSnapshotIsiPath(context.Background(), isiPath, "test-snapshot", accessZone)
	assert.Nil(t, err)
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

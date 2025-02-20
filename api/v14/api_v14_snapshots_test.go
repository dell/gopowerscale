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

package v14

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetWritableSnapshots_ErrorCase1(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("unauthorized")).Once()
	_, err := GetIsiWritableSnapshots(ctx, client)
	assert.NotNil(t, err)
}

func TestGetWritableSnapshots_ErrorCase2(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotQueryResponse)
		*resp = &IsiWritableSnapshotQueryResponse{
			Writable: []*IsiWritableSnapshot{
				{
					ID: 100,
				},
			},
			Total:  2048,
			Resume: "resume",
		}
	}).Once()

	client.On("Get", anyArgs...).Return(errors.New("invalid")).Once()
	_, err := GetIsiWritableSnapshots(ctx, client)
	assert.NotNil(t, err)
}

func TestGetWritableSnapshots(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotQueryResponse)
		*resp = &IsiWritableSnapshotQueryResponse{
			Writable: []*IsiWritableSnapshot{
				{
					ID: 100,
				},
			},
			Total:  2048,
			Resume: "resume",
		}
	}).Once()

	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotQueryResponse)
		*resp = &IsiWritableSnapshotQueryResponse{
			Writable: []*IsiWritableSnapshot{
				{
					ID: 1000000,
				},
			},
			Total:  1,
			Resume: "",
		}
	}).Once()

	result, err := GetIsiWritableSnapshots(ctx, client)
	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, int64(100), result[0].ID)
}

func TestGetWritableSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotQueryResponse)
		*resp = &IsiWritableSnapshotQueryResponse{
			Writable: []*IsiWritableSnapshot{
				{
					ID: 100,
				},
			},
		}
	}).Once()

	result, err := GetIsiWritableSnapshot(ctx, client, "/ifs/data1")
	assert.Nil(t, err)
	assert.Equal(t, int64(100), result.ID)

	client.On("Get", anyArgs...).Return(errors.New("not found")).Once()
	result, err = GetIsiWritableSnapshot(ctx, client, "/ifs/data1")
	assert.NotNil(t, err)
}

func TestCreateWritableSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	snapshotPath := "/ifs/data1"
	sourceSnapshot := "snapshot_source_1"
	destinationPath := "/ifs/data2"

	client.On("Post", anyArgs[0:7]...).Return(nil).Run(func(args mock.Arguments) {
		body := args.Get(5).(map[string]string)
		assert.Equal(t, sourceSnapshot, body["src_snap"])
		assert.Equal(t, destinationPath, body["dst_path"])

		resp := args.Get(6).(**IsiWritableSnapshot)
		*resp = &IsiWritableSnapshot{
			ID:           100,
			SrcPath:      snapshotPath,
			DstPath:      destinationPath,
			SrcSnap:      sourceSnapshot,
			State:        WritableSnapshotStateActive,
			LogSize:      100,
			PhysicalSize: 200,
		}
	}).Once()

	result, err := CreateWritableSnapshot(ctx, client, sourceSnapshot, destinationPath)
	assert.Nil(t, err)
	assert.Equal(t, int64(100), result.ID)
	assert.Equal(t, snapshotPath, result.SrcPath)
	assert.Equal(t, destinationPath, result.DstPath)
	assert.Equal(t, sourceSnapshot, result.SrcSnap)
	assert.Equal(t, WritableSnapshotStateActive, result.State)
	assert.Equal(t, int64(100), result.LogSize)
	assert.Equal(t, int64(200), result.PhysicalSize)

	// Test case: error in API call
	client.On("Post", anyArgs[0:7]...).Return(errors.New("API call failed")).Once()

	result, err = CreateWritableSnapshot(ctx, client, sourceSnapshot, destinationPath)
	assert.ErrorContains(t, err, "API call failed")
	assert.Nil(t, result)
}

func TestRemoveWritableSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		tgt := args.Get(1).(string)
		assert.Equal(t, "platform/14/snapshot/writable/ifs/data1", tgt)
	}).Once()

	err := RemoveWritableSnapshot(ctx, client, "/ifs/data1")
	assert.Nil(t, err)
}

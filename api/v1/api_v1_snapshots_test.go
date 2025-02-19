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

package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetIsiSnapshots(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Run(nil).Once()
	_, err := GetIsiSnapshots(ctx, client)
	assert.Equal(t, errors.New("error"), err)

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**GetIsiSnapshotsResp)
		*resp = &GetIsiSnapshotsResp{
			SnapshotList: []*IsiSnapshot{
				{
					ID:   1,
					Name: "test_snapshot",
					Path: "/path/to/snapshot/",
				},
			},
			Total:  1,
			Resume: "resume",
		}
	}).Once()
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**GetIsiSnapshotsResp)
		*resp = &GetIsiSnapshotsResp{
			SnapshotList: []*IsiSnapshot{
				{
					ID:   1,
					Name: "test_snapshot",
					Path: "/path/to/snapshot/",
				},
			},
			Total:  1,
			Resume: "",
		}
	}).Once()

	_, err = GetIsiSnapshots(ctx, client)
	assert.Equal(t, nil, err)
}

func TestGetIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	var x int64
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiSnapshot(ctx, client, x)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiSnapshotByIdentity(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiSnapshotByIdentity(ctx, client, "")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestCreateIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	_, err := CreateIsiSnapshot(ctx, client, "", "name")
	assert.Equal(t, errors.New("no path set"), err)

	client.On("Post", anyArgs...).Return(nil).Twice()
	_, err = CreateIsiSnapshot(ctx, client, "path", "name")
	assert.Equal(t, nil, err)
}

func TestRemoveIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	var id int64
	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := RemoveIsiSnapshot(ctx, client, id)
	assert.Equal(t, nil, err)
}

func TestCopyIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Put", anyArgs...).Return(nil).Once()
	_, err := CopyIsiSnapshot(ctx, client, "", "", "", "", "")
	assert.Equal(t, nil, err)
}

func TestCopyIsiSnapshotWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*GetIsiZonesResp)
		*resp = GetIsiZonesResp{
			Zones: []*IsiZone{
				{
					ID: "test-id",
				},
			},
		}
	}).Once()
	client.On("Put", anyArgs...).Return(nil).Once()
	_, err := CopyIsiSnapshotWithIsiPath(ctx, client, "", "", "", "", "", "")
	assert.Equal(t, nil, err)
}

func TestGetIsiSnapshotFolderWithSize(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*GetIsiZonesResp)
		*resp = GetIsiZonesResp{
			Zones: []*IsiZone{
				{
					ID: "test-id",
				},
			},
		}
	}).Once()
	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetIsiSnapshotFolderWithSize(ctx, client, "", "", "", "")
	assert.Equal(t, nil, err)
}

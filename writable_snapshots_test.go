/*
Copyright (c) 2026 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package gopowerscale

import (
	"context"
	"errors"
	"testing"

	apiv14 "github.com/dell/gopowerscale/api/v14"
	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateWritableSnapshot_Success(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	snapshotID := int64(12345)
	destinationPath := "/ifs/data/writable-vol-1"

	// Mock successful creation
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv14.IsiWritableSnapshotResponse)
		*resp = &apiv14.IsiWritableSnapshotResponse{
			ID:       24315648,
			DstPath:  destinationPath,
			SrcID:    snapshotID,
			SrcPath:  "/ifs/source",
			SrcSnap:  "Snapshot: test",
			Created:  1783420362,
			State:    "active",
			LogSize:  0,
			PhysSize: 2048,
		}
	}).Once()

	ws, err := client.CreateWritableSnapshot(context.Background(), destinationPath, "12345")
	assert.NoError(t, err)
	assert.NotNil(t, ws)
	assert.Equal(t, int64(24315648), ws.ID)
	assert.Equal(t, destinationPath, ws.DstPath)
	assert.Equal(t, "active", ws.State)
}

func TestCreateWritableSnapshot_Errors(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test license error
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(errors.New("license not enabled")).Once()
	_, err := client.CreateWritableSnapshot(context.Background(), "/ifs/data/writable-vol-1", "12345")
	assert.Error(t, err)

	// Test snapshot not found
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(errors.New("snapshot not found")).Once()
	_, err = client.CreateWritableSnapshot(context.Background(), "/ifs/data/writable-vol-1", "12345")
	assert.Error(t, err)

	// Test internal error
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(errors.New("internal server error")).Once()
	_, err = client.CreateWritableSnapshot(context.Background(), "/ifs/data/writable-vol-1", "12345")
	assert.Error(t, err)
}

func TestGetWritableSnapshot_Success(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	path := "/ifs/writeable"

	// Mock successful GET with wrapper response
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv14.IsiWritableSnapshotListResponse)
		*resp = &apiv14.IsiWritableSnapshotListResponse{
			Writable: []apiv14.IsiWritableSnapshotResponse{
				{
					ID:       24315648,
					DstPath:  path,
					SrcID:    3754,
					SrcPath:  "/ifs/sky",
					SrcSnap:  "Snapshot: 2026Jul22, 5:05 PM",
					Created:  1783420362,
					State:    "active",
					LogSize:  0,
					PhysSize: 2048,
				},
			},
		}
	}).Once()

	ws, err := client.GetWritableSnapshot(context.Background(), path)
	assert.NoError(t, err)
	assert.NotNil(t, ws)
	assert.Equal(t, int64(24315648), ws.ID)
	assert.Equal(t, path, ws.DstPath)
}

func TestGetWritableSnapshot_Error(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test not found
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("not found")).Once()
	_, err := client.GetWritableSnapshot(context.Background(), "/ifs/writeable")
	assert.Error(t, err)
}

func TestDeleteWritableSnapshot_Success(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	path := "/ifs/writeable"

	// Mock successful DELETE
	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(nil).Once()

	err := client.DeleteWritableSnapshot(context.Background(), path)
	assert.NoError(t, err)
}

func TestDeleteWritableSnapshot_Errors(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	// Test dependency error
	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(errors.New("dependency error")).Once()
	err := client.DeleteWritableSnapshot(context.Background(), "/ifs/writeable")
	assert.Error(t, err)

	// Test not found
	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(errors.New("not found")).Once()
	err = client.DeleteWritableSnapshot(context.Background(), "/ifs/writeable")
	assert.Error(t, err)
}

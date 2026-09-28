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

package v14

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateIsiWritableSnapshot_Success(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	request := &IsiWritableSnapshotRequest{
		SourceSnapshot:  "12345",
		DestinationPath: "/ifs/data/writable-vol-1",
	}

	// Mock successful POST response
	client.On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**IsiWritableSnapshotResponse)
		*resp = &IsiWritableSnapshotResponse{
			ID:       24315648,
			DstPath:  "/ifs/data/writable-vol-1",
			SrcID:    12345,
			SrcPath:  "/ifs/source",
			SrcSnap:  "Snapshot: test",
			Created:  1783420362,
			State:    "active",
			LogSize:  0,
			PhysSize: 2048,
		}
	}).Once()

	resp, err := CreateIsiWritableSnapshot(ctx, client, request)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(24315648), resp.ID)
	assert.Equal(t, "/ifs/data/writable-vol-1", resp.DstPath)
	assert.Equal(t, "active", resp.State)
}

func TestCreateIsiWritableSnapshot_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test nil request
	_, err := CreateIsiWritableSnapshot(ctx, client, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "request cannot be nil")

	// Test empty DestinationPath
	_, err = CreateIsiWritableSnapshot(ctx, client, &IsiWritableSnapshotRequest{
		SourceSnapshot:  "12345",
		DestinationPath: "",
	})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "destination_path cannot be empty")
}

func TestCreateIsiWritableSnapshot_HTTPErrors(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	request := &IsiWritableSnapshotRequest{
		SourceSnapshot:  "12345",
		DestinationPath: "/ifs/data/writable-vol-1",
	}

	// Test HTTP 404 - snapshot not found
	client.On("Post", anyArgs...).Return(errors.New("snapshot not found")).Once()
	_, err := CreateIsiWritableSnapshot(ctx, client, request)
	assert.Error(t, err)

	// Test HTTP 403 - license error
	client.On("Post", anyArgs...).Return(errors.New("license error")).Once()
	_, err = CreateIsiWritableSnapshot(ctx, client, request)
	assert.Error(t, err)

	// Test HTTP 500 - internal error
	client.On("Post", anyArgs...).Return(errors.New("internal server error")).Once()
	_, err = CreateIsiWritableSnapshot(ctx, client, request)
	assert.Error(t, err)
}

func TestGetIsiWritableSnapshot_Success(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Mock successful GET response with wrapper structure
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotListResponse)
		*resp = &IsiWritableSnapshotListResponse{
			Writable: []IsiWritableSnapshotResponse{
				{
					ID:       24315648,
					DstPath:  "/ifs/writeable",
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

	resp, err := GetIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(24315648), resp.ID)
	assert.Equal(t, "/ifs/writeable", resp.DstPath)
	assert.Equal(t, "active", resp.State)
}

func TestGetIsiWritableSnapshot_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test empty path
	_, err := GetIsiWritableSnapshot(ctx, client, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path cannot be empty")
}

func TestGetIsiWritableSnapshot_EmptyArray(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Mock response with empty writable array
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotListResponse)
		*resp = &IsiWritableSnapshotListResponse{
			Writable: []IsiWritableSnapshotResponse{},
		}
	}).Once()

	_, err := GetIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "writable snapshot not found")
}

func TestGetIsiWritableSnapshot_HTTPError(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test HTTP 404
	client.On("Get", anyArgs...).Return(errors.New("not found")).Once()
	_, err := GetIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.Error(t, err)
}

func TestGetIsiWritableSnapshot_NilResponse(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Mock response with nil list response
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotListResponse)
		*resp = nil
	}).Once()

	_, err := GetIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected nil response")
}

func TestDeleteIsiWritableSnapshot_Success(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Mock successful DELETE
	client.On("Delete", anyArgs...).Return(nil).Once()

	err := DeleteIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.NoError(t, err)
}

func TestDeleteIsiWritableSnapshot_ValidationErrors(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test empty path
	err := DeleteIsiWritableSnapshot(ctx, client, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path cannot be empty")
}

func TestDeleteIsiWritableSnapshot_HTTPErrors(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test HTTP 409 - dependency error
	client.On("Delete", anyArgs...).Return(errors.New("dependency error")).Once()
	err := DeleteIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.Error(t, err)

	// Test HTTP 404 - not found
	client.On("Delete", anyArgs...).Return(errors.New("not found")).Once()
	err = DeleteIsiWritableSnapshot(ctx, client, "/ifs/writeable")
	assert.Error(t, err)
}

func TestGetIsiWritableSnapshot_PathTrimming(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	expectedPath := "platform/14/snapshot/writable/ifs/data/csi/vol1"
	client.On("Get", ctx, expectedPath, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiWritableSnapshotListResponse)
		*resp = &IsiWritableSnapshotListResponse{
			Writable: []IsiWritableSnapshotResponse{
				{ID: 1, DstPath: "/ifs/data/csi/vol1", State: "active"},
			},
		}
	}).Once()

	_, err := GetIsiWritableSnapshot(ctx, client, "/ifs/data/csi/vol1/")
	assert.NoError(t, err)
}

func TestDeleteIsiWritableSnapshot_PathTrimming(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	expectedPath := "platform/14/snapshot/writable/ifs/data/csi/vol1"
	client.On("Delete", ctx, expectedPath, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()

	err := DeleteIsiWritableSnapshot(ctx, client, "/ifs/data/csi/vol1/")
	assert.NoError(t, err)
}

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
package v12

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestListSmbShares(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	zone := ""
	params := ListV12SmbSharesParams{
		Zone: &zone,
	}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ListSmbShares(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetSmbShare(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	zone := ""
	params := GetV12SmbShareParams{
		Zone: &zone,
	}
	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetSmbShare(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestCreateSmbShare(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	zone := ""
	params := CreateV12SmbShareRequest{
		Zone: &zone,
	}
	client.On("Post", anyArgs...).Return(nil).Once()
	_, err := CreateSmbShare(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestUpdateSmbShare(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	zone := ""
	params := UpdateV12SmbShareRequest{
		Zone: &zone,
	}
	client.On("Put", anyArgs...).Return(nil).Once()
	err := UpdateSmbShare(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestDeleteSmbShare(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	zone := ""
	params := DeleteV12SmbShareRequest{
		Zone: &zone,
	}
	client.On("Delete", anyArgs...).Return(nil).Once()
	err := DeleteSmbShare(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

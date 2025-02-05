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
package v4

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestListNfsExports(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	check := false
	params := ListV4NfsExportsParams{
		Check: &check,
	}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ListNfsExports(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetNfsExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	params := GetV2NfsExportRequest{
		V2NFSExportID: "",
	}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetNfsExport(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestCreateNfsExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	force := false
	params := CreateV4NfsExportRequest{
		Force: &force,
	}
	client.On("Post", anyArgs...).Return(nil).Once()
	_, err := CreateNfsExport(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestUpdateNfsExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	force := false
	params := UpdateV4NfsExportRequest{
		Force: &force,
	}
	client.On("Put", anyArgs...).Return(nil).Once()
	err := UpdateNfsExport(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestDeleteNfsExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	params := DeleteV4NfsExportRequest{
		V2NFSExportID: "",
	}
	client.On("Delete", anyArgs...).Return(nil).Once()
	err := DeleteNfsExport(ctx, params, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

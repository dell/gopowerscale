/*
Copyright (c) 2023-2025 Dell Inc, or its subsidiaries.

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

	v12 "github.com/dell/goisilon/api/v12"
	"github.com/dell/goisilon/mocks"
	"github.com/dell/goisilon/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListAllSmbSharesWithStructParams(t *testing.T) {
	limit := int32(1)
	firstPageResume := "resume_token"

	// Test case 1: Mock the call to return a "not found" error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err := client.ListALlSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assert.NotNil(t, err)

	// Test case 2: Mock the first call to return a non-nil Resume value
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{
			Shares: []openapi.V12SmbShareExtended{
				{Name: "share1"},
			},
			Resume: &firstPageResume,
		}
	}).Once()

	// Test case 3: Mock the second call to return an error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("mock error")).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{}
	}).Once()

	_, err = client.ListALlSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assert.NotNil(t, err)

	// Test case 4: Mock the third call to return a nil Resume value
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{
			Shares: []openapi.V12SmbShareExtended{
				{Name: "share2"},
			},
			Resume: nil,
		}
	}).Once()

	_, err = client.ListALlSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assert.Nil(t, err)

	// Test case 5: Mock the first call to return a non-nil Resume value and the second call to return a nil Resume value
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{
			Shares: []openapi.V12SmbShareExtended{
				{Name: "share1"},
			},
			Resume: &firstPageResume,
		}
	}).Once()

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{
			Shares: []openapi.V12SmbShareExtended{
				{Name: "share2"},
			},
			Resume: nil,
		}
	}).Once()

	_, err = client.ListALlSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assert.Nil(t, err)
}

func TestListSmbSharesWithStructParams(t *testing.T) {
	// use limit to test pagination, would still output all shares
	limit := int32(1)
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbShares)
		*resp = openapi.V12SmbShares{}
	}).Once()
	_, err := client.ListSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assert.Nil(t, err)
}

func TestGetSmbShareWithStructParams(t *testing.T) {
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V12SmbSharesExtended)
		*resp = openapi.V12SmbSharesExtended{}
	}).Once()
	_, err := client.GetSmbShareWithStructParams(defaultCtx, v12.GetV12SmbShareParams{})
	assert.Nil(t, err)
}

func TestCreateSmbShareWithStructParams(t *testing.T) {
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err := client.CreateSmbShareWithStructParams(defaultCtx, v12.CreateV12SmbShareRequest{})
	assert.Nil(t, err)
}

func TestDeleteSmbShareWithStructParams(t *testing.T) {
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err := client.DeleteSmbShareWithStructParams(defaultCtx, v12.DeleteV12SmbShareRequest{})
	assert.Nil(t, err)
}

func TestUpdateSmbShareWithStructParams(t *testing.T) {
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.UpdateSmbShareWithStructParams(defaultCtx, v12.UpdateV12SmbShareRequest{})
	assert.Nil(t, err)
}

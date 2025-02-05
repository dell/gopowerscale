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

func TestGetIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	userName := "userName"
	var uid int32
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiUser(ctx, client, &userName, &uid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiUserListResp)
		*resp = &IsiUserListResp{
			Users: []*IsiUser{
				{
					ID:   "test-id",
					Name: userName,
				},
			},
		}
	}).Once()
	resp, err := GetIsiUser(ctx, client, &userName, &uid)
	resp.ID = "test-id"
	resp.Name = userName
	assert.Equal(t, nil, err)
}

func TestGetIsiUserList(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryMemberOf := false
	queryDomain := "queryDomain"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiUserList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone,
		&queryProvider, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiUserListRespResume)
		*resp = &IsiUserListRespResume{
			Users: []*IsiUser{
				{
					ID:   "test-id",
					Name: "test-name",
				},
			},
		}
	}).Once()
	_, err = GetIsiUserList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone,
		&queryProvider, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)

	assert.Equal(t, nil, err)
}

func TestCreateIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32

	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := CreateIsiUser(ctx, client, queryNamePrefix, &queryCached, &queryZone,
		&queryProvider, &queryZone, &queryZone, &queryZone, &queryZone, &queryZone,
		&queryZone, &queryLimit, &queryLimit, &queryLimit, &queryResolveNames,
		&queryResolveNames, &queryResolveNames, &queryResolveNames)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestUpdateIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32

	client.On("Put", anyArgs...).Return(errors.New("error found")).Twice()
	err := UpdateIsiUser(ctx, client, &queryNamePrefix, &queryLimit, &queryCached,
		&queryProvider, &queryZone, &queryZone, &queryZone, &queryZone, &queryZone,
		&queryZone, &queryZone, &queryLimit, &queryLimit, &queryLimit,
		&queryResolveNames, &queryResolveNames, &queryResolveNames, &queryResolveNames)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestDeleteIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryNamePrefix := "queryNamePrefix"

	var queryLimit int32

	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := DeleteIsiUser(ctx, client, &queryNamePrefix, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

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

func TestGetIsiGroupList(t *testing.T) {
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
	_, err := GetIsiGroupList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone, &queryProvider,
		&queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiGroupListRespResume)
		*resp = &IsiGroupListRespResume{
			Groups: []*IsiGroup{
				{
					ID:   "test-id",
					Name: "test-name",
				},
			},
			Resume: "resume",
		}
	}).Once()
	client.On("Get", anyArgs...).Return(errors.New("error found")).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiGroupListRespResume)
		*resp = &IsiGroupListRespResume{
			Groups: []*IsiGroup{
				{
					ID:   "test-id",
					Name: "test-name",
				},
			},
			Resume: "",
		}
	}).Once()
	_, err = GetIsiGroupList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone, &queryProvider,
		&queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	assert.Error(t, err)
}

func TestGetIsiGroupMembers(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiGroupMembers(ctx, client, &groupName, &gid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestAddIsiGroupMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}
	err := AddIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	assert.Equal(t, errors.New("member type is wrong, only support user and group"), err)

	authMember.Type = fileGroupTypeUser
	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	err = AddIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestRemoveIsiGroupMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "user",
	}

	authMember.Type = fileGroupTypeUser
	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := RemoveIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestCreateIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	queryForce := false
	queryZone := "queryZone"
	var gid int32
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "user",
	}
	authMemberItem := []IsiAuthMemberItem{authMember}

	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := CreateIsiGroup(ctx, client, name, &gid, authMemberItem, &queryForce, &queryZone, &queryZone)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestUpdateIsiGroupGID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	queryZone := "queryZone"
	var gid int32
	var newgid int32

	client.On("Put", anyArgs...).Return(errors.New("error found")).Twice()
	err := UpdateIsiGroupGID(ctx, client, &groupName, &gid, newgid, &queryZone, &queryZone)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestDeleteIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32

	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := DeleteIsiGroup(ctx, client, &groupName, &gid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32

	client.On("Get", anyArgs...).Return(errors.New("error found")).Once()
	_, err := GetIsiGroup(ctx, client, &groupName, &gid)
	assert.Equal(t, errors.New("error found"), err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiGroupListResp)
		*resp = &IsiGroupListResp{
			Groups: []*IsiGroup{
				{
					ID:   "0",
					Name: "name",
				},
			},
		}
	}).Once()
	_, err = GetIsiGroup(ctx, client, &groupName, &gid)
	assert.Equal(t, nil, err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiGroupListResp)
		*resp = &IsiGroupListResp{
			Groups: []*IsiGroup{},
		}
	}).Once()
	_, err = GetIsiGroup(ctx, client, &groupName, &gid)
	assert.Error(t, err)
}

func TestGetIsiGroupMemberListWithResume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiGroupMemberListRespResume)
		*resp = &IsiGroupMemberListRespResume{
			Resume: "resume",
			Members: []*IsiAccessItemFileGroup{
				{
					ID:   "0",
					Name: "name",
				},
			},
		}
	}).Once()
	_, err := getIsiGroupMemberListWithResume(ctx, client, "groupName", "resume")
	assert.Equal(t, nil, err)
}

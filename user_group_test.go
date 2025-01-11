/*
Copyright (c) 2023 Dell Inc, or its subsidiaries.

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

	//"strconv"

	"testing"

	//api "github.com/dell/goisilon/api/v1"
	api "github.com/dell/goisilon/api/v1"
	apiv1 "github.com/dell/goisilon/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dell/goisilon/mocks"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

// Test GetAllGroups() and GetGroupsWithFilter()
func TestGroupGet(t *testing.T) {
	// Get groups with filter
	queryCached := false
	queryResolveNames := false
	queryMemberOf := false
	var queryLimit int32 = 1000

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroupListRespResume)
		*resp = &apiv1.IsiGroupListRespResume{}
	}).Twice() // Expect the Get method to be called twice

	_, err := client.GetGroupsWithFilter(defaultCtx, nil, nil, nil, nil, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	assert.Nil(t, err)

	_, err = client.GetAllGroups(defaultCtx)
	assert.Nil(t, err)
}

func TestGroupCreate(t *testing.T) {
	userName := "test_user_group_member"
	groupName := "test_group_create_options"
	gid := int32(100000)
	queryForce := true
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiUser)
		*resp = &apiv1.IsiUser{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiUser)
		*resp = &apiv1.IsiUser{}
	}).Once() // Expect the Post method to be called once

	client.API.(*mocks.Client).On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the Delete method

	client.API.(*mocks.Client).On("DeleteUserByNameOrUID", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the DeleteUserByNameOrUID method

	_, err := client.CreateUserByName(defaultCtx, userName)
	assert.Nil(t, err)
	client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiUserListResp)
		*resp = &apiv1.IsiUserListResp{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiUserListResp)
		*resp = &apiv1.IsiUserListResp{}
	}).Once() // Expect the Post method to be called once

	client.API.(*mocks.Client).On("GetUserByNameOrUID", mock.Anything, mock.Anything, mock.Anything).Return(&apiv1.IsiUser{
		Name: "test_user_group_member",
		UID:  apiv1.IsiAccessItemFileGroup{ID: "USER:test_user_group_member"},
	}, nil).Once()

	_, err = client.GetUserByNameOrUID(defaultCtx, &userName, nil)
	assertError(t, err)
	uid := 12345
	uid32 := int32(uid)
	member := []api.IsiAuthMemberItem{{Name: &userName, ID: &uid32, Type: "user"}}

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 5 {
			resp, ok := args.Get(5).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 6 {
			resp, ok := args.Get(6).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Post method to be called once

	_, err = client.CreateGroupWithOptions(defaultCtx, groupName, &gid, member, &queryForce, nil, nil)
	assertNoError(t, err)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 5 {
			resp, ok := args.Get(5).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 6 {
			resp, ok := args.Get(6).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Post method to be called once

	client.API.(*mocks.Client).On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the Delete method

	err = client.DeleteGroupByNameOrGID(defaultCtx, nil, &gid)
	assertNoError(t, err)
}

func TestGroupUpdate(t *testing.T) {
	groupName := "test_group_create_update"
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroup)
		*resp = &apiv1.IsiGroup{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiGroup)
		*resp = &apiv1.IsiGroup{}
	}).Once() // Expect the Post method to be called once

	_, err := client.CreatGroupByName(defaultCtx, groupName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroupListResp)
		*resp = &apiv1.IsiGroupListResp{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiGroupListResp)
		*resp = &apiv1.IsiGroupListResp{}
	}).Once()

	_, err = client.GetGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertError(t, err)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroupListResp)
		*resp = &apiv1.IsiGroupListResp{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Put", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiUpdateGroupReq)
		*resp = apiv1.IsiUpdateGroupReq{}
	}).Once() // Expect the Put method to be called once

	newGid := int32(10000)
	err = client.UpdateIsiGroupGIDByNameOrUID(defaultCtx, &groupName, nil, newGid, nil, nil)
	assert.Nil(t, err)
}

// Test GetGroupMembers(), AddGroupMember() and RemoveGroupMember()
func TestGroupMemberAdd(t *testing.T) {
	groupName := "test_group_add_member"

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroup)
		*resp = &apiv1.IsiGroup{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiGroup)
		*resp = &apiv1.IsiGroup{}
	}).Once()
	_, err := client.CreatGroupByName(defaultCtx, groupName)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroupListResp)
		*resp = &apiv1.IsiGroupListResp{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiGroupListResp)
		*resp = &apiv1.IsiGroupListResp{}
	}).Once()

	_, err = client.GetGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertError(t, err)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiGroupMemberListRespResume)
		*resp = &apiv1.IsiGroupMemberListRespResume{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiGroupMemberListRespResume)
		*resp = &apiv1.IsiGroupMemberListRespResume{}
	}).Once()

	_, err = client.GetGroupMembers(defaultCtx, &groupName, nil)
	assertNoError(t, err)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiUser)
		*resp = &apiv1.IsiUser{}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(**apiv1.IsiUser)
		*resp = &apiv1.IsiUser{}
	}).Once() // Expect the Post method to be called once

	userName := "test_user_group_add_member"
	_, err = client.CreateUserByName(defaultCtx, userName)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 5 {
			resp, ok := args.Get(5).(**apiv1.IsiAccessItemFileGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiAccessItemFileGroup{}
			}
		}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 6 {
			resp, ok := args.Get(6).(**apiv1.IsiAccessItemFileGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiAccessItemFileGroup{}
			}
		}
	}).Once() // Expect the Post method to be called once

	client.API.(*mocks.Client).On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the Delete method

	client.API.(*mocks.Client).On("DeleteUserByNameOrUID", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the DeleteUserByNameOrUID method

	userMember := &apiv1.IsiAuthMemberItem{Name: &userName, Type: "user"}
	err = client.AddGroupMember(defaultCtx, &groupName, nil, *userMember)
	assertNoError(t, err)

	client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)

	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 5 {
			resp, ok := args.Get(5).(**apiv1.IsiGroupMemberListRespResume)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroupMemberListRespResume{}
			}
		}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 6 {
			resp, ok := args.Get(6).(**apiv1.IsiGroupMemberListRespResume)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroupMemberListRespResume{}
			}
		}
	}).Once() // Expect the Post method to be called once
	_, err = client.GetGroupMembers(defaultCtx, &groupName, nil)
	assertNoError(t, err)

	userMember1 := &apiv1.IsiAuthMemberItem{Name: &userName, Type: "user"}
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).Calls = nil

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 5 {
			resp, ok := args.Get(5).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Get method to be called once

	client.API.(*mocks.Client).On("Post", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		if len(args) > 6 {
			resp, ok := args.Get(6).(**apiv1.IsiGroup)
			if ok && resp != nil {
				*resp = &apiv1.IsiGroup{}
			}
		}
	}).Once() // Expect the Post method to be called once

	client.API.(*mocks.Client).On("Delete", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil).Once() // Mock the Delete method

	err = client.RemoveGroupMember(defaultCtx, &groupName, nil, *userMember1) // Dereference the pointer
	assertNoError(t, err)
}

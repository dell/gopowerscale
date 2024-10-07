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
	"fmt"
	"strconv"
	"testing"

	api "github.com/dell/goisilon/api/v1"
	"github.com/stretchr/testify/assert"
)

// Test GetAllGroups() and GetGroupsWithFilter()
func TestGroupGet(t *testing.T) {
	// Get groups with filter
	queryCached := false
	queryResolveNames := false
	queryMemberOf := false
	var queryLimit int32 = 1000

	groups, err := client.GetGroupsWithFilter(defaultCtx, nil, nil, nil, nil, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err != nil {
		panic(err)
	}

	// Get All groups
	allGroups, err := client.GetAllGroups(defaultCtx)
	if err != nil {
		panic(err)
	}

	assert.True(t, len(allGroups) >= len(groups))
}

// Test GetGroupByNameOrGID(), CreateGroupWithOptions() and DeleteGroupByNameOrGID()
func TestGroupCreate(t *testing.T) {
	userName := "test_user_group_member"
	groupName := "test_group_create_options"
	gid := int32(100000)
	queryForce := true

	_, err = client.CreateUserByName(defaultCtx, userName)
	if err != nil {
		panic(err)
	}
	defer client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)
	user, err := client.GetUserByNameOrUID(defaultCtx, &userName, nil)
	if err != nil {
		panic(err)
	}
	uid, err := strconv.ParseInt(user.Uid.Id[4:], 10, 32)
	// #nosec G115
	uid32 := int32(uid)
	if err != nil {
		panic(err)
	}
	member := []api.IsiAuthMemberItem{{Name: &userName, Id: &uid32, Type: "user"}}

	_, err = client.CreateGroupWithOptions(defaultCtx, groupName, &gid, member, &queryForce, nil, nil)
	assertNoError(t, err)

	group, err := client.GetGroupByNameOrGID(defaultCtx, nil, &gid)
	assertNoError(t, err)
	assertNotNil(t, group)
	assertEqual(t, fmt.Sprintf("GID:%d", gid), group.Gid.Id)

	err = client.DeleteGroupByNameOrGID(defaultCtx, nil, &gid)
	assertNoError(t, err)

	group, err = client.GetGroupByNameOrGID(defaultCtx, &groupName, &gid)
	assertError(t, err)
	assertNil(t, group)
}

// Test GetGroupByNameOrGID(), CreatGroupByName(), UpdateIsiGroupGIDByNameOrUID() and DeleteGroupByNameOrGID()
func TestGroupUpdate(t *testing.T) {
	groupName := "test_group_create_update"
	_, err := client.CreatGroupByName(defaultCtx, groupName)
	defer client.DeleteGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertNoError(t, err)

	group, err := client.GetGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertNoError(t, err)
	assertNotNil(t, group)

	newGid := int32(10000)
	err = client.UpdateIsiGroupGIDByNameOrUID(defaultCtx, &groupName, nil, newGid, nil, nil)
	assertNoError(t, err)

	groupNew, err := client.GetGroupByNameOrGID(defaultCtx, &groupName, &newGid)
	assertNoError(t, err)
	assertNotNil(t, groupNew)
	assertEqual(t, group.Dn, groupNew.Dn)
	assertEqual(t, group.Provider, groupNew.Provider)
	assertEqual(t, fmt.Sprintf("GID:%d", newGid), groupNew.Gid.Id)
	assertNotEqual(t, group.Gid.Id, groupNew.Gid.Id)
}

// Test GetGroupMembers(), AddGroupMember() and RemoveGroupMember()
func TestGroupMemberAdd(t *testing.T) {
	groupName := "test_group_add_member"
	_, err := client.CreatGroupByName(defaultCtx, groupName)
	defer client.DeleteGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertNoError(t, err)

	group, err := client.GetGroupByNameOrGID(defaultCtx, &groupName, nil)
	assertNoError(t, err)
	assertNotNil(t, group)

	members, err := client.GetGroupMembers(defaultCtx, &groupName, nil)
	assertNoError(t, err)
	assertEqual(t, 0, len(members))

	userName := "test_user_group_add_member"
	_, err = client.CreateUserByName(defaultCtx, userName)
	if err != nil {
		panic(err)
	}
	defer client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)

	userMember := api.IsiAuthMemberItem{Name: &userName, Type: "user"}
	err = client.AddGroupMember(defaultCtx, &groupName, nil, userMember)
	assertNoError(t, err)

	members, err = client.GetGroupMembers(defaultCtx, &groupName, nil)
	assertNoError(t, err)
	assertEqual(t, 1, len(members))
	assertEqual(t, userName, members[0].Name)

	// remove group member
	err = client.RemoveGroupMember(defaultCtx, &groupName, nil, userMember)
	assertNoError(t, err)

	members, err = client.GetGroupMembers(defaultCtx, &groupName, nil)
	assertNoError(t, err)
	assertEqual(t, 0, len(members))
}

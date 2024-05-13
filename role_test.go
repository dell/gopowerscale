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
	"strconv"
	"testing"

	api "github.com/dell/goisilon/api/v1"
)

// Test GetAllRoles() and GetRolesWithFilter()
func TestRoleGet(t *testing.T) {
	// Get roles with filter
	queryResolveNames := true
	var queryLimit int32 = 1000

	roles, err := client.GetRolesWithFilter(defaultCtx, &queryResolveNames, &queryLimit)
	if err != nil {
		panic(err)
	}

	// Get All the roles
	allRoles, err := client.GetAllRoles(defaultCtx)
	if err != nil {
		panic(err)
	}

	assertTrue(t, len(roles) >= len(allRoles))
}

// Test GetRoleByID(), AddRoleMember(), RemoveRoleMember() and IsRoleMemberOf()
func TestRoleMemberAdd(t *testing.T) {
	roleID := "SystemAdmin"
	userName := "test_user_roleMember"

	role, err := client.GetRoleByID(defaultCtx, roleID)
	if err != nil {
		panic(err)
	}
	assertNotNil(t, role)

	_, err = client.CreateUserByName(defaultCtx, userName)
	if err != nil {
		panic(err)
	}
	user, err := client.GetUserByNameOrUID(defaultCtx, &userName, nil)
	if err != nil {
		panic(err)
	}
	defer client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)

	memberUserWithName := api.IsiAuthMemberItem{
		Type: "user",
		Name: &userName,
	}

	isRoleMember, err := client.IsRoleMemberOf(defaultCtx, roleID, memberUserWithName)
	if err != nil {
		panic(err)
	}
	assertFalse(t, isRoleMember)

	err = client.AddRoleMember(defaultCtx, roleID, memberUserWithName)
	if err != nil {
		panic(err)
	}
	isRoleMember, err = client.IsRoleMemberOf(defaultCtx, roleID, memberUserWithName)
	if err != nil {
		panic(err)
	}
	assertTrue(t, isRoleMember)

	err = client.RemoveRoleMember(defaultCtx, roleID, memberUserWithName)
	if err != nil {
		panic(err)
	}
	isRoleMember, err = client.IsRoleMemberOf(defaultCtx, roleID, memberUserWithName)
	if err != nil {
		panic(err)
	}
	assertFalse(t, isRoleMember)

	// add/remove role member by uid
	uid, err := strconv.ParseInt(user.Uid.Id[4:], 10, 32)
	uid32 := int32(uid)
	if err != nil {
		panic(err)
	}
	memberUserWithUID := api.IsiAuthMemberItem{
		Type: "user",
		Id:   &uid32,
	}
	err = client.AddRoleMember(defaultCtx, roleID, memberUserWithUID)
	if err != nil {
		panic(err)
	}
	isRoleMember, err = client.IsRoleMemberOf(defaultCtx, roleID, memberUserWithUID)
	if err != nil {
		panic(err)
	}
	assertTrue(t, isRoleMember)

	err = client.RemoveRoleMember(defaultCtx, roleID, memberUserWithUID)
	if err != nil {
		panic(err)
	}
	isRoleMember, err = client.IsRoleMemberOf(defaultCtx, roleID, memberUserWithUID)
	if err != nil {
		panic(err)
	}
	assertFalse(t, isRoleMember)
}

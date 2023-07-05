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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test GetAllUsers() and GetUsersWithFilter()
func TestUserGet(t *testing.T) {
	// Get users witl filter
	queryNamePrefix := "admin"
	queryCached := false
	queryResolveNames := false
	queryMemberOf := false
	var queryLimit int32 = 1000

	users, err := client.GetUsersWithFilter(defaultCtx, &queryNamePrefix, nil, nil, nil, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err != nil {
		panic(err)
	}

	// Get All the users
	allUsers, err := client.GetAllUsers(defaultCtx)
	if err != nil {
		panic(err)
	}

	assert.True(t, len(allUsers) >= len(users))
}

// Test GetUserByNameOrUID(), CreateUser(), CreateUserWithOptions() and DeleteUserByNameOrUID()
func TestUserCreate(t *testing.T) {

	userName := "test_user_create_options"
	uid := int32(100000)
	email := "test.dell.com"
	pw := "testPW"

	_, err = client.CreateUserWithOptions(defaultCtx, userName, &uid, nil, nil, nil, &email, nil, &pw, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	assertNoError(t, err)

	user, err := client.GetUserByNameOrUID(defaultCtx, nil, &uid)
	assertNoError(t, err)
	assertNotNil(t, user)
	assertEqual(t, fmt.Sprintf("UID:%d", uid), user.Uid.Id)
	assertEqual(t, email, user.Email)

	err = client.DeleteUserByNameOrUID(defaultCtx, nil, &uid)
	assertNoError(t, err)

	user, err = client.GetUserByNameOrUID(defaultCtx, nil, &uid)
	assertError(t, err)
	assertNil(t, user)
}

// Test GetUserByNameOrUID(), CreateUser(), UpdateUserByNameOrUID() and DeleteUserByNameOrUID()
func TestUserUpdate(t *testing.T) {
	userName := "test_user_create_update"

	_, err := client.CreateUserByName(defaultCtx, userName)
	defer client.DeleteUserByNameOrUID(defaultCtx, &userName, nil)
	assertNoError(t, err)

	user, err := client.GetUserByNameOrUID(defaultCtx, &userName, nil)
	assertNoError(t, err)
	assertNotNil(t, user)
	assertEqual(t, user.Email, "")

	email := "test.dell.com"
	pw := "testPW"
	newUid := int32(100000)
	queryForce := true
	err = client.UpdateUserByNameOrUID(defaultCtx, &userName, nil, &queryForce, nil, nil, &email, nil, &pw, nil, nil, nil, &newUid, nil, nil, nil, nil, nil, nil)
	assertNoError(t, err)
	userNew, err := client.GetUserByNameOrUID(defaultCtx, &userName, &newUid)
	assertNoError(t, err)
	assertNotNil(t, userNew)
	assertEqual(t, user.Dn, userNew.Dn)
	assertEqual(t, user.HomeDirectory, userNew.HomeDirectory)
	assertEqual(t, user.Provider, userNew.Provider)
	assertEqual(t, email, userNew.Email)
	assertNotEqual(t, user.Email, userNew.Uid.Id)
	assertEqual(t, fmt.Sprintf("UID:%d", newUid), userNew.Uid.Id)
	assertNotEqual(t, user.Uid.Id, userNew.Uid.Id)
}

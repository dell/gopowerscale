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
	"context"
	"fmt"
	"testing"

	apiv1 "github.com/dell/goisilon/api/v1"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const userPath string = "platform/1/auth/users"

func TestGetUserByNameOrUID(t *testing.T) {
	ctx := context.Background()
	name := "testuser"
	uid := int32(1001)
	expectedUser := &apiv1.IsiUser{Name: "testuser"}
	client := &Client{API: new(mocks.Client)}
	// Initialize the mocked client
	client.API.(*mocks.Client).On("GetIsiUser", ctx, mock.Anything, mock.Anything).Return(expectedUser, nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		// Ensure we are dereferencing the pointer correctly
		resp := args.Get(5).(**apiv1.IsiUserListResp)
		if resp != nil {
			*resp = &apiv1.IsiUserListResp{
				Users: []*apiv1.IsiUser{expectedUser},
			}
		}
	}).Once()

	user, err := client.GetUserByNameOrUID(ctx, &name, &uid)
	assert.NoError(t, err)
	assert.Equal(t, User(expectedUser), user)
}

func TestGetAllUsers(t *testing.T) {
	ctx := context.Background()
	expectedUsers := UserList{&apiv1.IsiUser{Name: "testuser1"}, &apiv1.IsiUser{Name: "testuser2"}}
	client := &Client{API: new(mocks.Client)}
	client.API.(*mocks.Client).On("GetIsiUserList", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiUserListRespResume)
		*resp = &apiv1.IsiUserListRespResume{
			Users: []*apiv1.IsiUser{
				{Name: "testuser1"},
				{Name: "testuser2"},
			},
		}
	})
	users, err := client.GetAllUsers(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedUsers, users)
}

func TestCreateUserByName(t *testing.T) {
	ctx := context.Background()
	name := "testuser"
	id := "1001"
	client := &Client{API: new(mocks.Client)}

	expectedUserReq := &apiv1.IsiUserReq{
		Name: name,
	}

	expectedUserResp := &apiv1.IsiUser{
		ID: id,
	}

	// Mock client setup
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		userPath,
		"",
		mock.Anything,
		mock.Anything,
		expectedUserReq,
		mock.Anything,
	).Run(func(args mock.Arguments) {
		arg := args.Get(6).(**apiv1.IsiUser)
		*arg = expectedUserResp
	}).Return(nil).Once()

	// Call the function under test
	result, err := client.CreateUserByName(ctx, name)

	// Assert that expectations were met and results are as expected
	assert.NoError(t, err)
	assert.Equal(t, id, result)
}

func TestUpdateUserByNameOrUID(t *testing.T) {
	ctx := context.Background()
	name := "testuser"
	uid := int32(123)
	queryForce := true
	queryZone := "zone1"
	queryProvider := "provider1"
	email := "test@example.com"
	homeDirectory := "/home/testuser"
	password := "newpassword"
	fullName := "Test User"
	shell := "/bin/bash"
	primaryGroupName := "users"
	authUserID := "UID:123"
	newUID := int32(456)
	expiry := int32(0)
	primaryGroupID := int32(789)
	enabled := true
	passwordExpires := true
	promptPasswordChange := true
	unlock := true

	// Prepare mock client
	client := &Client{API: &mocks.Client{}}

	client.API.(*mocks.Client).On(
		"Put",
		ctx,
		userPath,
		authUserID,
		mock.Anything,
		mock.Anything,
		&apiv1.IsiUpdateUserReq{
			Email:           &email,
			Enabled:         &enabled,
			Expiry:          &expiry,
			Gecos:           &fullName,
			HomeDirectory:   &homeDirectory,
			Password:        &password,
			PasswordExpires: &passwordExpires,
			PrimaryGroup: &apiv1.IsiAccessItemFileGroup{
				Type: "group",
				ID:   fmt.Sprintf("GID:%d", &primaryGroupID),
				Name: primaryGroupName,
			},
			PromptPasswordChange: &promptPasswordChange,
			Shell:                &shell,
			UID:                  &newUID,
			Unlock:               &unlock,
		},
		nil,
	).Return(nil).Once()

	// Call the function under test
	err := client.UpdateUserByNameOrUID(
		ctx, &name, &uid,
		&queryForce, &queryZone, &queryProvider,
		&email, &homeDirectory, &password, &fullName, &shell, &primaryGroupName,
		&newUID, &primaryGroupID, &expiry, &enabled, &passwordExpires, &promptPasswordChange, &unlock,
	)

	// Assert no error
	assert.NoError(t, err)
}

func TestDeleteUserByNameOrUID(t *testing.T) {
	ctx := context.Background()
	name := "testuser"
	uid := int32(123)

	// Prepare mock client
	client := &Client{API: &mocks.Client{}}

	// Mock authUserID
	authUserID := "UID:123"

	// Mock setup
	client.API.(*mocks.Client).On(
		"Delete",
		ctx,
		userPath,
		authUserID,
		mock.Anything,
		mock.Anything,
		nil,
	).Return(nil).Once()

	// Call the function under test
	err := client.DeleteUserByNameOrUID(ctx, &name, &uid)

	// Assert no error
	assert.NoError(t, err)
}

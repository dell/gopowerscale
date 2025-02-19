/*
Copyright (c) 2024-2025 Dell Inc, or its subsidiaries.

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

var roleMemberPath = "platform/1/auth/roles/%s/members"

func TestGetRoleByID(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	// Test case: Role exists
	client.API.(*mocks.Client).On("GetIsiRole", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiRoleListResp)
		*resp = &apiv1.IsiRoleListResp{
			Roles: []*apiv1.IsiRole{{Name: "testRole", ID: "roleID"}},
		}
	}).Once()
	role, err := client.GetRoleByID(ctx, "roleID")
	assert.Nil(t, err)
	assert.Equal(t, "testRole", role.Name)

	// Test case: Role does not exist
	client.API.(*mocks.Client).On("GetIsiRole", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	role, err = client.GetRoleByID(ctx, "roleID")
	assert.NotNil(t, err)
	assert.Nil(t, role)
}

func TestGetAllRoles(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	// Test case: Roles exist
	client.API.(*mocks.Client).On("GetIsiRoleList", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiRoleListRespResume)
		*resp = &apiv1.IsiRoleListRespResume{
			Roles: []*apiv1.IsiRole{{Name: "role1"}, {Name: "role2"}},
		}
	}).Once()
	roles, err := client.GetAllRoles(ctx)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(roles))

	// Test case: No roles found
	client.API.(*mocks.Client).On("GetIsiRoleList", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	roles, err = client.GetAllRoles(ctx)
	assert.NotNil(t, err)
	assert.Nil(t, roles)
}

func TestIsRoleMemberOf(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	roleID := "roleID"
	memberName := "testMember"
	memberID := int32(123)
	member := apiv1.IsiAuthMemberItem{
		Name: &memberName,
		ID:   &memberID,
		Type: "user",
	}

	// Test case: Member is part of the role
	client.API.(*mocks.Client).On("GetIsiRole", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiRoleListResp)
		*resp = &apiv1.IsiRoleListResp{
			Roles: []*apiv1.IsiRole{{
				ID: roleID,
				Members: []apiv1.IsiAccessItemFileGroup{
					{Name: memberName},
				},
			}},
		}
	}).Once()
	isMember, err := client.IsRoleMemberOf(ctx, roleID, member)
	assert.Nil(t, err)
	assert.True(t, isMember)

	// Test case: Member is not part of the role
	client.API.(*mocks.Client).On("GetIsiRole", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiRoleListResp)
		*resp = &apiv1.IsiRoleListResp{
			Roles: []*apiv1.IsiRole{{Name: memberName, ID: roleID}},
		}
	}).Once()
	isMember, err = client.IsRoleMemberOf(ctx, roleID, member)
	assert.Nil(t, err)
	assert.False(t, isMember)

	// Test case: GetIsiRole returns an error
	client.API.(*mocks.Client).On("GetIsiRole", anyArgs...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	isMember, err = client.IsRoleMemberOf(ctx, roleID, member)
	assert.NotNil(t, err)
	assert.False(t, isMember)
}

func TestAddRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	roleID := "roleID"
	memberName := "testMember"
	memberType := "user"
	member := apiv1.IsiAuthMemberItem{
		Name: &memberName,
		Type: memberType,
	}

	expectedData := &apiv1.IsiAccessItemFileGroup{
		Type: memberType,
		Name: memberName,
	}

	// Expect the Post method to be called with the specified parameters
	client.API.(*mocks.Client).On("Post", ctx, fmt.Sprintf(roleMemberPath, roleID), "", mock.Anything, mock.Anything, expectedData, nil).Return(nil).Once()

	// Call the AddRoleMember method
	err := client.AddRoleMember(ctx, roleID, member)

	// Assert that no error is returned
	assert.Nil(t, err)
}

func TestRemoveRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	roleID := "roleID"
	memberName := "testMember"
	memberType := "user"
	memberID := int32(123)
	member := apiv1.IsiAuthMemberItem{
		Name: &memberName,
		ID:   &memberID,
		Type: memberType,
	}

	// Simulate the authMemberID resolution within the test scope
	authMemberID := fmt.Sprintf("UID:%d", memberID)
	memberPath := fmt.Sprintf(roleMemberPath, roleID)

	// Expect the Delete method to be called with the specified parameters
	client.API.(*mocks.Client).On("Delete", ctx, memberPath, authMemberID, mock.Anything, mock.Anything, nil).Return(nil).Once()

	// Call the RemoveRoleMember method
	err := client.RemoveRoleMember(ctx, roleID, member)

	// Assert that no error is returned
	assert.Nil(t, err)
}

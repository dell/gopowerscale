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

func TestGetIsiRole(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error found")).Once()
	_, err := GetIsiRole(ctx, client, "")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiRoleListResp)
		*resp = &IsiRoleListResp{
			Roles: []*IsiRole{
				{
					ID: "id",
				},
			},
		}
	}).Once()
	_, err = GetIsiRole(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetIsiRoleList(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	x := false
	var y int32 = 5
	client.On("Get", anyArgs...).Return(errors.New("error found")).Once()
	_, err := GetIsiRoleList(ctx, client, &x, &y)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiRoleListRespResume)
		*resp = &IsiRoleListRespResume{
			Roles: []*IsiRole{
				{
					ID: "test",
				},
			},
			Resume: "",
		}
	}).Once()
	_, err = GetIsiRoleList(ctx, client, &x, &y)
	assert.Equal(t, nil, err)
}

func TestAddIsiRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}

	err := AddIsiRoleMember(ctx, client, "", authMember)
	assert.Equal(t, errors.New("member type is wrong, only support user and group"), err)

	authMember.Type = fileGroupTypeUser
	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	err = AddIsiRoleMember(ctx, client, "", authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestRemoveIsiRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}
	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := RemoveIsiRoleMember(ctx, client, "", authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

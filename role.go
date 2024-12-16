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
	"context"
	"fmt"
	"strings"

	api "github.com/dell/goisilon/api/v1"
)

// User maps to an Isilon Role.
type Role *api.IsiRole

// UserList maps to a set of roles.
type RoleList []*api.IsiRole

// GetRoleByID returns a specific role by role id.
func (c *Client) GetRoleByID(ctx context.Context, id string) (Role, error) {
	return api.GetIsiRole(ctx, c.API, id)
}

// GetAllRoles returns all roles on the cluster.
func (c *Client) GetAllRoles(ctx context.Context) (RoleList, error) {
	return c.GetRolesWithFilter(ctx, nil, nil)
}

// GetRolesWithFilter returns roles on the cluster with Optional filter: resolveNames or limit.
func (c *Client) GetRolesWithFilter(ctx context.Context, queryResolveNames *bool, queryLimit *int32) (RoleList, error) {
	return api.GetIsiRoleList(ctx, c.API, queryResolveNames, queryLimit)
}

// IsMemberOf checks if the Uer/Group is member of role.
func (c *Client) IsRoleMemberOf(ctx context.Context, roleID string, member api.IsiAuthMemberItem) (bool, error) {
	role, err := api.GetIsiRole(ctx, c.API, roleID)
	if err != nil {
		return false, err
	}

	for _, m := range role.Members {
		if member.Name != nil && m.Name == *member.Name ||
			member.ID != nil && m.ID == fmt.Sprintf("%sID:%d", strings.ToUpper(member.Type)[0:1], *member.ID) {
			return true, nil
		}
	}

	return false, nil
}

// AddRoleMember adds a member to a specific role,
// Required: roleId, memberType, and memberName/memberId
// memberType can be user/group.
func (c *Client) AddRoleMember(ctx context.Context, roleID string, member api.IsiAuthMemberItem) error {
	return api.AddIsiRoleMember(ctx, c.API, roleID, member)
}

// AddRoleMember removes a member from a specific role,
// Required: roleId, memberType, and memberName/memberId
// memberType can be user/group.
func (c *Client) RemoveRoleMember(ctx context.Context, roleID string, member api.IsiAuthMemberItem) error {
	return api.RemoveIsiRoleMember(ctx, c.API, roleID, member)
}

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
package v1

import (
	"context"
	"fmt"
	"strings"

	"github.com/dell/goisilon/api"
)

// GetIsiRole queries the role by role-id.
func GetIsiRole(ctx context.Context, client api.Client, roleID string) (role *IsiRole, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/roles/<role-id>

	var roleResp *isiRoleListResp
	if err = client.Get(ctx, rolePath, roleID, nil, nil, &roleResp); err != nil {
		return
	}

	if roleResp.Roles != nil && len(roleResp.Roles) > 0 {
		role = roleResp.Roles[0]
		return
	}

	return nil, fmt.Errorf("role not found: %s", roleID)
}

// GetIsiRoleList queries all roles on the cluster, filter by limit or resolveNames.
func GetIsiRoleList(ctx context.Context, client api.Client, queryResolveNames *bool, queryLimit *int32) (roles []*IsiRole, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/roles?resolve_names=&limit=
	values := api.OrderedValues{}
	if queryResolveNames != nil {
		values.StringAdd("resolve_names", fmt.Sprintf("%t", *queryResolveNames))
	}
	if queryLimit != nil {
		values.StringAdd("limit", fmt.Sprintf("%d", *queryLimit))
	}

	var roleListResp *IsiRoleListRespResume
	// First call without Resume param
	if err = client.Get(ctx, rolePath, "", values, nil, &roleListResp); err != nil {
		return
	}

	for {
		roles = append(roles, roleListResp.Roles...)
		if roleListResp.Resume == "" {
			break
		}

		if roleListResp, err = getIsiRoleListWithResume(ctx, client, roleListResp.Resume); err != nil {
			return
		}
	}

	return
}

// getIsiRoleListWithResume queries the next page roles based on resume token.
func getIsiRoleListWithResume(ctx context.Context, client api.Client, resume string) (roles *IsiRoleListRespResume, err error) {
	err = client.Get(ctx, rolePath, "", api.OrderedValues{{[]byte("resume"), []byte(resume)}}, nil, &roles)
	return
}

// AddIsiRoleMember adds a member to the role, member can be user/group.
func AddIsiRoleMember(ctx context.Context, client api.Client, roleID string, member IsiAuthMemberItem) error {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/roles/{role-id}/members
	//					{
	//					 	"type":"user",
	//					 	"name":"user123",
	//						"id":"UID:1522"
	//					}

	memberType := strings.ToLower(member.Type)
	if memberType != fileGroupTypeUser && memberType != fileGroupTypeGroup {
		return fmt.Errorf("member type is wrong, only support %s and %s", fileGroupTypeUser, fileGroupTypeGroup)
	}

	data := &IsiAccessItemFileGroup{
		Type: memberType,
	}
	if member.Id != nil {
		data.Id = fmt.Sprintf("%sID:%d", strings.ToUpper(memberType)[0:1], *member.Id)
	}
	if member.Name != nil {
		data.Name = *member.Name
	}

	return client.Post(ctx, fmt.Sprintf(roleMemberPath, roleID), "", nil, nil, data, nil)
}

// RemoveIsiRoleMember remove a member from the role, member can be user/group.
func RemoveIsiRoleMember(ctx context.Context, client api.Client, roleID string, member IsiAuthMemberItem) error {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/roles/{role-id}/members/<member-id>

	authMemberID, err := getAuthMemberId(member.Type, member.Name, member.Id)
	if err != nil {
		return err
	}

	return client.Delete(ctx, fmt.Sprintf(roleMemberPath, roleID), authMemberID, nil, nil, nil)
}

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

// GetIsiGroup queries the group by group-id.
func GetIsiGroup(ctx context.Context, client api.Client, groupName *string, gid *int32) (group *IsiGroup, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/groups/<group-id>

	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return
	}

	var groupListResp *isiGroupListResp
	if err = client.Get(ctx, groupPath, authGroupId, nil, nil, &groupListResp); err != nil {
		return
	}

	if groupListResp.Groups != nil && len(groupListResp.Groups) > 0 {
		group = groupListResp.Groups[0]
		return
	}

	return nil, fmt.Errorf("group not found: %s", authGroupId)
}

// GetIsiGroupList queries all groups on the cluster with filter,
// filter by namePrefix, domain, zone, provider, cached, resolveNames, memberOf, zone and limit.
func GetIsiGroupList(ctx context.Context, client api.Client,
	queryNamePrefix, queryDomain, queryZone, queryProvider *string,
	queryCached, queryResolveNames, queryMemberOf *bool,
	queryLimit *int32,
) (groups []*IsiGroup, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/groups?limit=&domain=&cached=&resolve_names=&query_member_of=&zone=&provider=&filter=
	values := api.OrderedValues{}
	if queryCached != nil {
		values.StringAdd("cached", fmt.Sprintf("%t", *queryCached))
	}
	if queryDomain != nil {
		values.StringAdd("domain", *queryDomain)
	}
	if queryNamePrefix != nil {
		values.StringAdd("filter", *queryNamePrefix)
	}
	if queryResolveNames != nil {
		values.StringAdd("resolve_names", fmt.Sprintf("%t", *queryResolveNames))
	}
	if queryMemberOf != nil {
		values.StringAdd("query_member_of", fmt.Sprintf("%t", *queryMemberOf))
	}
	if queryZone != nil {
		values.StringAdd("zone", *queryZone)
	}
	if queryProvider != nil {
		values.StringAdd("provider", *queryProvider)
	}
	if queryLimit != nil {
		values.StringAdd("limit", fmt.Sprintf("%d", *queryLimit))
	}

	var groupListResp *IsiGroupListRespResume
	// First call without Resume param
	if err = client.Get(ctx, groupPath, "", values, nil, &groupListResp); err != nil {
		return
	}

	for {
		groups = append(groups, groupListResp.Groups...)
		if groupListResp.Resume == "" {
			break
		}

		if groupListResp, err = getIsiGroupListWithResume(ctx, client, groupListResp.Resume); err != nil {
			return nil, err
		}
	}
	return
}

// getIsiGroupListWithResume queries the next page groups based on resume token.
func getIsiGroupListWithResume(ctx context.Context, client api.Client, resume string) (groups *IsiGroupListRespResume, err error) {
	err = client.Get(ctx, groupPath, "", api.OrderedValues{{[]byte("resume"), []byte(resume)}}, nil, &groups)
	return
}

// GetIsiGroupMembers retrieves the members of a group.
func GetIsiGroupMembers(ctx context.Context, client api.Client, groupName *string, gid *int32) (members []*IsiAccessItemFileGroup, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/groups/{group-id}/members

	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return
	}

	var groupMemberListResp *IsiGroupMemberListRespResume
	// First call without Resume param
	if err = client.Get(ctx, fmt.Sprintf(groupMemberPath, authGroupId), "", nil, nil, &groupMemberListResp); err != nil {
		return
	}

	for {
		members = append(members, groupMemberListResp.Members...)
		if groupMemberListResp.Resume == "" {
			break
		}

		if groupMemberListResp, err = getIsiGroupMemberListWithResume(ctx, client, authGroupId, groupMemberListResp.Resume); err != nil {
			return nil, err
		}
	}
	return
}

// getIsiGroupMemberListWithResume queries the next page group members based on resume token.
func getIsiGroupMemberListWithResume(ctx context.Context, client api.Client, groupId, resume string) (members *IsiGroupMemberListRespResume, err error) {
	err = client.Get(ctx, fmt.Sprintf(groupMemberPath, groupId), "", api.OrderedValues{{[]byte("resume"), []byte(resume)}}, nil, &members)
	return
}

// AddIsiGroupMember adds a member to the group, member can be a user/group.
func AddIsiGroupMember(ctx context.Context, client api.Client, groupName *string, gid *int32, member IsiAuthMemberItem) error {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/groups/{group-id}/members
	//					{
	//					 	"type":"user",
	//					 	"name":"user123",
	//						"id":"UID:1522"
	//					}

	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return err
	}

	memberType := strings.ToLower(member.Type)
	if memberType != fileGroupTypeUser && memberType != fileGroupTypeGroup {
		return fmt.Errorf("member type is wrong, only support %s and %s", fileGroupTypeUser, fileGroupTypeGroup)
	}

	data := &IsiAccessItemFileGroup{}
	if member.Id != nil {
		data.Id = fmt.Sprintf("%sID:%d", strings.ToUpper(memberType)[0:1], *member.Id)
	}
	if member.Name != nil && *member.Name != "" {
		data.Name = *member.Name
	}

	return client.Post(ctx, fmt.Sprintf(groupMemberPath, authGroupId), "", nil, nil, data, nil)
}

// RemoveIsiGroupMember remove a member from the group, member can be user/group.
func RemoveIsiGroupMember(ctx context.Context, client api.Client, groupName *string, gid *int32, member IsiAuthMemberItem) error {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/groups/{group-id}/members/<member-id>

	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return err
	}

	authMemberId, err := getAuthMemberId(member.Type, member.Name, member.Id)
	if err != nil {
		return err
	}

	return client.Delete(ctx, fmt.Sprintf(groupMemberPath, authGroupId), authMemberId, nil, nil, nil)
}

// CreateIsiGroup creates a new group.
func CreateIsiGroup(ctx context.Context, client api.Client,
	name string, gid *int32, members []IsiAuthMemberItem,
	queryForce *bool, queryZone, queryProvider *string,
) (string, error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/auth/groups?force=&zone=&provider=
	// 				{
	// 					"gid": "int",
	// 					"name": "str",
	// 					"members":
	// 					[
	//					{
	// 						"type":"str",
	// 						"id":"str",
	// 						"name":"str"
	// 					},
	// 					],
	//				}

	values := api.OrderedValues{}
	if queryForce != nil {
		values.StringAdd("force", fmt.Sprintf("%t", *queryForce))
	}
	if queryZone != nil {
		values.StringAdd("zone", *queryZone)
	}
	if queryProvider != nil {
		values.StringAdd("provider", *queryProvider)
	}

	var memberList []IsiAccessItemFileGroup
	for _, m := range members {

		memberType := strings.ToLower(m.Type)
		if memberType != fileGroupTypeUser && memberType != fileGroupTypeGroup {
			return "", fmt.Errorf("member type is wrong, only support %s and %s", fileGroupTypeUser, fileGroupTypeGroup)
		}

		member := &IsiAccessItemFileGroup{
			Type: memberType,
		}

		if m.Id != nil {
			member.Id = fmt.Sprintf("%sID:%d", strings.ToUpper(memberType)[0:1], *m.Id)
		} else if m.Name != nil && *m.Name != "" {
			member.Name = *m.Name
		} else {
			continue
		}

		memberList = append(memberList, *member)
	}

	data := &IsiGroupReq{
		Name:    name,
		Gid:     gid,
		Members: memberList,
	}

	var groupResp *IsiGroup
	if err := client.Post(ctx, groupPath, "", values, nil, data, &groupResp); err != nil {
		return "", err
	}

	return groupResp.Id, nil
}

// UpdateIsiGroupGID updates the group's gid.
func UpdateIsiGroupGID(ctx context.Context, client api.Client, groupName *string, gid *int32, newGid int32,
	queryZone, queryProvider *string,
) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/auth/groups/<group-id>?force=true&zone=&provider=
	// 				{
	// 					"gid": int
	//				}

	// set force to true to change group's gid
	values := api.OrderedValues{{[]byte("force"), []byte("true")}}

	if queryZone != nil {
		values.StringAdd("zone", *queryZone)
	}

	if queryProvider != nil {
		values.StringAdd("provider", *queryProvider)
	}

	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return
	}

	return client.Put(ctx, groupPath, authGroupId, values, nil, &IsiUpdateGroupReq{newGid}, nil)
}

// DeleteIsiGroup removes the group by group-id.
func DeleteIsiGroup(ctx context.Context, client api.Client, groupName *string, gid *int32) (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/auth/groups/<group-id>
	authGroupId, err := getAuthMemberId(fileGroupTypeGroup, groupName, gid)
	if err != nil {
		return
	}

	return client.Delete(ctx, groupPath, authGroupId, nil, nil, nil)
}

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

	api "github.com/dell/goisilon/api/v1"
)

// Group maps to an Isilon Group.
type Group *api.IsiGroup

// GroupList maps to a set of groups.
type GroupList []*api.IsiGroup

// GroupMemberList maps to a set of group members.
type GroupMemberList []*api.IsiAccessItemFileGroup

// GetGroupByNameOrGID returns a specific group by group name or gid.
func (c *Client) GetGroupByNameOrGID(ctx context.Context, name *string, gid *int32) (Group, error) {
	return api.GetIsiGroup(ctx, c.API, name, gid)
}

// GetAllGroups returns all groups on the cluster.
func (c *Client) GetAllGroups(ctx context.Context) (GroupList, error) {
	return c.GetGroupsWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, nil)
}

// GetGroupsWithFilter returns groups on the cluster with
// Optional filter: namePrefix, domain, zone, provider, cached, resolveNames, memberOf, zone and limit.
func (c *Client) GetGroupsWithFilter(ctx context.Context,
	queryNamePrefix, queryDomain, queryZone, queryProvider *string,
	queryCached, queryResolveNames, queryMemberOf *bool,
	queryLimit *int32,
) (GroupList, error) {
	return api.GetIsiGroupList(ctx, c.API, queryNamePrefix, queryDomain, queryZone, queryProvider, queryCached, queryResolveNames, queryMemberOf, queryLimit)
}

// GetGroupMembers retrieves the members of a group by name or gid.
func (c *Client) GetGroupMembers(ctx context.Context, name *string, gid *int32) (GroupMemberList, error) {
	return api.GetIsiGroupMembers(ctx, c.API, name, gid)
}

// AddGroupMember adds a member to a specific group,
// Required: groupName/gid, member,
// member can be a user or group.
func (c *Client) AddGroupMember(ctx context.Context, name *string, gid *int32, member api.IsiAuthMemberItem) error {
	return api.AddIsiGroupMember(ctx, c.API, name, gid, member)
}

// RemoveGroupMember removes a member from a specific group,
// Required: groupName/gid, member,
// member can be a user or group.
func (c *Client) RemoveGroupMember(ctx context.Context, name *string, gid *int32, member api.IsiAuthMemberItem) error {
	return api.RemoveIsiGroupMember(ctx, c.API, name, gid, member)
}

// CreatGroupByName creates a new group with name.
func (c *Client) CreatGroupByName(ctx context.Context, name string) (string, error) {
	return c.CreateGroupWithOptions(ctx, name, nil, nil, nil, nil, nil)
}

// CreateGroupWithOptions creates a new group with name(required), gid(optional), and members(optional),
// Optional filter: force, zone and provider.
func (c *Client) CreateGroupWithOptions(
	ctx context.Context, name string, gid *int32, members []api.IsiAuthMemberItem,
	queryForce *bool, queryZone, queryProvider *string,
) (string, error) {
	return api.CreateIsiGroup(ctx, c.API, name, gid, members, queryForce, queryZone, queryProvider)
}

// UpdateIsiGroupGIDByNameOrUID modifies a specific group's gid by group name or gid with
// Required: newGid,
// Optional filter: zone and provider.
func (c *Client) UpdateIsiGroupGIDByNameOrUID(
	ctx context.Context, name *string, gid *int32, newGid int32, queryZone, queryProvider *string,
) error {
	return api.UpdateIsiGroupGID(ctx, c.API, name, gid, newGid, queryZone, queryProvider)
}

// DeleteGroupByNameOrGID deletes a specific group by group name or gid.
func (c *Client) DeleteGroupByNameOrGID(ctx context.Context, name *string, gid *int32) error {
	return api.DeleteIsiGroup(ctx, c.API, name, gid)
}

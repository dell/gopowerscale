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

// User maps to an Isilon User.
type User *api.IsiUser

// UserList maps to a set of users
type UserList []*api.IsiUser

// GetUserByNameOrUID returns a specific user by user name or uid.
func (c *Client) GetUserByNameOrUID(ctx context.Context, name *string, uid *int32) (User, error) {
	return api.GetIsiUser(ctx, c.API, name, uid)
}

// GetAllUsers returns all users on the cluster
func (c *Client) GetAllUsers(ctx context.Context) (UserList, error) {
	return c.GetUsersWithFilter(ctx, nil, nil, nil, nil, nil, nil, nil, nil)
}

// GetUsersWithFilter returns users on the cluster with
// Optional filter: namePrefix, domain, zone, provider, cached, resolveNames, memberOf, zone and limit.
func (c *Client) GetUsersWithFilter(ctx context.Context,
	queryNamePrefix, queryDomain, queryZone, queryProvider *string,
	queryCached, queryResolveNames, queryMemberOf *bool,
	queryLimit *int32,
) (UserList, error) {
	return api.GetIsiUserList(ctx, c.API, queryNamePrefix, queryDomain, queryZone, queryProvider, queryCached, queryResolveNames, queryMemberOf, queryLimit)
}

// CreateUserByName creates a new user with name.
func (c *Client) CreateUserByName(ctx context.Context, name string) (string, error) {
	return c.CreateUserWithOptions(ctx, name, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
}

// CreateUserWithOptions creates a new user with name(required) and Uid (Optional),
// Optional filter: force, zone and provider,
// Optional configuration: email, homeDirectory, password, primaryGroupName, fullName, shell
// ,uid, primaryGroupId, expiry, enabled, passwordExpires, promptPasswordChange and unlock.
func (c *Client) CreateUserWithOptions(
	ctx context.Context, name string, uid *int32, queryForce *bool, queryZone, queryProvider *string,
	email, homeDirectory, password, fullName, shell, primaryGroupName *string,
	primaryGroupID, expiry *int32, enabled, passwordExpires, promptPasswordChange, unlock *bool,
) (string, error) {
	return api.CreateIsiUser(
		ctx, c.API, name,
		queryForce, queryZone, queryProvider,
		email, homeDirectory, password, primaryGroupName, fullName, shell,
		uid, primaryGroupID, expiry, enabled, passwordExpires, promptPasswordChange, unlock)
}

// UpdateUserByName modifies a specific user by user name or uid with
// Optional filter: force, zone and provider,
// Optional configuration: email, homeDirectory, password, primaryGroupName, fullName, shell
// newUid, primaryGroupId, expiry, enabled, passwordExpires, promptPasswordChange and unlock.
func (c *Client) UpdateUserByNameOrUID(
	ctx context.Context, name *string, uid *int32,
	queryForce *bool, queryZone, queryProvider *string,
	email, homeDirectory, password, fullName, shell, primaryGroupName *string,
	newUID, primaryGroupID, expiry *int32, enabled, passwordExpires, promptPasswordChange, unlock *bool,
) error {
	return api.UpdateIsiUser(
		ctx, c.API, name, uid,
		queryForce, queryZone, queryProvider,
		email, homeDirectory, password, primaryGroupName, fullName, shell,
		newUID, primaryGroupID, expiry, enabled, passwordExpires, promptPasswordChange, unlock)
}

// DeleteUserByNameOrUID deletes a specific user by user name or uid.
func (c *Client) DeleteUserByNameOrUID(ctx context.Context, name *string, uid *int32) error {
	return api.DeleteIsiUser(ctx, c.API, name, uid)
}

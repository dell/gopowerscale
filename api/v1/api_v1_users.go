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

	"github.com/dell/goisilon/api"
)

// GetIsiUser queries the user by user user-id.
func GetIsiUser(ctx context.Context, client api.Client, userName *string, uid *int32) (user *IsiUser, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/users/<user-id>

	authUserID, err := getAuthMemberId(fileGroupTypeUser, userName, uid)
	if err != nil {
		return
	}

	var userListResp *isiUserListResp
	if err = client.Get(ctx, userPath, authUserID, nil, nil, &userListResp); err != nil {
		return
	}

	if userListResp.Users != nil && len(userListResp.Users) > 0 {
		user = userListResp.Users[0]
		return
	}

	return nil, fmt.Errorf("user not found: %s", authUserID)
}

// GetIsiUserList queries all users on the cluster,
// filter by namePrefix, domain, zone, provider, cached, resolveNames, memberOf, zone and limit.
func GetIsiUserList(ctx context.Context, client api.Client,
	queryNamePrefix, queryDomain, queryZone, queryProvider *string,
	queryCached, queryResolveNames, queryMemberOf *bool,
	queryLimit *int32,
) (users []*IsiUser, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/users?limit=&cached=&resolve_names=&query_member_of=&zone=&provider=
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

	var userListResp *IsiUserListRespResume
	// First call without Resume param
	if err = client.Get(ctx, userPath, "", values, nil, &userListResp); err != nil {
		return
	}
	for {
		users = append(users, userListResp.Users...)
		if userListResp.Resume == "" {
			break
		}

		if userListResp, err = getIsiUserListWithResume(ctx, client, userListResp.Resume); err != nil {
			return nil, err
		}
	}
	return
}

// getIsiUserListWithResume queries the next page users based on resume token.
func getIsiUserListWithResume(ctx context.Context, client api.Client, resume string) (users *IsiUserListRespResume, err error) {
	err = client.Get(ctx, userPath, "", api.OrderedValues{{[]byte("resume"), []byte(resume)}}, nil, &users)
	return
}

// CreateIsiUser creates a new user.
func CreateIsiUser(ctx context.Context, client api.Client, name string,
	queryForce *bool, queryZone, queryProvider *string,
	email, homeDirectory, password, primaryGroupName, fullName, shell *string,
	uid, primaryGroupId, expiry *int32, enabled, passwordExpires, promptPasswordChange, unlock *bool,
) (string, error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/auth/users?force=&zone=&provider=
	// 				{
	//					"email": "str",
	// 					"enabled": "bool",
	// 					"expiry": "int",
	//					"gecos": "str",
	// 					"home_directory": "str",
	// 					"password": "str",
	// 					"password_expires": "bool",
	// 					"primary_group":
	//					 {
	// 						"type":"str",
	// 						"id":"str",
	// 						"name":"str"
	// 					 },
	// 					"prompt_password_change": "bool",
	// 					"shell": "str",
	// 					"uid": "int",
	// 					"unlock": "bool",
	// 					"name": "str"
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
	var primaryGroup *IsiAccessItemFileGroup
	if primaryGroupId != nil || primaryGroupName != nil {
		primaryGroup = &IsiAccessItemFileGroup{
			Type: "group",
		}
		if primaryGroupId != nil {
			primaryGroup.Id = fmt.Sprintf("GID:%d", primaryGroupId)
		}
		if primaryGroupName != nil {
			primaryGroup.Name = *primaryGroupName
		}
	}

	data := &IsiUserReq{
		Email:                email,
		Enabled:              enabled,
		Expiry:               expiry,
		Gecos:                fullName,
		HomeDirectory:        homeDirectory,
		Name:                 name,
		Password:             password,
		PasswordExpires:      passwordExpires,
		PrimaryGroup:         primaryGroup,
		PromptPasswordChange: promptPasswordChange,
		Shell:                shell,
		Uid:                  uid,
		Unlock:               unlock,
	}

	var userResp *IsiUser
	if err := client.Post(ctx, userPath, "", values, nil, data, &userResp); err != nil {
		return "", err
	}

	return userResp.Id, nil
}

// UpdateIsiUser updates the user.
func UpdateIsiUser(ctx context.Context, client api.Client, userName *string, uid *int32,
	queryForce *bool, queryZone, queryProvider *string,
	email, homeDirectory, password, primaryGroupName, fullName, shell *string,
	newUid, primaryGroupId, expiry *int32, enabled, passwordExpires, promptPasswordChange, unlock *bool,
) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/auth/users/<user-id>?force=&zone=&provider=
	// 				{
	//					"email": "str",
	// 					"enabled": "bool",
	// 					"expiry": "int",
	//					"gecos": "str",
	// 					"home_directory": "str",
	// 					"password": "str",
	// 					"password_expires": "bool",
	// 					"primary_group":
	//					 {
	// 						"type":"str",
	// 						"id":"str",
	// 						"name":"str"
	// 					 },
	// 					"prompt_password_change": "bool",
	// 					"shell": "str",
	// 					"uid": "int",
	// 					"unlock": "bool"
	//				}
	authUserID, err := getAuthMemberId(fileGroupTypeUser, userName, uid)
	if err != nil {
		return
	}

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

	var primaryGroup *IsiAccessItemFileGroup
	if primaryGroupId != nil || primaryGroupName != nil {
		primaryGroup = &IsiAccessItemFileGroup{
			Type: "group",
		}
		if primaryGroupId != nil {
			primaryGroup.Id = fmt.Sprintf("GID:%d", primaryGroupId)
		}
		if primaryGroupName != nil {
			primaryGroup.Name = *primaryGroupName
		}
	}

	data := &IsiUpdateUserReq{
		Email:                email,
		Enabled:              enabled,
		Expiry:               expiry,
		Gecos:                fullName,
		HomeDirectory:        homeDirectory,
		Password:             password,
		PasswordExpires:      passwordExpires,
		PrimaryGroup:         primaryGroup,
		PromptPasswordChange: promptPasswordChange,
		Shell:                shell,
		Uid:                  newUid,
		Unlock:               unlock,
	}

	return client.Put(ctx, userPath, authUserID, values, nil, data, nil)
}

// DeleteIsiUser removes the user by user-id.
func DeleteIsiUser(ctx context.Context, client api.Client, userName *string, uid *int32) (err error) {
	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/auth/users/<user-id>

	authUserID, err := getAuthMemberId(fileGroupTypeUser, userName, uid)
	if err != nil {
		return
	}

	return client.Delete(ctx, userPath, authUserID, nil, nil, nil)
}

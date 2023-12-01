/*
Copyright (c) 2022 Dell Inc, or its subsidiaries.

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

	api "github.com/dell/goisilon/api/v2"
)

// ACL is an Isilon Access Control List used for managing an object's security.
type ACL *api.ACL

// GetVolumeACL returns the ACL for a volume.
func (c *Client) GetVolumeACL(
	ctx context.Context,
	volumeName string,
) (ACL, error) {
	return api.ACLInspect(ctx, c.API, volumeName)
}

// SetVolumeOwnerToCurrentUser sets the owner for a volume to the user that
// was used to connect to the API.
func (c *Client) SetVolumeOwnerToCurrentUser(
	ctx context.Context,
	volumeName string,
) error {
	return c.SetVolumeOwner(ctx, volumeName, c.API.User())
}

// SetVolumeOwner sets the owner for a volume.
func (c *Client) SetVolumeOwner(
	ctx context.Context,
	volumeName, userName string,
) error {
	mode := api.FileMode(0o777)

	return api.ACLUpdate(
		ctx,
		c.API,
		volumeName,
		&api.ACL{
			Action:        &api.PActionTypeReplace,
			Authoritative: &api.PAuthoritativeTypeMode,
			Owner: &api.Persona{
				ID: &api.PersonaID{
					ID:   userName,
					Type: api.PersonaIDTypeUser,
				},
			},
			Mode: &mode,
		})
}

// SetVolumeMode sets the permissions to the specified mode (chmod)
func (c *Client) SetVolumeMode(
	ctx context.Context,
	volumeName string, mode int,
) error {
	filemode := api.FileMode(mode)

	return api.ACLUpdate(
		ctx,
		c.API,
		volumeName,
		&api.ACL{
			Action:        &api.PActionTypeReplace,
			Authoritative: &api.PAuthoritativeTypeMode,
			Mode:          &filemode,
		})
}

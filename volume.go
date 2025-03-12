/*
Copyright (c) 2019-2022 Dell Inc, or its subsidiaries.

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
	"os"
	"path"
	"strings"
	"sync"

	log "github.com/akutz/gournal"
	apiv1 "github.com/dell/goisilon/api/v1"
	apiv2 "github.com/dell/goisilon/api/v2"
)

// Volume represents an Isilon Volume (namespace API).
type Volume *apiv1.IsiVolume

// VolumeChildren is a list of a container's children.
type VolumeChildren apiv2.ContainerChildList

// VolumeChildrenMap returns a map of all descendent children of a container, where the key is the path.
type VolumeChildrenMap map[string]*apiv2.ContainerChild

// GetVolume returns a specific volume by name or ID
func (c *Client) GetVolume(
	ctx context.Context, id, name string,
) (Volume, error) {
	if id != "" {
		name = id
	}
	volume, err := apiv1.GetIsiVolume(ctx, c.API, name)
	if err != nil {
		return nil, err
	}
	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: volume.AttributeMap}
	return isiVolume, nil
}

// GetVolumeWithIsiPath returns a specific volume by name or ID with isiPath
func (c *Client) GetVolumeWithIsiPath(
	ctx context.Context, isiPath, id, name string,
) (Volume, error) {
	if id != "" {
		name = id
	}
	volume, err := apiv1.GetIsiVolumeWithIsiPath(ctx, c.API, isiPath, name)
	if err != nil {
		return nil, err
	}
	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: volume.AttributeMap}
	return isiVolume, nil
}

// IsVolumeExistent checks whether a volume already exists
func (c *Client) IsVolumeExistent(
	ctx context.Context, id, name string,
) bool {
	// Need change here
	if id != "" {
		name = id
	}
	// query the volume without using the metadata parameter, use whether an error (typically, JSONError instance with "404 Not Found" status code) is returned to indicate whether the volume already exists.
	err := apiv1.GetIsiVolumeWithoutMetadata(ctx, c.API, name)

	if err == nil {
		log.Debug(ctx, "the query of volume (id '%s', name '%s') did not return an error, regard the volume as existent.", id, name)
	} else {
		log.Debug(ctx, "the query of volume (id '%s', name '%s') returned an error, regard the volume as non-existent. error : '%v'", id, name, err)
	}

	return err == nil
}

// IsVolumeExistentWithIsiPath checks whether a volume already exists with isiPath
func (c *Client) IsVolumeExistentWithIsiPath(
	ctx context.Context, isiPath, id, name string,
) bool {
	// Need change here
	if id != "" {
		name = id
	}
	// query the volume without using the metadata parameter, use whether an error (typically, JSONError instance with "404 Not Found" status code) is returned to indicate whether the volume already exists.
	err := apiv1.GetIsiVolumeWithoutMetadataWithIsiPath(ctx, c.API, isiPath, name)

	if err == nil {
		log.Debug(ctx, "the query of volume (id '%s', name '%s') did not return an error, regard the volume as existent.", id, name)
	} else {
		log.Debug(ctx, "the query of volume (id '%s', name '%s') returned an error, regard the volume as non-existent. error : '%v'", id, name, err)
	}

	return err == nil
}

// GetVolumes returns a list of volumes
func (c *Client) GetVolumes(ctx context.Context) ([]Volume, error) {
	volumes, err := apiv1.GetIsiVolumes(ctx, c.API)
	if err != nil {
		return nil, err
	}
	var isiVolumes []Volume
	for _, volume := range volumes.Children {
		newVolume := &apiv1.IsiVolume{Name: volume.Name}
		isiVolumes = append(isiVolumes, newVolume)
	}
	return isiVolumes, nil
}

// CreateVolume creates a volume
func (c *Client) CreateVolume(
	ctx context.Context, name string,
) (Volume, error) {
	_, err := apiv1.CreateIsiVolume(ctx, c.API, name)
	if err != nil {
		return nil, err
	}

	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

// CreateVolumeWithIsipath creates a volume with isiPath
func (c *Client) CreateVolumeWithIsipath(
	ctx context.Context, isiPath, name, isiVolumePathPermissions string,
) (Volume, error) {
	_, err := apiv1.CreateIsiVolumeWithIsiPath(ctx, c.API, isiPath, name, isiVolumePathPermissions)
	if err != nil {
		return nil, err
	}

	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

// CreateVolumeWithIsipathMetaData creates a volume with isiPath
func (c *Client) CreateVolumeWithIsipathMetaData(
	ctx context.Context, isiPath, name, isiVolumePathPermissions string, metadata map[string]string,
) (Volume, error) {
	_, err := apiv1.CreateIsiVolumeWithIsiPathMetaData(ctx, c.API, isiPath, name, isiVolumePathPermissions, metadata)
	if err != nil {
		return nil, err
	}

	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

// CreateVolumeNoACL creates a volume without setting ACL
func (c *Client) CreateVolumeNoACL(
	ctx context.Context, name string,
) (Volume, error) {
	_, err := apiv1.CreateIsiVolumeWithACL(ctx, c.API, name, "0777")
	if err != nil {
		return nil, err
	}

	isiVolume := &apiv1.IsiVolume{Name: name, AttributeMap: nil}
	return isiVolume, nil
}

// DeleteVolume deletes a volume
func (c *Client) DeleteVolume(
	ctx context.Context, name string,
) error {
	_, err := apiv1.DeleteIsiVolume(ctx, c.API, name)
	return err
}

// DeleteVolumeWithIsiPath deletes a volume with isiPath
func (c *Client) DeleteVolumeWithIsiPath(
	ctx context.Context, isiPath, name string,
) error {
	_, err := apiv1.DeleteIsiVolumeWithIsiPath(ctx, c.API, isiPath, name)
	return err
}

// ConcurrentHTTPConnections is the number of allowed concurrent HTTP
// connections for API functions that attempt to send multiple API calls at
// once.
var ConcurrentHTTPConnections = 2

func newConcurrentHTTPChan() chan bool {
	c := make(chan bool, ConcurrentHTTPConnections)
	for i := 0; i < ConcurrentHTTPConnections; i++ {
		c <- true
	}
	return c
}

// ForceDeleteVolume force deletes a volume by resetting the ownership of
// all descendent directories to the current user prior to issuing a delete
// call.
func (c *Client) ForceDeleteVolume(ctx context.Context, name string) error {
	var (
		user       = c.API.User()
		vpl        = len(c.API.VolumesPath()) + 1
		errs       = make(chan error)
		queryDone  = make(chan int)
		childPaths = make(chan string)
		setACLWait = &sync.WaitGroup{}
		setACLChan = newConcurrentHTTPChan()
		setACLDone = make(chan int)
		mode       = apiv2.FileMode(0o755)
		acl        = &apiv2.ACL{
			Action:        &apiv2.PActionTypeReplace,
			Authoritative: &apiv2.PAuthoritativeTypeMode,
			Owner: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   user,
					Type: apiv2.PersonaIDTypeUser,
				},
			},
			Mode: &mode,
		}
	)

	go func() {
		queryChan, queryErrs := apiv2.ContainerChildrenGetQuery(
			ctx, c.API, name, 1000, -1, "container", "ASC",
			[]string{"container_path", "name"},
			[]string{"owner", "name", "container_path"})

		go func() {
			if err := <-queryErrs; err != nil {
				errs <- err
				close(errs)
			}
		}()

		go func() {
			for child := range queryChan {
				if strings.EqualFold(user, *child.Owner) {
					continue
				}
				setACLWait.Add(1)
				go func(s string) {
					childPaths <- s
				}(path.Join(*child.Path, *child.Name)[vpl:])
			}
			close(queryDone)
		}()
	}()

	go func() {
		for childPath := range childPaths {
			go func(childPath string) {
				<-setACLChan
				if err := apiv2.ACLUpdate(
					ctx,
					c.API,
					childPath,
					acl); err != nil {
					go func(childPath string) {
						childPaths <- childPath
					}(childPath)
				} else {
					setACLWait.Done()
				}
				setACLChan <- true
			}(childPath)
		}
	}()

	go func() {
		<-queryDone
		setACLWait.Wait()
		close(setACLDone)
	}()

	select {
	case <-setACLDone:
	case err := <-errs:
		if err != nil {
			return err
		}
	}

	return c.DeleteVolume(ctx, name)
}

// CopyVolume creates a volume based on an existing volume
func (c *Client) CopyVolume(
	ctx context.Context, src, dest string,
) (Volume, error) {
	_, err := apiv1.CopyIsiVolume(ctx, c.API, src, dest)
	if err != nil {
		return nil, err
	}

	return c.GetVolume(ctx, dest, dest)
}

// CopyVolumeWithIsiPath creates a volume based on an existing volume
func (c *Client) CopyVolumeWithIsiPath(
	ctx context.Context, isiPath, src, dest string,
) (Volume, error) {
	res, err := apiv1.CopyIsiVolumeWithIsiPath(ctx, c.API, isiPath, src, dest)
	if err != nil {
		return nil, err
	}
	if res != nil && res.Success == false {
		log.Error(ctx, "error encountered while cloning volume. error : '%v'", res)
		return nil, err
	}

	return c.GetVolumeWithIsiPath(ctx, isiPath, dest, dest)
}

// ExportVolume exports a volume
func (c *Client) ExportVolume(
	ctx context.Context, name string,
) (int, error) {
	return c.Export(ctx, name)
}

// ExportVolumeWithZone exports a volume in the specified access zone
func (c *Client) ExportVolumeWithZone(
	ctx context.Context, name, zone, description string,
) (int, error) {
	return c.ExportWithZone(ctx, name, zone, description)
}

// ExportVolumeWithZoneAndPath exports a volume in the specified access zone and path
func (c *Client) ExportVolumeWithZoneAndPath(
	ctx context.Context, path, zone, description string,
) (int, error) {
	return c.ExportWithZoneAndPath(ctx, path, zone, description)
}

// UnexportVolume stops exporting a volume
func (c *Client) UnexportVolume(
	ctx context.Context, name string,
) error {
	return c.Unexport(ctx, name)
}

// QueryVolumeChildren retrieves a list of all of a volume's descendent files
// and directories.
func (c *Client) QueryVolumeChildren(
	ctx context.Context, name string,
) (VolumeChildrenMap, error) {
	return apiv2.ContainerChildrenMapAll(ctx, c.API, name)
}

// CreateVolumeDir creates a directory inside a volume.
func (c *Client) CreateVolumeDir(
	ctx context.Context,
	volumeName, dirPath string,
	fileMode os.FileMode,
	overwrite, recursive bool,
) error {
	return apiv2.ContainerCreateDir(
		ctx, c.API, volumeName, dirPath,
		apiv2.FileMode(fileMode), overwrite, recursive)
}

// GetVolumeExportMap returns a map that relates Volumes to their corresponding
// Exports. This function uses an Export's "clients" property to define the
// relationship. The flag "includeRootClients" can be set to "true" in order to
// also inspect the "root_clients" property of an Export when determining the
// Volume-to-Export relationship.
func (c *Client) GetVolumeExportMap(
	ctx context.Context,
	includeRootClients bool,
) (map[Volume]Export, error) {
	volumes, err := c.GetVolumes(ctx)
	if err != nil {
		return nil, err
	}
	exports, err := c.GetExports(ctx)
	if err != nil {
		return nil, err
	}

	volToExpMap := map[Volume]Export{}

	for _, v := range volumes {
		vp := c.API.VolumePath(v.Name)
		for _, e := range exports {
			if e.Clients == nil {
				continue
			}
			for _, p := range *e.Clients {
				if vp == p {
					if _, ok := volToExpMap[v]; ok {
						log.WithFields(map[string]interface{}{
							"volumeName": v.Name,
							"volumePath": vp,
						}).Info(ctx, "vol-ex client map already defined")
						break
					}
					volToExpMap[v] = e
				}
			}
			if !includeRootClients || e.RootClients == nil {
				continue
			}
			for _, p := range *e.RootClients {
				if vp == p {
					if _, ok := volToExpMap[v]; ok {
						log.WithFields(map[string]interface{}{
							"volumeName": v.Name,
							"volumePath": vp,
						}).Info(ctx, "vol-ex root client map already defined")
						break
					}
					volToExpMap[v] = e
				}
			}
		}
	}

	return volToExpMap, nil
}

// GetVolumeSize returns the total size of a snapshot folder
func (c *Client) GetVolumeSize(ctx context.Context,
	isiPath, name string,
) (int64, error) {
	folder, err := apiv1.GetIsiVolumeWithSize(ctx, c.API, isiPath, name)
	if err != nil {
		return 0, err
	}
	var totalSize int64
	totalSize = 0
	for _, attr := range folder.AttributeMap {
		totalSize += attr.Size
	}
	return totalSize, nil
}

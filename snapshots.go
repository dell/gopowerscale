/*
 Copyright (c) 2019 Dell Inc, or its subsidiaries.

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
	api "github.com/dell/goisilon/api/v1"
	"path"
	"strconv"
	"strings"
)

const namespacePath = "namespace"

// SnapshotList represents a list of Isilon snapshots.
type SnapshotList []*api.IsiSnapshot

// Snapshot represents an Isilon snapshot.
type Snapshot *api.IsiSnapshot

// GetSnapshots returns a list of snapshots from the cluster.
func (c *Client) GetSnapshots(ctx context.Context) (SnapshotList, error) {
	snapshots, err := api.GetIsiSnapshots(ctx, c.API)
	if err != nil {
		return nil, err
	}

	return snapshots.SnapshotList, nil
}

// GetSnapshotsByPath returns a list of snapshots covering the supplied path.
func (c *Client) GetSnapshotsByPath(
	ctx context.Context, path string) (SnapshotList, error) {

	snapshots, err := api.GetIsiSnapshots(ctx, c.API)
	if err != nil {
		return nil, err
	}
	// find all the snapshots with the same path
	snapshotsWithPath := make(SnapshotList, 0, len(snapshots.SnapshotList))
	for _, snapshot := range snapshots.SnapshotList {
		if snapshot.Path == c.API.VolumePath(path) {
			snapshotsWithPath = append(snapshotsWithPath, snapshot)
		}
	}
	return snapshotsWithPath, nil
}

// GetSnapshot returns a snapshot matching id, or if that is not found, matching name
func (c *Client) GetSnapshot(
	ctx context.Context, id int64, name string) (Snapshot, error) {

	// if we have an id, use it to find the snapshot
	snapshot, err := api.GetIsiSnapshot(ctx, c.API, id)
	if err == nil {
		return snapshot, nil
	}

	// there's no id or it didn't match, iterate through all snapshots and match
	// based on name
	if name == "" {
		return nil, err
	}
	snapshotList, err := c.GetSnapshots(ctx)
	if err != nil {
		return nil, err
	}

	for _, snapshot = range snapshotList {
		if snapshot.Name == name {
			return snapshot, nil
		}
	}

	return nil, nil
}

// CreateSnapshot creates a snapshot called name of the given path.
func (c *Client) CreateSnapshot(
	ctx context.Context, volName, snapshotName string) (Snapshot, error) {
	return api.CreateIsiSnapshot(ctx, c.API, c.API.VolumePath(volName), snapshotName)
}

// CreateSnapshotWithPath creates a snapshot by snapshot name and the path of volume.
func (c *Client) CreateSnapshotWithPath(
	ctx context.Context, path, snapshotName string) (Snapshot, error) {
	return api.CreateIsiSnapshot(ctx, c.API, path, snapshotName)
}

// RemoveSnapshot removes the snapshot by id, or failing that, the snapshot matching name.
func (c *Client) RemoveSnapshot(
	ctx context.Context, id int64, name string) error {

	snapshot, err := c.GetSnapshot(ctx, id, name)
	if err != nil {
		return err
	}

	return api.RemoveIsiSnapshot(ctx, c.API, snapshot.Id)
}

// CopySnapshot copies all files/directories in a snapshot to a new directory.
func (c *Client) CopySnapshot(
	ctx context.Context,
	sourceID int64, sourceName, destinationName string) (Volume, error) {

	snapshot, err := c.GetSnapshot(ctx, sourceID, sourceName)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, fmt.Errorf("Snapshot doesn't exist: (%d, %s)", sourceID, sourceName)
	}

	_, err = api.CopyIsiSnapshot(
		ctx, c.API, snapshot.Name,
		path.Base(snapshot.Path), destinationName)
	if err != nil {
		return nil, err
	}

	return c.GetVolume(ctx, destinationName, destinationName)
}

// CopySnapshotWithIsiPath copies all files/directories in a snapshot with isiPath to a new directory.
func (c *Client) CopySnapshotWithIsiPath(
	ctx context.Context,
	isiPath, snapshotSourceVolumeIsiPath string,
	sourceID int64,
	sourceName, destinationName string) (Volume, error) {

	snapshot, err := c.GetIsiSnapshotByIdentity(ctx, strconv.FormatInt(sourceID, 10))
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, fmt.Errorf("Snapshot doesn't exist: (%d, %s)", sourceID, sourceName)
	}

	_, err = api.CopyIsiSnapshotWithIsiPath(
		ctx, c.API, isiPath, snapshotSourceVolumeIsiPath, snapshot.Name,
		path.Base(snapshot.Path), destinationName)
	if err != nil {
		return nil, err
	}

	return c.GetVolumeWithIsiPath(ctx, isiPath, destinationName, destinationName)
}

// GetIsiSnapshotByIdentity query a snapshot by ID or name
// param identity string: name or id
func (c *Client) GetIsiSnapshotByIdentity(
	ctx context.Context, identity string) (Snapshot, error) {

	return api.GetIsiSnapshotByIdentity(ctx, c.API, identity)
}

// IsSnapshotExistent checks if a snapshot already exists
// param identity string: name or id
func (c *Client) IsSnapshotExistent(
	ctx context.Context, identity string) bool {

	snapshot, _ := api.GetIsiSnapshotByIdentity(ctx, c.API, identity)
	return snapshot != nil
}

// GetSnapshotFolderSize returns the total size of a snapshot folder
func (c *Client) GetSnapshotFolderSize(ctx context.Context,
	isiPath, name string) (int64, error) {

	snapshot, err := c.GetIsiSnapshotByIdentity(ctx, name)
	if err != nil {
		return 0, err
	}
	if snapshot == nil {
		return 0, fmt.Errorf("Snapshot doesn't exist: '%s'", name)
	}

	folder, err := api.GetIsiSnapshotFolderWithSize(ctx, c.API, isiPath, name, path.Base(snapshot.Path))
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

// GetSnapshotIsiPath returns the snapshot directory path
func (c *Client) GetSnapshotIsiPath(
	ctx context.Context,
	isiPath, snapshotId string) (string, error) {

	snapshot, err := c.GetIsiSnapshotByIdentity(ctx, snapshotId)
	if err != nil {
		return "", err
	}
	if snapshot == nil {
		return "", fmt.Errorf("Snapshot doesn't exist for snapshot id: (%s)", snapshotId)
	}

	snapshotPath := api.GetRealVolumeSnapshotPathWithIsiPath(isiPath, snapshot.Name)
	snapshotPath = path.Join(snapshotPath, path.Base(snapshot.Path))

	parts := strings.SplitN(snapshotPath, namespacePath, 2)
	if len(parts) < 2 {
		return "", fmt.Errorf("Snapshot doesn't exist for snapshot id: (%s)", snapshotId)
	}

	return parts[1], nil
}

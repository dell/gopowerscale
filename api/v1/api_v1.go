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
package v1

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dell/goisilon/api"
)

const (
	namespacePath     = "namespace"
	exportsPath       = "platform/1/protocols/nfs/exports"
	quotaPath         = "platform/1/quota/quotas"
	snapshotsPath     = "platform/1/snapshot/snapshots"
	zonesPath         = "platform/1/zones"
	snapshotParentDir = ".snapshot"
	userPath          = "platform/1/auth/users"
	rolePath          = "platform/1/auth/roles"
	roleMemberPath    = "platform/1/auth/roles/%s/members"
	groupPath         = "platform/1/auth/groups"
	groupMemberPath   = "platform/1/auth/groups/%s/members"
)

var debug, _ = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))

func realNamespacePath(client api.Client) string {
	return path.Join(namespacePath, client.VolumesPath())
}

func realexportsPath(client api.Client) string {
	return path.Join(exportsPath, client.VolumesPath())
}

func realVolumeSnapshotPath(client api.Client, name, zonePath, _ string) string {
	// Isi path is different from zone path
	volumeSnapshotPath := strings.Join([]string{zonePath, snapshotParentDir}, "/")
	if strings.Compare(zonePath, client.VolumesPath()) != 0 {
		parts := strings.SplitN(realNamespacePath(client), "/ifs", 2)
		return path.Join(parts[0], volumeSnapshotPath, name, parts[1])
	}
	return path.Join(volumeSnapshotPath, name)
}

// GetAbsoluteSnapshotPath get the absolute path of a snapshot
func GetAbsoluteSnapshotPath(c api.Client, snapshotName, volumeName, zonePath string) string {
	volumeSnapshotPath := strings.Join([]string{zonePath, snapshotParentDir}, "/")
	absoluteVolumePath := c.VolumePath(volumeName)
	return path.Join(volumeSnapshotPath, snapshotName, strings.TrimLeft(absoluteVolumePath, "/ifs/"))
}

// GetRealNamespacePathWithIsiPath gets the real namespace path by the combination of namespace and isiPath
func GetRealNamespacePathWithIsiPath(isiPath string) string {
	return path.Join(namespacePath, isiPath)
}

// GetRealVolumeSnapshotPathWithIsiPath gets the real volume snapshot path by using
// the isiPath in the parameter rather than use the default one in the client object
func GetRealVolumeSnapshotPathWithIsiPath(isiPath, zonePath, name, accessZone string) string {
	volumeSnapshotPath := strings.Join([]string{zonePath, snapshotParentDir}, "/")
	if accessZone == "System" {
		parts := strings.SplitN(GetRealNamespacePathWithIsiPath(isiPath), "/ifs", 2)
		if len(parts) == 2 {
			return path.Join(parts[0], volumeSnapshotPath, name, parts[1])
		}
	}
	// if Isi path is different then zone path get remaining isiPath
	_, remainIsiPath, found := strings.Cut(isiPath, zonePath)
	if found {
		return path.Join(namespacePath, zonePath, snapshotParentDir, name, remainIsiPath)
	}
	return path.Join(namespacePath, zonePath, snapshotParentDir, name)
}

// getAuthMemberID returns actual auth id, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel',
// memberType can be user/group.
func getAuthMemberID(memberType string, memberName *string, memberID *int32) (authMemberID string, err error) {
	memberType = strings.ToLower(memberType)
	if memberType != fileGroupTypeUser && memberType != fileGroupTypeGroup {
		return "", fmt.Errorf("member type is wrong, only support %s and %s", fileGroupTypeUser, fileGroupTypeGroup)
	}

	if memberName != nil && *memberName != "" {
		authMemberID = fmt.Sprintf("%s:%s", strings.ToUpper(memberType), *memberName)
	}

	if memberID != nil {
		authMemberID = fmt.Sprintf("%sID:%d", strings.ToUpper(memberType)[0:1], *memberID)
	}
	return
}

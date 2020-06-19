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
package v1

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dell/goisilon/api"
)

const (
	namespacePath       = "namespace"
	exportsPath         = "platform/1/protocols/nfs/exports"
	quotaPath           = "platform/1/quota/quotas"
	snapshotsPath       = "platform/1/snapshot/snapshots"
	volumesnapshotsPath = "/ifs/.snapshot"
	zonesPath           = "platform/1/zones"
)

var (
	debug, _ = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
)

func realNamespacePath(client api.Client) string {
	return path.Join(namespacePath, client.VolumesPath())
}

func realexportsPath(client api.Client) string {
	return path.Join(exportsPath, client.VolumesPath())
}

func realVolumeSnapshotPath(client api.Client, name string) string {
	parts := strings.SplitN(realNamespacePath(client), "/ifs", 2)
	return path.Join(parts[0], volumesnapshotsPath, name, parts[1])
}

// GetAbsoluteSnapshotPath get the absolute path of a snapshot
func GetAbsoluteSnapshotPath(c api.Client, snapshotName, volumeName string) string {
	absoluteVolumePath := c.VolumePath(volumeName)
	return path.Join(volumesnapshotsPath, snapshotName, strings.TrimLeft(absoluteVolumePath, "/ifs/"))
}

// GetRealNamespacePathWithIsiPath gets the real namespace path by the combination of namespace and isiPath
func GetRealNamespacePathWithIsiPath(isiPath string) string {
	return path.Join(namespacePath, isiPath)
}

// GetRealVolumeSnapshotPathWithIsiPath gets the real volume snapshot path by using
// the isiPath in the parameter rather than use the default one in the client object
func GetRealVolumeSnapshotPathWithIsiPath(isiPath string, name string) string {
	parts := strings.SplitN(GetRealNamespacePathWithIsiPath(isiPath), "/ifs", 2)
	return path.Join(parts[0], volumesnapshotsPath, name, parts[1])
}

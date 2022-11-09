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
package v2

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/dell/goisilon/api"
)

const (
	namespacePath       = "namespace"
	exportsPath         = "platform/2/protocols/nfs/exports"
	quotaPath           = "platform/2/quota/quotas"
	snapshotsPath       = "platform/2/snapshot/snapshots"
	volumeSnapshotsPath = "/ifs/.snapshot"
)

var (
	debug, _   = strconv.ParseBool(os.Getenv("GOISILON_DEBUG"))
	colonBytes = []byte{byte(':')}
)

func realNamespacePath(c api.Client) string {
	return path.Join(namespacePath, c.VolumesPath())
}

func realExportsPath(c api.Client) string {
	return path.Join(exportsPath, c.VolumesPath())
}

func realVolumeSnapshotPath(c api.Client, name string) string {
	parts := strings.SplitN(realNamespacePath(c), "/ifs/", 2)
	return path.Join(parts[0], volumeSnapshotsPath, name, parts[1])
}

// GetAbsoluteSnapshotPath get the absolute path of a snapshot
func GetAbsoluteSnapshotPath(c api.Client, snapshotName, volumeName string) string {
	absoluteVolumePath := c.VolumePath(volumeName)
	return path.Join(volumeSnapshotsPath, snapshotName, strings.TrimLeft(absoluteVolumePath, "/ifs/"))
}

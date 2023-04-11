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

type IsiVolume struct {
	Name         string `json:"name"`
	AttributeMap []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"attrs"`
}

// Isi PAPI volume JSON structs
type VolumeName struct {
	Name string `json:"name"`
}

type getIsiVolumesResp struct {
	Children []*VolumeName `json:"children"`
}

// Isi PAPI Volume ACL JSON structs
type Ownership struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type AclRequest struct {
	Authoritative string     `json:"authoritative"`
	Action        string     `json:"action"`
	Owner         *Ownership `json:"owner"`
	Group         *Ownership `json:"group,omitempty"`
}

// Isi PAPI volume attributes JSON struct
type getIsiVolumeAttributesResp struct {
	AttributeMap []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"attrs"`
}

// Isi PAPI volume size JSON struct
type getIsiVolumeSizeResp struct {
	AttributeMap []struct {
		Name string `json:"name"`
		Size int64  `json:"size"`
	} `json:"children"`
}

// Isi PAPI volume JSON struct
type getIsiVolumeResp struct {
	AttributeMap []struct {
		Name string `json:"name"`
	} `json:"children"`
}

// Isi PAPI export path JSON struct
type ExportPathList struct {
	Paths  []string `json:"paths"`
	MapAll struct {
		User   string   `json:"user"`
		Groups []string `json:"groups,omitempty"`
	} `json:"map_all"`
}

// Isi PAPI export clients JSON struct
type ExportClientList struct {
	Clients []string `json:"clients"`
}

// Isi PAPI export Id JSON struct
type postIsiExportResp struct {
	Id int `json:"id"`
}

// Isi PAPI export attributes JSON structs
type IsiExport struct {
	Id      int      `json:"id"`
	Paths   []string `json:"paths"`
	Clients []string `json:"clients"`
}

type getIsiExportsResp struct {
	ExportList []*IsiExport `json:"exports"`
}

// Isi PAPI snapshot path JSON struct
type SnapshotPath struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"`
}

// Isi PAPI snapshot JSON struct
type IsiSnapshot struct {
	Created       int64   `json:"created"`
	Expires       int64   `json:"expires"`
	HasLocks      bool    `json:"has_locks"`
	Id            int64   `json:"id"`
	Name          string  `json:"name"`
	Path          string  `json:"path"`
	PctFilesystem float64 `json:"pct_filesystem"`
	PctReserve    float64 `json:"pct_reserve"`
	Schedule      string  `json:"schedule"`
	ShadowBytes   int64   `json:"shadow_bytes"`
	Size          int64   `json:"size"`
	State         string  `json:"state"`
	TargetId      int64   `json:"target_it"`
	TargetName    string  `json:"target_name"`
}

type getIsiSnapshotsResp struct {
	SnapshotList []*IsiSnapshot `json:"snapshots"`
	Total        int64          `json:"total"`
	Resume       string         `json:"resume"`
}

type isiThresholds struct {
	Advisory             int64       `json:"advisory"`
	AdvisoryExceeded     bool        `json:"advisory_exceeded"`
	AdvisoryLastExceeded interface{} `json:"advisory_last_exceeded"`
	Hard                 int64       `json:"hard"`
	HardExceeded         bool        `json:"hard_exceeded"`
	HardLastExceeded     interface{} `json:"hard_last_exceeded"`
	Soft                 int64       `json:"soft"`
	SoftExceeded         bool        `json:"soft_exceeded"`
	SoftLastExceeded     interface{} `json:"soft_last_exceeded"`
	SoftGrace            int64       `json:"soft_grace"`
}

type IsiQuota struct {
	Container                 bool          `json:"container,omitempty"`
	Enforced                  bool          `json:"enforced,omitempty"`
	Id                        string        `json:"id"`
	IncludeSnapshots          bool          `json:"include_snapshots,omitempty"`
	Linked                    interface{}   `json:"linked,omitempty"`
	Notifications             string        `json:"notifications,omitempty"`
	Path                      string        `json:"path,omitempty"`
	Persona                   interface{}   `json:"persona,omitempty"`
	Ready                     bool          `json:"ready,omitempty"`
	Thresholds                isiThresholds `json:"thresholds,omitempty"`
	ThresholdsIncludeOverhead bool          `json:"thresholds_include_overhead,omitempty"`
	Type                      string        `json:"type,omitempty"`
	Usage                     struct {
		Inodes   int64 `json:"inodes"`
		Logical  int64 `json:"logical"`
		Physical int64 `json:"physical"`
	} `json:"usage"`
}

type isiThresholdsReq struct {
	Advisory  interface{} `json:"advisory,omitempty"`
	Hard      interface{} `json:"hard"`
	Soft      interface{} `json:"soft,omitempty"`
	SoftGrace interface{} `json:"soft_grace,omitempty"`
}

type IsiQuotaReq struct {
	Enforced                  bool             `json:"enforced"`
	IncludeSnapshots          bool             `json:"include_snapshots"`
	Path                      string           `json:"path"`
	Thresholds                isiThresholdsReq `json:"thresholds"`
	ThresholdsIncludeOverhead bool             `json:"thresholds_include_overhead"`
	Type                      string           `json:"type"`
	Container                 bool             `json:"container"`
}

type IsiUpdateQuotaReq struct {
	Enforced                  bool             `json:"enforced"`
	Thresholds                isiThresholdsReq `json:"thresholds"`
	ThresholdsIncludeOverhead bool             `json:"thresholds_include_overhead"`
}

type isiQuotaListResp struct {
	Quotas []IsiQuota `json:"quotas"`
}

type IsiQuotaListRespResume struct {
	Quotas []*IsiQuota `json:"quotas,omitempty"`
	Resume string      `json:"resume,omitempty"`
}

// getIsiZonesResp returns an array of all related access zones
type getIsiZonesResp struct {
	Zones []*IsiZone `json:"zones"`
}

// IsiZone contains information of an access zone
type IsiZone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
}

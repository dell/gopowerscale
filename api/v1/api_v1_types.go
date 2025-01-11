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

const (
	fileGroupTypeUser      = "user"
	fileGroupTypeGroup     = "group"
	fileGroupTypeWellKnown = "wellknown"
)

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

type GetIsiVolumesResp struct {
	Children []*VolumeName `json:"children"`
}

// Isi PAPI Volume ACL JSON structs
type Ownership struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ACLRequest struct {
	Authoritative string     `json:"authoritative"`
	Action        string     `json:"action"`
	Owner         *Ownership `json:"owner"`
	Group         *Ownership `json:"group,omitempty"`
}

// Isi PAPI volume attributes JSON struct
type GetIsiVolumeAttributesResp struct {
	AttributeMap []struct {
		Name  string      `json:"name"`
		Value interface{} `json:"value"`
	} `json:"attrs"`
}

// Isi PAPI volume size JSON struct
type GetIsiVolumeSizeResp struct {
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
	ID int `json:"id"`
}

// Isi PAPI export attributes JSON structs
type IsiExport struct {
	ID      int      `json:"id"`
	Paths   []string `json:"paths"`
	Clients []string `json:"clients"`
}

type GetIsiExportsResp struct {
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
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Path          string  `json:"path"`
	PctFilesystem float64 `json:"pct_filesystem"`
	PctReserve    float64 `json:"pct_reserve"`
	Schedule      string  `json:"schedule"`
	ShadowBytes   int64   `json:"shadow_bytes"`
	Size          int64   `json:"size"`
	State         string  `json:"state"`
	TargetID      int64   `json:"target_it"`
	TargetName    string  `json:"target_name"`
}

type GetIsiSnapshotsResp struct {
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
	ID                        string        `json:"id"`
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

type IsiCopySnapshotResp struct {
	Errors []struct {
		ErrorErc string `json:"error_src"`
		Message  string `json:"message"`
		Source   string `json:"source"`
		Target   string `json:"target"`
	} `json:"copy_errors"`
	Success bool `json:"success"`
}

// IsiAccessItemFileGroup Specifies the persona of the file group.
type IsiAccessItemFileGroup struct {
	// Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
	ID string `json:"id,omitempty"`
	// Specifies the persona name, which must be combined with a type.
	Name string `json:"name,omitempty"`
	// Specifies the type of persona, which must be combined with a name. Values can be user, group or wellknown
	Type string `json:"type,omitempty"`
}

// IsiAuthMemberItem Specifies the persona of the group member. Member can be user or group.
type IsiAuthMemberItem struct {
	ID   *int32  `json:"id,omitempty"`
	Name *string `json:"name,omitempty"`
	Type string  `json:"type"`
}

type IsiUser struct {
	// Specifies the distinguished name for the user.
	Dn string `json:"dn"`
	// Specifies the DNS domain.
	DNSDomain string `json:"dns_domain"`
	// Specifies the domain that the object is part of.
	Domain string `json:"domain"`
	// Specifies an email address.
	Email string `json:"email"`
	// True, if the authenticated user is enabled.
	Enabled bool `json:"enabled"`
	// True, if the authenticated user has expired.
	Expired bool `json:"expired"`
	// Specifies the Unix Epoch time at which the authenticated user will expire.
	Expiry int32 `json:"expiry"`
	// Specifies the GECOS value, which is usually the full name.
	Gecos string `json:"gecos"`
	// True, if the GID was generated.
	GeneratedGid bool `json:"generated_gid"`
	// True, if the UID was generated.
	GeneratedUID bool `json:"generated_uid"`
	// True, if the UPN was generated.
	GeneratedUpn bool                   `json:"generated_upn"`
	Gid          IsiAccessItemFileGroup `json:"gid"`
	// Specifies a home directory for the user.
	HomeDirectory string `json:"home_directory"`
	// Specifies the user or group ID.
	ID string `json:"id"`
	// If true, indicates that the account is locked.
	Locked bool `json:"locked"`
	// Specifies the maximum time in seconds allowed before the password expires.
	MaxPasswordAge int32 `json:"max_password_age"`
	// Specifies the groups that this user or group are members of.
	MemberOf            []IsiAccessItemFileGroup `json:"member_of"`
	Name                string                   `json:"name"`
	OnDiskGroupIdentity IsiAccessItemFileGroup   `json:"on_disk_group_identity"`
	OnDiskUserIdentity  IsiAccessItemFileGroup   `json:"on_disk_user_identity"`
	// If true, the password has expired.
	PasswordExpired bool `json:"password_expired"`
	// If true, the password is allowed to expire.
	PasswordExpires bool `json:"password_expires"`
	// Specifies the time in Unix Epoch seconds that the password will expire.
	PasswordExpiry int32 `json:"password_expiry"`
	// Specifies the last time the password was set.
	PasswordLastSet int32                  `json:"password_last_set"`
	PrimaryGroupSid IsiAccessItemFileGroup `json:"primary_group_sid"`
	// Prompts the user to change their password at the next login.
	PromptPasswordChange bool `json:"prompt_password_change"`
	// Specifies the authentication provider that the object belongs to.
	Provider string `json:"provider"`
	// Specifies a user or group name.
	SamAccountName string `json:"sam_account_name"`
	// Specifies a path to the shell for the user.
	Shell string                 `json:"shell"`
	Sid   IsiAccessItemFileGroup `json:"sid"`
	// Specifies the object type.
	Type string                 `json:"type"`
	UID  IsiAccessItemFileGroup `json:"uid"`
	// Specifies a principal name for the user.
	Upn string `json:"upn"`
	// Specifies whether the password for the user can be changed.
	UserCanChangePassword bool `json:"user_can_change_password"`
}

type IsiUserReq struct {
	// Specifies an email address for the user.
	Email *string `json:"email,omitempty"`
	// If true, the authenticated user is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Specifies the Unix Epoch time when the auth user will expire.
	Expiry *int32 `json:"expiry,omitempty"`
	// Specifies the GECOS value, which is usually the full name.
	Gecos *string `json:"gecos,omitempty"`
	// Specifies a home directory for the user.
	HomeDirectory *string `json:"home_directory,omitempty"`
	// Specifies a user name.
	Name string `json:"name"`
	// Changes the password for the user.
	Password *string `json:"password,omitempty"`
	// If true, the password should expire.
	PasswordExpires *bool                   `json:"password_expires,omitempty"`
	PrimaryGroup    *IsiAccessItemFileGroup `json:"primary_group,omitempty"`
	// If true, prompts the user to change their password at the next login.
	PromptPasswordChange *bool `json:"prompt_password_change,omitempty"`
	// Specifies the shell for the user.
	Shell *string `json:"shell,omitempty"`
	// Specifies a numeric user identifier.
	UID *int32 `json:"uid,omitempty"`
	// If true, the user account should be unlocked.
	Unlock *bool `json:"unlock,omitempty"`
}

type IsiUpdateUserReq struct {
	// Specifies an email address for the user.
	Email *string `json:"email,omitempty"`
	// If true, the authenticated user is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Specifies the Unix Epoch time when the auth user will expire.
	Expiry *int32 `json:"expiry,omitempty"`
	// Specifies the GECOS value, which is usually the full name.
	Gecos *string `json:"gecos,omitempty"`
	// Specifies a home directory for the user.
	HomeDirectory *string `json:"home_directory,omitempty"`
	// Changes the password for the user.
	Password *string `json:"password,omitempty"`
	// If true, the password should expire.
	PasswordExpires *bool                   `json:"password_expires,omitempty"`
	PrimaryGroup    *IsiAccessItemFileGroup `json:"primary_group,omitempty"`
	// If true, prompts the user to change their password at the next login.
	PromptPasswordChange *bool `json:"prompt_password_change,omitempty"`
	// Specifies the shell for the user.
	Shell *string `json:"shell,omitempty"`
	// Specifies a numeric user identifier.
	UID *int32 `json:"uid,omitempty"`
	// If true, the user account should be unlocked.
	Unlock *bool `json:"unlock,omitempty"`
}

type IsiUserListResp struct {
	Users []*IsiUser `json:"users,omitempty"`
}

type IsiUserListRespResume struct {
	Resume string     `json:"resume,omitempty"`
	Users  []*IsiUser `json:"users,omitempty"`
}

// IsiRolePrivilegeItem Specifies the system-defined privilege that may be granted to users.
type IsiRolePrivilegeItem struct {
	// Specifies the ID of the privilege.
	ID string `json:"id"`
	// Specifies the name of the privilege.
	Name *string `json:"name,omitempty"`
	// True, if the privilege is read-only.
	ReadOnly *bool `json:"read_only,omitempty"`
}

type IsiRole struct {
	// Specifies the description of the role.
	Description string `json:"description"`
	// Specifies the users or groups that have this role.
	Members []IsiAccessItemFileGroup `json:"members"`
	// Specifies the name of the role.
	Name string `json:"name"`
	// Specifies the privileges granted by this role.
	Privileges []IsiRolePrivilegeItem `json:"privileges"`
	// Specifies the ID of the role.
	ID string `json:"id"`
}

type isiRoleListResp struct {
	Roles []*IsiRole `json:"roles,omitempty"`
}

type IsiRoleListRespResume struct {
	Resume string     `json:"resume,omitempty"`
	Roles  []*IsiRole `json:"roles,omitempty"`
	Total  *int32     `json:"total,omitempty"`
}

// IsiGroup Specifies configuration properties for a group.
type IsiGroup struct {
	// Specifies the distinguished name for the user.
	Dn string `json:"dn"`
	// Specifies the DNS domain.
	DNSDomain string `json:"dns_domain"`
	// Specifies the domain that the object is part of.
	Domain string `json:"domain"`
	// If true, the GID was generated.
	GeneratedGid bool                   `json:"generated_gid"`
	Gid          IsiAccessItemFileGroup `json:"gid"`
	// Specifies the user or group ID.
	ID string `json:"id"`
	// Specifies the groups that this user or group are members of.
	MemberOf []IsiAccessItemFileGroup `json:"member_of"`
	// Specifies a user or group name.
	Name string `json:"name"`
	// ObjectHistory []V1AuthGroupObjectHistoryItem `json:"object_history,omitempty"`
	// Specifies the authentication provider that the object belongs to.
	Provider string `json:"provider"`
	// Specifies a user or group name.
	SamAccountName string                 `json:"sam_account_name"`
	Sid            IsiAccessItemFileGroup `json:"sid"`
	// Specifies the object type.
	Type string `json:"type"`
}

// IsiGroupReq Specifies the configuration properties for a group.
type IsiGroupReq struct {
	// Specifies the numeric group identifier.
	Gid *int32 `json:"gid,omitempty"`
	// Specifies the members of the group.
	Members []IsiAccessItemFileGroup `json:"members,omitempty"`
	// Specifies the group name.
	Name string `json:"name"`
}

// IsiUpdateGroupReq Specifies the configuration properties for a group.
type IsiUpdateGroupReq struct {
	// Specifies the numeric group identifier.
	Gid int32 `json:"gid,omitempty"`
}

type IsiGroupMemberListRespResume struct {
	Members []*IsiAccessItemFileGroup `json:"members,omitempty"`
	Resume  string                    `json:"resume,omitempty"`
}

type IsiGroupListResp struct {
	Groups []*IsiGroup `json:"groups,omitempty"`
}

type IsiGroupListRespResume struct {
	// Provide this token as the 'resume' query argument to continue listing results.
	Groups []*IsiGroup `json:"groups,omitempty"`
	// Provide this token as the 'resume' query argument to continue listing results.
	Resume string `json:"resume,omitempty"`
}

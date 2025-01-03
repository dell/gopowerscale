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

package openapi

// V12SmbShare struct for V12SmbShare
type V12SmbShare struct {
	// Only enumerate files and folders the requesting user has access to.
	AccessBasedEnumeration *bool `json:"access_based_enumeration,omitempty"`
	// Access-based enumeration on only the root directory of the share.
	AccessBasedEnumerationRootOnly *bool `json:"access_based_enumeration_root_only,omitempty"`
	// Allow deletion of read-only files in the share.
	AllowDeleteReadonly *bool `json:"allow_delete_readonly,omitempty"`
	// Allows users to execute files they have read rights for.
	AllowExecuteAlways *bool `json:"allow_execute_always,omitempty"`
	// Allow automatic expansion of variables for home directories.
	AllowVariableExpansion *bool `json:"allow_variable_expansion,omitempty"`
	// Automatically create home directories.
	AutoCreateDirectory *bool `json:"auto_create_directory,omitempty"`
	// Share is visible in net view and the browse list.
	Browsable *bool `json:"browsable,omitempty"`
	// Persistent open timeout for the share.
	CaTimeout *int32 `json:"ca_timeout,omitempty"`
	// Specify the level of write-integrity on continuously available shares.
	CaWriteIntegrity *string `json:"ca_write_integrity,omitempty"`
	// Level of change notification alerts on the share.
	ChangeNotify *string `json:"change_notify,omitempty"`
	// Specify if persistent opens are allowed on the share.
	ContinuouslyAvailable *bool `json:"continuously_available,omitempty"`
	// Create path if does not exist.
	CreatePath *bool `json:"create_path,omitempty"`
	// Create permissions for new files and directories in share.
	CreatePermissions *string `json:"create_permissions,omitempty"`
	// Client-side caching policy for the shares.
	CscPolicy *string `json:"csc_policy,omitempty"`
	// Description for this SMB share.
	Description *string `json:"description,omitempty"`
	// Directory create mask bits.
	DirectoryCreateMask *int32 `json:"directory_create_mask,omitempty"`
	// Directory create mode bits.
	DirectoryCreateMode *int32 `json:"directory_create_mode,omitempty"`
	// File create mask bits.
	FileCreateMask *int32 `json:"file_create_mask,omitempty"`
	// File create mode bits.
	FileCreateMode *int32 `json:"file_create_mode,omitempty"`
	// Specifies the list of file extensions.
	FileFilterExtensions []string `json:"file_filter_extensions,omitempty"`
	// Specifies if filter list is for deny or allow. Default is deny.
	FileFilterType *string `json:"file_filter_type,omitempty"`
	// Enables file filtering on this zone.
	FileFilteringEnabled *bool `json:"file_filtering_enabled,omitempty"`
	// Hide files and directories that begin with a period '.'.
	HideDotFiles *bool `json:"hide_dot_files,omitempty"`
	// An ACL expressing which hosts are allowed access. A deny clause must be the final entry.
	HostACL []string `json:"host_acl,omitempty"`
	// Specify the condition in which user access is done as the guest account.
	ImpersonateGuest *string `json:"impersonate_guest,omitempty"`
	// User account to be used as guest account.
	ImpersonateUser *string `json:"impersonate_user,omitempty"`
	// Set the inheritable ACL on the share path.
	InheritablePathACL *bool `json:"inheritable_path_acl,omitempty"`
	// Specifies the wchar_t starting point for automatic byte mangling.
	MangleByteStart *int32 `json:"mangle_byte_start,omitempty"`
	// Character mangle map.
	MangleMap []string `json:"mangle_map,omitempty"`
	// Share name.
	Name string `json:"name"`
	// Support NTFS ACLs on files and directories.
	NtfsACLSupport *bool `json:"ntfs_acl_support,omitempty"`
	// Support oplocks.
	Oplocks *bool `json:"oplocks,omitempty"`
	// Path of share within /ifs.
	Path string `json:"path"`
	// Specifies an ordered list of permission modifications.
	Permissions []V1SmbSharePermission `json:"permissions,omitempty"`
	// Allow account to run as root.
	RunAsRoot []V1AuthAccessAccessItemFileGroup `json:"run_as_root,omitempty"`
	// Enables SMB3 encryption for the share.
	Smb3EncryptionEnabled *bool `json:"smb3_encryption_enabled,omitempty"`
	// Enables sparse file.
	SparseFile *bool `json:"sparse_file,omitempty"`
	// Specifies if persistent opens would do strict lockout on the share.
	StrictCaLockout *bool `json:"strict_ca_lockout,omitempty"`
	// Handle SMB flush operations.
	StrictFlush *bool `json:"strict_flush,omitempty"`
	// Specifies whether byte range locks contend against SMB I/O.
	StrictLocking *bool `json:"strict_locking,omitempty"`
	// Name of the access zone to which to move this SMB share.
	Zone *string `json:"zone,omitempty"`
}

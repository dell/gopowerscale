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

// V2NfsExport Specifies configuration values for NFS exports.
type V2NfsExport struct {
	// True if all directories under the specified paths are mountable.
	AllDirs *bool `json:"all_dirs,omitempty"`
	// Specifies the block size returned by the NFS statfs procedure.
	BlockSize *int32 `json:"block_size,omitempty"`
	// True if the client can set file times through the NFS set attribute request. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CanSetTime *bool `json:"can_set_time,omitempty"`
	// True if the case is ignored for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CaseInsensitive *bool `json:"case_insensitive,omitempty"`
	// True if the case is preserved for file names. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	CasePreserving *bool `json:"case_preserving,omitempty"`
	// True if the superuser can change file ownership. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	ChownRestricted *bool `json:"chown_restricted,omitempty"`
	// Specifies the clients with root access to the export.
	Clients []string `json:"clients,omitempty"`
	// True if NFS  commit  requests execute asynchronously.
	CommitAsynchronous *bool `json:"commit_asynchronous,omitempty"`
	// Specifies the user-defined string that is used to identify the export.
	Description *string `json:"description,omitempty"`
	// Specifies the preferred size for directory read operations. This value is used to advise the client of optimal settings for the server, but is not enforced.
	DirectoryTransferSize *int32 `json:"directory_transfer_size,omitempty"`
	// Specifies the default character set encoding of the clients connecting to the export, unless otherwise specified.
	Encoding *string `json:"encoding,omitempty"`
	// Specifies the reported maximum number of links to a file. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	LinkMax    *int32             `json:"link_max,omitempty"`
	MapAll     *V2NfsExportMapAll `json:"map_all,omitempty"`
	MapFailure *V2NfsExportMapAll `json:"map_failure,omitempty"`
	// True if user mappings query the OneFS user database. When set to false, user mappings only query local authentication.
	MapFull *bool `json:"map_full,omitempty"`
	// True if incoming user IDs (UIDs) are mapped to users in the OneFS user database. When set to false, incoming UIDs are applied directly to file operations.
	MapLookupUID *bool              `json:"map_lookup_uid,omitempty"`
	MapNonRoot   *V2NfsExportMapAll `json:"map_non_root,omitempty"`
	// Determines whether searches for users specified in 'map_all', 'map_root' or 'map_nonroot' are retried if the search fails.
	MapRetry *bool              `json:"map_retry,omitempty"`
	MapRoot  *V2NfsExportMapAll `json:"map_root,omitempty"`
	// Specifies the maximum file size for any file accessed from the export. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	MaxFileSize *int64 `json:"max_file_size,omitempty"`
	// Specifies the reported maximum length of a file name. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NameMaxSize *int32 `json:"name_max_size,omitempty"`
	// True if long file names result in an error. This parameter does not affect server behavior, but is included to accommodate legacy client requirements.
	NoTruncate *bool `json:"no_truncate,omitempty"`
	// Specifies the paths under /ifs that are exported.
	Paths []string `json:"paths"`
	// True if the export is set to read-only.
	ReadOnly *bool `json:"read_only,omitempty"`
	// Specifies the clients with read-only access to the export.
	ReadOnlyClients []string `json:"read_only_clients,omitempty"`
	// Specifies the maximum buffer size that clients should use on NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMaxSize *int32 `json:"read_transfer_max_size,omitempty"`
	// Specifies the preferred multiple size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferMultiple *int32 `json:"read_transfer_multiple,omitempty"`
	// Specifies the preferred size for NFS read requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	ReadTransferSize *int32 `json:"read_transfer_size,omitempty"`
	// Specifies the clients with both read and write access to the export, even when the export is set to read-only.
	ReadWriteClients []string `json:"read_write_clients,omitempty"`
	// True if 'readdirplus' requests are enabled. Enabling this property might improve network performance and is only available for NFSv3.
	Readdirplus *bool `json:"readdirplus,omitempty"`
	// Sets the number of directory entries that are prefetched when a 'readdirplus' request is processed. (Deprecated.)
	ReaddirplusPrefetch *int32 `json:"readdirplus_prefetch,omitempty"`
	// Limits the size of file identifiers returned by NFSv3+ to 32-bit values (may require remount).
	Return32bitFileIDs *bool `json:"return_32bit_file_ids,omitempty"`
	// Clients that have root access to the export.
	RootClients []string `json:"root_clients,omitempty"`
	// Specifies the authentication types that are supported for this export.
	SecurityFlavors []string `json:"security_flavors,omitempty"`
	// True if set attribute operations execute asynchronously.
	SetattrAsynchronous *bool `json:"setattr_asynchronous,omitempty"`
	// Specifies the snapshot for all mounts.
	Snapshot *string `json:"snapshot,omitempty"`
	// True if symlinks are supported. This value is used to advise the client of optimal settings for the server, but is not enforced.
	Symlinks *bool `json:"symlinks,omitempty"`
	// Specifies the resolution of all time values that are returned to the clients
	TimeDelta *float32 `json:"time_delta,omitempty"`
	// Specifies the action to be taken when an NFSv3+ datasync write is requested.
	WriteDatasyncAction *string `json:"write_datasync_action,omitempty"`
	// Specifies the stability disposition returned when an NFSv3+ datasync write is processed.
	WriteDatasyncReply *string `json:"write_datasync_reply,omitempty"`
	// Specifies the action to be taken when an NFSv3+ filesync write is requested.
	WriteFilesyncAction *string `json:"write_filesync_action,omitempty"`
	// Specifies the stability disposition returned when an NFSv3+ filesync write is processed.
	WriteFilesyncReply *string `json:"write_filesync_reply,omitempty"`
	// Specifies the maximum buffer size that clients should use on NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMaxSize *int32 `json:"write_transfer_max_size,omitempty"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferMultiple *int32 `json:"write_transfer_multiple,omitempty"`
	// Specifies the preferred multiple size for NFS write requests. This value is used to advise the client of optimal settings for the server, but is not enforced.
	WriteTransferSize *int32 `json:"write_transfer_size,omitempty"`
	// Specifies the action to be taken when an NFSv3+ unstable write is requested.
	WriteUnstableAction *string `json:"write_unstable_action,omitempty"`
	// Specifies the stability disposition returned when an NFSv3+ unstable write is processed.
	WriteUnstableReply *string `json:"write_unstable_reply,omitempty"`
	// Specifies the zone in which the export is valid.
	Zone *string `json:"zone,omitempty"`
}

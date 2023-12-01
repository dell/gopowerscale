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

// V2NfsExportMapAll Specifies the users and groups to which non-root and root clients are mapped.
type V2NfsExportMapAll struct {
	// True if the user mapping is applied.
	Enabled      *bool                            `json:"enabled,omitempty"`
	PrimaryGroup *V1AuthAccessAccessItemFileGroup `json:"primary_group,omitempty"`
	// Specifies persona properties for the secondary user group. A persona consists of either a type and name, or an ID.
	SecondaryGroups []V2NfsExportMapAllSecondaryGroupsInner `json:"secondary_groups,omitempty"`
	User            *V1AuthAccessAccessItemFileGroup        `json:"user,omitempty"`
}

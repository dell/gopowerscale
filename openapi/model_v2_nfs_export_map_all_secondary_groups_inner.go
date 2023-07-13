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

// V2NfsExportMapAllSecondaryGroupsInner Specifies properties for a persona, which consists of either a 'type' and a 'name' or an 'ID'.
type V2NfsExportMapAllSecondaryGroupsInner struct {
	// Specifies the serialized form of a persona, which can be 'UID:0', 'USER:name', 'GID:0', 'GROUP:wheel', or 'SID:S-1-1'.
	Id *string `json:"id,omitempty"`
	// Specifies the persona name, which must be combined with a type.
	Name NullableString `json:"name,omitempty"`
	// Specifies the type of persona, which must be combined with a name.
	Type NullableString `json:"type,omitempty"`
}

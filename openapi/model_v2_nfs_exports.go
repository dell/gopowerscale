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

// V2NfsExports struct for V2NfsExports
type V2NfsExports struct {
	// An identifier for a set of exports.
	Digest *string `json:"digest,omitempty"`
	Exports []V2NfsExportExtended `json:"exports,omitempty"`
	// Provide this token as the 'resume' query argument to continue listing results.
	Resume *string `json:"resume,omitempty"`
	// Total number of items available.
	Total *int32 `json:"total,omitempty"`
}

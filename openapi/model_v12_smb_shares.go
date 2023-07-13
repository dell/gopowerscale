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

// V12SmbShares struct for V12SmbShares
type V12SmbShares struct {
	// An identifier for a set of shares.
	Digest *string `json:"digest,omitempty"`
	// Provide this token as the 'resume' query argument to continue listing results.
	Resume *string `json:"resume,omitempty"`
	Shares []V12SmbShareExtended `json:"shares"`
	// Total number of items available.
	Total *int32 `json:"total,omitempty"`
}

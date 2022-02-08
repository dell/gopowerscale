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
package goisilon

import "testing"

// Test if the zone returns correctly matched the name parsed in
func TestGetZoneByName(t *testing.T) {
	// Get local serial
	name := "csi0zone"
	zone, err := client.GetZoneByName(defaultCtx, name)
	if err != nil {
		panic(err)
	}
	if zone.Name != name {
		panic("Not match")
	}
	println("Test get zone by name complete, the path is: " + zone.Path)
}

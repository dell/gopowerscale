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

func TestGetStatistics(*testing.T) {
	keyArray := []string{"ifs.bytes.avail", "ifs.bytes.total"}
	stats, err := client.GetStatistics(defaultCtx, keyArray)
	if err != nil || len(stats.StatsList) != 2 {
		panic("Couldn't get statistics.")
	}
	if stats.StatsList[0].Value <= 0 {
		panic("Statistics returned bad value.")
	}
}

// Test if the local serial can be returned normally
func TestGetLocalSerial(t *testing.T) {
	// Get local serial
	localSerial, err := client.GetLocalSerial(defaultCtx)
	if err != nil {
		panic(err)
	}
	println(localSerial)

}

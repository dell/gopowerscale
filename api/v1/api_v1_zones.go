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

import (
	"context"

	"github.com/dell/goisilon/api"
)

// GetZoneByName returns a specific access zone which matches the name parsed in
func GetZoneByName(ctx context.Context,
	client api.Client,
	name string) (*IsiZone, error) {
	var resp getIsiZonesResp
	// PAPI call: GET https://1.2.3.4:8080/platform/1/zones/zone
	err := client.Get(ctx, zonesPath, name, nil, nil, &resp)
	if err != nil {
		return nil, err
	}
	if resp.Zones == nil {
		return nil, nil
	}
	return resp.Zones[0], nil
}

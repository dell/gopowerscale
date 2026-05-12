/*
Copyright (c) 2022-2025 Dell Inc, or its subsidiaries.

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
package gopowerscale

import (
	"context"
	"testing"

	apiv1 "github.com/dell/gopowerscale/api/v1"
	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test if the zone returns correctly matched the name parsed in
func TestGetZoneByName(t *testing.T) {
	ctx := context.Background()
	zoneName := "csi0zone"
	expectedZone := &apiv1.IsiZone{
		Name: zoneName,
		ID:   "test-id",
		Path: "/ifs/" + zoneName,
	}
	client.API.(*mocks.Client).On("GetZoneByName", mock.Anything, zoneName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: []*apiv1.IsiZone{expectedZone},
		}
	}).Once()
	zone, err := client.GetZoneByName(ctx, zoneName)
	assert.Nil(t, err)
	assert.Equal(t, expectedZone, zone)
}

// Test if GetIsiZoneList returns all access zones correctly
func TestGetIsiZoneList(t *testing.T) {
	ctx := context.Background()
	expectedZones := []*apiv1.IsiZone{
		{
			Name: "System",
			ID:   "system-id",
			Path: "/ifs",
		},
		{
			Name: "csi0zone",
			ID:   "zone-id",
			Path: "/ifs/csi0zone",
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.GetIsiZonesResp)
		*resp = apiv1.GetIsiZonesResp{
			Zones: expectedZones,
		}
	}).Once()
	zonesResp, err := client.GetIsiZoneList(ctx)
	assert.Nil(t, err)
	assert.NotNil(t, zonesResp)
	assert.Equal(t, expectedZones, zonesResp.Zones)
}

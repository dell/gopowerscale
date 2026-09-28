/*
Copyright (c) 2026 Dell Inc. or its subsidiaries. All Rights Reserved.

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

package v17

import (
	"context"
	"fmt"
	"testing"

	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestGetLicenseByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Test successful license retrieval
	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*LicensesResponse)
		resp.Licenses = []License{
			{
				ExpiredAlert:  false,
				ExpiringAlert: false,
				ID:            "SNAPSHOTIQ",
				Name:          "SNAPSHOTIQ",
				Status:        "Licensed",
				Tiers: []LicenseTier{
					{
						LicensedDriveCapacity: 24,
						LicensedNodeCount:     3,
						Platform:              "",
						Tier:                  "16",
						UsedDriveCapacity:     24,
						UsedNodeCount:         3,
					},
				},
			},
		}
	}).Once()

	license, err := GetLicenseByID(ctx, client, "SNAPSHOTIQ")
	assert.Nil(t, err)
	assert.NotNil(t, license)
	assert.Equal(t, "SNAPSHOTIQ", license.ID)
	assert.Equal(t, "Licensed", license.Status)
	assert.Len(t, license.Tiers, 1)
	assert.Equal(t, 3, license.Tiers[0].LicensedNodeCount)

	// Test empty license ID
	_, err = GetLicenseByID(ctx, client, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "license ID cannot be empty")

	// Test API error
	client.On("Get", anyArgs[0:6]...).Return(fmt.Errorf("API error")).Once()
	_, err = GetLicenseByID(ctx, client, "INVALID")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API error")

	// Test empty licenses response
	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*LicensesResponse)
		resp.Licenses = []License{}
	}).Once()
	_, err = GetLicenseByID(ctx, client, "NONEXISTENT")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "no license found")
}

func TestLicenseTypes(t *testing.T) {
	// Test LicenseTier struct
	tier := LicenseTier{
		LicensedDriveCapacity: 100,
		LicensedNodeCount:     5,
		Platform:              "PowerScale",
		Tier:                  "16",
		UsedDriveCapacity:     50,
		UsedNodeCount:         2,
	}
	assert.Equal(t, 100, tier.LicensedDriveCapacity)
	assert.Equal(t, 5, tier.LicensedNodeCount)
	assert.Equal(t, "PowerScale", tier.Platform)
	assert.Equal(t, "16", tier.Tier)
	assert.Equal(t, 50, tier.UsedDriveCapacity)
	assert.Equal(t, 2, tier.UsedNodeCount)

	// Test License struct
	license := License{
		ExpiredAlert:  false,
		ExpiringAlert: true,
		ID:            "TEST_LICENSE",
		Name:          "Test License",
		Status:        "Licensed",
		Tiers:         []LicenseTier{tier},
	}
	assert.False(t, license.ExpiredAlert)
	assert.True(t, license.ExpiringAlert)
	assert.Equal(t, "TEST_LICENSE", license.ID)
	assert.Equal(t, "Test License", license.Name)
	assert.Equal(t, "Licensed", license.Status)
	assert.Len(t, license.Tiers, 1)

	// Test LicensesResponse struct
	response := LicensesResponse{
		Licenses: []License{license},
	}
	assert.Len(t, response.Licenses, 1)
	assert.Equal(t, "TEST_LICENSE", response.Licenses[0].ID)
}

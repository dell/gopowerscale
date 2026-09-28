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

package gopowerscale

import (
	"fmt"
	"testing"

	apiV17 "github.com/dell/gopowerscale/api/v17"
	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetLicenseByID(t *testing.T) {
	// Test successful license retrieval
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiV17.LicensesResponse)
		resp.Licenses = []apiV17.License{
			{
				ExpiredAlert:  false,
				ExpiringAlert: false,
				ID:            "SNAPSHOTIQ",
				Name:          "SNAPSHOTIQ",
				Status:        "Licensed",
				Tiers: []apiV17.LicenseTier{
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

	license, err := client.GetLicenseByID(defaultCtx, "SNAPSHOTIQ")
	assert.Nil(t, err)
	assert.NotNil(t, license)
	assert.Equal(t, "SNAPSHOTIQ", license.ID)
	assert.Equal(t, "Licensed", license.Status)
	assert.Len(t, license.Tiers, 1)
	assert.Equal(t, 3, license.Tiers[0].LicensedNodeCount)

	// Test API error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("API error")).Once()
	_, err = client.GetLicenseByID(defaultCtx, "INVALID")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "API error")

	// Test empty license ID
	_, err = client.GetLicenseByID(defaultCtx, "")
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "license ID cannot be empty")
}

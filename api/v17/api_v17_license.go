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

	"github.com/dell/gopowerscale/api"
)

const (
	licensePath = "platform/17/license/licenses"
)

// LicenseTier represents tier information for a license
type LicenseTier struct {
	LicensedDriveCapacity int    `json:"licensed_drive_capacity,omitempty"`
	LicensedNodeCount     int    `json:"licensed_node_count,omitempty"`
	Platform              string `json:"platform,omitempty"`
	Tier                  string `json:"tier,omitempty"`
	UsedDriveCapacity     int    `json:"used_drive_capacity,omitempty"`
	UsedNodeCount         int    `json:"used_node_count,omitempty"`
}

// License represents a PowerScale license
type License struct {
	ExpiredAlert  bool          `json:"expired_alert,omitempty"`
	ExpiringAlert bool          `json:"expiring_alert,omitempty"`
	ID            string        `json:"id,omitempty"`
	Name          string        `json:"name,omitempty"`
	Status        string        `json:"status,omitempty"`
	Tiers         []LicenseTier `json:"tiers,omitempty"`
}

// LicensesResponse is the wrapper response for license queries
type LicensesResponse struct {
	Licenses []License `json:"licenses,omitempty"`
}

// GetLicenseByID retrieves license details by license ID
func GetLicenseByID(
	ctx context.Context,
	client api.Client,
	licenseID string,
) (*License, error) {
	if licenseID == "" {
		return nil, fmt.Errorf("license ID cannot be empty")
	}

	// PAPI call: GET https://1.2.3.4:8080/platform/17/license/licenses/{license-id}
	var response LicensesResponse
	err := client.Get(ctx, licensePath, licenseID, nil, nil, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Licenses) == 0 {
		return nil, fmt.Errorf("no license found with ID %s", licenseID)
	}

	return &response.Licenses[0], nil
}

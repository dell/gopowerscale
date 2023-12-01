/*
Copyright (c) 2022 Dell Inc, or its subsidiaries.

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

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test both GetQuota() and SetQuota()
func TestQuotaGetSet(t *testing.T) {
	volumeName := "test_quota_get_set"
	quotaSize := int64(1234567)
	var softLimit, advisoryLimit, softGracePrd int64

	// Setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)
	defer client.ClearQuota(defaultCtx, volumeName)

	// Make sure there is no quota yet
	quota, err := client.GetQuota(defaultCtx, volumeName)
	if quota != nil {
		panic(fmt.Sprintf("Quota should be nil: %v", quota))
	}
	if err == nil {
		panic("GetQuota should return an error when there isn't a quota.")
	}

	// Set the quota
	quotaID, err := client.SetQuotaSize(defaultCtx, volumeName, quotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}

	// Make sure the quota was set
	quota, err = client.GetQuotaByID(defaultCtx, quotaID)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic("Quota should not be nil")
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Unexpected new quota.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}
}

// Test GetAllQuotas()
func TestAllQuotasGet(t *testing.T) {
	// Get All the quotas
	quotas, err := client.GetAllQuotas(defaultCtx)
	if err != nil {
		panic(err)
	}
	assertNotNil(t, quotas)
}

// Test UpdateQuota()
func TestQuotaUpdate(t *testing.T) {
	volumeName := "test_quota_update"
	quotaSize := int64(1234567)
	updatedQuotaSize := int64(22345000)
	var softLimit, advisoryLimit, softGracePrd int64

	// Setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)
	defer client.ClearQuota(defaultCtx, volumeName)
	// Set the quota
	quotaID, err := client.SetQuotaSize(defaultCtx, volumeName, quotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}
	// Make sure the quota is initialized
	quota, err := client.GetQuotaByID(defaultCtx, quotaID)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Initial quota not set properly.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

	// Update the quota
	err = client.UpdateQuotaSizeByID(defaultCtx, quotaID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}

	// Make sure the quota is updated
	quota, err = client.GetQuotaByID(defaultCtx, quotaID)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Updated quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != updatedQuotaSize {
		panic(fmt.Sprintf("Updated quota not set properly.  Expected: %d Actual: %d", updatedQuotaSize, quota.Thresholds.Hard))
	}
}

// Test ClearQuota()
func TestQuotaClear(t *testing.T) {
	volumeName := "test_quota_clear"
	quotaSize := int64(1234567)
	var softLimit, advisoryLimit, softGracePrd int64

	// Setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)
	defer client.ClearQuota(defaultCtx, volumeName)
	// Set the quota
	quotaID, err := client.SetQuotaSize(defaultCtx, volumeName, quotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}
	// Make sure the quota is initialized
	quota, err := client.GetQuotaByID(defaultCtx, quotaID)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Initial quota not set properly.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

	// Clear the quota
	err = client.ClearQuota(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}

	// Make sure the quota is gone
	quota, err = client.GetQuotaByID(defaultCtx, quotaID)
	if err == nil {
		panic("Attempting to get a cleared quota should return an error but returned nil")
	}
	if quota != nil {
		panic(fmt.Sprintf("Cleared quota should be nil: %v", quota))
	}
}

// Test ClearQuotaByID()
func TestQuotaClearByID(t *testing.T) {
	volumeName := "test_quota_clear_by_id"
	quotaSize := int64(1234567)
	var softLimit, advisoryLimit, softGracePrd int64

	// Setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)
	defer client.ClearQuota(defaultCtx, volumeName)
	// Set the quota
	var quotaID string
	if quotaID, err = client.SetQuotaSize(defaultCtx, volumeName, quotaSize, softLimit, advisoryLimit, softGracePrd); err != nil {
		panic(err)
	}

	var quota Quota
	// Make sure the quota is initialized
	if quota, err = client.GetQuotaByID(defaultCtx, quotaID); err != nil {
		panic(err)
	}

	// Clear the quota
	if err := client.ClearQuotaByID(defaultCtx, quotaID); err != nil {
		panic(err)
	}

	// Make sure the quota is gone
	if quota, err = client.GetQuotaByID(defaultCtx, quotaID); err == nil {
		panic("Attempting to get a cleared quota should return an error but returned nil")
	} else if quota != nil {
		panic(fmt.Sprintf("Cleared quota should be nil: %v", quota))
	}
}

// Test IsQuotaLicenseActivated()
func TestIsQuotaLicenseActivated(t *testing.T) {
	t.Log("start TestIsQuotaLicenseActivated")

	isActivated, _ := client.IsQuotaLicenseActivated(defaultCtx)

	assert.True(t, isActivated)
}

// Test TestQuotaUpdateByID()
func TestQuotaUpdateByID(t *testing.T) {
	volumeName := "test_quota_update"
	quotaSize := int64(1234567)
	updatedQuotaSize := int64(22345000)
	var softLimit, advisoryLimit, softGracePrd int64

	// Setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	if err != nil {
		panic(err)
	}
	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)
	defer client.ClearQuota(defaultCtx, volumeName)
	// Set the quota
	id, err := client.SetQuotaSize(defaultCtx, volumeName, quotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}
	// Make sure the quota is initialized
	quota, err := client.GetQuotaByID(defaultCtx, id)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != quotaSize {
		panic(fmt.Sprintf("Initial quota not set properly.  Expected: %d Actual: %d", quotaSize, quota.Thresholds.Hard))
	}

	// Update the quota
	err = client.UpdateQuotaSizeByID(defaultCtx, quota.Id, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	if err != nil {
		panic(err)
	}

	// Make sure the quota is updated
	quota, err = client.GetQuotaByID(defaultCtx, id)
	if err != nil {
		panic(err)
	}
	if quota == nil {
		panic(fmt.Sprintf("Updated quota should not be nil: %v", quota))
	}
	if quota.Thresholds.Hard != updatedQuotaSize {
		panic(fmt.Sprintf("Updated quota not set properly.  Expected: %d Actual: %d", updatedQuotaSize, quota.Thresholds.Hard))
	}
}

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

	apiv1 "github.com/dell/goisilon/api/v1"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var quotaSize = int64(1234567)
var softLimit, advisoryLimit, softGracePrd int64
var quotaID string
var resume string
var ID string

// Test both GetQuota() and SetQuota()
func TestGetQuota(t *testing.T) {
	volumeName := "test_quota_get_set"

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiQuotaListResp)
		*resp = apiv1.IsiQuotaListResp{
			Quotas: []apiv1.IsiQuota{{}},
		}
	}).Once()

	// Make sure there is no quota yet
	_, err := client.GetQuota(defaultCtx, volumeName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()

	// Make sure there is no quota yet
	_, err = client.GetQuota(defaultCtx, volumeName)
	assert.NotNil(t, err)
}

// Test GetAllQuotas()
func TestGetAllQuotas(t *testing.T) {

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiQuotaListRespResume)
		*resp = &apiv1.IsiQuotaListRespResume{
			Quotas: []*apiv1.IsiQuota{},
			Resume: "",
		}
	}).Once()
	// Get All the quotas
	_, err := client.GetAllQuotas(defaultCtx)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Get All the quotas
	_, err = client.GetAllQuotas(defaultCtx)
	assert.NotNil(t, err)
}

func TestGetQuotasWithResume(t *testing.T) {

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiQuotaListRespResume)
		*resp = apiv1.IsiQuotaListRespResume{
			Quotas: []*apiv1.IsiQuota{},
			Resume: "",
		}
	}).Once()
	// Clear the quota
	_, err = client.GetQuotasWithResume(defaultCtx, resume)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Clear the quota
	_, err = client.GetQuotasWithResume(defaultCtx, resume)
	assert.NotNil(t, err)
}
func TestGetQuotaByID(t *testing.T) {

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiQuotaListResp)
		*resp = apiv1.IsiQuotaListResp{
			Quotas: []apiv1.IsiQuota{{}},
		}
	}).Once()

	// Clear the quota
	_, err = client.GetQuotaByID(defaultCtx, ID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Clear the quota
	_, err = client.GetQuotaByID(defaultCtx, ID)
	assert.NotNil(t, err)
}

func TestGetQuotawithPath(t *testing.T) {

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiQuotaListResp)
		*resp = apiv1.IsiQuotaListResp{
			Quotas: []apiv1.IsiQuota{{}},
		}
	}).Once()

	// Clear the quota
	_, err = client.GetQuotaWithPath(defaultCtx, ID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Clear the quota
	_, err = client.GetQuotaWithPath(defaultCtx, ID)
	assert.NotNil(t, err)
}

// Test UpdateQuota()
func TestUpdateQuotaSizeByID(t *testing.T) {

	updatedQuotaSize := int64(22345000)
	var softLimit, advisoryLimit, softGracePrd int64

	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()

	// Update the quota
	err = client.UpdateQuotaSizeByID(defaultCtx, quotaID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("not found")).Once()
	// Update the quota
	err = client.UpdateQuotaSizeByID(defaultCtx, quotaID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.NotNil(t, err)
}

// Test ClearQuota()
func TestClearQuota(t *testing.T) {

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	// Clear the quota
	err = client.ClearQuota(defaultCtx, volumeName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Clear the quota
	err = client.ClearQuota(defaultCtx, volumeName)
	assert.NotNil(t, err)
}

func TestClearQuotaByID(t *testing.T) {

	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	// Clear the quota
	err := client.ClearQuotaByID(defaultCtx, quotaID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// Clear the quota
	err = client.ClearQuotaByID(defaultCtx, quotaID)
	assert.NotNil(t, err)
}

func TestIsQuotaLicenseActivated(t *testing.T) {

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	_, err := client.IsQuotaLicenseActivated(defaultCtx)
	assert.Nil(t, err)

	// client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	// _, err = client.IsQuotaLicenseActivated(defaultCtx)
	// assert.NotNil(t, err)
}

// Test TestQuotaUpdateByID()
func TestQuotaUpdateByID(_ *testing.T) {
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
	err = client.UpdateQuotaSizeByID(defaultCtx, quota.ID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
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

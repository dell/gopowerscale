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
package goisilon

import (
	"fmt"
	"testing"

	apiv1 "github.com/dell/goisilon/api/v1"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	quotaSize                              = int64(1234567)
	softLimit, advisoryLimit, softGracePrd int64
	quotaID, name, zone                    string
	resume                                 string
	ID                                     string
	size                                   int64 = 22345000
	container                              bool
)

func TestGetQuota(t *testing.T) {
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv1.IsiQuotaListResp)
		*resp = apiv1.IsiQuotaListResp{
			Quotas: []apiv1.IsiQuota{{}},
		}
	}).Once()
	_, err := client.GetQuota(defaultCtx, volumeName)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	_, err = client.GetQuota(defaultCtx, volumeName)
	assert.NotNil(t, err)
}

func TestGetAllQuotas(t *testing.T) {
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.IsiQuotaListRespResume)
		*resp = &apiv1.IsiQuotaListRespResume{
			Quotas: []*apiv1.IsiQuota{},
			Resume: "",
		}
	}).Once()
	_, err := client.GetAllQuotas(defaultCtx)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
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
	_, err = client.GetQuotasWithResume(defaultCtx, resume)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
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

	_, err = client.GetQuotaByID(defaultCtx, ID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
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

	_, err = client.GetQuotaWithPath(defaultCtx, ID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	_, err = client.GetQuotaWithPath(defaultCtx, ID)
	assert.NotNil(t, err)
}

func TestCreateQuota(t *testing.T) {
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err = client.CreateQuota(defaultCtx, quotaID, container, size, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)
}

func TestCreateQuotaWithPath(t *testing.T) {
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err = client.CreateQuotaWithPath(defaultCtx, quotaID, container, size, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)
}

func TestSetQuotaSize(t *testing.T) {
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	_, err = client.SetQuotaSize(defaultCtx, name, size, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)
}

func TestUpdateQuotaSizeByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	updatedQuotaSize := int64(22345000)
	var softLimit, advisoryLimit, softGracePrd int64

	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.UpdateQuotaSizeByID(defaultCtx, quotaID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("not found")).Once()
	err = client.UpdateQuotaSizeByID(defaultCtx, quotaID, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.NotNil(t, err)
}

func TestUpdateQuotaSize(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	updatedQuotaSize := int64(22345000)
	quotaPath := "/platform/1/quota/quotas/" + quotaID
	var softLimit, advisoryLimit, softGracePrd int64

	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).
		Return(nil).Run(func(args mock.Arguments) {
		// Create a response for IsiQuotaListResp
		resp := args.Get(5).(*apiv1.IsiQuotaListResp)
		*resp = apiv1.IsiQuotaListResp{
			Quotas: []apiv1.IsiQuota{{
				ID: quotaID,
			}},
		}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.UpdateQuotaSize(defaultCtx, quotaPath, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("/ifs/data").Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(fmt.Errorf("not found")).Once()
	err = client.UpdateQuotaSize(defaultCtx, quotaPath, updatedQuotaSize, softLimit, advisoryLimit, softGracePrd)
	assert.NotNil(t, err)
}

func TestClearQuota(t *testing.T) {
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err = client.ClearQuota(defaultCtx, volumeName)
	assert.Nil(t, err)
}

func TestClearQuotaWithPath(t *testing.T) {
	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(nil).Once()
	err = client.ClearQuotaWithPath(defaultCtx, volumeName)
	assert.Nil(t, err)
}

func TestClearQuotaByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err := client.ClearQuotaByIDWithZone(defaultCtx, quotaID, zone)
	assert.Nil(t, err)
}

func TestClearQuotaByID(t *testing.T) {
	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(nil).Once()
	err := client.ClearQuotaByID(defaultCtx, quotaID)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Delete", anyArgs[0:6]...).Return(fmt.Errorf("not found")).Once()
	err = client.ClearQuotaByID(defaultCtx, quotaID)
	assert.NotNil(t, err)
}

func TestIsQuotaLicenseActivated(t *testing.T) {
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	_, err := client.IsQuotaLicenseActivated(defaultCtx)
	assert.Nil(t, err)
}

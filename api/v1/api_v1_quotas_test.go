/*
Copyright (c) 2025 Dell Inc, or its subsidiaries.

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
	"errors"
	"testing"

	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestModifyIsiQuotaByID_WrongType_SoftLimit(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	params := map[string]interface{}{
		"SoftLimit": "200", // string instead of int64
	}

	err := ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SoftLimit must be int64 or nil")

	anyArgs8 := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
	client.AssertNotCalled(t, "DoWithHeaders", anyArgs8...)
}

func TestModifyIsiQuotaByID_WrongType_SoftGracePrd(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	params := map[string]interface{}{
		"SoftGracePrd": "300", // string instead of int64
	}

	err := ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "SoftGracePrd must be int64 or nil")

	anyArgs8 := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
	client.AssertNotCalled(t, "DoWithHeaders", anyArgs8...)
}

func TestModifyIsiQuotaByID_UnsupportedKey(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Unsupported key should error before making any HTTP call
	params := map[string]interface{}{
		"HardLimit": int64(100),
	}

	err := ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported param key: HardLimit")

	// Ensure no request was attempted
	anyArgs8 := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
	client.AssertNotCalled(t, "DoWithHeaders", anyArgs8...)
}

func TestModifyIsiQuotaByID_WrongType(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	// Wrong type for AdvisoryLimit should return an error and not call HTTP
	params := map[string]interface{}{
		"AdvisoryLimit": "100", // string instead of int64
	}

	err := ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "AdvisoryLimit must be int64 or nil")

	// Ensure no request was attempted
	anyArgs8 := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}
	client.AssertNotCalled(t, "DoWithHeaders", anyArgs8...)
}

func TestGetIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListResp)
		*resp = IsiQuotaListResp{
			Quotas: []IsiQuota{},
		}
	}).Once()
	client.On("Get", anyArgs...).Return(errors.New("Quota not found: ")).Run(nil).Once()
	_, err := GetIsiQuota(ctx, client, "")
	assert.Equal(t, errors.New("Quota not found: "), err)

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err = GetIsiQuota(ctx, client, "")
	assert.Equal(t, errors.New("Quota not found: "), err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListResp)
		*resp = IsiQuotaListResp{
			Quotas: []IsiQuota{
				{
					ID: "test",
				},
			},
		}
	}).Once()
	client.On("Get", anyArgs...).Return(nil).Run(nil).Once()
	_, err = GetIsiQuota(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetAllIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiQuotaListRespResume)
		*resp = &IsiQuotaListRespResume{
			Quotas: []*IsiQuota{},
			Resume: "",
		}
	}).Once()
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**IsiQuotaListRespResume)
		*resp = &IsiQuotaListRespResume{
			Quotas: []*IsiQuota{
				{
					ID: "test",
				},
			},
			Resume: "resume",
		}
	}).Once()
	_, err := GetAllIsiQuota(ctx, client)
	assert.Equal(t, nil, err)

	client.On("Get", anyArgs...).Return(errors.New("error")).Twice()
	_, err = GetAllIsiQuota(ctx, client)
	assert.Equal(t, errors.New("error"), err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = GetAllIsiQuota(ctx, client)
	assert.Equal(t, errors.New("error"), err)
}

func TestGetIsiQuotaWithResume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err := GetIsiQuotaWithResume(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListRespResume)
		*resp = IsiQuotaListRespResume{
			Quotas: []*IsiQuota{
				{
					ID:   "test",
					Path: "/test",
				},
			},
			Resume: "resume",
		}
	}).Once()
	_, err = GetIsiQuotaWithResume(ctx, client, "/test")
	assert.Equal(t, nil, err)
}

func TestGetIsiQuotaByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err := GetIsiQuotaByID(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListResp)
		*resp = IsiQuotaListResp{
			Quotas: []IsiQuota{
				{
					ID: "test-id",
				},
			},
		}
	}).Once()
	_, err = GetIsiQuotaByID(ctx, client, "test-id")
	assert.Equal(t, nil, err)

	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListResp)
		*resp = IsiQuotaListResp{
			Quotas: []IsiQuota{},
		}
	}).Once()
	_, err = GetIsiQuotaByID(ctx, client, "test-id")
	assert.Equal(t, errors.New("Quota not found: test-id"), err)
}

func TestSetIsiQuotaHardThreshold(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Post", anyArgs...).Return(nil).Twice()
	_, err := SetIsiQuotaHardThreshold(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
}

func TestUpdateIsiQuotaHardThreshold(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	err := UpdateIsiQuotaHardThreshold(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, errors.New("Quota not found: "), err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*IsiQuotaListResp)
		*resp = IsiQuotaListResp{
			Quotas: []IsiQuota{
				{
					ID: "test",
				},
			},
		}
	}).Once()
	client.On("Put", anyArgs...).Return(nil).Once()
	err = UpdateIsiQuotaHardThreshold(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
}

func TestUpdateIsiQuotaHardThresholdByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	err := UpdateIsiQuotaHardThresholdByID(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
}

func TestModifyIsiQuotaByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	params := map[string]interface{}{
		"AdvisoryLimit": int64(10),
		"SoftLimit":     int64(20),
		"SoftGracePrd":  int64(30),
	}

	paramsClear := map[string]interface{}{
		"AdvisoryLimit": int64(0),
		"SoftLimit":     int64(0),
		"SoftGracePrd":  int64(0),
	}

	anyArgs8 := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

	// Success case with values
	client.On("DoWithHeaders", anyArgs8...).Return(nil).Once()
	err := ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Equal(t, nil, err)

	// Success case with clearing limits
	client.On("DoWithHeaders", anyArgs8...).Return(nil).Once()
	err = ModifyIsiQuotaByID(ctx, client, "test-id", paramsClear)
	assert.Equal(t, nil, err)

	// Success case with nil limits
	paramsNil := map[string]interface{}{
		"AdvisoryLimit": nil,
		"SoftLimit":     nil,
		"SoftGracePrd":  nil,
	}
	client.On("DoWithHeaders", anyArgs8...).Return(nil).Once()
	err = ModifyIsiQuotaByID(ctx, client, "test-id", paramsNil)
	assert.Equal(t, nil, err)

	// Error case
	client.On("DoWithHeaders", anyArgs8...).Return(errors.New("error")).Once()
	err = ModifyIsiQuotaByID(ctx, client, "test-id", params)
	assert.Equal(t, errors.New("error"), err)
}

func TestDeleteIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuota(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiQuotaByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuotaByID(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiQuotaByIDWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuotaByIDWithZone(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

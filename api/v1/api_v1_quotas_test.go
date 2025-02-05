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

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

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
}

func TestUpdateIsiQuotaHardThresholdByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	err := UpdateIsiQuotaHardThresholdByID(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
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

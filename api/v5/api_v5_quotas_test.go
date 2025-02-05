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

package v5

import (
	"context"
	"fmt"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestIsQuotaLicenseActivated(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	licenseStatus := QuotaLicenseStatus{
		value: "Expired",
	}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := IsQuotaLicenseActivated(ctx, client)
	assert.Equal(t, nil, err)

	client.On("Get", ctx, mock.Anything).Return("unlicensed", nil).Twice()
	_, err = IsQuotaLicenseActivated(ctx, client)
	assert.Nil(t, err)

	client.On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*QuotaLicense)
		resp.STATUS = licenseStatus.toString()
	}).Once()
	_, err = IsQuotaLicenseActivated(ctx, client)
	assert.Nil(t, err)
}

func TestIsQuotaLicenseStatusValid(t *testing.T) {
	licenseStatus := QuotaLicenseStatus{
		value: "Expired",
	}
	value := isQuotaLicenseStatusValid(licenseStatus.toString())
	assert.Equal(t, true, value)
}

func TestGetIsiQuotaLicense(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(fmt.Errorf("not found")).Twice()
	_, err := GetIsiQuotaLicense(context.Background(), client)
	assert.NotNil(t, err)
}

func TestGetIsiQuotaLicenseStatus(t *testing.T) {
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	_, err := getIsiQuotaLicenseStatus(context.Background(), client)
	assert.NotNil(t, err)
}

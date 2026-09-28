/*
Copyright (c) 2022-2026 Dell Inc, or its subsidiaries.

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

func TestGetIsiNodePools(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	anyArgs := []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

	// Test success case
	client.On("Get", anyArgs...).Return(nil).Once()
	resp, err := GetIsiNodePools(ctx, client)
	assert.NoError(t, err)
	assert.NotNil(t, resp)

	// Test error case
	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(errors.New("error in get isi node pools")).Once()
	resp, err = GetIsiNodePools(ctx, client)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

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
)

func TestGetZoneByName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetZoneByName(ctx, client, "name")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err = GetZoneByName(ctx, client, "name")
	assert.Equal(t, nil, err)
}

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
package v3

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestGetIsiClusterNode(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiClusterNode(ctx, client, 0)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetIsiClusterNodes(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiClusterNodes(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetIsiClusterIdentity(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiClusterIdentity(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetIsiClusterConfig(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiClusterConfig(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestIsIOInProgress(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := IsIOInProgress(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetIsiFloatStats(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiFloatStats(ctx, client, []string{})
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestGetIsiStats(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiStats(ctx, client, []string{})
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

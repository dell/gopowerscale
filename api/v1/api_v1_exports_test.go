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

func TestExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	err := Export(ctx, client, "")
	assert.Equal(t, errors.New("no path set"), err)

	client.On("User", anyArgs...).Return("").Twice()
	client.On("Group", anyArgs...).Return("").Twice()
	client.On("Post", anyArgs...).Return(nil).Twice()
	err = Export(ctx, client, "path")
	assert.Equal(t, nil, err)
}

func TestSetExportClients(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	err := SetExportClients(ctx, client, 0, []string{""})
	assert.Equal(t, nil, err)
}

func TestUnexport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	err := Unexport(ctx, client, 0)
	assert.Equal(t, errors.New("no path Id set"), err)

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err = Unexport(ctx, client, 1)
	assert.Equal(t, nil, err)
}

func TestGetIsiExports(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiExports(ctx, client)
	assert.Equal(t, nil, err)
}

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

package v2

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestContainerChildList_MarshalJSON(t *testing.T) {
	name1 := "child1"
	name2 := "child2"
	path1 := "/path/to/child1"
	path2 := "/path/to/child2"
	size1 := 123
	size2 := 456
	fileType := "type"
	owner := "owner"
	group := "group"
	mode := FileMode(0o755)

	tests := []struct {
		name     string
		input    ContainerChildList
		expected string // Expected JSON string
	}{
		{
			name: "Non-empty list",
			input: ContainerChildList{
				&ContainerChild{Name: &name1, Path: &path1, Size: &size1, Type: &fileType, Owner: &owner, Group: &group, Mode: &mode},
				&ContainerChild{Name: &name2, Path: &path2, Size: &size2, Type: &fileType, Owner: &owner, Group: &group, Mode: &mode},
			},
			expected: `{"children":[{"name":"child1","container_path":"/path/to/child1","size":123,"type":"type","owner":"owner","group":"group","mode":"0755"}, {"name":"child2","container_path":"/path/to/child2","size":456,"type":"type","owner":"owner","group":"group","mode":"0755"}]}`,
		},
		{
			name:     "Empty list",
			input:    ContainerChildList{},
			expected: `{}`,
		},
		{
			name:     "Nil list",
			input:    nil,
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.input.MarshalJSON()
			assert.NoError(t, err)
			assert.JSONEq(t, tt.expected, string(result))
		})
	}
}

func TestContainerChildList_UnmarshalJSON(t *testing.T) {
	name1 := "child1"
	path1 := "/path/to/child1"
	size1 := 123
	fileType := "type"
	owner := "owner"
	group := "group"
	mode := FileMode(0o755)

	tests := []struct {
		name        string
		input       string
		expected    ContainerChildList
		expectError bool
	}{
		{
			name:  "Non-empty list",
			input: `{"children":[{"name":"child1","container_path":"/path/to/child1","size":123,"type":"type","owner":"owner","group":"group","mode":"0755"}]}`,
			expected: ContainerChildList{
				&ContainerChild{Name: &name1, Path: &path1, Size: &size1, Type: &fileType, Owner: &owner, Group: &group, Mode: &mode},
			},
			expectError: false,
		},
		{
			name:        "Empty list",
			input:       `{"children":[]}`,
			expected:    ContainerChildList{},
			expectError: false,
		},
		{
			name:        "Nil input",
			input:       `{"children":null}`,
			expected:    nil,
			expectError: false,
		},
		{
			name:        "Invalid JSON",
			input:       `{"children":`,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ContainerChildList
			err := json.Unmarshal([]byte(tt.input), &result)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

func TestContainerChildrenMapAll(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, err := ContainerChildrenMapAll(context.Background(), client, "")
	assert.NoError(t, err)
}

func TestContainerChildrenGetAll(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, err := ContainerChildrenGetAll(context.Background(), client, "")
	assert.NoError(t, err)
}

func TestContainerChildrenGetQuery(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, errChan := ContainerChildrenGetQuery(context.Background(), client, "", 0, 0, "", "", nil, nil)
	assert.NoError(t, <-errChan)

	// when objectType is not empty
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, errChan = ContainerChildrenGetQuery(context.Background(), client, "", 0, 0, "", "test-object", nil, []string{"a", "b"})
	assert.NoError(t, <-errChan)

	// when len(sort) > 0
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, errChan = ContainerChildrenGetQuery(context.Background(), client, "", 0, 0, "", "", []string{"a", "b"}, nil)
	assert.NoError(t, <-errChan)

	// when len(sort) > 0 and sortDir is not "asc" or "desc"
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Get", anyArgs...).Return(nil).Once()

	_, errChan = ContainerChildrenGetQuery(context.Background(), client, "", 0, 0, "", "asc", []string{"a", "b"}, nil)
	assert.NoError(t, <-errChan)
}

func TestContainerChildrenPostQuery(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Post", anyArgs...).Return(nil).Once()

	_, err := ContainerChildrenPostQuery(context.Background(), client, "", 0, 0, nil)
	assert.NoError(t, err)
}

func TestContainerCreateDir(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Put", anyArgs...).Return(nil).Once()

	err := ContainerCreateDir(context.Background(), client, "", "", 0, false, false)
	assert.NoError(t, err)

	// overwrite && recursive both true
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Put", anyArgs...).Return(nil).Once()

	err = ContainerCreateDir(context.Background(), client, "", "", 0, true, true)
	assert.NoError(t, err)

	// when eitther overwrite or recursive is true
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Put", anyArgs...).Return(nil).Once()

	err = ContainerCreateDir(context.Background(), client, "", "", 0, false, true)
	assert.NoError(t, err)
}

func TestContainerCreateFile(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Put", anyArgs...).Return(nil).Once()

	err := ContainerCreateFile(context.Background(), client, "", "", 0, 0, nil, false)
	assert.NoError(t, err)
}

func TestContainerChildDelete(t *testing.T) {
	client := &mocks.Client{}
	client.On("VolumesPath", anyArgs...).Return(testVolumePath).Once()
	client.On("Delete", anyArgs...).Return(nil).Once()

	err := ContainerChildDelete(context.Background(), client, "", true)
	assert.NoError(t, err)
}

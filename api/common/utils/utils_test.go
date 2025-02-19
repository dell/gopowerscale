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
package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsStringInSlice(t *testing.T) {
	list := []string{"hello", "world", "jason"}

	assert.True(t, IsStringInSlice("world", list))
	assert.False(t, IsStringInSlice("mary", list))
	assert.False(t, IsStringInSlice("harry", nil))
}

func TestIsStringInSlices(t *testing.T) {
	list := []string{"hello", "world", "jason"}

	assert.True(t, IsStringInSlices("world", list))
	assert.False(t, IsStringInSlices("mary", list))
	assert.False(t, IsStringInSlices("harry", nil))
}

func TestRemoveStringFromSlice(t *testing.T) {
	list := []string{"hello", "world", "jason"}

	result := RemoveStringFromSlice("hello", list)

	assert.Equal(t, 3, len(list))
	assert.Equal(t, 2, len(result))

	result = RemoveStringFromSlice("joe", list)

	assert.Equal(t, 3, len(result))
}

func TestRemoveStringsFromSlice(t *testing.T) {
	list := []string{"hello", "world", "jason"}

	filterList := []string{"hello", "there", "chap", "world"}

	result := RemoveStringsFromSlice(filterList, list)

	assert.Equal(t, 1, len(result))
}

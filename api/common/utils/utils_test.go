/*
 Copyright (c) 2022 Dell Inc, or its subsidiaries.

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
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	fmt.Print("executing TestMain\n")
}

func TestIsStringInSlice(t *testing.T) {

	var list = []string{"hello", "world", "jason"}

	assert.True(t, IsStringInSlice("world", list))
	assert.False(t, IsStringInSlice("mary", list))
	assert.False(t, IsStringInSlice("harry", nil))
}

func TestRemoveStringFromSlice(t *testing.T) {

	var list = []string{"hello", "world", "jason"}

	var result = RemoveStringFromSlice("hello", list)

	assert.Equal(t, 3, len(list))
	assert.Equal(t, 2, len(result))

	result = RemoveStringFromSlice("joe", list)

	assert.Equal(t, 3, len(result))
}

func TestRemoveStringsFromSlice(t *testing.T) {

	var list = []string{"hello", "world", "jason"}

	var filterList = []string{"hello", "there", "chap", "world"}

	var result = RemoveStringsFromSlice(filterList, list)

	assert.Equal(t, 1, len(result))
}

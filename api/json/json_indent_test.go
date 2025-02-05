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

package json

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndent(t *testing.T) {
	// Test case for simple JSON object
	src := []byte(`{"key1":"value1","key2":"value2"}`)
	dst := &bytes.Buffer{}
	prefix := ""
	indent := "    "
	err := Indent(dst, src, prefix, indent)
	assert.NoError(t, err)
	expected := `{
		"key1": "value1",
		"key2": "value2"
	}`
	// assert.Equal(t, expected, dst.String())

	// Test case for nested JSON object
	src = []byte(`{"key1": {"nestedKey1": "nestedValue1", "nestedKey2": "nestedValue2"}, "key2": "value2"}`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.NoError(t, err)
	expected = `{
		"key1": {
			"nestedKey1": "nestedValue1",
			"nestedKey2": "nestedValue2"
		},
		"key2": "value2"
	}`
	// assert.Equal(t, expected, dst.String())

	// Test case for JSON array
	src = []byte(`[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.NoError(t, err)
	expected = `[
		{
			"key1": "value1",
			"key2": "value2"
		},
		{
			"key3": "value3",
			"key4": "value4"
		}
	]`
	// assert.Equal(t, expected, dst.String())

	// Test case for invalid JSON
	src = []byte(`{"key1":"value1","key2":"value2"`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.Error(t, err)
	assert.Equal(t, "", dst.String())

	// Test case for empty JSON object
	src = []byte(`{}`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.NoError(t, err)
	expected = `{}`
	assert.Equal(t, expected, dst.String())

	// Test case for empty JSON array
	src = []byte(`[]`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.NoError(t, err)
	expected = `[]`
	assert.Equal(t, expected, dst.String())
}

func TestCompact(t *testing.T) {
	// Test case for simple JSON object
	src := []byte(`{"key1":"value1","key2":"value2"}`)
	dst := &bytes.Buffer{}
	escape := true
	err := compact(dst, src, escape)
	assert.NoError(t, err)
	expected := `{"key1":"value1","key2":"value2"}`
	assert.Equal(t, expected, dst.String())

	// Test case for nested JSON object
	src = []byte(`{"key1": {"nestedKey1": "nestedValue1", "nestedKey2": "nestedValue2"}, "key2": "value2"}`)
	dst.Reset()
	err = compact(dst, src, escape)
	assert.NoError(t, err)
	expected = `{"key1":{"nestedKey1":"nestedValue1","nestedKey2":"nestedValue2"},"key2":"value2"}`
	assert.Equal(t, expected, dst.String())

	// Test case for JSON array
	src = []byte(`[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]`)
	dst.Reset()
	err = compact(dst, src, escape)
	assert.NoError(t, err)
	expected = `[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]`
	assert.Equal(t, expected, dst.String())
	err = Compact(dst, src)

	// Test case for invalid JSON
	src = []byte(`{"key1":"value1","key2":"value2"`)
	dst.Reset()
	err = compact(dst, src, escape)
	assert.Error(t, err)
	assert.Equal(t, "", dst.String())

	src = []byte(`{"<":"value1","key2":"value2"`)
	dst.Reset()
	err = compact(dst, src, escape)
	assert.Error(t, err)
	assert.Equal(t, "", dst.String())

	src = []byte("Hello\xE2\x80\xA8World")
	err = compact(dst, src, false)

	dst.Reset()

	src = []byte("Hello\xE2\x80\xA9World")
	err = compact(dst, src, false)

	dst.Reset()

	src = []byte("Hello<World>")
	err = compact(dst, src, true)
}

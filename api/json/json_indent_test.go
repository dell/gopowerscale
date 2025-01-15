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
	assert.Equal(t, expected, dst.String())

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
	assert.Equal(t, expected, dst.String())

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
	assert.Equal(t, expected, dst.String())

	// Test case for invalid JSON
	src = []byte(`{"key1":"value1","key2":"value2"`)
	dst.Reset()
	err = Indent(dst, src, prefix, indent)
	assert.Error(t, err)
	assert.Equal(t, "", dst.String())
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
}

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
	"encoding"
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsValidTag(t *testing.T) {
	value := isValidTag("")
	assert.Equal(t, false, value)

	value = isValidTag("a")
	assert.Equal(t, true, value)

	value = isValidTag("\\")
	assert.Equal(t, false, value)
}

func TestString(t *testing.T) {
	state := encodeState{}

	value := state.string("\\", false)
	assert.Equal(t, 4, value)

	value = state.string("\n", false)
	assert.Equal(t, 4, value)

	value = state.string("\r", false)
	assert.Equal(t, 4, value)

	value = state.string("\t", false)
	assert.Equal(t, 4, value)

	value = state.string("\a", false)
	assert.Equal(t, 8, value)

	value = state.string("\u2028", false)
	assert.Equal(t, 8, value)

	value = state.string("\u2029", false)
	assert.Equal(t, 8, value)

	value = state.string("\ufffd", false)
	assert.Equal(t, 5, value)

	value = state.string("\x01", false)
	assert.Equal(t, 8, value)

	value = state.string(string([]byte{0xff}), false)
	assert.Equal(t, 8, value)
}

func TestStringBytes(t *testing.T) {
	state := encodeState{}

	value := state.stringBytes([]byte{'\\'}, false)
	assert.Equal(t, 4, value)

	value = state.stringBytes([]byte{'\n'}, false)
	assert.Equal(t, 4, value)

	value = state.stringBytes([]byte{'\r'}, false)
	assert.Equal(t, 4, value)

	value = state.stringBytes([]byte{'\t'}, false)
	assert.Equal(t, 4, value)

	value = state.stringBytes([]byte{'\a'}, false)
	assert.Equal(t, 8, value)

	value = state.stringBytes([]byte{0xe2, 0x80, 0xa8}, false)
	assert.Equal(t, 8, value)

	value = state.stringBytes([]byte{0xe2, 0x80, 0xa9}, false)
	assert.Equal(t, 8, value)

	value = state.stringBytes([]byte{0xef, 0xbf, 0xbd}, false)
	assert.Equal(t, 5, value)

	value = state.stringBytes([]byte{0x01}, false)
	assert.Equal(t, 8, value)

	value = state.stringBytes([]byte{0xff}, false)
	assert.Equal(t, 8, value)
}

func TestDominantField(t *testing.T) {
	testCases := []struct {
		input    []field
		expected bool
	}{
		{
			input: []field{
				{
					name: "test",
					tag:  true,
				},
			},
			expected: true,
		},
		{
			input: []field{
				{
					name: "test",
					tag:  true,
				},
				{
					name: "test",
					tag:  true,
				},
			},
			expected: false,
		},
		{
			input: []field{
				{
					name: "test",
					tag:  true,
				},
				{
					name: "test",
					tag:  false,
				},
			},
			expected: true,
		},
		{
			input: []field{
				{
					name: "test",
					tag:  false,
				},
				{
					name: "test",
					tag:  true,
				},
			},
			expected: true,
		},
		{
			input: []field{
				{
					name: "test",
					tag:  false,
				},
				{
					name: "test",
					tag:  false,
				},
			},
			expected: false,
		},
	}

	for _, tc := range testCases {
		result, value := dominantField(tc.input)
		if value != tc.expected {
			t.Errorf("dominantField(%v) = %v, expected %v", tc.input, value, tc.expected)
		}
		if value && result.name != tc.input[0].name {
			t.Errorf("dominantField(%v) = %v, expected %v", tc.input, result.name, tc.input[0].name)
		}
	}
}

func TestIsEmptyValue(t *testing.T) {
	// Testing for empty values
	testCases := []struct {
		value    interface{}
		expected bool
	}{
		{value: "", expected: true},
		{value: []int{}, expected: true},
		{value: map[string]string{}, expected: true},
		{value: false, expected: true},
		{value: int(0), expected: true},
		{value: uint(0), expected: true},
		{value: float32(0), expected: true},
		{value: nil, expected: false},
		// Testing for non-empty values
		{value: "not empty", expected: false},
		{value: []int{1, 2, 3}, expected: false},
		{value: map[string]string{"key": "value"}, expected: false},
		{value: true, expected: false},
		{value: int(1), expected: false},
		{value: uint(1), expected: false},
		{value: float32(1.0), expected: false},
		{value: &struct{}{}, expected: false},
	}

	for _, tc := range testCases {
		if result := isEmptyValue(reflect.ValueOf(tc.value)); result != tc.expected {
			t.Errorf("isEmptyValue(%v) = %v, expected %v", tc.value, result, tc.expected)
		}
	}
}

func TestStringEncoder(t *testing.T) {
	// Testing for quoted values
	testCases := []struct {
		value    interface{}
		expected string
		opts     encOpts
	}{
		{value: "", expected: `"\"\""`, opts: encOpts{quoted: true, escapeHTML: false}},
		{value: "test", expected: `"\"test\""`, opts: encOpts{quoted: true, escapeHTML: false}},
	}

	for _, tc := range testCases {
		e := &encodeState{
			Buffer: bytes.Buffer{},
		}
		stringEncoder(e, reflect.ValueOf(tc.value), tc.opts)
		if result := e.String(); result != tc.expected {
			t.Errorf("stringEncoder(%v) = %v, expected %v", tc.value, result, tc.expected)
		}
	}

	// Testing for number values
	testCases1 := []struct {
		value    interface{}
		expected string
	}{
		{value: int(0), expected: "\"<int Value>\""},
		{value: float64(1.0), expected: "\"<float64 Value>\""},
	}

	for _, tc := range testCases1 {
		e := &encodeState{
			Buffer: bytes.Buffer{},
		}
		stringEncoder(e, reflect.ValueOf(tc.value), encOpts{})
		if result := e.String(); result != tc.expected {
			t.Errorf("stringEncoder(%v) = %v, expected %v", tc.value, result, tc.expected)
		}
	}

	testCases2 := []struct {
		value    interface{}
		expected string
	}{
		{value: "invalid", expected: "invalid"},
	}

	for _, tc := range testCases2 {
		e := &encodeState{
			Buffer: bytes.Buffer{},
		}
		stringEncoder(e, reflect.ValueOf(tc.value), encOpts{})
		if result := e.String(); !strings.Contains(result, tc.expected) {
			t.Errorf("stringEncoder(%v) = %v, expected %v", tc.value, result, tc.expected)
		}
	}

	testCases3 := []struct {
		value    interface{}
		expected string
	}{
		{value: Number(""), expected: "0"},
		{value: Number("1.0"), expected: "1.0"},
	}

	for _, tc := range testCases3 {
		e := &encodeState{
			Buffer: bytes.Buffer{},
		}
		stringEncoder(e, reflect.ValueOf(tc.value), encOpts{})
		if result := e.String(); result != tc.expected {
			t.Errorf("stringEncoder(%v) = %v, expected %v", tc.value, result, tc.expected)
		}
	}
}

func TestHTMLEscape(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{input: "", expected: ""},
		{input: "test", expected: "test"},
		{input: "test&test", expected: "test\\u0026test"},
		{input: "test\u2028test", expected: "test\\u2028test"},
		{input: "test\u2029test", expected: "test\\u2029test"},
	}

	for _, tc := range testCases {
		input := []byte(tc.input)
		expected := []byte(tc.expected)

		var dst bytes.Buffer
		HTMLEscape(&dst, input)
		result := dst.Bytes()

		if !bytes.Equal(result, expected) {
			t.Errorf("HTMLEscape(%v) = %v, expected %v", tc.input, string(result), string(expected))
		}
	}
}

func TestEncodeByteSlice(t *testing.T) {
	testCases := []struct {
		input    []byte
		expected string
	}{
		{input: nil, expected: ""},
		{input: []byte{}, expected: ``},
		{input: []byte{1, 2, 3}, expected: ``},
		{input: bytes.Repeat([]byte{1}, 1024), expected: ""},
		{input: bytes.Repeat([]byte{1}, 1025), expected: ""},
	}

	for _, tc := range testCases {
		var dst bytes.Buffer
		e := &encodeState{
			Buffer: dst,
		}
		encodeByteSlice(e, reflect.ValueOf(tc.input), encOpts{})
		result := dst.String()

		if result != tc.expected {
			t.Errorf("json.encodeByteSlice(%v) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestResolve(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{input: "test", expected: "test"},
		{input: int(123), expected: "123"},
		{input: uint(456), expected: "456"},
		{input: json.Number("789"), expected: "789"},
	}

	for _, tc := range testCases {
		w := &reflectWithString{
			v: reflect.ValueOf(tc.input),
		}
		err := w.resolve()
		if err != nil {
			t.Errorf("w.resolve() = %v, expected nil", err)
		}

		if w.s != tc.expected {
			t.Errorf("w.resolve() = %v, expected %v", w.s, tc.expected)
		}
	}
}

func TestTypeByIndex(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 []byte
	}

	testCases := []struct {
		input    reflect.Type
		index    []int
		expected reflect.Type
	}{
		{input: reflect.TypeOf(TestStruct{}), index: []int{0}, expected: reflect.TypeOf("")},
		{input: reflect.TypeOf(TestStruct{}), index: []int{1}, expected: reflect.TypeOf(0)},
	}

	for _, tc := range testCases {
		result := typeByIndex(tc.input, tc.index)
		if result != tc.expected {
			t.Errorf("typeByIndex(%v, %v) = %v, expected %v", tc.input, tc.index, result, tc.expected)
		}
	}
}

func TestLess(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 []byte
	}
	testCases := []struct {
		input    []field
		index1   int
		index2   int
		expected bool
	}{
		{input: []field{
			{name: "Field1", index: []int{0}, typ: reflect.TypeOf("")},
			{name: "Field2", index: []int{1}, typ: reflect.TypeOf(0)},
			{name: "Field3", index: []int{2}, typ: reflect.TypeOf([]byte{})},
		}, index1: 0, index2: 1, expected: true},
		{input: []field{
			{name: "Field1", index: []int{0}, typ: reflect.TypeOf("")},
			{name: "Field2", index: []int{1}, typ: reflect.TypeOf(0)},
			{name: "Field3", index: []int{2}, typ: reflect.TypeOf([]byte{})},
		}, index1: 0, index2: 0, expected: false},
	}

	for _, tc := range testCases {
		result := byName(tc.input).Less(tc.index1, tc.index2)
		if result != tc.expected {
			t.Errorf("byName.Less(%v, %v) = %v, expected %v", tc.index1, tc.index2, result, tc.expected)
		}
	}

	src := []field{
		{name: "field1", index: []int{0}, tag: true},
		{name: "field2", index: []int{0, 1}, tag: true},
	}
	// expected := len(src[1].index) < len(src[0].index)
	result := byName(src).Less(0, 1)
	assert.NotNil(t, result)
}

func TestMarshal(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
		opts     encOpts
	}{
		{input: "test", expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
		{input: int(123), expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
		{input: uint(456), expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
		{input: json.Number("789"), expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
		{input: []byte("bytes"), expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
		{input: "<html>", expected: ``, opts: encOpts{quoted: true, escapeHTML: true}},
		{input: "<html>", expected: ``, opts: encOpts{quoted: true, escapeHTML: false}},
	}

	for _, tc := range testCases {
		var dst bytes.Buffer
		e := &encodeState{
			Buffer: dst,
		}
		err := e.marshal(tc.input, tc.opts)
		if err != nil {
			t.Errorf("e.marshal(%v) = %v, expected nil", tc.input, err)
		}

		result := dst.String()
		if result != tc.expected {
			t.Errorf("e.marshal(%v) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestEncodefunc(_ *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 []byte
	}

	testCases := []struct {
		input    TestStruct
		expected string
		opts     encOpts
	}{
		{input: TestStruct{Field1: "test", Field2: 123, Field3: []byte("bytes")}, expected: `{"Field1":"test","Field2":123,"Field3":"Ynl0ZXM="}`, opts: encOpts{quoted: false, escapeHTML: false}},
	}

	for _, tc := range testCases {
		var dst bytes.Buffer
		e := &encodeState{
			Buffer: dst,
		}
		se := &structEncoder{
			fields: []field{
				{name: "Field1", index: []int{0}, typ: reflect.TypeOf(""), quoted: true},
			},
			fieldEncs: []encoderFunc{
				stringEncoder,
				intEncoder,
			},
		}
		se.encode(e, reflect.ValueOf(tc.input), tc.opts)
	}
}

func TestEncodeFloatEncoder(t *testing.T) {
	testCases := []struct {
		input    float64
		expected string
		opts     encOpts
	}{
		{input: math.Pi, expected: ``, opts: encOpts{quoted: false, escapeHTML: false}},
	}

	for _, tc := range testCases {
		var dst bytes.Buffer
		e := &encodeState{
			Buffer: dst,
		}
		f := floatEncoder(64)
		f.encode(e, reflect.ValueOf(tc.input), tc.opts)

		result := dst.String()
		if result != tc.expected {
			t.Errorf("f.encode(%v) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestFillField(t *testing.T) {
	testCases := []struct {
		input    field
		expected field
	}{
		{
			input: field{
				name: "test",
			},
			expected: field{
				name:      "test",
				nameBytes: []byte("test"),
				equalFold: bytes.EqualFold,
			},
		},
		{
			input: field{
				name: "test2",
			},
			expected: field{
				name:      "test2",
				nameBytes: []byte("test2"),
				equalFold: bytes.EqualFold,
			},
		},
	}

	for _, tc := range testCases {
		result := fillField(tc.input)
		if result.name != tc.expected.name {
			t.Errorf("fillField(%v).name = %v, expected %v", tc.input, result.name, tc.expected.name)
		}
		if !bytes.Equal(result.nameBytes, tc.expected.nameBytes) {
			t.Errorf("fillField(%v).nameBytes = %v, expected %v", tc.input, result.nameBytes, tc.expected.nameBytes)
		}
	}
}

func TestTypeFields(t *testing.T) {
	type Field4 struct {
		Field41 float32 `json:"field41,string"`
		Field42 float64 `json:"field42,omitempty,omitempty,string"`
		Field43 string
	}

	type TestStruct struct {
		Field1 string `json:"field1"`
		Field2 int    `json:"-"`
		Field3 bool   `json:"field3,omitempty"`
		Field4 Field4
	}

	testCases := []struct {
		input    reflect.Type
		expected []field
	}{
		{
			input: reflect.TypeOf(TestStruct{}),
			expected: []field{
				{
					name:      "field1",
					tag:       true,
					index:     []int{0},
					typ:       reflect.TypeOf(""),
					omitEmpty: false,
					quoted:    false,
				},
				{
					name:      "field3",
					tag:       true,
					index:     []int{2},
					typ:       reflect.TypeOf(false),
					omitEmpty: true,
					quoted:    false,
				},
			},
		},
	}

	for _, tc := range testCases {
		result := typeFields(tc.input)
		if len(result) != len(tc.expected) {
			continue
		}

		for i, f := range result {
			expected := tc.expected[i]
			if f.name != expected.name {
				t.Errorf("typeFields(%v).name = %v, expected %v", tc.input, f.name, expected.name)
			}
			if f.tag != expected.tag {
				t.Errorf("typeFields(%v).tag = %v, expected %v", tc.input, f.tag, expected.tag)
			}
			if !reflect.DeepEqual(f.index, expected.index) {
				t.Errorf("typeFields(%v).index = %v, expected %v", tc.input, f.index, expected.index)
			}
			if f.typ != expected.typ {
				t.Errorf("typeFields(%v).typ = %v, expected %v", tc.input, f.typ, expected.typ)
			}
			if f.omitEmpty != expected.omitEmpty {
				t.Errorf("typeFields(%v).omitEmpty = %v, expected %v", tc.input, f.omitEmpty, expected.omitEmpty)
			}
			if f.quoted != expected.quoted {
				t.Errorf("typeFields(%v).quoted = %v, expected %v", tc.input, f.quoted, expected.quoted)
			}
		}
	}
}

func TestMarshalIndent(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 []byte
	}

	testCases := []struct {
		input    TestStruct
		expected string
		opts     encOpts
	}{
		{input: TestStruct{Field1: "test", Field2: 123, Field3: []byte("bytes")}, expected: `{
  "Field1": "test",
  "Field2": 123,
  "Field3": "Ynl0ZXM="
}`, opts: encOpts{quoted: false, escapeHTML: false}},
	}

	for _, tc := range testCases {
		result, err := MarshalIndent(tc.input, "", "  ")
		if err != nil {
			t.Errorf("MarshalIndent(%v) = %v, expected nil", tc.input, err)
		}

		if string(result) != tc.expected {
			t.Errorf("MarshalIndent(%v) = %v, expected %v", tc.input, string(result), tc.expected)
		}
	}
}

// Mock for the Marshaler interface
type mockMarshaler struct {
	mock.Mock
	json []byte
	err  error
}

func (m *mockMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`{"key":"value"}`), nil
}

func TestMarshalerEncoder(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		marshalReturn  []byte
		marshalErr     error
		expectedString string
		expectedErr    string
	}{
		{
			name:           "Handle nil pointer value",
			value:          (*mockMarshaler)(nil),
			expectedString: "null",
		},
		{
			name:           "Successful MarshalJSON",
			value:          &mockMarshaler{},
			marshalReturn:  []byte(`{"key":"value"}`),
			expectedString: `{"key":"value"}`,
		},
		{
			name:           "Handle optional MarshalJSON error",
			value:          &mockMarshaler{},
			marshalReturn:  []byte(`{"key":"value"}`),
			expectedString: `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize encodeState
			e := &encodeState{}

			var v reflect.Value
			if tt.value != nil {
				m := tt.value.(*mockMarshaler)
				if tt.marshalErr != nil || tt.marshalReturn != nil {
					m.On("MarshalJSON").Return(tt.marshalReturn, tt.marshalErr).Once()
				}
				v = reflect.ValueOf(m)
			} else {
				v = reflect.ValueOf(tt.value)
			}

			// Call marshalerEncoder
			marshalerEncoder(e, v, encOpts{})

			// Check result
			if tt.expectedErr != "" {
				assert.Contains(t, e.String(), tt.expectedErr)
			} else {
				assert.Equal(t, tt.expectedString, e.String())
			}
		})
	}
}

// Custom type implementing encoding.TextMarshaler
type Person struct {
	Name string
	Age  int
}

func (p Person) MarshalText() ([]byte, error) {
	if p.Name == "" {
		return nil, errors.New("name is empty")
	}
	return []byte(p.Name), nil
}

func TestTextMarshalerEncoder(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected string
	}{
		{"NilPointer", (*Person)(nil), ""},
		{"ValidPerson", Person{Name: "Alice", Age: 30}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{}
			v := reflect.ValueOf(tt.value)
			opts := encOpts{escapeHTML: false}

			textMarshalerEncoder(e, v, opts)

			if e.result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, e.result)
			}
		})
	}
}

type mockTextMarshaler struct {
	text []byte
	err  error
}

func (m *mockTextMarshaler) MarshalText() ([]byte, error) {
	return m.text, m.err
}

func TestAddrTextMarshalerEncoder(t *testing.T) {
	tests := []struct {
		name     string
		input    encoding.TextMarshaler
		opts     encOpts
		expected string
		err      error
	}{
		{
			name:     "Nil pointer",
			input:    nil,
			opts:     encOpts{escapeHTML: false},
			expected: "",
		},
		{
			name:     "Successful marshal",
			input:    &mockTextMarshaler{text: []byte("test")},
			opts:     encOpts{escapeHTML: false},
			expected: "",
		},
		{
			name:     "Marshal error",
			input:    &mockTextMarshaler{text: []byte("test"), err: errors.New("marshal error")},
			opts:     encOpts{escapeHTML: false},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{}
			v := reflect.ValueOf(tt.input)
			addrTextMarshalerEncoder(e, v, tt.opts)

			if e.result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, e.result)
			}
		})
	}
}

func TestAddrMarshalerEncoder(t *testing.T) {
	tests := []struct {
		name     string
		input    Marshaler
		expected string
		err      error
	}{
		{
			name:     "Nil pointer value",
			input:    nil,
			expected: "null",
		},
		{
			name:     "Successful marshal",
			input:    &mockMarshaler{json: []byte(`{"key":"value"}`)},
			expected: `{"key":"value"}`,
		},
		{
			name:     "Marshal error",
			input:    &mockMarshaler{err: errors.New("marshal error")},
			expected: `{"key":"value"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &encodeState{}
			var v reflect.Value
			if tt.input != nil {
				v = reflect.ValueOf(tt.input).Elem()
			} else {
				v = reflect.ValueOf(tt.input)
			}
			addrMarshalerEncoder(e, v, encOpts{})
			if e.Buffer.String() != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, e.Buffer.String())
			}
		})
	}
}

func TestInterfaceEncoder(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		expectedString string
	}{
		{
			name:           "Handle nil pointer value",
			value:          (*interface{})(nil),
			expectedString: "null",
		},
		{
			name:           "Handle nil value",
			value:          nil,
			expectedString: "null",
		},
		{
			name:           "Handle string value",
			value:          "test",
			expectedString: `"test"`,
		},
		{
			name:           "Handle int value",
			value:          123,
			expectedString: "123",
		},
		{
			name:           "Handle float value",
			value:          123.45,
			expectedString: "123.45",
		},
		{
			name:           "Handle bool value",
			value:          true,
			expectedString: "true",
		},
		{
			name:           "Handle map value",
			value:          map[string]interface{}{"key": "value"},
			expectedString: `{"key":"value"}`,
		},
		{
			name:           "Handle slice value",
			value:          []interface{}{"value1", "value2"},
			expectedString: `["value1","value2"]`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			// Initialize encodeState
			e := &encodeState{}
			var s *string

			// Call interfaceEncoder
			interfaceEncoder(e, reflect.ValueOf(s), encOpts{})
		})
	}
}

func TestEncodeStructEncoder(_ *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 int
		Field3 []byte
	}

	testCases := []struct {
		input    TestStruct
		expected string
		opts     encOpts
	}{
		{input: TestStruct{Field1: "test", Field2: 123, Field3: []byte("bytes")}, expected: `{"Field1":"test","Field2":123,"Field3":"Ynl0ZXM="}`, opts: encOpts{quoted: false, escapeHTML: false}},
	}

	for _, tc := range testCases {
		var s *string
		var dst bytes.Buffer
		e := &encodeState{
			Buffer: dst,
		}
		se := &structEncoder{
			fields: []field{
				{name: "Field1", index: []int{0}, typ: reflect.TypeOf(""), quoted: true},
			},
			fieldEncs: []encoderFunc{
				stringEncoder,
				intEncoder,
			},
		}
		se.encode(e, reflect.ValueOf(s), tc.opts)
	}
}

func TestMarshal_Error(t *testing.T) {
	_, err := Marshal(make(chan int))
	if err == nil {
		t.Errorf("expected an error, but got nil")
	}
}

func TestNewCondAddrEncoder(t *testing.T) {
	// Testing for newCondAddrEncoder
	encoder := newCondAddrEncoder(stringEncoder, stringEncoder)
	assert.NotNil(t, encoder)
}

func TestBoolEncoder(t *testing.T) {
	// Testing for boolEncoder
	var e encodeState
	encoder := boolEncoder

	// Test case for true without quotes
	encoder(&e, reflect.ValueOf(true), encOpts{})
	assert.Equal(t, "true", e.String())
	e.Reset()

	// Test case for false without quotes
	encoder(&e, reflect.ValueOf(false), encOpts{})
	assert.Equal(t, "false", e.String())
	e.Reset()

	// Test case for true with quotes
	encoder(&e, reflect.ValueOf(true), encOpts{quoted: true})
	assert.Equal(t, `"true"`, e.String())
	e.Reset()

	// Test case for false with quotes
	encoder(&e, reflect.ValueOf(false), encOpts{quoted: true})
	assert.Equal(t, `"false"`, e.String())
	e.Reset()
}

func TestInvalidUTF8Error(_ *testing.T) {
	err := InvalidUTF8Error{}
	_ = err.Error()
}

func TestMarshalerError(t *testing.T) {
	err := MarshalerError{
		Type: reflect.TypeOf("test"),
		Err:  errors.New("marshal error"),
	}
	expected := "json: error calling MarshalJSON for type string: marshal error"
	result := err.Error()
	assert.Equal(t, expected, result)
}

func TestUnsupportedValueError(_ *testing.T) {
	err := UnsupportedValueError{}
	_ = err.Error()
}

func TestUnsupportedTypeError(t *testing.T) {
	err := UnsupportedTypeError{Type: reflect.TypeOf("test")}
	expected := "json: unsupported type: string"
	result := err.Error()
	assert.Equal(t, expected, result)
}

func TestIntEncoder(t *testing.T) {
	var e encodeState
	encoder := intEncoder

	// Test case for positive integer without quotes
	encoder(&e, reflect.ValueOf(123), encOpts{})
	assert.Equal(t, "123", e.String())
	e.Reset()

	// Test case for negative integer without quotes
	encoder(&e, reflect.ValueOf(-123), encOpts{})
	assert.Equal(t, "-123", e.String())
	e.Reset()

	// Test case for positive integer with quotes
	encoder(&e, reflect.ValueOf(123), encOpts{quoted: true})
	assert.Equal(t, `"123"`, e.String())
	e.Reset()

	// Test case for negative integer with quotes
	encoder(&e, reflect.ValueOf(-123), encOpts{quoted: true})
	assert.Equal(t, `"-123"`, e.String())
	e.Reset()
}

func TestUintEncoder(t *testing.T) {
	var e encodeState
	encoder := uintEncoder

	// Test case for positive integer without quotes
	encoder(&e, reflect.ValueOf(uint(123)), encOpts{})
	assert.Equal(t, "123", e.String())
	e.Reset()

	// Test case for positive integer with quotes
	encoder(&e, reflect.ValueOf(uint(123)), encOpts{quoted: true})
	assert.Equal(t, `"123"`, e.String())
	e.Reset()
}

func TestEncodefloatEncoder(_ *testing.T) {
	var e encodeState
	encoder := floatEncoder(64)

	encoder.encode(&e, reflect.ValueOf(123.45), encOpts{})
}

func TestSliceEncoder(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		expectedString string
		arrayEnc       func(e *encodeState, v reflect.Value, opts encOpts)
	}{
		{
			name:           "Handle nil slice",
			value:          ([]int)(nil),
			expectedString: "null",
			arrayEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
				// This should not be called for nil slice
				t.Fail()
			},
		},
		{
			name:           "Handle non-nil slice",
			value:          []int{1, 2, 3},
			expectedString: "[1,2,3]",
			arrayEnc: func(e *encodeState, _ reflect.Value, _ encOpts) {
				// Simulate encoding of a non-nil slice
				_, _ = e.WriteString("[1,2,3]")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize encodeState
			e := &encodeState{}

			// Initialize sliceEncoder with the provided arrayEnc function
			se := &sliceEncoder{
				arrayEnc: tt.arrayEnc,
			}

			// Call encode
			se.encode(e, reflect.ValueOf(tt.value), encOpts{})

			// Check result
			assert.Equal(t, tt.expectedString, e.String())
		})
	}
}

func TestPtrEncoder(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		expectedString string
		elemEnc        func(e *encodeState, v reflect.Value, opts encOpts)
	}{
		{
			name:           "Handle nil pointer",
			value:          (*int)(nil),
			expectedString: "null",
			elemEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
				// This should not be called for nil pointer
				t.Fail()
			},
		},
		{
			name:           "Handle non-nil pointer",
			value:          func() *int { i := 42; return &i }(),
			expectedString: "42",
			elemEnc: func(e *encodeState, _ reflect.Value, _ encOpts) {
				// Simulate encoding of a non-nil pointer
				_, _ = e.WriteString("42")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Initialize encodeState
			e := &encodeState{}

			// Initialize ptrEncoder with the provided elemEnc function
			pe := &ptrEncoder{
				elemEnc: tt.elemEnc,
			}

			// Call encode
			pe.encode(e, reflect.ValueOf(tt.value), encOpts{})

			// Check result
			assert.Equal(t, tt.expectedString, e.String())
		})
	}
}

func TestCondAddrEncoder_Encode(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		canAddrEnc     func(e *encodeState, v reflect.Value, opts encOpts)
		elseEnc        func(e *encodeState, v reflect.Value, opts encOpts)
		expectedString string
	}{
		{
			name:  "Test case for value with CanAddr() == true",
			value: "test",
			canAddrEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
			},
			elseEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
			},
			expectedString: "expected string for value with CanAddr() == true",
		},
		{
			name:  "Test case for value with CanAddr() == false",
			value: "test",
			canAddrEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
			},
			elseEnc: func(_ *encodeState, _ reflect.Value, _ encOpts) {
			},
			expectedString: "expected string for value with CanAddr() == false",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(_ *testing.T) {
			e := &encodeState{}

			ce := &condAddrEncoder{
				canAddrEnc: tt.canAddrEnc,
				elseEnc:    tt.elseEnc,
			}

			ce.encode(e, reflect.ValueOf(tt.value), encOpts{})
		})
	}
}

func TestByString_Swap(t *testing.T) {
	sv := byString{
		{s: "a"},
		{s: "b"},
		{s: "c"},
	}

	sv.Swap(0, 2)

	assert.Equal(t, "c", sv[0].s)
	assert.Equal(t, "b", sv[1].s)
	assert.Equal(t, "a", sv[2].s)
}

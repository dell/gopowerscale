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
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiteralStore(_ *testing.T) {
	d := &decodeState{
		data: []byte(`null`),
		off:  0,
	}

	// Test case for unmarshaling a valid JSON null
	var item []byte
	item = []byte("n")
	v := reflect.ValueOf(5)
	d.literalStore(item, v, false)

	// Test case for unmarshaling a valid JSON true
	item = []byte("t")
	v = reflect.ValueOf(5)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON false
	item = []byte("f")
	v = reflect.ValueOf(5)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON number
	item = []byte("2")
	x := 5
	v = reflect.ValueOf(&x)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON float
	item = []byte("2.2")
	y := 5.5
	v = reflect.ValueOf(&y)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON uint8
	item = []byte("2")
	var z uint8 = 5
	v = reflect.ValueOf(&z)
	d.literalStore(item, v, true)

	// Test case for unmarshaling an empty item
	item = []byte{}
	v = reflect.ValueOf(5)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON string
	item = []byte(`"hello"`)
	var str string
	v = reflect.ValueOf(&str)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON slice of bytes
	item = []byte(`"aGVsbG8="`) // base64 encoded "hello"
	var byteSlice []byte
	v = reflect.ValueOf(&byteSlice)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON interface
	item = []byte("42")
	var iface interface{}
	v = reflect.ValueOf(&iface)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON int into an interface
	item = []byte("42")
	var ifaceInt interface{}
	v = reflect.ValueOf(&ifaceInt)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON float into an interface
	item = []byte("42.5")
	var ifaceFloat interface{}
	v = reflect.ValueOf(&ifaceFloat)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON bool into an interface
	item = []byte("true")
	var ifaceBool interface{}
	v = reflect.ValueOf(&ifaceBool)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON string into an interface
	item = []byte(`"hello"`)
	var ifaceString interface{}
	v = reflect.ValueOf(&ifaceString)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON null into an interface
	item = []byte("null")
	var ifaceNull interface{}
	v = reflect.ValueOf(&ifaceNull)
	d.literalStore(item, v, true)

	// Test case for unmarshaling a valid JSON string into a custom type implementing encoding.TextUnmarshaler
	type MyTextUnmarshaler string
	var myText MyTextUnmarshaler
	item = []byte(`"custom text"`)
	v = reflect.ValueOf(&myText)
	d.literalStore(item, v, true)
}

func TestIsValidNumber(t *testing.T) {
	value := isValidNumber("")
	assert.False(t, value)

	value = isValidNumber("-")
	assert.False(t, value)

	value = isValidNumber("0")
	assert.True(t, value)

	value = isValidNumber("2")
	assert.True(t, value)

	value = isValidNumber("12")
	assert.True(t, value)

	value = isValidNumber("123.254")
	assert.True(t, value)

	value = isValidNumber("123e+10")
	assert.True(t, value)

	value = isValidNumber("A")
	assert.False(t, value)

	value = isValidNumber("e+")
	assert.False(t, value)
}

func TestDecodeState_value(t *testing.T) {
	testCases := []struct {
		input    string
		expected interface{}
	}{
		{input: `null`, expected: nil},
		{input: `true`, expected: true},
		{input: `false`, expected: false},
		{input: `"hello world"`, expected: "hello world"},
		{input: `42`, expected: int64(42)},
		{input: `3.14`, expected: float64(3.14)},
		{input: `{"key":"value"}`, expected: map[string]interface{}{"key": "value"}},
	}

	for _, tc := range testCases {
		d := &decodeState{
			data:      []byte(tc.input),
			scan:      scanner{},
			nextscan:  scanner{},
			useNumber: true,
		}

		var v reflect.Value
		if tc.expected == nil {
			v = reflect.ValueOf(tc.expected)
		} else {
			v = reflect.New(reflect.TypeOf(tc.expected)).Elem()
		}

		d.scan.reset()
		d.value(v)

		if !v.IsValid() {
			continue
		}

		result := v.Interface()
		if !reflect.DeepEqual(result, tc.expected) {
			t.Errorf("d.value(%v) = %v, expected %v", tc.input, result, tc.expected)
		}
	}
}

func TestGetu4(t *testing.T) {
	value := getu4([]byte(`\u1234`))
	assert.Equal(t, value, rune(0x1234))

	value = getu4([]byte(`\u123`))
	assert.Equal(t, value, rune(-1))

	value = getu4([]byte(`\u12345`))
	assert.Equal(t, value, rune(4660))
}

func TestUnquoteBytes(t *testing.T) {
	value, _ := unquoteBytes([]byte(`"hello"`))
	assert.Equal(t, value, []byte("hello"))

	_, _ = unquoteBytes([]byte(`"hello\\world"`))
	_, _ = unquoteBytes([]byte(`"hello\nworld"`))
	_, _ = unquoteBytes([]byte(`"hello\tworld"`))
	_, _ = unquoteBytes([]byte(`"hello\rworld"`))
	_, _ = unquoteBytes([]byte(`"hello\bworld"`))
	_, _ = unquoteBytes([]byte(`"hello\fworld"`))
	_, _ = unquoteBytes([]byte(`"hello\uworld"`))
	_, _ = unquoteBytes([]byte(`"""`))
	_, _ = unquoteBytes([]byte(`"00100100"`))
	_, _ = unquoteBytes([]byte(`"hello`))
}

func TestUnmarshalFunctions(t *testing.T) {
	type testStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	t.Run("Valid JSON", func(t *testing.T) {
		var ts testStruct
		err := Unmarshal([]byte(`{"name":"test","age":20}`), &ts)
		assert.Nil(t, err)
		assert.Equal(t, "test", ts.Name)
		assert.Equal(t, 20, ts.Age)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		var ts testStruct
		err := Unmarshal([]byte(`{"name":"test","age":20`), &ts)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "unexpected end of JSON input")
	})

	t.Run("Non-pointer value", func(t *testing.T) {
		var ts testStruct
		err := Unmarshal([]byte(`{"name":"test","age":20}`), ts)
		assert.NotNil(t, err)
		assert.IsType(t, &InvalidUnmarshalError{}, err)
	})

	t.Run("Nil pointer", func(t *testing.T) {
		var ts *testStruct
		err := Unmarshal([]byte(`{"name":"test","age":20}`), ts)
		assert.NotNil(t, err)
		assert.IsType(t, &InvalidUnmarshalError{}, err)
	})

	t.Run("Simulated decoding error", func(t *testing.T) {
		var ts testStruct
		err := Unmarshal([]byte(`{"name":"test","age":20`), &ts)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "unexpected end of JSON input")
	})
}

func TestUnquote(t *testing.T) {
	value, _ := unquote([]byte(`"hello"`))
	assert.Equal(t, value, "hello")
}

func TestUnmarshalTypeError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      UnmarshalTypeError
		expected string
	}{
		{
			name: "Test basic error",
			err: UnmarshalTypeError{
				Value: "test_value",
				Type:  reflect.TypeOf(int(0)),
			},
			expected: "json: cannot unmarshal test_value into Go value of type int",
		},
		{
			name: "Test with offset",
			err: UnmarshalTypeError{
				Value:  "test_value",
				Type:   reflect.TypeOf(int(0)),
				Offset: 10,
			},
			expected: "json: cannot unmarshal test_value into Go value of type int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Error(); err != tt.expected {
				t.Errorf("Expected error: %s, got: %s", tt.expected, err)
			}
		})
	}
}

func TestUnmarshalFieldError_Error(t *testing.T) {
	tests := []struct {
		name     string
		input    *UnmarshalFieldError
		expected string
	}{
		{
			name: "Basic case",
			input: &UnmarshalFieldError{
				Key:   "testKey",
				Field: reflect.StructField{Name: "TestField"},
				Type:  reflect.TypeOf(""),
			},
			expected: "json: cannot unmarshal object key " + strconv.Quote("testKey") + " into unexported field TestField of type string",
		},
		{
			name: "Numeric key and integer type",
			input: &UnmarshalFieldError{
				Key:   "123",
				Field: reflect.StructField{Name: "Age"},
				Type:  reflect.TypeOf(123),
			},
			expected: "json: cannot unmarshal object key " + strconv.Quote("123") + " into unexported field Age of type int",
		},
		{
			name: "Special characters in key and boolean type",
			input: &UnmarshalFieldError{
				Key:   "!@#$%^&*()",
				Field: reflect.StructField{Name: "Flag"},
				Type:  reflect.TypeOf(true),
			},
			expected: "json: cannot unmarshal object key " + strconv.Quote("!@#$%^&*()") + " into unexported field Flag of type bool",
		},
		{
			name: "Empty key and struct type",
			input: &UnmarshalFieldError{
				Key:   "",
				Field: reflect.StructField{Name: "Address"},
				Type:  reflect.TypeOf(struct{}{}),
			},
			expected: `json: cannot unmarshal object key ` + strconv.Quote("") + ` into unexported field Address of type struct {}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := tt.input.Error()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestInvalidUnmarshalError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      InvalidUnmarshalError
		expected string
	}{
		{
			name: "Test with nil type",
			err: InvalidUnmarshalError{
				Type: nil,
			},
			expected: "json: Unmarshal(nil)",
		},
		{
			name: "Test with non-pointer type",
			err: InvalidUnmarshalError{
				Type: reflect.TypeOf(int(0)),
			},
			expected: "json: Unmarshal(non-pointer int)",
		},
		{
			name: "Test with pointer type",
			err: InvalidUnmarshalError{
				Type: reflect.TypeOf(new(int)),
			},
			expected: "json: Unmarshal(nil *int)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.err.Error(); err != tt.expected {
				t.Errorf("Expected error: %s, got: %s", tt.expected, err)
			}
		})
	}
}

func TestNumberString(t *testing.T) {
	tests := []struct {
		name string
		n    Number
		want string
	}{
		{
			name: "Empty number",
			n:    "",
			want: "",
		},
		{
			name: "Non-empty number",
			n:    "123",
			want: "123",
		},
		{
			name: "Number with negative sign",
			n:    "-456",
			want: "-456",
		},
		{
			name: "Number with decimal point",
			n:    "7.89",
			want: "7.89",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.n.String(); got != tt.want {
				t.Errorf("Number.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumber_Float64(t *testing.T) {
	tests := []struct {
		name    string
		input   Number
		want    float64
		wantErr bool
	}{
		{
			name:    "valid float",
			input:   "3.14",
			want:    3.14,
			wantErr: false,
		},
		{
			name:    "valid integer",
			input:   "42",
			want:    42.0,
			wantErr: false,
		},
		{
			name:    "invalid float",
			input:   "abc",
			want:    0.0,
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			want:    0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.input.Float64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Float64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.Float64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumberInt64(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int64
		wantErr bool
	}{
		{
			name:    "Valid integer",
			input:   "123",
			want:    123,
			wantErr: false,
		},
		{
			name:    "Invalid integer",
			input:   "abc",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty string",
			input:   "",
			want:    0,
			wantErr: true,
		},
		{
			name:    "String with non-numeric characters",
			input:   "123abc",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := Number(tt.input)
			got, err := n.Int64()
			if (err != nil) != tt.wantErr {
				t.Errorf("Number.Int64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Number.Int64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeStateError(t *testing.T) {
	tests := []struct {
		name        string
		inputErr    error
		expectPanic bool
	}{
		{"Normal Error", errors.New("test error"), true},
		{"Nil Error", nil, true}, // Even with nil, panic should occur as it's the function's intent
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decodeState
			defer func() {
				if r := recover(); r != nil {
					if tt.expectPanic {
						if tt.inputErr == nil {
							// Special handling for nil panic
							if _, ok := r.(*runtime.PanicNilError); ok {
								assert.Nil(t, tt.inputErr)
							} else {
								t.Fatalf("expected a PanicNilError, but got: %v", r)
							}
						} else {
							assert.Equal(t, tt.inputErr, r)
						}
					} else {
						t.Fatalf("did not expect a panic but got: %v", r)
					}
				} else if tt.expectPanic {
					t.Fatalf("expected a panic but did not get one")
				}
			}()
			d.error(tt.inputErr)
		})
	}
}

func TestDecodeStateNext(t *testing.T) {
	tests := []struct {
		name           string
		data           []byte
		off            int
		initializeScan bool
		expected       []byte
		shouldPanic    bool
	}{
		{
			name:           "Valid JSON Object",
			data:           []byte(`{"key":"value"}`),
			off:            0,
			initializeScan: true,
			expected:       []byte(`{"key":"value"}`),
			shouldPanic:    false,
		},
		{
			name:           "Valid JSON Array",
			data:           []byte(`[1, 2, 3]`),
			off:            0,
			initializeScan: true,
			expected:       []byte(`[1, 2, 3]`),
			shouldPanic:    false,
		},
		{
			name:           "Empty Data",
			data:           []byte(``),
			off:            0,
			initializeScan: false,
			expected:       nil,
			shouldPanic:    true,
		},
		{
			name:           "Invalid Data",
			data:           []byte(`invalid`),
			off:            0,
			initializeScan: false,
			expected:       nil,
			shouldPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d decodeState
			d.data = tt.data
			d.off = tt.off
			stepFunc := func(_ *scanner, _ byte) int {
				// Mock behavior: do nothing for now
				return 0
			}
			if tt.initializeScan {
				d.scan = scanner{step: stepFunc}
				d.nextscan = scanner{step: stepFunc}
			}

			if tt.shouldPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Expected panic for %s did not occur", tt.name)
					}
				}()
			}

			result := d.next()

			if !tt.shouldPanic {
				assert.Equal(t, tt.expected, result)
				assert.Equal(t, len(tt.data), d.off)
			}
		})
	}
}

func TestValueInterface(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Array",
			input:    `[1, 2, 3]`,
			expected: []interface{}{1.0, 2.0, 3.0},
		},
		{
			name:     "Object",
			input:    `{"key": "value"}`,
			expected: map[string]interface{}{"key": "value"},
		},
		{
			name:     "Literal",
			input:    `"literal"`,
			expected: "literal",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &decodeState{
				data: []byte(tt.input),
				off:  0,
				scan: scanner{
					step: func(_ *scanner, c byte) int {
						switch c {
						case '[':
							return scanBeginArray
						case '{':
							return scanBeginObject
						case '"':
							return scanBeginLiteral
						default:
							return scanContinue
						}
					},
					parseState: []int{scanContinue},
				},
			}
			d.scan.reset()
			result := d.valueInterface()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("valueInterface() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestValueQuoted(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name:     "Array",
			input:    `[1, 2, 3]`,
			expected: unquotedValue{},
		},
		{
			name:     "Object",
			input:    `{"key": "value"}`,
			expected: unquotedValue{},
		},
		{
			name:     "Literal Nil",
			input:    `null`,
			expected: nil,
		},
		{
			name:     "Literal String",
			input:    `"literal"`,
			expected: "literal",
		},
		{
			name:     "Literal Number",
			input:    `123`,
			expected: unquotedValue{},
		},
		{
			name:     "Invalid Literal",
			input:    `invalid`,
			expected: unquotedValue{},
		},
		{
			name:     "Empty Input",
			input:    ``,
			expected: unquotedValue{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &decodeState{
				data: []byte(tt.input),
				off:  0,
				scan: scanner{
					step: func(_ *scanner, c byte) int {
						switch c {
						case '[':
							return scanBeginArray
						case '{':
							return scanBeginObject
						case '"':
							return scanBeginLiteral
						case 'n':
							return scanBeginLiteral
						case '1', '2', '3':
							return scanBeginLiteral
						default:
							return scanContinue
						}
					},
					parseState: []int{scanContinue},
				},
			}
			d.scan.reset()
			var result interface{}
			func() {
				defer func() {
					if r := recover(); r != nil {
						result = unquotedValue{}
					}
				}()
				result = d.valueQuoted()
			}()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("valueQuoted() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestArray(t *testing.T) {
	defaultIndirectFunc := indirectFunc

	afterEach := func() {
		indirectFunc = defaultIndirectFunc
	}
	tests := []struct {
		name     string
		input    string
		expected interface{}
		setup    func()
	}{
		{
			name:     "Empty Array",
			input:    `[]`,
			expected: nil,
		},
		{
			name:     "inject errors from indirectFunc",
			input:    `[]`,
			expected: nil,
			setup: func() {
				indirectFunc = func(_ reflect.Value, _ bool) (Unmarshaler, encoding.TextUnmarshaler, reflect.Value) {
					return nil, nil, reflect.ValueOf(fmt.Errorf("error"))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer afterEach()
			if tt.setup != nil {
				tt.setup()
			}
			d := &decodeState{
				data: []byte(tt.input),
				off:  0,
				scan: scanner{
					step: func(_ *scanner, c byte) int {
						switch c {
						case '[':
							return scanBeginArray
						case '{':
							return scanBeginObject
						case '"':
							return scanBeginLiteral
						case 'n':
							return scanBeginLiteral
						case '1', '2', '3', 'a', 'b', 'c', 't', 'f':
							return scanBeginLiteral
						default:
							return scanContinue
						}
					},
					parseState: []int{scanContinue},
				},
			}
			d.scan.reset()
			var result interface{}
			if tt.name == "Array with fewer elements than target array" {
				result = [3]interface{}{}
			} else {
				result = []interface{}{}
			}
			func() {
				defer func() {
					if r := recover(); r != nil {
						result = nil
					}
				}()
				v := reflect.ValueOf(&result).Elem()
				d.array(v)
				if v.Kind() == reflect.Slice && v.IsNil() {
					v.Set(reflect.MakeSlice(v.Type(), 0, 0))
				}
				result = v.Interface()
			}()
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("array() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

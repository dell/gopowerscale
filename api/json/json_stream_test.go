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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToken(t *testing.T) {
	dec := &Decoder{
		r:          bytes.NewBufferString(`[1,2,3]`),
		tokenState: tokenTopValue,
		tokenStack: []int{},
	}
	token, err := dec.Token()
	assert.NoError(t, err)
	assert.Equal(t, Delim('['), token)

	dec.r = bytes.NewBufferString(`{"key1": "value1", "key2": "value2"}`)
	dec.tokenState = tokenTopValue
	dec.tokenStack = []int{}
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec.r = bytes.NewBufferString(`"value"`)
	dec.tokenState = tokenArrayComma
	dec.tokenStack = []int{}
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec.tokenState = tokenObjectComma
	_, err = dec.Token()
	assert.NotEqual(t, nil, err)

	dec.tokenState = tokenObjectColon
	_, err = dec.Token()
	assert.NotEqual(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`[1,2,3]`),
		tokenState: tokenObjectComma,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Equal(t, tokenObjectComma, dec.tokenState)

	dec = &Decoder{
		r:          bytes.NewBufferString(`]1,2,3]`),
		tokenState: tokenObjectComma,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Equal(t, tokenObjectComma, dec.tokenState)

	dec.tokenState = tokenArrayStart
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`{1,2,3}`),
		tokenState: tokenObjectComma,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Equal(t, tokenObjectComma, dec.tokenState)

	dec.tokenState = tokenArrayComma
	_, err = dec.Token()
	assert.NotEqual(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`}1,2,3}`),
		tokenState: tokenObjectComma,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec.tokenState = tokenArrayStart
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`"1,2,3"`),
		tokenState: tokenObjectStart,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`:1,2,3"`),
		tokenState: tokenObjectStart,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.NotEqual(t, nil, err)

	dec.tokenState = tokenObjectColon
	_, err = dec.Token()
	assert.Equal(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`{1,2,3}`),
		tokenState: tokenTopValue,
	}
	token, err = dec.Token()
	assert.NoError(t, err)
	assert.Equal(t, Delim('{'), token)

	dec = &Decoder{
		r:          bytes.NewBufferString(`}1,2,3}`),
		tokenState: tokenObjectKey,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Error(t, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`,2,3]`),
		tokenState: tokenObjectValue,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	assert.Error(t, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`"1,2,3"`),
		tokenState: tokenObjectStart,
		tokenStack: []int{tokenTopValue},
	}
	_, err = dec.Token()
	err1 := dec.Decode(nil)
	assert.Nil(t, err)
	assert.Error(t, err1)
}

func TestTokenPrepareForDecode(t *testing.T) {
	dec := &Decoder{
		r:          bytes.NewBufferString(`[1,2,3]`),
		tokenState: tokenArrayComma,
		tokenStack: []int{},
	}
	err := dec.tokenPrepareForDecode()
	assert.NotEqual(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`,2,3]`),
		tokenState: tokenArrayComma,
		tokenStack: []int{},
	}
	err = dec.tokenPrepareForDecode()
	assert.Equal(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`[1,2,3]`),
		tokenState: tokenObjectColon,
		tokenStack: []int{},
	}
	err = dec.tokenPrepareForDecode()
	assert.NotEqual(t, nil, err)

	dec = &Decoder{
		r:          bytes.NewBufferString(`:1,2,3]`),
		tokenState: tokenObjectColon,
		tokenStack: []int{},
	}
	err = dec.tokenPrepareForDecode()
	assert.Equal(t, nil, err)
}

func TestTokenError(t *testing.T) {
	dec := &Decoder{
		r:          bytes.NewBufferString(`[1,2,3]`),
		tokenState: tokenObjectValue,
		tokenStack: []int{},
	}
	var c byte = '5'
	token, _ := dec.tokenError(c)
	assert.Equal(t, nil, token)

	dec.tokenState = tokenArrayComma
	token, _ = dec.tokenError(c)
	assert.Equal(t, nil, token)

	dec.tokenState = tokenObjectKey
	token, _ = dec.tokenError(c)
	assert.Equal(t, nil, token)

	dec.tokenState = tokenObjectColon
	token, _ = dec.tokenError(c)
	assert.Equal(t, nil, token)

	dec.tokenState = tokenObjectComma
	token, _ = dec.tokenError(c)
	assert.Equal(t, nil, token)

	dec.tokenState = tokenTopValue
	token, _ = dec.tokenError(c)
	assert.Equal(t, nil, token)
}

func TestEncode(t *testing.T) {
	// Test case for encoding a simple JSON object
	enc := &Encoder{
		w:          bytes.NewBuffer(nil),
		escapeHTML: false,
	}
	err := enc.Encode(map[string]interface{}{
		"key1": "value1",
		"key2": "value2",
	})
	assert.NoError(t, err)
	expected := `{"key1":"value1","key2":"value2"}
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode(map[string]interface{}{
		"key1": map[string]interface{}{
			"nestedKey1": "nestedValue1",
			"nestedKey2": "nestedValue2",
		},
		"key2": "value2",
	})
	assert.NoError(t, err)
	expected = `{"key1":{"nestedKey1":"nestedValue1","nestedKey2":"nestedValue2"},"key2":"value2"}
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode([]interface{}{
		map[string]interface{}{
			"key1": "value1",
			"key2": "value2",
		},
		map[string]interface{}{
			"key3": "value3",
			"key4": "value4",
		},
	})
	assert.NoError(t, err)
	expected = `[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	// Test case for encoding a string
	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode("test string")
	assert.NoError(t, err)
	expected = `"test string"
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	// Test case for encoding a number
	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode(123)
	assert.NoError(t, err)
	expected = `123
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	// Test case for encoding a boolean
	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode(true)
	assert.NoError(t, err)
	expected = `true
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	// Test case for encoding a nil value
	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode(nil)
	assert.NoError(t, err)
	expected = `null
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())

	// Test case for encoding an error
	enc.w.(*bytes.Buffer).Reset()
	err = enc.Encode(errors.New("test error"))

	enc = &Encoder{
		w:            bytes.NewBuffer(nil),
		escapeHTML:   false,
		indentPrefix: "prefix",
	}
	err = enc.Encode(map[string]interface{}{
		"key1": "value1",
	})
	assert.NoError(t, err)
	expected = `{
prefix"key1": "value1"
prefix}
`
	assert.Equal(t, expected, enc.w.(*bytes.Buffer).String())
}

func TestNonSpace(t *testing.T) {
	var b byte = 'a'
	value := nonSpace([]byte{b})
	assert.Equal(t, true, value)
}

func TestSetIndent(_ *testing.T) {
	enc := &Encoder{
		w:          bytes.NewBuffer(nil),
		escapeHTML: false,
	}
	enc.SetIndent("", "")
}

func TestContains(t *testing.T) {
	var tag tagOptions
	value := tag.Contains("")
	assert.Equal(t, false, value)

	tag = "string"
	value = tag.Contains("string")
	assert.Equal(t, true, value)

	tag = "string,string1,string2"
	value = tag.Contains("string1")
	assert.Equal(t, true, value)
}

func TestParseTag(t *testing.T) {
	value, _ := parseTag("string,string1")
	assert.Equal(t, "string", value)

	value, _ = parseTag("string")
	assert.Equal(t, "string", value)
}

func TestStringMethod(t *testing.T) {
	d := Delim('a')
	value := d.String()
	assert.Equal(t, "a", value)
}

func TestBuffered(t *testing.T) {
	b := Decoder{}
	value := b.Buffered()
	assert.NotNil(t, value)
}

func TestNewDecoder(t *testing.T) {
	r := bytes.NewBufferString(`[1,2,3]`)
	dec := NewDecoder(r)
	assert.NotNil(t, dec)
}

func TestNewEncoder(t *testing.T) {
	w := bytes.NewBuffer(nil)
	enc := NewEncoder(w)
	assert.NotNil(t, enc)
}

func TestSetEscapeHTML(t *testing.T) {
	enc := &Encoder{
		w:          bytes.NewBuffer(nil),
		escapeHTML: false,
	}
	enc.SetEscapeHTML(true)
	assert.Equal(t, true, enc.escapeHTML)
}

func TestMarshalJSON(t *testing.T) {
	m := RawMessage{}
	_, err := m.MarshalJSON()
	assert.Equal(t, nil, err)
}

func TestMore(t *testing.T) {
	dec := &Decoder{
		r: bytes.NewBufferString(`[1,2,3]`),
	}
	value := dec.More()
	assert.Equal(t, true, value)
}

func TestUseNumber(_ *testing.T) {
	dec := &Decoder{
		r: bytes.NewBufferString(`[1,2,3]`),
	}
	dec.UseNumber()
}

func TestRawMessageUnmarshalJSON(t *testing.T) {
	tests := []struct {
		name          string
		rawMessage    *RawMessage
		data          []byte
		expected      RawMessage
		expectedError error
	}{
		{
			name:          "Nil RawMessage",
			rawMessage:    nil,
			data:          []byte(`{"key": "value"}`),
			expected:      nil,
			expectedError: errors.New("json.RawMessage: UnmarshalJSON on nil pointer"),
		},
		{
			name:          "Valid RawMessage",
			rawMessage:    new(RawMessage),
			data:          []byte(`{"key": "value"}`),
			expected:      RawMessage(`{"key": "value"}`),
			expectedError: nil,
		},
		{
			name:          "Empty Data",
			rawMessage:    new(RawMessage),
			data:          []byte(``),
			expected:      RawMessage(nil),
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rawMessage.UnmarshalJSON(tt.data)
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, *tt.rawMessage)
			}
		})
	}
}

func TestClearOffset(t *testing.T) {
	tests := []struct {
		name           string
		inputError     error
		expectedMsg    string
		expectedOffset int64
	}{
		{
			name:           "Clear Offset for SyntaxError",
			inputError:     &SyntaxError{"test error", 100},
			expectedMsg:    "test error",
			expectedOffset: 0,
		},
		{
			name:           "Non-SyntaxError",
			inputError:     errors.New("non syntax error"),
			expectedMsg:    "non syntax error",
			expectedOffset: 100, // original offset should remain unchanged
		},
		{
			name:           "Nil Error",
			inputError:     nil,
			expectedMsg:    "",
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.inputError == nil {
				assert.NotPanics(t, func() { clearOffset(tt.inputError) })
			} else {
				clearOffset(tt.inputError)
				if syntaxErr, ok := tt.inputError.(*SyntaxError); ok {
					assert.Equal(t, tt.expectedOffset, syntaxErr.Offset)
					assert.Equal(t, tt.expectedMsg, syntaxErr.Error())
				} else {
					assert.EqualError(t, tt.inputError, tt.expectedMsg)
				}
			}
		})
	}
}

func TestDecodeError(t *testing.T) {
	dec := &Decoder{
		buf: []byte("dummy"),
	}
	dec.err = errors.New("json: error")
	expected := "json: error"
	err := dec.Decode(nil)
	assert.Error(t, err)
	assert.EqualError(t, err, expected)
}

func TestReadValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    int
		wantErr bool
	}{
		{
			name:    "Valid JSON Object",
			input:   `{"key":"value"}`,
			want:    15,
			wantErr: false,
		},
		{
			name:    "Valid JSON Array",
			input:   `[1, 2, 3]`,
			want:    9,
			wantErr: false,
		},
		{
			name:    "Invalid Data",
			input:   "[invalid](cci:1://invalid)",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dec := &Decoder{
				buf: []byte(tt.input),
			}

			got, err := dec.readValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("readValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("readValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

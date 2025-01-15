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

func TestUnmarshalJSON(t *testing.T) {
	m := RawMessage{}
	err := m.UnmarshalJSON([]byte(`{"key1":"value1","key2":"value2"}`))
	assert.Equal(t, nil, err)

	m = RawMessage(`key1`)
	err = m.UnmarshalJSON([]byte(`{"key1":"value1","key2":"value2"}`))
	assert.Equal(t, nil, err)
}

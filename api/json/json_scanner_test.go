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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateBeginValue(t *testing.T) {
	s := &scanner{step: stateBeginValue}
	var c byte = '-'
	v := stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = '0'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = 't'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = 'n'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = 'f'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = '5'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	c = '!'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanError, v)

	s.reset()
	c = '!'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanError, v)
}

func TestNextValue(t *testing.T) {
	// Test case for valid JSON object
	data := []byte(`{"key1":"value1","key2":"value2"}`)
	scan := &scanner{step: stateBeginValue}
	value, rest, err := nextValue(data, scan)
	assert.NoError(t, err)
	expectedValue := []byte(`{"key1":"value1","key2":"value2"}`)
	expectedRest := []byte{}
	assert.Equal(t, expectedValue, value)
	assert.Equal(t, expectedRest, rest)

	// Test case for valid JSON array
	data = []byte(`[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]`)
	scan.reset()
	value, rest, err = nextValue(data, scan)
	assert.NoError(t, err)
	expectedValue = []byte(`[{"key1":"value1","key2":"value2"},{"key3":"value3","key4":"value4"}]`)
	expectedRest = []byte{}
	assert.Equal(t, expectedValue, value)
	assert.Equal(t, expectedRest, rest)

	// Test case for invalid JSON
	data = []byte(`{"key1":"value1","key2":"value2"`)
	scan.reset()
	value, rest, err = nextValue(data, scan)
	assert.Error(t, err)
	assert.Equal(t, []byte(nil), value)
	assert.Equal(t, []byte(nil), rest)

	// Test case for empty JSON
	data = []byte(`{}`)
	scan.reset()
	value, rest, err = nextValue(data, scan)
	assert.NoError(t, err)
	expectedValue = []byte(`{}`)
	expectedRest = []byte{}
	assert.Equal(t, expectedValue, value)
	assert.Equal(t, expectedRest, rest)

	data = []byte(`a`)
	value, rest, err = nextValue(data, scan)
	assert.Error(t, err)
	expectedValue = []byte(nil)
	expectedRest = []byte(nil)
	assert.Equal(t, expectedValue, value)
	assert.Equal(t, expectedRest, rest)
}

func TestStateInStringEsc(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = 'b'
	v := stateInStringEsc(s, c)
	assert.Equal(t, scanContinue, v)

	// Test case for 'u' (unicode escape)
	s.reset()
	c = 'u'
	v = stateInStringEsc(s, c)
	assert.Equal(t, scanContinue, v)

	// Test case for invalid character
	s.reset()
	c = '!'
	v = stateInStringEsc(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateInStringEscU(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '8'
	v := stateInStringEscU(s, c)
	assert.Equal(t, scanContinue, v)

	v = stateInStringEscU1(s, c)
	assert.Equal(t, scanContinue, v)

	v = stateInStringEscU12(s, c)
	assert.Equal(t, scanContinue, v)

	v = stateInStringEscU123(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateInStringEscU(s, c)
	assert.Equal(t, scanError, v)

	v = stateInStringEscU1(s, c)
	assert.Equal(t, scanError, v)

	v = stateInStringEscU12(s, c)
	assert.Equal(t, scanError, v)

	v = stateInStringEscU123(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateNeg(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '0'
	v := stateNeg(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '5'
	v = stateNeg(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateNeg(s, c)
	assert.Equal(t, scanError, v)
}

func TestState1(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '1'
	v := state1(s, c)
	assert.Equal(t, scanContinue, v)
}

func TestState0(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '.'
	v := state0(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = 'e'
	v = state0(s, c)
	assert.Equal(t, scanContinue, v)
}

func TestStateDot(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '5'
	v := stateDot(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateDot(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateDot0(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '5'
	v := stateDot0(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = 'e'
	v = stateDot0(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateDot0(s, c)
	assert.Equal(t, 10, v)
}

func TestStateE(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '+'
	v := stateE(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateE(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateESign(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '5'
	v := stateESign(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateESign(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateE0(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '5'
	v := stateE0(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateE0(s, c)
	assert.Equal(t, 10, v)
}

func TestStateT(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = 'r'
	v := stateT(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateT(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateTr(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = 'u'
	v := stateTr(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateTr(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateTru(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = 'e'
	v := stateTru(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateTru(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateFa(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'l'
	v := stateFa(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateFa(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateFal(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 's'
	v := stateFal(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateFal(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateFals(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'e'
	v := stateFals(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateFals(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateN(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'u'
	v := stateN(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateN(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateNul(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'l'
	v := stateNul(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateNul(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateNu(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'l'
	v := stateNu(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateNu(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateInString(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = '\\'
	v := stateInString(s, c)
	assert.Equal(t, scanContinue, v)
}

func TestStateF(t *testing.T) {
	s := &scanner{step: stateInStringEsc}
	var c byte = 'a'
	v := stateF(s, c)
	assert.Equal(t, scanContinue, v)

	s.reset()
	c = '!'
	v = stateF(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateBeginValueOrEmpty(t *testing.T) {
	s := &scanner{step: stateBeginValue}
	var c byte = ' '
	v := stateBeginValueOrEmpty(s, c)
	assert.Equal(t, scanSkipSpace, v)

	s.reset()
	s = &scanner{step: stateBeginValue}
	c = ']'
	v = stateBeginValueOrEmpty(s, c)
	assert.Equal(t, scanEnd, v)
}

func TestStateBeginString(t *testing.T) {
	s := &scanner{step: stateBeginString}
	var c byte = ' '
	v := stateBeginString(s, c)
	assert.Equal(t, scanSkipSpace, v)

	s.reset()
	s = &scanner{step: stateBeginString}
	c = '"'
	v = stateBeginString(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	s.reset()
	s = &scanner{step: stateBeginString}
	c = '['
	v = stateBeginString(s, c)
	assert.Equal(t, scanError, v)
}

func TestStateError(t *testing.T) {
	s := &scanner{}
	result := stateError(s, 'a')
	assert.Equal(t, scanError, result, "Expected stateError to return scanError")
}

func TestQuoteChar(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected string
	}{
		{
			name:     "Single quote",
			input:    '\'',
			expected: `'\''`,
		},
		{
			name:     "Double quote",
			input:    '"',
			expected: `'"'`,
		},
		{
			name:     "Alphanumeric character",
			input:    'a',
			expected: `'a'`,
		},
		{
			name:     "Digit",
			input:    '1',
			expected: `'1'`,
		},
		{
			name:     "Special character",
			input:    '\\',
			expected: `'\\'`,
		},
		{
			name:     "Whitespace",
			input:    ' ',
			expected: `' '`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := quoteChar(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

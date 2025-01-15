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

	// Test case for zero
	s.reset()
	c = '0'
	v = stateBeginValue(s, c)
	assert.Equal(t, scanBeginLiteral, v)

	// // Test case for 't' (true)
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
}

func TestStateInStringEsc(t *testing.T) {
	// Test case for backspace
	s := &scanner{step: stateInStringEsc}
	var c byte = 'b'
	v := stateInStringEsc(s, c)
	assert.Equal(t, scanContinue, v)
	// assert.Equal(t, stateInString, s.step)

	// Test case for 'u' (unicode escape)
	s.reset()
	c = 'u'
	v = stateInStringEsc(s, c)
	assert.Equal(t, scanContinue, v)
	// assert.Equal(t, stateInStringEscU, s.step)

	// Test case for invalid character
	s.reset()
	c = '!'
	v = stateInStringEsc(s, c)
	assert.Equal(t, scanError, v)
	// assert.Equal(t, stateInStringEsc, s.step)
	// assert.Error(t, s.err)
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

	// s.reset()
	// c = '!'
	// v = state1(s, c)
	// assert.Equal(t, scanError, v)
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

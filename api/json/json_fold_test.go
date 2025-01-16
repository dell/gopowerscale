package json

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFoldFunc(t *testing.T) {

	s := []byte("abc")
	_ = foldFunc(s)

	s = []byte("K")
	_ = foldFunc(s)
}

func TestEqualFoldRight(t *testing.T) {
	s := []byte("abc")
	t1 := []byte("abc")
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("aKs")
	t1 = []byte("aks")
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("123")
	t1 = []byte("123")
	assert.True(t, equalFoldRight(s, t1))
}

func TestASCIIEqualFold(t *testing.T) {
	s := []byte("abc")
	t1 := []byte("abc")
	assert.True(t, asciiEqualFold(s, t1))

	s = []byte("abc")
	t1 = []byte("ZW")
	assert.False(t, asciiEqualFold(s, t1))

	s = []byte("123")
	t1 = []byte("123")
	assert.True(t, asciiEqualFold(s, t1))

	t1 = []byte("1234")
	assert.False(t, asciiEqualFold(s, t1))
}

func TestSimpleLetterEqualFold(t *testing.T) {
	s := []byte("abc")
	t1 := []byte("abc")
	assert.True(t, simpleLetterEqualFold(s, t1))

	s = []byte("abc")
	t1 = []byte("abcd")
	assert.False(t, simpleLetterEqualFold(s, t1))
}

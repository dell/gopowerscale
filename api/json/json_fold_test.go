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

func TestFoldFunc(_ *testing.T) {
	s := []byte("abc")
	_ = foldFunc(s)

	s = []byte("K")
	_ = foldFunc(s)

	s = []byte("123")
	_ = foldFunc(s)

	s = []byte("abc\x80")
	_ = foldFunc(s)

	s = []byte("aKs")
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

	s = []byte("123")
	t1 = []byte("")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("123")
	t1 = []byte("abc")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("")
	t1 = []byte("abc")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("s")
	t1 = []byte("s")
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("k")
	t1 = []byte("k")
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("s")
	t1 = []byte(string(smallLongEss))
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("S")
	t1 = []byte(string(smallLongEss))
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("k")
	t1 = []byte(string(kelvin))
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("K")
	t1 = []byte(string(kelvin))
	assert.True(t, equalFoldRight(s, t1))

	s = []byte("s")
	t1 = []byte("k")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("k")
	t1 = []byte("s")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("s")
	t1 = []byte("x")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("S")
	t1 = []byte("x")
	assert.False(t, equalFoldRight(s, t1))

	s = []byte("x")
	t1 = []byte{0xC3, 0xA9}
	assert.False(t, equalFoldRight(s, t1))
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

	s = []byte("abc")
	t1 = []byte("ZWR")
	assert.False(t, asciiEqualFold(s, t1))

	s = []byte("123")
	t1 = []byte("abc")
	assert.False(t, asciiEqualFold(s, t1))
}

func TestSimpleLetterEqualFold(t *testing.T) {
	s := []byte("abc")
	t1 := []byte("abc")
	assert.True(t, simpleLetterEqualFold(s, t1))

	s = []byte("abc")
	t1 = []byte("abcd")
	assert.False(t, simpleLetterEqualFold(s, t1))

	s = []byte("abc")
	t1 = []byte("abX")
	assert.False(t, simpleLetterEqualFold(s, t1))
}

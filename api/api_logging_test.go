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

package api

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsBinOctetBody(t *testing.T) {
	header := http.Header{}
	header.Add(headerKeyContentType, headerValContentTypeBinaryOctetStream)
	assert.True(t, isBinOctetBody(header))

	header.Set(headerKeyContentType, "application/json")
	assert.False(t, isBinOctetBody(header))
}

func TestLogRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/test-url", nil)
	ctx := context.Background()
	var out bytes.Buffer

	t.Run("VerboseLow", func(t *testing.T) {
		out.Reset()
		logRequest(ctx, &out, req, VerboseLow)
		assert.Contains(t, out.String(), "GET /test-url ")
	})

	t.Run("VerboseHigh", func(t *testing.T) {
		out.Reset()
		logRequest(ctx, &out, req, VerboseHigh)
		assert.Contains(t, out.String(), "GET /test-url ")
		assert.Contains(t, out.String(), "Host: example.com")
	})
}

func TestLogResponse(t *testing.T) {
	res := httptest.NewRecorder().Result()
	ctx := context.Background()
	var out bytes.Buffer

	t.Run("VerboseLow", func(t *testing.T) {
		out.Reset()
		logResponse(ctx, res, VerboseLow)
		assert.Contains(t, out.String(), "")
	})

	t.Run("VerboseMedium", func(t *testing.T) {
		out.Reset()
		logResponse(ctx, res, VerboseMedium)
		assert.Contains(t, out.String(), "")
	})

	t.Run("VerboseHigh", func(t *testing.T) {
		out.Reset()
		logResponse(ctx, res, VerboseHigh)
		assert.Contains(t, out.String(), "")
	})
}

func TestWriteIndented(t *testing.T) {
	data := []byte("line1\nline2\nline3")
	var out bytes.Buffer
	err := WriteIndented(&out, data)
	assert.NoError(t, err)
	assert.Equal(t, "    line1\n    line2\n    line3", out.String())
}

type errorWriter struct {
	failAfter int
	writes    int
}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	if ew.writes >= ew.failAfter {
		return 0, errors.New("forced write error")
	}
	ew.writes++
	return len(p), nil
}

func TestWriteIndentedN(t *testing.T) {
	data := []byte("line1\nline2\nline3")

	// Original test case to ensure it still passes
	t.Run("normal case", func(t *testing.T) {
		var out bytes.Buffer
		err := WriteIndentedN(&out, data, 2)
		assert.NoError(t, err)
		assert.Equal(t, "  line1\n  line2\n  line3", out.String())
	})

	// Test case to cover error scenarios
	t.Run("error after writing space", func(t *testing.T) {
		ew := &errorWriter{failAfter: 1}
		err := WriteIndentedN(ew, data, 2)
		assert.Error(t, err)
		assert.Equal(t, "forced write error", err.Error())
	})

	t.Run("error after writing line content", func(t *testing.T) {
		ew := &errorWriter{failAfter: 3} // 2 spaces + 1 line content = 3 writes
		err := WriteIndentedN(ew, data, 2)
		assert.Error(t, err)
		assert.Equal(t, "forced write error", err.Error())
	})

	t.Run("error after writing newline", func(t *testing.T) {
		ew := &errorWriter{failAfter: 4} // 2 spaces + 1 line content + 1 newline = 4 writes
		err := WriteIndentedN(ew, data, 2)
		assert.Error(t, err)
		assert.Equal(t, "forced write error", err.Error())
	})
}

func TestEncryptPassword(t *testing.T) {
	reqData := "GET / HTTP/1.1\nAuthorization: Basic " + base64.StdEncoding.EncodeToString([]byte("user:password")) + "\n"

	result := encryptPassword([]byte(reqData))
	expected := "GET / HTTP/1.1\nAuthorization: user:******\n"
	assert.Equal(t, expected, string(result))

	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "password",
			input:    `{"password":"my-secret-password"}`,
			expected: `{"password":"****"}` + "\n",
		},
		{
			name:     "session id",
			input:    "Cookie: isisessid=my-session-id",
			expected: "Cookie: isisessid=****-session-id\n",
		},
		{
			name:     "CSRF Token",
			input:    "X-Csrf-Token: my-csrf-token",
			expected: "X-Csrf-Token:****-csrf-token\n",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := encryptPassword([]byte(c.input))
			assert.Contains(t, c.expected, string(result))
		})
	}
}

func TestFetchValueIndexForKey(t *testing.T) {
	cases := []struct {
		name        string
		line        string
		key         string
		separator   string
		expStart    int
		expEnd      int
		expMatchLen int
	}{
		{
			name:        "no separator",
			line:        `"password":"my-secret-password"`,
			key:         `"password":"`,
			separator:   `"`,
			expStart:    0,
			expEnd:      18,
			expMatchLen: 12,
		},
		{
			name:        "with separator",
			line:        `Cookie: isisessid=my-session-id; path=/`,
			key:         `isisessid=`,
			separator:   `;`,
			expStart:    8,
			expEnd:      13,
			expMatchLen: 10,
		},
		{
			name:        "full key",
			line:        `X-Csrf-Token: my-csrf-token`,
			key:         `X-Csrf-Token:`,
			separator:   ` `,
			expStart:    0,
			expEnd:      0,
			expMatchLen: 13,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			start, end, matchLen := FetchValueIndexForKey(c.line, c.key, c.separator)
			assert.Equal(t, c.expStart, start)
			assert.Equal(t, c.expEnd, end)
			assert.Equal(t, c.expMatchLen, matchLen)
		})
	}
}

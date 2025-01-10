/*
Copyright (c) 2022 Dell Inc, or its subsidiaries.

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
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type EmptyMockBody struct{}
type MockBody struct {
	ReadFunc  func(p []byte) (n int, err error)
	CloseFunc func() error
}

func (m *MockBody) Read(p []byte) (n int, err error) {
	return m.ReadFunc(p)
}

func (m *MockBody) Close() error {
	return m.CloseFunc()
}

func assertLen(t *testing.T, obj interface{}, expLen int) {
	if !assert.Len(t, obj, expLen) {
		t.FailNow()
	}
}

func assertError(t *testing.T, err error) {
	if !assert.Error(t, err) {
		t.FailNow()
	}
}

func assertNoError(t *testing.T, err error) {
	if !assert.NoError(t, err) {
		t.FailNow()
	}
}

func assertNil(t *testing.T, i interface{}) {
	if !assert.Nil(t, i) {
		t.FailNow()
	}
}

func assertNotNil(t *testing.T, i interface{}) {
	if !assert.NotNil(t, i) {
		t.FailNow()
	}
}
func TestNew(t *testing.T) {

	ctx := context.Background()
	hostname := "example.com"
	username := "testuser"
	password := "testpassword"
	groupname := "testgroup"
	verboseLogging := uint(1)
	authType := uint8(42)
	authType = authTypeBasic

	// Create a mock ClientOptions
	opts := &ClientOptions{
		VolumesPath:             "test/volumes",
		VolumesPathPermissions:  "test/permissions",
		IgnoreUnresolvableHosts: true,
		Timeout:                 10 * time.Second,
		Insecure:                true,
	}

	// Call the function
	c, _ := New(ctx, hostname, username, password, groupname, verboseLogging, authType, opts)
	assert.Equal(t, nil, c)

	c, err := New(ctx, "", username, password, groupname, verboseLogging, authType, opts)
	assert.Equal(t, errors.New("missing endpoint, username, or password"), err)

	authType = 2
	c, _ = New(ctx, hostname, username, password, groupname, verboseLogging, authType, opts)
	assert.Equal(t, nil, c)

	authType = authTypeSessionBased
	c, _ = New(ctx, hostname, username, password, groupname, verboseLogging, authType, opts)
	assert.Equal(t, nil, c)

	opts = &ClientOptions{
		VolumesPath:             "test/volumes",
		VolumesPathPermissions:  "test/permissions",
		IgnoreUnresolvableHosts: true,
		Timeout:                 10 * time.Second,
		Insecure:                false,
	}
	c, _ = New(ctx, hostname, username, password, groupname, verboseLogging, authType, opts)
	assert.Equal(t, nil, c)
}

func TestDoAndGetResponseBody(t *testing.T) {
	// Create a mock client
	c := &client{
		hostname: "https://example.com",
		http:     http.DefaultClient,
	}
	ctx := context.Background()

	res := &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(nil),
	}

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()
	c.hostname = server.URL
	headers := map[string]string{
		"Content-Type": "application/json",
		"Accept":       "application/json",
	}
	orderedValues := [][][]byte{
		{
			[]byte("value1"),
			[]byte("value2"),
		},
	}
	res, _, err := c.DoAndGetResponseBody(ctx, http.MethodGet, "api/v1/endpoint", "", orderedValues, headers, EmptyMockBody{})
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)

	body := &MockBody{
		ReadFunc: func(_ []byte) (n int, err error) {
			return 0, io.EOF
		},
		CloseFunc: func() error {
			return nil
		},
	}
	res, _, err = c.DoAndGetResponseBody(ctx, http.MethodGet, "api/v1/endpoint", "ID", orderedValues, headers, body)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestAuthenticate(t *testing.T) {
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()
	username := "testuser"
	password := "testpassword"
	endpoint := "https://example.com"

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/session/1/session/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Authentication successful"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	err := c.authenticate(ctx, username, password, endpoint)
	assert.Equal(t, errors.New("authenticate error. response-"), err)
	assert.Equal(t, "", c.GetReferer())

	// Create a mock server for 201 response code
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/session/1/session/", r.URL.Path)

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message":"Authentication successful"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	err = c.authenticate(ctx, username, password, endpoint)
	assert.Equal(t, "", c.GetReferer())

	// create a mock server for 401 response code
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/session/1/session/", r.URL.Path)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"message":"Authentication failed"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	err = c.authenticate(ctx, username, password, endpoint)
	assert.EqualError(t, err, "authentication failed. unable to login to powerscale. verify username and password")

}

func TestExecuteWithRetryAuthenticate(t *testing.T) {
	// Create a mock client
	c := &client{
		http:     http.DefaultClient,
		authType: authTypeSessionBased,
		username: "testuser",
		password: "testpassword",
		hostname: "https://example.com",
	}
	ctx := context.Background()
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint/", r.URL.String())
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	c.hostname = server.URL
	headers := map[string]string{
		"Content-Type": "text/html",
	}
	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, "api/v1/endpoint", "", nil, headers, nil, nil)
	assert.NotEqual(t, "authentication successful", err)
}

func TestParseJSONHTMLError(t *testing.T) {
	// Create a mock HTML response

	htmlResponse := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       io.NopCloser(bytes.NewBufferString("<html><head><title>Bad Request</title></head><body><h1>Bad Request</h1></body></html>")),
		Header:     http.Header{"Content-Type": {"text/html"}},
	}

	err := parseJSONHTMLError(htmlResponse)
	expectedError := &HTMLError{
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
	}
	assert.Equal(t, expectedError, err)
}

func TestDoWithHeaders(t *testing.T) {
	// Create a mock client
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "/api/v1/endpoint/", r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	resp := &struct {
		Message string `json:"message"`
	}{}

	err := c.DoWithHeaders(ctx, http.MethodGet, "api/v1/endpoint", "", nil, nil, nil, resp)
	assert.NoError(t, err)
	expectedResp := &struct {
		Message string `json:"message"`
	}{
		Message: "Success",
	}
	assert.Equal(t, expectedResp, resp)
}

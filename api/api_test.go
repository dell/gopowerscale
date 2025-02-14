/*
Copyright (c) 2022-2025 Dell Inc, or its subsidiaries.

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

	"encoding/json"
	"github.com/stretchr/testify/assert"
	"strings"
)

type (
	EmptyMockBody struct{}
	MockBody      struct {
		ReadFunc  func(p []byte) (n int, err error)
		CloseFunc func() error
	}
)

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

func newMockHTTPServer(handleReq func(http.ResponseWriter, *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleReq(w, r)
	}))
}

func TestNew(t *testing.T) {

	getReqHandler := func(serverVersion string) func(http.ResponseWriter, *http.Request) {
		if serverVersion != "" {
			return func(w http.ResponseWriter, _) {
				res := &apiVerResponse{Latest: &serverVersion}
				w.WriteHeader(http.StatusOK)
				body, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, err = w.Write(body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
			}
		}
		return func(w http.ResponseWriter, _) {
			w.WriteHeader(http.StatusOK)
		}
	}

	serverURL := "server.URL"

	testData := []struct {
		testName       string
		hostname       string
		username       string
		password       string
		groupName      string
		verboseLogging uint
		authType       uint8
		opts           *ClientOptions
		reqHandler     func(http.ResponseWriter, *http.Request)
		expectedErr    string
	}{
		{
			testName:    "Negative: empty call params",
			expectedErr: "missing endpoint, username, or password",
		},
		{
			testName: "Negative: bad hostname",
			hostname: "test",
			username: "testuser",
			password: "testpassword",
			authType: 42, // unknown auth type should default to basic
			opts: &ClientOptions{
				VolumesPath:             "test/volumes",
				VolumesPathPermissions:  "test/permissions",
				IgnoreUnresolvableHosts: true,
				Timeout:                 10 * time.Second,
				Insecure:                true,
			},
			expectedErr: "unsupported protocol scheme",
		},
		{
			testName:       "Negative: empty server response",
			hostname:       serverURL,
			username:       "testuser",
			password:       "testpassword",
			groupName:      "testgroup",
			verboseLogging: 1,
			authType:       authTypeSessionBased,
			opts: &ClientOptions{
				Insecure: false,
			},
			reqHandler:  getReqHandler(""),
			expectedErr: "OneFS releases older than",
		},
		{
			testName:    "Negative: malformed major version in response",
			hostname:    serverURL,
			username:    "testuser",
			password:    "testpassword",
			reqHandler:  getReqHandler("a.3"),
			expectedErr: "strconv.ParseUint: parsing ",
		},
		{
			testName:    "Negative: malformed minor version in response",
			hostname:    serverURL,
			username:    "testuser",
			password:    "testpassword",
			reqHandler:  getReqHandler("8.b"),
			expectedErr: "strconv.ParseUint: parsing ",
		},
		{
			testName:    "Positive: correct version in response",
			hostname:    serverURL,
			username:    "testuser",
			password:    "testpassword",
			reqHandler:  getReqHandler("8.3"),
			expectedErr: "",
		},
	}

	for _, td := range testData {
		t.Run(td.testName, func(t *testing.T) {
			if td.reqHandler != nil {
				server := newMockHTTPServer(td.reqHandler)
				if td.hostname == serverURL {
					td.hostname = server.URL
				}
				defer server.Close()
			}
			c, err := New(
				context.Background(),
				td.hostname,
				td.username,
				td.password,
				td.groupName,
				td.verboseLogging,
				td.authType,
				td.opts)
			if td.expectedErr != "" {
				assert.ErrorContains(t, err, td.expectedErr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, c)
			}
		})
	}
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
		authType: authTypeBasic,
		username: "testuser",
		password: "testpassword",
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if strings.HasPrefix(r.URL.Path, "/bad-auth-") {
				var res *JSONError
				if strings.HasSuffix(r.URL.Path, "401/") {
					res = &JSONError{StatusCode: http.StatusUnauthorized, Err: []Error{{Message: "Unauthorized", Code: "401"}}}
					w.WriteHeader(http.StatusUnauthorized)
				} else if strings.HasSuffix(r.URL.Path, "400/") {
					res = &JSONError{StatusCode: http.StatusBadRequest, Err: []Error{{Message: "Bad Request", Code: "400"}}}
					w.WriteHeader(http.StatusBadRequest)
				} else {
					res = &JSONError{StatusCode: http.StatusNotFound, Err: []Error{{Message: "Unknown URL", Code: "404"}}}
					w.WriteHeader(http.StatusNotFound)
				}
				body, err := json.Marshal(res)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				_, err = w.Write(body)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			} else if strings.HasPrefix(r.URL.Path, "/bad-html-auth-") {
				var body string
				w.Header().Set("Content-Type", "text/html")
				if strings.HasSuffix(r.URL.Path, "401/") {
					body = "<html><head><title>HTML error 401 title</title></head><body></body></html>"
					w.WriteHeader(http.StatusUnauthorized)
				} else if strings.HasSuffix(r.URL.Path, "400/") {
					body = "<html><head><title>HTML error 400 title</title></head><body></body></html>"
					w.WriteHeader(http.StatusBadRequest)
				} else {
					body = "<html><head><title>HTML error title</title></head><body></body></html>"
					w.WriteHeader(http.StatusNotFound)
				}
				_, err := w.Write([]byte(body))
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
				return
			} else if r.URL.Path == "/good-path/" {
				w.WriteHeader(http.StatusOK)
				return
			}
		} else if r.Method == http.MethodPost {
			// Authentication successful
			w.Header().Set(isiSessCsrfToken, "isisessid=123;isicsrf=abc;")
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	c.hostname = server.URL
	headers := map[string]string{
		"Content-Type": "text/html",
	}

	err := c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/good-path", "", nil, headers, nil, nil)
	assert.NoError(t, err)

	c.authType = authTypeSessionBased

	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/good-path", "", nil, headers, nil, nil)
	assert.NoError(t, err)

	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/bad-auth-401", "", nil, headers, nil, nil)
	assert.Error(t, err)

	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/bad-auth-400", "", nil, headers, nil, nil)
	assert.Error(t, err)

	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/bad-html-auth-401", "", nil, headers, nil, nil)
	assert.Error(t, err)

	err = c.executeWithRetryAuthenticate(ctx, http.MethodGet, "/bad-html-auth-400", "", nil, headers, nil, nil)
	assert.Error(t, err)
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

func TestClient_APIVersion(t *testing.T) {
	c := &client{apiVersion: 1}
	assert.Equal(t, uint8(1), c.APIVersion())
}

func TestClient_User(t *testing.T) {
	c := &client{username: "testuser"}
	assert.Equal(t, "testuser", c.User())
}

func TestClient_Group(t *testing.T) {
	c := &client{groupname: "testgroup"}
	assert.Equal(t, "testgroup", c.Group())
}

func TestClient_VolumesPath(t *testing.T) {
	c := &client{volumePath: "/mnt/volumes"}
	assert.Equal(t, "/mnt/volumes", c.VolumesPath())
}

func TestClient_VolumePath(t *testing.T) {
	c := &client{volumePath: "/mnt/volumes"}
	assert.Equal(t, "/mnt/volumes/volume1", c.VolumePath("volume1"))
}

func TestHTMLError_Error(t *testing.T) {
	err := &HTMLError{Message: "HTML error message"}
	assert.Equal(t, "HTML error message", err.Error())
}

func TestClient_SetAuthToken(t *testing.T) {
	c := &client{}
	c.SetAuthToken("testcookie")
	assert.Equal(t, "testcookie", c.sessionCredentials.sessionCookies)
}

func TestClient_SetCSRFToken(t *testing.T) {
	c := &client{}
	c.SetCSRFToken("testcsrf")
	assert.Equal(t, "testcsrf", c.sessionCredentials.sessionCSRF)
}

func TestClient_SetReferer(t *testing.T) {
	c := &client{}
	c.SetReferer("testreferer")
	assert.Equal(t, "testreferer", c.sessionCredentials.referer)
}

func TestClient_GetCSRFToken(t *testing.T) {
	c := &client{}
	c.GetCSRFToken()
	assert.Equal(t, "", c.sessionCredentials.sessionCSRF)
}

func TestParseJSONHTMLError(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedErr    error
		expectedStatus int
	}{
		{
			name:           "HTML error response",
			contentType:    "text/html",
			body:           `<html><head><title>HTML error title</title></head><body><h1>HTML error message</h1></body></html>`,
			expectedErr:    &HTMLError{Message: "HTML error message"},
			expectedStatus: 401,
		},
		{
			name:           "HTML error without h1",
			contentType:    "text/html",
			body:           `<html><head><title>HTML error title</title></head><body></body></html>`,
			expectedErr:    &HTMLError{Message: "HTML error title"},
			expectedStatus: 403,
		},
		{
			name:           "Invalid JSON",
			contentType:    "application/json",
			body:           `{invalid json`,
			expectedErr:    &JSONError{Err: []Error{{Message: "invalid character 'i' looking for beginning of object key string"}}},
			expectedStatus: 400,
		},
		{
			name:           "Invalid HTML",
			contentType:    "text/html",
			body:           `<html>`,
			expectedErr:    &HTMLError{Message: ""},
			expectedStatus: 500,
		},
		{
			name:           "JSON error with empty message",
			contentType:    "application/json",
			body:           `{"errors":[{"message":""}]}`,
			expectedErr:    &JSONError{Err: []Error{{Message: "400"}}},
			expectedStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := bytes.NewBufferString(tt.body)
			resp := httptest.NewRecorder()
			resp.Body = body
			resp.Header().Set("Content-Type", tt.contentType)
			resp.Code = tt.expectedStatus

			err := parseJSONHTMLError(resp.Result())

			if tt.expectedErr != nil {
				assert.NotNil(t, err)

				switch expected := tt.expectedErr.(type) {
				case *JSONError:
					assert.Contains(t, err.Error(), expected.Error())
				default:
					assert.IsType(t, expected, err)
					assert.EqualError(t, err, expected.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestClient_Put(t *testing.T) {
	// Create a mock client
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodPut, r.Method)
		assert.Equal(t, "/PUT/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	body := map[string]string{
		"Content-Type": "application/json",
	}
	resp := &struct {
		Message string `json:"message"`
	}{}
	// Call the Put method
	err := c.Put(ctx, http.MethodPut, "api/v1/endpoint", nil, nil, body, resp)
	assert.NoError(t, err)
}

func TestClient_Post(t *testing.T) {
	// Create a mock client
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "/POST/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	body := map[string]string{
		"Content-Type": "application/json",
	}
	resp := &struct {
		Message string `json:"message"`
	}{}
	// Call the Post method
	err := c.Post(ctx, http.MethodPost, "api/v1/endpoint", nil, nil, body, resp)

	// Assertions
	assert.NoError(t, err)
}

func TestClient_Delete(t *testing.T) {
	// Create a mock client
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, http.MethodDelete, r.Method)
		assert.Equal(t, "/DELETE/api/v1/endpoint", r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	resp := &struct {
		Message string `json:"message"`
	}{}
	// Call the Delete method
	err := c.Delete(ctx, http.MethodDelete, "api/v1/endpoint", nil, nil, resp)

	// Assertions
	assert.NoError(t, err)
}

func TestClient_Do(t *testing.T) {
	// Create a mock client
	c := &client{
		http: http.DefaultClient,
	}
	ctx := context.Background()

	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the request is as expected
		assert.Equal(t, "method", r.Method)
		assert.Equal(t, "/api/v1/endpoint/", r.URL.String())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message":"Success"}`))
	}))
	defer server.Close()
	c.hostname = server.URL
	resp := &struct {
		Message string `json:"message"`
	}{}
	// Call the Do method
	err := c.Do(ctx, "method", "api/v1/endpoint", "", nil, resp, resp)

	// Assertions
	assert.NoError(t, err)
}

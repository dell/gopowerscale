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
package goisilon

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	os.Setenv("GOISILON_INSECURE", "true")
	os.Setenv("GOISILON_UNRESOLVABLE_HOSTS", "false")
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "8.0.0.0", "message": "success"}`))
	}))
	defer mockServer.Close()

	client, err := NewClient(context.Background())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	assert.Nil(t, client)

	os.Setenv("GOISILON_INSECURE", "true")
	os.Unsetenv("GOISILON_UNRESOLVABLE_HOSTS")
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "8.0.0.0", "message": "success"}`))
	}))
	defer mockServer.Close()

	client, err = NewClient(context.Background())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	assert.Nil(t, client)

	os.Unsetenv("GOISILON_INSECURE")
	os.Unsetenv("GOISILON_UNRESOLVABLE_HOSTS")
	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "8.0.0.0", "message": "success"}`))
	}))
	defer mockServer.Close()

	client, err = NewClient(context.Background())
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
	assert.Nil(t, client)

	os.Setenv("GOISILON_INSECURE", "true")
	os.Setenv("GOISILON_UNRESOLVABLE_HOSTS", "false")
	os.Setenv("GOISILON_AUTHTYPE", "1")
	os.Setenv("GOISILON_ENDPOINT", mockServer.URL)
	os.Setenv("GOISILON_USERNAME", "user")
	os.Setenv("GOISILON_GROUP", "group")
	os.Setenv("GOISILON_PASSWORD", "pass")
	os.Setenv("GOISILON_VOLUMEPATH", "/path")
	os.Setenv("GOISILON_VOLUMEPATH_PERMISSIONS", "0777")

	mockServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "8.0.0.0", "message": "success"}`))
	}))
	defer mockServer.Close()

	client, _ = NewClient(context.Background())
	assert.Nil(t, client)
}

func TestNewClientWithArgs(t *testing.T) {
	os.Setenv("GOISILON_TIMEOUT", "30s")
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"version": "8.0.0.0", "message": "success"}`))
	}))
	defer mockServer.Close()

	client, _ := NewClientWithArgs(context.Background(), mockServer.URL, true, 1, "user", "group", "pass", "/path", "0777", false, 1)
	assert.Nil(t, client)
}

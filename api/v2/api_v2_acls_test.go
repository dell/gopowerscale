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
package v2

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestACLInspect(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	client.On("VolumesPath", anyArgs...).Return("").Once()
	_, err := ACLInspect(ctx, client, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestACLUpdate(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	var authoritativeType AuthoritativeType = 5
	acl := ACL{
		Authoritative: &authoritativeType,
	}
	client.On("Put", anyArgs...).Return(nil).Once()
	client.On("VolumesPath", anyArgs...).Return("").Once()
	err := ACLUpdate(ctx, client, "", &acl)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestParseAuthoritativeType(t *testing.T) {
	authType := ParseAuthoritativeType(authoritativeTypeACLStr)
	assert.Equal(t, AuthoritativeTypeACL, authType)

	authType = ParseAuthoritativeType(authoritativeTypeModeStr)
	assert.Equal(t, AuthoritativeTypeMode, authType)

	authType = ParseAuthoritativeType("")
	assert.Equal(t, AuthoritativeTypeUnknown, authType)
}

func TestParseActionType(t *testing.T) {
	actionType := ParseActionType(ActionTypeReplaceStr)
	assert.Equal(t, ActionTypeReplace, actionType)

	actionType = ParseActionType(ActionTypeUpdateStr)
	assert.Equal(t, ActionTypeUpdate, actionType)

	actionType = ParseActionType("")
	assert.Equal(t, ActionTypeUnknown, actionType)
}

func TestParseFileMode(t *testing.T) {
	_, err := ParseFileMode("755")
	assert.Equal(t, nil, err)

	_, err = ParseFileMode("75")
	assert.Equal(t, errInvalidFileMode, err)
}

// UT for MarshalJSON()
func TestAuthoritativeType_MarshalJSON(t *testing.T) {
	jsonTests := []struct {
		name     string
		value    AuthoritativeType
		expected string
	}{
		{"Marshal AuthoritativeTypeUnknown", AuthoritativeTypeUnknown, `"unknown"`},
		{"Marshal AuthoritativeTypeACL", AuthoritativeTypeACL, `"acl"`},
		{"Marshal AuthoritativeTypeMode", AuthoritativeTypeMode, `"mode"`},
	}
	for _, tt := range jsonTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestActionType_MarshalJSON(t *testing.T) {
	jsonTests := []struct {
		name     string
		value    ActionType
		expected string
	}{
		{"Marshal ActionTypeUnknown", ActionTypeUnknown, `"unknown"`},
		{"Marshal ActionTypeReplace", ActionTypeReplace, `"replace"`},
		{"Marshal ActionTypeUpdate", ActionTypeUpdate, `"update"`},
	}
	for _, tt := range jsonTests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.value)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

// UT for UnmarshalJSON()
func TestActionType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    ActionType
		expectError bool
	}{
		{"Unmarshal replace", []byte(`"replace"`), ActionTypeReplace, false},
		{"Unmarshal update", []byte(`"update"`), ActionTypeUpdate, false},
		{"Unmarshal unknown", []byte(`"unknown"`), ActionTypeUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result ActionType
			err := json.Unmarshal([]byte(tt.input), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAuthoritativeType_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name        string
		input       []byte
		expected    AuthoritativeType
		expectError bool
	}{
		{"Unmarshal acl", []byte(`"acl"`), AuthoritativeTypeACL, false},
		{"Unmarshal mode", []byte(`"mode"`), AuthoritativeTypeMode, false},
		{"Unmarshal unknown", []byte(`"unknown"`), AuthoritativeTypeUnknown, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result AuthoritativeType
			err := json.Unmarshal([]byte(tt.input), &result)
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

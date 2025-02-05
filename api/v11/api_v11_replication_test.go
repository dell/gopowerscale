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

package v11

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestGetPolicyByName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetPolicyByName(ctx, client, "")
	assert.Equal(t, errors.New("successful code returned, but policy  not found"), err)

	client.On("Get", anyArgs...).Return(errors.New("unable to get policy")).Twice()
	_, err = GetPolicyByName(ctx, client, "")
	assert.Error(t, err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**Policies)
		*resp = &Policies{
			Policy: []Policy{
				{
					ID:   "test-id",
					Name: "test-name",
				},
			},
		}
	}).Twice()
	_, err = GetPolicyByName(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetTargetPolicyByName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetTargetPolicyByName(ctx, client, "")
	assert.Equal(t, nil, err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**TargetPolicies)
		*resp = &TargetPolicies{
			Policy: []TargetPolicy{
				{
					ID:   "test-id",
					Name: "test-name",
				},
			},
		}
	}).Twice()
	_, err = GetTargetPolicyByName(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetReport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetReport(ctx, client, "")
	assert.Equal(t, errors.New("no reports found with report name "), err)

	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**Reports)
		*resp = &Reports{
			Reports: []Report{
				{
					ID:      "test-id",
					JobID:   12345,
					State:   "completed",
					EndTime: 1609459200,
					Errors:  []string{"error1", "error2"},
				},
			},
		}
	}).Twice()
	_, err = GetReport(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetReportsByPolicyName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetReportsByPolicyName(ctx, client, "", 0)
	assert.Equal(t, errors.New("no reports found for policy "), err)
}

func TestGetJobsByPolicyName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetJobsByPolicyName(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeleteTargetPolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteTargetPolicy(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeletePolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeletePolicy(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestCreatePolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Post", anyArgs...).Return(nil).Twice()
	err := CreatePolicy(ctx, client, "", "", "", "", "", 0, false)
	assert.Equal(t, nil, err)
}

func TestUpdatePolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	policy := Policy{
		ID: "",
	}
	client.On("Put", anyArgs...).Return(nil).Twice()
	err := UpdatePolicy(ctx, client, &policy)
	assert.Equal(t, nil, err)
}

func TestResolvePolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	resolvePolicy := ResolvePolicyReq{
		ID: "",
	}
	client.On("Put", anyArgs...).Return(nil).Twice()
	err := ResolvePolicy(ctx, client, &resolvePolicy)
	assert.Equal(t, nil, err)
}

func TestResetPolicy(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Post", anyArgs...).Return(nil).Twice()
	err := ResetPolicy(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestStartSyncIQJob(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	jobRequest := JobRequest{
		ID: "",
	}

	client.On("Post", anyArgs...).Return(nil).Twice()
	_, err := StartSyncIQJob(ctx, client, &jobRequest)
	assert.Equal(t, nil, err)
}

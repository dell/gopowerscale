/*
Copyright (c) 2021-2024 Dell Inc, or its subsidiaries.
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
	"fmt"
	"testing"

	apiv11 "github.com/dell/goisilon/api/v11"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	policiesPath       = "/platform/11/sync/policies/"
	targetPoliciesPath = "/platform/11/sync/target/policies/"
	jobsPath           = "/platform/11/sync/jobs/"
)

const resolveErrorToIgnore = "The policy was not conflicted, so no change was made"

func TestGetPolicyByName(t *testing.T) {
	ctx := context.Background()
	policyID := "test-policy"
	expectedPolicy := &apiv11.Policy{
		ID: policyID,
	}
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Policies)
		*resp = &apiv11.Policies{
			Policy: []apiv11.Policy{
				*expectedPolicy,
			},
		}
	}).Once()
	policy, err := client.GetPolicyByName(ctx, policyID)
	assert.Nil(t, err)
	assert.Equal(t, Policy(expectedPolicy), policy)
}

func TestGetTargetPolicyByName(t *testing.T) {
	ctx := context.Background()
	policyID := "test-target-policy"
	expectedPolicy := &apiv11.TargetPolicy{
		ID:                      policyID,
		Name:                    "Test Policy",
		SourceClusterGUID:       "source-cluster-guid",
		LastJobState:            apiv11.JobState("running"),
		TargetPath:              "/target/path",
		SourceHost:              "source-host",
		LastSourceCoordinatorIP: "192.168.1.1",
		FailoverFailbackState:   apiv11.FailoverFailbackState("writes_disabled"),
	}

	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*expectedPolicy,
			},
		}
	}).Once()
	policy, err := client.GetTargetPolicyByName(ctx, policyID)
	assert.Nil(t, err)
	assert.Equal(t, TargetPolicy(expectedPolicy), policy)
}

func TestCreatePolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	name := "test-policy"
	rpo := 5
	sourcePath := "/source/path"
	targetPath := "/target/path"
	targetHost := "target-host"
	targetCert := "target-cert"
	enabled := true

	expectedData := &apiv11.Policy{
		Action:     "sync",
		Name:       name,
		Enabled:    enabled,
		TargetPath: targetPath,
		SourcePath: sourcePath,
		TargetHost: targetHost,
		JobDelay:   rpo,
		TargetCert: targetCert,
		Schedule:   "when-source-modified",
	}

	var policyResp apiv11.Policy

	// Set up expectations
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		policiesPath,
		"",
		mock.Anything,
		mock.Anything,
		expectedData,
		&policyResp,
	).Return(nil).Once()

	// Call the CreatePolicy method
	err := client.CreatePolicy(ctx, name, rpo, sourcePath, targetPath, targetHost, targetCert, enabled)

	// Assertions
	assert.Nil(t, err)
}

func TestDeletePolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	name := "test-policy"
	// Define response for Delete operation
	var resp string
	// Set up expectations
	client.API.(*mocks.Client).On(
		"Delete",
		ctx,
		policiesPath,
		name,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()
	// Call the DeletePolicy method
	err := client.DeletePolicy(ctx, name)
	// Assertions
	assert.Nil(t, err)
}

func TestDeleteTargetPolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	id := "test-target-policy"
	// Define response for Delete operation
	var resp string
	// Set up expectations
	client.API.(*mocks.Client).On(
		"Delete",
		ctx,
		targetPoliciesPath,
		id,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()
	// Call the DeleteTargetPolicy method
	err := client.DeleteTargetPolicy(ctx, id)
	// Assertions
	assert.Nil(t, err)
}

func TestBreakAssociation(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	targetPolicyName := "test-target-policy"
	targetPolicyID := "test-target-policy-id"

	expectedTargetPolicy := &apiv11.TargetPolicy{
		ID: targetPolicyID,
	}

	// Mock GetTargetPolicyByName method
	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, targetPolicyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*expectedTargetPolicy,
			},
		}
	}).Once()

	// Mock DeleteTargetPolicy method
	client.API.(*mocks.Client).On(
		"Delete",
		ctx,
		targetPoliciesPath,
		targetPolicyID,
		mock.Anything,
		mock.Anything,
		mock.Anything,
	).Return(nil).Once()

	// Call the BreakAssociation method
	err := client.BreakAssociation(ctx, targetPolicyName)

	// Assertions
	assert.Nil(t, err)
}

func TestResetPolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"

	// Define response for Post operation
	var resp apiv11.Policy

	// Set up expectations
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		policiesPath,
		name+"/reset",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()

	// Call the ResetPolicy method
	err := client.ResetPolicy(ctx, name)

	// Assertions
	assert.Nil(t, err)
}

func TestEnablePolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"
	policyID := "policy-id"
	schedule := "when-source-modified"

	pp := &apiv11.Policy{
		ID:       policyID,
		Enabled:  false,
		Schedule: schedule,
	}

	updatedPolicy := &apiv11.Policy{
		Enabled:  true,
		Schedule: schedule,
	}

	// Mock GetPolicyByName method
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Policies)
		*resp = &apiv11.Policies{
			Policy: []apiv11.Policy{
				*pp,
			},
		}
	}).Once()

	// Mock UpdatePolicy method
	client.API.(*mocks.Client).On(
		"Put",
		ctx,
		policiesPath,
		policyID,
		mock.Anything,
		mock.Anything,
		updatedPolicy,
		nil,
	).Return(nil).Once()

	// Call the EnablePolicy method
	err := client.EnablePolicy(ctx, name)

	// Assertions
	assert.Nil(t, err)

}

func TestDisablePolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"
	policyID := "policy-id"
	schedule := "when-source-modified"

	pp := &apiv11.Policy{
		ID:       policyID,
		Enabled:  true,
		Schedule: schedule,
	}

	updatedPolicy := &apiv11.Policy{
		Enabled:  false,
		Schedule: schedule,
	}

	// Mock GetPolicyByName method
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Policies)
		*resp = &apiv11.Policies{
			Policy: []apiv11.Policy{
				*pp,
			},
		}
	}).Once()

	// Mock UpdatePolicy method
	client.API.(*mocks.Client).On(
		"Put",
		ctx,
		policiesPath,
		policyID,
		mock.Anything,
		mock.Anything,
		updatedPolicy,
		nil,
	).Return(nil).Once()

	// Call the EnablePolicy method
	err := client.DisablePolicy(ctx, name)

	// Assertions
	assert.Nil(t, err)

}

func TestSetPolicyEnabledField(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"
	policyID := "policy-id"
	schedule := "when-source-modified"

	pp := &apiv11.Policy{
		ID:       policyID,
		Enabled:  false,
		Schedule: schedule,
	}

	updatedPolicy := &apiv11.Policy{
		Enabled:  true,
		Schedule: schedule,
	}

	// Mock GetPolicyByName method
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Policies)
		*resp = &apiv11.Policies{
			Policy: []apiv11.Policy{
				*pp,
			},
		}
	}).Once()

	// Mock UpdatePolicy method (simulating the Put method call)
	client.API.(*mocks.Client).On(
		"Put",
		ctx,
		policiesPath,
		policyID,
		mock.Anything,
		mock.Anything,
		updatedPolicy,
		mock.Anything,
	).Return(nil).Once()

	// Call the SetPolicyEnabledField method
	err := client.SetPolicyEnabledField(ctx, name, true)

	// Assertions
	assert.Nil(t, err)
}

func TestModifyPolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"
	policyID := "policy-id"
	enabled := true

	pp := &apiv11.Policy{
		ID:      policyID,
		Enabled: enabled,
	}

	t.Run("Manual Schedule", func(t *testing.T) {
		updatedPolicy := &apiv11.Policy{
			Enabled:  enabled,
			Schedule: "",
		}

		// Mock GetPolicyByName method
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*pp,
				},
			}
		}).Once()

		// Mock UpdatePolicy method (simulating the Put method call)
		client.API.(*mocks.Client).On(
			"Put",
			ctx,
			policiesPath,
			policyID,
			mock.Anything,
			mock.Anything,
			updatedPolicy,
			mock.Anything,
		).Return(nil).Once()

		// Call the ModifyPolicy method with manual schedule
		err := client.ModifyPolicy(ctx, name, "", 0)

		// Assertions
		assert.Nil(t, err)
	})

	t.Run("When Source Modified Schedule", func(t *testing.T) {
		rpo := 10
		updatedPolicy := &apiv11.Policy{
			Enabled:  enabled,
			Schedule: "when-source-modified",
			JobDelay: rpo,
		}

		// Mock GetPolicyByName method
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*pp,
				},
			}
		}).Once()

		// Mock UpdatePolicy method (simulating the Put method call)
		client.API.(*mocks.Client).On(
			"Put",
			ctx,
			policiesPath,
			policyID,
			mock.Anything,
			mock.Anything,
			updatedPolicy,
			mock.Anything,
		).Return(nil).Once()

		// Call the ModifyPolicy method with "when-source-modified" schedule
		err := client.ModifyPolicy(ctx, name, "when-source-modified", rpo)

		// Assertions
		assert.Nil(t, err)
	})
}

func TestResolvePolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	name := "test-policy"
	policyID := "policy-id"
	enabled := true
	schedule := "daily"

	pp := &apiv11.Policy{
		ID:       policyID,
		Enabled:  enabled,
		Schedule: schedule,
	}

	resolvePolicyReq := &apiv11.ResolvePolicyReq{
		Conflicted: false,
		Enabled:    enabled,
		Schedule:   schedule,
	}

	t.Run("ResolvePolicy_Success", func(t *testing.T) {
		// Mock GetPolicyByName method
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*pp,
				},
			}
		}).Once()

		// Mock ResolvePolicy method (simulating the Put method call)
		client.API.(*mocks.Client).On(
			"Put",
			ctx,
			policiesPath,
			policyID,
			mock.Anything,
			mock.Anything,
			resolvePolicyReq,
			mock.Anything,
		).Return(nil).Once()

		// Call the ResolvePolicy method with no errors
		err := client.ResolvePolicy(ctx, name)

		// Assertions
		assert.Nil(t, err)
	})

	t.Run("ResolvePolicy_HandlesSpecificError", func(t *testing.T) {
		// Mock GetPolicyByName method
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*pp,
				},
			}
		}).Once()

		// Mock ResolvePolicy method to return a specific error that should be ignored
		client.API.(*mocks.Client).On(
			"Put",
			ctx,
			policiesPath,
			policyID,
			mock.Anything,
			mock.Anything,
			resolvePolicyReq,
			mock.Anything,
		).Return(fmt.Errorf(resolveErrorToIgnore)).Once()
		// Call the ResolvePolicy method and expect the specific error to be ignored
		err := client.ResolvePolicy(ctx, name)

		// Assertions
		assert.Nil(t, err)
	})
}

func TestAllowWrites(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"

	targetPolicy := &apiv11.TargetPolicy{
		ID:                    policyName,
		FailoverFailbackState: WritesDisabled,
	}

	writeEnabledtargetPolicy := &apiv11.TargetPolicy{
		ID:                    policyName,
		FailoverFailbackState: WritesEnabled,
	}

	// Mock GetTargetPolicyByName method
	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*targetPolicy,
			},
		}
	}).Once()

	// Define response for Post operation
	var resp apiv11.Job
	// Set up expectations
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		jobsPath,
		"",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()

	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*writeEnabledtargetPolicy,
			},
		}
	}).Maybe()

	// Call the AllowWrites method
	err := client.AllowWrites(ctx, policyName)

	// Assertions
	assert.Nil(t, err)
}

func TestDisallowWrites(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"

	targetPolicy := &apiv11.TargetPolicy{
		ID:                    policyName,
		FailoverFailbackState: WritesEnabled,
	}

	writeDisabledtargetPolicy := &apiv11.TargetPolicy{
		ID:                    policyName,
		FailoverFailbackState: WritesDisabled,
	}

	// Mock GetTargetPolicyByName method
	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*targetPolicy,
			},
		}
	}).Once()

	// Define response for Post operation
	var resp apiv11.Job

	// Set up expectations
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		jobsPath,
		"",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()

	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*writeDisabledtargetPolicy,
			},
		}
	}).Maybe()

	// Call the AllowWrites method
	err := client.DisallowWrites(ctx, policyName)

	// Assertions
	assert.Nil(t, err)
}

func TestResyncPrep(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"

	// Define response for Post operation
	var resp apiv11.Job

	// Set up expectations
	client.API.(*mocks.Client).On(
		"Post",
		ctx,
		jobsPath,
		"",
		mock.Anything,
		mock.Anything,
		mock.Anything,
		&resp,
	).Return(nil).Once()

	// Call the ResyncPrep method
	err := client.ResyncPrep(ctx, policyName)

	// Assertions
	assert.Nil(t, err)
}

func TestRunActionForPolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	action := apiv11.AllowWrite
	jobRequest := &apiv11.JobRequest{Action: action, ID: policyName}
	expectedJob := &apiv11.Job{}

	// Mock Post method
	client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, jobRequest, expectedJob).Return(nil).Once()

	// Call the RunActionForPolicy method
	job, err := client.RunActionForPolicy(ctx, policyName, action)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedJob, job)
}

func TestStartSyncIQJob(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	jobRequest := &apiv11.JobRequest{
		ID:     "policy-id",
		Action: apiv11.ResyncPrep,
	}
	expectedJob := &apiv11.Job{}

	// Expect the Post method to be called with the specified parameters
	client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, jobRequest, expectedJob).Return(nil).Once()

	// Call the StartSyncIQJob method
	job, err := client.StartSyncIQJob(ctx, jobRequest)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedJob, job)
}

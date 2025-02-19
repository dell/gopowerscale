/*
Copyright (c) 2021-2025 Dell Inc, or its subsidiaries.
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
	"errors"
	"fmt"
	"testing"

	"github.com/dell/goisilon/api"
	apiv11 "github.com/dell/goisilon/api/v11"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	policiesPath       = "/platform/11/sync/policies/"
	targetPoliciesPath = "/platform/11/sync/target/policies/"
	jobsPath           = "/platform/11/sync/jobs/"
	reportsPath        = "/platform/11/sync/reports"
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
	expectedError := errors.New("target policy not found")

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

	// Mock GetTargetPolicyByName method error scenario
	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, targetPolicyName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(expectedError).Once()

	// Call the BreakAssociation method for error scenario
	err = client.BreakAssociation(ctx, targetPolicyName)

	// Assertions for error scenario
	assert.NotNil(t, err)
	assert.Equal(t, expectedError, err)
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

	// Scenario 1: GetPolicyByName returns an error
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, name).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()
	err := client.SetPolicyEnabledField(ctx, name, true)
	assert.NotNil(t, err)
	assert.Equal(t, "policy not found", err.Error())

	// Scenario 2: Policy already has the same Enabled value
	client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Policies)
		*resp = &apiv11.Policies{
			Policy: []apiv11.Policy{
				*pp,
			},
		}
	}).Once()
	err = client.SetPolicyEnabledField(ctx, name, false)
	assert.Nil(t, err)

	// Scenario 3: Successfully update the policy
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
	err = client.SetPolicyEnabledField(ctx, name, true)

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

	t.Run("When GetPolicyByName return an error", func(t *testing.T) {
		// Scenario 1: GetPolicyByName returns an error
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, name).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()

		// Call the ModifyPolicy method
		err := client.ModifyPolicy(ctx, name, "", 0)

		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

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

	t.Run("When GetPolicyByName return an error", func(t *testing.T) {
		// Scenario 1: GetPolicyByName returns an error
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, name).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()

		// Call the ResolvePolicy method
		err := client.ResolvePolicy(ctx, name)

		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

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

	t.Run("GetTargetPolicyByName returns an error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()
		err = client.AllowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("FailoverFailbackState == WritesEnabled", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*writeEnabledtargetPolicy,
				},
			}
		}).Once()

		err = client.AllowWrites(ctx, policyName)
		assert.Nil(t, err)
	})

	// Define response for Post operation
	var resp apiv11.Job

	t.Run("RunActionForPolicy returns an error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*targetPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On(
			"Post",
			ctx,
			jobsPath,
			"",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			&resp,
		).Return(errors.New("policy not found")).Once()
		err = client.AllowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("WaitForTargetPolicyCondition returns error", func(t *testing.T) {
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

		client.API.(*mocks.Client).On("WaitForTargetPolicyCondition", ctx, policyName, WritesEnabled).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("wait condition failed")).Once()

		err = client.AllowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "wait condition failed", err.Error())
	})

	t.Run("AllowWrites method final scenario", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*targetPolicy,
				},
			}
		}).Once()

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

		err = client.AllowWrites(ctx, policyName)
		assert.Nil(t, err)
	})
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

	t.Run("GetTargetPolicyByName returns an error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()
		err := client.DisallowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("FailoverFailbackState == WritesDisabled", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*writeDisabledtargetPolicy,
				},
			}
		}).Once()

		err := client.DisallowWrites(ctx, policyName)
		assert.Nil(t, err)
	})

	// Define response for Post operation
	var resp apiv11.Job

	t.Run("RunActionForPolicy returns an error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*targetPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On(
			"Post",
			ctx,
			jobsPath,
			"",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			&resp,
		).Return(errors.New("policy not found")).Once()
		err := client.DisallowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("WaitForTargetPolicyCondition returns error", func(t *testing.T) {
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

		client.API.(*mocks.Client).On("WaitForTargetPolicyCondition", ctx, policyName, WritesDisabled).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("wait condition failed")).Once()

		err := client.DisallowWrites(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "wait condition failed", err.Error())
	})

	t.Run("DisallowWrites method final scenario", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.TargetPolicies)
			*resp = &apiv11.TargetPolicies{
				Policy: []apiv11.TargetPolicy{
					*targetPolicy,
				},
			}
		}).Once()

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

		err := client.DisallowWrites(ctx, policyName)
		assert.Nil(t, err)
	})
}

func TestResyncPrep(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"

	// Define response for Post operation
	var resp apiv11.Job

	t.Run("Post returns an error", func(t *testing.T) {
		client.API.(*mocks.Client).On(
			"Post",
			ctx,
			jobsPath,
			"",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			&resp,
		).Return(errors.New("policy not found")).Once()
		// Call the ResyncPrep method
		err := client.ResyncPrep(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, "policy not found", err.Error())
	})

	t.Run("Post succeeds", func(t *testing.T) {
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
	})
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

func TestGetReport(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	reportName := "test-report"
	expectedReport := &apiv11.Report{
		ID:      "report-id",
		JobID:   12345,
		State:   "completed",
		EndTime: 1609459200,
		Errors:  []string{"error1", "error2"},
	}
	reports := &apiv11.Reports{
		Reports: []apiv11.Report{*expectedReport},
	}

	// Mock GetReport method
	client.API.(*mocks.Client).On("GetReport", mock.Anything, reportName).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Reports)
		*resp = reports
	}).Once()

	// Call the GetReport method
	report, err := client.GetReport(ctx, reportName)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedReport, report)
}

func TestGetReportsByPolicyName(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	reportsForPolicy := 5
	expectedReports := &apiv11.Reports{
		Reports: []apiv11.Report{
			{
				Policy:  apiv11.Policy{Name: policyName},
				ID:      "report-id-1",
				JobID:   12345,
				State:   "completed",
				EndTime: 1609459200,
				Errors:  []string{"error1", "error2"},
			},
			{
				Policy:  apiv11.Policy{},
				ID:      "report-id-2",
				JobID:   12346,
				State:   "completed",
				EndTime: 1609459201,
				Errors:  []string{},
			},
		},
	}

	// Mock GetReportsByPolicyName method
	client.API.(*mocks.Client).On("GetReportsByPolicyName", mock.Anything, policyName, reportsForPolicy).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.Reports)
		*resp = expectedReports
	}).Once()

	// Call the GetReportsByPolicyName method
	reports, err := client.GetReportsByPolicyName(ctx, policyName, reportsForPolicy)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedReports, reports)
}

func TestWaitForPolicyEnabledFieldCondition(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	policyID := "policy-id"

	// Initially returning policy with Enabled field set to false
	firstPolicy := &apiv11.Policy{
		ID:      policyID,
		Enabled: false,
	}

	// Eventually returning policy with Enabled field set to true
	enabledPolicy := &apiv11.Policy{
		ID:      policyID,
		Enabled: true,
	}

	t.Run("return false, err", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("get policy error")).Once()
		// Call the WaitForPolicyEnabledFieldCondition method
		err := client.WaitForPolicyEnabledFieldCondition(ctx, policyName, true)

		// Assertions
		assert.NotNil(t, err)
		assert.Equal(t, "get policy error", err.Error())
	})
	t.Run("return pollErr", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, defaultPoll) // Using a short timeout to force a timeout error
		defer cancel()

		// Set up the mocks to always return the same firstPolicy with Enabled = false.
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Times(2)
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*firstPolicy,
				},
			}
		}).Twice()

		// Call the WaitForPolicyEnabledFieldCondition method with shorter timeout
		err := client.WaitForPolicyEnabledFieldCondition(ctx, policyName, true)

		// Assertions
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "timed out waiting for the condition")
	})

	t.Run("PollImmediate returns true eventually", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*firstPolicy,
				},
			}
		}).Once()

		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()

		// Call the WaitForPolicyEnabledFieldCondition method
		err := client.WaitForPolicyEnabledFieldCondition(ctx, policyName, true)

		// Assertions
		assert.Nil(t, err)
	})
}

func TestWaitForNoActiveJobs(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	activeJobs := []apiv11.Job{
		{ID: "active-job-1"},
		{ID: "active-job-2"},
	}
	noActiveJobs := []apiv11.Job{}

	t.Run("return false, err", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("get jobs error")).Once()

		// Call the WaitForNoActiveJobs method
		err := client.WaitForNoActiveJobs(ctx, policyName)

		// Assertions
		assert.NotNil(t, err)
		assert.Equal(t, "get jobs error", err.Error())
	})

	t.Run("return pollErr", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, defaultPoll) // Using a short timeout to force a timeout error
		defer cancel()

		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Twice()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: activeJobs,
			}
		}).Twice()

		// Call the WaitForNoActiveJobs method
		err := client.WaitForNoActiveJobs(ctx, policyName)

		// Assertions
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "timed out waiting for the condition")
	})

	t.Run("PollImmediate returns true eventually", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: activeJobs,
			}
		}).Once()

		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noActiveJobs,
			}
		}).Once()

		// Call the WaitForNoActiveJobs method
		err := client.WaitForNoActiveJobs(ctx, policyName)

		// Assertions
		assert.Nil(t, err)
	})
}

func TestWaitForPolicyLastJobState(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	initialState := apiv11.SCHEDULED
	targetState := apiv11.RUNNING

	initialPolicy := &apiv11.Policy{
		ID:           "policy-id",
		LastJobState: initialState,
	}
	finalPolicy := &apiv11.Policy{
		ID:           "policy-id",
		LastJobState: targetState,
	}

	t.Run("return false, err", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("get policy error")).Once()

		// Call the WaitForPolicyLastJobState method
		err := client.WaitForPolicyLastJobState(ctx, policyName, targetState)

		// Assertions
		assert.NotNil(t, err)
		assert.Equal(t, "get policy error", err.Error())
	})

	t.Run("return pollErr", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(ctx, defaultPoll) // Using a short timeout to force a timeout error
		defer cancel()

		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Twice()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*initialPolicy,
				},
			}
		}).Twice()

		// Call the WaitForPolicyLastJobState method
		err := client.WaitForPolicyLastJobState(ctx, policyName, targetState)

		// Assertions
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "timed out waiting for the condition")
	})

	t.Run("PollImmediate returns true eventually", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*initialPolicy,
				},
			}
		}).Once()

		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*finalPolicy,
				},
			}
		}).Once()

		// Call the WaitForPolicyLastJobState method
		err := client.WaitForPolicyLastJobState(ctx, policyName, targetState)

		// Assertions
		assert.Nil(t, err)
	})
}

func TestWaitForTargetPolicyCondition(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	initialState := apiv11.WritesDisabled
	targetState := apiv11.WritesEnabled

	initialPolicy := &apiv11.TargetPolicy{
		ID:                    "policy-id",
		FailoverFailbackState: initialState,
	}
	finalPolicy := &apiv11.TargetPolicy{
		ID:                    "policy-id",
		FailoverFailbackState: targetState,
	}

	// Mock GetTargetPolicyByName method
	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, initialState).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*initialPolicy,
			},
		}
	}).Once()

	client.API.(*mocks.Client).On("GetTargetPolicyByName", mock.Anything, targetState).Return("", nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv11.TargetPolicies)
		*resp = &apiv11.TargetPolicies{
			Policy: []apiv11.TargetPolicy{
				*finalPolicy,
			},
		}
	}).Once()

	// Call the WaitForTargetPolicyCondition method
	err := client.WaitForTargetPolicyCondition(ctx, policyName, targetState)

	// Assertions
	assert.Nil(t, err)
}

func TestSyncPolicy(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}
	policyName := "test-policy"
	policyID := "policy-id"
	maxRetries := 3 // Set the maxRetries variable for the test context
	enabledPolicy := &apiv11.Policy{
		ID:      policyName,
		Enabled: true,
	}
	disabledPolicy := &apiv11.Policy{
		ID:      policyName,
		Enabled: false,
	}
	resolvePolicyReq := &apiv11.ResolvePolicyReq{
		Conflicted: false,
		Enabled:    true,
	}
	runningJob := apiv11.Job{
		ID:     "running-job",
		Action: apiv11.SYNC,
	}
	expectedReports := &apiv11.Reports{
		Reports: []apiv11.Report{
			{
				Policy:  *enabledPolicy,
				ID:      "report-id-1",
				JobID:   12345,
				State:   "completed",
				EndTime: 1609459200,
				Errors:  []string{"error1", "error2"},
			},
			{
				Policy:  *disabledPolicy,
				ID:      "report-id-2",
				JobID:   12346,
				State:   "completed",
				EndTime: 1609459201,
				Errors:  []string{},
			},
		},
	}
	noJobs := []apiv11.Job{}
	successfulJob := &apiv11.Job{ID: "successful-job"}

	t.Run("GetPolicyByName returns error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("policy not found")).Once()

		err := client.SyncPolicy(ctx, policyName)

		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "policy not found")
	})
	t.Run("SyncPolicy handles GetJobsByPolicyName error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("get jobs error")).Once()
		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "get jobs error")
	})
	t.Run("SyncPolicy with enabled policy and no running jobs", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(nil).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.Nil(t, err)
	})

	t.Run("SyncPolicy with disabled policy", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*disabledPolicy,
				},
			}
		}).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot run sync on disabled policy %s", policyName))
	})

	t.Run("SyncPolicy with running SYNC jobs", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: []apiv11.Job{runningJob},
			}
		}).Once()
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(nil).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.Nil(t, err)
	})

	t.Run("SyncPolicy handles WaitForNoActiveJobs error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(nil).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("wait error")).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "wait error")
	})

	t.Run("Max Retries Reached for Retryable Error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: []apiv11.Job{},
			}
		}).Once()

		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(errors.New(retryablePolicyError)).Times(maxRetries)

		client.API.(*mocks.Client).On("GetReportsByPolicyName", mock.Anything, policyName, 1).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Reports)
			*resp = &apiv11.Reports{
				Reports: []apiv11.Report{
					{
						Errors: []string{retryablePolicyError},
					},
				},
			}
		}).Times(maxRetries - 1)

		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New(retryablePolicyError)).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), retryablePolicyError)
	})

	t.Run("SyncPolicy with maxRetries reached and no retryable error found", func(t *testing.T) {
		policyName := "test-policy"
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()

		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(errors.New(retryablePolicyError)).Times(maxRetries)

		client.API.(*mocks.Client).On("GetReportsByPolicyName", mock.Anything, policyName, 1).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("fake report error")).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "fake report error")
	})

	t.Run("SyncPolicy with retryable error and reports retrieval error", func(t *testing.T) {
		failureWithRetryableError := errors.New(retryablePolicyError)

		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(failureWithRetryableError).Once()
		client.API.(*mocks.Client).On("GetReportsByPolicyName", mock.Anything, policyName, 1).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("error while retrieving reports for failed sync job")).Once()
		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error while retrieving reports")
	})

	t.Run("SyncPolicy with retryable error when starting sync job", func(t *testing.T) {
		failureOnceWithRetryableError := errors.New(retryablePolicyError)
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: noJobs,
			}
		}).Once()
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, mock.Anything).Return(failureOnceWithRetryableError).Once()
		// Mock GetReportsByPolicyName method
		client.API.(*mocks.Client).On("GetReportsByPolicyName", mock.Anything, policyName, 1).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Reports)
			*resp = expectedReports
		}).Once()
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyID).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()

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
		client.API.(*mocks.Client).On("Post", ctx, jobsPath, "", mock.Anything, mock.Anything, &apiv11.JobRequest{ID: policyName}, successfulJob).Return(nil).Once()
		client.API.(*mocks.Client).On("WaitForNoActiveJobs", mock.Anything, policyName).Return(nil).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
	})

	t.Run("SyncPolicy logs and returns error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: []apiv11.Job{
					runningJob,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("get jobs error")).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "get jobs error")
	})

	t.Run("SyncPolicy calls WaitForNoActiveJobs and returns error", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetPolicyByName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Policies)
			*resp = &apiv11.Policies{
				Policy: []apiv11.Policy{
					*enabledPolicy,
				},
			}
		}).Once()
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = &apiv11.Jobs{
				Job: []apiv11.Job{runningJob},
			}
		}).Once()
		client.API.(*mocks.Client).On("WaitForNoActiveJobs", mock.Anything, policyName).Return(nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("wait error")).Once()

		err := client.SyncPolicy(ctx, policyName)
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "wait error")
	})
}

func TestGetJobsByPolicyName(t *testing.T) {
	ctx := context.Background()
	client := &Client{API: new(mocks.Client)}

	policyName := "test-policy"
	jobs := []apiv11.Job{
		{ID: "job-1", Action: "sync"},
		{ID: "job-2", Action: "sync"},
	}
	jobsResponse := &apiv11.Jobs{
		Job: jobs,
	}

	t.Run("Successfully retrieve jobs", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
			resp := args.Get(5).(**apiv11.Jobs)
			*resp = jobsResponse
		}).Once()
		result, err := client.GetJobsByPolicyName(ctx, policyName)
		assert.Nil(t, err)
		assert.Equal(t, jobs, result)
	})

	t.Run("Handle 404 error and return empty job list", func(t *testing.T) {
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(&api.JSONError{StatusCode: 404}).Once()
		result, err := client.GetJobsByPolicyName(ctx, policyName)
		assert.Nil(t, err)
		assert.Empty(t, result)
	})

	t.Run("Handle other errors", func(t *testing.T) {
		testError := errors.New("test error")
		client.API.(*mocks.Client).On("GetJobsByPolicyName", mock.Anything, policyName).Return("", nil).Once()
		client.API.(*mocks.Client).On("Get", anyArgs...).Return(testError).Once()
		result, err := client.GetJobsByPolicyName(ctx, policyName)
		assert.NotNil(t, err)
		assert.Equal(t, testError, err)
		assert.Nil(t, result)
	})
}

func TestFilterReports(t *testing.T) {
	reports := []apiv11.Report{
		{
			ID:      "report-1",
			State:   "completed",
			Errors:  []string{},
			EndTime: 1609459200,
		},
		{
			ID:      "report-2",
			State:   "failed",
			Errors:  []string{"error1"},
			EndTime: 1609459201,
		},
		{
			ID:      "report-3",
			State:   "completed",
			Errors:  []string{},
			EndTime: 1609459202,
		},
	}

	t.Run("Filter completed reports", func(t *testing.T) {
		filterFunc := func(report apiv11.Report) bool {
			return report.State == "completed"
		}
		expectedFiltered := []apiv11.Report{
			reports[0],
			reports[2],
		}

		filteredReports := FilterReports(reports, filterFunc)
		assert.Equal(t, expectedFiltered, filteredReports)
	})

	t.Run("Filter failed reports", func(t *testing.T) {
		filterFunc := func(report apiv11.Report) bool {
			return report.State == "failed"
		}
		expectedFiltered := []apiv11.Report{
			reports[1],
		}

		filteredReports := FilterReports(reports, filterFunc)
		assert.Equal(t, expectedFiltered, filteredReports)
	})

	t.Run("Filter reports with no errors", func(t *testing.T) {
		filterFunc := func(report apiv11.Report) bool {
			return len(report.Errors) == 0
		}
		expectedFiltered := []apiv11.Report{
			reports[0],
			reports[2],
		}

		filteredReports := FilterReports(reports, filterFunc)
		assert.Equal(t, expectedFiltered, filteredReports)
	})

	t.Run("Filter reports with specific end time", func(t *testing.T) {
		specificEndTime := int64(1609459201)
		filterFunc := func(report apiv11.Report) bool {
			return report.EndTime == specificEndTime
		}
		expectedFiltered := []apiv11.Report{
			reports[1],
		}

		filteredReports := FilterReports(reports, filterFunc)
		assert.Equal(t, expectedFiltered, filteredReports)
	})
}

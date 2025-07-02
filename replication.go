package goisilon

/*
Copyright (c) 2021-2023 Dell Inc, or its subsidiaries.

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

import (
	"context"
	"fmt"
	"strings"
	"time"

	log "github.com/akutz/gournal"
	"github.com/dell/goisilon/api/common/utils/poll"
	apiv11 "github.com/dell/goisilon/api/v11"
)

const (
	defaultPoll    = 5 * time.Second
	defaultTimeout = 10 * time.Minute // set high timeout, we expect to be canceled via context before
)

const (
	retryInterval        = 15 * time.Second
	maxRetries           = 20 // with 15 sec, it is 4 retries per min. For 5 minutes, it is 20 retries
	retryablePolicyError = "is in an error state. Please resolve it and retry"
	retryableReportError = "A new quota domain that has not finished QuotaScan has been found"
)

const (
	ResyncPrep           apiv11.JobAction             = "resync_prep"
	AllowWrite           apiv11.JobAction             = "allow_write"
	AllowWriteRevert     apiv11.JobAction             = "allow_write_revert"
	Test                 apiv11.JobAction             = "test"
	SCHEDULED            apiv11.JobState              = "scheduled"
	RUNNING              apiv11.JobState              = "running"
	PAUSED               apiv11.JobState              = "paused"
	FINISHED             apiv11.JobState              = "finished"
	FAILED               apiv11.JobState              = "failed"
	CANCELED             apiv11.JobState              = "canceled"
	NeedsAttention       apiv11.JobState              = "needs_attention"
	SKIPPED              apiv11.JobState              = "skipped"
	PENDING              apiv11.JobState              = "pending"
	UNKNOWN              apiv11.JobState              = "unknown"
	WritesDisabled       apiv11.FailoverFailbackState = "writes_disabled"
	EnablingWrites       apiv11.FailoverFailbackState = "enabling_writes"
	WritesEnabled        apiv11.FailoverFailbackState = "writes_enabled"
	DisablingWrites      apiv11.FailoverFailbackState = "disabling_writes"
	CreatingResyncPolicy apiv11.FailoverFailbackState = "creating_resync_policy"
	ResyncPolicyCreated  apiv11.FailoverFailbackState = "resync_policy_created"
)

// Policy is an Isilon Policy
type Policy *apiv11.Policy

type TargetPolicy *apiv11.TargetPolicy

// GetPolicyByName returns a policy with the provided ID.
func (c *Client) GetPolicyByName(ctx context.Context, id string) (Policy, error) {
	return apiv11.GetPolicyByName(ctx, c.API, id)
}

func (c *Client) GetTargetPolicyByName(ctx context.Context, id string) (TargetPolicy, error) {
	return apiv11.GetTargetPolicyByName(ctx, c.API, id)
}

func (c *Client) CreatePolicy(ctx context.Context, name string, rpo int, sourcePath string, targetPath string, targetHost string, targetCert string, enabled bool) error {
	return apiv11.CreatePolicy(ctx, c.API, name, sourcePath, targetPath, targetHost, targetCert, rpo, enabled)
}

func (c *Client) DeletePolicy(ctx context.Context, name string) error {
	return apiv11.DeletePolicy(ctx, c.API, name)
}

func (c *Client) DeleteTargetPolicy(ctx context.Context, id string) error {
	return apiv11.DeleteTargetPolicy(ctx, c.API, id)
}

func (c *Client) BreakAssociation(ctx context.Context, targetPolicyName string) error {
	tp, err := apiv11.GetTargetPolicyByName(ctx, c.API, targetPolicyName)
	if err != nil {
		return err
	}

	return c.DeleteTargetPolicy(ctx, tp.ID)
}

func (c *Client) ResetPolicy(ctx context.Context, name string) error {
	return apiv11.ResetPolicy(ctx, c.API, name)
}

func (c *Client) EnablePolicy(ctx context.Context, name string) error {
	return c.SetPolicyEnabledField(ctx, name, true)
}

func (c *Client) DisablePolicy(ctx context.Context, name string) error {
	return c.SetPolicyEnabledField(ctx, name, false)
}

func (c *Client) SetPolicyEnabledField(ctx context.Context, name string, value bool) error {
	pp, err := c.GetPolicyByName(ctx, name)
	if err != nil {
		return err
	}
	if pp == nil {
		return nil
	}

	if pp.Enabled == value {
		return nil
	}

	p := &apiv11.Policy{
		ID:       pp.ID,
		Enabled:  value,
		Schedule: pp.Schedule, // keep existing schedule, otherwise it will be cleared
	}

	return apiv11.UpdatePolicy(ctx, c.API, p)
}

// Modifies the policy schedule and job delay.
// scdeule - can be either empty string (manual) or when-source-modified
// rpo - can be 0 (manual) or in seconds (when-source-modified)
func (c *Client) ModifyPolicy(ctx context.Context, name string, schedule string, rpo int) error {
	pp, err := c.GetPolicyByName(ctx, name)
	if err != nil {
		return err
	}

	p := &apiv11.Policy{
		ID:      pp.ID,
		Enabled: pp.Enabled, // keep existing enabled state, otherwise it will be cleared
	}

	if schedule == "" { // manual
		p.Schedule = schedule
	} else if schedule == "when-source-modified" {
		p.Schedule = schedule
		p.JobDelay = rpo
	}

	return apiv11.UpdatePolicy(ctx, c.API, p)
}

// Resolves the policy that is in error state due to QuotaScan requirement.
func (c *Client) ResolvePolicy(ctx context.Context, name string) error {
	pp, err := c.GetPolicyByName(ctx, name)
	if err != nil {
		return err
	}

	p := &apiv11.ResolvePolicyReq{
		ID:         pp.ID,
		Conflicted: false,
		Enabled:    pp.Enabled,  // keep existing enabled state, otherwise it will be cleared
		Schedule:   pp.Schedule, // keep existing schedule, otherwise it will be cleared
	}

	return apiv11.ResolvePolicy(ctx, c.API, p)
}

func (c *Client) AllowWrites(ctx context.Context, policyName string) error {
	targetPolicy, err := c.GetTargetPolicyByName(ctx, policyName)
	if err != nil {
		return err
	}
	if targetPolicy.FailoverFailbackState == WritesEnabled {
		return nil
	}

	_, err = c.RunActionForPolicy(ctx, policyName, apiv11.AllowWrite)
	if err != nil {
		return err
	}

	err = c.WaitForTargetPolicyCondition(ctx, policyName, WritesEnabled)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DisallowWrites(ctx context.Context, policyName string) error {
	targetPolicy, err := c.GetTargetPolicyByName(ctx, policyName)
	if err != nil {
		return err
	}
	if targetPolicy.FailoverFailbackState == WritesDisabled {
		return nil
	}

	_, err = c.RunActionForPolicy(ctx, policyName, apiv11.AllowWriteRevert)
	if err != nil {
		return err
	}

	err = c.WaitForTargetPolicyCondition(ctx, policyName, WritesDisabled)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ResyncPrep(ctx context.Context, policyName string) error {
	_, err := c.RunActionForPolicy(ctx, policyName, apiv11.ResyncPrep)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) RunActionForPolicy(ctx context.Context, policyName string, action apiv11.JobAction) (*apiv11.Job, error) {
	job := &apiv11.JobRequest{
		ID:     policyName,
		Action: action,
	}

	return apiv11.StartSyncIQJob(ctx, c.API, job)
}

func (c *Client) StartSyncIQJob(ctx context.Context, job *apiv11.JobRequest) (*apiv11.Job, error) {
	return apiv11.StartSyncIQJob(ctx, c.API, job)
}

func (c *Client) GetReport(ctx context.Context, reportName string) (*apiv11.Report, error) {
	return apiv11.GetReport(ctx, c.API, reportName)
}

func (c *Client) GetReportsByPolicyName(ctx context.Context, policyName string, reportsForPolicy int) (*apiv11.Reports, error) {
	return apiv11.GetReportsByPolicyName(ctx, c.API, policyName, reportsForPolicy)
}

func (c *Client) WaitForPolicyEnabledFieldCondition(ctx context.Context, policyName string, enabled bool) error {
	pollErr := poll.PollImmediateWithContext(ctx, defaultPoll, defaultTimeout,
		func(iCtx context.Context) (bool, error) {
			p, err := c.GetPolicyByName(iCtx, policyName)
			if err != nil {
				return false, err
			}

			if p.Enabled != enabled {
				return false, nil
			}

			return true, nil
		})

	if pollErr != nil {
		return pollErr
	}

	return nil
}

func (c *Client) WaitForNoActiveJobs(ctx context.Context, policyName string) error {
	pollErr := poll.PollImmediateWithContext(ctx, defaultPoll, defaultTimeout,
		func(iCtx context.Context) (bool, error) {
			p, err := c.GetJobsByPolicyName(iCtx, policyName)
			if err != nil {
				return false, err
			}

			if len(p) != 0 {
				return false, nil
			}

			return true, nil
		})

	if pollErr != nil {
		return pollErr
	}

	return nil
}

func (c *Client) WaitForPolicyLastJobState(ctx context.Context, policyName string, state apiv11.JobState) error {
	pollErr := poll.PollImmediateWithContext(ctx, defaultPoll, defaultTimeout,
		func(iCtx context.Context) (bool, error) {
			p, err := c.GetPolicyByName(iCtx, policyName)
			if err != nil {
				return false, err
			}

			if p.LastJobState != state {
				return false, nil
			}

			return true, nil
		})

	if pollErr != nil {
		return pollErr
	}

	return nil
}

func (c *Client) WaitForTargetPolicyCondition(ctx context.Context, policyName string, condition apiv11.FailoverFailbackState) error {
	pollErr := poll.PollImmediateWithContext(ctx, defaultPoll, defaultTimeout,
		func(iCtx context.Context) (bool, error) {
			tp, err := c.GetTargetPolicyByName(iCtx, policyName)
			if err != nil {
				return false, err
			}

			if tp.FailoverFailbackState != condition {
				return false, nil
			}

			return true, nil
		})

	if pollErr != nil {
		return pollErr
	}

	return nil
}

func (c *Client) SyncPolicy(ctx context.Context, policyName string) error {
	// get all running
	// if running - wait for it and succeed
	// if no running - start new - wait for it and succeed

	var isRunning bool

	policy, err := c.GetPolicyByName(ctx, policyName)
	if err != nil {
		return err
	}
	if !policy.Enabled {
		return fmt.Errorf("cannot run sync on disabled policy %s", policyName)
	}

	runningJobs, err := c.GetJobsByPolicyName(ctx, policyName)
	if err != nil {
		log.Info(ctx, err.Error())
		return err
	}
	for _, i := range runningJobs {
		if i.Action == apiv11.SYNC {
			// running job detected. Wait for it to complete.
			isRunning = true
		}
	}
	if isRunning {
		log.Info(ctx, "found active jobs, waiting for completion")
		err = c.WaitForNoActiveJobs(ctx, policyName)
		if err != nil {
			return err
		}
		return nil
	}
	jobReq := &apiv11.JobRequest{
		ID: policyName,
	}
	log.Info(ctx, "found no active sync jobs, starting a new one")

	// workaround for PowerScale KB article
	// https://www.dell.com/support/kbdoc/en-us/000019414/quotas-on-synciq-source-directories
	for i := 0; i < maxRetries; i++ {
		_, err := c.StartSyncIQJob(ctx, jobReq)
		if err == nil {
			break
		}
		if strings.Contains(err.Error(), retryablePolicyError) {
			if i+1 == maxRetries {
				return err
			}

			reports, err := c.GetReportsByPolicyName(ctx, policyName, 1)
			if err != nil {
				return fmt.Errorf("error while retrieving reports for failed sync job %s %s", policyName, err.Error())
			}
			if !(len(reports.Reports) > 0 && len(reports.Reports[0].Errors) > 0 &&
				strings.Contains(reports.Reports[0].Errors[0], retryableReportError)) {
				return fmt.Errorf("found no retryable error in reports for failed sync job %s", policyName)
			}

			log.Info(ctx, "Sync job failed with error: %s. %v of %v - retrying in %v...",
				reports.Reports[0].Errors[0], i+1, maxRetries, retryInterval)
			time.Sleep(retryInterval)

			// Resolve policy with error before retrying
			err = c.ResolvePolicy(ctx, policyName)
			if err != nil {
				return err
			}
		} else { // not a retryable error
			return err
		}
	}

	time.Sleep(3 * time.Second)
	err = c.WaitForNoActiveJobs(ctx, policyName)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) GetJobsByPolicyName(ctx context.Context, policyName string) ([]apiv11.Job, error) {
	return apiv11.GetJobsByPolicyName(ctx, c.API, policyName)
}

func FilterReports(values []apiv11.Report, filterFunc func(apiv11.Report) bool) []apiv11.Report {
	filtered := make([]apiv11.Report, 0)
	for _, v := range values {
		if filterFunc(v) {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

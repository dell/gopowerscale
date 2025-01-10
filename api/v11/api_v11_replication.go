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
package v11

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dell/goisilon/api"
)

const (
	policiesPath       = "/platform/11/sync/policies/"
	targetPoliciesPath = "/platform/11/sync/target/policies/"
	jobsPath           = "/platform/11/sync/jobs/"
	reportsPath        = "/platform/11/sync/reports"
)

type JobAction string

const (
	ResyncPrep       JobAction = "resync_prep"
	AllowWrite       JobAction = "allow_write"
	AllowWriteRevert JobAction = "allow_write_revert"
	Test             JobAction = "test"
)

type RunningJobAction string

const (
	SYNC RunningJobAction = "sync"
	COPY RunningJobAction = "copy"
)

type JobState string

const (
	SCHEDULED      JobState = "scheduled"
	RUNNING        JobState = "running"
	PAUSED         JobState = "paused"
	FINISHED       JobState = "finished"
	FAILED         JobState = "failed"
	CANCELED       JobState = "canceled"
	NeedsAttention JobState = "needs_attention"
	SKIPPED        JobState = "skipped"
	PENDING        JobState = "pending"
	UNKNOWN        JobState = "unknown"
)

type FailoverFailbackState string

const (
	WritesDisabled       FailoverFailbackState = "writes_disabled"
	EnablingWrites       FailoverFailbackState = "enabling_writes"
	WritesEnabled        FailoverFailbackState = "writes_enabled"
	DisablingWrites      FailoverFailbackState = "disabling_writes"
	CreatingResyncPolicy FailoverFailbackState = "creating_resync_policy"
	ResyncPolicyCreated  FailoverFailbackState = "resync_policy_created"
)

const resolveErrorToIgnore = "The policy was not conflicted, so no change was made"

var (
	policyNameArg       = []byte("policy_name")
	sortArg             = []byte("sort")
	reportsPerPolicyArg = []byte("reports_per_policy")
)

// Policy contains the CloudIQ policy info.
type Policy struct {
	Action       string   `json:"action,omitempty"`
	ID           string   `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	Enabled      bool     `json:"enabled"`
	TargetPath   string   `json:"target_path,omitempty"`
	SourcePath   string   `json:"source_root_path,omitempty"`
	TargetHost   string   `json:"target_host,omitempty"`
	TargetCert   string   `json:"target_certificate_id,omitempty"`
	JobDelay     int      `json:"job_delay,omitempty"`
	Schedule     string   `json:"schedule"`
	LastJobState JobState `json:"last_job_state,omitempty"`
}

type ResolvePolicyReq struct {
	ID         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Enabled    bool   `json:"enabled"`
	Conflicted bool   `json:"conflicted"`
	Schedule   string `json:"schedule"`
}

type Policies struct {
	Policy []Policy `json:"policies,omitempty"`
}

type Reports struct {
	Reports []Report `json:"reports,omitempty"`
}

type TargetPolicy struct {
	ID                      string                `json:"id,omitempty"`
	Name                    string                `json:"name,omitempty"`
	SourceClusterGUID       string                `json:"source_cluster_guid,omitempty"`
	LastJobState            JobState              `json:"last_job_state,omitempty"`
	TargetPath              string                `json:"target_path,omitempty"`
	SourceHost              string                `json:"source_host,omitempty"`
	LastSourceCoordinatorIP string                `json:"last_source_coordinator_ip,omitempty"`
	FailoverFailbackState   FailoverFailbackState `json:"failover_failback_state,omitempty"`
}

type TargetPolicies struct {
	Policy []TargetPolicy `json:"policies,omitempty"`
}

type JobRequest struct {
	Action       JobAction `json:"action,omitempty"`
	ID           string    `json:"id,omitempty"` // ID or Name of policy
	SkipFailover bool      `json:"skip_failover,omitempty"`
	SkipMap      bool      `json:"skip_map,omitempty"`
	SkipCopy     bool      `json:"skip_copy,omitempty"`
}

type Job struct {
	Action RunningJobAction `json:"policy_action,omitempty"`
	ID     string           `json:"id,omitempty"` // ID or Name of policy
}

type Jobs struct {
	Job []Job `json:"jobs,omitempty"`
}

type Report struct {
	Policy  Policy   `json:"policy,omitempty"`
	ID      string   `json:"id"`
	JobID   int64    `json:"job_id"`
	State   JobState `json:"state,omitempty"`
	EndTime int64    `json:"end_time"`
	Errors  []string `json:"errors"`
}

// GetPolicyByName returns policy by name
func GetPolicyByName(ctx context.Context, client api.Client, name string) (policy *Policy, err error) {
	p := &Policies{}
	err = client.Get(ctx, policiesPath, name, nil, nil, &p)
	if err != nil {
		return nil, err
	} else if len(p.Policy) == 0 {
		return nil, fmt.Errorf("successful code returned, but policy %s not found", name)
	}
	return &p.Policy[0], nil
}

func GetTargetPolicyByName(ctx context.Context, client api.Client, name string) (policy *TargetPolicy, err error) {
	p := &TargetPolicies{}
	err = client.Get(ctx, targetPoliciesPath, name, nil, nil, &p)
	if err != nil || len(p.Policy) == 0 {
		return nil, err
	}
	return &p.Policy[0], nil
}

func CreatePolicy(ctx context.Context, client api.Client, name string, sourcePath string, targetPath string, targetHost string, targetCert string, rpo int, enabled bool) error {
	var policyResp Policy
	body := &Policy{
		Action:     "sync",
		ID:         "",
		Name:       name,
		Enabled:    enabled,
		TargetPath: targetPath,
		SourcePath: sourcePath,
		TargetHost: targetHost,
		JobDelay:   rpo,
		TargetCert: targetCert,
		Schedule:   "when-source-modified",
	}
	return client.Post(ctx, policiesPath, "", nil, nil, body, &policyResp)
}

func DeletePolicy(ctx context.Context, client api.Client, name string) error {
	resp := ""
	return client.Delete(ctx, policiesPath, name, nil, nil, &resp)
}

func DeleteTargetPolicy(ctx context.Context, client api.Client, id string) error {
	resp := ""
	return client.Delete(ctx, targetPoliciesPath, id, nil, nil, &resp)
}

func UpdatePolicy(ctx context.Context, client api.Client, policy *Policy) error {
	id := policy.ID
	policy.ID = ""

	return client.Put(ctx, policiesPath, id, nil, nil, policy, nil)
}

func ResolvePolicy(ctx context.Context, client api.Client, policy *ResolvePolicyReq) error {
	id := policy.ID
	policy.ID = ""

	err := client.Put(ctx, policiesPath, id, nil, nil, policy, nil)
	if err != nil && !strings.Contains(err.Error(), resolveErrorToIgnore) {
		return err
	}
	return nil
}

func ResetPolicy(ctx context.Context, client api.Client, name string) error {
	resp := Policy{}
	return client.Post(ctx, policiesPath, name+"/reset", nil, nil, nil, &resp)
}

func StartSyncIQJob(ctx context.Context, client api.Client, job *JobRequest) (*Job, error) {
	var jobResp Job
	return &jobResp, client.Post(ctx, jobsPath, "", nil, nil, job, &jobResp)
}

func GetReport(ctx context.Context, client api.Client, reportName string) (*Report, error) {
	r := &Reports{}
	err := client.Get(ctx, reportsPath, reportName, nil, nil, &r)
	if err != nil {
		return nil, err
	}
	if len(r.Reports) == 0 {
		return nil, fmt.Errorf("no reports found with report name %s", reportName)
	}
	return &r.Reports[0], nil
}

func GetReportsByPolicyName(ctx context.Context, client api.Client, policyName string, reportsForPolicy int) (*Reports, error) {
	r := &Reports{}
	err := client.Get(ctx, reportsPath, "",
		api.OrderedValues{
			{policyNameArg, []byte(policyName)},
			{sortArg, []byte("end_time")},
			{reportsPerPolicyArg, []byte(strconv.Itoa(reportsForPolicy))},
		},
		nil, &r)
	if err != nil {
		return nil, err
	}

	if len(r.Reports) == 0 {
		return nil, fmt.Errorf("no reports found for policy %s", policyName)
	}

	return r, nil
}

func GetJobsByPolicyName(ctx context.Context, client api.Client, policyName string) ([]Job, error) {
	j := &Jobs{}
	err := client.Get(ctx, jobsPath, policyName, nil, nil, &j)
	if err != nil {
		if e, ok := err.(*api.JSONError); ok {
			if e.StatusCode == 404 {
				return []Job{}, nil
			}
		}
		return nil, err
	}
	return j.Job, nil
}

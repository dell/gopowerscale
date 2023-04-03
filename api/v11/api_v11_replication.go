/*
Copyright (c) 2022-2023 Dell Inc, or its subsidiaries.

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

	"github.com/dell/goisilon/api"
)

const (
	policiesPath       = "/platform/11/sync/policies/"
	targetPoliciesPath = "/platform/11/sync/target/policies/"
	jobsPath           = "/platform/11/sync/jobs/"
	reportsPath        = "/platform/11/sync/reports"
)

type JOB_ACTION string

const (
	RESYNC_PREP        JOB_ACTION = "resync_prep"
	ALLOW_WRITE        JOB_ACTION = "allow_write"
	ALLOW_WRITE_REVERT JOB_ACTION = "allow_write_revert"
	TEST               JOB_ACTION = "test"
)

type RUNNING_JOB_ACTION string

const (
	SYNC RUNNING_JOB_ACTION = "sync"
	COPY RUNNING_JOB_ACTION = "copy"
)

type JOB_STATE string

const (
	SCHEDULED       JOB_STATE = "scheduled"
	RUNNING         JOB_STATE = "running"
	PAUSED          JOB_STATE = "paused"
	FINISHED        JOB_STATE = "finished"
	FAILED          JOB_STATE = "failed"
	CANCELED        JOB_STATE = "canceled"
	NEEDS_ATTENTION JOB_STATE = "needs_attention"
	SKIPPED         JOB_STATE = "skipped"
	PENDING         JOB_STATE = "pending"
	UNKNOWN         JOB_STATE = "unknown"
)

type FAILOVER_FAILBACK_STATE string

const (
	WRITES_DISABLED        FAILOVER_FAILBACK_STATE = "writes_disabled"
	ENABLING_WRITES        FAILOVER_FAILBACK_STATE = "enabling_writes"
	WRITES_ENABLED         FAILOVER_FAILBACK_STATE = "writes_enabled"
	DISABLING_WRITES       FAILOVER_FAILBACK_STATE = "disabling_writes"
	CREATING_RESYNC_POLICY FAILOVER_FAILBACK_STATE = "creating_resync_policy"
	RESYNC_POLICY_CREATED  FAILOVER_FAILBACK_STATE = "resync_policy_created"
)

var policyNameArg = []byte("policy_name")
var sortArg = []byte("sort")
var reportsPerPolicyArg = []byte("reports_per_policy")

// Policy contains the CloudIQ policy info.
type Policy struct {
	Action       string    `json:"action,omitempty"`
	Id           string    `json:"id,omitempty"`
	Name         string    `json:"name,omitempty"`
	Enabled      bool      `json:"enabled"`
	Conflicted   bool      `json:"conflicted"`
	TargetPath   string    `json:"target_path,omitempty"`
	SourcePath   string    `json:"source_root_path,omitempty"`
	TargetHost   string    `json:"target_host,omitempty"`
	TargetCert   string    `json:"target_certificate_id,omitempty"`
	JobDelay     int       `json:"job_delay,omitempty"`
	Schedule     string    `json:"schedule"`
	LastJobState JOB_STATE `json:"last_job_state,omitempty"`
}

type Policies struct {
	Policy []Policy `json:"policies,omitempty"`
}

type Reports struct {
	Reports []Report `json:"reports,omitempty"`
}

type TargetPolicy struct {
	Id                      string                  `json:"id,omitempty"`
	Name                    string                  `json:"name,omitempty"`
	SourceClusterGuid       string                  `json:"source_cluster_guid,omitempty"`
	LastJobState            JOB_STATE               `json:"last_job_state,omitempty"`
	TargetPath              string                  `json:"target_path,omitempty"`
	SourceHost              string                  `json:"source_host,omitempty"`
	LastSourceCoordinatorIp string                  `json:"last_source_coordinator_ip,omitempty"`
	FailoverFailbackState   FAILOVER_FAILBACK_STATE `json:"failover_failback_state,omitempty"`
}

type TargetPolicies struct {
	Policy []TargetPolicy `json:"policies,omitempty"`
}

type JobRequest struct {
	Action       JOB_ACTION `json:"action,omitempty"`
	Id           string     `json:"id,omitempty"` // ID or Name of policy
	SkipFailover bool       `json:"skip_failover,omitempty"`
	SkipMap      bool       `json:"skip_map,omitempty"`
	SkipCopy     bool       `json:"skip_copy,omitempty"`
}

type Job struct {
	Action RUNNING_JOB_ACTION `json:"policy_action,omitempty"`
	Id     string             `json:"id,omitempty"` // ID or Name of policy
}

type Jobs struct {
	Job []Job `json:"jobs,omitempty"`
}

type Report struct {
	Policy  Policy    `json:"policy,omitempty"`
	Id      string    `json:"id"`
	JobId   int64     `json:"job_id"`
	State   JOB_STATE `json:"state,omitempty"`
	EndTime int64     `json:"end_time"`
	Errors  []string  `json:"errors"`
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
		Id:         "",
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
	id := policy.Id
	policy.Id = ""

	return client.Put(ctx, policiesPath, id, nil, nil, policy, nil)
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

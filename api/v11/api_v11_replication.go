package v11

import (
	"context"
	"fmt"
	"github.com/dell/goisilon/api"
)

const (
	policiesPath = "/platform/11/sync/policies/"
)

// Policy contains the CloudIQ policy info.
type Policy struct {
	Action     string `json:"action,omitempty"`
	Id         string `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Enabled    bool   `json:"enabled,omitempty"`
	TargetPath string `json:"target_path,omitempty"`
	SourcePath string `json:"source_root_path,omitempty"`
	TargetHost string `json:"target_host,omitempty"`
	JobDelay   int    `json:"job_delay,omitempty"`
	Schedule   string `json:"schedule,omitempty"`
}

type Policies struct {
	Policy []Policy `json:"policies,omitempty"`
}

// GetPolicyByName returns policy by name
func GetPolicyByName(
	ctx context.Context,
	client api.Client, name string) (policy *Policy, err error) {
	p := &Policies{}
	err = client.Get(ctx, policiesPath, name, nil, nil, &p)
	if err != nil {
		return nil, err
	} else if len(p.Policy) == 0 {
		return nil, fmt.Errorf("successful code returned, but policy %s not found", name)
	}
	return &p.Policy[0], nil
}

func CreatePolicy(ctx context.Context, client api.Client, name string, sourcePath string, targetPath string, targetHost string, rpo int) error {
	resp := ""
	body := &Policy{
		Action:     "sync",
		Id:         "",
		Name:       name,
		Enabled:    true,
		TargetPath: targetPath,
		SourcePath: sourcePath,
		TargetHost: targetHost,
		JobDelay:   rpo,
		Schedule:   "when-source-modified",
	}
	return client.Post(ctx, policiesPath, "", nil, nil, body, resp)
}

func DeletePolicy(ctx context.Context, client api.Client, name string) error {
	resp := ""
	return client.Delete(ctx, policiesPath, name, nil, nil, &resp)
}
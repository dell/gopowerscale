package goisilon

import (
	"context"
	apiv11 "github.com/dell/goisilon/api/v11"
)

// Policy is an Isilon Policy
type Policy *apiv11.Policy

// GetPolicyByName returns a policy with the provided ID.
func (c *Client) GetPolicyByName(ctx context.Context, id string) (Policy, error) {
	return apiv11.GetPolicyByName(ctx,c.API,id)
}

func (c *Client) CreatePolicy(ctx context.Context, name string, rpo int, sourcePath string, targetPath string, targetHost string) error{
	return apiv11.CreatePolicy(ctx,c.API,name,sourcePath,targetPath,targetHost,rpo)
}
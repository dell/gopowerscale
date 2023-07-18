package v14

import (
	"context"
	"github.com/dell/goisilon/api"
)

// GetIsiClusterAcs queries ACS status of OneFS cluster
func GetIsiClusterAcs(
	ctx context.Context,
	client api.Client) (clusterAcs *IsiClusterAcs, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/14/cluster/acs
	// This will return ACS status.
	var clusterAcsResp IsiClusterAcs
	err = client.Get(ctx, clusterAcsPath, "", nil, nil, &clusterAcsResp)
	if err != nil {
		return nil, err
	}

	return &clusterAcsResp, nil
}

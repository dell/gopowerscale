package v7

import (
	"context"

	"github.com/dell/goisilon/api"
)

// GetIsiClusterInternalNetworks queries internal networks settings
func GetIsiClusterInternalNetworks(
	ctx context.Context,
	client api.Client,
) (clusterInternalNetworks *IsiClusterInternalNetworks, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/7/cluster/internal-networks
	// This will return the internal networks settings
	var clusterInternalNetworksResp IsiClusterInternalNetworks
	err = client.Get(ctx, clusterInternalNetworksPath, "", nil, nil, &clusterInternalNetworksResp)
	if err != nil {
		return nil, err
	}

	return &clusterInternalNetworksResp, nil
}

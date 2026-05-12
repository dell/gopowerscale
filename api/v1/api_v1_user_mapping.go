package v1

import (
	"context"

	"github.com/dell/gopowerscale/api"
)

// GetIsiUserMapping queries the onefs array for user mappings, filter by user name and access zone.
func GetIsiUserMapping(ctx context.Context, client api.Client, user, accessZone *string) (userLookupResp *UserLookup, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/auth/mapping/users/lookup?user=&zone=
	values := api.OrderedValues{}
	if user != nil {
		values.StringAdd("user", *user)
	}
	if accessZone != nil {
		values.StringAdd("zone", *accessZone)
	}

	if err = client.Get(ctx, userLookupPath, "", values, nil, &userLookupResp); err != nil {
		return userLookupResp, err
	}

	return userLookupResp, err
}

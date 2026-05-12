package gopowerscale

import (
	"context"

	api "github.com/dell/gopowerscale/api/v1"
)

type UserMappingLookup *api.UserLookup

func (c *Client) GetUserMapping(ctx context.Context,
	userName, accessZoneDomain *string,
) (UserMappingLookup, error) {
	UserMappingLookup, err := api.GetIsiUserMapping(ctx, c.API, userName, accessZoneDomain)
	return UserMappingLookup, err
}

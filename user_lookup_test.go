package gopowerscale

import (
	"context"
	"fmt"
	"testing"

	apiv1 "github.com/dell/gopowerscale/api/v1"
	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.UserLookup)
		*resp = &apiv1.UserLookup{}
	}).Once()
	userName := "user"
	accessZone := "accesszone"
	_, err = client.GetUserMapping(context.Background(), &userName, &accessZone)
	assert.Nil(t, err)

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(fmt.Errorf("not found")).Once()
	userName = "wronguser"
	accessZone = "wrongaccesszone"
	_, err = client.GetUserMapping(context.Background(), &userName, &accessZone)
	assert.NotNil(t, err)
}

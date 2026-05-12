package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/gopowerscale/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetIsiUserMapping(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	userName := "wronguser"
	accessZone := "wrongaccesszone"
	_, err := GetIsiUserMapping(ctx, client, &userName, &accessZone)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Twice()
	userName = "test"
	accessZone = "teszone"
	_, err = GetIsiUserMapping(ctx, client, &userName, &accessZone)
	assert.Equal(t, nil, err)
}

package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetZoneByName(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetZoneByName(ctx, client, "name")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
	client.ExpectedCalls = nil
	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err = GetZoneByName(ctx, client, "name")
	assert.Equal(t, nil, err)
}

package v7

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestGetIsiClusterInternalNetworks(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiClusterInternalNetworks(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

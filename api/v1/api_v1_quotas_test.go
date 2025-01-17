package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestGetIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiQuota(ctx, client, "")
	assert.Equal(t, errors.New("Quota not found: "), err)
}

func TestGetAllIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Twice()
	_, err := GetAllIsiQuota(ctx, client)
	assert.Equal(t, errors.New("error"), err)
}

func TestGetIsiQuotaWithResume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Twice()
	_, err := GetIsiQuotaWithResume(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestGetIsiQuotaByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error")).Twice()
	_, err := GetIsiQuotaByID(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestSetIsiQuotaHardThreshold(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Post", anyArgs...).Return(nil).Twice()
	_, err := SetIsiQuotaHardThreshold(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
}

func TestUpdateIsiQuotaHardThreshold(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	err := UpdateIsiQuotaHardThreshold(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, errors.New("Quota not found: "), err)
}

func TestUpdateIsiQuotaHardThresholdByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	err := UpdateIsiQuotaHardThresholdByID(ctx, client, "", 5, 0, 0, 0)
	assert.Equal(t, nil, err)
}

func TestDeleteIsiQuota(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuota(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiQuotaByID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuotaByID(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiQuotaByIDWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := DeleteIsiQuotaByIDWithZone(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

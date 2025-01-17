package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetIsiSnapshots(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiSnapshots(ctx, client)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	var x int64 = 0
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiSnapshot(ctx, client, x)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiSnapshotByIdentity(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiSnapshotByIdentity(ctx, client, "")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestCreateIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	_, err := CreateIsiSnapshot(ctx, client, "", "name")
	assert.Equal(t, errors.New("no path set"), err)

	client.On("Post", anyArgs...).Return(nil).Twice()
	_, err = CreateIsiSnapshot(ctx, client, "path", "name")
	assert.Equal(t, nil, err)
}

func TestRemoveIsiSnapshot(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	var id int64 = 0
	client.On("Delete", anyArgs...).Return(nil).Twice()
	err := RemoveIsiSnapshot(ctx, client, id)
	assert.Equal(t, nil, err)
}

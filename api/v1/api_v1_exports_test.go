package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	err := Export(ctx, client, "")
	assert.Equal(t, errors.New("no path set"), err)

	client.On("User", anyArgs...).Return("").Twice()
	client.On("Group", anyArgs...).Return("").Twice()
	client.On("Post", anyArgs...).Return(nil).Twice()
	err = Export(ctx, client, "path")
	assert.Equal(t, nil, err)
}

func TestSetExportClients(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	err := SetExportClients(ctx, client, 0, []string{""})
	assert.Equal(t, nil, err)
}

func TestUnexport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	err := Unexport(ctx, client, 0)
	assert.Equal(t, errors.New("no path Id set"), err)

	client.On("Delete", anyArgs...).Return(nil).Twice()
	err = Unexport(ctx, client, 1)
	assert.Equal(t, nil, err)
}

func TestGetIsiExports(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiExports(ctx, client)
	assert.Equal(t, nil, err)
}

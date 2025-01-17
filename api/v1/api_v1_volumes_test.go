package v1

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestCreateIsiVolume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Put", anyArgs...).Return(nil).Twice()
	_, err := CreateIsiVolume(ctx, client, "name")
	assert.Equal(t, nil, err)
}

func TestCreateIsiVolumeWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Put", anyArgs...).Return(nil).Twice()
	_, err := CreateIsiVolumeWithIsiPath(ctx, client, "path", "name", "")
	assert.Equal(t, nil, err)
}

func TestCreateIsiVolumeWithIsiPathMetaData(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	metadata := make(map[string]string)
	metadata["path"] = "/path/path"
	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Put", anyArgs...).Return(nil).Twice()
	_, err := CreateIsiVolumeWithIsiPathMetaData(ctx, client, "path", "name", "", metadata)
	assert.Equal(t, nil, err)
}

func TestGetIsiVolumeWithSize(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiVolumeWithSize(ctx, client, "path", "name")
	assert.Equal(t, nil, err)
}

func TestCopyIsiVolumeWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Put", anyArgs...).Return(nil).Twice()
	_, err := CopyIsiVolumeWithIsiPath(ctx, client, "", "", "")
	assert.Equal(t, nil, err)
}

func TestCopyIsiVolume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Put", anyArgs...).Return(nil).Twice()
	_, err := CopyIsiVolume(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiVolumeWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Twice()
	_, err := DeleteIsiVolumeWithIsiPath(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

func TestDeleteIsiVolume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Delete", anyArgs...).Return(nil).Twice()
	_, err := DeleteIsiVolume(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetIsiVolumeWithoutMetadataWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	err := GetIsiVolumeWithoutMetadataWithIsiPath(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

func TestGetIsiVolumeWithoutMetadata(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Get", anyArgs...).Return(nil).Twice()
	err := GetIsiVolumeWithoutMetadata(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetIsiVolumeWithIsiPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiVolumeWithIsiPath(ctx, client, "", "")
	assert.Equal(t, nil, err)
}

func TestGetIsiVolume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiVolume(ctx, client, "")
	assert.Equal(t, nil, err)
}

func TestGetIsiVolumes(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("VolumesPath", anyArgs...).Return("").Twice()
	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := GetIsiVolumes(ctx, client)
	assert.Equal(t, nil, err)
}

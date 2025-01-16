package v5

import (
	"context"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestIsQuotaLicenseActivated(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Twice()
	_, err := IsQuotaLicenseActivated(ctx, client)
	assert.Equal(t, nil, err)
}

func TestIsQuotaLicenseStatusValid(t *testing.T) {
	licenseStatus := QuotaLicenseStatus{
		value: "Expired",
	}
	value := isQuotaLicenseStatusValid(licenseStatus.toString())
	assert.Equal(t, true, value)
}

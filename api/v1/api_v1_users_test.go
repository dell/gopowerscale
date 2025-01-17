package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	userName := "userName"
	var uid int32 = 0
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiUser(ctx, client, &userName, &uid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiUserList(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryMemberOf := false
	queryDomain := "queryDomain"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32 = 0

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiUserList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone,
		&queryProvider, &queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestCreateIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32 = 0

	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := CreateIsiUser(ctx, client, queryNamePrefix, &queryCached, &queryZone,
		&queryProvider, &queryZone, &queryZone, &queryZone, &queryZone, &queryZone,
		&queryZone, &queryLimit, &queryLimit, &queryLimit, &queryResolveNames,
		&queryResolveNames, &queryResolveNames, &queryResolveNames)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestUpdateIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryProvider := "queryProvider"
	queryNamePrefix := "queryNamePrefix"
	queryZone := "queryZone"
	queryCached := false
	queryResolveNames := false
	var queryLimit int32 = 0

	client.On("Put", anyArgs...).Return(errors.New("error found")).Twice()
	err := UpdateIsiUser(ctx, client, &queryNamePrefix, &queryLimit, &queryCached,
		&queryProvider, &queryZone, &queryZone, &queryZone, &queryZone, &queryZone,
		&queryZone, &queryZone, &queryLimit, &queryLimit, &queryLimit,
		&queryResolveNames, &queryResolveNames, &queryResolveNames, &queryResolveNames)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestDeleteIsiUser(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	queryNamePrefix := "queryNamePrefix"

	var queryLimit int32 = 0

	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := DeleteIsiUser(ctx, client, &queryNamePrefix, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

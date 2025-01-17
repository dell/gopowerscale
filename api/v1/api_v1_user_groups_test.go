package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetIsiGroupList(t *testing.T) {
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
	_, err := GetIsiGroupList(ctx, client, &queryNamePrefix, &queryDomain, &queryZone, &queryProvider,
		&queryCached, &queryResolveNames, &queryMemberOf, &queryLimit)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiGroupMembers(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32 = 0

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiGroupMembers(ctx, client, &groupName, &gid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestAddIsiGroupMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32 = 0
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}
	err := AddIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	assert.Equal(t, errors.New("member type is wrong, only support user and group"), err)

	authMember.Type = fileGroupTypeUser
	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	err = AddIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestRemoveIsiGroupMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32 = 0
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "user",
	}

	authMember.Type = fileGroupTypeUser
	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := RemoveIsiGroupMember(ctx, client, &groupName, &gid, authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestCreateIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	queryForce := false
	queryZone := "queryZone"
	var gid int32 = 0
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "user",
	}
	authMemberItem := []IsiAuthMemberItem{authMember}

	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := CreateIsiGroup(ctx, client, name, &gid, authMemberItem, &queryForce, &queryZone, &queryZone)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestUpdateIsiGroupGID(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	queryZone := "queryZone"
	var gid int32 = 0
	var newgid int32 = 0

	client.On("Put", anyArgs...).Return(errors.New("error found")).Twice()
	err := UpdateIsiGroupGID(ctx, client, &groupName, &gid, newgid, &queryZone, &queryZone)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestDeleteIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32 = 0

	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := DeleteIsiGroup(ctx, client, &groupName, &gid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiGroup(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	groupName := "groupName"
	var gid int32 = 0

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiGroup(ctx, client, &groupName, &gid)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

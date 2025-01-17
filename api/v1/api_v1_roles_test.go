package v1

import (
	"context"
	"errors"
	"testing"

	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetIsiRole(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiRole(ctx, client, "")
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestGetIsiRoleList(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	x := false
	var y int32 = 5
	client.On("Get", anyArgs...).Return(errors.New("error found")).Twice()
	_, err := GetIsiRoleList(ctx, client, &x, &y)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestAddIsiRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}

	err := AddIsiRoleMember(ctx, client, "", authMember)
	assert.Equal(t, errors.New("member type is wrong, only support user and group"), err)

	authMember.Type = fileGroupTypeUser
	client.On("Post", anyArgs...).Return(errors.New("error found")).Twice()
	err = AddIsiRoleMember(ctx, client, "", authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

func TestRemoveIsiRoleMember(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	name := "name"
	var id int32 = 12
	authMember := IsiAuthMemberItem{
		ID:   &id,
		Name: &name,
		Type: "type",
	}
	client.On("Delete", anyArgs...).Return(errors.New("error found")).Twice()
	err := RemoveIsiRoleMember(ctx, client, "", authMember)
	if err == nil {
		assert.Equal(t, "Test case failed", err)
	}
}

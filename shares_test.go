/*
Copyright (c) 2023 Dell Inc, or its subsidiaries.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package goisilon

import (
	"testing"

	v12 "github.com/dell/goisilon/api/v12"
	"github.com/dell/goisilon/openapi"
	"github.com/stretchr/testify/assert"
)

func TestClient_SmbShareWithStructParams(t *testing.T) {
	// use limit to test pagination, would still output all shares
	limit := int32(1)
	_, err := client.ListALlSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assertNil(t, err)
}

func TestClient_SmbShareLifeCycleWithStructParams(t *testing.T) {
	trusteeID := "SID:S-1-1-0"
	trusteeName := "Everyone"
	trusteeType := "wellknown"
	shareName := "tf_share"
	createResponse, err := client.CreateSmbShareWithStructParams(defaultCtx, v12.CreateV12SmbShareRequest{
		V12SmbShare: &openapi.V12SmbShare{
			Name: shareName,
			Path: "/ifs/data",
			Permissions: []openapi.V1SmbSharePermission{{
				Permission:     "full",
				PermissionType: "allow",
				Trustee: openapi.V1AuthAccessAccessItemFileGroup{
					Id:   &trusteeID,
					Name: &trusteeName,
					Type: &trusteeType,
				},
			}},
		},
	})
	assertNil(t, err)
	assert.Equal(t, shareName, createResponse.Id)

	// Test list SMB
	limit := int32(1)
	getShares, err := client.ListSmbSharesWithStructParams(defaultCtx, v12.ListV12SmbSharesParams{
		Limit: &limit,
	})
	assertNil(t, err)
	assert.Equal(t, 1, len(getShares.Shares))

	// Test get SMB
	getShare, err := client.GetSmbShareWithStructParams(defaultCtx, v12.GetV12SmbShareParams{
		V12SmbShareId: shareName,
	})
	assertNil(t, err)
	assert.NotZero(t, len(getShare.Shares))
	assert.Equal(t, shareName, getShare.Shares[0].Name)

	// Test update SMB
	updateCaTimeout := int32(112)
	err = client.UpdateSmbShareWithStructParams(defaultCtx, v12.UpdateV12SmbShareRequest{
		V12SmbShareId: shareName,
		V12SmbShare: &openapi.V12SmbShareExtendedExtended{
			CaTimeout: &updateCaTimeout,
		},
	})
	assertNil(t, err)
	getShare, err = client.GetSmbShareWithStructParams(defaultCtx, v12.GetV12SmbShareParams{
		V12SmbShareId: shareName,
	})
	assertNil(t, err)
	assert.NotZero(t, len(getShare.Shares))
	assert.Equal(t, &updateCaTimeout, getShare.Shares[0].CaTimeout)

	// Test Delete SMB
	err = client.DeleteSmbShareWithStructParams(defaultCtx, v12.DeleteV12SmbShareRequest{
		V12SmbShareId: shareName,
	})
	assertNil(t, err)
	// ensure smb is cleaned
	getShare, err = client.GetSmbShareWithStructParams(defaultCtx, v12.GetV12SmbShareParams{
		V12SmbShareId: shareName,
	})
	assertNotNil(t, err)
}

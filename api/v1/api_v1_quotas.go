/*
 Copyright (c) 2022 Dell Inc, or its subsidiaries.

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
package v1

import (
	"context"
	"errors"
	"fmt"

	"github.com/dell/goisilon/api"
)

// GetIsiQuota queries the quota for a directory
func GetIsiQuota(
	ctx context.Context,
	client api.Client,
	path string) (quota *IsiQuota, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/1/quota/quotas
	// This will list out all quotas on the cluster

	var quotaResp isiQuotaListResp
	err = client.Get(ctx, quotaPath, "", nil, nil, &quotaResp)
	if err != nil {
		return nil, err
	}

	// find the specific quota we are looking for
	for _, quota := range quotaResp.Quotas {
		if quota.Path == path {
			return &quota, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Quota not found: %s", path))
}

// GetAllIsiQuota queries all quotas on the cluster
func GetAllIsiQuota(
	ctx context.Context,
	client api.Client) (quotas []*IsiQuota, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/1/quota/quotas

	var quotaResp *IsiQuotaListRespResume

	// First call without Resume param
	if err := client.Get(ctx, quotaPath, "", nil, nil, &quotaResp); err != nil {
		return nil, err
	}
	for {
		for _, q := range quotaResp.Quotas {
			quotas = append(quotas, q)
		}
		if quotaResp.Resume == "" {
			break
		}

		if quotaResp, err = GetIsiQuotaWithResume(ctx, client,
			quotaResp.Resume); err != nil {
			return nil, err
		}
	}

	return quotas, nil
}

// GetIsiQuotaWithResume queries the next page quotas based on resume token
func GetIsiQuotaWithResume(
	ctx context.Context,
	client api.Client, resume string) (quotas *IsiQuotaListRespResume, err error) {

	var quotaResp IsiQuotaListRespResume
	err = client.Get(ctx, quotaPath, "",
		api.OrderedValues{
			{[]byte("resume"), []byte(resume)},
		}, nil, &quotaResp)
	if err != nil {
		return nil, err
	}

	return &quotaResp, nil
}

// GetIsiQuotaByID get the Quota instance by ID
func GetIsiQuotaByID(
	ctx context.Context,
	client api.Client,
	ID string) (quota *IsiQuota, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/1/quota/quotas/igSJAAEAAAAAAAAAAAAAQH0RAAAAAAAA
	// This will list the quota by id on the cluster

	var quotaResp isiQuotaListResp
	err = client.Get(ctx, quotaPath, ID, nil, nil, &quotaResp)
	if err != nil {
		return nil, err
	}

	if quotaResp.Quotas != nil && len(quotaResp.Quotas) > 0 {
		quota = &quotaResp.Quotas[0]
		return quota, nil
	}

	return quota, fmt.Errorf("Quota not found: %s", ID)
}

// TODO: Add a means to set/update more than just the hard threshold

// CreateIsiQuota creates a hard directory quota on given path
func CreateIsiQuota(
	ctx context.Context,
	client api.Client,
	path string, container bool, size int64) (string, error) {

	// PAPI call: POST https://1.2.3.4:8080/platform/1/quota/quotas
	//             { "enforced" : true,
	//               "include_snapshots" : false,
	//               "path" : "/ifs/volumes/volume_name",
	//               "container" : true,
	//               "thresholds_include_overhead" : false,
	//               "type" : "directory",
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	var data = &IsiQuotaReq{
		Enforced:                  true,
		IncludeSnapshots:          false,
		Path:                      path,
		Container:                 container,
		ThresholdsIncludeOverhead: false,
		Type:                      "directory",
		Thresholds:                isiThresholdsReq{Advisory: nil, Hard: size, Soft: nil},
	}

	var quotaResp IsiQuota
	err := client.Post(ctx, quotaPath, "", nil, nil, data, &quotaResp)
	return quotaResp.Id, err
}

// SetIsiQuotaHardThreshold sets the hard threshold of a quota for a directory
// This is really just CreateIsiQuota() with container set to false
func SetIsiQuotaHardThreshold(
	ctx context.Context,
	client api.Client,
	path string, size int64) (string, error) {

	return CreateIsiQuota(ctx, client, path, false, size)
}

// UpdateIsiQuotaHardThreshold modifies the hard threshold of a quota for a directory
func UpdateIsiQuotaHardThreshold(
	ctx context.Context,
	client api.Client,
	path string, size int64) (err error) {

	// PAPI call: PUT https://1.2.3.4:8080/platform/1/quota/quotas/Id
	//             { "enforced" : true,
	//               "thresholds_include_overhead" : false,
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	var data = &IsiUpdateQuotaReq{
		Enforced:                  true,
		ThresholdsIncludeOverhead: false,
		Thresholds:                isiThresholdsReq{Advisory: nil, Hard: size, Soft: nil},
	}

	quota, err := GetIsiQuota(ctx, client, path)
	if err != nil {
		return err
	}

	var quotaResp IsiQuota
	err = client.Put(ctx, quotaPath, quota.Id, nil, nil, data, &quotaResp)
	return err
}

// UpdateIsiQuotaHardThresholdByID modifies the hard threshold of a quota for a directory
func UpdateIsiQuotaHardThresholdByID(
	ctx context.Context,
	client api.Client,
	ID string, size int64) (err error) {

	// PAPI call: PUT https://1.2.3.4:8080/platform/1/quota/quotas/Id
	//             { "enforced" : true,
	//               "thresholds_include_overhead" : false,
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	var data = &IsiUpdateQuotaReq{
		Enforced:                  true,
		ThresholdsIncludeOverhead: false,
		Thresholds:                isiThresholdsReq{Advisory: nil, Hard: size, Soft: nil},
	}

	var quotaResp IsiQuota
	err = client.Put(ctx, quotaPath, ID, nil, nil, data, &quotaResp)
	return err
}

var byteArrPath = []byte("path")
var byteArrID = []byte("id")

// DeleteIsiQuota removes the quota for a directory
func DeleteIsiQuota(
	ctx context.Context,
	client api.Client,
	path string) (err error) {

	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/quota/quotas?path=/path/to/volume
	// This will remove a the quota on a volume

	return client.Delete(
		ctx,
		quotaPath,
		"",
		api.OrderedValues{{byteArrPath, []byte(path)}},
		nil,
		nil)
}

// DeleteIsiQuotaByID removes the quota for a directory by quota id
func DeleteIsiQuotaByID(
	ctx context.Context,
	client api.Client,
	id string) (err error) {

	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/quota/quotas/AABpAQEAAAAAAAAAAAAAQA0AAAAAAAAA
	// This will remove a the quota on a volume by the quota id

	return client.Delete(
		ctx,
		quotaPath,
		id,
		nil,
		nil,
		nil)
}

// DeleteIsiQuotaByIDWithZone removes the quota for a directory by quota id with access zone
func DeleteIsiQuotaByIDWithZone(
	ctx context.Context,
	client api.Client,
	id, zone string) (err error) {

	// PAPI call: DELETE https://1.2.3.4:8080/platform/1/quota/quotas/AABpAQEAAAAAAAAAAAAAQA0AAAAAAAAA
	// This will remove a the quota on a volume by the quota id

	return client.Delete(
		ctx,
		quotaPath,
		id,
		api.OrderedValues{
			{[]byte("zone"), []byte(zone)},
		},
		nil,
		nil)
}

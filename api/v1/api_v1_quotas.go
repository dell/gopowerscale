/*
Copyright (c) 2022-2023 Dell Inc, or its subsidiaries.

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
	"fmt"

	"github.com/dell/goisilon/api"
)

// GetIsiQuota queries the quota for a directory
func GetIsiQuota(
	ctx context.Context,
	client api.Client,
	path string,
) (quota *IsiQuota, err error) {
	// PAPI call: GET https://1.2.3.4:8080/platform/1/quota/quotas?path=/path/to/volume
	// This will list the quota by path on the cluster

	var quotaResp isiQuotaListResp
	pathWithQueryParam := quotaPath + "?path=" + path
	err = client.Get(ctx, pathWithQueryParam, "", nil, nil, &quotaResp)
	if err != nil {
		return nil, err
	}

	if quotaResp.Quotas != nil && len(quotaResp.Quotas) > 0 {
		quota = &quotaResp.Quotas[0]
		return quota, nil
	}

	return nil, fmt.Errorf("Quota not found: %s", path)
}

// GetAllIsiQuota queries all quotas on the cluster
func GetAllIsiQuota(
	ctx context.Context,
	client api.Client,
) (quotas []*IsiQuota, err error) {
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
	client api.Client, resume string,
) (quotas *IsiQuotaListRespResume, err error) {
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
	ID string,
) (quota *IsiQuota, err error) {
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
	path string, container bool, size, softLimit, advisoryLimit, softGracePrd int64,
	includeSnapshots bool,
) (string, error) {
	// PAPI call: POST https://1.2.3.4:8080/platform/1/quota/quotas
	//             { "enforced" : true,
	//               "include_snapshots" : true,
	//               "path" : "/ifs/volumes/volume_name",
	//               "container" : true,
	//               "thresholds_include_overhead" : false,
	//               "type" : "directory",
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	// body={'path': '/ifs/data/quotatest', 'thresholds': {'soft_grace': 86400L, 'soft': 1048576L}, 'include_snapshots': True, 'force': False, 'type': 'directory'}
	// softGrace := 86400U
	thresholds := isiThresholdsReq{Advisory: advisoryLimit, Hard: size, Soft: softLimit, SoftGrace: softGracePrd}
	if advisoryLimit == 0 {
		thresholds.Advisory = nil
	}
	if softLimit == 0 {
		thresholds.Soft = nil
	}
	if softGracePrd == 0 {
		thresholds.SoftGrace = nil
	}
	data := &IsiQuotaReq{
		Enforced:                  true,
		IncludeSnapshots:          includeSnapshots,
		Path:                      path,
		Container:                 container,
		ThresholdsIncludeOverhead: false,
		Type:                      "directory",
		Thresholds:                thresholds,
	}

	var quotaResp IsiQuota
	err := client.Post(ctx, quotaPath, "", nil, nil, data, &quotaResp)
	return quotaResp.ID, err
}

// SetIsiQuotaHardThreshold sets the hard threshold of a quota for a directory
// This is really just CreateIsiQuota() with container set to false
func SetIsiQuotaHardThreshold(
	ctx context.Context,
	client api.Client,
	path string, size, softLimit, advisoryLimit, softGracePrd int64,
) (string, error) {
	return CreateIsiQuota(ctx, client, path, false, size, softLimit, advisoryLimit, softGracePrd, true)
}

// UpdateIsiQuotaHardThreshold modifies the hard threshold of a quota for a directory
func UpdateIsiQuotaHardThreshold(
	ctx context.Context,
	client api.Client,
	path string, size, softLimit, advisoryLimit, softGracePrd int64,
) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/quota/quotas/Id
	//             { "enforced" : true,
	//               "thresholds_include_overhead" : false,
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	thresholds := isiThresholdsReq{Advisory: advisoryLimit, Hard: size, Soft: softLimit, SoftGrace: softGracePrd}
	if advisoryLimit == 0 {
		thresholds.Advisory = nil
	}
	if softLimit == 0 {
		thresholds.Soft = nil
	}
	if softGracePrd == 0 {
		thresholds.SoftGrace = nil
	}

	data := &IsiUpdateQuotaReq{
		Enforced:                  true,
		ThresholdsIncludeOverhead: false,
		Thresholds:                thresholds,
	}

	quota, err := GetIsiQuota(ctx, client, path)
	if err != nil {
		return err
	}

	var quotaResp IsiQuota
	err = client.Put(ctx, quotaPath, quota.ID, nil, nil, data, &quotaResp)
	return err
}

// UpdateIsiQuotaHardThresholdByID modifies the hard threshold of a quota for a directory
func UpdateIsiQuotaHardThresholdByID(
	ctx context.Context,
	client api.Client,
	ID string, size, softLimit, advisoryLimit, softGracePrd int64,
) (err error) {
	// PAPI call: PUT https://1.2.3.4:8080/platform/1/quota/quotas/Id
	//             { "enforced" : true,
	//               "thresholds_include_overhead" : false,
	//               "thresholds" : { "advisory" : null,
	//                                "hard" : 1234567890,
	//                                "soft" : null
	//                              }
	//             }
	thresholds := isiThresholdsReq{Advisory: advisoryLimit, Hard: size, Soft: softLimit, SoftGrace: softGracePrd}
	if advisoryLimit == 0 {
		thresholds.Advisory = nil
	}
	if softLimit == 0 {
		thresholds.Soft = nil
	}
	if softGracePrd == 0 {
		thresholds.SoftGrace = nil
	}
	data := &IsiUpdateQuotaReq{
		Enforced:                  true,
		ThresholdsIncludeOverhead: false,
		Thresholds:                thresholds,
	}

	var quotaResp IsiQuota
	err = client.Put(ctx, quotaPath, ID, nil, nil, data, &quotaResp)
	return err
}

var (
	byteArrPath = []byte("path")
	byteArrID   = []byte("id")
)

// DeleteIsiQuota removes the quota for a directory
func DeleteIsiQuota(
	ctx context.Context,
	client api.Client,
	path string,
) (err error) {
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
	id string,
) (err error) {
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
	id, zone string,
) (err error) {
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

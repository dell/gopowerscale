/* 
 Copyright (c) 2019 Dell Inc, or its subsidiaries.

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
package v5

import (
	"context"
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/dell/goisilon/api"
)

const (
	quotaLicensePath = "platform/5/quota/license"
)

// SmartQuotas license statuses
var unlicensed = QuotaLicenseStatus{"Unlicensed"}
var licensed = QuotaLicenseStatus{"Licensed"}
var expired = QuotaLicenseStatus{"Expired"}
var evaluation = QuotaLicenseStatus{"Evaluation"}
var evaluationExpired = QuotaLicenseStatus{"Evaluation Expired"}

var validQuotaLicenseStatus = [5]QuotaLicenseStatus{unlicensed, licensed, expired, evaluation, evaluationExpired}

// QuotaLicense contains the SmartQuotas license info.
type QuotaLicense struct {
	DaysToExpiry int    `json:"days_to_expiry,omitempty"`
	Expiration   string `json:"expiration,omitempty"`
	ID           string `json:"id,omitempty"`
	NAME         string `json:"name,omitempty"`
	STATUS       string `json:"status,omitempty"`
}

// GetIsiQuotaLicense retrieves the SmartQuotas license info
func GetIsiQuotaLicense(
	ctx context.Context,
	client api.Client) (lic *QuotaLicense, err error) {

	// PAPI call: GET https://1.2.3.4:8080/platform/5/quota/license
	// This will return the SmartQuotas license info

	var quotaLicense QuotaLicense
	err = client.Get(ctx, quotaLicensePath, "", nil, nil, &quotaLicense)
	if err != nil {
		return nil, err
	}

	return &quotaLicense, nil
}

func getIsiQuotaLicenseStatus(
	ctx context.Context,
	client api.Client) (status string, err error) {

	lic, e := GetIsiQuotaLicense(ctx, client)

	if e != nil {
		return "", e
	}

	if lic.STATUS == "" {
		return "", errors.New("SmartQuotas license status is empty")
	}

	log.Debugf("SmartQuotas license status retrieved : '%s'", lic.STATUS)

	if !isQuotaLicenseStatusValid(lic.STATUS) {
		return "", fmt.Errorf("unknown SmartQuotas license status '%s'", lic.STATUS)
	}

	return lic.STATUS, nil
}

// IsQuotaLicenseActivated checks if SmartQuotas has been activated (either licensed or in evaluation)
func IsQuotaLicenseActivated(ctx context.Context,
	client api.Client) (bool, error) {

	status, err := getIsiQuotaLicenseStatus(ctx, client)

	if err != nil {

		log.Errorf("error encountered when retrieving SmartQuotas license info, cannot determine whether SmartQuotas is activated. error : '%v'", err)
		return false, nil
	}

	isActivated := status == licensed.toString() || status == evaluation.toString()

	return isActivated, nil
}

func isQuotaLicenseStatusValid(status string) bool {

	isStatusValid := false

	for _, stat := range validQuotaLicenseStatus {
		if stat.toString() == status {
			isStatusValid = true
		}
	}

	return isStatusValid
}

// QuotaLicenseStatus represents a SmartQuotas license status
type QuotaLicenseStatus struct {
	value string
}

func (s *QuotaLicenseStatus) toString() string {

	return s.value
}

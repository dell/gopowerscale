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
package v3

type isiStats struct {
	ID        int    `json:"devid"`
	Error     string `json:"error"`
	ErrorCode int    `json:"error_code"`
	Key       string `json:"key"`
	Time      int64  `json:"time"`
	Value     int64  `json:"value"`
}

type isiFloatStats struct {
	ID        int     `json:"devid"`
	Error     string  `json:"error"`
	ErrorCode int     `json:"error_code"`
	Key       string  `json:"key"`
	Time      int64   `json:"time"`
	Value     float64 `json:"value"`
}

// IsiStatsResp PAPI stats response attributes JSON structure
type IsiStatsResp struct {
	StatsList []*isiStats `json:"stats"`
}

// IsiFloatStatsResp PAPI stats response float attributes JSON structure
type IsiFloatStatsResp struct {
	StatsList []*isiFloatStats `json:"stats"`
}

// IsiClusterConfig returns the configuration information of cluster.
// TODO add other variables refers to the request cluster/config
type IsiClusterConfig struct {
	Description string       `json:"description"`
	Devices     []*IsiDevice `json:"devices"`
	GUID        string       `json:"guid"`
	JoinMode    string       `json:"join_mode"`
	LocalDevId  int64        `json:"local_devid"`
	LocalLnn    int64        `json:"local_lnn"`
	LocalSerial string       `json:"local_serial"`
	Name        string       `json:"name"`
}

// IsiDevice refers to device information of a cluster
type IsiDevice struct {
	DevId int64  `json:"devid"`
	GUID  string `json:"guid"`
	IsUp  bool   `json:"is_up"`
	Lnn   int64  `json:"lnn"`
}

type isiClientList struct {
	Protocol   string `json:"protocol"`
	RemoteAddr string `json:"remote_addr"`
	RemoteName string `json:"remote_name"`
}
type ExportClientList struct {
	ClientsList []*isiClientList `json:"client"`
}

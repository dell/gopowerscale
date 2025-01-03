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

package v7

// IsiClusterInternalNetworksFailoverIPAddresse  Specifies range of IP addresses where 'low' is starting address and 'high' is the end address.' Both 'low' and 'high' addresses are inclusive to the range.
type IsiClusterInternalNetworksFailoverIPAddresse struct {
	// IPv4 address in the format: xxx.xxx.xxx.xxx
	High string `json:"high"`
	// IPv4 address in the format: xxx.xxx.xxx.xxx
	Low string `json:"low"`
}

// IsiClusterInternalNetworks Configuration fields for internal networks.
type IsiClusterInternalNetworks struct {
	// Array of IP address ranges to be used to configure the internal failover network of the OneFS cluster.
	FailoverIPAddresses []IsiClusterInternalNetworksFailoverIPAddresse `json:"failover_ip_addresses,omitempty"`
	// Status of failover network.
	FailoverStatus *string `json:"failover_status,omitempty"`
	// Network fabric used for the primary network int-a.
	IntAFabric *string `json:"int_a_fabric,omitempty"`
	// Array of IP address ranges to be used to configure the internal int-a network of the OneFS cluster.
	IntAIpAddresses []IsiClusterInternalNetworksFailoverIPAddresse `json:"int_a_ip_addresses,omitempty"`
	// Maximum Transfer Unit (MTU) of the primary network int-a.
	IntAMtu *int32 `json:"int_a_mtu,omitempty"`
	// Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.
	IntAPrefixLength *int32 `json:"int_a_prefix_length,omitempty"`
	// Status of the primary network int-a.
	IntAStatus *string `json:"int_a_status,omitempty"`
	// Network fabric used for the failover network.
	IntBFabric *string `json:"int_b_fabric,omitempty"`
	// Array of IP address ranges to be used to configure the internal int-b network of the OneFS cluster.
	IntBIpAddresses []IsiClusterInternalNetworksFailoverIPAddresse `json:"int_b_ip_addresses,omitempty"`
	// Maximum Transfer Unit (MTU) of the failover network int-b.
	IntBMtu *int32 `json:"int_b_mtu,omitempty"`
	// Prefixlen specifies the length of network bits used in an IP address. This field is the right-hand part of the CIDR notation representing the subnet mask.
	IntBPrefixLength *int32 `json:"int_b_prefix_length,omitempty"`
}

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

// IsiClusterConfigOnefsVersion struct for IsiClusterConfigOnefsVersion
type IsiClusterConfigOnefsVersion struct {
	// OneFS build string.
	Build string `json:"build"`
	// Kernel release number.
	Release string `json:"release"`
	// OneFS build number.
	Revision string `json:"revision"`
	// Kernel release type.
	Type string `json:"type"`
	// Kernel full version information.
	Version string `json:"version"`
}

// IsiClusterConfigTimezone The cluster timezone settings.
type IsiClusterConfigTimezone struct {
	// Timezone abbreviation.
	Abbreviation string `json:"abbreviation,omitempty"`
	// Customer timezone information.
	Custom string `json:"custom,omitempty"`
	// Timezone full name.
	Name string `json:"name,omitempty"`
	// Timezone hierarchical name.
	Path string `json:"path,omitempty"`
}

// IsiClusterConfig returns the configuration information of cluster.
// TODO add other variables refers to the request cluster/config
type IsiClusterConfig struct {
	Description  string                        `json:"description"`
	Devices      []*IsiDevice                  `json:"devices"`
	GUID         string                        `json:"guid"`
	JoinMode     string                        `json:"join_mode"`
	LocalDevID   int64                         `json:"local_devid"`
	LocalLnn     int64                         `json:"local_lnn"`
	LocalSerial  string                        `json:"local_serial"`
	Name         string                        `json:"name"`
	OnefsVersion *IsiClusterConfigOnefsVersion `json:"onefs_version,omitempty"`
	Timezone     *IsiClusterConfigTimezone     `json:"timezone,omitempty"`
}

// IsiDevice refers to device information of a cluster
type IsiDevice struct {
	DevID int64  `json:"devid"`
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

// IsiClusterIdentityLogon The information displayed when a user logs in to the cluster.
type IsiClusterIdentityLogon struct {
	// The message of the day.
	Motd string `json:"motd"`
	// The header to the message of the day.
	MotdHeader string `json:"motd_header"`
}

// IsiClusterIdentity Unprivileged cluster information for display when logging in.
type IsiClusterIdentity struct {
	// A description of the cluster.
	Description string                  `json:"description"`
	Logon       IsiClusterIdentityLogon `json:"logon"`
	// The name of the cluster.
	Name string `json:"name"`
}

// IsiClusterNodes struct for IsiClusterNodes
type IsiClusterNodes struct {
	// A list of errors encountered by the individual nodes involved in this request, or an empty list if there were no errors.
	Errors []IsiNodeStatusError `json:"errors,omitempty"`
	// The responses from the individual nodes involved in this request.
	Nodes []IsiClusterNode `json:"nodes,omitempty"`
	// The total number of nodes responding.
	Total *int32 `json:"total,omitempty"`
}

// IsiNodeStatusError An object describing a single error.
type IsiNodeStatusError struct {
	// The error code.
	Code *string `json:"code,omitempty"`
	// The field with the error if applicable.
	Field *string `json:"field,omitempty"`
	// Node ID (Device Number) of a node.
	ID *int32 `json:"id,omitempty"`
	// Logical Node Number (LNN) of a node.
	Lnn *int32 `json:"lnn,omitempty"`
	// The error message.
	Message *string `json:"message,omitempty"`
	// HTTP Status code returned by this node.
	Status *int32 `json:"status,omitempty"`
}

// IsiClusterNode Node information.
type IsiClusterNode struct {
	// List of the drives in this node.
	Drives []IsiClusterNodeDrive `json:"drives,omitempty"`
	// Error message, if the HTTP status returned from this node was not 200.
	Error    *string                 `json:"error,omitempty"`
	Hardware *IsiClusterNodeHardware `json:"hardware,omitempty"`
	// Node ID (Device Number) of a node.
	ID *int32 `json:"id,omitempty"`
	// Logical Node Number (LNN) of a node.
	Lnn        *int32                    `json:"lnn,omitempty"`
	Partitions *V10ClusterNodePartitions `json:"partitions,omitempty"`
	Sensors    *IsiClusterNodeSensors    `json:"sensors,omitempty"`
	State      *IsiClusterNodeState      `json:"state,omitempty"`
	Status     *IsiClusterNodeStatus     `json:"status,omitempty"`
}

// IsiClusterNodeDrive Drive information.
type IsiClusterNodeDrive struct {
	// Numerical representation of this drive's bay.
	Baynum *int32 `json:"baynum,omitempty"`
	// Number of blocks on this drive.
	Blocks *int64 `json:"blocks,omitempty"`
	// The chassis number which contains this drive.
	Chassis *int32 `json:"chassis,omitempty"`
	// This drive's device name.
	Devname  *string                      `json:"devname,omitempty"`
	Firmware *IsiClusterNodeDriveFirmware `json:"firmware,omitempty"`
	// Drive_d's handle representation for this driveIf we fail to retrieve the handle for this drive from drive_d: -1
	Handle *int32 `json:"handle,omitempty"`
	// String representation of this drive's interface type.
	InterfaceType *string `json:"interface_type,omitempty"`
	// This drive's logical drive number in IFS.
	Lnum *int32 `json:"lnum,omitempty"`
	// String representation of this drive's physical location.
	Locnstr *string `json:"locnstr,omitempty"`
	// Size of a logical block on this drive.
	LogicalBlockLength *int32 `json:"logical_block_length,omitempty"`
	// String representation of this drive's media type.
	MediaType *string `json:"media_type,omitempty"`
	// This drive's manufacturer and model.
	Model *string `json:"model,omitempty"`
	// Size of a physical block on this drive.
	PhysicalBlockLength *int32 `json:"physical_block_length,omitempty"`
	// Indicates whether this drive is physically present in the node.
	Present *bool `json:"present,omitempty"`
	// This drive's purpose in the DRV state machine.
	Purpose *string `json:"purpose,omitempty"`
	// Description of this drive's purpose.
	PurposeDescription *string `json:"purpose_description,omitempty"`
	// Serial number for this drive.
	Serial *string `json:"serial,omitempty"`
	// This drive's state as presented to the UI.
	UIState *string `json:"ui_state,omitempty"`
	// The drive's 'worldwide name' from its NAA identifiers.
	Wwn *string `json:"wwn,omitempty"`
	// This drive's x-axis grid location.
	XLoc *int32 `json:"x_loc,omitempty"`
	// This drive's y-axis grid location.
	YLoc *int32 `json:"y_loc,omitempty"`
}

// IsiClusterNodeDriveFirmware Drive firmware information.
type IsiClusterNodeDriveFirmware struct {
	// This drive's current firmware revision
	CurrentFirmware *string `json:"current_firmware,omitempty"`
	// This drive's desired firmware revision.
	DesiredFirmware *string `json:"desired_firmware,omitempty"`
}

// IsiClusterNodeHardware Node hardware identifying information (static).
type IsiClusterNodeHardware struct {
	// Name of this node's chassis.
	Chassis *string `json:"chassis,omitempty"`
	// Chassis code of this node (1U, 2U, etc.).
	ChassisCode *string `json:"chassis_code,omitempty"`
	// Number of chassis making up this node.
	ChassisCount *string `json:"chassis_count,omitempty"`
	// Class of this node (storage, accelerator, etc.).
	Class *string `json:"class,omitempty"`
	// Node configuration ID.
	ConfigurationID *string `json:"configuration_id,omitempty"`
	// Manufacturer and model of this node's CPU.
	CPU *string `json:"cpu,omitempty"`
	// Manufacturer and model of this node's disk controller.
	DiskController *string `json:"disk_controller,omitempty"`
	// Manufacturer and model of this node's disk expander.
	DiskExpander *string `json:"disk_expander,omitempty"`
	// Family code of this node (X, S, NL, etc.).
	FamilyCode *string `json:"family_code,omitempty"`
	// Manufacturer, model, and device id of this node's flash drive.
	FlashDrive *string `json:"flash_drive,omitempty"`
	// Generation code of this node.
	GenerationCode *string `json:"generation_code,omitempty"`
	// PowerScale hardware generation name.
	Hwgen *string `json:"hwgen,omitempty"`
	// Version of this node's PowerScale Management Board.
	ImbVersion *string `json:"imb_version,omitempty"`
	// Infiniband card type.
	Infiniband *string `json:"infiniband,omitempty"`
	// Version of the LCD panel.
	LcdVersion *string `json:"lcd_version,omitempty"`
	// Manufacturer and model of this node's motherboard.
	Motherboard *string `json:"motherboard,omitempty"`
	// Description of all this node's network interfaces.
	NetInterfaces *string `json:"net_interfaces,omitempty"`
	// Manufacturer and model of this node's NVRAM board.
	Nvram *string `json:"nvram,omitempty"`
	// Description strings for each power supply on this node.
	Powersupplies []string `json:"powersupplies,omitempty"`
	// Number of processors and cores on this node.
	Processor *string `json:"processor,omitempty"`
	// PowerScale product name.
	Product *string `json:"product,omitempty"`
	// Size of RAM in bytes.
	RAM *int64 `json:"ram,omitempty"`
	// Serial number of this node.
	SerialNumber *string `json:"serial_number,omitempty"`
	// Series of this node (X, I, NL, etc.).
	Series *string `json:"series,omitempty"`
	// Storage class of this node (storage or diskless).
	StorageClass *string `json:"storage_class,omitempty"`
}

// IsiClusterNodeState Node state information (reported and modifiable).
type IsiClusterNodeState struct {
	Readonly     map[string]interface{}           `json:"readonly,omitempty"`
	Servicelight *IsiClusterNodeStateServicelight `json:"servicelight,omitempty"`
	Smartfail    *IsiClusterNodeStateSmartfail    `json:"smartfail,omitempty"`
}

// IsiClusterNodeStateServicelight Node service light state.
type IsiClusterNodeStateServicelight struct {
	// The node service light state (True = on).
	Enabled bool `json:"enabled"`
}

// IsiClusterNodeStateSmartfail Node smartfail state.
type IsiClusterNodeStateSmartfail struct {
	// This node is smartfailed (soft_devs).
	Smartfailed *bool `json:"smartfailed,omitempty"`
}

// V10ClusterNodePartitions Node partition information.
type V10ClusterNodePartitions struct {
	// Count of how many partitions are included.
	Count *int32 `json:"count,omitempty"`
	// Partition information.
	Partitions []IsiClusterNodePartition `json:"partitions,omitempty"`
}

// IsiClusterNodePartition Node partition information.
type IsiClusterNodePartition struct {
	// The block size used for the reported partition information.
	BlockSize *int32 `json:"block_size,omitempty"`
	// Total blocks on this file system partition.
	Capacity *int32 `json:"capacity,omitempty"`
	// Comma separated list of devices used for this file system partition.
	ComponentDevices *string `json:"component_devices,omitempty"`
	// Directory on which this partition is mounted.
	MountPoint *string `json:"mount_point,omitempty"`
	// Used blocks on this file system partition, expressed as a percentage.
	PercentUsed *string                        `json:"percent_used,omitempty"`
	Statfs      *IsiClusterNodePartitionStatfs `json:"statfs,omitempty"`
	// Used blocks on this file system partition.
	Used *int32 `json:"used,omitempty"`
}

// IsiClusterNodePartitionStatfs System partition details as provided by statfs(2).
type IsiClusterNodePartitionStatfs struct {
	// Free blocks available to non-superuser on this partition.
	FBavail *int32 `json:"f_bavail,omitempty"`
	// Free blocks on this partition.
	FBfree *int32 `json:"f_bfree,omitempty"`
	// Total data blocks on this partition.
	FBlocks *int32 `json:"f_blocks,omitempty"`
	// Filesystem fragment size; block size in OneFS.
	FBsize *int32 `json:"f_bsize,omitempty"`
	// Free file nodes avail to non-superuser.
	FFfree *int32 `json:"f_ffree,omitempty"`
	// Total file nodes in filesystem.
	FFiles *int32 `json:"f_files,omitempty"`
	// Mount exported flags.
	FFlags *int32 `json:"f_flags,omitempty"`
	// File system type name.
	FFstypename *string `json:"f_fstypename,omitempty"`
	// Optimal transfer block size.
	FIosize *int32 `json:"f_iosize,omitempty"`
	// Names of devices this partition is mounted from.
	FMntfromname *string `json:"f_mntfromname,omitempty"`
	// Directory this partition is mounted to.
	FMntonname *string `json:"f_mntonname,omitempty"`
	// Maximum filename length.
	FNamemax *int32 `json:"f_namemax,omitempty"`
	// UID of user that mounted the filesystem.
	FOwner *int32 `json:"f_owner,omitempty"`
	// Type of filesystem.
	FType *int32 `json:"f_type,omitempty"`
	// statfs() structure version number.
	FVersion *int32 `json:"f_version,omitempty"`
}

// IsiClusterNodeSensors Node sensor information (hardware reported).
type IsiClusterNodeSensors struct {
	// This node's sensor information.
	Sensors []IsiClusterNodeSensor `json:"sensors,omitempty"`
}

// IsiClusterNodeSensor Node sensor information.
type IsiClusterNodeSensor struct {
	// The count of values in this sensor group.
	Count *int32 `json:"count,omitempty"`
	// The name of this sensor group.
	Name *string `json:"name,omitempty"`
	// The list of specific sensor value info in this sensor group.
	Values []IsiClusterNodeSensorValue `json:"values,omitempty"`
}

// IsiClusterNodeSensorValue Specific sensor value info.
type IsiClusterNodeSensorValue struct {
	// The descriptive name of this sensor.
	Desc *string `json:"desc,omitempty"`
	// The identifier name of this sensor.
	Name *string `json:"name,omitempty"`
	// The units of this sensor.
	Units *string `json:"units,omitempty"`
	// The value of this sensor.
	Value *string `json:"value,omitempty"`
}

// IsiClusterNodeStatus Node status information (hardware reported).
type IsiClusterNodeStatus struct {
	Batterystatus *IsiClusterNodeStatusBatterystatus `json:"batterystatus,omitempty"`
	// Storage capacity of this node.
	Capacity      []IsiClusterNodeStatusCapacityItem `json:"capacity,omitempty"`
	CPU           *IsiClusterNodeStatusCPU           `json:"cpu,omitempty"`
	Nvram         *IsiClusterNodeStatusNvram         `json:"nvram,omitempty"`
	Powersupplies *IsiClusterNodeStatusPowersupplies `json:"powersupplies,omitempty"`
	// OneFS release.
	Release *string `json:"release,omitempty"`
	// Seconds this node has been online.
	Uptime *int32 `json:"uptime,omitempty"`
	// OneFS version.
	Version *string `json:"version,omitempty"`
}

// IsiClusterNodeStatusBatterystatus Battery status information.
type IsiClusterNodeStatusBatterystatus struct {
	// The last battery test time for battery 1.
	LastTestTime1 *string `json:"last_test_time1,omitempty"`
	// The last battery test time for battery 2.
	LastTestTime2 *string `json:"last_test_time2,omitempty"`
	// The next checkup for battery 1.
	NextTestTime1 *string `json:"next_test_time1,omitempty"`
	// The next checkup for battery 2.
	NextTestTime2 *string `json:"next_test_time2,omitempty"`
	// Node has battery status.
	Present *bool `json:"present,omitempty"`
	// The result of the last battery test for battery 1.
	Result1 *string `json:"result1,omitempty"`
	// The result of the last battery test for battery 2.
	Result2 *string `json:"result2,omitempty"`
	// The status of battery 1.
	Status1 *string `json:"status1,omitempty"`
	// The status of battery 2.
	Status2 *string `json:"status2,omitempty"`
	// Node supports battery status.
	Supported *bool `json:"supported,omitempty"`
}

// IsiClusterNodeStatusCapacityItem Node capacity information.
type IsiClusterNodeStatusCapacityItem struct {
	// Total device storage bytes.
	Bytes *int64 `json:"bytes,omitempty"`
	// Total device count.
	Count *int32 `json:"count,omitempty"`
	// Device type.
	Type *string `json:"type,omitempty"`
}

// IsiClusterNodeStatusCpu CPU status information for this node.
type IsiClusterNodeStatusCPU struct {
	// Manufacturer model description of this CPU.
	Model *string `json:"model,omitempty"`
	// CPU overtemp state.
	Overtemp *string `json:"overtemp,omitempty"`
	// Type of processor and core of this CPU.
	Proc *string `json:"proc,omitempty"`
	// CPU throttling (expressed as a percentage).
	SpeedLimit *string `json:"speed_limit,omitempty"`
}

// IsiClusterNodeStatusNvram Node NVRAM information.
type IsiClusterNodeStatusNvram struct {
	// This node's NVRAM battery status information.
	Batteries []IsiClusterNodeStatusNvramBattery `json:"batteries,omitempty"`
	// This node's NVRAM battery count. On failure: -1, otherwise 1 or 2.
	BatteryCount *int32 `json:"battery_count,omitempty"`
	// This node's NVRAM battery charge status, as a color.
	ChargeStatus *string `json:"charge_status,omitempty"`
	// This node's NVRAM battery charge status, as a number. Error or not supported: -1. BR_BLACK: 0. BR_GREEN: 1. BR_YELLOW: 2. BR_RED: 3.
	ChargeStatusNumber *int32 `json:"charge_status_number,omitempty"`
	// This node's NVRAM device name with path.
	Device *string `json:"device,omitempty"`
	// This node has NVRAM.
	Present *bool `json:"present,omitempty"`
	// This node has NVRAM with flash storage.
	PresentFlash *bool `json:"present_flash,omitempty"`
	// The size of the NVRAM, in bytes.
	PresentSize *int32 `json:"present_size,omitempty"`
	// This node's NVRAM type.
	PresentType *string `json:"present_type,omitempty"`
	// This node's current ship mode state for NVRAM batteries. If not supported or on failure: -1. Disabled: 0. Enabled: 1.
	ShipMode *int32 `json:"ship_mode,omitempty"`
	// This node supports NVRAM.
	Supported *bool `json:"supported,omitempty"`
	// This node supports NVRAM with flash storage.
	SupportedFlash *bool `json:"supported_flash,omitempty"`
	// The maximum size of the NVRAM, in bytes.
	SupportedSize *int64 `json:"supported_size,omitempty"`
	// This node's supported NVRAM type.
	SupportedType *string `json:"supported_type,omitempty"`
}

// IsiClusterNodeStatusNvramBattery NVRAM battery status information.
type IsiClusterNodeStatusNvramBattery struct {
	// The current status color of the NVRAM battery.
	Color *string `json:"color,omitempty"`
	// Identifying index for the NVRAM battery.
	ID *int32 `json:"id,omitempty"`
	// The current status message of the NVRAM battery.
	Status *string `json:"status,omitempty"`
	// The current voltage of the NVRAM battery.
	Voltage *string `json:"voltage,omitempty"`
}

// IsiClusterNodeStatusPowersupplies Information about this node's power supplies.
type IsiClusterNodeStatusPowersupplies struct {
	// Count of how many power supplies are supported.
	Count *int32 `json:"count,omitempty"`
	// Count of how many power supplies have failed.
	Failures *int32 `json:"failures,omitempty"`
	// Does this node have a CFF power supply.
	HasCff *bool `json:"has_cff,omitempty"`
	// A descriptive status string for this node's power supplies.
	Status *string `json:"status,omitempty"`
	// List of this node's power supplies.
	Supplies []IsiClusterNodeStatusPowersuppliesSupply `json:"supplies,omitempty"`
	// Does this node support CFF power supplies.
	SupportsCff *bool `json:"supports_cff,omitempty"`
}

// IsiClusterNodeStatusPowersuppliesSupply Power supply information.
type IsiClusterNodeStatusPowersuppliesSupply struct {
	// Which node chassis is this power supply in.
	Chassis *int32 `json:"chassis,omitempty"`
	// The current firmware revision of this power supply.
	Firmware *string `json:"firmware,omitempty"`
	// Is this power supply in a failure state.
	Good *string `json:"good,omitempty"`
	// Identifying index for this power supply.
	ID int32 `json:"id"`
	// Complete identifying string for this power supply.
	Name *string `json:"name,omitempty"`
	// A descriptive status string for this power supply.
	Status *string `json:"status,omitempty"`
	// The type of this power supply.
	Type *string `json:"type,omitempty"`
}

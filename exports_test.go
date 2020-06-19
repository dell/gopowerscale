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
package goisilon

import (
	"context"
	"fmt"
	"sort"
	"testing"

	log "github.com/akutz/gournal"
	"github.com/dell/goisilon/api"
	"github.com/stretchr/testify/assert"
)

var exportForClient int

func TestExportsList(t *testing.T) {
	volumeName1 := "test_get_exports1"
	volumeName2 := "test_get_exports2"
	volumeName3 := "test_get_exports3"

	// Identify all exports currently on the cluster
	exportMap := make(map[int]string)
	exports, err := client.GetExports(context.Background())
	assertNoError(t, err)

	for _, export := range exports {
		exportMap[export.ID] = (*export.Paths)[0]
	}
	initialExportCount := len(exports)

	var (
		vol      Volume
		exportID int
	)

	// Add the test exports
	vol, err = client.CreateVolume(defaultCtx, volumeName1)
	assertNoError(t, err)
	assertNotNil(t, vol)
	volumeName1 = vol.Name
	volumePath1 := client.API.VolumePath(volumeName1)
	t.Logf("created volume: %s", volumeName1)

	vol, err = client.CreateVolume(defaultCtx, volumeName2)
	assertNoError(t, err)
	assertNotNil(t, vol)
	volumeName2 = vol.Name
	volumePath2 := client.API.VolumePath(volumeName2)
	t.Logf("created volume: %s", volumeName2)

	vol, err = client.CreateVolume(defaultCtx, volumeName3)
	assertNoError(t, err)
	assertNotNil(t, vol)
	volumeName3 = vol.Name
	volumePath3 := client.API.VolumePath(volumeName3)
	t.Logf("created volume: %s", volumeName3)

	exportID, err = client.Export(defaultCtx, volumeName1)
	assertNoError(t, err)
	t.Logf("created export: %d", exportID)

	exportID, err = client.Export(defaultCtx, volumeName2)
	assertNoError(t, err)
	t.Logf("created export: %d", exportID)

	exportID, err = client.Export(defaultCtx, volumeName3)
	assertNoError(t, err)
	t.Logf("created export: %d", exportID)

	// make sure we clean up when we're done
	defer client.Unexport(defaultCtx, volumeName1)
	defer client.Unexport(defaultCtx, volumeName2)
	defer client.Unexport(defaultCtx, volumeName3)
	defer client.DeleteVolume(defaultCtx, volumeName1)
	defer client.DeleteVolume(defaultCtx, volumeName2)
	defer client.DeleteVolume(defaultCtx, volumeName3)

	// Get the updated export list
	exports, err = client.GetExports(defaultCtx)
	assertNoError(t, err)

	// Verify that the new exports are there as well as all the old exports.
	if !assert.Equal(t, initialExportCount+3, len(exports)) {
		t.FailNow()
	}

	// Remove the original exports and add the new ones.  In the end, we should only have the
	// exports we just created and nothing more.
	for _, export := range exports {
		if _, found := exportMap[export.ID]; found == true {
			// this export was exported prior to the test start
			delete(exportMap, export.ID)
		} else {
			// this export is new
			exportMap[export.ID] = (*export.Paths)[0]
		}
	}

	if !assert.Len(t, exportMap, 3) {
		t.FailNow()
	}

	volumeBitmap := 0
	for _, path := range exportMap {
		if path == volumePath1 {
			volumeBitmap++
		} else if path == volumePath2 {
			volumeBitmap += 2
		} else if path == volumePath3 {
			volumeBitmap += 4
		}
	}

	assert.Equal(t, 7, volumeBitmap)
}

func TestExportCreate(t *testing.T) {
	volumeName := "test_create_export"
	volumePath := client.API.VolumePath(volumeName)

	// setup the test
	_, err := client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.Unexport(defaultCtx, volumeName)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify the volume isn't already exported
	export, err := client.GetExportByName(defaultCtx, volumeName)
	assertNoError(t, err)
	assertNil(t, export)

	desc := "description for test"

	// export the volume
	_, err = client.ExportWithZone(defaultCtx, volumeName, "System", desc)
	assertNoError(t, err)

	// verify the volume has been exported
	export, err = client.GetExportByName(defaultCtx, volumeName)
	assertNoError(t, err)
	assertNotNil(t, export)
	assert.Equal(t, desc, export.Description, "unexpected description of the export")

	found := false
	for _, path := range *export.Paths {
		if path == volumePath {
			found = true
			break
		}
	}
	assert.True(t, found)
}

func TestExportDelete(t *testing.T) {
	volumeName := "test_unexport_volume"

	// initialize the export
	_, err := client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	_, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify the volume is exported
	export, err := client.GetExportByName(defaultCtx, volumeName)
	assertNoError(t, err)
	assertNotNil(t, export)

	// Unexport the volume
	err = client.Unexport(defaultCtx, volumeName)
	assertNoError(t, err)

	// verify the volume is no longer exported
	export, err = client.GetExportByName(defaultCtx, volumeName)
	assertNoError(t, err)
	assertNil(t, export)
}

func TestExportNonRootMapping(t *testing.T) {
	testUserMapping(
		t,
		"test_export_non_root_mapping",
		client.GetNonRootMappingByID,
		client.EnableNonRootMappingByID,
		client.DisableNonRootMappingByID)
}

func TestExportFailureMapping(t *testing.T) {
	testUserMapping(
		t,
		"test_export_failure_mapping",
		client.GetFailureMappingByID,
		client.EnableFailureMappingByID,
		client.DisableFailureMappingByID)
}

func TestExportRootMapping(t *testing.T) {
	testUserMapping(
		t,
		"test_export_root_mapping",
		client.GetRootMappingByID,
		client.EnableRootMappingByID,
		client.DisableRootMappingByID)
}

func testUserMapping(
	t *testing.T,
	volumeName string,
	getMap func(ctx context.Context, id int) (UserMapping, error),
	enaMap func(ctx context.Context, id int, user string) error,
	disMap func(ctx context.Context, id int) error) {

	var (
		err      error
		exportID int
		userMap  UserMapping
	)

	// initialize the export
	_, err = client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	exportID, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.UnexportByID(defaultCtx, exportID)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify the existing mapping is mapped to nobody
	userMap, err = getMap(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, userMap)
	assertNotNil(t, userMap.User)
	assertNotNil(t, userMap.User.ID)
	assertNotNil(t, userMap.User.ID.ID)
	assert.Equal(t, "nobody", userMap.User.ID.ID)

	// update the user mapping to root
	err = enaMap(defaultCtx, exportID, "root")
	assertNoError(t, err)

	// verify the user mapping is mapped to root
	userMap, err = getMap(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, userMap)
	assertNotNil(t, userMap.User)
	assertNotNil(t, userMap.User.ID)
	assertNotNil(t, userMap.User.ID.ID)
	assert.Equal(t, "root", userMap.User.ID.ID)

	// disable the user mapping
	err = disMap(defaultCtx, exportID)
	assertNoError(t, err)

	// verify the user mapping is disabled
	userMap, err = getMap(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, userMap.Enabled)
	assert.False(t, *userMap.Enabled)
}

var (
	getClients = func(ctx context.Context, e Export) []string {
		return *e.Clients
	}
	getRootClients = func(ctx context.Context, e Export) []string {
		return *e.RootClients
	}
)

func TestExportClientsGet(t *testing.T) {
	testExportClientsGet(
		t,
		"test_get_export_clients",
		client.GetExportClientsByID,
		client.SetExportClientsByID)
}

func TestExportClientsSet(t *testing.T) {
	testExportClientsSet(
		t,
		"test_set_export_clients",
		getClients,
		client.SetExportClientsByID)
}

func TestExportClientsAdd(t *testing.T) {
	testExportClientsAdd(
		t,
		"test_add_export_clients",
		getClients,
		client.SetExportClientsByID,
		client.AddExportClientsByExportID)
}

func TestAddExportClientsByID(t *testing.T) {
  
        // Add the test exports
        volumeName1 := "test_get_exports1"
        vol, err := client.CreateVolume(defaultCtx, volumeName1)
        assertNoError(t, err)
        assertNotNil(t, vol)
        volumeName1 = vol.Name
        t.Logf("created volume: %s", volumeName1)

        exportID, err := client.Export(defaultCtx, volumeName1)
        assertNoError(t, err)
        t.Logf("created export: %d", exportID)

        exportForClient = exportID
	export, _ := client.GetExportByID(defaultCtx, exportID)

	fmt.Printf("export '%d' has \n%-20v: '%v'\n%-20v: '%v'\n%-20v: '%v'\n", exportID, "clients", *export.Clients, "read_only_cilents", *export.ReadOnlyClients, "read_write_cilents", *export.ReadWriteClients)

	testAddExportClientsByID(t, exportID, export, client.AddExportClientsByID)
	testAddExportClientsByID(t, exportID, export, client.AddExportReadOnlyClientsByID)
	testAddExportClientsByID(t, exportID, export, client.AddExportReadWriteClientsByID)

	export, _ = client.GetExportByID(defaultCtx, exportID)

	fmt.Printf("now export '%d' has \n%-20v: '%v'\n%-20v: '%v'\n%-20v: '%v'\n", exportID, "clients", *export.Clients, "read_only_cilents", *export.ReadOnlyClients, "read_write_cilents", *export.ReadWriteClients)

	assert.Contains(t, *export.Clients, "192.168.1.112")
	assert.Contains(t, *export.ReadOnlyClients, "192.168.1.112")
	assert.Contains(t, *export.ReadWriteClients, "192.168.1.112")
}

func testAddExportClientsByID(t *testing.T, exportID int, export Export, addExportClientsByID func(
	ctx context.Context, id int, clients []string) error) {

	var clientsToAdd = []string{"192.168.1.110", "192.168.1.110", "192.168.1.111", "192.168.1.112", "192.168.1.113"}

	log.Debug(defaultCtx, "add '%v' to '%v' for export '%d'", clientsToAdd, *export.Clients, exportID)

	err = addExportClientsByID(defaultCtx, exportID, clientsToAdd)
	assert.Nil(t, err)
}

func TestRemoveExportClientsByID(t *testing.T) {

	testRemoveExportClients(t, client.RemoveExportClientsByID, nil)
}

func TestRemoveExportClientsByName(t *testing.T) {

	testRemoveExportClients(t, nil, client.RemoveExportClientsByName)
        volumeName1 := "test_get_exports1"
        // make sure we clean up when we're done
        defer client.Unexport(defaultCtx, volumeName1)
        defer client.DeleteVolume(defaultCtx, volumeName1)
}

func testRemoveExportClients(t *testing.T,
	removeExportClientsByIDFunc func(ctx context.Context, id int, clientsToRemove []string) error,
	removeExportClientsByNameFunc func(ctx context.Context, name string, clientsToRemove []string) error) {

        volumeName1 := "test_get_exports1"

        exportID := exportForClient
	export, _ := client.GetExportByName(defaultCtx, volumeName1)
	exportName := volumeName1

	fmt.Printf("export '%d' has \n%-20v: '%v'\n%-20v: '%v'\n%-20v: '%v'\n", exportID, "clients", *export.Clients, "read_only_cilents", *export.ReadOnlyClients, "read_write_cilents", *export.ReadWriteClients)

	var clientsToRemove = []string{"192.168.1.110", "192.168.1.110", "192.168.1.111", "192.168.1.116", "k8s-node-1.lab.acme.com"}

	log.Debug(defaultCtx, "remove '%v' from '%v' for export '%d'", clientsToRemove, *export.Clients, exportID)

	if removeExportClientsByIDFunc != nil {
		err = removeExportClientsByIDFunc(defaultCtx, exportID, clientsToRemove)
		assert.Nil(t, err)
	} else {
		err = removeExportClientsByNameFunc(defaultCtx, exportName, clientsToRemove)
		assert.Nil(t, err)
	}

	export, _ = client.GetExportByID(defaultCtx, exportID)

	fmt.Printf("now export '%d' has \n%-20v: '%v'\n%-20v: '%v'\n%-20v: '%v'\n", exportID, "clients", *export.Clients, "read_only_cilents", *export.ReadOnlyClients, "read_write_cilents", *export.ReadWriteClients)

	assert.Contains(t, *export.Clients, "192.168.1.112")
	assert.NotContains(t, *export.Clients, "192.168.1.110")
	assert.NotContains(t, *export.Clients, "192.168.1.111")

	assert.Contains(t, *export.ReadOnlyClients, "192.168.1.112")
	assert.NotContains(t, *export.ReadOnlyClients, "192.168.1.110")
	assert.NotContains(t, *export.ReadOnlyClients, "192.168.1.111")

	assert.Contains(t, *export.ReadWriteClients, "192.168.1.112")
	assert.NotContains(t, *export.ReadWriteClients, "192.168.1.110")
	assert.NotContains(t, *export.ReadWriteClients, "192.168.1.111")
}

func TestExportClientsClear(t *testing.T) {
	testExportClientsClear(
		t,
		"test_clear_export_clients",
		getClients,
		client.SetExportClientsByID,
		client.ClearExportClientsByID)
}

func TestExportRootClientsGet(t *testing.T) {
	testExportClientsGet(
		t,
		"test_get_export_root_clients",
		client.GetExportRootClientsByID,
		client.SetExportRootClientsByID)
}

func TestExportRootClientsSet(t *testing.T) {
	testExportClientsSet(
		t,
		"test_set_export_root_clients",
		getRootClients,
		client.SetExportRootClientsByID)
}

func TestExportRootClientsAdd(t *testing.T) {
	testExportClientsAdd(
		t,
		"test_add_export_root_clients",
		getRootClients,
		client.SetExportRootClientsByID,
		client.AddExportRootClientsByID)
}

func TestExportRootClientsClear(t *testing.T) {
	testExportClientsClear(
		t,
		"test_clear_export_root_clients",
		getRootClients,
		client.SetExportRootClientsByID,
		client.ClearExportRootClientsByID)
}

func testExportClientsGet(
	t *testing.T,
	volumeName string,
	getClients func(ctx context.Context, id int) ([]string, error),
	setClients func(ctx context.Context, id int, clients ...string) error) {

	var (
		err            error
		exportID       int
		clientList     = []string{"1.2.3.4", "1.2.3.5"}
		currentClients []string
	)

	// initialize the export
	_, err = client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	exportID, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.UnexportByID(defaultCtx, exportID)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// set the export client
	err = setClients(defaultCtx, exportID, clientList...)
	assertNoError(t, err)

	// test getting the client list
	currentClients, err = getClients(defaultCtx, exportID)
	assertNoError(t, err)

	// verify we received the correct clients
	assert.Equal(t, len(clientList), len(currentClients))

	sort.Strings(currentClients)
	sort.Strings(clientList)

	for i := range currentClients {
		assert.Equal(t, currentClients[i], clientList[i])
	}
}

func testExportClientsSet(
	t *testing.T,
	volumeName string,
	getClients func(ctx context.Context, e Export) []string,
	setClients func(ctx context.Context, id int, clients ...string) error) {

	var (
		err            error
		export         Export
		exportID       int
		currentClients []string
		clientList     = []string{"1.2.3.4", "1.2.3.5"}
	)

	sort.Strings(clientList)

	// initialize the export
	_, err = client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	exportID, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.UnexportByID(defaultCtx, exportID)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify we aren't already exporting the volume to any of the clients
	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	for _, currentClient := range getClients(defaultCtx, export) {
		for _, newClient := range clientList {
			assert.NotEqual(t, currentClient, newClient)
		}
	}

	// test setting the export client
	err = setClients(defaultCtx, exportID, clientList...)
	assertNoError(t, err)

	// verify the export client was set
	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	currentClients = getClients(defaultCtx, export)
	assert.Equal(t, len(clientList), len(currentClients))

	sort.Strings(currentClients)
	for i := range currentClients {
		assert.Equal(t, currentClients[i], clientList[i])
	}
}

func testExportClientsAdd(
	t *testing.T,
	volumeName string,
	getClients func(ctx context.Context, e Export) []string,
	setClients func(ctx context.Context, id int, clients ...string) error,
	addClients func(ctx context.Context, id int, clients ...string) error) {

	var (
		err            error
		export         Export
		exportID       int
		currentClients []string
		clientList     = []string{"1.2.3.4", "1.2.3.5"}
		addedClients   = []string{"1.2.3.6", "1.2.3.7"}
		allClients     = append(clientList, addedClients...)
	)

	sort.Strings(clientList)
	sort.Strings(allClients)

	// initialize the export
	_, err = client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	exportID, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.UnexportByID(defaultCtx, exportID)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify we aren't already exporting the volume to any of the clients
	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	for _, currentClient := range getClients(defaultCtx, export) {
		for _, newClient := range clientList {
			assert.NotEqual(t, currentClient, newClient)
		}
	}

	// test setting the export client
	err = setClients(defaultCtx, exportID, clientList...)
	assertNoError(t, err)

	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	currentClients = getClients(defaultCtx, export)
	assert.Equal(t, len(clientList), len(currentClients))

	sort.Strings(currentClients)
	for i := range currentClients {
		assert.Equal(t, currentClients[i], clientList[i])
	}

	// verify that added clients are added to the list
	err = addClients(defaultCtx, exportID, addedClients...)
	assertNoError(t, err)

	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	currentClients = getClients(defaultCtx, export)
	assert.Equal(t, len(allClients), len(currentClients))

	sort.Strings(currentClients)
	for i := range currentClients {
		assert.Equal(t, currentClients[i], allClients[i])
	}
}

func testExportClientsClear(
	t *testing.T,
	volumeName string,
	getClients func(ctx context.Context, e Export) []string,
	setClients func(ctx context.Context, id int, clients ...string) error,
	nilClients func(ctx context.Context, id int) error) {

	var (
		err            error
		export         Export
		exportID       int
		currentClients []string
		clientList     = []string{"1.2.3.4", "1.2.3.5"}
	)

	sort.Strings(clientList)

	// initialize the export
	_, err = client.CreateVolume(defaultCtx, volumeName)
	assertNoError(t, err)

	exportID, err = client.Export(defaultCtx, volumeName)
	assertNoError(t, err)

	// make sure we clean up when we're done
	defer client.UnexportByID(defaultCtx, exportID)
	defer client.DeleteVolume(defaultCtx, volumeName)

	// verify we are exporting the volume
	err = setClients(defaultCtx, exportID, clientList...)
	assertNoError(t, err)

	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	currentClients = getClients(defaultCtx, export)
	assert.Equal(t, len(clientList), len(currentClients))

	sort.Strings(currentClients)
	for i := range currentClients {
		assert.Equal(t, currentClients[i], clientList[i])
	}

	// test clearing the export client
	err = nilClients(defaultCtx, exportID)
	assertNoError(t, err)

	// verify the export client was cleared
	export, err = client.GetExportByID(defaultCtx, exportID)
	assertNoError(t, err)
	assertNotNil(t, export)

	assert.Len(t, getClients(defaultCtx, export), 0)
}

func TestGetExportsWithPagination(t *testing.T) {
	// This test makes assumption that the number of exports is no less than 2
	limit := "2"
	params := api.OrderedValues{
		{[]byte("limit"), []byte(limit)},
	}
	exports, err := client.GetExportsWithParams(defaultCtx, params)
	if err != nil {
		panic(err)
	}
	limitCallID := exports.Exports[0].ID

	// Test to get the next page
	resume := exports.Resume
	params = api.OrderedValues{
		{[]byte("resume"), []byte(resume)},
	}
	exports, err = client.GetExportsWithParams(defaultCtx, params)
	if err != nil {
		if resume == "" {
			panic("The last call got the last page")
		} else {
			panic(err)
		}
	}
	resumeCallID := exports.Exports[0].ID

	// Test if the results are different
	if limitCallID == resumeCallID {
		panic("Resume call didn't get the exports of the next page")
	}

	// Test to get exports based on the resume parameter only
	exports, err = client.GetExportsWithResume(defaultCtx, resume)
	if err != nil {
		panic(err)
	}

	// Test to get exports based on the limit parameter only
	exports, err = client.GetExportsWithLimit(defaultCtx, limit)
	if err != nil {
		panic(err)
	}
}

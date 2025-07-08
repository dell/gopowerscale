/*
Copyright (c) 2022-2025 Dell Inc, or its subsidiaries.

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
	"errors"
	"fmt"
	"sort"
	"strconv"
	"testing"

	log "github.com/akutz/gournal"
	api "github.com/dell/goisilon/api"
	apiv1 "github.com/dell/goisilon/api/v1"
	apiv2 "github.com/dell/goisilon/api/v2"
	apiv4 "github.com/dell/goisilon/api/v4"
	"github.com/dell/goisilon/mocks"
	"github.com/dell/goisilon/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var exportForClient int

func TestGetExports(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export list is empty
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	exports, err := client.GetExports(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, len(exports))

	// Test case: Export list is not empty
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{
			&apiv2.Export{
				ID:    1,
				Paths: &[]string{"/path1"},
			},
			&apiv2.Export{
				ID:    2,
				Paths: &[]string{"/path2"},
			},
		}
	}).Once()
	exports, err = client.GetExports(context.Background())
	assert.NoError(t, err)
	assert.NotNil(t, exports)
	assert.Equal(t, 2, len(exports))
	assert.Equal(t, int(1), exports[0].ID)
	assert.Equal(t, "/path1", (*exports[0].Paths)[0])
	assert.Equal(t, int(2), exports[1].ID)
	assert.Equal(t, "/path2", (*exports[1].Paths)[0])

	// Test case: API error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("API error")).Once()
	exports, err = client.GetExports(context.Background())
	assert.NotNil(t, err)
	assert.Nil(t, exports)
}

func TestGetExportByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export list is empty
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	export, err := client.GetExportByID(context.Background(), 0)
	assert.NoError(t, err)
	assert.Nil(t, export)

	// Test case: Export list is not empty
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{
			&apiv2.Export{
				ID:    1,
				Paths: &[]string{"/path1"},
			},
		}
	}).Once()
	export, err = client.GetExportByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, export)
	assert.Equal(t, int(1), export.ID)
	assert.Equal(t, "/path1", (*export.Paths)[0])

	// Test case: API error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("API error")).Once()
	export, err = client.GetExportByID(context.Background(), 2)
	assert.NotNil(t, err)
	assert.Nil(t, export)
}

func TestGetExportByName(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	name := "test_volume"
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/path1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/path1").Once()
	export, err := client.GetExportByName(defaultCtx, name)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, export.ID)
	assert.Equal(t, expectedExport.Paths, export.Paths)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	export, err = client.GetExportByName(defaultCtx, name)
	assert.NoError(t, err)
	assert.Nil(t, export)

	// Test case: Error returned from ExportsList
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	export, err = client.GetExportByName(defaultCtx, name)
	assert.ErrorIs(t, err, expectedErr)
	assert.Nil(t, export)
}

func TestGetExportByNameWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	name := "test_export"
	zone := "test_zone"
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/path1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/path1").Once()
	export, err := client.GetExportByNameWithZone(defaultCtx, name, zone)
	assert.NoError(t, err)
	assert.NotNil(t, export)
	if export != nil {
		assert.Equal(t, expectedExport.ID, export.ID)
		assert.Equal(t, expectedExport.Paths, export.Paths)
	}

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	export, err = client.GetExportByNameWithZone(defaultCtx, name, zone)
	assert.NoError(t, err)
	assert.Nil(t, export)

	// Test case: Error returned by ExportsListWithZone
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	export, err = client.GetExportByNameWithZone(defaultCtx, name, zone)
	assert.Equal(t, expectedErr, err)
	assert.Nil(t, export)
}

func TestExport(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Volume already exported
	testExportID := 1
	testVolumeName := "test_volume"
	testVolumePath := "/path/to/test_volume"
	expectedExport := &apiv2.Export{
		ID:    testExportID,
		Paths: &[]string{testVolumePath},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return(testVolumePath).Once()
	exportID, err := client.Export(context.Background(), testVolumeName)
	assert.NoError(t, err)
	assert.Equal(t, testExportID, exportID)

	// Test case: Volume not exported
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return(testVolumePath).Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(*apiv2.Export)
		*resp = *expectedExport
	}).Once()
	exportID, err = client.Export(context.Background(), testVolumeName)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, exportID)

	// Test case: Error in IsExported
	testIsExportedError := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(testIsExportedError).Once()
	_, err = client.Export(context.Background(), testVolumeName)
	assert.ErrorIs(t, err, testIsExportedError)

	// Test case: Error in ExportCreate
	testExportCreateError := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return(testVolumePath).Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(testExportCreateError).Once()
	_, err = client.Export(context.Background(), testVolumeName)
	assert.ErrorIs(t, err, testExportCreateError)
}

func TestExportWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export volume with valid name and zone
	name := "test_volume"
	zone := "test_zone"
	description := "test_description"
	expectedExport := &apiv2.Export{
		Paths:       &[]string{"/path1"},
		Description: description,
	}
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/path1").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(*apiv2.Export)
		*resp = *expectedExport
	}).Once()
	exportID, err := client.ExportWithZone(defaultCtx, name, zone, description)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, exportID)

	// Test case: Export volume with invalid name
	name = ""
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	exportID, err = client.ExportWithZone(defaultCtx, name, zone, description)
	assert.NoError(t, err)
	assert.Equal(t, exportID, 0)

	// Test case: Export volume with invalid zone
	zone = ""
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	exportID, err = client.ExportWithZone(defaultCtx, name, zone, description)
	assert.Error(t, err, "zone cannot be empty")
	assert.Equal(t, exportID, 0)

	// Test case: Export volume with error from ExportCreateWithZone
	zone = "test_zone"
	testExportCreateWithZoneError := errors.New("test error")
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(testExportCreateWithZoneError).Once()
	exportID, err = client.ExportWithZone(defaultCtx, name, zone, description)
	assert.ErrorIs(t, err, testExportCreateWithZoneError)
	assert.Equal(t, exportID, 0)
}

func TestExportWithZoneAndPath(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Exporting a volume with a valid path, zone, and description
	path := "test_path"
	zone := "test_zone"
	description := "test_description"
	expectedExport := &apiv2.Export{
		ID:          1,
		Paths:       &[]string{path},
		Description: description,
	}
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(*apiv2.Export)
		*resp = *expectedExport
	}).Once()
	exportID, err := client.ExportWithZoneAndPath(defaultCtx, path, zone, description)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, exportID)

	// Test case: Exporting a volume with an invalid path
	invalidPath := ""
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(*apiv2.Export)
		*resp = *expectedExport
	}).Once()
	_, err = client.ExportWithZoneAndPath(defaultCtx, invalidPath, zone, description)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, exportID)

	// Test case: Exporting a volume with an invalid zone
	invalidZone := ""
	_, err = client.ExportWithZoneAndPath(defaultCtx, path, invalidZone, description)
	assert.Error(t, err, "zone cannot be empty")

	// Test case: Exporting a volume with an invalid description
	path = "test_path"
	zone = "test_zone"
	invalidDescription := ""
	expectedExport = &apiv2.Export{
		ID:          2,
		Paths:       &[]string{path},
		Description: invalidDescription,
	}
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(6).(*apiv2.Export)
		*resp = *expectedExport
	}).Once()
	exportID, err = client.ExportWithZoneAndPath(defaultCtx, path, zone, invalidDescription)
	assert.NoError(t, err)
	assert.Equal(t, expectedExport.ID, exportID)
}

func TestGetRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	mapping, err := client.GetRootMapping(context.Background(), "export1")
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/").Once()
	mapping, err = client.GetRootMapping(context.Background(), "export2")
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetRootMapping(context.Background(), "path3")
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestGetRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	mapping, err := client.GetRootMappingByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	mapping, err = client.GetRootMappingByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetRootMappingByID(context.Background(), 3)
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestEnableRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable root mapping for an existing export
	name := "export1"
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableRootMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Enable root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.EnableRootMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableRootMapping(defaultCtx, name, user)
	assert.Error(t, err, expectedErr)
}

func TestEnableRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable root mapping for an existing export
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableRootMappingByID(defaultCtx, 1, user)
	assert.NoError(t, err)

	// Test case: Enable root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.EnableRootMappingByID(defaultCtx, 2, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableRootMappingByID(defaultCtx, 3, user)
	assert.Error(t, err, expectedErr)
}

func TestDisableRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable root mapping for an existing export
	name := "export1"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableRootMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Disable root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.DisableRootMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableRootMapping(defaultCtx, name)
	assert.Error(t, err, expectedErr)
}

func TestDisableRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable root mapping for an existing export
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableRootMappingByID(defaultCtx, 1)
	assert.NoError(t, err)

	// Test case: Disable root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.DisableRootMappingByID(defaultCtx, 2)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableRootMappingByID(defaultCtx, 3)
	assert.Error(t, err, expectedErr)
}

func TestGetNonRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapNonRoot: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	mapping, err := client.GetNonRootMapping(context.Background(), "export1")
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/").Once()
	mapping, err = client.GetNonRootMapping(context.Background(), "export2")
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetNonRootMapping(context.Background(), "path3")
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestGetNonRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapNonRoot: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	mapping, err := client.GetNonRootMappingByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	mapping, err = client.GetNonRootMappingByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetNonRootMappingByID(context.Background(), 3)
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestEnableNonRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable non root mapping for an existing export
	name := "export1"
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableNonRootMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Enable non root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.EnableNonRootMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableNonRootMapping(defaultCtx, name, user)
	assert.Error(t, err, expectedErr)
}

func TestEnableNonRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable non root mapping for an existing export
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableNonRootMappingByID(defaultCtx, 1, user)
	assert.NoError(t, err)

	// Test case: Enable non root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.EnableNonRootMappingByID(defaultCtx, 2, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableNonRootMappingByID(defaultCtx, 3, user)
	assert.Error(t, err, expectedErr)
}

func TestDisableNonRootMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable non root mapping for an existing export
	name := "export1"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableNonRootMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Disable non root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.DisableNonRootMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableNonRootMapping(defaultCtx, name)
	assert.Error(t, err, expectedErr)
}

func TestDisableNonRootMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable non root mapping for an existing export
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableNonRootMappingByID(defaultCtx, 1)
	assert.NoError(t, err)

	// Test case: Disable non root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.DisableNonRootMappingByID(defaultCtx, 2)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableNonRootMappingByID(defaultCtx, 3)
	assert.Error(t, err, expectedErr)
}

func TestGetFailureMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapFailure: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	mapping, err := client.GetFailureMapping(context.Background(), "export1")
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/").Once()
	mapping, err = client.GetFailureMapping(context.Background(), "export2")
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetFailureMapping(context.Background(), "path3")
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestGetFailureMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export found
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/export1"},
		MapFailure: &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   "user1",
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	mapping, err := client.GetFailureMappingByID(context.Background(), 1)
	assert.NoError(t, err)
	assert.NotNil(t, mapping)
	assert.Equal(t, "user1", mapping.User.ID.ID)

	// Test case: Export not found
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{}
	}).Once()
	mapping, err = client.GetFailureMappingByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Nil(t, mapping)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	mapping, err = client.GetFailureMappingByID(context.Background(), 3)
	assert.Error(t, err, expectedErr)
	assert.Nil(t, mapping)
}

func TestEnableFailureMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable non root mapping for an existing export
	name := "export1"
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableFailureMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Enable non root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.EnableFailureMapping(defaultCtx, name, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableFailureMapping(defaultCtx, name, user)
	assert.Error(t, err, expectedErr)
}

func TestEnableFailureMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Enable non root mapping for an existing export
	user := "test_user"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.EnableFailureMappingByID(defaultCtx, 1, user)
	assert.NoError(t, err)

	// Test case: Enable non root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.EnableFailureMappingByID(defaultCtx, 2, user)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.EnableFailureMappingByID(defaultCtx, 3, user)
	assert.Error(t, err, expectedErr)
}

func TestDisableFailureMapping(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable non root mapping for an existing export
	name := "export1"
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableFailureMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Disable non root mapping for a non-existing export
	name = "export2"
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.DisableFailureMapping(defaultCtx, name)
	assert.NoError(t, err)

	// Test case: Error getting export
	name = ""
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableFailureMapping(defaultCtx, name)
	assert.Error(t, err, expectedErr)
}

func TestDisableFailureMappingByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Disable non root mapping for an existing export
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		MapRoot: &apiv2.UserMapping{},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.DisableFailureMappingByID(defaultCtx, 1)
	assert.NoError(t, err)

	// Test case: Disable non root mapping for a non-existing export
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	err = client.DisableFailureMappingByID(defaultCtx, 2)
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.DisableFailureMappingByID(defaultCtx, 3)
	assert.Error(t, err, expectedErr)
}

func TestGetExportClients(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1", "client2"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	clients, err := client.GetExportClients(defaultCtx, "export1")
	assert.NoError(t, err)
	assert.Equal(t, len(*ex.Clients), len(clients))
	assert.Equal(t, *ex.Clients, clients)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	clients, err = client.GetExportClients(defaultCtx, "export1")
	assert.NoError(t, err)
	assert.Nil(t, clients)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	clients, err = client.GetExportClients(defaultCtx, "export2")
	assert.NoError(t, err)
	assert.Nil(t, clients)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	clients, err = client.GetExportClients(defaultCtx, "export3")
	assert.Error(t, err, expectedErr)
	assert.Nil(t, clients)
}

func TestGetExportClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1", "client2"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	clients, err := client.GetExportClientsByID(defaultCtx, 1)
	assert.NoError(t, err)
	assert.Equal(t, len(*ex.Clients), len(clients))
	assert.Equal(t, *ex.Clients, clients)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	clients, err = client.GetExportClientsByID(defaultCtx, 1)
	assert.NoError(t, err)
	assert.Nil(t, clients)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	clients, err = client.GetExportClientsByID(defaultCtx, 2)
	assert.NoError(t, err)
	assert.Nil(t, clients)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	clients, err = client.GetExportClientsByID(defaultCtx, 3)
	assert.Error(t, err, expectedErr)
	assert.Nil(t, clients)
}

func TestAddExportClients(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportClients(defaultCtx, "export1", "client2")
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClients(defaultCtx, "export1", "client2")
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.AddExportClients(defaultCtx, "export2", "client2")
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportClients(defaultCtx, "export3", "client2")
	assert.Error(t, err, expectedErr)
}

func TestAddExportClientsByExportID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportClientsByExportID(defaultCtx, 1, "client2")
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByExportID(defaultCtx, 1, "client2")
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByExportID(defaultCtx, 2, "client2")
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportClientsByExportID(defaultCtx, 3, "client2")
	assert.Error(t, err, expectedErr)
}

func TestAddExportClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByID(defaultCtx, 2, []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportClientsByID(defaultCtx, 3, []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportReadOnlyClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:              1,
		Paths:           &[]string{"/export1"},
		ReadOnlyClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportReadOnlyClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.ReadOnlyClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadOnlyClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadOnlyClientsByID(defaultCtx, 2, []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportReadOnlyClientsByID(defaultCtx, 3, []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportReadWriteClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:               1,
		Paths:            &[]string{"/export1"},
		ReadWriteClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportReadWriteClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.ReadWriteClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadWriteClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadWriteClientsByID(defaultCtx, 2, []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportReadWriteClientsByID(defaultCtx, 3, []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportClientsByExportIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportClientsByExportIDWithZone(defaultCtx, 1, "zone1", false, "client2")
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByExportIDWithZone(defaultCtx, 1, "zone1", false, "client2")
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByExportIDWithZone(defaultCtx, 2, "zone1", false, "client2")
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportClientsByExportIDWithZone(defaultCtx, 3, "zone1", false, "client2")
	assert.Error(t, err, expectedErr)
}

func TestAddExportClientsByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.Clients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportClientsByIDWithZone(defaultCtx, 2, "zone1", []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportClientsByIDWithZone(defaultCtx, 3, "zone1", []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportRootClientsByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:          1,
		Paths:       &[]string{"/export1"},
		RootClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportRootClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.RootClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportRootClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportRootClientsByIDWithZone(defaultCtx, 2, "zone1", []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportRootClientsByIDWithZone(defaultCtx, 3, "zone1", []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportReadOnlyClientsByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:              1,
		Paths:           &[]string{"/export1"},
		ReadOnlyClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportReadOnlyClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.ReadOnlyClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadOnlyClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadOnlyClientsByIDWithZone(defaultCtx, 2, "zone1", []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportReadOnlyClientsByIDWithZone(defaultCtx, 3, "zone1", []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestAddExportReadWriteClientsByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:               1,
		Paths:            &[]string{"/export1"},
		ReadWriteClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportReadWriteClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.ReadWriteClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadWriteClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportReadWriteClientsByIDWithZone(defaultCtx, 2, "zone1", []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportReadWriteClientsByIDWithZone(defaultCtx, 3, "zone1", []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestRemoveExportClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:               1,
		Paths:            &[]string{"/export1"},
		Clients:          &[]string{"client1"},
		RootClients:      &[]string{"client2"},
		ReadOnlyClients:  &[]string{"client3"},
		ReadWriteClients: &[]string{"client4"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.RemoveExportClientsByID(defaultCtx, 1, []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.RemoveExportClientsByID(defaultCtx, 2, []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.RemoveExportClientsByID(defaultCtx, 3, []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func TestRemoveExportClientsByIDWithZone(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:               1,
		Paths:            &[]string{"/export1"},
		Clients:          &[]string{"client1"},
		RootClients:      &[]string{"client2"},
		ReadOnlyClients:  &[]string{"client3"},
		ReadWriteClients: &[]string{"client4"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.RemoveExportClientsByIDWithZone(defaultCtx, 1, "zone1", []string{"client2"}, false)
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.RemoveExportClientsByIDWithZone(defaultCtx, 2, "zone1", []string{"client2"}, false)
	assert.Error(t, err, "Export instance is nil, abort calling exportAddClients")

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.RemoveExportClientsByIDWithZone(defaultCtx, 3, "zone1", []string{"client2"}, false)
	assert.Error(t, err, expectedErr)
}

func testExportsList(t *testing.T) {
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

func testExportCreate(t *testing.T) {
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

func testExportDelete(t *testing.T) {
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

func testExportNonRootMapping(t *testing.T) {
	testUserMapping(
		t,
		"test_export_non_root_mapping",
		client.GetNonRootMappingByID,
		client.EnableNonRootMappingByID,
		client.DisableNonRootMappingByID)
}

func testExportFailureMapping(t *testing.T) {
	testUserMapping(
		t,
		"test_export_failure_mapping",
		client.GetFailureMappingByID,
		client.EnableFailureMappingByID,
		client.DisableFailureMappingByID)
}

func testExportRootMapping(t *testing.T) {
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
	disMap func(ctx context.Context, id int) error,
) {
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
	getClients = func(_ context.Context, e Export) []string {
		return *e.Clients
	}
	getRootClients = func(_ context.Context, e Export) []string {
		return *e.RootClients
	}
)

func testAddExportClientsByID(t *testing.T, exportID int, export Export, addExportClientsByID func(
	ctx context.Context, id int, clients []string, ignoreUnresolvableHosts bool) error,
) {
	clientsToAdd := []string{"192.168.1.110", "192.168.1.110", "192.168.1.111", "192.168.1.112", "192.168.1.113"}

	log.Debug(defaultCtx, "add '%v' to '%v' for export '%d'", clientsToAdd, *export.Clients, exportID)

	err = addExportClientsByID(defaultCtx, exportID, clientsToAdd, false)
	assert.NoError(t, err)
}

func testRemoveExportClientsByID(t *testing.T) {
	testRemoveExportClients(t, client.RemoveExportClientsByID, nil)
}

func testRemoveExportClientsByName(t *testing.T) {
	testRemoveExportClients(t, nil, client.RemoveExportClientsByName)
	volumeName1 := "test_get_exports1"
	// make sure we clean up when we're done
	defer client.Unexport(defaultCtx, volumeName1)
	defer client.DeleteVolume(defaultCtx, volumeName1)
}

func testRemoveExportClients(t *testing.T,
	removeExportClientsByIDFunc func(ctx context.Context, id int, clientsToRemove []string, ignoreUnresolvableHosts bool) error,
	removeExportClientsByNameFunc func(ctx context.Context, name string, clientsToRemove []string, ignoreUnresolvableHosts bool) error,
) {
	volumeName1 := "test_get_exports1"

	exportID := exportForClient
	export, _ := client.GetExportByName(defaultCtx, volumeName1)
	exportName := volumeName1

	fmt.Printf("export '%d' has \n%-20v: '%v'\n%-20v: '%v'\n%-20v: '%v'\n", exportID, "clients", *export.Clients, "read_only_cilents", *export.ReadOnlyClients, "read_write_cilents", *export.ReadWriteClients)

	clientsToRemove := []string{"192.168.1.110", "192.168.1.110", "192.168.1.111", "192.168.1.116", "k8s-node-1.lab.acme.com"}

	log.Debug(defaultCtx, "remove '%v' from '%v' for export '%d'", clientsToRemove, *export.Clients, exportID)

	if removeExportClientsByIDFunc != nil {
		err = removeExportClientsByIDFunc(defaultCtx, exportID, clientsToRemove, false)
		assert.NoError(t, err)
	} else {
		err = removeExportClientsByNameFunc(defaultCtx, exportName, clientsToRemove, false)
		assert.NoError(t, err)
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

func testExportRootClientsGet(t *testing.T) {
	testExportClientsGet(
		t,
		"test_get_export_root_clients",
		client.GetExportRootClientsByID,
		client.SetExportRootClientsByID)
}

func testExportRootClientsSet(t *testing.T) {
	testExportClientsSet(
		t,
		"test_set_export_root_clients",
		getRootClients,
		client.SetExportRootClientsByID)
}

func testExportRootClientsAdd(t *testing.T) {
	testExportClientsAdd(
		t,
		"test_add_export_root_clients",
		getRootClients,
		client.SetExportRootClientsByID,
		client.AddExportRootClientsByID)
}

func testExportRootClientsClear(t *testing.T) {
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
	setClients func(ctx context.Context, id int, clients ...string) error,
) {
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
	setClients func(ctx context.Context, id int, clients ...string) error,
) {
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
	addClients func(ctx context.Context, id int, clients ...string) error,
) {
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
	nilClients func(ctx context.Context, id int) error,
) {
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

func testGetExportsWithPagination(_ *testing.T) {
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
		}
		panic(err)
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

func testClientListExportsWithStructParams(t *testing.T) {
	// use limit to test pagination, would still output all shares
	limit := int32(600)
	_, err := client.ListAllExportsWithStructParams(defaultCtx, apiv4.ListV4NfsExportsParams{Limit: &limit})
	assertNil(t, err)
}

func testClientExportLifeCycleWithStructParams(t *testing.T) {
	exportName := "tf_nfs_export_test"
	defer client.DeleteVolume(defaultCtx, exportName)

	vol, err := client.CreateVolume(defaultCtx, exportName)
	assertNoError(t, err)
	assertNotNil(t, vol)
	defer client.DeleteVolume(defaultCtx, exportName)

	fullPath := client.API.VolumePath(exportName)
	res, err := client.CreateExportWithStructParams(defaultCtx, apiv4.CreateV4NfsExportRequest{
		V4NfsExport: &openapi.V2NfsExport{
			Paths: []string{fullPath},
		},
	})
	assert.NotZero(t, res.ID)
	assertNil(t, err)

	// Test list Exports
	limit := int32(1)
	listExports, err := client.ListExportsWithStructParams(defaultCtx, apiv4.ListV4NfsExportsParams{Limit: &limit})
	assertNil(t, err)
	assert.Equal(t, 1, len(listExports.Exports))

	exportID := strconv.Itoa(int(res.ID))
	// Test getExport
	getExport, err := client.GetExportWithStructParams(defaultCtx, apiv4.GetV2NfsExportRequest{
		V2NFSExportID: exportID,
	})
	assertNil(t, err)
	assert.Equal(t, 1, len(getExport.Exports))
	assert.Equal(t, res.ID, *(getExport.Exports[0].ID))

	// Test getExport
	readOnly := true
	err = client.UpdateExportWithStructParams(defaultCtx, apiv4.UpdateV4NfsExportRequest{
		V2NFSExportID: exportID,
		V2NfsExport: &openapi.V2NfsExportExtendedExtended{
			ReadOnly: &readOnly,
		},
	})
	assertNil(t, err)
	getUpdatedExport, err := client.GetExportWithStructParams(defaultCtx, apiv4.GetV2NfsExportRequest{
		V2NFSExportID: exportID,
	})
	assert.Equal(t, true, *(getUpdatedExport.Exports[0].ReadOnly))

	// Test delete export
	err = client.DeleteExportWithStructParams(defaultCtx, apiv4.DeleteV4NfsExportRequest{V2NFSExportID: exportID})
	assertNil(t, err)
	_, err = client.GetExportWithStructParams(defaultCtx, apiv4.GetV2NfsExportRequest{
		V2NFSExportID: exportID,
	})
	assertNotNil(t, err)
}

func TestClient_SetExportClients(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	err := client.SetExportClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.ErrorContains(t, err, "Could not find export")

	// test when getting export returns nil
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()

	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Twice()

	err = client.SetExportClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)

	// test when export is returned
	exports := apiv2.ExportList{
		{
			ID: 1,
			Paths: &[]string{
				"test/vol/path",
				"test/vol/path2",
			},
		},
	}
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = exports
	}).Once()
	err = client.SetExportClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)
}

func TestSetExportClientsByID(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetExportClientsByID(ctx, 1, "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)
}

func TestSetExportClientsByIDWithZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetExportClientsByIDWithZone(ctx, 1, "1.2.3.4", true, "test")
	assert.NoError(t, err)
}

func TestClearExportClients(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	err := client.ClearExportClients(ctx, "test")
	assert.NoError(t, err)
}

func TestClearExportClientsByID(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	err := client.ClearExportClientsByID(ctx, 1)
	assert.NoError(t, err)
}

func TestGetExportRootClients(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// Test return nil export
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("/export1").Times(3)
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetExportRootClients(ctx, "test")
	assert.NoError(t, err)

	// Test returning an export w/ no root clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	_, err = client.GetExportRootClients(ctx, "/export1")
	assert.NoError(t, err)

	// Test returning an export w/ no root clients
	ex = &apiv2.Export{
		ID:          1,
		Paths:       &[]string{"/export1"},
		Clients:     &[]string{"client1"},
		RootClients: &[]string{"client2"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	resp, err := client.GetExportRootClients(ctx, "/export1")
	assert.NoError(t, err)
	assertEqual(t, resp, []string{"client2"})

	// Test returning export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	_, err = client.GetExportRootClients(ctx, "/export1")
	assert.ErrorContains(t, err, "Could not find export")
}

func TestGetExportRootClientsByID(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// Test return nil export
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("/export1").Times(3)
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	_, err := client.GetExportRootClientsByID(ctx, 1)
	assert.NoError(t, err)

	// Test returning an export w/ no root clients
	ex := &apiv2.Export{
		ID:      1,
		Paths:   &[]string{"/export1"},
		Clients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	_, err = client.GetExportRootClientsByID(ctx, 1)
	assert.NoError(t, err)

	// Test returning an export w/ no root clients
	ex = &apiv2.Export{
		ID:          1,
		Paths:       &[]string{"/export1"},
		Clients:     &[]string{"client1"},
		RootClients: &[]string{"client2"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	resp, err := client.GetExportRootClientsByID(ctx, 1)
	assert.NoError(t, err)
	assertEqual(t, resp, []string{"client2"})

	// Test returning export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	_, err = client.GetExportRootClientsByID(ctx, 1)
	assert.ErrorContains(t, err, "Could not find export")
}

func TestAddExportRootClients(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has root clients
	ex := &apiv2.Export{
		ID:          1,
		Paths:       &[]string{"/export1"},
		RootClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportRootClients(defaultCtx, "export1", "client2")
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.RootClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export1").Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportRootClients(defaultCtx, "export1", "client2")
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/export2").Once()
	err = client.AddExportRootClients(defaultCtx, "export2", "client2")
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportRootClients(defaultCtx, "export3", "client2")
	assert.Error(t, err, expectedErr)
}

func TestAddExportRootClientsByID(t *testing.T) {
	client.API.(*mocks.Client).ExpectedCalls = nil
	// Test case: Export exists and has clients
	ex := &apiv2.Export{
		ID:          1,
		Paths:       &[]string{"/export1"},
		RootClients: &[]string{"client1"},
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.AddExportRootClientsByID(defaultCtx, 1, "client2")
	assert.NoError(t, err)

	// Test case: Export exists but has no clients
	ex.RootClients = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{ex}
	}).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportRootClientsByID(defaultCtx, 1, "client2")
	assert.NoError(t, err)

	// Test case: Export does not exist
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err = client.AddExportRootClientsByID(defaultCtx, 2, "client2")
	assert.NoError(t, err)

	// Test case: Error getting export
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	err = client.AddExportRootClientsByID(defaultCtx, 3, "client2")
	assert.Error(t, err, expectedErr)
}

func TestClient_SetExportRootClients(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	err := client.SetExportRootClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.ErrorContains(t, err, "Could not find export")

	// test when getting export returns nil
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()

	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Twice()

	err = client.SetExportRootClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)

	// test when export is returned
	exports := apiv2.ExportList{
		{
			ID: 1,
			Paths: &[]string{
				"test/vol/path",
				"test/vol/path2",
			},
		},
	}
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = exports
	}).Once()
	err = client.SetExportRootClients(ctx, "test", "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)
}

func TestSetExportRootClientsByID(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	err := client.SetExportRootClientsByID(ctx, 1, "1.2.3.4", "1.2.3.5")
	assert.NoError(t, err)
}

func TestClearExportRootClients(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	err := client.ClearExportRootClients(ctx, "1.2.3.5")
	assert.NoError(t, err)
}

func TestClearExportRootClientsByID(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("Put", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	err := client.ClearExportRootClientsByID(ctx, 1)
	assert.NoError(t, err)
}

func TestUnexport(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	err := client.Unexport(defaultCtx, "1.2.3.5")
	assert.ErrorContains(t, err, "Could not find export")

	// test when export is returned
	exports := apiv2.ExportList{
		{
			ID: 1,
			Paths: &[]string{
				"test/vol/path",
				"test/vol/path2",
			},
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = exports
	}).Once()

	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(nil).Once()
	err = client.Unexport(defaultCtx, "test/vol/path")
	assertNoError(t, err)
}

func TestIsExportedWithZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// True case
	expectedExport := &apiv2.Export{
		ID:    1,
		Paths: &[]string{"/path1"},
		Zone:  "zone1",
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{expectedExport}
	}).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs[0:6]...).Return("/path1").Twice()

	// client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("/path1").Once()
	isExported, id, err := client.IsExportedWithZone(ctx, "/path1", "zone1")

	assert.NoError(t, err)
	assert.True(t, isExported)
	assertEqual(t, id, 1)

	// False case
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	isExported, id, err = client.IsExportedWithZone(ctx, "/path1", "zone1")
	assert.NoError(t, err)
	assert.False(t, isExported)
	assertEqual(t, id, 0)

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("Could not find export")).Once()
	isExported, id, err = client.IsExportedWithZone(ctx, "/path1", "zone1")
	assert.ErrorContains(t, err, "Could not find export")
	assert.False(t, isExported)
	assertEqual(t, id, 0)
}

func TestUnexportWithZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Twice()
	err := client.UnexportWithZone(defaultCtx, "test/vol/path", "zone1")
	assert.ErrorContains(t, err, "Could not find export")

	// test when export is returned
	exports := apiv2.ExportList{
		{
			ID: 1,
			Paths: &[]string{
				"test/vol/path",
				"test/vol/path2",
			},
			Zone: "zone1",
		},
	}
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = exports
	}).Once()

	client.API.(*mocks.Client).On("Delete", anyArgs...).Return(nil).Once()
	err = client.UnexportWithZone(defaultCtx, "test/vol/path", "zone1")
	assertNoError(t, err)

	// test when no export is returned
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Once()
	err = client.UnexportWithZone(defaultCtx, "test/vol/path", "zone1")
	assertNoError(t, err)
}

func TestGetExportsWithParams(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	exports := apiv2.Exports{
		Digest: "test",
		Exports: apiv2.ExportList{
			{
				ID: 1,
				Paths: &[]string{
					"test/vol/path",
					"test/vol/path2",
				},
			},
		},
		Resume: "1",
		Total:  1,
	}

	values := api.NewOrderedValues([][]string{
		{"key1", "value1"},
	})
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.Exports)
		*resp = exports
	}).Once()

	gotExports, err := client.GetExportsWithParams(defaultCtx, values)
	assertEqual(t, exports.Exports, gotExports.Exports)
	assertNoError(t, err)

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("Could not find export")).Once()
	_, err = client.GetExportsWithParams(defaultCtx, values)
	assert.ErrorContains(t, err, "Could not find export")
}

func TestGetExportsWithResume(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	exports := apiv2.Exports{
		Digest: "test",
		Exports: apiv2.ExportList{
			{
				ID: 1,
				Paths: &[]string{
					"test/vol/path",
					"test/vol/path2",
				},
			},
		},
		Resume: "1",
		Total:  1,
	}

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.Exports)
		*resp = exports
	}).Once()

	gotExports, err := client.GetExportsWithResume(defaultCtx, "1")
	assertEqual(t, exports.Exports, gotExports.Exports)
	assertNoError(t, err)

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("Could not find export")).Once()
	_, err = client.GetExportsWithResume(defaultCtx, "maybe")
	assert.ErrorContains(t, err, "Could not find export")
}

func TestGetExportsWithLimit(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	exports := apiv2.Exports{
		Digest: "test",
		Exports: apiv2.ExportList{
			{
				ID: 1,
				Paths: &[]string{
					"test/vol/path",
					"test/vol/path2",
				},
			},
		},
		Resume: "1",
		Total:  1,
	}

	client.API.(*mocks.Client).On("Get", anyArgs...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.Exports)
		*resp = exports
	}).Once()

	gotExports, err := client.GetExportsWithLimit(defaultCtx, "limit")
	assertEqual(t, exports.Exports, gotExports.Exports)
	assertNoError(t, err)

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(errors.New("Could not find export")).Once()
	_, err = client.GetExportsWithLimit(defaultCtx, "limit")
	assert.ErrorContains(t, err, "Could not find export")
}

func TestExportSnapshotWithZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	client.API.(*mocks.Client).On("VolumePath", anyArgs...).Return("test/vol/path").Once()
	client.API.(*mocks.Client).On("Post", anyArgs...).Return(nil).Once()
	i, err := client.ExportSnapshotWithZone(defaultCtx, "test-snap", "test-vol", "test-zone", "test description")
	assertEqual(t, i, 0)
	assertNoError(t, err)
}

func TestGetExportWithPath(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	expectedExport := apiv2.Export{
		ID:    1,
		Paths: &[]string{"/test/vol/path"},
		Zone:  "zone1",
	}

	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{&expectedExport}
	}).Once()
	gotExport, err := client.GetExportWithPath(defaultCtx, "test/vol/path")
	assertNoError(t, err)
	assert.Equal(t, expectedExport, *gotExport)
}

func TestGetExportWithPathAndZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	expectedExport := apiv2.Export{
		ID:    1,
		Paths: &[]string{"/test/vol/path"},
		Zone:  "zone1",
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{&expectedExport}
	}).Once()
	gotExport, err := client.GetExportWithPathAndZone(defaultCtx, "test/vol/path", "zone1")
	assertNoError(t, err)
	assert.Equal(t, expectedExport, *gotExport)
}

func TestGetExportByIDWithZone(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	expectedExport := apiv2.Export{
		ID:    1,
		Paths: &[]string{"/test/vol/path"},
		Zone:  "zone1",
	}
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*apiv2.ExportList)
		*resp = apiv2.ExportList{&expectedExport}
	}).Once()
	gotExport, err := client.GetExportByIDWithZone(defaultCtx, 1, "zone1")
	assertNoError(t, err)
	assert.Equal(t, expectedExport, *gotExport)
}

func TestListAllExportsWithStructParams(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}

	zone := "zone1"
	digest := "test-digest"
	resume := "test-resume"
	id := int32(1)
	params := apiv4.ListV4NfsExportsParams{
		Zone: &zone,
	}

	// test when getting export returns error
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	_, err := client.ListAllExportsWithStructParams(defaultCtx, params)
	assert.ErrorContains(t, err, "Could not find export")

	expectedExport := openapi.V2NfsExports{
		Digest: &digest,
		Exports: []openapi.V2NfsExportExtended{
			{
				ID:    &id,
				Paths: []string{"/test/vol/path"},
				Zone:  &zone,
			},
		},
		Resume: &resume,
	}

	emptyExports := openapi.V2NfsExports{
		Digest:  &digest,
		Exports: nil,
		Resume:  nil,
	}
	// test when export is returned
	// first get should return expectedExport
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V2NfsExports)
		*resp = expectedExport
	}).Once()

	// second get should return empty exports list, to avoid infinite loop
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V2NfsExports)
		*resp = emptyExports
	}).Once()

	gotExport, err := client.ListAllExportsWithStructParams(defaultCtx, params)
	assertNoError(t, err)
	assertEqual(t, expectedExport.Exports, gotExport)

	// try returning error on second listNFS call
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(*openapi.V2NfsExports)
		*resp = expectedExport
	}).Once()
	client.API.(*mocks.Client).On("Get", anyArgs...).Return(errors.New("Could not find export")).Once()
	_, err = client.ListAllExportsWithStructParams(defaultCtx, params)
	assert.ErrorContains(t, err, "Could not find export")
}
func TestGetExportsCountAttachedToNode(t *testing.T) {
	client := &Client{}
	client.API = &mocks.Client{}
	ctx := context.Background()

	// Test case: API returns an error
	expectedErr := errors.New("test error")
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(expectedErr).Once()
	_, err := client.GetExportsCountAttachedToNode(ctx, "nodeip")
	assert.Error(t, err, expectedErr)

	// Test case: One export with matching client
	client.API.(*mocks.Client).ExpectedCalls = nil
	client.API.(*mocks.Client).On("Get", anyArgs[0:6]...).Return(nil).Run(func(args mock.Arguments) {
		resp := args.Get(5).(**apiv1.GetIsiExportsResp)
		*resp = &apiv1.GetIsiExportsResp{
			ExportList: []*apiv1.IsiExport{
				{Clients: []string{"nodeip"}},
			},
		}
	}).Once()
	count, err := client.GetExportsCountAttachedToNode(ctx, "nodeip")
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)
}

/*
Copyright (c) 2019-2022 Dell Inc, or its subsidiaries.

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

	apiv4 "github.com/dell/goisilon/api/v4"
	"github.com/dell/goisilon/openapi"

	api "github.com/dell/goisilon/api"
	"github.com/dell/goisilon/api/common/utils"
	apiv2 "github.com/dell/goisilon/api/v2"
)

// ExportList is a list of Isilon Exports.
type ExportList []*apiv2.Export

// Export is an Isilon Export
type Export *apiv2.Export

// Exports is the whole struct of return JSON from the exports REST API
type Exports *apiv2.Exports

// UserMapping maps to the ISI <user-mapping> type.
type UserMapping *apiv2.UserMapping

// GetExports returns a list of all exports on the cluster
func (c *Client) GetExports(ctx context.Context) (ExportList, error) {
	return apiv2.ExportsList(ctx, c.API)
}

// GetExportByID returns an export with the provided ID.
func (c *Client) GetExportByID(ctx context.Context, id int) (Export, error) {
	return apiv2.ExportInspect(ctx, c.API, id)
}

// GetExportByName returns the first export with a path for the provided
// volume name.
func (c *Client) GetExportByName(
	ctx context.Context, name string,
) (Export, error) {
	exports, err := apiv2.ExportsList(ctx, c.API)
	if err != nil {
		return nil, err
	}
	path := c.API.VolumePath(name)
	for _, ex := range exports {
		for _, p := range *ex.Paths {
			if p == path {
				return ex, nil
			}
		}
	}
	return nil, nil
}

// GetExportByNameWithZone returns the first export with a path for the provided
// volume name in the given zone.
func (c *Client) GetExportByNameWithZone(
	ctx context.Context, name, zone string,
) (Export, error) {
	exports, err := apiv2.ExportsListWithZone(ctx, c.API, zone)
	if err != nil {
		return nil, err
	}
	path := c.API.VolumePath(name)
	for _, ex := range exports {
		for _, p := range *ex.Paths {
			if p == path {
				return ex, nil
			}
		}
	}
	return nil, nil
}

// Export the volume with a given name on the cluster
func (c *Client) Export(ctx context.Context, name string) (int, error) {
	ok, id, err := c.IsExported(ctx, name)
	if err != nil {
		return 0, err
	}
	if ok {
		return id, nil
	}

	paths := []string{c.API.VolumePath(name)}

	return apiv2.ExportCreate(
		ctx, c.API,
		&apiv2.Export{Paths: &paths})
}

// ExportWithZone exports the volume with a given name and zone on the cluster
func (c *Client) ExportWithZone(ctx context.Context, name, zone, description string) (int, error) {
	// Removed the call to c.IsExportedWithZone(ctx, name, zone) to check if the path has already been exported:
	// 1. the POST /platform/2/protocols/nfs/exports API will return 500 error if the path is already exported, so there won't be false positive
	// 2. c.IsExportedWithZone(ctx, name, zone) iterates through the full list of exports, which could be expensive in a scaled environment, the potential pagination on the result set could also add to the complexity

	paths := []string{c.API.VolumePath(name)}

	return apiv2.ExportCreateWithZone(
		ctx, c.API,
		&apiv2.Export{Paths: &paths, Description: description},
		zone)
}

// ExportWithZoneAndPath exports the volume with a given name, zone and path on the cluster
func (c *Client) ExportWithZoneAndPath(ctx context.Context, path, zone, description string) (int, error) {
	paths := []string{path}

	return apiv2.ExportCreateWithZone(
		ctx, c.API,
		&apiv2.Export{Paths: &paths, Description: description},
		zone)
}

// GetRootMapping returns the root mapping for an Export.
func (c *Client) GetRootMapping(
	ctx context.Context, name string,
) (UserMapping, error) {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapRoot, nil
}

// GetRootMappingByID returns the root mapping for an Export.
func (c *Client) GetRootMappingByID(
	ctx context.Context, id int,
) (UserMapping, error) {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapRoot, nil
}

// EnableRootMapping enables the root mapping for an Export.
func (c *Client) EnableRootMapping(
	ctx context.Context, name, user string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapRoot: ex.MapRoot}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapRoot },
		func(e Export, m UserMapping) { e.MapRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// EnableRootMappingByID enables the root mapping for an Export.
func (c *Client) EnableRootMappingByID(
	ctx context.Context, id int, user string,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapRoot: ex.MapRoot}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapRoot },
		func(e Export, m UserMapping) { e.MapRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableRootMapping disables the root mapping for an Export.
func (c *Client) DisableRootMapping(
	ctx context.Context, name string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapRoot: ex.MapRoot}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapRoot },
		func(e Export, m UserMapping) { e.MapRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableRootMappingByID disables the root mapping for an Export.
func (c *Client) DisableRootMappingByID(
	ctx context.Context, id int,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapRoot: ex.MapRoot}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapRoot },
		func(e Export, m UserMapping) { e.MapRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// GetNonRootMapping returns the map_non_root mapping for an Export.
func (c *Client) GetNonRootMapping(
	ctx context.Context, name string,
) (UserMapping, error) {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapNonRoot, nil
}

// GetNonRootMappingByID returns the map_non_root mapping for an Export.
func (c *Client) GetNonRootMappingByID(
	ctx context.Context, id int,
) (UserMapping, error) {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapNonRoot, nil
}

// EnableNonRootMapping enables the map_non_root mapping for an Export.
func (c *Client) EnableNonRootMapping(
	ctx context.Context, name, user string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapNonRoot: ex.MapNonRoot}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapNonRoot },
		func(e Export, m UserMapping) { e.MapNonRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// EnableNonRootMappingByID enables the map_non_root mapping for an Export.
func (c *Client) EnableNonRootMappingByID(
	ctx context.Context, id int, user string,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapNonRoot: ex.MapNonRoot}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapNonRoot },
		func(e Export, m UserMapping) { e.MapNonRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableNonRootMapping disables the map_non_root mapping for an Export.
func (c *Client) DisableNonRootMapping(
	ctx context.Context, name string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapNonRoot: ex.MapNonRoot}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapNonRoot },
		func(e Export, m UserMapping) { e.MapNonRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableNonRootMappingByID disables the map_non_root mapping for an Export.
func (c *Client) DisableNonRootMappingByID(
	ctx context.Context, id int,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapNonRoot: ex.MapNonRoot}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapNonRoot },
		func(e Export, m UserMapping) { e.MapNonRoot = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// GetFailureMapping returns the map_failure mapping for an Export.
func (c *Client) GetFailureMapping(
	ctx context.Context, name string,
) (UserMapping, error) {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapFailure, nil
}

// GetFailureMappingByID returns the map_failure mapping for an Export.
func (c *Client) GetFailureMappingByID(
	ctx context.Context, id int,
) (UserMapping, error) {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	return ex.MapFailure, nil
}

// EnableFailureMapping enables the map_failure mapping for an Export.
func (c *Client) EnableFailureMapping(
	ctx context.Context, name, user string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapFailure: ex.MapFailure}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapFailure },
		func(e Export, m UserMapping) { e.MapFailure = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// EnableFailureMappingByID enables the map_failure mapping for an Export.
func (c *Client) EnableFailureMappingByID(
	ctx context.Context, id int, user string,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapFailure: ex.MapFailure}

	setUserMapping(
		nex,
		user,
		true,
		func(e Export) UserMapping { return e.MapFailure },
		func(e Export, m UserMapping) { e.MapFailure = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableFailureMapping disables the map_failure mapping for an Export.
func (c *Client) DisableFailureMapping(
	ctx context.Context, name string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapFailure: ex.MapFailure}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapFailure },
		func(e Export, m UserMapping) { e.MapFailure = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

// DisableFailureMappingByID disables the map_failure mapping for an Export.
func (c *Client) DisableFailureMappingByID(
	ctx context.Context, id int,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}

	nex := &apiv2.Export{ID: ex.ID, MapFailure: ex.MapFailure}

	setUserMapping(
		nex,
		"nobody",
		false,
		func(e Export) UserMapping { return e.MapFailure },
		func(e Export, m UserMapping) { e.MapFailure = m })

	return apiv2.ExportUpdate(ctx, c.API, nex)
}

func setUserMapping(
	ex Export,
	user string,
	enabled bool,
	getMapping func(Export) UserMapping,
	setMapping func(Export, UserMapping),
) {
	m := getMapping(ex)
	if m == nil {
		m = &apiv2.UserMapping{
			User: &apiv2.Persona{
				ID: &apiv2.PersonaID{
					ID:   user,
					Type: apiv2.PersonaIDTypeUser,
				},
			},
		}
		setMapping(ex, m)
		return
	}

	if m.Enabled != nil || !enabled {
		m.Enabled = &enabled
	}

	if m.User == nil {
		m.User = &apiv2.Persona{
			ID: &apiv2.PersonaID{
				ID:   user,
				Type: apiv2.PersonaIDTypeUser,
			},
		}
		return
	}

	u := m.User
	if u.ID != nil {
		u.ID.ID = user
		return
	}

	u.Name = &user
}

// GetExportClients returns an Export's clients property.
func (c *Client) GetExportClients(
	ctx context.Context, name string,
) ([]string, error) {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	if ex.Clients == nil {
		return nil, nil
	}
	return *ex.Clients, nil
}

// GetExportClientsByID returns an Export's clients property.
func (c *Client) GetExportClientsByID(
	ctx context.Context, id int,
) ([]string, error) {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	if ex.Clients == nil {
		return nil, nil
	}
	return *ex.Clients, nil
}

// AddExportClients adds to the Export's clients property.
func (c *Client) AddExportClients(
	ctx context.Context, name string, clients ...string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}
	addClients := ex.Clients
	if addClients == nil {
		addClients = &clients
	} else {
		*addClients = append(*addClients, clients...)
	}
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: ex.ID, Clients: addClients})
}

// AddExportClientsByExportID adds to the Export's clients property.
func (c *Client) AddExportClientsByExportID(
	ctx context.Context, id int, clients ...string,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}
	addClients := ex.Clients
	if addClients == nil {
		addClients = &clients
	} else {
		*addClients = append(*addClients, clients...)
	}
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: ex.ID, Clients: addClients})
}

// AddExportClientsByID adds to the Export's clients property.
func (c *Client) AddExportClientsByID(
	ctx context.Context, id int, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}

	return c.exportAddClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportReadOnlyClientsByID adds to the Export's read-only clients property.
func (c *Client) AddExportReadOnlyClientsByID(
	ctx context.Context, id int, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}

	return c.exportAddReadOnlyClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportReadWriteClientsByID adds to the Export's read-write clients property.
func (c *Client) AddExportReadWriteClientsByID(
	ctx context.Context, id int, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}

	return c.exportAddReadWriteClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportClientsByExportIDWithZone adds to the Export's clients with access zone property.
func (c *Client) AddExportClientsByExportIDWithZone(
	ctx context.Context, id int, zone string, ignoreUnresolvableHosts bool, clients ...string,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}
	if export == nil {
		return nil
	}
	addClients := export.Clients
	if addClients == nil {
		addClients = &clients
	} else {
		*addClients = append(*addClients, clients...)
	}
	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, Clients: addClients}, export.Zone, ignoreUnresolvableHosts)
}

// AddExportClientsByIDWithZone adds to the Export's clients property.
func (c *Client) AddExportClientsByIDWithZone(
	ctx context.Context, id int, zone string, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}
	return c.exportAddClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportRootClientsByIDWithZone adds to the Export's clients property.
func (c *Client) AddExportRootClientsByIDWithZone(
	ctx context.Context, id int, zone string, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}
	return c.exportAddRootClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportReadOnlyClientsByIDWithZone adds to the Export's read-only clients property.
func (c *Client) AddExportReadOnlyClientsByIDWithZone(
	ctx context.Context, id int, zone string, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}

	return c.exportAddReadOnlyClients(ctx, export, clients, ignoreUnresolvableHosts)
}

// AddExportReadWriteClientsByIDWithZone adds to the Export's read-write clients property.
func (c *Client) AddExportReadWriteClientsByIDWithZone(
	ctx context.Context, id int, zone string, clients []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}

	return c.exportAddReadWriteClients(ctx, export, clients, ignoreUnresolvableHosts)
}

func (c *Client) exportAddClients(ctx context.Context, export Export, clientsToAdd []string, ignoreUnresolvableHosts bool) error {
	if export == nil {
		return errors.New("Export instance is nil, abort calling exportAddClients")
	}
	updatedClients := c.getUpdatedClients(ctx, export.ID, export.Clients, clientsToAdd)
	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, Clients: updatedClients}, export.Zone, ignoreUnresolvableHosts)
}

func (c *Client) exportAddRootClients(ctx context.Context, export Export, clientsToAdd []string, ignoreUnresolvableHosts bool) error {
	if export == nil {
		return errors.New("Export instance is nil, abort calling exportAddRootClients")
	}
	updatedClients := c.getUpdatedClients(ctx, export.ID, export.RootClients, clientsToAdd)
	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, RootClients: updatedClients}, export.Zone, ignoreUnresolvableHosts)
}

func (c *Client) exportAddReadOnlyClients(ctx context.Context, export Export, clientsToAdd []string, ignoreUnresolvableHosts bool) error {
	if export == nil {
		return errors.New("Export instance is nil, abort calling exportAddReadOnlyClients")
	}

	updatedReadOnlyClients := c.getUpdatedClients(ctx, export.ID, export.ReadOnlyClients, clientsToAdd)

	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, ReadOnlyClients: updatedReadOnlyClients}, export.Zone, ignoreUnresolvableHosts)
}

func (c *Client) exportAddReadWriteClients(ctx context.Context, export Export, clientsToAdd []string, ignoreUnresolvableHosts bool) error {
	if export == nil {
		return errors.New("Export instance is nil, abort calling exportAddReadWriteClients")
	}
	updatedReadWriteClients := c.getUpdatedClients(ctx, export.ID, export.ReadWriteClients, clientsToAdd)

	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, ReadWriteClients: updatedReadWriteClients}, export.Zone, ignoreUnresolvableHosts)
}

func (c *Client) getUpdatedClients(ctx context.Context, exportID int, clients *[]string, clientsToAdd []string) *[]string {
	if clients == nil {
		clients = &clientsToAdd
	} else {
		// ensure uniqueness, if the client to be added is already in, skip it
		clientsToAdd = utils.RemoveStringsFromSlice(*clients, clientsToAdd)
		*clients = append(*clients, clientsToAdd...)
	}

	return clients
}

// RemoveExportClientsByID removes the given clients from the Export's clients/read_only_clients/read_write_clients properties.
func (c *Client) RemoveExportClientsByID(
	ctx context.Context, id int, clientsToRemove []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	return c.removeExportClients(ctx, export, clientsToRemove, ignoreUnresolvableHosts)
}

// RemoveExportClientsByIDWithZone removes the given clients from the
// Export's clients/read_only_clients/read_write_clients properties in a specified access zone.
func (c *Client) RemoveExportClientsByIDWithZone(
	ctx context.Context, id int, zone string, clientsToRemove []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByIDWithZone(ctx, id, zone)
	if err != nil {
		return err
	}
	return c.removeExportClients(ctx, export, clientsToRemove, ignoreUnresolvableHosts)
}

// RemoveExportClientsByName removes the given clients from the Export's clients/read_only_clients/read_write_clients properties.
func (c *Client) RemoveExportClientsByName(
	ctx context.Context, name string, clientsToRemove []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}

	return c.removeExportClients(ctx, export, clientsToRemove, ignoreUnresolvableHosts)
}

// RemoveExportClientsWithPathAndZone removes the given clients from
// the Export's clients/read_only_clients/read_write_clients/root_clients properties with export path and access zone.
func (c *Client) RemoveExportClientsWithPathAndZone(
	ctx context.Context, path, zone string, clientsToRemove []string, ignoreUnresolvableHosts bool,
) error {
	export, err := c.GetExportWithPathAndZone(ctx, path, zone)
	if err != nil {
		return err
	}

	return c.removeExportClients(ctx, export, clientsToRemove, ignoreUnresolvableHosts)
}

func (c *Client) removeExportClients(ctx context.Context, export Export, clientsToRemove []string, ignoreUnresolvableHosts bool) error {
	if export == nil {
		return errors.New("Export instance is nil, abort calling ExportRemoveClients")
	}

	clients := export.Clients
	readOnlyClients := export.ReadOnlyClients
	readWriteClients := export.ReadWriteClients
	rootClients := export.RootClients

	*clients = c.removeClients(clientsToRemove, *clients)
	*readOnlyClients = c.removeClients(clientsToRemove, *readOnlyClients)
	*readWriteClients = c.removeClients(clientsToRemove, *readWriteClients)
	*rootClients = c.removeClients(clientsToRemove, *rootClients)
	return apiv2.ExportUpdateWithZone(
		ctx, c.API, &apiv2.Export{ID: export.ID, Clients: clients, ReadOnlyClients: readOnlyClients, ReadWriteClients: readWriteClients, RootClients: rootClients}, export.Zone, ignoreUnresolvableHosts)
}

func (c *Client) removeClients(clientsToRemove []string, sourceClients []string) []string {
	if sourceClients == nil {
		return nil
	}
	return utils.RemoveStringsFromSlice(clientsToRemove, sourceClients)
}

// SetExportClients sets the Export's clients property.
func (c *Client) SetExportClients(
	ctx context.Context, name string, clients ...string,
) error {
	ok, id, err := c.IsExported(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return apiv2.ExportUpdate(ctx, c.API, &apiv2.Export{ID: id, Clients: &clients})
}

// SetExportClientsByID sets the Export's clients property.
func (c *Client) SetExportClientsByID(
	ctx context.Context, id int, clients ...string,
) error {
	return apiv2.ExportUpdate(ctx, c.API, &apiv2.Export{ID: id, Clients: &clients})
}

// SetExportClientsByIDWithZone sets the Export's clients with access zone property.
func (c *Client) SetExportClientsByIDWithZone(
	ctx context.Context, id int, zone string, ignoreUnresolvableHosts bool, clients ...string,
) error {
	return apiv2.ExportUpdateWithZone(ctx, c.API, &apiv2.Export{ID: id, Clients: &clients}, zone, ignoreUnresolvableHosts)
}

// ClearExportClients sets the Export's clients property to nil.
func (c *Client) ClearExportClients(
	ctx context.Context, name string,
) error {
	return c.SetExportClients(ctx, name, []string{}...)
}

// ClearExportClientsByID sets the Export's clients property to nil.
func (c *Client) ClearExportClientsByID(
	ctx context.Context, id int,
) error {
	return c.SetExportClientsByID(ctx, id, []string{}...)
}

// GetExportRootClients returns an Export's root_clients property.
func (c *Client) GetExportRootClients(
	ctx context.Context, name string,
) ([]string, error) {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	if ex.RootClients == nil {
		return nil, nil
	}
	return *ex.RootClients, nil
}

// GetExportRootClientsByID returns an Export's clients property.
func (c *Client) GetExportRootClientsByID(
	ctx context.Context, id int,
) ([]string, error) {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex == nil {
		return nil, nil
	}
	if ex.RootClients == nil {
		return nil, nil
	}
	return *ex.RootClients, nil
}

// AddExportRootClients adds to the Export's root_clients property.
func (c *Client) AddExportRootClients(
	ctx context.Context, name string, clients ...string,
) error {
	ex, err := c.GetExportByName(ctx, name)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}
	addClients := ex.RootClients
	if addClients == nil {
		addClients = &clients
	} else {
		*addClients = append(*addClients, clients...)
	}
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: ex.ID, RootClients: addClients})
}

// AddExportRootClientsByID adds to the Export's root_clients property.
func (c *Client) AddExportRootClientsByID(
	ctx context.Context, id int, clients ...string,
) error {
	ex, err := c.GetExportByID(ctx, id)
	if err != nil {
		return err
	}
	if ex == nil {
		return nil
	}
	addClients := ex.RootClients
	if addClients == nil {
		addClients = &clients
	} else {
		*addClients = append(*addClients, clients...)
	}
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: ex.ID, RootClients: addClients})
}

// SetExportRootClients sets the Export's root_clients property.
func (c *Client) SetExportRootClients(
	ctx context.Context, name string, clients ...string,
) error {
	ok, id, err := c.IsExported(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: id, RootClients: &clients})
}

// SetExportRootClientsByID sets the Export's clients property.
func (c *Client) SetExportRootClientsByID(
	ctx context.Context, id int, clients ...string,
) error {
	return apiv2.ExportUpdate(
		ctx, c.API, &apiv2.Export{ID: id, RootClients: &clients})
}

// ClearExportRootClients sets the Export's root_clients property to nil.
func (c *Client) ClearExportRootClients(
	ctx context.Context, name string,
) error {
	return c.SetExportRootClients(ctx, name, []string{}...)
}

// ClearExportRootClientsByID sets the Export's clients property to nil.
func (c *Client) ClearExportRootClientsByID(
	ctx context.Context, id int,
) error {
	return c.SetExportRootClientsByID(ctx, id, []string{}...)
}

// Unexport stops exporting a given volume from the cluster.
func (c *Client) Unexport(
	ctx context.Context, name string,
) error {
	ok, id, err := c.IsExported(ctx, name)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return c.UnexportByID(ctx, id)
}

// UnexportWithZone stops exporting a given volume in the given from the cluster.
func (c *Client) UnexportWithZone(
	ctx context.Context, name, zone string,
) error {
	ok, id, err := c.IsExportedWithZone(ctx, name, zone)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	return c.UnexportByIDWithZone(ctx, id, zone)
}

// UnexportByID unexports an Export by its ID.
func (c *Client) UnexportByID(
	ctx context.Context, id int,
) error {
	return apiv2.Unexport(ctx, c.API, id)
}

// UnexportByIDWithZone unexports an Export by its ID and zone.
func (c *Client) UnexportByIDWithZone(
	ctx context.Context, id int, zone string,
) error {
	return apiv2.UnexportWithZone(ctx, c.API, id, zone)
}

// IsExported returns a flag and export ID if the provided volume name is
// already exported.
func (c *Client) IsExported(
	ctx context.Context, name string,
) (bool, int, error) {
	export, err := c.GetExportByName(ctx, name)
	if err != nil {
		return false, 0, err
	}
	if export == nil {
		return false, 0, nil
	}
	return true, export.ID, nil
}

// IsExportedWithZone returns a flag and export ID if the provided volume name in the
// specified zone isalready exported.
func (c *Client) IsExportedWithZone(
	ctx context.Context, name, zone string,
) (bool, int, error) {
	export, err := c.GetExportByNameWithZone(ctx, name, zone)
	if err != nil {
		return false, 0, err
	}
	if export == nil {
		return false, 0, nil
	}
	return true, export.ID, nil
}

// GetExportsWithParams returns exports based on the parameters
func (c *Client) GetExportsWithParams(
	ctx context.Context, params api.OrderedValues,
) (Exports, error) {
	exports, err := apiv2.ExportsListWithParams(ctx, c.API, params)
	if err != nil {
		return nil, err
	}
	return exports, nil
}

// GetExportsWithResume returns the next page of exports
// based on the resume token from the previous call.
func (c *Client) GetExportsWithResume(
	ctx context.Context, resume string,
) (Exports, error) {
	exports, err := apiv2.ExportsListWithResume(ctx, c.API, resume)
	if err != nil {
		return nil, err
	}
	return exports, nil
}

// GetExportsWithLimit returns a number of exports in the default sequence
// and the number is the parameter limit.
func (c *Client) GetExportsWithLimit(
	ctx context.Context, limit string,
) (Exports, error) {
	exports, err := apiv2.ExportsListWithLimit(ctx, c.API, limit)
	if err != nil {
		return nil, err
	}
	return exports, nil
}

// ExportSnapshotWithZone exports the given snapshot and zone on the cluster
func (c *Client) ExportSnapshotWithZone(ctx context.Context, snapshotName, volumeName, zone, description string) (int, error) {
	path := apiv2.GetAbsoluteSnapshotPath(c.API, snapshotName, volumeName)
	return c.ExportPathWithZone(ctx, path, zone, description)
}

// ExportPathWithZone exports the given path and zone on the cluster
func (c *Client) ExportPathWithZone(ctx context.Context, path, zone, description string) (int, error) {
	paths := []string{path}
	return apiv2.ExportCreateWithZone(
		ctx, c.API,
		&apiv2.Export{Paths: &paths, Description: description},
		zone)
}

// GetExportWithPath gets the export with target path
func (c *Client) GetExportWithPath(
	ctx context.Context, path string,
) (Export, error) {
	return apiv2.GetExportWithPath(ctx, c.API, path)
}

// GetExportWithPathAndZone gets the export with target path and access zone
func (c *Client) GetExportWithPathAndZone(
	ctx context.Context, path, zone string,
) (Export, error) {
	return apiv2.GetExportWithPathAndZone(ctx, c.API, path, zone)
}

// GetExportByIDWithZone gets the export by export id and access zone
func (c *Client) GetExportByIDWithZone(ctx context.Context, id int, zone string) (Export, error) {
	return apiv2.GetExportByIDWithZone(ctx, c.API, id, zone)
}

// ListAllExportsWithStructParams lists all the exports with parameters
func (c *Client) ListAllExportsWithStructParams(ctx context.Context, params apiv4.ListV4NfsExportsParams) ([]openapi.V2NfsExportExtended, error) {
	var result []openapi.V2NfsExportExtended
	exports, err := apiv4.ListNfsExports(ctx, params, c.API)
	result = exports.Exports
	if err != nil {
		return nil, err
	}
	for exports.Resume != nil {
		resumeParam := apiv4.ListV4NfsExportsParams{Resume: exports.Resume}
		exports, err = apiv4.ListNfsExports(ctx, resumeParam, c.API)
		if err != nil {
			return nil, err
		}
		result = append(result, exports.Exports...)
	}
	return result, nil
}

// ListExportsWithStructParams lists all the exports with parameters
func (c *Client) ListExportsWithStructParams(ctx context.Context, params apiv4.ListV4NfsExportsParams) (*openapi.V2NfsExports, error) {
	return apiv4.ListNfsExports(ctx, params, c.API)
}

// GetExportWithStructParams list specific export with parameters
func (c *Client) GetExportWithStructParams(ctx context.Context, params apiv4.GetV2NfsExportRequest) (*openapi.V2NfsExportsExtended, error) {
	return apiv4.GetNfsExport(ctx, params, c.API)
}

// CreateExportWithStructParams create export with parameters
func (c *Client) CreateExportWithStructParams(ctx context.Context, params apiv4.CreateV4NfsExportRequest) (*openapi.Createv3EventEventResponse, error) {
	return apiv4.CreateNfsExport(ctx, params, c.API)
}

// DeleteExportWithStructParams delete export with parameters
func (c *Client) DeleteExportWithStructParams(ctx context.Context, params apiv4.DeleteV4NfsExportRequest) error {
	return apiv4.DeleteNfsExport(ctx, params, c.API)
}

// UpdateExportWithStructParams update export with parameters
func (c *Client) UpdateExportWithStructParams(ctx context.Context, params apiv4.UpdateV4NfsExportRequest) error {
	return apiv4.UpdateNfsExport(ctx, params, c.API)
}

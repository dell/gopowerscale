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
package v2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/dell/goisilon/api"
	"github.com/dell/goisilon/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/dell/goisilon/api/json"
)

var anyArgs = []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything}

func TestExportEncodeJSON(t *testing.T) {
	clients := []string{}
	ex := &Export{ID: 3, Clients: &clients}
	buf, err := json.Marshal(ex)
	if err != nil {
		t.Fatal(err)
	}
	s := string(buf)
	if !assert.Equal(t, `{"clients":[]}`, s) {
		t.FailNow()
	}
	t.Log(s)
}

func TestExportDecodeJSON(t *testing.T) {
	j := `{"id":3,"clients":[]}`
	var ex Export
	if err := json.Unmarshal([]byte(j), &ex); err != nil {
		t.Fatal(err)
	}
	fmt.Fprintf(os.Stdout, "%+v\n", ex)
}

func TestPersonaIDTypeMarshal(t *testing.T) {
	pidt := PersonaIDTypeUser
	assert.Equal(t, "user", pidt.String())

	buf, err := json.Marshal(pidt)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `"user"`, string(buf))

	assert.Equal(t, PersonaIDTypeUser, ParsePersonaIDType("user"))
	assert.Equal(t, PersonaIDTypeUser, ParsePersonaIDType("USER"))

	assert.Equal(t, PersonaIDTypeGroup, ParsePersonaIDType("group"))
	assert.Equal(t, PersonaIDTypeGroup, ParsePersonaIDType("GROUP"))

	assert.Equal(t, PersonaIDTypeUID, ParsePersonaIDType("uid"))
	assert.Equal(t, PersonaIDTypeUID, ParsePersonaIDType("UID"))

	assert.Equal(t, PersonaIDTypeGID, ParsePersonaIDType("gid"))
	assert.Equal(t, PersonaIDTypeGID, ParsePersonaIDType("GID"))

	if err := json.Unmarshal([]byte(`"user"`), &pidt); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, PersonaIDTypeUser, pidt)
}

func TestPersonaIDMarshal(t *testing.T) {
	pid := &PersonaID{
		ID:   "akutz",
		Type: PersonaIDTypeUser,
	}

	buf, err := json.Marshal(pid)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, `"user:akutz"`, string(buf))
}

func TestOneExportListMarshal(t *testing.T) {
	testAllExportListMarshal(t, getOneExportJSON)
}

func TestAllExportListMarshal(t *testing.T) {
	testAllExportListMarshal(t, getAllExportsJSON)
}

func TestAllExportListMarshal2(t *testing.T) {
	testAllExportListMarshal(t, getAllExports2JSON)
}

func TestAllExportListMarshal3(t *testing.T) {
	testAllExportListMarshal(t, getAllExports3JSON)
}

func testAllExportListMarshal(t *testing.T, list []byte) {
	var exList ExportList
	if err := json.Unmarshal(list, &exList); err != nil {
		t.Fatal(err)
	}

	buf, err := json.Marshal(exList)
	if err != nil {
		t.Fatal(err)
	}

	map1 := map[string]interface{}{}
	if err := json.Unmarshal(buf, &map1); err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal(buf, &exList); err != nil {
		t.Fatal(err)
	}

	buf, err = json.Marshal(exList)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(buf))

	map2 := map[string]interface{}{}
	if err := json.Unmarshal(buf, &map2); err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(t, map1, map2)
}

func TestExportsList(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportsList(ctx, client)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportsList(ctx, client)
	assert.Equal(t, errors.New("error"), err)
}

func TestExportsListWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportsListWithZone(ctx, client, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportsListWithZone(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestExportInspect(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportInspect(ctx, client, 0)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportInspect(ctx, client, 0)
	assert.Equal(t, errors.New("error"), err)
}

func TestExportCreate(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	export := Export{
		ID: 0,
	}
	client.On("Post", anyArgs...).Return(nil).Once()
	_, err := ExportCreate(ctx, client, &export)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Post", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportCreate(ctx, client, &export)
	assert.Equal(t, errors.New("error"), err)
}

func TestExportCreateWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	export := Export{
		ID: 0,
	}

	_, err := ExportCreateWithZone(ctx, client, &export, "")
	assert.Equal(t, errors.New("zone cannot be empty"), err)

	client.On("Post", anyArgs...).Return(nil).Once()
	_, err = ExportCreateWithZone(ctx, client, &export, "zone")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
	client.On("Post", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportCreateWithZone(ctx, client, &export, "zone")
	assert.Equal(t, errors.New("error"), err)
}

func TestSetExportRootClients(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	client.On("Put", anyArgs...).Return(nil).Once()
	err := SetExportRootClients(ctx, client, 0, "addrs")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestExportUpdateWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	export := Export{
		ID: 0,
	}
	client.On("Put", anyArgs...).Return(nil).Once()
	err := ExportUpdateWithZone(ctx, client, &export, "", true)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestUnexport(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Once()
	err := Unexport(ctx, client, 0)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestUnexportWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Delete", anyArgs...).Return(nil).Once()
	err := UnexportWithZone(ctx, client, 0, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

func TestExportsListWithResume(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportsListWithResume(ctx, client, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportsListWithResume(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestExportsListWithLimit(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportsListWithLimit(ctx, client, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = ExportsListWithLimit(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestGetExportWithPath(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetExportWithPath(ctx, client, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}

	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = GetExportWithPath(ctx, client, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestGetExportWithPathAndZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetExportWithPathAndZone(ctx, client, "", "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = GetExportWithPathAndZone(ctx, client, "", "")
	assert.Equal(t, errors.New("error"), err)
}

func TestGetExportByIDWithZone(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}

	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := GetExportByIDWithZone(ctx, client, 0, "")
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
	client.On("Get", anyArgs...).Return(errors.New("error")).Once()
	_, err = GetExportByIDWithZone(ctx, client, 0, "")
	assert.Equal(t, errors.New("error"), err)
}

func TestExportsListWithParams(t *testing.T) {
	ctx := context.Background()
	client := &mocks.Client{}
	orderedValues := api.NewOrderedValues([][]string{
		{"detail", "owner", "group"},
		{"info", "?"},
	})
	client.On("Get", anyArgs...).Return(nil).Once()
	_, err := ExportsListWithParams(ctx, client, orderedValues)
	if err != nil {
		assert.Equal(t, "Test scenario failed", err)
	}
}

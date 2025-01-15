package v2

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePersonaType(t *testing.T) {
	personaType := ParsePersonaType("user")
	assert.Equal(t, PersonaTypeUser, personaType)

	personaType = ParsePersonaType("group")
	assert.Equal(t, PersonaTypeGroup, personaType)

	personaType = ParsePersonaType("wellknown")
	assert.Equal(t, PersonaTypeWellKnown, personaType)

	personaType = ParsePersonaType("unknown")
	assert.Equal(t, PersonaTypeUnknown, personaType)
}

func TestPersonaTypeString(t *testing.T) {
	personaType := PersonaTypeUser
	str := personaType.String()
	assert.Equal(t, "user", str)

	personaType = PersonaTypeGroup
	str = personaType.String()
	assert.Equal(t, "group", str)

	personaType = PersonaTypeWellKnown
	str = personaType.String()
	assert.Equal(t, "wellknown", str)

	personaType = PersonaTypeUnknown
	str = personaType.String()
	assert.Equal(t, "unknown", str)
}

func TestPersonaTypeUnmarshalJSON(t *testing.T) {
	var personaType PersonaType
	err := json.Unmarshal([]byte(`"user"`), &personaType)
	assert.NoError(t, err)
	assert.Equal(t, PersonaTypeUser, personaType)

	err = json.Unmarshal([]byte(`"group"`), &personaType)
	assert.NoError(t, err)
	assert.Equal(t, PersonaTypeGroup, personaType)

	err = json.Unmarshal([]byte(`"wellknown"`), &personaType)
	assert.NoError(t, err)
	assert.Equal(t, PersonaTypeWellKnown, personaType)

	err = json.Unmarshal([]byte(`"unknown"`), &personaType)
	assert.NoError(t, err)
	assert.Equal(t, PersonaTypeUnknown, personaType)

	err = json.Unmarshal([]byte(`"invalid"`), &personaType)
	assert.NoError(t, err)
	assert.Equal(t, PersonaTypeUnknown, personaType)
}

func TestPersonaTypeMarshalJSON(t *testing.T) {
	personaType := PersonaTypeUser
	data, err := personaType.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte(`"user"`), data)
}

func TestParsePersonaIDType(t *testing.T) {
	personaIDType := ParsePersonaIDType("SID")
	assert.Equal(t, PersonaIDTypeSID, personaIDType)

	personaIDType = ParsePersonaIDType("xyz")
	assert.Equal(t, PersonaIDTypeUnknown, personaIDType)
}

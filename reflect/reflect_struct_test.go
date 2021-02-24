package reflect_test

import (
	"encoding/json"
	"testing"

	"github.com/gothing/draft/reflect"
	"github.com/stretchr/testify/assert"
)

type StructWithJSONTag struct {
	Private     string `json:"-"`
	AccessToken string
	ClientID    string `json:"ClientID"`
}

type StructComposed struct {
	StructWithJSONTag
	Value  string
	hidden string
}

func TestStructWithJSONTag(t *testing.T) {
	v := reflect.Get(StructWithJSONTag{}, reflect.Options{
		SnakeCase: true,
	})

	assert.Equal(t, []string{"access_token", "ClientID"}, v.Keys())
}

type StructWithRequiredTags struct {
	RequiredField  int    `required:"true"`
	RequiredField1 string `required:"true,grp"`
	RequiredField2 uint64 `required:"true,grp"`
	OptionalField  uint   `required:"false"`
	ExtraField     bool
}

func TestStructWithRequiredTag(t *testing.T) {
	v := reflect.Get(
		StructWithRequiredTags{},
		reflect.Options{
			SnakeCase: true,
		},
	)

	jsonObj, err := json.Marshal(v)
	assert.NoError(t, err)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(jsonObj, &m))

	fields := []struct {
		Name          string
		Required      bool
		RequiredGroup string
	}{
		{
			Name:     "required_field",
			Required: true,
		},
		{
			Name:          "required_field1",
			Required:      true,
			RequiredGroup: "grp",
		},
		{
			Name:          "required_field2",
			Required:      true,
			RequiredGroup: "grp",
		},
		{
			Name: "optional_field",
		},
		{
			Name: "extra_field",
		},
	}

	for _, field := range fields {
		nested := m["nested"].([]interface{})
		for _, f := range nested {
			f := f.(map[string]interface{})
			if f["name"] == field.Name {
				assert.Equal(t, field.Required, f["required"])
				assert.Equal(t, field.RequiredGroup, f["required_group"])
			}
		}
	}
}

func xTestStructComposed(t *testing.T) {
	v := reflect.Get(StructComposed{}, reflect.Options{
		SnakeCase: true,
	})

	assert.Equal(t, []string{"access_token", "ClientID", "value"}, v.Keys())
}

package reflect_test

import (
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

func TestStructComposed(t *testing.T) {
	v := reflect.Get(StructComposed{}, reflect.Options{
		SnakeCase: true,
	})

	assert.Equal(t, []string{"access_token", "ClientID", "value"}, v.Keys())
}

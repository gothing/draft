package reflect_test

import (
	"testing"

	"github.com/gothing/draft/reflect"
	"github.com/stretchr/testify/assert"
)

type MyStringType string
type UserID uint64

func (v UserID) TypeDescription() string {
	return "Uniq user ID"
}

type UserObject struct {
	ID     UserID
	Exists bool `comment:"User Exists?"`
	Flags  UserFlags
	Detail *UserDetail
}

type UserFlags struct {
	Active  bool
	IsAdmin bool
}

type UserDetail struct {
	Name string
}

func (v *UserDetail) TypeDescription() string {
	return "User detail object"
}

func TestNil(t *testing.T) {
	v := reflect.Get(nil, reflect.Options{})
	assert.Equal(t, v.Type, "nil")
	assert.Equal(t, v.Type, v.MetaType)
	assert.Equal(t, v.Value, nil)
}

func TestString(t *testing.T) {
	v := reflect.Get("foo", reflect.Options{})
	assert.Equal(t, "string", v.Type)
	assert.Equal(t, v.Type, v.MetaType)
	assert.Equal(t, v.Value, "foo")
}

func TestTyped(t *testing.T) {
	v := reflect.Get(MyStringType("foo"), reflect.Options{})
	assert.Equal(t, "string", v.Type)
	assert.Equal(t, "MyStringType", v.MetaType)
}

func TestStruct(t *testing.T) {
	v := reflect.Get(UserObject{ID: 123}, reflect.Options{})

	assert.Equal(t, "struct", v.Type)
	assert.Equal(t, "UserObject", v.MetaType)
	assert.Equal(t, []string{"ID", "Exists", "Flags", "Detail"}, v.Keys())

	// ID
	assert.Equal(t, "ID", v.Nested[0].Name)
	assert.Equal(t, "Uniq user ID", v.Nested[0].Comment)
	assert.Equal(t, "uint64", v.Nested[0].Type)
	assert.Equal(t, "UserID", v.Nested[0].MetaType)

	// Exists
	assert.Equal(t, "User Exists?", v.Nested[1].Comment)

	// Flags
	assert.Equal(t, 2, len(v.Nested[2].Nested))
	assert.Equal(t, []string{"Active", "IsAdmin"}, v.Nested[2].Keys())

	// Detail
	assert.Equal(t, "User detail object", v.Nested[3].Comment)
	assert.Equal(t, []string{"Active", "IsAdmin"}, v.Nested[2].Keys())
}

func TestRef(t *testing.T) {
	x := &UserObject{ID: 123}

	v := reflect.Get(x, reflect.Options{})
	assert.Equal(t, "struct", v.Type, "ref")

	(func(v interface{}) {
		ref := reflect.Get(x, reflect.Options{})
		assert.Equal(t, "struct", ref.Type, "interface")
	})(x)
}

func TestInterface(t *testing.T) {
	var x interface{} = nil

	x = &UserObject{ID: 123}
	v := reflect.Get(x, reflect.Options{})
	assert.Equal(t, "struct", v.Type, "ref")

	(func(v interface{}) {
		ref := reflect.Get(x, reflect.Options{})
		assert.Equal(t, "struct", ref.Type, "interface")
	})(x)
}

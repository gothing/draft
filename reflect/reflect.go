package reflect

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

// Options -
type Options struct {
	SnakeCase bool
}

// Item -
type Item struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"-"`
	Type     string      `json:"type"`
	MetaType string      `json:"meta_type"`
	Tags     string      `json:"tags"`
	Comment  string      `json:"comment"`
	Required bool        `json:"required"`
	Nested   []Item      `json:"nested"`
}

// Get -
func Get(v interface{}, o Options) Item {
	item := Item{Type: "nil", MetaType: "nil"}
	if v == nil {
		return item
	}

	typeRef := reflect.TypeOf(v)
	valueRef := reflect.ValueOf(v)
	initItem(o, &item, typeRef, valueRef)

	return item
}

func initItem(
	o Options,
	item *Item,
	typeRef reflect.Type,
	valRef reflect.Value,
) {
	if reflect.Interface == typeRef.Kind() {
		valRef = valRef.Elem()
		typeRef = valRef.Type()
	}

	if item.Name == "" {
		if reflect.Ptr == typeRef.Kind() {
			typeRef = typeRef.Elem()
		}

		if valRef.Kind() == reflect.Ptr {
			valRef = valRef.Elem()
		}
	}

	item.Type = typeRef.Kind().String()
	item.MetaType = typeRef.Name()

	if valRef.IsValid() {
		item.Value = valRef.Interface()
	}

	switch typeRef.Kind() {
	case reflect.Struct:
		nested := make([]Item, typeRef.NumField())

		for i := 0; i < typeRef.NumField(); i++ {
			item := &nested[i]
			f := typeRef.Field(i)
			fv := valRef.FieldByName(f.Name)

			if o.SnakeCase {
				item.Name = strcase.ToSnake(f.Name)
			} else {
				item.Name = f.Name
			}

			item.Tags = string(f.Tag)
			item.Comment = f.Tag.Get("comment")
			item.Required = f.Tag.Get("required") == "true"

			td := fv.MethodByName("TypeDescription")
			if td.IsValid() {
				c := td.Call([]reflect.Value{})[0].String()
				if item.Comment == "" {
					item.Comment = c
				} else {
					item.Comment = strings.Replace(item.Comment, "{super}", c, -1)
				}
			}

			initItem(o, item, f.Type, fv)
		}

		item.Nested = nested
	}
}

// Keys -
func (item Item) Keys() []string {
	keys := make([]string, 0, len(item.Nested))
	for _, v := range item.Nested {
		if v.Name != "" {
			keys = append(keys, v.Name)
		}
	}

	return keys
}

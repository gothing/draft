package reflect

import (
	"reflect"
	"github.com/iancoleman/strcase"
)

// Options -
type Options struct {
	SnakeCase bool
}

// Item -
type Item struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	MetaType string `json:"meta_type"`
	Tags     string `json:"tags"`
	Comment  string `json:"comment"`
	Nested   []Item `json:"nested"`
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
	item.Type = typeRef.Kind().String()
	item.MetaType = typeRef.Name()

	switch typeRef.Kind() {
	case reflect.Struct:
		nested := make([]Item, typeRef.NumField())

		for i := 0; i < typeRef.NumField(); i++ {
			item := &nested[i]
			f := typeRef.Field(i)
			fv := valRef.FieldByName(f.Name)
			td := fv.MethodByName("TypeDescription")

			item.Name = f.Name
			if o.SnakeCase {
				item.Name = strcase.ToSnake(item.Name)
			}

			item.Tags = string(f.Tag)
			item.Comment = f.Tag.Get("comment")

			if td.IsValid() {
				item.Comment = td.Call([]reflect.Value{})[0].String()
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

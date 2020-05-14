package reflect

import (
	"reflect"
	"regexp"
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
	Enum     interface{} `json:"enum"`
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
	isNilVal := false

	// fmt.Println("-----", item.Name, "-----")
	// fmt.Println("1)", item, "->", typeRef.Kind(), "=", valRef.Kind(), "->", typeRef, "->", valRef)
	if typeRef.Kind() == reflect.Interface {
		valRef = valRef.Elem()
		typeRef = valRef.Type()
	}

	// fmt.Println("2)", item, "->", typeRef.Kind(), "=", valRef.Kind(), "->", typeRef, "->", valRef)
	if typeRef.Kind() == reflect.Ptr {
		typeRef = typeRef.Elem()
	}

	// fmt.Println("3)", item, "->", typeRef.Kind(), "=", valRef.Kind(), "->", typeRef, "->", valRef)
	if valRef.Kind() == reflect.Ptr {
		valRef = valRef.Elem()
		if !valRef.IsValid() {
			isNilVal = true
			valRef = reflect.Zero(typeRef)
		}
	}

	// fmt.Println("4)", item, "->", typeRef.Kind(), "=", valRef.Kind(), "->", typeRef, "->", valRef)
	item.Type = typeRef.Kind().String()
	item.MetaType = typeRef.Name()

	if valRef.IsValid() {
		if isNilVal {
			item.Value = nil
		} else {
			item.Value = valRef.Interface()
		}
	}

	switch typeRef.Kind() {
	case reflect.Slice:
		item.MetaType = typeRef.Elem().Name()
		tev := reflect.Zero(typeRef.Elem()).MethodByName("TypeEnumValues")
		if tev.IsValid() {
			item.Enum = tev.Call([]reflect.Value{})[0].Interface()
		}

	case reflect.Struct:
		item.Nested = initNested(o, typeRef, valRef)
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

var reIsPrivate = regexp.MustCompile(`^[a-z]`)

func initNested(o Options, typeRef reflect.Type, valRef reflect.Value) []Item {
	nested := make([]Item, 0, typeRef.NumField())

	for i := 0; i < typeRef.NumField(); i++ {
		item := &Item{}
		f := typeRef.Field(i)
		jsonTag := f.Tag.Get("json")

		if f.Anonymous {
			sub := initNested(o, f.Type, valRef.FieldByName(f.Name))
			nested = append(nested, sub...)
			continue
		}

		if jsonTag == "-" || reIsPrivate.MatchString(f.Name) {
			continue
		} else if jsonTag != "" {
			item.Name = jsonTag
		} else if o.SnakeCase {
			item.Name = strcase.ToSnake(f.Name)
		} else {
			item.Name = f.Name
		}

		nested = append(nested, *item)

		// fmt.Println(" -", item.Name)
		fv := valRef.FieldByName(f.Name)
		zfv := fv
		if f.Type.Kind() == reflect.Ptr {
			zfv = reflect.Zero(f.Type.Elem())
		}

		item.Tags = string(f.Tag)
		item.Comment = f.Tag.Get("comment")
		item.Required = f.Tag.Get("required") == "true"

		td := zfv.MethodByName("TypeDescription")
		if td.IsValid() {
			c := td.Call([]reflect.Value{})[0].String()
			if item.Comment == "" {
				item.Comment = c
			} else {
				item.Comment = strings.ReplaceAll(item.Comment, "{super}", c)
			}
		}

		tev := zfv.MethodByName("TypeEnumValues")
		if tev.IsValid() {
			item.Enum = tev.Call([]reflect.Value{})[0].Interface()
		}

		initItem(o, item, f.Type, fv)
	}

	return nested
}

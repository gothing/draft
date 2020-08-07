package dts

import (
	"fmt"
	"strings"

	"github.com/gothing/draft/reflect"

	"github.com/gothing/draft"
	"github.com/iancoleman/strcase"
)

type dtsObjectType map[string]interface{}

// Gen -
func Gen(s *draft.Scheme) {
	// scheme := s.ToJSON()
	code := make([]string, 0)

	code = append(code, dtsGenAccess()...)
}

func dtsGenAccess() []string {
	code := make([]string, 0)

	for _, r := range draft.GetDocConfig().Rights {
		name := strcase.ToCamel(`EndpointAccess_` + string(r.ID))
		code = append(code, `export type `+name+` = {`)
		for _, e := range r.Extra {
			code = append(code, fmt.Sprintf("\t%q:", e.Name))
			code = append(code, "\t\tparams: "+dtsGenObjectType(e.Params, "\t\t\t"))
			code = append(code, "\t};")
		}
		code = append(code, `};`)
	}

	return code
}

func dtsGenObjectType(v interface{}, ind string) string {
	ref := reflect.Get(v, reflect.Options{SnakeCase: true})
	code := []string{"{"}

	for _, v := range ref.Nested {
		code = append(code, ind+"// "+v.Comment)
		if len(v.Nested) > 0 {

		} else {
			code = append(code, ind+v.Name+": "+dtsGenJSType(v))
		}
	}

	return strings.Join(code, "\n") + "\n" + ind + "};"
}

func dtsGenJSType(v reflect.Item) string {
	switch v.Type {
	case "int", "int32", "int64", "uint", "uint32", "uint64", "float":
		return "number"

	case "bool":
		return "boolean"

	case "string":
		return "string"
	}

	return "any"
}

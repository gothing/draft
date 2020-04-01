package draft

// MethodType -
type MethodType string

// Method -
var Method = struct {
	GET     MethodType
	POST    MethodType
	PUT     MethodType
	PATCH   MethodType
	DELETE  MethodType
	HEAD    MethodType
	OPTIONS MethodType
}{
	"GET",
	"POST",
	"PUT",
	"PATCH",
	"DELETE",
	"HEAD",
	"OPTIONS",
}

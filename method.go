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

// MimeType -
type MimeType string

// Mime -
var Mime = struct {
	Any                MimeType
	JSON               MimeType
	FormData           MimeType
	XWWWFormUrlencoded MimeType
}{
	"",
	"application/json",
	"multipart/form-data",
	"application/x-www-form-urlencoded",
}

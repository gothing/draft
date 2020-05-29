package draft

// AccessType -
type AccessType string

// Access — парва доступка к апишке
var Access = struct {
	All          AccessType
	Auth         AccessType
	AuthOptional AccessType
	Cookie       AccessType
	AuthCookie   AccessType
	Signature    AccessType
}{
	"all",
	"auth",
	"auth:optional",
	"cookie",
	"auth:cookie",
	"signature",
}

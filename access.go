package draft

// AccessType -
type AccessType string

// Access — парва доступка к апишке
var Access = struct {
	All          AccessType
	Auth         AccessType
	Anonymous    AccessType
	AuthOptional AccessType
	Cookie       AccessType
	AuthCookie   AccessType
	Signature    AccessType
	BasicAuth    AccessType
	OAuth2       AccessType
	Token        AccessType
	Session      AccessType
}{
	"all",
	"auth",
	"anonymous",
	"auth:optional",
	"cookie",
	"auth:cookie",
	"signature",
	"basic:auth",
	"oauth2",
	"token",
	"session",
}

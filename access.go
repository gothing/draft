package draft

// AccessType -
type AccessType string

// Access — парва доступка к апишке
var Access = struct {
	All          AccessType
	Auth         AccessType
	AuthOptional AccessType
	SWASign      AccessType
}{
	"all",
	"auth",
	"auth:optional",
	"swa:sign",
}

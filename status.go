package draft

// StatusType -
type StatusType string

// Status — возможные статусы ответа
var Status = struct {
	Processing                 StatusType
	OK                         StatusType
	Accepted                   StatusType
	NonAuthoritative           StatusType
	Partial                    StatusType
	Move                       StatusType
	Found                      StatusType
	NotModified                StatusType
	Invalid                    StatusType
	PaymentRequired            StatusType
	Denied                     StatusType
	NotFound                   StatusType
	Unacceptable               StatusType
	Timeout                    StatusType
	Conflict                   StatusType
	ExpectationFailed          StatusType
	Unprocessable              StatusType
	FailedDependency           StatusType
	Locked                     StatusType
	UpgradeRequired            StatusType
	ManyRequests               StatusType
	RetryWith                  StatusType
	UnavailableForLegalReasons StatusType
	Fail                       StatusType
	NotImplemented             StatusType
	Unavaliable                StatusType
	Insufficient               StatusType
}{
	"processing",
	"ok",
	"accepted",
	"non_authoritative",
	"partial",
	"move",
	"found",
	"notmodified",
	"invalid",
	"payment_required",
	"denied",
	"notfound",
	"unacceptable",
	"timeout",
	"conflict",
	"expectation_failed",
	"unprocessable",
	"failed_dependency",
	"locked",
	"upgrade_required",
	"many_requests",
	"retry_with",
	"unavailable_for_legal_reasons",
	"fail",
	"not_implemented",
	"unavaliable",
	"insufficient",
}

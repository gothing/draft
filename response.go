package draft

// Response -
type Response struct {
	Status StatusType
	Body   interface{}
}

// NewResponse -
func NewResponse(s StatusType, b interface{}) *Response {
	return &Response{
		Status: s,
		Body:   b,
	}
}

// NewInvalidResponseByMissedParams -
func NewInvalidResponseByMissedParams(list []string) *Response {
	body := make(map[string]map[string]string)
	for _, n := range list {
		body[n] = map[string]string{
			"error": "required",
			"value": "",
		}
	}

	return NewResponse(Status.Invalid, body)
}

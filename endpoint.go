package draft

import (
	"encoding/json"
	"net/http"

	"github.com/gothing/draft/reflect"
	"github.com/iancoleman/orderedmap"
)

// endpoint — интерфейс «конца»
type endpoint interface {
	GetEndpointMock(r *Request) interface{}
	ValidateEndpointMockRequest(r *Request) *Response
	InitEndpointScheme(s *Scheme)
}

// Endpoint — базовые методы «конца»
type Endpoint struct {
	endpointCtrl   endpoint
	endpointScheme *Scheme
}

// Init -
func (e *Endpoint) Init(ctrl endpoint) {
	if e.endpointScheme == nil {
		scheme := new(Scheme)
		e.endpointCtrl = ctrl
		e.endpointScheme = scheme

		ctrl.InitEndpointScheme(scheme)
	}
}

// GetScheme -
func (e *Endpoint) GetScheme() *Scheme {
	return e.endpointScheme
}

// ValidateEndpointMockRequest -
func (e Endpoint) ValidateEndpointMockRequest(r *Request) *Response {
	return nil
}

// GetEndpointMock -
func (e Endpoint) GetEndpointMock(r *Request) interface{} {
	cases := e.endpointScheme.Cases()
	for _, c := range cases {
		if Status.OK == c.Status {
			m := r.Params.GetMissed(c.Params)
			if len(m) > 0 {
				return NewInvalidResponseByMissedParams(m)
			}

			return NewResponse(Status.OK, c.Body)
		}
	}
	return nil
}

// ServeHTTP -
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	s := e.endpointScheme
	if s == nil {
		http.Error(w, "Endpoint not initialized", http.StatusInternalServerError)
		return
	}

	var resp interface{}
	req := ParseRequest(r)
	resp = e.endpointCtrl.ValidateEndpointMockRequest(req)
	if resp.(*Response) == nil {
		resp = e.endpointCtrl.GetEndpointMock(req)
	}

	bytes, err := json.Marshal(MockToResponse(resp))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Intercept -
type Intercept struct {
	Body interface{}
}

// MockToResponse -
func MockToResponse(v interface{}) interface{} {
	ref := reflect.Get(v, reflect.Options{
		SnakeCase: true,
	})

	return prepareMock(ref)
}

func prepareMock(item reflect.Item) interface{} {
	switch item.Type {
	case "struct":
		obj := orderedmap.New()
		for _, entry := range item.Nested {
			obj.Set(entry.Name, prepareMock(entry))
		}
		return obj

	default:
		return item.Value
	}
}

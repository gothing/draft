package draft

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gothing/draft/reflect"
	"github.com/iancoleman/orderedmap"
)

// EndpointAPI — интерфейс «конца»
type EndpointAPI interface {
	http.Handler
	URL() string
	GetHandler() http.HandlerFunc
	EndpointServeHTTP(http.ResponseWriter, *http.Request)
	InitEndpoint(ctrl EndpointAPI)
	InitEndpointScheme(s *Scheme)
	GetEndpointMock(r *Request) interface{}
	EndpointHandle(r *http.Request) ([]byte, error)
	ValidateEndpointMockRequest(r *Request) *Response
	GetScheme() *Scheme
}

// Endpoint — базовые методы «конца»
type Endpoint struct {
	Handler        http.HandlerFunc
	endpointCtrl   EndpointAPI
	endpointScheme *Scheme
}

// URL -
func (e *Endpoint) URL() string {
	return e.GetScheme().url
}

// GetHandler -
func (e *Endpoint) GetHandler() http.HandlerFunc {
	return e.Handler
}

// InitEndpoint -
func (e *Endpoint) InitEndpoint(ctrl EndpointAPI) {
	if e.endpointScheme == nil {
		scheme := &Scheme{
			defAccess: Access.All,
			defMethod: Method.POST,
		}
		e.endpointCtrl = ctrl
		e.endpointScheme = scheme

		ctrl.InitEndpointScheme(scheme)
	}
}

// InitEndpointScheme -
func (e *Endpoint) InitEndpointScheme(s *Scheme) {
	log.Println("InitEndpointScheme — not implemented")
}

// GetScheme -
func (e *Endpoint) GetScheme() *Scheme {
	return e.endpointScheme
}

// ValidateEndpointMockRequest -
func (e *Endpoint) ValidateEndpointMockRequest(r *Request) *Response {
	return nil
}

// GetEndpointMock -
func (e *Endpoint) GetEndpointMock(r *Request) interface{} {
	idx := -1
	weight := -1
	cases := e.endpointScheme.Cases()
	missed := []string{}

	for i, c := range cases {
		if Status.OK == c.Status {
			w := 0
			ref := reflect.Get(c.Params, reflect.Options{SnakeCase: true})
			m := make([]string, 0, len(ref.Nested))

			for _, item := range ref.Nested {
				if !r.Params.Has(item.Name) && item.Required {
					m = append(m, item.Name)
				}

				if r.Params.Get(item.Name) == fmt.Sprintf("%v", item.Value) {
					w++
				}
			}

			if len(m) > 0 && len(missed) == 0 {
				missed = m
			}

			if len(m) == 0 && w > weight {
				idx = i
				weight = w
			}
		}
	}

	if idx != -1 {
		c := cases[idx]
		return NewResponse(c.Status, c.Body)
	}

	if len(missed) > 0 {
		return NewInvalidResponseByMissedParams(missed)
	}

	return nil
}

// ServeHTTP -
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler := e.GetHandler()
	if handler == nil {
		e.EndpointServeHTTP(w, r)
	} else {
		handler(w, r)
	}
}

// EndpointServeHTTP -
func (e *Endpoint) EndpointServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("GODRAFT-Time", time.Now().Format(time.RFC3339))

	bytes, err := e.endpointCtrl.EndpointHandle(r)
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

// EndpointHandle -
func (e *Endpoint) EndpointHandle(r *http.Request) ([]byte, error) {
	s := e.endpointScheme
	if s == nil {
		return nil, fmt.Errorf("Endpoint not initialized")
	}

	req := ParseRequest(r)
	if strings.Contains(req.Path, "/godraft:scheme/") {
		return json.Marshal(s.ToJSON())
	}

	var resp interface{}
	resp = e.endpointCtrl.ValidateEndpointMockRequest(req)
	if resp.(*Response) == nil {
		resp = e.endpointCtrl.GetEndpointMock(req)
	}

	bytes, err := json.Marshal(MockToResponse(resp))
	if err != nil {
		return nil, err
	}

	return bytes, nil
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
		if item.Value == nil {
			return nil
		}

		obj := orderedmap.New()
		for _, entry := range item.Nested {
			obj.Set(entry.Name, prepareMock(entry))
		}
		return obj

	default:
		return item.Value
	}
}

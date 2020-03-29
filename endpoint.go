package draft

import (
	"encoding/json"
	"net/http"
)

// endpoint — интерфейс «конца»
type endpoint interface {
	InitEndpoint() error
	InitEndpointScheme(s *Scheme)
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Endpoint — базовые методы «конца»
type Endpoint struct {
	endpoint
	endpointScheme *Scheme
}

// getEndpointScheme -
func (e *Endpoint) getEndpointScheme() *Scheme {
	if e.endpointScheme == nil {
		scheme := &Scheme{}
		e.InitEndpointScheme(scheme)
		e.endpointScheme = scheme
	}

	return e.endpointScheme
}

// InitEndpoint -
func (e *Endpoint) InitEndpoint() (*Intercept, error) {
	_ = e.getEndpointScheme()

	return nil, nil
}

// ServeHTTP -
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	s := e.getEndpointScheme()
	bytes, err := json.Marshal(s.GetCaseByStatus(Status.OK))
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

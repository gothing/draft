package draft

import (
	"encoding/json"
	"net/http"
)

// endpoint — интерфейс «конца»
type endpoint interface {
	InitEndpointScheme(s *Scheme)
}

// Endpoint — базовые методы «конца»
type Endpoint struct {
	endpointScheme *Scheme
}

// Init -
func (e *Endpoint) Init(api endpoint) {
	if e.endpointScheme == nil {
		scheme := &Scheme{}
		api.InitEndpointScheme(scheme)
		e.endpointScheme = scheme
	}
}

// GetScheme -
func (e *Endpoint) GetScheme() *Scheme {
	return e.endpointScheme
}

// ServeHTTP -
func (e *Endpoint) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	s := e.endpointScheme
	if s == nil {
		http.Error(w, "Endpoint not initialized", http.StatusInternalServerError)
		return
	}

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

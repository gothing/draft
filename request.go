package draft

import (
	"net/http"
	"net/url"

	"github.com/gothing/draft/reflect"
)

// Request -
type Request struct {
	req    *http.Request
	Path   string
	Params *RequestParams
}

// RequestParams -
type RequestParams struct {
	query url.Values
	form  url.Values
}

// ParseRequest -
func ParseRequest(r *http.Request) *Request {
	p := &RequestParams{
		query: r.URL.Query(),
		form:  nil,
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err == nil {
			p.form = r.Form
		}
	}

	return &Request{
		req:    r,
		Path:   r.URL.Path,
		Params: p,
	}
}

// Has -
func (p *RequestParams) Has(key string) bool {
	if p.form != nil {
		if _, ok := p.form[key]; ok {
			return ok
		}
	}

	_, ok := p.query[key]
	return ok
}

// Get -
func (p *RequestParams) Get(key string) string {
	if p.form != nil {
		if v, ok := p.form[key]; ok {
			return v[0]
		}
	}

	return p.query.Get(key)
}

// GetMissed -
func (p *RequestParams) GetMissed(v interface{}) []string {
	ref := reflect.Get(v, reflect.Options{SnakeCase: true})
	missed := make([]string, 0, len(ref.Nested))

	for _, item := range ref.Nested {
		if !p.Has(item.Name) {
			missed = append(missed, item.Name)
		}
	}

	return missed
}

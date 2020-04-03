package draft

import (
	"fmt"
	"net/http"
	"strings"
)

// APIService -
type APIService struct {
	http.Handler
	routes map[string]endpoint
}

// Compose -
func Compose(endpoints ...endpoint) *APIService {
	routes := make(map[string]endpoint)
	for _, ctrl := range endpoints {
		ctrl.InitEndpoint(ctrl)
		s := ctrl.GetScheme()
		routes[s.url] = ctrl
	}

	return &APIService{
		routes: routes,
	}
}

// ServeHTTP -
func (api *APIService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.Contains(path, "/godraft:doc") {
		RenderDOC(api, w, r)
		return
	}

	path = strings.Replace(path, "/godraft/", "/", 1)
	path = strings.Replace(path, "/godraft:scheme/", "/", 1)

	ctrl, exists := api.routes[path]
	if exists {
		if path != r.URL.Path {
			ctrl.EndpointServeHTTP(w, r)
		} else {
			ctrl.ServeHTTP(w, r)
		}
	} else {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(fmt.Sprintf("'%s' not found", path)))
	}
}

// URLs -
func (api *APIService) URLs() []string {
	list := make([]string, 0, len(api.routes))
	for u := range api.routes {
		list = append(list, u)
	}
	return list
}

// Handle -
func (api *APIService) Handle(pattern string, handler http.Handler) {
	if handler == nil {
		http.Handle(pattern, api)
	} else {
		http.Handle(pattern, handler)
	}

	http.Handle("/godraft"+pattern, api)
	http.Handle("/godraft:scheme"+pattern, api)
}

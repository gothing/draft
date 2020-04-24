package draft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// APIService -
type APIService struct {
	http.Handler
	routes      map[string]endpoint
	rootGroup   *apiGroupEntry
	activeGroup *apiGroupEntry
}

type apiGroupEntry struct {
	Type        string           `json:"type"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Scheme      *JSONScheme      `json:"scheme"`
	Entries     []*apiGroupEntry `json:"entries"`
}

// Compose -
func Compose(endpoints ...endpoint) *APIService {
	routes := make(map[string]endpoint)
	for _, ctrl := range endpoints {
		ctrl.InitEndpoint(ctrl)
		s := ctrl.GetScheme()
		routes[s.url] = ctrl
	}

	root := createGroupEntry("G", "#root", "")
	srv := &APIService{
		rootGroup:   root,
		activeGroup: root,
		routes:      routes,
	}

	http.Handle("/godraft:scheme/", srv)

	return srv
}

func (api *APIService) getDodraftScheme() *apiGroupEntry {
	if len(api.rootGroup.Entries) == 0 {
		for p := range api.routes {
			api.rootGroup.Entries = append(api.rootGroup.Entries, createGroupEntry("E", p, ""))
		}
	}

	return api.rootGroup.init(api)
}

// ServeHTTP -
func (api *APIService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.Contains(path, "/godraft:doc") {
		RenderDOC(api, w, r)
		return
	}

	if path == "/godraft:scheme/" {
		if len(api.rootGroup.Entries) == 0 {
			for p := range api.routes {
				api.rootGroup.Entries = append(api.rootGroup.Entries, createGroupEntry("E", p, ""))
			}
		}

		result, _ := json.Marshal(api.getDodraftScheme())
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(result)
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

func createGroupEntry(t, name, description string) *apiGroupEntry {
	return &apiGroupEntry{
		Type:        t,
		Name:        name,
		Description: description,
		Entries:     make([]*apiGroupEntry, 0),
	}
}

// Group -
func (api *APIService) Group(name, description string, executer func()) {
	pg := api.activeGroup
	api.activeGroup = createGroupEntry("G", name, description)

	executer()

	api.activeGroup = pg
}

// GroupHR -
func (api *APIService) GroupHR() {
	api.activeGroup.Entries = append(api.activeGroup.Entries, createGroupEntry("HR", "", ""))
}

// Handle -
func (api *APIService) Handle(pattern string, handler http.Handler) {
	api.activeGroup.Entries = append(api.activeGroup.Entries, createGroupEntry("E", pattern, ""))

	if handler == nil {
		http.Handle(pattern, api)
	} else {
		http.Handle(pattern, handler)
	}

	http.Handle("/godraft"+pattern, api)
	http.Handle("/godraft:scheme"+pattern, api)
}

func (e *apiGroupEntry) init(api *APIService) *apiGroupEntry {
	if e.Scheme == nil && e.Type == "E" {
		if ctrl, ok := api.routes[e.Name]; ok {
			s := ctrl.GetScheme().ToJSON()
			e.Scheme = &s
		}
	}

	for _, ne := range e.Entries {
		ne.init(api)
	}

	return e
}

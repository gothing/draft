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

func (api *APIService) getGodraftScheme() *apiGroupEntry {
	if len(api.rootGroup.Entries) == 0 {
		for p := range api.routes {
			api.rootGroup.Entries = append(api.rootGroup.Entries, createGroupEntry("E", p, ""))
		}
	}

	return api.rootGroup.init(api)
}

// Group -
type Group struct {
	name  string
	items []endpoint
}

// Compose -
func Compose(name string, items ...endpoint) Group {
	return Group{name, items}
}

// Add -
func (api *APIService) Add(g Group) {
	api.Group(g.name, "", func() {
		for _, item := range g.items {
			api.Handle(item, item.GetHandler())
		}
	})
}

// ListenAndServe -
func (api *APIService) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, api)
}

// ServeHTTP -
func (api *APIService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if strings.Contains(path, "/godraft:doc") {
		RenderDOC(api, w, r)
		return
	}

	if path == "/godraft:scheme/" {
		result, _ := json.Marshal(api.getGodraftScheme())
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(result)
		return
	}

	path = strings.Replace(path, "/godraft/", "/", 1)
	path = strings.Replace(path, "/godraft:scheme/", "/", 1)

	ctrl, exists := api.routes[path]
	if exists {
		if path != r.URL.Path || ctrl.GetHandler() == nil {
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
	parent := api.activeGroup
	api.activeGroup = createGroupEntry("G", name, description)

	parent.Entries = append(parent.Entries, api.activeGroup)
	executer()

	api.activeGroup = parent
}

// GroupHR -
func (api *APIService) GroupHR() {
	api.activeGroup.Entries = append(api.activeGroup.Entries, createGroupEntry("HR", "", ""))
}

// Handle -
func (api *APIService) Handle(endpoint endpoint, handler http.Handler) {
	endpoint.InitEndpoint(endpoint)

	scheme := endpoint.GetScheme()
	pattern := scheme.url

	api.activeGroup.Entries = append(api.activeGroup.Entries, createGroupEntry("E", pattern, ""))
	api.routes[pattern] = endpoint

	if handler == nil {
		http.Handle(pattern, api)
	} else {
		http.Handle(pattern, handler)
	}

	http.Handle("/godraft"+pattern, api)
	http.Handle("/godraft:scheme"+pattern, api)
}

// Create -
func Create() *APIService {
	root := createGroupEntry("G", "#root", "")
	srv := &APIService{
		rootGroup:   root,
		activeGroup: root,
		routes:      make(map[string]endpoint),
	}
	return srv
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

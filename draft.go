package draft

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// APIService -
type APIService struct {
	http.Handler
	config         Config
	routes         map[string]apiServiceRoute
	rootGroup      *apiGroupEntry
	activeGroup    *apiGroupEntry
	endpointClient *http.Client
}

type apiServiceRoute struct {
	handler http.Handler
	ctrl    EndpointAPI
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

// GroupHandler -
type GroupHandler interface {
	http.Handler
	Routes() []string
}

// Group -
type Group struct {
	name  string
	items []EndpointAPI
}

// Compose -
func Compose(name string, items ...EndpointAPI) Group {
	return Group{name, items}
}

// Add -
func (api *APIService) Add(g Group, groupHandlers ...GroupHandler) {
	api.Group(g.name, "", func() {
		for _, item := range g.items {
			item.InitEndpoint(item)
			api.Handle(
				item,
				findGroupHandler(item.URL(), groupHandlers, item.GetHandler()),
			)
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

	if api.config.DevMode {
		if strings.Contains(path, "/godraft:request/") {
			doDraftRequest(api, w, r)
			return
		}

		if strings.Contains(path, "/godraft:doc") {
			RenderDOC(api, w, r)
			return
		}

		if path == "/godraft:scheme/" {
			result, err := json.Marshal(api.getGodraftScheme())
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Content-Type", "application/json; charset=utf-8")

			if err != nil {
				w.Header().Set("X-GODraft-Scheme-Marshal-Error", err.Error())
				w.Write([]byte(fmt.Sprintf(`{"error":%q}`, err.Error())))
				return
			}

			w.Write(result)
			return
		}

		path = strings.Replace(path, "/godraft/", "/", 1)
		path = strings.Replace(path, "/godraft:scheme/", "/", 1)
	}

	route, exists := api.routes[path]
	if exists {
		if path != r.URL.Path || !isHTTPHandler(route.handler) {
			if api.config.DevMode {
				route.ctrl.ServeHTTP(w, r)
				return
			}
		} else {
			route.handler.ServeHTTP(w, r)
			return
		}
	}

	w.WriteHeader(404)
	_, _ = w.Write([]byte(fmt.Sprintf("'%s' not found", path)))
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
func (api *APIService) Handle(endpoint EndpointAPI, handler http.Handler) {
	endpoint.InitEndpoint(endpoint)
	pattern := endpoint.URL()

	api.activeGroup.Entries = append(api.activeGroup.Entries, createGroupEntry("E", pattern, ""))
	api.routes[pattern] = apiServiceRoute{
		handler: handler,
		ctrl:    endpoint,
	}

	if api.config.MockMode == MockEnable {
		if isHTTPHandler(handler) {
			http.Handle(pattern, handler)
		} else if api.config.DevMode {
			http.Handle(pattern, api)
		}
	}

	if api.config.DevMode {
		http.Handle("/godraft"+pattern, api)
		http.Handle("/godraft:doc"+pattern, api)
		http.Handle("/godraft:docs"+pattern, api)
		http.Handle("/godraft:scheme"+pattern, api)
	}
}

// ugly!
var draftHandled = false

// Config -
type Config struct {
	DevMode      bool
	ClientConfig ClientConfig
	MockMode     MockMode
}

// ClientConfig -
type ClientConfig struct {
	RequestTimeout time.Duration
	SkipVerifyCert bool
}

// MockMode -
type MockMode int

const (
	MockEnable MockMode = iota
	MockDisable
)

// Create -
func Create(cfg Config) *APIService {
	root := createGroupEntry("G", "#root", "")
	srv := &APIService{
		config:      cfg,
		rootGroup:   root,
		activeGroup: root,
		routes:      make(map[string]apiServiceRoute),
		endpointClient: &http.Client{
			Timeout: cfg.ClientConfig.RequestTimeout,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: cfg.ClientConfig.SkipVerifyCert,
				},
			},
		},
	}

	if cfg.DevMode && !draftHandled {
		draftHandled = true
		http.Handle("/godraft:doc/", srv)
		http.Handle("/godraft:docs/", srv)
		http.Handle("/godraft:scheme/", srv)
		http.Handle("/godraft:request/", srv)
	}

	return srv
}

func (e *apiGroupEntry) init(api *APIService) *apiGroupEntry {
	if e.Scheme == nil && e.Type == "E" {
		if route, ok := api.routes[e.Name]; ok {
			s := route.ctrl.GetScheme().ToJSON()
			e.Scheme = &s
		}
	}

	for _, ne := range e.Entries {
		ne.init(api)
	}

	return e
}

func isHTTPHandler(handler http.Handler) bool {
	if handler != nil {
		if v, ok := handler.(http.HandlerFunc); ok {
			return v != nil
		}
		if v, ok := handler.(http.Handler); ok {
			return v != nil
		}
	}

	return false
}

func findGroupHandler(u string, list []GroupHandler, def http.HandlerFunc) http.HandlerFunc {
	for _, gh := range list {
		if gh == nil {
			continue
		}

		for _, gu := range gh.Routes() {
			if gu == u {
				return gh.ServeHTTP
			}
		}
	}

	return def
}

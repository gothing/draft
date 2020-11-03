package draft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/gothing/draft/reflect"
)

var (
	reGodraftConfig  = regexp.MustCompile(`window\.__GODRAFT_CONFIG__`)
	reGodraftFixSrc  = regexp.MustCompile(` src="\/`)
	reGodraftFixHref = regexp.MustCompile(` href="\/`)

	pureDocConfig = DocConfig{}
)

// DocConfig -
type DocConfig struct {
	FrontURL       string         `json:"front_url"`
	ActiveGroup    string         `json:"active_group"`
	Groups         []DocGroup     `json:"groups"`
	Projects       []DocProject   `json:"projects"`
	Rights         []DocAccess    `json:"rights"`
	RequestFactory RequestFactory `json:"-"`
}

// DocGroup -
type DocGroup struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Entries []string `json:"entries"`
}

// DocProject -
type DocProject struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Descrition string `json:"descrition"`
	Host       string `json:"host"`
	HostRC     string `json:"host_rc"`
	HostDEV    string `json:"host_dev"`
}

// DocAccess -
type DocAccess struct {
	ID          AccessType       `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Badge       string           `json:"badge"`
	Extra       []DocAccessExtra `json:"extra"`
	ReqPrepare  RequestPrepare   `json:"-"`
}

// DocAccessExtra -
type DocAccessExtra struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Headers     interface{}    `json:"headers"`
	Cookies     interface{}    `json:"cookies"`
	Params      interface{}    `json:"params"`
	ReqPrepare  RequestPrepare `json:"-"`
}

type docFrontConfig struct {
	DocConfig
	Rights []docFrontAccess `json:"rights"`
}

type docFrontAccess struct {
	DocAccess
	Extra []docFrontAccessExtra `json:"extra"`
}

type docFrontAccessExtra struct {
	DocAccessExtra
	Headers docFrontAccessExtraReflectItem `json:"headers"`
	Cookies docFrontAccessExtraReflectItem `json:"cookies"`
	Params  docFrontAccessExtraReflectItem `json:"params"`
}

type docFrontAccessExtraReflectItem struct {
	Value   interface{}  `json:"value"`
	Reflect reflect.Item `json:"reflect"`
}

// RenderDOC -
func RenderDOC(api *APIService, w http.ResponseWriter, r *http.Request) {
	frontURL := pureDocConfig.FrontURL

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if frontURL == "" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`<h1 style="color: red">GODARFT :: DOC <— not configured.</h1>`))
		return
	}

	resp, err := http.Get(pureDocConfig.FrontURL)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	scheme := "http:"
	if r.TLS != nil {
		scheme = "https:"
	}

	if len(pureDocConfig.Groups) == 0 {
		id := strings.ReplaceAll(r.Host, ":", "-")
		url := fmt.Sprintf("%s/%s/godraft:scheme/", scheme, r.Host)
		pureDocConfig.ActiveGroup = id
		pureDocConfig.Groups = []DocGroup{
			{
				ID:      id,
				Name:    r.Host,
				Entries: []string{url},
			},
		}
	}

	frontConfig := docFrontConfig{
		DocConfig: pureDocConfig,
		Rights:    preapreFrontRights(pureDocConfig.Rights),
	}

	config, _ := json.MarshalIndent(frontConfig, "", " ")

	body = reGodraftConfig.ReplaceAll(body, config)
	body = reGodraftFixSrc.ReplaceAll(body, []byte(fmt.Sprintf(` src="%s`, frontURL)))
	body = reGodraftFixHref.ReplaceAll(body, []byte(fmt.Sprintf(` href="%s`, frontURL)))
	w.Write(body)
}

// SetupDoc -
func SetupDoc(c DocConfig) {
	pureDocConfig = c
}

// GetDocConfig -
func GetDocConfig() DocConfig {
	return pureDocConfig
}

func preapreFrontRights(rights []DocAccess) []docFrontAccess {
	if rights == nil {
		return make([]docFrontAccess, 0)
	}

	list := make([]docFrontAccess, 0, len(rights))

	for _, a := range rights {
		list = append(list, docFrontAccess{
			DocAccess: a,
			Extra:     prepareFrontAccessExtra(a.Extra),
		})
	}

	return list
}

func prepareFrontAccessExtra(extra []DocAccessExtra) []docFrontAccessExtra {
	if extra == nil {
		return make([]docFrontAccessExtra, 0)
	}

	opts := reflect.Options{SnakeCase: true}
	list := make([]docFrontAccessExtra, 0, len(extra))

	for _, e := range extra {
		list = append(list, docFrontAccessExtra{
			DocAccessExtra: e,

			Headers: docFrontAccessExtraReflectItem{
				Value:   e.Headers,
				Reflect: reflect.Get(e.Headers, opts),
			},

			Cookies: docFrontAccessExtraReflectItem{
				Value:   e.Cookies,
				Reflect: reflect.Get(e.Cookies, opts),
			},

			Params: docFrontAccessExtraReflectItem{
				Value:   e.Params,
				Reflect: reflect.Get(e.Params, opts),
			},
		})
	}

	return list
}

func findDocProject(id string) *DocProject {
	for _, v := range pureDocConfig.Projects {
		if v.ID == id {
			return &v
		}
	}

	return nil
}

func findDocAccess(id AccessType) *DocAccess {
	for _, v := range pureDocConfig.Rights {
		if v.ID == id {
			return &v
		}
	}

	return nil
}

package draft

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

var (
	reGodraftConfig  = regexp.MustCompile(`window\.__GODRAFT_CONFIG__`)
	reGodraftFixSrc  = regexp.MustCompile(` src="\/`)
	reGodraftFixHref = regexp.MustCompile(` href="\/`)

	frontServer = ""
)

type docFrontProject struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Entries []string `json:"entries"`
}

type docFrontCfg struct {
	ActiveProject string            `json:"activeProject"`
	Sitemap       []docFrontProject `json:"sitemap"`
}

// RenderDOC -
func RenderDOC(api *APIService, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if frontServer == "" {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`<h1 style="color: red">GODARFT :: DOC <— not configured.</h1>`))
		return
	}

	resp, err := http.Get(frontServer)
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

	id := strings.ReplaceAll(r.Host, ":", "-")
	url := fmt.Sprintf("%s//%s/godraft:scheme/", scheme, r.Host)
	config, _ := json.Marshal(docFrontCfg{
		id,
		[]docFrontProject{
			{id, id, []string{url}},
		},
	})

	body = reGodraftConfig.ReplaceAll(body, config)
	body = reGodraftFixSrc.ReplaceAll(body, []byte(fmt.Sprintf(` src="%s`, frontServer)))
	body = reGodraftFixHref.ReplaceAll(body, []byte(fmt.Sprintf(` href="%s`, frontServer)))
	w.Write(body)
}

// SetDocServer -
func SetDocServer(u string) {
	frontServer = u
}

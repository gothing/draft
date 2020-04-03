package draft

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/gothing/draft/reflect"
	"github.com/iancoleman/orderedmap"
	"github.com/russross/blackfriday"
)

var (
	reRemoveDocSlug = regexp.MustCompile(`^/godraft:doc(.*?)/+`)
	reFirstSpaces   = regexp.MustCompile(`^[\s\t]+`)
)

// RenderDOC -
func RenderDOC(api *APIService, w http.ResponseWriter, r *http.Request) {
	path := reRemoveDocSlug.ReplaceAllString(r.URL.Path, "/")

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	name := r.Host
	input := []string{"# [GODRAFT](/godraft:doc/) :: " + name}

	endpoint, exists := api.routes[path]
	if exists {
		input = append(input, renderEndpointDoc(endpoint.GetScheme())...)
	} else {
		input = append(input, renderIndexDoc(api)...)
	}

	output := blackfriday.Run([]byte(strings.Join(input, "\n")))

	w.Write([]byte(`
	<!DOCTYPE html>
	<html>
	<head>
		<meta charset="utf-8"/>
		<title>GODRAFT :: ` + name + `</title>
		<link
			rel="stylesheet"
			href="https://cdnjs.cloudflare.com/ajax/libs/github-markdown-css/4.0.0/github-markdown.min.css"
		/>
	</head>
	<body>
		<div class="markdown-body" style="margin: 20px">
	`))
	w.Write([]byte(output))
	w.Write([]byte(`
		</div>
	</body>
	</html>`))
}

func renderIndexDoc(api *APIService) []string {
	lines := []string{}

	for _, r := range api.routes {
		s := r.GetScheme()
		lines = append(lines, fmt.Sprintf(
			" - [%s](.%s) — %s",
			s.url,
			s.url,
			s.name,
		))
	}

	return lines
}

func renderEndpointDoc(s *Scheme) []string {
	lines := []string{
		"## " + s.url + " — " + s.name,
		s.descr,
	}

	json := s.ToJSON()
	for _, c := range json.Cases {
		lines = append(lines,
			`<div style="
				border-radius: 6px;
				border: 1px solid #ccc;
				padding: 0 15px 10px;
				margin: 0 10px 10px;
			">`,
			"<h3>"+getDocStatusIcon(c.Status)+" "+c.Name+"</h3>",
			rednerDocDescr(c.Description),
			renderDocRequest(json.Detail[c.Status].Request, c),
			renderDocResponse(json.Detail[c.Status].Response, c),
			"</div>",
		)
	}

	return lines
}

func smartDocTrim(v string) string {
	lines := strings.Split(v, "\n")

	switch len(lines) {
	case 0:
		return ""

	case 1, 2:
		return strings.TrimSpace(v)

	default:
		m := reFirstSpaces.FindString(lines[1])
		l := len(m)

		for i, v := range lines {
			if len(v) >= l {
				lines[i] = v[l:]
			}
		}

		return strings.Join(lines, "\n")
	}
}

func getDocStatusIcon(s StatusType) string {
	c := "#FF1100"
	if s == Status.OK {
		c = "#28C75D"
	}

	v := `
		color: #fff;
		display: inline-block;
		padding: 0 5px 2px;
		border-radius: 4px;
		background: ` + c + `;
	`
	return `<span style="` + v + `">` + string(s) + `</span>`
}

func renderDocRequest(r *JSONSchemeRequest, c *SchemeCase) string {
	str := []string{`<div style="
		margin: 0 -15px;
		padding: 15px 20px 5px;
		background: #F0F0F0;
	">`}

	p, ok := c.Params.(*orderedmap.OrderedMap)
	if ok {
		for _, k := range p.Keys() {
			raw, _ := p.Get(k)
			v := r.Params[k]
			if isDefaultDocValue(v, raw) {
				continue
			}

			str = append(str, fmt.Sprintf(
				`<div style="margin-bottom: 10px">%s<br/>
				<b>%s</b> = %v</div>`,
				renderDocComment(v),
				v.Name,
				renderJSONValue(v, raw),
			))
		}
	}

	str = append(str, "</div>")
	return strings.Join(str, "\n")
}

func renderDocResponse(r *JSONSchemeResponse, c *SchemeCase) string {
	str := []string{
		`<div style="
			margin: 0 -15px;
			padding: 15px 20px 5px;
		">`,
		renderDocJSON(r.Body, c.Body),
		`</div>`,
	}

	return strings.Join(str, "\n")
}

func renderDocJSON(d map[string]reflect.Item, raw interface{}) string {
	if raw == nil {
		return ""
	}

	json, ok := raw.(*orderedmap.OrderedMap)
	if !ok {
		return fmt.Sprintf(`<span style="color: #008079">"%v"</span>`, raw)
	}

	str := []string{`{<br/>`}
	keys := json.Keys()

	for i, key := range keys {
		glue := ","
		prop := d[key]
		val, _ := json.Get(key)
		if len(keys)-i == 1 {
			glue = ""
		}

		jv := renderJSONValue(prop, val)
		if prop.Type == "struct" && val != nil {
			jv = jv[0:len(jv)-5] + glue + "<br/>"
		} else {
			jv += glue
		}

		str = append(str, fmt.Sprintf(
			`<div style="padding-left: 20px; padding-bottom: 5px;">%s<br/>
			"<span>%s</span>": %s
			</div>`,
			renderDocComment(prop),
			prop.Name,
			jv,
		))
	}

	str = append(str, "}<br/>")

	return strings.Join(str, "\n")
}

func renderJSONValue(prop reflect.Item, raw interface{}) string {
	if raw == nil {
		return `<b style="color: #5856D6">null</b>`
	}

	switch prop.Type {
	case "string":
		return fmt.Sprintf(`<span style="color: #008079">"%v"</span>`, raw)

	case "bool":
		return fmt.Sprintf(`<span style="color: #FF2D54">%v</span>`, raw)

	case "struct":
		d := make(map[string]reflect.Item)
		for _, p := range prop.Nested {
			d[p.Name] = p
		}
		return renderDocJSON(d, raw)
	}

	return fmt.Sprintf(`<span style="color: #005FF9">%v</span>`, raw)
}

func isDefaultDocValue(v reflect.Item, raw interface{}) bool {
	switch v.Type {
	case "bool":
		return raw.(bool) == false
	}

	return false
}

func rednerDocDescr(v string) string {
	v = smartDocTrim(v)
	if v == "" {
		return ""
	}

	return `<div style="margin-bottom: 20px;">` + v + `</div>`
}

func renderDocComment(prop reflect.Item) string {
	t := prop.Type

	return fmt.Sprintf(
		`<span style="color: #999">/\* %s. <b>%s</b> \*/</span>`,
		prop.Comment,
		t,
	)
}

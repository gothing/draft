package draft

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// RequestFactory -
type RequestFactory func(RequestFactoryParams, RequestPrepare) (*http.Request, error)

// RequestPrepare -
type RequestPrepare func(req *http.Request) error

// RequestFactoryParams -
type RequestFactoryParams struct {
	Project     string     `json:"project"`
	Access      AccessType `json:"access"`
	AccessExtra string     `json:"access_extra"`
	Method      MethodType `json:"method"`
	Scheme      string     `json:"scheme"`
	Host        string     `json:"host"`
	Path        string     `json:"path"`
	Values      url.Values `json:"values"`
}

const (
	HeaderRequestID = "X-Request-Id"
)

// DefaultRequestFactory -
func DefaultRequestFactory(params RequestFactoryParams, prepare RequestPrepare) (*http.Request, error) {
	req, err := NewHTTPRequest(params)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = params.Values.Encode()

	if prepare != nil {
		err = prepare(req)
		if err != nil {
			return nil, err
		}
	}

	return req, nil
}

// NewHTTPRequest -
func NewHTTPRequest(params RequestFactoryParams) (*http.Request, error) {
	req := &http.Request{
		Method:     string(params.Method),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	req.Header = make(http.Header)
	req.Header.Set(HeaderRequestID, NewRequestID())

	scheme := params.Scheme
	if scheme == "" {
		scheme = "https"
	}

	rawURL := scheme + "://" + params.Host + params.Path
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, fmt.Errorf("url '%s' parse failed: %s", rawURL, err)
	}

	req.URL = u
	req.Host = u.Host

	return req, nil
}

// GetRequestPrepare -
func GetRequestPrepare(params RequestFactoryParams) RequestPrepare {
	for _, access := range pureDocConfig.Rights {
		if access.ID == params.Access {
			for _, extra := range access.Extra {
				if extra.Name == params.AccessExtra && extra.ReqPrepare != nil {
					return extra.ReqPrepare
				}
			}

			return access.ReqPrepare
		}
	}

	return nil
}

// NewRequestID -
func NewRequestID() string {
	id := make([]byte, 5)
	rand.Read(id)
	return hex.EncodeToString(id)
}

type requestFactoryResponse struct {
	General         requestFactoryResponseGeneral `json:"general"`
	ResponseHeaders http.Header                   `json:"response_headers"`
	RequestHeaders  http.Header                   `json:"request_headers"`
	QueryParams     url.Values                    `json:"query_params"`
	ResponseBody    interface{}                   `json:"response_body"`
}

type requestFactoryResponseGeneral struct {
	RequestURL    string `json:"request_url"`
	RequestMethod string `json:"request_method"`
	StatusCode    int    `json:"status_code"`
	StatusText    string `json:"status_text"`
	// RemoteAddress string `json:"remote_address"`
}

func doDraftRequest(api *APIService, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	params := RequestFactoryParams{}
	err := json.Unmarshal([]byte(r.URL.Query().Get("data")), &params)
	if err != nil {
		writeRequestFactoryError(w, "PARSE_JSON_REQUEST_DATA", err)
		return
	}

	requestFactory := DefaultRequestFactory
	if pureDocConfig.RequestFactory != nil {
		requestFactory = pureDocConfig.RequestFactory
	}

	req, err := requestFactory(params, GetRequestPrepare(params))
	if err != nil {
		writeRequestFactoryError(w, "CREATE_REQUEST", err)
		return
	}

	result := requestFactoryResponse{
		General:        requestFactoryResponseGeneral{req.URL.String(), req.Method, 0, ""},
		RequestHeaders: req.Header,
		QueryParams:    req.URL.Query(),
	}

	reqResp, err := api.endpointClient.Do(req)
	if err != nil {
		writeRequestFactoryError(w, "DO_REQUEST", err)
		return
	}

	result.General.StatusCode = reqResp.StatusCode
	result.General.StatusText = reqResp.Status
	result.ResponseHeaders = reqResp.Header

	defer reqResp.Body.Close()
	body, err := ioutil.ReadAll(reqResp.Body)

	if err != nil {
		result.ResponseBody = requestFactoryError{"READ_RESPONSE_BODY", err.Error()}
	} else if strings.Contains(reqResp.Header.Get("Content-Type"), "application/json") {
		r := make(map[string]interface{})
		err := json.Unmarshal(body, &r)
		if err != nil {
			result.ResponseBody = requestFactoryError{
				"PARSE_JSON_RESPONSE_BODY",
				fmt.Sprintf("parse %q failed: %s", string(body), err),
			}
		} else {
			result.ResponseBody = r
		}
	} else {
		result.ResponseBody = string(body)
	}

	jsonResult, err := json.Marshal(result)
	if err != nil {
		writeRequestFactoryError(w, "RESPONSE_MARSHAL", err)
		return
	}

	w.Write(jsonResult)
}

type requestFactoryError struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

func writeRequestFactoryError(w http.ResponseWriter, t string, err error) {
	log.Printf("[godraft:request] [warn] type: %q, message: %s", t, err)
	json, _ := json.Marshal(requestFactoryError{t, err.Error()})
	w.Write(json)
}

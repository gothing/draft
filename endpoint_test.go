package draft_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gothing/draft"
	"github.com/gothing/draft/types"
	"github.com/stretchr/testify/assert"
)

type UserEndpoint struct {
	draft.Endpoint
}

type UserEndpointBody struct {
	ID    types.UserID
	Flags UserFlags
}

type UserFlags struct {
	IsAdmin bool
	Deleted bool
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.URL("/api/v1/user")
	s.Case(draft.Status.OK, "Wow!", func() {
		s.Body(UserEndpointBody{
			ID: 20976,
		})
	})
}

type UserEndpoint2 struct {
	draft.Endpoint
}

func (ue *UserEndpoint2) InitEndpointScheme(s *draft.Scheme) {
	s.URL("/api/v1/user2")
}

type UserEndpoint3 struct {
	draft.Endpoint
}

func (ue *UserEndpoint3) InitEndpointScheme(s *draft.Scheme) {
	s.URL("/api/v1/user3")
}

func TestEndpoint(t *testing.T) {
	ue := &UserEndpoint{}
	api := draft.Create(draft.Config{DevMode: true})
	group := draft.Compose("test", ue)

	api.Add(group, nil)

	r := httptest.NewRequest("GET", "http://gothing/api/v1/user", nil)
	w := httptest.NewRecorder()

	assert.Equal(t, []string{"/api/v1/user"}, api.URLs(), "url")
	assert.Equal(t, ue.GetScheme().GetCaseByStatus(draft.Status.OK).Name, "Wow!")

	api.ServeHTTP(w, r)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, nil, err, "error")

	json := `{"status":"ok","body":{"id":20976,"flags":{"is_admin":false,"deleted":false}}}`
	assert.Equal(t, json, string(body), "Mock")
}

func TestEndpointWithHandler(t *testing.T) {
	ue := &UserEndpoint2{
		Endpoint: draft.Endpoint{
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK:" + r.URL.Path))
			},
		},
	}
	api := draft.Create(draft.Config{DevMode: true})
	group := draft.Compose("test", ue)

	api.Add(group, nil)

	r := httptest.NewRequest("GET", "http://gothing/api/v1/user2", nil)
	w := httptest.NewRecorder()

	api.ServeHTTP(w, r)
	resp := w.Result()
	body, err := ioutil.ReadAll(resp.Body)
	assert.Equal(t, nil, err, "error")

	result := `OK:/api/v1/user2`
	assert.Equal(t, result, string(body), "result")
}

type testGroupHandler struct {
}

func (tgh *testGroupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("WOW:" + r.URL.Path))
}

func (tgh *testGroupHandler) Routes() []string {
	return []string{"/api/v1/user3"}
}

// func TestEndpointWithGroupHandler(t *testing.T) {
// 	ue := &UserEndpoint3{}
// 	api := draft.Create()
// 	group := draft.Compose("test", ue)

// 	api.Add(group, &testGroupHandler{})

// 	r := httptest.NewRequest("GET", "http://gothing/api/v1/user3", nil)
// 	w := httptest.NewRecorder()

// 	api.ServeHTTP(w, r)
// 	resp := w.Result()
// 	body, err := ioutil.ReadAll(resp.Body)
// 	assert.Equal(t, nil, err, "error")

// 	result := `OK:/api/v1/user3`
// 	assert.Equal(t, result, string(body), "result")
// }

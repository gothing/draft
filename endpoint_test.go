package draft_test

import (
	"io/ioutil"
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

func (ue *UserEndpoint) Init() *UserEndpoint {
	ue.Endpoint.Init(ue)
	return ue
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.Case(draft.Status.OK, "Wow!", func() {
		s.Body(UserEndpointBody{
			ID: 20976,
		})
	})
}

func TestEndpoint(t *testing.T) {
	r := httptest.NewRequest("GET", "http://gothing/draft/user", nil)
	w := httptest.NewRecorder()

	ue := new(UserEndpoint).Init()
	assert.Equal(t, ue.GetScheme().GetCaseByStatus(draft.Status.OK).Name, "Wow!")

	ue.ServeHTTP(w, r)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)

	json := `{"status":"ok","body":{"id":20976,"flags":{"is_admin":false,"deleted":false}}}`
	assert.Equal(t, json, string(body), "Mock")
}

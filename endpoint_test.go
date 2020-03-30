package draft_test

import (
	"testing"

	"github.com/gothing/draft"
	"github.com/stretchr/testify/assert"
)

type UserEndpoint struct {
	draft.Endpoint
}

func (ue *UserEndpoint) Init() *UserEndpoint {
	ue.Endpoint.Init(ue)
	return ue
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.Case(draft.Status.OK, "Wow!", func() {

	})
}

func TestEndpoint(t *testing.T) {
	ue := new(UserEndpoint).Init()
	assert.Equal(t, ue.GetScheme().GetCaseByStatus(draft.Status.OK).Name, "Wow!")
}

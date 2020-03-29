package draft_test

import (
	"testing"

	"github.com/rubaxa/draft"
	"github.com/rubaxa/draft/types"
	"github.com/stretchr/testify/assert"
)

func TestScheme(t *testing.T) {
	s := &draft.Scheme{}
	s.Name("Foo")
	s.Description("Bar")
	s.Access(draft.Access.All)
	s.Case(draft.Status.OK, "Success", func() {
		s.Params(struct {
			Login types.Email `required:"true"`
		}{"qux@mail.ru"})
		s.Body(struct {
			Exists bool `comment:"Exists?"`
		}{})
	})

	r := s.ToJSON()
	assert.Equal(t, "login", r.Detail["ok"].Request.Params["login"].Name)
}

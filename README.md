draft
-----


#### Usage

```go
type API struct {
	User *UserEndpoint
}

type UserEndpoint struct {
	draft.Endpoint
}

func (ue *UserEndpoint) Init() *Endpoint {
	ue.Endpoint.Init(ue)
	return ue
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.Description("...")
	s.Case(draft.Status.OK, "User object by ID", func () {
		// ...
	})
}

func InitAPI() *API {
	return &API{
		User: new(UserEndpoint).Init(),
	}
}
```

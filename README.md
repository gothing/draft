draft
-----


### Usage

#### Setup API Endpoint

```go
type API struct {
	User *UserEndpoint
}

type UserEndpoint struct {
	draft.Endpoint
}

type UserEndpointParams struct {
	ID types.UserID `required:"true"`
}

type UserEndpointBody struct {
	ID types.UserID
}

func (ue *UserEndpoint) Init() *Endpoint {
	ue.Endpoint.Init(ue)
	return ue
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.Description("...")
	s.Case(draft.Status.OK, "User object by ID", func () {
		id := types.GenUserID()

		// Request params
		s.Params(&UserEndpointParams{
			ID: id,
		})

		// Response body
		s.Body(&UserEndpointBody{
			ID: id,
		})
	})
}

func InitAPI() *API {
	return &API{
		User: new(UserEndpoint).Init(),
	}
}
```

#### Add http handle

```go
func main() {
	api := InitAPI()
	http.Handle("/user", api.User)
}
```

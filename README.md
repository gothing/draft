draft
-----


### Usage

#### Setup API Endpoint

```go
type UserEndpoint struct {
	draft.Endpoint
}

type UserEndpointParams struct {
	ID types.UserID `required:"true"`
}

type UserEndpointBody struct {
	ID types.UserID
}

func (ue *UserEndpoint) InitEndpointScheme(s *draft.Scheme) {
	s.URL("/api/v1/user")
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
```

#### Add http handle

```go
func main() {
	draftSrv := draft.Create(darft.Config{
		DevMode: true, // Not production
	})
	userAPI := draft.Compose(
		"User API",
		new(UserEndpoint),
	)

draftSrv.Add(userAPI, nil)
	draftSrv.ListenAndServe(srv)
}
```

---

draft
-----
Self-documenting code

### Installing

```sh
go get github.com/gothing/draft
```

---

### Communication & Articles

- 📣https://t.me/joinchat/C0HLDxwwuQrH-lWBhBDbTA
- 📝https://medium.com/@ibnRubaXa/godraft-93307bb56794


---

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

#### Listen And Serve without implementation

But with mocks!

```go
func main() {
	draftSrv := draft.Create(darft.Config{
		DevMode: true, // Not production
	})
	userAPI := draft.Compose(
		"User API",
		new(UserEndpoint),
	)

	draftSrv.Add(userAPI)
	draftSrv.ListenAndServe(srv)
}
```

#### Listen And Serve with implementation

⚠️ See second argument for `draftSrv.Add`

```go
func main() {
	draftSrv := draft.Create(darft.Config{
		DevMode: true, // Not production
	})
	userAPI := draft.Compose(
		"User API",
		new(UserEndpoint),
	)

	draftSrv.Add(userAPI, MyDaemon{}) // <----- see here
	draftSrv.ListenAndServe(srv)
}

// MyDaemon - Example struct
type MyDaemon struct {}

// Routes - List of supported endpoints
func (d MyDaemon) Routes() []string {
	return []string{"/api/v1/user"} // rest of endpoints will remain as mock
}

// ServeHTTP - classic HTTP handler
func (d MyDaemon) ServeHTTP(w ResponseWriter, r *Request) {
	// ...
}
```

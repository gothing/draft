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

---

### Configuration

```go
func init() {
	draft.SetupDoc(draft.DocConfig{
		// Публичный UI, который `godraft` выкачивает и вставлчяет него конфигурацию,
		// так что нет никакой опасности, что схемы-концов утекут «наружу», но если вы
		// фанат интима, то клонируйте https://github.com/gothing/draft-front и поднемайте
		// где вам спокойней
		FrontURL: "https://gothing.github.io/draft-front/",

		// Выбранная активная группа
		ActiveGroup: "demo",

		// Сами группы и их схемы, как видите там может быть массив из разных истоникой,
		// т.е. `godraft` можно просто поднять отдельно, а не в месте с вашим кодом (но я не советую)
		Groups: []draft.DocGroup{
			// ID    Название         Источники описания «концов»
			{"demo", "Demo", []string{"http://localhost:2047/godraft:scheme/"}},
		},

		// Проекты и их найтсройки
		Projects: []draft.DocProject{
			{
				ID:      "auth",
				Name:    "Auth",
				Host:    "auth.mail.ru",
				HostRC:  "test.auth.mail.ru",
				HostDEV: "localhost:2047",
			},
		},

		// Права доступа
		Rights: []draft.DocAccess{
			{
				// Идентификатор, который вы используете в коде
				ID:   draft.Access.Auth,

				// Название, которое отображается в интефейсе
				Name: "Auth",

				// Детализация прав доступа (какие хедеры и/или параметры нужно передать дополнительно)
				Extra: []draft.DocAccessExtra{
					{
						Name: "mPOP",
						Headers: struct {
							Cookie string `json:"cookie" required:"true" comment:"mPOP-куки"`
						}{"Mpop=...;"},
						Params: struct {
							Token string `json:"token" required:"true" comment:"mPOP-токен"`
						}{"854724ce05861c2ce336e279039444a9%3A5441407e0..."},
					},
					{
						Name:        "OAuth",
						Description: "Читать [тут](https://oauth.net/2/)",
						Params: struct {
							AccessToken string `json:"access_token" required:"true" comment:"OAuth-токен"`
						}{"36ee693610a344929218133291cd27ca..."},
					},
				},
			},
		},
	})
}
```

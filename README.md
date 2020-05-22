draft
-----
Self-documenting code

[GODRAFT](https://repository-images.githubusercontent.com/251119592/10ed4a00-9b8a-11ea-9b08-ae5f7916de6f)

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

### Фабрика запросов

Первое что нужно понимать, она работает сразу из коробки, даже ничего настраивать не нужно, нажимаете кнопку в `[>]` в интерфейса, заполняете значения параметров запрос (по умолчанию они уже подставлены на основе выбранного `case`) и нажимаете Execute, изян 💁🏻‍♂️, запрос уходит в выбранный endpoint, а на выходе получается красивый JSON.

Но, всё становится интересней, когда доступ к «концу» ограничен «правами» (требуется авторизация или хитрая подписать на основе секретов от сервера), например:

Есть некий метод создания поста и доступен он только для авторизованного юзера:

```go
s.Url("/api/v1/post/create")
s.Access(draft.Access.Auth) // <--- вот это ограничение

s.Case(draft.Status.OK, "Успех", func () {
	s.Params(CreateParams{
		Title: "....",
		Content: "...",
	})

	s.Body(CreateBody{
		ID: types.GenID(),
	})
})
```

По умолчанию «Фабрика» ничего не знает о том как эти права интерпретировать, поэтому для начала в конфиге godraft нужно определить, как же работают у нас права:

```go
Rights: []draft.DocAccess{
	{
		ID: draft.Access.Auth,
		Name: "Auth",
		Extra: []draft.DocAccessExtra{
			{
				Name: "OAuth2",
				Params: struct {
					AccessToken string `json:"access_token" required:"true" comment:"OAuth2-токен"`
				}{"36e..."},
				ReqPrepare: prepareOAuth2Request, // <--- 🏭
			},
		},
	},
},
```

Фабрика, по полученым данным найдет `draft.Access.Auth` -> `OAuth2` и если у него определен метод ReqPrepare применит его к запросу:

```go
func prepareOAuth2Request(req *http.Request) error {
	// модифицируем `req` добавляя в него `access_token`
	// (например тестового юзера)
}
```

Вот и всё 🚀
_(благодаря этой штуке я могу через итерфейс (!) дергать интернал методы и не знать беды, если вы понимаете о чем я говорю 😏)_

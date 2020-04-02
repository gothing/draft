package types

import "net/url"

type Token string

func (v Token) TypeDescription() string {
	return "Токен"
}

type AGToken string

func (v AGToken) TypeDescription() string {
	return "Autogen-токен"
}

func GenAGToken() AGToken {
	return "AG_qaGCTRgWF9IBZezivsPEl"
}

type AuthToken struct {
	Token   AGToken `json:"token"`
	URL     string  `json:"url" comment:"Авторизационный URL"`
	Expires int     `json:"expires" comment:"Время жизни (в секундах)"`
}

func (v AuthToken) TypeDescription() string {
	return "Авторизационная сущность"
}

func GenAuthToken() AuthToken {
	t := GenAGToken()
	u, _ := url.Parse("https://auth.mail.ru/cgi-bin/auth")

	q := u.Query()
	q.Set("token", string(t))
	q.Set("Login", string(GenEmail()))
	q.Set("Page", "https://octavius.mail.ru/")
	u.RawQuery = q.Encode()

	return AuthToken{
		Token:   t,
		URL:     u.String(),
		Expires: 10 * 60,
	}
}

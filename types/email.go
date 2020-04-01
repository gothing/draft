package types

import (
	"strings"
)

// Email — email юзера
type Email string

// TypeDescription —
func (v Email) TypeDescription() string {
	return "Email пользователя"
}

// TypeValidity -
func (v Email) TypeValidity() TypeValidity {
	return TypeValidity{OK: strings.Contains(string(v), "@")}
}

// GenEmail -
func GenEmail() Email {
	return "fast.test@list.ru"
}

// GenCorpEmail -
func GenCorpEmail() Email {
	return "k.lebedev@corp.mail.ru"
}

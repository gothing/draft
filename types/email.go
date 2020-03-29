package types

import (
	"strings"
)

// Email — email юзера
type Email string

// TypeDescription —  gjkexbnm описание типа
func (v Email) TypeDescription() string {
	return "Email пользователя"
}

// TypeMockValue —
func (v Email) TypeMockValue() interface{} {
	return "fast.test@list.ru"
}

// TypeValidity -
func (v Email) TypeValidity() TypeValidity {
	return TypeValidity{OK: strings.Contains(string(v), "@")}
}

package types

import "github.com/sethvargo/go-password/password"

type Password string

func (v Password) TypeDescription() string {
	return "Пароль"
}

func GenPassword() Password {
	v, _ := password.Generate(12, 6, 6, false, false)
	return Password(v)
}

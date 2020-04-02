package types

type Login string

func (v Login) TypeDescription() string {
	return "Login пользователя (он же email, телефон или соц. аккаунт)"
}

func GenLogin() Login {
	return "fast.test@mail.ru"
}

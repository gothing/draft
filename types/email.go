package types

// Email — email юзера
type Email string

// TypeDescription —
func (v Email) TypeDescription() string {
	return "Email пользователя"
}

// XEmail —
type XEmail string

// TypeDescription —
func (v XEmail) TypeDescription() string {
	return "Email пользователя, указывает какую из учетных записей из мульти сессии использовать"
}

// GenEmail -
func GenEmail() Email {
	return "fast.test@list.ru"
}

// GenXEmail -
func GenXEmail() Email {
	return GenEmail()
}

// GenCorpEmail -
func GenCorpEmail() Email {
	return "k.lebedev@corp.mail.ru"
}

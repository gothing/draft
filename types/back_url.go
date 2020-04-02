package types

type BackURL string

func (v BackURL) TypeDescription() string {
	return "Страница (только почтовые проекты), куда нужно перенаправить юзера после авторизации (по умолчанию на https://e.mail.ru)"
}

func GenBackURL() BackURL {
	return "https://octavius.mail.ru/"
}

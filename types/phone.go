package types

type Phone string

func (v Phone) TypeDescription() string {
	return "Телефон"
}

func GenPhone() Phone {
	return "79096381667"
}

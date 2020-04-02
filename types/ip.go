package types

type IP string

func (v IP) TypeDescription() string {
	return "IP юзера"
}

func GenIP() IP {
	return "172.20.96.191"
}

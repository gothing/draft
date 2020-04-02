package types

type PNGBase64 string

func (v PNGBase64) TypeDescription() string {
	return "Base64 PNG-изображения"
}

func GenPNGBase64() PNGBase64 {
	return "data:image/png;base64,..."
}

package gen

import (
	"math/rand"
	
	"github.com/gothing/draft/types"
)

func ID() types.ID {
	return types.ID(rand.Int63() + 1)
}

func UserID() types.UserID {
	return types.UserID(rand.Int63() + 1)
}

func CorpEmail() types.Email {
	return "k.lebedev@corp.mail.ru"
}

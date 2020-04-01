package types

import (
	"math/rand"
)

// Bool -
type Bool = bool

// GenBool -
func GenBool() bool {
	return rand.Intn(2) == 1
}

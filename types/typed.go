package types

import (
	"math/rand"
	"time"
)

// TypedValue -
type TypedValue interface {
	TypeDescription() string
	TypeValidity() TypeValidity
}

// TypeValidity -
type TypeValidity struct {
	OK bool
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

package types

import "math/rand"

// UserID — User ID
type UserID uint64

// TypeDescription —
func (v UserID) TypeDescription() string {
	return "User ID"
}

// TypeValidity -
func (v UserID) TypeValidity() TypeValidity {
	return TypeValidity{OK: v > 0}
}

// GenUserID -
func GenUserID() UserID {
	return UserID(rand.Intn(100000) + 1)
}

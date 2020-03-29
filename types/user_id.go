package types

// UserID — ID
type UserID uint64

// TypeDescription —  
func (v UserID) TypeDescription() string {
	return "User ID"
}

// TypeValidity -
func (v UserID) TypeValidity() TypeValidity {
	return TypeValidity{OK: v > 0}
}

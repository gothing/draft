package types

// ID — ID
type ID uint64

// TypeDescription —  
func (v ID) TypeDescription() string {
	return "ID"
}

// TypeValidity -
func (v ID) TypeValidity() TypeValidity {
	return TypeValidity{OK: v > 0}
}

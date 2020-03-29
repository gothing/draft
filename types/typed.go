package types

// TypedValue -
type TypedValue interface {
	TypeDescription() string
	TypeValidity() TypeValidity
}

// TypeValidity -
type TypeValidity struct {
	OK bool
}

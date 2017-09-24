package typing

// Typing represents a type in expressive
type Typing interface {
	Equals(typing Typing) bool
	String() string
}

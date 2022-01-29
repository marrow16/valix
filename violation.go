package valix

// Violation contains information about an encountered validation violation
type Violation struct {
	// Property is the name of the property that failed validation
	Property string
	// Path is the path to the property that failed validation (in JSON format, e.g. "foo.bar[0].baz")
	Path string
	// Message is the violation message
	Message string
}

// NewViolation creates a new violation with the specified property, path and message
func NewViolation(property string, path string, msg string) *Violation {
	return &Violation{
		Property: property,
		Path:     path,
		Message:  msg,
	}
}

// NewEmptyViolation creates a new violation with the specified message (path and property are blank)
func NewEmptyViolation(msg string) *Violation {
	return &Violation{
		Property: "",
		Path:     "",
		Message:  msg,
	}
}

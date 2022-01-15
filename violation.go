package valix

// Violation contains information about an encountered validation violation
type Violation struct {
	// the name of the property that failed validation
	PropertyName string
	// to path to the property that failed validation
	Path string
	// the violation message
	Message string
}

func NewViolation(propertyName string, path string, msg string) *Violation {
	return &Violation{
		PropertyName: propertyName,
		Path:         path,
		Message:      msg,
	}
}

func newEmptyViolation(msg string) *Violation {
	return &Violation{
		PropertyName: "",
		Path:         "",
		Message:      msg,
	}
}

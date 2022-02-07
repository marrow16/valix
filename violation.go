package valix

// Violation contains information about an encountered validation violation
type Violation struct {
	// Property is the name of the property that failed validation
	Property string
	// Path is the path to the property that failed validation (in JSON format, e.g. "foo.bar[0].baz")
	Path string
	// Message is the violation message
	Message string
	// BadRequest is a flag indicating that the request could not be validated because
	// the payload was not JSON.  This effectively allows the caller of validation to determine
	// whether to respond with `400 Bad Request` or `422 Unprocessable Entity`
	//
	// Such violations are only added by the Validator.RequestValidate method where:
	//
	// * the request has an empty body
	//
	// * the request body does not parse (unmarshal) as JSON
	//
	// * the request body is JSON null (i.e. a request body containing just 'null') and the
	// Validator.AllowNull is set to false
	BadRequest bool
}

// NewEmptyViolation creates a new violation with the specified message (path and property are blank)
func NewEmptyViolation(msg string) *Violation {
	return &Violation{
		Property: "",
		Path:     "",
		Message:  msg,
	}
}

// NewViolation creates a new violation with the specified property, path and message
func NewViolation(property string, path string, msg string) *Violation {
	return &Violation{
		Property: property,
		Path:     path,
		Message:  msg,
	}
}

// NewBadRequestViolation creates a new violation with BadRequest flag set (path and property are blank)
func NewBadRequestViolation(msg string) *Violation {
	return &Violation{
		Property:   "",
		Path:       "",
		Message:    msg,
		BadRequest: true,
	}
}

package valix

// Constraint is the interface for all validation constraints on a property and object
//
// Custom constraints must implement this interface with the method:
type Constraint interface {
	// Validate Validates the constraint for a given value
	Validate(value interface{}, vcx *ValidatorContext) (bool, string)
	// GetMessage returns the actual message for the constraint
	GetMessage() string
}

type Validate func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string)

type CustomConstraint struct {
	validate Validate
	Message  string
}

// NewCustomConstraint Creates a custom Constraint which uses the supplied Validate function
func NewCustomConstraint(validate Validate, message string) *CustomConstraint {
	return &CustomConstraint{validate: validate, Message: message}
}
func (v *CustomConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	return v.validate(value, vcx, v)
}
func (v *CustomConstraint) GetMessage() string {
	return v.Message
}

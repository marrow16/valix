package valix

// Constraint is the interface for all validation constraints on a property and object
type Constraint interface {
	// Check the constraint against a given value
	Check(value interface{}, vcx *ValidatorContext) (bool, string)
	// GetMessage returns the actual message for the constraint
	//
	// This method is required so that any documenting functionality can determine
	// the constraint message without having to actually run the constraint
	GetMessage() string
}

// Check function signature for custom constraints
type Check func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string)

// CustomConstraint is a constraint that can declared on the fly and implements the Constraint interface
type CustomConstraint struct {
	CheckFunc Check
	Message   string
}

// NewCustomConstraint Creates a custom Constraint which uses the supplied Check function
func NewCustomConstraint(check Check, message string) *CustomConstraint {
	return &CustomConstraint{CheckFunc: check, Message: message}
}
func (v *CustomConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return v.CheckFunc(value, vcx, v)
}
func (v *CustomConstraint) GetMessage() string {
	return v.Message
}

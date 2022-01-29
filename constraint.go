package valix

// Constraint is the interface for all validation constraints on a property and object
//
// Custom constraints must implement this interface with the method:
//    Validate(property string, value interface{}, ctx *valix.Context) (bool, string)
type Constraint interface {
	// Validate Validates the constraint for a given value
	Validate(value interface{}, ctx *Context) (bool, string)
}

type Validate func(value interface{}, ctx *Context) (bool, string)

type customConstraint struct {
	validate Validate
}

// CustomConstraint Creates a custom Constraint which uses the supplied Validate function
func CustomConstraint(validate Validate) *customConstraint {
	return &customConstraint{validate: validate}
}
func (v *customConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	return v.validate(value, ctx)
}

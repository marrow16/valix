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
func (c *CustomConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return c.CheckFunc(value, vcx, c)
}
func (c *CustomConstraint) GetMessage() string {
	return c.Message
}

// ConstraintSet is a constraint that contains other constraints
//
// The contained constraints are checked sequentially but the overall
// set stops on the first failing constraint
type ConstraintSet struct {
	// Constraints is the slice of constraints within the set
	Constraints Constraints
	// Message is the violation message to be used if any of the constraints fail
	//
	// If the message is empty, the message from the failing contained constraint is used
	Message string
}

func (c *ConstraintSet) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	for _, cc := range c.Constraints {
		// don't use the `value` arg because contained constraints could change it...
		if ok, msg := cc.Check(vcx.CurrentValue(), vcx); !ok {
			if c.Message == "" {
				return false, msg
			} else {
				return false, c.GetMessage()
			}
		}
		if !vcx.continueAll || !vcx.continuePty() {
			break
		}
	}
	return true, ""
}
func (c *ConstraintSet) GetMessage() string {
	if c.Message == "" {
		for _, sc := range c.Constraints {
			if msg := sc.GetMessage(); msg != "" {
				return msg
			}
		}
	}
	return c.Message
}

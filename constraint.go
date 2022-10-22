package valix

// Constraint is the interface for all validation constraints on a property and object
type Constraint interface {
	// Check the constraint against a given value
	Check(value interface{}, vcx *ValidatorContext) (passed bool, message string)
	// GetMessage returns the actual message for the constraint
	//
	// This method is required so that any documenting functionality can determine
	// the constraint message without having to actually run the constraint
	GetMessage(tcx I18nContext) string
}

// Check function signature for custom constraints
type Check func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string)

type Conditional interface {
	MeetsConditions(vcx *ValidatorContext) bool
}

func isConditional(c Constraint) (Conditional, bool) {
	cc, ok := c.(Conditional)
	return cc, ok
}

func isCheckRequired(c Constraint, vcx *ValidatorContext) bool {
	if cc, ok := isConditional(c); ok {
		return cc.MeetsConditions(vcx)
	}
	return true
}

// CustomConstraint is a constraint that can declared on the fly and implements the Constraint interface
type CustomConstraint struct {
	CheckFunc Check
	Message   string
}

// NewCustomConstraint Creates a custom Constraint which uses the supplied Check function
func NewCustomConstraint(check Check, message string) *CustomConstraint {
	return &CustomConstraint{CheckFunc: check, Message: message}
}

// Check implements the Constraint.Check and calls the CustomConstraint.CheckFunc
func (c *CustomConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return c.CheckFunc(v, vcx, c)
}

// GetMessage implements the Constraint.GetMessage
func (c *CustomConstraint) GetMessage(tcx I18nContext) string {
	return obtainI18nContext(tcx).TranslateMessage(c.Message)
}

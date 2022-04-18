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

// ConstraintSet is a constraint that contains other constraints
//
// The contained constraints are checked sequentially but the overall
// set stops on the first failing constraint
type ConstraintSet struct {
	// Constraints is the slice of constraints within the set
	Constraints Constraints
	// when set to true, OneOf specifies that the constraint set should pass just one of
	// the contained constraints (rather than all of them)
	OneOf bool
	// Message is the violation message to be used if any of the constraints fail
	//
	// If the message is empty, the message from the first failing contained constraint is used
	Message string
	// Stop when set to true, prevents further validation checks on the property if this constraint set fails
	Stop bool
}

const (
	// used for marshaling, unmarshalling and tag parsing of constraint set...
	constraintSetName               = "ConstraintSet"
	constraintSetFieldConstraints   = "Constraints"
	constraintSetFieldOneOf         = "OneOf"
	constraintSetFieldMessage       = constraintPtyNameMessage
	constraintSetFieldStop          = constraintPtyNameStop
	fmtMsgConstraintSetDefaultAllOf = "Constraint set must pass all of %[1]d undisclosed validations"
	fmtMsgConstraintSetDefaultOneOf = "Constraint set must pass one of %[1]d undisclosed validations"
)

// Check implements the Constraint.Check and checks the constraints within the set
func (c *ConstraintSet) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.OneOf {
		finalOk := false
		firstMsg := ""
		for _, cc := range c.Constraints {
			if ok, msg := cc.Check(v, vcx); ok {
				finalOk = true
				break
			} else if firstMsg == "" {
				firstMsg = msg
			}
			if !vcx.continueAll || !vcx.continuePty() {
				break
			}
		}
		if finalOk {
			return true, ""
		}
		vcx.CeaseFurtherIf(c.Stop)
		if c.Message == "" && firstMsg != "" {
			return false, firstMsg
		}
		return false, c.GetMessage(vcx)
	}
	for _, cc := range c.Constraints {
		if ok, msg := cc.Check(v, vcx); !ok {
			if c.Message == "" && msg != "" {
				vcx.CeaseFurtherIf(c.Stop)
				return false, msg
			}
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		if !vcx.continueAll || !vcx.continuePty() {
			break
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *ConstraintSet) GetMessage(tcx I18nContext) string {
	if c.Message == "" {
		for _, sc := range c.Constraints {
			if msg := sc.GetMessage(tcx); msg != "" {
				return msg
			}
		}
		// if we get here then the message is still empty!...
		if c.OneOf {
			return obtainI18nContext(tcx).TranslateFormat(fmtMsgConstraintSetDefaultOneOf, len(c.Constraints))
		}
		return obtainI18nContext(tcx).TranslateFormat(fmtMsgConstraintSetDefaultAllOf, len(c.Constraints))
	}
	return obtainI18nContext(tcx).TranslateMessage(c.Message)
}

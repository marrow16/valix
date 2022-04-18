package valix

// FailingConstraint is a utility constraint that always fails
type FailingConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, StopAll stops the entire validation
	StopAll bool
}

// Check implements Constraint.Check
func (c *FailingConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	vcx.CeaseFurtherIf(c.Stop)
	if c.StopAll {
		vcx.Stop()
	}
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *FailingConstraint) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgFailure)
}

// FailWhen is a utility constraint that fails when specified conditions are met
type FailWhen struct {
	// the conditions under which to fail
	Conditions []string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, StopAll stops the entire validation
	StopAll bool
}

// Check implements Constraint.Check
func (c *FailWhen) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if vcx.meetsWhenConditions(c.Conditions) {
		vcx.CeaseFurtherIf(c.Stop)
		if c.StopAll {
			vcx.Stop()
		}
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *FailWhen) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgFailure)
}

// SetConditionFrom constraint is a utility constraint that can be used to set a condition in the
// ValidatorContext from string value of the property (to which this constraint is added)
//
// Note: It will only set a condition if the property value is a string!
type SetConditionFrom struct {
	// Parent by default, conditions are set on the current property or object - but specifying
	// true for this field means the condition is set on the parent object too
	Parent bool
	// Global setting this field to true means the condition is set for the entire
	// validator context
	Global bool
	// Prefix is any prefix to be appended to the string value
	Prefix string
	// Mapping converts the string value to alternate values (if the value is not found in the map
	// then the original value is used
	Mapping map[string]string
}

// Check implements Constraint.Check
func (c *SetConditionFrom) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if alt, aOk := c.Mapping[str]; aOk {
			str = alt
		}
		if c.Global {
			vcx.SetGlobalCondition(c.Prefix + str)
		} else if c.Parent {
			vcx.SetParentCondition(c.Prefix + str)
		} else {
			vcx.SetCondition(c.Prefix + str)
		}
	}
	return true, c.GetMessage(vcx)
}

func (c *SetConditionFrom) GetMessage(tcx I18nContext) string {
	return ""
}

// SetConditionProperty constraint is a utility constraint that can be used to set a condition in the
// ValidatorContext from string value of a specified property within the object to which this
// constraint is attached
//
// This constraint is normally only used in Validator.Constraints
//
// Note: It will only set a condition if the specified property value is a string!
type SetConditionProperty struct {
	// PropertyName is the name of the property to extract the condition value from
	PropertyName string `v8n:"default"`
	// Prefix is any prefix to be appended to the string value
	Prefix string
	// Mapping converts the string value to alternate values (if the value is not found in the map
	// then the original value is used
	Mapping map[string]string
}

// Check implements Constraint.Check
func (c *SetConditionProperty) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if m, ok := v.(map[string]interface{}); ok {
		if raw, ok := m[c.PropertyName]; ok {
			if str, ok := raw.(string); ok {
				if alt, aOk := c.Mapping[str]; aOk {
					str = alt
				}
				vcx.SetCondition(c.Prefix + str)
			}
		}
	}
	return true, c.GetMessage(vcx)
}

func (c *SetConditionProperty) GetMessage(tcx I18nContext) string {
	return ""
}

const (
	// used for marshaling, unmarshalling and tag parsing...
	constraintVariableProperty = "VariablePropertyConstraint"
)

// VariablePropertyConstraint is a special constraint that allows properties with varying names to be validated
//
// For example, if the following JSON needs to be validated:
//    {
//        "FOO": {
//            "amount": 123
//        },
//        "BAR": {
//            "amount": 345
//        }
//    }
// but the property names "FOO" and "BAR" could be anything - then the following validator would work:
//    validator := &valix.Validator{
//        IgnoreUnknownProperties: true, // Important!
//        Constraints: Constraints{
//            &VariablePropertyConstraint{
//                ObjectValidator: &Validator{
//                    Properties: valix.Properties{
//                        "amount": {
//                            Type:      valix.JsonInteger,
//                            Mandatory: true,
//                            NotNull:   true,
//                        },
//                    },
//                },
//            },
//        },
//    }
type VariablePropertyConstraint struct {
	// NameConstraints is a slice of Constraint items that are checked against the property name
	NameConstraints Constraints
	// ObjectValidator is the Validator to use for the property value
	ObjectValidator *Validator
	// AllowNull when set to true, allows property values to be null (and not checked against the ObjectValidator)
	AllowNull bool
}

func (c *VariablePropertyConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if m, ok := v.(map[string]interface{}); ok {
		for name, value := range m {
			nameOk := true
			for _, nc := range c.NameConstraints {
				if nOk, msg := nc.Check(name, vcx); !nOk {
					nameOk = false
					if msg != "" {
						vcx.addViolationPropertyForCurrent(name, msg, CodeInvalidPropertyName, name)
					} else {
						vcx.addViolationPropertyForCurrent(name, msgInvalidPropertyName, CodeInvalidPropertyName, name)
					}
					break
				}
			}
			if nameOk {
				if value == nil {
					if !c.AllowNull {
						vcx.addViolationPropertyForCurrent(name, msgValueCannotBeNull, CodeValueCannotBeNull)
					}
				} else if c.ObjectValidator != nil {
					vcx.pushPathProperty(name, value, c.ObjectValidator)
					if mValue, mOk := value.(map[string]interface{}); mOk {
						c.ObjectValidator.validate(mValue, vcx)
					} else {
						vcx.addUnTranslatedViolationForCurrent(msgPropertyValueMustBeObject, CodePropertyValueMustBeObject)
					}
					vcx.popPath()
				}
				if !vcx.continueAll {
					break
				}
			}
		}
	}
	return true, c.GetMessage(vcx)
}

func (c *VariablePropertyConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

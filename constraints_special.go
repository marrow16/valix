package valix

import (
	"encoding/json"
	"fmt"
	"math"
)

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
// ValidatorContext from the value of the property (to which this constraint is added)
//
// Note: It will only set a condition if the property value is a string!
type SetConditionFrom struct {
	// Parent by default, conditions are set on the current property or object - but specifying
	// true for this field means the condition is set on the parent object too
	Parent bool
	// Global setting this field to true means the condition is set for the entire
	// validator context
	Global bool
	// Prefix is any prefix to be appended to the condition token
	Prefix string
	// Mapping converts the string value to alternate values (if the value is not found in the map
	// then the original value is used
	Mapping map[string]string
	// NullToken is the condition token used if the value of the property is null/nil.  If this field is not set
	// and the property value is null at validation - then a condition token of "null" is used
	NullToken string
	// Format is an optional format string for dealing with non-string property values
	Format string
}

// Check implements Constraint.Check
func (c *SetConditionFrom) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	useToken := defaultString(c.NullToken, "null")
	if v != nil {
		if str, ok := v.(string); ok {
			useToken = str
		} else {
			useToken = fmt.Sprintf(defaultString(c.Format, "%v"), v)
		}
	}
	if alt, ok := c.Mapping[useToken]; ok {
		useToken = alt
	}
	if c.Global {
		vcx.SetGlobalCondition(c.Prefix + useToken)
	} else if c.Parent {
		vcx.SetParentCondition(c.Prefix + useToken)
	} else {
		vcx.SetCondition(c.Prefix + useToken)
	}
	return true, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *SetConditionFrom) GetMessage(tcx I18nContext) string {
	return ""
}

// SetConditionOnType constraint is a utility constraint that can be used to set a condition in the
// ValidatorContext indicating the type of the property value to which this constraint is added.
//
// The condition token set will be at least one of the following:
//   "type_null", "type_object", "type_array", "type_boolean", "type_string", "type_integer", "type_number", "type_unknown"
// Note that an int value will set both "type_number" and "type_integer" (because an int is both)
//
// On detecting a value represented by a json.Number - the "type_number" will always be set.  And this may also be
// complimented by the "type_integer" (if the json.Number holds an int value)
//
// Also, with json.Number values, the following condition tokens may also be set
//   "type_invalid_number", "type_nan", "type_inf"
// ("type_invalid_number" indicating that the json.Number could not be parsed to either int or float)
//
// When handling json.Number values, This constraint can be used in conjunction with a following FailWhen constraint
// to enforce failures in case of Inf, Nan or unparseable
type SetConditionOnType struct{}

var conditionTypeTokens = []string{"null", "object", "array", "boolean", "string", "integer", "number", "invalid_number", "nan", "inf", "unknown"}

// Check implements Constraint.Check
func (c *SetConditionOnType) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	const prefix = "type_"
	for _, clr := range conditionTypeTokens {
		vcx.ClearCondition(prefix + clr)
	}
	if v == nil {
		vcx.SetCondition(prefix + "null")
	} else {
		switch vt := v.(type) {
		case map[string]interface{}:
			vcx.SetCondition(prefix + "object")
		case []interface{}:
			vcx.SetCondition(prefix + "array")
		case bool:
			vcx.SetCondition(prefix + "boolean")
		case string:
			vcx.SetCondition(prefix + "string")
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			// set both integer and number (an int is both an integer and a number)
			vcx.SetCondition(prefix + "number")
			vcx.SetCondition(prefix + "integer")
		case float32, float64:
			vcx.SetCondition(prefix + "number")
		case json.Number:
			vcx.SetCondition(prefix + "number")
			if _, err := vt.Int64(); err == nil {
				vcx.SetCondition(prefix + "integer")
			} else if fv, err := vt.Float64(); err == nil {
				if math.IsNaN(fv) {
					vcx.SetCondition(prefix + "nan")
				} else if math.IsInf(fv, 0) {
					vcx.SetCondition(prefix + "inf")
				}
			} else {
				vcx.SetCondition(prefix + "invalid_number")
			}
		default:
			vcx.SetCondition(prefix + "unknown")
		}
	}
	return true, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *SetConditionOnType) GetMessage(tcx I18nContext) string {
	return ""
}

// SetConditionProperty constraint is a utility constraint that can be used to set a condition in the
// ValidatorContext from the value of a specified property within the object to which this
// constraint is attached
//
// This constraint is normally only used in Validator.Constraints
//
// Note: The property value can be of any type (inc. null) or, indeed, the property may be missing
type SetConditionProperty struct {
	// PropertyName is the name of the property to extract the condition value from
	PropertyName string `v8n:"default"`
	// Prefix is any prefix to be appended to the condition token
	Prefix string
	// Mapping converts the token value to alternate values (if the value is not found in the map
	// then the original value is used)
	Mapping map[string]string
	// NullToken is the condition token used if the value of the property specified is null/nil.  If this field is not set
	// and the property value is null at validation - then a condition token of "null" is used
	NullToken string
	// MissingToken is the condition token used if the property specified is missing.  If this field is not set
	// and the property is missing at validation - then a condition token of "missing" is used
	MissingToken string
	// Format is an optional format string for dealing with non-string property values
	Format string
}

// Check implements Constraint.Check
func (c *SetConditionProperty) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if m, ok := v.(map[string]interface{}); ok {
		useToken := defaultString(c.MissingToken, "missing")
		if raw, ok := m[c.PropertyName]; ok {
			if raw == nil {
				useToken = defaultString(c.NullToken, "null")
			} else if str, ok := raw.(string); ok {
				useToken = str
			} else {
				useToken = fmt.Sprintf(defaultString(c.Format, "%v"), raw)
			}
		}
		if alt, ok := c.Mapping[useToken]; ok {
			useToken = alt
		}
		vcx.SetCondition(c.Prefix + useToken)
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

const (
	// NoMessage is a special message used by constraints to indicate no message rather than default message
	NoMessage = "[NO_MESSAGE]"
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

const (
	// used for marshaling & unmarshalling...
	conditionalConstraintName = "ConditionalConstraint"
)

// ConditionalConstraint is a special constraint that wraps another constraint - but the wrapped
// constraint is only checked when the specified when conditions are met
type ConditionalConstraint struct {
	// When is the condition tokens that determine when the wrapped constraint is checked
	When Conditions
	// Constraint is the wrapped constraint
	Constraint Constraint
}

func (c *ConditionalConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if vcx.meetsWhenConditions(c.When) {
		return c.Constraint.Check(v, vcx)
	}
	return true, c.GetMessage(vcx)
}

func (c *ConditionalConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

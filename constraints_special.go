package valix

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
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

// FailWith is a utility constraint that fails when specified others property expression evaluates to true
type FailWith struct {
	// Others is the others expression to be evaluated to determine whether the constraint should fail
	Others OthersExpr `v8n:"default"`
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
func (c *FailWith) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.Others != nil {
		if curr, ancestryVals, ok := vcx.ancestorValueObject(0); ok {
			if c.Others.Evaluate(curr, ancestryVals, vcx) {
				vcx.CeaseFurtherIf(c.Stop)
				if c.StopAll {
					vcx.Stop()
				}
				return false, c.GetMessage(vcx)
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *FailWith) GetMessage(tcx I18nContext) string {
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

// SetConditionIf is a special constraint that wraps another constraint and sets a condition based on whether
// that wrapped constraint is ok or fails
//
// Note: the wrapped constraint cannot add any violations and cannot stop the validation (i.e. it is called 'silently')
type SetConditionIf struct {
	// is the wrapped constraint to be checked
	//
	// If this is nil, the SetOk condition is always set
	//
	// Note: the wrapped constraint cannot add any violations and cannot stop the validation (i.e. it is called 'silently')
	Constraint Constraint
	// is the condition to set if the wrapped constraint is ok
	//
	// Note: if this is an empty string - no condition is set
	SetOk string
	// is the condition to set if the wrapped constraint fails
	//
	// Note: if this is an empty string - no condition is set
	SetFail string
	// Parent by default, conditions are set on the current property or object - but specifying
	// true for this field means the condition is set on the parent object too
	Parent bool
	// Global setting this field to true means the condition is set for the entire
	// validator context
	Global bool
}

// Check implements Constraint.Check
func (c *SetConditionIf) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	setCond := c.SetOk
	if c.Constraint != nil {
		vcx.Lock()
		if ok, _ := c.Constraint.Check(v, vcx); !ok {
			setCond = c.SetFail
		}
		vcx.UnLock()
	}
	if setCond != "" {
		if c.Global {
			vcx.SetGlobalCondition(setCond)
		} else if c.Parent {
			vcx.SetParentCondition(setCond)
		} else {
			vcx.SetCondition(setCond)
		}
	}
	return true, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *SetConditionIf) GetMessage(tcx I18nContext) string {
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
	conditionalConstraintName      = "ConditionalConstraint"
	arrayconditionalConstraintName = "ArrayConditionalConstraint"
)

// ConditionalConstraint is a special constraint that wraps another constraint - but the wrapped
// constraint is only checked when the specified when condition is met
type ConditionalConstraint struct {
	// When is the condition tokens that determine when the wrapped constraint is checked
	When Conditions
	// Others is the others expression to be evaluated to determine when the wrapped constraint is checked
	Others OthersExpr
	// Constraint is the wrapped constraint
	Constraint Constraint
	// FailNotMet specifies that the conditional constraint should fail if the conditions are not met
	//
	// By default, if the conditions are not met the conditional constraint passes (without calling the wrapped constraint)
	FailNotMet bool
	// NotMetMessage is the message used when FailNotMet is set and the conditions are not met
	NotMetMessage string
}

func (c *ConditionalConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.MeetsConditions(vcx) {
		return c.Constraint.Check(v, vcx)
	}
	return !c.FailNotMet, c.GetMessage(vcx)
}

func (c *ConditionalConstraint) GetMessage(tcx I18nContext) string {
	return ternary(c.FailNotMet).string(c.NotMetMessage, "")
}

func (c *ConditionalConstraint) MeetsConditions(vcx *ValidatorContext) bool {
	if c.Constraint != nil && vcx.meetsWhenConditions(c.When) {
		if c.Others != nil {
			if curr, ancestryVals, ok := vcx.ancestorValueObject(0); ok {
				return c.Others.Evaluate(curr, ancestryVals, vcx)
			}
		} else {
			return true
		}
	}
	return false
}

// ArrayConditionalConstraint is a special constraint that wraps another constraint - but the wrapped
// constraint is only checked when the specified array condition is met (see When property)
type ArrayConditionalConstraint struct {
	// When is the special token denoting the array condition on which the wrapped constraint is to be checked
	//
	// The special token can be one of:
	//
	// * "first" - when the array item is the first
	//
	// * "!first" - when the array item is not the first
	//
	// * "last" - when the array item is the last
	//
	// * "!last" - when the array item is not the last
	//
	// * "%n" - when the modulus n of the array index is zero
	//
	// * ">n" - when the array index is greater than n
	//
	// * "<n" - when the array index is less than n
	//
	// * "n" - when the array index is n
	When string
	// Ancestry is ancestry depth at which to obtain the current array index information
	//
	// Note: the ancestry level is only for arrays in the object tree (and does not need to include other levels).
	// Therefore, by default the value is 0 (zero) - which means the last encountered array
	Ancestry   uint
	Constraint Constraint
}

func (c *ArrayConditionalConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if index, max, ok := vcx.AncestryIndex(c.Ancestry); ok {
		check := false
		switch c.When {
		case "first":
			check = index == 0
		case "!first":
			check = index != 0
		case "last":
			check = index == max
		case "!last":
			check = index != max
		default:
			if len(c.When) > 0 {
				pfx := c.When[:1]
				nOffs := 0
				if strings.ContainsAny(pfx, "%><") {
					nOffs = 1
				}
				if n, err := strconv.Atoi(c.When[nOffs:]); err == nil {
					switch pfx {
					case "%":
						check = (index % n) == 0
					case ">":
						check = index > n
					case "<":
						check = index < n
					default:
						check = index == n
					}
				}
			}
		}
		if check && c.Constraint != nil && isCheckRequired(c.Constraint, vcx) {
			return c.Constraint.Check(v, vcx)
		}
	}
	return true, c.GetMessage(vcx)
}

func (c *ArrayConditionalConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

// ConstraintSet is a constraint that contains other constraints
//
// The contained constraints are checked sequentially but the overall
// set stops on the first failing constraint
type ConstraintSet struct {
	// Constraints is the slice of constraints within the set
	Constraints Constraints `v8n:"default"`
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
		return c.checkOneOf(v, vcx)
	} else {
		return c.checkAllOf(v, vcx)
	}
}

func (c *ConstraintSet) checkAllOf(v interface{}, vcx *ValidatorContext) (bool, string) {
	for _, cc := range c.Constraints {
		if isCheckRequired(cc, vcx) {
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
	}
	return true, ""
}

func (c *ConstraintSet) checkOneOf(v interface{}, vcx *ValidatorContext) (bool, string) {
	finalOk := false
	firstMsg := ""
	for _, cc := range c.Constraints {
		if isCheckRequired(cc, vcx) {
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

// IsNull is a utility constraint to check that a value is null
//
// Normally, null checking would be performed by the PropertyValidator.NotNull setting - however,
// it may be the case that, under certain conditions, null is the required value
type IsNull struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements the Constraint.Check and checks the constraints within the set
func (c *IsNull) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if v != nil {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *IsNull) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNull)
}

// IsNotNull is a utility constraint to check that a value is not null
//
// Normally, null checking would be performed by the PropertyValidator.NotNull setting - however,
// it may be the case that null is only disallowed under certain conditions
type IsNotNull struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements the Constraint.Check and checks the constraints within the set
func (c *IsNotNull) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if v == nil {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *IsNotNull) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValueCannotBeNull)
}

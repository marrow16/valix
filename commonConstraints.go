package valix

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

const (
	messageValueNotEmptyString    = "Value must not be empty string"
	messageValueNotBlankString    = "Value must not be a blank string"
	messageValueNoControlChars    = "Value must not contain control characters"
	messageInvalidPattern         = "Value has invalid pattern"
	messageValueAtLeast           = "Value length must be at least %d"
	messageValueNotMore           = "Value length must not exceed %d"
	messageValueMinMax            = "Value length must be between %d and %d (inclusive)"
	messageValuePositive          = "Value must be positive"
	messageValuePositiveOrZero    = messageValuePositive + " or zero"
	messageValueNegative          = "Value must be negative"
	messageValueNegativeOrZero    = messageValueNegative + " or zero"
	messageValueGte               = "Value must be greater than or equal to %f"
	messageValueLte               = "Value must be less than or equal to %f"
	messageValueRange             = "Value must be between %f and %f (inclusive)"
	messageArrayElementType       = "Array value elements must be of type %s"
	messageArrayElementTypeOrNull = "Array value elements must be of type %s or null"
	messageValueValidUuid         = "Value must be a valid UUID"
	messageUuidMinVersion         = "Value must be a valid UUID (minimum version %d)"
	messageUuidCorrectVer         = "Value must be a valid UUID (version %d)"
	uuidRegexpPattern             = "([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})"
	iso8601FullPattern            = "(\\d{4}-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?)"
)

var (
	uuidRegexp = *regexp.MustCompile(uuidRegexpPattern)
)

// StringNotEmptyConstraint to check that string value is not empty (i.e. not "")
type StringNotEmptyConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNotEmptyConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringNotEmptyConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValueNotEmptyString)
}

// StringNotBlankConstraint to check that string value is not blank (i.e. that after removing leading and
// trailing whitespace the value is not an empty string)
type StringNotBlankConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNotBlankConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if len(strings.Trim(str, " \t\n\r")) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringNotBlankConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValueNotBlankString)
}

// StringNoControlCharsConstraint to check that a string does not contain any control characters (i.e. chars < 32)
type StringNoControlCharsConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNoControlCharsConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		for _, ch := range str {
			if ch < 32 {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *StringNoControlCharsConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValueNoControlChars)
}

// StringPatternConstraint to check that a string matches a given regexp pattern
type StringPatternConstraint struct {
	// the regexp pattern that the string value must match
	Regexp regexp.Regexp
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringPatternConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if !c.Regexp.MatchString(str) {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringPatternConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageInvalidPattern)
}

// StringMinLengthConstraint to check that a string has a minimum length
type StringMinLengthConstraint struct {
	// the minimum length value
	Value int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringMinLengthConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) < c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringMinLengthConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageValueAtLeast, c.Value))
}

// StringMaxLengthConstraint to check that a string has a maximum length
type StringMaxLengthConstraint struct {
	// the maximum length value
	Value int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringMaxLengthConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) > c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringMaxLengthConstraint) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageValueNotMore, c.Value))
}

// StringLengthConstraint to check that a string has a minimum and maximum length
type StringLengthConstraint struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringLengthConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringLengthConstraint) GetMessage() string {
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageValueMinMax, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageValueAtLeast, c.Minimum))
}

// LengthConstraint to check that a property value has minimum and maximum length
//
// This constraint can be used for string, object and array property values
type LengthConstraint struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *LengthConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, okS := value.(string); okS {
		if len(str) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if m, okM := value.(map[string]interface{}); okM {
		if len(m) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(m) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if a, okA := value.([]interface{}); okA {
		if len(a) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(a) > c.Maximum {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *LengthConstraint) GetMessage() string {
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageValueMinMax, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageValueAtLeast, c.Minimum))
}

// PositiveConstraint to check that a numeric value is positive (exc. zero)
type PositiveConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *PositiveConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f <= 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i <= 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv <= 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *PositiveConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValuePositive)
}

// PositiveOrZeroConstraint to check that a numeric value is positive or zero
type PositiveOrZeroConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *PositiveOrZeroConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i < 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv < 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *PositiveOrZeroConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValuePositiveOrZero)
}

// NegativeConstraint to check that a numeric value is negative
type NegativeConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *NegativeConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f >= 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i >= 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv >= 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *NegativeConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValueNegative)
}

// NegativeOrZeroConstraint to check that a numeric value is negative or zero
type NegativeOrZeroConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *NegativeOrZeroConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f > 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i > 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv > 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *NegativeOrZeroConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValueNegativeOrZero)
}

// MinimumConstraint to check that a numeric value is greater than or equal to a specified minimum
type MinimumConstraint struct {
	// the minimum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *MinimumConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < c.Value {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Value {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv < c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *MinimumConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageValueGte, c.Value))
}

// MaximumConstraint to check that a numeric value is less than or equal to a specified maximum
type MaximumConstraint struct {
	// the maximum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *MaximumConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f > c.Value {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) > c.Value {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv > c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *MaximumConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageValueLte, c.Value))
}

// RangeConstraint to check that a numeric value is within a specified minimum and maximum range
type RangeConstraint struct {
	// the minimum value of the range (inclusive)
	Minimum float64
	// the maximum value of the range (inclusive)
	Maximum float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *RangeConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < c.Minimum {
			return false, c.GetMessage()
		} else if f > c.Maximum {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Minimum {
			return false, c.GetMessage()
		} else if float64(i) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil {
			if fv < c.Minimum {
				return false, c.GetMessage()
			} else if fv > c.Maximum {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *RangeConstraint) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageValueRange, c.Minimum, c.Maximum))
}

// ArrayOfConstraint to check each element in an array value is of the correct type
type ArrayOfConstraint struct {
	// the type to check for each item (use PropertyType values)
	Type string
	// whether to allow null items in the array
	AllowNullElement bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *ArrayOfConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if a, ok := value.([]interface{}); ok {
		for _, elem := range a {
			if elem == nil {
				if !c.AllowNullElement {
					return false, c.GetMessage()
				}
			} else if !checkValueType(elem, c.Type) {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *ArrayOfConstraint) GetMessage() string {
	if c.AllowNullElement {
		return defaultMessage(c.Message, fmt.Sprintf(messageArrayElementTypeOrNull, c.Type))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageArrayElementType, c.Type))
}

// StringValidUuidConstraint to check that a string value is a valid UUID
type StringValidUuidConstraint struct {
	// the minimum UUID version (optional - if zero this is not checked)
	MinVersion uint8
	// the specific UUID version (optional - if zero this is not checked)
	SpecificVersion uint8
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c StringValidUuidConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	if str, ok := value.(string); ok {
		if !uuidRegexp.MatchString(str) {
			return false, c.GetMessage()
		}
		var version = str[14] - 48
		if c.MinVersion > 0 && version < c.MinVersion {
			return false, c.GetMessage()
		}
		if c.SpecificVersion > 0 && version != c.SpecificVersion {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c StringValidUuidConstraint) GetMessage() string {
	if c.SpecificVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidCorrectVer, c.SpecificVersion))
	} else if c.MinVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidMinVersion, c.MinVersion))
	}
	return defaultMessage(c.Message, messageValueValidUuid)
}

func defaultMessage(msg string, def string) string {
	if len(msg) == 0 {
		return def
	}
	return msg
}

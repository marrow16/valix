package valix

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	messageValueNotEmptyString = "Value must not be empty string"
	messageValueNotBlankString = "Value must not be a blank string"
	messageValueNoControlChars = "Value must not contain control characters"
	messageInvalidPattern      = "Value has invalid pattern"
	messageValueAtLeast        = "Value must have at least %d %s"
	messageValueNotMore        = "Value must not have more than %d %s"
	messageValuePositive       = "Value must be positive"
	messageValuePositiveOrZero = messageValuePositive + " or zero"
	messageValueNegative       = "Value must be negative"
	messageValueNegativeOrZero = messageValueNegative + " or zero"
	messageValueGte            = "Value must be greater than or equal to %f"
	messageValueLte            = "Value must be less than or equal to %f"
	messageArrayElementNull    = "Array value must not contain null elements (at index %d)"
	messageArrayElementType    = "Array value elements must be of type %s (at index %d)"
	messageValueNotValidUuid   = "Value must be a valid UUID"
	messageUuidMinVersion      = "Value UUID below minimum version (expected minimum version %d)"
	messageUuidIncorrectVer    = "Value UUID incorrect version (expected version %d)"
	wordCharacter              = "character"
	wordElement                = "element"
	wordProperty               = "property"
	wordProperties             = "properties"
	uuidRegexpPattern          = "([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})"
	iso8601FullPattern         = "(\\d{4}-\\d\\d-\\d\\dT\\d\\d:\\d\\d:\\d\\d(\\.\\d+)?(([+-]\\d\\d:\\d\\d)|Z)?)"
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
			return false, defaultMessage(c.Message, messageValueNotEmptyString)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValueNotBlankString)
		}
	}
	return true, ""
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
				return false, defaultMessage(c.Message, messageValueNoControlChars)
			}
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageInvalidPattern)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueAtLeast, c.Value, pluralize(c.Value, wordCharacter, "")))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueNotMore, c.Value, pluralize(c.Value, wordCharacter, "")))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueAtLeast, c.Minimum, pluralize(c.Minimum, wordCharacter, "")))
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueNotMore, c.Maximum, pluralize(c.Maximum, wordCharacter, "")))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueAtLeast, c.Minimum, pluralize(c.Minimum, wordCharacter, "")))
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, defaultMessage(c.Message,
				fmt.Sprintf(messageValueNotMore, c.Maximum, pluralize(c.Maximum, wordCharacter, "")))
		}
	} else if m, okM := value.(map[string]interface{}); okM {
		if len(m) < c.Minimum {
			return false, fmt.Sprintf(messageValueAtLeast, c.Minimum, pluralize(c.Minimum, wordProperty, wordProperties))
		} else if c.Maximum > 0 && len(m) > c.Maximum {
			return false, fmt.Sprintf(messageValueNotMore, c.Maximum, pluralize(c.Maximum, wordProperty, wordProperties))
		}
	} else if a, okA := value.([]interface{}); okA {
		if len(a) < c.Minimum {
			return false, fmt.Sprintf(messageValueAtLeast, c.Minimum, pluralize(c.Minimum, wordElement, ""))
		} else if c.Maximum > 0 && len(a) > c.Maximum {
			return false, fmt.Sprintf(messageValueNotMore, c.Maximum, pluralize(c.Maximum, wordElement, ""))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValuePositive)
		}
	} else if i, ok2 := value.(int); ok2 {
		if i <= 0 {
			return false, defaultMessage(c.Message, messageValuePositive)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValuePositiveOrZero)
		}
	} else if i, ok2 := value.(int); ok2 {
		if i < 0 {
			return false, defaultMessage(c.Message, messageValuePositiveOrZero)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValueNegative)
		}
	} else if i, ok2 := value.(int); ok2 {
		if i >= 0 {
			return false, defaultMessage(c.Message, messageValueNegative)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValueNegativeOrZero)
		}
	} else if i, ok2 := value.(int); ok2 {
		if i > 0 {
			return false, defaultMessage(c.Message, messageValueNegativeOrZero)
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueGte, c.Value))
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Value {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueGte, c.Value))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueLte, c.Value))
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) > c.Value {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueLte, c.Value))
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueGte, c.Minimum))
		} else if f > c.Maximum {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueLte, c.Maximum))
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Minimum {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueGte, c.Minimum))
		} else if float64(i) > c.Maximum {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageValueLte, c.Maximum))
		}
	}
	return true, ""
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
		for i, elem := range a {
			if elem == nil {
				if !c.AllowNullElement {
					return false, defaultMessage(c.Message, fmt.Sprintf(messageArrayElementNull, i))
				}
			} else if !checkValueType(elem, c.Type) {
				return false, defaultMessage(c.Message, fmt.Sprintf(messageArrayElementType, c.Type, i))
			}
		}
	}
	return true, ""
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
			return false, defaultMessage(c.Message, messageValueNotValidUuid)
		}
		var version = str[14] - 48
		if c.MinVersion > 0 && version < c.MinVersion {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageUuidMinVersion, c.MinVersion))
		}
		if c.SpecificVersion > 0 && version != c.SpecificVersion {
			return false, defaultMessage(c.Message, fmt.Sprintf(messageUuidIncorrectVer, c.SpecificVersion))
		}
	}
	return true, ""
}

func pluralize(val int, singular string, plural string) string {
	if val == 1 {
		return singular
	} else if len(plural) == 0 {
		return singular + "s"
	}
	return plural
}

func defaultMessage(msg string, def string) string {
	if len(msg) == 0 {
		return def
	}
	return msg
}

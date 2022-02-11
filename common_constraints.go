package valix

import (
	"encoding/json"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	messageNotEmptyString           = "String value must not be an empty string"
	messageNotBlankString           = "String value must not be a blank string"
	messageNoControlChars           = "String value must not contain control characters"
	messageInvalidPattern           = "String value must have valid pattern"
	messageInvalidToken             = "String value must be valid token - \"%s\""
	messageInvalidCharacters        = "String Value must not have invalid characters"
	messageStringMinLen             = "String value length must be at least %d characters"
	messageStringMaxLen             = "String value length must not exceed %d characters"
	messageStringExactLen           = "String value length must be %d characters"
	messageStringMinMaxLen          = "String value length must be between %d and %d (inclusive)"
	messageUnicodeNormalization     = "String value must be correct normalization form"
	messageUnicodeNormalizationNFC  = messageUnicodeNormalization + " NFC"
	messageUnicodeNormalizationNFKC = messageUnicodeNormalization + " NFKC"
	messageUnicodeNormalizationNFD  = messageUnicodeNormalization + " NFD"
	messageUnicodeNormalizationNFKD = messageUnicodeNormalization + " NFKD"
	messageMinLen                   = "Value length must be at least %d"
	messageExactLen                 = "Value length must be %d"
	messageMinMax                   = "Value length must be between %d and %d (inclusive)"
	messagePositive                 = "Value must be positive"
	messagePositiveOrZero           = messagePositive + " or zero"
	messageNegative                 = "Value must be negative"
	messageNegativeOrZero           = messageNegative + " or zero"
	messageGte                      = "Value must be greater than or equal to %f"
	messageLte                      = "Value must be less than or equal to %f"
	messageRange                    = "Value must be between %f and %f (inclusive)"
	messageArrayElementType         = "JsonArray value elements must be of type %s"
	messageArrayElementTypeOrNull   = "JsonArray value elements must be of type %s or null"
	messageValidUuid                = "Value must be a valid UUID"
	messageUuidMinVersion           = "Value must be a valid UUID (minimum version %d)"
	messageUuidCorrectVer           = "Value must be a valid UUID (version %d)"
	uuidRegexpPattern               = "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"
	messageValidCardNumber          = "Value must be a valid card number"
)

var (
	uuidRegexp = *regexp.MustCompile(uuidRegexpPattern)
)

// StringNotEmpty to check that string value is not empty (i.e. not "")
type StringNotEmpty struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringNotEmpty) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotEmpty) GetMessage() string {
	return defaultMessage(c.Message, messageNotEmptyString)
}

// StringNotBlank to check that string value is not blank (i.e. that after removing leading and
// trailing whitespace the value is not an empty string)
type StringNotBlank struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringNotBlank) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(strings.Trim(str, " \t\n\r")) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotBlank) GetMessage() string {
	return defaultMessage(c.Message, messageNotBlankString)
}

// StringNoControlCharacters to check that a string does not contain any control characters (i.e. chars < 32)
type StringNoControlCharacters struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringNoControlCharacters) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		for _, ch := range str {
			if ch < 32 {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNoControlCharacters) GetMessage() string {
	return defaultMessage(c.Message, messageNoControlChars)
}

// StringPattern to check that a string matches a given regexp pattern
type StringPattern struct {
	// the regexp pattern that the string value must match
	Regexp regexp.Regexp
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringPattern) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !c.Regexp.MatchString(str) {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringPattern) GetMessage() string {
	return defaultMessage(c.Message, messageInvalidPattern)
}

// StringValidToken checks that a string matches one of a pre-defined list of tokens
type StringValidToken struct {
	// Tokens is the set of allowed tokens for the string
	Tokens []string
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidToken) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		for _, t := range c.Tokens {
			if str == t {
				return true, ""
			}
		}
		return false, c.GetMessage()
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidToken) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageInvalidToken, strings.Join(c.Tokens, "\",\"")))
}

var _Bmp = &unicode.RangeTable{
	R16: []unicode.Range16{
		{0x0000, 0xffff, 1},
	},
}
var _Smp = &unicode.RangeTable{
	R32: []unicode.Range32{
		{0x10000, 0x1ffff, 1},
	},
}
var _Sip = &unicode.RangeTable{
	R32: []unicode.Range32{
		{0x20000, 0x2ffff, 1},
	},
}
var (
	// UnicodeBMP is the Unicode BMP (Basic Multilingual Plane)
	//
	// For use with StringCharacters
	UnicodeBMP = _Bmp
	// UnicodeSMP is the Unicode SMP (Supplementary Multilingual Plane)
	//
	// For use with StringCharacters
	UnicodeSMP = _Smp
	// UnicodeSIP is the Unicode SIP (Supplementary Ideographic Plane)
	//
	// For use with StringCharacters
	UnicodeSIP = _Sip
)

// StringCharacters to check that a string contains only allowable characters (and does not contain any disallowed characters)
type StringCharacters struct {
	// AllowRanges the ranges of characters (runes) that are allowed - each character
	// must be in at least one of these
	AllowRanges []*unicode.RangeTable
	// DisallowRanges the ranges of characters (runes) that are not allowed - if any character
	// is in any of these ranges then the constraint is violated
	DisallowRanges []*unicode.RangeTable
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringCharacters) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		runes := []rune(str)
		allowedCount := -1
		for i, r := range runes {
			for _, dr := range c.DisallowRanges {
				if unicode.Is(dr, r) {
					return false, c.GetMessage()
				}
			}
			for _, ar := range c.AllowRanges {
				if unicode.Is(ar, r) {
					allowedCount++
					break
				}
			}
			if i != allowedCount {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringCharacters) GetMessage() string {
	return defaultMessage(c.Message, messageInvalidCharacters)
}

// StringMinLength to check that a string has a minimum length
type StringMinLength struct {
	// the minimum length value
	Value int
	// UseRuneLen if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringMinLength) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l < c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMinLength) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageStringMinLen, c.Value))
}

// StringMaxLength to check that a string has a maximum length
type StringMaxLength struct {
	// the maximum length value
	Value int
	// UseRuneLen if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringMaxLength) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l > c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMaxLength) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageStringMaxLen, c.Value))
}

// StringLength to check that a string has a minimum and maximum length
type StringLength struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// UseRuneLen if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringLength) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && l > c.Maximum {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringLength) GetMessage() string {
	if c.Minimum == c.Maximum {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageStringExactLen, c.Minimum))
	}
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageStringMinMaxLen, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageStringMinLen, c.Minimum))
}

// StringValidUnicodeNormalization to check that a string has the correct Unicode normalization form
type StringValidUnicodeNormalization struct {
	// Form is the normalization form required - i.e. norm.NFC, norm.NFKC, norm.NFD or norm.NFKD
	//
	// (from package "golang.org/x/text/unicode/norm")
	Form norm.Form
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidUnicodeNormalization) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !c.Form.IsNormalString(str) {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidUnicodeNormalization) GetMessage() string {
	switch c.Form {
	case norm.NFKC:
		return defaultMessage(c.Message, messageUnicodeNormalizationNFKC)
	case norm.NFD:
		return defaultMessage(c.Message, messageUnicodeNormalizationNFD)
	case norm.NFKD:
		return defaultMessage(c.Message, messageUnicodeNormalizationNFKD)
	}
	return defaultMessage(c.Message, messageUnicodeNormalizationNFC)
}

// Length to check that a property value has minimum and maximum length
//
// This constraint can be used for string, object and array property values - however, if
// checking string lengths where actual Unicode length needs to be checked, it is better
// to use StringLength with UseRuneLength set to true
type Length struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Length) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Length) GetMessage() string {
	if c.Minimum == c.Maximum {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageExactLen, c.Minimum))
	}
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageMinMax, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageMinLen, c.Minimum))
}

// Positive to check that a numeric value is positive (exc. zero)
type Positive struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Positive) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Positive) GetMessage() string {
	return defaultMessage(c.Message, messagePositive)
}

// PositiveOrZero to check that a numeric value is positive or zero
type PositiveOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *PositiveOrZero) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *PositiveOrZero) GetMessage() string {
	return defaultMessage(c.Message, messagePositiveOrZero)
}

// Negative to check that a numeric value is negative
type Negative struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Negative) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Negative) GetMessage() string {
	return defaultMessage(c.Message, messageNegative)
}

// NegativeOrZero to check that a numeric value is negative or zero
type NegativeOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *NegativeOrZero) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *NegativeOrZero) GetMessage() string {
	return defaultMessage(c.Message, messageNegativeOrZero)
}

// Minimum to check that a numeric value is greater than or equal to a specified minimum
type Minimum struct {
	// the minimum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Minimum) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Minimum) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageGte, c.Value))
}

// Maximum to check that a numeric value is less than or equal to a specified maximum
type Maximum struct {
	// the maximum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Maximum) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Maximum) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageLte, c.Value))
}

// Range to check that a numeric value is within a specified minimum and maximum range
type Range struct {
	// the minimum value of the range (inclusive)
	Minimum float64
	// the maximum value of the range (inclusive)
	Maximum float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *Range) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *Range) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageRange, c.Minimum, c.Maximum))
}

// ArrayOf to check each element in an array value is of the correct type
type ArrayOf struct {
	// the type to check for each item (use Type values)
	Type JsonType
	// whether to allow null items in the array
	AllowNullElement bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *ArrayOf) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *ArrayOf) GetMessage() string {
	if c.AllowNullElement {
		return defaultMessage(c.Message, fmt.Sprintf(messageArrayElementTypeOrNull, c.Type))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageArrayElementType, c.Type))
}

// StringValidUuid to check that a string value is a valid UUID
type StringValidUuid struct {
	// the minimum UUID version (optional - if zero this is not checked)
	MinVersion uint8
	// the specific UUID version (optional - if zero this is not checked)
	SpecificVersion uint8
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidUuid) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
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

// GetMessage implements the Constraint.GetMessage
func (c *StringValidUuid) GetMessage() string {
	if c.SpecificVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidCorrectVer, c.SpecificVersion))
	} else if c.MinVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidMinVersion, c.MinVersion))
	}
	return defaultMessage(c.Message, messageValidUuid)
}

const (
	iso8601FullPattern             = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601NoOffsPattern           = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?)$"
	iso8601NoMillisPattern         = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601MinPattern              = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2})$"
	iso8601DateOnlyPattern         = "^(\\d{4}-\\d{2}-\\d{2})$"
	iso8601FullLayout              = "2006-01-02T15:04:05.999999999Z07:00"
	iso8601NoOffLayout             = "2006-01-02T15:04:05.999999999"
	iso8601NoMillisLayout          = "2006-01-02T15:04:05Z07:00"
	iso8601MinLayout               = "2006-01-02T15:04:05"
	iso8601DateOnlyLayout          = "2006-01-02"
	messageValidISODatetime        = "Value must be a valid date/time string"
	messageValidISODate            = "Value must be a valid date string (format: YYYY-MM-DD)"
	messageDatetimeFormatFull      = " (format: YYYY-MM-DDThh:mm:ss.sss[Z|+-hh:mm])"
	messageDatetimeFormatNoOffs    = " (format: YYYY-MM-DDThh:mm:ss.sss)"
	messageDatetimeFormatNoMillis  = " (format: YYYY-MM-DDThh:mm:ss[Z|+-hh:mm])"
	messageDatetimeFormatMin       = " (format: YYYY-MM-DDThh:mm:ss)"
	messageDatetimeFuture          = "Value must be a valid date/time in the future"
	messageDatetimeFutureOrPresent = "Value must be a valid date/time in the future or present"
	messageDatetimePast            = "Value must be a valid date/time in the past"
	messageDatetimePastOrPresent   = "Value must be a valid date/time in the past or present"
)

var (
	iso8601FullRegex     = regexp.MustCompile(iso8601FullPattern)
	iso8601NoOffsRegex   = regexp.MustCompile(iso8601NoOffsPattern)
	iso8601NoMillisRegex = regexp.MustCompile(iso8601NoMillisPattern)
	iso8601MinRegex      = regexp.MustCompile(iso8601MinPattern)
	iso8601DateOnlyRegex = regexp.MustCompile(iso8601DateOnlyPattern)
)

// StringValidISODatetime checks that a string value is a valid ISO8601 Date/time format
type StringValidISODatetime struct {
	// NoOffset specifies, if set to true, that time offsets are not permitted
	NoOffset bool
	// NoMillis specifies, if set to true, that seconds cannot have decimal places
	NoMillis bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidISODatetime) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		useRegex := iso8601FullRegex
		useLayout := iso8601FullLayout
		if c.NoOffset && c.NoMillis {
			useRegex = iso8601MinRegex
			useLayout = iso8601MinLayout
		} else if c.NoOffset {
			useRegex = iso8601NoOffsRegex
			useLayout = iso8601NoOffLayout
		} else if c.NoMillis {
			useRegex = iso8601NoMillisRegex
			useLayout = iso8601NoMillisLayout
		}
		if !useRegex.MatchString(str) {
			return false, c.GetMessage()
		}
		// and attempt to parse it (it may match regex but datetime could still be invalid)...
		if _, err := time.Parse(useLayout, str); err != nil {
			if pErr, ok := err.(*time.ParseError); ok && pErr.Message == "" && strings.HasSuffix(useLayout, pErr.LayoutElem) {
				// time.Parse is pretty dumb when it comes to timezones - if it's in the layout but not the string it fails
				// so remove the bit it doesn't like (in the layout) and try again...
				useLayout = useLayout[0 : len(useLayout)-len(pErr.LayoutElem)]
				if _, err := time.Parse(useLayout, str); err != nil {
					return false, c.GetMessage()
				}
			} else {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODatetime) GetMessage() string {
	if c.NoOffset && c.NoMillis {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatMin)
	} else if c.NoOffset {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatNoOffs)
	} else if c.NoMillis {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatNoMillis)
	}
	return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatFull)
}

// StringValidISODate checks that a string value is a valid ISO8601 Date format (excluding time)
type StringValidISODate struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidISODate) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !iso8601DateOnlyRegex.MatchString(str) {
			return false, c.GetMessage()
		}
		// and attempt to parse it (it may match regex but date could still be invalid)...
		if _, err := time.Parse(iso8601DateOnlyLayout, str); err != nil {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODate) GetMessage() string {
	return defaultMessage(c.Message, messageValidISODate)
}

// DatetimeFuture checks that a datetime/date (represented as string or time.Time) is in the future
type DatetimeFuture struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *DatetimeFuture) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || !dt.After(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && !dt.After(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFuture) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimeFuture)
}

// DatetimeFutureOrPresent checks that a datetime/date (represented as string or time.Time) is in the future or present
type DatetimeFutureOrPresent struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *DatetimeFutureOrPresent) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || dt.Before(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && dt.Before(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFutureOrPresent) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimeFutureOrPresent)
}

// DatetimePast checks that a datetime/date (represented as string or time.Time) is in the past
type DatetimePast struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *DatetimePast) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || !dt.Before(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && !dt.Before(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePast) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimePast)
}

// DatetimePastOrPresent checks that a datetime/date (represented as string or time.Time) is in the past or present
type DatetimePastOrPresent struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *DatetimePastOrPresent) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || dt.After(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && dt.After(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePastOrPresent) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimePastOrPresent)
}

// StringValidCardNumber checks that a string contains a valid card number according
// to Luhn Algorithm and checking that card number is 14 to 19 digits
type StringValidCardNumber struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

// Check implements Constraint.Check
func (c *StringValidCardNumber) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	const digitMinChar = '0'
	const digitMaxChar = '9'
	if str, ok := value.(string); ok {
		l := len(str)
		if l < 14 || l > 19 {
			return false, c.GetMessage()
		}
		var checkSum uint8 = 0
		doubling := false
		for i := l - 1; i >= 0; i-- {
			ch := str[i]
			if (ch < digitMinChar) || (ch > digitMaxChar) {
				return false, c.GetMessage()
			}
			digit := ch - digitMinChar
			if doubling && digit > 4 {
				digit = (digit * 2) - 9
			} else if doubling {
				digit = digit * 2
			}
			checkSum = checkSum + digit
			doubling = !doubling
		}
		if !(checkSum%10 == 0) {
			return false, c.GetMessage()
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidCardNumber) GetMessage() string {
	return defaultMessage(c.Message, messageValidCardNumber)
}

func stringToTime(str string) (*time.Time, bool) {
	parseLayout := ""
	if iso8601DateOnlyRegex.MatchString(str) {
		parseLayout = iso8601DateOnlyLayout
	} else if iso8601MinRegex.MatchString(str) {
		parseLayout = iso8601MinLayout
	} else if iso8601NoMillisRegex.MatchString(str) {
		parseLayout = iso8601NoMillisLayout
	} else if iso8601NoOffsRegex.MatchString(str) {
		parseLayout = iso8601NoOffLayout
	} else if iso8601FullRegex.MatchString(str) {
		parseLayout = iso8601FullLayout
	} else {
		return nil, false
	}
	result, err := time.Parse(parseLayout, str)
	return &result, err == nil
}

func defaultMessage(msg string, def string) string {
	if len(msg) == 0 {
		return def
	}
	return msg
}

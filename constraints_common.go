package valix

import (
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

// StringNotEmpty constraint to check that string value is not empty (i.e. not "")
type StringNotEmpty struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringNotEmpty) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if len(str) == 0 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotEmpty) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNotEmptyString)
}

// StringNotBlank constraint to check that string value is not blank (i.e. that after removing leading and
// trailing whitespace the value is not an empty string)
type StringNotBlank struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringNotBlank) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if len(strings.Trim(str, " \t\n\r")) == 0 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotBlank) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNotBlankString)
}

// StringNoControlCharacters constraint to check that a string does not contain any control characters (i.e. chars < 32)
type StringNoControlCharacters struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringNoControlCharacters) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		for _, ch := range str {
			if ch < 32 {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.GetMessage(vcx)
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNoControlCharacters) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNoControlChars)
}

// StringLowercase constraint to check that a string has only lowercase letters
type StringLowercase struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringLowercase) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if str != strings.ToLower(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringLowercase) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgStringLowercase)
}

// StringUppercase constraint to check that a string has only uppercase letters
type StringUppercase struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringUppercase) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if str != strings.ToUpper(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringUppercase) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgStringUppercase)
}

// StringPattern constraint to check that a string matches a given regexp pattern
type StringPattern struct {
	// the regexp pattern that the string value must match
	Regexp regexp.Regexp `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringPattern) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !c.Regexp.MatchString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringPattern) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidPattern)
}

// StringPresetPattern constraint to check that a string matches a given preset pattern
//
// Preset patterns are defined in PatternPresets (add your own where required)
//
// Messages for the preset patterns are defined in PatternPresetMessages
//
// If the preset pattern requires some extra validation beyond the regexp match, then add
// a checker to the PatternPresetPostPatternChecks variable
type StringPresetPattern struct {
	// the preset token (which must exist in the PatternPresets map)
	//
	// If the specified preset token does not exist - the constraint fails!
	Preset string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringPresetPattern) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if p, ok := presetsRegistry.get(c.Preset); ok {
			if !p.check(str) {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.getMessage(vcx, p.msg)
			}
		} else {
			vcx.CeaseFurtherIf(c.Stop)
			return false, vcx.TranslateFormat(fmtMsgUnknownPresetPattern, c.Preset)
		}
	}
	return true, ""
}

func (c *StringPresetPattern) getMessage(tcx I18nContext, msg string) string {
	if c.Message != "" {
		return obtainI18nContext(tcx).TranslateMessage(c.Message)
	} else if msg != "" {
		return obtainI18nContext(tcx).TranslateMessage(msg)
	}
	// still an empty message...
	return obtainI18nContext(tcx).TranslateMessage(msgValidPattern)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringPresetPattern) GetMessage(tcx I18nContext) string {
	if c.Message != "" {
		return obtainI18nContext(tcx).TranslateMessage(c.Message)
	} else if p, ok := presetsRegistry.get(c.Preset); ok && p.msg != "" {
		return obtainI18nContext(tcx).TranslateMessage(p.msg)
	}
	return obtainI18nContext(tcx).TranslateMessage(msgValidPattern)
}

// StringValidToken constraint checks that a string matches one of a pre-defined list of tokens
type StringValidToken struct {
	// the set of allowed tokens for the string
	Tokens []string `v8n:"default"`
	// set to true to make the token check case in-sensitive
	IgnoreCase bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidToken) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		lstr := strings.ToLower(str)
		for _, t := range c.Tokens {
			if str == t || (c.IgnoreCase && lstr == strings.ToLower(t)) {
				return true, ""
			}
		}
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidToken) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgValidToken, strings.Join(c.Tokens, "\",\""))
}

// StringCharacters constraint to check that a string contains only allowable characters (and does not contain any disallowed characters)
type StringCharacters struct {
	// the ranges of characters (runes) that are allowed - each character
	// must be in at least one of these
	AllowRanges []*unicode.RangeTable
	// the ranges of characters (runes) that are not allowed - if any character
	// is in any of these ranges then the constraint is violated
	DisallowRanges []*unicode.RangeTable
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringCharacters) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		runes := []rune(str)
		allowedCount := -1
		for i, r := range runes {
			for _, dr := range c.DisallowRanges {
				if unicode.Is(dr, r) {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			}
			for _, ar := range c.AllowRanges {
				if unicode.Is(ar, r) {
					allowedCount++
					break
				}
			}
			if i != allowedCount {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.GetMessage(vcx)
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringCharacters) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgInvalidCharacters)
}

// StringMinLength constraint to check that a string has a minimum length
type StringMinLength struct {
	// the minimum length value
	Value int `v8n:"default"`
	// when set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringMinLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l < c.Value || (c.ExclusiveMin && l == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMinLength) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgStringMinLenExc, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringMinLen, c.Value)
}

// StringMaxLength constraint to check that a string has a maximum length
type StringMaxLength struct {
	// the maximum length value
	Value int `v8n:"default"`
	// when set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringMaxLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l > c.Value || (c.ExclusiveMax && l == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMaxLength) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgStringMaxLenExc, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringMaxLen, c.Value)
}

// StringExactLength constraint to check that a string has an exact length
type StringExactLength struct {
	// the exact length expected
	Value int `v8n:"default"`
	// if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringExactLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringExactLength) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgStringExactLen, c.Value)
}

// StringLength constraint to check that a string has a minimum and maximum length
type StringLength struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// if set to true, uses the rune length (true Unicode length) to check length of string
	UseRuneLen bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		l := len(str)
		if c.UseRuneLen {
			l = len([]rune(str))
		}
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringLength) GetMessage(tcx I18nContext) string {
	if c.Maximum > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgStringMinMaxLen, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
	} else if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgStringMinLenExc, c.Minimum)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringMinLen, c.Minimum)
}

// StringValidUnicodeNormalization constraint to check that a string has the correct Unicode normalization form
type StringValidUnicodeNormalization struct {
	// the normalization form required - i.e. norm.NFC, norm.NFKC, norm.NFD or norm.NFKD
	//
	// (from package "golang.org/x/text/unicode/norm")
	Form norm.Form `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidUnicodeNormalization) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !c.Form.IsNormalString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidUnicodeNormalization) GetMessage(tcx I18nContext) string {
	switch c.Form {
	case norm.NFKC:
		return defaultMessage(tcx, c.Message, msgUnicodeNormalizationNFKC)
	case norm.NFD:
		return defaultMessage(tcx, c.Message, msgUnicodeNormalizationNFD)
	case norm.NFKD:
		return defaultMessage(tcx, c.Message, msgUnicodeNormalizationNFKD)
	}
	return defaultMessage(tcx, c.Message, msgUnicodeNormalizationNFC)
}

// Length constraint to check that a property value has minimum and maximum length
//
// This constraint can be used for object, array and string property values - however, if
// checking string lengths it is better to use the StringLength constraint
//
// * when checking array values, the number of elements in the array is checked
//
// * when checking object values, the number of properties in the object is checked
//
// * when checking string values, the length of the string is checked
type Length struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Length) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, okS := v.(string); okS {
		l := len(str)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if m, okM := v.(map[string]interface{}); okM {
		l := len(m)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if a, okA := v.([]interface{}); okA {
		l := len(a)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Length) GetMessage(tcx I18nContext) string {
	if c.Maximum > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgMinMax, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
	} else if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgMinLenExc, c.Minimum)
	}
	return defaultMessage(tcx, c.Message, fmtMsgMinLen, c.Minimum)
}

// LengthExact constraint to check that a property value has a specific length
//
// This constraint can be used for object, array and string property values - however, if
// checking string lengths it is better to use the StringExactLength constraint
//
// * when checking array values, the number of elements in the array is checked
//
// * when checking object values, the number of properties in the object is checked
//
// * when checking string values, the length of the string is checked
type LengthExact struct {
	// the length to check
	Value int `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LengthExact) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, okS := v.(string); okS {
		l := len(str)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if m, okM := v.(map[string]interface{}); okM {
		l := len(m)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if a, okA := v.([]interface{}); okA {
		l := len(a)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *LengthExact) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgExactLen, c.Value)
}

// Positive constraint to check that a numeric value is positive (exc. zero)
type Positive struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Positive) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f <= 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Positive) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgPositive)
}

// PositiveOrZero constraint to check that a numeric value is positive or zero
type PositiveOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *PositiveOrZero) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f < 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *PositiveOrZero) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgPositiveOrZero)
}

// Negative constraint to check that a numeric value is negative
type Negative struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Negative) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f >= 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Negative) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNegative)
}

// NegativeOrZero constraint to check that a numeric value is negative or zero
type NegativeOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NegativeOrZero) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f > 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NegativeOrZero) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNegativeOrZero)
}

// Minimum constraint to check that a numeric value is greater than or equal to a specified minimum
type Minimum struct {
	// the minimum value
	Value float64 `v8n:"default"`
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Minimum) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f < c.Value || (c.ExclusiveMin && f == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Minimum) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgGt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgGte, c.Value)
}

// Maximum constraint to check that a numeric value is less than or equal to a specified maximum
type Maximum struct {
	// the maximum value
	Value float64 `v8n:"default"`
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Maximum) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f > c.Value || (c.ExclusiveMax && f == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Maximum) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgLt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgLte, c.Value)
}

// Range constraint to check that a numeric value is within a specified minimum and maximum range
type Range struct {
	// the minimum value of the range
	Minimum float64
	// the maximum value of the range
	Maximum float64
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Range) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f < c.Minimum || (c.ExclusiveMin && f == c.Minimum) ||
			f > c.Maximum || (c.ExclusiveMax && f == c.Maximum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Range) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgRange, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
}

// MinimumInt constraint to check that an integer numeric value is greater than or equal to a specified minimum
type MinimumInt struct {
	// the minimum value
	Value int64 `v8n:"default"`
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *MinimumInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i < c.Value || (c.ExclusiveMin && i == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MinimumInt) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgGt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgGte, c.Value)
}

// MaximumInt constraint to check that an integer value is less than or equal to a specified maximum
type MaximumInt struct {
	// the maximum value
	Value int64 `v8n:"default"`
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *MaximumInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i > c.Value || (c.ExclusiveMax && i == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MaximumInt) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgLt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgLte, c.Value)
}

// RangeInt constraint to check that an integer value is within a specified minimum and maximum range
type RangeInt struct {
	// the minimum value of the range
	Minimum int64
	// the maximum value of the range
	Maximum int64
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *RangeInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i < c.Minimum || (c.ExclusiveMin && i == c.Minimum) ||
			i > c.Maximum || (c.ExclusiveMax && i == c.Maximum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *RangeInt) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgRange, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
}

// MultipleOf constraint to check that an integer value is a multiple of a specific number
//
// Note: this constraint will check values that are float or json Number - but the
// check will fail if either of these is not a 'whole number'
type MultipleOf struct {
	// the multiple of value to check
	Value int64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *MultipleOf) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i%c.Value != 0 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MultipleOf) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgMultipleOf, c.Value)
}

// ArrayOf constraint to check each element in an array value is of the correct type
type ArrayOf struct {
	// the type to check for each item (use Type values)
	Type string `v8n:"default"`
	// whether to allow null items in the array
	AllowNullElement bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *ArrayOf) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if a, ok := v.([]interface{}); ok {
		if chkType, tOk := JsonTypeFromString(c.Type); tOk {
			for _, elem := range a {
				if elem == nil {
					if !c.AllowNullElement {
						vcx.CeaseFurtherIf(c.Stop)
						return false, c.GetMessage(vcx)
					}
				} else if !checkValueType(elem, chkType) {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *ArrayOf) GetMessage(tcx I18nContext) string {
	if c.AllowNullElement {
		return defaultMessage(tcx, c.Message, fmtMsgArrayElementTypeOrNull, c.Type)
	}
	return defaultMessage(tcx, c.Message, fmtMsgArrayElementType, c.Type)
}

// ArrayUnique constraint to check each element in an array value is unique
type ArrayUnique struct {
	// whether to ignore null items in the array
	IgnoreNulls bool `v8n:"default"`
	// whether uniqueness is case in-insensitive (for string elements)
	IgnoreCase bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *ArrayUnique) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if a, ok := v.([]interface{}); ok {
		list := make([]interface{}, 0, len(a))
		for _, iv := range a {
			if !(iv == nil && c.IgnoreNulls) {
				if !isUniqueCompare(iv, c.IgnoreCase, &list) {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *ArrayUnique) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgArrayUnique)
}

// StringValidUuid constraint to check that a string value is a valid UUID
type StringValidUuid struct {
	// the minimum UUID version (optional - if zero this is not checked)
	MinVersion uint8
	// the specific UUID version (optional - if zero this is not checked)
	SpecificVersion uint8
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidUuid) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !uuidRegexp.MatchString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		var version = str[14] - 48
		if c.MinVersion > 0 && version < c.MinVersion {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		if c.SpecificVersion > 0 && version != c.SpecificVersion {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidUuid) GetMessage(tcx I18nContext) string {
	if c.SpecificVersion > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgUuidCorrectVer, c.SpecificVersion)
	} else if c.MinVersion > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgUuidMinVersion, c.MinVersion)
	}
	return defaultMessage(tcx, c.Message, msgValidUuid)
}

// StringValidISODatetime constraint checks that a string value is a valid ISO8601 Date/time format
type StringValidISODatetime struct {
	// specifies, if set to true, that time offsets are not permitted
	NoOffset bool
	// specifies, if set to true, that seconds cannot have decimal places
	NoMillis bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidISODatetime) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
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
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		// and attempt to parse it (it may match regex but datetime could still be invalid)...
		if _, err := time.Parse(useLayout, str); err != nil {
			if pErr, ok := err.(*time.ParseError); ok && pErr.Message == "" && strings.HasSuffix(useLayout, pErr.LayoutElem) {
				// time.Parse is pretty dumb when it comes to timezones - if it's in the layout but not the string it fails
				// so remove the bit it doesn't like (in the layout) and try again...
				useLayout = useLayout[0 : len(useLayout)-len(pErr.LayoutElem)]
				if _, err := time.Parse(useLayout, str); err != nil {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			} else {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.GetMessage(vcx)
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODatetime) GetMessage(tcx I18nContext) string {
	if c.NoOffset && c.NoMillis {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatMin)
	} else if c.NoOffset {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatNoOffs)
	} else if c.NoMillis {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatNoMillis)
	}
	return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatFull)
}

// StringValidISODate constraint checks that a string value is a valid ISO8601 Date format (excluding time)
type StringValidISODate struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidISODate) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !iso8601DateOnlyRegex.MatchString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		// and attempt to parse it (it may match regex but date could still be invalid)...
		if _, err := time.Parse(iso8601DateOnlyLayout, str); err != nil {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODate) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidISODate)
}

// DatetimeFuture constraint checks that a datetime/date (represented as string or time.Time) is in the future
type DatetimeFuture struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeFuture) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if dt, ok2 := stringToDatetime(str, c.ExcTime); !ok2 || !dt.After(truncateDate(time.Now(), c.ExcTime)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if dt, ok2 := v.(time.Time); ok2 && !truncateDate(dt, c.ExcTime).After(truncateDate(time.Now(), c.ExcTime)) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFuture) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimeFuture)
}

// DatetimeFutureOrPresent constraint checks that a datetime/date (represented as string or time.Time) is in the future or present
type DatetimeFutureOrPresent struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeFutureOrPresent) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if dt, ok2 := stringToDatetime(str, c.ExcTime); !ok2 || dt.Before(truncateDate(time.Now(), c.ExcTime)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if dt, ok2 := v.(time.Time); ok2 && truncateDate(dt, c.ExcTime).Before(truncateDate(time.Now(), c.ExcTime)) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFutureOrPresent) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimeFutureOrPresent)
}

// DatetimePast constraint checks that a datetime/date (represented as string or time.Time) is in the past
type DatetimePast struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimePast) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if dt, ok2 := stringToDatetime(str, c.ExcTime); !ok2 || !dt.Before(truncateDate(time.Now(), c.ExcTime)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if dt, ok2 := v.(time.Time); ok2 && !truncateDate(dt, c.ExcTime).Before(truncateDate(time.Now(), c.ExcTime)) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePast) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimePast)
}

// DatetimePastOrPresent constraint checks that a datetime/date (represented as string or time.Time) is in the past or present
type DatetimePastOrPresent struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimePastOrPresent) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if dt, ok2 := stringToDatetime(str, c.ExcTime); !ok2 || !dt.Before(truncateDate(time.Now(), c.ExcTime)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if dt, ok2 := v.(time.Time); ok2 && truncateDate(dt, c.ExcTime).After(truncateDate(time.Now(), c.ExcTime)) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePastOrPresent) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimePastOrPresent)
}

// StringValidCardNumber constraint checks that a string contains a valid card number according
// to Luhn Algorithm and checking that card number is 14 to 19 digits
type StringValidCardNumber struct {
	// if set to true, AllowSpaces accepts space separators in the card number (but must appear between each 4 digits)
	AllowSpaces bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidCardNumber) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	const digitMinChar = '0'
	const digitMaxChar = '9'
	if str, ok := v.(string); ok {
		buffer := []byte(str)
		l := len(buffer)
		if c.AllowSpaces {
			stripBuffer := make([]byte, 0, len(str))
			for i, by := range buffer {
				if by == ' ' {
					if (i+1)%5 != 0 || i+1 == l {
						vcx.CeaseFurtherIf(c.Stop)
						return false, c.GetMessage(vcx)
					}
				} else {
					stripBuffer = append(stripBuffer, by)
				}
			}
			buffer = stripBuffer
			l = len(buffer)
		}
		if l < 12 || l > 19 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		var checkSum uint8 = 0
		doubling := false
		for i := l - 1; i >= 0; i-- {
			ch := buffer[i]
			if (ch < digitMinChar) || (ch > digitMaxChar) {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.GetMessage(vcx)
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
		if checkSum%10 != 0 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidCardNumber) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidCardNumber)
}

// StringValidEmail constraint checks that a string contains a valid email address (does not
// verify the email address!)
//
// NB. Uses mail.ParseAddress to check valid email address
type StringValidEmail struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidEmail) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if _, err := mail.ParseAddress(str); err != nil {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidEmail) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidEmail)
}

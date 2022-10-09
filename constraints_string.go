package valix

import (
	"encoding/json"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"unicode"
)

// StringCharacters constraint to check that a string contains only allowable characters (and does not contain any disallowed characters)
type StringCharacters struct {
	// the ranges of characters (runes) that are allowed - each character
	// must be in at least one of these
	AllowRanges []unicode.RangeTable
	// the ranges of characters (runes) that are not allowed - if any character
	// is in any of these ranges then the constraint is violated
	DisallowRanges []unicode.RangeTable
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringCharacters) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringCharacters) checkString(str string, vcx *ValidatorContext) bool {
	runes := []rune(str)
	allowedCount := -1
	for i, r := range runes {
		for _, dr := range c.DisallowRanges {
			if unicode.Is(&dr, r) {
				return false
			}
		}
		for _, ar := range c.AllowRanges {
			if unicode.Is(&ar, r) {
				allowedCount++
				break
			}
		}
		if i != allowedCount {
			return false
		}
	}
	return true
}

// GetMessage implements the Constraint.GetMessage
func (c *StringCharacters) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgInvalidCharacters)
}

// StringContains constraint to check that a string contains with a given value
type StringContains struct {
	// the value to check that the string contains
	Value string `v8n:"default"`
	// multiple additional values that the string may contain
	Values []string
	// whether the check is case-insensitive (by default, the check is case-sensitive)
	CaseInsensitive bool
	// whether the check is NOT-ed (i.e. checks that the string does not contain)
	Not bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringContains) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringContains) checkString(str string, vcx *ValidatorContext) bool {
	ckStr := caseInsensitive(str, c.CaseInsensitive)
	contains := c.Value != "" && strings.Contains(ckStr, caseInsensitive(c.Value, c.CaseInsensitive))
	if !contains {
		for _, s := range c.Values {
			contains = s != "" && strings.Contains(ckStr, caseInsensitive(s, c.CaseInsensitive))
			if contains {
				break
			}
		}
	}
	return (!c.Not && contains) || (c.Not && !contains)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringContains) GetMessage(tcx I18nContext) string {
	possibles := make([]string, 0)
	if c.Value != "" {
		possibles = append(possibles, "'"+c.Value+"'")
	}
	for _, s := range c.Values {
		if s != "" {
			possibles = append(possibles, "'"+s+"'")
		}
	}
	possiblesStr := strings.Join(possibles, ",")
	if c.Not {
		return defaultMessage(tcx, c.Message, fmtMsgStringNotContains, possiblesStr)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringContains, possiblesStr)
}

// StringEndsWith constraint to check that a string ends with a given suffix
type StringEndsWith struct {
	// the value to check that the string ends with
	Value string `v8n:"default"`
	// multiple additional values that the string may end with
	Values []string
	// whether the check is case-insensitive (by default, the check is case-sensitive)
	CaseInsensitive bool
	// whether the check is NOT-ed (i.e. checks that the string does not end with)
	Not bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringEndsWith) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringEndsWith) checkString(str string, vcx *ValidatorContext) bool {
	ckStr := caseInsensitive(str, c.CaseInsensitive)
	endsWith := c.Value != "" && strings.HasSuffix(ckStr, caseInsensitive(c.Value, c.CaseInsensitive))
	if !endsWith {
		for _, s := range c.Values {
			endsWith = s != "" && strings.HasSuffix(ckStr, caseInsensitive(s, c.CaseInsensitive))
			if endsWith {
				break
			}
		}
	}
	return (!c.Not && endsWith) || (c.Not && !endsWith)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringEndsWith) GetMessage(tcx I18nContext) string {
	possibles := make([]string, 0)
	if c.Value != "" {
		possibles = append(possibles, "'"+c.Value+"'")
	}
	for _, s := range c.Values {
		if s != "" {
			possibles = append(possibles, "'"+s+"'")
		}
	}
	possiblesStr := strings.Join(possibles, ",")
	if c.Not {
		return defaultMessage(tcx, c.Message, fmtMsgStringNotEndsWith, possiblesStr)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringEndsWith, possiblesStr)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringExactLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringExactLength) checkString(str string, vcx *ValidatorContext) bool {
	l := len(str)
	if c.UseRuneLen {
		l = len([]rune(str))
	}
	return l == c.Value
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringLength) checkString(str string, vcx *ValidatorContext) bool {
	l := len(str)
	if c.UseRuneLen {
		l = len([]rune(str))
	}
	return (l > c.Minimum || (!c.ExclusiveMin && l == c.Minimum)) &&
		(c.Maximum <= 0 || l < c.Maximum || (!c.ExclusiveMax && l == c.Maximum))
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

// StringLowercase constraint to check that a string has only lowercase letters
type StringLowercase struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringLowercase) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringLowercase) checkString(str string, vcx *ValidatorContext) bool {
	return str == strings.ToLower(str)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringLowercase) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgStringLowercase)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringMaxLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringMaxLength) checkString(str string, vcx *ValidatorContext) bool {
	l := len(str)
	if c.UseRuneLen {
		l = len([]rune(str))
	}
	return l < c.Value || (!c.ExclusiveMax && l == c.Value)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMaxLength) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgStringMaxLenExc, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringMaxLen, c.Value)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringMinLength) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringMinLength) checkString(str string, vcx *ValidatorContext) bool {
	l := len(str)
	if c.UseRuneLen {
		l = len([]rune(str))
	}
	return l > c.Value || (!c.ExclusiveMin && l == c.Value)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringMinLength) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgStringMinLenExc, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringMinLen, c.Value)
}

// StringNoControlCharacters constraint to check that a string does not contain any control characters (i.e. chars < 32)
type StringNoControlCharacters struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringNoControlCharacters) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringNoControlCharacters) checkString(str string, vcx *ValidatorContext) bool {
	for _, ch := range str {
		if ch < 32 {
			return false
		}
	}
	return true
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNoControlCharacters) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNoControlChars)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringNotBlank) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringNotBlank) checkString(str string, vcx *ValidatorContext) bool {
	return len(strings.Trim(str, " \t\n\r")) != 0
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotBlank) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNotBlankString)
}

// StringNotEmpty constraint to check that string value is not empty (i.e. not "")
type StringNotEmpty struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringNotEmpty) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringNotEmpty) checkString(str string, vcx *ValidatorContext) bool {
	return len(str) > 0
}

// GetMessage implements the Constraint.GetMessage
func (c *StringNotEmpty) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNotEmptyString)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringPattern) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringPattern) checkString(str string, vcx *ValidatorContext) bool {
	return c.Regexp.MatchString(str)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringPresetPattern) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if p, ok := presetsRegistry.get(c.Preset); ok {
			if !p.Check(str) {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.getMessage(vcx, p.GetMessage())
			}
		} else {
			vcx.CeaseFurtherIf(c.Stop)
			return false, vcx.TranslateFormat(fmtMsgUnknownPresetPattern, c.Preset)
		}
	} else if c.Strict {
		if p, ok := presetsRegistry.get(c.Preset); ok {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.getMessage(vcx, p.GetMessage())
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
	} else if p, ok := presetsRegistry.get(c.Preset); ok {
		if msg := p.GetMessage(); msg != "" {
			return obtainI18nContext(tcx).TranslateMessage(msg)
		}
	}
	return obtainI18nContext(tcx).TranslateMessage(msgValidPattern)
}

// StringStartsWith constraint to check that a string starts with a given prefix
type StringStartsWith struct {
	// the value to check that the string starts with
	Value string `v8n:"default"`
	// multiple additional values that the string may start with
	Values []string
	// whether the check is case-insensitive (by default, the check is case-sensitive)
	CaseInsensitive bool
	// whether the check is NOT-ed (i.e. checks that the string does not start with)
	Not bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringStartsWith) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringStartsWith) checkString(str string, vcx *ValidatorContext) bool {
	ckStr := caseInsensitive(str, c.CaseInsensitive)
	startsWith := c.Value != "" && strings.HasPrefix(ckStr, caseInsensitive(c.Value, c.CaseInsensitive))
	if !startsWith {
		for _, s := range c.Values {
			startsWith = s != "" && strings.HasPrefix(ckStr, caseInsensitive(s, c.CaseInsensitive))
			if startsWith {
				break
			}
		}
	}
	return (!c.Not && startsWith) || (c.Not && !startsWith)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringStartsWith) GetMessage(tcx I18nContext) string {
	possibles := make([]string, 0)
	if c.Value != "" {
		possibles = append(possibles, "'"+c.Value+"'")
	}
	for _, s := range c.Values {
		if s != "" {
			possibles = append(possibles, "'"+s+"'")
		}
	}
	possiblesStr := strings.Join(possibles, ",")
	if c.Not {
		return defaultMessage(tcx, c.Message, fmtMsgStringNotStartsWith, possiblesStr)
	}
	return defaultMessage(tcx, c.Message, fmtMsgStringStartsWith, possiblesStr)
}

// StringUppercase constraint to check that a string has only uppercase letters
type StringUppercase struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringUppercase) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringUppercase) checkString(str string, vcx *ValidatorContext) bool {
	return str == strings.ToUpper(str)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringUppercase) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgStringUppercase)
}

// StringValidJson constraint checks that a string is valid json
type StringValidJson struct {
	DisallowNullJson bool
	DisallowValue    bool
	DisallowArray    bool
	DisallowObject   bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidJson) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidJson) checkString(str string, vcx *ValidatorContext) bool {
	var v interface{}
	fail := true
	if err := json.Unmarshal([]byte(str), &v); err == nil {
		switch v.(type) {
		case map[string]interface{}:
			fail = c.DisallowObject
		case []interface{}:
			fail = c.DisallowArray
		case nil:
			fail = c.DisallowNullJson
		default:
			fail = c.DisallowValue
		}
	}
	return !fail
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidJson) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgStringValidJson)
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidToken) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidToken) checkString(str string, vcx *ValidatorContext) bool {
	lstr := strings.ToLower(str)
	for _, t := range c.Tokens {
		if str == t || (c.IgnoreCase && lstr == strings.ToLower(t)) {
			return true
		}
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidToken) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgValidToken, strings.Join(c.Tokens, "\",\""))
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
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidUnicodeNormalization) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return checkStringConstraint(v, vcx, c, c.Strict, c.Stop)
}

func (c *StringValidUnicodeNormalization) checkString(str string, vcx *ValidatorContext) bool {
	return c.Form.IsNormalString(str)
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

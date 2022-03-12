package valix

import (
	"golang.org/x/text/unicode/norm"
	"strings"
)

type baseNoMsg struct {
}

func (c baseNoMsg) GetMessage() string {
	return ""
}

// StringTrim constraint trims a string value
type StringTrim struct {
	// Cutset is the string containing codepoints to bve trimmed - or if blank, will default
	// to removing tabs and spaces
	Cutset string
	baseNoMsg
}

// Check implements Constraint.Check
func (c *StringTrim) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if c.Cutset == "" {
			vcx.SetCurrentValue(strings.Trim(str, " \t"))
		} else {
			vcx.SetCurrentValue(strings.Trim(str, c.Cutset))
		}
	}
	return true, c.GetMessage()
}

// StringNormalizeUnicode constraint sets the unicode normalization of a string value
type StringNormalizeUnicode struct {
	// Form is the normalization form required - i.e. norm.NFC, norm.NFKC, norm.NFD or norm.NFKD
	//
	// (from package "golang.org/x/text/unicode/norm")
	Form norm.Form
	baseNoMsg
}

// Check implements Constraint.Check
func (c *StringNormalizeUnicode) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		vcx.SetCurrentValue(c.Form.String(str))
	}
	return true, c.GetMessage()
}

// SetConditionFrom constraint is a utility constraint that can be used to set a condition in the
// ValidatorContext from string value of the property (to which this constraint is added)
//
// Note: It will only set a condition if the property value is a string!
type SetConditionFrom struct {
	baseNoMsg
}

// Check implements Constraint.Check
func (c *SetConditionFrom) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		vcx.SetCondition(str)
	}
	return true, c.GetMessage()
}

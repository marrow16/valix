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

// StringTrim trims a string value
type StringTrim struct {
	// Cutset is the string containing codepoints to bve trimmed - or if blank, will default
	// to removing tabs and spaces
	Cutset string
	baseNoMsg
}

// Check implements Constraint.Check
func (c *StringTrim) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if c.Cutset == "" {
			vcx.SetCurrentValue(strings.Trim(str, " \t"))
		} else {
			vcx.SetCurrentValue(strings.Trim(str, c.Cutset))
		}
	}
	return true, c.GetMessage()
}

// StringNormalizeUnicode sets the unicode normalization of a string value
type StringNormalizeUnicode struct {
	// Form is the normalization form required - i.e. norm.NFC, norm.NFKC, norm.NFD or norm.NFKD
	//
	// (from package "golang.org/x/text/unicode/norm")
	Form norm.Form
	baseNoMsg
}

// Check implements Constraint.Check
func (c *StringNormalizeUnicode) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		vcx.SetCurrentValue(c.Form.String(str))
	}
	return true, c.GetMessage()
}

package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
	"unicode"
)

func TestCanCreateCustomConstraint(t *testing.T) {
	cc := NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage()
	}, "")
	require.NotNil(t, cc)
}

func TestCustomConstraintStoresMessage(t *testing.T) {
	const testMsg = "TEST MESSAGE"
	cc := NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage()
	}, testMsg)
	require.Equal(t, testMsg, cc.GetMessage())
}

func TestCustomConstraint(t *testing.T) {
	const testMsg = "Value must be greater than 'B'"
	validator := buildFooValidator(JsonString,
		NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, cc.GetMessage()
			}
			return true, ""
		}, testMsg), false)
	obj := jsonObject(`{
		"foo": "OK is greater than B"
	}`)

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = "A"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = "B"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = "Ba"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestConstraintSet(t *testing.T) {
	const msg = "String value length must be between 16 and 64 chars; must be letters (upper or lower), digits or underscores; must start with an uppercase letter"
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringTrim{},
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if len(str) == 0 || str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage()
					}
				}
				return true, ""
			}, ""),
			&StringCharacters{
				AllowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: msg,
	}
	require.Equal(t, msg, set.GetMessage())

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	// this should be ok (even though longer than 64 it gets trimmed)...
	obj["foo"] = " Abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ      "
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// some more not oks...
	obj["foo"] = "abcdefghijklmnopqrstuvwxyz" // not starts with capital
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "Abc." // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "        " // empty after trim
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
}

func TestConstraintSetNoMsg(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringTrim{},
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage()
					}
				}
				return true, ""
			}, "Must start with a capital"),
			&StringCharacters{
				AllowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: "",
	}
	// message should first sub-constraint non-empty message...
	require.Equal(t, messageNotEmptyString, set.GetMessage())

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 16, 64), violations[0].Message)

	// this should be ok (even though longer than 64 it gets trimmed)...
	obj["foo"] = " Abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ      "
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// some more not oks...
	obj["foo"] = "abcdefghijklmnopqrstuvwxyz" // not starts with capital
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Must start with a capital", violations[0].Message)
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 16, 64), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 16, 64), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)
	obj["foo"] = "        " // empty after trim
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)
}

func TestConstraintSetCeases(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringTrim{},
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if str[0] < 'A' || str[0] > 'Z' {
						vcx.CeaseFurther()
					}
				}
				return true, ""
			}, "Must start with a capital"),
			&StringCharacters{
				AllowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []*unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: "",
	}
	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "this does not start with capital"
	}`)

	// the custom constraint ceased further and passed - so validation passes...
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// this should be ok (even though longer than 64 it gets trimmed)...
	obj["foo"] = " Abcdefghijklmnopqrstuvwxyz0123456789_ABCDEFGHIJKLMNOPQRSTUVWXYZ      "
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// some more not oks...
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 16, 64), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 16, 64), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)
	obj["foo"] = "        " // empty after trim
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)
}

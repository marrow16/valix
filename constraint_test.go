package valix

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestCanCreateCustomConstraint(t *testing.T) {
	cc := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage(vcx)
	}, "")
	require.NotNil(t, cc)
}

func TestCustomConstraintStoresMessage(t *testing.T) {
	const testMsg = "TEST MESSAGE"
	cc := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage(vcx)
	}, testMsg)
	require.Equal(t, testMsg, cc.GetMessage(nil))
}

func TestCustomConstraint(t *testing.T) {
	const testMsg = "Value must be greater than 'B'"
	validator := buildFooValidator(JsonString,
		NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, cc.GetMessage(vcx)
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
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if len(str) == 0 || str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage(vcx)
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
	require.Equal(t, msg, set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

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
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage(vcx)
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
	// message should return first sub-constraint with non-empty message...
	require.Equal(t, msgNotEmptyString, set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmptyString, violations[0].Message)
}

func TestConstraintSetCeases(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{
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

	// some more not oks...
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmptyString, violations[0].Message)
}

func TestConstraintSetOneOf(t *testing.T) {
	constraint1 := &testConstraint{passes: false, stops: false}
	constraint2 := &testConstraint{passes: true, stops: false}
	set := &ConstraintSet{
		OneOf:       true,
		Constraints: Constraints{constraint1, constraint2},
	}
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 2), set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "anything"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	// the second passing doesn't get reached if first stops...
	constraint1.stops = true
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 2), violations[0].Message)

	constraint1.stops = false
	constraint1.msg = "first message"
	constraint2.msg = "second message"
	constraint2.passes = false
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "first message", violations[0].Message)
}

func TestConstraintSetDefaultMessage(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{&testConstraint{}},
	}
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultAllOf, 1), set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "anything"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultAllOf, 1), violations[0].Message)

	set.OneOf = true
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 1), set.GetMessage(nil))
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 1), violations[0].Message)
}

type testConstraint struct {
	passes bool
	msg    string
	stops  bool
}

func (c *testConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	vcx.CeaseFurtherIf(c.stops)
	return c.passes, c.msg
}
func (c *testConstraint) GetMessage(tcx I18nContext) string {
	return c.msg
}

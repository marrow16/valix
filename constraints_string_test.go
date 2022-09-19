package valix

import (
	"fmt"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestStringCharactersConstraint(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringCharacters{
			AllowRanges: []*unicode.RangeTable{
				unicode.Upper,
			},
		}, false)
	obj := jsonObject(`{
		"foo": "xxx"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)

	obj["foo"] = "ABC"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringCharactersConstraintWithDisallows(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringCharacters{
			AllowRanges: []*unicode.RangeTable{
				unicode.Upper,
			},
			DisallowRanges: []*unicode.RangeTable{
				unicode.Hex_Digit,
			},
		}, false)
	obj := jsonObject(`{
		"foo": "ABC"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)

	obj["foo"] = "GHI"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringCharactersConstraintWithPlanes(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringCharacters{
			AllowRanges: []*unicode.RangeTable{
				UnicodeBMP, UnicodeSMP,
			},
		}, false)
	obj := jsonObject(`{
		"foo": "\ud840\udc06 < surrogate representation of u+20006 - a Chinese char in SIP (Supplementary Ideographic Plane)"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)

	// now try again by allowing SIP (Supplementary Ideographic Plane)...
	validator = buildFooValidator(JsonString,
		&StringCharacters{
			AllowRanges: []*unicode.RangeTable{
				UnicodeBMP, UnicodeSMP, UnicodeSIP,
			},
		}, false)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringContains(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: "_XX_"}, false)
	obj := jsonObject(`{
		"foo": "test_XX_foo"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_xx_foo"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_CONTAINS"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringContains, "'_XX_'"), violations[0].Message)
}

func TestStringContainsCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: "_XX_", CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_xx_foo"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_XX_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_CONTAINS"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringContains, "'_XX_'"), violations[0].Message)
}

func TestStringContainsNot(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: "_XX_", Not: true}, false)
	obj := jsonObject(`{
		"foo": "test_XX_foo"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotContains, "'_XX_'"), violations[0].Message)

	obj["foo"] = "test_xx_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringContainsNotCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: "_XX_", Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_XX_foo"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotContains, "'_XX_'"), violations[0].Message)

	obj["foo"] = "test_xx_foo"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringContainsMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Values: []string{"_YY_", "_ZZ_"}}, false)
	obj := jsonObject(`{
		"foo": "test_YY_foo"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_ZZ_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_zz_foo"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_CONTAINS"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringContains, "'_YY_','_ZZ_'"), violations[0].Message)
}

func TestStringContainsMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Values: []string{"_YY_", "_ZZ_"}, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_yy_foo"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_YY_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_zz_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_CONTAINS"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringContains, "'_YY_','_ZZ_'"), violations[0].Message)
}

func TestStringContainsNotMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Values: []string{"_YY_", "_ZZ_"}, Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_ZZ_foo"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotContains, "'_YY_','_ZZ_'"), violations[0].Message)

	obj["foo"] = "test_zz_foo"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringContainsNotMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Values: []string{"_YY_", "_ZZ_"}, Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "foo_ZZ_test"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotContains, "'_YY_','_ZZ_'"), violations[0].Message)

	obj["foo"] = "foo_zz_test"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestStringContainsEmptyAlwaysFails(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: ""}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	// the message will be strange (garbage in, garbage out)
	require.Equal(t, fmt.Sprintf(fmtMsgStringContains, ""), violations[0].Message)
}

func TestStringContainsNotEmptyAlwaysPasses(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringContains{Value: "", Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
}

func TestStringEndsWith(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: "_XX"}, false)
	obj := jsonObject(`{
		"foo": "test_XX"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_xx"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_ENDS_WITH"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringEndsWith, "'_XX'"), violations[0].Message)
}

func TestStringEndsWithCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: "_XX", CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_xx"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_XX"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_ENDS_WITH"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringEndsWith, "'_XX'"), violations[0].Message)
}

func TestStringEndsWithNot(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: "_XX", Not: true}, false)
	obj := jsonObject(`{
		"foo": "test_XX"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotEndsWith, "'_XX'"), violations[0].Message)

	obj["foo"] = "test_xx"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringEndsWithNotCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: "_XX", Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_XX"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotEndsWith, "'_XX'"), violations[0].Message)

	obj["foo"] = "test_xx"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringEndsWithMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Values: []string{"_YY", "_ZZ"}}, false)
	obj := jsonObject(`{
		"foo": "test_YY"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_ZZ"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_zz"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_ENDS_WITH"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringEndsWith, "'_YY','_ZZ'"), violations[0].Message)
}

func TestStringEndsWithMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Values: []string{"_YY", "_ZZ"}, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test_yy"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_YY"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_zz"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_ENDS_WITH"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringEndsWith, "'_YY','_ZZ'"), violations[0].Message)
}

func TestStringEndsWithNotMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Values: []string{"_YY", "_ZZ"}, Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_ZZ"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotEndsWith, "'_YY','_ZZ'"), violations[0].Message)

	obj["foo"] = "test_zz"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringEndsWithNotMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Values: []string{"_YY", "_ZZ"}, Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test_ZZ"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotEndsWith, "'_YY','_ZZ'"), violations[0].Message)

	obj["foo"] = "test_zz"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestStringEndsWithEmptyAlwaysFails(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: ""}, false)
	obj := jsonObject(`{
		"foo": "XX_test"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	// the message will be strange (garbage in, garbage out)
	require.Equal(t, fmt.Sprintf(fmtMsgStringEndsWith, ""), violations[0].Message)
}

func TestStringEndsWithNotEmptyAlwaysPasses(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringEndsWith{Value: "", Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
}

func TestStringExactLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringExactLength{Value: 2}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringExactLen, 2), violations[0].Message)

	obj["foo"] = "Abc"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringExactLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringExactLengthWithRuneLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringExactLength{Value: 2, UseRuneLen: true}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringExactLen, 2), violations[0].Message)

	obj["foo"] = "Abc"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringExactLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// NB. "\ud840\udc06" is surrogate representation of u+20006
	obj = jsonObject(`{
		"foo": "A\ud840\udc06"
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// and again with rune length turned off...
	validator = buildFooValidator(JsonString,
		&StringExactLength{Value: 2, UseRuneLen: false}, false)
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestStringLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringLength{Minimum: 2, Maximum: 3}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = "Abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// and without maximum length...
	validator = buildFooValidator(JsonString,
		&StringLength{Minimum: 2}, false)
	obj = jsonObject(`{
		"foo": "A"
	}`)

	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinLen, 2), violations[0].Message)
}

func TestStringLengthMinOnlyExc(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringLength{Minimum: 3, ExclusiveMin: true}, false)
	obj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinLenExc, 3), violations[0].Message)
}

func TestStringLengthExc(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringLength{Minimum: 3, Maximum: 5, ExclusiveMin: true, ExclusiveMax: true}, false)
	obj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 3, tokenExclusive, 5, tokenExclusive), violations[0].Message)

	obj["foo"] = "Abcde"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 3, tokenExclusive, 5, tokenExclusive), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringLengthWithRuneLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringLength{Minimum: 1, Maximum: 1}, false)
	// NB. "\ud840\udc06" is surrogate representation of u+20006
	obj := jsonObject(`{
		"foo": "\ud840\udc06"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 1, tokenInclusive, 1, tokenInclusive), violations[0].Message)

	// now try again but using rune length (actual Unicode length)...
	validator = buildFooValidator(JsonString,
		&StringLength{Minimum: 1, Maximum: 1, UseRuneLen: true}, false)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringLowercase(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringLowercase{}, false)
	obj := jsonObject(`{
		"foo": "Abc\t\n\r"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgStringLowercase, violations[0].Message)

	obj["foo"] = "abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMaxLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringMaxLength{Value: 2}, false)
	obj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMaxLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMaxLengthExclusive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringMaxLength{Value: 2, ExclusiveMax: true}, false)
	obj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMaxLenExc, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "A"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMaxLengthWithRuneLength(t *testing.T) {
	vWithUnicode := buildFooValidator(JsonString,
		&StringMaxLength{Value: 1, UseRuneLen: true}, false)
	vWithoutUnicode := buildFooValidator(JsonString,
		&StringMaxLength{Value: 1, UseRuneLen: false}, false)
	// NB. "\ud834\udd22" is surrogate representation of u+1D122 (musical symbol F clef)
	obj := jsonObject(`{
		"foo": "\ud834\udd22"
	}`)

	ok, _ := vWithUnicode.Validate(obj)
	require.True(t, ok)
	ok, violations := vWithoutUnicode.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMaxLen, 1), violations[0].Message)
}

func TestStringMinLength(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringMinLength{Value: 2}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMinLengthExclusive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringMinLength{Value: 2, ExclusiveMin: true}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinLenExc, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "Abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMinLengthWithRuneLength(t *testing.T) {
	vWithUnicode := buildFooValidator(JsonString,
		&StringMinLength{Value: 2, UseRuneLen: true}, false)
	vWithoutUnicode := buildFooValidator(JsonString,
		&StringMinLength{Value: 2, UseRuneLen: false}, false)
	// NB. "\ud834\udd22" is surrogate representation of u+1D122 (musical symbol F clef)
	obj := jsonObject(`{
		"foo": "\ud834\udd22"
	}`)

	ok, violations := vWithUnicode.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinLen, 2), violations[0].Message)
	ok, _ = vWithoutUnicode.Validate(obj)
	require.True(t, ok)
}

func TestStringNoControlChars(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringNoControlCharacters{}, false)
	obj := jsonObject(`{
		"foo": "Abc\t\n\r"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNoControlChars, violations[0].Message)

	obj["foo"] = "Abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringNotBlank(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringNotBlank{}, false)
	obj := jsonObject(`{
		"foo": " \t\n\r "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotBlankString, violations[0].Message)

	obj["foo"] = " bar \t\r\n"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringNotEmpty(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringNotEmpty{}, false)
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmptyString, violations[0].Message)

	obj["foo"] = "bar"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringPattern(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringPattern{
			Regexp: *regexp.MustCompile("([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})"),
		}, false)
	obj := jsonObject(`{
		"foo": "xxx"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidPattern, violations[0].Message)

	obj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringPresetPattern(t *testing.T) {
	constraint := &StringPresetPattern{
		Preset: "ISBN",
	}
	require.Equal(t, msgPresetISBN, constraint.GetMessage(nil))

	validator := buildFooValidator(JsonAny, constraint, false)
	obj := jsonObject(`{
		"foo": "9780201616224"
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = "does not match pattern"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPresetISBN, violations[0].Message)

	obj["foo"] = "9780201616229" // bad check digit
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPresetISBN, violations[0].Message)

	// with overridden message...
	constraint.Message = "OVERRIDDEN"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "OVERRIDDEN", violations[0].Message)
	require.Equal(t, "OVERRIDDEN", constraint.GetMessage(nil))

	constraint.Preset = "unknown preset"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgUnknownPresetPattern, "unknown preset"), violations[0].Message)

	// defaulted message when preset message does not exist...
	constraint.Message = ""
	msg := constraint.GetMessage(nil)
	require.Equal(t, msgValidPattern, msg)

	// defaulted message when preset has no message...
	constraint.Message = ""
	presetsRegistry.register("fooey", &patternPreset{regex: regexp.MustCompile("zzz")})
	defer func() {
		presetsRegistry.reset()
	}()
	constraint.Preset = "fooey"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidPattern, violations[0].Message)
}

func TestStringStartsWith(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: "XX_"}, false)
	obj := jsonObject(`{
		"foo": "XX_test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "xx_test"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_START_WITH_XX"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringStartsWith, "'XX_'"), violations[0].Message)
}

func TestStringStartsWithCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: "XX_", CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "xx_test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "XX_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_START_WITH_XX"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringStartsWith, "'XX_'"), violations[0].Message)
}

func TestStringStartsWithNot(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: "XX_", Not: true}, false)
	obj := jsonObject(`{
		"foo": "XX_test"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotStartsWith, "'XX_'"), violations[0].Message)

	obj["foo"] = "xx_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringStartsWithNotCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: "XX_", Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "XX_test"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotStartsWith, "'XX_'"), violations[0].Message)

	obj["foo"] = "xx_test"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = "test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringStartsWithMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Values: []string{"YY_", "ZZ_"}}, false)
	obj := jsonObject(`{
		"foo": "YY_test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "ZZ_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "zz_test"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "NOT_START_WITH_XX"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringStartsWith, "'YY_','ZZ_'"), violations[0].Message)
}

func TestStringStartsWithMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Values: []string{"YY_", "ZZ_"}, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "yy_test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "YY_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "zz_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "NOT_START_WITH_XX"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringStartsWith, "'YY_','ZZ_'"), violations[0].Message)
}

func TestStringStartsWithNotMultiValue(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Values: []string{"YY_", "ZZ_"}, Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "ZZ_test"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotStartsWith, "'YY_','ZZ_'"), violations[0].Message)

	obj["foo"] = "zz_test"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringStartsWithNotMultiValueCaseInsensitive(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Values: []string{"YY_", "ZZ_"}, Not: true, CaseInsensitive: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "ZZ_test"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringNotStartsWith, "'YY_','ZZ_'"), violations[0].Message)

	obj["foo"] = "zz_test"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestStringStartsWithEmptyAlwaysFails(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: ""}, false)
	obj := jsonObject(`{
		"foo": "XX_test"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	// the message will be strange (garbage in, garbage out)
	require.Equal(t, fmt.Sprintf(fmtMsgStringStartsWith, ""), violations[0].Message)
}

func TestStringStartsWithNotEmptyAlwaysPasses(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringStartsWith{Value: "", Not: true}, false)
	obj := jsonObject(`{
		"foo": "test"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
}

func TestStringUppercase(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringUppercase{}, false)
	obj := jsonObject(`{
		"foo": "Abc\t\n\r"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgStringUppercase, violations[0].Message)

	obj["foo"] = "ABC"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringValidJson(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidJson{}

	ok, msg := c.Check("invalid json!", vcx)
	require.False(t, ok)
	require.Equal(t, msgStringValidJson, msg)

	ok, _ = c.Check("null", vcx)
	require.True(t, ok)
	c.DisallowNullJson = true
	ok, _ = c.Check("null", vcx)
	require.False(t, ok)

	ok, _ = c.Check("\"\"", vcx)
	require.True(t, ok)
	ok, _ = c.Check("true", vcx)
	require.True(t, ok)
	ok, _ = c.Check("false", vcx)
	require.True(t, ok)
	ok, _ = c.Check("0", vcx)
	require.True(t, ok)
	c.DisallowValue = true
	ok, _ = c.Check("\"\"", vcx)
	require.False(t, ok)
	ok, _ = c.Check("true", vcx)
	require.False(t, ok)
	ok, _ = c.Check("false", vcx)
	require.False(t, ok)
	ok, _ = c.Check("0", vcx)
	require.False(t, ok)

	ok, _ = c.Check("{\"foo\":true}", vcx)
	require.True(t, ok)
	c.DisallowObject = true
	ok, _ = c.Check("{\"foo\":true}", vcx)
	require.False(t, ok)

	ok, _ = c.Check("[true,0,\"foo\"]", vcx)
	require.True(t, ok)
	c.DisallowArray = true
	ok, _ = c.Check("[true,0,\"foo\"]", vcx)
	require.False(t, ok)
}

func TestStringValidToken(t *testing.T) {
	validTokens := []string{"AAA", "BBB", "CCC"}
	constraint := &StringValidToken{
		Tokens: validTokens,
	}
	validator := buildFooValidator(JsonAny, constraint, false)
	obj := jsonObject(`{
		"foo": "xxx"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValidToken, strings.Join(validTokens, "\",\"")), violations[0].Message)

	obj["foo"] = validTokens[0]
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with invalid type (constraint should ignore)...
	obj["foo"] = true
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with ignore case...
	obj["foo"] = strings.ToLower(validTokens[0])
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	constraint.IgnoreCase = true
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringValidUnicodeNormalization(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFC}, false)
	// NB. "\u0063\u0327" is 'c' followed by combining cedilla
	obj := jsonObject(`{
		"foo": "\u0063\u0327"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnicodeNormalizationNFC, violations[0].Message)

	obj["foo"] = "\u00e7" // u+00E7 is 'c' with cedilla
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFD}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnicodeNormalizationNFD, violations[0].Message)
}

func TestStringValidUnicodeNormalizationK(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFKC}, false)
	// NB. "\u0063\u0327" is 'c' followed by combining cedilla
	obj := jsonObject(`{
		"foo": "\u0063\u0327"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnicodeNormalizationNFKC, violations[0].Message)

	obj["foo"] = "\u00e7" // u+00E7 is 'c' with cedilla
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFKD}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnicodeNormalizationNFKD, violations[0].Message)
}

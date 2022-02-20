package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
	"regexp"
	"strings"
	"testing"
	"time"
	"unicode"
)

func TestStringNotEmpty(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringNotEmpty{}, false)
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)

	obj["foo"] = "bar"
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
	require.Equal(t, messageNotBlankString, violations[0].Message)

	obj["foo"] = " bar \t\r\n"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
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
	require.Equal(t, messageNoControlChars, violations[0].Message)

	obj["foo"] = "Abc"
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
	require.Equal(t, messageInvalidPattern, violations[0].Message)

	obj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringValidToken(t *testing.T) {
	validTokens := []string{"AAA", "BBB", "CCC"}
	validator := buildFooValidator(JsonAny,
		&StringValidToken{
			Tokens: validTokens,
		}, false)
	obj := jsonObject(`{
		"foo": "xxx"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageInvalidToken, strings.Join(validTokens, "\",\"")), violations[0].Message)

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
}

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
	require.Equal(t, messageInvalidCharacters, violations[0].Message)

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
	require.Equal(t, messageInvalidCharacters, violations[0].Message)

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
	require.Equal(t, messageInvalidCharacters, violations[0].Message)

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
	require.Equal(t, messageUnicodeNormalizationNFC, violations[0].Message)

	obj["foo"] = "\u00e7" // u+00E7 is 'c' with cedilla
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFD}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageUnicodeNormalizationNFD, violations[0].Message)
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
	require.Equal(t, messageUnicodeNormalizationNFKC, violations[0].Message)

	obj["foo"] = "\u00e7" // u+00E7 is 'c' with cedilla
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(JsonString,
		&StringValidUnicodeNormalization{Form: norm.NFKD}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageUnicodeNormalizationNFKD, violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageStringMinLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
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
	require.Equal(t, fmt.Sprintf(messageStringMinLen, 2), violations[0].Message)
	ok, _ = vWithoutUnicode.Validate(obj)
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
	require.Equal(t, fmt.Sprintf(messageStringMaxLen, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
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

	ok, violations := vWithUnicode.Validate(obj)
	require.True(t, ok)
	ok, violations = vWithoutUnicode.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMaxLen, 1), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 2, 3), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 2, 3), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(messageStringMinLen, 2), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageStringExactLen, 1), violations[0].Message)

	// now try again but using rune length (actual Unicode length)...
	validator = buildFooValidator(JsonString,
		&StringLength{Minimum: 1, Maximum: 1, UseRuneLen: true}, false)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthWithString(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&Length{Minimum: 2, Maximum: 3}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = "Abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// and without max...
	validator = buildFooValidator(JsonString,
		&Length{Minimum: 2}, false)
	obj = jsonObject(`{
		"foo": "A"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinLen, 2), violations[0].Message)
}

func TestLengthWithObject(t *testing.T) {
	validator := buildFooValidator(JsonObject,
		&Length{Minimum: 2, Maximum: 3}, false)
	obj := jsonObject(`{
		"foo": {
			"bar": null
		}
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
		"quz": nil,
	}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
	}
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthWithArray(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&Length{Minimum: 2, Maximum: 3}, false)
	obj := jsonObject(`{
		"foo": ["bar"]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz", "quz"}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz"}
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthConstraintWithSameMinMax(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&Length{Minimum: 2, Maximum: 2}, false)
	obj := jsonObject(`{
		"foo": ["bar"]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageExactLen, 2), violations[0].Message)
}

func TestPositive(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&Positive{}, false)
	obj := jsonObject(`{
		"foo": -1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)

	obj["foo"] = 1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)
	obj["foo"] = json.Number("1")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestPositiveOrZero(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&PositiveOrZero{}, false)
	obj := jsonObject(`{
		"foo": -1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)

	obj["foo"] = 0.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int..
	obj["foo"] = -1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegative(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&Negative{}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)

	obj["foo"] = -1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)
	obj["foo"] = -1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)
	obj["foo"] = json.Number("-1")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegativeOrZero(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&NegativeOrZero{}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)

	obj["foo"] = 0.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = -1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMinimum(t *testing.T) {
	testMsg := "Must not be less than 2"
	validator := buildFooValidator(JsonNumber,
		&Minimum{Value: 2, Message: testMsg}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("2")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMaximum(t *testing.T) {
	testMsg := "Must not be greater than 2"
	validator := buildFooValidator(JsonNumber,
		&Maximum{Value: 2, Message: testMsg}, false)
	obj := jsonObject(`{
		"foo": 3
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 3
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("3")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("2")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestRange(t *testing.T) {
	testMsg := "Must be between 2 and 3 (inclusive)"
	validator := buildFooValidator(JsonNumber,
		&Range{Minimum: 2, Maximum: 3, Message: testMsg}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 3.0001
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 4
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("4")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("2")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestArrayOf(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&ArrayOf{Type: "string", AllowNullElement: false}, false)
	obj := jsonObject(`{
		"foo": ["ok", false]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementType, JsonString), violations[0].Message)

	obj = jsonObject(`{
		"foo": [null]
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementType, JsonString), violations[0].Message)

	obj = jsonObject(`{
		"foo": ["ok", "ok2"]
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(JsonArray,
		&ArrayOf{Type: "string", AllowNullElement: true}, false)
	obj = jsonObject(`{
		"foo": [1, "ok2"]
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementTypeOrNull, JsonString), violations[0].Message)
	obj = jsonObject(`{
		"foo": [null, "ok2"]
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringValidUuid(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidUuid{MinVersion: 4}, false)
	obj := jsonObject(`{
		"foo": "not a uuid"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidMinVersion, 4), violations[0].Message)

	obj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "db15398d-328f-3d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidMinVersion, 4), violations[0].Message)

	validator = buildFooValidator(JsonString,
		&StringValidUuid{SpecificVersion: 4}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidCorrectVer, 4), violations[0].Message)

	validator = buildFooValidator(JsonString,
		&StringValidUuid{}, false)
	obj = jsonObject(`{
		"foo": "not a uuid"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValidUuid), violations[0].Message)
}

func TestStringValidISODatetime(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidISODatetime{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02T18:19:20.12345+01:00 but not a datetime with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODatetime+messageDatetimeFormatFull, violations[0].Message)

	obj["foo"] = "2022-02-02T18:19:20.123+01:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-13-02T18:19:20.12345+01:00"
	//                  ^^ 13th month?
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODatetime+messageDatetimeFormatFull, violations[0].Message)
}

func TestStringValidISODatetimeWithDifferentSettings(t *testing.T) {
	vFull := buildFooValidator(JsonString,
		&StringValidISODatetime{}, false)
	vNoNano := buildFooValidator(JsonString,
		&StringValidISODatetime{NoMillis: true}, false)
	vNoOffs := buildFooValidator(JsonString,
		&StringValidISODatetime{NoOffset: true}, false)
	vMin := buildFooValidator(JsonString,
		&StringValidISODatetime{NoOffset: true, NoMillis: true}, false)

	testCases := []struct {
		testValue  string
		okFull     bool
		okNoMillis bool
		okNoOffs   bool
		okMin      bool
	}{
		{testValue: "", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20", okFull: true, okNoMillis: true, okNoOffs: true, okMin: true},
		{testValue: "2022-02-02T18:19:20.12345+01:00", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.1234567890+01:00", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345-01:00", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345Z", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345Z+01:00", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20+01:00", okFull: true, okNoMillis: true, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.123456", okFull: true, okNoMillis: false, okNoOffs: true, okMin: false},
		{testValue: "2022-02-02T18:19:20Z", okFull: true, okNoMillis: true, okNoOffs: false, okMin: false},
		// bad dates/times...
		{testValue: "2022-13-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-41T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T25:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:60:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:60.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		// too many digits in various places...
		{testValue: "20222-01-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-012-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-012T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T189:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:190:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:201.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:20.1234567891Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Value: %s", testCase.testValue), func(t *testing.T) {
			obj["foo"] = testCase.testValue
			ok, violations := vFull.Validate(obj)
			require.Equal(t, testCase.okFull, ok)
			if !testCase.okFull {
				require.Equal(t, 1, len(violations))
				require.Equal(t, messageValidISODatetime+messageDatetimeFormatFull, violations[0].Message)
			}
			ok, violations = vNoNano.Validate(obj)
			require.Equal(t, testCase.okNoMillis, ok)
			if !testCase.okNoMillis {
				require.Equal(t, 1, len(violations))
				require.Equal(t, messageValidISODatetime+messageDatetimeFormatNoMillis, violations[0].Message)
			}
			ok, violations = vNoOffs.Validate(obj)
			require.Equal(t, testCase.okNoOffs, ok)
			if !testCase.okNoOffs {
				require.Equal(t, 1, len(violations))
				require.Equal(t, messageValidISODatetime+messageDatetimeFormatNoOffs, violations[0].Message)
			}
			ok, violations = vMin.Validate(obj)
			require.Equal(t, testCase.okMin, ok)
			if !testCase.okMin {
				require.Equal(t, 1, len(violations))
				require.Equal(t, messageValidISODatetime+messageDatetimeFormatMin, violations[0].Message)
			}
		})
	}
}

func TestStringValidISODate(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidISODate{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02 but not a date with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODate, violations[0].Message)

	obj["foo"] = "2022-02-02"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-13-02"
	//                  ^^ 13th month?
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODate, violations[0].Message)

	// should also fail with time specified...
	obj["foo"] = "2022-13-02T18:19:20"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODate, violations[0].Message)
}

var variousDatetimeFormats = []string{
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999Z",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func TestDatetimeFuture(t *testing.T) {
	pastTime := time.Now().Add(0 - (5 * time.Minute))
	validator := buildFooValidator(JsonAny,
		&DatetimeFuture{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFuture, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = pastTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFuture, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = pastTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, messageDatetimeFuture, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFuture, violations[0].Message)
}

func TestDatetimeFutureOrPresent(t *testing.T) {
	pastTime := time.Now().Add(0 - (5 * time.Minute))
	validator := buildFooValidator(JsonAny,
		&DatetimeFutureOrPresent{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFutureOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = pastTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFutureOrPresent, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = pastTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, messageDatetimeFutureOrPresent, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFutureOrPresent, violations[0].Message)
}

func TestDatetimePast(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	validator := buildFooValidator(JsonAny,
		&DatetimePast{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePast, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = futureTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePast, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = futureTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, messageDatetimePast, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePast, violations[0].Message)
}

func TestDatetimePastOrPresent(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	validator := buildFooValidator(JsonAny,
		&DatetimePastOrPresent{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePastOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = futureTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePastOrPresent, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = futureTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, messageDatetimePastOrPresent, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePastOrPresent, violations[0].Message)
}

func TestStringValidCardNumberConstraint(t *testing.T) {
	validator := buildFooValidator(JsonAny,
		&StringValidCardNumber{}, false)
	obj := map[string]interface{}{
		"foo": "",
	}
	testCardNumbers := map[string]bool{
		// valid VISA...
		"4902498374064506":    true,
		"4556494687321500":    true,
		"4969975508508776718": true,
		// valid MasterCard...
		"5385096580406173": true,
		"5321051022936318": true,
		"5316711656334695": true,
		// valid American Express (AMEX)...
		"344349836405221": true,
		"370579317661267": true,
		"379499960274519": true,
		// valid Discover...
		"6011816780702703":    true,
		"6011803416863109":    true,
		"6011658115940081575": true,
		// valid JCB...
		"3537599036021738":    true,
		"3535975622588649":    true,
		"3538751375142304859": true,
		// valid Diners Club (North America)...
		"5449933541584900": true,
		"5424932328289252": true,
		"5443511725058507": true,
		// valid Diners Club (Carte Blanche)...
		"30533389154463": true,
		"30225631860175": true,
		"30132515376577": true,
		// valid Diners Club (International)...
		"36023499511897": true,
		"36895205164933": true,
		"36614132300415": true,
		// valid Maestro...
		"5020799867464796": true,
		"5893155499331362": true,
		"6763838557832695": true,
		// valid Visa Electron...
		"4917958531215104": true,
		"4913912408530396": true,
		"4917079458677141": true,
		// valid InstaPayment...
		"6387294734923401": true,
		"6382441564848878": true,
		"6371830528023664": true,
		// valid all zeroes...
		"00000000000000":      true,
		"000000000000000":     true,
		"0000000000000000":    true,
		"00000000000000000":   true,
		"000000000000000000":  true,
		"0000000000000000000": true,
		// invalid all zeroes...
		"0000000000000":        false,
		"00000000000000000000": false,
		// invalids...
		"1234567890123":        false, // too short
		"12345678901234567890": false, // too long
		"4902498374064505":     false,
		"4969975508508776717":  false,
		"5385096580406171":     false,
		"344349836405220":      false,
		"6011816780702709":     false,
		"6011658115940081576":  false,
		"3537599036021730":     false,
		"3538751375142304850":  false,
		"5449933541584901":     false,
		"30533389154466":       false,
		"36023499511898":       false,
		"5020799867464791":     false,
		"4917958531215105":     false,
		"6387294734923400":     false,
		// invalid bad digits...
		"12345678901234x": false,
	}
	for ccn, expect := range testCardNumbers {
		t.Run(fmt.Sprintf("CardNumber:\"%s\"", ccn), func(t *testing.T) {
			obj["foo"] = ccn
			ok, _ := validator.Validate(obj)
			require.Equal(t, expect, ok)
		})
	}
}

func buildFooValidator(propertyType JsonType, constraint Constraint, notNull bool) *Validator {
	return &Validator{
		Properties: Properties{
			"foo": {
				Type:        propertyType,
				NotNull:     notNull,
				Mandatory:   true,
				Constraints: Constraints{constraint},
			},
		},
	}
}

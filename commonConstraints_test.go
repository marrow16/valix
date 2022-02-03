package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestStringNotEmpty(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNotEmptyConstraint{}, false)
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotEmptyString, violations[0].Message)

	obj["foo"] = "bar"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringNotBlank(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNotBlankConstraint{}, false)
	obj := jsonObject(`{
		"foo": " \t\n\r "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNotBlankString, violations[0].Message)

	obj["foo"] = " bar \t\r\n"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringNoControlChars(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNoControlCharsConstraint{}, false)
	obj := jsonObject(`{
		"foo": "Abc\t\n\r"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNoControlChars, violations[0].Message)

	obj["foo"] = "Abc"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringPattern(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringPatternConstraint{
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMinLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringMinLengthConstraint{Value: 2}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageAtLeast, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringMaxLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringMaxLengthConstraint{Value: 2}, false)
	obj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageNotMore, 2), violations[0].Message)

	obj["foo"] = "Ab"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringLengthConstraint{Minimum: 2, Maximum: 3}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// and without maximum length...
	validator = buildFooValidator(PropertyType.String,
		&StringLengthConstraint{Minimum: 2}, false)
	obj = jsonObject(`{
		"foo": "A"
	}`)

	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageAtLeast, 2), violations[0].Message)
}

func TestLengthWithString(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// and without max...
	validator = buildFooValidator(PropertyType.String,
		&LengthConstraint{Minimum: 2}, false)
	obj = jsonObject(`{
		"foo": "A"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageAtLeast, 2), violations[0].Message)
}

func TestLengthWithObject(t *testing.T) {
	validator := buildFooValidator(PropertyType.Object,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthWithArray(t *testing.T) {
	validator := buildFooValidator(PropertyType.Array,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestPositive(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&PositiveConstraint{}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositive, violations[0].Message)
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestPositiveOrZero(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&PositiveOrZeroConstraint{}, false)
	obj := jsonObject(`{
		"foo": -1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 1.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int..
	obj["foo"] = -1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messagePositiveOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegative(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&NegativeConstraint{}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)
	obj["foo"] = -1
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegative, violations[0].Message)
	obj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegativeOrZero(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&NegativeOrZeroConstraint{}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = -1.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageNegativeOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestMinimum(t *testing.T) {
	testMsg := "Must not be less than 2"
	validator := buildFooValidator(PropertyType.Number,
		&MinimumConstraint{Value: 2, Message: testMsg}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 2.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestMaximum(t *testing.T) {
	testMsg := "Must not be greater than 2"
	validator := buildFooValidator(PropertyType.Number,
		&MaximumConstraint{Value: 2, Message: testMsg}, false)
	obj := jsonObject(`{
		"foo": 3
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = 2.0
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 3
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("3")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestRange(t *testing.T) {
	testMsg := "Must be between 2 and 3 (inclusive)"
	validator := buildFooValidator(PropertyType.Number,
		&RangeConstraint{Minimum: 2, Maximum: 3, Message: testMsg}, false)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
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
	ok, violations = validator.Validate(obj)
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
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestArrayOf(t *testing.T) {
	validator := buildFooValidator(PropertyType.Array,
		&ArrayOfConstraint{Type: PropertyType.String, AllowNullElement: false}, false)
	obj := jsonObject(`{
		"foo": ["ok", false]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementType, PropertyType.String), violations[0].Message)

	obj = jsonObject(`{
		"foo": [null]
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementType, PropertyType.String), violations[0].Message)

	obj = jsonObject(`{
		"foo": ["ok", "ok2"]
	}`)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	validator = buildFooValidator(PropertyType.Array,
		&ArrayOfConstraint{Type: PropertyType.String, AllowNullElement: true}, false)
	obj = jsonObject(`{
		"foo": [1, "ok2"]
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementTypeOrNull, PropertyType.String), violations[0].Message)
	obj = jsonObject(`{
		"foo": [null, "ok2"]
	}`)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

func TestStringValidUuid(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringValidUuidConstraint{MinVersion: 4}, false)
	obj := jsonObject(`{
		"foo": "not a uuid"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidMinVersion, 4), violations[0].Message)

	obj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = "db15398d-328f-3d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidMinVersion, 4), violations[0].Message)

	validator = buildFooValidator(PropertyType.String,
		&StringValidUuidConstraint{SpecificVersion: 4}, false)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidCorrectVer, 4), violations[0].Message)

	validator = buildFooValidator(PropertyType.String,
		&StringValidUuidConstraint{}, false)
	obj = jsonObject(`{
		"foo": "not a uuid"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValidUuid), violations[0].Message)
}

func TestStringValidISODatetime(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringValidISODatetimeConstraint{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02T18:19:20.12345+01:00 but not a datetime with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODatetime+messageDatetimeFormatFull, violations[0].Message)

	obj["foo"] = "2022-02-02T18:19:20.123+01:00"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-13-02T18:19:20.12345+01:00"
	//                  ^^ 13th month?
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODatetime+messageDatetimeFormatFull, violations[0].Message)
}

func TestStringValidISODatetimeWithDifferentSettings(t *testing.T) {
	vFull := buildFooValidator(PropertyType.String,
		&StringValidISODatetimeConstraint{}, false)
	vNoNano := buildFooValidator(PropertyType.String,
		&StringValidISODatetimeConstraint{NoMillis: true}, false)
	vNoOffs := buildFooValidator(PropertyType.String,
		&StringValidISODatetimeConstraint{NoOffset: true}, false)
	vMin := buildFooValidator(PropertyType.String,
		&StringValidISODatetimeConstraint{NoOffset: true, NoMillis: true}, false)

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
	validator := buildFooValidator(PropertyType.String,
		&StringValidISODateConstraint{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02 but not a date with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValidISODate, violations[0].Message)

	obj["foo"] = "2022-02-02"
	ok, violations = validator.Validate(obj)
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
	validator := buildFooValidator("",
		&DatetimeFutureConstraint{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFuture, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, violations = validator.Validate(obj)
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
	validator := buildFooValidator("",
		&DatetimeFutureOrPresentConstraint{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimeFutureOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, violations = validator.Validate(obj)
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
	validator := buildFooValidator("",
		&DatetimePastConstraint{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePast, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, violations = validator.Validate(obj)
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
	validator := buildFooValidator("",
		&DatetimePastOrPresentConstraint{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageDatetimePastOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, violations = validator.Validate(obj)
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

func buildFooValidator(propertyType string, constraint Constraint, notNull bool) *Validator {
	return &Validator{
		Properties: Properties{
			"foo": {
				PropertyType: propertyType,
				NotNull:      notNull,
				Mandatory:    true,
				Constraints:  Constraints{constraint},
			},
		},
	}
}

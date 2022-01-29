package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"regexp"
	"strings"
	"testing"
)

func TestStringNotEmpty(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNotEmptyConstraint{}, false)
	jobj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNotEmptyString, violations[0].Message)

	jobj["foo"] = "bar"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringNotBlank(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNotBlankConstraint{}, false)
	jobj := jsonObject(`{
		"foo": " \t\n\r "
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNotBlankString, violations[0].Message)

	jobj["foo"] = " bar \t\r\n"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringNoControlChars(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringNoControlCharsConstraint{}, false)
	jobj := jsonObject(`{
		"foo": "Abc\t\n\r"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNoControlChars, violations[0].Message)

	jobj["foo"] = "Abc"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringPattern(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringPatternConstraint{
			Regexp: *regexp.MustCompile("([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})"),
		}, false)
	jobj := jsonObject(`{
		"foo": "xxx"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageInvalidPattern, violations[0].Message)

	jobj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringMinLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringMinLengthConstraint{Value: 2}, false)
	jobj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueAtLeast, 2, "characters"), violations[0].Message)

	jobj["foo"] = "Ab"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringMaxLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringMaxLengthConstraint{Value: 2}, false)
	jobj := jsonObject(`{
		"foo": "Abc"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueNotMore, 2, "characters"), violations[0].Message)

	jobj["foo"] = "Ab"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringLength(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringLengthConstraint{Minimum: 2, Maximum: 3}, false)
	jobj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueAtLeast, 2, "characters"), violations[0].Message)

	jobj["foo"] = "Abcd"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueNotMore, 3, "characters"), violations[0].Message)

	jobj["foo"] = "Abc"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestLengthWithString(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
	jobj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueAtLeast, 2, "characters"), violations[0].Message)

	jobj["foo"] = "Abcd"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueNotMore, 3, "characters"), violations[0].Message)

	jobj["foo"] = "Abc"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestLengthWithObject(t *testing.T) {
	validator := buildFooValidator(PropertyType.Object,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
	jobj := jsonObject(`{
		"foo": {
			"bar": null
		}
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueAtLeast, 2, "properties"), violations[0].Message)

	jobj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
		"quz": nil,
	}
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueNotMore, 3, "properties"), violations[0].Message)

	jobj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
	}
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestLengthWithArray(t *testing.T) {
	validator := buildFooValidator(PropertyType.Array,
		&LengthConstraint{Minimum: 2, Maximum: 3}, false)
	jobj := jsonObject(`{
		"foo": ["bar"]
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueAtLeast, 2, "elements"), violations[0].Message)

	jobj["foo"] = []interface{}{"foo", "bar", "baz", "quz"}
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueNotMore, 3, "elements"), violations[0].Message)

	jobj["foo"] = []interface{}{"foo", "bar", "baz"}
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestPositive(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&PositiveConstraint{}, false)
	jobj := jsonObject(`{
		"foo": -1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositive, violations[0].Message)

	jobj["foo"] = 0.0
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositive, violations[0].Message)

	jobj["foo"] = 1.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 0
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositive, violations[0].Message)
	jobj["foo"] = 1
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("0")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositive, violations[0].Message)
	jobj["foo"] = json.Number("1")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestPositiveOrZero(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&PositiveOrZeroConstraint{}, false)
	jobj := jsonObject(`{
		"foo": -1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositiveOrZero, violations[0].Message)

	jobj["foo"] = 0.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = 1.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int..
	jobj["foo"] = -1
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositiveOrZero, violations[0].Message)
	jobj["foo"] = 0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValuePositiveOrZero, violations[0].Message)
	jobj["foo"] = json.Number("0")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestNegative(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&NegativeConstraint{}, false)
	jobj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegative, violations[0].Message)

	jobj["foo"] = 0.0
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegative, violations[0].Message)

	jobj["foo"] = -1.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 0
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegative, violations[0].Message)
	jobj["foo"] = -1
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("0")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegative, violations[0].Message)
	jobj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestNegativeOrZero(t *testing.T) {
	validator := buildFooValidator(PropertyType.Number,
		&NegativeOrZeroConstraint{}, false)
	jobj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegativeOrZero, violations[0].Message)

	jobj["foo"] = 0.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = -1.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 1
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegativeOrZero, violations[0].Message)
	jobj["foo"] = 0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("1")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNegativeOrZero, violations[0].Message)
	jobj["foo"] = json.Number("0")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestMinimum(t *testing.T) {
	testMsg := "Must not be less than 2"
	validator := buildFooValidator(PropertyType.Number,
		&MinimumConstraint{Value: 2, Message: testMsg}, false)
	jobj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	jobj["foo"] = 2.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 1
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = 2
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("1")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = json.Number("2")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestMaximum(t *testing.T) {
	testMsg := "Must not be greater than 2"
	validator := buildFooValidator(PropertyType.Number,
		&MaximumConstraint{Value: 2, Message: testMsg}, false)
	jobj := jsonObject(`{
		"foo": 3
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	jobj["foo"] = 2.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 3
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = 2
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("3")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = json.Number("2")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestRange(t *testing.T) {
	testMsg := "Must be between 2 and 3 (inclusive)"
	validator := buildFooValidator(PropertyType.Number,
		&RangeConstraint{Minimum: 2, Maximum: 3, Message: testMsg}, false)
	jobj := jsonObject(`{
		"foo": 1
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	jobj["foo"] = 3.0001
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	jobj["foo"] = 2.0
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	jobj["foo"] = nil
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with int...
	jobj["foo"] = 4
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = 1
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = 2
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)

	// test with json number...
	jobj["foo"] = json.Number("4")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = json.Number("1")
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	jobj["foo"] = json.Number("2")
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestArrayOf(t *testing.T) {
	validator := buildFooValidator(PropertyType.Array,
		&ArrayOfConstraint{Type: PropertyType.String, AllowNullElement: false}, false)
	jobj := jsonObject(`{
		"foo": ["ok", false]
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementType, PropertyType.String, 1), violations[0].Message)

	jobj = jsonObject(`{
		"foo": [null]
	}`)
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageArrayElementNull, 0), violations[0].Message)

	jobj = jsonObject(`{
		"foo": ["ok", "ok2"]
	}`)
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func TestStringValidUuid(t *testing.T) {
	validator := buildFooValidator(PropertyType.String,
		&StringValidUuidConstraint{MinVersion: 4}, false)
	jobj := jsonObject(`{
		"foo": "not a uuid"
	}`)

	ok, violations := validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, messageValueNotValidUuid, violations[0].Message)

	jobj["foo"] = "db15398d-328f-4d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
	jobj["foo"] = "db15398d-328f-3d16-be2f-f38e8f2d0a79"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidMinVersion, 4), violations[0].Message)

	validator = buildFooValidator(PropertyType.String,
		&StringValidUuidConstraint{SpecificVersion: 4}, false)
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageUuidIncorrectVer, 4), violations[0].Message)
}

func TestCustomConstraint(t *testing.T) {
	msg := "Value must be greater than 'B'"
	validator := buildFooValidator(PropertyType.String,
		CustomConstraint(func(value interface{}, ctx *Context) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, msg
			}
			return true, ""
		}), false)
	jobj := jsonObject(`{
		"foo": "OK is greater than B"
	}`)

	ok, violations := validator.Validate(jobj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	jobj["foo"] = "A"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	jobj["foo"] = "B"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	jobj["foo"] = "Ba"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

func buildFooValidator(propertyType string, constraint Constraint, notNull bool) *Validator {
	return &Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: propertyType,
				NotNull:      notNull,
				Mandatory:    true,
				Constraints:  []Constraint{constraint},
			},
		},
	}
}

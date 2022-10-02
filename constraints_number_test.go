package valix

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

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

func TestMaximum_Strict(t *testing.T) {
	c := &Maximum{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestMaximumExc(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&Maximum{Value: 2, ExclusiveMax: true}, false)
	obj := jsonObject(`{
		"foo": 2
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)

	obj["foo"] = 1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)
	obj["foo"] = json.Number("1")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMaximumInt(t *testing.T) {
	testMsg := "Must not be greater than 2"
	validator := buildFooValidator(JsonNumber,
		&MaximumInt{Value: 2, Message: testMsg}, false)
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

func TestMaximumInt_Strict(t *testing.T) {
	c := &MaximumInt{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestMaximumIntExc(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&MaximumInt{Value: 2, ExclusiveMax: true}, false)
	obj := jsonObject(`{
		"foo": 2
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)

	obj["foo"] = 1.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)
	obj["foo"] = json.Number("1")
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

func TestMinimum_Strict(t *testing.T) {
	c := &Minimum{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestMinimumExc(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&Minimum{Value: 2, ExclusiveMin: true}, false)
	obj := jsonObject(`{
		"foo": 2
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)

	obj["foo"] = 3.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)
	obj["foo"] = json.Number("3")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMinimumInt(t *testing.T) {
	testMsg := "Must not be less than 2"
	validator := buildFooValidator(JsonNumber,
		&MinimumInt{Value: 2, Message: testMsg}, false)
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

func TestMinimumInt_Strict(t *testing.T) {
	c := &MinimumInt{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestMinimumIntExc(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&MinimumInt{Value: 2, ExclusiveMin: true}, false)
	obj := jsonObject(`{
		"foo": 2
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)

	obj["foo"] = 3.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with int...
	obj["foo"] = 2
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("2")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)
	obj["foo"] = json.Number("3")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMultipleOf(t *testing.T) {
	validator := buildFooValidator(JsonNumber,
		&MultipleOf{Value: 5}, false)
	obj := jsonObject(`{
		"foo": 5
	}`)

	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 10
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("15")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 5.000000000001
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMultipleOf, 5), violations[0].Message)

	obj["foo"] = 6
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMultipleOf, 5), violations[0].Message)

	obj["foo"] = json.Number("16")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMultipleOf, 5), violations[0].Message)

	obj["foo"] = json.Number("20.00000000")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestMultipleOf_Strict(t *testing.T) {
	c := &MultipleOf{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
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
	require.Equal(t, msgNegative, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNegative, violations[0].Message)

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
	require.Equal(t, msgNegative, violations[0].Message)
	obj["foo"] = -1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNegative, violations[0].Message)
	obj["foo"] = json.Number("-1")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegative_Strict(t *testing.T) {
	c := &Negative{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
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
	require.Equal(t, msgNegativeOrZero, violations[0].Message)

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
	require.Equal(t, msgNegativeOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNegativeOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestNegativeOrZero_Strict(t *testing.T) {
	c := &NegativeOrZero{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
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
	require.Equal(t, msgPositive, violations[0].Message)

	obj["foo"] = 0.0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPositive, violations[0].Message)

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
	require.Equal(t, msgPositive, violations[0].Message)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("0")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPositive, violations[0].Message)
	obj["foo"] = json.Number("1")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestPositive_Strict(t *testing.T) {
	c := &Positive{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
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
	require.Equal(t, msgPositiveOrZero, violations[0].Message)

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
	require.Equal(t, msgPositiveOrZero, violations[0].Message)
	obj["foo"] = 0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// test with json number...
	obj["foo"] = json.Number("-1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPositiveOrZero, violations[0].Message)
	obj["foo"] = json.Number("0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestPositiveOrZero_Strict(t *testing.T) {
	c := &PositiveOrZero{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
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

func TestRange_Strict(t *testing.T) {
	c := &Range{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestRangeInt(t *testing.T) {
	testMsg := "Must be between 2 and 3 (inclusive)"
	validator := buildFooValidator(JsonNumber,
		&RangeInt{Minimum: 2, Maximum: 3, Message: testMsg}, false)
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

func TestRangeInt_Strict(t *testing.T) {
	c := &RangeInt{}
	validator := buildFooValidator(JsonAny, c, false)
	obj := map[string]interface{}{
		"foo": "not a number",
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	c.Strict = true
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

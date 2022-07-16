package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestEqualsOtherAndNotEqualsOther(t *testing.T) {
	vEquals := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&EqualsOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	vNotEquals := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&NotEqualsOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 1,
		"bar": 1
	}`)

	testCases := []struct {
		foo    interface{}
		bar    interface{}
		expect bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			nil,
			true,
			false,
		},
		{
			false,
			nil,
			false,
		},
		{
			true,
			true,
			true,
		},
		{
			true,
			false,
			false,
		},
		{
			1.0,
			1.0,
			true,
		},
		{
			json.Number("xxx"),
			json.Number("1"),
			false,
		},
		{
			json.Number("1"),
			json.Number("xxx"),
			false,
		},
		{
			json.Number("xxx"),
			1.0,
			false,
		},
		{
			json.Number("xxx"),
			1,
			false,
		},
		{
			1.0,
			1.0000001,
			false,
		},
		{
			1,
			1,
			true,
		},
		{
			1,
			2,
			false,
		},
		{
			1,
			1.0,
			true,
		},
		{
			1.0,
			1,
			true,
		},
		{
			1,
			json.Number("1.0"),
			true,
		},
		{
			json.Number("1.0"),
			1,
			true,
		},
		{
			1.0,
			json.Number("1.0"),
			true,
		},
		{
			json.Number("1.0"),
			1.0,
			true,
		},
		{
			json.Number("1"),
			json.Number("1.0"),
			true,
		},
		{
			json.Number("1.0"),
			json.Number("1.0000000001"),
			false,
		},
		{
			1,
			2,
			false,
		},
		{
			json.Number("1.0001"),
			json.Number("1.0001"),
			true,
		},
		{
			json.Number("1.0001"),
			json.Number("1.0002"),
			false,
		},
		{
			map[string]interface{}{},
			map[string]interface{}{},
			true,
		},
		{
			map[string]interface{}{"foo": nil},
			map[string]interface{}{"foo": nil},
			true,
		},
		{
			map[string]interface{}{"foo": true, "bar": false},
			map[string]interface{}{"bar": false, "foo": true},
			true,
		},
		{
			map[string]interface{}{"foo": 1},
			map[string]interface{}{"foo": 1.0},
			true,
		},
		{
			map[string]interface{}{"foo": 1.0},
			map[string]interface{}{"foo": 1},
			true,
		},
		{
			map[string]interface{}{"foo": json.Number("1")},
			map[string]interface{}{"foo": json.Number("1.0")},
			true,
		},
		{
			map[string]interface{}{"foo": 1},
			map[string]interface{}{"foo": json.Number("1.0")},
			true,
		},
		{
			map[string]interface{}{"foo": 1.0},
			map[string]interface{}{"foo": json.Number("1.0")},
			true,
		},
		{
			map[string]interface{}{"foo": json.Number("1.0")},
			map[string]interface{}{"foo": 1},
			true,
		},
		{
			map[string]interface{}{"foo": json.Number("1.0")},
			map[string]interface{}{"foo": 1.0},
			true,
		},
		{
			[]interface{}{},
			[]interface{}{},
			true,
		},
		{
			[]interface{}{"xxx"},
			[]interface{}{"xxx"},
			true,
		},
		{
			[]interface{}{nil},
			[]interface{}{nil},
			true,
		},
		{
			[]interface{}{1},
			[]interface{}{1},
			true,
		},
		{
			[]interface{}{1.0},
			[]interface{}{1.0},
			true,
		},
		{
			[]interface{}{json.Number("1.0")},
			[]interface{}{json.Number("1.0")},
			true,
		},
		{
			[]interface{}{json.Number("1")},
			[]interface{}{json.Number("1.0")},
			true,
		},
		{
			[]interface{}{json.Number("1.0")},
			[]interface{}{json.Number("1")},
			true,
		},
		{
			[]interface{}{1},
			[]interface{}{1.0},
			true,
		},
		{
			[]interface{}{1.0},
			[]interface{}{1},
			true,
		},
		{
			[]interface{}{1.0},
			[]interface{}{json.Number("1")},
			true,
		},
		{
			[]interface{}{json.Number("1")},
			[]interface{}{1.0},
			true,
		},
		{
			[]interface{}{1},
			[]interface{}{json.Number("1.0001")},
			false,
		},
		{
			[]interface{}{json.Number("1.0001")},
			[]interface{}{1},
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Equals[%d]", i+1), func(t *testing.T) {
			obj["foo"] = tc.foo
			obj["bar"] = tc.bar
			ok, violations := vEquals.Validate(obj)
			if tc.expect {
				require.True(t, ok)
			} else {
				require.False(t, ok)
				require.Equal(t, 1, len(violations))
				require.Equal(t, fmt.Sprintf(fmtMsgEqualsOther, "bar"), violations[0].Message)
			}
			ok, violations = vNotEquals.Validate(obj)
			if !tc.expect {
				require.True(t, ok)
			} else {
				require.False(t, ok)
				require.Equal(t, 1, len(violations))
				require.Equal(t, fmt.Sprintf(fmtMsgNotEqualsOther, "bar"), violations[0].Message)
			}
		})
	}
}

func TestGreaterThan(t *testing.T) {
	validator := buildFooValidator(JsonNumber, &GreaterThan{Value: 2}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)

	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestGreaterThanOrEqual(t *testing.T) {
	validator := buildFooValidator(JsonNumber, &GreaterThanOrEqual{Value: 2}, false)
	obj := jsonObject(`{
		"foo": 1
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGte, 2.0), violations[0].Message)

	obj["foo"] = 1.9999999999
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("1.9999999999")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLessThan(t *testing.T) {
	validator := buildFooValidator(JsonNumber, &LessThan{Value: 2}, false)
	obj := jsonObject(`{
		"foo": 3
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLt, 2.0), violations[0].Message)

	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 1.999999
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("2.000001")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("1.999999")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLessThanOrEqual(t *testing.T) {
	validator := buildFooValidator(JsonNumber, &LessThanOrEqual{Value: 2}, false)
	obj := jsonObject(`{
		"foo": 3
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLte, 2.0), violations[0].Message)

	obj["foo"] = 1.9999999999
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("1.9999999999")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestGreaterThanOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonNumber,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThanOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 1,
		"bar": 2
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)

	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 1
	obj["bar"] = 2
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["bar"] = 2.0
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["bar"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["bar"] = json.Number("not numeric")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["bar"] = "not numeric"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["bar"] = nil
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	delete(obj, "bar")
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestGreaterThanOrEqualOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonNumber,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThanOrEqualOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 1,
		"bar": 2
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGteOther, "bar"), violations[0].Message)

	obj["foo"] = 1.9999999999
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("1.9999999999")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLessThanOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonNumber,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&LessThanOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 3,
		"bar": 2
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)

	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 1.999999
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("2.000001")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = json.Number("1.999999")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLessThanOrEqualOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonNumber,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&LessThanOrEqualOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 3,
		"bar": 2
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLteOther, "bar"), violations[0].Message)

	obj["foo"] = 1.9999999999
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.0
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 2.000001
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = 3
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj["foo"] = 2
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("1.9999999999")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.0")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	obj["foo"] = json.Number("2.00001")
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
}

func TestCompareNonNumerics(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&GreaterThan{Value: 2},
				},
			},
		},
	}
	obj := jsonObject(`{"foo": 1}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)

	obj["foo"] = "not numeric"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)

	obj["foo"] = json.Number("not numeric")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGt, 2.0), violations[0].Message)
}

func TestDatetimeGreaterThan(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeGreaterThan{Value: "2022-04-01T00:00:00"}, false)
	obj := jsonObject(`{"foo": "2022-03-31T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:01.1234"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = false
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)
}

func TestDatetimeGreaterThanExcTime(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeGreaterThan{Value: "2022-04-01T00:00:00", ExcTime: true}, false)
	obj := jsonObject(`{"foo": "2022-03-31T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:01.1234"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "1234567890"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGt, "2022-04-01T00:00:00"), violations[0].Message)
}

func TestDatetimeGreaterThanOrEqual(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeGreaterThanOrEqual{Value: "2022-04-01T00:00:00"}, false)
	obj := jsonObject(`{"foo": "2022-03-31T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGte, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-04-01T00:00:01.1234"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeGreaterThanOrEqualExcTime(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeGreaterThanOrEqual{Value: "2022-04-01T00:00:00", ExcTime: true}, false)
	obj := jsonObject(`{"foo": "2022-03-31T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtGte, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-04-01T00:00:01.1234"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThan(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeLessThan{Value: "2022-04-01T00:00:00"}, false)
	obj := jsonObject(`{"foo": "2022-04-02T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-03-31T23:59:59.9999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = time.Now().Add(0 - (time.Hour * 24 * 365 * 100))
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanExcTime(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeLessThan{Value: "2022-04-01T00:00:00", ExcTime: true}, false)
	obj := jsonObject(`{"foo": "2022-04-02T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLt, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-03-31T23:59:59.9999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = time.Now().Add(0 - (time.Hour * 24 * 365 * 100))
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOrEqual(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeLessThanOrEqual{Value: "2022-04-01T00:00:00"}, false)
	obj := jsonObject(`{"foo": "2022-04-02T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLte, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-03-31T23:59:59.9999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, violations = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = time.Now().Add(0 - (time.Hour * 24 * 365 * 100))
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOrEqualExcTime(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeLessThanOrEqual{Value: "2022-04-01T00:00:00", ExcTime: true}, false)
	obj := jsonObject(`{"foo": "2022-04-02T00:00:00"}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgDtLte, "2022-04-01T00:00:00"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-03-31T23:59:59.9999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, violations = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = time.Now().Add(0 - (time.Hour * 24 * 365 * 100))
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeGreaterThanOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeGreaterThanOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-01T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00.001"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["bar"] = time.Now().Add(0 - time.Hour)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	delete(obj, "bar")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)
}

func TestDatetimeGreaterThanOtherExcTime(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeGreaterThanOther{PropertyName: "bar", ExcTime: true},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-01T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-04-01T00:00:00.001"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["bar"] = time.Now().Add(0 - (25 * time.Hour))
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	delete(obj, "bar")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "bar"), violations[0].Message)
}

func TestDatetimeGreaterThanOrEqualOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeGreaterThanOrEqualOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-02T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGteOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-04-02T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["bar"] = time.Now().Add(0 - time.Hour)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeGreaterThanOrEqualOtherExcTime(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeGreaterThanOrEqualOther{PropertyName: "bar", ExcTime: true},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-02T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgGteOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-04-02T00:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	now := time.Now()
	obj["foo"] = now
	obj["bar"] = now
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeLessThanOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-01T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-03-31T23:59:59.999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)

	obj["bar"] = time.Now().Add(0 - time.Hour)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)

	obj["bar"] = time.Now().Add(time.Hour)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOtherExcTime(t *testing.T) {
	constraint := &DatetimeLessThanOther{PropertyName: "bar", ExcTime: true}
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					constraint,
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00",
		"bar": "2022-04-01T12:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-03-31T23:59:59.999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	dt := time.Date(2022, 04, 01, 12, 0, 0, 0, time.UTC)
	obj["foo"] = dt
	obj["bar"] = dt.Add(time.Hour)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLtOther, "bar"), violations[0].Message)
	constraint.ExcTime = false
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOrEqualOther(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&DatetimeLessThanOrEqualOther{PropertyName: "bar"},
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T12:00:00",
		"bar": "2022-04-01T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLteOther, "bar"), violations[0].Message)

	obj["foo"] = "2022-03-31T23:59:59.999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLteOther, "bar"), violations[0].Message)

	obj["bar"] = time.Now().Add(0 - time.Hour)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLteOther, "bar"), violations[0].Message)

	obj["bar"] = time.Now().Add(time.Hour)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeLessThanOrEqualOtherExcTime(t *testing.T) {
	constraint := &DatetimeLessThanOrEqualOther{PropertyName: "bar", ExcTime: true}
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					constraint,
				},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T12:00:00",
		"bar": "2022-04-01T00:00:00"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-03-31T23:59:59.999"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Now()
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgLteOther, "bar"), violations[0].Message)

	obj["bar"] = time.Now().Add(0 - time.Millisecond)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeTolerance(t *testing.T) {
	// no fields - so same day assumed...
	validator := buildFooValidator(JsonAny, &DatetimeTolerance{Value: "2022-04-01T00:00:00"}, false)
	obj := jsonObject(`{
		"foo": "2022-04-01T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = "2022-04-02T00:00:00"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as 2022-04-01T00:00:00", violations[0].Message)

	obj["foo"] = time.Date(2022, 4, 2, 0, 0, 0, 0, time.UTC)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as 2022-04-01T00:00:00", violations[0].Message)

	obj["foo"] = time.Date(2022, 4, 1, 12, 0, 0, 0, time.UTC)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// with bad date...
	obj["foo"] = "not a date"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as 2022-04-01T00:00:00", violations[0].Message)

	// with null...
	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as 2022-04-01T00:00:00", violations[0].Message)

	// and with ignore nulls...
	validator = buildFooValidator(JsonAny, &DatetimeTolerance{Value: "2022-04-01T00:00:00", IgnoreNull: true}, false)
	obj = jsonObject(`{
		"foo": null
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeToleranceToNow(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeToleranceToNow{
		Duration: -5,
		Unit:     "year",
	}, false)
	today := time.Now()
	fiveYearsAgo, _ := shiftDatetimeByYears(&today, -5)
	obj := map[string]interface{}{
		"foo": fiveYearsAgo.Format("2006-01-02T15:04:05.999999999"),
	}
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	sixYearsAgo, _ := shiftDatetimeByYears(&today, -6)
	obj["foo"] = sixYearsAgo.Format("2006-01-02T15:04:05.999999999")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must not be more than 5 years before now", violations[0].Message)
}

func TestDatetimeToleranceToNowDefault(t *testing.T) {
	// no fields set - so same day check...
	validator := buildFooValidator(JsonAny, &DatetimeToleranceToNow{}, false)
	obj := map[string]interface{}{
		"foo": time.Now().Format("2006-01-02T15:04:05.999999999"),
	}
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = time.Now().Add(24 * time.Hour).Format("2006-01-02T15:04:05.999999999")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as now", violations[0].Message)

	obj["foo"] = time.Now().Add(0 - (24 * time.Hour))
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as now", violations[0].Message)

	obj["foo"] = time.Now()
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// with bad date...
	obj["foo"] = "not a date"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as now", violations[0].Message)

	// with null...
	obj["foo"] = nil
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as now", violations[0].Message)

	// and with ignore null...
	validator = buildFooValidator(JsonAny, &DatetimeToleranceToNow{IgnoreNull: true}, false)
	obj = jsonObject(`{
		"foo": null
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeToleranceToNowAsAgeCheck(t *testing.T) {
	validator := buildFooValidator(JsonAny, &DatetimeToleranceToNow{
		Duration: -18,
		Unit:     "year",
		MinCheck: true,
	}, false)
	now := time.Now().UTC()
	seventeenYearsAgo, _ := shiftDatetimeBy(&now, -17, "year")
	obj := map[string]interface{}{
		"foo": seventeenYearsAgo.Format("2006-01-02"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be at least 18 years before now", violations[0].Message)

	eighteenYearsAgo, _ := shiftDatetimeBy(&now, -18, "year")
	obj["foo"] = eighteenYearsAgo.Format("2006-01-02")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	nineteenYearsAgo, _ := shiftDatetimeBy(&now, -19, "year")
	obj["foo"] = nineteenYearsAgo.Format("2006-01-02")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeToleranceToOther(t *testing.T) {
	// no fields - so same day...
	constraint := &DatetimeToleranceToOther{PropertyName: "bar"}
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:        JsonAny,
				Constraints: Constraints{constraint},
			},
			"bar": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-04-01T16:30:00",
		"bar": "2022-04-02T00:00:00"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Value must be same day as value of property 'bar'", violations[0].Message)

	obj["foo"] = "2022-04-02T12:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = time.Date(2022, 4, 2, 18, 19, 20, 21, time.UTC)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["bar"] = "bad date"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj["foo"] = "bad date"
	obj["bar"] = "2022-04-01T12:00:00"
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	// other as time...
	obj["foo"] = "2022-04-01T12:00:00"
	obj["bar"] = time.Date(2022, 4, 1, 16, 17, 18, 19, time.UTC)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// with ignore nulls...
	constraint.IgnoreNull = true
	obj["foo"] = nil
	obj["bar"] = "2022-04-01T12:00:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-04-01T12:00:00"
	obj["bar"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestCheckDatetimeTolerance(t *testing.T) {
	testCases := []struct {
		date   time.Time
		other  time.Time
		amount int64
		unit   string
		expect bool
	}{
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			time.Date(2021, 3, 31, 12, 13, 14, 15, time.UTC),
			0,
			"year",
			false,
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 1, 1, 12, 13, 14, 15, time.UTC),
			0,
			"year",
			true,
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			time.Date(2021, 3, 31, 12, 13, 14, 15, time.UTC),
			-1,
			"year",
			true,
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			time.Date(2023, 3, 31, 12, 13, 14, 15, time.UTC),
			1,
			"year",
			true,
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			time.Date(2023, 4, 1, 12, 13, 14, 15, time.UTC),
			1,
			"year",
			false,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 4, 1, 16, 13, 14, 15, time.UTC),
			0,
			"day",
			true,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 4, 2, 16, 13, 14, 15, time.UTC),
			0,
			"day",
			false,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 3, 31, 16, 13, 14, 15, time.UTC),
			0,
			"day",
			false,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 4, 2, 11, 13, 14, 15, time.UTC),
			1,
			"day",
			true,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 4, 2, 12, 30, 14, 15, time.UTC),
			1,
			"day",
			false,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 3, 31, 12, 30, 14, 15, time.UTC),
			-1,
			"day",
			true,
		},
		{
			time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC),
			time.Date(2022, 3, 31, 12, 0, 14, 15, time.UTC),
			-1,
			"day",
			false,
		},
	}
	for i, tc := range testCases {
		ltgtsame := "SAME"
		if tc.amount < 0 {
			ltgtsame = "BEFORE"
		} else if tc.amount > 0 {
			ltgtsame = "AFTER"
		}
		if !tc.expect {
			ltgtsame = "!" + ltgtsame
		}
		t.Run(fmt.Sprintf("[%d]%s=%d[%s]_\"%s\"_\"%s\"", i+1, ltgtsame, tc.amount, tc.unit,
			tc.date.Format("2006-01-02T15:04:05.999999999"), tc.other.Format("2006-01-02T15:04:05.999999999")),
			func(t *testing.T) {
				ok := checkDatetimeTolerance(&tc.date, &tc.other, tc.amount, tc.unit, false)
				require.Equal(t, tc.expect, ok)
			})
	}
}

func TestCheckDatetimeToleranceMin(t *testing.T) {
	// same...
	value := time.Date(2022, 7, 15, 12, 13, 14, 15, time.UTC)
	other := time.Date(2027, 7, 15, 12, 13, 14, 15, time.UTC)
	result := checkDatetimeTolerance(&value, &other, 5, "year", true)
	require.False(t, result)

	now := time.Date(2022, 7, 15, 12, 13, 14, 15, time.UTC)
	fourYearsAgo := time.Date(2018, 7, 15, 12, 13, 14, 15, time.UTC)
	result = checkDatetimeTolerance(&now, &fourYearsAgo, -5, "year", true)
	require.False(t, result)
	sixYearsAgo := time.Date(2016, 7, 15, 12, 13, 14, 15, time.UTC)
	result = checkDatetimeTolerance(&now, &sixYearsAgo, -5, "year", true)
	require.True(t, result)

	fourYearsTime := time.Date(2026, 7, 15, 12, 13, 14, 15, time.UTC)
	result = checkDatetimeTolerance(&now, &fourYearsTime, 5, "year", true)
	require.False(t, result)
	sixYearsTime := time.Date(2028, 7, 15, 12, 13, 14, 15, time.UTC)
	result = checkDatetimeTolerance(&now, &sixYearsTime, 5, "year", true)
	require.True(t, result)
}

func TestCheckDatetimeToleranceFailWithBadUnit(t *testing.T) {
	dt := time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)

	ok := checkDatetimeTolerance(&dt, &dt, -1, "unknown", false)
	require.False(t, ok)
	ok = checkDatetimeTolerance(&dt, &dt, 0, "", false) // same day assumed
	require.True(t, ok)
}

func TestShiftDatetimeBy(t *testing.T) {
	testCases := []struct {
		date   time.Time
		unit   string
		amount int64
		expect bool
		result string
	}{
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"unknown",
			1,
			false,
			"",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"", // default to day
			1,
			true,
			"2022-04-01T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"", // default to day
			-1,
			true,
			"2022-03-30T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"millennium",
			-1,
			true,
			"1022-03-31T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"millennium",
			-1,
			true,
			"1024-02-29T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"millennium",
			1,
			true,
			"3024-02-29T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"century",
			-1,
			true,
			"1924-02-29T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"century",
			1,
			true,
			"2124-02-29T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"decade",
			-1,
			true,
			"2014-02-28T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"decade",
			1,
			true,
			"2034-02-28T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"year",
			-1,
			true,
			"2021-03-31T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"year",
			1,
			true,
			"2023-03-31T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"year",
			-1,
			true,
			"2023-02-28T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"year",
			1,
			true,
			"2025-02-28T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"year",
			-10,
			true,
			"2014-02-28T12:13:14.000000015",
		},
		{
			time.Date(2024, 2, 29, 12, 13, 14, 15, time.UTC),
			"year",
			-8,
			true,
			"2016-02-29T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"year",
			-10,
			true,
			"2012-03-31T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"month",
			-1,
			true,
			"2022-02-28T12:13:14.000000015",
		},
		{
			time.Date(2022, 1, 31, 12, 13, 14, 15, time.UTC),
			"month",
			1,
			true,
			"2022-02-28T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"month",
			-13,
			true,
			"2021-02-28T12:13:14.000000015",
		},
		{
			time.Date(2022, 1, 31, 12, 13, 14, 15, time.UTC),
			"month",
			13,
			true,
			"2023-02-28T12:13:14.000000015",
		},
		{
			time.Date(2022, 1, 31, 12, 13, 14, 15, time.UTC),
			"month",
			1,
			true,
			"2022-02-28T12:13:14.000000015",
		},
		{
			time.Date(2024, 3, 31, 12, 13, 14, 15, time.UTC),
			"month",
			-1,
			true,
			"2024-02-29T12:13:14.000000015",
		},
		{
			time.Date(2022, 10, 31, 12, 13, 14, 15, time.UTC),
			"month",
			-1,
			true,
			"2022-09-30T12:13:14.000000015",
		},
		{
			time.Date(2022, 10, 31, 12, 13, 14, 15, time.UTC),
			"month",
			1,
			true,
			"2022-11-30T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"week",
			-1,
			true,
			"2022-03-24T12:13:14.000000015",
		},
		{
			time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC),
			"week",
			1,
			true,
			"2022-04-07T12:13:14.000000015",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]ShiftTime:by:%d,unit:%s,date:%s", i+1, tc.amount, tc.unit, tc.date.Format("2006-01-02T15:04:05.999999999")),
			func(t *testing.T) {
				shfted, ok := shiftDatetimeBy(&tc.date, tc.amount, tc.unit)
				if tc.expect {
					require.True(t, ok)
					require.Equal(t, tc.result, shfted.Format("2006-01-02T15:04:05.999999999"))
				} else {
					require.False(t, ok)
				}
			})
	}
}

func TestCheckDatetimeToleranceSames(t *testing.T) {
	dt := time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other := time.Date(2922, 1, 1, 12, 13, 14, 15, time.UTC)
	ok := checkDatetimeTolerance(&dt, &other, 0, "millennium", false)
	require.True(t, ok)
	other = time.Date(1922, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "millennium", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2000, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "century", false)
	require.True(t, ok)
	other = time.Date(2122, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "century", false)
	require.False(t, ok)

	dt = time.Date(2020, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2029, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "decade", false)
	require.True(t, ok)
	other = time.Date(2031, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "decade", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2022, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "year", false)
	require.True(t, ok)
	other = time.Date(2023, 1, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "year", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2022, 3, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "month", false)
	require.True(t, ok)
	other = time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "month", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2022, 3, 28, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "week", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 27, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "week", false)
	require.False(t, ok)
	other = time.Date(2022, 4, 3, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "week", false)
	require.True(t, ok)
	other = time.Date(2022, 4, 4, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "week", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "day", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 30, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "day", false)
	require.False(t, ok)
	other = time.Date(2022, 4, 1, 12, 13, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "day", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 0, 14, 15, time.UTC)
	other = time.Date(2022, 3, 31, 12, 30, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "hour", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 13, 30, 14, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "hour", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 0, 15, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 30, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "minute", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 12, 14, 15, 15, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "minute", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 0, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 14, 1000, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "second", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 12, 13, 15, 1000, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "second", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 111000000, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 14, 111222333, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "milli", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 12, 13, 14, 112000000, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "milli", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 111222000, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 14, 111222333, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "micro", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 12, 13, 14, 111223000, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "micro", false)
	require.False(t, ok)

	dt = time.Date(2022, 3, 31, 12, 13, 14, 111222333, time.UTC)
	other = time.Date(2022, 3, 31, 12, 13, 14, 111222333, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "nano", false)
	require.True(t, ok)
	other = time.Date(2022, 3, 31, 12, 13, 14, 111222331, time.UTC)
	ok = checkDatetimeTolerance(&dt, &other, 0, "nano", false)
	require.False(t, ok)

	ok = checkDatetimeTolerance(&dt, &dt, 0, "UNKNOWN", false)
	require.False(t, ok)
}

func TestDatetimeToleranceMessages(t *testing.T) {
	constraint := &DatetimeTolerance{Value: "2022-04-01T18:19:00"}

	msg := constraint.GetMessage(nil)
	require.Equal(t, "Value must be same day as 2022-04-01T18:19:00", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 century after 2022-04-01T18:19:00", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 centuries after 2022-04-01T18:19:00", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 year before 2022-04-01T18:19:00", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 years before 2022-04-01T18:19:00", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 century after 2022-04-01T18:19:00", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 centuries after 2022-04-01T18:19:00", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 year before 2022-04-01T18:19:00", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 years before 2022-04-01T18:19:00", msg)

	constraint.Message = "test"
	msg = constraint.GetMessage(nil)
	require.Equal(t, "test", msg)
}

func TestDatetimeToleranceToNowMessages(t *testing.T) {
	constraint := &DatetimeToleranceToNow{}

	msg := constraint.GetMessage(nil)
	require.Equal(t, "Value must be same day as now", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 century after now", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 centuries after now", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 year before now", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 years before now", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 century after now", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 centuries after now", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 year before now", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 years before now", msg)

	constraint.Message = "test"
	msg = constraint.GetMessage(nil)
	require.Equal(t, "test", msg)
}

func TestDatetimeToleranceToOtherMessages(t *testing.T) {
	constraint := &DatetimeToleranceToOther{PropertyName: "bar"}

	msg := constraint.GetMessage(nil)
	require.Equal(t, "Value must be same day as value of property 'bar'", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 century after value of property 'bar'", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 centuries after value of property 'bar'", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 1 year before value of property 'bar'", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must not be more than 2 years before value of property 'bar'", msg)

	constraint.Unit = "century"
	constraint.Duration = 1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 century after value of property 'bar'", msg)
	constraint.Duration = 2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 centuries after value of property 'bar'", msg)

	constraint.Unit = "year"
	constraint.Duration = -1
	constraint.MinCheck = true
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 1 year before value of property 'bar'", msg)
	constraint.Duration = -2
	msg = constraint.GetMessage(nil)
	require.Equal(t, "Value must be at least 2 years before value of property 'bar'", msg)

	constraint.Message = "test"
	msg = constraint.GetMessage(nil)
	require.Equal(t, "test", msg)
}

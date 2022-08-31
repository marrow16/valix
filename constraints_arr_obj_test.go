package valix

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestArrayOf(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&ArrayOf{Type: "string", AllowNullElement: false}, false)
	obj := jsonObject(`{
		"foo": ["ok", false]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgArrayElementType, JsonString), violations[0].Message)

	obj = jsonObject(`{
		"foo": [null]
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgArrayElementType, JsonString), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(fmtMsgArrayElementTypeOrNull, JsonString), violations[0].Message)
	obj = jsonObject(`{
		"foo": [null, "ok2"]
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestArrayUnique(t *testing.T) {
	constraint := &ArrayUnique{}
	validator := buildFooValidator(JsonArray, constraint, false)
	obj := jsonObject(`{
		"foo": ["foo", "Foo", false, true, 1, 1.1, null, {"foo": false}, {"foo": true}, ["aaa"], ["bbb"]]
	}`)

	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj = jsonObject(`{"foo": ["aaa", "aaa"]}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgArrayUnique, violations[0].Message)

	constraint.IgnoreCase = true
	obj = jsonObject(`{"foo": ["Aaa", "aAA"]}`)
	ok, _ = validator.Validate(obj)
	require.False(t, ok)

	obj = jsonObject(`{"foo": [null, null, "foo"]}`)
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	constraint.IgnoreNulls = true
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// objects...
	obj = jsonObject(`{"foo": [{"foo": true}, {"foo": true}]}`)
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj = jsonObject(`{"foo": [{"foo": true}, {"foo": false}]}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// arrays..
	obj = jsonObject(`{"foo": [["aaa"], ["aaa"]]}`)
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	obj = jsonObject(`{"foo": [["aaa"], ["bbb"]]}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// json various numerics...
	testCases := []struct {
		arr        []interface{}
		expectFail bool
	}{
		{
			[]interface{}{json.Number("xxx"), 1},
			false,
		},
		{
			[]interface{}{json.Number("1.0"), json.Number("1")},
			true,
		},
		{
			[]interface{}{json.Number("1.0"), 1},
			true,
		},
		{
			[]interface{}{json.Number("1.0"), 1.0},
			true,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("ArrayUniqueNumerics[%d]", i+2), func(t *testing.T) {
			obj["foo"] = tc.arr
			ok, violations = validator.Validate(obj)
			if tc.expectFail {
				require.False(t, ok)
				require.Equal(t, 1, len(violations))
				require.Equal(t, msgArrayUnique, violations[0].Message)
			} else {
				require.True(t, ok)
			}
		})
	}
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
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(fmtMsgMinLen, 2), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
		"quz": nil,
	}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz", "quz"}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz"}
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthMinOnlyExc(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&Length{Minimum: 3, ExclusiveMin: true}, false)
	obj := jsonObject(`{
		"foo": ["bar"]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinLenExc, 3), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz", "quz"}
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = []interface{}{"foo", "bar", "baz"}
	ok, _ = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinLenExc, 3), violations[0].Message)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthExactWithString(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&LengthExact{Value: 3}, false)
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

	obj["foo"] = "Abcd"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

	obj["foo"] = "Abc"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestLengthExactWithObject(t *testing.T) {
	validator := buildFooValidator(JsonObject,
		&LengthExact{Value: 3}, false)
	obj := jsonObject(`{
		"foo": {
			"bar": null
		}
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

	obj["foo"] = map[string]interface{}{
		"foo": nil,
		"bar": nil,
		"baz": nil,
		"quz": nil,
	}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

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

func TestLengthExactWithArray(t *testing.T) {
	validator := buildFooValidator(JsonArray,
		&LengthExact{Value: 3}, false)
	obj := jsonObject(`{
		"foo": ["bar"]
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz", "quz"}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgExactLen, 3), violations[0].Message)

	obj["foo"] = []interface{}{"foo", "bar", "baz"}
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestNotEmpty(t *testing.T) {
	validator := buildFooValidator(JsonAny,
		&NotEmpty{}, false)
	obj := jsonObject(`{
		"foo": []
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmpty, violations[0].Message)
	obj = jsonObject(`{
		"foo": ["bar"]
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj = jsonObject(`{
		"foo": {}
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmpty, violations[0].Message)
	obj = jsonObject(`{
		"foo": {"foo": "bar"}
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj = jsonObject(`{
		"foo": ""
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmpty, violations[0].Message)

	obj = jsonObject(`{
		"foo": "bar"
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

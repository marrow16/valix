package valix

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTypeCheckString(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "xxx"
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)
}

func TestTypeCheckNumber(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonNumber,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 1
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = float64(1)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = json.Number("1.5")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonNumber), violations[0].Message)
}

func TestTypeCheckInteger(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonInteger,
			},
		},
	}
	obj := jsonObject(`{
		"foo": 1
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = float64(1)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = json.Number("1")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = json.Number("1.1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonInteger), violations[0].Message)
}

func TestTypeCheckBoolean(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonBoolean,
			},
		},
	}
	obj := jsonObject(`{
		"foo": true
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)
}

func TestTypeCheckObject(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonObject,
			},
		},
	}
	obj := jsonObject(`{
		"foo": {}
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)
}

func TestTypeCheckArray(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
			},
		},
	}
	obj := jsonObject(`{
		"foo": []
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonArray), violations[0].Message)
}

func TestTypeCheckAny(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
			},
		},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = json.Number("1.1")
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

func TestJsonTypeToString(t *testing.T) {
	testValues := map[JsonType]string{
		JsonString:  jsonTypeTokenString,
		JsonNumber:  jsonTypeTokenNumber,
		JsonInteger: jsonTypeTokenInteger,
		JsonBoolean: jsonTypeTokenBoolean,
		JsonObject:  jsonTypeTokenObject,
		JsonArray:   jsonTypeTokenArray,
		JsonAny:     jsonTypeTokenAny,
	}
	for k, v := range testValues {
		t.Run(fmt.Sprintf("JsonType_Token_%s", v), func(t *testing.T) {
			str := k.String()
			require.Equal(t, v, str)
		})
	}
	jt := JsonType(-1)
	str := jt.String()
	require.Equal(t, "", str)
}

func TestJsonTypeFromString(t *testing.T) {
	testValues := map[string]JsonType{
		jsonTypeTokenString:  JsonString,
		jsonTypeTokenNumber:  JsonNumber,
		jsonTypeTokenInteger: JsonInteger,
		jsonTypeTokenBoolean: JsonBoolean,
		jsonTypeTokenObject:  JsonObject,
		jsonTypeTokenArray:   JsonArray,
		jsonTypeTokenAny:     JsonAny,
	}
	for k, v := range testValues {
		t.Run(fmt.Sprintf("JsonType_Token_%s", k), func(t *testing.T) {
			av, ok := JsonTypeFromString(k)
			require.True(t, ok)
			require.Equal(t, v, av)
		})
	}
	_, ok := JsonTypeFromString("xxx")
	require.False(t, ok)
}

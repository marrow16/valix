package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonNumber), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonNumber), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = json.Number("1.1")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonInteger), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = 0
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonArray), violations[0].Message)
	obj["foo"] = "xxx"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonArray), violations[0].Message)
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
	jt := JsonString
	require.Equal(t, jsonTypeTokenString, jt.String())
	jt = JsonNumber
	require.Equal(t, jsonTypeTokenNumber, jt.String())
	jt = JsonInteger
	require.Equal(t, jsonTypeTokenInteger, jt.String())
	jt = JsonBoolean
	require.Equal(t, jsonTypeTokenBoolean, jt.String())
	jt = JsonObject
	require.Equal(t, jsonTypeTokenObject, jt.String())
	jt = JsonArray
	require.Equal(t, jsonTypeTokenArray, jt.String())
	jt = JsonAny
	require.Equal(t, jsonTypeTokenAny, jt.String())
	jt = 99
	require.Equal(t, jsonTypeTokenUndefined, jt.String())
}

func TestJsonTypeFromString(t *testing.T) {
	token := jsonTypeTokenString
	require.Equal(t, JsonString, JsonTypeFromString(token))
	token = jsonTypeTokenNumber
	require.Equal(t, JsonNumber, JsonTypeFromString(token))
	token = jsonTypeTokenInteger
	require.Equal(t, JsonInteger, JsonTypeFromString(token))
	token = jsonTypeTokenBoolean
	require.Equal(t, JsonBoolean, JsonTypeFromString(token))
	token = jsonTypeTokenObject
	require.Equal(t, JsonObject, JsonTypeFromString(token))
	token = jsonTypeTokenArray
	require.Equal(t, JsonArray, JsonTypeFromString(token))
	token = jsonTypeTokenAny
	require.Equal(t, JsonAny, JsonTypeFromString(token))
	token = "???"
	require.Equal(t, JsonTypeUndefined, JsonTypeFromString(token))
}

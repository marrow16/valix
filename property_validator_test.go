package valix

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

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

func TestTypeCheckDatetime(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonDatetime,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "2022-05-04T10:11:12"
	}`)
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = time.Now()
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = &time.Time{}
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = Time{}
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	obj["foo"] = &Time{}
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonDatetime), violations[0].Message)
	obj["foo"] = 1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonDatetime), violations[0].Message)
	obj["foo"] = 1.1
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonDatetime), violations[0].Message)
	obj["foo"] = []interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonDatetime), violations[0].Message)
	obj["foo"] = map[string]interface{}{}
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonDatetime), violations[0].Message)
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
		JsonString:   jsonTypeTokenString,
		JsonDatetime: jsonTypeTokenDatetime,
		JsonNumber:   jsonTypeTokenNumber,
		JsonInteger:  jsonTypeTokenInteger,
		JsonBoolean:  jsonTypeTokenBoolean,
		JsonObject:   jsonTypeTokenObject,
		JsonArray:    jsonTypeTokenArray,
		JsonAny:      jsonTypeTokenAny,
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
		jsonTypeTokenString:   JsonString,
		jsonTypeTokenDatetime: JsonDatetime,
		jsonTypeTokenNumber:   JsonNumber,
		jsonTypeTokenInteger:  JsonInteger,
		jsonTypeTokenBoolean:  JsonBoolean,
		jsonTypeTokenObject:   JsonObject,
		jsonTypeTokenArray:    JsonArray,
		jsonTypeTokenAny:      JsonAny,
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

func TestPropertyValidator_SetType(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, JsonAny, pv.Type)

	pv.SetType(JsonString)
	require.Equal(t, JsonString, pv.Type)
}

func TestPropertyValidator_SetNullable(t *testing.T) {
	pv := &PropertyValidator{NotNull: true}
	require.True(t, pv.NotNull)

	pv.SetNullable()
	require.False(t, pv.NotNull)
}

func TestPropertyValidator_SetNotNullable(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.NotNull)

	pv.SetNotNullable()
	require.True(t, pv.NotNull)
}

func TestPropertyValidator_SetMandatory(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.Mandatory)

	pv.SetMandatory()
	require.True(t, pv.Mandatory)
}

func TestPropertyValidator_SetRequired(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.Mandatory)

	pv.SetRequired()
	require.True(t, pv.Mandatory)
}

func TestPropertyValidator_SetOptional(t *testing.T) {
	pv := &PropertyValidator{Mandatory: true}
	require.True(t, pv.Mandatory)

	pv.SetOptional()
	require.False(t, pv.Mandatory)
}

func TestPropertyValidator_AddMandatoryWhens(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.MandatoryWhen))

	pv.AddMandatoryWhens("foo", "bar")
	require.Equal(t, 2, len(pv.MandatoryWhen))
}

func TestPropertyValidator_AddConstraints(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	pv.AddConstraints(&StringNotEmpty{}, &StringNotBlank{})
	require.Equal(t, 2, len(pv.Constraints))
}

func TestPropertyValidator_SetObjectValidator(t *testing.T) {
	pv := &PropertyValidator{}
	require.Nil(t, pv.ObjectValidator)

	pv.SetObjectValidator(&Validator{})
	require.NotNil(t, pv.ObjectValidator)
}

func TestPropertyValidator_SetOrder(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, pv.Order)

	pv.SetOrder(1)
	require.Equal(t, 1, pv.Order)
}

func TestPropertyValidator_AddWhenConditions(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.WhenConditions))

	pv.AddWhenConditions("foo", "bar")
	require.Equal(t, 2, len(pv.WhenConditions))
}

func TestPropertyValidator_AddUnwantedConditions(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.UnwantedConditions))

	pv.AddUnwantedConditions("foo", "bar")
	require.Equal(t, 2, len(pv.UnwantedConditions))
}

func TestPropertyValidator_SetRequiredWith(t *testing.T) {
	pv := &PropertyValidator{}
	require.Nil(t, pv.RequiredWith)

	pv.SetRequiredWith(MustParseExpression("foo && bar"))
	require.NotNil(t, pv.RequiredWith)
	require.Equal(t, 2, len(pv.RequiredWith))
}

func TestPropertyValidator_SetRequiredWithMessage(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, "", pv.RequiredWithMessage)

	pv.SetRequiredWithMessage("fooey")
	require.Equal(t, "fooey", pv.RequiredWithMessage)
}

func TestPropertyValidator_SetUnwantedWith(t *testing.T) {
	pv := &PropertyValidator{}
	require.Nil(t, pv.UnwantedWith)

	pv.SetUnwantedWith(MustParseExpression("foo && bar"))
	require.NotNil(t, pv.UnwantedWith)
	require.Equal(t, 2, len(pv.UnwantedWith))
}

func TestPropertyValidator_SetUnwantedWithMessage(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, "", pv.UnwantedWithMessage)

	pv.SetUnwantedWithMessage("fooey")
	require.Equal(t, "fooey", pv.UnwantedWithMessage)
}

func TestPropertyValidator_Validate(t *testing.T) {
	pv := &PropertyValidator{}

	ok, violations := pv.Validate(nil)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	pv = &PropertyValidator{
		Type:      JsonString,
		Mandatory: true,
		NotNull:   true,
		Constraints: Constraints{
			&StringNotBlank{},
		},
	}
	ok, violations = pv.Validate(nil)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueCannotBeNull, violations[0].Message)
	require.Equal(t, "", violations[0].Property)

	ok, violations = pv.Validate("")
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotBlankString, violations[0].Message)
	require.Equal(t, "", violations[0].Property)

	ok, violations = pv.Validate(1)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, "string"), violations[0].Message)
	require.Equal(t, "", violations[0].Property)
}

func TestNewPropertyValidator(t *testing.T) {
	pv, err := NewPropertyValidator()
	require.Nil(t, err)
	require.NotNil(t, pv)

	pv, err = NewPropertyValidator("")
	require.Nil(t, err)
	require.NotNil(t, pv)

	pv, err = NewPropertyValidator("]),'])")
	require.NotNil(t, err)
	require.Equal(t, "unopened parenthesis at position 0", err.Error())

	_, err = NewPropertyValidator("blah")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "blah"), err.Error())

	_, err = NewPropertyValidator("blah,blah")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "blah"), err.Error())

	pv, err = NewPropertyValidator("mandatory,notNull", "type:"+jsonTypeTokenString, "&StringNotEmpty{Message:'fooey'},&StringNotBlank{Message:'fooey2',Stop:true}", "only")
	require.Nil(t, err)
	require.True(t, pv.Mandatory)
	require.True(t, pv.NotNull)
	require.Equal(t, JsonString, pv.Type)
	require.True(t, pv.Only)
	require.Equal(t, 2, len(pv.Constraints))
	c1, ok := pv.Constraints[0].(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "fooey", c1.Message)
	c2, ok := pv.Constraints[1].(*StringNotBlank)
	require.True(t, ok)
	require.Equal(t, "fooey2", c2.Message)
	require.True(t, c2.Stop)
}

func TestCreatePropertyValidator(t *testing.T) {
	require.Panics(t, func() {
		_ = CreatePropertyValidator("]),'])")

	})

	pv := CreatePropertyValidator("mandatory,notNull", "type:"+jsonTypeTokenString, "&StringNotEmpty{Message:'fooey'},&StringNotBlank{Message:'fooey2',Stop:true}", "only")
	require.True(t, pv.Mandatory)
	require.True(t, pv.NotNull)
	require.Equal(t, JsonString, pv.Type)
	require.True(t, pv.Only)
	require.Equal(t, 2, len(pv.Constraints))
	c1, ok := pv.Constraints[0].(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "fooey", c1.Message)
	c2, ok := pv.Constraints[1].(*StringNotBlank)
	require.True(t, ok)
	require.Equal(t, "fooey2", c2.Message)
	require.True(t, c2.Stop)
}

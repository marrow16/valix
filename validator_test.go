package valix

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// test validator - validates JSON...
// {
//   "name": <<string,mandatory,not-null,length 1-255>>
//   "age": <<int,mandatory,not-null,positive-or-zero>>
// }
var personValidator = &Validator{
	IgnoreUnknownProperties: false,
	AllowArray:              true,
	Properties: Properties{
		"name": {
			Type:      JsonString,
			NotNull:   true,
			Mandatory: true,
			Constraints: Constraints{
				&StringLength{Minimum: 1, Maximum: 255},
			},
		},
		"age": {
			Type:      JsonInteger,
			NotNull:   true,
			Mandatory: true,
			Constraints: Constraints{
				&PositiveOrZero{},
			},
		},
	},
}

func TestEmptyValidatorWorks(t *testing.T) {
	v := Validator{}
	o := jsonObject(`{}`)
	ok, violations := v.Validate(o)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

var addPersonToGroupValidator = &Validator{
	IgnoreUnknownProperties: false,
	Properties: Properties{
		"person": {
			Type:            JsonObject,
			ObjectValidator: personValidator,
		},
		"group": {
			Type:      JsonString,
			NotNull:   true,
			Mandatory: true,
			Constraints: Constraints{
				&StringLength{Minimum: 1, Maximum: 255},
			},
		},
	},
}

func TestValidatorWithReUseWorks(t *testing.T) {
	o := jsonObject(`{
		"person": {
			"name": "",
			"age": -1
		},
		"group": ""
	}`)
	ok, violations := addPersonToGroupValidator.Validate(o)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	require.Equal(t, CodePropertyConstraintFail, violations[0].Codes[0])
	require.Equal(t, CodePropertyConstraintFail, violations[1].Codes[0])
	require.Equal(t, CodePropertyConstraintFail, violations[2].Codes[0])
}

func TestMissingPropertyDetection(t *testing.T) {
	o := jsonObject(`{
		"name": "Bilbo"
	}`)
	ok, violations := personValidator.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgMissingProperty, violations[0].Message)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, CodeMissingProperty, violations[0].Codes[0])
	require.Equal(t, "", violations[0].Path)
}

func TestUnknownPropertyDetection(t *testing.T) {
	o := jsonObject(`{
		"name": "Bilbo",
		"age": 16,
		"unknown_property": true
	}`)
	ok, violations := personValidator.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnknownProperty, violations[0].Message)
	require.Equal(t, "unknown_property", violations[0].Property)
	require.Equal(t, CodeUnknownProperty, violations[0].Codes[0])
	require.Equal(t, "", violations[0].Path)
}

func TestValidationOfArray(t *testing.T) {
	a := jsonArray(`[
		{
			"name": "",
			"age": 16
		},
		{
			"name": "Gandalf",
			"age": -1
		}
	]`)
	ok, violations := personValidator.ValidateArrayOf(a)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, "[0]", violations[0].Path)
	require.Equal(t, "name", violations[0].Property)
	require.Equal(t, "[1]", violations[1].Path)
	require.Equal(t, "age", violations[1].Property)
}

func TestValidationOfArrayFailsWithNonObjectElement(t *testing.T) {
	a := jsonArray(`[
		{
			"name": "Bilbo",
			"age": 16
		},
		"this_should_be_an_object"
	]`)
	ok, violations := personValidator.ValidateArrayOf(a)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgArrayElementMustBeObject, violations[0].Message)
	require.Equal(t, CodeArrayElementMustBeObject, violations[0].Codes[0])
	require.Equal(t, 1, violations[0].Codes[1])
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "[1]", violations[0].Property)
}

func TestValidatorWithObjectConstraint(t *testing.T) {
	v := Validator{
		IgnoreUnknownProperties: true,
		Constraints: Constraints{
			&Length{Minimum: 2, Maximum: 3},
		},
	}
	o := jsonObject(`{
		"foo": null
	}`)
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgMinMax, 2, tokenInclusive, 3, tokenInclusive), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
	require.Equal(t, CodeValidatorConstraintFail, violations[0].Codes[0])
}

func TestRequestValidation(t *testing.T) {
	body := strings.NewReader(`{"name": "", "age": -1}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := personValidator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.False(t, violations[0].BadRequest)
	require.False(t, violations[1].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationUsingJsonNumber(t *testing.T) {
	validator := Validator{
		UseNumber: true,
		Properties: Properties{
			"foo": {
				Type:        JsonNumber,
				Constraints: Constraints{&Positive{}},
			},
			"bar": {
				Type:        JsonNumber,
				Constraints: Constraints{&Negative{}},
			},
		},
	}
	body := strings.NewReader(`{"foo": -1, "bar": 1}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := validator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.False(t, violations[0].BadRequest)
	require.False(t, violations[1].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationWithArray(t *testing.T) {
	body := strings.NewReader(`[{"name": "", "age": -1}]`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := personValidator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.False(t, violations[0].BadRequest)
	require.False(t, violations[1].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationWithJsonNullBody(t *testing.T) {
	body := strings.NewReader(`null`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := personValidator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyNotJsonNull, violations[0].Message)
	require.Equal(t, CodeRequestBodyNotJsonNull, violations[0].Codes[0])
	require.True(t, violations[0].BadRequest)
	require.Nil(t, obj)
}

func TestRequestValidationWithNoBody(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := personValidator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyEmpty, violations[0].Message)
	require.Equal(t, CodeRequestBodyEmpty, violations[0].Codes[0])
	require.True(t, violations[0].BadRequest)
	require.Nil(t, obj)
}

func TestRequestValidationFailsWithArrayWhenArrayNotAllowed(t *testing.T) {
	body := strings.NewReader(`[{"foo": true}]`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	v := Validator{AllowArray: false}
	ok, violations, obj := v.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyNotJsonArray, violations[0].Message)
	require.Equal(t, CodeRequestBodyNotJsonArray, violations[0].Codes[0])
	require.False(t, violations[0].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationFailsWithObjectWhenObjectNotAllowed(t *testing.T) {
	body := strings.NewReader(`{"foo": true}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	v := Validator{AllowArray: true, DisallowObject: true}
	ok, violations, obj := v.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyExpectedJsonArray, violations[0].Message)
	require.Equal(t, CodeRequestBodyExpectedJsonArray, violations[0].Codes[0])
	require.False(t, violations[0].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationFailsWhenExpectingObject(t *testing.T) {
	body := strings.NewReader(`false`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	v := Validator{AllowArray: false, DisallowObject: true}
	ok, violations, obj := v.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyExpectedJsonObject, violations[0].Message)
	require.Equal(t, CodeRequestBodyExpectedJsonObject, violations[0].Codes[0])
	require.True(t, violations[0].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationFailsWhenObjectDisallowed(t *testing.T) {
	body := strings.NewReader(`{"foo": true}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	v := Validator{AllowArray: false, DisallowObject: true}
	ok, violations, obj := v.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyNotJsonObject, violations[0].Message)
	require.Equal(t, CodeRequestBodyNotJsonObject, violations[0].Codes[0])
	require.False(t, violations[0].BadRequest)
	require.NotNil(t, obj)
}

func TestRequestValidationWithBadJson(t *testing.T) {
	body := strings.NewReader(`this is bad json!`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, obj := personValidator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnableToDecodeRequest, violations[0].Message)
	require.Equal(t, CodeUnableToDecodeRequest, violations[0].Codes[0])
	require.True(t, violations[0].BadRequest)
	require.Nil(t, obj)
}

func TestValidatorStopsOnConstraint(t *testing.T) {
	v := Validator{
		IgnoreUnknownProperties: false,
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
				vcx.Stop()
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				// this should not get checked because the constraint stopped...
				NotNull: true,
			},
		},
	}
	o := jsonObject(`{"foo": null}`)

	ok, violations := v.Validate(o)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

func TestValidatorStopsOnArrayElement(t *testing.T) {
	testMsg := "Test message"
	v := Validator{
		AllowArray: true,
		Properties: Properties{
			"foo": {
				Type: JsonBoolean,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.Stop()
						return false, cc.GetMessage(vcx)
					}, testMsg),
				},
			},
		},
	}
	a := jsonArray(`[{"foo": true},{"foo": "this should be bool but won't get checked"}]`)

	ok, violations := v.ValidateArrayOf(a)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
}

func TestValidatorStopsOnFirst(t *testing.T) {
	v := Validator{
		StopOnFirst: true,
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				NotNull:   true,
			},
			"bar": {
				Type:      JsonNumber,
				Mandatory: true,
				NotNull:   true,
			},
		},
	}
	o := jsonObject(`{"foo": null, "unknown": true}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))

	v.StopOnFirst = false
	ok, violations = v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
}

func TestPropertyValueObjectValidatorFailsForObjectOrArray(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				ObjectValidator: &Validator{
					DisallowObject: false,
					AllowArray:     true,
				},
			},
		},
	}
	o := jsonObject(`{"foo": "should be object or array"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueMustBeObjectOrArray, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForObjectOnly(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				ObjectValidator: &Validator{
					DisallowObject: false,
					AllowArray:     false,
				},
			},
		},
	}
	o := jsonObject(`{"foo": "should be object"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueMustBeObject, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForArrayOnly(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				ObjectValidator: &Validator{
					DisallowObject: true,
					AllowArray:     true,
				},
			},
		},
	}
	o := jsonObject(`{"foo": "should be array"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueMustBeArray, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForNeitherObjectNorArray(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				ObjectValidator: &Validator{
					DisallowObject: true,
					AllowArray:     false,
				},
			},
		},
	}
	o := jsonObject(`{"foo": "should be array"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPropertyObjectValidatorError, violations[0].Message)
}

func TestSubPropertyValidation(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"person": {
				NotNull:         true,
				Mandatory:       true,
				ObjectValidator: personValidator,
			},
		},
	}
	o := jsonObject(`{
		"person": {
			"name": "Bilbo",
			"age": 16
		}
	}`)
	ok, violations := v.Validate(o)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	o["person"] = jsonObject(`{
		"name": "",
		"age": -1
	}`)
	ok, violations = v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
}

func TestSubPropertyAsArrayValidation(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"person": {
				NotNull:         true,
				Mandatory:       true,
				ObjectValidator: personValidator,
			},
		},
	}
	o := jsonObject(`{
		"person": [
			{
				"name": "Bilbo",
				"age": 16
			},
			{
				"name": "Gandalf",
				"age": 100
			}
		]
	}`)
	ok, violations := v.Validate(o)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	o["person"] = jsonObject(`{
		"name": "",
		"age": -1
	}`)
	ok, violations = v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	// because properties are iterated over a map, we can't predict the order - so let's sort them...
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, "person", violations[0].Path)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, msgPositiveOrZero, violations[0].Message)
	require.Equal(t, "person", violations[1].Path)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 1, tokenInclusive, 255, tokenInclusive), violations[1].Message)

	o["person"] = []interface{}{
		jsonObject(`{
			"name": "Foo",
			"age": -1
		}`),
		jsonObject(`{
			"name": "",
			"age": -1
		}`),
	}
	ok, violations = v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, "person[0]", violations[0].Path)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, "person[1]", violations[1].Path)
	require.Equal(t, "age", violations[1].Property)
	require.Equal(t, "person[1]", violations[2].Path)
	require.Equal(t, "name", violations[2].Property)
}

func TestCheckPropertyTypeString(t *testing.T) {
	v := buildFooPropertyTypeValidator(JsonString, true)
	o := jsonObject(`{"foo": 1}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonString), violations[0].Message)

	o["foo"] = true
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = []interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = map[string]interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = nil
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = "str"
	ok, _ = v.Validate(o)
	require.True(t, ok)
}

func TestCheckPropertyTypeNumber(t *testing.T) {
	v := buildFooPropertyTypeValidator(JsonNumber, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonNumber), violations[0].Message)

	o["foo"] = true
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = []interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = map[string]interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = nil
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = 1.0
	ok, _ = v.Validate(o)
	require.True(t, ok)

	o["foo"] = 1
	ok, _ = v.Validate(o)
	require.True(t, ok)
}

func TestCheckPropertyTypeBoolean(t *testing.T) {
	v := buildFooPropertyTypeValidator(JsonBoolean, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonBoolean), violations[0].Message)

	o["foo"] = 1.0
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = []interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = map[string]interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = nil
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = true
	ok, _ = v.Validate(o)
	require.True(t, ok)

	o["foo"] = false
	ok, _ = v.Validate(o)
	require.True(t, ok)
}

func TestCheckPropertyTypeObject(t *testing.T) {
	v := buildFooPropertyTypeValidator(JsonObject, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonObject), violations[0].Message)

	o["foo"] = 1.0
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = []interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = true
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = nil
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = map[string]interface{}{}
	ok, _ = v.Validate(o)
	require.True(t, ok)
}

func TestCheckPropertyTypeArray(t *testing.T) {
	v := buildFooPropertyTypeValidator(JsonArray, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, JsonArray), violations[0].Message)

	o["foo"] = 1.0
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = map[string]interface{}{}
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = true
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = nil
	ok, _ = v.Validate(o)
	require.False(t, ok)

	o["foo"] = []interface{}{}
	ok, _ = v.Validate(o)
	require.True(t, ok)
}

func TestValidatorStops(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				Type:    JsonString,
				NotNull: true,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.Stop()
						return true, ""
					}, ""),
					// the following constraint should never be checked (because the prev constraint stops)
					&StringNotEmpty{},
				},
			},
		},
	}
	o := jsonObject(`{"foo": ""}`)

	ok, violations := v.Validate(o)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// check that the NotNull still works...
	o["foo"] = nil
	ok, violations = v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueCannotBeNull, violations[0].Message)
}

func TestValidatorManuallyAddedViolation(t *testing.T) {
	msg := "Something went wrong"
	v := Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.AddViolation(NewViolation("", "", msg))
						return true, ""
					}, ""),
					&StringNotEmpty{},
				},
			},
		},
	}
	o := jsonObject(`{"foo": ""}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, msg, violations[0].Message)
	require.Equal(t, msgNotEmptyString, violations[1].Message)
}

func TestValidatorManuallyAddedCurrentViolation(t *testing.T) {
	msg := "Something went wrong"
	v := Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.AddViolationForCurrent(msg, true)
						return true, ""
					}, ""),
					&StringNotEmpty{},
				},
			},
		},
	}
	o := jsonObject(`{"foo": ""}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, msg, violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, msgNotEmptyString, violations[1].Message)
	require.Equal(t, "", violations[1].Path)
	require.Equal(t, "foo", violations[1].Property)
}

func TestValidatorContextPathing(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				Constraints: Constraints{&myCustomConstraint{expectedPath: ""}},
				ObjectValidator: &Validator{
					Properties: Properties{
						"bar": {
							Constraints: Constraints{&myCustomConstraint{expectedPath: "foo"}},
							ObjectValidator: &Validator{
								Properties: Properties{
									"baz": {
										Constraints: Constraints{&myCustomConstraint{expectedPath: "foo.bar"}},
										ObjectValidator: &Validator{
											Properties: Properties{
												"qux": {
													Constraints: Constraints{&myCustomConstraint{expectedPath: "foo.bar.baz"}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	o := jsonObject(`{
		"foo": {
			"bar": {
				"baz": {
					"qux": false
				}
			}
		}
	}`)
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))
	for _, violation := range violations {
		require.Equal(t, violation.Message, violation.Path)
	}
}

func TestValidatorContextPathingOnArrays(t *testing.T) {
	v := Validator{
		AllowArray: true,
		Properties: Properties{
			"foo": {
				Constraints: Constraints{&myCustomConstraint{expectedPath: ""}},
				ObjectValidator: &Validator{
					AllowArray: true,
					Properties: Properties{
						"bar": {
							Constraints: Constraints{&myCustomConstraint{expectedPath: "foo[0]"}},
							ObjectValidator: &Validator{
								AllowArray: true,
								Properties: Properties{
									"baz": {
										Constraints: Constraints{&myCustomConstraint{expectedPath: "foo[0].bar[0]"}},
										ObjectValidator: &Validator{
											AllowArray: true,
											Properties: Properties{
												"qux": {
													Constraints: Constraints{&myCustomConstraint{expectedPath: "foo[0].bar[0].baz[0]"}},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	o := jsonObject(`{
		"foo": [
			{
				"bar": [
					{
						"baz": [
							{
								"qux": false
							}
						]
					}
				]
			}
		]
	}`)
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))
	for _, violation := range violations {
		require.Equal(t, violation.Message, violation.Path)
	}
}

func TestCeaseFurtherWorks(t *testing.T) {
	v := Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.CeaseFurther()
						return true, ""
					}, ""),
					// these constraints should not be run because of the CeaseFurther (above)
					&StringNotEmpty{},
				},
				ObjectValidator: &Validator{
					Constraints: Constraints{
						NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
							// this constraint should not be run because of the CeaseFurther (above)
							return false, "Never run"
						}, ""),
					},
				},
			},
			"bar": {
				Type: JsonString,
				// not null will be checked - because this is a different property
				NotNull: true,
			},
		},
	}
	o := jsonObject(`{
		"foo": "  ",
		"bar": null
	}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

func TestValidateReader(t *testing.T) {
	validator := &Validator{
		IgnoreUnknownProperties: true,
		DisallowObject:          false,
		AllowArray:              true,
	}

	r := strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj := validator.ValidateReader(r)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	validator.IgnoreUnknownProperties = false
	r = strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.False(t, violations[0].BadRequest)

	r = strings.NewReader(`NOT JSON`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.True(t, violations[0].BadRequest)
}

func TestValidateReaderCanReadObjectOrArray(t *testing.T) {
	validator := &Validator{
		IgnoreUnknownProperties: true,
		AllowNullJson:           false,
		DisallowObject:          false,
		AllowArray:              false,
	}

	r := strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj := validator.ValidateReader(r)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)
	_, ok = obj.(map[string]interface{})
	require.True(t, ok)

	r = strings.NewReader(`false`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgExpectedJsonObject, violations[0].Message)
	require.Equal(t, CodeExpectedJsonObject, violations[0].Codes[0])

	r = strings.NewReader(`null`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonNull, violations[0].Message)
	require.Equal(t, CodeNotJsonNull, violations[0].Codes[0])

	r = strings.NewReader(`[{"foo": "bar"}]`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonArray, violations[0].Message)
	require.Equal(t, CodeNotJsonArray, violations[0].Codes[0])

	validator.AllowArray = true
	validator.DisallowObject = false
	r = strings.NewReader(`[{"foo": "bar"}]`)
	ok, violations, obj = validator.ValidateReader(r)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)
	_, ok = obj.([]interface{})
	require.True(t, ok)

	validator.AllowArray = true
	validator.DisallowObject = true
	r = strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgExpectedJsonArray, violations[0].Message)
	require.Equal(t, CodeExpectedJsonArray, violations[0].Codes[0])

	validator.AllowArray = false
	validator.DisallowObject = true
	r = strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonObject, violations[0].Message)
	require.Equal(t, CodeNotJsonObject, violations[0].Codes[0])
}

func TestValidateString(t *testing.T) {
	validator := &Validator{
		IgnoreUnknownProperties: true,
		DisallowObject:          false,
		AllowArray:              true,
	}

	ok, violations, obj := validator.ValidateString(`{"foo": { "bar": "baz" }}`)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	validator.IgnoreUnknownProperties = false
	ok, violations, obj = validator.ValidateString(`{"foo": { "bar": "baz" }}`)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))

	require.False(t, violations[0].BadRequest)

	ok, violations, obj = validator.ValidateString(`NOT JSON`)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.True(t, violations[0].BadRequest)
}

type intoTestStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
	Sub struct {
		Foo    string `json:"foo"`
		SubSub struct {
			Foo string `json:"foo"`
			Bar int    `json:"bar"`
		} `json:"subSub"`
	} `json:"sub"`
}

var validatorForInto = MustCompileValidatorFor(
	intoTestStruct{},
	&ValidatorForOptions{
		IgnoreUnknownProperties: false,
	})

func TestValidateReaderInto(t *testing.T) {
	r := strings.NewReader(`{
		"foo": "blah",
		"bar": 1,
		"sub": {
			"foo": "blah blah",
			"subSub": {
				"foo": "blah blah blah",
				"bar": 2
			}
		}
	}`)
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	// check that it got read into...
	require.Equal(t, "blah", myObj.Foo)
	require.Equal(t, 1, myObj.Bar)
	require.Equal(t, "blah blah", myObj.Sub.Foo)
	require.Equal(t, "blah blah blah", myObj.Sub.SubSub.Foo)
	require.Equal(t, 2, myObj.Sub.SubSub.Bar)
}

func TestValidateStringInto(t *testing.T) {
	s := `{
		"foo": "blah",
		"bar": 1,
		"sub": {
			"foo": "blah blah",
			"subSub": {
				"foo": "blah blah blah",
				"bar": 2
			}
		}
	}`
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateStringInto(s, myObj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	// check that it got read into...
	require.Equal(t, "blah", myObj.Foo)
	require.Equal(t, 1, myObj.Bar)
	require.Equal(t, "blah blah", myObj.Sub.Foo)
	require.Equal(t, "blah blah blah", myObj.Sub.SubSub.Foo)
	require.Equal(t, 2, myObj.Sub.SubSub.Bar)
}

func TestRequestValidateInto(t *testing.T) {
	body := strings.NewReader(`{
		"foo": "blah",
		"bar": 1,
		"sub": {
			"foo": "blah blah",
			"subSub": {
				"foo": "blah blah blah",
				"bar": 2
			}
		}
	}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	// check that it got read into...
	require.Equal(t, "blah", myObj.Foo)
	require.Equal(t, 1, myObj.Bar)
	require.Equal(t, "blah blah", myObj.Sub.Foo)
	require.Equal(t, "blah blah blah", myObj.Sub.SubSub.Foo)
	require.Equal(t, 2, myObj.Sub.SubSub.Bar)
}

func TestRequestValidateIntoFailsWithBadJsonBody(t *testing.T) {
	body := strings.NewReader(`NOT JSON`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnableToDecodeRequest, violations[0].Message)
	require.Equal(t, CodeUnableToDecodeRequest, violations[0].Codes[0])
	require.Nil(t, obj)
}

func TestRequestValidateIntoFailsWithEmptyBody(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgRequestBodyEmpty, violations[0].Message)
	require.Equal(t, CodeRequestBodyEmpty, violations[0].Codes[0])
	require.Nil(t, obj)
}

func TestRequestValidateIntoFailsWithErrorReader(t *testing.T) {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Body = &mockErrorReaderCloser{}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgErrorReading, violations[0].Message)
	require.Equal(t, CodeErrorReading, violations[0].Codes[0])
	require.Nil(t, obj)
}

func TestRequestValidateIntoFailsWithValidation(t *testing.T) {
	body := strings.NewReader(`{"xxx": "unexpected property"}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnknownProperty, violations[0].Message)
	require.Equal(t, "xxx", violations[0].Property)
	require.NotNil(t, obj)
}

type unknownStruct struct {
	Foo string `json:"foo"`
}

func TestRequestValidateIntoFailsWithWhenIntoStructDifferent(t *testing.T) {
	body := strings.NewReader(`{
		"foo": "blah",
		"bar": 1,
		"sub": {
			"foo": "blah blah",
			"subSub": {
				"foo": "blah blah blah",
				"bar": 2
			}
		}
	}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	myObj := &unknownStruct{}
	ok, violations, obj := validatorForInto.RequestValidateInto(req, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgErrorUnmarshall, violations[0].Message)
	require.Equal(t, CodeErrorUnmarshall, violations[0].Codes[0])
	require.NotNil(t, obj)
}

func TestValidateReaderIntoFailsWithBadJson(t *testing.T) {
	r := strings.NewReader(`BAD JSON`)
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnableToDecode, violations[0].Message)
	require.Equal(t, CodeUnableToDecode, violations[0].Codes[0])
	require.Nil(t, obj)
}

func TestValidateReaderIntoFailsValidation(t *testing.T) {
	r := strings.NewReader(`{"xxx": "unexpected property"}`)
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnknownProperty, violations[0].Message)
	require.Equal(t, "xxx", violations[0].Property)
	require.NotNil(t, obj)
}

func TestValidateReaderIntoFailsWhenIntoStructDifferent(t *testing.T) {
	r := strings.NewReader(`{
		"foo": "blah",
		"bar": 1,
		"sub": {
			"foo": "blah blah",
			"subSub": {
				"foo": "blah blah blah",
				"bar": 2
			}
		}
	}`)
	myObj := &unknownStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgErrorUnmarshall, violations[0].Message)
	require.Equal(t, CodeErrorUnmarshall, violations[0].Codes[0])
	require.NotNil(t, obj)
}

func TestValidatorOrderedPropertyChecking(t *testing.T) {
	v := &Validator{
		OrderedPropertyChecks: true,
		Properties: Properties{
			"aaa": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
			},
			"bbb": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
			},
			"ccc": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
			},
		},
	}
	obj := jsonObject(`{"aaa": 1, "ccc": false, "bbb": null}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	require.Equal(t, "aaa", violations[0].Property)
	require.Equal(t, "bbb", violations[1].Property)
	require.Equal(t, "ccc", violations[2].Property)

	// now change the ordering...
	v.Properties["ccc"].Order = 0
	v.Properties["bbb"].Order = 1
	v.Properties["aaa"].Order = 2
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	require.Equal(t, "ccc", violations[0].Property)
	require.Equal(t, "bbb", violations[1].Property)
	require.Equal(t, "aaa", violations[2].Property)

	// and check stop on first works with ordering...
	v.StopOnFirst = true
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "ccc", violations[0].Property)

	// and validate with no properties (for coverage!)...
	v.StopOnFirst = false
	obj = jsonObject(`{}`)
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
}

type unMarshallable struct {
	Name string `json:"name" v8n:"notNull,mandatory,constraints:[StringNoControlCharacters{},StringLength{Minimum: 1, Maximum: 255}]"`
	Age  int    `json:"age" v8n:"type:Integer,notNull,mandatory,constraint:PositiveOrZero{}"`
}

var unMarshallableValidator *Validator = nil

func init() {
	unMarshallableValidator = MustCompileValidatorFor(unMarshallable{}, nil)
}
func (s *unMarshallable) UnmarshalJSON(data []byte) error {
	// if we don't create a copy of the struct - json.Decoder goes into a loop and stack overflows
	type copyStruct unMarshallable
	cpy := copyStruct{}
	err := unMarshallableValidator.ValidateInto(data, &cpy)
	if err != nil {
		return err
	}
	*s = unMarshallable(cpy)
	return nil
}

func TestUnMarshallableValidator(t *testing.T) {
	un := &unMarshallable{}
	data := []byte(`{"name":"charlie","age":10}`)
	err := json.Unmarshal(data, un)
	require.Nil(t, err)
	require.Equal(t, "charlie", un.Name)
	require.Equal(t, 10, un.Age)

	data = []byte(`{"name":null,"age":-1}`)
	err = json.Unmarshal(data, un)
	require.NotNil(t, err)
	require.Equal(t, msgPositiveOrZero, err.Error()) // messages get sorted!
	vErr, ok := err.(*ValidationError)
	require.True(t, ok)
	require.NotNil(t, vErr)
	require.False(t, vErr.IsBadRequest)
	require.Equal(t, 2, len(vErr.Violations))
	require.Equal(t, "age", vErr.Violations[0].Property)
	require.Equal(t, "", vErr.Violations[0].Path)
	require.Equal(t, msgPositiveOrZero, vErr.Violations[0].Message)
	require.Equal(t, "name", vErr.Violations[1].Property)
	require.Equal(t, "", vErr.Violations[1].Path)
	require.Equal(t, msgValueCannotBeNull, vErr.Violations[1].Message)

	data = []byte(`null`)
	err = json.Unmarshal(data, un)
	require.NotNil(t, err)
	require.Equal(t, msgNotJsonNull, err.Error())
	vErr, ok = err.(*ValidationError)
	require.True(t, ok)
	require.NotNil(t, vErr)
	require.True(t, vErr.IsBadRequest)
	require.Equal(t, 0, len(vErr.Violations))
}

func TestValidatorMeetsWhenConditions(t *testing.T) {
	testIgnore := false
	v := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if testIgnore {
					vcx.SetCondition("TEST_IGNORE")
				}
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringNotEmpty{},
				},
				WhenConditions: []string{"!TEST_IGNORE"},
			},
		},
	}
	obj := jsonObject(`{}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgMissingProperty, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)

	testIgnore = true
	ok, violations = v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// and again with ordered
	v.OrderedPropertyChecks = true
	testIgnore = false
	ok, _ = v.Validate(obj)
	require.False(t, ok)
	testIgnore = true
	ok, _ = v.Validate(obj)
	require.True(t, ok)
}

func TestValidatorMeetsUnwantedConditions(t *testing.T) {
	testUnwanted := true
	v := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if testUnwanted {
					vcx.SetCondition("TEST_UNWANTED")
				}
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringNotEmpty{},
				},
				UnwantedConditions: []string{"TEST_UNWANTED"},
			},
		},
	}
	obj := jsonObject(`{"foo": "bar"}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnwantedProperty, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, CodeUnwantedProperty, violations[0].Codes[0])

	testUnwanted = false
	ok, violations = v.Validate(obj)
	require.True(t, ok)

	// and again with ordered
	v.OrderedPropertyChecks = true
	testUnwanted = true
	ok, _ = v.Validate(obj)
	require.False(t, ok)
	testUnwanted = false
	ok, _ = v.Validate(obj)
	require.True(t, ok)
}

func TestPresumptiveOrderedPropertyChecks(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Order: 0,
			},
		},
	}
	require.False(t, v.OrderedPropertyChecks)
	require.False(t, v.IsOrderedPropertyChecks())

	v.OrderedPropertyChecks = true
	require.True(t, v.IsOrderedPropertyChecks())

	v.OrderedPropertyChecks = false
	v.Properties["foo"].Order = -1
	require.False(t, v.OrderedPropertyChecks)
	require.True(t, v.IsOrderedPropertyChecks())
}

type mockErrorReader struct {
}

func (m *mockErrorReader) Read(p []byte) (int, error) {
	return 0, errors.New("whoops")
}

type mockErrorReaderCloser struct {
}

func (m *mockErrorReaderCloser) Read(p []byte) (int, error) {
	return 0, errors.New("whoops")
}
func (m *mockErrorReaderCloser) Close() error {
	return nil
}

func TestValidateReaderIntoFailsWithErrorReader(t *testing.T) {
	r := &mockErrorReader{}
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgErrorReading, violations[0].Message)
	require.Equal(t, CodeErrorReading, violations[0].Codes[0])
	require.Nil(t, obj)
}

func TestValidatorPolymorphism(t *testing.T) {
	catProperties := Properties{
		"body": {
			Type:      JsonString,
			Mandatory: true,
			NotNull:   true,
			Constraints: Constraints{
				&StringValidToken{
					Tokens: []string{"foreign", "semi-foreign", "moderate", "cobby", "lean", "muscular", "large", "normal"},
				},
			},
		},
		"coatLength": {
			Type:      JsonString,
			Mandatory: true,
			NotNull:   true,
			Constraints: Constraints{
				&StringValidToken{
					Tokens: []string{"short", "semi-long", "long", "rex", "hairless"},
				},
			},
		},
	}
	dogProperties := Properties{
		"breed": {
			Type:      JsonString,
			Mandatory: true,
			NotNull:   true,
		},
	}
	rabbitProperties := Properties{
		"isFluffy": {
			Type:      JsonBoolean,
			Mandatory: true,
			NotNull:   true,
		},
		"isFriendly": {
			Type:      JsonBoolean,
			Mandatory: true,
			NotNull:   true,
		},
	}
	petValidator := &Validator{
		Constraints: Constraints{
			&SetConditionProperty{
				PropertyName: "petType",
			},
		},
		Properties: Properties{
			"petType": {
				Type:      JsonString,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&StringValidToken{
						Tokens: []string{"cat", "dog", "rabbit"},
					},
				},
			},
		},
		ConditionalVariants: ConditionalVariants{
			{
				Properties:     catProperties,
				WhenConditions: []string{"cat"},
			},
			{
				Properties:     dogProperties,
				WhenConditions: []string{"dog"},
			},
			{
				Properties:     rabbitProperties,
				WhenConditions: []string{"rabbit"},
			},
		},
	}

	obj := jsonObject(`{
		"petType": "dog",
		"breed": "Staffordshire Bull Terrier"
	}`)
	ok, violations := petValidator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj = jsonObject(`{
		"petType": "dog",
		"body": "cobby"
	}`)
	ok, violations = petValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, msgInvalidProperty, violations[0].Message)
	require.Equal(t, "body", violations[0].Property)
	require.Equal(t, CodeInvalidProperty, violations[0].Codes[0])
	require.Equal(t, msgMissingProperty, violations[1].Message)
	require.Equal(t, "breed", violations[1].Property)

	obj = jsonObject(`{
		"petType": "cat",
		"body": "cobby"
	}`)
	ok, violations = petValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgMissingProperty, violations[0].Message)
	require.Equal(t, "coatLength", violations[0].Property)

	obj = jsonObject(`{
		"petType": "cat",
		"body": "cobby",
		"coatLength": "short"
	}`)
	ok, _ = petValidator.Validate(obj)
	require.True(t, ok)

	obj = jsonObject(`{
		"petType": "rabbit",
		"body": "cobby",
		"coatLength": "short"
	}`)
	ok, violations = petValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))
}

func TestConditionalVariantsStopOnFirst(t *testing.T) {
	constraintPasses := false
	v := &Validator{
		StopOnFirst:           true,
		OrderedPropertyChecks: true,
		Constraints: Constraints{
			&SetConditionProperty{PropertyName: "foo"},
		},
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringValidToken{Tokens: []string{"bar"}},
				},
			},
			"baz": {
				Type:    JsonString,
				NotNull: true,
			},
		},
		ConditionalVariants: ConditionalVariants{
			{
				WhenConditions: []string{"bar"},
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						return constraintPasses, "constraint fail"
					}, ""),
				},
				Properties: Properties{
					"bar": {
						Type:    JsonString,
						NotNull: true,
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": "bar",
		"bar": null,
		"baz": 1,
		"unknown": true
	}`)
	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "constraint fail", violations[0].Message)

	constraintPasses = true
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnknownProperty, violations[0].Message)
	require.Equal(t, "unknown", violations[0].Property)

	obj2 := jsonObject(`{
		"foo": "bar",
		"bar": null,
		"baz": 1
	}`)
	ok, violations = v.Validate(obj2)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueCannotBeNull, violations[0].Message)
	require.Equal(t, "bar", violations[0].Property)

	// and without stop on first...
	constraintPasses = false
	v.StopOnFirst = false
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))
}

func TestConditionalNestedVariants(t *testing.T) {
	v := &Validator{
		Constraints: Constraints{
			&SetConditionProperty{PropertyName: "foo"},
		},
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringValidToken{Tokens: []string{"bar"}},
				},
			},
		},
		ConditionalVariants: ConditionalVariants{
			{
				WhenConditions: []string{"bar"},
				Constraints: Constraints{
					&SetConditionProperty{PropertyName: "bar"},
				},
				Properties: Properties{
					"bar": {
						Type:      JsonString,
						NotNull:   true,
						Mandatory: true,
					},
				},
				ConditionalVariants: ConditionalVariants{
					{
						WhenConditions: []string{"baz"},
						Properties: Properties{
							"baz": {
								Type:      JsonString,
								NotNull:   true,
								Mandatory: true,
							},
						},
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": "bar",
		"bar": "baz",
		"baz": "xxx"
	}`)
	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = "xxx"
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, msgUnknownProperty, violations[0].Message)
	require.Equal(t, "bar", violations[0].Property)
	require.Equal(t, msgUnknownProperty, violations[1].Message)
	require.Equal(t, "baz", violations[1].Property)
	require.Equal(t, "foo", violations[2].Property)
}

func TestConditionsSetFromRequest(t *testing.T) {
	var grabVcx *ValidatorContext
	validator := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
				grabVcx = vcx
				return true, ""
			}, ""),
		},
	}

	body := strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, _, _ := validator.RequestValidate(req)
	require.True(t, ok)
	require.NotNil(t, grabVcx)

	stackItem := grabVcx.currentStackItem()
	require.Equal(t, 3, len(stackItem.conditions))
	require.True(t, stackItem.conditions["METHOD_POST"])
	require.True(t, stackItem.conditions["LANG_en"])
	require.True(t, stackItem.conditions["RGN_"])

	// and set language plus region...
	body = strings.NewReader(`{}`)
	req, err = http.NewRequest("PUT", "", body)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Language", "en-GB;q=0.5, fr-CA, en")
	grabVcx = nil
	ok, _, _ = validator.RequestValidate(req)
	require.True(t, ok)
	require.NotNil(t, grabVcx)

	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 3, len(stackItem.conditions))
	require.True(t, stackItem.conditions["METHOD_PUT"])
	require.True(t, stackItem.conditions["LANG_fr"])
	require.True(t, stackItem.conditions["RGN_CA"])

}

func TestValidatorInitialConditionsSet(t *testing.T) {
	var grabVcx *ValidatorContext
	validator := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
				grabVcx = vcx
				return true, ""
			}, ""),
		},
	}

	body := strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, _, _ := validator.RequestValidate(req, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem := grabVcx.currentStackItem()
	require.Equal(t, 5, len(stackItem.conditions)) // 3 set from request
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	ok, _ = validator.Validate(map[string]interface{}{}, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	ok, _ = validator.ValidateArrayOf([]interface{}{map[string]interface{}{}}, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	body = strings.NewReader(`{}`)
	ok, _, _ = validator.ValidateReader(body, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	data := []byte(`{}`)
	myObj := &intoTestStruct{}
	err = validator.ValidateInto(data, myObj, "FOO", "BAR")
	require.Nil(t, err)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	body = strings.NewReader(`{}`)
	ok, _, _ = validator.ValidateReaderInto(body, myObj, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	ok, _, _ = validator.ValidateString(`{}`, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])

	grabVcx = nil
	ok, _, _ = validator.ValidateStringInto(`{}`, myObj, "FOO", "BAR")
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	stackItem = grabVcx.currentStackItem()
	require.Equal(t, 2, len(stackItem.conditions))
	require.True(t, stackItem.conditions["FOO"])
	require.True(t, stackItem.conditions["BAR"])
}

func TestDifferentPOSTorPATCHMandatoryPropertyRequirements(t *testing.T) {
	validator := &Validator{
		Constraints: Constraints{
			&ConditionalConstraint{
				When: Conditions{"METHOD_PATCH"},
				Constraint: &Length{
					Minimum: 1,
					Message: "PATCH Object must not be empty",
				},
			},
		},
		Properties: Properties{
			"givenName": {
				Type:          JsonString,
				Mandatory:     true,
				MandatoryWhen: Conditions{"METHOD_POST"},
			},
			"familyName": {
				Type:          JsonString,
				Mandatory:     true,
				MandatoryWhen: Conditions{"METHOD_POST"},
			},
			"email": {
				Type:          JsonString,
				Mandatory:     true,
				MandatoryWhen: Conditions{"METHOD_POST"},
			},
		},
	}

	body := strings.NewReader(`{}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, _ := validator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))

	body = strings.NewReader(`{}`)
	req, err = http.NewRequest("PATCH", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, _ = validator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "PATCH Object must not be empty", violations[0].Message)

	body = strings.NewReader(`{"email": 1}`)
	req, err = http.NewRequest("PATCH", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, violations, _ = validator.RequestValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "email", violations[0].Property)

	body = strings.NewReader(`{"email": "foo@example.com"}`)
	req, err = http.NewRequest("PATCH", "", body)
	if err != nil {
		t.Fatal(err)
	}
	ok, _, _ = validator.RequestValidate(req)
	require.True(t, ok)
}

/* Utility structs & functions */
type myCustomConstraint struct {
	expectedPath string
}

func (c *myCustomConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return false, c.expectedPath
}
func (c *myCustomConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

func buildFooPropertyTypeValidator(jsonType JsonType, notNull bool) *Validator {
	return &Validator{
		Properties: Properties{
			"foo": {
				Type:      jsonType,
				NotNull:   notNull,
				Mandatory: true,
			},
		},
	}
}

func jsonObject(jsonStr string, flags ...bool) map[string]interface{} {
	r := strings.NewReader(jsonStr)
	decoder := json.NewDecoder(r)
	if len(flags) > 0 {
		if flags[0] {
			decoder.UseNumber()
		}
	}
	result := map[string]interface{}{}
	if err := decoder.Decode(&result); err != nil {
		panic(err)
	}
	return result
}

func jsonArray(jsonStr string, flags ...bool) []interface{} {
	r := strings.NewReader(jsonStr)
	decoder := json.NewDecoder(r)
	if len(flags) > 0 && flags[0] {
		decoder.UseNumber()
	}
	var result []interface{}
	if err := decoder.Decode(&result); err != nil {
		panic(err)
	}
	return result
}

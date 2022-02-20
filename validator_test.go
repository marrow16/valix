package valix

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"sort"
	"strings"
	"testing"
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
}

func TestMissingPropertyDetection(t *testing.T) {
	o := jsonObject(`{
		"name": "Bilbo"
	}`)
	ok, violations := personValidator.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(msgMissingProperty, "age"), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
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
	require.Equal(t, fmt.Sprintf(msgUnknownProperty, "unknown_property"), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
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
	sortViolationsByPathAndProperty(violations)
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
	require.Equal(t, fmt.Sprintf(msgArrayElementMustBeObject, 1), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
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
	require.Equal(t, fmt.Sprintf(messageMinMax, 2, 3), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
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
						return false, cc.GetMessage()
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
	require.Equal(t, messageValueMustBeObjectOrArray, violations[0].Message)
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
	require.Equal(t, messageValueMustBeObject, violations[0].Message)
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
	require.Equal(t, messageValueMustBeArray, violations[0].Message)
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
	require.Equal(t, messagePropertyObjectValidatorError, violations[0].Message)
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
	sortViolationsByPathAndProperty(violations)
	require.Equal(t, "person", violations[0].Path)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, messagePositiveOrZero, violations[0].Message)
	require.Equal(t, "person", violations[1].Path)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, fmt.Sprintf(messageStringMinMaxLen, 1, 255), violations[1].Message)

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
	sortViolationsByPathAndProperty(violations)
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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonString), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonNumber), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonBoolean), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonObject), violations[0].Message)

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
	require.Equal(t, fmt.Sprintf(messageValueExpectedType, JsonArray), violations[0].Message)

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
	require.Equal(t, messageValueCannotBeNull, violations[0].Message)
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
	require.Equal(t, messageNotEmptyString, violations[1].Message)
}

func TestValidatorManuallyAddedCurrentViolation(t *testing.T) {
	msg := "Something went wrong"
	v := Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
						vcx.AddViolationForCurrent(msg)
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
	require.Equal(t, messageNotEmptyString, violations[1].Message)
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

	r = strings.NewReader(`null`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonNull, violations[0].Message)

	r = strings.NewReader(`[{"foo": "bar"}]`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonArray, violations[0].Message)

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

	validator.AllowArray = false
	validator.DisallowObject = true
	r = strings.NewReader(`{"foo": { "bar": "baz" }}`)
	ok, violations, obj = validator.ValidateReader(r)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotJsonObject, violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(msgUnknownProperty, "xxx"), violations[0].Message)
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
	require.NotNil(t, obj)
}

func TestValidateReaderIntoFailsWithBadJson(t *testing.T) {
	r := strings.NewReader(`BAD JSON`)
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgUnableToDecode, violations[0].Message)
	require.Nil(t, obj)
}

func TestValidateReaderIntoFailsValidation(t *testing.T) {
	r := strings.NewReader(`{"xxx": "unexpected property"}`)
	myObj := &intoTestStruct{}
	ok, violations, obj := validatorForInto.ValidateReaderInto(r, myObj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(msgUnknownProperty, "xxx"), violations[0].Message)
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
	require.NotNil(t, obj)
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
	require.Nil(t, obj)
}

type myCustomConstraint struct {
	expectedPath string
}

func (c *myCustomConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return false, c.expectedPath
}
func (c *myCustomConstraint) GetMessage() string {
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

func sortViolationsByPathAndProperty(violations []*Violation) {
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].Path == violations[j].Path {
			return violations[i].Property < violations[j].Property
		}
		return violations[i].Path < violations[j].Path
	})
}

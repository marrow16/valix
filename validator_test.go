package valix

import (
	"encoding/json"
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
	Properties: map[string]*PropertyValidator{
		"name": {
			PropertyType: PropertyType.String,
			NotNull:      true,
			Mandatory:    true,
			Constraints: []Constraint{
				&StringLengthConstraint{Minimum: 1, Maximum: 255},
			},
		},
		"age": {
			PropertyType: PropertyType.Int,
			NotNull:      true,
			Mandatory:    true,
			Constraints: []Constraint{
				&PositiveOrZeroConstraint{},
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

func TestMissingPropertyDetection(t *testing.T) {
	o := jsonObject(`{
		"name": "Bilbo"
	}`)
	ok, violations := personValidator.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageMissingProperty, "age"), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(MessageUnknownProperty, "unknown_property"), violations[0].Message)
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
	require.Equal(t, fmt.Sprintf(MessageArrayElementMustBeObject, 1), violations[0].Message)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "", violations[0].Property)
}

func TestValidatorWithObjectConstraint(t *testing.T) {
	v := Validator{
		IgnoreUnknownProperties: true,
		Constraints: []Constraint{
			&LengthConstraint{Minimum: 2, Maximum: 3},
		},
	}
	o := jsonObject(`{
		"foo": null
	}`)
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(messageValueMinMax, 2, 3), violations[0].Message)
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
	require.NotNil(t, obj)
}

func TestRequestValidationUsingJsonNumber(t *testing.T) {
	validator := Validator{
		UseNumber: true,
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.Number,
				Constraints:  []Constraint{&PositiveConstraint{}},
			},
			"bar": {
				PropertyType: PropertyType.Number,
				Constraints:  []Constraint{&NegativeConstraint{}},
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
	require.Equal(t, MessageRequestBodyNotJsonNull, violations[0].Message)
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
	require.Equal(t, MessageRequestBodyEmpty, violations[0].Message)
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
	require.Equal(t, MessageRequestBodyNotJsonArray, violations[0].Message)
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
	require.Equal(t, MessageRequestBodyExpectedJsonArray, violations[0].Message)
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
	require.Equal(t, MessageRequestBodyExpectedJsonObject, violations[0].Message)
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
	require.Equal(t, MessageRequestBodyNotJsonObject, violations[0].Message)
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
	require.Equal(t, MessageUnableToDecode, violations[0].Message)
	require.Nil(t, obj)
}

func TestValidatorStopsOnConstraint(t *testing.T) {
	v := Validator{
		IgnoreUnknownProperties: false,
		Constraints: []Constraint{
			NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
				ctx.Stop()
				return true, ""
			}, ""),
		},
		Properties: map[string]*PropertyValidator{
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.Boolean,
				Constraints: []Constraint{
					NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
						ctx.Stop()
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				ObjectValidator: &Validator{
					AllowNull:      false,
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
	require.Equal(t, MessageValueMustBeObjectOrArray, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForObjectOnly(t *testing.T) {
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				ObjectValidator: &Validator{
					AllowNull:      false,
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
	require.Equal(t, MessageValueMustBeObject, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForArrayOnly(t *testing.T) {
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				ObjectValidator: &Validator{
					AllowNull:      false,
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
	require.Equal(t, MessageValueMustBeArray, violations[0].Message)
}

func TestPropertyValueObjectValidatorFailsForNeitherObjectNorArray(t *testing.T) {
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				ObjectValidator: &Validator{
					AllowNull:      false,
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
	require.Equal(t, MessagePropertyObjectValidatorError, violations[0].Message)
}

func TestSubPropertyValidation(t *testing.T) {
	v := Validator{
		Properties: map[string]*PropertyValidator{
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
		Properties: map[string]*PropertyValidator{
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
	require.Equal(t, messageValuePositiveOrZero, violations[0].Message)
	require.Equal(t, "person", violations[1].Path)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, fmt.Sprintf(messageValueMinMax, 1, 255), violations[1].Message)

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
	v := buildFooPropertyTypeValidator(PropertyType.String, true)
	o := jsonObject(`{"foo": 1}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageValueExpectedType, PropertyType.String), violations[0].Message)

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
	v := buildFooPropertyTypeValidator(PropertyType.Number, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageValueExpectedType, PropertyType.Number), violations[0].Message)

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
	v := buildFooPropertyTypeValidator(PropertyType.Boolean, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageValueExpectedType, PropertyType.Boolean), violations[0].Message)

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
	v := buildFooPropertyTypeValidator(PropertyType.Object, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageValueExpectedType, PropertyType.Object), violations[0].Message)

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
	v := buildFooPropertyTypeValidator(PropertyType.Array, true)
	o := jsonObject(`{"foo": "abc"}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(MessageValueExpectedType, PropertyType.Array), violations[0].Message)

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
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.String,
				NotNull:      true,
				Constraints: []Constraint{
					NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
						ctx.Stop()
						return true, ""
					}, ""),
					// the following constraint should never be checked (because the prev constraint stops)
					&StringNotEmptyConstraint{},
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
	require.Equal(t, MessageValueCannotBeNull, violations[0].Message)
}

func TestValidatorManuallyAddedViolation(t *testing.T) {
	msg := "Something went wrong"
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.String,
				Constraints: []Constraint{
					NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
						ctx.AddViolation(NewViolation("", "", msg))
						return true, ""
					}, ""),
					&StringNotEmptyConstraint{},
				},
			},
		},
	}
	o := jsonObject(`{"foo": ""}`)

	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, msg, violations[0].Message)
	require.Equal(t, messageValueNotEmptyString, violations[1].Message)
}

func TestValidatorManuallyAddedCurrentViolation(t *testing.T) {
	msg := "Something went wrong"
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.String,
				Constraints: []Constraint{
					NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
						ctx.AddViolationForCurrent(msg)
						return true, ""
					}, ""),
					&StringNotEmptyConstraint{},
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
	require.Equal(t, messageValueNotEmptyString, violations[1].Message)
	require.Equal(t, "", violations[1].Path)
	require.Equal(t, "foo", violations[1].Property)
}

func TestValidatorContextPathing(t *testing.T) {
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				Constraints: []Constraint{&myCustomConstraint{expectedPath: ""}},
				ObjectValidator: &Validator{
					Properties: map[string]*PropertyValidator{
						"bar": {
							Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo"}},
							ObjectValidator: &Validator{
								Properties: map[string]*PropertyValidator{
									"baz": {
										Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo.bar"}},
										ObjectValidator: &Validator{
											Properties: map[string]*PropertyValidator{
												"qux": {
													Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo.bar.baz"}},
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				Constraints: []Constraint{&myCustomConstraint{expectedPath: ""}},
				ObjectValidator: &Validator{
					AllowArray: true,
					Properties: map[string]*PropertyValidator{
						"bar": {
							Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo[0]"}},
							ObjectValidator: &Validator{
								AllowArray: true,
								Properties: map[string]*PropertyValidator{
									"baz": {
										Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo[0].bar[0]"}},
										ObjectValidator: &Validator{
											AllowArray: true,
											Properties: map[string]*PropertyValidator{
												"qux": {
													Constraints: []Constraint{&myCustomConstraint{expectedPath: "foo[0].bar[0].baz[0]"}},
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.String,
				Constraints: []Constraint{
					NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
						ctx.CeaseFurther()
						return true, ""
					}, ""),
					// these constraints should not be run because of the CeaseFurther (above)
					&StringNotEmptyConstraint{},
				},
				ObjectValidator: &Validator{
					Constraints: []Constraint{
						NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
							// this constraint should not be run because of the CeaseFurther (above)
							return false, "Never run"
						}, ""),
					},
				},
			},
			"bar": {
				PropertyType: PropertyType.String,
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

type myCustomConstraint struct {
	expectedPath string
}

func (c *myCustomConstraint) Validate(value interface{}, ctx *Context) (bool, string) {
	return false, c.expectedPath
}
func (c *myCustomConstraint) GetMessage() string {
	return ""
}

func buildFooPropertyTypeValidator(propertyType string, notNull bool) *Validator {
	return &Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: propertyType,
				NotNull:      notNull,
				Mandatory:    true,
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

package valix

import (
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/text/unicode/norm"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	vJson = `{
		"allowArray": true,
		"allowNullJson": true,
		"conditionalVariants": [
			{
				"constraints": [
					{
						"fields": {
							"Mapping": null,
							"Prefix": "",
							"PropertyName": "foo"
						},
						"name": "SetConditionProperty"
					}
				],
				"properties": {
					"foo": {
						"mandatory": true,
						"notNull": false,
						"order": 0,
						"type": "string"
					}
				},
				"whenConditions": [
					"cv_foo"
				]
			},
			{
				"constraints": [
					{
						"fields": {
							"Mapping": null,
							"Prefix": "",
							"PropertyName": "foo"
						},
						"name": "SetConditionProperty"
					}
				],
				"properties": {
					"foo": {
						"mandatory": true,
						"notNull": false,
						"order": 0,
						"type": "string"
					}
				},
				"whenConditions": [
					"!cv_foo"
				]
			}
		],
		"constraints": [
			{
				"fields": {
					"Maximum": 2,
					"Message": "message 1",
					"Minimum": 1,
					"Stop": true
				},
				"name": "Length"
			},
			{
				"fields": {
					"Constraints": [
						{
							"fields": {
								"Maximum": 4,
								"Message": "message 2",
								"Minimum": 3,
								"Stop": true
							},
							"name": "Range"
						}
					],
					"Message": "message 3",
					"Stop": true
				},
				"name": "ConstraintSet"
			}
		],
		"disallowObject": true,
		"ignoreUnknownProperties": true,
		"oasInfo": {
			"deprecated": true,
			"description": "",
			"example": "",
			"format": "",
			"title": ""
		},
		"orderedPropertyChecks": true,
		"properties": {
			"arr": {
				"constraints": [
					{
						"fields": {
							"Ancestry": 16,
							"Constraint": {
								"fields": {
									"CaseInsensitive": false,
									"Message": "",
									"PropertyName": "foo",
									"Stop": false
								},
								"name": "StringGreaterThanOther"
							},
							"When": "%2"
						},
						"name": "ArrayConditionalConstraint"
					}
				],
				"mandatory": false,
				"notNull": false,
				"order": 0,
				"requiredWithMessage": "",
				"type": "string",
				"unwantedWithMessage": ""
			},
			"bar": {
				"mandatory": true,
				"notNull": true,
				"constraints": [
					{
						"fields": {
							"ExclusiveMax": false,
							"ExclusiveMin": false,
							"Maximum": 0,
							"Message": "",
							"Minimum": 1,
							"Stop": false
						},
						"name": "Length",
						"whenConditions": [
							"check_bar"
						]
					}
				],
				"objectValidator": {
					"allowArray": true,
					"allowNullJson": true,
					"disallowObject": true,
					"ignoreUnknownProperties": true,
					"orderedPropertyChecks": true,
					"properties": {},
					"stopOnFirst": true,
					"useNumber": true
				},
				"order": 0,
				"type": "object"
			},
			"foo": {
				"constraints": [
					{
						"fields": {
							"Message": "message 4",
							"Stop": true
						},
						"name": "StringNotEmpty"
					},
					{
						"fields": {
							"Message": "message 5",
							"Regexp": "^([A-Z]+)$",
							"Stop": true
						},
						"name": "StringPattern"
					}
				],
				"mandatory": true,
				"mandatoryWhen": ["want_foo"],
				"requiredWith": "(foo || bar) && !(foo && bar)",
				"requiredWithMessage": "only sometimes",
				"notNull": true,
				"oasInfo": {
					"deprecated": true,
					"description": "",
					"example": "",
					"format": "",
					"title": ""
				},
				"order": 1,
				"type": "string",
				"unwantedConditions": [
					"!when_foo"
				],
				"whenConditions": [
					"when_foo"
				]
			},
			"baz": {
				"constraints": [
					{
						"fields": {
							"AllowNull": true,
							"NameConstraints": [
								{
									"fields": {
										"Message": "",
										"Stop": false,
										"Tokens": [
											"FOO2",
											"BAR2",
											"BAZ2"
										]
									},
									"name": "StringValidToken"
								}
							],
							"ObjectValidator": {
								"allowArray": false,
								"allowNullJson": false,
								"disallowObject": false,
								"ignoreUnknownProperties": false,
								"orderedPropertyChecks": false,
								"properties": {
									"foo3": {
										"mandatory": false,
										"notNull": false,
										"order": 0,
										"type": "string"
									}
								},
								"stopOnFirst": true,
								"useNumber": false
							}
						},
						"name": "VariablePropertyConstraint"
					}
				],
				"mandatory": true,
				"notNull": true,
				"order": 0,
				"type": "object"
			}
		},
		"stopOnFirst": true,
		"useNumber": true,
		"whenConditions": [
			"test"
		]
	}`
)

func TestValidatorUnmarshal(t *testing.T) {
	v := Validator{}
	err := json.Unmarshal([]byte(vJson), &v)
	require.Nil(t, err)
	require.True(t, v.AllowNullJson)
	require.True(t, v.AllowArray)
	require.True(t, v.DisallowObject)
	require.True(t, v.UseNumber)
	require.True(t, v.IgnoreUnknownProperties)
	require.True(t, v.OrderedPropertyChecks)
	require.True(t, v.StopOnFirst)
	require.True(t, v.OasInfo.Deprecated)
	require.Equal(t, 1, len(v.WhenConditions))
	require.Equal(t, "test", v.WhenConditions[0])
	require.Equal(t, 2, len(v.Constraints))
	constraint1 := v.Constraints[0].(*Length)
	require.Equal(t, 1, constraint1.Minimum)
	require.Equal(t, 2, constraint1.Maximum)
	require.True(t, constraint1.Stop)
	require.Equal(t, "message 1", constraint1.Message)
	constraint2 := v.Constraints[1].(*ConstraintSet)
	require.True(t, constraint2.Stop)
	require.Equal(t, "message 3", constraint2.Message)
	require.Equal(t, 1, len(constraint2.Constraints))
	constraint3 := constraint2.Constraints[0].(*Range)
	require.True(t, constraint3.Stop)
	require.Equal(t, "message 2", constraint3.Message)
	require.Equal(t, float64(3), constraint3.Minimum)
	require.Equal(t, float64(4), constraint3.Maximum)

	require.Equal(t, 2, len(v.ConditionalVariants))
	cv := v.ConditionalVariants[0]
	require.Equal(t, 1, len(cv.WhenConditions))
	require.Equal(t, "cv_foo", cv.WhenConditions[0])

	require.Equal(t, 4, len(v.Properties))
	foo := v.Properties["foo"]
	require.Equal(t, 1, foo.Order)
	require.True(t, foo.Mandatory)
	require.Equal(t, 1, len(foo.MandatoryWhen))
	require.True(t, foo.NotNull)
	require.Equal(t, JsonString, foo.Type)
	require.True(t, foo.OasInfo.Deprecated)
	require.Equal(t, 1, len(foo.WhenConditions))
	require.Equal(t, "when_foo", foo.WhenConditions[0])
	require.Equal(t, 1, len(foo.UnwantedConditions))
	require.Equal(t, "!when_foo", foo.UnwantedConditions[0])
	require.Equal(t, 2, len(foo.Constraints))
	fooConstraint1 := foo.Constraints[0].(*StringNotEmpty)
	require.True(t, fooConstraint1.Stop)
	require.Equal(t, "message 4", fooConstraint1.Message)
	fooConstraint2 := foo.Constraints[1].(*StringPattern)
	require.True(t, fooConstraint2.Stop)
	require.Equal(t, "message 5", fooConstraint2.Message)
	require.Equal(t, "^([A-Z]+)$", fooConstraint2.Regexp.String())
	require.Equal(t, 2, len(foo.RequiredWith))
	require.Equal(t, "only sometimes", foo.RequiredWithMessage)

	bar := v.Properties["bar"]
	require.Equal(t, 1, len(bar.Constraints))
	barConstraint := bar.Constraints[0].(*ConditionalConstraint)
	require.Equal(t, 1, len(barConstraint.When))
	require.True(t, bar.Mandatory)
	require.True(t, bar.NotNull)
	require.Equal(t, 0, bar.Order)
	require.Equal(t, JsonObject, bar.Type)
	require.NotNil(t, bar.ObjectValidator)
	require.Equal(t, 0, len(bar.ObjectValidator.Properties))
	require.Equal(t, 0, len(bar.ObjectValidator.WhenConditions))
	require.Equal(t, 0, len(bar.ObjectValidator.Constraints))
	require.True(t, bar.ObjectValidator.AllowArray)
	require.True(t, bar.ObjectValidator.DisallowObject)
	require.True(t, bar.ObjectValidator.AllowNullJson)
	require.True(t, bar.ObjectValidator.IgnoreUnknownProperties)
	require.True(t, bar.ObjectValidator.OrderedPropertyChecks)
	require.True(t, bar.ObjectValidator.StopOnFirst)
	require.True(t, bar.ObjectValidator.UseNumber)

	baz := v.Properties["baz"]
	require.Equal(t, 1, len(baz.Constraints))
	constraint4 := baz.Constraints[0].(*VariablePropertyConstraint)
	require.True(t, constraint4.AllowNull)
	require.Equal(t, 1, len(constraint4.NameConstraints))
	require.True(t, constraint4.ObjectValidator.StopOnFirst)
	require.Equal(t, 1, len(constraint4.ObjectValidator.Properties))

	arr := v.Properties["arr"]
	require.Equal(t, 1, len(arr.Constraints))
	constraint5 := arr.Constraints[0].(*ArrayConditionalConstraint)
	require.Equal(t, uint(16), constraint5.Ancestry)
	require.Equal(t, "%2", constraint5.When)
	constraint6 := constraint5.Constraint.(*StringGreaterThanOther)
	require.Equal(t, "foo", constraint6.PropertyName)
}

func TestValidatorValidation(t *testing.T) {
	v := &Validator{}
	ok, violations, _ := ValidatorValidator.ValidateStringInto(vJson, v)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.True(t, v.StopOnFirst)
}

func TestValidatorValidationWithBadJson(t *testing.T) {
	testCases := []map[string]interface{}{
		{},
		{
			ptyNameProperties: nil,
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": nil,
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": false,
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType: 1,
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType: "not valid",
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:    "string",
					ptyNameNotNull: "should be bool",
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: "should be bool",
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:        "string",
					ptyNameNotNull:     true,
					ptyNameMandatory:   true,
					ptyNameConstraints: "should be array",
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:        "string",
					ptyNameNotNull:     true,
					ptyNameMandatory:   true,
					ptyNameConstraints: []string{},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						"should be an object",
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						nil, // null objects not allowed
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name": true, // should be string
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name": nil, // cannot be null
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name": "unknown_constraint_name",
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name":             "ConstraintSet",
							"unknown_property": true,
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name":   "ConstraintSet",
							"fields": true, // should be an object
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name":   "ConstraintSet",
							"fields": nil, // cannot be null
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name": "StringNotEmpty",
							"fields": map[string]interface{}{
								"Message":      "Fooey",
								"UnknownField": true,
							},
						},
					},
				},
			},
		},
		{
			ptyNameProperties: map[string]interface{}{
				"foo": map[string]interface{}{
					ptyNameType:      "string",
					ptyNameNotNull:   true,
					ptyNameMandatory: true,
					ptyNameConstraints: []interface{}{
						map[string]interface{}{
							"name": "StringNotEmpty",
							"fields": map[string]interface{}{
								"Message": false, // this should be a string
							},
						},
					},
				},
			},
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Test case: %d", i+1), func(t *testing.T) {
			ok, violations := ValidatorValidator.Validate(tc)
			require.False(t, ok)
			require.Equal(t, 1, len(violations))
		})
	}
}

func TestGetConstraintFieldNames(t *testing.T) {
	constraint := &StringValidUuid{}
	fields := getConstraintFieldNames(constraint)
	require.Equal(t, 4, len(fields))
}

func TestGetConstraintFieldNamesOnNonStruct(t *testing.T) {
	var constraint nonStructConstraint = ""
	fields := getConstraintFieldNames(&constraint)
	require.Equal(t, 0, len(fields))
}

type nonStructConstraint string

func (n *nonStructConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return false, ""
}
func (n *nonStructConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

func TestGetConstraintFieldNamesWithJsonTags(t *testing.T) {
	constraint := &constraintWithJsonTags{}
	fields := getConstraintFieldNames(constraint)
	require.Equal(t, 2, len(fields))
	require.True(t, fields["foo"])
	require.True(t, fields["bar"])
}

type constraintWithJsonTags struct {
	Foo string `json:"foo"`
	Bar string `json:"bar,omitempty"`
}

func (c *constraintWithJsonTags) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return false, ""
}
func (c *constraintWithJsonTags) GetMessage(tcx I18nContext) string {
	return ""
}

func TestUnmarshalPropertiesType(t *testing.T) {
	j := `{
		"Properties": {
			"foo": {
				"Type": 0
			}
		}
	}`
	v := Validator{}
	err := json.Unmarshal([]byte(j), &v)
	require.Nil(t, err)

	j = `{
		"Properties": {
			"foo": {
				"Type": -1
			}
		}
	}`
	v = Validator{}
	err = json.Unmarshal([]byte(j), &v)
	require.NotNil(t, err)

	j = `{
		"Properties": {
			"foo": {
				"Type": "string"
			}
		}
	}`
	v = Validator{}
	err = json.Unmarshal([]byte(j), &v)
	require.Nil(t, err)

	j = `{
		"Properties": {
			"foo": {
				"Type": "unknown"
			}
		}
	}`
	v = Validator{}
	err = json.Unmarshal([]byte(j), &v)
	require.NotNil(t, err)
}

func TestUnmarshalConstraints(t *testing.T) {
	j := `{
		"Constraints": [
			{
				"name": "StringNotEmpty",
				"fields": {
					"Message": "Overridden message",
					"Stop": true
				}
			}
		]
	}`
	v := Validator{}
	err := json.Unmarshal([]byte(j), &v)
	require.Nil(t, err)
	require.Equal(t, 1, len(v.Constraints))
	constraint := v.Constraints[0].(*StringNotEmpty)
	require.True(t, constraint.Stop)
	require.Equal(t, "Overridden message", constraint.Message)

	j = `{
		"Constraints": "this should be an array!"
	}`
	v = Validator{}
	err = json.Unmarshal([]byte(j), &v)
	require.NotNil(t, err)

	j = `{
		"Constraints": [
			"this should be an object!"
		]
	}`
	v = Validator{}
	err = json.Unmarshal([]byte(j), &v)
	require.NotNil(t, err)
}

func TestUnmarshalWithUnmarshallableConstraint(t *testing.T) {
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	constraintsRegistry.registerNamed(true, "Unmarshallable", &unmarshallableConstraint{})
	j := `{
		"Constraints": [
			{
				"name": "Unmarshallable"
			}
		]
	}`
	v := Validator{}
	err := json.Unmarshal([]byte(j), &v)
	require.NotNil(t, err)
	require.Equal(t, "this constraint cannot be unmarshalled", err.Error())
}

type unmarshallableConstraint struct {
}

func (n *unmarshallableConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return false, ""
}
func (n *unmarshallableConstraint) GetMessage(tcx I18nContext) string {
	return ""
}
func (n *unmarshallableConstraint) UnmarshalJSON(data []byte) error {
	return errors.New("this constraint cannot be unmarshalled")
}

func TestConstraintSetUnmarshalJSON(t *testing.T) {
	js := `{}`
	constraint := &ConstraintSet{}
	err := json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)

	js = `[]`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Message": "foo",
		"Stop": true
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, "foo", constraint.Message)
	require.True(t, constraint.Stop)

	js = `{
		"Message": 1
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, constraintSetFieldMessage, "string"), err.Error())

	js = `{
		"Stop": "xxx"
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, constraintSetFieldStop, "bool"), err.Error())

	js = `{
		"Constraints": []
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, 0, len(constraint.Constraints))

	js = `{
		"Constraints": {}
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, constraintSetFieldConstraints, "array"), err.Error())

	js = `{
		"Constraints": [
			{
				"name": "UNKNOWN_CONSTRAINT_NAME"
			}
		]
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgUnknownNamedConstraint, "UNKNOWN_CONSTRAINT_NAME"), err.Error())

	js = `{
		"Constraints": [
			{
				"name": 1
			}
		]
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, "name", "string"), err.Error())

	js = `{
		"Constraints": [
			{
				"name": "StringNotEmpty"
			}
		]
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, 1, len(constraint.Constraints))

	js = `{
		"Constraints": [
			{
				"name": "StringNotEmpty",
				"fields": "yyyy"
			}
		]
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, "fields", "object"), err.Error())

	js = `{
		"Constraints": [
			{
				"name": "StringNotEmpty",
				"fields": {
					"Message": 1
				}
			}
		]
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"OneOf": "this should be boolean"
	}`
	constraint = &ConstraintSet{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, constraintSetFieldOneOf, "bool"), err.Error())
}

func TestStringPatternUnmarshalJSON(t *testing.T) {
	js := `{}`
	constraint := &StringPattern{}
	err := json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)

	js = `{
		"Message": "foo",
		"Stop": true,
		"Regexp": "^([A-Z]+)$"
	}`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, "foo", constraint.Message)
	require.True(t, constraint.Stop)
	require.Equal(t, "^([A-Z]+)$", constraint.Regexp.String())

	js = `{
		"Regexp": ")][("
	}`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Regexp": 123
	}`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Message": true
	}`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Stop": "true"
	}`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `[]`
	constraint = &StringPattern{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
}

func TestStringValidUnicodeNormalizationUnmarshalJSON(t *testing.T) {
	js := `{}`
	constraint := &StringValidUnicodeNormalization{}
	err := json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)

	js = `{
		"Form": "NFC",
		"Message": "foo",
		"Stop": true
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, "foo", constraint.Message)
	require.True(t, constraint.Stop)
	require.Equal(t, norm.NFC, constraint.Form)

	js = `{
		"Form": "NFD"
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, norm.NFD, constraint.Form)

	js = `{
		"Form": "NFKC"
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, norm.NFKC, constraint.Form)

	js = `{
		"Form": "NFKD"
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	require.Equal(t, norm.NFKD, constraint.Form)

	js = `{
		"Form": 123
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Form": "XXX"
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `[]`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Message": 123
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Stop": 123
	}`
	constraint = &StringValidUnicodeNormalization{}
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
}

func TestArrayConditionalConstraintUnmarshalJSON(t *testing.T) {
	js := `{}`
	constraint := &ArrayConditionalConstraint{}
	err := json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Constraint": nil
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)
	js = `{
		"Constraint": 1
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"When": 1,
		"Constraint": {
			"name": "StringNotEmpty"
		}
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"Ancestry": "xxx",
		"Constraint": {
			"name": "StringNotEmpty"
		}
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"When": "%2",
		"Ancestry": 16,
		"Constraint": {
			"name": "bad name"
		}
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.NotNil(t, err)

	js = `{
		"When": "%2",
		"Ancestry": 16,
		"Constraint": {
			"name": "StringGreaterThanOther",
			"fields": {
				"PropertyName": "foo"
			}
		}
	}`
	err = json.Unmarshal([]byte(js), constraint)
	require.Nil(t, err)
	wrapped := constraint.Constraint.(*StringGreaterThanOther)
	require.Equal(t, "foo", wrapped.PropertyName)
}

func TestValidateUnmarshallingConstraintWithNonStructConstraint(t *testing.T) {
	var constraint nonStructConstraint = ""
	_, err := validateUnmarshallingConstraint(&constraint, map[string]interface{}{}, nil)
	require.NotNil(t, err)
	require.Equal(t, msgConstraintNotStruct, err.Error())
}

func TestUnmarshalConstraintWithBadWhens(t *testing.T) {
	// check good first...
	obj := jsonObject(`{
		"fields": {
			"Minimum": 1
		},
		"name": "Length",
			"whenConditions": [
				"check_bar"
		]
	}`)
	constraint, err := unmarshalConstraint(obj)
	require.Nil(t, err)
	condConstraint, ok := constraint.(*ConditionalConstraint)
	require.True(t, ok)
	require.Equal(t, 1, len(condConstraint.When))
	require.Equal(t, "check_bar", condConstraint.When[0])
	innerConstraint, ok := condConstraint.Constraint.(*Length)
	require.True(t, ok)
	require.Equal(t, 1, innerConstraint.Minimum)

	// now test bad whens...
	obj = jsonObject(`{
		"fields": {
			"Minimum": 1
		},
		"name": "Length",
			"whenConditions": "should be an array"
	}`)
	constraint, err = unmarshalConstraint(obj)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, ptyNameWhenConditions, "array"), err.Error())

	obj = jsonObject(`{
		"fields": {
			"Minimum": 1
		},
		"name": "Length",
			"whenConditions": [1]
	}`)
	constraint, err = unmarshalConstraint(obj)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(errMsgFieldExpectedType, ptyNameWhenConditions, "array of strings"), err.Error())
}

func TestUnmarshalPropertyValidatorWithBadWithExprs(t *testing.T) {
	pv := &PropertyValidator{}
	js := `{
		"requiredWith": ""
	}`
	err := json.Unmarshal([]byte(js), pv)
	require.Nil(t, err)
	require.NotNil(t, pv.RequiredWith)
	require.Equal(t, 0, len(pv.RequiredWith))

	pv.RequiredWith = nil
	js = `{
		"requiredWith": null
	}`
	err = json.Unmarshal([]byte(js), pv)
	require.Nil(t, err)
	require.NotNil(t, pv.RequiredWith)
	require.Equal(t, 0, len(pv.RequiredWith))

	js = `{
		"requiredWith": 0
	}`
	err = json.Unmarshal([]byte(js), pv)
	require.NotNil(t, err)

	// with bad expression (missing boolean operators)...
	js = `{
		"requiredWith": "foo bar baz"
	}`
	err = json.Unmarshal([]byte(js), pv)
	require.NotNil(t, err)
	require.True(t, strings.Contains(err.Error(), fmt.Sprintf("at position %d", 4)))

	pv.RequiredWith = nil
	js = `{
		"requiredWith": "foo || !bar"
	}`
	err = json.Unmarshal([]byte(js), pv)
	require.Nil(t, err)
	require.NotNil(t, pv.RequiredWith)
	require.Equal(t, 2, len(pv.RequiredWith))
	opt := pv.RequiredWith[0].(*OtherProperty)
	require.Equal(t, "foo", opt.Name)
	opt = pv.RequiredWith[1].(*OtherProperty)
	require.Equal(t, "bar", opt.Name)
	require.True(t, opt.Not)
	require.Equal(t, Or, opt.Op)
}

func TestPropertyValidatorValidatorFailsWithBadWiths(t *testing.T) {
	obj := jsonObject(`{
		"requiredWith": null
	}`)
	ok, violations := PropertyValidatorValidator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj[ptyNameRequiredWith] = 0
	ok, violations = PropertyValidatorValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, "string"), violations[0].Message)

	obj[ptyNameRequiredWith] = "foo bar"
	ok, violations = PropertyValidatorValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.True(t, strings.Contains(violations[0].Message, fmt.Sprintf("at position %d", 4)))
}

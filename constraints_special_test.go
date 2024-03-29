package valix

import (
	"encoding/json"
	"fmt"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestFailingConstraint(t *testing.T) {
	validator := &Validator{
		OrderedPropertyChecks: true,
		Properties: Properties{
			"foo": {
				Order: 0,
				Type:  JsonAny,
				Constraints: Constraints{
					&FailingConstraint{},
					&StringNotEmpty{},
				},
			},
			"bar": {
				Order: 1,
				Type:  JsonString,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "",
		"bar": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))

	validator.Properties["foo"].Constraints[0].(*FailingConstraint).Stop = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))

	validator.Properties["foo"].Constraints[0].(*FailingConstraint).StopAll = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

func TestFailWhen(t *testing.T) {
	validator := &Validator{
		OrderedPropertyChecks: true,
		Properties: Properties{
			"foo": {
				Order: 0,
				Type:  JsonString,
				Constraints: Constraints{
					&SetConditionFrom{},
					&FailWhen{
						Conditions: []string{"fail"},
					},
				},
			},
			"bar": {
				Order: 1,
				Type:  JsonString,
			},
		},
	}
	obj := jsonObject(`{
		"foo": "fail",
		"bar": 1
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))

	validator.Properties["foo"].Constraints[1].(*FailWhen).StopAll = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))

	obj = jsonObject(`{
		"foo": "ok"
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestFailWith(t *testing.T) {
	constraint := &FailWith{
		Others: MustParseExpression("(foo && bar) || (bar && baz)"),
	}
	validator := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": {
				Constraints: Constraints{
					constraint,
					&StringNotEmpty{},
				},
			},
		},
	}

	obj := jsonObject(`{
		"foo": "aaa",
		"bar": "aaa"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgFailure, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)

	obj = jsonObject(`{
		"foo": "aaa",
		"baz": "aaa"
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj = jsonObject(`{
		"foo": "",
		"bar": "aaa"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))

	constraint.StopAll = true
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

func TestSetConditionFromConstraint(t *testing.T) {
	conditionWasSet := false
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&SetConditionFrom{},
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						conditionWasSet = vcx.IsCondition("TEST_CONDITION_TOKEN")
						return true, ""
					}, ""),
				},
			},
		},
	}
	obj := map[string]interface{}{
		"foo": "TEST_CONDITION_TOKEN",
	}

	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.True(t, conditionWasSet)

	obj["foo"] = "bar"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.False(t, conditionWasSet)
}

func TestSetConditionFromConstraintWithMapping(t *testing.T) {
	conditionWasSet := false
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&SetConditionFrom{
						Mapping: map[string]string{
							"TEST_CONDITION_TOKEN": "TEST_CONDITION_TOKEN2",
						},
						Prefix: "xxx-",
					},
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						conditionWasSet = vcx.IsCondition("xxx-TEST_CONDITION_TOKEN2")
						return true, ""
					}, ""),
				},
			},
		},
	}
	obj := map[string]interface{}{
		"foo": "TEST_CONDITION_TOKEN",
	}

	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.True(t, conditionWasSet)

	obj["foo"] = "bar"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.False(t, conditionWasSet)
}

func TestSetConditionFromConstraintOnParent(t *testing.T) {
	conditionWasSet := false
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&SetConditionFrom{Parent: true},
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						conditionWasSet = vcx.pathStack[len(vcx.pathStack)-2].conditions["TEST_CONDITION_TOKEN"]
						return true, ""
					}, ""),
				},
			},
		},
	}
	obj := map[string]interface{}{
		"foo": "TEST_CONDITION_TOKEN",
	}

	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.True(t, conditionWasSet)

	obj["foo"] = "bar"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.False(t, conditionWasSet)
}

func TestSetConditionFromConstraintGlobal(t *testing.T) {
	conditionWasSet := false
	validator := &Validator{
		OrderedPropertyChecks: true,
		Properties: Properties{
			"foo": {
				Order:     1,
				Type:      JsonObject,
				NotNull:   true,
				Mandatory: true,
				ObjectValidator: &Validator{
					Properties: Properties{
						"subFoo": {
							Type:      JsonObject,
							NotNull:   true,
							Mandatory: true,
							ObjectValidator: &Validator{
								Properties: Properties{
									"subSubFoo": {
										Type:      JsonString,
										NotNull:   true,
										Mandatory: true,
										Constraints: Constraints{
											&SetConditionFrom{Global: true},
										},
									},
								},
							},
						},
					},
				},
			},
			"bar": {
				Order:     2,
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						conditionWasSet = vcx.pathStack[len(vcx.pathStack)-2].conditions["TEST_CONDITION_TOKEN"]
						return true, ""
					}, ""),
				},
			},
		},
	}

	obj := jsonObject(`{
		"foo": {
			"subFoo": {
				"subSubFoo": "TEST_CONDITION_TOKEN"
			}
		},
		"bar": "whatever"
	}`)

	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.True(t, conditionWasSet)

	obj = jsonObject(`{
		"foo": {
			"subFoo": {
				"subSubFoo": "not set"
			}
		},
		"bar": "whatever"
	}`)
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.False(t, conditionWasSet)
}

func TestSetConditionFromNullsAndNonString(t *testing.T) {
	var grabVcx *ValidatorContext
	constraint := &SetConditionFrom{
		Global: true,
	}
	validator := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				grabVcx = vcx
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					constraint,
				},
			},
		},
	}

	obj := jsonObject(`{"foo": null}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("null"))

	constraint.NullToken = "foo_is_null"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("foo_is_null"))

	obj["foo"] = 1
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("1"))

	constraint.Format = "is_%v"
	constraint.Prefix = "foo_"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("foo_is_1"))
}

func TestSetConditionPropertyConstraint(t *testing.T) {
	validator := &Validator{
		OrderedPropertyChecks: true,
		Constraints: Constraints{
			&SetConditionProperty{
				PropertyName: "type",
			},
		},
		Properties: Properties{
			"type": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&StringValidToken{
						Tokens: []string{"foo", "bar", "baz"},
					},
				},
			},
			"foo": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"foo"},
			},
			"bar": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"bar"},
			},
			"baz": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"baz"},
			},
		},
	}

	// "bar":true should not be checked (even though it's incorrect type)...
	obj := jsonObject(`{
		"type": "foo",
		"foo": "checked",
		"bar": true
	}`)

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

func TestSetConditionPropertyConstraintWithMapping(t *testing.T) {
	validator := &Validator{
		OrderedPropertyChecks: true,
		Constraints: Constraints{
			&SetConditionProperty{
				PropertyName: "type",
				Prefix:       "xxx-",
				Mapping: map[string]string{
					"foo": "foo2",
					"bar": "bar2",
					"baz": "baz2",
				},
			},
		},
		Properties: Properties{
			"type": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&StringValidToken{
						Tokens: []string{"foo", "bar", "baz"},
					},
				},
			},
			"foo": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"xxx-foo2"},
			},
			"bar": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"xxx-bar2"},
			},
			"baz": {
				Type:           JsonString,
				NotNull:        true,
				Mandatory:      true,
				WhenConditions: []string{"xxx-baz2"},
			},
		},
	}

	// "bar":true should not be checked (even though it's incorrect type)...
	obj := jsonObject(`{
		"type": "foo",
		"foo": "checked",
		"bar": true
	}`)

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

func TestSetConditionProperty(t *testing.T) {
	var grabVcx *ValidatorContext
	constraint := &SetConditionProperty{
		PropertyName: "foo",
	}
	validator := &Validator{
		Constraints: Constraints{
			constraint,
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				grabVcx = vcx
				return true, ""
			}, ""),
		},
		Properties: Properties{
			"foo": {
				Type: JsonAny,
			},
		},
	}

	obj := jsonObject(`{}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("missing"))

	constraint.MissingToken = "not_there"
	grabVcx = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("not_there"))

	obj["foo"] = nil
	grabVcx = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("null"))

	obj["foo"] = 1
	grabVcx = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("1"))

	constraint.Format = "is_%v"
	grabVcx = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("is_1"))

	constraint.Prefix = "foo_"
	grabVcx = nil
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
	require.NotNil(t, grabVcx)
	require.True(t, grabVcx.IsCondition("foo_is_1"))
}

func TestVariablePropertyNameValidation(t *testing.T) {
	obj := jsonObject(`{
		"FOO": {
			"amount": 123
		},
		"BAR": {
			"amount": 345
		}
	}`)
	ncv := &VariablePropertyConstraint{
		NameConstraints: Constraints{
			&StringValidToken{Tokens: []string{"FOO", "BAR", "BAZ", "QUX"}, Message: "BAD PROPERTY NAME"},
		},
		ObjectValidator: &Validator{
			Properties: Properties{
				"amount": {
					Type:      JsonInteger,
					Mandatory: true,
					NotNull:   true,
				},
			},
		},
	}
	v := &Validator{
		IgnoreUnknownProperties: true,
		Constraints: Constraints{
			ncv,
		},
	}

	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj = jsonObject(`{
		"FOO": {
			"amount": null
		},
		"BAR": {
		},
		"BAZ": null,
		"QUX": "this should be an object",
		"XXX": {}
	}`)
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 5, len(violations))
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, msgValueCannotBeNull, violations[0].Message)
	require.Equal(t, "BAZ", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, CodeValueCannotBeNull, violations[0].Codes[0])
	require.Equal(t, msgPropertyValueMustBeObject, violations[1].Message)
	require.Equal(t, "QUX", violations[1].Property)
	require.Equal(t, "", violations[1].Path)
	require.Equal(t, CodePropertyValueMustBeObject, violations[1].Codes[0])
	require.Equal(t, "BAD PROPERTY NAME", violations[2].Message)
	require.Equal(t, "XXX", violations[2].Property)
	require.Equal(t, "", violations[2].Path)
	require.Equal(t, CodeInvalidPropertyName, violations[2].Codes[0])
	require.Equal(t, msgMissingProperty, violations[3].Message)
	require.Equal(t, "amount", violations[3].Property)
	require.Equal(t, "BAR", violations[3].Path)
	require.Equal(t, CodeMissingProperty, violations[3].Codes[0])
	require.Equal(t, msgValueCannotBeNull, violations[4].Message)
	require.Equal(t, "amount", violations[4].Property)
	require.Equal(t, "FOO", violations[4].Path)
	require.Equal(t, CodeValueCannotBeNull, violations[4].Codes[0])

	// with defaulted message...
	// (by placing NoMessage as message for any name validator - the default message is used)
	ncv.NameConstraints[0] = &StringValidToken{Tokens: []string{"FOO", "BAR", "BAZ", "QUX"}, Message: NoMessage}
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 5, len(violations))
	SortViolationsByPathAndProperty(violations)
	require.Equal(t, msgInvalidPropertyName, violations[2].Message)
	require.Equal(t, "XXX", violations[2].Property)
	require.Equal(t, "", violations[2].Path)
	require.Equal(t, CodeInvalidPropertyName, violations[2].Codes[0])

	// and again allowing nulls...
	ncv.AllowNull = true
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))
}

func TestVariablePropertyNameValidationStopOnFirst(t *testing.T) {
	obj := jsonObject(`{
		"FOO": {
			"amount": 123
		},
		"BAR": {
			"amount": 345
		}
	}`)
	ncv := &VariablePropertyConstraint{
		NameConstraints: Constraints{
			&StringValidToken{Tokens: []string{"FOO", "BAR", "BAZ", "QUX"}, Message: "BAD PROPERTY NAME"},
		},
		ObjectValidator: &Validator{
			Properties: Properties{
				"amount": {
					Type:      JsonInteger,
					Mandatory: true,
					NotNull:   true,
				},
			},
		},
	}
	v := &Validator{
		IgnoreUnknownProperties: true,
		StopOnFirst:             true,
		Constraints: Constraints{
			ncv,
		},
	}

	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj = jsonObject(`{
			"FOO": {
				"amount": null
			},
			"BAR": {
			},
			"BAZ": null,
			"QUX": "this should be an object",
			"XXX": {}
		}`)
	ok, _ = v.Validate(obj)
	require.False(t, ok)
}

func TestSetConditionOnType(t *testing.T) {
	foundTypes := map[string]bool{}
	foundTypesTestConstraint := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
		for k := range foundTypes {
			delete(foundTypes, k)
		}
		for _, token := range conditionTypeTokens {
			if vcx.IsCondition("type_" + token) {
				foundTypes["type_"+token] = true
			}
		}
		return true, ""
	}, "")
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonAny,
				Constraints: Constraints{
					&SetConditionOnType{},
					foundTypesTestConstraint,
				},
			},
		},
	}

	testCases := []struct {
		fooValue         interface{}
		expectFoundTypes []string
	}{
		{
			nil,
			[]string{"type_null"},
		},
		{
			"a string",
			[]string{"type_string"},
		},
		{
			true,
			[]string{"type_boolean"},
		},
		{
			map[string]interface{}{},
			[]string{"type_object"},
		},
		{
			[]interface{}{},
			[]string{"type_array"},
		},
		{
			3,
			[]string{"type_number", "type_integer"},
		},
		{
			3.0,
			[]string{"type_number"},
		},
		{
			json.Number("3"),
			[]string{"type_number", "type_integer"},
		},
		{
			json.Number("3.0"),
			[]string{"type_number"},
		},
		{
			json.Number("NaN"),
			[]string{"type_number", "type_nan"},
		},
		{
			json.Number("Inf"),
			[]string{"type_number", "type_inf"},
		},
		{
			json.Number("xxx"),
			[]string{"type_number", "type_invalid_number"},
		},
		{
			struct{}{},
			[]string{"type_unknown"},
		},
	}
	obj := map[string]interface{}{}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%v", i+1, tc.fooValue), func(t *testing.T) {
			obj["foo"] = tc.fooValue
			ok, _ := validator.Validate(obj)
			require.True(t, ok)
			okTypes := map[string]bool{}
			for _, ft := range tc.expectFoundTypes {
				require.True(t, foundTypes[ft])
				okTypes[ft[5:]] = true
			}
			for _, other := range conditionTypeTokens {
				if !okTypes[other] {
					require.False(t, foundTypes["type_"+other])
				}
			}
		})
	}
}

func TestSetConditionOnTypeWithFailingWhen(t *testing.T) {
	validator := &Validator{
		UseNumber: true,
		Properties: Properties{
			"foo": {
				Type:      JsonNumber,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&SetConditionOnType{},
					&FailWhen{
						Conditions: Conditions{"type_invalid_number", "type_nan", "type_inf"},
						Message:    "Must not NaN or Inf",
					},
				},
			},
		},
	}
	obj := jsonObject(`{"foo": 1}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = json.Number("NaN")
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Must not NaN or Inf", violations[0].Message)

	obj["foo"] = json.Number("Inf")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Must not NaN or Inf", violations[0].Message)

	obj["foo"] = json.Number("+Inf")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Must not NaN or Inf", violations[0].Message)

	obj["foo"] = json.Number("xxx")
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgValueExpectedType, "number"), violations[0].Message)
}

func TestArrayConditionalConstraint(t *testing.T) {
	// also tests array pathing!!!
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&Length{Minimum: 1},
				},
				ObjectValidator: &Validator{
					AllowArray:     true,
					DisallowObject: true,
					Properties: Properties{
						"bar": {
							Type:      JsonString,
							Mandatory: true,
							Constraints: Constraints{
								&ArrayConditionalConstraint{
									When: "!first",
									Constraint: &StringGreaterThanOther{
										PropertyName: "..[-1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "!last",
									Constraint: &StringLessThanOther{
										PropertyName: "..[+1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "first",
									Constraint: &StringLessThanOther{
										PropertyName: "..[+1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "last",
									Constraint: &StringGreaterThanOther{
										PropertyName: "..[-1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: ">0",
									Constraint: &StringGreaterThanOther{
										PropertyName: "..[-1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "<1",
									Constraint: &StringLessThanOther{
										PropertyName: "..[+1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "0",
									Constraint: &StringLessThanOther{
										PropertyName: "..[+1].bar",
									},
								},
								&ArrayConditionalConstraint{
									When: "%2",
									Constraint: &StringLessThanOther{
										PropertyName: "..[+1].bar",
									},
								},
							},
						},
						"qux": {
							Type:      JsonArray,
							Mandatory: true,
							ObjectValidator: &Validator{
								AllowArray:     true,
								DisallowObject: true,
								Properties: Properties{
									"foo": {
										Type: JsonString,
										Constraints: Constraints{
											&ArrayConditionalConstraint{
												When: "!first",
												Constraint: &StringGreaterThanOther{
													PropertyName: "..[-1].foo",
												},
											},
											&StringGreaterThanOther{
												PropertyName: "....[0].bar",
											},
											&StringGreaterThanOrEqualOther{
												PropertyName: "....[0].qux[0].foo",
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
	obj := jsonObject(`{
		"foo": [
			{
				"bar": "A",
				"qux": [
					{
						"foo": "AA"
					},
					{
						"foo": "AB"
					}
				]
			},
			{
				"bar": "B",
				"qux": [
					{
						"foo": "BA"
					}
				]
			},
			{
				"bar": "C",
				"qux": [
					{
						"foo": "CA"
					}
				]
			},
			{
				"bar": "D",
				"qux": [
					{
						"foo": "DA"
					}
				]
			}
		]
	}`)

	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
}

func TestConditionalConstraint_When(t *testing.T) {
	v := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&ConditionalConstraint{
						Constraint: &StringNotBlank{},
						When:       Conditions{"METHOD_POST"},
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	ok, violations = v.Validate(obj, "METHOD_POST")
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotBlankString, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
}

func TestConditionalConstraint_Others(t *testing.T) {
	v := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&ConditionalConstraint{
						Constraint: &StringNotBlank{},
						Others:     MustParseExpression("bar && baz"),
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["bar"] = ""
	ok, violations = v.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["baz"] = ""
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotBlankString, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
}

func TestConditionalConstraint_FailNotMet(t *testing.T) {
	v := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&ConditionalConstraint{
						Constraint:    &StringNotBlank{},
						When:          Conditions{"METHOD_POST"},
						FailNotMet:    true,
						NotMetMessage: "Conditions not met",
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Conditions not met", violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)

	ok, violations = v.Validate(obj, "METHOD_POST")
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotBlankString, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
}

func TestSetConditionIf(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
				Constraints: Constraints{
					&SetConditionIf{
						Constraint: &StringUppercase{},
						SetOk:      "IS_UPPER",
						SetFail:    "IS_NOT_UPPER",
					},
					&ConditionalConstraint{
						When:       Conditions{"IS_UPPER"},
						Constraint: &FailingConstraint{Message: "IS UPPER FAIL"},
					},
					&ConditionalConstraint{
						When:       Conditions{"IS_NOT_UPPER"},
						Constraint: &FailingConstraint{Message: "IS NOT UPPER FAIL"},
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": "A"
	}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "IS UPPER FAIL", violations[0].Message)

	obj["foo"] = "a"
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "IS NOT UPPER FAIL", violations[0].Message)
}

func TestSetConditionIf_NoLeaks(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&SetConditionIf{
						Constraint: NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
							vcx.AddViolation(NewViolation("", "", "SHOULD NOT BE SEEN"))
							vcx.Stop()
							vcx.CeaseFurther()
							return false, ""
						}, ""),
						SetOk:   "YES",
						SetFail: "NO",
						Global:  true,
					},
					&FailingConstraint{Message: "HERE"},
				},
			},
			"bar": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&FailingConstraint{Message: "HERE"},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": "",
		"bar": ""
	}`)

	ok, violations := v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, "HERE", violations[0].Message)
	require.Equal(t, "HERE", violations[1].Message)

	v = &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&SetConditionIf{
						Constraint: NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
							vcx.AddViolation(NewViolation("", "", "SHOULD NOT BE SEEN"))
							vcx.Stop()
							vcx.CeaseFurther()
							return false, ""
						}, ""),
						SetOk:   "YES",
						SetFail: "NO",
						Parent:  true,
					},
					&FailingConstraint{Message: "HERE"},
				},
			},
			"bar": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&FailingConstraint{Message: "HERE"},
				},
			},
		},
	}
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, "HERE", violations[0].Message)
	require.Equal(t, "HERE", violations[1].Message)
}

func TestConstraintSet(t *testing.T) {
	const msg = "String value length must be between 16 and 64 chars; must be letters (upper or lower), digits or underscores; must start with an uppercase letter"
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if len(str) == 0 || str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage(vcx)
					}
				}
				return true, ""
			}, ""),
			&StringCharacters{
				AllowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: msg,
	}
	require.Equal(t, msg, set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	// some more not oks...
	obj["foo"] = "abcdefghijklmnopqrstuvwxyz" // not starts with capital
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "Abc." // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
	obj["foo"] = "        " // empty after trim
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)
}

func TestConstraintSetNoMsg(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if str[0] < 'A' || str[0] > 'Z' {
						return false, this.GetMessage(vcx)
					}
				}
				return true, ""
			}, "Must start with a capital"),
			&StringCharacters{
				AllowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: "",
	}
	// message should return first sub-constraint with non-empty message...
	require.Equal(t, msgNotEmptyString, set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "  this does not start with capital letter and has spaces (and punctuation)  "
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)

	// some more not oks...
	obj["foo"] = "abcdefghijklmnopqrstuvwxyz" // not starts with capital
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Must start with a capital", violations[0].Message)
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmptyString, violations[0].Message)
}

func TestConstraintSetCeases(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{
			&StringNotEmpty{},
			&StringLength{Minimum: 16, Maximum: 64},
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				if str, ok := value.(string); ok {
					if str[0] < 'A' || str[0] > 'Z' {
						vcx.CeaseFurther()
					}
				}
				return true, ""
			}, "Must start with a capital"),
			&StringCharacters{
				AllowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{'0', 'z', 1}}},
				},
				DisallowRanges: []unicode.RangeTable{
					{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
					{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
					{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
				},
			},
		},
		Message: "",
	}
	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "this does not start with capital"
	}`)

	// the custom constraint ceased further and passed - so validation passes...
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// some more not oks...
	obj["foo"] = "AbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyzAbcdefghijklmnopqrstuvwxyz" // too long
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc" // too short
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgStringMinMaxLen, 16, tokenInclusive, 64, tokenInclusive), violations[0].Message)
	obj["foo"] = "Abc.01234567890123456" // contains invalid char
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgInvalidCharacters, violations[0].Message)
	obj["foo"] = "" // empty string
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNotEmptyString, violations[0].Message)
}

func TestConstraintSetOneOf(t *testing.T) {
	constraint1 := &testConstraint{passes: false, stops: false}
	constraint2 := &testConstraint{passes: true, stops: false}
	set := &ConstraintSet{
		OneOf:       true,
		Constraints: Constraints{constraint1, constraint2},
	}
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 2), set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "anything"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	// the second passing doesn't get reached if first stops...
	constraint1.stops = true
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 2), violations[0].Message)

	constraint1.stops = false
	constraint1.msg = "first message"
	constraint2.msg = "second message"
	constraint2.passes = false
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "first message", violations[0].Message)
}

func TestConstraintSetDefaultMessage(t *testing.T) {
	set := &ConstraintSet{
		Constraints: Constraints{&testConstraint{}},
	}
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultAllOf, 1), set.GetMessage(nil))

	validator := buildFooValidator(JsonString, set, false)
	obj := jsonObject(`{
		"foo": "anything"
	}`)
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultAllOf, 1), violations[0].Message)

	set.OneOf = true
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 1), set.GetMessage(nil))
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgConstraintSetDefaultOneOf, 1), violations[0].Message)
}

func TestConstraintSetWithConditionals(t *testing.T) {
	cc1 := &ConditionalConstraint{
		When: []string{"AA"},
		Constraint: &StringNotEmpty{
			Message: "MSG_1",
		},
	}
	cc2 := &ConditionalConstraint{
		When: []string{"BB"},
		Constraint: &StringNotEmpty{
			Message: "MSG_2",
		},
	}
	set := &ConstraintSet{
		Constraints: Constraints{cc1, cc2},
	}
	vcx := newValidatorContext(nil, nil, false, nil)
	ok, _ := set.Check("", vcx)
	require.True(t, ok)

	vcx.SetCondition("BB")
	ok, msg := set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "MSG_2", msg)

	vcx.SetCondition("AA")
	ok, msg = set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "MSG_1", msg)

	vcx.ClearCondition("AA")
	vcx.ClearCondition("BB")
	set.OneOf = true
	ok, msg = set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "Constraint set must pass one of 2 undisclosed validations", msg)
	vcx.SetCondition("AA")
	ok, msg = set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "MSG_1", msg)
	vcx.ClearCondition("AA")
	vcx.SetCondition("BB")
	ok, msg = set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "MSG_2", msg)
	vcx.SetCondition("AA")
	ok, msg = set.Check("", vcx)
	require.False(t, ok)
	require.Equal(t, "MSG_1", msg)
}

func TestIsNull(t *testing.T) {
	validator := buildFooValidator(JsonString, &IsNull{}, false)
	obj := jsonObject(`{
		"foo": null
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "not null"
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgNull, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
}

func TestIsNotNull(t *testing.T) {
	validator := buildFooValidator(JsonString, &IsNotNull{}, false)
	obj := jsonObject(`{
		"foo": "not null"
	}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = nil
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValueCannotBeNull, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
}

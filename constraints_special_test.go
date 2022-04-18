package valix

import (
	"testing"

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
	// (by placing " " as message for any name validator - the default message is used)
	ncv.NameConstraints[0] = &StringValidToken{Tokens: []string{"FOO", "BAR", "BAZ", "QUX"}, Message: " "}
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

	// and again with stop on first...
	//	ncv.AllowNull = false
	//	v.StopOnFirst = true
	//	ok, violations = v.Validate(obj)
	//	require.False(t, ok)
	//	require.Equal(t, 1, len(violations))
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
	ok, violations = v.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

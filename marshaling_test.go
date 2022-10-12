package valix

import (
	"encoding/json"
	"regexp"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
)

func TestValidator_MarshalJSON(t *testing.T) {
	v := &Validator{
		IgnoreUnknownProperties: true,
		AllowArray:              true,
		DisallowObject:          true,
		StopOnFirst:             true,
		OrderedPropertyChecks:   true,
		AllowNullItems:          true,
		WhenConditions:          []string{"test"},
		Constraints: Constraints{
			&Length{Minimum: 1, Maximum: 2, Message: "message 1"},
			&ConstraintSet{
				Constraints: Constraints{
					&Range{Minimum: 3, Maximum: 4, Message: "message 2"},
				},
				Message: "message 3",
			},
		},
		ConditionalVariants: ConditionalVariants{
			&ConditionalVariant{
				Constraints: Constraints{
					&SetConditionProperty{PropertyName: "foo"},
				},
				Properties: Properties{
					"foo": {
						Type:      JsonString,
						Mandatory: true,
					},
				},
				WhenConditions: []string{"cv_foo"},
			},
			&ConditionalVariant{
				Constraints: Constraints{
					&SetConditionProperty{PropertyName: "foo"},
				},
				Properties: Properties{
					"foo": {
						Type:      JsonString,
						Mandatory: true,
					},
				},
				WhenConditions: []string{"!cv_foo"},
			},
		},
		OasInfo: &OasInfo{Deprecated: true},
		Properties: Properties{
			"foo": {
				Type:          JsonString,
				Mandatory:     true,
				MandatoryWhen: []string{"want_foo"},
				NotNull:       true,
				Constraints: Constraints{
					&StringNotEmpty{Message: "message 4"},
					&StringPattern{Message: "message 5", Regexp: *regexp.MustCompile("^([A-Z]+)$")},
				},
				Order:               1,
				WhenConditions:      []string{"when_foo"},
				UnwantedConditions:  []string{"!when_foo"},
				RequiredWith:        MustParseExpression("(foo || bar) && !(foo && bar)"),
				RequiredWithMessage: "only sometimes",
				UnwantedWith:        MustParseExpression("(foo || bar) && !(foo && bar)"),
				ObjectValidator:     nil,
				OasInfo:             &OasInfo{Deprecated: true},
			},
			"bar": {
				Type:      JsonObject,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&ConditionalConstraint{
						When:       Conditions{"check_bar"},
						Others:     MustParseExpression("foo && bar"),
						Constraint: &Length{Minimum: 1},
					},
				},
				ObjectValidator: &Validator{
					AllowArray:     true,
					DisallowObject: true,
				},
			},
			"baz": {
				Type:      JsonObject,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&VariablePropertyConstraint{
						NameConstraints: Constraints{
							&StringValidToken{Tokens: []string{"FOO2", "BAR2", "BAZ2"}},
						},
						ObjectValidator: &Validator{
							Properties: Properties{
								"foo3": {
									Type: JsonString,
								},
							},
						},
						AllowNull: true,
					},
				},
			},
			"qux": {
				Constraints: Constraints{
					&SetConditionIf{
						Constraint: &StringUppercase{},
						SetOk:      "IS_UPPER",
						SetFail:    "IS_NOT_UPPER",
					},
				},
			},
			"arr": {
				Type: JsonString,
				Constraints: Constraints{
					&ArrayConditionalConstraint{
						When: "%2",
						Constraint: &StringGreaterThanOther{
							PropertyName: "foo",
						},
					},
				},
			},
		},
	}

	b, err := json.Marshal(v)
	require.Nil(t, err)

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)
	require.Nil(t, err)
	require.NotNil(t, obj)
	//pretty, _ := json.MarshalIndent(obj, "", "\t")
	//println(string(pretty[:]))

	valid, violations := ValidatorValidator.Validate(obj)
	require.Equal(t, 0, len(violations))
	require.True(t, valid)

	// and check the JSON matches the validator...
	require.Equal(t, 13, len(obj))
	require.True(t, obj[ptyNameIgnoreUnknownProperties].(bool))
	require.True(t, obj[ptyNameAllowArray].(bool))
	require.True(t, obj[ptyNameDisallowObject].(bool))
	require.True(t, obj[ptyNameStopOnFirst].(bool))
	require.True(t, obj[ptyNameOrderedPropertyChecks].(bool))
	require.True(t, obj[ptyNameAllowNullItems].(bool))
	require.Equal(t, 1, len(obj[ptyNameWhenConditions].([]interface{})))
	require.Equal(t, "test", obj[ptyNameWhenConditions].([]interface{})[0])

	constraints := obj[ptyNameConstraints].([]interface{})
	require.Equal(t, 2, len(constraints))

	constraint := constraints[0].(map[string]interface{})
	require.Equal(t, "Length", constraint["name"])
	fields := constraint[ptyNameFields].(map[string]interface{})
	require.Equal(t, 6, len(fields))
	require.Equal(t, "message 1", fields["Message"])
	require.Equal(t, float64(1), fields["Minimum"])
	require.Equal(t, float64(2), fields["Maximum"])

	constraint = constraints[1].(map[string]interface{})
	require.Equal(t, constraintSetName, constraint["name"])
	fields = constraint[ptyNameFields].(map[string]interface{})
	require.Equal(t, 4, len(fields))
	require.Equal(t, "message 3", fields["Message"])
	slc := fields["Constraints"].([]interface{})
	require.Equal(t, 1, len(slc))
	constraint = slc[0].(map[string]interface{})
	require.Equal(t, "Range", constraint["name"])
	fields = constraint[ptyNameFields].(map[string]interface{})
	require.Equal(t, 7, len(fields))
	require.Equal(t, "message 2", fields["Message"])
	require.Equal(t, float64(3), fields["Minimum"])
	require.Equal(t, float64(4), fields["Maximum"])

	variants := obj[ptyNameConditionalVariants].([]interface{})
	require.Equal(t, 2, len(variants))
	variant := variants[0].(map[string]interface{})
	require.Equal(t, 3, len(variant))
	whens := variant[ptyNameWhenConditions].([]interface{})
	require.Equal(t, 1, len(whens))
	require.Equal(t, "cv_foo", whens[0])
	sub := variant[ptyNameProperties].(map[string]interface{})
	require.Equal(t, 1, len(sub))
	slc = variant[ptyNameConstraints].([]interface{})
	require.Equal(t, 1, len(slc))

	sub = obj[ptyNameOasInfo].(map[string]interface{})
	require.Equal(t, 5, len(sub))
	require.True(t, sub[ptyNameOasDeprecated].(bool))

	sub = obj[ptyNameProperties].(map[string]interface{})
	require.Equal(t, 5, len(sub))

	pty := sub["foo"].(map[string]interface{})
	require.Equal(t, 13, len(pty))
	require.Equal(t, "string", pty[ptyNameType])
	require.Equal(t, true, pty[ptyNameMandatory])
	require.Equal(t, true, pty[ptyNameNotNull])
	require.Equal(t, "(foo || bar) && !(foo && bar)", pty[ptyNameRequiredWith])
	require.Equal(t, "(foo || bar) && !(foo && bar)", pty[ptyNameUnwantedWith])
	require.Equal(t, float64(1), pty[ptyNameOrder])
	_, present := pty[ptyNameObjectValidator]
	require.False(t, present)
	slc = pty[ptyNameMandatoryWhen].([]interface{})
	require.Equal(t, 1, len(slc))
	require.Equal(t, "want_foo", slc[0])
	slc = pty[ptyNameWhenConditions].([]interface{})
	require.Equal(t, 1, len(slc))
	require.Equal(t, "when_foo", slc[0])
	slc = pty[ptyNameUnwantedConditions].([]interface{})
	require.Equal(t, 1, len(slc))
	require.Equal(t, "!when_foo", slc[0])
	constraints = pty[ptyNameConstraints].([]interface{})
	require.Equal(t, 2, len(constraints))
	constraint = constraints[0].(map[string]interface{})
	require.Equal(t, "StringNotEmpty", constraint["name"])
	fields = constraint[ptyNameFields].(map[string]interface{})
	require.Equal(t, 3, len(fields))
	require.Equal(t, "message 4", fields["Message"])
	constraint = constraints[1].(map[string]interface{})
	require.Equal(t, "StringPattern", constraint["name"])
	fields = constraint[ptyNameFields].(map[string]interface{})
	require.Equal(t, 3, len(fields))
	require.Equal(t, "message 5", fields["Message"])
	require.Equal(t, "^([A-Z]+)$", fields["Regexp"])
	subSub := obj[ptyNameOasInfo].(map[string]interface{})
	require.Equal(t, 5, len(subSub))
	require.True(t, subSub[ptyNameOasDeprecated].(bool))

	pty = sub["bar"].(map[string]interface{})
	require.Equal(t, 8, len(pty))
	require.Equal(t, "object", pty[ptyNameType])
	require.Equal(t, true, pty[ptyNameMandatory])
	require.Equal(t, true, pty[ptyNameNotNull])
	require.Equal(t, float64(0), pty[ptyNameOrder])
	_, present = pty[ptyNameObjectValidator]
	require.True(t, present)
	subSub = pty[ptyNameObjectValidator].(map[string]interface{})
	require.Equal(t, true, subSub[ptyNameAllowArray])
	require.Equal(t, true, subSub[ptyNameDisallowObject])
	constraints = pty[ptyNameConstraints].([]interface{})
	require.Equal(t, 1, len(constraints))
	c0 := constraints[0].(map[string]interface{})
	require.Equal(t, "foo && bar", c0[ptyNameOthersExpr])

	pty = sub["qux"].(map[string]interface{})
	constraints = pty[ptyNameConstraints].([]interface{})
	require.Equal(t, 1, len(constraints))
	c0 = constraints[0].(map[string]interface{})
	require.Equal(t, "SetConditionIf", c0["name"])
	require.NotEmpty(t, c0["fields"])
	fields = c0["fields"].(map[string]interface{})
	require.Equal(t, "IS_UPPER", fields["SetOk"])
	require.Equal(t, "IS_NOT_UPPER", fields["SetFail"])
	require.NotEmpty(t, fields["Constraint"])
	wc := fields["Constraint"].(map[string]interface{})
	require.Equal(t, "StringUppercase", wc["name"])
}

func TestValidator_MarshalJSON_FailsWithCustomConstraints(t *testing.T) {
	v := &Validator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				return true, ""
			}, ""),
		},
	}
	_, err := json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: unsupported type: valix.Check", err.Error())

	v = &Validator{
		Constraints: Constraints{
			&ConstraintSet{
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						return true, ""
					}, ""),
				},
			},
		},
	}
	_, err = json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: error calling MarshalJSON for type *valix.ConstraintSet: json: unsupported type: valix.Check", err.Error())

	v = &Validator{
		Properties: Properties{
			"foo": {
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						return true, ""
					}, ""),
				},
			},
		},
	}
	_, err = json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: unsupported type: valix.Check", err.Error())

	v = &Validator{
		Properties: Properties{
			"foo": {
				ObjectValidator: &Validator{
					Constraints: Constraints{
						NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
							return true, ""
						}, ""),
					},
				},
			},
		},
	}
	_, err = json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: unsupported type: valix.Check", err.Error())

	v = &Validator{
		ConditionalVariants: ConditionalVariants{
			{
				WhenConditions: []string{"foo"},
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
						return true, ""
					}, ""),
				},
			},
		},
	}
	_, err = json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: unsupported type: valix.Check", err.Error())

	v = &Validator{
		ConditionalVariants: ConditionalVariants{
			{
				WhenConditions: []string{"foo"},
				Properties: Properties{
					"foo": {
						Constraints: Constraints{
							NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
								return true, ""
							}, ""),
						},
					},
				},
			},
		},
	}
	_, err = json.Marshal(v)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.Validator: json: unsupported type: valix.Check", err.Error())

	cs := Constraints{
		NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
			return true, ""
		}, ""),
	}
	_, err = json.Marshal(cs)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type valix.Constraints: json: unsupported type: valix.Check", err.Error())
}

func TestPropertyValidator_MarshalJSON_WithAllCommonConstraints(t *testing.T) {
	all := defaultConstraints()
	list := make([]Constraint, 0, len(all))
	for _, c := range all {
		list = append(list, c)
	}
	pv := &PropertyValidator{
		Constraints: list,
	}
	b, err := json.Marshal(pv)
	require.Nil(t, err)

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)
	require.Nil(t, err)
	require.NotNil(t, obj)
	require.Equal(t, 7, len(obj))

	constraints := obj[ptyNameConstraints].([]interface{})
	require.Equal(t, len(all), len(constraints))
}

func TestPropertyValidator_MarshalJSON_FailsWithCustomConstraints(t *testing.T) {
	pv := &PropertyValidator{
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
				return true, ""
			}, ""),
		},
	}
	_, err := json.Marshal(pv)
	require.NotNil(t, err)
	require.Equal(t, "json: error calling MarshalJSON for type *valix.PropertyValidator: json: unsupported type: valix.Check", err.Error())
}

func TestConstraint_StringCharacters_MarshalJSON(t *testing.T) {
	c := &StringCharacters{
		AllowRanges: []unicode.RangeTable{
			UnicodeBMP,
			UnicodeSMP,
			UnicodeSIP,
		},
		DisallowRanges: []unicode.RangeTable{
			*unicode.Upper,
			*unicode.Title,
		},
		Message: "test message",
	}
	b, err := json.Marshal(c)
	require.Nil(t, err)

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)
	require.Nil(t, err)
	require.NotNil(t, obj)
}

func TestConstraint_StringPattern_MarshalJSON(t *testing.T) {
	c := &StringPattern{
		Regexp:  *regexp.MustCompile("^([A-Z]+)$"),
		Message: "test message",
	}
	b, err := json.Marshal(c)
	require.Nil(t, err)

	obj := map[string]interface{}{}
	err = json.Unmarshal(b, &obj)
	require.Nil(t, err)
	require.NotNil(t, obj)
	require.Equal(t, "test message", obj["Message"])
	require.Equal(t, "^([A-Z]+)$", obj["Regexp"])
}

func TestConstraint_ArrayConditionalConstraint_MarshalJSONFails(t *testing.T) {
	c := &ArrayConditionalConstraint{
		Constraint: &CustomConstraint{
			CheckFunc: func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
				return true, ""
			},
		},
	}
	_, err := json.Marshal(c)
	require.NotNil(t, err)
}

func TestAllConstraintsCanBeMarshalledAndUnmarshalled(t *testing.T) {
	pv := &PropertyValidator{
		Constraints: Constraints{},
	}
	v := Validator{
		Properties: Properties{
			"foo": pv,
		},
	}
	for _, c := range constraintsRegistry.namedConstraints {
		pv.Constraints = append(pv.Constraints, c)
	}
	b, err := json.Marshal(v)
	require.Nil(t, err)

	//obj := map[string]interface{}{}
	//err = json.Unmarshal(b, &obj)
	//pretty, _ := json.MarshalIndent(obj, "", "\t")
	//println(string(pretty[:]))

	uv := &Validator{}
	err = json.Unmarshal(b, &uv)
	require.Nil(t, err)
	upv := uv.Properties["foo"]
	require.Equal(t, len(pv.Constraints), len(upv.Constraints))
}

func TestConstraint_SetConditionIf_MarshalJSONFails(t *testing.T) {
	c := &SetConditionIf{
		Constraint: &CustomConstraint{
			CheckFunc: func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
				return true, ""
			},
		},
	}
	_, err := json.Marshal(c)
	require.NotNil(t, err)
}

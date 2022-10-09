package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"testing"
	"unicode"
)

func TestPropertyValidator_ToV8nTagString_Empty(t *testing.T) {
	pv := &PropertyValidator{}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s`, tagTokenType, jsonTypeTokenAny), str)
}

func TestPropertyValidator_ToV8nTagString_WithOptions(t *testing.T) {
	pv := &PropertyValidator{Type: JsonString, Order: -1}
	str := pv.ToV8nTagString(&V8nTagStringOptions{UnSpaced: true})
	require.Equal(t, fmt.Sprintf(`%s:%s,%s:%d`, tagTokenType, jsonTypeTokenString, tagTokenOrder, pv.Order), str)

	str = pv.ToV8nTagString(&V8nTagStringOptions{UnSpaced: false})
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:%d`, tagTokenType, jsonTypeTokenString, tagTokenOrder, pv.Order), str)
}

func TestPropertyValidator_ToV8nTagString_Order(t *testing.T) {
	pv := &PropertyValidator{
		Order: -1,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:%d`, tagTokenType, jsonTypeTokenAny, tagTokenOrder, pv.Order), str)
}

func TestPropertyValidator_ToV8nTagString_Required(t *testing.T) {
	pv := &PropertyValidator{
		Mandatory: true,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s`, tagTokenType, jsonTypeTokenAny, tagTokenRequired), str)
}

func TestPropertyValidator_ToV8nTagString_RequiredWhen(t *testing.T) {
	pv := &PropertyValidator{
		MandatoryWhen: []string{"foo"},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo`, tagTokenType, jsonTypeTokenAny, tagTokenRequired), str)

	pv = &PropertyValidator{
		MandatoryWhen: []string{"foo", "bar"},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:[foo,bar]`, tagTokenType, jsonTypeTokenAny, tagTokenRequired), str)
}

func TestPropertyValidator_ToV8nTagString_NotNull(t *testing.T) {
	pv := &PropertyValidator{
		NotNull: true,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s`, tagTokenType, jsonTypeTokenAny, tagTokenNotNull), str)
}

func TestPropertyValidator_ToV8nTagString_Only(t *testing.T) {
	pv := &PropertyValidator{
		Only: true,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s`, tagTokenType, jsonTypeTokenAny, tagTokenOnly), str)
}

func TestPropertyValidator_ToV8nTagString_OnlyWithMessage(t *testing.T) {
	pv := &PropertyValidator{
		Only:        true,
		OnlyMessage: `This doesn't have double quotes`,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s, %s:%s`, tagTokenType, jsonTypeTokenAny, tagTokenOnly, tagTokenOnlyMsg, `"This doesn't have double quotes"`), str)
}

func TestPropertyValidator_ToV8nTagString_OnlyWithMessageDoubleQuotes(t *testing.T) {
	pv := &PropertyValidator{
		Only:        true,
		OnlyMessage: `This does not have single quotes but has "doubles"`,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s, %s:'%s'`, tagTokenType, jsonTypeTokenAny, tagTokenOnly, tagTokenOnlyMsg, `This does not have single quotes but has "doubles"`), str)
}

func TestPropertyValidator_ToV8nTagString_OnlySingleCondition(t *testing.T) {
	pv := &PropertyValidator{
		OnlyConditions: []string{"foo"},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo`, tagTokenType, jsonTypeTokenAny, tagTokenOnly), str)
}

func TestPropertyValidator_ToV8nTagString_OnlyMultiConditions(t *testing.T) {
	pv := &PropertyValidator{
		OnlyConditions: []string{"foo", "bar"},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:[foo,bar]`, tagTokenType, jsonTypeTokenAny, tagTokenOnly), str)
}

func TestPropertyValidator_ToV8nTagString_StopOnFirst(t *testing.T) {
	pv := &PropertyValidator{
		StopOnFirst: true,
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s`, tagTokenType, jsonTypeTokenAny, tagTokenStopOnFirstAlt), str)
}

func TestPropertyValidator_ToV8nTagString_WhenConditions(t *testing.T) {
	pv := &PropertyValidator{
		WhenConditions: []string{"foo"},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo`, tagTokenType, jsonTypeTokenAny, tagTokenWhen), str)

	pv = &PropertyValidator{
		WhenConditions: []string{"foo", "bar"},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:[foo,bar]`, tagTokenType, jsonTypeTokenAny, tagTokenWhen), str)
}

func TestPropertyValidator_ToV8nTagString_UnwantedConditions(t *testing.T) {
	pv := &PropertyValidator{
		UnwantedConditions: []string{"foo"},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo`, tagTokenType, jsonTypeTokenAny, tagTokenUnwanted), str)

	pv = &PropertyValidator{
		UnwantedConditions: []string{"foo", "bar"},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:[foo,bar]`, tagTokenType, jsonTypeTokenAny, tagTokenUnwanted), str)
}

func TestPropertyValidator_ToV8nTagString_RequiredWith(t *testing.T) {
	pv := &PropertyValidator{
		RequiredWith: MustParseExpression("foo && !bar"),
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo && !bar`, tagTokenType, jsonTypeTokenAny, tagTokenRequiredWithAlt), str)

	pv = &PropertyValidator{
		RequiredWith:        MustParseExpression("foo && !bar"),
		RequiredWithMessage: "fooey",
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo && !bar, %s:'fooey'`, tagTokenType, jsonTypeTokenAny, tagTokenRequiredWithAlt, tagTokenRequiredWithAltMsg), str)
}

func TestPropertyValidator_ToV8nTagString_UnwantedWith(t *testing.T) {
	pv := &PropertyValidator{
		UnwantedWith: MustParseExpression("foo && !bar"),
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo && !bar`, tagTokenType, jsonTypeTokenAny, tagTokenUnwantedWithAlt), str)

	pv = &PropertyValidator{
		UnwantedWith:        MustParseExpression("foo && !bar"),
		UnwantedWithMessage: "fooey",
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo && !bar, %s:'fooey'`, tagTokenType, jsonTypeTokenAny, tagTokenUnwantedWithAlt, tagTokenUnwantedWithAltMsg), str)
}

func TestPropertyValidator_ToV8nTagString_ObjectValidator(t *testing.T) {
	pv := &PropertyValidator{
		ObjectValidator: &Validator{},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s`, tagTokenType, jsonTypeTokenAny), str)

	pv.ObjectValidator.IgnoreUnknownProperties = true
	pv.ObjectValidator.OrderedPropertyChecks = true
	pv.ObjectValidator.AllowNullItems = true
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s, %s, %s`, tagTokenType, jsonTypeTokenAny, tagTokenObjIgnoreUnknownProperties, tagTokenObjOrdered, tagTokenArrAllowNullItems), str)

	pv.ObjectValidator = &Validator{
		WhenConditions: Conditions{"foo"},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:foo`, tagTokenType, jsonTypeTokenAny, tagTokenObjWhen), str)

	pv.ObjectValidator = &Validator{
		WhenConditions: Conditions{"foo", "bar"},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:[foo,bar]`, tagTokenType, jsonTypeTokenAny, tagTokenObjWhen), str)

	pv.ObjectValidator = &Validator{
		Constraints: Constraints{
			&NotEmpty{},
			&LengthExact{Value: 1},
		},
	}
	str = pv.ToV8nTagString(nil)
	require.Equal(t, fmt.Sprintf(`%s:%s, %s:&NotEmpty{}, %s:&LengthExact{1}`, tagTokenType, jsonTypeTokenAny, tagTokenObjConstraint, tagTokenObjConstraint), str)
}

func TestPropertyValidator_ToV8nTagString_Constraint(t *testing.T) {
	basicFormatPlus := func(cStr string) string {
		if cStr == "" {
			return fmt.Sprintf(`%s:%s`, tagTokenType, jsonTypeTokenAny)
		}
		return fmt.Sprintf(`%s:%s,%s`, tagTokenType, jsonTypeTokenAny, cStr)
	}
	testCases := []struct {
		constraints Constraints
		options     *V8nTagStringOptions
		expect      string
	}{
		{
			nil,
			nil,
			basicFormatPlus(""),
		},
		{
			Constraints{},
			nil,
			basicFormatPlus(""),
		},
		{
			Constraints{&StringNotEmpty{}},
			nil,
			basicFormatPlus(" &StringNotEmpty{}"),
		},
		{
			Constraints{&StringNotEmpty{}},
			&V8nTagStringOptions{DiscardUnneededCurlies: true},
			basicFormatPlus(" &StringNotEmpty"),
		},
		{
			Constraints{
				&StringNotEmpty{
					Message: "fooey",
				},
			},
			&V8nTagStringOptions{DiscardUnneededCurlies: true},
			basicFormatPlus(" &StringNotEmpty{'fooey'}"),
		},
		{
			Constraints{
				&StringNotEmpty{
					Message: "fooey",
					Stop:    true,
				},
			},
			nil,
			basicFormatPlus(" &StringNotEmpty{Message:'fooey', Stop}"),
		},
		{
			Constraints{
				&StringNotEmpty{},
			},
			&V8nTagStringOptions{
				AbbreviateConstraintNames: true,
				DiscardUnneededCurlies:    true,
			},
			basicFormatPlus(" &strne"),
		},
		{
			Constraints{
				&StringPresetPattern{},
			},
			nil,
			basicFormatPlus(" &StringPresetPattern{}"),
		},
		{
			Constraints{
				&StringPresetPattern{
					Preset: PresetPublication,
				},
			},
			nil,
			basicFormatPlus(fmt.Sprintf(" &%s{}", PresetPublication)),
		},
		{
			Constraints{
				&StringPresetPattern{
					Preset: PresetPublication,
				},
			},
			&V8nTagStringOptions{
				DiscardUnneededCurlies: true,
			},
			basicFormatPlus(fmt.Sprintf(" &%s", PresetPublication)),
		},
		{
			Constraints{
				&ConditionalConstraint{
					When: Conditions{"foo", "bar"},
					Constraint: &StringPresetPattern{
						Preset: PresetBarcode,
					},
				},
			},
			nil,
			basicFormatPlus(" &[foo,bar]barcode{}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					Others: MustParseExpression("foo && bar"),
					Constraint: &StringPresetPattern{
						Preset: PresetBarcode,
					},
				},
			},
			nil,
			basicFormatPlus(" &<foo && bar>barcode{}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					Others: MustParseExpression("foo && bar"),
					Constraint: &StringPresetPattern{
						Preset: PresetBarcode,
					},
				},
			},
			&V8nTagStringOptions{
				NoUnwrapConditionalConstraints: true,
			},
			basicFormatPlus(" &ConditionalConstraint{Others:'foo && bar', Constraint:&barcode{}}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					Constraint: &StringPresetPattern{
						Preset: PresetBarcode,
					},
				},
			},
			nil,
			basicFormatPlus(" &barcode{}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					When:   Conditions{"foo", "bar"},
					Others: MustParseExpression("foo && bar"),
					Constraint: &StringPresetPattern{
						Preset: PresetBarcode,
					},
				},
			},
			nil,
			basicFormatPlus(" &ConditionalConstraint{When:[foo,bar], Others:'foo && bar', Constraint:&barcode{}}"),
		},
		{
			Constraints{
				&Range{
					Minimum: -1.23,
					Maximum: 4.56,
				},
			},
			nil,
			basicFormatPlus(" &Range{Minimum:-1.23, Maximum:4.56}"),
		},
		{
			Constraints{
				&RangeInt{
					Minimum:      1,
					Maximum:      2,
					ExclusiveMin: true,
					ExclusiveMax: false,
				},
			},
			nil,
			basicFormatPlus(" &RangeInt{Minimum:1, Maximum:2, ExclusiveMin}"),
		},
		{
			Constraints{
				&ArrayOf{
					Type: "string",
					Constraints: Constraints{
						&StringNotEmpty{},
						&StringUppercase{},
					},
				},
			},
			nil,
			basicFormatPlus(" &ArrayOf{Type:'string', Constraints:[&StringNotEmpty{}, &StringUppercase{}]}"),
		},
		{
			Constraints{
				&StringCharacters{},
			},
			nil,
			basicFormatPlus(" &StringCharacters{}"),
		},
		{
			Constraints{
				&StringCharacters{
					AllowRanges: []unicode.RangeTable{
						UnicodeBMP,
					},
				},
			},
			nil,
			basicFormatPlus(" &StringCharacters{AllowRanges:[{\"R16\":[{\"Lo\":0,\"Hi\":65535,\"Stride\":1}],\"R32\":null,\"LatinOffset\":0}]}"),
		},
		{
			Constraints{
				&StringPattern{
					Regexp: *regexp.MustCompile("^([A-Za-z0-9]*)$"),
				},
			},
			nil,
			basicFormatPlus(" &StringPattern{'^([A-Za-z0-9]*)$'}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					When: Conditions{"foo"},
					Constraint: &ConditionalConstraint{
						When: Conditions{"bar"},
						Constraint: &ConditionalConstraint{
							Others:     MustParseExpression("baz && qux"),
							Constraint: &StringNotEmpty{},
						},
					},
				},
			},
			nil,
			basicFormatPlus(" &ConditionalConstraint{When:[foo], Constraint:&ConditionalConstraint{When:[bar], Constraint:&<baz && qux>StringNotEmpty{}}}"),
		},
		{
			Constraints{
				&ConditionalConstraint{
					When: Conditions{"foo"},
					Constraint: &ConditionalConstraint{
						When: Conditions{"bar"},
						Constraint: &ConditionalConstraint{
							Others:     MustParseExpression("baz && qux"),
							Constraint: &StringNotEmpty{},
						},
					},
				},
			},
			&V8nTagStringOptions{NoUnwrapConditionalConstraints: true},
			basicFormatPlus(" &ConditionalConstraint{When:[foo], Constraint:&ConditionalConstraint{When:[bar], Constraint:&ConditionalConstraint{Others:'baz && qux', Constraint:&StringNotEmpty{}}}}"),
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s", i+1, tc.expect), func(t *testing.T) {
			pv := &PropertyValidator{
				Constraints: tc.constraints,
			}
			str := pv.ToV8nTagString(tc.options)
			require.Equal(t, tc.expect, str)

			_, err := NewPropertyValidator(str)
			require.Nil(t, err)
		})
	}
}

type ConstraintWithDefaultBoolTrue struct {
	Test bool `v8n:"default"`
}

func (c *ConstraintWithDefaultBoolTrue) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (c *ConstraintWithDefaultBoolTrue) GetMessage(tcx I18nContext) string {
	return ""
}

type ConstraintWithDefaultBoolFalse struct {
	Test bool `v8n:"default"`
}

func (c *ConstraintWithDefaultBoolFalse) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (c *ConstraintWithDefaultBoolFalse) GetMessage(tcx I18nContext) string {
	return ""
}

func TestPropertyValidator_ToV8nTagString_FieldAbbreviations(t *testing.T) {
	constraintsRegistry.register(true, &ConstraintWithDefaultBoolTrue{Test: true})
	constraintsRegistry.register(true, &ConstraintWithDefaultBoolFalse{Test: false})
	defer constraintsRegistry.reset()
	basicFormatPlus := func(cStr string) string {
		return fmt.Sprintf(`%s:%s,%s`, tagTokenType, jsonTypeTokenAny, cStr)
	}
	testCases := []struct {
		constraint Constraint
		min        uint
		expect     string
	}{
		{
			&Range{
				Minimum: 1,
				Maximum: 1,
			},
			10,
			basicFormatPlus(" &Range{Minimum:1, Maximum:1}"),
		},
		{
			&Range{
				Minimum: 1,
				Maximum: 1,
			},
			3,
			basicFormatPlus(" &Range{Min:1, Max:1}"),
		},
		{
			&Range{
				Minimum: 1,
				Maximum: 1,
			},
			2,
			basicFormatPlus(" &Range{Mi:1, Ma:1}"),
		},
		{
			&Range{
				Minimum: 1,
				Maximum: 1,
			},
			1,
			basicFormatPlus(" &Range{Mi:1, Ma:1}"),
		},
		{
			&Range{
				Minimum: 1,
				Maximum: 1,
			},
			0,
			basicFormatPlus(" &Range{Min:1, Max:1}"),
		},
		{
			&Range{
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			0,
			basicFormatPlus(" &Range{ExcMin, ExcMax}"),
		},
		{
			&Range{
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			1,
			basicFormatPlus(" &Range{ExMi, ExMa}"),
		},
		{
			&Range{
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			2,
			basicFormatPlus(" &Range{ExMi, ExMa}"),
		},
		{
			&Range{
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			3,
			basicFormatPlus(" &Range{ExcMin, ExcMax}"),
		},
		{
			&Range{
				ExclusiveMin: true,
				ExclusiveMax: true,
			},
			20,
			basicFormatPlus(" &Range{ExclusiveMin, ExclusiveMax}"),
		},
		{
			&NetIsIP{
				V4Only:           true,
				V6Only:           true,
				Resolvable:       true,
				DisallowLoopback: true,
				DisallowPrivate:  true,
				AllowLocalhost:   true,
			},
			0,
			basicFormatPlus(" &NetIsIP{V4Only, V6Only, Res, DisLoo, DisPri, AllLoc}"),
		},
		{
			&ConstraintWithDefaultBoolTrue{},
			0,
			basicFormatPlus(" &ConstraintWithDefaultBoolTrue{false}"),
		},
		{
			&ConstraintWithDefaultBoolTrue{Test: true},
			0,
			basicFormatPlus(" &ConstraintWithDefaultBoolTrue{}"),
		},
		{
			&ConstraintWithDefaultBoolFalse{},
			0,
			basicFormatPlus(" &ConstraintWithDefaultBoolFalse{}"),
		},
		{
			&ConstraintWithDefaultBoolFalse{Test: true},
			0,
			basicFormatPlus(" &ConstraintWithDefaultBoolFalse{true}"),
		},
	}
	options := &V8nTagStringOptions{
		AbbreviateFieldNames: true,
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]%s", i+1, tc.expect), func(t *testing.T) {
			options.MinimumFieldNameLength = tc.min
			pv := &PropertyValidator{
				Constraints: Constraints{tc.constraint},
			}
			str := pv.ToV8nTagString(options)
			_, err := NewPropertyValidator(str)
			require.Nil(t, err)

			require.Equal(t, tc.expect, str)
		})
	}
}

func TestPropertyValidator_ToV8nTagString_WithAllConstraints(t *testing.T) {
	for k, c := range defaultConstraints() {
		t.Run(k, func(t *testing.T) {
			pv := &PropertyValidator{
				Constraints: Constraints{c},
			}
			str := pv.ToV8nTagString(nil)
			name := reflect.TypeOf(c).Elem().Name()
			if ppc, ok := c.(*StringPresetPattern); ok && ppc.Preset != "" {
				require.Contains(t, str, "&"+ppc.Preset)
			} else {
				require.Contains(t, str, "&"+name)
			}
		})
	}
}

func TestPropertyValidator_ToV8nTagString_WithUnknownConstraint(t *testing.T) {
	pv := &PropertyValidator{
		Constraints: Constraints{
			&testV8nUnknownConstraint{},
		},
	}
	require.Panics(t, func() {
		pv.ToV8nTagString(nil)
	})
}

func TestPropertyValidator_ToV8nTagString_WithAliasedConstraint(t *testing.T) {
	constraintsRegistry.registerNamed(true, "unknown", &testV8nUnknownConstraint{})
	defer constraintsRegistry.reset()
	pv := &PropertyValidator{
		Constraints: Constraints{
			&testV8nUnknownConstraint{},
		},
	}
	str := pv.ToV8nTagString(nil)
	require.Equal(t, "type:any, &unknown{}", str)
}

type testV8nUnknownConstraint struct {
	Field1 string
}

func (c *testV8nUnknownConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return false, ""
}
func (c *testV8nUnknownConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

func TestPropertyValidator_ToV8nTagString_WithBadConstraintField(t *testing.T) {
	pv := &PropertyValidator{
		Constraints: Constraints{
			&testV8nBadConstraint{Field1: &SubStruct{}},
		},
	}
	defer constraintsRegistry.reset()
	constraintsRegistry.register(false, &testV8nBadConstraint{})
	require.Panics(t, func() {
		str := pv.ToV8nTagString(nil)
		println(str)
	})
}

type testV8nBadConstraint struct {
	Field1 *SubStruct
}
type SubStruct struct {
	Foo string
}

func (c *testV8nBadConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return false, ""
}
func (c *testV8nBadConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

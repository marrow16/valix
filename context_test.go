package valix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateContext(t *testing.T) {
	vcx := newValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
}

func TestContextPathing(t *testing.T) {
	vcx := newValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathProperty("baz", nil)
	require.Equal(t, "baz", vcx.CurrentProperty())
	require.Equal(t, "baz", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo.bar", vcx.CurrentPath())

	vcx.pushPathProperty("qux", nil)
	require.Equal(t, "qux", vcx.CurrentProperty())
	require.Equal(t, "qux", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo.bar.baz", vcx.CurrentPath())

	vcx.popPath()
	require.Equal(t, "baz", vcx.CurrentProperty())
	require.Equal(t, "baz", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo.bar", vcx.CurrentPath())
}

func TestContextIndexPathing(t *testing.T) {
	vcx := newValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathIndex(0, nil)
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", vcx.CurrentPath())

	vcx.popPath()
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathIndex(1, nil)
	require.Equal(t, 1, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 1, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", vcx.CurrentPath())
	vcx.pushPathIndex(2, nil)
	require.Equal(t, 2, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 2, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0][1]", vcx.CurrentPath())
}

func TestContextIndexPathingFromRoot(t *testing.T) {
	vcx := newValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathIndex(0, nil)
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "[0]", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "[0].foo", vcx.CurrentPath())

	vcx = newValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
	vcx.pushPathIndex(0, nil)
	vcx.pushPathIndex(1, nil)
	require.Equal(t, 1, vcx.CurrentProperty())
	require.Equal(t, 1, *vcx.CurrentArrayIndex())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, "[0]", vcx.CurrentPath())
}

func TestPathPopNeverFails(t *testing.T) {
	vcx := newValidatorContext(nil)

	// push 4...
	vcx.pushPathProperty("foo1", nil)
	vcx.pushPathProperty("foo2", nil)
	vcx.pushPathProperty("foo3", nil)
	vcx.pushPathProperty("foo4", nil)
	require.Equal(t, "foo1.foo2.foo3", vcx.CurrentPath())
	require.Equal(t, "foo4", vcx.CurrentProperty())
	require.Equal(t, "foo4", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())

	// and pop 6...
	vcx.popPath()
	require.Equal(t, "foo3", vcx.CurrentProperty())
	require.Equal(t, "foo3", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo1.foo2", vcx.CurrentPath())
	vcx.popPath()
	require.Equal(t, "foo2", vcx.CurrentProperty())
	require.Equal(t, "foo2", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo1", vcx.CurrentPath())
	vcx.popPath()
	require.Equal(t, "foo1", vcx.CurrentProperty())
	require.Equal(t, "foo1", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
	vcx.popPath()
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
	vcx.popPath()
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
	vcx.popPath()
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
}

func TestContext_CurrentDepth(t *testing.T) {
	vcx := newValidatorContext(nil)
	require.Equal(t, 0, vcx.CurrentDepth())

	vcx.pushPathProperty("foo", nil)
	require.Equal(t, 1, vcx.CurrentDepth())
	vcx.pushPathIndex(16, nil)
	require.Equal(t, 2, vcx.CurrentDepth())
	vcx.pushPathProperty("bar", nil)
	require.Equal(t, 3, vcx.CurrentDepth())

	vcx.popPath()
	require.Equal(t, 2, vcx.CurrentDepth())
	vcx.popPath()
	require.Equal(t, 1, vcx.CurrentDepth())
	vcx.popPath()
	require.Equal(t, 0, vcx.CurrentDepth())
	vcx.popPath()
	require.Equal(t, 0, vcx.CurrentDepth())
}

func TestContext_AncestorPath(t *testing.T) {
	vcx := newValidatorContext(nil)
	ap, apok := vcx.AncestorPath(0)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("foo", nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)
	// check ancestor too far is no ok...
	ap, apok = vcx.AncestorPath(1)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("bar", nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)

	vcx.pushPathProperty("baz", nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)

	vcx.pushPathProperty("qux", nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo.bar", *ap)
}

func TestContext_AncestorProperty(t *testing.T) {
	vcx := newValidatorContext(nil)
	ap, apok := vcx.AncestorProperty(0)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("foo", nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("bar", nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "foo", ap)

	vcx.pushPathProperty("baz", nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "bar", ap)

	vcx.pushPathIndex(16, nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "baz", ap)

	vcx.pushPathProperty("qux", nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, 16, ap)

	ap, apok = vcx.AncestorProperty(1)
	require.True(t, apok)
	require.Equal(t, "baz", ap)
	ap, apok = vcx.AncestorProperty(2)
	require.True(t, apok)
	require.Equal(t, "bar", ap)
	ap, apok = vcx.AncestorProperty(3)
	require.True(t, apok)
	require.Equal(t, "foo", ap)
	ap, apok = vcx.AncestorProperty(4)
	require.True(t, apok)
	require.Nil(t, ap)
	_, apok = vcx.AncestorProperty(5)
	require.False(t, apok)
}

func TestContext_AncestorPropertyName(t *testing.T) {
	vcx := newValidatorContext(nil)
	vcx.pushPathProperty("foo", nil)
	vcx.pushPathProperty("bar", nil)
	vcx.pushPathProperty("baz", nil)
	vcx.pushPathIndex(16, nil)
	vcx.pushPathProperty("qux", nil)

	ap, apok := vcx.AncestorPropertyName(0)
	require.False(t, apok)
	ai, apok := vcx.AncestorArrayIndex(0)
	require.True(t, apok)
	require.Equal(t, 16, *ai)
	ap, apok = vcx.AncestorPropertyName(1)
	require.True(t, apok)
	require.Equal(t, "baz", *ap)
	ap, apok = vcx.AncestorPropertyName(2)
	require.True(t, apok)
	require.Equal(t, "bar", *ap)
	ap, apok = vcx.AncestorPropertyName(3)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)
	_, apok = vcx.AncestorPropertyName(4)
	require.False(t, apok)
	_, apok = vcx.AncestorPropertyName(5)
	require.False(t, apok)
}

func TestContext_AncestorArrayIndex(t *testing.T) {
	vcx := newValidatorContext(nil)
	vcx.pushPathIndex(0, nil)
	vcx.pushPathIndex(1, nil)
	vcx.pushPathIndex(2, nil)
	vcx.pushPathProperty("foo", nil)
	vcx.pushPathIndex(3, nil)

	ap, apok := vcx.AncestorPropertyName(0)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)
	_, apok = vcx.AncestorArrayIndex(0)
	require.False(t, apok)
	ai, apok := vcx.AncestorArrayIndex(1)
	require.True(t, apok)
	require.Equal(t, 2, *ai)
	ai, apok = vcx.AncestorArrayIndex(2)
	require.True(t, apok)
	require.Equal(t, 1, *ai)
	ai, apok = vcx.AncestorArrayIndex(3)
	require.True(t, apok)
	require.Equal(t, 0, *ai)
	_, apok = vcx.AncestorArrayIndex(4)
	require.False(t, apok)
	_, apok = vcx.AncestorArrayIndex(5)
	require.False(t, apok)
}

func TestContext_AncestorValue(t *testing.T) {
	arrItem0 := map[string]interface{}{
		"qux": true,
	}
	bazVal := []interface{}{
		arrItem0,
	}
	barVal := map[string]interface{}{
		"baz": bazVal,
	}
	fooVal := map[string]interface{}{
		"bar": barVal,
	}
	o := map[string]interface{}{
		"foo": fooVal,
	}
	const testMsg = "TEST MESSAGE"
	finalTestConstraint := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		av, ok := vcx.AncestorValue(0)
		require.True(t, ok)
		require.Equal(t, arrItem0, av)
		av, ok = vcx.AncestorValue(1)
		require.True(t, ok)
		require.Equal(t, bazVal, av)
		av, ok = vcx.AncestorValue(2)
		require.True(t, ok)
		require.Equal(t, barVal, av)
		av, ok = vcx.AncestorValue(3)
		require.True(t, ok)
		require.Equal(t, fooVal, av)
		av, ok = vcx.AncestorValue(4)
		require.True(t, ok)
		require.Equal(t, o, av)
		_, ok = vcx.AncestorValue(5)
		require.False(t, ok)
		return false, cc.GetMessage()
	}, testMsg)
	v := Validator{
		Properties: Properties{
			"foo": {
				Mandatory: true,
				NotNull:   true,
				ObjectValidator: &Validator{
					Properties: Properties{
						"bar": {
							Mandatory: true,
							NotNull:   true,
							ObjectValidator: &Validator{
								Properties: Properties{
									"baz": {
										Mandatory: true,
										NotNull:   true,
										ObjectValidator: &Validator{
											AllowArray: true,
											Properties: Properties{
												"qux": {
													Type: JsonBoolean,
													Constraints: Constraints{
														finalTestConstraint,
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
		},
	}

	ok, violations := v.Validate(o)
	// check that validator actually hit our test constraint...
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
}

func TestContext_SetCurrentValue(t *testing.T) {
	o := jsonObject(`{
		"foo": {
			"bar": true
		}
	}`)
	const testMsg = "TEST MESSAGE"
	v := Validator{
		Properties: Properties{
			"foo": {
				Mandatory: true,
				NotNull:   true,
				ObjectValidator: &Validator{
					Properties: Properties{
						"bar": {
							Type:      JsonBoolean,
							Mandatory: true,
							NotNull:   true,
							Constraints: Constraints{
								NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
									vcx.SetCurrentValue(false)
									return false, cc.GetMessage()
								}, testMsg),
							},
						},
					},
				},
			},
		},
	}

	// check that validator actually hit our test constraint...
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	// and check that foo.bar got set to false...
	foo := o["foo"].(map[string]interface{})
	fooBar := foo["bar"].(bool)
	require.False(t, fooBar)
}

func TestContext_SetCurrentValueInArray(t *testing.T) {
	o := jsonObject(`{
		"foo": [
			{
				"bar": true
			},
			{
				"bar": true
			}
		]
	}`)
	const testMsg = "TEST MESSAGE"
	const testValue = "TEST VALUE"
	v := Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonArray,
				Mandatory: true,
				NotNull:   true,
				ObjectValidator: &Validator{
					AllowArray: true,
					Constraints: Constraints{
						NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
							vcx.SetCurrentValue(testValue)
							return false, cc.GetMessage()
						}, testMsg),
					},
					Properties: Properties{
						"bar": {
							Type: JsonBoolean,
						},
					},
				},
			},
		},
	}

	// check that validator actually hit our test constraint...
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
	require.Equal(t, testMsg, violations[1].Message)

	// and check the values of each foo item were changed...
	foo := o["foo"].([]interface{})
	require.Equal(t, testValue, foo[0])
	require.Equal(t, testValue, foo[1])
}

func TestContext_SetCurrentValueOnRootFails(t *testing.T) {
	o := jsonObject(`{"foo": "bar"}`)
	const testMsg = "TEST MESSAGE"
	v := Validator{
		IgnoreUnknownProperties: true,
		Constraints: Constraints{
			NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
				return vcx.SetCurrentValue(nil), cc.GetMessage()
			}, testMsg),
		},
	}

	// check that validator actually hit our test constraint...
	ok, violations := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)
}

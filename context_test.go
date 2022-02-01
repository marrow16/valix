package valix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateContext(t *testing.T) {
	ctx := newContext(nil)
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
}

func TestContextPathing(t *testing.T) {
	ctx := newContext(nil)
	require.Nil(t, ctx.CurrentProperty())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.CurrentProperty())
	require.Equal(t, "foo", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.CurrentProperty())
	require.Equal(t, "bar", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo", ctx.CurrentPath())

	ctx.pushPathProperty("baz", nil)
	require.Equal(t, "baz", ctx.CurrentProperty())
	require.Equal(t, "baz", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo.bar", ctx.CurrentPath())

	ctx.pushPathProperty("qux", nil)
	require.Equal(t, "qux", ctx.CurrentProperty())
	require.Equal(t, "qux", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo.bar.baz", ctx.CurrentPath())

	ctx.popPath()
	require.Equal(t, "baz", ctx.CurrentProperty())
	require.Equal(t, "baz", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo.bar", ctx.CurrentPath())
}

func TestContextIndexPathing(t *testing.T) {
	ctx := newContext(nil)
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.CurrentProperty())
	require.Equal(t, "foo", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathIndex(0, nil)
	require.Equal(t, 0, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, 0, *ctx.CurrentArrayIndex())
	require.Equal(t, "foo", ctx.CurrentPath())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.CurrentProperty())
	require.Equal(t, "bar", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", ctx.CurrentPath())

	ctx.popPath()
	require.Equal(t, 0, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, 0, *ctx.CurrentArrayIndex())
	require.Equal(t, "foo", ctx.CurrentPath())

	ctx.pushPathIndex(1, nil)
	require.Equal(t, 1, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, 1, *ctx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", ctx.CurrentPath())
	ctx.pushPathIndex(2, nil)
	require.Equal(t, 2, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, 2, *ctx.CurrentArrayIndex())
	require.Equal(t, "foo[0][1]", ctx.CurrentPath())
}

func TestContextIndexPathingFromRoot(t *testing.T) {
	ctx := newContext(nil)
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathIndex(0, nil)
	require.Equal(t, 0, ctx.CurrentProperty())
	require.Equal(t, 0, *ctx.CurrentArrayIndex())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, "", ctx.CurrentPath())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.CurrentProperty())
	require.Equal(t, "foo", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "[0]", ctx.CurrentPath())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.CurrentProperty())
	require.Equal(t, "bar", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "[0].foo", ctx.CurrentPath())

	ctx = newContext(nil)
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
	ctx.pushPathIndex(0, nil)
	ctx.pushPathIndex(1, nil)
	require.Equal(t, 1, ctx.CurrentProperty())
	require.Equal(t, 1, *ctx.CurrentArrayIndex())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Equal(t, "[0]", ctx.CurrentPath())
}

func TestPathPopNeverFails(t *testing.T) {
	ctx := newContext(nil)

	// push 4...
	ctx.pushPathProperty("foo1", nil)
	ctx.pushPathProperty("foo2", nil)
	ctx.pushPathProperty("foo3", nil)
	ctx.pushPathProperty("foo4", nil)
	require.Equal(t, "foo1.foo2.foo3", ctx.CurrentPath())
	require.Equal(t, "foo4", ctx.CurrentProperty())
	require.Equal(t, "foo4", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())

	// and pop 6...
	ctx.popPath()
	require.Equal(t, "foo3", ctx.CurrentProperty())
	require.Equal(t, "foo3", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo1.foo2", ctx.CurrentPath())
	ctx.popPath()
	require.Equal(t, "foo2", ctx.CurrentProperty())
	require.Equal(t, "foo2", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "foo1", ctx.CurrentPath())
	ctx.popPath()
	require.Equal(t, "foo1", ctx.CurrentProperty())
	require.Equal(t, "foo1", *ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
	ctx.popPath()
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
	ctx.popPath()
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
	ctx.popPath()
	require.Nil(t, ctx.CurrentProperty())
	require.Nil(t, ctx.CurrentPropertyName())
	require.Nil(t, ctx.CurrentArrayIndex())
	require.Equal(t, "", ctx.CurrentPath())
}

func TestContext_CurrentDepth(t *testing.T) {
	ctx := newContext(nil)
	require.Equal(t, 0, ctx.CurrentDepth())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, 1, ctx.CurrentDepth())
	ctx.pushPathIndex(16, nil)
	require.Equal(t, 2, ctx.CurrentDepth())
	ctx.pushPathProperty("bar", nil)
	require.Equal(t, 3, ctx.CurrentDepth())

	ctx.popPath()
	require.Equal(t, 2, ctx.CurrentDepth())
	ctx.popPath()
	require.Equal(t, 1, ctx.CurrentDepth())
	ctx.popPath()
	require.Equal(t, 0, ctx.CurrentDepth())
	ctx.popPath()
	require.Equal(t, 0, ctx.CurrentDepth())
}

func TestContext_AncestorPath(t *testing.T) {
	ctx := newContext(nil)
	ap, apok := ctx.AncestorPath(0)
	require.False(t, apok)
	require.Nil(t, ap)

	ctx.pushPathProperty("foo", nil)
	ap, apok = ctx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)
	// check ancestor too far is no ok...
	ap, apok = ctx.AncestorPath(1)
	require.False(t, apok)
	require.Nil(t, ap)

	ctx.pushPathProperty("bar", nil)
	ap, apok = ctx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)

	ctx.pushPathProperty("baz", nil)
	ap, apok = ctx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)

	ctx.pushPathProperty("qux", nil)
	ap, apok = ctx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo.bar", *ap)
}

func TestContext_AncestorProperty(t *testing.T) {
	ctx := newContext(nil)
	ap, apok := ctx.AncestorProperty(0)
	require.False(t, apok)
	require.Nil(t, ap)

	ctx.pushPathProperty("foo", nil)
	ap, apok = ctx.AncestorProperty(0)
	require.True(t, apok)
	require.Nil(t, ap)

	ctx.pushPathProperty("bar", nil)
	ap, apok = ctx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "foo", ap)

	ctx.pushPathProperty("baz", nil)
	ap, apok = ctx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "bar", ap)

	ctx.pushPathIndex(16, nil)
	ap, apok = ctx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "baz", ap)

	ctx.pushPathProperty("qux", nil)
	ap, apok = ctx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, 16, ap)

	ap, apok = ctx.AncestorProperty(1)
	require.True(t, apok)
	require.Equal(t, "baz", ap)
	ap, apok = ctx.AncestorProperty(2)
	require.True(t, apok)
	require.Equal(t, "bar", ap)
	ap, apok = ctx.AncestorProperty(3)
	require.True(t, apok)
	require.Equal(t, "foo", ap)
	ap, apok = ctx.AncestorProperty(4)
	require.True(t, apok)
	require.Nil(t, ap)
	ap, apok = ctx.AncestorProperty(5)
	require.False(t, apok)
}

func TestContext_AncestorPropertyName(t *testing.T) {
	ctx := newContext(nil)
	ctx.pushPathProperty("foo", nil)
	ctx.pushPathProperty("bar", nil)
	ctx.pushPathProperty("baz", nil)
	ctx.pushPathIndex(16, nil)
	ctx.pushPathProperty("qux", nil)

	ap, apok := ctx.AncestorPropertyName(0)
	require.False(t, apok)
	ai, apok := ctx.AncestorArrayIndex(0)
	require.True(t, apok)
	require.Equal(t, 16, *ai)
	ap, apok = ctx.AncestorPropertyName(1)
	require.True(t, apok)
	require.Equal(t, "baz", *ap)
	ap, apok = ctx.AncestorPropertyName(2)
	require.True(t, apok)
	require.Equal(t, "bar", *ap)
	ap, apok = ctx.AncestorPropertyName(3)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)
	ap, apok = ctx.AncestorPropertyName(4)
	require.False(t, apok)
	ap, apok = ctx.AncestorPropertyName(5)
	require.False(t, apok)
}

func TestContext_AncestorArrayIndex(t *testing.T) {
	ctx := newContext(nil)
	ctx.pushPathIndex(0, nil)
	ctx.pushPathIndex(1, nil)
	ctx.pushPathIndex(2, nil)
	ctx.pushPathProperty("foo", nil)
	ctx.pushPathIndex(3, nil)

	ap, apok := ctx.AncestorPropertyName(0)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)
	_, apok = ctx.AncestorArrayIndex(0)
	require.False(t, apok)
	ai, apok := ctx.AncestorArrayIndex(1)
	require.True(t, apok)
	require.Equal(t, 2, *ai)
	ai, apok = ctx.AncestorArrayIndex(2)
	require.True(t, apok)
	require.Equal(t, 1, *ai)
	ai, apok = ctx.AncestorArrayIndex(3)
	require.True(t, apok)
	require.Equal(t, 0, *ai)
	ai, apok = ctx.AncestorArrayIndex(4)
	require.False(t, apok)
	ai, apok = ctx.AncestorArrayIndex(5)
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
	finalTestConstraint := NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
		av, ok := ctx.AncestorValue(0)
		require.True(t, ok)
		require.Equal(t, arrItem0, av)
		av, ok = ctx.AncestorValue(1)
		require.True(t, ok)
		require.Equal(t, bazVal, av)
		av, ok = ctx.AncestorValue(2)
		require.True(t, ok)
		require.Equal(t, barVal, av)
		av, ok = ctx.AncestorValue(3)
		require.True(t, ok)
		require.Equal(t, fooVal, av)
		av, ok = ctx.AncestorValue(4)
		require.True(t, ok)
		require.Equal(t, o, av)
		av, ok = ctx.AncestorValue(5)
		require.False(t, ok)
		return false, cc.GetMessage()
	}, testMsg)
	v := Validator{
		Properties: map[string]*PropertyValidator{
			"foo": {
				Mandatory: true,
				NotNull:   true,
				ObjectValidator: &Validator{
					Properties: map[string]*PropertyValidator{
						"bar": {
							Mandatory: true,
							NotNull:   true,
							ObjectValidator: &Validator{
								Properties: map[string]*PropertyValidator{
									"baz": {
										Mandatory: true,
										NotNull:   true,
										ObjectValidator: &Validator{
											AllowArray: true,
											Properties: map[string]*PropertyValidator{
												"qux": {
													PropertyType: PropertyType.Boolean,
													Constraints: []Constraint{
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				Mandatory: true,
				NotNull:   true,
				ObjectValidator: &Validator{
					Properties: map[string]*PropertyValidator{
						"bar": {
							PropertyType: PropertyType.Boolean,
							Mandatory:    true,
							NotNull:      true,
							Constraints: []Constraint{
								NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
									ctx.SetCurrentValue(false)
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
		Properties: map[string]*PropertyValidator{
			"foo": {
				PropertyType: PropertyType.Array,
				Mandatory:    true,
				NotNull:      true,
				ObjectValidator: &Validator{
					AllowArray: true,
					Constraints: []Constraint{
						NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
							ctx.SetCurrentValue(testValue)
							return false, cc.GetMessage()
						}, testMsg),
					},
					Properties: map[string]*PropertyValidator{
						"bar": {
							PropertyType: PropertyType.Boolean,
						},
					},
				},
			},
		},
	}

	// check that validator actually hit our test constraint...
	ok, violaions := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 2, len(violaions))
	require.Equal(t, testMsg, violaions[0].Message)
	require.Equal(t, testMsg, violaions[1].Message)

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
		Constraints: []Constraint{
			NewCustomConstraint(func(value interface{}, ctx *Context, cc *CustomConstraint) (bool, string) {
				return ctx.SetCurrentValue(nil), cc.GetMessage()
			}, testMsg),
		},
	}

	// check that validator actually hit our test constraint...
	ok, violaions := v.Validate(o)
	require.False(t, ok)
	require.Equal(t, 1, len(violaions))
	require.Equal(t, testMsg, violaions[0].Message)
}

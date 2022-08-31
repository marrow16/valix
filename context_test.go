package valix

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateContext(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
}

func TestContextPathing(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil, nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil, nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathProperty("baz", nil, nil)
	require.Equal(t, "baz", vcx.CurrentProperty())
	require.Equal(t, "baz", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo.bar", vcx.CurrentPath())

	vcx.pushPathProperty("qux", nil, nil)
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
	vcx := newEmptyValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil, nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathIndex(0, nil, nil)
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil, nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", vcx.CurrentPath())

	vcx.popPath()
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo", vcx.CurrentPath())

	vcx.pushPathIndex(1, nil, nil)
	require.Equal(t, 1, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 1, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0]", vcx.CurrentPath())
	vcx.pushPathIndex(2, nil, nil)
	require.Equal(t, 2, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, 2, *vcx.CurrentArrayIndex())
	require.Equal(t, "foo[0][1]", vcx.CurrentPath())
}

func TestContextIndexPathingFromRoot(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathIndex(0, nil, nil)
	require.Equal(t, 0, vcx.CurrentProperty())
	require.Equal(t, 0, *vcx.CurrentArrayIndex())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, "", vcx.CurrentPath())

	vcx.pushPathProperty("foo", nil, nil)
	require.Equal(t, "foo", vcx.CurrentProperty())
	require.Equal(t, "foo", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "[0]", vcx.CurrentPath())

	vcx.pushPathProperty("bar", nil, nil)
	require.Equal(t, "bar", vcx.CurrentProperty())
	require.Equal(t, "bar", *vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "[0].foo", vcx.CurrentPath())

	vcx = newEmptyValidatorContext(nil)
	require.Nil(t, vcx.CurrentProperty())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Nil(t, vcx.CurrentArrayIndex())
	require.Equal(t, "", vcx.CurrentPath())
	vcx.pushPathIndex(0, nil, nil)
	vcx.pushPathIndex(1, nil, nil)
	require.Equal(t, 1, vcx.CurrentProperty())
	require.Equal(t, 1, *vcx.CurrentArrayIndex())
	require.Nil(t, vcx.CurrentPropertyName())
	require.Equal(t, "[0]", vcx.CurrentPath())
}

func TestPathPopNeverFails(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	// push 4...
	vcx.pushPathProperty("foo1", nil, nil)
	vcx.pushPathProperty("foo2", nil, nil)
	vcx.pushPathProperty("foo3", nil, nil)
	vcx.pushPathProperty("foo4", nil, nil)
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
	vcx := newEmptyValidatorContext(nil)
	require.Equal(t, 0, vcx.CurrentDepth())

	vcx.pushPathProperty("foo", nil, nil)
	require.Equal(t, 1, vcx.CurrentDepth())
	vcx.pushPathIndex(16, nil, nil)
	require.Equal(t, 2, vcx.CurrentDepth())
	vcx.pushPathProperty("bar", nil, nil)
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
	vcx := newEmptyValidatorContext(nil)
	ap, apok := vcx.AncestorPath(0)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("foo", nil, nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)
	// check ancestor too far is no ok...
	ap, apok = vcx.AncestorPath(1)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("bar", nil, nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "", *ap)

	vcx.pushPathProperty("baz", nil, nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo", *ap)

	vcx.pushPathProperty("qux", nil, nil)
	ap, apok = vcx.AncestorPath(0)
	require.True(t, apok)
	require.Equal(t, "foo.bar", *ap)
}

func TestContext_AncestorProperty(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	ap, apok := vcx.AncestorProperty(0)
	require.False(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("foo", nil, nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Nil(t, ap)

	vcx.pushPathProperty("bar", nil, nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "foo", ap)

	vcx.pushPathProperty("baz", nil, nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "bar", ap)

	vcx.pushPathIndex(16, nil, nil)
	ap, apok = vcx.AncestorProperty(0)
	require.True(t, apok)
	require.Equal(t, "baz", ap)

	vcx.pushPathProperty("qux", nil, nil)
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
	vcx := newEmptyValidatorContext(nil)
	vcx.pushPathProperty("foo", nil, nil)
	vcx.pushPathProperty("bar", nil, nil)
	vcx.pushPathProperty("baz", nil, nil)
	vcx.pushPathIndex(16, nil, nil)
	vcx.pushPathProperty("qux", nil, nil)

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
	vcx := newEmptyValidatorContext(nil)
	vcx.pushPathIndex(0, nil, nil)
	vcx.pushPathIndex(1, nil, nil)
	vcx.pushPathIndex(2, nil, nil)
	vcx.pushPathProperty("foo", nil, nil)
	vcx.pushPathIndex(3, nil, nil)

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
		return false, cc.GetMessage(vcx)
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

func TestValidatorConditionSet(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	require.False(t, vcx.IsCondition("TEST"))

	vcx.SetCondition("TEST")
	require.True(t, vcx.IsCondition("TEST"))
}

func TestValidatorConditionSetWithExclamation(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	require.False(t, vcx.IsCondition("TEST"))

	vcx.SetCondition("TEST")
	require.True(t, vcx.IsCondition("TEST"))

	vcx.SetCondition("!TEST")
	require.False(t, vcx.IsCondition("TEST"))
}

func TestValidatorContextMeetsWhenConditions(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	vcx.SetCondition("TEST1")
	require.True(t, vcx.IsCondition("TEST1"))
	vcx.SetCondition("TEST2")
	require.True(t, vcx.IsCondition("TEST2"))

	require.True(t, vcx.meetsWhenConditions([]string{}))

	require.True(t, vcx.meetsWhenConditions([]string{"TEST1"}))
	require.True(t, vcx.meetsWhenConditions([]string{"TEST2"}))
	require.True(t, vcx.meetsWhenConditions([]string{"TEST1", "TEST2"}))
	require.False(t, vcx.meetsWhenConditions([]string{"!TEST1"}))
	require.False(t, vcx.meetsWhenConditions([]string{"TEST1", "!TEST1"}))
	require.False(t, vcx.meetsWhenConditions([]string{"FOO"}))

	require.True(t, vcx.meetsWhenConditions([]string{"!FOO"}))
	require.True(t, vcx.meetsWhenConditions([]string{"TEST1", "!FOO"}))
	require.False(t, vcx.meetsWhenConditions([]string{"TEST1", "!TEST2"}))
}

func TestValidatorContextMeetsUnwantedConditions(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)
	require.True(t, vcx.meetsUnwantedConditions([]string{}))

	vcx.SetCondition("TEA")
	require.True(t, vcx.meetsUnwantedConditions([]string{"!TEA"}))
	require.False(t, vcx.meetsUnwantedConditions([]string{"TEA"}))
	require.True(t, vcx.meetsUnwantedConditions([]string{"COFFEE"}))
	require.False(t, vcx.meetsUnwantedConditions([]string{"!COFFEE"}))
	vcx.SetCondition("COFFEE")
	require.True(t, vcx.meetsUnwantedConditions([]string{"!TEA"}))
	require.False(t, vcx.meetsUnwantedConditions([]string{"TEA"}))
	require.False(t, vcx.meetsUnwantedConditions([]string{"COFFEE"}))
	require.True(t, vcx.meetsUnwantedConditions([]string{"!COFFEE"}))
}

func TestValidatorContextConditionsPassedDownStack(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	vcx.pushPathProperty("test1", nil, nil)
	vcx.SetCondition("TEA")
	vcx.pushPathProperty("test2", nil, nil)
	vcx.pushPathProperty("test3", nil, nil)

	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
}

func TestValidatorContextSetParentCondition(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	vcx.pushPathProperty("test1", nil, nil)
	vcx.pushPathProperty("test2", nil, nil)
	vcx.pushPathProperty("test3", nil, nil)
	vcx.SetParentCondition("TEA")

	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
}

func TestValidatorContextSetParentNotCondition(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	vcx.pushPathProperty("test1", nil, nil)
	vcx.SetCondition("TEA")
	vcx.pushPathProperty("test2", nil, nil)
	vcx.pushPathProperty("test3", nil, nil)
	vcx.SetParentCondition("!TEA")

	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
}

func TestValidatorContextSetGlobalCondition(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	vcx.pushPathProperty("test1", nil, nil)
	vcx.pushPathProperty("test2", nil, nil)
	vcx.pushPathProperty("test3", nil, nil)
	vcx.SetGlobalCondition("TEA")

	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.True(t, vcx.IsCondition("TEA"))
}

func TestValidatorContextSetGlobalNotCondition(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	vcx.pushPathProperty("test1", nil, nil)
	vcx.SetCondition("TEA")
	require.True(t, vcx.IsCondition("TEA"))
	vcx.pushPathProperty("test2", nil, nil)
	vcx.pushPathProperty("test3", nil, nil)
	vcx.SetGlobalCondition("!TEA")

	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
	vcx.popPath()
	require.False(t, vcx.IsCondition("TEA"))
}

func TestObtainI18nContextAlwaysNonNil(t *testing.T) {
	i18ctx := obtainI18nContext(nil)
	require.NotNil(t, i18ctx)

	defer func() {
		DefaultI18nProvider = &defaultI18nProvider{}
	}()
	DefaultI18nProvider = &dummyI18nProvider{}
	i18ctx = obtainI18nContext(nil)
	require.NotNil(t, i18ctx)
}

type dummyI18nProvider struct{}

func (i *dummyI18nProvider) ContextFromRequest(r *http.Request) I18nContext {
	return nil
}
func (i *dummyI18nProvider) DefaultContext() I18nContext {
	return nil
}

func TestValuesAncestry(t *testing.T) {
	obj3 := map[string]interface{}{
		"level": 3,
	}
	obj2 := map[string]interface{}{
		"foo":   obj3,
		"level": 2,
	}
	arr2 := []interface{}{obj2}
	obj1 := map[string]interface{}{
		"foo":   arr2,
		"level": 1,
	}
	rootObj := map[string]interface{}{
		"foo":   obj1,
		"level": 0,
	}
	vcx := newValidatorContext(rootObj, nil, false, nil)

	vcx.pushPathProperty("foo", obj1, nil)
	vcx.pushPathProperty("foo", arr2, nil)
	vcx.pushPathIndex(0, obj2, nil)
	vcx.pushPathProperty("foo", obj3, nil)

	values := vcx.ValuesAncestry()
	require.Equal(t, 5, len(values))

	require.Equal(t, rootObj, values[4])
	av, _ := vcx.AncestorValue(3)
	require.Equal(t, rootObj, av)

	require.Equal(t, obj1, values[3])
	av, _ = vcx.AncestorValue(2)
	require.Equal(t, obj1, av)

	require.Equal(t, arr2, values[2])
	av, _ = vcx.AncestorValue(1)
	require.Equal(t, arr2, av)

	require.Equal(t, obj2, values[1])
	av, _ = vcx.AncestorValue(0)
	require.Equal(t, obj2, av)

	require.Equal(t, obj3, values[0])
	av = vcx.CurrentValue()
	require.Equal(t, obj3, av)
}

func TestAncestryIndex(t *testing.T) {
	obj3 := map[string]interface{}{
		"level": 3,
	}
	obj2 := map[string]interface{}{
		"foo":   obj3,
		"level": 2,
	}
	arr2 := []interface{}{obj2, obj2, obj2, obj2}
	obj1 := map[string]interface{}{
		"foo":   arr2,
		"level": 1,
	}
	rootObj := map[string]interface{}{
		"foo":   obj1,
		"level": 0,
	}

	vcx := newValidatorContext(rootObj, nil, false, nil)
	// no pathing as yet...
	_, _, ok := vcx.AncestryIndex(0)
	require.False(t, ok)

	vcx.pushPathProperty("foo", obj1, nil)
	_, _, ok = vcx.AncestryIndex(0)
	require.False(t, ok)

	vcx.pushPathProperty("foo", arr2, nil)
	_, _, ok = vcx.AncestryIndex(0)
	require.False(t, ok)

	vcx.pushPathIndex(1, obj2, nil)
	index, max, ok := vcx.AncestryIndex(0)
	require.True(t, ok)
	require.Equal(t, 1, index)
	require.Equal(t, 3, max)

	_, _, ok = vcx.AncestryIndex(1)
	require.False(t, ok)

	// test the odd situation (which should never occur) where things are pushed with incoherent order (i.e. an index into an object)...
	vcx = newValidatorContext(rootObj, nil, false, nil)
	vcx.pushPathIndex(16, obj2, nil)
	index, max, ok = vcx.AncestryIndex(0)
	require.True(t, ok)
	require.Equal(t, 16, index)
	require.Equal(t, -1, max)
}

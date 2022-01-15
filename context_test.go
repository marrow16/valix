package valix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateContext(t *testing.T) {
	ctx := newContext(nil, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
}

func TestContextPathing(t *testing.T) {
	ctx := newContext(nil, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.PropertyName())
	require.Equal(t, "foo", ctx.Path())

	ctx.pushPathProperty("baz", nil)
	require.Equal(t, "baz", ctx.PropertyName())
	require.Equal(t, "foo.bar", ctx.Path())

	ctx.pushPathProperty("qux", nil)
	require.Equal(t, "qux", ctx.PropertyName())
	require.Equal(t, "foo.bar.baz", ctx.Path())

	ctx.popPath()
	require.Equal(t, "baz", ctx.PropertyName())
	require.Equal(t, "foo.bar", ctx.Path())
}

func TestContextIndexPathing(t *testing.T) {
	ctx := newContext(nil, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())

	ctx.pushPathIndex(0, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "foo[0]", ctx.Path())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.PropertyName())
	require.Equal(t, "foo[0]", ctx.Path())

	ctx.popPath()
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "foo[0]", ctx.Path())

	ctx.pushPathIndex(0, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "foo[0][0]", ctx.Path())
}

func TestContextIndexPathingFromRoot(t *testing.T) {
	ctx := newContext(nil, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())

	ctx.pushPathIndex(0, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "[0]", ctx.Path())

	ctx.pushPathProperty("foo", nil)
	require.Equal(t, "foo", ctx.PropertyName())
	require.Equal(t, "[0]", ctx.Path())

	ctx.pushPathProperty("bar", nil)
	require.Equal(t, "bar", ctx.PropertyName())
	require.Equal(t, "[0].foo", ctx.Path())

	ctx = newContext(nil, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
	ctx.pushPathIndex(0, nil)
	ctx.pushPathIndex(0, nil)
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "[0][0]", ctx.Path())
}

func TestPathPopNeverFails(t *testing.T) {
	ctx := newContext(nil, nil)

	// push 4...
	ctx.pushPathProperty("foo1", nil)
	ctx.pushPathProperty("foo2", nil)
	ctx.pushPathProperty("foo3", nil)
	ctx.pushPathProperty("foo4", nil)
	require.Equal(t, "foo1.foo2.foo3", ctx.Path())
	require.Equal(t, "foo4", ctx.PropertyName())

	// and pop 6...
	ctx.popPath()
	require.Equal(t, "foo3", ctx.PropertyName())
	require.Equal(t, "foo1.foo2", ctx.Path())
	ctx.popPath()
	require.Equal(t, "foo2", ctx.PropertyName())
	require.Equal(t, "foo1", ctx.Path())
	ctx.popPath()
	require.Equal(t, "foo1", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
	ctx.popPath()
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
	ctx.popPath()
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
	ctx.popPath()
	require.Equal(t, "", ctx.PropertyName())
	require.Equal(t, "", ctx.Path())
}

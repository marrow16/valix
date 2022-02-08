package valix

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCanCreateCustomConstraint(t *testing.T) {
	cc := NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage()
	}, "")
	require.NotNil(t, cc)
}

func TestCustomConstraintStoresMessage(t *testing.T) {
	const testMsg = "TEST MESSAGE"
	cc := NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage()
	}, testMsg)
	require.Equal(t, testMsg, cc.GetMessage())
}

func TestCustomConstraint(t *testing.T) {
	const testMsg = "Value must be greater than 'B'"
	validator := buildFooValidator(JsonString,
		NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, cc.GetMessage()
			}
			return true, ""
		}, testMsg), false)
	obj := jsonObject(`{
		"foo": "OK is greater than B"
	}`)

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	obj["foo"] = "A"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = "B"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, testMsg, violations[0].Message)

	obj["foo"] = "Ba"
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
}

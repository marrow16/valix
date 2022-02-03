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
	msg := "Value must be greater than 'B'"
	validator := buildFooValidator(PropertyType.String,
		NewCustomConstraint(func(value interface{}, ctx *ValidatorContext, cc *CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, cc.GetMessage()
			}
			return true, ""
		}, msg), false)
	jobj := jsonObject(`{
		"foo": "OK is greater than B"
	}`)

	ok, violations := validator.Validate(jobj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	jobj["foo"] = "A"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	jobj["foo"] = "B"
	ok, violations = validator.Validate(jobj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msg, violations[0].Message)

	jobj["foo"] = "Ba"
	ok, violations = validator.Validate(jobj)
	require.True(t, ok)
}

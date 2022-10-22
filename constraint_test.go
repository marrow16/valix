package valix

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestCanCreateCustomConstraint(t *testing.T) {
	cc := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage(vcx)
	}, "")
	require.NotNil(t, cc)
}

func TestCustomConstraintStoresMessage(t *testing.T) {
	const testMsg = "TEST MESSAGE"
	cc := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
		return false, cc.GetMessage(vcx)
	}, testMsg)
	require.Equal(t, testMsg, cc.GetMessage(nil))
}

func TestCustomConstraint(t *testing.T) {
	const testMsg = "Value must be greater than 'B'"
	validator := buildFooValidator(JsonString,
		NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, cc *CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				return strings.Compare(str, "B") > 0, cc.GetMessage(vcx)
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
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestIsConditional(t *testing.T) {
	c := &ConditionalConstraint{}
	_, isCond := isConditional(c)
	require.True(t, isCond)

	c2 := &StringNotEmpty{}
	_, isCond = isConditional(c2)
	require.False(t, isCond)
}

type testConstraint struct {
	passes bool
	msg    string
	stops  bool
}

func (c *testConstraint) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	vcx.CeaseFurtherIf(c.stops)
	return c.passes, c.msg
}
func (c *testConstraint) GetMessage(tcx I18nContext) string {
	return c.msg
}

func buildFooValidator(propertyType JsonType, constraint Constraint, notNull bool) *Validator {
	return &Validator{
		Properties: Properties{
			"foo": {
				Type:        propertyType,
				NotNull:     notNull,
				Mandatory:   true,
				Constraints: Constraints{constraint},
			},
		},
	}
}

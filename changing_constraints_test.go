package valix

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
	"testing"
)

func TestStringTrimConstraint(t *testing.T) {
	validator := buildFooValidator(JsonString, &StringTrim{}, true)
	obj := map[string]interface{}{
		"foo": "   \t",
	}

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "", obj["foo"])
}

func TestStringTrimConstraintWithCutset(t *testing.T) {
	validator := buildFooValidator(JsonString, &StringTrim{Cutset: "AC"}, true)
	obj := map[string]interface{}{
		"foo": "AAABCCC",
	}

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "B", obj["foo"])
}

func TestStringTrimConstraintWithFollowingLengthConstraint(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&StringTrim{},
					&StringLength{Maximum: 1},
				},
			},
		},
	}
	obj := map[string]interface{}{
		"foo": "   A   ",
	}

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "A", obj["foo"])
}

func TestStringNormalizeUnicodeConstraint(t *testing.T) {
	validator := buildFooValidator(JsonString, &StringNormalizeUnicode{Form: norm.NFC}, true)
	obj := map[string]interface{}{
		"foo": "\u0063\u0327", // is 'c' (u+0063) followed by combining cedilla (U+0327)
	}
	// before normalization the string should be 3 bytes (1 for 'c' and 2 for the combining cedilla)
	require.Equal(t, 3, len(obj["foo"].(string)))

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, 2, len(obj["foo"].(string)))
	require.Equal(t, "\u00e7", obj["foo"])
}

func TestNormalizeUnicodeConstraintWithFollowingLengthConstraint(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type:      JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: Constraints{
					&StringNormalizeUnicode{Form: norm.NFC},
					&StringLength{Maximum: 1, UseRuneLen: true},
				},
			},
		},
	}
	obj := map[string]interface{}{
		"foo": "\u0063\u0327", // is 'c' (u+0063) followed by combining cedilla (U+0327)
	}

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, 2, len(obj["foo"].(string)))
	require.Equal(t, "\u00e7", obj["foo"])
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

package valix

import (
	"github.com/stretchr/testify/require"
	"golang.org/x/text/unicode/norm"
	"testing"
)

func TestStringTrimConstraint(t *testing.T) {
	validator := buildFooValidator(PropertyType.String, &StringTrim{}, true)
	obj := map[string]interface{}{
		"foo": "   \t",
	}

	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "", obj["foo"])
}

func TestStringTrimConstraintWithCutset(t *testing.T) {
	validator := buildFooValidator(PropertyType.String, &StringTrim{Cutset: "AC"}, true)
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
				PropertyType: PropertyType.String,
				NotNull:      true,
				Mandatory:    true,
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
	validator := buildFooValidator(PropertyType.String, &StringNormalizeUnicode{Form: norm.NFC}, true)
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
				PropertyType: PropertyType.String,
				NotNull:      true,
				Mandatory:    true,
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

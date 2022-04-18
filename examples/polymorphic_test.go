package examples

import (
	"testing"

	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
)

var (
	polymorphicBeverageValidator = &valix.Validator{
		Constraints: valix.Constraints{
			&valix.SetConditionProperty{
				PropertyName: "type",
			},
		},
		Properties: valix.Properties{
			"type": {
				Type:      valix.JsonString,
				Mandatory: true,
				NotNull:   true,
				Constraints: valix.Constraints{
					&valix.StringValidToken{Tokens: []string{"tea", "coffee", "soft"}},
				},
			},
			"quantity": {
				Type:      valix.JsonInteger,
				Mandatory: true,
				NotNull:   true,
				Constraints: valix.Constraints{
					&valix.Positive{},
				},
			},
		},
		ConditionalVariants: valix.ConditionalVariants{
			{
				WhenConditions: []string{"tea"},
				Properties: valix.Properties{
					"blend": {
						Type:      valix.JsonString,
						Mandatory: true,
						NotNull:   true,
						Constraints: valix.Constraints{
							&valix.StringValidToken{Tokens: []string{"Earl Grey", "English Breakfast", "Masala Chai"}},
						},
					},
				},
			},
			{
				WhenConditions: []string{"coffee"},
				Properties: valix.Properties{
					"roast": {
						Type:      valix.JsonString,
						Mandatory: true,
						NotNull:   true,
						Constraints: valix.Constraints{
							&valix.StringValidToken{Tokens: []string{"light", "medium", "dark"}},
						},
					},
				},
			},
			{
				WhenConditions: []string{"soft"},
				Constraints: valix.Constraints{
					&valix.SetConditionProperty{
						PropertyName: "brand",
						Prefix:       "soft-",
						Mapping: map[string]string{
							"Coca Cola": "coke",
							"Fanta":     "fanta",
							"Tango":     "tango",
						},
					},
				},
				Properties: valix.Properties{
					"brand": {
						Type:      valix.JsonString,
						Mandatory: true,
						NotNull:   true,
						Constraints: valix.Constraints{
							&valix.StringValidToken{Tokens: []string{"Coca Cola", "Fanta", "Tango"}},
						},
					},
					"flavor": {
						Type:      valix.JsonString,
						Mandatory: false,
						NotNull:   false,
					},
				},
				ConditionalVariants: valix.ConditionalVariants{
					{
						WhenConditions: []string{"soft-coke"},
						Properties: valix.Properties{
							"flavor": {
								Type:      valix.JsonString,
								Mandatory: false,
								NotNull:   true,
								Constraints: valix.Constraints{
									&valix.StringValidToken{Tokens: []string{"Regular", "Diet", "Zero", "Cherry"}},
								},
							},
						},
					},
					{
						WhenConditions: []string{"soft-fanta"},
						Properties: valix.Properties{
							"flavor": {
								Type:      valix.JsonString,
								Mandatory: true,
								NotNull:   true,
								Constraints: valix.Constraints{
									&valix.StringValidToken{Tokens: []string{"Orange", "Pineapple", "Strawberry", "Grape", "Pina Colada", "Dragon Fruit"}},
								},
							},
						},
					},
					{
						WhenConditions: []string{"soft-tango"},
						Properties: valix.Properties{
							"flavor": {
								Type:      valix.JsonString,
								Mandatory: true,
								NotNull:   true,
								Constraints: valix.Constraints{
									&valix.StringValidToken{Tokens: []string{"Orange", "Apple", "Strawberry", "Watermelon", "Tropical"}},
								},
							},
						},
					},
				},
			},
		},
	}
)

func TestValidator(t *testing.T) {
	obj := map[string]interface{}{
		"type":     "soft",
		"quantity": 1,
		"brand":    "Tango",
		"flavor":   "Unknown",
	}
	ok, violations := polymorphicBeverageValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "flavor", violations[0].Property)

	obj["flavor"] = "Apple"
	ok, _ = polymorphicBeverageValidator.Validate(obj)
	require.True(t, ok)

	obj["type"] = "tea"
	ok, violations = polymorphicBeverageValidator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	valix.SortViolationsByPathAndProperty(violations)
	require.Equal(t, "blend", violations[0].Property)
	require.Equal(t, "Missing property", violations[0].Message)
	require.Equal(t, "brand", violations[1].Property)
	require.Equal(t, "Invalid property", violations[1].Message)
	require.Equal(t, "flavor", violations[2].Property)
	require.Equal(t, "Invalid property", violations[2].Message)
}

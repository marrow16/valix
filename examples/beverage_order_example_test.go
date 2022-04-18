package examples

import (
	"testing"

	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
)

/*
type InitialBeverageOrder struct {
	Type     string `json:"type" v8n:"notNull,required,&StringValidToken{Tokens:['tea','coffee']}"`
	Quantity int    `json:"quantity" v8n:"notNull,required,&Positive{}"`
	// only relevant to type="tea"...
	Blend string `json:"blend" v8n:"notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
	// only relevant to type="coffee"...
	Roast string `json:"roast" v8n:"notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}

var InitialBeverageOrderValidator = valix.MustCompileValidatorFor(InitialBeverageOrder{}, nil)
*/

var BeverageOrderValidator = valix.MustCompileValidatorFor(BeverageOrder{}, nil)

type BeverageOrder struct {
	Type string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{Parent:true}"`
	//                                             sets the condition from the value of property 'type' ^^^
	Quantity int `json:"quantity" v8n:"notNull,required,&Positive{}"`
	// only relevant to type="tea"...
	Blend string `json:"blend" v8n:"when:tea,notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
	// only relevant to type="coffee"...
	Roast string `json:"roast" v8n:"when:coffee,notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}

func TestBeverageOrderValidatorIsCorrect(t *testing.T) {
	require.NotNil(t, BeverageOrderValidator)
	require.False(t, BeverageOrderValidator.OrderedPropertyChecks)
	require.True(t, BeverageOrderValidator.IsOrderedPropertyChecks())
	require.Equal(t, -1, BeverageOrderValidator.Properties["type"].Order)
	require.Equal(t, 1, len(BeverageOrderValidator.Properties["blend"].WhenConditions))
	require.Equal(t, 1, len(BeverageOrderValidator.Properties["roast"].WhenConditions))
}

func TestBeverageOrderValidator(t *testing.T) {
	json := `{
		"type": "tea",
		"quantity": 1,
		"blend": "Earl Grey"
	}`
	req := &BeverageOrder{}
	ok, violations, _ := BeverageOrderValidator.ValidateStringInto(json, req)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "tea", req.Type)
	require.Equal(t, "Earl Grey", req.Blend)
	require.Equal(t, "", req.Roast)

	json = `{
		"type": "coffee",
		"quantity": 1,
		"roast": "medium"
	}`
	req = &BeverageOrder{}
	ok, violations, _ = BeverageOrderValidator.ValidateStringInto(json, req)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "coffee", req.Type)
	require.Equal(t, "medium", req.Roast)
	require.Equal(t, "", req.Blend)
}

func TestBeverageOrderValidatorWithBadType(t *testing.T) {
	json := `{
		"type": "unknown",
		"quantity": 1,
		"blend": "wont be checked!",
		"roast": "wont be checked!"
	}`
	req := &BeverageOrder{}
	ok, violations, _ := BeverageOrderValidator.ValidateStringInto(json, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "String value must be valid token - \"tea\",\"coffee\"", violations[0].Message)
}

func TestBeverageOrderValidatorTeaWithMissingBlend(t *testing.T) {
	json := `{
		"type": "tea",
		"quantity": 1
	}`
	req := &BeverageOrder{}
	ok, violations, _ := BeverageOrderValidator.ValidateStringInto(json, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Missing property", violations[0].Message)
	require.Equal(t, "blend", violations[0].Property)
}

func TestBeverageOrderValidatorCoffeeWithMissingRoast(t *testing.T) {
	json := `{
		"type": "coffee",
		"quantity": 1
	}`
	req := &BeverageOrder{}
	ok, violations, _ := BeverageOrderValidator.ValidateStringInto(json, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Missing property", violations[0].Message)
	require.Equal(t, "roast", violations[0].Property)
}

var BeverageOrderStrictValidator = valix.MustCompileValidatorFor(BeverageOrderStrict{}, nil)

// BeverageOrderStrict is a stricter version of BeverageOrder - where if appropriate type has not been set then
// the 'blend' and 'roast' properties should not be present (because they have the `unwanted` v8n tag set)
type BeverageOrderStrict struct {
	Type string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{Parent:true}"`
	//                                                sets the condition from the value of property 'type' ^^^
	Quantity int `json:"quantity" v8n:"notNull,required,&Positive{}"`
	// Blend is only relevant when type="tea"...
	Blend string `json:"blend" v8n:"when:tea,unwanted:!tea,notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
	//                                                   ^^^ when type != "tea" this property is unwanted
	// Roast is only relevant when type="coffee"...
	Roast string `json:"roast" v8n:"when:coffee,unwanted:!coffee,notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
	//                                                      ^^^ when type != "coffee" this property is unwanted
}

func TestBeverageOrderStrictValidatorIsCorrect(t *testing.T) {
	require.NotNil(t, BeverageOrderStrictValidator)
	require.False(t, BeverageOrderStrictValidator.OrderedPropertyChecks)
	require.True(t, BeverageOrderStrictValidator.IsOrderedPropertyChecks())
	require.Equal(t, -1, BeverageOrderStrictValidator.Properties["type"].Order)
	require.Equal(t, 1, len(BeverageOrderStrictValidator.Properties["blend"].WhenConditions))
	require.Equal(t, 1, len(BeverageOrderStrictValidator.Properties["roast"].WhenConditions))
	require.Equal(t, 1, len(BeverageOrderStrictValidator.Properties["blend"].UnwantedConditions))
	require.Equal(t, 1, len(BeverageOrderStrictValidator.Properties["roast"].UnwantedConditions))
}

func TestBeverageOrderStrictValidator(t *testing.T) {
	json := `{
		"type": "tea",
		"quantity": 1,
		"blend": "Earl Grey",
		"roast": "medium"
	}`
	req := &BeverageOrderStrict{}
	ok, violations, _ := BeverageOrderStrictValidator.ValidateStringInto(json, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Property must not be present", violations[0].Message)
	require.Equal(t, "roast", violations[0].Property)

	json = `{
		"type": "coffee",
		"quantity": 1,
		"roast": "medium",
		"blend": "Earl Grey"
	}`
	req = &BeverageOrderStrict{}
	ok, violations, _ = BeverageOrderStrictValidator.ValidateStringInto(json, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Property must not be present", violations[0].Message)
	require.Equal(t, "blend", violations[0].Property)
}

func TestBeverageOrderStrictValidatorExpressedAsCode(t *testing.T) {
	v := &valix.Validator{
		Properties: valix.Properties{
			"type": {
				Type:      valix.JsonString,
				NotNull:   true,
				Mandatory: true,
				Order:     -1,
				Constraints: valix.Constraints{
					&valix.StringValidToken{Tokens: []string{"tea", "coffee"}},
					&valix.SetConditionFrom{Parent: true},
				},
			},
			"quantity": {
				Type:        valix.JsonInteger,
				NotNull:     true,
				Mandatory:   true,
				Constraints: valix.Constraints{&valix.Positive{}},
			},
			"blend": {
				Type:      valix.JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: valix.Constraints{
					&valix.StringValidToken{Tokens: []string{"Earl Grey", "English Breakfast", "Masala Chai"}},
				},
				WhenConditions:     []string{"tea"},
				UnwantedConditions: []string{"!tea"},
			},
			"roast": {
				Type:      valix.JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: valix.Constraints{
					&valix.StringValidToken{Tokens: []string{"light", "medium", "dark"}},
				},
				WhenConditions:     []string{"coffee"},
				UnwantedConditions: []string{"!coffee"},
			},
		},
	}

	json := `{
		"type": "tea",
		"quantity": 1,
		"blend": "Earl Grey",
		"roast": "medium"
	}`
	ok, violations, _ := v.ValidateString(json)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Property must not be present", violations[0].Message)
	require.Equal(t, "roast", violations[0].Property)

	json = `{
		"type": "tea",
		"quantity": 1,
		"blend": "Earl Grey"
	}`
	ok, _, _ = v.ValidateString(json)
	require.True(t, ok)

	json = `{
		"type": "coffee",
		"quantity": 1,
		"roast": "medium",
		"blend": "Earl Grey"
	}`
	ok, violations, _ = v.ValidateString(json)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Property must not be present", violations[0].Message)
	require.Equal(t, "blend", violations[0].Property)

	json = `{
		"type": "coffee",
		"quantity": 1,
		"roast": "medium"
	}`
	ok, _, _ = v.ValidateString(json)
	require.True(t, ok)
}

// Package valix
/*
- Go package for validating requests

Check requests in the form of *http.Request, `map[string]interface{}` or `[]interface{}`

Validators can be expressed effectively as, for example:
	var personValidator = &valix.Validator{
		IgnoreUnknownProperties: false,
		Properties: valix.Properties{
			"name": {
				PropertyType: valix.PropertyType.String,
				NotNull:      true,
				Mandatory:    true,
				Constraints:  valix.Constraints{
					&valix.StringLength{Minimum: 1, Maximum: 255},
				},
			},
			"age": {
				PropertyType: valix.PropertyType.Int,
				NotNull:      true,
				Mandatory:    true,
				Constraints:  valix.Constraints{
					&valix.PositiveOrZero{},
				},
			},
		},
	}

Validators can re-use common property validators, for example re-using the `personValidator` above:
	var addPersonToGroupValidator = &valix.Validator{
		IgnoreUnknownProperties: false,
		Properties: valix.Properties{
			"person": {
				PropertyType: valix.PropertyType.Object,
				ObjectValidator: personValidator,
			},
			"group": {
				PropertyType: valix.PropertyType.String,
				NotNull:      true,
				Mandatory:    true,
				Constraints:  valix.Constraints{
					&valix.StringLength{Minimum: 1, Maximum: 255},
				},
			},
		},
	}
*/
package valix

// Package valix - Go package for validating requests
/*

Check requests in the form of *http.Request, `map[string]interface{}` or `[]interface{}`

Validators can be created from existing structs, for example:
	type AddPersonRequest struct {
		Name string `json:"name" v8n:"notNull,mandatory,&StringNoControlCharacters{},&StringLength{Minimum: 1, Maximum: 255}"`
		Age int `json:"age" v8n:"type:Integer,notNull,mandatory,&PositiveOrZero{}"`
	}
	var AddPersonRequestValidator = valix.MustCompileValidatorFor(AddPersonRequest{}, nil)


Or validators can be expressed effectively in code, for example:
	var personValidator = &valix.Validator{
		IgnoreUnknownProperties: false,
		Properties: valix.Properties{
			"name": {
				Type: valix.Type.JsonString,
				NotNull:      true,
				Mandatory:    true,
				Constraints:  valix.Constraints{
					&valix.StringLength{Minimum: 1, Maximum: 255},
				},
			},
			"age": {
				Type: valix.Type.Int,
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
				Type: valix.Type.JsonObject,
				ObjectValidator: personValidator,
			},
			"group": {
				Type: valix.Type.JsonString,
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

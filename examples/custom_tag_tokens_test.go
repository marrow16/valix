package examples

import (
	"strings"
	"testing"

	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
)

type CustomTagToken struct{}

func (c *CustomTagToken) Handle(token string, hasValue bool, tokenValue string, pv *valix.PropertyValidator, propertyName string, fieldName string) error {
	pv.Mandatory = true
	pv.NotNull = true
	pv.Constraints = append(pv.Constraints, &valix.StringNotEmpty{})
	return nil
}

func TestExampleCustomTagToken(t *testing.T) {
	// register custom tag token "my_mnnne" (mandatory, not null, not empty)...
	valix.RegisterCustomTagToken("my_mnnne", &CustomTagToken{})

	type ExampleTagTokenExt struct {
		Foo string `json:"foo" v8n:"my_mnnne"`
	}
	// and build validator for struct that uses the custom tag token...
	validator := valix.MustCompileValidatorFor(ExampleTagTokenExt{}, nil)

	// and now validate using the validator...
	req := &ExampleTagTokenExt{}

	reader := strings.NewReader(`{
		"foo": "here"
	}`)
	ok, _, _ := validator.ValidateReaderInto(reader, req)
	require.True(t, ok)

	reader = strings.NewReader(`{
		"foo": ""
	}`)
	ok, violations, _ := validator.ValidateReaderInto(reader, req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

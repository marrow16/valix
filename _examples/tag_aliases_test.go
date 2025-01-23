package examples

import (
	"strings"
	"testing"

	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
)

func TestTagAliases(t *testing.T) {
	// register a tag token alias "mnnne" (mandatory, not null, string not empty)...
	valix.RegisterTagTokenAlias("mnnne", "mandatory, notNull, &StringNotEmpty{}")

	// use the alias in a `v8n` tag...
	type ExampleTagAliasUsage struct {
		Foo string `json:"foo" v8n:"$mnnne"`
	}
	validator := valix.MustCompileValidatorFor(ExampleTagAliasUsage{}, nil)

	eg := &ExampleTagAliasUsage{}

	reader := strings.NewReader(`{
		"foo": "here"
	}`)
	ok, _, _ := validator.ValidateReaderInto(reader, eg)
	require.True(t, ok)

	reader = strings.NewReader(`{
		"foo": ""
	}`)
	ok, violations, _ := validator.ValidateReaderInto(reader, eg)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
}

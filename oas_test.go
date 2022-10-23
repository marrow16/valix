package valix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidatorForWithOasTag(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo" oas:"description:'THIS IS DESC!',title:'THIS IS TITLE',format:'THIS IS FORMAT',example:'THIS IS EXAMPLE',deprecated"`
	}
	v, err := ValidatorFor(myStruct{}, nil)
	require.NoError(t, err)
	require.NotNil(t, v)

	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.True(t, pv.OasInfo.Deprecated)
	require.Equal(t, "THIS IS DESC!", pv.OasInfo.Description)
	require.Equal(t, "THIS IS TITLE", pv.OasInfo.Title)
	require.Equal(t, "THIS IS FORMAT", pv.OasInfo.Format)
	require.Equal(t, "THIS IS EXAMPLE", pv.OasInfo.Example)
}

func TestValidatorForWithOasTagUnknownTokenFails(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo" oas:"UNKNOWN_TOKEN"`
	}
	_, err := ValidatorFor(myStruct{}, nil)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgOasUnknownTokenInTag, "UNKNOWN_TOKEN"), err.Error())
}

func TestValidatorForWithBadOasTag(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo" oas:"title:'<unclosed single quote"`
	}
	_, err := ValidatorFor(myStruct{}, nil)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 6), err.Error())
}

func TestOasTagItemExpectedString(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addOasTagItem("desc:xxx")
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgOasExpectedString, "desc"), err.Error())
}

func TestOasTagItemUnexpectedColon(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addOasTagItem("deprecated:xxx")
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgOasUnexpectedColon, "deprecated"), err.Error())
}

package examples

import (
	"encoding/json"
	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagMandatory(t *testing.T) {
	type Example struct {
		Foo string `v8n:"mandatory"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["Foo"].Mandatory)
}

func TestTagNotNull(t *testing.T) {
	type Example struct {
		Foo string `v8n:"notNull"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["Foo"].NotNull)
}

func TestTagNullable(t *testing.T) {
	type Example struct {
		Foo string `v8n:"nullable"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.False(t, v.Properties["Foo"].NotNull)
}

func TestTagOptional(t *testing.T) {
	type Example struct {
		Foo string `v8n:"optional"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.False(t, v.Properties["Foo"].Mandatory)
}

func TestTagOrder(t *testing.T) {
	type Example struct {
		Foo string `v8n:"order:0"`
		Bar string `v8n:"order:1"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 0, v.Properties["Foo"].Order)
	require.Equal(t, 1, v.Properties["Bar"].Order)
}

func TestTagRequired(t *testing.T) {
	type Example struct {
		Foo string `v8n:"required"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["Foo"].Mandatory)
}

func TestTagType(t *testing.T) {
	type Example struct {
		Foo json.Number `v8n:"type:integer"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, valix.JsonInteger, v.Properties["Foo"].Type)
}

func TestTagConstraint(t *testing.T) {
	type Example struct {
		Foo string `v8n:"constraint:StringMaxLength{Value:255}"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["Foo"].Constraints))
}

func TestShorthandConstraint(t *testing.T) {
	type Example struct {
		Foo string `v8n:"&StringMaxLength{Value:255}"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["Foo"].Constraints))
}

func TestTagConstraints(t *testing.T) {
	type Example struct {
		Foo string `v8n:"constraints:[StringNotEmpty{},StringNoControlCharacters{}]"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 2, len(v.Properties["Foo"].Constraints))
}

func TestTagWhen(t *testing.T) {
	type Example struct {
		Foo string `v8n:"when:YES_FOO"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["Foo"].WhenConditions))
}

func TestTagWhenMultiple(t *testing.T) {
	type Example struct {
		Foo string `v8n:"when:[YES_FOO,NO_FOO]"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 2, len(v.Properties["Foo"].WhenConditions))
}

func TestTagWhenSingleMultipleTimes(t *testing.T) {
	type Example struct {
		Foo string `v8n:"when:YES_FOO,when:NO_FOO"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 2, len(v.Properties["Foo"].WhenConditions))
}

func TestTagUnwanted(t *testing.T) {
	type Example struct {
		Foo string `v8n:"unwanted:NO_FOO"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["Foo"].UnwantedConditions))
}

func TestTagUnwantedMultiple(t *testing.T) {
	type Example struct {
		Foo string `v8n:"unwanted:[NO_FOO,NO_ANY]"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 2, len(v.Properties["Foo"].UnwantedConditions))
}

func TestTagUnwantedSingleMultipleTimes(t *testing.T) {
	type Example struct {
		Foo string `v8n:"unwanted:NO_FOO,unwanted:NO_ANY"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 2, len(v.Properties["Foo"].UnwantedConditions))
}

func TestObjTagIgnoreUnknownProperties(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string
		} `json:"subObj" v8n:"obj.ignoreUnknownProperties"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["subObj"].ObjectValidator.IgnoreUnknownProperties)
}

func TestObjTagunkUownProperties(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string
		} `json:"subObj" v8n:"obj.unknownProperties:false"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.False(t, v.Properties["subObj"].ObjectValidator.IgnoreUnknownProperties)
}

func TestObjTagConstraint(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string
		} `json:"subObj" v8n:"obj.constraint:Length{Minimum:1,Maximum:16}"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["subObj"].ObjectValidator.Constraints))
}

func TestObjTagWhen(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string
		} `json:"subObj" v8n:"obj.when:YES_SUB"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 1, len(v.Properties["subObj"].ObjectValidator.WhenConditions))
}

func TestObjTagOrdered(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string `v8n:"order:0"`
			Bar string `v8n:"order:1"`
		} `v8n:"obj.ordered"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["SubObj"].ObjectValidator.OrderedPropertyChecks)
}

func TestObjTagOrdered2nd(t *testing.T) {
	type Example struct {
		SubObj struct {
			Foo string `json:"foo"`
			Bar string `json:"bar"`
		} `v8n:"obj.ordered"`
	}
	v, err := valix.ValidatorFor(Example{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["SubObj"].ObjectValidator.OrderedPropertyChecks)
}

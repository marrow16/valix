package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPropertiesRepoInitialised(t *testing.T) {
	require.NotNil(t, propertiesRepo)
	require.Equal(t, 0, len(propertiesRepo.properties))
}

func TestCanRegisterProperties(t *testing.T) {
	defer propertiesRepo.reset()

	require.Equal(t, 0, len(propertiesRepo.properties))

	RegisterProperties(Properties{
		"foo": {
			Type: JsonString,
		},
		"bar": nil,
	})
	require.Equal(t, 2, len(propertiesRepo.properties))
	require.Equal(t, JsonString, propertiesRepo.properties["foo"].Type)
	require.NotNil(t, propertiesRepo.properties["bar"])
	require.Equal(t, JsonAny, propertiesRepo.properties["bar"].Type)
}

func TestPropertiesRepoClears(t *testing.T) {
	propertiesRepo.reset()
	defer propertiesRepo.reset()

	require.Equal(t, 0, len(propertiesRepo.properties))

	RegisterProperties(Properties{
		"foo": {
			Type: JsonString,
		},
	})
	require.Equal(t, 1, len(propertiesRepo.properties))

	PropertiesRepoClear()
	require.Equal(t, 0, len(propertiesRepo.properties))
}

func TestPropertiesRepoPanicsSet(t *testing.T) {
	defer propertiesRepo.reset()

	require.True(t, propertiesRepo.panics)
	PropertiesRepoPanics(false)
	require.False(t, propertiesRepo.panics)
}

func TestPropertiesRepoFetchPanics(t *testing.T) {
	defer func() {
		propertiesRepo.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(propertiesRepoPanicMsg, "foo"), r.(error).Error())
		}
	}()

	ptys := propertiesRepo.fetch(Properties{
		"foo": nil,
	})
	require.Equal(t, 1, len(ptys))
}

func TestPropertiesRepoFetchNotPanics(t *testing.T) {
	defer propertiesRepo.reset()

	PropertiesRepoPanics(false)

	ptys := propertiesRepo.fetch(Properties{
		"foo": nil,
	})
	require.Equal(t, 1, len(ptys))
}

func TestPropertiesRepoFetch(t *testing.T) {
	defer propertiesRepo.reset()

	PropertiesRepoPanics(false)

	barProperty := &PropertyValidator{
		Type: JsonString,
	}
	ptys := propertiesRepo.fetch(Properties{
		"foo": nil,
		"bar": barProperty,
	})
	require.Equal(t, 2, len(ptys))
	require.Equal(t, JsonAny, ptys["foo"].Type)
	require.Equal(t, barProperty, ptys["bar"])

	fooProperty := &PropertyValidator{
		Type: JsonString,
	}
	RegisterProperties(Properties{
		"foo": fooProperty,
	})
	ptys = propertiesRepo.fetch(Properties{
		"foo": nil,
		"bar": barProperty,
	})
	require.Equal(t, 2, len(ptys))
	require.Equal(t, fooProperty, ptys["foo"])
	require.Equal(t, barProperty, ptys["bar"])
}

func TestValidatorNotPanicsWithNilPropertyValidator(t *testing.T) {
	defer propertiesRepo.reset()
	propertiesRepo.panics = false

	validator := &Validator{
		Properties: Properties{
			"foo": nil,
		},
	}
	obj := jsonObject(`{"foo": null}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
}

func TestValidatorPanicsWithNilPropertyValidator(t *testing.T) {
	defer func() {
		propertiesRepo.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(propertiesRepoPanicMsg, "foo"), r.(error).Error())
		}
	}()
	propertiesRepo.panics = true

	validator := &Validator{
		Properties: Properties{
			"foo": nil,
		},
	}
	obj := jsonObject(`{"foo": null}`)
	ok, _ := validator.Validate(obj)
	require.True(t, ok)
}

func TestGetNamed(t *testing.T) {
	propertiesRepo.reset()
	defer propertiesRepo.reset()

	pv := propertiesRepo.getNamed("foo")
	require.Nil(t, pv)

	RegisterProperties(Properties{
		"foo": {
			Type: JsonString,
		},
	})
	pv = propertiesRepo.getNamed("foo")
	require.NotNil(t, pv)
	require.Equal(t, JsonString, pv.Type)
}

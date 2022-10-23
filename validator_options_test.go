package valix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

type test struct{}

func TestValidatorForOptions(t *testing.T) {
	v, err := ValidatorFor(test{}, &ValidatorForOptions{
		IgnoreUnknownProperties: true,
		Constraints:             Constraints{&Length{}},
		AllowNullJson:           true,
		AllowArray:              true,
		DisallowObject:          true,
		StopOnFirst:             true,
		UseNumber:               true,
		OrderedPropertyChecks:   true,
		OasInfo:                 &OasInfo{},
	})
	require.NoError(t, err)
	require.NotNil(t, v)
	require.True(t, v.IgnoreUnknownProperties)
	require.Equal(t, 1, len(v.Constraints))
	require.True(t, v.AllowNullJson)
	require.True(t, v.AllowArray)
	require.True(t, v.DisallowObject)
	require.True(t, v.StopOnFirst)
	require.True(t, v.UseNumber)
	require.True(t, v.OrderedPropertyChecks)
	require.NotNil(t, v.OasInfo)
}

func TestOptionIgnoreUnknownProperties(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionIgnoreUnknownProperties)
	require.NoError(t, err)
	require.True(t, v.IgnoreUnknownProperties)

	v, err = ValidatorFor(test{}, OptionIgnoreUnknownProperties, OptionDisallowUnknownProperties)
	require.NoError(t, err)
	require.False(t, v.IgnoreUnknownProperties)
}

func TestOptionConstraints(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionConstraints())
	require.NoError(t, err)
	require.Equal(t, 0, len(v.Constraints))

	v, err = ValidatorFor(test{}, OptionConstraints(&Length{}))
	require.NoError(t, err)
	require.Equal(t, 1, len(v.Constraints))

	v, err = ValidatorFor(test{}, OptionConstraints(&Length{}, &Length{}))
	require.NoError(t, err)
	require.Equal(t, 2, len(v.Constraints))

	v, err = ValidatorFor(test{}, OptionConstraints(&Length{}), OptionConstraints(&Length{}))
	require.NoError(t, err)
	require.Equal(t, 2, len(v.Constraints))
}

func TestOptionAllowNullJson(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionAllowNullJson)
	require.NoError(t, err)
	require.True(t, v.AllowNullJson)

	v, err = ValidatorFor(test{}, OptionAllowNullJson, OptionDisallowNullJson)
	require.NoError(t, err)
	require.False(t, v.AllowNullJson)
}

func TestOptionAllowArray(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionAllowArray)
	require.NoError(t, err)
	require.True(t, v.AllowArray)

	v, err = ValidatorFor(test{}, OptionAllowArray, OptionDisallowArray)
	require.NoError(t, err)
	require.False(t, v.AllowArray)
}

func TestOptionAllowObject(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionDisallowObject)
	require.NoError(t, err)
	require.True(t, v.DisallowObject)

	v, err = ValidatorFor(test{}, OptionDisallowObject, OptionAllowObject)
	require.NoError(t, err)
	require.False(t, v.DisallowObject)
}

func TestOptionStopOnFirst(t *testing.T) {
	v, err := ValidatorFor(test{}, OptionStopOnFirst)
	require.NoError(t, err)
	require.True(t, v.StopOnFirst)

	v, err = ValidatorFor(test{}, OptionStopOnFirst, OptionDontStopOnFirst)
	require.NoError(t, err)
	require.False(t, v.StopOnFirst)
}

func TestOptionUseNumber(t *testing.T) {
	v, err := ValidatorFor(test{})
	require.NoError(t, err)
	require.False(t, v.UseNumber)

	v, err = ValidatorFor(test{}, OptionUseNumber)
	require.NoError(t, err)
	require.True(t, v.UseNumber)
}

func TestOptionOrderedPropertyChecks(t *testing.T) {
	v, err := ValidatorFor(test{})
	require.NoError(t, err)
	require.False(t, v.OrderedPropertyChecks)

	v, err = ValidatorFor(test{}, OptionOrderedPropertyChecks)
	require.NoError(t, err)
	require.True(t, v.OrderedPropertyChecks)

	v, err = ValidatorFor(test{}, OptionOrderedPropertyChecks, OptionUnOrderedPropertyChecks)
	require.NoError(t, err)
	require.False(t, v.OrderedPropertyChecks)
}

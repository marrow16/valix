package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"math"
	"testing"
)

func TestCoerceToFloat(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expect   bool
		isNumber bool
	}{
		{
			float64(1),
			true,
			true,
		},
		{
			math.NaN(),
			false,
			true,
		},
		{
			math.Inf(0),
			true,
			true,
		},
		{
			math.Inf(-1),
			true,
			true,
		},
		{
			float32(1),
			true,
			true,
		},
		{
			1,
			true,
			true,
		},
		{
			int8(1),
			true,
			true,
		},
		{
			int16(1),
			true,
			true,
		},
		{
			int32(1),
			true,
			true,
		},
		{
			int64(1),
			true,
			true,
		},
		{
			math.MaxInt64,
			true,
			true,
		},
		{
			uint(1),
			true,
			true,
		},
		{
			uint8(1),
			true,
			true,
		},
		{
			uint16(1),
			true,
			true,
		},
		{
			uint32(1),
			true,
			true,
		},
		{
			uint64(1),
			true,
			true,
		},
		{
			json.Number("1"),
			true,
			true,
		},
		{
			json.Number("1.1e-2"),
			true,
			true,
		},
		{
			json.Number("Inf"),
			true,
			true,
		},
		{
			json.Number("NaN"),
			false,
			true,
		},
		{
			json.Number("xxx"),
			false,
			true, // even though it's not a valid number - it's still a number type
		},
		{
			nil,
			false,
			false,
		},
		{
			"str",
			false,
			false,
		},
		{
			false,
			false,
			false,
		},
		{
			[]string{"x"},
			false,
			false,
		},
		{
			map[string]interface{}{"x": nil},
			false,
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]CoerceToFloat:%[2]v(%[2]T)", i+1, tc.v), func(t *testing.T) {
			_, ok, isNumber := coerceToFloat(tc.v)
			require.Equal(t, tc.expect, ok)
			require.Equal(t, tc.isNumber, isNumber)
		})
	}
}

func TestCoerceToInt(t *testing.T) {
	testCases := []struct {
		v        interface{}
		expect   bool
		isNumber bool
	}{
		{
			float64(1),
			true,
			true,
		},
		{
			math.SmallestNonzeroFloat64,
			false,
			true,
		},
		{
			float32(1),
			true,
			true,
		},
		{
			math.SmallestNonzeroFloat32,
			false,
			true,
		},
		{
			math.Inf(0),
			false,
			true,
		},
		{
			math.Inf(-1),
			false,
			true,
		},
		{
			math.NaN(),
			false,
			true,
		},
		{
			1,
			true,
			true,
		},
		{
			int8(1),
			true,
			true,
		},
		{
			int16(1),
			true,
			true,
		},
		{
			int32(1),
			true,
			true,
		},
		{
			int64(1),
			true,
			true,
		},
		{
			uint(1),
			true,
			true,
		},
		{
			uint8(1),
			true,
			true,
		},
		{
			uint16(1),
			true,
			true,
		},
		{
			uint32(1),
			true,
			true,
		},
		{
			uint64(1),
			true,
			true,
		},
		{
			json.Number("1"),
			true,
			true,
		},
		{
			json.Number("1.23e5"),
			true,
			true,
		},
		{
			json.Number("123e-1"),
			false,
			true,
		},
		{
			json.Number("1.1e-2"),
			false,
			true,
		},
		{
			json.Number("Inf"),
			false,
			true,
		},
		{
			json.Number("NaN"),
			false,
			true,
		},
		{
			json.Number("xxx"),
			false,
			true, // even though it's not a valid number - it's still a number type
		},
		{
			nil,
			false,
			false,
		},
		{
			"str",
			false,
			false,
		},
		{
			false,
			false,
			false,
		},
		{
			[]string{"x"},
			false,
			false,
		},
		{
			map[string]interface{}{"x": nil},
			false,
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("[%d]CoerceToInt:%[2]v(%[2]T)", i+1, tc.v), func(t *testing.T) {
			_, ok, isNumber := coerceToInt(tc.v)
			require.Equal(t, tc.expect, ok)
			require.Equal(t, tc.isNumber, isNumber)
		})
	}
}

func TestIfTernary(t *testing.T) {
	require.Equal(t, "yes", ternary(true).string("yes", "no"))
	require.Equal(t, "no", ternary(false).string("yes", "no"))
	require.Equal(t, 1, ternary(true).int(1, 2))
	require.Equal(t, 2, ternary(false).int(1, 2))
}

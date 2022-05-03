package valix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConditions_Clone(t *testing.T) {
	src := Conditions{"foo"}
	dst := src.Clone()
	require.Equal(t, len(src), len(dst))

	src[0] = "x"
	require.NotEqual(t, src[0], dst[0])

	src = append(src, "bar")
	require.NotEqual(t, len(src), len(dst))

	src = nil
	dst = src.Clone()
	require.Nil(t, dst)
}

func TestConstraints_Clone(t *testing.T) {
	constraint := &StringNotEmpty{}
	src := Constraints{constraint, nil}
	dst := src.Clone()
	require.Equal(t, len(src), len(dst))

	// the actual constraint items are NOT cloned...
	require.Equal(t, src[0], dst[0])
	constraint.Message = "Fooey"
	require.Equal(t, "Fooey", src[0].(*StringNotEmpty).Message)
	require.Equal(t, "Fooey", dst[0].(*StringNotEmpty).Message)

	src = nil
	dst = src.Clone()
	require.Nil(t, dst)
}

func TestOasInfo_Clone(t *testing.T) {
	src := &OasInfo{Title: "foo"}
	dst := src.Clone()
	require.Equal(t, src, dst)
	require.Equal(t, src.Title, dst.Title)
	dst.Title = "bar"
	require.NotEqual(t, src.Title, dst.Title)

	dst2 := cloneOasInfo(src)
	require.NotNil(t, dst2)
	require.NotEqual(t, dst, dst2)

	src = nil
	dst = cloneOasInfo(src)
	require.Nil(t, dst)
}

func TestOthersExpr_Clone(t *testing.T) {
	item := &OtherProperty{Name: "foo"}
	src := OthersExpr{item}
	dst := src.Clone()
	require.Equal(t, len(src), len(dst))

	// the actual items are NOT cloned...
	require.Equal(t, src[0], dst[0])

	src = nil
	dst = src.Clone()
	require.Nil(t, dst)
}

func TestConditionalVariant_Clone(t *testing.T) {
	src := &ConditionalVariant{
		WhenConditions: Conditions{"foo"},
	}
	dst := src.Clone()
	require.False(t, src == dst)

	require.Equal(t, len(src.WhenConditions), len(dst.WhenConditions))
	src.WhenConditions = append(src.WhenConditions, "bar")
	require.NotEqual(t, len(src.WhenConditions), len(dst.WhenConditions))

	dst2 := cloneConditionalVariant(src)
	require.NotNil(t, dst2)
	require.False(t, dst == dst2)

	src = nil
	dst = cloneConditionalVariant(src)
	require.Nil(t, dst)
}

func TestConditionalVariants_Clone(t *testing.T) {
	variant := &ConditionalVariant{
		WhenConditions: Conditions{"foo"},
	}
	src := ConditionalVariants{variant}
	dst := src.Clone()
	require.Equal(t, len(src), len(dst))
	require.Equal(t, len(src[0].WhenConditions), len(dst[0].WhenConditions))
	src[0].WhenConditions = append(src[0].WhenConditions, "bar")
	require.NotEqual(t, len(src[0].WhenConditions), len(dst[0].WhenConditions))

	require.Equal(t, src[0].WhenConditions[0], dst[0].WhenConditions[0])
	src[0].WhenConditions[0] = "bar"
	require.NotEqual(t, src[0].WhenConditions[0], dst[0].WhenConditions[0])

	src = nil
	dst = src.Clone()
	require.Nil(t, dst)
}

func TestProperties_Clone(t *testing.T) {
	src := Properties{
		"foo": {
			Type: JsonString,
		},
		"bar": nil,
	}
	dst := src.Clone()
	require.Equal(t, len(src), len(dst))
	src["baz"] = nil
	require.NotEqual(t, len(src), len(dst))

	dst["foo"].Type = JsonAny
	require.Equal(t, JsonString, src["foo"].Type)

	src = nil
	dst = src.Clone()
	require.Nil(t, dst)
}

func TestPropertyValidator_Clone(t *testing.T) {
	src := &PropertyValidator{
		Mandatory:      true,
		WhenConditions: Conditions{"foo"},
	}
	dst := src.Clone()
	require.False(t, src == dst)
	require.Equal(t, len(src.WhenConditions), len(dst.WhenConditions))
	require.Equal(t, src.Mandatory, dst.Mandatory)
	src.WhenConditions = append(src.WhenConditions, "bar")
	src.Mandatory = false
	require.NotEqual(t, len(src.WhenConditions), len(dst.WhenConditions))
	require.NotEqual(t, src.Mandatory, dst.Mandatory)
}

func TestValidator_Clone(t *testing.T) {
	src := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": nil,
		},
	}
	dst := src.Clone()
	require.False(t, src == dst)

	require.Equal(t, src.IgnoreUnknownProperties, dst.IgnoreUnknownProperties)
	require.Equal(t, len(src.Properties), len(dst.Properties))
	src.IgnoreUnknownProperties = !src.IgnoreUnknownProperties
	src.Properties["bar"] = nil
	require.NotEqual(t, src.IgnoreUnknownProperties, dst.IgnoreUnknownProperties)
	require.NotEqual(t, len(src.Properties), len(dst.Properties))

	dst2 := cloneValidator(src)
	require.NotNil(t, dst2)
	require.False(t, dst == dst2)

	src = nil
	dst = cloneValidator(src)
	require.Nil(t, dst)
}

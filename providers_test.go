package valix

import (
	"encoding/json"
	"github.com/stretchr/testify/require"
	"reflect"
	"strings"
	"testing"
)

func TestDecoderProvider(t *testing.T) {
	defer func() {
		DefaultDecoderProvider = &defaultDecoderProvider{}
	}()
	jStr := `{
		"foo": 1
	}`

	d := getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), nil)
	obj := map[string]interface{}{}
	err := d.Decode(&obj)
	require.NoError(t, err)
	require.Equal(t, float64(1), obj["foo"])

	d = getDefaultDecoderProvider().NewDecoder(strings.NewReader(jStr), true)
	obj = map[string]interface{}{}
	err = d.Decode(&obj)
	require.NoError(t, err)
	jn, ok := obj["foo"].(json.Number)
	require.True(t, ok)
	require.Equal(t, "1", jn.String())

	v := &Validator{
		UseNumber: false,
	}
	d = getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), v)
	obj = map[string]interface{}{}
	err = d.Decode(&obj)
	require.NoError(t, err)
	require.Equal(t, float64(1), obj["foo"])

	v.UseNumber = true
	d = getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), v)
	obj = map[string]interface{}{}
	err = d.Decode(&obj)
	require.NoError(t, err)
	jn, ok = obj["foo"].(json.Number)
	require.True(t, ok)
	require.Equal(t, "1", jn.String())

	v.IgnoreUnknownProperties = false
	d = getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), v)
	type unknownPtyTest struct {
		Bar string
	}
	ts := &unknownPtyTest{}
	err = d.Decode(ts)
	require.Error(t, err)
	require.Equal(t, "json: unknown field \"foo\"", err.Error())

	v.IgnoreUnknownProperties = true
	d = getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), v)
	ts = &unknownPtyTest{}
	err = d.Decode(ts)
	require.NoError(t, err)

	// and nil replacement protection...
	DefaultDecoderProvider = nil
	d = getDefaultDecoderProvider().NewDecoderFor(strings.NewReader(jStr), v)
	ts = &unknownPtyTest{}
	err = d.Decode(ts)
	require.NoError(t, err)
}

func TestDefaultPropertyNameProvider(t *testing.T) {
	defer func() {
		DefaultPropertyNameProvider = &defaultPropertyNameProvider{}
	}()
	type namesStruct struct {
		Foo string `json:"foo,omitempty"`
	}
	ns := namesStruct{}
	ty := reflect.TypeOf(ns)
	fld := ty.Field(0)
	name, ok := getDefaultPropertyNameProvider().NameFor(fld)
	require.True(t, ok)
	require.Equal(t, "foo", name)

	type namesStruct2 struct {
		Foo string
	}
	ns2 := namesStruct2{}
	ty = reflect.TypeOf(ns2)
	fld = ty.Field(0)
	name, ok = getDefaultPropertyNameProvider().NameFor(fld)
	require.True(t, ok)
	require.Equal(t, "Foo", name)

	// try replacement...
	DefaultPropertyNameProvider = &testAutoLowerFirstLetter{}
	name, ok = getDefaultPropertyNameProvider().NameFor(fld)
	require.True(t, ok)
	require.Equal(t, "foo", name)

	// and nil replacement protection...
	DefaultPropertyNameProvider = nil
	name, ok = getDefaultPropertyNameProvider().NameFor(fld)
	require.True(t, ok)
	require.Equal(t, "Foo", name)
}

type testAutoLowerFirstLetter struct{}

func (dpnp *testAutoLowerFirstLetter) NameFor(field reflect.StructField) (name string, ok bool) {
	result := field.Name
	if tag, ok := field.Tag.Lookup(tagNameJson); ok {
		if cAt := strings.Index(tag, ","); cAt != -1 {
			result = tag[0:cAt]
		} else {
			result = tag
		}
	} else {
		result = strings.ToLower(result[:1]) + result[1:]
	}
	return result, true
}

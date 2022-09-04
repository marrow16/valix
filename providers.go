package valix

import (
	"encoding/json"
	"io"
	"reflect"
	"strings"
)

// DefaultDecoderProvider is the decoder provider used by Validator - replace with your own if necessary
var DefaultDecoderProvider DecoderProvider = &defaultDecoderProvider{}

func getDefaultDecoderProvider() DecoderProvider {
	if DefaultDecoderProvider != nil {
		return DefaultDecoderProvider
	}
	return &defaultDecoderProvider{}
}

// DecoderProvider is the interface needed for replacing the DefaultDecoderProvider
type DecoderProvider interface {
	NewDecoder(r io.Reader, useNumber bool) *json.Decoder
	NewDecoderFor(r io.Reader, validator *Validator) *json.Decoder
}

type defaultDecoderProvider struct{}

func (ddp *defaultDecoderProvider) NewDecoder(r io.Reader, useNumber bool) *json.Decoder {
	d := json.NewDecoder(r)
	if useNumber {
		d.UseNumber()
	}
	return d
}

func (ddp *defaultDecoderProvider) NewDecoderFor(r io.Reader, validator *Validator) *json.Decoder {
	if validator == nil {
		return ddp.NewDecoder(r, false)
	}
	d := ddp.NewDecoder(r, validator.UseNumber)
	if !validator.IgnoreUnknownProperties {
		d.DisallowUnknownFields()
	}
	return d
}

// DefaultPropertyNameProvider is the property name provider used by Validator - replace with your own if necessary
var DefaultPropertyNameProvider PropertyNameProvider = &defaultPropertyNameProvider{}

func getDefaultPropertyNameProvider() PropertyNameProvider {
	if DefaultPropertyNameProvider != nil {
		return DefaultPropertyNameProvider
	}
	return &defaultPropertyNameProvider{}
}

type PropertyNameProvider interface {
	// NameFor provides the property name for a given struct field (return ok false if default field name is to be used)
	NameFor(field reflect.StructField) (name string, ok bool)
}

type defaultPropertyNameProvider struct{}

func (dpnp *defaultPropertyNameProvider) NameFor(field reflect.StructField) (name string, ok bool) {
	result := field.Name
	if tag, ok := field.Tag.Lookup(tagNameJson); ok {
		if cAt := strings.Index(tag, ","); cAt != -1 {
			result = tag[0:cAt]
		} else {
			result = tag
		}
	}
	return result, true
}

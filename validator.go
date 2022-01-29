package valix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

const (
	MessageRequestBodyEmpty              = "Request body is empty"
	MessageUnableToDecode                = "Unable to decode request body as JSON"
	MessageRequestBodyNotJsonNull        = "Request body must not be JSON null"
	MessageRequestBodyNotJsonArray       = "Request body must not be a JSON array"
	MessageRequestBodyNotJsonObject      = "Request body must not be a JSON object"
	MessageRequestBodyExpectedJsonArray  = "Request body expected to be JSON array"
	MessageRequestBodyExpectedJsonObject = "Request body must be a JSON object"
	MessageArrayElementMustBeObject      = "Array element [%d] must be an object"
	MessageMissingProperty               = "Missing property '%s'"
	MessageUnknownProperty               = "Unknown property '%s'"
)

type Validator struct {
	// IgnoreUnknownProperties is whether to ignore unknown properties (default false)
	//
	// Set this to `true` if you want to allow unknown properties
	IgnoreUnknownProperties bool
	// Properties is the map of property names (key) and PropertyValidator (value)
	Properties map[string]*PropertyValidator
	// Constraints is an optional slice of Constraint items to be checked on the object/array
	//
	// * These are checked in the order specified and prior to property validator & unknown property checks
	Constraints    []Constraint
	AllowArray     bool
	DisallowObject bool
	AllowNull      bool
	// UseNumber forces RequestValidate method to use json.Number when decoding request body
	UseNumber bool
}

// RequestValidate Performs validation on the request body of the supplied http.Request
func (v *Validator) RequestValidate(r *http.Request) (bool, []*Violation, interface{}) {
	ctx := newContext(r, v)
	var obj interface{} = nil
	if r.Body == nil {
		ctx.AddViolation(NewEmptyViolation(MessageRequestBodyEmpty))
	} else {
		decoder := json.NewDecoder(r.Body)
		if v.UseNumber {
			decoder.UseNumber()
		}
		obj = reflect.Interface
		if err := decoder.Decode(&obj); err != nil {
			obj = nil
			ctx.AddViolation(NewEmptyViolation(MessageUnableToDecode))
		} else if obj == nil {
			if !v.AllowNull {
				ctx.AddViolation(NewEmptyViolation(MessageRequestBodyNotJsonNull))
			}
		} else {
			// determine whether body is a map (object) or an array...
			if arr, isArr := obj.([]interface{}); isArr {
				if v.AllowArray {
					v.validateArrayOf(arr, ctx)
				} else {
					ctx.AddViolation(NewEmptyViolation(MessageRequestBodyNotJsonArray))
				}
			} else if m, isMap := obj.(map[string]interface{}); isMap {
				if v.DisallowObject {
					if v.AllowArray {
						ctx.AddViolation(NewEmptyViolation(MessageRequestBodyExpectedJsonArray))
					} else {
						ctx.AddViolation(NewEmptyViolation(MessageRequestBodyNotJsonObject))
					}
				} else {
					v.validate(m, ctx)
				}
			} else {
				ctx.AddViolation(NewEmptyViolation(MessageRequestBodyExpectedJsonObject))
			}
		}
	}
	return ctx.ok, ctx.violations, obj
}

// Validate Performs validation on the supplied object
func (v *Validator) Validate(obj map[string]interface{}) (bool, []*Violation) {
	ctx := newContext(obj, v)
	v.validate(obj, ctx)
	return ctx.ok, ctx.violations
}

// ValidateArrayOf Performs validation on each element of the supplied array
func (v *Validator) ValidateArrayOf(arr []interface{}) (bool, []*Violation) {
	ctx := newContext(arr, v)
	v.validateArrayOf(arr, ctx)
	return ctx.ok, ctx.violations
}

func (v *Validator) validate(obj map[string]interface{}, ctx *Context) {
	if v.Constraints != nil {
		for _, constraint := range v.Constraints {
			if ok, msg := constraint.Validate(obj, ctx); !ok {
				ctx.AddViolationForCurrent(msg)
			}
			if !ctx.continueAll {
				return
			}
		}
	}
	v.checkUnknownProperties(obj, ctx)
	v.checkProperties(obj, ctx)
}

func (v *Validator) validateArrayOf(arr []interface{}, ctx *Context) {
	for i, elem := range arr {
		ctx.pushPathIndex(i, elem)
		if obj, itemOk := elem.(map[string]interface{}); itemOk {
			v.validate(obj, ctx)
			ctx.popPath()
		} else {
			ctx.popPath()
			ctx.AddViolationForCurrent(fmt.Sprintf(MessageArrayElementMustBeObject, i))
		}
		if !ctx.continueAll {
			return
		}
	}
}

func (v *Validator) checkProperties(obj map[string]interface{}, ctx *Context) {
	for propertyName, pv := range v.Properties {
		if actualValue, pOk := obj[propertyName]; !pOk {
			if pv.Mandatory {
				ctx.AddViolationForCurrent(fmt.Sprintf(MessageMissingProperty, propertyName))
			}
		} else {
			ctx.pushPathProperty(propertyName, actualValue)
			pv.validate(actualValue, ctx)
			ctx.popPath()
		}
		if !ctx.continueAll {
			return
		}
	}
}

func (v *Validator) checkUnknownProperties(obj map[string]interface{}, ctx *Context) {
	if !v.IgnoreUnknownProperties {
		for propertyName := range obj {
			if _, hasK := v.Properties[propertyName]; !hasK {
				ctx.AddViolationForCurrent(fmt.Sprintf(MessageUnknownProperty, propertyName))
			}
		}
	}
}

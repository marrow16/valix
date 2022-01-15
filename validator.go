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
	MessageRequestBodyNotJsonArray       = "Request body must not be JSON array"
	MessageRequestBodyNotJsonObject      = "Request body must not be JSON object"
	MessageRequestBodyExpectedJsonArray  = "Request body expected to be JSON array"
	MessageRequestBodyExpectedJsonObject = "Request body must be a JSON object"
	MessageArrayElementMustBeObject      = "Array element [%d] must be an object"
	MessageMissingProperty               = "Missing property '%s'"
	MessageUnknownProperty               = "Unknown property '%s'"
)

type Validator struct {
	IgnoreUnknownProperties bool
	Properties              map[string]*PropertyValidator
	Constraints             []Constraint
	AllowArray              bool
	DisallowObject          bool
	AllowNull               bool
}

// RequestValidate Performs validation on the request body of the supplied http.Request
func (v *Validator) RequestValidate(r *http.Request) (bool, []*Violation, interface{}) {
	ctx := newContext(r, v)
	var obj interface{} = nil
	if r.Body == nil {
		ctx.AddViolation(newEmptyViolation(MessageRequestBodyEmpty))
	} else {
		decoder := json.NewDecoder(r.Body)
		obj = reflect.Interface
		if err := decoder.Decode(&obj); err != nil {
			obj = nil
			ctx.AddViolation(newEmptyViolation(MessageUnableToDecode))
		} else if obj == nil {
			if !v.AllowNull {
				ctx.AddViolation(newEmptyViolation(MessageRequestBodyNotJsonNull))
			}
		} else {
			// determine whether body is a map (object) or an array...
			if arr, isArr := obj.([]interface{}); isArr {
				if v.AllowArray {
					v.validateArrayOf(arr, ctx)
				} else {
					ctx.AddViolation(newEmptyViolation(MessageRequestBodyNotJsonArray))
				}
			} else if m, isMap := obj.(map[string]interface{}); isMap {
				if v.DisallowObject {
					if v.AllowArray {
						ctx.AddViolation(newEmptyViolation(MessageRequestBodyExpectedJsonArray))
					} else {
						ctx.AddViolation(newEmptyViolation(MessageRequestBodyNotJsonObject))
					}
				} else {
					v.validate(m, ctx)
				}
			} else {
				ctx.AddViolation(newEmptyViolation(MessageRequestBodyExpectedJsonObject))
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
			if !ctx.continues {
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
		if !ctx.continues {
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
		if !ctx.continues {
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

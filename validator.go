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

type Properties map[string]*PropertyValidator
type Constraints []Constraint

type Validator struct {
	// IgnoreUnknownProperties is whether to ignore unknown properties (default false)
	//
	// Set this to `true` if you want to allow unknown properties
	IgnoreUnknownProperties bool
	// Properties is the map of property names (key) and PropertyValidator (value)
	Properties Properties
	// Constraints is an optional slice of Constraint items to be checked on the object/array
	//
	// * These are checked in the order specified and prior to property validator & unknown property checks
	Constraints    Constraints
	AllowArray     bool
	DisallowObject bool
	AllowNull      bool
	// UseNumber forces RequestValidate method to use json.Number when decoding request body
	UseNumber bool
}

// RequestValidate Performs validation on the request body of the supplied http.Request
func (v *Validator) RequestValidate(r *http.Request) (bool, []*Violation, interface{}) {
	vcx := newValidatorContext(r)
	var obj interface{} = nil
	if r.Body == nil {
		vcx.AddViolation(NewBadRequestViolation(MessageRequestBodyEmpty))
	} else {
		decoder := json.NewDecoder(r.Body)
		if v.UseNumber {
			decoder.UseNumber()
		}
		obj = reflect.Interface
		if err := decoder.Decode(&obj); err != nil {
			obj = nil
			vcx.AddViolation(NewBadRequestViolation(MessageUnableToDecode))
		} else if obj == nil {
			if !v.AllowNull {
				vcx.AddViolation(NewBadRequestViolation(MessageRequestBodyNotJsonNull))
			}
		} else {
			// determine whether body is a map (object) or an array...
			if arr, isArr := obj.([]interface{}); isArr {
				if v.AllowArray {
					v.validateArrayOf(arr, vcx)
				} else {
					vcx.AddViolation(NewEmptyViolation(MessageRequestBodyNotJsonArray))
				}
			} else if m, isMap := obj.(map[string]interface{}); isMap {
				if v.DisallowObject {
					if v.AllowArray {
						vcx.AddViolation(NewEmptyViolation(MessageRequestBodyExpectedJsonArray))
					} else {
						vcx.AddViolation(NewEmptyViolation(MessageRequestBodyNotJsonObject))
					}
				} else {
					v.validate(m, vcx)
				}
			} else {
				vcx.AddViolation(NewBadRequestViolation(MessageRequestBodyExpectedJsonObject))
			}
		}
	}
	return vcx.ok, vcx.violations, obj
}

// Validate Performs validation on the supplied object
func (v *Validator) Validate(obj map[string]interface{}) (bool, []*Violation) {
	vcx := newValidatorContext(obj)
	v.validate(obj, vcx)
	return vcx.ok, vcx.violations
}

// ValidateArrayOf Performs validation on each element of the supplied array
func (v *Validator) ValidateArrayOf(arr []interface{}) (bool, []*Violation) {
	vcx := newValidatorContext(arr)
	v.validateArrayOf(arr, vcx)
	return vcx.ok, vcx.violations
}

func (v *Validator) validate(obj map[string]interface{}, vcx *ValidatorContext) {
	if v.Constraints != nil {
		for _, constraint := range v.Constraints {
			if ok, msg := constraint.Validate(obj, vcx); !ok {
				vcx.AddViolationForCurrent(msg)
			}
			if !vcx.continueAll {
				return
			}
		}
	}
	v.checkUnknownProperties(obj, vcx)
	v.checkProperties(obj, vcx)
}

func (v *Validator) validateArrayOf(arr []interface{}, vcx *ValidatorContext) {
	for i, elem := range arr {
		vcx.pushPathIndex(i, elem)
		if obj, itemOk := elem.(map[string]interface{}); itemOk {
			v.validate(obj, vcx)
			vcx.popPath()
		} else {
			vcx.popPath()
			vcx.AddViolationForCurrent(fmt.Sprintf(MessageArrayElementMustBeObject, i))
		}
		if !vcx.continueAll {
			return
		}
	}
}

func (v *Validator) checkProperties(obj map[string]interface{}, vcx *ValidatorContext) {
	for propertyName, pv := range v.Properties {
		if actualValue, pOk := obj[propertyName]; !pOk {
			if pv.Mandatory {
				vcx.AddViolationForCurrent(fmt.Sprintf(MessageMissingProperty, propertyName))
			}
		} else {
			vcx.pushPathProperty(propertyName, actualValue)
			pv.validate(actualValue, vcx)
			vcx.popPath()
		}
		if !vcx.continueAll {
			return
		}
	}
}

func (v *Validator) checkUnknownProperties(obj map[string]interface{}, vcx *ValidatorContext) {
	if !v.IgnoreUnknownProperties {
		for propertyName := range obj {
			if _, hasK := v.Properties[propertyName]; !hasK {
				vcx.AddViolationForCurrent(fmt.Sprintf(MessageUnknownProperty, propertyName))
			}
		}
	}
}

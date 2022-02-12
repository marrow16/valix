package valix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
)

const (
	messageRequestBodyEmpty              = "Request body is empty"
	messageUnableToDecode                = "Unable to decode request body as JSON"
	messageRequestBodyNotJsonNull        = "Request body must not be JSON null"
	messageRequestBodyNotJsonArray       = "Request body must not be a JSON array"
	messageRequestBodyNotJsonObject      = "Request body must not be a JSON object"
	messageRequestBodyExpectedJsonArray  = "Request body expected to be JSON array"
	messageRequestBodyExpectedJsonObject = "Request body must be a JSON object"
	messageArrayElementMustBeObject      = "JsonArray element [%d] must be an object"
	messageMissingProperty               = "Missing property '%s'"
	messageUnknownProperty               = "Unknown property '%s'"
)

// Properties type used by Validator.Properties
type Properties map[string]*PropertyValidator

// Constraints type used by Validator.Constraints and PropertyValidator.Constraints
type Constraints []Constraint

// Validator is the validator against which requests, maps and slices can be checked
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
	Constraints Constraints
	// AllowArray denotes, when true (default is false), that this validator will allow a JSON array - where each
	// item in the array can be validated as an object
	AllowArray bool
	// DisallowObject denotes, when set to true, that this validator will disallow JSON objects - i.e. that it
	// expects JSON arrays (in which case the AllowArray should also be set to true)
	DisallowObject bool
	// AllowNullJson forces RequestValidate to accept a request body that is null JSON (i.e. a body containing just `null`)
	AllowNullJson bool
	// UseNumber forces RequestValidate method to use json.Number when decoding request body
	UseNumber bool
}

// RequestValidate Performs validation on the request body of the supplied http.Request
func (v *Validator) RequestValidate(r *http.Request) (bool, []*Violation, interface{}) {
	vcx := newValidatorContext(r)
	ok, obj := v.decodeRequestBody(r, vcx)
	if ok {
		v.requestBodyValidate(vcx, obj)
	}
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) decodeRequestBody(r *http.Request, vcx *ValidatorContext) (bool, interface{}) {
	if r.Body == nil {
		vcx.AddViolation(NewBadRequestViolation(messageRequestBodyEmpty))
		return false, nil
	}
	decoder := json.NewDecoder(r.Body)
	if v.UseNumber {
		decoder.UseNumber()
	}
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx.AddViolation(NewBadRequestViolation(messageUnableToDecode))
		return false, nil
	}
	return true, obj
}

func (v *Validator) requestBodyValidate(vcx *ValidatorContext, obj interface{}) {
	if obj != nil {
		// determine whether body is a map (object) or an array...
		if arr, isArr := obj.([]interface{}); isArr {
			if v.AllowArray {
				v.validateArrayOf(arr, vcx)
			} else {
				vcx.AddViolation(NewEmptyViolation(messageRequestBodyNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject {
				if v.AllowArray {
					vcx.AddViolation(NewEmptyViolation(messageRequestBodyExpectedJsonArray))
				} else {
					vcx.AddViolation(NewEmptyViolation(messageRequestBodyNotJsonObject))
				}
			} else {
				v.validate(m, vcx)
			}
		} else {
			vcx.AddViolation(NewBadRequestViolation(messageRequestBodyExpectedJsonObject))
		}
	} else if !v.AllowNullJson {
		vcx.AddViolation(NewBadRequestViolation(messageRequestBodyNotJsonNull))
	}
}

// Validate performs validation on the supplied object
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
			if ok, msg := constraint.Check(obj, vcx); !ok {
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
			vcx.AddViolationForCurrent(fmt.Sprintf(messageArrayElementMustBeObject, i))
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
				vcx.AddViolationForCurrent(fmt.Sprintf(messageMissingProperty, propertyName))
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
				vcx.AddViolationForCurrent(fmt.Sprintf(messageUnknownProperty, propertyName))
			}
		}
	}
}

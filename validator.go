package valix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const (
	msgUnableToDecode     = "Unable to decode as JSON"
	msgNotJsonNull        = "JSON must not be JSON null"
	msgNotJsonArray       = "JSON must not be JSON array"
	msgNotJsonObject      = "JSON must not be JSON object"
	msgExpectedJsonArray  = "JSON expected to be JSON array"
	msgExpectedJsonObject = "JSON expected to be JSON object"
	msgErrorReading       = "Unexpected error reading reader"
	msgErrorUnmarshall    = "Unexpected error during unmarshalling"

	msgRequestBodyEmpty              = "Request body is empty"
	msgUnableToDecodeRequest         = "Unable to decode request body as JSON"
	msgRequestBodyNotJsonNull        = "Request body must not be JSON null"
	msgRequestBodyNotJsonArray       = "Request body must not be JSON array"
	msgRequestBodyNotJsonObject      = "Request body must not be JSON object"
	msgRequestBodyExpectedJsonArray  = "Request body expected to be JSON array"
	msgRequestBodyExpectedJsonObject = "Request body expected to be JSON object"
	msgArrayElementMustBeObject      = "JsonArray element [%d] must be an object"
	msgMissingProperty               = "Missing property '%s'"
	msgUnknownProperty               = "Unknown property '%s'"
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
//
// If the validation of the request body fails, false is returned and the returned violations
// give the reason(s) for the validation failure.
//
// If the validation is successful, the validated JSON (object or array) is also returned -
// as represented by (if the body was a JSON object)
//   map[string]interface{}
// or as represented by (if the body was a JSON array)
//   []interface{}
func (v *Validator) RequestValidate(r *http.Request) (bool, []*Violation, interface{}) {
	tmpVcx := newValidatorContext(nil)
	ok, obj := v.decodeRequestBody(r.Body, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj)
	v.requestBodyValidate(vcx, obj)
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) decodeRequestBody(r io.Reader, vcx *ValidatorContext) (bool, interface{}) {
	if r == nil {
		vcx.AddViolation(NewBadRequestViolation(msgRequestBodyEmpty))
		return false, nil
	}
	decoder := v.createDecoder(r)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx.AddViolation(NewBadRequestViolation(msgUnableToDecodeRequest))
		return false, nil
	}
	return true, obj
}

func (v *Validator) createDecoder(r io.Reader) *json.Decoder {
	decoder := json.NewDecoder(r)
	if v.UseNumber {
		decoder.UseNumber()
	}
	return decoder
}

func (v *Validator) requestBodyValidate(vcx *ValidatorContext, obj interface{}) {
	if obj != nil {
		// determine whether body is a map (object) or a slice (array)...
		if arr, isArr := obj.([]interface{}); isArr {
			if v.AllowArray {
				v.validateArrayOf(arr, vcx)
			} else {
				vcx.AddViolation(NewEmptyViolation(msgRequestBodyNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject && v.AllowArray {
				vcx.AddViolation(NewEmptyViolation(msgRequestBodyExpectedJsonArray))
			} else if v.DisallowObject {
				vcx.AddViolation(NewEmptyViolation(msgRequestBodyNotJsonObject))
			} else {
				v.validate(m, vcx)
			}
		} else {
			vcx.AddViolation(NewBadRequestViolation(msgRequestBodyExpectedJsonObject))
		}
	} else if !v.AllowNullJson {
		vcx.AddViolation(NewBadRequestViolation(msgRequestBodyNotJsonNull))
	}
}

// Validate performs validation on the supplied JSON object
//
// Where the JSON object is represented as an unmarshalled
//   map[string]interface{}
func (v *Validator) Validate(obj map[string]interface{}) (bool, []*Violation) {
	vcx := newValidatorContext(obj)
	v.validate(obj, vcx)
	return vcx.ok, vcx.violations
}

// ValidateArrayOf Performs validation on each element of the supplied JSON array
//
// Where the JSON array is represented as an unmarshalled
//   []interface{}
// and each item of the slice is expected to be a JSON object represented as an unmarshalled
//   map[string]interface{}
func (v *Validator) ValidateArrayOf(arr []interface{}) (bool, []*Violation) {
	vcx := newValidatorContext(arr)
	v.validateArrayOf(arr, vcx)
	return vcx.ok, vcx.violations
}

// ValidateReader performs validation on the supplied reader (representing JSON)
func (v *Validator) ValidateReader(r io.Reader) (bool, []*Violation, interface{}) {
	decoder := v.createDecoder(r)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx := newValidatorContext(nil)
		vcx.AddViolation(NewBadRequestViolation(msgUnableToDecode))
		return vcx.ok, vcx.violations, nil
	}
	vcx := newValidatorContext(obj)
	v.validateObjectOrArray(vcx, obj)
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) validateObjectOrArray(vcx *ValidatorContext, obj interface{}) {
	if obj != nil {
		// determine whether obj is a map (object) or a slice (array)...
		if arr, isArr := obj.([]interface{}); isArr {
			if v.AllowArray {
				v.validateArrayOf(arr, vcx)
			} else {
				vcx.AddViolation(NewEmptyViolation(msgNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject && v.AllowArray {
				vcx.AddViolation(NewEmptyViolation(msgExpectedJsonArray))
			} else if v.DisallowObject {
				vcx.AddViolation(NewEmptyViolation(msgNotJsonObject))
			} else {
				v.validate(m, vcx)
			}
		} else {
			vcx.AddViolation(NewBadRequestViolation(msgExpectedJsonObject))
		}
	} else if !v.AllowNullJson {
		vcx.AddViolation(NewBadRequestViolation(msgNotJsonNull))
	}
}

// ValidateReaderInto performs validation on the supplied reader (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) ValidateReaderInto(r io.Reader, value interface{}) (bool, []*Violation, interface{}) {
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		errVcx := newValidatorContext(nil)
		errVcx.AddViolation(NewBadRequestViolation(msgErrorReading))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	decoder := v.createDecoder(initialReader)
	var obj interface{} = reflect.Interface
	if dErr := decoder.Decode(&obj); dErr != nil {
		errVcx := newValidatorContext(nil)
		errVcx.AddViolation(NewBadRequestViolation(msgUnableToDecode))
		return false, errVcx.violations, nil
	}
	vcx := newValidatorContext(obj)
	v.requestBodyValidate(vcx, obj)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder = v.createDecoder(intoReader)
	if !v.IgnoreUnknownProperties {
		decoder.DisallowUnknownFields()
	}
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(NewBadRequestViolation(msgErrorUnmarshall))
	}
	return vcx.ok, vcx.violations, obj
}

// ValidateString performs validation on the supplied string (representing JSON)
func (v *Validator) ValidateString(s string) (bool, []*Violation, interface{}) {
	return v.ValidateReader(strings.NewReader(s))
}

// ValidateStringInto performs validation on the supplied string (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) ValidateStringInto(s string, value interface{}) (bool, []*Violation, interface{}) {
	return v.ValidateReaderInto(strings.NewReader(s), value)
}

// RequestValidateInto performs validation on the request body (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) RequestValidateInto(req *http.Request, value interface{}) (bool, []*Violation, interface{}) {
	if req.Body == nil {
		errVcx := newValidatorContext(nil)
		errVcx.AddViolation(NewBadRequestViolation(msgRequestBodyEmpty))
		return false, errVcx.violations, nil
	}
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errVcx := newValidatorContext(nil)
		errVcx.AddViolation(NewBadRequestViolation(msgErrorReading))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	tmpVcx := newValidatorContext(nil)
	ok, obj := v.decodeRequestBody(initialReader, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj)
	v.requestBodyValidate(vcx, obj)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder := v.createDecoder(intoReader)
	if !v.IgnoreUnknownProperties {
		decoder.DisallowUnknownFields()
	}
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(NewBadRequestViolation(msgErrorUnmarshall))
	}
	return vcx.ok, vcx.violations, obj
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
			vcx.AddViolationForCurrent(fmt.Sprintf(msgArrayElementMustBeObject, i))
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
				vcx.AddViolationForCurrent(fmt.Sprintf(msgMissingProperty, propertyName))
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
				vcx.AddViolationForCurrent(fmt.Sprintf(msgUnknownProperty, propertyName))
			}
		}
	}
}

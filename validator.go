package valix

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"sort"
	"strings"
)

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
	// StopOnFirst if set, instructs the validator to stop at the first violation found
	StopOnFirst bool
	// UseNumber forces RequestValidate method to use json.Number when decoding request body
	UseNumber bool
	// OrderedPropertyChecks determines whether properties should be checked in order - when set to true, properties
	// are sorted by PropertyValidator.Order and property name
	//
	// When this is set to false (default) properties are checked in the order in which they appear in the properties map -
	// which is unpredictable
	//
	// Note: If any of the properties in the validator has PropertyValidator.Order set to a non-zero value
	// then ordered property checks are also performed
	OrderedPropertyChecks bool
	// WhenConditions is the condition tokens that dictate under which conditions this validator is to be checked
	//
	// Condition tokens can be set and unset during validation to allow polymorphism of validation
	// (see ValidatorContext.SetCondition & ValidatorContext.ClearCondition)
	WhenConditions []string
	// ConditionalVariants represents a slice of ConditionalVariant items - where the first one that has the condition
	// satisfied is used (only one is ever used!)
	//
	// If none of the conditionals is satisfied, the validation falls back to using the
	// parent (this) Validator
	//
	// Condition tokens can be set and unset during validation to allow polymorphism of validation
	// (see ValidatorContext.SetCondition & ValidatorContext.ClearCondition)
	ConditionalVariants ConditionalVariants
	// OasInfo is additional information (for OpenAPI Specification) - used for generating and reading OAS
	OasInfo *OasInfo
}

const (
	msgUnableToDecode                 = "Unable to decode as JSON"
	CodeUnableToDecode                = 40001
	msgNotJsonNull                    = "JSON must not be JSON null"
	CodeNotJsonNull                   = 40002
	msgNotJsonArray                   = "JSON must not be JSON array"
	CodeNotJsonArray                  = 42200
	msgNotJsonObject                  = "JSON must not be JSON object"
	CodeNotJsonObject                 = 42201
	msgExpectedJsonArray              = "JSON expected to be JSON array"
	CodeExpectedJsonArray             = 42202
	msgExpectedJsonObject             = "JSON expected to be JSON object"
	CodeExpectedJsonObject            = 40003
	msgErrorReading                   = "Unexpected error reading reader"
	CodeErrorReading                  = 40004
	msgErrorUnmarshall                = "Unexpected error during unmarshalling"
	CodeErrorUnmarshall               = 40005
	msgRequestBodyEmpty               = "Request body is empty"
	CodeRequestBodyEmpty              = 40006
	msgUnableToDecodeRequest          = "Unable to decode request body as JSON"
	CodeUnableToDecodeRequest         = 40007
	msgRequestBodyNotJsonNull         = "Request body must not be JSON null"
	CodeRequestBodyNotJsonNull        = 40008
	msgRequestBodyNotJsonArray        = "Request body must not be JSON array"
	CodeRequestBodyNotJsonArray       = 42203
	msgRequestBodyNotJsonObject       = "Request body must not be JSON object"
	CodeRequestBodyNotJsonObject      = 42204
	msgRequestBodyExpectedJsonArray   = "Request body expected to be JSON array"
	CodeRequestBodyExpectedJsonArray  = 42205
	msgRequestBodyExpectedJsonObject  = "Request body expected to be JSON object"
	CodeRequestBodyExpectedJsonObject = 40009
	msgArrayElementMustBeObject       = "JSON array element must be an object"
	CodeArrayElementMustBeObject      = 42206
	msgMissingProperty                = "Missing property"
	CodeMissingProperty               = 42207
	msgUnwantedProperty               = "Property must not be present"
	CodeUnwantedProperty              = 42208
	msgUnknownProperty                = "Unknown property"
	CodeUnknownProperty               = 42209
	msgInvalidProperty                = "Invalid property"
	CodeInvalidProperty               = 42210
	msgInvalidPropertyName            = "Invalid property name"
	CodeInvalidPropertyName           = 42217
	msgPropertyValueMustBeObject      = "Property value must be an object"
	CodePropertyValueMustBeObject     = 42218
	CodeValidatorConstraintFail       = 42298
)

// DefaultDecoderProvider is the decoder provider used by Validator - replace with your own if necessary
var DefaultDecoderProvider DecoderProvider = &defaultDecoderProvider{}

// DecoderProvider is the interface needed for replacing the DefaultDecoderProvider
type DecoderProvider interface {
	NewDecoder(r io.Reader, useNumber bool) *json.Decoder
}

type defaultDecoderProvider struct{}

func (ddp *defaultDecoderProvider) NewDecoder(r io.Reader, useNumber bool) *json.Decoder {
	d := json.NewDecoder(r)
	if useNumber {
		d.UseNumber()
	}
	return d
}

// Properties type used by Validator.Properties
type Properties map[string]*PropertyValidator

// Constraints type used by Validator.Constraints and PropertyValidator.Constraints
type Constraints []Constraint

// ConditionalVariants type used by Validator.ConditionalVariants
type ConditionalVariants []*ConditionalVariant

// ConditionalVariant represents the condition(s) under which to use a specific variant Validator
type ConditionalVariant struct {
	// WhenConditions is the condition tokens that determine when this variant is used
	WhenConditions []string
	Constraints    Constraints
	Properties     Properties
	// ConditionalVariants is any descendant conditional variants
	ConditionalVariants ConditionalVariants
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
	tmpVcx := newEmptyValidatorContext(obtainI18nProvider().ContextFromRequest(r))
	ok, obj := v.decodeRequestBody(r.Body, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().ContextFromRequest(r))
	v.requestBodyValidate(vcx, obj)
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) decodeRequestBody(r io.Reader, vcx *ValidatorContext) (bool, interface{}) {
	if r == nil {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgRequestBodyEmpty), CodeRequestBodyEmpty, nil))
		return false, nil
	}
	decoder := DefaultDecoderProvider.NewDecoder(r, v.UseNumber)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgUnableToDecodeRequest), CodeUnableToDecodeRequest, err))
		return false, nil
	}
	return true, obj
}

func (v *Validator) requestBodyValidate(vcx *ValidatorContext, obj interface{}) {
	if obj != nil {
		// determine whether body is a map (object) or a slice (array)...
		if arr, isArr := obj.([]interface{}); isArr {
			if v.AllowArray {
				v.validateArrayOf(arr, vcx)
			} else {
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgRequestBodyNotJsonArray), CodeRequestBodyNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject && v.AllowArray {
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgRequestBodyExpectedJsonArray), CodeRequestBodyExpectedJsonArray))
			} else if v.DisallowObject {
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgRequestBodyNotJsonObject), CodeRequestBodyNotJsonObject))
			} else {
				v.validate(m, vcx)
			}
		} else {
			vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgRequestBodyExpectedJsonObject), CodeRequestBodyExpectedJsonObject, nil))
		}
	} else if !v.AllowNullJson {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgRequestBodyNotJsonNull), CodeRequestBodyNotJsonNull, nil))
	}
}

// Validate performs validation on the supplied JSON object
//
// Where the JSON object is represented as an unmarshalled
//   map[string]interface{}
func (v *Validator) Validate(obj map[string]interface{}) (bool, []*Violation) {
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
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
	vcx := newValidatorContext(arr, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	v.validateArrayOf(arr, vcx)
	return vcx.ok, vcx.violations
}

// ValidateReader performs validation on the supplied reader (representing JSON)
func (v *Validator) ValidateReader(r io.Reader) (bool, []*Violation, interface{}) {
	decoder := DefaultDecoderProvider.NewDecoder(r, v.UseNumber)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx := newEmptyValidatorContext(obtainI18nProvider().DefaultContext())
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgUnableToDecode), CodeUnableToDecode, err))
		return vcx.ok, vcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
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
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgNotJsonArray), CodeNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject && v.AllowArray {
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgExpectedJsonArray), CodeExpectedJsonArray))
			} else if v.DisallowObject {
				vcx.AddViolation(NewEmptyViolation(vcx.TranslateMessage(msgNotJsonObject), CodeNotJsonObject))
			} else {
				v.validate(m, vcx)
			}
		} else {
			vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgExpectedJsonObject), CodeExpectedJsonObject, nil))
		}
	} else if !v.AllowNullJson {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgNotJsonNull), CodeNotJsonNull, nil))
	}
}

type ValidationError struct {
	Message      string
	Violations   []*Violation
	IsBadRequest bool
}

func (ve *ValidationError) Error() string {
	return ve.Message
}

// ValidateInto performs validation on the supplied data (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
//
// If validation is unsuccessful (i.e. any violations) this method returns a ValidationError
func (v *Validator) ValidateInto(data []byte, value interface{}) error {
	r := bytes.NewReader(data)
	ok, violations, _ := v.ValidateReaderInto(r, value)
	if ok {
		return nil
	}
	isBad := false
	msg := msgErrorUnmarshall
	if l := len(violations); l == 1 {
		msg = violations[0].Message
		isBad = violations[0].BadRequest
		if isBad {
			violations = []*Violation{}
		}
	} else if l > 0 {
		isBad = violations[0].BadRequest
		SortViolationsByPathAndProperty(violations)
		msg = violations[0].Message
	}
	return &ValidationError{
		Message:      msg,
		Violations:   violations,
		IsBadRequest: isBad,
	}
}

// ValidateReaderInto performs validation on the supplied reader (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) ValidateReaderInto(r io.Reader, value interface{}) (bool, []*Violation, interface{}) {
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		errVcx := newEmptyValidatorContext(obtainI18nProvider().DefaultContext())
		errVcx.AddViolation(NewBadRequestViolation(errVcx.TranslateMessage(msgErrorReading), CodeErrorReading, err))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	decoder := DefaultDecoderProvider.NewDecoder(initialReader, v.UseNumber)
	var obj interface{} = reflect.Interface
	if dErr := decoder.Decode(&obj); dErr != nil {
		errVcx := newEmptyValidatorContext(nil)
		errVcx.AddViolation(NewBadRequestViolation(errVcx.TranslateMessage(msgUnableToDecode), CodeUnableToDecode, dErr))
		return false, errVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	v.validateObjectOrArray(vcx, obj)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder = DefaultDecoderProvider.NewDecoder(intoReader, v.UseNumber)
	if !v.IgnoreUnknownProperties {
		decoder.DisallowUnknownFields()
	}
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgErrorUnmarshall), CodeErrorUnmarshall, err))
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
	i18ctx := obtainI18nProvider().ContextFromRequest(req)
	if req.Body == nil {
		errVcx := newEmptyValidatorContext(i18ctx)
		errVcx.AddViolation(NewBadRequestViolation(errVcx.TranslateMessage(msgRequestBodyEmpty), CodeRequestBodyEmpty, nil))
		return false, errVcx.violations, nil
	}
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errVcx := newEmptyValidatorContext(i18ctx)
		errVcx.AddViolation(NewBadRequestViolation(errVcx.TranslateMessage(msgErrorReading), CodeErrorReading, err))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	tmpVcx := newEmptyValidatorContext(i18ctx)
	ok, obj := v.decodeRequestBody(initialReader, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, i18ctx)
	v.requestBodyValidate(vcx, obj)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder := DefaultDecoderProvider.NewDecoder(intoReader, v.UseNumber)
	if !v.IgnoreUnknownProperties {
		decoder.DisallowUnknownFields()
	}
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(NewBadRequestViolation(vcx.TranslateMessage(msgErrorUnmarshall), CodeErrorUnmarshall, err))
	}
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) validate(obj map[string]interface{}, vcx *ValidatorContext) {
	if checkConstraints(obj, vcx, v.Constraints) {
		return
	}
	if has, variant := getMatchingVariant(vcx, v.ConditionalVariants); has {
		v.variantValidate(obj, vcx, *variant, v.ConditionalVariants, v.Properties.clone(), v.Properties.clone())
		return
	}
	if checkUnknownProperties(obj, vcx, v.IgnoreUnknownProperties, v.Properties, nil) {
		return
	}
	v.checkProperties(obj, vcx, v.Properties)
}

func (v *Validator) variantValidate(obj map[string]interface{}, vcx *ValidatorContext, variant ConditionalVariant, variants ConditionalVariants, properties Properties, others Properties) {
	if checkConstraints(obj, vcx, variant.Constraints) {
		return
	}
	if has, subVariant := getMatchingVariant(vcx, variant.ConditionalVariants); has {
		// add all properties for the current variant...
		for k, pv := range variant.Properties {
			properties[k] = pv
		}
		// collect properties from other variants (need to know if they are known when doing unknown property checks)...
		for _, other := range variant.ConditionalVariants {
			for k, pv := range other.Properties {
				others[k] = pv
			}
		}
		v.variantValidate(obj, vcx, *subVariant, variant.ConditionalVariants, properties, others)
		return
	}
	for k, pv := range variant.Properties {
		properties[k] = pv
	}
	for _, other := range variants {
		for k, pv := range other.Properties {
			others[k] = pv
		}
	}
	if checkUnknownProperties(obj, vcx, v.IgnoreUnknownProperties, properties, others) {
		return
	}
	v.checkProperties(obj, vcx, properties)
}

func (v *Validator) checkProperties(obj map[string]interface{}, vcx *ValidatorContext, properties Properties) (stops bool) {
	if v.IsOrderedPropertyChecks() {
		sorted := sortProperties(properties)
		for _, p := range sorted {
			actualValue, present := obj[p.name]
			if present && !vcx.meetsUnwantedConditions(p.pv.UnwantedConditions) {
				vcx.addViolationPropertyForCurrent(p.name, msgUnwantedProperty, CodeUnwantedProperty, p.name)
			} else if vcx.meetsWhenConditions(p.pv.WhenConditions) {
				if !present {
					if p.pv.Mandatory {
						vcx.addViolationPropertyForCurrent(p.name, msgMissingProperty, CodeMissingProperty, p.name)
					}
				} else {
					vcx.pushPathProperty(p.name, actualValue, p.pv)
					p.pv.validate(actualValue, vcx)
					vcx.popPath()
				}
			}
			if !vcx.continueAll {
				return true
			}
		}
	} else {
		for propertyName, pv := range properties {
			actualValue, present := obj[propertyName]
			if present && !vcx.meetsUnwantedConditions(pv.UnwantedConditions) {
				vcx.addViolationPropertyForCurrent(propertyName, msgUnwantedProperty, CodeUnwantedProperty, propertyName)
			} else if vcx.meetsWhenConditions(pv.WhenConditions) {
				if !present {
					if pv.Mandatory {
						vcx.addViolationPropertyForCurrent(propertyName, msgMissingProperty, CodeMissingProperty, propertyName)
					}
				} else {
					vcx.pushPathProperty(propertyName, actualValue, pv)
					pv.validate(actualValue, vcx)
					vcx.popPath()
				}
			}
			if !vcx.continueAll {
				return true
			}
		}
	}
	return false
}

func getMatchingVariant(vcx *ValidatorContext, variants ConditionalVariants) (has bool, variant *ConditionalVariant) {
	has = false
	variant = nil
	for _, cv := range variants {
		if vcx.meetsWhenConditions(cv.WhenConditions) {
			has = true
			variant = cv
			break
		}
	}
	return
}

func (v *Validator) validateArrayOf(arr []interface{}, vcx *ValidatorContext) {
	for i, elem := range arr {
		vcx.pushPathIndex(i, elem, v)
		if obj, itemOk := elem.(map[string]interface{}); itemOk {
			v.validate(obj, vcx)
		} else {
			vcx.addUnTranslatedViolationForCurrent(msgArrayElementMustBeObject, CodeArrayElementMustBeObject, i)
		}
		vcx.popPath()
		if !vcx.continueAll {
			return
		}
	}
}

func checkUnknownProperties(obj map[string]interface{}, vcx *ValidatorContext, ignoreUnknowns bool, properties Properties, others Properties) (stops bool) {
	if !ignoreUnknowns {
		for propertyName := range obj {
			if _, has := properties[propertyName]; !has {
				msg := msgUnknownProperty
				code := CodeUnknownProperty
				if _, other := others[propertyName]; other {
					// adjust the message and code - because the property is known but not valid...
					msg = msgInvalidProperty
					code = CodeInvalidProperty
				}
				vcx.addViolationPropertyForCurrent(propertyName, msg, code, propertyName)
				if !vcx.continueAll {
					return true
				}
			}
		}
	}
	return false
}

func checkConstraints(obj map[string]interface{}, vcx *ValidatorContext, constraints Constraints) (stops bool) {
	for i, constraint := range constraints {
		if ok, msg := constraint.Check(obj, vcx); !ok {
			// the message is already translated by the constraint
			vcx.addTranslatedViolationForCurrent(msg, CodeValidatorConstraintFail, i)
		}
		if !vcx.continueAll {
			return true
		}
	}
	return false
}

func (v *Validator) IsOrderedPropertyChecks() bool {
	if !v.OrderedPropertyChecks {
		for _, pv := range v.Properties {
			if pv.Order != 0 {
				return true
			}
		}
		return false
	}
	return true
}

type orderedProperty struct {
	name string
	pv   PropertyValidator
}

func sortProperties(properties Properties) []orderedProperty {
	result := make([]orderedProperty, len(properties))
	i := 0
	for k, pv := range properties {
		result[i] = orderedProperty{name: k, pv: *pv}
		i++
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].pv.Order == result[j].pv.Order {
			return result[i].name < result[j].name
		}
		return result[i].pv.Order < result[j].pv.Order
	})
	return result
}

func (p Properties) clone() Properties {
	result := make(Properties, len(p))
	for k, v := range p {
		result[k] = v
	}
	return result
}

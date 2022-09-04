package valix

import (
	"bytes"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
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
	// AllowNullItems [for arrays only] denotes whether null array items are allowed
	AllowNullItems bool
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
	WhenConditions Conditions
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
	msgPropertyRequiredWhen           = "Property is required under certain criteria"
	CodePropertyRequiredWhen          = 42219
	msgPropertyUnwantedWhen           = "Property must not be present under certain criteria"
	CodePropertyUnwantedWhen          = 42220
	msgArrayElementMustNotBeNull      = "JSON array element must not be null"
	CodeArrayElementMustNotBeNull     = 42221
	CodeValidatorConstraintFail       = 42298
)

// Properties type used by Validator.Properties
type Properties map[string]*PropertyValidator

// Constraints type used by Validator.Constraints and PropertyValidator.Constraints
type Constraints []Constraint

// ConditionalVariants type used by Validator.ConditionalVariants
type ConditionalVariants []*ConditionalVariant

// ConditionalVariant represents the condition(s) under which to use a specific variant Validator
type ConditionalVariant struct {
	// WhenConditions is the condition tokens that determine when this variant is used
	WhenConditions Conditions
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
func (v *Validator) RequestValidate(req *http.Request, initialConditions ...string) (bool, []*Violation, interface{}) {
	tmpVcx := newEmptyValidatorContext(obtainI18nProvider().ContextFromRequest(req))
	ok, obj := v.decodeRequestBody(req.Body, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().ContextFromRequest(req))
	vcx.setConditionsFromRequest(req)
	vcx.setInitialConditions(initialConditions...)
	v.validateObjectOrArray(vcx, obj, true)
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) decodeRequestBody(r io.Reader, vcx *ValidatorContext) (bool, interface{}) {
	if r == nil {
		vcx.AddViolation(newBadRequestViolation(vcx, msgRequestBodyEmpty, CodeRequestBodyEmpty, nil))
		return false, nil
	}
	decoder := getDefaultDecoderProvider().NewDecoder(r, v.UseNumber)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx.AddViolation(newBadRequestViolation(vcx, msgUnableToDecodeRequest, CodeUnableToDecodeRequest, err))
		return false, nil
	}
	return true, obj
}

// Validate performs validation on the supplied JSON object
//
// Where the JSON object is represented as an unmarshalled
//   map[string]interface{}
func (v *Validator) Validate(obj map[string]interface{}, initialConditions ...string) (bool, []*Violation) {
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	vcx.setInitialConditions(initialConditions...)
	v.validate(obj, vcx)
	return vcx.ok, vcx.violations
}

// ValidateArrayOf Performs validation on each element of the supplied JSON array
//
// Where the JSON array is represented as an unmarshalled
//   []interface{}
// and each item of the slice is expected to be a JSON object represented as an unmarshalled
//   map[string]interface{}
func (v *Validator) ValidateArrayOf(arr []interface{}, initialConditions ...string) (bool, []*Violation) {
	vcx := newValidatorContext(arr, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	vcx.setInitialConditions(initialConditions...)
	v.validateArrayOf(arr, vcx)
	return vcx.ok, vcx.violations
}

// ValidateReader performs validation on the supplied reader (representing JSON)
func (v *Validator) ValidateReader(r io.Reader, initialConditions ...string) (bool, []*Violation, interface{}) {
	decoder := getDefaultDecoderProvider().NewDecoder(r, v.UseNumber)
	var obj interface{} = reflect.Interface
	if err := decoder.Decode(&obj); err != nil {
		vcx := newEmptyValidatorContext(obtainI18nProvider().DefaultContext())
		vcx.AddViolation(newBadRequestViolation(vcx, msgUnableToDecode, CodeUnableToDecode, err))
		return vcx.ok, vcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	vcx.setInitialConditions(initialConditions...)
	v.validateObjectOrArray(vcx, obj, false)
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) validateObjectOrArray(vcx *ValidatorContext, obj interface{}, isRequest bool) {
	var violation *Violation = nil
	if obj != nil {
		// determine whether body is a map (object) or a slice (array)...
		if arr, isArr := obj.([]interface{}); isArr {
			if v.AllowArray {
				v.validateArrayOf(arr, vcx)
			} else {
				violation = newEmptyViolation(vcx,
					ternary(isRequest).string(msgRequestBodyNotJsonArray, msgNotJsonArray),
					ternary(isRequest).int(CodeRequestBodyNotJsonArray, CodeNotJsonArray))
			}
		} else if m, isMap := obj.(map[string]interface{}); isMap {
			if v.DisallowObject && v.AllowArray {
				violation = newEmptyViolation(vcx,
					ternary(isRequest).string(msgRequestBodyExpectedJsonArray, msgExpectedJsonArray),
					ternary(isRequest).int(CodeRequestBodyExpectedJsonArray, CodeExpectedJsonArray))
			} else if v.DisallowObject {
				violation = newEmptyViolation(vcx,
					ternary(isRequest).string(msgRequestBodyNotJsonObject, msgNotJsonObject),
					ternary(isRequest).int(CodeRequestBodyNotJsonObject, CodeNotJsonObject))
			} else {
				v.validate(m, vcx)
			}
		} else {
			violation = newBadRequestViolation(vcx,
				ternary(isRequest).string(msgRequestBodyExpectedJsonObject, msgExpectedJsonObject),
				ternary(isRequest).int(CodeRequestBodyExpectedJsonObject, CodeExpectedJsonObject))
		}
	} else if !v.AllowNullJson {
		violation = newBadRequestViolation(vcx,
			ternary(isRequest).string(msgRequestBodyNotJsonNull, msgNotJsonNull),
			ternary(isRequest).int(CodeRequestBodyNotJsonNull, CodeNotJsonNull))
	}
	if violation != nil {
		vcx.AddViolation(violation)
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
func (v *Validator) ValidateInto(data []byte, value interface{}, initialConditions ...string) error {
	r := bytes.NewReader(data)
	ok, violations, _ := v.ValidateReaderInto(r, value, initialConditions...)
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
func (v *Validator) ValidateReaderInto(r io.Reader, value interface{}, initialConditions ...string) (bool, []*Violation, interface{}) {
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(r)
	if err != nil {
		errVcx := newEmptyValidatorContext(obtainI18nProvider().DefaultContext())
		errVcx.AddViolation(newBadRequestViolation(errVcx, msgErrorReading, CodeErrorReading, err))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	decoder := getDefaultDecoderProvider().NewDecoder(initialReader, v.UseNumber)
	var obj interface{} = reflect.Interface
	if dErr := decoder.Decode(&obj); dErr != nil {
		errVcx := newEmptyValidatorContext(nil)
		errVcx.AddViolation(newBadRequestViolation(errVcx, msgUnableToDecode, CodeUnableToDecode, dErr))
		return false, errVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, obtainI18nProvider().DefaultContext())
	vcx.setInitialConditions(initialConditions...)
	v.validateObjectOrArray(vcx, obj, false)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder = getDefaultDecoderProvider().NewDecoderFor(intoReader, v)
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(newBadRequestViolation(vcx, msgErrorUnmarshall, CodeErrorUnmarshall, err))
	}
	return vcx.ok, vcx.violations, obj
}

// ValidateString performs validation on the supplied string (representing JSON)
func (v *Validator) ValidateString(s string, initialConditions ...string) (bool, []*Violation, interface{}) {
	return v.ValidateReader(strings.NewReader(s), initialConditions...)
}

// ValidateStringInto performs validation on the supplied string (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) ValidateStringInto(s string, value interface{}, initialConditions ...string) (bool, []*Violation, interface{}) {
	return v.ValidateReaderInto(strings.NewReader(s), value, initialConditions...)
}

// RequestValidateInto performs validation on the request body (representing JSON)
// and, if validation successful, attempts to unmarshall the JSON into the supplied value
func (v *Validator) RequestValidateInto(req *http.Request, value interface{}, initialConditions ...string) (bool, []*Violation, interface{}) {
	i18ctx := obtainI18nProvider().ContextFromRequest(req)
	if req.Body == nil {
		errVcx := newEmptyValidatorContext(i18ctx)
		errVcx.AddViolation(newBadRequestViolation(i18ctx, msgRequestBodyEmpty, CodeRequestBodyEmpty, nil))
		return false, errVcx.violations, nil
	}
	// we'll need to read the reader twice - first into our representation (for validation) and then into the value
	buffer, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errVcx := newEmptyValidatorContext(i18ctx)
		errVcx.AddViolation(newBadRequestViolation(i18ctx, msgErrorReading, CodeErrorReading, err))
		return false, errVcx.violations, nil
	}
	initialReader := bytes.NewReader(buffer)
	tmpVcx := newEmptyValidatorContext(i18ctx)
	ok, obj := v.decodeRequestBody(initialReader, tmpVcx)
	if !ok {
		return false, tmpVcx.violations, nil
	}
	vcx := newValidatorContext(obj, v, v.StopOnFirst, i18ctx)
	vcx.setConditionsFromRequest(req)
	vcx.setInitialConditions(initialConditions...)
	v.validateObjectOrArray(vcx, obj, true)
	if !vcx.ok {
		return false, vcx.violations, obj
	}
	// now read into the provided value...
	intoReader := bytes.NewReader(buffer)
	decoder := getDefaultDecoderProvider().NewDecoderFor(intoReader, v)
	err = decoder.Decode(value)
	if err != nil {
		vcx.AddViolation(newBadRequestViolation(vcx, msgErrorUnmarshall, CodeErrorUnmarshall, err))
	}
	return vcx.ok, vcx.violations, obj
}

func (v *Validator) validate(obj map[string]interface{}, vcx *ValidatorContext) {
	if checkConstraints(obj, vcx, v.Constraints) {
		return
	}
	if has, variant := getMatchingVariant(vcx, v.ConditionalVariants); has {
		v.variantValidate(obj, vcx, *variant, v.ConditionalVariants, v.Properties.Clone(), v.Properties.Clone())
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

func (v *Validator) checkProperties(obj map[string]interface{}, vcx *ValidatorContext, properties Properties) {
	//func (v *Validator) checkProperties(obj map[string]interface{}, vcx *ValidatorContext, properties Properties) (stops bool) {
	names, pvs := v.orderedProperties(properties)
	// before we check the property values, we need to check property required/not required...
	names, pvs, cont := checkPropertiesRequiredWithWithout(obj, vcx, names, pvs)
	if !cont {
		return
	}
	for i, propertyName := range names {
		pv := pvs[i]
		actualValue, present := obj[propertyName]
		if present && !vcx.meetsUnwantedConditions(pv.UnwantedConditions) {
			vcx.addViolationPropertyForCurrent(propertyName, msgUnwantedProperty, CodeUnwantedProperty, propertyName)
		} else if vcx.meetsWhenConditions(pv.WhenConditions) {
			if !present {
				if pv.Mandatory && (len(pv.MandatoryWhen) == 0 || vcx.meetsWhenConditions(pv.MandatoryWhen)) {
					vcx.addViolationPropertyForCurrent(propertyName, msgMissingProperty, CodeMissingProperty, propertyName)
				}
			} else {
				vcx.pushPathProperty(propertyName, actualValue, pv)
				pv.validate(actualValue, vcx)
				vcx.popPath()
			}
		}
		if !vcx.continueAll {
			return
		}
	}
}

func checkPropertiesRequiredWithWithout(obj map[string]interface{}, vcx *ValidatorContext, names []string, pvs []*PropertyValidator) ([]string, []*PropertyValidator, bool) {
	rNames := names
	rPvs := pvs
	for i := 0; i < len(rNames); {
		propertyName := rNames[i]
		pv := rPvs[i]
		_, exists := obj[propertyName]
		if !exists && len(pv.RequiredWith) > 0 && pv.RequiredWith.Evaluate(obj, vcx.ValuesAncestry(), vcx) {
			vcx.addViolationPropertyForCurrent(propertyName,
				ternary(pv.RequiredWithMessage == "").string(msgPropertyRequiredWhen, pv.RequiredWithMessage),
				CodePropertyRequiredWhen)
			// it's not present but required - leave it in the list
			i++
		}
		if exists && len(pv.UnwantedWith) > 0 && pv.UnwantedWith.Evaluate(obj, vcx.ValuesAncestry(), vcx) {
			vcx.addViolationPropertyForCurrent(propertyName,
				ternary(pv.UnwantedWithMessage == "").string(msgPropertyUnwantedWhen, pv.UnwantedWithMessage),
				CodePropertyUnwantedWhen)
			// if it's present but unwanted, need to remove it from the list so the value isn't also checked
			rNames = append(rNames[:i], rNames[i+1:]...)
			rPvs = append(rPvs[:i], rPvs[i+1:]...)
		} else {
			i++
		}
		if !vcx.continueAll {
			break
		}
	}
	return rNames, rPvs, vcx.continueAll
}

func (v *Validator) orderedProperties(properties Properties) ([]string, []*PropertyValidator) {
	needsSorting := v.OrderedPropertyChecks
	useProperties := propertiesRepo.fetch(properties)
	names := make([]string, len(useProperties))
	validators := make([]*PropertyValidator, len(useProperties))
	i := 0
	for pn, pv := range useProperties {
		names[i] = pn
		validators[i] = pv
		needsSorting = needsSorting || pv.Order != 0
		i++
	}
	if needsSorting {
		quickSortOrdersAndNames(validators, names)
	}
	return names, validators
}

func quickSortOrdersAndNames(pvs []*PropertyValidator, ns []string) {
	if len(ns) < 2 {
		return
	}
	left, right := 0, len(ns)-1
	pivot := rand.Int() % len(ns)

	ns[pivot], ns[right] = ns[right], ns[pivot]
	pvs[pivot], pvs[right] = pvs[right], pvs[pivot]

	for i, _ := range ns {
		if (pvs[i].Order < pvs[right].Order) || (pvs[i].Order == pvs[right].Order && ns[i] < ns[right]) {
			ns[left], ns[i] = ns[i], ns[left]
			pvs[left], pvs[i] = pvs[i], pvs[left]
			left++
		}
	}

	ns[left], ns[right] = ns[right], ns[left]
	pvs[left], pvs[right] = pvs[right], pvs[left]

	quickSortOrdersAndNames(pvs[:left], ns[:left])
	quickSortOrdersAndNames(pvs[left+1:], ns[left+1:])
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
		if elem == nil {
			if !v.AllowNullItems {
				vcx.addUnTranslatedViolationForCurrent(msgArrayElementMustNotBeNull, CodeArrayElementMustNotBeNull, i)
			}
		} else if obj, itemOk := elem.(map[string]interface{}); itemOk {
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

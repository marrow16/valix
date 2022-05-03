package valix

import (
	"encoding/json"
	"math"
	"strings"
)

const (
	msgValueCannotBeNull             = "Value cannot be null"
	CodeValueCannotBeNull            = 42211
	fmtMsgValueExpectedType          = "Value expected to be of type %[1]s"
	CodeValueExpectedType            = 42212
	msgValueMustBeObject             = "Value must be an object"
	CodeValueMustBeObject            = 42213
	msgValueMustBeArray              = "Value must be an array"
	CodeValueMustBeArray             = 42214
	msgValueMustBeObjectOrArray      = "Value must be an object or array"
	CodeValueMustBeObjectOrArray     = 42215
	msgPropertyObjectValidatorError  = "Object validator error - does not allow object or array!"
	CodePropertyObjectValidatorError = 42216
	CodePropertyConstraintFail       = 42299
)

// JsonType is the type for JSON values
type JsonType int

const (
	// JsonAny matches any JSON value type
	JsonAny JsonType = iota
	// JsonString checks JSON value type is a string
	JsonString
	// JsonNumber checks JSON value type is a number
	JsonNumber
	// JsonInteger checks JSON value type is a number (that is or can be expressed as an int)
	JsonInteger
	// JsonBoolean checks JSON value type is a boolean
	JsonBoolean
	// JsonObject checks JSON value type is an object
	JsonObject
	// JsonArray checks JSON value type is an array
	JsonArray
)

const (
	jsonTypeTokenString  = "string"
	jsonTypeTokenNumber  = "number"
	jsonTypeTokenInteger = "integer"
	jsonTypeTokenBoolean = "boolean"
	jsonTypeTokenObject  = "object"
	jsonTypeTokenArray   = "array"
	jsonTypeTokenAny     = "any"
	jsonTypeTokensList   = "\"" + jsonTypeTokenAny + "\",\"" + jsonTypeTokenString +
		"\",\"" + jsonTypeTokenNumber + "\",\"" + jsonTypeTokenInteger +
		"\",\"" + jsonTypeTokenObject +
		"\",\"" + jsonTypeTokenObject + "\",\"" + jsonTypeTokenArray
)

func (jt JsonType) String() string {
	result := ""
	switch jt {
	case JsonString:
		result = jsonTypeTokenString
	case JsonNumber:
		result = jsonTypeTokenNumber
	case JsonInteger:
		result = jsonTypeTokenInteger
	case JsonBoolean:
		result = jsonTypeTokenBoolean
	case JsonObject:
		result = jsonTypeTokenObject
	case JsonArray:
		result = jsonTypeTokenArray
	case JsonAny:
		result = jsonTypeTokenAny
	}
	return result
}

func JsonTypeFromString(str string) (JsonType, bool) {
	result := JsonType(-1)
	ok := false
	switch strings.ToLower(str) {
	case jsonTypeTokenString:
		result = JsonString
		ok = true
	case jsonTypeTokenNumber:
		result = JsonNumber
		ok = true
	case jsonTypeTokenInteger:
		result = JsonInteger
		ok = true
	case jsonTypeTokenBoolean:
		result = JsonBoolean
		ok = true
	case jsonTypeTokenObject:
		result = JsonObject
		ok = true
	case jsonTypeTokenArray:
		result = JsonArray
		ok = true
	case jsonTypeTokenAny:
		result = JsonAny
		ok = true
	}
	return result, ok
}

// PropertyValidator is the individual validator for properties
type PropertyValidator struct {
	// Type specifies the property type to be checked (i.e. one of Type)
	//
	// If this value is not one of Type (or an empty string), then the property type is not checked
	Type JsonType
	// NotNull specifies that the value of the property may not be null
	NotNull bool
	// Mandatory specifies that the property must be present
	Mandatory bool
	// MandatoryWhen is complimentary to the Mandatory property - and limits the conditions under which the property is
	// seen as mandatory (if this is empty and Mandatory is set - then the property is always mandatory)
	//
	// Note: If the Mandatory property is not set to true - this property has no effect!
	MandatoryWhen Conditions
	// Constraints is a slice of Constraint items and are checked in the order they are specified
	Constraints Constraints
	// ObjectValidator is checked, if specified, after all Constraints are checked
	ObjectValidator *Validator
	// Order is the order in which the property is checked (see Validator.OrderedPropertyChecks)
	//
	// Note: setting any property with Order other than 0 (zero) will force the validator to use ordered property checks
	// (i.e. as if Validator.OrderedPropertyChecks had been set to true)
	Order int
	// WhenConditions is the condition tokens that dictate under which conditions this validator is to be checked
	//
	// Condition tokens can be set and unset during validation to allow polymorphism of validation
	// (see ValidatorContext.SetCondition & ValidatorContext.ClearCondition)
	WhenConditions Conditions
	// UnwantedConditions is the condition tokens that dictate when the property should not be present
	UnwantedConditions Conditions
	// RequiredWith is an expression of when this property is required according to the presence of other properties
	//
	// Use MustParseExpression or ParseExpression to build the expression - or build in code directly using
	// combinations of OthersExpr, OtherProperty and OtherGrouping
	RequiredWith OthersExpr
	// RequiredWithMessage is the violation message to use when the RequiredWith fails (if this string is empty, then
	// the default message is used)
	RequiredWithMessage string
	// UnwantedWith is an expression of when this property is unwanted according to the presence of other properties
	//
	// Use MustParseExpression or ParseExpression to build the expression - or build in code directly using
	// combinations of OthersExpr, OtherProperty and OtherGrouping
	UnwantedWith OthersExpr
	// UnwantedWithMessage is the violation message to use when the UnwantedWith fails (if this string is empty, then
	// the default message is used)
	UnwantedWithMessage string
	// OasInfo is additional information (for OpenAPI Specification)
	OasInfo *OasInfo
}

func (pv *PropertyValidator) validate(value interface{}, vcx *ValidatorContext) {
	if value == nil && pv.NotNull {
		vcx.addUnTranslatedViolationForCurrent(msgValueCannotBeNull, CodeValueCannotBeNull)
		return
	}
	// only check the type if non-nil...
	if value == nil || pv.checkType(value, vcx) {
		pv.checkConstraints(vcx)
		if vcx.continueAll && vcx.continuePty() {
			pv.checkObjectValidation(value, vcx)
		}
	}
}

func (pv *PropertyValidator) checkType(actualValue interface{}, vcx *ValidatorContext) bool {
	ok := checkValueType(actualValue, pv.Type)
	if !ok {
		vcx.addTranslatedViolationForCurrent(vcx.TranslateFormat(fmtMsgValueExpectedType, vcx.TranslateToken(pv.Type.String())), CodeValueExpectedType)
	}
	return ok
}

func checkValueType(value interface{}, t JsonType) bool {
	ok := true
	switch t {
	case JsonString:
		_, ok = value.(string)
		break
	case JsonBoolean:
		_, ok = value.(bool)
		break
	case JsonNumber:
		ok = checkNumeric(value, false)
		break
	case JsonInteger:
		ok = checkNumeric(value, true)
		break
	case JsonObject:
		_, ok = value.(map[string]interface{})
		break
	case JsonArray:
		_, ok = value.([]interface{})
		break
	}
	return ok
}

func checkNumeric(value interface{}, isInt bool) bool {
	var ok = false
	if fVal, fOk := value.(float64); fOk {
		ok = !isInt || (math.Trunc(fVal) == fVal)
	} else if nVal, nOk := value.(json.Number); nOk {
		// using json.JsonNumber.Float64() to parse - as this still allows for e notation (e.g. "0.1e1" is a valid int)
		if f, err := nVal.Float64(); err == nil {
			ok = !isInt || (math.Trunc(f) == f)
		}
	} else {
		_, ok = value.(int)
	}
	return ok
}

func (pv *PropertyValidator) checkObjectValidation(value interface{}, vcx *ValidatorContext) {
	if pv.ObjectValidator != nil && vcx.meetsWhenConditions(pv.ObjectValidator.WhenConditions) {
		if !pv.ObjectValidator.DisallowObject && pv.ObjectValidator.AllowArray {
			// can be object or array...
			if !pv.subValidateObjectOrArray(value, vcx) {
				vcx.addUnTranslatedViolationForCurrent(msgValueMustBeObjectOrArray, CodeValueMustBeObjectOrArray)
			}
		} else if !pv.ObjectValidator.DisallowObject {
			// can only be an object...
			if !pv.subValidateObject(value, vcx) {
				vcx.addUnTranslatedViolationForCurrent(msgValueMustBeObject, CodeValueMustBeObject)
			}
		} else if pv.ObjectValidator.AllowArray {
			// can only be an array...
			if !pv.subValidateArray(value, vcx) {
				vcx.addUnTranslatedViolationForCurrent(msgValueMustBeArray, CodeValueMustBeArray)
			}
		} else {
			// something seriously wrong here because the object validator doesn't allow an object or an array!...
			vcx.addUnTranslatedViolationForCurrent(msgPropertyObjectValidatorError, CodePropertyObjectValidatorError)
			vcx.Stop()
		}
	}
}

func (pv *PropertyValidator) subValidateObjectOrArray(actualValue interface{}, vcx *ValidatorContext) bool {
	return pv.subValidateObject(actualValue, vcx) || pv.subValidateArray(actualValue, vcx)
}

func (pv *PropertyValidator) subValidateObject(actualValue interface{}, vcx *ValidatorContext) bool {
	if o, ok := actualValue.(map[string]interface{}); ok {
		pv.ObjectValidator.validate(o, vcx)
		return true
	}
	return false
}

func (pv *PropertyValidator) subValidateArray(actualValue interface{}, vcx *ValidatorContext) bool {
	if a, ok := actualValue.([]interface{}); ok {
		pv.ObjectValidator.validateArrayOf(a, vcx)
		return true
	}
	return false
}

func (pv *PropertyValidator) checkConstraints(vcx *ValidatorContext) {
	if pv.Constraints != nil {
		v := vcx.CurrentValue()
		for i, constraint := range pv.Constraints {
			if ok, msg := constraint.Check(v, vcx); !ok {
				// the message is already translated by the constraint!...
				vcx.addTranslatedViolationForCurrent(msg, CodePropertyConstraintFail, i)
			}
			if !vcx.continueAll || !vcx.continuePty() {
				return
			}
		}
	}
}

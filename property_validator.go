package valix

import (
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

const (
	messageValueCannotBeNull            = "Value cannot be null"
	messageValueExpectedType            = "Value expected to be of type %s"
	messageValueMustBeObject            = "Value must be an object"
	messageValueMustBeArray             = "Value must be an array"
	messageValueMustBeObjectOrArray     = "Value must be an object or array"
	messagePropertyObjectValidatorError = "CurrentProperty object validator error - does not allow object or array!"
)

// JsonType is the type for JSON values
type JsonType int

const (
	JsonTypeUndefined JsonType = iota
	// JsonAny matches any JSON value type
	JsonAny
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
	jsonTypeTokenUndefined = "undefined"
	jsonTypeTokenString    = "string"
	jsonTypeTokenNumber    = "number"
	jsonTypeTokenInteger   = "integer"
	jsonTypeTokenBoolean   = "boolean"
	jsonTypeTokenObject    = "object"
	jsonTypeTokenArray     = "array"
	jsonTypeTokenAny       = "any"
)

func (jt JsonType) String() string {
	result := jsonTypeTokenUndefined
	switch jt {
	case JsonString:
		result = jsonTypeTokenString
		break
	case JsonNumber:
		result = jsonTypeTokenNumber
		break
	case JsonInteger:
		result = jsonTypeTokenInteger
		break
	case JsonBoolean:
		result = jsonTypeTokenBoolean
		break
	case JsonObject:
		result = jsonTypeTokenObject
		break
	case JsonArray:
		result = jsonTypeTokenArray
		break
	case JsonAny:
		result = jsonTypeTokenAny
		break
	}
	return result
}

func JsonTypeFromString(str string) JsonType {
	result := JsonTypeUndefined
	switch strings.ToLower(str) {
	case jsonTypeTokenString:
		result = JsonString
		break
	case jsonTypeTokenNumber:
		result = JsonNumber
		break
	case jsonTypeTokenInteger:
		result = JsonInteger
		break
	case jsonTypeTokenBoolean:
		result = JsonBoolean
		break
	case jsonTypeTokenObject:
		result = JsonObject
		break
	case jsonTypeTokenArray:
		result = JsonArray
		break
	case jsonTypeTokenAny:
		result = JsonAny
		break
	}
	return result
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
	// Constraints is a slice of Constraint items and are checked in the order they are specified
	Constraints Constraints
	// ObjectValidator is checked, if specified, after all Constraints are checked
	ObjectValidator *Validator
}

func (pv *PropertyValidator) validate(actualValue interface{}, vcx *ValidatorContext) {
	if actualValue == nil {
		if pv.NotNull {
			vcx.AddViolationForCurrent(messageValueCannotBeNull)
		}
	} else if pv.checkType(actualValue, vcx) {
		// don't pass the actualValue down further - because it may change!
		pv.checkConstraints(vcx)
		if vcx.continueAll && vcx.continuePty() {
			pv.checkObjectValidation(vcx)
		}
	}
}

func (pv *PropertyValidator) checkType(actualValue interface{}, vcx *ValidatorContext) bool {
	ok := checkValueType(actualValue, pv.Type)
	if !ok {
		vcx.AddViolationForCurrent(fmt.Sprintf(messageValueExpectedType, pv.Type))
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

func (pv *PropertyValidator) checkObjectValidation(vcx *ValidatorContext) {
	if pv.ObjectValidator != nil {
		if !pv.ObjectValidator.DisallowObject && pv.ObjectValidator.AllowArray {
			// can be object or array...
			if !pv.subValidateObjectOrArray(vcx.CurrentValue(), vcx) {
				vcx.AddViolationForCurrent(messageValueMustBeObjectOrArray)
			}
		} else if !pv.ObjectValidator.DisallowObject {
			// can only be an object...
			if !pv.subValidateObject(vcx.CurrentValue(), vcx) {
				vcx.AddViolationForCurrent(messageValueMustBeObject)
			}
		} else if pv.ObjectValidator.AllowArray {
			// can only be an array...
			if !pv.subValidateArray(vcx.CurrentValue(), vcx) {
				vcx.AddViolationForCurrent(messageValueMustBeArray)
			}
		} else {
			// something seriously wrong here because the object validator doesn't allow an object or an array!...
			vcx.AddViolationForCurrent(messagePropertyObjectValidatorError)
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
		for _, constraint := range pv.Constraints {
			// re-get the current value each time because it may have changed
			if ok, msg := constraint.Check(vcx.CurrentValue(), vcx); !ok {
				vcx.AddViolationForCurrent(msg)
			}
			if !vcx.continueAll || !vcx.continuePty() {
				return
			}
		}
	}
}

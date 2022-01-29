package valix

import (
	"encoding/json"
	"fmt"
	"math"
)

const (
	MessageValueCannotBeNull            = "Value cannot be null"
	MessageValueExpectedType            = "Value expected to be of type %s"
	MessageValueMustBeObject            = "Value must be an object"
	MessageValueMustBeArray             = "Value must be an array"
	MessageValueMustBeObjectOrArray     = "Value must be an object or array"
	MessagePropertyObjectValidatorError = "CurrentProperty object validator error - does not allow object or array!"
)

var PropertyType = newPropertyTypesEnum()

func newPropertyTypesEnum() *propertyType {
	return &propertyType{
		String:  "string",
		Number:  "number",
		Int:     "int",
		Boolean: "boolean",
		Object:  "object",
		Array:   "array",
	}
}

type propertyType struct {
	String  string
	Number  string
	Int     string
	Boolean string
	Object  string
	Array   string
}

type PropertyValidator struct {
	PropertyType string
	NotNull      bool
	Mandatory    bool
	// Constraints are checked in the order they are specified
	Constraints []Constraint
	// ObjectValidator is checked, if specified, after all Constraints are checked
	ObjectValidator *Validator
}

func (pv *PropertyValidator) validate(actualValue interface{}, ctx *Context) {
	if actualValue == nil {
		if pv.NotNull {
			ctx.AddViolationForCurrent(MessageValueCannotBeNull)
		}
	} else if pv.checkType(actualValue, ctx) {
		pv.checkConstraints(actualValue, ctx)
		if ctx.continueAll && ctx.continuePty() {
			pv.checkObjectValidation(actualValue, ctx)
		}
	}
}

func (pv *PropertyValidator) checkType(actualValue interface{}, ctx *Context) bool {
	ok := checkValueType(actualValue, pv.PropertyType)
	if !ok {
		ctx.AddViolationForCurrent(fmt.Sprintf(MessageValueExpectedType, pv.PropertyType))
	}
	return ok
}

func checkValueType(value interface{}, t string) bool {
	ok := true
	switch t {
	case PropertyType.String:
		_, ok = value.(string)
		break
	case PropertyType.Boolean:
		_, ok = value.(bool)
		break
	case PropertyType.Number:
		ok = checkNumeric(value, false)
		break
	case PropertyType.Int:
		ok = checkNumeric(value, true)
		break
	case PropertyType.Object:
		_, ok = value.(map[string]interface{})
		break
	case PropertyType.Array:
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
		// using json.Number.Float64() to parse - as this still allows for e notation (e.g. "0.1e1" is a valid int)
		if f, err := nVal.Float64(); err == nil {
			ok = !isInt || (math.Trunc(f) == f)
		}
	} else {
		_, ok = value.(int)
	}
	return ok
}

func (pv *PropertyValidator) checkObjectValidation(actualValue interface{}, ctx *Context) {
	if pv.ObjectValidator != nil {
		if !pv.ObjectValidator.DisallowObject && pv.ObjectValidator.AllowArray {
			// can be object or array...
			if !pv.subValidateObjectOrArray(actualValue, ctx) {
				ctx.AddViolationForCurrent(MessageValueMustBeObjectOrArray)
			}
		} else if !pv.ObjectValidator.DisallowObject {
			// can only be an object...
			if !pv.subValidateObject(actualValue, ctx) {
				ctx.AddViolationForCurrent(MessageValueMustBeObject)
			}
		} else if pv.ObjectValidator.AllowArray {
			// can only be an array...
			if !pv.subValidateArray(actualValue, ctx) {
				ctx.AddViolationForCurrent(MessageValueMustBeArray)
			}
		} else {
			// something seriously wrong here because the object validator doesn't allow an object or an array!...
			ctx.AddViolationForCurrent(MessagePropertyObjectValidatorError)
			ctx.Stop()
		}
	}
}

func (pv *PropertyValidator) subValidateObjectOrArray(actualValue interface{}, ctx *Context) bool {
	return pv.subValidateObject(actualValue, ctx) || pv.subValidateArray(actualValue, ctx)
}

func (pv *PropertyValidator) subValidateObject(actualValue interface{}, ctx *Context) bool {
	if o, ok := actualValue.(map[string]interface{}); ok {
		pv.ObjectValidator.validate(o, ctx)
		return true
	}
	return false
}

func (pv *PropertyValidator) subValidateArray(actualValue interface{}, ctx *Context) bool {
	if a, ok := actualValue.([]interface{}); ok {
		pv.ObjectValidator.validateArrayOf(a, ctx)
		return true
	}
	return false
}

func (pv *PropertyValidator) checkConstraints(actualValue interface{}, ctx *Context) {
	if pv.Constraints != nil {
		for _, constraint := range pv.Constraints {
			if ok, msg := constraint.Validate(actualValue, ctx); !ok {
				ctx.AddViolationForCurrent(msg)
			}
			if !ctx.continueAll || !ctx.continuePty() {
				return
			}
		}
	}
}

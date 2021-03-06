package valix

import (
	"encoding/json"
	"math"
	"strings"
	"time"
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
	// JsonDatetime checks JSON value type is a string and is a valid parseable datetime
	JsonDatetime
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
	jsonTypeTokenString   = "string"
	jsonTypeTokenDatetime = "datetime"
	jsonTypeTokenNumber   = "number"
	jsonTypeTokenInteger  = "integer"
	jsonTypeTokenBoolean  = "boolean"
	jsonTypeTokenObject   = "object"
	jsonTypeTokenArray    = "array"
	jsonTypeTokenAny      = "any"
	jsonTypeTokensList    = "\"" + jsonTypeTokenAny + "\",\"" + jsonTypeTokenString + "\",\"" + jsonTypeTokenDatetime +
		"\",\"" + jsonTypeTokenNumber + "\",\"" + jsonTypeTokenInteger +
		"\",\"" + jsonTypeTokenObject +
		"\",\"" + jsonTypeTokenObject + "\",\"" + jsonTypeTokenArray
)

func (jt JsonType) String() string {
	result := ""
	switch jt {
	case JsonString:
		result = jsonTypeTokenString
	case JsonDatetime:
		result = jsonTypeTokenDatetime
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
	case jsonTypeTokenDatetime:
		result = JsonDatetime
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
	case JsonDatetime:
		ok = false
		switch av := value.(type) {
		case time.Time, *time.Time, Time, *Time:
			ok = true
		case string:
			_, ok = stringToDatetime(av, false)
		}
	case JsonBoolean:
		_, ok = value.(bool)
	case JsonNumber:
		ok = checkNumeric(value, false)
	case JsonInteger:
		ok = checkNumeric(value, true)
	case JsonObject:
		_, ok = value.(map[string]interface{})
	case JsonArray:
		_, ok = value.([]interface{})
	}
	return ok
}

func checkNumeric(value interface{}, isInt bool) bool {
	ok := false
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

// public alteration methods...

// SetType sets the expected type for the property validator
func (pv *PropertyValidator) SetType(t JsonType) *PropertyValidator {
	pv.Type = t
	return pv
}

// SetNullable sets the property validator to allow nulls
func (pv *PropertyValidator) SetNullable() *PropertyValidator {
	pv.NotNull = false
	return pv
}

// SetNotNullable sets the property validator to disallow nulls
func (pv *PropertyValidator) SetNotNullable() *PropertyValidator {
	pv.NotNull = true
	return pv
}

// SetMandatory sets the property is mandatory (required) for the property validator
func (pv *PropertyValidator) SetMandatory() *PropertyValidator {
	pv.Mandatory = true
	return pv
}

// SetOptional sets the property is optional for the property validator
func (pv *PropertyValidator) SetOptional() *PropertyValidator {
	pv.Mandatory = false
	return pv
}

// SetRequired same as SetMandatory
func (pv *PropertyValidator) SetRequired() *PropertyValidator {
	pv.Mandatory = true
	return pv
}

// AddMandatoryWhens adds mandatory when condition token(s) to the property validator
func (pv *PropertyValidator) AddMandatoryWhens(c ...string) *PropertyValidator {
	pv.MandatoryWhen = append(pv.MandatoryWhen, c...)
	return pv
}

// AddConstraints adds constraint(s) to the property validator
func (pv *PropertyValidator) AddConstraints(c ...Constraint) *PropertyValidator {
	pv.Constraints = append(pv.Constraints, c...)
	return pv
}

// SetObjectValidator sets the object validator for the property validator
func (pv *PropertyValidator) SetObjectValidator(v *Validator) *PropertyValidator {
	pv.ObjectValidator = v
	return pv
}

// SetOrder sets the property check order for the property validator
func (pv *PropertyValidator) SetOrder(order int) *PropertyValidator {
	pv.Order = order
	return pv
}

// AddWhenConditions adds when condition token(s) to the property validator
func (pv *PropertyValidator) AddWhenConditions(c ...string) *PropertyValidator {
	pv.WhenConditions = append(pv.WhenConditions, c...)
	return pv
}

// AddUnwantedConditions adds when condition token(s) to the property validator
func (pv *PropertyValidator) AddUnwantedConditions(c ...string) *PropertyValidator {
	pv.UnwantedConditions = append(pv.UnwantedConditions, c...)
	return pv
}

// SetRequiredWith sets the required with expression for the property validator
func (pv *PropertyValidator) SetRequiredWith(expr OthersExpr) *PropertyValidator {
	pv.RequiredWith = expr
	return pv
}

// SetRequiredWithMessage sets the required with message for the property validator
func (pv *PropertyValidator) SetRequiredWithMessage(msg string) *PropertyValidator {
	pv.RequiredWithMessage = msg
	return pv
}

// SetUnwantedWith sets the unwanted with expression for the property validator
func (pv *PropertyValidator) SetUnwantedWith(expr OthersExpr) *PropertyValidator {
	pv.UnwantedWith = expr
	return pv
}

// SetUnwantedWithMessage sets the unwanted with message for the property validator
func (pv *PropertyValidator) SetUnwantedWithMessage(msg string) *PropertyValidator {
	pv.UnwantedWithMessage = msg
	return pv
}

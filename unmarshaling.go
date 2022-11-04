package valix

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	ptyNameIgnoreUnknownProperties = "ignoreUnknownProperties"
	ptyNameProperties              = "properties"
	ptyNameConstraints             = "constraints"
	ptyNameAllowArray              = "allowArray"
	ptyNameDisallowObject          = "disallowObject"
	ptyNameAllowNullJson           = "allowNullJson"
	ptyNameAllowNullItems          = "allowNullItems"
	ptyNameStopOnFirst             = "stopOnFirst"
	ptyNameUseNumber               = "useNumber"
	ptyNameOrderedPropertyChecks   = "orderedPropertyChecks"
	ptyNameWhenConditions          = "whenConditions"
	ptyNameOthersExpr              = "othersExpr"
	ptyNameMandatoryWhen           = "mandatoryWhen"
	ptyNameConditionalVariants     = "conditionalVariants"
	ptyNameOasInfo                 = "oasInfo"
	ptyNameType                    = "type"
	ptyNameMandatory               = "mandatory"
	ptyNameNotNull                 = "notNull"
	ptyNameOrder                   = "order"
	ptyNameUnwantedConditions      = "unwantedConditions"
	ptyNameRequiredWith            = "requiredWith"
	ptyNameRequiredWithMessage     = "requiredWithMessage"
	ptyNameUnwantedWith            = "unwantedWith"
	ptyNameUnwantedWithMessage     = "unwantedWithMessage"
	ptyNameObjectValidator         = "objectValidator"
	ptyNameName                    = "name"
	ptyNameFields                  = "fields"
)

const (
	errMsgUnknownNamedConstraint       = "unknown constraint '%s'"
	errMsgFieldExpectedType            = "field '%s' expected type %s"
	errMsgCannotParseExpr              = "cannot parse other expression '%s' - %s"
	errMsgConstraintExpectedObject     = "constraint [%d] expected to be an object"
	msgUnknownField                    = "Unknown field '%s'"
	msgConstraintNotStruct             = "constraint not a struct"
	msgPropertyValidatorNotNull        = "Property validator for property '%s' cannot be null"
	msgPropertyValidatorExpectedObject = "Property validator for property '%s' must be an object"
	msgConstraintNameString            = "Constraint name must be a string"
	msgUnknownConstraintName           = "Unknown constraint name '%s'"
	msgFieldExpectedType               = "Field '%s' expected type %s (found type %s)"
	msgFieldsUnmarshalable             = "Unable to unmarshal fields"
)

var ValidatorValidator *Validator
var PropertyValidatorValidator *Validator
var constraintValidator *Validator
var conditionalVariantValidator *Validator
var oasInfoValidator *Validator

func init() {
	oasInfoValidator = &Validator{
		IgnoreUnknownProperties: false,
		OrderedPropertyChecks:   true,
		Properties: Properties{
			ptyNameOasDescription: {
				Type: JsonString,
			},
			ptyNameOasTitle: {
				Type: JsonString,
			},
			ptyNameOasFormat: {
				Type: JsonString,
			},
			ptyNameOasExample: {
				Type: JsonString,
			},
			ptyNameOasDeprecated: {
				Type: JsonBoolean,
			},
		},
	}
	constraintValidator = &Validator{
		OrderedPropertyChecks: true,
		AllowArray:            true,
		DisallowObject:        true,
		Properties: Properties{
			ptyNameName: {
				Order:     0,
				Type:      JsonString,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&SetConditionFrom{Parent: true, Prefix: "constraint_name:"},
					NewCustomConstraint(constraintNameCheck, ""),
				},
			},
			ptyNameFields: {
				Order:     1,
				Type:      JsonObject,
				Mandatory: false,
				NotNull:   true,
				ObjectValidator: &Validator{
					IgnoreUnknownProperties: true, // don't know what the fields are without looking at the constraint
					Constraints: Constraints{
						NewCustomConstraint(constraintFieldsCheck, ""),
					},
				},
			},
			ptyNameWhenConditions: {
				Order:     2,
				Type:      JsonArray,
				Mandatory: false,
				NotNull:   false,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameOthersExpr: {
				Order:     3,
				Type:      JsonString,
				Mandatory: false,
				NotNull:   false,
				Constraints: Constraints{
					NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (passed bool, message string) {
						passed = false
						if expr, ok := value.(string); ok {
							if _, err := ParseExpression(expr); err == nil {
								passed = true
							}
						}
						return
					}, fmt.Sprintf("Must be a valid properties expression")),
				},
			},
		},
	}
	PropertyValidatorValidator = &Validator{
		IgnoreUnknownProperties: false,
		OrderedPropertyChecks:   true,
		Properties: Properties{
			ptyNameType: {
				Type:      JsonString,
				Mandatory: false,
				NotNull:   true,
				Constraints: Constraints{
					&StringValidToken{
						Tokens: []string{
							JsonAny.String(),
							JsonString.String(),
							JsonNumber.String(),
							JsonInteger.String(),
							JsonBoolean.String(),
							JsonObject.String(),
							JsonArray.String(),
						},
					},
				},
			},
			ptyNameNotNull: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameMandatory: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameMandatoryWhen: {
				Type:      JsonArray,
				Mandatory: false,
				NotNull:   true,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameConstraints: {
				Type:            JsonArray,
				Mandatory:       false,
				NotNull:         true,
				ObjectValidator: constraintValidator,
				Constraints: Constraints{
					&ArrayOf{Type: JsonObject.String(), AllowNullElement: false, Stop: true},
				},
			},
			ptyNameObjectValidator: {
				Type:            JsonObject,
				Mandatory:       false,
				NotNull:         false,
				ObjectValidator: ValidatorValidator,
			},
			ptyNameOrder: {
				Type:      JsonInteger,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameWhenConditions: {
				Type:      JsonArray,
				Mandatory: false,
				NotNull:   true,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameUnwantedConditions: {
				Type:      JsonArray,
				Mandatory: false,
				NotNull:   true,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameRequiredWith: {
				Type:      JsonString,
				Mandatory: false,
				NotNull:   false,
				Constraints: Constraints{
					&parseCheckPropertiesExpression{pty: ptyNameRequiredWith},
				},
			},
			ptyNameRequiredWithMessage: {
				Type:      JsonString,
				Mandatory: false,
				NotNull:   false,
			},
			ptyNameUnwantedWith: {
				Type:      JsonString,
				Mandatory: false,
				NotNull:   false,
				Constraints: Constraints{
					&parseCheckPropertiesExpression{pty: ptyNameUnwantedWith},
				},
			},
			ptyNameUnwantedWithMessage: {
				Type:      JsonString,
				Mandatory: false,
				NotNull:   false,
			},
			ptyNameOasInfo: {
				Type:            JsonObject,
				Mandatory:       false,
				NotNull:         false,
				ObjectValidator: oasInfoValidator,
			},
		},
	}
	conditionalVariantValidator = &Validator{
		IgnoreUnknownProperties: false,
		AllowArray:              true,
		DisallowObject:          true,
		Properties: Properties{
			ptyNameWhenConditions: {
				Type:      JsonArray,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameProperties: {
				Type:      JsonObject,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					NewCustomConstraint(validatorPropertiesCheck, ""),
				},
			},
			ptyNameConstraints: {
				Type:            JsonArray,
				Mandatory:       false,
				NotNull:         false,
				ObjectValidator: constraintValidator,
				Constraints: Constraints{
					&ArrayOf{Type: JsonObject.String(), AllowNullElement: false, Stop: true},
				},
			},
		},
	}
	ValidatorValidator = &Validator{
		IgnoreUnknownProperties: false,
		OrderedPropertyChecks:   true,
		Properties: Properties{
			ptyNameIgnoreUnknownProperties: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameProperties: {
				Type:      JsonObject,
				Mandatory: true,
				NotNull:   true,
				Constraints: Constraints{
					NewCustomConstraint(validatorPropertiesCheck, ""),
				},
			},
			ptyNameConstraints: {
				Type:            JsonArray,
				Mandatory:       false,
				NotNull:         false,
				ObjectValidator: constraintValidator,
				Constraints: Constraints{
					&ArrayOf{Type: JsonObject.String(), AllowNullElement: false, Stop: true},
				},
			},
			ptyNameAllowArray: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameDisallowObject: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameAllowNullJson: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameAllowNullItems: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameStopOnFirst: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameUseNumber: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameOrderedPropertyChecks: {
				Type:      JsonBoolean,
				Mandatory: false,
				NotNull:   true,
			},
			ptyNameWhenConditions: {
				Type:      JsonArray,
				Mandatory: false,
				NotNull:   false,
				Constraints: Constraints{
					&ArrayOf{Type: JsonString.String(), AllowNullElement: false},
				},
			},
			ptyNameConditionalVariants: {
				Type:            JsonArray,
				Mandatory:       false,
				NotNull:         true,
				ObjectValidator: conditionalVariantValidator,
			},
			ptyNameOasInfo: {
				Type:            JsonObject,
				Mandatory:       false,
				NotNull:         false,
				ObjectValidator: oasInfoValidator,
			},
		},
		AllowArray:     false,
		DisallowObject: false,
		AllowNullJson:  false,
	}
	// prevent initialisation loop...
	PropertyValidatorValidator.Properties[ptyNameObjectValidator].ObjectValidator = ValidatorValidator
}

type parseCheckPropertiesExpression struct {
	pty string
}

func (ck *parseCheckPropertiesExpression) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		_, _, err := parseExpression(str)
		if err != nil {
			return false, fmt.Sprintf(ck.GetMessage(vcx)+" - %s", err.Error())
		}
	}
	return true, ""
}

func (ck *parseCheckPropertiesExpression) GetMessage(tcx I18nContext) string {
	return fmt.Sprintf("Property '%s' must be a valid properties expression", ck.pty)
}

func validatorPropertiesCheck(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
	result := false
	if m, ok := value.(map[string]interface{}); ok {
		result = true
		for k, v := range m {
			vcx.pushPathProperty(k, v, nil)
			if v == nil {
				vcx.addTranslatedViolationForCurrent(vcx.TranslateFormat(msgPropertyValidatorNotNull, k))
			} else if mv, mvOk := v.(map[string]interface{}); mvOk {
				PropertyValidatorValidator.validate(mv, vcx)
			} else {
				vcx.addTranslatedViolationForCurrent(vcx.TranslateFormat(msgPropertyValidatorExpectedObject, k))
			}
			vcx.popPath()
		}
	}
	return result, "Error checking validator properties"
}

func constraintNameCheck(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
	result := false
	msg := msgConstraintNameString
	if name, ok := value.(string); ok {
		if name != constraintSetName && name != conditionalConstraintName && name != arrayconditionalConstraintName &&
			name != constraintVariableProperty && !constraintsRegistry.has(name) {
			msg = vcx.TranslateFormat(msgUnknownConstraintName, name)
		} else {
			result = true
		}
	}
	return result, msg
}

func constraintFieldsCheck(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
	if mv, mOk := value.(map[string]interface{}); mOk {
		parent, pOk := vcx.AncestorValue(0)
		if pOk {
			if parentConstraint, pcOk := parent.(map[string]interface{}); pcOk {
				constraintName := parentConstraint[ptyNameName].(string)
				var constraint Constraint
				constraintFound := false
				if constraintName == constraintSetName {
					constraint = &ConstraintSet{}
					constraintFound = true
				} else if constraintName == constraintVariableProperty {
					constraint = &VariablePropertyConstraint{}
					constraintFound = true
				} else {
					constraint, constraintFound = constraintsRegistry.get(constraintName)
					// don't need to report any violation here - as not found constraint is dealt with elsewhere
				}
				if constraintFound {
					fieldNames := getConstraintFieldNames(constraint)
					fieldsOk := true
					for k, v := range mv {
						vcx.pushPathProperty(k, v, nil)
						if !fieldNames[k] {
							fieldsOk = true
							vcx.addTranslatedViolationForCurrent(vcx.TranslateFormat(msgUnknownField, k))
						}
						vcx.popPath()
					}
					if fieldsOk {
						_, err := validateUnmarshallingConstraint(constraint, mv, nil, nil)
						if err != nil {
							msg := msgFieldsUnmarshalable
							if ute, is := err.(*json.UnmarshalTypeError); is {
								msg = vcx.TranslateFormat(msgFieldExpectedType, ute.Field, ute.Type.Name(), ute.Value)
							}
							vcx.addTranslatedViolationForCurrent(msg)
						}
					}
				}
			}
		}
	}
	return true, ""
}

func validateUnmarshallingConstraint(constraint Constraint, v map[string]interface{}, whens []string, othersExpr OthersExpr) (Constraint, error) {
	ty := reflect.TypeOf(constraint)
	if ty.Kind() != reflect.Ptr || ty.Elem().Kind() != reflect.Struct {
		return nil, errors.New(msgConstraintNotStruct)
	}
	ty = ty.Elem()
	newC := reflect.New(ty)
	intf := newC.Interface()
	data, _ := json.Marshal(v)
	err := json.Unmarshal(data, intf)
	if err != nil {
		return nil, err
	}
	if whens != nil || othersExpr != nil {
		return &ConditionalConstraint{
			When:       whens,
			Others:     othersExpr,
			Constraint: newC.Interface().(Constraint),
		}, nil
	}
	return newC.Interface().(Constraint), nil
}

func getConstraintFieldNames(c Constraint) map[string]bool {
	result := map[string]bool{}
	ty := reflect.TypeOf(c)
	if ty.Kind() != reflect.Ptr || ty.Elem().Kind() != reflect.Struct {
		// if it isn't a pointer to a struct we can't get fields...
		return result
	}
	ty = ty.Elem()
	count := ty.NumField()
	for f := 0; f < count; f++ {
		fld := ty.Field(f)
		fldName := fld.Name
		if tag, ok := fld.Tag.Lookup(tagNameJson); ok {
			if cAt := strings.Index(tag, ","); cAt != -1 {
				fldName = tag[0:cAt]
			} else {
				fldName = tag
			}
		}
		result[fldName] = true
	}
	return result
}

func (c *ConstraintSet) UnmarshalJSON(data []byte) error {
	c.Constraints = Constraints{}
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	if raw, ok := obj[constraintPtyNameMessage]; ok {
		if v, ok := raw.(string); ok {
			c.Message = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintPtyNameMessage, "string")
		}
	}
	if raw, ok := obj[constraintPtyNameStop]; ok {
		if v, ok := raw.(bool); ok {
			c.Stop = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintPtyNameStop, "bool")
		}
	}
	if raw, ok := obj[constraintSetFieldOneOf]; ok {
		if v, ok := raw.(bool); ok {
			c.Stop = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintSetFieldOneOf, "bool")
		}
	}
	if raw, ok := obj[constraintSetFieldConstraints]; ok {
		if v, ok := raw.([]interface{}); ok {
			if cs, err := unmarshalConstraints(v); err != nil {
				return err
			} else {
				c.Constraints = cs
			}
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintSetFieldConstraints, "array")
		}
	}
	return nil
}

func (cs *Constraints) UnmarshalJSON(data []byte) error {
	arr := make([]interface{}, 0)
	err := json.Unmarshal(data, &arr)
	if err != nil {
		return err
	}
	constraints, err := unmarshalConstraints(arr)
	if err != nil {
		return err
	}
	*cs = constraints
	return nil
}

func unmarshalConstraints(cs []interface{}) (Constraints, error) {
	result := Constraints{}
	for i, c := range cs {
		if v, ok := c.(map[string]interface{}); ok {
			constraint, err := unmarshalConstraint(v)
			if err != nil {
				return nil, err
			}
			result = append(result, constraint)
		} else {
			return nil, fmt.Errorf(errMsgConstraintExpectedObject, i)
		}
	}
	return result, nil
}

func unmarshalConstraint(c map[string]interface{}) (Constraint, error) {
	constraint, err := unmarshalGetConstraint(c)
	if err != nil {
		return nil, err
	}
	var whens []string
	var othersExpr OthersExpr
	if rawWhens, ok := c[ptyNameWhenConditions]; ok {
		if slc, ok := rawWhens.([]interface{}); ok {
			whens = make([]string, len(slc))
			for i, v := range slc {
				if str, ok := v.(string); ok {
					whens[i] = str
				} else {
					return nil, fmt.Errorf(errMsgFieldExpectedType, ptyNameWhenConditions, "array of strings")
				}
			}
		} else {
			return nil, fmt.Errorf(errMsgFieldExpectedType, ptyNameWhenConditions, "array")
		}
	}
	if rawExpr, ok := c[ptyNameOthersExpr]; ok {
		if cExpr, ok := rawExpr.(string); ok {
			if expr, err := ParseExpression(cExpr); err == nil {
				othersExpr = expr
			} else {
				return nil, fmt.Errorf(errMsgCannotParseExpr, cExpr, err.Error())
			}
		} else {
			return nil, fmt.Errorf(errMsgFieldExpectedType, ptyNameOthersExpr, "string")
		}
	}
	if fs, ok := c[ptyNameFields]; ok && fs != nil {
		if fields, fOk := fs.(map[string]interface{}); fOk {
			return validateUnmarshallingConstraint(constraint, fields, whens, othersExpr)
		} else {
			return nil, fmt.Errorf(errMsgFieldExpectedType, ptyNameFields, "object")
		}
	} else {
		return validateUnmarshallingConstraint(constraint, map[string]interface{}{}, whens, othersExpr)
	}
}

func unmarshalGetConstraint(c map[string]interface{}) (Constraint, error) {
	constraintName := ""
	if raw, ok := c[ptyNameName]; ok {
		if v, ok := raw.(string); ok {
			constraintName = v
		} else {
			return nil, fmt.Errorf(errMsgFieldExpectedType, ptyNameName, "string")
		}
	}
	var constraint Constraint
	if constraintName == constraintSetName {
		constraint = &ConstraintSet{}
	} else if constraintName == constraintVariableProperty {
		constraint = &VariablePropertyConstraint{}
	} else if constraintName == arrayconditionalConstraintName {
		constraint = &ArrayConditionalConstraint{}
	} else {
		constraint, _ = constraintsRegistry.get(constraintName)
	}
	if constraint == nil {
		return nil, fmt.Errorf(errMsgUnknownNamedConstraint, constraintName)
	}
	return constraint, nil
}

func (jt *JsonType) UnmarshalJSON(data []byte) error {
	str := string(data[:])
	if strings.HasPrefix(str, `"`) && strings.HasPrefix(str, `"`) {
		if v, ok := JsonTypeFromString(str[1 : len(str)-1]); ok {
			*jt = v
			return nil
		}
	} else if i, err := strconv.ParseInt(str, 10, 32); err == nil {
		jtv := JsonType(i)
		if jtv >= JsonAny && jtv <= JsonArray {
			*jt = jtv
			return nil
		}
	}
	return fmt.Errorf("value for JsonType expected string (%s) or integer (%d to %d)", jsonTypeTokensList, JsonAny, JsonArray)
}

func (c *StringPattern) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	err := json.Unmarshal(data, &obj)
	if err != nil {
		return err
	}
	if raw, ok := obj[constraintPtyNameMessage]; ok {
		if v, ok := raw.(string); ok {
			c.Message = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintPtyNameMessage, "string")
		}
	}
	if raw, ok := obj[constraintPtyNameStop]; ok {
		if v, ok := raw.(bool); ok {
			c.Stop = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, constraintPtyNameStop, "bool")
		}
	}
	if raw, ok := obj["Regexp"]; ok {
		if v, ok := raw.(string); ok {
			rx, err := regexp.Compile(v)
			if err != nil {
				return err
			}
			c.Regexp = *rx
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Regexp", "regexp string")
		}
	}
	return nil
}

func (c *ArrayConditionalConstraint) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	_ = json.Unmarshal(data, &obj)
	if raw, ok := obj["When"]; ok {
		if v, ok := raw.(string); ok {
			c.When = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "When", "string")
		}
	}
	if raw, ok := obj["Ancestry"]; ok {
		if v, ok := raw.(float64); ok {
			c.Ancestry = uint(v)
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Ancestry", "int")
		}
	}
	if raw, ok := obj["Constraint"]; ok && raw != nil {
		if v, ok := raw.(map[string]interface{}); ok {
			if wrapped, err := unmarshalConstraint(v); err == nil {
				c.Constraint = wrapped
			} else {
				return err
			}
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Constraint", "object")
		}
	}
	return nil
}

func (o *OthersExpr) UnmarshalJSON(data []byte) error {
	expr := ""
	err := json.Unmarshal(data, &expr)
	if err != nil {
		return err
	}
	newO, err := ParseExpression(expr)
	if err != nil {
		return fmt.Errorf("unable to parse properties expression - %w", err)
	}
	*o = newO
	return nil
}

func (c *SetConditionIf) UnmarshalJSON(data []byte) error {
	obj := map[string]interface{}{}
	_ = json.Unmarshal(data, &obj)
	if raw, ok := obj["Constraint"]; ok && raw != nil {
		if v, ok := raw.(map[string]interface{}); ok {
			if wrapped, err := unmarshalConstraint(v); err == nil {
				c.Constraint = wrapped
			} else {
				return err
			}
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Constraint", "object")
		}
	}
	if raw, ok := obj["SetOk"]; ok {
		if v, ok := raw.(string); ok {
			c.SetOk = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "SetOk", "string")
		}
	}
	if raw, ok := obj["SetFail"]; ok {
		if v, ok := raw.(string); ok {
			c.SetFail = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "SetFail", "string")
		}
	}
	if raw, ok := obj["Parent"]; ok {
		if v, ok := raw.(bool); ok {
			c.Parent = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Parent", "bool")
		}
	}
	if raw, ok := obj["Global"]; ok {
		if v, ok := raw.(bool); ok {
			c.Global = v
		} else {
			return fmt.Errorf(errMsgFieldExpectedType, "Global", "bool")
		}
	}
	return nil
}

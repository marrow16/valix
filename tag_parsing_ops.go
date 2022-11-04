package valix

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	tagTokenNotNull            = "notNull"
	tagTokenNullable           = "nullable"
	tagTokenMandatory          = "mandatory"
	tagTokenRequired           = "required"
	tagTokenOptional           = "optional"
	tagTokenOnly               = "only"
	tagTokenOnlyMsg            = "only_msg"
	tagTokenType               = "type"
	tagTokenConstraint         = "constraint"
	tagTokenConstraints        = "constraints"
	tagTokenConstraintsPrefix  = tagTokenConstraints + ":"
	tagTokenOrder              = "order"
	tagTokenWhen               = "when"
	tagTokenUnwanted           = "unwanted"
	tagTokenRequiredWith       = "required_with"
	tagTokenRequiredWithAlt    = "+"
	tagTokenRequiredWithMsg    = "required_with_msg"
	tagTokenRequiredWithAltMsg = "+msg"
	tagTokenUnwantedWith       = "unwanted_with"
	tagTokenUnwantedWithAlt    = "-"
	tagTokenUnwantedWithMsg    = "unwanted_with_msg"
	tagTokenUnwantedWithAltMsg = "-msg"
	tagTokenStopOnFirst        = "stop_on_first"
	tagTokenStopOnFirstAlt     = "stop1st"
	// object level tag items...
	tagTokenObjPrefix                  = "obj."
	tagTokenObjIgnoreUnknownProperties = tagTokenObjPrefix + "ignoreUnknownProperties"
	tagTokenObjUnknownProperties       = tagTokenObjPrefix + "unknownProperties" // true/false
	tagTokenObjConstraint              = tagTokenObjPrefix + tagTokenConstraint
	tagTokenObjOrdered                 = tagTokenObjPrefix + "ordered"
	tagTokenObjWhen                    = tagTokenObjPrefix + tagTokenWhen
	tagTokenObjNo                      = tagTokenObjPrefix + "no"
	// array level tag items...
	tagTokenArrPrefix         = "arr."
	tagTokenArrAllowNullItems = tagTokenArrPrefix + "allowNulls"
)

var tagExpectsColon = map[string]bool{
	tagTokenNotNull:  false,
	tagTokenNullable: false,
	// tagTokenMandatory, tagTokenRequired: either
	tagTokenOptional:       false,
	tagTokenStopOnFirst:    false,
	tagTokenStopOnFirstAlt: false,
	// tagTokenOnly: either
	tagTokenOnlyMsg:                    true,
	tagTokenType:                       true,
	tagTokenOrder:                      true,
	tagTokenConstraint:                 true,
	tagTokenObjConstraint:              true,
	tagTokenWhen:                       true,
	tagTokenUnwanted:                   true,
	tagTokenRequiredWith:               true,
	tagTokenRequiredWithAlt:            true,
	tagTokenUnwantedWith:               true,
	tagTokenUnwantedWithAlt:            true,
	tagTokenRequiredWithMsg:            true,
	tagTokenRequiredWithAltMsg:         true,
	tagTokenUnwantedWithMsg:            true,
	tagTokenUnwantedWithAltMsg:         true,
	tagTokenObjIgnoreUnknownProperties: false,
	tagTokenObjUnknownProperties:       true,
	tagTokenObjOrdered:                 false,
	tagTokenObjWhen:                    true,
	tagTokenObjNo:                      false,
	tagTokenArrAllowNullItems:          false,
}

type tagTokenOperation func(pv *PropertyValidator, hasColon bool, tagValue string) error

var tagOpRequiredWith = func(pv *PropertyValidator, hasColon bool, tagValue string) error {
	if expr, err := ParseExpression(tagValue); err != nil {
		return err
	} else if pv.RequiredWith == nil {
		pv.RequiredWith = expr
	} else {
		for _, x := range expr {
			pv.RequiredWith = append(pv.RequiredWith, x)
		}
	}
	return nil
}

var tagOpUnwantedWith = func(pv *PropertyValidator, hasColon bool, tagValue string) error {
	if expr, err := ParseExpression(tagValue); err != nil {
		return err
	} else if pv.UnwantedWith == nil {
		pv.UnwantedWith = expr
	} else {
		for _, x := range expr {
			pv.UnwantedWith = append(pv.UnwantedWith, x)
		}
	}
	return nil
}

var tagOpRequiredWithMsg = func(pv *PropertyValidator, hasColon bool, tagValue string) error {
	if unq, ok := isQuotedStr(tagValue); ok {
		pv.RequiredWithMessage = unq
	} else {
		pv.RequiredWithMessage = tagValue
	}
	return nil
}

var tagOpUnwantedWithMsg = func(pv *PropertyValidator, hasColon bool, tagValue string) error {
	if unq, ok := isQuotedStr(tagValue); ok {
		pv.UnwantedWithMessage = unq
	} else {
		pv.UnwantedWithMessage = tagValue
	}
	return nil
}

var tagTokenOperations = map[string]tagTokenOperation{
	tagTokenNotNull: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.NotNull = true
		return nil
	},
	tagTokenNullable: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.NotNull = false
		return nil
	},
	tagTokenMandatory: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		return pv.setTagMandatoryWhen(hasColon, tagValue)
	},
	tagTokenRequired: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		return pv.setTagMandatoryWhen(hasColon, tagValue)
	},
	tagTokenOptional: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.Mandatory = false
		return nil
	},
	tagTokenStopOnFirst: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.StopOnFirst = true
		return nil
	},
	tagTokenStopOnFirstAlt: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.StopOnFirst = true
		return nil
	},
	tagTokenOnly: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.Only = true
		if hasColon {
			return addConditions(&pv.OnlyConditions, tagValue, true)
		}
		return nil
	},
	tagTokenOnlyMsg: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if unq, ok := isQuotedStr(tagValue); ok {
			pv.OnlyMessage = unq
		} else {
			pv.OnlyMessage = tagValue
		}
		return nil
	},
	tagTokenType: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		ty, ok := JsonTypeFromString(tagValue)
		if !ok {
			return fmt.Errorf(msgUnknownPropertyType, tagValue)
		}
		pv.Type = ty
		return nil
	},
	tagTokenOrder: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		v, err := strconv.ParseInt(tagValue, 10, 32)
		if err != nil {
			return fmt.Errorf(msgUnknownTagValue, tagTokenOrder, "int", tagValue)
		}
		pv.Order = int(v)
		return nil
	},
	tagTokenConstraint: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		return pv.addConstraint(tagValue)
	},
	tagTokenObjConstraint: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenObjConstraint)
		}
		return pv.ObjectValidator.addConstraint(tagValue)
	},
	tagTokenWhen: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		return addConditions(&pv.WhenConditions, tagValue, true)
	},
	tagTokenUnwanted: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		return addConditions(&pv.UnwantedConditions, tagValue, true)
	},
	tagTokenRequiredWith:       tagOpRequiredWith,
	tagTokenRequiredWithAlt:    tagOpRequiredWith,
	tagTokenUnwantedWith:       tagOpUnwantedWith,
	tagTokenUnwantedWithAlt:    tagOpUnwantedWith,
	tagTokenRequiredWithMsg:    tagOpRequiredWithMsg,
	tagTokenRequiredWithAltMsg: tagOpRequiredWithMsg,
	tagTokenUnwantedWithMsg:    tagOpUnwantedWithMsg,
	tagTokenUnwantedWithAltMsg: tagOpUnwantedWithMsg,
	tagTokenObjIgnoreUnknownProperties: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenObjIgnoreUnknownProperties)
		}
		pv.ObjectValidator.IgnoreUnknownProperties = true
		return nil
	},
	tagTokenObjUnknownProperties: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenObjUnknownProperties)
		}
		b, err := strconv.ParseBool(tagValue)
		if err != nil {
			return fmt.Errorf(msgUnknownTagValue, tagTokenObjUnknownProperties, "boolean", tagValue)
		}
		pv.ObjectValidator.IgnoreUnknownProperties = b
		return nil
	},
	tagTokenObjOrdered: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenObjOrdered)
		}
		pv.ObjectValidator.OrderedPropertyChecks = true
		return nil
	},
	tagTokenObjWhen: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenObjWhen)
		}
		if isBracedStr(tagValue, true) {
			if tokens, err := parseCommas(tagValue[1 : len(tagValue)-1]); err == nil {
				for _, token := range tokens {
					if unq, ok := isQuotedStr(token); ok {
						pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, unq)
					} else {
						pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, token)
					}
				}
			} else {
				return err
			}
		} else if unq, ok := isQuotedStr(tagValue); ok {
			pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, unq)
		} else {
			pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, tagValue)
		}
		return nil
	},
	tagTokenObjNo: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		pv.ObjectValidator = nil
		return nil
	},
	tagTokenArrAllowNullItems: func(pv *PropertyValidator, hasColon bool, tagValue string) error {
		if pv.ObjectValidator == nil {
			return fmt.Errorf(msgPropertyNotObject, tagTokenArrAllowNullItems)
		}
		pv.ObjectValidator.AllowNullItems = true
		return nil
	},
}

func isQuotedStr(str string) (string, bool) {
	if strings.HasPrefix(str, `"`) && strings.HasSuffix(str, `"`) {
		return strings.ReplaceAll(str[1:len(str)-1], `""`, `"`), true
	} else if strings.HasPrefix(str, `'`) && strings.HasSuffix(str, `'`) {
		return strings.ReplaceAll(str[1:len(str)-1], `''`, `'`), true
	}
	return "", false
}

func isBracedStr(str string, allowCurly bool) bool {
	return (strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) ||
		(allowCurly && strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}"))
}

package valix

import (
	"encoding/json"
	"reflect"

	"golang.org/x/text/unicode/norm"
)

func (v *Validator) MarshalJSON() ([]byte, error) {
	j, err := v.toJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(j)
}

func (v *Validator) toJSON() (map[string]interface{}, error) {
	properties := make(map[string]interface{}, len(v.Properties))
	for k, pv := range v.Properties {
		pvj, err := pv.toJSON()
		if err != nil {
			return nil, err
		}
		properties[k] = pvj
	}
	result := map[string]interface{}{
		ptyNameIgnoreUnknownProperties: v.IgnoreUnknownProperties,
		ptyNameAllowArray:              v.AllowArray,
		ptyNameDisallowObject:          v.DisallowObject,
		ptyNameAllowNullJson:           v.AllowNullJson,
		ptyNameAllowNullItems:          v.AllowNullItems,
		ptyNameStopOnFirst:             v.StopOnFirst,
		ptyNameUseNumber:               v.UseNumber,
		ptyNameOrderedPropertyChecks:   v.OrderedPropertyChecks,
		ptyNameProperties:              properties,
	}
	if len(v.Constraints) > 0 {
		cs, err := constraintsToJson(v.Constraints)
		if err != nil {
			return nil, err
		}
		result[ptyNameConstraints] = cs
	}
	if len(v.WhenConditions) > 0 {
		arr := make([]string, len(v.WhenConditions))
		for i, s := range v.WhenConditions {
			arr[i] = s
		}
		result[ptyNameWhenConditions] = arr
	}
	if len(v.ConditionalVariants) > 0 {
		if cvs, err := v.ConditionalVariants.toJson(); err != nil {
			return nil, err
		} else {
			result[ptyNameConditionalVariants] = cvs
		}
	}
	if v.OasInfo != nil {
		result[ptyNameOasInfo] = v.OasInfo.toJson()
	}
	return result, nil
}

func (pv *PropertyValidator) MarshalJSON() ([]byte, error) {
	j, err := pv.toJSON()
	if err != nil {
		return nil, err
	}
	return json.Marshal(j)
}

func (pv *PropertyValidator) toJSON() (map[string]interface{}, error) {
	result := map[string]interface{}{
		ptyNameType:                pv.Type.String(),
		ptyNameMandatory:           pv.Mandatory,
		ptyNameNotNull:             pv.NotNull,
		ptyNameOrder:               pv.Order,
		ptyNameRequiredWithMessage: pv.RequiredWithMessage,
		ptyNameUnwantedWithMessage: pv.UnwantedWithMessage,
	}
	if len(pv.Constraints) > 0 {
		cs, err := constraintsToJson(pv.Constraints)
		if err != nil {
			return nil, err
		}
		result[ptyNameConstraints] = cs
	}
	if len(pv.MandatoryWhen) > 0 {
		arr := make([]string, len(pv.MandatoryWhen))
		for i, s := range pv.MandatoryWhen {
			arr[i] = s
		}
		result[ptyNameMandatoryWhen] = arr
	}
	if len(pv.WhenConditions) > 0 {
		arr := make([]string, len(pv.WhenConditions))
		for i, s := range pv.WhenConditions {
			arr[i] = s
		}
		result[ptyNameWhenConditions] = arr
	}
	if len(pv.UnwantedConditions) > 0 {
		arr := make([]string, len(pv.UnwantedConditions))
		for i, s := range pv.UnwantedConditions {
			arr[i] = s
		}
		result[ptyNameUnwantedConditions] = arr
	}
	if pv.RequiredWith != nil {
		result[ptyNameRequiredWith] = pv.RequiredWith.String()
	}
	if pv.UnwantedWith != nil {
		result[ptyNameUnwantedWith] = pv.UnwantedWith.String()
	}
	if pv.ObjectValidator != nil {
		ov, err := pv.ObjectValidator.toJSON()
		if err != nil {
			return nil, err
		}
		result[ptyNameObjectValidator] = ov
	}
	if pv.OasInfo != nil {
		result[ptyNameOasInfo] = pv.OasInfo.toJson()
	}
	return result, nil
}

func (cvs *ConditionalVariants) toJson() ([]interface{}, error) {
	result := make([]interface{}, 0, len(*cvs))
	for _, cv := range *cvs {
		if cvj, err := cv.toJson(); err != nil {
			return nil, err
		} else {
			result = append(result, cvj)
		}
	}
	return result, nil
}

func (cv *ConditionalVariant) toJson() (map[string]interface{}, error) {
	result := map[string]interface{}{}
	arr := make([]string, len(cv.WhenConditions))
	for i, s := range cv.WhenConditions {
		arr[i] = s
	}
	result[ptyNameWhenConditions] = arr
	if len(cv.Constraints) > 0 {
		cs, err := constraintsToJson(cv.Constraints)
		if err != nil {
			return nil, err
		}
		result[ptyNameConstraints] = cs
	}
	properties := make(map[string]interface{}, len(cv.Properties))
	for k, pv := range cv.Properties {
		pvj, err := pv.toJSON()
		if err != nil {
			return nil, err
		}
		properties[k] = pvj
	}
	result[ptyNameProperties] = properties
	return result, nil
}

func (cs Constraints) MarshalJSON() ([]byte, error) {
	constraints := make([]interface{}, len(cs))
	for i, v := range cs {
		if m, err := constraintToJson(v); err != nil {
			return nil, err
		} else {
			constraints[i] = m
		}
	}
	return json.Marshal(constraints)
}

func constraintsToJson(constraints Constraints) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, len(constraints))
	for i, c := range constraints {
		cj, err := constraintToJson(c)
		if err != nil {
			return nil, err
		}
		result[i] = cj
	}
	return result, nil
}

func constraintToJson(constraint Constraint) (map[string]interface{}, error) {
	constraintName := reflect.TypeOf(constraint).Elem().Name()
	useConstraint := constraint
	isConditional := false
	var conditions []string
	othersExpr := ""
	if conditional, ok := constraint.(*ConditionalConstraint); ok {
		isConditional = true
		useConstraint = conditional.Constraint
		if useConstraint != nil {
			constraintName = reflect.TypeOf(useConstraint).Elem().Name()
		}
		conditions = make([]string, len(conditional.When))
		for i, s := range conditional.When {
			conditions[i] = s
		}
		if conditional.Others != nil {
			othersExpr = conditional.Others.String()
		}
	}
	by, err := json.Marshal(useConstraint)
	if err != nil {
		return nil, err
	}
	args := make(map[string]interface{})
	_ = json.Unmarshal(by, &args)
	result := map[string]interface{}{
		ptyNameName:   constraintName,
		ptyNameFields: args,
	}
	if isConditional {
		result[ptyNameWhenConditions] = conditions
		result[ptyNameOthersExpr] = othersExpr
	}
	return result, nil
}

func (oas *OasInfo) toJson() map[string]interface{} {
	return map[string]interface{}{
		ptyNameOasDescription: oas.Description,
		ptyNameOasTitle:       oas.Title,
		ptyNameOasFormat:      oas.Format,
		ptyNameOasExample:     oas.Example,
		ptyNameOasDeprecated:  oas.Deprecated,
	}
}

func (c *ConstraintSet) MarshalJSON() ([]byte, error) {
	constraints := make([]interface{}, len(c.Constraints))
	for i, v := range c.Constraints {
		m, err := constraintToJson(v)
		if err != nil {
			return nil, err
		}
		constraints[i] = m
	}
	j := map[string]interface{}{
		constraintSetFieldConstraints: constraints,
		constraintSetFieldOneOf:       c.OneOf,
		constraintSetFieldMessage:     c.Message,
		constraintSetFieldStop:        c.Stop,
	}
	return json.Marshal(j)
}

func (c *StringPattern) MarshalJSON() ([]byte, error) {
	j := map[string]interface{}{
		"Regexp":  c.Regexp.String(),
		"Message": c.Message,
		"Stop":    c.Stop,
	}
	return json.Marshal(j)
}

func (c *StringValidUnicodeNormalization) MarshalJSON() ([]byte, error) {
	strForm := ""
	switch c.Form {
	case norm.NFC:
		strForm = "NFC"
	case norm.NFD:
		strForm = "NFD"
	case norm.NFKC:
		strForm = "NFKC"
	case norm.NFKD:
		strForm = "NFKD"
	}
	j := map[string]interface{}{
		"Form":    strForm,
		"Message": c.Message,
		"Stop":    c.Stop,
	}
	return json.Marshal(j)
}

func (c *ArrayConditionalConstraint) MarshalJSON() ([]byte, error) {
	j := map[string]interface{}{
		"When":       c.When,
		"Ancestry":   c.Ancestry,
		"Constraint": nil,
	}
	if c.Constraint != nil {
		if cj, err := constraintToJson(c.Constraint); err == nil {
			j["Constraint"] = cj
		} else {
			return nil, err
		}
	}
	return json.Marshal(j)
}

func (c *SetConditionIf) MarshalJSON() ([]byte, error) {
	j := map[string]interface{}{
		"SetOk":      c.SetOk,
		"SetFail":    c.SetFail,
		"Parent":     c.Parent,
		"Global":     c.Global,
		"Constraint": nil,
	}
	if c.Constraint != nil {
		if cj, err := constraintToJson(c.Constraint); err == nil {
			j["Constraint"] = cj
		} else {
			return nil, err
		}
	}
	return json.Marshal(j)
}

package valix

func (v *Validator) Clone() *Validator {
	return &Validator{
		IgnoreUnknownProperties: v.IgnoreUnknownProperties,
		Properties:              v.Properties.Clone(),
		Constraints:             v.Constraints.Clone(),
		AllowArray:              v.AllowArray,
		DisallowObject:          v.DisallowObject,
		AllowNullJson:           v.AllowNullJson,
		StopOnFirst:             v.StopOnFirst,
		UseNumber:               v.UseNumber,
		OrderedPropertyChecks:   v.OrderedPropertyChecks,
		WhenConditions:          v.WhenConditions.Clone(),
		ConditionalVariants:     v.ConditionalVariants.Clone(),
		OasInfo:                 cloneOasInfo(v.OasInfo),
	}
}

func cloneValidator(src *Validator) *Validator {
	if src == nil {
		return nil
	}
	return src.Clone()
}

func (pv *PropertyValidator) Clone() *PropertyValidator {
	return &PropertyValidator{
		Type:                pv.Type,
		NotNull:             pv.NotNull,
		Mandatory:           pv.Mandatory,
		MandatoryWhen:       pv.MandatoryWhen.Clone(),
		Constraints:         pv.Constraints.Clone(),
		ObjectValidator:     cloneValidator(pv.ObjectValidator),
		Order:               pv.Order,
		WhenConditions:      pv.WhenConditions.Clone(),
		UnwantedConditions:  pv.UnwantedConditions.Clone(),
		RequiredWith:        pv.RequiredWith.Clone(),
		RequiredWithMessage: pv.RequiredWithMessage,
		UnwantedWith:        pv.UnwantedWith.Clone(),
		UnwantedWithMessage: pv.UnwantedWithMessage,
		OasInfo:             cloneOasInfo(pv.OasInfo),
	}
}

func (src Properties) Clone() Properties {
	if src == nil {
		return nil
	}
	result := Properties{}
	for k, v := range src {
		if v == nil {
			result[k] = nil
		} else {
			result[k] = v.Clone()
		}
	}
	return result
}

func (src ConditionalVariants) Clone() ConditionalVariants {
	if src == nil {
		return nil
	}
	result := make(ConditionalVariants, 0, len(src))
	for _, v := range src {
		result = append(result, cloneConditionalVariant(v))
	}
	return result
}

func cloneConditionalVariant(src *ConditionalVariant) *ConditionalVariant {
	if src == nil {
		return nil
	}
	return src.Clone()
}

func (src *ConditionalVariant) Clone() *ConditionalVariant {
	return &ConditionalVariant{
		WhenConditions:      src.WhenConditions.Clone(),
		Constraints:         src.Constraints.Clone(),
		Properties:          src.Properties.Clone(),
		ConditionalVariants: src.ConditionalVariants.Clone(),
	}
}

func (src Conditions) Clone() Conditions {
	if src == nil {
		return nil
	}
	result := make(Conditions, len(src))
	copy(result, src)
	return result
}

func (src Constraints) Clone() Constraints {
	if src == nil {
		return nil
	}
	result := make(Constraints, len(src))
	copy(result, src)
	return result
}

func (src OthersExpr) Clone() OthersExpr {
	if src == nil {
		return nil
	}
	result := make(OthersExpr, len(src))
	copy(result, src)
	return result
}

func cloneOasInfo(src *OasInfo) *OasInfo {
	if src == nil {
		return nil
	}
	return src.Clone()
}

func (oas *OasInfo) Clone() *OasInfo {
	return &OasInfo{
		Description: oas.Description,
		Title:       oas.Title,
		Format:      oas.Format,
		Example:     oas.Example,
		Deprecated:  oas.Deprecated,
	}
}

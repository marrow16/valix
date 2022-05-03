package valix

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagTokenNotNull            = "notNull"
	tagTokenNullable           = "nullable"
	tagTokenMandatory          = "mandatory"
	tagTokenRequired           = "required"
	tagTokenOptional           = "optional"
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
	// object level tag items...
	tagTokenObjPrefix                  = "obj."
	tagTokenObjIgnoreUnknownProperties = tagTokenObjPrefix + "ignoreUnknownProperties"
	tagTokenObjUnknownProperties       = tagTokenObjPrefix + "unknownProperties" // true/false
	tagTokenObjConstraint              = tagTokenObjPrefix + tagTokenConstraint
	tagTokenObjOrdered                 = tagTokenObjPrefix + "ordered"
	tagTokenObjWhen                    = tagTokenObjPrefix + tagTokenWhen
)

const (
	msgV8nPrefix                    = "tag " + tagNameV8n + " - "
	msgUnknownPropertyType          = msgV8nPrefix + "unknown property type '%s'"
	msgUnknownTokenInTag            = msgV8nPrefix + "unknown token '%s'"
	msgUnexpectedColon              = msgV8nPrefix + "unexpected ':' colon after token '%s'"
	msgExpectedColon                = msgV8nPrefix + "expected ':' colon after token '%s'"
	msgConstraintsFormat            = msgV8nPrefix + "must specify constraints in the format '&name{}' (found \"%s\")"
	msgConditionalConstraintsFormat = msgV8nPrefix + "must specify conditional constraints in the format '&[token,...]name{}' (found \"%s\")"
	msgUnknownConstraint            = msgV8nPrefix + "contains unknown constraint '%s'"
	msgCannotCreateConstraint       = msgV8nPrefix + "cannot create constraint '%s{}' (on non-struct constraint)"
	msgConstraintFieldUnknown       = msgV8nPrefix + "constraint '%s{}' field '%s' is unknown or not assignable"
	msgConstraintFieldInvalidValue  = msgV8nPrefix + "constraint '%s{}' field '%s' cannot be assigned with value specified"
	msgConstraintFieldNotExported   = msgV8nPrefix + "constraint '%s{}' has unexported field '%s' - so no fields can be specified as args"
	msgConstraintArgsParseError     = msgV8nPrefix + "constraint '%s{}' - args parsing error (%s)"
	msgTagFldMissingColon           = msgV8nPrefix + "constraint '%s{}` field missing ':' separator"
	msgConstraintNoDefaultValue     = msgV8nPrefix + "constraint '%s{}` does not have tagged default field"
	msgUnknownTagValue              = msgV8nPrefix + "token '%s' expected %s value (found \"%s\")"
	msgPropertyNotObject            = msgV8nPrefix + "token '%s' cannot be used on non object/array field"
	msgUnclosed                     = "unclosed parenthesis or quote started at position %d"
	msgUnopened                     = "unopened parenthesis at position %d"
	msgWrapped                      = "field '%s' (property '%s') - %s"
)

var (
	// grab some type information (for types used in common constraints)
	rx         = regexp.MustCompile("([a])")
	regexpKind = reflect.TypeOf(rx).Elem().Kind()
)

func (pv *PropertyValidator) processV8nTag(fieldName string, propertyName string, fld reflect.StructField) error {
	if tag, ok := fld.Tag.Lookup(tagNameV8n); ok {
		return pv.processV8nTagValue(fieldName, propertyName, tag)
	}
	return nil
}

func (pv *PropertyValidator) processV8nTagValue(fieldName string, propertyName string, tagValue string) error {
	tagItems, err := parseCommas(tagValue)
	if err != nil {
		return fmt.Errorf(msgWrapped, fieldName, propertyName, err.Error())
	}
	tagItems, err = tagAliasesRepo.resolve(tagItems)
	if err != nil {
		return fmt.Errorf(msgWrapped, fieldName, propertyName, err.Error())
	}
	if err := pv.processTagItems(fieldName, propertyName, tagItems); err != nil {
		return fmt.Errorf(msgWrapped, fieldName, propertyName, err.Error())
	}
	return nil
}

func (pv *PropertyValidator) processTagItems(fieldName string, propertyName string, tagItems []string) error {
	for _, tagItem := range tagItems {
		if tagItem != "" {
			cs, is, err := isConstraintsList(tagItem)
			if err != nil {
				return err
			} else if is {
				for _, c := range cs {
					if e2 := pv.addConstraint(c); e2 != nil {
						return e2
					}
				}
			} else if e3 := pv.addTagItem(fieldName, propertyName, tagItem); e3 != nil {
				return e3
			}
		}
	}
	return nil
}

func isConstraintsList(tagItem string) ([]string, bool, error) {
	if strings.HasPrefix(tagItem, tagTokenConstraintsPrefix) {
		value := strings.Trim(tagItem[len(tagTokenConstraintsPrefix):], " ")
		if isBracedStr(value, false) {
			if list, err := parseCommas(value[1 : len(value)-1]); err != nil {
				return nil, false, err
			} else {
				return list, true, nil
			}
		} else {
			return []string{value}, true, nil
		}
	}
	return nil, false, nil
}

func (pv *PropertyValidator) addTagItem(fieldName string, propertyName string, tagItem string) (result error) {
	result = nil
	tagToken := tagItem
	tagValue := ""
	hasColon := false
	if cAt := firstValidColonAt(tagItem); cAt != -1 {
		hasColon = true
		tagToken = strings.Trim(tagItem[0:cAt], " ")
		tagValue = strings.Trim(tagItem[cAt+1:], " ")
	}
	colonErr := false
	noColonErr := false
	switch tagToken {
	case tagTokenNotNull:
		colonErr = hasColon
		pv.NotNull = true
		break
	case tagTokenNullable:
		colonErr = hasColon
		pv.NotNull = false
		break
	case tagTokenMandatory, tagTokenRequired:
		if hasColon {
			result = pv.setTagMandatoryWhen(tagValue)
		}
		pv.Mandatory = true
		break
	case tagTokenOptional:
		colonErr = hasColon
		pv.Mandatory = false
		break
	case tagTokenType:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagType(tagValue)
		}
		break
	case tagTokenOrder:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagOrder(tagValue)
		}
		break
	case tagTokenConstraint:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.addConstraint(tagValue)
		}
		break
	case tagTokenObjConstraint:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagObjConstraint(tagValue)
		}
		break
	case tagTokenWhen:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagWhen(tagValue)
		}
		break
	case tagTokenUnwanted:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagUnwanted(tagValue)
		}
		break
	case tagTokenRequiredWith, tagTokenRequiredWithAlt:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.addRequiredWith(tagValue)
		}
	case tagTokenUnwantedWith, tagTokenUnwantedWithAlt:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.addUnwantedWith(tagValue)
		}
	case tagTokenRequiredWithMsg, tagTokenRequiredWithAltMsg:
		noColonErr = !hasColon
		if !noColonErr {
			if isQuotedStr(tagValue, true) {
				pv.RequiredWithMessage = tagValue[1 : len(tagValue)-1]
			} else {
				pv.RequiredWithMessage = tagValue
			}
		}
	case tagTokenUnwantedWithMsg, tagTokenUnwantedWithAltMsg:
		noColonErr = !hasColon
		if !noColonErr {
			if isQuotedStr(tagValue, true) {
				pv.UnwantedWithMessage = tagValue[1 : len(tagValue)-1]
			} else {
				pv.UnwantedWithMessage = tagValue
			}
		}
	case tagTokenObjIgnoreUnknownProperties:
		colonErr = hasColon
		if !colonErr {
			if pv.ObjectValidator == nil {
				result = fmt.Errorf(msgPropertyNotObject, tagToken)
			} else {
				pv.ObjectValidator.IgnoreUnknownProperties = true
			}
		}
		break
	case tagTokenObjUnknownProperties:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagObjUnknownProperties(tagValue)
		}
		break
	case tagTokenObjOrdered:
		colonErr = hasColon
		if !colonErr {
			if pv.ObjectValidator == nil {
				result = fmt.Errorf(msgPropertyNotObject, tagToken)
			} else {
				pv.ObjectValidator.OrderedPropertyChecks = true
			}
		}
		break
	case tagTokenObjWhen:
		noColonErr = !hasColon
		if !noColonErr {
			result = pv.setTagObjWhen(tagValue)
		}
		break
	default:
		if strings.HasPrefix(tagItem, "&") {
			// if tagName starts with '&' we can assume it's a constraint
			if err := pv.addConstraint(tagItem); err != nil {
				result = err
			}
		} else if ok, err := customTagTokenRegistry.handle(tagToken, hasColon, tagValue, pv, propertyName, fieldName); ok {
			result = err
		} else {
			result = fmt.Errorf(msgUnknownTokenInTag, tagToken)
		}
	}
	if result == nil {
		if colonErr {
			result = fmt.Errorf(msgUnexpectedColon, tagToken)
		} else if noColonErr {
			result = fmt.Errorf(msgExpectedColon, tagToken)
		}
	}
	return
}

func (pv *PropertyValidator) setTagType(tagValue string) error {
	ty, ok := JsonTypeFromString(tagValue)
	if !ok {
		return fmt.Errorf(msgUnknownPropertyType, tagValue)
	}
	pv.Type = ty
	return nil
}

func (pv *PropertyValidator) setTagOrder(tagValue string) error {
	v, err := strconv.ParseInt(tagValue, 10, 32)
	if err != nil {
		return fmt.Errorf(msgUnknownTagValue, tagTokenOrder, "int", tagValue)
	}
	pv.Order = int(v)
	return nil
}

func addConditions(conditions *Conditions, tagValue string, allowCurly bool) error {
	if isBracedStr(tagValue, allowCurly) {
		if tokens, err := parseCommas(tagValue[1 : len(tagValue)-1]); err == nil {
			for _, token := range tokens {
				if isQuotedStr(token, true) {
					*conditions = append(*conditions, token[1:len(token)-1])
				} else {
					*conditions = append(*conditions, token)
				}
			}
		} else {
			return err
		}
	} else if isQuotedStr(tagValue, true) {
		*conditions = append(*conditions, tagValue[1:len(tagValue)-1])
	} else {
		*conditions = append(*conditions, tagValue)
	}
	return nil
}

func (pv *PropertyValidator) setTagMandatoryWhen(tagValue string) error {
	return addConditions(&pv.MandatoryWhen, tagValue, true)
}

func (pv *PropertyValidator) setTagObjConstraint(tagValue string) error {
	if pv.ObjectValidator == nil {
		return fmt.Errorf(msgPropertyNotObject, tagTokenObjConstraint)
	}
	return pv.ObjectValidator.addConstraint(tagValue)
}

func (pv *PropertyValidator) setTagWhen(tagValue string) error {
	return addConditions(&pv.WhenConditions, tagValue, true)
}

func (pv *PropertyValidator) setTagUnwanted(tagValue string) error {
	return addConditions(&pv.UnwantedConditions, tagValue, true)
}

func (pv *PropertyValidator) setTagObjUnknownProperties(tagValue string) error {
	if pv.ObjectValidator == nil {
		return fmt.Errorf(msgPropertyNotObject, tagTokenObjUnknownProperties)
	}
	b, err := strconv.ParseBool(tagValue)
	if err != nil {
		return fmt.Errorf(msgUnknownTagValue, tagTokenObjUnknownProperties, "boolean", tagValue)
	}
	pv.ObjectValidator.IgnoreUnknownProperties = b
	return nil
}

func (pv *PropertyValidator) setTagObjWhen(tagValue string) error {
	if pv.ObjectValidator == nil {
		return fmt.Errorf(msgPropertyNotObject, tagTokenObjWhen)
	}
	if isBracedStr(tagValue, true) {
		if tokens, err := parseCommas(tagValue[1 : len(tagValue)-1]); err == nil {
			for _, token := range tokens {
				if isQuotedStr(token, true) {
					pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, token[1:len(token)-1])
				} else {
					pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, token)
				}
			}
		} else {
			return err
		}
	} else if isQuotedStr(tagValue, true) {
		pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, tagValue[1:len(tagValue)-1])
	} else {
		pv.ObjectValidator.WhenConditions = append(pv.ObjectValidator.WhenConditions, tagValue)
	}
	return nil
}

func (pv *PropertyValidator) addConstraint(tagValue string) error {
	c, err := buildConstraintFromTagValue(tagValue)
	if err != nil {
		return err
	}
	pv.Constraints = append(pv.Constraints, c)
	return nil
}

func (pv *PropertyValidator) addRequiredWith(tagValue string) error {
	expr, err := ParseExpression(tagValue)
	if err != nil {
		return err
	}
	if pv.RequiredWith == nil {
		pv.RequiredWith = expr
	} else {
		for _, x := range expr {
			pv.RequiredWith = append(pv.RequiredWith, x)
		}
	}
	return nil
}

func (pv *PropertyValidator) addUnwantedWith(tagValue string) error {
	expr, err := ParseExpression(tagValue)
	if err != nil {
		return err
	}
	if pv.UnwantedWith == nil {
		pv.UnwantedWith = expr
	} else {
		for _, x := range expr {
			pv.UnwantedWith = append(pv.UnwantedWith, x)
		}
	}
	return nil
}

func (v *Validator) addConstraint(tagValue string) error {
	c, err := buildConstraintFromTagValue(tagValue)
	if err != nil {
		return err
	}
	v.Constraints = append(v.Constraints, c)
	return nil
}

func buildConstraintFromTagValue(tagValue string) (Constraint, error) {
	useValue := strings.Trim(tagValue, " ")
	if strings.HasPrefix(useValue, "&") {
		useValue = useValue[1:]
	}
	isConditional := false
	var conditions Conditions
	if strings.HasPrefix(useValue, "[") {
		isConditional = true
		closeAt := strings.Index(useValue, "]")
		if closeAt == -1 {
			return nil, fmt.Errorf(msgConditionalConstraintsFormat, tagValue)
		}
		list := useValue[:closeAt+1]
		if err := addConditions(&conditions, list, false); err != nil {
			return nil, err
		}
		useValue = useValue[closeAt+1:]
	}
	constraintName := useValue
	argsStr := ""
	if curlyOpenAt := strings.Index(useValue, "{"); curlyOpenAt != -1 {
		if !strings.HasSuffix(useValue, "}") {
			return nil, fmt.Errorf(msgConstraintsFormat, tagValue)
		}
		argsStr = strings.Trim(useValue[curlyOpenAt+1:len(useValue)-1], " ")
		constraintName = useValue[0:curlyOpenAt]
	}
	//	if strings.HasPrefix(constraintName, "&") {
	//		constraintName = constraintName[1:]
	//	}
	if constraintName == constraintSetName {
		return buildConstraintSetWithArgs(argsStr)
	}
	c, ok := constraintsRegistry.get(constraintName)
	if !ok {
		return nil, fmt.Errorf(msgUnknownConstraint, constraintName)
	}
	// check if the tag value has any args specified...
	if argsStr == "" {
		// no args within curly braces, so it's safe to re-use the registered constraint...
		if isConditional {
			return &ConditionalConstraint{
				Constraint: c,
				When:       conditions,
			}, nil
		}
		return c, nil
	}
	newC, cErr := rebuildConstraintWithArgs(constraintName, c, argsStr)
	if cErr != nil {
		return nil, cErr
	}
	if isConditional {
		return &ConditionalConstraint{
			Constraint: newC,
			When:       conditions,
		}, nil
	}
	return newC, nil
}

func buildConstraintSetWithArgs(argsStr string) (Constraint, error) {
	args, err := parseCommas(argsStr)
	if err != nil {
		return nil, err
	}
	result := &ConstraintSet{
		Message:     "",
		Constraints: Constraints{},
	}
	for _, arg := range args {
		colonAt := firstValidColonAt(arg)
		if colonAt == -1 {
			return nil, fmt.Errorf(msgTagFldMissingColon, constraintSetName)
		}
		fieldName := strings.Trim(arg[0:colonAt], " ")
		value := strings.Trim(arg[colonAt+1:], " ")
		if fieldName == constraintSetFieldMessage {
			if isQuotedStr(value, true) {
				result.Message = value[1 : len(value)-1]
			} else {
				return nil, fmt.Errorf(msgConstraintFieldInvalidValue, constraintSetName, constraintSetFieldMessage)
			}
		} else if fieldName == constraintSetFieldConstraints {
			if isBracedStr(value, true) {
				tagValues, _ := parseCommas(value[1 : len(value)-1])
				for _, tagValue := range tagValues {
					if constraint, cErr := buildConstraintFromTagValue(tagValue); cErr != nil {
						return nil, cErr
					} else {
						result.Constraints = append(result.Constraints, constraint)
					}
				}
			} else {
				return nil, fmt.Errorf(msgConstraintFieldInvalidValue, constraintSetName, constraintSetFieldConstraints)
			}
		} else if fieldName == constraintSetFieldStop {
			if b, pErr := strconv.ParseBool(value); pErr == nil {
				result.Stop = b
			} else {
				return nil, fmt.Errorf(msgConstraintFieldInvalidValue, constraintSetName, constraintSetFieldStop)
			}
		} else if fieldName == constraintSetFieldOneOf {
			if b, pErr := strconv.ParseBool(value); pErr == nil {
				result.OneOf = b
			} else {
				return nil, fmt.Errorf(msgConstraintFieldInvalidValue, constraintSetName, constraintSetFieldOneOf)
			}
		} else {
			return nil, fmt.Errorf(msgConstraintFieldUnknown, constraintSetName, fieldName)
		}
	}
	return result, nil
}

func rebuildConstraintWithArgs(cName string, c Constraint, argsStr string) (Constraint, error) {
	args, err := argsStringToArgs(cName, argsStr)
	if err != nil {
		return nil, err
	}
	ty := reflect.TypeOf(c)
	// even though the original constraint implements Constraint - it still needs to be a struct
	// so that we can reconstruct the struct with args...
	if ty.Kind() != reflect.Ptr || ty.Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf(msgCannotCreateConstraint, cName)
	}
	ty = ty.Elem()
	newC := reflect.New(ty)
	var result = newC.Interface().(Constraint)
	// clone the original constraint fields into new constraint...
	orgV := reflect.ValueOf(c)
	count := ty.NumField()
	for f := 0; f < count; f++ {
		fn := ty.Field(f).Name
		fv := orgV.Elem().FieldByName(fn)
		fld := newC.Elem().FieldByName(fn)
		if fld.Kind() != reflect.Invalid && fld.CanSet() {
			fld.Set(fv)
		} else {
			return nil, fmt.Errorf(msgConstraintFieldNotExported, cName, fn)
		}
	}
	// now overwrite any specified args into the constraint fields...
	for argName, argVal := range args {
		if argName == "" && len(args) == 1 {
			if defName, err := getConstraintDefaultValueField(cName, ty); err != nil {
				return nil, err
			} else {
				argName = defName
			}
		}
		fld := newC.Elem().FieldByName(argName)
		if fld.Kind() != reflect.Invalid && fld.CanSet() {
			if !safeSet(fld, argVal) {
				return nil, fmt.Errorf(msgConstraintFieldInvalidValue, cName, argName)
			}
		} else {
			return nil, fmt.Errorf(msgConstraintFieldUnknown, cName, argName)
		}
	}
	return result, nil
}

func getConstraintDefaultValueField(cName string, ty reflect.Type) (string, error) {
	// scan fields to see if any are tagged as `v8n:"default"`...
	count := ty.NumField()
	if count == 1 {
		// only has one field, so that must be the default field...
		return ty.Field(0).Name, nil
	}
	for f := 0; f < count; f++ {
		sf := ty.Field(f)
		if sf.Tag.Get(tagNameV8n) == "default" {
			return sf.Name, nil
		}
	}
	return "", fmt.Errorf(msgConstraintNoDefaultValue, cName)
}

func safeSet(fv reflect.Value, valueStr string) (result bool) {
	result = false
	switch fv.Kind() {
	case reflect.String:
		if isQuotedStr(valueStr, true) {
			fv.SetString(valueStr[1 : len(valueStr)-1])
			result = true
		}
		break
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if i, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
			fv.SetInt(i)
			result = true
		}
		break
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if i, err := strconv.ParseUint(valueStr, 10, 64); err == nil {
			fv.SetUint(i)
			result = true
		}
		break
	case reflect.Float64, reflect.Float32:
		if f, err := strconv.ParseFloat(valueStr, 64); err == nil {
			fv.SetFloat(f)
			result = true
		}
		break
	case reflect.Bool:
		if b, err := strconv.ParseBool(valueStr); err == nil {
			fv.SetBool(b)
			result = true
		}
		break
	case reflect.Slice:
		if isBracedStr(valueStr, true) {
			if items, ok := itemsToSlice(fv.Type(), valueStr); ok {
				fv.Set(items)
				result = true
			}
		}
		break
	case regexpKind:
		if isQuotedStr(valueStr, true) {
			rxv := reflect.ValueOf(regexp.MustCompile(valueStr[1 : len(valueStr)-1])).Elem()
			fv.Set(rxv)
			result = true
		}
	}
	return
}

func itemsToSlice(itemType reflect.Type, arrayStr string) (result reflect.Value, ok bool) {
	ok = false
	if strItems, err := parseCommas(arrayStr[1 : len(arrayStr)-1]); err == nil {
		result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
		ik := itemType.Elem().Kind()
		switch ik {
		case reflect.String:
			ok = true
			for i, vu := range strItems {
				v := strings.Trim(vu, " ")
				if isQuotedStr(v, true) {
					result.Index(i).SetString(v[1 : len(v)-1])
				} else {
					ok = false
					break
				}
			}
			break
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			ok = true
			for i, vu := range strItems {
				v := strings.Trim(vu, " ")
				if iv, e := strconv.ParseInt(v, 10, 64); e == nil {
					result.Index(i).SetInt(iv)
				} else {
					ok = false
					break
				}
			}
			break
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			ok = true
			for i, vu := range strItems {
				v := strings.Trim(vu, " ")
				if iv, e := strconv.ParseUint(v, 10, 64); e == nil {
					result.Index(i).SetUint(iv)
				} else {
					ok = false
					break
				}
			}
			break
		case reflect.Float64, reflect.Float32:
			ok = true
			for i, vu := range strItems {
				v := strings.Trim(vu, " ")
				if fv, e := strconv.ParseFloat(v, 64); e == nil {
					result.Index(i).SetFloat(fv)
				} else {
					ok = false
					break
				}
			}
			break
		case reflect.Bool:
			ok = true
			for i, vu := range strItems {
				v := strings.Trim(vu, " ")
				if bv, e := strconv.ParseBool(v); e == nil {
					result.Index(i).SetBool(bv)
				} else {
					ok = false
					break
				}
			}
			break
		}
	}
	return
}

func argsStringToArgs(cName string, str string) (map[string]string, error) {
	rawArgs, err := parseCommas(str)
	if err != nil {
		return nil, fmt.Errorf(msgConstraintArgsParseError, cName, err)
	}
	result := map[string]string{}
	for _, arg := range rawArgs {
		argName := strings.Trim(arg, " ")
		cAt := firstValidColonAt(argName)
		if cAt == -1 {
			if len(rawArgs) == 1 {
				result[""] = argName
				break
			}
			return nil, fmt.Errorf(msgTagFldMissingColon, cName)
		}
		argValue := strings.Trim(argName[cAt+1:], " ")
		argName = strings.Trim(argName[:cAt], " ")
		result[argName] = argValue
	}
	return result, nil
}

func firstValidColonAt(str string) int {
	result := -1
	for i, ch := range str {
		if ch == ':' {
			result = i
			break
		} else if !(ch == ' ' || ch == '_' || ch == '.' || ch == '+' || ch == '-' ||
			(ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')) {
			break
		}
	}
	return result
}

func parseCommas(str string) ([]string, error) {
	result := make([]string, 0, len(str)/2)
	runes := []rune(str)
	stk := &delimiterStack{current: nil, stack: []*delimiter{}}
	lastTokenAt := 0
	for i, r := range runes {
		switch r {
		case ',':
			if !stk.inAny() {
				result = append(result, strings.Trim(string(runes[lastTokenAt:i]), " "))
				lastTokenAt = i + 1
			}
			break
		case '"', '\'', '[', ']', '(', ')', '{', '}':
			if err := stk.delimiter(r, i); err != nil {
				return nil, err
			}
			break
		}
	}
	if stk.inAny() {
		return nil, fmt.Errorf(msgUnclosed, stk.current.pos)
	}
	if lastTokenAt < len(runes) {
		result = append(result, strings.Trim(string(runes[lastTokenAt:]), " "))
	}
	return result, nil
}

type delimiterStack struct {
	current *delimiter
	stack   []*delimiter
}
type delimiter struct {
	open    rune
	pos     int
	isQuote bool
}

func (ds *delimiterStack) delimiter(ch rune, pos int) error {
	switch ch {
	case '"', '\'':
		if ds.current != nil && ds.current.open == ch {
			ds.pop()
		} else if !ds.inQuote() {
			ds.push(ch, pos)
		}
		break
	case '(', '[', '{':
		if !ds.inQuote() {
			ds.push(ch, pos)
		}
		break
	case ')':
		if !ds.inQuote() {
			if ds.current == nil || ds.current.open != '(' {
				return fmt.Errorf(msgUnopened, pos)
			}
			ds.pop()
		}
		break
	case ']':
		if !ds.inQuote() {
			if ds.current == nil || ds.current.open != '[' {
				return fmt.Errorf(msgUnopened, pos)
			}
			ds.pop()
		}
		break
	case '}':
		if !ds.inQuote() {
			if ds.current == nil || ds.current.open != '{' {
				return fmt.Errorf(msgUnopened, pos)
			}
			ds.pop()
		}
		break
	}
	return nil
}

func (ds *delimiterStack) push(ch rune, pos int) {
	if ds.current != nil {
		ds.stack = append(ds.stack, ds.current)
	}
	ds.current = &delimiter{open: ch, pos: pos, isQuote: ch == '"' || ch == '\''}
}

func (ds *delimiterStack) pop() {
	if len(ds.stack) > 0 {
		ds.current = ds.stack[len(ds.stack)-1]
		ds.stack = ds.stack[0 : len(ds.stack)-1]
	} else {
		ds.current = nil
	}
}
func (ds *delimiterStack) inAny() bool {
	return ds.current != nil
}
func (ds *delimiterStack) inQuote() bool {
	return ds.current != nil && ds.current.isQuote
}

func isQuotedStr(str string, allowSingles bool) bool {
	return (strings.HasPrefix(str, "\"") && strings.HasSuffix(str, "\"")) ||
		(allowSingles && strings.HasPrefix(str, "'") && strings.HasSuffix(str, "'"))
}

func isBracedStr(str string, allowCurly bool) bool {
	return (strings.HasPrefix(str, "[") && strings.HasSuffix(str, "]")) ||
		(allowCurly && strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}"))
}

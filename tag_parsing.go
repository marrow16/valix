package valix

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	msgV8nPrefix                    = "tag " + tagNameV8n + " - "
	msgUnknownPropertyType          = msgV8nPrefix + "unknown property type '%s'"
	msgUnknownTokenInTag            = msgV8nPrefix + "unknown token '%s'"
	msgUnexpectedColon              = msgV8nPrefix + "unexpected ':' colon after token '%s'"
	msgExpectedColon                = msgV8nPrefix + "expected ':' colon after token '%s'"
	msgConstraintsFormat            = msgV8nPrefix + "must specify constraints in the format '&name{}' (found \"%s\")"
	msgConditionalConstraintsFormat = msgV8nPrefix + "must specify conditional constraints in the format '&[token,...]name{}' or '&<expr>name{}' (found \"%s\")"
	msgConditionalExpr              = msgV8nPrefix + "invalid other properties expression \"%s\" - %s"
	msgUnknownConstraint            = msgV8nPrefix + "contains unknown constraint '%s'"
	msgCannotCreateConstraint       = msgV8nPrefix + "cannot create constraint '%s{}' (on non-struct constraint)"
	msgConstraintFieldUnknown       = msgV8nPrefix + "constraint '%s{}' field '%s' is unknown or not assignable"
	msgConstraintFieldInvalidValue  = msgV8nPrefix + "constraint '%s{}' field '%s' cannot be assigned with value specified"
	msgConstraintFieldNotExported   = msgV8nPrefix + "constraint '%s{}' has unexported field '%s' - so no fields can be specified as args"
	msgConstraintArgsParseError     = msgV8nPrefix + "constraint '%s{}' - args parsing error (%s)"
	msgUnknownTagValue              = msgV8nPrefix + "token '%s' expected %s value (found \"%s\")"
	msgPropertyNotObject            = msgV8nPrefix + "token '%s' cannot be used on non object/array field"
	msgUnclosed                     = "unclosed parenthesis or quote started at position %d"
	msgUnopened                     = "unopened parenthesis at position %d"
	msgWrapped                      = "field '%s' (property '%s') - %s"
)

var (
	// grab some type information (for types used in common constraints)
	rx             = regexp.MustCompile("([a])")
	regexpKind     = reflect.TypeOf(rx).Elem().Kind()
	constraintKind reflect.Kind
	otherKind      reflect.Kind
)

func init() {
	type dummyConstraint struct {
		Constraint  Constraint
		Constraints Constraints
		Other       Other
	}
	dmy := dummyConstraint{}
	v := reflect.ValueOf(dmy)
	f := v.FieldByName("Constraint")
	constraintKind = f.Kind()
	f = v.FieldByName("Other")
	otherKind = f.Kind()
}

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
	tagToken, tagValue, hasColon, err := splitTagItem(tagItem)
	if err != nil {
		return err
	}
	if op, ok := tagTokenOperations[tagToken]; ok {
		result = op(pv, hasColon, tagValue)
	} else {
		result = pv.defaultTagItemProcess(tagItem, tagToken, tagValue, hasColon, propertyName, fieldName)
	}
	return
}

func splitTagItem(tagItem string) (tagToken string, tagValue string, hasColon bool, err error) {
	tagToken = tagItem
	tagValue = ""
	hasColon = false
	err = nil
	if cAt := firstValidColonAt(tagItem); cAt != -1 {
		hasColon = true
		tagToken = strings.Trim(tagItem[0:cAt], " ")
		tagValue = strings.Trim(tagItem[cAt+1:], " ")
	}
	if expectColon, known := tagExpectsColon[tagToken]; known {
		if expectColon && !hasColon {
			err = fmt.Errorf(msgExpectedColon, tagToken)
		} else if !expectColon && hasColon {
			err = fmt.Errorf(msgUnexpectedColon, tagToken)
		}
	}
	return
}

func (pv *PropertyValidator) defaultTagItemProcess(tagItem, tagToken, tagValue string, hasColon bool, propertyName, fieldName string) error {
	if strings.HasPrefix(tagItem, "&") {
		// if tagName starts with '&' we can assume it's a constraint
		return pv.addConstraint(tagItem)
	} else if ok, err := customTagTokenRegistry.handle(tagToken, hasColon, tagValue, pv, propertyName, fieldName); ok {
		return err
	}
	return fmt.Errorf(msgUnknownTokenInTag, tagToken)
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

func (pv *PropertyValidator) setTagMandatoryWhen(hasValue bool, tagValue string) error {
	pv.Mandatory = true
	if hasValue {
		return addConditions(&pv.MandatoryWhen, tagValue, true)
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
	var others OthersExpr
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
	} else if strings.HasPrefix(useValue, "<") {
		isConditional = true
		closeAt := strings.Index(useValue, ">")
		if closeAt == -1 {
			return nil, fmt.Errorf(msgConditionalConstraintsFormat, tagValue)
		}
		if expr, err := ParseExpression(useValue[1:closeAt]); err != nil {
			return nil, fmt.Errorf(msgConditionalExpr, useValue[1:closeAt], err.Error())
		} else {
			others = expr
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
				Others:     others,
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
			Others:     others,
		}, nil
	}
	return newC, nil
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
	fields := make(map[string]reflect.Value, count)
	var defaultField *reflect.Value = nil
	defaultFieldName := ""
	var firstField *reflect.Value = nil
	firstFieldName := ""
	for f := 0; f < count; f++ {
		afld := ty.Field(f)
		fn := afld.Name
		fv := orgV.Elem().FieldByName(fn)
		fld := newC.Elem().FieldByName(fn)
		if fld.Kind() != reflect.Invalid && fld.CanSet() {
			fields[fn] = fld
			if defaultField == nil && afld.Tag.Get(tagNameV8n) == "default" {
				defaultField = &fld
				defaultFieldName = fn
			}
			if firstField == nil {
				firstField = &fld
				firstFieldName = fn
			}
			fld.Set(fv)
		} else {
			return nil, fmt.Errorf(msgConstraintFieldNotExported, cName, fn)
		}
	}
	if defaultField == nil && firstField != nil && len(fields) == 1 {
		defaultField = firstField
		defaultFieldName = firstFieldName
	}
	// now overwrite any specified args into the constraint fields...
	err = overwriteFieldArgs(cName, args, fields, defaultField, defaultFieldName)
	return result, err
}

func overwriteFieldArgs(cName string, args []argHolder, fields map[string]reflect.Value, defaultField *reflect.Value, defaultFieldName string) error {
	if len(args) == 1 && !args[0].hasValue && defaultField != nil {
		// only one arg, it has no value and we have a default field...
		if fld, ok := secondGuessField(args[0].name, fields); ok {
			if !safeSet(fld, "", false) {
				return fmt.Errorf(msgConstraintFieldInvalidValue, cName, args[0].name)
			}
		} else {
			if !safeSet(*defaultField, args[0].name, false) {
				return fmt.Errorf(msgConstraintFieldInvalidValue, cName, defaultFieldName)
			}
		}
	} else {
		for _, arg := range args {
			if fld, ok := secondGuessField(arg.name, fields); ok {
				if !safeSet(fld, arg.value, arg.hasValue) {
					return fmt.Errorf(msgConstraintFieldInvalidValue, cName, arg.name)
				}
			} else {
				return fmt.Errorf(msgConstraintFieldUnknown, cName, arg.name)
			}
		}
	}
	return nil
}

func secondGuessField(name string, fields map[string]reflect.Value) (reflect.Value, bool) {
	if fld, ok := fields[name]; ok {
		return fld, true
	}
	wds := camelToWords(name)
	lcName := strings.ToLower(name)
	candidates, singleCandidates, fld, singleFld := countFieldNameCandidates(lcName, len(wds) == 1, fields)
	if singleCandidates == 1 {
		return singleFld, true
	} else if candidates == 0 {
		for n, f := range fields {
			if abbreviateName(strings.ToLower(n)) == lcName {
				fld = f
				candidates = 1
				break
			}
			if len(wds) > 1 {
				fwds := camelToWords(n)
				if len(wds) == len(fwds) {
					matches := 0
					for i, w := range fwds {
						if w == wds[i] || strings.HasPrefix(w, wds[i]) || abbreviateName(w) == wds[i] {
							matches++
						}
					}
					if matches == len(wds) {
						fld = f
						candidates = 1
						break
					}
				}
			}
		}
	}
	return fld, candidates == 1
}

func countFieldNameCandidates(name string, singleWord bool, fields map[string]reflect.Value) (candidates, singleCandidates int, fld, singleFld reflect.Value) {
	candidates = 0
	singleCandidates = 0
	for n, f := range fields {
		ln := strings.ToLower(n)
		if ln == name {
			fld = f
			candidates = 1
			break
		} else if singleWord && strings.HasPrefix(ln, name) {
			singleFld = f
			singleCandidates++
			candidates++
		} else if strings.Contains(ln, name) {
			fld = f
			candidates++
		}
	}
	return
}

var vowelReplacer = strings.NewReplacer("a", "", "e", "", "i", "", "o", "", "u", "")

func abbreviateName(name string) string {
	result := name[:1] + vowelReplacer.Replace(name[1:])
	shift := 0
	buf := []byte(result)
	end := len(buf)
	for i := 1; i < (end - 1); i++ {
		if buf[i-1] == buf[i] {
			buf[i-shift] = buf[i+1]
			shift++
		}
	}
	return string(buf[:end-shift])
}

func camelToWords(str string) []string {
	if str == strings.ToUpper(str) {
		return []string{strings.ToLower(str)}
	}
	buf := []byte(str)
	lastCap := 0
	onCap := buf[0] >= 'A' && buf[0] <= 'Z'
	result := make([]string, 0, len(str)/3)
	for i := 1; i < len(buf); i++ {
		ch := buf[i]
		if ch < 'A' || ch > 'Z' {
			// not cap
			onCap = false
		} else if !onCap {
			// new cap start
			result = append(result, strings.ToLower(string(buf[lastCap:i])))
			onCap = true
			lastCap = i
		}
	}
	if lastCap < len(buf) {
		result = append(result, strings.ToLower(string(buf[lastCap:])))
	}
	return result
}

func safeSet(fv reflect.Value, valueStr string, hasValue bool) (result bool) {
	result = false
	switch fv.Kind() {
	case reflect.String:
		result = safeSetString(fv, valueStr)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result = safeSetInt(fv, valueStr)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result = safeSetUint(fv, valueStr)
	case reflect.Float64, reflect.Float32:
		result = safeSetFloat(fv, valueStr)
	case reflect.Bool:
		result = safeSetBool(fv, valueStr, hasValue)
	case reflect.Slice:
		result = safeSetSlice(fv, valueStr)
	case regexpKind:
		result = safeSetRegexp(fv, valueStr)
	case constraintKind:
		result = safeSetConstraint(fv, valueStr)
	}
	return
}

func safeSetString(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if isQuotedStr(valueStr, true) {
		fv.SetString(valueStr[1 : len(valueStr)-1])
		result = true
	}
	return
}

func safeSetInt(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if i, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		fv.SetInt(i)
		result = true
	}
	return
}

func safeSetUint(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if i, err := strconv.ParseUint(valueStr, 10, 64); err == nil {
		fv.SetUint(i)
		result = true
	}
	return
}

func safeSetFloat(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if f, err := strconv.ParseFloat(valueStr, 64); err == nil {
		fv.SetFloat(f)
		result = true
	}
	return
}

func safeSetBool(fv reflect.Value, valueStr string, hasValue bool) (result bool) {
	result = false
	if !hasValue && valueStr != "false" {
		fv.SetBool(true)
		result = true
	} else if b, err := strconv.ParseBool(valueStr); err == nil {
		fv.SetBool(b)
		result = true
	}
	return
}

func safeSetSlice(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if isBracedStr(valueStr, true) {
		if items, ok := itemsToSlice(fv.Type(), valueStr); ok {
			fv.Set(items)
			result = true
		}
	} else if fv.Type().Elem().Kind() == otherKind {
		useValue := valueStr
		if isQuotedStr(valueStr, true) {
			useValue = valueStr[1 : len(valueStr)-1]
		}
		if expr, err := ParseExpression(useValue); err == nil {
			vx := reflect.ValueOf(expr)
			fv.Set(vx)
			result = true
		}
	}
	return
}

func safeSetRegexp(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if isQuotedStr(valueStr, true) {
		if rx, err := regexp.Compile(valueStr[1 : len(valueStr)-1]); err == nil {
			rxv := reflect.ValueOf(rx).Elem()
			fv.Set(rxv)
			result = true
		}
	}
	return
}

func safeSetConstraint(fv reflect.Value, valueStr string) (result bool) {
	result = false
	if c, err := buildConstraintFromTagValue(valueStr); err == nil {
		fv.Set(reflect.ValueOf(c))
		result = true
	}
	return
}

func itemsToSlice(itemType reflect.Type, arrayStr string) (result reflect.Value, ok bool) {
	if strItems, err := parseCommas(arrayStr[1 : len(arrayStr)-1]); err == nil {
		switch itemType.Elem().Kind() {
		case reflect.String:
			result, ok = itemsToSliceOfStrings(strItems, itemType)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result, ok = itemsToSliceOfInts(strItems, itemType)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			result, ok = itemsToSliceOfUints(strItems, itemType)
		case reflect.Float64, reflect.Float32:
			result, ok = itemsToSliceOfFloats(strItems, itemType)
		case reflect.Bool:
			result, ok = itemsToSliceOfBools(strItems, itemType)
		case constraintKind:
			result, ok = itemsToSliceOfConstraints(strItems, itemType)
		default:
			result, ok = itemsToSliceOther(strItems, itemType)
		}
	}
	return
}

func itemsToSliceOther(arrayItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	result = reflect.MakeSlice(itemType, len(arrayItems), len(arrayItems))
	ok = true
	for i, vu := range arrayItems {
		v := strings.Trim(vu, " ")
		itm := reflect.New(itemType.Elem())
		itmI := itm.Interface()
		if err := json.Unmarshal([]byte(v), &itmI); err != nil {
			ok = false
			break
		}
		result.Index(i).Set(reflect.ValueOf(itmI).Elem())
	}
	return
}

func itemsToSliceOfStrings(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		v := strings.Trim(vu, " ")
		if isQuotedStr(v, true) {
			result.Index(i).SetString(v[1 : len(v)-1])
		} else {
			result.Index(i).SetString(v)
		}
	}
	return
}

func itemsToSliceOfInts(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		if iv, e := strconv.ParseInt(strings.Trim(vu, " "), 10, 64); e == nil {
			result.Index(i).SetInt(iv)
		} else {
			ok = false
			break
		}
	}
	return
}

func itemsToSliceOfUints(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		if iv, e := strconv.ParseUint(strings.Trim(vu, " "), 10, 64); e == nil {
			result.Index(i).SetUint(iv)
		} else {
			ok = false
			break
		}
	}
	return
}

func itemsToSliceOfFloats(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		if fv, e := strconv.ParseFloat(strings.Trim(vu, " "), 64); e == nil {
			result.Index(i).SetFloat(fv)
		} else {
			ok = false
			break
		}
	}
	return
}

func itemsToSliceOfBools(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		if bv, e := strconv.ParseBool(strings.Trim(vu, " ")); e == nil {
			result.Index(i).SetBool(bv)
		} else {
			ok = false
			break
		}
	}
	return
}

func itemsToSliceOfConstraints(strItems []string, itemType reflect.Type) (result reflect.Value, ok bool) {
	ok = true
	result = reflect.MakeSlice(itemType, len(strItems), len(strItems))
	for i, vu := range strItems {
		if c, err := buildConstraintFromTagValue(strings.Trim(vu, " ")); err == nil {
			result.Index(i).Set(reflect.ValueOf(c))
		} else {
			ok = false
			break
		}
	}
	return
}

type argHolder struct {
	name     string
	value    string
	hasValue bool
}

func argsStringToArgs(cName string, str string) ([]argHolder, error) {
	rawArgs, err := parseCommas(str)
	if err != nil {
		return nil, fmt.Errorf(msgConstraintArgsParseError, cName, err)
	}
	result := make([]argHolder, 0, len(rawArgs))
	for _, arg := range rawArgs {
		argName := strings.Trim(arg, " ")
		if cAt := firstValidColonAt(argName); cAt == -1 {
			result = append(result, argHolder{
				name:     argName,
				hasValue: false,
			})
		} else {
			result = append(result, argHolder{
				name:     strings.Trim(argName[:cAt], " "),
				value:    strings.Trim(argName[cAt+1:], " "),
				hasValue: true,
			})
		}
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
	result := make([]string, 0, len(str)/4)
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
		case '"', '\'', '[', ']', '(', ')', '{', '}':
			if err := stk.delimiter(r, i); err != nil {
				return nil, err
			}
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
		ds.delimiterQuote(ch, pos)
	case '(', '[', '{':
		ds.delimiterOpenBrace(ch, pos)
	case ')':
		return ds.delimiterCloseBrace(pos, '(')
	case ']':
		return ds.delimiterCloseBrace(pos, '[')
	case '}':
		return ds.delimiterCloseBrace(pos, '{')
	}
	return nil
}

func (ds *delimiterStack) delimiterQuote(ch rune, pos int) {
	if ds.current != nil && ds.current.open == ch {
		ds.pop()
	} else if !ds.inQuote() {
		ds.push(ch, pos)
	}
}

func (ds *delimiterStack) delimiterOpenBrace(ch rune, pos int) {
	if !ds.inQuote() {
		ds.push(ch, pos)
	}
}

func (ds *delimiterStack) delimiterCloseBrace(pos int, opened rune) error {
	if !ds.inQuote() {
		if ds.current == nil || ds.current.open != opened {
			return fmt.Errorf(msgUnopened, pos)
		}
		ds.pop()
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

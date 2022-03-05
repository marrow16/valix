package valix

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	tagItemNotNull           = "notNull"
	tagItemMandatory         = "mandatory"
	tagItemOptional          = "optional"
	tagItemType              = "type"
	tagItemConstraint        = "constraint"
	tagItemConstraints       = "constraints"
	tagItemConstraintsPrefix = tagItemConstraints + ":"
	tagItemOrder             = "order"
	// object level tag items...
	tagItemObjPrefix                  = "obj."
	tagItemObjIgnoreUnknownProperties = tagItemObjPrefix + "ignoreUnknownProperties"
	tagItemObjUnknownProperties       = tagItemObjPrefix + "unknownProperties" // true/false
	tagItemObjConstraint              = tagItemObjPrefix + tagItemConstraint
	tagItemObjOrdered                 = tagItemObjPrefix + "ordered"
)

const (
	msgV8nPrefix                    = "tag " + tagNameV8n + " - "
	msgUnknownPropertyType          = msgV8nPrefix + "unknown property type '%s'"
	msgUnknownTokenInTag            = msgV8nPrefix + "unknown token '%s'"
	msgConstraintsFormat            = msgV8nPrefix + "must specify constraints in the format '&name{}' (found \"%s\")"
	msgUnknownConstraint            = msgV8nPrefix + "contains unknown constraint '%s'"
	msgCannotCreateConstraint       = msgV8nPrefix + "cannot create constraint '%s{}' (on non-struct constraint)"
	msgConstraintFieldNotAssignable = msgV8nPrefix + "constraint '%s{}' field '%s' is not assignable"
	msgConstraintFieldInvalidValue  = msgV8nPrefix + "constraint '%s{}' field '%s' cannot be assigned with value specified"
	msgConstraintArgsParseError     = msgV8nPrefix + "constraint '%s{}' - args parsing error (%s)"
	msgTagFldMissingColon           = msgV8nPrefix + "constraint '%s{}` field missing ':' separator"
	msgUnknownTagValue              = msgV8nPrefix + "token '%s' expected %s value (found \"%s\")"
	msgUnclosed                     = "unclosed parenthesis or quote started at position %d"
	msgUnopened                     = "unopened parenthesis at position %d"
)

var (
	// grab some type information (for types used in common constraints)
	rx         = regexp.MustCompile("([a])")
	regexpKind = reflect.TypeOf(rx).Elem().Kind()
)

func (pv *PropertyValidator) processV8nTag(fld reflect.StructField) error {
	if tag, ok := fld.Tag.Lookup(tagNameV8n); ok {
		tagItems, err := parseCommas(tag)
		if err != nil {
			return err
		}
		return pv.processTagItems(tagItems)
	}
	return nil
}

func (pv *PropertyValidator) processTagItems(tagItems []string) error {
	for _, ti := range tagItems {
		cs, is, err := isConstraintsList(ti)
		if err != nil {
			return err
		} else if is {
			for _, c := range cs {
				if e2 := pv.addConstraint(c); e2 != nil {
					return e2
				}
			}
		} else if e3 := pv.addTagItem(ti); e3 != nil {
			return e3
		}
	}
	return nil
}

func isConstraintsList(tagItem string) ([]string, bool, error) {
	if strings.HasPrefix(tagItem, tagItemConstraintsPrefix) {
		value := strings.Trim(tagItem[len(tagItemConstraintsPrefix):], " ")
		if strings.HasPrefix(value, "[") && strings.HasSuffix(value, "]") {
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

func (pv *PropertyValidator) addTagItem(tagItem string) (result error) {
	result = nil
	tagName := tagItem
	tagValue := ""
	if cAt := firstValidColonAt(tagItem); cAt != -1 {
		tagName = strings.Trim(tagItem[0:cAt], " ")
		tagValue = strings.Trim(tagItem[cAt+1:], " ")
	}
	switch tagName {
	case tagItemNotNull:
		pv.NotNull = true
		break
	case tagItemMandatory:
		pv.Mandatory = true
		break
	case tagItemOptional:
		pv.Mandatory = false
		break
	case tagItemType:
		if ty := JsonTypeFromString(tagValue); ty == JsonTypeUndefined {
			result = errors.New(fmt.Sprintf(msgUnknownPropertyType, tagValue))
		} else {
			pv.Type = ty
		}
		break
	case tagItemOrder:
		if v, err := strconv.ParseInt(tagValue, 10, 32); err != nil {
			result = errors.New(fmt.Sprintf(msgUnknownTagValue, tagItemOrder, "int", tagValue))
		} else {
			pv.Order = int(v)
		}
		break
	case tagItemConstraint:
		if err := pv.addConstraint(tagValue); err != nil {
			result = err
		}
		break
	case tagItemObjConstraint:
		if err := pv.ObjectValidator.addConstraint(tagValue); err != nil {
			result = err
		}
		break
	case tagItemObjIgnoreUnknownProperties:
		pv.ObjectValidator.IgnoreUnknownProperties = true
		break
	case tagItemObjUnknownProperties:
		if b, err := strconv.ParseBool(tagValue); err == nil {
			pv.ObjectValidator.IgnoreUnknownProperties = b
		} else {
			result = errors.New(fmt.Sprintf(msgUnknownTagValue, tagItemObjUnknownProperties, "boolean", tagValue))
		}
		break
	case tagItemObjOrdered:
		pv.ObjectValidator.OrderedPropertyChecks = true
		break
	default:
		if strings.HasPrefix(tagItem, "&") && strings.HasSuffix(tagItem, "}") {
			// if tagName starts with '&' we can assume it's a constraint
			if err := pv.addConstraint(tagItem); err != nil {
				result = err
			}
		} else {
			result = errors.New(fmt.Sprintf(msgUnknownTokenInTag, tagName))
		}
	}
	return
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
	curlyOpenAt := strings.Index(useValue, "{")
	if curlyOpenAt == -1 || !strings.HasSuffix(useValue, "}") {
		return nil, errors.New(fmt.Sprintf(msgConstraintsFormat, tagValue))
	}
	constraintName := useValue[0:curlyOpenAt]
	if strings.HasPrefix(constraintName, "&") {
		constraintName = constraintName[1:]
	}
	c, ok := registry.get(constraintName)
	if !ok {
		return nil, errors.New(fmt.Sprintf(msgUnknownConstraint, constraintName))
	}
	// check if the tag value has any args specified...
	inCurly := strings.Trim(useValue[curlyOpenAt+1:len(useValue)-1], " ")
	newC, cErr := rebuildConstraintWithArgs(constraintName, c, inCurly)
	if cErr != nil {
		return nil, cErr
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
		return nil, errors.New(fmt.Sprintf(msgCannotCreateConstraint, cName))
	}
	ty = ty.Elem()
	newC := reflect.New(ty)
	var result = newC.Interface().(Constraint)
	// clone the original constraint fields into new constraint...
	orgV := reflect.ValueOf(c)
	for f := 0; f < ty.NumField(); f++ {
		fn := ty.Field(f).Name
		fv := orgV.Elem().FieldByName(fn)
		fld := newC.Elem().FieldByName(fn)
		if fld.Kind() != reflect.Invalid && fld.CanSet() {
			fld.Set(fv)
		}
	}
	// now overwrite any specified args into the constraint fields...
	for argName, argVal := range args {
		fld := newC.Elem().FieldByName(argName)
		if fld.Kind() != reflect.Invalid && fld.CanSet() {
			if !safeSet(fld, argVal) {
				return nil, errors.New(fmt.Sprintf(msgConstraintFieldInvalidValue, cName, argName))
			}
		} else {
			return nil, errors.New(fmt.Sprintf(msgConstraintFieldNotAssignable, cName, argName))
		}
	}
	return result, nil
}

func safeSet(fv reflect.Value, valueStr string) (result bool) {
	result = false
	switch fv.Kind() {
	case reflect.String:
		if (strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"")) ||
			(strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'")) {
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
		if (strings.HasPrefix(valueStr, "[") && strings.HasSuffix(valueStr, "]")) ||
			(strings.HasPrefix(valueStr, "{") && strings.HasSuffix(valueStr, "}")) {
			if items, ok := itemsToSlice(fv.Type(), valueStr); ok {
				fv.Set(items)
				result = true
			}
		}
		break
	case regexpKind:
		if (strings.HasPrefix(valueStr, "\"") && strings.HasSuffix(valueStr, "\"")) ||
			(strings.HasPrefix(valueStr, "'") && strings.HasSuffix(valueStr, "'")) {
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
				if (strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"")) ||
					(strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) {
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
		return nil, errors.New(fmt.Sprintf(msgConstraintArgsParseError, cName, err))
	}
	result := map[string]string{}
	for _, arg := range rawArgs {
		argName := strings.Trim(arg, " ")
		cAt := firstValidColonAt(argName)
		if cAt == -1 {
			return nil, errors.New(fmt.Sprintf(msgTagFldMissingColon, cName))
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
		} else if !(ch == ' ' || ch == '_' || ch == '.' || (ch >= '0' && ch <= '9') || (ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z')) {
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
				part := string(runes[lastTokenAt:i])
				result = append(result, part)
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
		return nil, errors.New(fmt.Sprintf(msgUnclosed, stk.current.pos))
	}
	if lastTokenAt < len(runes) {
		str := string(runes[lastTokenAt:])
		result = append(result, str)
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
				return errors.New(fmt.Sprintf(msgUnopened, pos))
			}
			ds.pop()
		}
		break
	case ']':
		if !ds.inQuote() {
			if ds.current == nil || ds.current.open != '[' {
				return errors.New(fmt.Sprintf(msgUnopened, pos))
			}
			ds.pop()
		}
		break
	case '}':
		if !ds.inQuote() {
			if ds.current == nil || ds.current.open != '{' {
				return errors.New(fmt.Sprintf(msgUnopened, pos))
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

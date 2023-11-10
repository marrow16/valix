package valix

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"reflect"
	"regexp"
	"strings"
)

// V8nTagStringOptions is used by PropertyValidator.ToV8nTagString to control the output v8n tag format
type V8nTagStringOptions struct {
	// if set to true, constraint names are abbreviated in the resulting v8n tag string
	AbbreviateConstraintNames bool
	// if set to true, constraint field names are abbreviated in the resulting v8n tag string
	AbbreviateFieldNames bool
	// the minimum (desired) length for abbreviated constraint field names (the default 0 is the same as setting to 3)
	MinimumFieldNameLength uint
	// if set to true, constraints where no fields are set have the trailing {} dropped
	DiscardUnneededCurlies bool
	// if set to true, conditional constraints are not unwrapped (i.e. not reduced to their short form)
	NoUnwrapConditionalConstraints bool
	// if set to true, does not put spaces between each v8n tag token
	UnSpaced bool
}

// ToV8nTagString converts the property validator to its v8n tag representation string
//
// Note: This method can panic if unable to convert one of the constraints to a v8n tag string.
// Although this should not occur with any built-in constraints it may happen on custom implemented constraints
func (pv *PropertyValidator) ToV8nTagString(options *V8nTagStringOptions) string {
	useOptions := &V8nTagStringOptions{}
	if options != nil {
		useOptions = options
	}
	parts := pv.v8nBasics(*useOptions)
	return strings.Join(parts, ternary(useOptions.UnSpaced).string(",", ", "))
}

func (pv *PropertyValidator) v8nBasics(options V8nTagStringOptions) []string {
	result := []string{v8nTagAndValue(tagTokenType, pv.Type.String())}
	result = append(result, pv.v8nNotNull()...)
	result = append(result, pv.v8nOrder()...)
	result = append(result, pv.v8nMandatory()...)
	result = append(result, pv.v8nOnly()...)
	result = append(result, pv.v8nStopOnFirst()...)
	result = append(result, pv.v8nWhenConditions()...)
	result = append(result, pv.v8nUnwantedConditions()...)
	result = append(result, pv.v8nRequiredWith()...)
	result = append(result, pv.v8nUnwantedWith()...)
	result = append(result, pv.v8nConstraints(options)...)
	result = append(result, pv.v8nObjectValidator(options)...)

	return result
}

func (pv *PropertyValidator) v8nOrder() (result []string) {
	if pv.Order != 0 {
		result = []string{fmt.Sprintf("%s:%d", tagTokenOrder, pv.Order)}
	}
	return
}

func (pv *PropertyValidator) v8nMandatory() (result []string) {
	if len(pv.MandatoryWhen) > 0 {
		result = []string{v8nTagAndOneOrMany(tagTokenRequired, pv.MandatoryWhen)}
	} else if pv.Mandatory {
		result = []string{tagTokenRequired}
	}
	return
}

func (pv *PropertyValidator) v8nNotNull() (result []string) {
	if pv.NotNull {
		result = []string{tagTokenNotNull}
	}
	return
}

func (pv *PropertyValidator) v8nOnly() (result []string) {
	if len(pv.OnlyConditions) > 0 {
		result = append(result, v8nTagAndOneOrMany(tagTokenOnly, pv.OnlyConditions))
	} else if pv.Only {
		result = append(result, tagTokenOnly)
	}
	if pv.OnlyMessage != "" && (pv.Only || len(pv.OnlyConditions) > 0) {
		result = append(result, v8nTagAndValue(tagTokenOnlyMsg, safeQuotes(pv.OnlyMessage)))
	}
	return
}

func (pv *PropertyValidator) v8nStopOnFirst() (result []string) {
	if pv.StopOnFirst {
		result = []string{tagTokenStopOnFirstAlt}
	}
	return
}

func (pv *PropertyValidator) v8nWhenConditions() (result []string) {
	if len(pv.WhenConditions) > 0 {
		result = []string{v8nTagAndOneOrMany(tagTokenWhen, pv.WhenConditions)}
	}
	return
}

func (pv *PropertyValidator) v8nUnwantedConditions() (result []string) {
	if len(pv.UnwantedConditions) > 0 {
		result = []string{v8nTagAndOneOrMany(tagTokenUnwanted, pv.UnwantedConditions)}
	}
	return
}

func (pv *PropertyValidator) v8nRequiredWith() (result []string) {
	if pv.RequiredWith != nil && len(pv.RequiredWith) > 0 {
		result = []string{v8nTagAndValue(tagTokenRequiredWithAlt, pv.RequiredWith.String())}
		if pv.RequiredWithMessage != "" {
			result = append(result, fmt.Sprintf(v8nTagAndValue(tagTokenRequiredWithAltMsg, safeQuotes(pv.RequiredWithMessage))))
		}
	}
	return
}

func (pv *PropertyValidator) v8nUnwantedWith() (result []string) {
	if pv.UnwantedWith != nil && len(pv.UnwantedWith) > 0 {
		result = []string{v8nTagAndValue(tagTokenUnwantedWithAlt, pv.UnwantedWith.String())}
		if pv.UnwantedWithMessage != "" {
			result = append(result, v8nTagAndValue(tagTokenUnwantedWithAltMsg, safeQuotes(pv.UnwantedWithMessage)))
		}
	}
	return
}

func (pv *PropertyValidator) v8nConstraints(options V8nTagStringOptions) (result []string) {
	for _, constraint := range pv.Constraints {
		if constraint != nil {
			result = append(result, v8nConstraintToString(constraint, options, false))
		}
	}
	return
}

func v8nConstraintToString(constraint Constraint, options V8nTagStringOptions, noAmp bool) string {
	if cstr, is := v8nIsConditionalConstraint(constraint, options); is {
		return cstr
	}
	pfx := ternary(noAmp).string("", "&")

	found, raw, name, unchanged := v8nConstraintFromRegistry(constraint, options.AbbreviateConstraintNames)
	if !found {
		panic(fmt.Sprintf("cannot find constraint '%s' in registry", name))
	}
	if unchanged {
		return pfx + name + ternary(options.DiscardUnneededCurlies).string("", "{}")
	}
	rawV := reflect.ValueOf(raw).Elem()
	rawTy := reflect.TypeOf(raw).Elem()
	vflds := reflect.VisibleFields(rawTy)
	actV := reflect.ValueOf(constraint).Elem()
	diffs := make([]v8nConstraintField, 0, len(vflds))
	allNames := map[string]bool{}
	for _, fld := range vflds {
		if fld.IsExported() {
			allNames[fld.Name] = true
			af := actV.Field(fld.Index[0])
			rf := rawV.Field(fld.Index[0])
			if !cmp.Equal(af.Interface(), rf.Interface(), regexOption) {
				diff := v8nConstraintField{
					name:     fld.Name,
					value:    af.Interface(),
					rValue:   af,
					orgValue: rf.Interface(),
				}
				if df, ok := fld.Tag.Lookup(tagNameV8n); ok && df == "default" {
					diff.isDefault = true
				}
				diffs = append(diffs, diff)
			}
		}
	}
	if len(diffs) == 0 {
		return pfx + name + ternary(options.DiscardUnneededCurlies).string("", "{}")
	}
	return pfx + name + "{" + v8nFieldsToString(diffs, allNames, options) + "}"
}

func v8nConstraintFromRegistry(constraint Constraint, abbr bool) (found bool, result Constraint, name string, unchanged bool) {
	name = reflect.TypeOf(constraint).Elem().Name()
	if preset, ok := constraint.(*StringPresetPattern); ok {
		if preC, ok := constraintsRegistry.get(preset.Preset); ok {
			found = true
			result = preC
			name = preset.Preset
			unchanged = !preset.Stop && !preset.Strict && preset.Message == ""
			return
		}
	}
	result, found = constraintsRegistry.get(name)
	if found && abbr {
		if abbrName, ok := abbreviatedConstraintNames[name]; ok {
			if abbrC, ok := constraintsRegistry.get(abbrName); ok {
				name = abbrName
				result = abbrC
			}
		}
	}
	if !found {
		// still didn't find it - it might be registered with an alias in the registry...
		if alias, c, ok := constraintsRegistry.search(name); ok {
			name = alias
			result = c
			found = true
		}
	}
	return
}

func v8nFieldsToString(diffs []v8nConstraintField, allNames map[string]bool, options V8nTagStringOptions) string {
	if len(diffs) == 1 && diffs[0].isDefault {
		if v, colReqd := diffs[0].valueToString(options); colReqd {
			return v
		}
		return "true"
	}
	list := make([]string, len(diffs))
	for i, f := range diffs {
		list[i] = f.toString(allNames, options)
	}
	return strings.Join(list, ternary(options.UnSpaced).string(",", ", "))
}

type v8nConstraintField struct {
	name      string
	isDefault bool
	value     interface{}
	rValue    reflect.Value
	orgValue  interface{}
}

func (f v8nConstraintField) toString(allNames map[string]bool, options V8nTagStringOptions) string {
	aName := f.reduceFieldName(allNames, options)
	if v, reqColon := f.valueToString(options); reqColon {
		return v8nTagAndValue(aName, v)
	} else {
		return aName
	}
}

func (f v8nConstraintField) reduceFieldName(allNames map[string]bool, options V8nTagStringOptions) string {
	if !options.AbbreviateFieldNames {
		return f.name
	}
	return shortestUniqueFieldName(f.name, allNames, options.MinimumFieldNameLength)
}

func (f v8nConstraintField) valueToString(options V8nTagStringOptions) (string, bool) {
	switch tv := f.value.(type) {
	case string:
		return safeQuotes(tv), true
	case bool:
		org := false
		if ov, ok := f.orgValue.(bool); ok {
			org = ov
		}
		if !org && tv {
			return "", false
		}
		return fmt.Sprintf("%v", tv), true
	case float32, float64, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%v", tv), true
	case OthersExpr:
		return "'" + tv.String() + "'", true
	case Conditions:
		return "[" + strings.Join(tv, ",") + "]", true
	case Constraint:
		return v8nConstraintToString(tv, options, false), true
	case Constraints:
		return "[" + v8nConstraintsToString(tv, options) + "]", true
	case regexp.Regexp:
		return `'` + tv.String() + `'`, true
	}
	switch f.rValue.Kind() {
	case reflect.Slice:
		if data, err := json.Marshal(f.value); err == nil {
			return string(data[:]), true
		}
	}
	panic(fmt.Sprintf("unable to convert constraint field '%s' to v8n tag string", f.name))
}

func shortestUniqueFieldName(name string, allNames map[string]bool, minLen uint) (result string) {
	min := int(minLen)
	if min == 0 {
		min = 3
	}
	result = name
	if min >= len(name) {
		return
	}
	wds := camelToWords(name)
	if len(wds) == 1 {
		for i := min; i < len(name); i++ {
			sn := name[:i]
			if isFieldNameUnique(sn, allNames) {
				result = sn
				break
			}
		}
	} else {
		for i := min; i < len(name); i++ {
			if isFieldNameCamelUnique(wds, i, allNames) {
				res := make([]string, len(wds))
				for j, wd := range wds {
					res[j] = strings.ToUpper(wd[:1]) + wd[1:i]
				}
				result = strings.Join(res, "")
				break
			}
		}
	}
	return
}

func isFieldNameUnique(name string, allNames map[string]bool) bool {
	matches := 0
	ln := strings.ToLower(name)
	for kn := range allNames {
		if lkn := strings.ToLower(kn); lkn == ln || strings.HasPrefix(lkn, ln) {
			matches++
		}
	}
	return matches == 1
}

func isFieldNameCamelUnique(camels []string, ln int, allNames map[string]bool) bool {
	matches := 0
	for kn := range allNames {
		wds := camelToWords(kn)
		if len(wds) == len(camels) {
			matchWords := 0
			for i := 0; i < len(wds); i++ {
				if len(camels[i]) >= ln && strings.HasPrefix(wds[i], camels[i][:ln]) {
					matchWords++
				}
			}
			if matchWords == len(wds) {
				matches++
			}
		}
	}
	return matches == 1
}

func v8nConstraintsToString(constraints Constraints, options V8nTagStringOptions) string {
	result := make([]string, 0, len(constraints))
	for _, constraint := range constraints {
		if constraint != nil {
			result = append(result, v8nConstraintToString(constraint, options, false))
		}
	}
	return strings.Join(result, ternary(options.UnSpaced).string(",", ", "))
}

func v8nIsConditionalConstraint(constraint Constraint, options V8nTagStringOptions) (string, bool) {
	if conditional, ok := constraint.(*ConditionalConstraint); ok && conditional.Constraint != nil {
		if len(conditional.When) == 0 && len(conditional.Others) == 0 {
			// was pointless conditional...
			return v8nConstraintToString(conditional.Constraint, options, false), true
		} else if _, dblWrap := conditional.Constraint.(*ConditionalConstraint); dblWrap {
			// cannot unwrap double wrapped...
			return "", false
		} else if len(conditional.When) != 0 && len(conditional.Others) != 0 {
			// cannot unwrap it because both When and Others are specified...
			return "", false
		} else if !options.NoUnwrapConditionalConstraints {
			if len(conditional.When) != 0 {
				return "&[" + strings.Join(conditional.When, ",") + "]" + v8nConstraintToString(conditional.Constraint, options, true), true
			}
			return "&<" + conditional.Others.String() + ">" + v8nConstraintToString(conditional.Constraint, options, true), true
		}
	}
	return "", false
}

func (pv *PropertyValidator) v8nObjectValidator(options V8nTagStringOptions) (result []string) {
	if pv.ObjectValidator != nil {
		if pv.ObjectValidator.IgnoreUnknownProperties {
			result = append(result, tagTokenObjIgnoreUnknownProperties)
		}
		if pv.ObjectValidator.OrderedPropertyChecks {
			result = append(result, tagTokenObjOrdered)
		}
		if pv.ObjectValidator.AllowNullItems {
			result = append(result, tagTokenArrAllowNullItems)
		}
		if len(pv.ObjectValidator.WhenConditions) == 1 {
			result = append(result, tagTokenObjWhen+":"+pv.ObjectValidator.WhenConditions[0])

		} else if len(pv.ObjectValidator.WhenConditions) > 0 {
			result = append(result, tagTokenObjWhen+":["+strings.Join(pv.ObjectValidator.WhenConditions, ",")+"]")
		}
		for _, c := range pv.ObjectValidator.Constraints {
			if c != nil {
				result = append(result, tagTokenObjConstraint+":"+v8nConstraintToString(c, options, false))
			}
		}
	}
	return
}

func v8nTagAndValue(tag, val string) string {
	return fmt.Sprintf(tag + ":" + val)
}

func v8nTagAndOneOrMany(tag string, vals []string) string {
	if len(vals) == 1 {
		return v8nTagAndValue(tag, vals[0])
	}
	return v8nTagAndValue(tag, "["+strings.Join(vals, ",")+"]")
}

func safeQuotes(str string) string {
	if !strings.Contains(str, `'`) {
		return "'" + str + "'"
	}
	return `"` + strings.ReplaceAll(str, `"`, `\"`) + `"`
}

var abbreviatedConstraintNames = map[string]string{
	"ArrayConditionalConstraint":      "acond",
	"ArrayDistinctProperty":           "adistinctp",
	"DatetimeYearsOld":                "age",
	"ArrayOf":                         "aof",
	"ArrayUnique":                     "aunique",
	"ConditionalConstraint":           "cond",
	"SetConditionFrom":                "cfrom",
	"SetConditionIf":                  "cif",
	"StringContains":                  "contains",
	"SetConditionProperty":            "cpty",
	"SetConditionOnType":              "ctype",
	"DatetimeDayOfWeek":               "dtdow",
	"DatetimeFuture":                  "dtfuture",
	"DatetimeFutureOrPresent":         "dtfuturep",
	"DatetimeGreaterThan":             "dtgt",
	"DatetimeGreaterThanOrEqual":      "dtgte",
	"DatetimeGreaterThanOrEqualOther": "dtgteo",
	"DatetimeGreaterThanOther":        "dtgto",
	"DatetimeLessThan":                "dtlt",
	"DatetimeLessThanOrEqual":         "dtlte",
	"DatetimeLessThanOrEqualOther":    "dtlteo",
	"DatetimeLessThanOther":           "dtlto",
	"DatetimePast":                    "dtpast",
	"DatetimePastOrPresent":           "dtpastp",
	"DatetimeRange":                   "dtrange",
	"DatetimeTimeOfDayRange":          "dttodrange",
	"DatetimeTolerance":               "dttol",
	"DatetimeToleranceToNow":          "dttolnow",
	"DatetimeToleranceToOther":        "dttolother",
	"StringEndsWith":                  "ends",
	"EqualsOther":                     "eqo",
	"FailingConstraint":               "fail",
	"FailWhen":                        "failw",
	"FailWith":                        "failwith",
	"GreaterThan":                     "gt",
	"GreaterThanOrEqual":              "gte",
	"GreaterThanOrEqualOther":         "gteo",
	"GreaterThanOther":                "gto",
	"NetIsCIDR":                       "isCIDR",
	"NetIsHostname":                   "isHostname",
	"NetIsIP":                         "isIP",
	"NetIsMac":                        "isMac",
	"NetIsTCP":                        "isTCP",
	"NetIsTld":                        "isTld",
	"NetIsUDP":                        "isUDP",
	"NetIsURI":                        "isURI",
	"NetIsURL":                        "isURL",
	"Length":                          "len",
	"LengthExact":                     "lenx",
	"LessThan":                        "lt",
	"LessThanOrEqual":                 "lte",
	"LessThanOrEqualOther":            "lteo",
	"LessThanOther":                   "lto",
	"Maximum":                         "max",
	"MaximumInt":                      "maxi",
	"Minimum":                         "min",
	"MinimumInt":                      "mini",
	"Negative":                        "neg",
	"NegativeOrZero":                  "negz",
	"NotEqualsOther":                  "neqo",
	"NotEmpty":                        "notempty",
	"Positive":                        "pos",
	"PositiveOrZero":                  "posz",
	"Range":                           "range",
	"RangeInt":                        "rangei",
	"ConstraintSet":                   "set",
	"StringStartsWith":                "starts",
	"StringValidCurrencyCode":         "strccy",
	"StringCharacters":                "strchars",
	"StringValidCountryCode":          "strcountry",
	"StringValidEmail":                "stremail",
	"StringGreaterThan":               "strgt",
	"StringGreaterThanOrEqual":        "strgte",
	"StringGreaterThanOrEqualOther":   "strgteo",
	"StringGreaterThanOther":          "strgto",
	"StringValidISODate":              "strisod",
	"StringValidISODatetime":          "strisodt",
	"StringValidISODuration":          "strisodur",
	"StringValidJson":                 "strjson",
	"StringValidLanguageCode":         "strlang",
	"StringLength":                    "strlen",
	"StringLowercase":                 "strlower",
	"StringLessThan":                  "strlt",
	"StringLessThanOrEqual":           "strlte",
	"StringLessThanOrEqualOther":      "strlteo",
	"StringLessThanOther":             "strlto",
	"StringMaxLength":                 "strmax",
	"StringMinLength":                 "strmin",
	"StringNotBlank":                  "strnb",
	"StringNotEmpty":                  "strne",
	"StringNoControlCharacters":       "strnocc",
	"StringPattern":                   "strpatt",
	"StringPresetPattern":             "strpreset",
	"StringValidToken":                "strtoken",
	"StringValidTimezone":             "strtz",
	"StringValidUnicodeNormalization": "struninorm",
	"StringUppercase":                 "strupper",
	"StringValidUuid":                 "struuid",
	"StringValidCardNumber":           "strvcn",
	"StringExactLength":               "strxlen",
	"MultipleOf":                      "xof",
}

var regexOption = cmp.FilterValues(regexCompareFilter, cmp.Comparer(regexComparator))

func regexCompareFilter(v1, v2 interface{}) bool {
	switch v1.(type) {
	case regexp.Regexp:
		switch v2.(type) {
		case regexp.Regexp:
			return true
		}
	}
	return false
}

func regexComparator(av1, av2 interface{}) (result bool) {
	switch rx1 := av1.(type) {
	case regexp.Regexp:
		switch rx2 := av2.(type) {
		case regexp.Regexp:
			result = rx1.String() == rx2.String()
		}
	}
	return
}

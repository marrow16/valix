package valix

import (
	"encoding/json"
	"github.com/google/go-cmp/cmp"
	"golang.org/x/text/unicode/norm"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	// Note: In formats, explicit argument indices are used - this is to aid i18n translations
	msgNotEmpty                  = "Value must not be empty"
	msgNotEmptyString            = "String value must not be an empty string"
	msgNotBlankString            = "String value must not be a blank string"
	msgNoControlChars            = "String value must not contain control characters"
	msgValidPattern              = "String value must have valid pattern"
	fmtMsgUnknownPresetPattern   = "Unknown preset pattern '%[1]s'"
	fmtMsgValidToken             = "String value must be valid token - \"%[1]s\""
	msgInvalidCharacters         = "String value must not have invalid characters"
	fmtMsgStringMinLen           = "String value length must be at least %[1]d characters"
	fmtMsgStringMinLenExc        = "String value length must be greater than %[1]d characters"
	fmtMsgStringMaxLen           = "String value length must not exceed %[1]d characters"
	fmtMsgStringMaxLenExc        = "String value length must be less than %[1]d characters"
	fmtMsgStringExactLen         = "String value length must be %[1]d characters"
	fmtMsgStringMinMaxLen        = "String value length must be between %[1]d (%[2]s) and %[3]d (%[4]s)"
	msgStringLowercase           = "String value must contain only lowercase letters"
	msgStringUppercase           = "String value must contain only uppercase letters"
	msgStringValidJson           = "String value must be valid JSON"
	msgUnicodeNormalization      = "String value must be correct normalization form"
	msgUnicodeNormalizationNFC   = "String value must be correct normalization form NFC"
	msgUnicodeNormalizationNFKC  = "String value must be correct normalization form NFKC"
	msgUnicodeNormalizationNFD   = "String value must be correct normalization form NFD"
	msgUnicodeNormalizationNFKD  = "String value must be correct normalization form NFKD"
	fmtMsgStringContains         = "String must contain %[1]s"
	fmtMsgStringNotContains      = "String must not contain %[1]s"
	fmtMsgStringStartsWith       = "String value must start with %[1]s"
	fmtMsgStringNotStartsWith    = "String value must not start with %[1]s"
	fmtMsgStringEndsWith         = "String value must end with %[1]s"
	fmtMsgStringNotEndsWith      = "String value must not end with %[1]s"
	fmtMsgMinLen                 = "Value length must be at least %[1]d"
	fmtMsgMinLenExc              = "Value length must be greater than %[1]d"
	fmtMsgExactLen               = "Value length must be %[1]d"
	fmtMsgMinMax                 = "Value length must be between %[1]d (%[2]s) and %[3]d (%[4]s)"
	msgPositive                  = "Value must be positive"
	msgPositiveOrZero            = "Value must be positive or zero"
	msgNegative                  = "Value must be negative"
	msgNegativeOrZero            = "Value must be negative or zero"
	msgNull                      = "Value must be null"
	fmtMsgGt                     = "Value must be greater than %[1]v"
	fmtMsgGte                    = "Value must be greater than or equal to %[1]v"
	fmtMsgLt                     = "Value must be less than %[1]v"
	fmtMsgLte                    = "Value must be less than or equal to %[1]v"
	fmtMsgStrGt                  = "Value must be greater than '%[1]s'"
	fmtMsgStrGte                 = "Value must be greater than or equal to '%[1]s'"
	fmtMsgStrLt                  = "Value must be less than '%[1]s'"
	fmtMsgStrLte                 = "Value must be less than or equal to '%[1]s'"
	fmtMsgRange                  = "Value must be between %[1]v (%[2]s) and %[3]v (%[4]s)"
	fmtMsgMultipleOf             = "Value must be a multiple of %[1]d"
	fmtMsgArrayElementType       = "Array elements must be of type %[1]s"
	fmtMsgArrayElementTypeOrNull = "Array elements must be of type %[1]s or null"
	msgArrayUnique               = "Array elements must be unique"
	msgValidUuid                 = "Value must be a valid UUID"
	fmtMsgUuidMinVersion         = "Value must be a valid UUID (minimum version %[1]d)"
	fmtMsgUuidCorrectVer         = "Value must be a valid UUID (version %[1]d)"
	msgValidCardNumber           = "Value must be a valid card number"
	msgValidCountryCode          = "Value must be a valid ISO-3166 country code"
	msgValidCurrencyCode         = "Value must be a valid ISO-4217 currency code"
	msgValidEmail                = "Value must be an email address"
	msgValidLanguageCode         = "Value must be a valid language code"
	fmtMsgEqualsOther            = "Value must equal the value of property '%[1]s'"
	fmtMsgNotEqualsOther         = "Value must not equal the value of property '%[1]s'"
	fmtMsgGtOther                = "Value must be greater than value of property '%[1]s'"
	fmtMsgGteOther               = "Value must be greater than or equal to value of property '%[1]s'"
	fmtMsgLtOther                = "Value must be less than value of property '%[1]s'"
	fmtMsgLteOther               = "Value must be less than or equal to value of property '%[1]s'"
	fmtMsgDtGt                   = "Value must be after '%[1]s'"
	fmtMsgDtGte                  = "Value must be after or equal to '%[1]s'"
	fmtMsgDtLt                   = "Value must be before '%[1]s'"
	fmtMsgDtLte                  = "Value must be before or equal to '%[1]s'"
	msgFailure                   = "Validation failed"
)

func defaultMessage(tcx I18nContext, msg string, def string, defArgs ...interface{}) string {
	if msg == "" {
		if len(defArgs) > 0 {
			return obtainI18nContext(tcx).TranslateFormat(def, defArgs...)
		}
		return obtainI18nContext(tcx).TranslateMessage(def)
	} else if msg == NoMessage {
		return ""
	}
	return obtainI18nContext(tcx).TranslateMessage(msg)
}

const (
	constraintPtyNameMessage = "Message"
	constraintPtyNameStop    = "Stop"
)

const (
	tokenInclusive = "inclusive"
	tokenExclusive = "exclusive"
)

func incExc(tcx I18nContext, exc bool) string {
	return obtainI18nContext(tcx).TranslateToken(ternary(exc).string(tokenExclusive, tokenInclusive))
}

var (
	uuidRegexp = regexp.MustCompile(uuidRegexpPattern)
)

const (
	uuidRegexpPattern                 = "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"
	iso8601FullPattern                = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601NoOffsPattern              = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?)$"
	iso8601NoMillisPattern            = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601MinPattern                 = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2})$"
	iso8601DateOnlyPattern            = "^(\\d{4}-\\d{2}-\\d{2})$"
	iso8601FullLayout                 = "2006-01-02T15:04:05.999999999Z07:00"
	iso8601NoOffLayout                = "2006-01-02T15:04:05.999999999"
	iso8601NoMillisLayout             = "2006-01-02T15:04:05Z07:00"
	iso8601MinLayout                  = "2006-01-02T15:04:05"
	iso8601DateOnlyLayout             = "2006-01-02"
	msgValidISODate                   = "Value must be a valid date string (format: YYYY-MM-DD)"
	msgValidISODatetimeFormatFull     = "Value must be a valid date/time string (format: YYYY-MM-DDThh:mm:ss.sss[Z|+-hh:mm])"
	msgValidISODatetimeFormatNoOffs   = "Value must be a valid date/time string (format: YYYY-MM-DDThh:mm:ss.sss)"
	msgValidISODatetimeFormatNoMillis = "Value must be a valid date/time string (format: YYYY-MM-DDThh:mm:ss[Z|+-hh:mm])"
	msgValidISODatetimeFormatMin      = "Value must be a valid date/time string (format: YYYY-MM-DDThh:mm:ss)"
	msgDatetimeFuture                 = "Value must be a valid date/time in the future"
	msgDatetimeFutureOrPresent        = "Value must be a valid date/time in the future or present"
	msgDatetimePast                   = "Value must be a valid date/time in the past"
	msgDatetimePastOrPresent          = "Value must be a valid date/time in the past or present"
	msgValidTimezone                  = "Value must be a valid timezone"
	msgDatetimeDayOfWeek              = "Value must be a valid day of week"
	msgValidISODuration               = "Value must be a valid ISO 8601 duration"
)

var (
	iso8601FullRegex     = regexp.MustCompile(iso8601FullPattern)
	iso8601NoOffsRegex   = regexp.MustCompile(iso8601NoOffsPattern)
	iso8601NoMillisRegex = regexp.MustCompile(iso8601NoMillisPattern)
	iso8601MinRegex      = regexp.MustCompile(iso8601MinPattern)
	iso8601DateOnlyRegex = regexp.MustCompile(iso8601DateOnlyPattern)
)

func stringToDatetime(str string, truncTime bool) (*time.Time, bool) {
	if truncTime {
		if len(str) < 10 {
			return nil, false
		}
		dateStr := str[:10]
		if iso8601DateOnlyRegex.MatchString(dateStr) {
			result, err := time.Parse(iso8601DateOnlyLayout, dateStr)
			return &result, err == nil
		}
		return nil, false
	}
	parseLayout := ""
	if iso8601DateOnlyRegex.MatchString(str) {
		parseLayout = iso8601DateOnlyLayout
	} else if iso8601MinRegex.MatchString(str) {
		parseLayout = iso8601MinLayout
	} else if iso8601NoMillisRegex.MatchString(str) {
		parseLayout = iso8601NoMillisLayout
	} else if iso8601NoOffsRegex.MatchString(str) {
		parseLayout = iso8601NoOffLayout
	} else if iso8601FullRegex.MatchString(str) {
		parseLayout = iso8601FullLayout
	} else {
		return nil, false
	}
	result, err := time.Parse(parseLayout, str)
	return &result, err == nil
}

func truncateTime(t time.Time, truncTime bool) time.Time {
	if truncTime {
		return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
	}
	return t
}

var (
	// UnicodeBMP is a unicode.RangeTable that represents the Unicode BMP (Basic Multilingual Plane)
	//
	// For use with StringCharacters constraint
	UnicodeBMP = unicode.RangeTable{
		R16: []unicode.Range16{
			{0x0000, 0xffff, 1},
		},
	}
	// UnicodeSMP is a unicode.RangeTable that represents the Unicode SMP (Supplementary Multilingual Plane)
	//
	// For use with StringCharacters constraint
	UnicodeSMP = unicode.RangeTable{
		R32: []unicode.Range32{
			{0x10000, 0x1ffff, 1},
		},
	}
	// UnicodeSIP is a unicode.RangeTable that represents the Unicode SIP (Supplementary Ideographic Plane)
	//
	// For use with StringCharacters constraint
	UnicodeSIP = unicode.RangeTable{
		R32: []unicode.Range32{
			{0x20000, 0x2ffff, 1},
		},
	}
)

var rangeTableNames = map[string]*unicode.RangeTable{
	"BMP": &UnicodeBMP,
	"SMP": &UnicodeSMP,
	"SIP": &UnicodeSIP,
}

func lookupRangeTableName(name string) (*unicode.RangeTable, bool) {
	if r, ok := rangeTableNames[name]; ok {
		return r, true
	}
	if strings.HasPrefix(name, "Category-") && len(name) > 9 {
		return lookupUnicodeRangeTable(name[9:], unicode.Categories)
	} else if strings.HasPrefix(name, "Script-") && len(name) > 7 {
		return lookupUnicodeRangeTable(name[7:], unicode.Scripts)
	} else if strings.HasPrefix(name, "Property-") && len(name) > 9 {
		return lookupUnicodeRangeTable(name[9:], unicode.Properties)
	} else if strings.HasPrefix(name, "FoldCategory-") && len(name) > 13 {
		return lookupUnicodeRangeTable(name[13:], unicode.FoldCategory)
	} else if strings.HasPrefix(name, "FoldScript-") && len(name) > 11 {
		return lookupUnicodeRangeTable(name[11:], unicode.FoldScript)
	}
	return nil, false
}

func lookupUnicodeRangeTable(name string, in map[string]*unicode.RangeTable) (*unicode.RangeTable, bool) {
	r, ok := in[name]
	return r, ok
}

var (
	// used by cmp.Equal to test equality of JSON numerics (i.e. float64, in and json.Number)
	jsonNumericsOption = cmp.FilterValues(jsonNumericCompareFilter, cmp.Comparer(jsonNumericComparator))
)

const (
	// used by cmp.Equal and jsonNumericsOption - to make equality check deterministic and symmetric
	numericJsonNumber = 1
	numericFloat      = 2
	numericInt        = 3
)

func getOtherProperty(propertyName string, vcx *ValidatorContext) (interface{}, bool) {
	ancestryLevel, down := getOtherPropertyPath(propertyName)
	if from, idx, ok := vcx.ancestorValueAndIndex(ancestryLevel); ok {
		return propertyWalkDown(from, down, 0, idx)
	}
	return nil, false
}

func propertyWalkDown(from interface{}, down []string, on int, currIdx int) (interface{}, bool) {
	switch oa := from.(type) {
	case map[string]interface{}:
		return propertyWalkDownMap(oa, down, on, currIdx)
	case []interface{}:
		return propertyWalkDownSlice(oa, down, on, currIdx)
	}
	return nil, false
}

func propertyWalkDownMap(from map[string]interface{}, down []string, on int, currIdx int) (interface{}, bool) {
	if v, ok := from[down[on]]; ok {
		if on == len(down)-1 {
			return v, true
		} else {
			return propertyWalkDown(v, down, on+1, currIdx)
		}
	}
	return nil, false
}

func propertyWalkDownSlice(from []interface{}, down []string, on int, currIdx int) (interface{}, bool) {
	pty := down[on]
	if strings.HasPrefix(pty, "[") && strings.HasSuffix(pty, "]") {
		pty = pty[1 : len(pty)-1]
		if idx, err := strconv.Atoi(pty); err == nil {
			if currIdx != -1 && on == 0 && (strings.HasPrefix(pty, "+") || strings.HasPrefix(pty, "-")) {
				idx = currIdx + idx
			}
			if idx >= 0 && idx < len(from) {
				v := from[idx]
				if on == len(down)-1 {
					return v, true
				} else {
					return propertyWalkDown(v, down, on+1, currIdx)
				}
			}
		}
	}
	return nil, false
}

func getOtherPropertyPath(propertyName string) (uint, []string) {
	level, path := splicePropertyPath(propertyName)
	// fix up any array index indicators...
	actualPath := make([]string, 0, len(path))
	for _, part := range path {
		if strings.HasSuffix(part, "]") {
			if at := strings.Index(part, "["); at == 0 {
				actualPath = append(actualPath, part)
			} else if at != -1 {
				actualPath = append(actualPath, part[:at], part[at:])
			}
		} else {
			actualPath = append(actualPath, part)
		}
	}
	return level, actualPath
}

func splicePropertyPath(propertyName string) (uint, []string) {
	if !strings.Contains(propertyName, ".") {
		return 0, []string{propertyName}
	}
	pth := strings.Split(propertyName, ".")
	if !strings.HasPrefix(propertyName, ".") {
		return 0, pth
	}
	level := -1
	down := make([]string, 0, len(pth))
	atStart := true
	for _, part := range pth {
		if part == "" && atStart {
			level++
		} else {
			atStart = false
			down = append(down, part)
		}
	}
	return uint(level), down
}

func getOtherPropertyDatetime(propertyName string, vcx *ValidatorContext, truncTime bool, allowNull bool) (*time.Time, bool) {
	if other, ok := getOtherProperty(propertyName, vcx); ok {
		if other == nil {
			return nil, allowNull
		}
		if dt, ok := isTime(other, truncTime); ok {
			return &dt, true
		}
	}
	return nil, false
}

func getOtherPropertyString(propertyName string, vcx *ValidatorContext) (string, bool) {
	if other, ok := getOtherProperty(propertyName, vcx); ok {
		if str, sok := other.(string); sok {
			return str, true
		}
	}
	return "", false
}

func stringCompare(a string, b string, caseInsensitive bool) int {
	if caseInsensitive {
		return strings.Compare(strings.ToLower(a), strings.ToLower(b))
	}
	return strings.Compare(a, b)
}

func isTime(v interface{}, truncTime bool) (time.Time, bool) {
	switch dv := v.(type) {
	case *time.Time:
		return truncateTime(*dv, truncTime), true
	case time.Time:
		return truncateTime(dv, truncTime), true
	case *Time:
		return truncateTime(dv.Time, truncTime), true
	case Time:
		return truncateTime(dv.Time, truncTime), true
	case string:
		if dt, ok := stringToDatetime(dv, truncTime); ok {
			return *dt, true
		}
	case *string:
		if dt, ok := stringToDatetime(*dv, truncTime); ok {
			return *dt, true
		}
	}
	return time.Time{}, false
}

func compareNumerics(v1, v2 interface{}) (int, bool) {
	if f1, ok, _ := coerceToFloat(v1); ok {
		if f2, ok, _ := coerceToFloat(v2); ok {
			if f1 > f2 {
				return 1, true
			} else if f1 < f2 {
				return -1, true
			}
			return 0, true
		}
	}
	return 0, false
}

func typedEquals(v1, v2 interface{}) bool {
	if v1 == nil && v2 == nil {
		return true
	} else if v1 == nil || v2 == nil {
		return false
	}
	return cmp.Equal(v1, v2, jsonNumericsOption)
}

func jsonNumericCompareFilter(v1, v2 interface{}) bool {
	if isN1, n1Type := isNumericType(v1); isN1 {
		if isN2, n2type := isNumericType(v2); isN2 {
			return n1Type == numericJsonNumber || n1Type != n2type
		}
	}
	return false
}

func jsonNumericComparator(av1, av2 interface{}) bool {
	v1 := av1
	v2 := av2
	_, vt1 := isNumericType(av1)
	_, vt2 := isNumericType(av2)
	if vt2 < vt1 {
		// make this symmetric & deterministic...
		v1 = av2
		v2 = av1
		swap := vt1
		vt1 = vt2
		vt2 = swap
	}
	switch vt1 {
	case numericJsonNumber:
		jn1 := v1.(json.Number)
		switch vt2 {
		case numericJsonNumber:
			if i1, ok1 := coerceJsonNumberToInt(jn1); ok1 {
				if i2, ok2 := coerceJsonNumberToInt(v2.(json.Number)); ok2 {
					return i1 == i2
				}
			}
			if f1, ok1 := coerceJsonNumberToFloat(jn1); ok1 {
				if f2, ok2 := coerceJsonNumberToFloat(v2.(json.Number)); ok2 {
					return f1 == f2
				}
			}
		case numericFloat:
			if f1, ok1 := coerceJsonNumberToFloat(jn1); ok1 {
				return f1 == v2.(float64)
			}
		case numericInt:
			if i1, ok1 := coerceJsonNumberToInt(jn1); ok1 {
				return i1 == int64(v2.(int))
			}
		}
	case numericFloat:
		f1 := v1.(float64)
		switch vt2 {
		case numericInt:
			return f1 == float64(v2.(int))
		}
	}
	return false
}

func isNumericType(v interface{}) (bool, int) {
	switch v.(type) {
	case json.Number:
		return true, numericJsonNumber
	case float64:
		return true, numericFloat
	case int:
		return true, numericInt
	}
	return false, 0
}

func coerceJsonNumberToInt(jn json.Number) (int64, bool) {
	if i, err := jn.Int64(); err != nil {
		// the int parse of number failed, but that doesn't mean the number isn't an int...
		if f, ok := coerceJsonNumberToFloat(jn); ok {
			if ft := math.Trunc(f); ft == f {
				return int64(ft), true
			}
		}
	} else {
		return i, true
	}
	return 0, false
}

func coerceJsonNumberToFloat(jn json.Number) (float64, bool) {
	if f, err := jn.Float64(); err == nil {
		return f, true
	}
	return 0, false
}

// used by ArrayUnique to check uniqueness of items
func isUniqueCompare(v interface{}, ignoreCase bool, list *[]interface{}) bool {
	result := true
	useV := v
	if ignoreCase {
		if vs, ok := v.(string); ok {
			useV = strings.ToLower(vs)
		}
	}
	for _, ov := range *list {
		if cmp.Equal(ov, useV, jsonNumericsOption) {
			result = false
			break
		}
	}
	if result {
		*list = append(*list, useV)
	}
	return result
}

// Duration represents a parsed ISO8601 duration
type Duration struct {
	Negative bool
	Years    *float64
	Months   *float64
	Weeks    *float64
	Days     *float64
	Hours    *float64
	Minutes  *float64
	Seconds  *float64
}

const (
	durationNumPatt  = `(?:[+-]?)\d+(?:[.,]\d+)?`
	durationFullPatt = `^(?P<negative>[\+\-])?P(?:(?P<year>` + durationNumPatt + `)Y)?(?:(?P<month>` + durationNumPatt + `)M)?(?:(?P<day>` + durationNumPatt + `)D)?` +
		`(?:T(?:(?P<hour>` + durationNumPatt + `)H)?(?:(?P<minute>` + durationNumPatt + `)M)?(?:(?P<second>` + durationNumPatt + `)S)?)?$`
	durationWeekPatt = `^(?P<negative>[\+\-])?P(?:(?P<week>` + durationNumPatt + `)W)$`
)

var (
	iso8601DurationFullRegexp = regexp.MustCompile(durationFullPatt)
	iso8601DurationWeekRegexp = regexp.MustCompile(durationWeekPatt)
)

func ParseDuration(str string) (*Duration, bool) {
	if matches := iso8601DurationWeekRegexp.FindStringSubmatch(str); matches != nil {
		return parseDurationParts(iso8601DurationWeekRegexp, matches)
	} else if matches := iso8601DurationFullRegexp.FindStringSubmatch(str); matches != nil {
		return parseDurationParts(iso8601DurationFullRegexp, matches)
	}
	return nil, false
}

func parseDurationParts(rx *regexp.Regexp, matches []string) (*Duration, bool) {
	result := &Duration{}
	ok := true
	some := false
	for i, nm := range rx.SubexpNames() {
		pt := matches[i]
		if i == 0 || nm == "" || pt == "" {
			continue
		}
		switch nm {
		case "negative":
			result.Negative = pt == "-"
		case "year":
			result.Years, ok = parseDurationNumber(pt)
			some = true
		case "month":
			result.Months, ok = parseDurationNumber(pt)
			some = true
		case "week":
			result.Weeks, ok = parseDurationNumber(pt)
			some = true
		case "day":
			result.Days, ok = parseDurationNumber(pt)
			some = true
		case "hour":
			result.Hours, ok = parseDurationNumber(pt)
			some = true
		case "minute":
			result.Minutes, ok = parseDurationNumber(pt)
			some = true
		case "second":
			result.Seconds, ok = parseDurationNumber(pt)
			some = true
		}
	}
	if ok && some {
		return result, true
	}
	return nil, false
}

func parseDurationNumber(str string) (fv *float64, result bool) {
	if f, err := strconv.ParseFloat(str, 64); err == nil {
		fv = &f
		result = true
	} else if strings.Contains(str, ",") {
		if f, err := strconv.ParseFloat(strings.ReplaceAll(str, ",", "."), 64); err == nil {
			fv = &f
			result = true
		}
	}
	return
}

type stringConstraint interface {
	Constraint
	checkString(str string, vcx *ValidatorContext) bool
}

func checkStringConstraint(v interface{}, vcx *ValidatorContext, c stringConstraint, strict bool, stop bool) (bool, string) {
	if str, ok := v.(string); ok {
		if !c.checkString(str, vcx) {
			vcx.CeaseFurtherIf(stop)
			return false, c.GetMessage(vcx)
		}
	} else if strict {
		vcx.CeaseFurtherIf(stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

func getUnicodeNormalisationForm(f string) (form norm.Form, ok bool) {
	ok = true
	switch strings.ToUpper(f) {
	case "NFC":
		form = norm.NFC
	case "NFKC":
		form = norm.NFKC
	case "NFD":
		form = norm.NFD
	case "NFKD":
		form = norm.NFKD
	default:
		ok = false
	}
	return
}

type dateCompareConstraint interface {
	Constraint
	compareDates(value time.Time, other time.Time) bool
}

func checkDateCompareConstraint(value string, other interface{}, vcx *ValidatorContext, c dateCompareConstraint, excTime bool, stop bool) (bool, string) {
	if cdt, ok := stringToDatetime(value, excTime); ok {
		if dt, ok := isTime(other, excTime); ok && c.compareDates(*cdt, dt) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(stop)
	return false, c.GetMessage(vcx)
}

func checkDateCompareOtherPropertyConstraint(v interface{}, otherPtyName string, vcx *ValidatorContext, c dateCompareConstraint, excTime bool, stop bool) (bool, string) {
	if other, ok := getOtherPropertyDatetime(otherPtyName, vcx, excTime, false); ok {
		if dt, ok := isTime(v, excTime); ok && c.compareDates(dt, *other) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(stop)
	return false, c.GetMessage(vcx)
}

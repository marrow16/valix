package valix

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// StringValidISODatetime constraint checks that a string value is a valid ISO8601 Date/time format
type StringValidISODatetime struct {
	// specifies, if set to true, that time offsets are not permitted
	NoOffset bool
	// specifies, if set to true, that seconds cannot have decimal places
	NoMillis bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidISODatetime) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		useRegex := iso8601FullRegex
		useLayout := iso8601FullLayout
		if c.NoOffset && c.NoMillis {
			useRegex = iso8601MinRegex
			useLayout = iso8601MinLayout
		} else if c.NoOffset {
			useRegex = iso8601NoOffsRegex
			useLayout = iso8601NoOffLayout
		} else if c.NoMillis {
			useRegex = iso8601NoMillisRegex
			useLayout = iso8601NoMillisLayout
		}
		if !useRegex.MatchString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		// and attempt to parse it (it may match regex but datetime could still be invalid)...
		if _, err := time.Parse(useLayout, str); err != nil {
			if pErr, ok := err.(*time.ParseError); ok && pErr.Message == "" && strings.HasSuffix(useLayout, pErr.LayoutElem) {
				// time.Parse is pretty dumb when it comes to timezones - if it's in the layout but not the string it fails
				// so remove the bit it doesn't like (in the layout) and try again...
				useLayout = useLayout[0 : len(useLayout)-len(pErr.LayoutElem)]
				if _, err := time.Parse(useLayout, str); err != nil {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			} else {
				vcx.CeaseFurtherIf(c.Stop)
				return false, c.GetMessage(vcx)
			}
		}
	} else if c.Strict {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODatetime) GetMessage(tcx I18nContext) string {
	if c.NoOffset && c.NoMillis {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatMin)
	} else if c.NoOffset {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatNoOffs)
	} else if c.NoMillis {
		return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatNoMillis)
	}
	return defaultMessage(tcx, c.Message, msgValidISODatetimeFormatFull)
}

// StringValidISODate constraint checks that a string value is a valid ISO8601 Date format (excluding time)
type StringValidISODate struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
	// when set to true, fails if the value being checked is not a correct type
	Strict bool
}

// Check implements Constraint.Check
func (c *StringValidISODate) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if !iso8601DateOnlyRegex.MatchString(str) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
		// and attempt to parse it (it may match regex but date could still be invalid)...
		if _, err := time.Parse(iso8601DateOnlyLayout, str); err != nil {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if c.Strict {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODate) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidISODate)
}

// StringValidTimezone constraint checks that a string value is a valid timezone
//
// NB. If both LocationOnly and OffsetOnly are set to true - this constraint will always fail!
type StringValidTimezone struct {
	// allows location only
	LocationOnly bool
	// allows offset only
	OffsetOnly bool
	// if set, allows offset to be a numeric value
	AllowNumeric bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

var timeoffsetRegex = regexp.MustCompile("^[+-]?\\d{1,2}(:\\d{2})?$")

// Check implements Constraint.Check
func (c *StringValidTimezone) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if !(c.LocationOnly && c.OffsetOnly) {
		if str, ok := v.(string); ok && str != "" && strings.ToLower(str) != "local" {
			if c.checkString(str) {
				return true, ""
			}
		} else if c.AllowNumeric && !c.LocationOnly {
			if iv, ok, isNumber := coerceToInt(v); ok && isNumber {
				if iv >= -12 && iv <= 14 {
					return true, ""
				}
			}
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *StringValidTimezone) checkString(str string) bool {
	if !c.LocationOnly && timeoffsetRegex.MatchString(str) {
		var h, m int
		if cAt := strings.IndexByte(str, ':'); cAt != -1 {
			h, _ = strconv.Atoi(str[:cAt])
			m, _ = strconv.Atoi(str[cAt+1:])
		} else {
			h, _ = strconv.Atoi(str)
		}
		if h >= -12 && h <= 14 && m < 60 {
			return true
		}
	}
	if !c.OffsetOnly {
		if _, err := time.LoadLocation(str); err == nil {
			return true
		}
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidTimezone) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidTimezone)
}

// DatetimeDayOfWeek checks that a date (represented as string or time.Time) is an allowed day of the week
type DatetimeDayOfWeek struct {
	// is the allowed days (of the week) expressed as a string of allowed week day numbers (in any order)
	//
	// Where 0 = Sunday, e.g. "06" (or "60") allows Sunday or Saturday
	//
	// or to allow only 'working days' of the week - "12345"
	Days string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeDayOfWeek) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, false); ok {
		if strings.Contains(c.Days, strconv.Itoa(int(dt.Weekday()))) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeDayOfWeek) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimeDayOfWeek)
}

// DatetimeFuture constraint checks that a datetime/date (represented as string or time.Time) is in the future
type DatetimeFuture struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeFuture) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, c.ExcTime); ok {
		if dt.After(truncateTime(time.Now(), c.ExcTime)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFuture) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimeFuture)
}

// DatetimeFutureOrPresent constraint checks that a datetime/date (represented as string or time.Time) is in the future or present
type DatetimeFutureOrPresent struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeFutureOrPresent) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, c.ExcTime); ok {
		if dt.After(truncateTime(time.Now(), c.ExcTime)) || dt.Equal(truncateTime(time.Now(), c.ExcTime)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeFutureOrPresent) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimeFutureOrPresent)
}

// DatetimePast constraint checks that a datetime/date (represented as string or time.Time) is in the past
type DatetimePast struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimePast) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, c.ExcTime); ok {
		if dt.Before(truncateTime(time.Now(), c.ExcTime)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePast) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimePast)
}

// DatetimePastOrPresent constraint checks that a datetime/date (represented as string or time.Time) is in the past or present
type DatetimePastOrPresent struct {
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimePastOrPresent) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, c.ExcTime); ok {
		now := truncateTime(time.Now(), c.ExcTime)
		if dt.Before(now) || dt.Equal(now) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimePastOrPresent) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgDatetimePastOrPresent)
}

// DatetimeRange constraint checks that a datetime/date (represented as string or time.Time) is within a specified range
type DatetimeRange struct {
	// is the minimum datetime/date (if this is empty, then no minimum check is performed)
	Minimum string
	// is the maximum datetime/date (if this is empty, then no maximum check is performed)
	Maximum string
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeRange) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, c.ExcTime); ok {
		if c.checkMin(dt) && c.checkMax(dt) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *DatetimeRange) checkMin(dt time.Time) bool {
	if c.Minimum == "" {
		return true
	}
	if cdt, ok := stringToDatetime(c.Minimum, c.ExcTime); ok {
		return dt.After(*cdt) || (!c.ExclusiveMin && dt.Equal(*cdt))
	}
	return false
}

func (c *DatetimeRange) checkMax(dt time.Time) bool {
	if c.Maximum == "" {
		return true
	}
	if cdt, ok := stringToDatetime(c.Maximum, c.ExcTime); ok {
		return dt.Before(*cdt) || (!c.ExclusiveMax && dt.Equal(*cdt))
	}
	return false
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeRange) GetMessage(tcx I18nContext) string {
	if c.Message != "" {
		return c.Message
	} else if c.Minimum != "" && c.Maximum != "" {
		return defaultMessage(tcx, c.Message, fmtMsgRange, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
	} else if c.Minimum != "" {
		return defaultMessage(tcx, "", ternary(c.ExclusiveMin).string(fmtMsgDtGt, fmtMsgDtGte), c.Minimum)
	} else if c.Maximum != "" {
		return defaultMessage(tcx, "", ternary(c.ExclusiveMax).string(fmtMsgDtLt, fmtMsgDtLte), c.Maximum)
	}
	return msgValidISODatetimeFormatFull
}

// DatetimeTimeOfDayRange constraint checks that a datetime (represented as string or time.Time) is within a specified time of day range
type DatetimeTimeOfDayRange struct {
	// is the minimum time of day (if this is empty, then no minimum check is performed)
	Minimum string
	// is the maximum time of day (if this is empty, then no maximum check is performed)
	Maximum string
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeTimeOfDayRange) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if dt, ok := isTime(v, false); ok {
		if c.checkMin(dt) && c.checkMax(dt) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *DatetimeTimeOfDayRange) checkMin(dt time.Time) bool {
	if c.Minimum == "" {
		return true
	}
	if cdt, ok := stringToDatetime(c.timeToDatetime(c.Minimum), false); ok {
		return timeOfDayCompare(dt, *cdt) > ternary(c.ExclusiveMin).int(0, -1)
	}
	return false
}

func (c *DatetimeTimeOfDayRange) checkMax(dt time.Time) bool {
	if c.Maximum == "" {
		return true
	}
	if cdt, ok := stringToDatetime(c.timeToDatetime(c.Maximum), false); ok {
		return timeOfDayCompare(dt, *cdt) < ternary(c.ExclusiveMax).int(0, 1)
	}
	return false
}

func timeOfDayCompare(a, b time.Time) int {
	if hc := intCompare(a.Hour(), b.Hour()); hc == 0 {
		if mc := intCompare(a.Minute(), b.Minute()); mc == 0 {
			if sc := intCompare(a.Second(), b.Second()); sc == 0 {
				return intCompare(a.Nanosecond(), b.Nanosecond())
			} else {
				return sc
			}
		} else {
			return mc
		}
	} else {
		return hc
	}
}

func (c *DatetimeTimeOfDayRange) timeToDatetime(t string) string {
	return "2000-01-01T" + c.timePadRight(t, false)
}

func (c *DatetimeTimeOfDayRange) timePadRight(t string, printable bool) string {
	l := len(t)
	if l >= 8 || (l == 5 && printable) {
		return t
	} else if l == 1 {
		return "0" + t + (":00:00" + ternary(printable).string("", ".000000000"))[l-1:]
	}
	return t + ("00:00:00" + ternary(printable).string("", ".000000000"))[l:]
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeTimeOfDayRange) GetMessage(tcx I18nContext) string {
	if c.Message != "" {
		return c.Message
	} else if c.Minimum != "" && c.Maximum != "" {
		return defaultMessage(tcx, c.Message, fmtMsgRange, c.timePadRight(c.Minimum, true), incExc(tcx, c.ExclusiveMin), c.timePadRight(c.Maximum, true), incExc(tcx, c.ExclusiveMax))
	} else if c.Minimum != "" {
		return defaultMessage(tcx, "", ternary(c.ExclusiveMin).string(fmtMsgDtGt, fmtMsgDtGte), c.timePadRight(c.Minimum, true))
	} else if c.Maximum != "" {
		return defaultMessage(tcx, "", ternary(c.ExclusiveMax).string(fmtMsgDtLt, fmtMsgDtLte), c.timePadRight(c.Maximum, true))
	}
	return msgValidISODatetimeFormatFull
}

// DatetimeYearsOld constraint checks that a date (datetime represented as string or time.Time) meets the specified
// minimum and/or maximum years-old.  Can also be used to simply check a minimum age or maximum age
//
// Notes:
//
// * If the value being checked contains a time (hh:mm:ss), it is ignored (very few people know, or are expected to
// specify, their exact time of birth)
//
// * If the value being checked is in the future - this constraint fails
//
// * If both Minimum and Maximum are set to zero (or less) then no check is performed
type DatetimeYearsOld struct {
	// is the minimum age (not checked if this value is zero or less)
	Minimum int
	// is the maximum age (not checked if this value is zero or less)
	Maximum int
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// if set to true, ExclusiveMax specifies the maximum value is exclusive
	ExclusiveMax bool
	// if set to true, only checks the minimum/maximum age against the current year - i.e. the current age is calculated
	// based on 23:59:59.999999999 at 31st December of the current year
	ThisYear bool
	// is an optional string representing a threshold date at which the age is calculated
	//
	// If this is specified, the year part is ignored (the current year is always used)
	//
	// Note: if specified, this also overrides the ThisYear flag
	ThresholdDate string
	// if set, adjusts the way leapday birthdays are age calculated
	//
	// By default, leapday birthdays are taken as 1st March when the current year is not a leap year
	//
	// Setting LeapdayAdjust to true means that leapday birthdays are taken as 28th Feb
	LeapdayAdjust bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeYearsOld) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.Minimum <= 0 && c.Maximum <= 0 {
		// if neither min nor max are set, then no checks...
		return true, ""
	}
	if dob, ok := isTime(v, true); ok {
		age := c.calculateAge(time.Now(), dob)
		if age >= 0 && (c.Minimum <= 0 || age > c.Minimum || (!c.ExclusiveMin && age == c.Minimum)) &&
			(c.Maximum <= 0 || age < c.Maximum || (!c.ExclusiveMax && age == c.Maximum)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

func (c *DatetimeYearsOld) calculateAge(now, dob time.Time) int {
	useNow := now
	if c.ThresholdDate != "" {
		if td, ok := stringToDatetime(c.ThresholdDate, false); ok {
			useNow = time.Date(now.Year(), td.Month(), td.Day(), td.Hour(), td.Minute(), td.Second(), td.Nanosecond(), time.UTC)
		}
	} else if c.ThisYear {
		useNow = time.Date(now.Year(), 12, 31, 23, 59, 59, 999999999, time.UTC)
	}
	age := useNow.Year() - dob.Year()
	dobM := dob.Month()
	dobD := dob.Day()
	if dobM == 2 && dobD == 29 {
		dob = time.Date(now.Year(), dobM, dobD, 0, 0, 0, 0, time.UTC)
		dobM = dob.Month()
		dobD = dob.Day()
		if dobM != 29 && c.LeapdayAdjust {
			dobM = 2
			dobD = 29 // even though 29th Feb may not be valid for this year
		}
	}
	if useNow.Month() < dobM || (useNow.Month() == dobM && useNow.Day() < dobD) {
		age--
	}
	return age
}

const (
	fmtMsgDtAgeMin          = "Age must be over %[1]d years old"
	fmtMsgDtAgeMinOrOver    = "Age must be %[1]d years old or over"
	fmtMsgDtAgeMax          = "Age must be under %[1]d years old"
	fmtMsgDtAgeMaxOrUnder   = "Age must be %[1]d years old or under"
	fmtMsgDtAgeMinExcMaxExc = "Age must be over %[1]d years old and under %[2]d years old"
	fmtMsgDtAgeMinMax       = "Age must be between %[1]d years old and %[2]d years old"
	fmtMsgDtAgeMinMaxExc    = "Age must be %[1]d years old or over and under %[2]d years old"
	fmtMsgDtAgeMinExcMax    = "Age must be between over %[1]d years old and %[2]d years old or under"
)

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeYearsOld) GetMessage(tcx I18nContext) string {
	if c.Message != "" || (c.Minimum <= 0 && c.Maximum <= 0) {
		return obtainI18nContext(tcx).TranslateMessage(c.Message)
	} else if c.Minimum > 0 && c.Maximum > 0 {
		if c.ExclusiveMin && c.ExclusiveMax {
			return defaultMessage(tcx, "", fmtMsgDtAgeMinExcMaxExc, c.Minimum, c.Maximum)
		} else if !c.ExclusiveMin && !c.ExclusiveMax {
			return defaultMessage(tcx, "", fmtMsgDtAgeMinMax, c.Minimum, c.Maximum)
		} else if c.ExclusiveMax {
			return defaultMessage(tcx, "", fmtMsgDtAgeMinMaxExc, c.Minimum, c.Maximum)
		}
		return defaultMessage(tcx, "", fmtMsgDtAgeMinExcMax, c.Minimum, c.Maximum)
	} else if c.Minimum > 0 {
		return defaultMessage(tcx, "", ternary(c.ExclusiveMin).string(fmtMsgDtAgeMin, fmtMsgDtAgeMinOrOver), c.Minimum)
	}
	return defaultMessage(tcx, "", ternary(c.ExclusiveMax).string(fmtMsgDtAgeMax, fmtMsgDtAgeMaxOrUnder), c.Maximum)

}

// StringValidISODuration constraint checks that a string value is a valid ISO8601 Duration
type StringValidISODuration struct {
	// if set, disallows negative durations (e.g. "-P1Y")
	DisallowNegative bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *StringValidISODuration) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := v.(string); ok {
		if dur, ok := ParseDuration(str); ok && (!c.DisallowNegative || !dur.Negative) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *StringValidISODuration) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgValidISODuration)
}

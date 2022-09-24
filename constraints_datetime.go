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

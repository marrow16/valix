package valix

import (
	"time"
)

// EqualsOther constraint to check that a property value equals the value of another named property
type EqualsOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *EqualsOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if typedEquals(v, other) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *EqualsOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgEqualsOther, c.PropertyName)
}

// NotEqualsOther constraint to check that a property value not equals the value of another named property
type NotEqualsOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NotEqualsOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if typedEquals(v, other) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NotEqualsOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgNotEqualsOther, c.PropertyName)
}

// GreaterThan constraint to check that a numeric value is greater than a specified value
//
// Note: This constraint is stricter than Minimum, Maximum and Range constraints in that if
// the property value is not a numeric then this constraint fails
type GreaterThan struct {
	// the value to compare against
	Value float64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *GreaterThan) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if comp, ok := compareNumerics(c.Value, v); ok && comp < 0 {
		return true, ""
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *GreaterThan) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGt, c.Value)
}

// GreaterThanOrEqual constraint to check that a numeric value is greater than or equal to a specified value
//
// Note: This constraint is stricter than Minimum, Maximum and Range constraints in that if
// the property value is not a numeric then this constraint fails
type GreaterThanOrEqual struct {
	// the value to compare against
	Value float64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *GreaterThanOrEqual) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if comp, ok := compareNumerics(c.Value, v); ok && comp <= 0 {
		return true, ""
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *GreaterThanOrEqual) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGte, c.Value)
}

// LessThan constraint to check that a numeric value is less than a specified value
//
// Note: This constraint is stricter than Minimum, Maximum and Range constraints in that if
// the property value is not a numeric then this constraint fails
type LessThan struct {
	// the value to compare against
	Value float64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LessThan) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if comp, ok := compareNumerics(c.Value, v); ok && comp > 0 {
		return true, ""
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *LessThan) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLt, c.Value)
}

// LessThanOrEqual constraint to check that a numeric value is less than or equal to a specified value
//
// Note: This constraint is stricter than Minimum, Maximum and Range constraints in that if
// the property value is not a numeric then this constraint fails
type LessThanOrEqual struct {
	// the value to compare against
	Value float64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LessThanOrEqual) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if comp, ok := compareNumerics(c.Value, v); ok && comp >= 0 {
		return true, ""
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *LessThanOrEqual) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLte, c.Value)
}

// GreaterThanOther constraint to check that a numeric value is greater than another named property value
//
// Note: this constraint is strict - if either the current or other property is not numeric then this
// constraint fails
type GreaterThanOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *GreaterThanOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if comp, ok := compareNumerics(other, v); ok && comp < 0 {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *GreaterThanOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGtOther, c.PropertyName)
}

// GreaterThanOrEqualOther constraint to check that a numeric value is greater than or equal to
// another named property value
//
// Note: this constraint is strict - if either the current or other property is not numeric then this
// constraint fails
type GreaterThanOrEqualOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *GreaterThanOrEqualOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if comp, ok := compareNumerics(other, v); ok && comp <= 0 {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *GreaterThanOrEqualOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGteOther, c.PropertyName)
}

// LessThanOther constraint to check that a numeric value is less than another named property value
//
// Note: this constraint is strict - if either the current or other property is not numeric then this
// constraint fails
type LessThanOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LessThanOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if comp, ok := compareNumerics(other, v); ok && comp > 0 {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *LessThanOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLtOther, c.PropertyName)
}

// LessThanOrEqualOther constraint to check that a numeric value is less than or equal to
// another named property value
//
// Note: this constraint is strict - if either the current or other property is not numeric then this
// constraint fails
type LessThanOrEqualOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LessThanOrEqualOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherProperty(c.PropertyName, vcx); ok {
		if comp, ok := compareNumerics(other, v); ok && comp >= 0 {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *LessThanOrEqualOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLteOther, c.PropertyName)
}

// DatetimeGreaterThan constraint to check that a date/time (as an ISO string) value is greater than a specified value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeGreaterThan struct {
	// the value to compare against (a string representation of date or datetime in ISO format)
	Value string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeGreaterThan) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if cdt, ok := stringToDatetime(c.Value, c.ExcTime); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && dt.After(*cdt) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeGreaterThan) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgDtGt, c.Value)
}

// DatetimeGreaterThanOrEqual constraint to check that a date/time (as an ISO string) value is greater than or equal to a specified value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeGreaterThanOrEqual struct {
	// the value to compare against (a string representation of date or datetime in ISO format)
	Value string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeGreaterThanOrEqual) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if cdt, ok := stringToDatetime(c.Value, c.ExcTime); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && (dt.After(*cdt) || dt.Equal(*cdt)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeGreaterThanOrEqual) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgDtGte, c.Value)
}

// DatetimeLessThan constraint to check that a date/time (as an ISO string) value is less than a specified value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeLessThan struct {
	// the value to compare against (a string representation of date or datetime in ISO format)
	Value string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeLessThan) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if cdt, ok := stringToDatetime(c.Value, c.ExcTime); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && dt.Before(*cdt) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeLessThan) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgDtLt, c.Value)
}

// DatetimeLessThanOrEqual constraint to check that a date/time (as an ISO string) value is less than or equal to a specified value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeLessThanOrEqual struct {
	// the value to compare against (a string representation of date or datetime in ISO format)
	Value string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeLessThanOrEqual) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if cdt, ok := stringToDatetime(c.Value, c.ExcTime); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && (dt.Before(*cdt) || dt.Equal(*cdt)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeLessThanOrEqual) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgDtLte, c.Value)
}

// DatetimeGreaterThanOther constraint to check that a date/time (as an ISO string) value is greater than another named property value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeGreaterThanOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeGreaterThanOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherPropertyDatetime(c.PropertyName, vcx, c.ExcTime, false); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && dt.After(*other) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeGreaterThanOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGtOther, c.PropertyName)
}

// DatetimeGreaterThanOrEqualOther constraint to check that a date/time (as an ISO string) value is greater than or equal to another named property value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeGreaterThanOrEqualOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeGreaterThanOrEqualOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherPropertyDatetime(c.PropertyName, vcx, c.ExcTime, false); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && (dt.After(*other) || dt.Equal(*other)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeGreaterThanOrEqualOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgGteOther, c.PropertyName)
}

// DatetimeLessThanOther constraint to check that a date/time (as an ISO string) value is less than another named property value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeLessThanOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeLessThanOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherPropertyDatetime(c.PropertyName, vcx, c.ExcTime, false); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && dt.Before(*other) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeLessThanOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLtOther, c.PropertyName)
}

// DatetimeLessThanOrEqualOther constraint to check that a date/time (as an ISO string) value is less than or equal to another named property value
//
// Note: this constraint is strict - if either of the compared values is not a valid ISO datetime then this
// constraint fails
type DatetimeLessThanOrEqualOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeLessThanOrEqualOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if other, ok := getOtherPropertyDatetime(c.PropertyName, vcx, c.ExcTime, false); ok {
		if dt, ok := isTime(v, c.ExcTime); ok && (dt.Before(*other) || dt.Equal(*other)) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeLessThanOrEqualOther) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgLteOther, c.PropertyName)
}

// DatetimeTolerance constraint to check that a date/time (as an ISO string) value meets a tolerance against a specified value
//
// Note: this constraint is strict - if the property value is not a valid ISO datetime then
// this constraint fails
type DatetimeTolerance struct {
	// the value to compare against (a string representation of date or datetime in ISO format)
	Value string `v8n:"default"`
	// the tolerance duration amount - which can be positive, negative or zero
	//
	// For negative values, this is the maximum duration into the past
	//
	// For positive values, this is the maximum duration into the future
	//
	// Note: If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	// specified.  For example, if the Duration is zero and the Unit is specified as "year" then this constraint
	// will check the same year
	Duration int64
	// is the string token specifying the unit in which the Duration is measured
	//
	// this can be "millennium", "century", "decade", "year", "month", "week", "day",
	// "hour", "min", "sec" or "milli" (millisecond), "micro" (microsecond) or "nano" (nanosecond)
	//
	// Note: if this is empty, then "day" is assumed.  If the token is invalid - this constraint fails!
	Unit string
	// when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
	MinCheck bool
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// when set to true, IgnoreNull makes the constraint less strict by ignoring null values
	IgnoreNull bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeTolerance) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.IgnoreNull && v == nil {
		return true, ""
	}
	useExcTime := c.ExcTime && c.Duration != 0
	if dt, ok := isTime(v, useExcTime); ok {
		if cdt, ok := stringToDatetime(c.Value, useExcTime); ok && checkDatetimeTolerance(&dt, cdt, c.Duration, c.Unit, c.MinCheck) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeTolerance) GetMessage(tcx I18nContext) string {
	useTcx := obtainI18nContext(tcx)
	if c.Message != "" {
		return useTcx.TranslateMessage(c.Message)
	}
	useUnit := defaultToleranceUnit(c.Unit)
	if c.Duration == 0 {
		return useTcx.TranslateFormat(fmtMsgDtToleranceFixedSame, useTcx.TranslateToken(useUnit), c.Value)
	}
	useAmount := c.Duration
	if useAmount < 0 {
		useAmount = 0 - useAmount
	}
	if useAmount > 1 {
		useUnit = useUnit + "..."
	}
	if c.MinCheck {
		useFmt := fmtMsgDtToleranceFixedMinAfter
		if c.Duration < 0 {
			useFmt = fmtMsgDtToleranceFixedMinBefore
		}
		return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit), c.Value)
	}
	useFmt := fmtMsgDtToleranceFixedMaxAfter
	if c.Duration < 0 {
		useFmt = fmtMsgDtToleranceFixedMaxBefore
	}
	return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit), c.Value)
}

func defaultToleranceUnit(unit string) string {
	if unit == "" {
		return "day"
	}
	return unit
}

const (
	fmtMsgDtToleranceFixedSame      = "Value must be same %[1]s as %[2]s"
	fmtMsgDtToleranceFixedMaxAfter  = "Value must not be more than %[1]d %[2]s after %[3]s"
	fmtMsgDtToleranceFixedMaxBefore = "Value must not be more than %[1]d %[2]s before %[3]s"
	fmtMsgDtToleranceFixedMinAfter  = "Value must be at least %[1]d %[2]s after %[3]s"
	fmtMsgDtToleranceFixedMinBefore = "Value must be at least %[1]d %[2]s before %[3]s"
)

// DatetimeToleranceToNow constraint to check that a date/time (as an ISO string) value meets a tolerance against the current time
//
// Note: this constraint is strict - if the property value is not a valid ISO datetime then
// this constraint fails
type DatetimeToleranceToNow struct {
	// the tolerance duration amount - which can be positive, negative or zero
	//
	// For negative values, this is the maximum duration into the past
	//
	// For positive values, this is the maximum duration into the future
	//
	// Note: If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	// specified.  For example, if the Duration is zero and the Unit is specified as "year" then this constraint
	// will check the same year
	Duration int64
	// is the string token specifying the unit in which the Duration is measured
	//
	// this can be "millennium", "century", "decade", "year", "month", "week", "day",
	// "hour", "min", "sec" or "milli" (millisecond), "micro" (microsecond) or "nano" (nanosecond)
	//
	// Note: if this is empty, then "day" is assumed.  If the token is invalid - this constraint fails!
	Unit string
	// when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
	MinCheck bool
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// when set to true, IgnoreNull makes the constraint less strict by ignoring null values
	IgnoreNull bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeToleranceToNow) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.IgnoreNull && v == nil {
		return true, ""
	}
	useExcTime := c.ExcTime && c.Duration != 0
	if dt, ok := isTime(v, useExcTime); ok {
		now := truncateTime(time.Now(), useExcTime)
		if checkDatetimeTolerance(&now, &dt, c.Duration, c.Unit, c.MinCheck) {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeToleranceToNow) GetMessage(tcx I18nContext) string {
	useTcx := obtainI18nContext(tcx)
	if c.Message != "" {
		return useTcx.TranslateMessage(c.Message)
	}
	useUnit := defaultToleranceUnit(c.Unit)
	if c.Duration == 0 {
		return useTcx.TranslateFormat(fmtMsgDtToleranceNowSame, useTcx.TranslateToken(useUnit))
	}
	useAmount := c.Duration
	if useAmount < 0 {
		useAmount = 0 - useAmount
	}
	if useAmount > 1 {
		useUnit = useUnit + "..."
	}
	if c.MinCheck {
		useFmt := fmtMsgDtToleranceNowMinAfter
		if c.Duration < 0 {
			useFmt = fmtMsgDtToleranceNowMinBefore
		}
		return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit))
	}
	useFmt := fmtMsgDtToleranceNowMaxAfter
	if c.Duration < 0 {
		useFmt = fmtMsgDtToleranceNowMaxBefore
	}
	return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit))
}

const (
	fmtMsgDtToleranceNowSame      = "Value must be same %[1]s as now"
	fmtMsgDtToleranceNowMaxAfter  = "Value must not be more than %[1]d %[2]s after now"
	fmtMsgDtToleranceNowMaxBefore = "Value must not be more than %[1]d %[2]s before now"
	fmtMsgDtToleranceNowMinAfter  = "Value must be at least %[1]d %[2]s after now"
	fmtMsgDtToleranceNowMinBefore = "Value must be at least %[1]d %[2]s before now"
)

// DatetimeToleranceToOther constraint to check that a date/time (as an ISO string) value meets a tolerance against the
// value of another named property value
//
// Note: this constraint is strict - if the property value is not a valid ISO datetime then
// this constraint fails
type DatetimeToleranceToOther struct {
	// the property name of the other value to compare against
	PropertyName string `v8n:"default"`
	// the tolerance duration amount - which can be positive, negative or zero
	//
	// For negative values, this is the maximum duration into the past
	//
	// For positive values, this is the maximum duration into the future
	//
	// Note: If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	// specified.  For example, if the Duration is zero and the Unit is specified as "year" then this constraint
	// will check the same year
	Duration int64
	// is the string token specifying the unit in which the Duration is measured
	//
	// this can be "millennium", "century", "decade", "year", "month", "week", "day",
	// "hour", "min", "sec" or "milli" (millisecond), "micro" (microsecond) or "nano" (nanosecond)
	//
	// Note: if this is empty, then "day" is assumed.  If the token is invalid - this constraint fails!
	Unit string
	// when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
	MinCheck bool
	// when set to true, excludes the time when comparing
	//
	// Note: This also excludes the effect of any timezone offsets specified in either of the compared values
	ExcTime bool
	// when set to true, IgnoreNull makes the constraint less strict by ignoring null values
	//
	// NB. ignoring nulls applies to both the property being checked and the other named property
	IgnoreNull bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *DatetimeToleranceToOther) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if c.IgnoreNull && v == nil {
		return true, ""
	}
	useExcTime := c.ExcTime && c.Duration != 0
	if other, ok := getOtherPropertyDatetime(c.PropertyName, vcx, useExcTime, c.IgnoreNull); ok {
		if other != nil {
			if dt, ok := isTime(v, useExcTime); ok && checkDatetimeTolerance(&dt, other, c.Duration, c.Unit, c.MinCheck) {
				return true, ""
			}
		} else if c.IgnoreNull {
			return true, ""
		}
	}
	vcx.CeaseFurtherIf(c.Stop)
	return false, c.GetMessage(vcx)
}

// GetMessage implements the Constraint.GetMessage
func (c *DatetimeToleranceToOther) GetMessage(tcx I18nContext) string {
	useTcx := obtainI18nContext(tcx)
	if c.Message != "" {
		return useTcx.TranslateMessage(c.Message)
	}
	useUnit := defaultToleranceUnit(c.Unit)
	if c.Duration == 0 {
		return useTcx.TranslateFormat(fmtMsgDtToleranceOtherSame, useTcx.TranslateToken(useUnit), c.PropertyName)
	}
	useAmount := c.Duration
	if useAmount < 0 {
		useAmount = 0 - useAmount
	}
	if useAmount > 1 {
		useUnit = useUnit + "..."
	}
	if c.MinCheck {
		useFmt := fmtMsgDtToleranceOtherMinAfter
		if c.Duration < 0 {
			useFmt = fmtMsgDtToleranceOtherMinBefore
		}
		return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit), c.PropertyName)
	}
	useFmt := fmtMsgDtToleranceOtherMaxAfter
	if c.Duration < 0 {
		useFmt = fmtMsgDtToleranceOtherMaxBefore
	}
	return useTcx.TranslateFormat(useFmt, useAmount, useTcx.TranslateToken(useUnit), c.PropertyName)
}

const (
	fmtMsgDtToleranceOtherSame      = "Value must be same %[1]s as value of property '%[2]s'"
	fmtMsgDtToleranceOtherMaxAfter  = "Value must not be more than %[1]d %[2]s after value of property '%[3]s'"
	fmtMsgDtToleranceOtherMaxBefore = "Value must not be more than %[1]d %[2]s before value of property '%[3]s'"
	fmtMsgDtToleranceOtherMinAfter  = "Value must be at least %[1]d %[2]s after value of property '%[3]s'"
	fmtMsgDtToleranceOtherMinBefore = "Value must be at least %[1]d %[2]s before value of property '%[3]s'"
)

func checkDatetimeTolerance(value *time.Time, other *time.Time, amount int64, unit string, minCheck bool) bool {
	if amount == 0 {
		return checkDatetimeSame(value, other, unit)
	}
	shifted, ok := shiftDatetimeBy(value, amount, unit)
	if !ok {
		return false
	}
	result := false
	if minCheck {
		if shifted.Equal(*other) {
			// same - so min tolerance cannot be ok...
			result = false
		} else if amount < 0 {
			result = shifted.After(*other)
		} else {
			result = shifted.Before(*other)
		}
	} else {
		if shifted.Equal(*other) {
			// same - so max tolerance must be ok...
			result = true
		} else if amount < 0 {
			result = !shifted.After(*other)
		} else {
			result = !shifted.Before(*other)
		}
	}
	return result
}

func checkDatetimeSame(value *time.Time, other *time.Time, unit string) bool {
	switch unit {
	case "millennium", "millen", "mille":
		return (value.Year() / 1000) == (other.Year() / 1000)
	case "century":
		return (value.Year() / 100) == (other.Year() / 100)
	case "decade":
		return (value.Year() / 10) == (other.Year() / 10)
	case "year":
		return value.Year() == other.Year()
	case "month":
		return value.Year() == other.Year() && value.Month() == other.Month()
	case "day", "":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day()
	case "week":
		vy, vwk := value.ISOWeek()
		oy, owk := other.ISOWeek()
		return vy == oy && vwk == owk
	case "hour":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour()
	case "min", "minute":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour() && value.Minute() == other.Minute()
	case "sec", "second":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour() && value.Minute() == other.Minute() && value.Second() == other.Second()
	case "milli", "millisecond":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour() && value.Minute() == other.Minute() && value.Second() == other.Second() &&
			(value.Nanosecond()/1000000) == (other.Nanosecond()/1000000)
	case "micro", "microsecond":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour() && value.Minute() == other.Minute() && value.Second() == other.Second() &&
			(value.Nanosecond()/1000) == (other.Nanosecond()/1000)
	case "nano", "nanosecond":
		return value.Year() == other.Year() && value.Month() == other.Month() && value.Day() == other.Day() &&
			value.Hour() == other.Hour() && value.Minute() == other.Minute() && value.Second() == other.Second() &&
			value.Nanosecond() == other.Nanosecond()
	}
	return false
}

func shiftDatetimeBy(t *time.Time, amount int64, unit string) (*time.Time, bool) {
	switch unit {
	case "millennium", "millen", "mille":
		return shiftDatetimeByYears(t, amount*1000)
	case "century":
		return shiftDatetimeByYears(t, amount*100)
	case "decade":
		return shiftDatetimeByYears(t, amount*10)
	case "year":
		return shiftDatetimeByYears(t, amount)
	case "month":
		return shiftDatetimeByMonths(t, amount)
	}
	unitAmount, ok := timeToleranceUnits[unit]
	if !ok {
		return nil, false
	}
	actualAmount := amount * unitAmount
	result := t.Add(time.Duration(actualAmount))
	return &result, true
}

func shiftDatetimeByYears(t *time.Time, amount int64) (*time.Time, bool) {
	result := time.Date(t.Year()+int(amount), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	if result.Day() != t.Day() {
		startDay := result.Day()
		for {
			result = result.Add(0 - (time.Hour * 24))
			if result.Day() > startDay {
				break
			}
		}
	}
	return &result, true
}

func shiftDatetimeByMonths(t *time.Time, amount int64) (*time.Time, bool) {
	result := time.Date(t.Year(), time.Month(int(t.Month())+int(amount)), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), t.Location())
	if result.Day() != t.Day() {
		startDay := result.Day()
		for {
			result = result.Add(0 - (time.Hour * 24))
			if result.Day() > startDay {
				break
			}
		}
	}
	return &result, true
}

var timeToleranceUnits = map[string]int64{
	"":            int64(time.Hour * 24),
	"day":         int64(time.Hour * 24),
	"week":        int64(time.Hour * 24 * 7),
	"hour":        int64(time.Hour),
	"min":         int64(time.Minute),
	"minute":      int64(time.Minute),
	"sec":         int64(time.Second),
	"second":      int64(time.Second),
	"milli":       int64(time.Millisecond),
	"millisecond": int64(time.Millisecond),
	"micro":       int64(time.Microsecond),
	"microsecond": int64(time.Microsecond),
	"nano":        int64(time.Nanosecond),
	"nanosecond":  int64(time.Nanosecond),
}

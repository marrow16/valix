package valix

// Maximum constraint to check that a numeric value is less than or equal to a specified maximum
type Maximum struct {
	// the maximum value
	Value float64 `v8n:"default"`
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
func (c *Maximum) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f > c.Value || (c.ExclusiveMax && f == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Maximum) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgLt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgLte, c.Value)
}

// MaximumInt constraint to check that an integer value is less than or equal to a specified maximum
type MaximumInt struct {
	// the maximum value
	Value int64 `v8n:"default"`
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
func (c *MaximumInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i > c.Value || (c.ExclusiveMax && i == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MaximumInt) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMax {
		return defaultMessage(tcx, c.Message, fmtMsgLt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgLte, c.Value)
}

// Minimum constraint to check that a numeric value is greater than or equal to a specified minimum
type Minimum struct {
	// the minimum value
	Value float64 `v8n:"default"`
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Minimum) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f < c.Value || (c.ExclusiveMin && f == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Minimum) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgGt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgGte, c.Value)
}

// MinimumInt constraint to check that an integer numeric value is greater than or equal to a specified minimum
type MinimumInt struct {
	// the minimum value
	Value int64 `v8n:"default"`
	// if set to true, ExclusiveMin specifies the minimum value is exclusive
	ExclusiveMin bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *MinimumInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i < c.Value || (c.ExclusiveMin && i == c.Value) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MinimumInt) GetMessage(tcx I18nContext) string {
	if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgGt, c.Value)
	}
	return defaultMessage(tcx, c.Message, fmtMsgGte, c.Value)
}

// MultipleOf constraint to check that an integer value is a multiple of a specific number
//
// Note: this constraint will check values that are float or json Number - but the
// check will fail if either of these is not a 'whole number'
type MultipleOf struct {
	// the multiple of value to check
	Value int64 `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *MultipleOf) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i%c.Value != 0 {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *MultipleOf) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgMultipleOf, c.Value)
}

// Negative constraint to check that a numeric value is negative
type Negative struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Negative) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f >= 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Negative) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNegative)
}

// NegativeOrZero constraint to check that a numeric value is negative or zero
type NegativeOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *NegativeOrZero) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f > 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *NegativeOrZero) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgNegativeOrZero)
}

// Positive constraint to check that a numeric value is positive (exc. zero)
type Positive struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *Positive) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f <= 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Positive) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgPositive)
}

// PositiveOrZero constraint to check that a numeric value is positive or zero
type PositiveOrZero struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string `v8n:"default"`
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *PositiveOrZero) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber && (!ok || f < 0) {
		vcx.CeaseFurtherIf(c.Stop)
		return false, c.GetMessage(vcx)
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *PositiveOrZero) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgPositiveOrZero)
}

// Range constraint to check that a numeric value is within a specified minimum and maximum range
type Range struct {
	// the minimum value of the range
	Minimum float64
	// the maximum value of the range
	Maximum float64
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
func (c *Range) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok, isNumber := coerceToFloat(v); isNumber {
		if !ok || f < c.Minimum || (c.ExclusiveMin && f == c.Minimum) ||
			f > c.Maximum || (c.ExclusiveMax && f == c.Maximum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Range) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgRange, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
}

// RangeInt constraint to check that an integer value is within a specified minimum and maximum range
type RangeInt struct {
	// the minimum value of the range
	Minimum int64
	// the maximum value of the range
	Maximum int64
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
func (c *RangeInt) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if i, ok, isNumber := coerceToInt(v); isNumber {
		if !ok || i < c.Minimum || (c.ExclusiveMin && i == c.Minimum) ||
			i > c.Maximum || (c.ExclusiveMax && i == c.Maximum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *RangeInt) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgRange, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
}

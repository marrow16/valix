package valix

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	messageNotEmptyString         = "Value must not be empty string"
	messageNotBlankString         = "Value must not be a blank string"
	messageNoControlChars         = "Value must not contain control characters"
	messageInvalidPattern         = "Value has invalid pattern"
	messageAtLeast                = "Value length must be at least %d"
	messageNotMore                = "Value length must not exceed %d"
	messageMinMax                 = "Value length must be between %d and %d (inclusive)"
	messagePositive               = "Value must be positive"
	messagePositiveOrZero         = messagePositive + " or zero"
	messageNegative               = "Value must be negative"
	messageNegativeOrZero         = messageNegative + " or zero"
	messageGte                    = "Value must be greater than or equal to %f"
	messageLte                    = "Value must be less than or equal to %f"
	messageRange                  = "Value must be between %f and %f (inclusive)"
	messageArrayElementType       = "Array value elements must be of type %s"
	messageArrayElementTypeOrNull = "Array value elements must be of type %s or null"
	messageValidUuid              = "Value must be a valid UUID"
	messageUuidMinVersion         = "Value must be a valid UUID (minimum version %d)"
	messageUuidCorrectVer         = "Value must be a valid UUID (version %d)"
	uuidRegexpPattern             = "^([a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12})$"
)

var (
	uuidRegexp = *regexp.MustCompile(uuidRegexpPattern)
)

// StringNotEmptyConstraint to check that string value is not empty (i.e. not "")
type StringNotEmptyConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNotEmptyConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringNotEmptyConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageNotEmptyString)
}

// StringNotBlankConstraint to check that string value is not blank (i.e. that after removing leading and
// trailing whitespace the value is not an empty string)
type StringNotBlankConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNotBlankConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(strings.Trim(str, " \t\n\r")) == 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringNotBlankConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageNotBlankString)
}

// StringNoControlCharsConstraint to check that a string does not contain any control characters (i.e. chars < 32)
type StringNoControlCharsConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringNoControlCharsConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		for _, ch := range str {
			if ch < 32 {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *StringNoControlCharsConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageNoControlChars)
}

// StringPatternConstraint to check that a string matches a given regexp pattern
type StringPatternConstraint struct {
	// the regexp pattern that the string value must match
	Regexp regexp.Regexp
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringPatternConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !c.Regexp.MatchString(str) {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringPatternConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageInvalidPattern)
}

// StringMinLengthConstraint to check that a string has a minimum length
type StringMinLengthConstraint struct {
	// the minimum length value
	Value int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringMinLengthConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) < c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringMinLengthConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageAtLeast, c.Value))
}

// StringMaxLengthConstraint to check that a string has a maximum length
type StringMaxLengthConstraint struct {
	// the maximum length value
	Value int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringMaxLengthConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) > c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringMaxLengthConstraint) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageNotMore, c.Value))
}

// StringLengthConstraint to check that a string has a minimum and maximum length
type StringLengthConstraint struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringLengthConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if len(str) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringLengthConstraint) GetMessage() string {
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageMinMax, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageAtLeast, c.Minimum))
}

// LengthConstraint to check that a property value has minimum and maximum length
//
// This constraint can be used for string, object and array property values
type LengthConstraint struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *LengthConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, okS := value.(string); okS {
		if len(str) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(str) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if m, okM := value.(map[string]interface{}); okM {
		if len(m) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(m) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if a, okA := value.([]interface{}); okA {
		if len(a) < c.Minimum {
			return false, c.GetMessage()
		} else if c.Maximum > 0 && len(a) > c.Maximum {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *LengthConstraint) GetMessage() string {
	if c.Maximum > 0 {
		return defaultMessage(c.Message,
			fmt.Sprintf(messageMinMax, c.Minimum, c.Maximum))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageAtLeast, c.Minimum))
}

// PositiveConstraint to check that a numeric value is positive (exc. zero)
type PositiveConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *PositiveConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f <= 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i <= 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv <= 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *PositiveConstraint) GetMessage() string {
	return defaultMessage(c.Message, messagePositive)
}

// PositiveOrZeroConstraint to check that a numeric value is positive or zero
type PositiveOrZeroConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *PositiveOrZeroConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i < 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv < 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *PositiveOrZeroConstraint) GetMessage() string {
	return defaultMessage(c.Message, messagePositiveOrZero)
}

// NegativeConstraint to check that a numeric value is negative
type NegativeConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *NegativeConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f >= 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i >= 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv >= 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *NegativeConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageNegative)
}

// NegativeOrZeroConstraint to check that a numeric value is negative or zero
type NegativeOrZeroConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *NegativeOrZeroConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f > 0 {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if i > 0 {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv > 0 {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *NegativeOrZeroConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageNegativeOrZero)
}

// MinimumConstraint to check that a numeric value is greater than or equal to a specified minimum
type MinimumConstraint struct {
	// the minimum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *MinimumConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < c.Value {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Value {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv < c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *MinimumConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageGte, c.Value))
}

// MaximumConstraint to check that a numeric value is less than or equal to a specified maximum
type MaximumConstraint struct {
	// the maximum value
	Value float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *MaximumConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f > c.Value {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) > c.Value {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil && fv > c.Value {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *MaximumConstraint) GetMessage() string {
	return defaultMessage(c.Message, fmt.Sprintf(messageLte, c.Value))
}

// RangeConstraint to check that a numeric value is within a specified minimum and maximum range
type RangeConstraint struct {
	// the minimum value of the range (inclusive)
	Minimum float64
	// the maximum value of the range (inclusive)
	Maximum float64
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *RangeConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if f, ok := value.(float64); ok {
		if f < c.Minimum {
			return false, c.GetMessage()
		} else if f > c.Maximum {
			return false, c.GetMessage()
		}
	} else if i, ok2 := value.(int); ok2 {
		if float64(i) < c.Minimum {
			return false, c.GetMessage()
		} else if float64(i) > c.Maximum {
			return false, c.GetMessage()
		}
	} else if n, ok3 := value.(json.Number); ok3 {
		if fv, fe := n.Float64(); fe == nil {
			if fv < c.Minimum {
				return false, c.GetMessage()
			} else if fv > c.Maximum {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *RangeConstraint) GetMessage() string {
	return defaultMessage(c.Message,
		fmt.Sprintf(messageRange, c.Minimum, c.Maximum))
}

// ArrayOfConstraint to check each element in an array value is of the correct type
type ArrayOfConstraint struct {
	// the type to check for each item (use PropertyType values)
	Type string
	// whether to allow null items in the array
	AllowNullElement bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *ArrayOfConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if a, ok := value.([]interface{}); ok {
		for _, elem := range a {
			if elem == nil {
				if !c.AllowNullElement {
					return false, c.GetMessage()
				}
			} else if !checkValueType(elem, c.Type) {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *ArrayOfConstraint) GetMessage() string {
	if c.AllowNullElement {
		return defaultMessage(c.Message, fmt.Sprintf(messageArrayElementTypeOrNull, c.Type))
	}
	return defaultMessage(c.Message,
		fmt.Sprintf(messageArrayElementType, c.Type))
}

// StringValidUuidConstraint to check that a string value is a valid UUID
type StringValidUuidConstraint struct {
	// the minimum UUID version (optional - if zero this is not checked)
	MinVersion uint8
	// the specific UUID version (optional - if zero this is not checked)
	SpecificVersion uint8
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringValidUuidConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !uuidRegexp.MatchString(str) {
			return false, c.GetMessage()
		}
		var version = str[14] - 48
		if c.MinVersion > 0 && version < c.MinVersion {
			return false, c.GetMessage()
		}
		if c.SpecificVersion > 0 && version != c.SpecificVersion {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringValidUuidConstraint) GetMessage() string {
	if c.SpecificVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidCorrectVer, c.SpecificVersion))
	} else if c.MinVersion > 0 {
		return defaultMessage(c.Message, fmt.Sprintf(messageUuidMinVersion, c.MinVersion))
	}
	return defaultMessage(c.Message, messageValidUuid)
}

const (
	iso8601FullPattern             = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601NoOffsPattern           = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(\\.\\d+)?)$"
	iso8601NoMillisPattern         = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}(([+-]\\d{2}:\\d{2})|Z)?)$"
	iso8601MinPattern              = "^(\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2})$"
	iso8601DateOnlyPattern         = "^(\\d{4}-\\d{2}-\\d{2})$"
	iso8601FullLayout              = "2006-01-02T15:04:05.999999999Z07:00"
	iso8601NoOffLayout             = "2006-01-02T15:04:05.999999999"
	iso8601NoMillisLayout          = "2006-01-02T15:04:05Z07:00"
	iso8601MinLayout               = "2006-01-02T15:04:05"
	iso8601DateOnlyLayout          = "2006-01-02"
	messageValidISODatetime        = "Value must be a valid date/time string"
	messageValidISODate            = "Value must be a valid date string (format: YYYY-MM-DD)"
	messageDatetimeFormatFull      = " (format: YYYY-MM-DDThh:mm:ss.sss[Z|+-hh:mm])"
	messageDatetimeFormatNoOffs    = " (format: YYYY-MM-DDThh:mm:ss.sss)"
	messageDatetimeFormatNoMillis  = " (format: YYYY-MM-DDThh:mm:ss[Z|+-hh:mm])"
	messageDatetimeFormatMin       = " (format: YYYY-MM-DDThh:mm:ss)"
	messageDatetimeFuture          = "Value must be a valid date/time in the future"
	messageDatetimeFutureOrPresent = "Value must be a valid date/time in the future or present"
	messageDatetimePast            = "Value must be a valid date/time in the past"
	messageDatetimePastOrPresent   = "Value must be a valid date/time in the past or present"
)

var (
	iso8601FullRegex     = regexp.MustCompile(iso8601FullPattern)
	iso8601NoOffsRegex   = regexp.MustCompile(iso8601NoOffsPattern)
	iso8601NoMillisRegex = regexp.MustCompile(iso8601NoMillisPattern)
	iso8601MinRegex      = regexp.MustCompile(iso8601MinPattern)
	iso8601DateOnlyRegex = regexp.MustCompile(iso8601DateOnlyPattern)
)

// StringValidISODatetimeConstraint checks that a string value is a valid ISO8601 Date/time format
type StringValidISODatetimeConstraint struct {
	// NoOffset specifies, if set to true, that time offsets are not permitted
	NoOffset bool
	// NoMillis specifies, if set to true, that seconds cannot have decimal places
	NoMillis bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringValidISODatetimeConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
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
			return false, c.GetMessage()
		}
		// and attempt to parse it (it may match regex but datetime could still be invalid)...
		if _, err := time.Parse(useLayout, str); err != nil {
			if perr, ok := err.(*time.ParseError); ok && perr.Message == "" && strings.HasSuffix(useLayout, perr.LayoutElem) {
				// time.Parse is pretty dumb when it comes to timezones - if it's in the layout but not the string it fails
				// so remove the bit it doesn't like (in the layout) and try again...
				useLayout = useLayout[0 : len(useLayout)-len(perr.LayoutElem)]
				if _, err := time.Parse(useLayout, str); err != nil {
					return false, c.GetMessage()
				}
			} else {
				return false, c.GetMessage()
			}
		}
	}
	return true, ""
}
func (c *StringValidISODatetimeConstraint) GetMessage() string {
	if c.NoOffset && c.NoMillis {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatMin)
	} else if c.NoOffset {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatNoOffs)
	} else if c.NoMillis {
		return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatNoMillis)
	}
	return defaultMessage(c.Message, messageValidISODatetime+messageDatetimeFormatFull)
}

// StringValidISODateConstraint checks that a string value is a valid ISO8601 Date format (excluding time)
type StringValidISODateConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *StringValidISODateConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if !iso8601DateOnlyRegex.MatchString(str) {
			return false, c.GetMessage()
		}
		// and attempt to parse it (it may match regex but date could still be invalid)...
		if _, err := time.Parse(iso8601DateOnlyLayout, str); err != nil {
			return false, c.GetMessage()
		}
	}
	return true, ""
}
func (c *StringValidISODateConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageValidISODate)
}

// DatetimeFutureConstraint checks that a datetime/data (represented as string or time.Time) is in the future
type DatetimeFutureConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *DatetimeFutureConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || !dt.After(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && !dt.After(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}
func (c *DatetimeFutureConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimeFuture)
}

// DatetimeFutureOrPresentConstraint checks that a datetime/data (represented as string or time.Time) is in the future or present
type DatetimeFutureOrPresentConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *DatetimeFutureOrPresentConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || dt.Before(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && dt.Before(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}
func (c *DatetimeFutureOrPresentConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimeFutureOrPresent)
}

// DatetimePastConstraint checks that a datetime/data (represented as string or time.Time) is in the past
type DatetimePastConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *DatetimePastConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || !dt.Before(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && !dt.Before(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}
func (c *DatetimePastConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimePast)
}

// DatetimePastOrPresentConstraint checks that a datetime/data (represented as string or time.Time) is in the past or present
type DatetimePastOrPresentConstraint struct {
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
}

func (c *DatetimePastOrPresentConstraint) Validate(value interface{}, vcx *ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		if dt, ok2 := stringToTime(str); !ok2 || dt.After(time.Now()) {
			return false, c.GetMessage()
		}
	} else if dt, ok2 := value.(time.Time); ok2 && dt.After(time.Now()) {
		return false, c.GetMessage()
	}
	return true, ""
}
func (c *DatetimePastOrPresentConstraint) GetMessage() string {
	return defaultMessage(c.Message, messageDatetimePastOrPresent)
}

func stringToTime(str string) (*time.Time, bool) {
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

func defaultMessage(msg string, def string) string {
	if len(msg) == 0 {
		return def
	}
	return msg
}

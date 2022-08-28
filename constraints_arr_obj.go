package valix

// ArrayOf constraint to check each element in an array value is of the correct type
type ArrayOf struct {
	// the type to check for each item (use Type values)
	Type string `v8n:"default"`
	// whether to allow null items in the array
	AllowNullElement bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *ArrayOf) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if a, ok := v.([]interface{}); ok {
		if chkType, tOk := JsonTypeFromString(c.Type); tOk {
			for _, elem := range a {
				if elem == nil {
					if !c.AllowNullElement {
						vcx.CeaseFurtherIf(c.Stop)
						return false, c.GetMessage(vcx)
					}
				} else if !checkValueType(elem, chkType) {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *ArrayOf) GetMessage(tcx I18nContext) string {
	if c.AllowNullElement {
		return defaultMessage(tcx, c.Message, fmtMsgArrayElementTypeOrNull, c.Type)
	}
	return defaultMessage(tcx, c.Message, fmtMsgArrayElementType, c.Type)
}

// ArrayUnique constraint to check each element in an array value is unique
type ArrayUnique struct {
	// whether to ignore null items in the array
	IgnoreNulls bool `v8n:"default"`
	// whether uniqueness is case in-insensitive (for string elements)
	IgnoreCase bool
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *ArrayUnique) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if a, ok := v.([]interface{}); ok {
		list := make([]interface{}, 0, len(a))
		for _, iv := range a {
			if !(iv == nil && c.IgnoreNulls) {
				if !isUniqueCompare(iv, c.IgnoreCase, &list) {
					vcx.CeaseFurtherIf(c.Stop)
					return false, c.GetMessage(vcx)
				}
			}
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *ArrayUnique) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, msgArrayUnique)
}

// Length constraint to check that a property value has minimum and maximum length
//
// This constraint can be used for object, array and string property values - however, if
// checking string lengths it is better to use the StringLength constraint
//
// * when checking array values, the number of elements in the array is checked
//
// * when checking object values, the number of properties in the object is checked
//
// * when checking string values, the length of the string is checked
type Length struct {
	// the minimum length
	Minimum int
	// the maximum length (only checked if this value is > 0)
	Maximum int
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
func (c *Length) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, okS := v.(string); okS {
		l := len(str)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if m, okM := v.(map[string]interface{}); okM {
		l := len(m)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if a, okA := v.([]interface{}); okA {
		l := len(a)
		if l < c.Minimum || (c.ExclusiveMin && l == c.Minimum) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		} else if c.Maximum > 0 && (l > c.Maximum || (c.ExclusiveMax && l == c.Maximum)) {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *Length) GetMessage(tcx I18nContext) string {
	if c.Maximum > 0 {
		return defaultMessage(tcx, c.Message, fmtMsgMinMax, c.Minimum, incExc(tcx, c.ExclusiveMin), c.Maximum, incExc(tcx, c.ExclusiveMax))
	} else if c.ExclusiveMin {
		return defaultMessage(tcx, c.Message, fmtMsgMinLenExc, c.Minimum)
	}
	return defaultMessage(tcx, c.Message, fmtMsgMinLen, c.Minimum)
}

// LengthExact constraint to check that a property value has a specific length
//
// This constraint can be used for object, array and string property values - however, if
// checking string lengths it is better to use the StringExactLength constraint
//
// * when checking array values, the number of elements in the array is checked
//
// * when checking object values, the number of properties in the object is checked
//
// * when checking string values, the length of the string is checked
type LengthExact struct {
	// the length to check
	Value int `v8n:"default"`
	// the violation message to be used if the constraint fails (see Violation.Message)
	//
	// (if the Message is an empty string then the default violation message is used)
	Message string
	// when set to true, Stop prevents further validation checks on the property if this constraint fails
	Stop bool
}

// Check implements Constraint.Check
func (c *LengthExact) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	if str, okS := v.(string); okS {
		l := len(str)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if m, okM := v.(map[string]interface{}); okM {
		l := len(m)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	} else if a, okA := v.([]interface{}); okA {
		l := len(a)
		if l != c.Value {
			vcx.CeaseFurtherIf(c.Stop)
			return false, c.GetMessage(vcx)
		}
	}
	return true, ""
}

// GetMessage implements the Constraint.GetMessage
func (c *LengthExact) GetMessage(tcx I18nContext) string {
	return defaultMessage(tcx, c.Message, fmtMsgExactLen, c.Value)
}
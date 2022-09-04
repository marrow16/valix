package valix

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

const (
	CodeRequestQueryParamMultiNotAllowed = 40010
	msgQueryParamMultiNotAllowed         = "Query param may not be specified more than once"
	CodeRequestQueryParamInvalidType     = 40011
	fmtMsgQueryParamType                 = "Query param must be of type %[1]s"
)

// RequestQueryValidate Performs validation on the request query (http.Request.URL.Query) of the supplied http.Request
//
// If the validation of the request query fails, false is returned and the returned violations
// give the reason(s) for the validation failure.
//
// If the validation is successful, the validated query (as JSON object) is also returned
func (v *Validator) RequestQueryValidate(req *http.Request, initialConditions ...string) (bool, []*Violation, interface{}) {
	i18ctx := obtainI18nProvider().ContextFromRequest(req)
	if obj, violations := v.queryParamsToObject(req, i18ctx); len(violations) == 0 {
		vcx := newValidatorContext(obj, v, v.StopOnFirst, i18ctx)
		vcx.setConditionsFromRequest(req)
		vcx.setInitialConditions(initialConditions...)
		v.validateObjectOrArray(vcx, obj, true)
		return vcx.ok, vcx.violations, obj
	} else {
		return false, violations, nil
	}
}

// RequestQueryValidateInto performs validation on the request query (http.Request.URL.Query) of the supplied http.Request
// and, if validation successful, attempts to unmarshall the query params into the supplied value
func (v *Validator) RequestQueryValidateInto(req *http.Request, value interface{}, initialConditions ...string) (bool, []*Violation, interface{}) {
	i18ctx := obtainI18nProvider().ContextFromRequest(req)
	if obj, violations := v.queryParamsToObject(req, i18ctx); len(violations) == 0 {
		vcx := newValidatorContext(obj, v, v.StopOnFirst, i18ctx)
		vcx.setConditionsFromRequest(req)
		vcx.setInitialConditions(initialConditions...)
		v.validateObjectOrArray(vcx, obj, true)
		if !vcx.ok {
			return vcx.ok, vcx.violations, obj
		}
		// now read into the provided value...
		buffer, _ := json.Marshal(obj)
		intoReader := bytes.NewReader(buffer)
		decoder := getDefaultDecoderProvider().NewDecoderFor(intoReader, v)
		err := decoder.Decode(value)
		if err != nil {
			vcx.AddViolation(newBadRequestViolation(vcx, msgErrorUnmarshall, CodeErrorUnmarshall, err))
		}
		return vcx.ok, vcx.violations, obj
	} else {
		return false, violations, nil
	}
}

func (v *Validator) queryParamsToObject(req *http.Request, i18ctx I18nContext) (map[string]interface{}, []*Violation) {
	result := map[string]interface{}{}
	tmpVcx := newEmptyValidatorContext(i18ctx)
	values := req.URL.Query()
	for k, vs := range values {
		if pty, ok := v.Properties[k]; ok {
			if useV, violations := convertQueryParamValues(vs, k, pty, i18ctx); len(violations) == 0 {
				result[k] = useV
			} else {
				tmpVcx.violations = append(tmpVcx.violations, violations...)
			}
		} else if len(vs) == 1 {
			if vs[0] == "" {
				result[k] = true
			} else {
				result[k] = vs[0]
			}
		} else {
			allEmpty := vs[0] == ""
			rvs := make([]interface{}, len(vs))
			for i, vsv := range vs {
				rvs[i] = vsv
				allEmpty = allEmpty && vsv == ""
			}
			if allEmpty {
				rvb := make([]interface{}, len(rvs))
				for i := range vs {
					rvb[i] = true
				}
				result[k] = rvb
			} else {
				result[k] = rvs
			}
		}
	}
	return result, tmpVcx.violations
}

func convertQueryParamValues(values []string, name string, pty *PropertyValidator, i18ctx I18nContext) (interface{}, []*Violation) {
	var result interface{} = nil
	violations := make([]*Violation, 0)
	switch pty.Type {
	case JsonArray:
		elementType := JsonAny
		// need to check what element type it's expecting - only way to do this (currently) is to sniff at constraints
		// and see if ArrayOf has been used...
		for _, c := range pty.Constraints {
			switch aoc := c.(type) {
			case *ArrayOf:
				if aot, ok := JsonTypeFromString(aoc.Type); ok {
					elementType = aot
					break
				}
			}
		}
		arrResult := make([]interface{}, len(values))
		for i, ev := range values {
			r, violation := convertQueryParamValue(ev, name, elementType, i18ctx)
			arrResult[i] = r
			if violation != nil {
				violations = append(violations, violation)
			}
		}
		result = arrResult
	case JsonAny:
		// don't do any conversions...
		arrResult := make([]interface{}, len(values))
		for i, ev := range values {
			arrResult[i] = ev
		}
		result = arrResult
	default:
		if len(values) > 1 {
			violation := NewViolation(name, "", defaultMessage(i18ctx, "", msgQueryParamMultiNotAllowed), CodeRequestQueryParamMultiNotAllowed)
			violation.BadRequest = true
			violations = append(violations, violation)
		} else {
			r, violation := convertQueryParamValue(values[0], name, pty.Type, i18ctx)
			result = r
			if violation != nil {
				violations = append(violations, violation)
			}
		}
	}
	return result, violations
}

func convertQueryParamValue(value string, name string, t JsonType, i18ctx I18nContext) (interface{}, *Violation) {
	var result interface{} = nil
	var violation *Violation = nil
	switch t {
	case JsonString, JsonAny:
		result = value
	case JsonNumber, JsonInteger:
		result = json.Number(value)
	case JsonDatetime:
		if dt, ok := stringToDatetime(value, false); ok {
			result = dt
		} else {
			// don't raise a violation here (other type validation checks should spot it)
			// just add the raw string value...
			result = value
		}
	case JsonBoolean:
		if value == "" {
			result = true
		} else if b, err := strconv.ParseBool(value); err == nil {
			result = b
		} else {
			violation = NewViolation(name, "", defaultMessage(i18ctx, "", fmtMsgQueryParamType, "boolean"), CodeRequestQueryParamInvalidType)
			violation.BadRequest = true
		}
	case JsonObject:
		if value != "" {
			// this is an odd one - we'll have to assume the intention was that the query param value is an encoded json object...
			obj := map[string]interface{}{}
			if err := json.Unmarshal([]byte(value), &obj); err == nil {
				result = obj
			} else {
				violation = NewViolation(name, "", defaultMessage(i18ctx, "", fmtMsgQueryParamType, "object"), CodeRequestQueryParamInvalidType)
				violation.BadRequest = true
			}
		}
	case JsonArray:
		if value != "" {
			// and this is also another odd one - at first level this won't be reached, so it's only multi-valued where element types are array...
			// so assume that each item is an encoded json array...
			arr := make([]interface{}, 0)
			if err := json.Unmarshal([]byte(value), &arr); err == nil {
				result = arr
			} else {
				violation = NewViolation(name, "", defaultMessage(i18ctx, "", fmtMsgQueryParamType, "array"), CodeRequestQueryParamInvalidType)
				violation.BadRequest = true
			}
		}
	}
	return result, violation
}

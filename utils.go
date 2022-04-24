package valix

import (
	"encoding/json"
	"math"
	"reflect"
)

func coerceToFloat(value interface{}) (f float64, ok bool, isNumber bool) {
	if value == nil {
		return 0, false, false
	}
	ok = true
	isNumber = true
	switch v := value.(type) {
	case float64:
		f = v
	case float32:
		f = float64(v)
	case int, int8, int16, int32, int64:
		f = float64(reflect.ValueOf(value).Int())
	case uint, uint8, uint16, uint32, uint64:
		f = float64(reflect.ValueOf(value).Uint())
	case json.Number:
		if jf, err := v.Float64(); err == nil {
			f = jf
		} else {
			ok = false
		}
	default:
		ok = false
		isNumber = false
	}
	ok = ok && !math.IsNaN(f)
	return
}

func coerceToInt(value interface{}) (i int64, ok bool, isNumber bool) {
	if value == nil {
		return 0, false, false
	}
	ok = true
	isNumber = true
	switch v := value.(type) {
	case int64:
		i = v
	case int, int8, int16, int32:
		i = reflect.ValueOf(v).Int()
	case uint, uint8, uint16, uint32, uint64:
		i = int64(reflect.ValueOf(v).Uint())
	case float32, float64:
		f := reflect.ValueOf(value).Float()
		if !math.IsNaN(f) && !math.IsInf(f, 0) {
			i = int64(f)
			ok = math.Trunc(f) == f
		} else {
			ok = false
		}
	case json.Number:
		if ji, err := v.Int64(); err == nil {
			i = ji
		} else if jf, err := v.Float64(); err == nil {
			if !math.IsNaN(jf) && !math.IsInf(jf, 0) {
				i = int64(jf)
				ok = math.Trunc(jf) == jf
			} else {
				ok = false
			}
		} else {
			ok = false
		}
	default:
		ok = false
		isNumber = false
	}
	return
}

func defaultString(str string, def string) string {
	if str != "" {
		return str
	}
	return def
}

// ternary operations

type ternary bool

func (b ternary) string(t, f string) string {
	if b {
		return t
	}
	return f
}

func (b ternary) int(t, f int) int {
	if b {
		return t
	}
	return f
}

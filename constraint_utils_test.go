package valix

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
	"unicode"
)

func TestDefaultMessage(t *testing.T) {
	vcx := newEmptyValidatorContext(nil)

	msg := defaultMessage(vcx, "Foo", fmtMsgGtOther, "other")
	require.Equal(t, "Foo", msg)

	msg = defaultMessage(vcx, "", fmtMsgGtOther, "other")
	require.Equal(t, fmt.Sprintf(fmtMsgGtOther, "other"), msg)

	msg = defaultMessage(vcx, NoMessage, fmtMsgGtOther, "other")
	require.Equal(t, "", msg)

	msg = defaultMessage(vcx, "", msgFailure)
	require.Equal(t, msgFailure, msg)

	i18ctx := newDefaultI18nContext("fr", "")
	vcx = newEmptyValidatorContext(i18ctx)
	msg = defaultMessage(vcx, "", msgFailure)
	require.NotEqual(t, msgFailure, msg)
}
func TestIncExc(t *testing.T) {
	i18ctx := newDefaultI18nContext("en", "")
	msg := incExc(i18ctx, true)
	require.Equal(t, tokenExclusive, msg)
	msg = incExc(i18ctx, false)
	require.Equal(t, tokenInclusive, msg)

	i18ctx = newDefaultI18nContext("fr", "")
	msg = incExc(i18ctx, false)
	require.NotEqual(t, tokenInclusive, msg)
}

func TestStringToDatetime(t *testing.T) {
	testCases := map[string]bool{
		"":                                     false,
		"not a date":                           false,
		"2022-05-04T10:11:12":                  true,
		"2022-05-04T30:11:12":                  false,
		"2022-05-04T10:60:12":                  false,
		"2022-05-04T10:11:60":                  false,
		"2022-05-04T10:11:12Z":                 true,
		"2022-05-04T10:11:12+12:00":            true,
		"2022-05-04T10:11:12.13456+12:00":      true,
		"2022-05-04T10:11:12.134567890+12:00":  true,
		"2022-05-04T10:11:12.1345678901+12:00": false,
		"2022-05-04T10:11:12.1":                true,
		"2022-05-04":                           true,
		"2022-05-04T":                          false,
		"2022-05-04T1":                         false,
		"2022-05-04T10":                        false,
		"2022-13-04":                           false,
		"2022-05-33":                           false,
	}
	for str, expectOk := range testCases {
		t.Run(fmt.Sprintf("stringToDatetime(%s)", str), func(t *testing.T) {
			_, ok := stringToDatetime(str, false)
			require.Equal(t, expectOk, ok)
		})
	}
}

func TestStringToDatetimeWithTruncation(t *testing.T) {
	str := ""
	_, ok := stringToDatetime(str, true)
	require.False(t, ok)
	str = "2022005004"
	_, ok = stringToDatetime(str, true)
	require.False(t, ok)

	str = "2022-05-04"
	dt, ok := stringToDatetime(str, true)
	require.True(t, ok)
	require.Equal(t, 0, dt.Hour())
	require.Equal(t, 0, dt.Minute())
	require.Equal(t, 0, dt.Second())
	require.Equal(t, 0, dt.Nanosecond())
	zn, zo := dt.Zone()
	require.Equal(t, time.UTC.String(), zn)
	require.Equal(t, 0, zo)

	str = "2022-05-04T10:11:12.1234+07:00"
	dt, ok = stringToDatetime(str, true)
	require.True(t, ok)
	require.Equal(t, 0, dt.Hour())
	require.Equal(t, 0, dt.Minute())
	require.Equal(t, 0, dt.Second())
	require.Equal(t, 0, dt.Nanosecond())
	zn, zo = dt.Zone()
	require.Equal(t, time.UTC.String(), zn)
	require.Equal(t, 0, zo)
}

func TestTruncateTime(t *testing.T) {
	loc := time.FixedZone("Foo", -3600)
	dt := time.Date(2022, 5, 4, 10, 11, 12, 123, loc)
	tdt := truncateTime(dt, false)
	require.Equal(t, dt, tdt)

	tdt = truncateTime(dt, true)
	require.NotEqual(t, dt, tdt)
	require.Equal(t, 0, tdt.Hour())
	require.Equal(t, 0, tdt.Minute())
	require.Equal(t, 0, tdt.Second())
	require.Equal(t, 0, tdt.Nanosecond())
	zn, zo := tdt.Zone()
	require.Equal(t, time.UTC.String(), zn)
	require.Equal(t, 0, zo)
}

func TestGetOtherProperty(t *testing.T) {
	obj := map[string]interface{}{
		"foo": "foo value",
	}
	vcx := newValidatorContext(obj, nil, false, nil)

	_, ok := getOtherProperty("foo", vcx)
	require.False(t, ok)

	vcx.pushPathProperty("bar", nil, nil)
	ov, ok := getOtherProperty("foo", vcx)
	require.True(t, ok)
	require.NotNil(t, ov)
	str := ov.(string)
	require.Equal(t, "foo value", str)

	delete(obj, "foo")
	_, ok = getOtherProperty("foo", vcx)
	require.False(t, ok)
}

func TestGetOtherPropertyPathed(t *testing.T) {
	obj := map[string]interface{}{
		"foo": "foo value",
		"bar": map[string]interface{}{
			"foo": "bar.foo value",
			"baz": map[string]interface{}{
				"foo": "bar.baz.foo value",
				"qux": map[string]interface{}{
					"foo": "bar.baz.qux.foo value",
				},
			},
		},
	}
	vcx := newValidatorContext(obj, nil, false, nil)
	bar := obj["bar"].(map[string]interface{})
	vcx.pushPathProperty("bar", bar, nil)

	// test walking down...
	v, ok := getOtherProperty("foo", vcx)
	require.True(t, ok)
	require.Equal(t, "foo value", v)
	v, ok = getOtherProperty(".foo", vcx)
	require.True(t, ok)
	require.Equal(t, "foo value", v)

	v, ok = getOtherProperty(".bar.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.foo value", v)
	v, ok = getOtherProperty("bar.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.foo value", v)

	v, ok = getOtherProperty(".bar.baz.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.foo value", v)

	v, ok = getOtherProperty(".bar.baz.qux.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.qux.foo value", v)

	_, ok = getOtherProperty(".bar.baz.qux.foo.oops", vcx)
	require.False(t, ok)
	_, ok = getOtherProperty(".bar.baz....", vcx)
	require.False(t, ok)

	// setup for walking up...
	baz := bar["baz"].(map[string]interface{})
	vcx.pushPathProperty("baz", baz, nil)
	qux := baz["qux"].(map[string]interface{})
	vcx.pushPathProperty("qux", qux, nil)
	finalFoo := qux["foo"]
	vcx.pushPathProperty("foo", finalFoo, nil)

	// test walking up...
	v, ok = getOtherProperty("foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.qux.foo value", v)
	v, ok = getOtherProperty(".foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.qux.foo value", v)

	v, ok = getOtherProperty("..foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.foo value", v)

	v, ok = getOtherProperty("...foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.foo value", v)

	v, ok = getOtherProperty("....foo", vcx)
	require.True(t, ok)
	require.Equal(t, "foo value", v)

	_, ok = getOtherProperty(".....foo", vcx)
	require.False(t, ok)

	// test walking up and then down...
	v, ok = getOtherProperty("....bar.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.foo value", v)

	v, ok = getOtherProperty("....bar.baz.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.foo value", v)

	v, ok = getOtherProperty("....bar.baz.qux.foo", vcx)
	require.True(t, ok)
	require.Equal(t, "bar.baz.qux.foo value", v)
}

func TestGetOtherPropertyDatetime(t *testing.T) {
	obj := map[string]interface{}{
		"foo": "2022-05-04T10:11:12",
	}
	vcx := newValidatorContext(obj, nil, false, nil)

	_, ok := getOtherPropertyDatetime("foo", vcx, false, false)
	require.False(t, ok)

	vcx.pushPathProperty("bar", nil, nil)
	dt, ok := getOtherPropertyDatetime("foo", vcx, false, false)
	require.True(t, ok)
	require.NotNil(t, dt)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 0, dt.Nanosecond())

	obj["foo"] = time.Date(2022, 5, 4, 10, 11, 12, 123, time.UTC)
	dt, ok = getOtherPropertyDatetime("foo", vcx, false, false)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 123, dt.Nanosecond())

	obj["foo"] = Time{time.Date(2022, 5, 4, 10, 11, 12, 456, time.UTC), ""}
	dt, ok = getOtherPropertyDatetime("foo", vcx, false, false)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 456, dt.Nanosecond())

	dt, ok = getOtherPropertyDatetime("foo", vcx, true, false)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 0, dt.Hour())
	require.Equal(t, 0, dt.Minute())
	require.Equal(t, 0, dt.Second())
	require.Equal(t, 0, dt.Nanosecond())

	obj["foo"] = nil
	_, ok = getOtherPropertyDatetime("foo", vcx, false, false)
	require.False(t, ok)
	_, ok = getOtherPropertyDatetime("foo", vcx, false, true)
	require.True(t, ok)

}

func TestIsTime(t *testing.T) {
	_, ok := isTime(nil, false)
	require.False(t, ok)

	_, ok = isTime(time.Now(), false)
	require.True(t, ok)

	now := time.Now()
	_, ok = isTime(&now, false)
	require.True(t, ok)

	dt, ok := isTime(Time{time.Date(2022, 5, 4, 10, 11, 12, 456, time.UTC), ""}, false)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 10, dt.Hour())
	require.Equal(t, 456, dt.Nanosecond())
	dt, ok = isTime(Time{time.Date(2022, 5, 4, 10, 11, 12, 456, time.UTC), ""}, true)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 0, dt.Hour())
	require.Equal(t, 0, dt.Nanosecond())

	dt, ok = isTime(&Time{time.Date(2022, 5, 4, 10, 11, 12, 456, time.UTC), ""}, false)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 10, dt.Hour())
	require.Equal(t, 456, dt.Nanosecond())
	dt, ok = isTime(&Time{time.Date(2022, 5, 4, 10, 11, 12, 456, time.UTC), ""}, true)
	require.True(t, ok)
	require.Equal(t, 2022, dt.Year())
	require.Equal(t, 0, dt.Hour())
	require.Equal(t, 0, dt.Nanosecond())

	_, ok = isTime("", false)
	require.False(t, ok)

	dt, ok = isTime("2022-05-04", false)
	require.True(t, ok)
	require.Equal(t, 4, dt.Day())
	require.Equal(t, 0, dt.Hour())
	dt, ok = isTime("2022-05-04T10:11:12.345", true)
	require.True(t, ok)
	require.Equal(t, 4, dt.Day())
	require.Equal(t, 0, dt.Hour())

	str := "2022-05-04"
	dt, ok = isTime(&str, false)
	require.True(t, ok)
	require.Equal(t, 4, dt.Day())
	require.Equal(t, 0, dt.Hour())

	_, ok = isTime(0, false)
	require.False(t, ok)
	_, ok = isTime(false, false)
	require.False(t, ok)
}

func TestCompareNumerics(t *testing.T) {
	_, ok := compareNumerics(nil, nil)
	require.False(t, ok)
	_, ok = compareNumerics(1, nil)
	require.False(t, ok)
	_, ok = compareNumerics(nil, 1)
	require.False(t, ok)

	cmp, ok := compareNumerics(1.0, 2)
	require.True(t, ok)
	require.Equal(t, -1, cmp)

	cmp, ok = compareNumerics(2, 1.5)
	require.True(t, ok)
	require.Equal(t, 1, cmp)

	cmp, ok = compareNumerics(2, 2.0)
	require.True(t, ok)
	require.Equal(t, 0, cmp)
}

func TestTypedEquals(t *testing.T) {
	require.True(t, typedEquals(nil, nil))
	require.False(t, typedEquals("", nil))
	require.False(t, typedEquals(nil, ""))

	require.True(t, typedEquals("", ""))
	require.False(t, typedEquals("", "a"))
	require.True(t, typedEquals(1.0, 1))

	require.True(t, typedEquals(1, json.Number("1.000")))
	require.True(t, typedEquals(json.Number("1"), json.Number("1.000")))
	require.True(t, typedEquals(json.Number("1.6"), json.Number("1.600")))
	require.False(t, typedEquals(json.Number("1.6"), 1.0))

	require.True(t, typedEquals(1, 1))
	require.False(t, typedEquals(1, 2))
}

func TestCoerceJsonNumberToFloat(t *testing.T) {
	f, ok := coerceJsonNumberToFloat(json.Number("1"))
	require.True(t, ok)
	require.Equal(t, 1.0, f)

	_, ok = coerceJsonNumberToFloat(json.Number("x"))
	require.False(t, ok)
}

func TestIsUniqueCompare(t *testing.T) {
	distinctList := make([]interface{}, 0)
	unique := isUniqueCompare("foo", true, &distinctList)
	require.True(t, unique)
	unique = isUniqueCompare("Foo", true, &distinctList)
	require.False(t, unique)
	unique = isUniqueCompare("Foo", false, &distinctList)
	require.True(t, unique)

	distinctList = make([]interface{}, 0)
	unique = isUniqueCompare(1.0, false, &distinctList)
	require.True(t, unique)
	unique = isUniqueCompare(json.Number("1"), false, &distinctList)
	require.False(t, unique)
}

func TestGetOtherPropertyPath(t *testing.T) {
	level, path := getOtherPropertyPath(".foo")
	require.Equal(t, uint(0), level)
	require.Equal(t, 1, len(path))
	require.Equal(t, "foo", path[0])

	level, path = getOtherPropertyPath("foo")
	require.Equal(t, uint(0), level)
	require.Equal(t, 1, len(path))
	require.Equal(t, "foo", path[0])

	level, path = getOtherPropertyPath("...foo")
	require.Equal(t, uint(2), level)
	require.Equal(t, 1, len(path))
	require.Equal(t, "foo", path[0])

	level, path = getOtherPropertyPath("...foo[0]")
	require.Equal(t, uint(2), level)
	require.Equal(t, 2, len(path))
	require.Equal(t, "foo", path[0])
	require.Equal(t, "[0]", path[1])

	level, path = getOtherPropertyPath("...foo.[0]")
	require.Equal(t, uint(2), level)
	require.Equal(t, 2, len(path))
	require.Equal(t, "foo", path[0])
	require.Equal(t, "[0]", path[1])
}

func TestPropertyWalkDown(t *testing.T) {
	subm := map[string]interface{}{"foo": "bar"}
	sl := []interface{}{"a", "b", "c", subm}
	mp := map[string]interface{}{"foo": sl}

	v, ok := propertyWalkDown(mp, []string{"foo", "[1]"}, 0, -1)
	require.True(t, ok)
	require.Equal(t, "b", v)

	v, ok = propertyWalkDown(mp, []string{"foo", "[3]", "foo"}, 0, -1)
	require.True(t, ok)
	require.Equal(t, "bar", v)

	v, ok = propertyWalkDown(sl, []string{"[2]"}, 0, -1)
	require.True(t, ok)
	require.Equal(t, "c", v)

	v, ok = propertyWalkDown(sl, []string{"[+1]"}, 0, 1)
	require.True(t, ok)
	require.Equal(t, "c", v)

	v, ok = propertyWalkDown(sl, []string{"[+1]", "foo"}, 0, 2)
	require.True(t, ok)
	require.Equal(t, "bar", v)
}

func TestPropertyWalkDownSlice(t *testing.T) {
	sl := []interface{}{"a", "b", "c"}

	// test with invalid array index token...
	v, ok := propertyWalkDownSlice(sl, []string{"not an [] index"}, 0, -1)
	require.False(t, ok)
	require.Nil(t, v)

	// test getting value...
	v, ok = propertyWalkDownSlice(sl, []string{"[1]"}, 0, -1)
	require.True(t, ok)
	require.Equal(t, "b", v)
}

func TestParseDuration(t *testing.T) {
	testCases := map[string]struct {
		expectOk bool
		negative bool
		years    float64
		months   float64
		weeks    float64
		days     float64
		hours    float64
		minutes  float64
		seconds  float64
	}{
		"":        {},
		"xxx":     {},
		"P":       {},
		"p":       {},
		"PW":      {},
		"PYMTHMS": {},
		"P1M1M":   {},
		"PNaNY":   {},
		"PInfY":   {},
		"PT":      {},
		"PTM":     {},
		"P1Y1W":   {},
		"P2W": {
			expectOk: true,
			weeks:    2,
		},
		"P1.5W": {
			expectOk: true,
			weeks:    1.5,
		},
		"P1,5W": {
			expectOk: true,
			weeks:    1.5,
		},
		"-P3W": {
			expectOk: true,
			negative: true,
			weeks:    3,
		},
		"P1Y": {
			expectOk: true,
			years:    1,
		},
		"P-1Y": {
			expectOk: true,
			years:    -1,
		},
		"-P1Y": {
			expectOk: true,
			years:    1,
			negative: true,
		},
		"-P-1Y": {
			expectOk: true,
			negative: true,
			years:    -1,
		},
		"P1M": {
			expectOk: true,
			months:   1,
		},
		"P2.5M": {
			expectOk: true,
			months:   2.5,
		},
		"P2,5M": {
			expectOk: true,
			months:   2.5,
		},
		"P1D": {
			expectOk: true,
			days:     1,
		},
		"P2.5D": {
			expectOk: true,
			days:     2.5,
		},
		"P2,5D": {
			expectOk: true,
			days:     2.5,
		},
		"PT1H": {
			expectOk: true,
			hours:    1,
		},
		"PT2.5H": {
			expectOk: true,
			hours:    2.5,
		},
		"PT2,5H": {
			expectOk: true,
			hours:    2.5,
		},
		"PT1M": {
			expectOk: true,
			minutes:  1,
		},
		"PT2.5M": {
			expectOk: true,
			minutes:  2.5,
		},
		"PT2,5M": {
			expectOk: true,
			minutes:  2.5,
		},
		"PT1S": {
			expectOk: true,
			seconds:  1,
		},
		"PT2.5S": {
			expectOk: true,
			seconds:  2.5,
		},
		"PT2,5S": {
			expectOk: true,
			seconds:  2.5,
		},
		"P1Y2M3DT4H5M6S": {
			expectOk: true,
			years:    1,
			months:   2,
			days:     3,
			hours:    4,
			minutes:  5,
			seconds:  6,
		},
		"P-1Y-2M-3DT-4H-5M-6S": {
			expectOk: true,
			years:    -1,
			months:   -2,
			days:     -3,
			hours:    -4,
			minutes:  -5,
			seconds:  -6,
		},
		"-P1Y2M3DT4H5M6S": {
			expectOk: true,
			negative: true,
			years:    1,
			months:   2,
			days:     3,
			hours:    4,
			minutes:  5,
			seconds:  6,
		},
	}
	for ds, tc := range testCases {
		t.Run(ds, func(t *testing.T) {
			dur, ok := ParseDuration(ds)
			require.Equal(t, tc.expectOk, ok)
			if tc.expectOk {
				require.NotNil(t, dur)
				require.Equal(t, tc.negative, dur.Negative)
				if tc.weeks != 0 {
					require.Equal(t, tc.weeks, *dur.Weeks)
				} else {
					require.Nil(t, dur.Weeks)
				}
				if tc.years != 0 {
					require.Equal(t, tc.years, *dur.Years)
				} else {
					require.Nil(t, dur.Years)
				}
				if tc.months != 0 {
					require.Equal(t, tc.months, *dur.Months)
				} else {
					require.Nil(t, dur.Months)
				}
				if tc.days != 0 {
					require.Equal(t, tc.days, *dur.Days)
				} else {
					require.Nil(t, dur.Days)
				}
				if tc.hours != 0 {
					require.Equal(t, tc.hours, *dur.Hours)
				} else {
					require.Nil(t, dur.Hours)
				}
				if tc.minutes != 0 {
					require.Equal(t, tc.minutes, *dur.Minutes)
				} else {
					require.Nil(t, dur.Minutes)
				}
				if tc.seconds != 0 {
					require.Equal(t, tc.seconds, *dur.Seconds)
				} else {
					require.Nil(t, dur.Seconds)
				}
			}
		})
	}
}

func TestLookupRangeTableName(t *testing.T) {
	_, ok := lookupRangeTableName("BMP")
	require.True(t, ok)
	_, ok = lookupRangeTableName("nope")
	require.False(t, ok)
	for k := range unicode.Categories {
		_, ok = lookupRangeTableName("Category-" + k)
		require.True(t, ok)
	}
	for k := range unicode.Scripts {
		_, ok = lookupRangeTableName("Script-" + k)
		require.True(t, ok)
	}
	for k := range unicode.Properties {
		_, ok = lookupRangeTableName("Property-" + k)
		require.True(t, ok)
	}
	for k := range unicode.FoldCategory {
		_, ok = lookupRangeTableName("FoldCategory-" + k)
		require.True(t, ok)
	}
	for k := range unicode.FoldScript {
		_, ok = lookupRangeTableName("FoldScript-" + k)
		require.True(t, ok)
	}
}

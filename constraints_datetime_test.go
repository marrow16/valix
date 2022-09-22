package valix

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestStringValidISODatetime(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidISODatetime{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02T18:19:20.12345+01:00 but not a datetime with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidISODatetimeFormatFull, violations[0].Message)

	obj["foo"] = "2022-02-02T18:19:20.123+01:00"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-13-02T18:19:20.12345+01:00"
	//                  ^^ 13th month?
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidISODatetimeFormatFull, violations[0].Message)
}

func TestStringValidISODatetimeWithDifferentSettings(t *testing.T) {
	vFull := buildFooValidator(JsonString,
		&StringValidISODatetime{}, false)
	vNoNano := buildFooValidator(JsonString,
		&StringValidISODatetime{NoMillis: true}, false)
	vNoOffs := buildFooValidator(JsonString,
		&StringValidISODatetime{NoOffset: true}, false)
	vMin := buildFooValidator(JsonString,
		&StringValidISODatetime{NoOffset: true, NoMillis: true}, false)

	testCases := []struct {
		testValue  string
		okFull     bool
		okNoMillis bool
		okNoOffs   bool
		okMin      bool
	}{
		{testValue: "", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20", okFull: true, okNoMillis: true, okNoOffs: true, okMin: true},
		{testValue: "2022-02-02T18:19:20.12345+01:00", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.1234567890+01:00", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345-01:00", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345Z", okFull: true, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.12345Z+01:00", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20+01:00", okFull: true, okNoMillis: true, okNoOffs: false, okMin: false},
		{testValue: "2022-02-02T18:19:20.123456", okFull: true, okNoMillis: false, okNoOffs: true, okMin: false},
		{testValue: "2022-02-02T18:19:20Z", okFull: true, okNoMillis: true, okNoOffs: false, okMin: false},
		// bad dates/times...
		{testValue: "2022-13-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-41T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T25:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:60:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:60.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		// too many digits in various places...
		{testValue: "20222-01-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-012-01T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-012T18:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T189:19:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:190:20.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:201.12345Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
		{testValue: "2022-01-01T18:19:20.1234567891Z", okFull: false, okNoMillis: false, okNoOffs: false, okMin: false},
	}
	obj := jsonObject(`{
		"foo": ""
	}`)
	for _, testCase := range testCases {
		t.Run(fmt.Sprintf("Value: %s", testCase.testValue), func(t *testing.T) {
			obj["foo"] = testCase.testValue
			ok, violations := vFull.Validate(obj)
			require.Equal(t, testCase.okFull, ok)
			if !testCase.okFull {
				require.Equal(t, 1, len(violations))
				require.Equal(t, msgValidISODatetimeFormatFull, violations[0].Message)
			}
			ok, violations = vNoNano.Validate(obj)
			require.Equal(t, testCase.okNoMillis, ok)
			if !testCase.okNoMillis {
				require.Equal(t, 1, len(violations))
				require.Equal(t, msgValidISODatetimeFormatNoMillis, violations[0].Message)
			}
			ok, violations = vNoOffs.Validate(obj)
			require.Equal(t, testCase.okNoOffs, ok)
			if !testCase.okNoOffs {
				require.Equal(t, 1, len(violations))
				require.Equal(t, msgValidISODatetimeFormatNoOffs, violations[0].Message)
			}
			ok, violations = vMin.Validate(obj)
			require.Equal(t, testCase.okMin, ok)
			if !testCase.okMin {
				require.Equal(t, 1, len(violations))
				require.Equal(t, msgValidISODatetimeFormatMin, violations[0].Message)
			}
		})
	}
}

func TestStringValidISODate(t *testing.T) {
	validator := buildFooValidator(JsonString,
		&StringValidISODate{}, false)
	obj := jsonObject(`{
		"foo": "2022-02-02 but not a date with this on the end"
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidISODate, violations[0].Message)

	obj["foo"] = "2022-02-02"
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	obj["foo"] = "2022-13-02"
	//                  ^^ 13th month?
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidISODate, violations[0].Message)

	// should also fail with time specified...
	obj["foo"] = "2022-13-02T18:19:20"
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgValidISODate, violations[0].Message)
}

var variousDatetimeFormats = []string{
	"2006-01-02T15:04:05.999999999-07:00",
	"2006-01-02T15:04:05.999999999Z",
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05-07:00",
	"2006-01-02T15:04:05Z",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

func TestDatetimeFuture(t *testing.T) {
	pastTime := time.Now().Add(0 - (5 * time.Hour))
	validator := buildFooValidator(JsonAny,
		&DatetimeFuture{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFuture, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = pastTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFuture, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = pastTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, msgDatetimeFuture, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFuture, violations[0].Message)
}

func TestDatetimeFutureExcTime(t *testing.T) {
	constraint := &DatetimeFuture{ExcTime: true}
	validator := buildFooValidator(JsonAny, constraint, false)
	futureTime := time.Now().Add(2 * time.Second)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFuture, violations[0].Message)

	constraint.ExcTime = false
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimeFutureOrPresent(t *testing.T) {
	pastTime := time.Now().Add(0 - (5 * time.Hour))
	validator := buildFooValidator(JsonAny,
		&DatetimeFutureOrPresent{}, false)
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFutureOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = pastTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFutureOrPresent, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = pastTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, msgDatetimeFutureOrPresent, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFutureOrPresent, violations[0].Message)
}

func TestDatetimeFutureOrPresentExcTime(t *testing.T) {
	constraint := &DatetimeFutureOrPresent{ExcTime: true}
	validator := buildFooValidator(JsonAny, constraint, false)
	pastTime := time.Now().Add(0 - (2 * time.Second))
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, _ := validator.Validate(obj)
	require.True(t, ok)

	constraint.ExcTime = false
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimeFutureOrPresent, violations[0].Message)
}

func TestDatetimePast(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	validator := buildFooValidator(JsonAny,
		&DatetimePast{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePast, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = futureTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePast, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = futureTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, msgDatetimePast, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePast, violations[0].Message)
}

func TestDatetimePastExcTime(t *testing.T) {
	constraint := &DatetimePast{ExcTime: true}
	validator := buildFooValidator(JsonAny, constraint, false)
	pastTime := time.Now().Add(0 - (2 * time.Second))
	obj := map[string]interface{}{
		"foo": pastTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePast, violations[0].Message)

	constraint.ExcTime = false
	ok, _ = validator.Validate(obj)
	require.True(t, ok)
}

func TestDatetimePastOrPresent(t *testing.T) {
	futureTime := time.Now().Add(24 * time.Hour)
	validator := buildFooValidator(JsonAny,
		&DatetimePastOrPresent{}, false)
	obj := map[string]interface{}{
		"foo": futureTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePastOrPresent, violations[0].Message)

	obj["foo"] = time.Now().Add(0 - time.Minute).Format("2006-01-02T15:04:05.000000000-07:00")
	ok, _ = validator.Validate(obj)
	require.True(t, ok)

	// check with actual time.Time...
	obj["foo"] = futureTime
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePastOrPresent, violations[0].Message)

	// check with varying formats...
	for _, layout := range variousDatetimeFormats {
		obj["foo"] = futureTime.Format(layout)
		ok, violations = validator.Validate(obj)
		require.False(t, ok)
		require.Equal(t, 1, len(violations))
		require.Equal(t, msgDatetimePastOrPresent, violations[0].Message)
	}

	// and finally with invalid datetime...
	obj["foo"] = ""
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePastOrPresent, violations[0].Message)
}

func TestDatetimePastOrPresentExcTime(t *testing.T) {
	constraint := &DatetimePastOrPresent{ExcTime: true}
	validator := buildFooValidator(JsonAny, constraint, false)
	testTime := time.Now().Add(2 * time.Second)
	obj := map[string]interface{}{
		"foo": testTime.Format("2006-01-02T15:04:05.000000000-07:00"),
	}
	ok, violations := validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	constraint.ExcTime = false
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgDatetimePastOrPresent, violations[0].Message)
}

func TestStringValidTimezone(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &StringValidTimezone{}
	testCases := []struct {
		value        string
		locationOnly bool
		offsetOnly   bool
		expected     bool
	}{
		{
			"UTC",
			true,
			false,
			true,
		},
		{
			"UTC",
			false,
			false,
			true,
		},
		{
			"Europe/London",
			true,
			false,
			true,
		},
		{
			"America/Los_Angeles",
			true,
			false,
			true,
		},
		{
			"Zulu",
			true,
			false,
			true,
		},
		{
			"UTC",
			false,
			true,
			false,
		},
		{
			"Europe/London",
			false,
			true,
			false,
		},
		{
			"America/Los_Angeles",
			false,
			true,
			false,
		},
		{
			"Zulu",
			false,
			true,
			false,
		},
		{
			"+00:00",
			false,
			false,
			true,
		},
		{
			"+00:00",
			false,
			true,
			true,
		},
		{
			"+00:00",
			true,
			false,
			false,
		},
		{
			"+00:00",
			true,
			true,
			false,
		},
		{
			"+14:00",
			false,
			false,
			true,
		},
		{
			"+14",
			false,
			false,
			true,
		},
		{
			"14",
			false,
			false,
			true,
		},
		{
			"+15:00",
			false,
			false,
			false,
		},
		{
			"15:00",
			false,
			false,
			false,
		},
		{
			"-12:00",
			false,
			false,
			true,
		},
		{
			"-13:00",
			false,
			false,
			false,
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("StringValidTimezone[%d]:\"%s\"", i+1, tc.value), func(t *testing.T) {
			c.LocationOnly = tc.locationOnly
			c.OffsetOnly = tc.offsetOnly
			ok, msg := c.Check(tc.value, vcx)
			if tc.expected {
				require.True(t, ok)
			} else {
				require.False(t, ok)
				require.Equal(t, msgValidTimezone, msg)
			}
		})
	}

	c.AllowNumeric = true
	c.LocationOnly = false
	c.OffsetOnly = false
	ok, _ := c.Check(-12, vcx)
	require.True(t, ok)
	ok, _ = c.Check(-13, vcx)
	require.False(t, ok)
	ok, _ = c.Check(14, vcx)
	require.True(t, ok)
	ok, _ = c.Check(15, vcx)
	require.False(t, ok)

	ok, _ = c.Check(0, vcx)
	require.True(t, ok)
	c.LocationOnly = true
	ok, _ = c.Check(0, vcx)
	require.False(t, ok)
}

func TestDatetimeDayOfWeek(t *testing.T) {
	vcx := newValidatorContext(nil, nil, false, nil)
	c := &DatetimeDayOfWeek{Days: "54321"}
	day := time.Hour * 24
	for i := time.Duration(-8); i < 9; i++ {
		dt := time.Now().Add(i * day)
		ok, msg := c.Check(dt, vcx)
		okS, _ := c.Check(dt.Format(time.RFC3339), vcx)
		if dt.Weekday() > 0 && dt.Weekday() < 6 {
			require.True(t, ok)
			require.True(t, okS)
		} else {
			require.False(t, ok)
			require.False(t, okS)
			require.Equal(t, msgDatetimeDayOfWeek, msg)
		}
	}
}

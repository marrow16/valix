package valix

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRequestQueryValidate(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"str": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringNotBlank{},
				},
			},
			"num": {
				Type:      JsonNumber,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThan{Value: 1},
				},
			},
			"int": {
				Type:      JsonInteger,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThan{Value: 1},
				},
			},
			"bool": {
				Type:      JsonBoolean,
				Mandatory: true,
			},
			"dt": {
				Type:      JsonDatetime,
				Mandatory: true,
				Constraints: Constraints{
					&DatetimeGreaterThan{Value: "2001-01-01T12:00:00"},
				},
			},
			"arr": {
				Type:      JsonArray,
				Mandatory: true,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenString,
						Constraints: Constraints{
							&StringGreaterThan{Value: "B"},
						},
					},
				},
			},
			"obj": {
				Type:      JsonObject,
				Mandatory: true,
				Constraints: Constraints{
					&NotEmpty{},
				},
			},
			"any": {
				Type:      JsonAny,
				Mandatory: true,
			},
		},
	}
	params := url.Values{
		"str":  []string{"string"},
		"num":  []string{"1.1"},
		"int":  []string{"2"},
		"bool": []string{"true"},
		"dt":   []string{"2022-09-03T11:00:00"},
		"arr":  []string{"element1", "element2"},
		"obj":  []string{"{\"foo\":\"bar\"}"},
		"any":  []string{"a"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)

	ok, violations, _ := v.RequestQueryValidate(req)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	params = url.Values{
		"str": []string{" "},
		"num": []string{"0.1"},
		"int": []string{"0"},
		//"bool": []string{"true"},
		"dt":  []string{"1999-09-03T11:00:00"},
		"arr": []string{"A", "A"},
		"obj": []string{"{}"},
		//"any":  []string{"a"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)

	ok, violations, _ = v.RequestQueryValidate(req)
	require.False(t, ok)
	require.Equal(t, 9, len(violations))

	// with bad query param...
	params = url.Values{
		"bool": []string{"not a valid boolean"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)

	ok, violations, _ = v.RequestQueryValidate(req)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "boolean"), violations[0].Message)
}

func TestRequestQueryValidateInto(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"str": {
				Type:      JsonString,
				Mandatory: true,
				Constraints: Constraints{
					&StringNotBlank{},
				},
			},
			"num": {
				Type:      JsonNumber,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThan{Value: 1},
				},
			},
			"int": {
				Type:      JsonInteger,
				Mandatory: true,
				Constraints: Constraints{
					&GreaterThan{Value: 1},
				},
			},
			"bool": {
				Type:      JsonBoolean,
				Mandatory: true,
			},
			"dt": {
				Type:      JsonDatetime,
				Mandatory: true,
				Constraints: Constraints{
					&DatetimeGreaterThan{Value: "2001-01-01T12:00:00"},
				},
			},
			"arr": {
				Type:      JsonArray,
				Mandatory: true,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenString,
						Constraints: Constraints{
							&StringGreaterThan{Value: "B"},
						},
					},
				},
			},
			"obj": {
				Type:      JsonObject,
				Mandatory: true,
				Constraints: Constraints{
					&NotEmpty{},
				},
			},
			"any": {
				Type:      JsonAny,
				Mandatory: true,
			},
		},
	}
	params := url.Values{
		"str":  []string{"string"},
		"num":  []string{"1.1"},
		"int":  []string{"2"},
		"bool": []string{"true"},
		"dt":   []string{"2022-09-03T11:00:00"},
		"arr":  []string{"element1", "element2"},
		"obj":  []string{"{\"foo\":\"bar\"}"},
		"any":  []string{"a"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)

	type paramsStruct struct {
		Str  string                 `json:"str"`
		Num  float64                `json:"num"`
		Int  int                    `json:"int"`
		Bool bool                   `json:"bool"`
		Dt   *time.Time             `json:"dt"`
		Arr  []string               `json:"arr"`
		Obj  map[string]interface{} `json:"obj"`
		Any  interface{}            `json:"any"`
	}
	reqParams := &paramsStruct{}

	ok, violations, _ := v.RequestQueryValidateInto(req, reqParams)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	params = url.Values{
		"str": []string{" "},
		"num": []string{"0.1"},
		"int": []string{"0"},
		//"bool": []string{"true"},
		"dt":  []string{"1999-09-03T11:00:00"},
		"arr": []string{"A", "A"},
		"obj": []string{"{}"},
		//"any":  []string{"a"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)

	ok, violations, _ = v.RequestQueryValidateInto(req, reqParams)
	require.False(t, ok)
	require.Equal(t, 9, len(violations))

	// with bad query param...
	params = url.Values{
		"bool": []string{"not a valid boolean"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)

	ok, violations, _ = v.RequestQueryValidateInto(req, reqParams)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "boolean"), violations[0].Message)

	// with forced encode failure...
	type badParamsStruct struct {
		Str string `json:"str"`
	}
	badReqParams := &badParamsStruct{}
	params = url.Values{
		"str":  []string{"string"},
		"num":  []string{"1.1"},
		"int":  []string{"2"},
		"bool": []string{"true"},
		"dt":   []string{"2022-09-03T11:00:00"},
		"arr":  []string{"element1", "element2"},
		"obj":  []string{"{\"foo\":\"bar\"}"},
		"any":  []string{"a"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	ok, violations, _ = v.RequestQueryValidateInto(req, badReqParams)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgErrorUnmarshall, violations[0].Message)
	require.True(t, violations[0].BadRequest)
}

func TestQueryParamsToObject(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"str": {
				Type: JsonString,
			},
			"num": {
				Type: JsonNumber,
			},
			"int": {
				Type: JsonInteger,
			},
			"bool": {
				Type: JsonBoolean,
			},
			"dt": {
				Type: JsonDatetime,
			},
			"arr": {
				Type: JsonArray,
			},
			"obj": {
				Type: JsonObject,
			},
			"any": {
				Type: JsonAny,
			},
		},
	}
	params := url.Values{
		"str":  []string{"string"},
		"num":  []string{"1.1"},
		"int":  []string{"2"},
		"bool": []string{"true"},
		"dt":   []string{"2022-09-03T11:00:00"},
		"arr":  []string{"element1", "element2"},
		"obj":  []string{"{\"foo\":\"bar\"}"},
		"any":  []string{"a"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)

	qpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.Equal(t, 8, len(qpo))
	require.Equal(t, "string", qpo["str"])
	jn, _ := qpo["num"].(json.Number)
	fv, _ := jn.Float64()
	require.Equal(t, 1.1, fv)
	jn, _ = qpo["int"].(json.Number)
	iv, _ := jn.Int64()
	require.Equal(t, int64(2), iv)
	m := qpo["obj"].(map[string]interface{})
	require.Equal(t, "bar", m["foo"])
	a := qpo["arr"].([]interface{})
	require.Equal(t, 2, len(a))
	require.Equal(t, "element1", a[0])
	require.Equal(t, "element2", a[1])
	any := qpo["any"].([]interface{})
	require.Equal(t, 1, len(any))
	require.Equal(t, "a", any[0])
}

func TestQueryParamsToObject_WithArrayElementTyped(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"arr": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenDatetime,
					},
				},
			},
		},
	}
	params := url.Values{
		"arr": []string{"2022-09-03T11:00:00", "2022-09-03T11:00:00"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)

	qpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.Equal(t, 1, len(qpo))
	a := qpo["arr"].([]interface{})
	require.Equal(t, 2, len(a))
	dt0 := a[0].(*time.Time)
	require.Equal(t, 2022, dt0.Year())
}

func TestQueryParamsToObject_NoProperties(t *testing.T) {
	v := &Validator{}
	params := url.Values{
		"foo": []string{"bar"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.NotEmpty(t, qpo)
	require.Equal(t, "bar", qpo["foo"])

	params = url.Values{
		"foo": []string{""},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations = v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.NotEmpty(t, qpo)
	require.Equal(t, true, qpo["foo"])

	params = url.Values{
		"foo": []string{"bar", "baz"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations = v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.NotEmpty(t, qpo)
	a := qpo["foo"].([]interface{})
	require.Equal(t, 2, len(a))
	require.Equal(t, "bar", a[0])

	params = url.Values{
		"foo": []string{"", ""},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations = v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.NotEmpty(t, qpo)
	a = qpo["foo"].([]interface{})
	require.Equal(t, 2, len(a))
	require.Equal(t, true, a[0])
	require.Equal(t, true, a[1])
}

func TestQueryParamsToObject_NoMultiValue(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonString,
			},
		},
	}
	params := url.Values{
		"foo": []string{"bar", "baz"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations := v.queryParamsToObject(req, nil)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgQueryParamMultiNotAllowed, violations[0].Message)
	require.Equal(t, "foo", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.True(t, violations[0].BadRequest)
	require.Equal(t, CodeRequestQueryParamMultiNotAllowed, violations[0].Codes[0])
}

func TestQueryParamsToObject_BadDate(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonDatetime,
			},
		},
	}
	params := url.Values{
		"foo": []string{"not a valid date"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.Equal(t, "not a valid date", qpo["foo"])

	v = &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenDatetime,
					},
				},
			},
		},
	}
	params = url.Values{
		"foo": []string{"not a valid date", "not a valid date"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations = v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	a := qpo["foo"].([]interface{})
	require.Equal(t, 2, len(a))
}

func TestQueryParamsToObject_BooleanEmpty(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonBoolean,
			},
		},
	}
	params := url.Values{
		"foo": []string{""},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.Equal(t, true, qpo["foo"])

	v = &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenBoolean,
					},
				},
			},
		},
	}
	params = url.Values{
		"foo": []string{"", ""},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	qpo, violations = v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	a := qpo["foo"].([]interface{})
	require.Equal(t, true, a[0])
	require.Equal(t, true, a[1])
}

func TestQueryParamsToObject_BadBoolean(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonBoolean,
			},
		},
	}
	params := url.Values{
		"foo": []string{"not a valid bool"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations := v.queryParamsToObject(req, nil)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "boolean"), violations[0].Message)

	v = &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenBoolean,
					},
				},
			},
		},
	}
	params = url.Values{
		"foo": []string{"not a valid bool", "not a valid bool"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations = v.queryParamsToObject(req, nil)
	require.Equal(t, 2, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "boolean"), violations[0].Message)
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "boolean"), violations[1].Message)
}

func TestQueryParamsToObject_BadObject(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonObject,
			},
		},
	}
	params := url.Values{
		"foo": []string{"not a valid object"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations := v.queryParamsToObject(req, nil)
	require.Equal(t, 1, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "object"), violations[0].Message)

	v = &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenObject,
					},
				},
			},
		},
	}
	params = url.Values{
		"foo": []string{"not a valid object", "not a valid object"},
	}
	qs = params.Encode()
	req, _ = http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations = v.queryParamsToObject(req, nil)
	require.Equal(t, 2, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "object"), violations[0].Message)
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "object"), violations[1].Message)
}

func TestQueryParamsToObject_ArrayOfArray(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenArray,
					},
				},
			},
		},
	}
	params := url.Values{
		"foo": []string{"[\"a\"]", "[\"b\"]"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	gpo, violations := v.queryParamsToObject(req, nil)
	require.Empty(t, violations)
	require.NotEmpty(t, gpo)
	a := gpo["foo"].([]interface{})
	require.Equal(t, 2, len(a))
	a0 := a[0].([]interface{})
	require.Equal(t, 1, len(a0))
	require.Equal(t, "a", a0[0])
}

func TestQueryParamsToObject_BadArray(t *testing.T) {
	v := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonArray,
				Constraints: Constraints{
					&ArrayOf{
						Type: jsonTypeTokenArray,
					},
				},
			},
		},
	}
	params := url.Values{
		"foo": []string{"not a valid array", "not a valid array"},
	}
	qs := params.Encode()
	req, _ := http.NewRequest("POST", "some.url?"+qs, nil)
	_, violations := v.queryParamsToObject(req, nil)
	require.Equal(t, 2, len(violations))
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "array"), violations[0].Message)
	require.Equal(t, fmt.Sprintf(fmtMsgQueryParamType, "array"), violations[1].Message)
}

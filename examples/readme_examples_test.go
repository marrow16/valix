package examples

import (
	"encoding/json"
	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode"
)

type AddPersonRequest struct {
	Name string `json:"name" v8n:"notNull,mandatory,constraints:[StringNoControlCharacters{},StringLength{Minimum: 1, Maximum: 255}]"`
	Age  int    `json:"age" v8n:"type:Integer,notNull,mandatory,constraint:PositiveOrZero{}"`
}

var AddPersonRequestValidator = valix.MustCompileValidatorFor(AddPersonRequest{}, nil)

/*
// Abbreviated forms of above
type AddPersonRequest struct {
	Name string `json:"name" v8n:"notNull,mandatory,&StringNoControlCharacters{},&StringLength{Minimum: 1, Maximum: 255}"`
	Age  int    `json:"age" v8n:"type:Integer,notNull,mandatory,&PositiveOrZero{}"`
}

var AddPersonRequestValidator = valix.MustCompileValidatorFor(AddPersonRequest{}, nil)
*/

// CreatePersonRequestValidator is the same as AddPersonRequestValidator (just expressed in code rather than as tags on struct)
var CreatePersonRequestValidator = &valix.Validator{
	IgnoreUnknownProperties: false,
	Properties: valix.Properties{
		"name": {
			Type:      valix.JsonString,
			NotNull:   true,
			Mandatory: true,
			Constraints: valix.Constraints{
				&valix.StringNoControlCharacters{},
				&valix.StringLength{Minimum: 1, Maximum: 255},
			},
		},
		"age": {
			Type:      valix.JsonInteger,
			NotNull:   true,
			Mandatory: true,
			Constraints: valix.Constraints{
				&valix.PositiveOrZero{},
			},
		},
	},
}

func TestCreatePersonRequestValidatorFails(t *testing.T) {
	body := strings.NewReader(`{"name": "", "age": -1}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AddPersonHandler(res, req)

	require.Equal(t, http.StatusUnprocessableEntity, res.Code)
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	resBody := map[string]interface{}{}
	_ = json.NewDecoder(res.Body).Decode(&resBody)
	require.Equal(t, "Request invalid", resBody["$error"])
	details := resBody["$details"].([]interface{})
	require.Equal(t, 2, len(details))
	first := details[0].(map[string]interface{})
	require.Equal(t, "Value must be positive or zero", first["message"])
	require.Equal(t, "", first["path"])
	require.Equal(t, "age", first["property"])
	second := details[1].(map[string]interface{})
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", second["message"])
	require.Equal(t, "", second["path"])
	require.Equal(t, "name", second["property"])
}

func TestCreatePersonRequestValidatorSucceeds(t *testing.T) {
	body := strings.NewReader(`{"name": "Bilbo Baggins", "age": 25}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AddPersonHandler(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	resBody := map[string]interface{}{}
	_ = json.NewDecoder(res.Body).Decode(&resBody)
	require.Equal(t, "Bilbo Baggins", resBody["name"])
	require.Equal(t, float64(25), resBody["age"])
}

func AddPersonHandler(w http.ResponseWriter, r *http.Request) {
	requestBody := &AddPersonRequest{}
	ok, violations, _ := CreatePersonRequestValidator.RequestValidateInto(r, requestBody)
	if !ok {
		// write an error response with full info - using violations information
		valix.SortViolationsByPathAndProperty(violations)
		errResponse := map[string]interface{}{
			"$error":   "Request invalid",
			"$details": violations,
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(errResponse)
		return
	}
	// the request will now be a validated struct
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(requestBody)
}

func TestCreatePerson2RequestValidatorFails(t *testing.T) {
	body := strings.NewReader(`{"name": "", "age": -1}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AddPersonHandler2(res, req)

	require.Equal(t, http.StatusUnprocessableEntity, res.Code)
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	resBody := map[string]interface{}{}
	_ = json.NewDecoder(res.Body).Decode(&resBody)
	require.Equal(t, "Request invalid", resBody["$error"])
	details := resBody["$details"].([]interface{})
	require.Equal(t, 2, len(details))
	first := details[0].(map[string]interface{})
	require.Equal(t, "Value must be positive or zero", first["message"])
	require.Equal(t, "", first["path"])
	require.Equal(t, "age", first["property"])
	second := details[1].(map[string]interface{})
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", second["message"])
	require.Equal(t, "", second["path"])
	require.Equal(t, "name", second["property"])
}

func TestCreatePerson2RequestValidatorSucceeds(t *testing.T) {
	body := strings.NewReader(`{"name": "Bilbo Baggins", "age": 25}`)
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}
	res := httptest.NewRecorder()

	AddPersonHandler2(res, req)

	require.Equal(t, http.StatusOK, res.Code)
	require.Equal(t, "application/json", res.Header().Get("Content-Type"))
	resBody := map[string]interface{}{}
	_ = json.NewDecoder(res.Body).Decode(&resBody)
	require.Equal(t, "Bilbo Baggins", resBody["name"])
	require.Equal(t, float64(25), resBody["age"])
}

func AddPersonHandler2(w http.ResponseWriter, r *http.Request) {
	ok, violations, obj := CreatePersonRequestValidator.RequestValidate(r)
	if !ok {
		// write an error response with full info - using violations information
		valix.SortViolationsByPathAndProperty(violations)
		errResponse := map[string]interface{}{
			"$error":   "Request invalid",
			"$details": violations,
		}
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(errResponse)
		return
	}
	// the 'obj' will now be a validated map[string]interface{}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(obj)
}

func TestValidateStringIntoStruct(t *testing.T) {
	str := `{
		"name": "",
		"age": -1
	}`
	req := &AddPersonRequest{}

	ok, violations, _ := AddPersonRequestValidator.ValidateStringInto(str, req)

	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	valix.SortViolationsByPathAndProperty(violations)
	require.Equal(t, "Value must be positive or zero", violations[0].Message)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", violations[1].Message)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, "", violations[1].Path)

	str = `{
		"name": "Bilbo Baggins",
		"age": 25
	}`
	ok, violations, _ = AddPersonRequestValidator.ValidateStringInto(str, req)

	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "Bilbo Baggins", req.Name)
	require.Equal(t, 25, req.Age)
}

func TestValidateReaderIntoStruct(t *testing.T) {
	reader := strings.NewReader(`{
		"name": "",
		"age": -1
	}`)
	req := &AddPersonRequest{}

	ok, violations, _ := AddPersonRequestValidator.ValidateReaderInto(reader, req)

	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	valix.SortViolationsByPathAndProperty(violations)
	require.Equal(t, "Value must be positive or zero", violations[0].Message)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", violations[1].Message)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, "", violations[1].Path)

	reader = strings.NewReader(`{
		"name": "Bilbo Baggins",
		"age": 25
	}`)
	ok, violations, _ = AddPersonRequestValidator.ValidateReaderInto(reader, req)

	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.Equal(t, "Bilbo Baggins", req.Name)
	require.Equal(t, 25, req.Age)
}

func TestValidateMap(t *testing.T) {
	req := map[string]interface{}{
		"name": "",
		"age":  -1,
	}

	ok, violations := AddPersonRequestValidator.Validate(req)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	valix.SortViolationsByPathAndProperty(violations)
	require.Equal(t, "Value must be positive or zero", violations[0].Message)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", violations[1].Message)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, "", violations[1].Path)

	req = map[string]interface{}{
		"name": "Bilbo Baggins",
		"age":  25,
	}
	ok, _ = AddPersonRequestValidator.Validate(req)
	require.True(t, ok)
}

func TestValidateSlice(t *testing.T) {
	req := []interface{}{
		map[string]interface{}{
			"name": "",
			"age":  -1,
		},
		map[string]interface{}{
			"name": "Bilbo Baggins",
			"age":  25,
		},
	}

	ok, violations := AddPersonRequestValidator.ValidateArrayOf(req)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	valix.SortViolationsByPathAndProperty(violations)
	require.Equal(t, "Value must be positive or zero", violations[0].Message)
	require.Equal(t, "age", violations[0].Property)
	require.Equal(t, "[0]", violations[0].Path)
	require.Equal(t, "String value length must be between 1 and 255 (inclusive)", violations[1].Message)
	require.Equal(t, "name", violations[1].Property)
	require.Equal(t, "[0]", violations[1].Path)

	req = []interface{}{
		map[string]interface{}{
			"name": "Frodo Baggins",
			"age":  20,
		},
		map[string]interface{}{
			"name": "Bilbo Baggins",
			"age":  25,
		},
	}
	ok, _ = AddPersonRequestValidator.ValidateArrayOf(req)
	require.True(t, ok)
}

var mySet = &valix.ConstraintSet{
	Constraints: valix.Constraints{
		&valix.StringTrim{},
		&valix.StringNotEmpty{},
		&valix.StringLength{Minimum: 16, Maximum: 64},
		valix.NewCustomConstraint(func(value interface{}, vcx *valix.ValidatorContext, this *valix.CustomConstraint) (bool, string) {
			if str, ok := value.(string); ok {
				if len(str) == 0 || str[0] < 'A' || str[0] > 'Z' {
					return false, this.GetMessage()
				}
			}
			return true, ""
		}, ""),
		&valix.StringCharacters{
			AllowRanges: []*unicode.RangeTable{
				{R16: []unicode.Range16{{'0', 'z', 1}}},
			},
			DisallowRanges: []*unicode.RangeTable{
				{R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
				{R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
				{R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
			},
		},
	},
	Message: "String value length must be between 16 and 64 chars; must be letters (upper or lower), digits or underscores; must start with an uppercase letter",
}

func TestMySet(t *testing.T) {
	validator := &valix.Validator{
		Properties: valix.Properties{
			"foo": {
				Type: valix.JsonString,
				Constraints: valix.Constraints{
					mySet,
				},
			},
		},
	}
	request := buildTestRequest(t, `{"foo": "A234567890123456_aaa"}`)

	ok, violations, obj := validator.RequestValidate(request)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	request = buildTestRequest(t, `{"foo": "a234567890123456_aaa"}`)
	ok, violations, obj = validator.RequestValidate(request)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.NotNil(t, obj)

	request = buildTestRequest(t, `{"foo": "Aaa"}`)
	ok, violations, obj = validator.RequestValidate(request)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.NotNil(t, obj)
}

func TestMyValidator(t *testing.T) {
	myValidator := &valix.Validator{
		IgnoreUnknownProperties: false,
		Properties: valix.Properties{
			"foo": {
				Type:      valix.JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: valix.Constraints{
					valix.NewCustomConstraint(func(value interface{}, ctx *valix.ValidatorContext, cc *valix.CustomConstraint) (bool, string) {
						if str, ok := value.(string); ok {
							return !strings.Contains(str, "foo"), cc.GetMessage()
						}
						return true, ""
					}, "Value must not contain \"foo\""),
				},
			},
		},
	}
	request := buildTestRequest(t, `{"foo": "bar"}`)

	ok, violations, obj := myValidator.RequestValidate(request)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	request = buildTestRequest(t, `{"foo": "foo bar"}`)

	ok, violations, obj = myValidator.RequestValidate(request)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.NotNil(t, obj)
}

type NoFoo struct {
}

func (c *NoFoo) Check(value interface{}, vcx *valix.ValidatorContext) (bool, string) {
	if str, ok := value.(string); ok {
		return !strings.Contains(str, "foo"), c.GetMessage()
	}
	return true, ""
}
func (c *NoFoo) GetMessage() string {
	return "Value must not contain \"foo\""
}

func TestNoFooCustomConstraint(t *testing.T) {
	myValidator := &valix.Validator{
		IgnoreUnknownProperties: false,
		Properties: valix.Properties{
			"foo": {
				Type:      valix.JsonString,
				NotNull:   true,
				Mandatory: true,
				Constraints: valix.Constraints{
					&NoFoo{},
				},
			},
		},
	}
	request := buildTestRequest(t, `{"foo": "bar"}`)

	ok, violations, obj := myValidator.RequestValidate(request)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))
	require.NotNil(t, obj)

	request = buildTestRequest(t, `{"foo": "foo bar"}`)

	ok, violations, obj = myValidator.RequestValidate(request)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.NotNil(t, obj)
}

func buildTestRequest(t *testing.T, body string) *http.Request {
	r := strings.NewReader(body)
	req, err := http.NewRequest("POST", "", r)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

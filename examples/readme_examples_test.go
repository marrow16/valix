package examples

import (
	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
	"net/http"
	"strings"
	"testing"
	"unicode"
)

var personValidator = &valix.Validator{
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

func TestAddPersonRequest(t *testing.T) {
	body := strings.NewReader(`{"name": "", "age": -1}`)
	request, err := http.NewRequest("POST", "", body)
	if err != nil {
		t.Fatal(err)
	}

	ok, violations, obj := personValidator.RequestValidate(request)
	require.False(t, ok)
	require.Equal(t, 2, len(violations))
	require.NotNil(t, obj)
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

type AddPerson struct {
	Name string `json:"name" v8n:"notNull,mandatory,constraints:[StringNoControlCharacters{},StringLength{Minimum: 1, Maximum: 255}]"`
	Age  int    `json:"age" v8n:"type:Integer,notNull,mandatory,constraint:PositiveOrZero{}"`
}

var AddPersonValidator = valix.MustCompileValidatorFor(AddPerson{}, nil)

type AddPersonAbbreviated struct {
	Name string `json:"name" v8n:"notNull,mandatory,&StringNoControlCharacters{},&StringLength{Minimum: 1, Maximum: 255}"`
	Age  int    `json:"age" v8n:"type:Integer,notNull,mandatory,&PositiveOrZero{}"`
}

var AddPersonAbbreviatedValidator = valix.MustCompileValidatorFor(AddPerson{}, nil)

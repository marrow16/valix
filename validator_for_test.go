package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestShouldPanicWithNonStruct(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, errMsgValidatorForStructOnly, r.(error).Error())
		}
	}()

	_ = MustCompileValidatorFor("", nil)
}

func TestMustCompileValidatorForEmptyStruct(t *testing.T) {
	myStruct := struct{}{}

	v := MustCompileValidatorFor(myStruct, nil)
	require.NotNil(t, v)
	require.False(t, v.IgnoreUnknownProperties)
	require.Equal(t, 0, len(v.Properties))
	require.Nil(t, v.Constraints)
	require.False(t, v.AllowArray)
	require.False(t, v.DisallowObject)
	require.False(t, v.UseNumber)
}

func TestShouldErrorWithNonStruct(t *testing.T) {
	_, err := ValidatorFor("", nil)
	require.NotNil(t, err)
	require.Equal(t, errMsgValidatorForStructOnly, err.Error())
}

func TestValidatorForWithEmptyStruct(t *testing.T) {
	myStruct := struct{}{}

	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.False(t, v.IgnoreUnknownProperties)
	require.Equal(t, 0, len(v.Properties))
	require.Nil(t, v.Constraints)
	require.False(t, v.AllowArray)
	require.False(t, v.DisallowObject)
	require.False(t, v.UseNumber)
}

func TestValidatorForWithEmptyStructAndOptions(t *testing.T) {
	myStruct := struct{}{}

	v, err := ValidatorFor(myStruct, &ValidatorForOptions{
		IgnoreUnknownProperties: true,
		Constraints:             Constraints{},
		AllowNullJson:           true,
		UseNumber:               true,
	})
	require.Nil(t, err)
	require.NotNil(t, v)
	require.True(t, v.IgnoreUnknownProperties)
	require.Equal(t, 0, len(v.Properties))
	require.NotNil(t, v.Constraints)
	require.Equal(t, 0, len(v.Constraints))
	require.False(t, v.AllowArray)
	require.False(t, v.DisallowObject)
	require.True(t, v.UseNumber)
}

func TestValidatorForWithJsonTag(t *testing.T) {
	myStruct := struct {
		Foo string `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	_, ok := v.Properties["foo"]
	require.True(t, ok)

	myStruct2 := struct {
		Foo int `json:"foo,omitempty"`
	}{}
	v, err = ValidatorFor(myStruct2, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	_, ok = v.Properties["foo"]
	require.True(t, ok)
}

func TestValidatorForDetectsTypeString(t *testing.T) {
	myStruct := struct {
		Foo string `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonString, pv.Type)
}

func TestValidatorForDetectsTypeNumber(t *testing.T) {
	testCases := []interface{}{
		struct {
			Foo int `json:"foo"`
		}{},
		struct {
			Foo int8 `json:"foo"`
		}{},
		struct {
			Foo int16 `json:"foo"`
		}{},
		struct {
			Foo int32 `json:"foo"`
		}{},
		struct {
			Foo int64 `json:"foo"`
		}{},
		struct {
			Foo uint `json:"foo"`
		}{},
		struct {
			Foo uint8 `json:"foo"`
		}{},
		struct {
			Foo uint16 `json:"foo"`
		}{},
		struct {
			Foo uint32 `json:"foo"`
		}{},
		struct {
			Foo uint64 `json:"foo"`
		}{},
		struct {
			Foo float32 `json:"foo"`
		}{},
		struct {
			Foo float64 `json:"foo"`
		}{},
	}
	for i, st := range testCases {
		t.Run(fmt.Sprintf("NumericTypeDetect[%d]", i), func(t *testing.T) {
			v, err := ValidatorFor(st, nil)
			require.Nil(t, err)
			require.NotNil(t, v)
			require.Equal(t, 1, len(v.Properties))
			pv, ok := v.Properties["foo"]
			require.True(t, ok)
			require.NotNil(t, pv)
			require.Equal(t, 0, len(pv.Constraints))
			require.Equal(t, JsonNumber, pv.Type)
		})
	}
}

func TestValidatorForDetectsTypeBoolean(t *testing.T) {
	myStruct := struct {
		Foo bool `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonBoolean, pv.Type)
}

func TestValidatorForDetectsTypeObject(t *testing.T) {
	myStruct := struct {
		Foo struct{ Sub string } `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonObject, pv.Type)
}

func TestValidatorForDetectsTypeMap(t *testing.T) {
	myStruct := struct {
		Foo map[string]interface{} `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonObject, pv.Type)

	myStruct2 := struct {
		Foo map[int]interface{} `json:"foo"`
	}{}
	v, err = ValidatorFor(myStruct2, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok = v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonTypeUndefined, pv.Type)
}

func TestValidatorForDetectsTypeArray(t *testing.T) {
	myStruct := struct {
		Foo []string `json:"foo"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.NotNil(t, pv)
	require.Equal(t, 0, len(pv.Constraints))
	require.Equal(t, JsonArray, pv.Type)
}

func TestValidatorForWithBadV8nTags(t *testing.T) {
	testCases := []interface{}{
		struct {
			Foo string `v8n:"unknown_token"`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis("`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis["`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis{"`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis)"`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis]"`
		}{},
		struct {
			Foo string `v8n:"unbalanced_parenthesis}"`
		}{},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("BadV8nTags[%d]", i), func(t *testing.T) {
			_, err := ValidatorFor(tc, nil)
			require.NotNil(t, err)
		})
	}
}

func TestValidatorForWithV8nTag(t *testing.T) {
	myStruct := struct {
		Foo string `json:"foo" v8n:"type:string,notNull,mandatory,constraints:[StringValidToken{Tokens:[\"A\",\"B\"]},StringLength{Minimum: 16, Maximum: 64, UseRuneLen:true, Message:\"Oh fooey\"},StringPattern{Regexp:\"^([a-fA-F0-9]{8})$\"}]"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, len(v.Properties))
	pv, ok := v.Properties["foo"]
	require.True(t, ok)
	require.Equal(t, 3, len(pv.Constraints))
}

func TestValidatorForWithOrderingTags(t *testing.T) {
	myStruct := struct {
		Foo struct {
			Aaa string `json:"aaa" v8n:"order:3"`
			Bbb string `json:"bbb" v8n:"order:2"`
			Ccc string `json:"ccc" v8n:"order:1"`
		} `json:"foo" v8n:"obj.ordered"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.True(t, v.Properties["foo"].ObjectValidator.OrderedPropertyChecks)
	require.Equal(t, 3, v.Properties["foo"].ObjectValidator.Properties["aaa"].Order)
	require.Equal(t, 2, v.Properties["foo"].ObjectValidator.Properties["bbb"].Order)
	require.Equal(t, 1, v.Properties["foo"].ObjectValidator.Properties["ccc"].Order)
}

func TestValidatorForWithBadOrderTagValue(t *testing.T) {
	myStruct := struct {
		Foo struct {
			Aaa string `json:"aaa" v8n:"order:not_a_number"`
		} `json:"foo" v8n:"obj.ordered"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, v)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTagValue, tagTokenOrder, "int", "not_a_number"), err.Error())
}

type subStruct struct {
	Foo struct {
		Bar string
	}
	FooBar string
	Int    int
}

func TestValidatorForWithNestedStruct(t *testing.T) {
	myStruct := struct {
		Struct      subStruct `json:"struct"`
		StringSlice []string
		StructSlice []subStruct            `json:"slice"`
		Map         map[string]interface{} `json:"map"`
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.Equal(t, 4, len(v.Properties))
}

func TestValidatorForWithNestedStructTagError(t *testing.T) {
	myStruct := struct {
		Sub1 struct {
			Sub2 struct{} `v8n:"BAD_TOKEN"`
		}
	}{}
	_, err := ValidatorFor(myStruct, nil)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "BAD_TOKEN"), err.Error())
}

func TestValidatorForWithNestedStruct_SetsValidatorTypeCorrectly(t *testing.T) {
	myStruct := struct {
		Sub1 struct {
			SubSub struct {
				Foo string `json:"sub1_sub_foo"`
			}
		}
		Sub2 []struct {
			SubSub struct {
				Foo string `json:"sub1_sub_foo"`
			}
		}
		Sub3 []string
	}{}
	v, err := ValidatorFor(myStruct, nil)
	require.Nil(t, err)
	require.NotNil(t, v)
	require.False(t, v.Properties["Sub1"].ObjectValidator.AllowArray)
	require.False(t, v.Properties["Sub1"].ObjectValidator.DisallowObject)
	require.True(t, v.Properties["Sub2"].ObjectValidator.AllowArray)
	require.True(t, v.Properties["Sub2"].ObjectValidator.DisallowObject)
	require.Nil(t, v.Properties["Sub3"].ObjectValidator)
}

type deepTestStruct struct {
	Foo string `json:"foo"`
	Bar int    `json:"bar"`
	Sub struct {
		Foo    string `json:"subFoo"`
		SubSub struct {
			Foo       string `json:"subSubFoo"`
			Bar       int    `json:"subSubBar" v8n:"type:Integer"`
			SubSubSub struct {
				Foo string `json:"subSubSubFoo"`
			} `json:"subSubSub"`
			SubSubArr []struct {
				Foo string `json:"subSubArrFoo"`
			} `json:"subSubArr"`
			SubSubSlice []string `json:"subSubSlice"`
		} `json:"subSub"`
	} `json:"sub"`
	SliceString []string `json:"sliceString"`
	Arr         []struct {
		Foo    string `json:"arrFoo"`
		ArrSub struct {
			Foo       string `json:"arrSubFoo"`
			Bar       int    `json:"arrSubBar" v8n:"type:Integer"`
			ArrSubSub struct {
				Foo string `json:"arrSubSubFoo"`
			} `json:"arrSubSub"`
			ArrSubArr []struct {
				Foo string `json:"arrSubArrFoo"`
			} `json:"arrSubArr"`
			ArrSubSlice []string `json:"arrSubSlice"`
		} `json:"arrSub"`
	} `json:"arr"`
}

func TestValidatorForFailsWithBadTagInSliceStruct(t *testing.T) {
	_, err := ValidatorFor(struct {
		Slice []struct {
			Foo string `v8n:"bad_token"`
		} `json:"slice"`
	}{}, nil)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "bad_token"), err.Error())
}

func TestValidatorForDeepStruct(t *testing.T) {
	v, err := ValidatorFor(deepTestStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	require.Equal(t, 5, len(v.Properties))
	require.Equal(t, JsonString, v.Properties["foo"].Type)
	require.Nil(t, v.Properties["foo"].ObjectValidator)
	require.Equal(t, JsonNumber, v.Properties["bar"].Type)
	require.Nil(t, v.Properties["bar"].ObjectValidator)
	require.Equal(t, JsonObject, v.Properties["sub"].Type)
	require.NotNil(t, v.Properties["sub"].ObjectValidator)
	require.Equal(t, JsonArray, v.Properties["sliceString"].Type)
	require.Nil(t, v.Properties["sliceString"].ObjectValidator)
	require.Equal(t, JsonArray, v.Properties["arr"].Type)
	require.NotNil(t, v.Properties["arr"].ObjectValidator)

	subV := v.Properties["sub"].ObjectValidator
	require.False(t, subV.DisallowObject)
	require.False(t, subV.AllowArray)
	require.Equal(t, 2, len(subV.Properties))
	require.Equal(t, JsonString, subV.Properties["subFoo"].Type)
	require.Nil(t, subV.Properties["subFoo"].ObjectValidator)
	require.Equal(t, JsonObject, subV.Properties["subSub"].Type)
	require.NotNil(t, subV.Properties["subSub"].ObjectValidator)

	subSubV := subV.Properties["subSub"].ObjectValidator
	require.Equal(t, 5, len(subSubV.Properties))
	require.Equal(t, JsonString, subSubV.Properties["subSubFoo"].Type)
	require.Nil(t, subSubV.Properties["subSubFoo"].ObjectValidator)
	require.Equal(t, JsonInteger, subSubV.Properties["subSubBar"].Type)
	require.Nil(t, subSubV.Properties["subSubFoo"].ObjectValidator)
	require.Equal(t, JsonObject, subSubV.Properties["subSubSub"].Type)
	require.NotNil(t, subSubV.Properties["subSubSub"].ObjectValidator)
	require.Equal(t, 1, len(subSubV.Properties["subSubSub"].ObjectValidator.Properties))
	require.Equal(t, JsonString, subSubV.Properties["subSubSub"].ObjectValidator.Properties["subSubSubFoo"].Type)
	require.Equal(t, JsonArray, subSubV.Properties["subSubArr"].Type)
	require.NotNil(t, subSubV.Properties["subSubArr"].ObjectValidator)
	require.Equal(t, 1, len(subSubV.Properties["subSubArr"].ObjectValidator.Properties))
	require.Equal(t, JsonString, subSubV.Properties["subSubArr"].ObjectValidator.Properties["subSubArrFoo"].Type)
	require.Equal(t, JsonArray, subSubV.Properties["subSubSlice"].Type)
	require.Nil(t, subSubV.Properties["subSubSlice"].ObjectValidator)

	arrV := v.Properties["arr"].ObjectValidator
	require.True(t, arrV.DisallowObject)
	require.True(t, arrV.AllowArray)
	require.Equal(t, 2, len(arrV.Properties))
	require.Equal(t, JsonString, arrV.Properties["arrFoo"].Type)
	require.Nil(t, arrV.Properties["arrFoo"].ObjectValidator)
	require.Equal(t, JsonObject, arrV.Properties["arrSub"].Type)
	require.NotNil(t, arrV.Properties["arrSub"].ObjectValidator)

	arrSubV := arrV.Properties["arrSub"].ObjectValidator
	require.Equal(t, 5, len(arrSubV.Properties))
	require.Equal(t, JsonString, arrSubV.Properties["arrSubFoo"].Type)
	require.Nil(t, arrSubV.Properties["arrSubFoo"].ObjectValidator)
	require.Equal(t, JsonInteger, arrSubV.Properties["arrSubBar"].Type)
	require.Nil(t, arrSubV.Properties["arrSubFoo"].ObjectValidator)
	require.Equal(t, JsonObject, arrSubV.Properties["arrSubSub"].Type)
	require.NotNil(t, arrSubV.Properties["arrSubSub"].ObjectValidator)
	require.Equal(t, 1, len(arrSubV.Properties["arrSubSub"].ObjectValidator.Properties))
	require.Equal(t, JsonString, arrSubV.Properties["arrSubSub"].ObjectValidator.Properties["arrSubSubFoo"].Type)
	require.Equal(t, JsonArray, arrSubV.Properties["arrSubArr"].Type)
	require.NotNil(t, arrSubV.Properties["arrSubArr"].ObjectValidator)
	require.Equal(t, 1, len(arrSubV.Properties["arrSubArr"].ObjectValidator.Properties))
	require.Equal(t, JsonString, arrSubV.Properties["arrSubArr"].ObjectValidator.Properties["arrSubArrFoo"].Type)
	require.Equal(t, JsonArray, arrSubV.Properties["arrSubSlice"].Type)
	require.Nil(t, arrSubV.Properties["arrSubSlice"].ObjectValidator)
}

type itemInSlice struct {
	Foo string `json:"foo" v8n:"notNull,mandatory"`
	Bar int    `json:"bar" v8n:"notNull,mandatory"`
}

var itemsValidator = MustCompileValidatorFor(itemInSlice{}, &ValidatorForOptions{AllowArray: true})

func TestValidatorForWithSlice(t *testing.T) {
	s := []itemInSlice{}
	json := `[
			{
				"foo": null
			},
			{
				"bar": null
			}
		]`

	ok, violations, _ := itemsValidator.ValidateStringInto(json, s)
	require.False(t, ok)
	require.Equal(t, 4, len(violations))

	SortViolationsByPathAndProperty(violations)
	require.Equal(t, fmt.Sprintf(msgMissingProperty, "bar"), violations[0].Message)
	require.Equal(t, "[0]", violations[0].Property)
	require.Equal(t, "", violations[0].Path)
	require.Equal(t, fmt.Sprintf(msgMissingProperty, "foo"), violations[1].Message)
	require.Equal(t, "[1]", violations[1].Property)
	require.Equal(t, "", violations[1].Path)
	require.Equal(t, messageValueCannotBeNull, violations[2].Message)
	require.Equal(t, "foo", violations[2].Property)
	require.Equal(t, "[0]", violations[2].Path)
	require.Equal(t, messageValueCannotBeNull, violations[3].Message)
	require.Equal(t, "bar", violations[3].Property)
	require.Equal(t, "[1]", violations[3].Path)
}

func TestNamedConstraintTagsUseCorrectDefaultFields(t *testing.T) {
	registry.reset()
	namedTestConstraint1 := &StringNotEmpty{Message: "Message 1"}
	namedTestConstraint2 := &StringNotEmpty{Message: "Message 2"}
	registry.registerNamed(true, "StringNotEmpty1", namedTestConstraint1)
	registry.registerNamed(true, "StringNotEmpty2", namedTestConstraint2)
	type MyStruct struct {
		Foo string `json:"foo" v8n:"&StringNotEmpty1{},&StringNotEmpty2{},&StringNotEmpty2{Message: 'Message 3'}"`
	}
	v, err := ValidatorFor(MyStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	str := `{"foo": ""}`
	ok, violations, _ := v.ValidateString(str)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
	require.Equal(t, "Message 1", violations[0].Message)
	require.Equal(t, "Message 2", violations[1].Message)
	require.Equal(t, "Message 3", violations[2].Message)
}

func TestConstraintTagUsingCustomConstraint(t *testing.T) {
	registry.reset()
	customConstraint := NewCustomConstraint(func(value interface{}, vcx *ValidatorContext, this *CustomConstraint) (bool, string) {
		return false, this.GetMessage()
	}, "My Custom Error")
	registry.registerNamed(true, "MyCustom", customConstraint)
	type myStruct struct {
		Foo string `json:"foo" v8n:"&MyCustom{Message: 'Overridden Message'}"`
	}
	v, err := ValidatorFor(myStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	str := `{"foo": ""}`
	ok, violations, _ := v.ValidateString(str)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Overridden Message", violations[0].Message)
}

type definedConstraintWithUnexportedField struct {
	msg       string
	TestField string
}

func (d *definedConstraintWithUnexportedField) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return false, d.GetMessage()
}
func (d *definedConstraintWithUnexportedField) GetMessage() string {
	return d.msg
}

func TestConstraintTagWithDefinedConstraintWithUnexportedField(t *testing.T) {
	registry.reset()
	registry.registerNamed(true, "MyDefined", &definedConstraintWithUnexportedField{msg: "TEST MESSAGE"})
	type myStruct struct {
		Foo string `json:"foo" v8n:"&MyDefined{}"`
	}
	v, err := ValidatorFor(myStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	type myStruct2 struct {
		Foo string `json:"foo" v8n:"&MyDefined{TestField: 'foo'}"`
	}
	v, err = ValidatorFor(myStruct2{}, nil)
	require.Nil(t, v)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldNotExported, "MyDefined", "msg"), err.Error())
}

func TestValidatorForStopsOnFirst(t *testing.T) {
	v, err := ValidatorFor(itemInSlice{}, &ValidatorForOptions{StopOnFirst: true})
	require.Nil(t, err)
	require.NotNil(t, v)
	json := `{
				"foo": null,
				"unknown": true
			}`

	s := itemInSlice{}
	ok, violations, _ := v.ValidateStringInto(json, s)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))

	v.StopOnFirst = false
	ok, violations, _ = v.ValidateStringInto(json, s)
	require.False(t, ok)
	require.Equal(t, 3, len(violations))
}

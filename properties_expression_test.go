package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestMustParseExpressionNoPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic")
		}
	}()

	_ = MustParseExpression("foo")
}

func TestMustParseExpressionPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	_ = MustParseExpression(")(")
}

func TestParseExpressions(t *testing.T) {
	testCases := map[string]struct {
		expectParse    bool
		countOrErrPosn int
	}{
		"(foo&&'bar') || (\"foo\"&& baz) || (bar && baz) && !(foo && bar && baz)": {
			true,
			4,
		},
		"((foo && bar) || (foo && baz) || (bar && baz)) && !(foo && bar && baz)": {
			true,
			2,
		},
		"": {
			true,
			0,
		},
		" ": {
			true,
			0,
		},
		".": {
			true,
			1,
		},
		"/.foo": {
			true,
			1,
		},
		"~condition": {
			true,
			1,
		},
		"?": {
			false,
			0,
		},
		"_": {
			true,
			1,
		},
		"\u00e9": {
			false,
			0,
		},
		"\u000c": {
			false,
			0,
		},
		"\t \n": {
			true,
			0,
		},
		")": {
			false,
			0,
		},
		"(foo": {
			false,
			3,
		},
		"!": {
			false,
			0,
		},
		"! foo": { // space between ! and name
			false,
			1,
		},
		"&": {
			false,
			0,
		},
		"&&": {
			false,
			0,
		},
		"|": {
			false,
			0,
		},
		"||": {
			false,
			0,
		},
		"^": {
			false,
			0,
		},
		"^^": {
			false,
			0,
		},
		"&& foo": {
			false,
			0,
		},
		"!&& foo": {
			false,
			1,
		},
		"|| foo": {
			false,
			0,
		},
		"foo bar": {
			false,
			4,
		},
		"'foo": {
			false,
			0,
		},
		" foo": {
			true,
			1,
		},
		"\tfoo": {
			true,
			1,
		},
		" foo ": {
			true,
			1,
		},
		"\tfoo\t": {
			true,
			1,
		},
		"foo &&": {
			false,
			5,
		},
		"foo && ": {
			false,
			6,
		},
		"foo !": {
			false,
			4,
		},
		"foo ! ": {
			false,
			5,
		},
		"(foo) !bar": {
			false,
			6,
		},
	}
	for str, tc := range testCases {
		t.Run(fmt.Sprintf("ParseExpression:\"%s\"", str), func(t *testing.T) {
			expr, err := ParseExpression(str)
			if tc.expectParse {
				require.Nil(t, err)
				require.Equal(t, tc.countOrErrPosn, len(expr))
			} else {
				require.NotNil(t, err)
				require.True(t, strings.Contains(err.Error(), fmt.Sprintf("at position %d)", tc.countOrErrPosn)), err.Error())
			}
		})
	}
}

func TestParseExpression(t *testing.T) {
	expStr := "(foo&&'bar') || (\"foo\"&& baz) || (bar && baz) && !(foo && bar && baz)"
	expr, err := ParseExpression(expStr)
	require.Nil(t, err)
	require.Equal(t, 4, len(expr))

	obj := jsonObject(`{
		"foo": "here"
	}`)
	result := expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	obj["bar"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	obj["baz"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	delete(obj, "foo")
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)
}

func TestEmptyExpressionAlwaysResolvesTrue(t *testing.T) {
	expr, err := ParseExpression("")
	require.Nil(t, err)

	result := expr.Evaluate(map[string]interface{}{}, nil, nil)
	require.True(t, result)
}

func TestAndOperator(t *testing.T) {
	expr, err := ParseExpression("foo && bar")
	require.Nil(t, err)

	obj := map[string]interface{}{}
	result := expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	obj["foo"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	obj["bar"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)
}

func TestOrOperator(t *testing.T) {
	expr, err := ParseExpression("foo || bar")
	require.Nil(t, err)

	obj := map[string]interface{}{}
	result := expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	obj["foo"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	obj["bar"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	delete(obj, "foo")
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	delete(obj, "bar")
	result = expr.Evaluate(obj, nil, nil)
	require.False(t, result)
}

func TestXorOperator(t *testing.T) {
	expr, err := ParseExpression("foo ^^ bar")
	require.Nil(t, err)

	obj := map[string]interface{}{}
	result := expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	obj["foo"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	obj["bar"] = "here"
	result = expr.Evaluate(obj, nil, nil)
	require.False(t, result)

	delete(obj, "foo")
	result = expr.Evaluate(obj, nil, nil)
	require.True(t, result)

	delete(obj, "bar")
	result = expr.Evaluate(obj, nil, nil)
	require.False(t, result)
}

func TestOthersOpIsAlwaysAnd(t *testing.T) {
	others := OthersExpr{}
	op := others.GetOperator()
	require.Equal(t, And, op)
}

func TestCanWalkDownProperties(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonObject,
			},
			"test": {
				Type: JsonAny,
				Constraints: Constraints{
					// check that the property gets hit...
					&FailingConstraint{Message: "Expected failure"},
				},
				RequiredWith: MustParseExpression("'foo.bar.baz'"),
			},
		},
	}
	obj := jsonObject(`{
		"foo": {
			"bar": {
				"baz": "present"
			}
		}
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPropertyRequiredWhen, violations[0].Message)
	require.Equal(t, "test", violations[0].Property)
	require.Equal(t, CodePropertyRequiredWhen, violations[0].Codes[0])

	obj = jsonObject(`{
		"foo": {
			"bar": null
		}
	}`)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// and just make sure we can hit the property...
	obj = jsonObject(`{
		"foo": {
			"bar": null
		},
		"test": "here"
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Expected failure", violations[0].Message)
}

func TestCanWalkUpProperties(t *testing.T) {
	validator := &Validator{
		Properties: Properties{
			"foo": {
				Type: JsonObject,
				ObjectValidator: &Validator{
					IgnoreUnknownProperties: true,
					Properties: Properties{
						"bar": {
							ObjectValidator: &Validator{
								Properties: Properties{
									"baz": {
										Constraints: Constraints{
											// check that the property gets hit...
											&FailingConstraint{Message: "Expected failure"},
										},
										RequiredWith: MustParseExpression("'..other_bar'"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": {
			"bar": {},
			"other_bar": "here"
		}
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPropertyRequiredWhen, violations[0].Message)
	require.Equal(t, "baz", violations[0].Property)
	require.Equal(t, CodePropertyRequiredWhen, violations[0].Codes[0])

	obj = jsonObject(`{
		"foo": {
			"bar": {}
		}
	}`)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// and just make sure we can hit the property...
	obj = jsonObject(`{
		"foo": {
			"bar": {
				"baz": "present"
			}
		}
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Expected failure", violations[0].Message)
}

func TestCanWalkFromRootProperties(t *testing.T) {
	validator := &Validator{
		IgnoreUnknownProperties: true,
		Properties: Properties{
			"foo": {
				Type: JsonObject,
				ObjectValidator: &Validator{
					Properties: Properties{
						"bar": {
							ObjectValidator: &Validator{
								Properties: Properties{
									"baz": {
										Constraints: Constraints{
											// check that the property gets hit...
											&FailingConstraint{Message: "Expected failure"},
										},
										RequiredWith: MustParseExpression("'/.other_foo.other_bar'"),
									},
								},
							},
						},
					},
				},
			},
		},
	}
	obj := jsonObject(`{
		"foo": {
			"bar": {}
		},
		"other_foo": {
			"other_bar": "present"
		}
	}`)

	ok, violations := validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, msgPropertyRequiredWhen, violations[0].Message)
	require.Equal(t, "baz", violations[0].Property)
	require.Equal(t, CodePropertyRequiredWhen, violations[0].Codes[0])

	obj = jsonObject(`{
		"foo": {
			"bar": {}
		},
		"other_foo": []
	}`)
	ok, violations = validator.Validate(obj)
	require.True(t, ok)
	require.Equal(t, 0, len(violations))

	// and just make sure we can hit the property...
	obj = jsonObject(`{
		"foo": {
			"bar": {
				"baz": "present"
			}
		}
	}`)
	ok, violations = validator.Validate(obj)
	require.False(t, ok)
	require.Equal(t, 1, len(violations))
	require.Equal(t, "Expected failure", violations[0].Message)
}

func TestOthersBuilding(t *testing.T) {
	others := OthersExpr{}
	require.Equal(t, 0, len(others))

	others.AddAndProperty("foo")
	require.Equal(t, 1, len(others))

	others.AddAndNotProperty("foo")
	require.Equal(t, 2, len(others))
}

func TestOthersFluentBuildingToStringAndParseBack(t *testing.T) {
	testCases := []struct {
		expr Other
		str  string
	}{
		{
			&OthersExpr{NewOtherProperty("foo")},
			"foo",
		},
		{
			&OthersExpr{NewOtherProperty("foo").NOTed()},
			"!foo",
		},
		{
			&OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar")},
			"foo && bar",
		},
		{
			&OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar").NOTed()},
			"foo && !bar",
		},
		{
			&OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar").ORed()},
			"foo || bar",
		},
		{
			&OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar").ORed().ANDed().NOTed().NOTed()},
			"foo && bar",
		},
		{
			(&OthersExpr{}).AddAndProperty("foo"),
			"foo",
		},
		{
			(&OthersExpr{}).AddNotProperty("foo"),
			"!foo",
		},
		{
			(&OthersExpr{}).AddAndProperty("foo./.bar"),
			"foo./.bar",
		},
		{
			(&OthersExpr{}).AddAndProperty("~condition"),
			"~condition",
		},
		{
			(&OthersExpr{}).AddAndProperty("foo\"bar"),
			"'foo\"bar'",
		},
		{
			(&OthersExpr{}).AddProperty("foo").AddOrNotProperty("bar"),
			"foo || !bar",
		},
		{
			(&OthersExpr{}).AddGroup((&OthersExpr{}).AddProperty("foo").AddOrProperty("bar")),
			"(foo || bar)",
		},
		{
			(&OthersExpr{}).AddGroup((&OthersExpr{}).AddProperty("foo").AddXorProperty("bar")),
			"(foo ^^ bar)",
		},
		{
			(&OthersExpr{}).AddGroup((&OthersExpr{}).AddProperty("foo").AddXorNotProperty("bar")),
			"(foo ^^ !bar)",
		},
		{
			(&OthersExpr{}).AddGroup(&OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar").ORed()}),
			"(foo || bar)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddOrGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) || (bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddXorGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) ^^ (bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddXorNotGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) ^^ !(bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo", "bar"),
				NewOtherGrouping("bar", "baz").ORed(),
			},
			"(foo && bar) || (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo", "bar"),
				NewOtherGrouping("bar", "baz").XORed(),
			},
			"(foo && bar) ^^ (bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddAndGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) && (bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddNotGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) && !(bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddAndNotGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) && !(bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo").AddAndProperty("bar")).
				AddOrNotGroup((&OthersExpr{}).AddProperty("bar").AddAndProperty("baz")),
			"(foo && bar) || !(bar && baz)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo")),
			"(foo)",
		},
		{
			(&OthersExpr{}).
				AddGroup((&OthersExpr{}).AddProperty("foo")),
			"(foo)",
		},
		{
			NewOtherGrouping("foo", "bar", "baz"),
			"(foo && bar && baz)",
		},
		{
			NewOtherGrouping("foo", Or, "bar", "baz"),
			"(foo || bar && baz)",
		},
		{
			NewOtherGrouping(NewOtherGrouping("foo", Or, "bar"), NewOtherGrouping(Or, "foo", "baz")),
			"((foo || bar) || (foo && baz))",
		},
		{
			NewOtherGrouping(NewOtherProperty("foo"), NewOtherProperty("bar").NOTed()),
			"(foo && !bar)",
		},
		{
			NewOtherGrouping(OtherProperty{Name: "foo"}, &OtherProperty{Name: "bar", Not: true}),
			"(foo && !bar)",
		},
		{
			NewOtherGrouping(NewOtherGrouping(NewOtherProperty("foo"), NewOtherProperty("bar").NOTed())),
			"((foo && !bar))",
		},
		{
			NewOtherGrouping(OtherGrouping{Of: OthersExpr{NewOtherProperty("foo"), NewOtherProperty("bar").NOTed()}}),
			"((foo && !bar))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo", "bar").NOTed(),
				NewOtherGrouping("bar", "baz").ORed(),
			},
			"!(foo && bar) || (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo", "bar").NOTed(),
				NewOtherGrouping("bar", "baz").ORed().ANDed(),
			},
			"!(foo && bar) && (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping().NOTed().AddProperty("foo").AddOrProperty("bar"),
				NewOtherGrouping("bar", "baz").ORed(),
			},
			"!(foo || bar) || (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping().NOTed().AddProperty("foo").AddXorProperty("bar"),
				NewOtherGrouping("bar", "baz").ORed(),
			},
			"!(foo ^^ bar) || (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping().NOTed().AddProperty("foo").AddXorNotProperty("bar"),
				NewOtherGrouping("bar", "baz").ORed(),
			},
			"!(foo ^^ !bar) || (bar && baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping().NOTed().AddProperty("foo").
					AddOrProperty("bar").
					AddNotProperty("baz"),
			},
			"!(foo || bar && !baz)",
		},
		{
			&OthersExpr{
				NewOtherGrouping().NOTed().AddProperty("foo").
					AddAndProperty("bar").
					AddAndNotProperty("baz").
					AddOrNotProperty("qux"),
			},
			"!(foo && bar && !baz || !qux)",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo && (bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddNotGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo && !(bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddAndGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo && (bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddAndNotGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo && !(bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddOrGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo || (bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddOrNotGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo || !(bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddXorGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo ^^ (bar || baz))",
		},
		{
			&OthersExpr{
				NewOtherGrouping("foo").AddXorNotGroup((&OthersExpr{}).AddProperty("bar").AddOrProperty("baz")),
			},
			"(foo ^^ !(bar || baz))",
		},
		{
			&OthersExpr{NewOtherProperty("foo").XORed(), NewOtherProperty("bar").XORed()},
			"foo ^^ bar",
		},
	}
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("TestCase[%d]:\"%s\"", i+1, tc.str), func(t *testing.T) {
			str := tc.expr.String()
			require.Equal(t, tc.str, str)
			_, err := ParseExpression(str)
			require.Nil(t, err)
		})
	}
	println(testCases)
}

func TestNewOtherGroupingPanicsWithBadArgs(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		}
	}()

	_ = NewOtherGrouping(nil)
}

func TestOtherPropertyWithEmptyNameAlwaysEvaluatesTrue(t *testing.T) {
	p := &OtherProperty{Name: ""}
	require.True(t, p.Evaluate(map[string]interface{}{}, nil, nil))

	p = &OtherProperty{Name: "."}
	require.True(t, p.Evaluate(map[string]interface{}{}, nil, nil))
	require.False(t, p.Evaluate(nil, nil, nil))
}

func TestOtherPropertyPathDetection(t *testing.T) {
	p := &OtherProperty{Name: ""}
	p.checkChanged()
	require.False(t, p.pathed)
	require.Equal(t, "", p.normalizedName)

	p = &OtherProperty{Name: "."}
	p.checkChanged()
	require.False(t, p.pathed)
	require.Equal(t, "", p.normalizedName)

	p = &OtherProperty{Name: ".foo"}
	p.checkChanged()
	require.False(t, p.pathed)
	require.Equal(t, "foo", p.normalizedName)

	p = &OtherProperty{Name: "..\\.fo\\.o.bar......"}
	p.checkChanged()
	require.True(t, p.pathed)
	require.Equal(t, 2, p.upPath)
	require.Equal(t, 2, len(p.downPath))
	require.Equal(t, ".fo.o", p.downPath[0])
	require.Equal(t, "bar", p.downPath[1])

	p = &OtherProperty{Name: "/.foo"}
	p.checkChanged()
	require.True(t, p.pathed)
	require.Equal(t, -1, p.upPath)
	require.Equal(t, 1, len(p.downPath))
	require.Equal(t, "foo", p.downPath[0])

	p = &OtherProperty{Name: "/\\.foo"}
	p.checkChanged()
	require.True(t, p.pathed)
	require.Equal(t, 0, p.upPath)
	require.Equal(t, 1, len(p.downPath))
	require.Equal(t, "/.foo", p.downPath[0])

}

func TestExpressionConditionCheck(t *testing.T) {
	expr, err := ParseExpression("bar && ~TEST_CONDITION")
	require.Nil(t, err)
	require.Equal(t, 2, len(expr))

	vcx := newEmptyValidatorContext(nil)
	obj := map[string]interface{}{
		"bar": "here",
	}

	result := expr.Evaluate(obj, nil, vcx)
	require.False(t, result)

	vcx.setInitialConditions("TEST_CONDITION")
	result = expr.Evaluate(obj, nil, vcx)
	require.True(t, result)
}

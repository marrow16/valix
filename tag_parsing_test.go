package valix

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseCommasWithBadTags(t *testing.T) {
	testStrs := []string{
		"(])",
		"[)]",
		"{]}",
		"(", "[", "{", "(()", "[[]", "{{}",
		")", "]", "}",
	}
	for i, ts := range testStrs {
		t.Run(fmt.Sprintf("Bad_Tag[%d]", i), func(t *testing.T) {
			_, err := parseCommas(ts)
			require.NotNil(t, err)
		})
	}
}

func TestParseV8nTagSimple(t *testing.T) {
	tagStr := "type:string,notNull,mandatory,constraints:[StringNotEmpty{Message: \"''''Foo\"}]"
	list, err := parseCommas(tagStr)
	require.Nil(t, err)
	require.Equal(t, 4, len(list))
}

func TestArgsStringToArgs(t *testing.T) {
	const constraintName = "TEST_CONSTRAINT"
	args, err := argsStringToArgs(constraintName, "str: \"(foo)\",bool:true,int: 0 ,float: 1.1")
	require.Nil(t, err)
	require.NotNil(t, args)
	require.Equal(t, 4, len(args))
	require.Equal(t, "str", args[0].name)
	require.Equal(t, `"(foo)"`, args[0].value)
	require.Equal(t, "bool", args[1].name)
	require.Equal(t, "true", args[1].value)
	require.Equal(t, "int", args[2].name)
	require.Equal(t, "0", args[2].value)
	require.Equal(t, "float", args[3].name)
	require.Equal(t, "1.1", args[3].value)

	_, err = argsStringToArgs(constraintName, "str: (,bool:true")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnclosed, 5)), err.Error())

	_, err = argsStringToArgs(constraintName, "str: ),bool:true")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnopened, 5)), err.Error())

	args, err = argsStringToArgs(constraintName, "xxx")
	require.Nil(t, err)
	require.Equal(t, 1, len(args))
	require.Equal(t, "xxx", args[0].name)

	// arg without value...
	args, err = argsStringToArgs(constraintName, "xxx:'',yyy")
	require.Nil(t, err)
	require.Equal(t, 2, len(args))
	require.Equal(t, "xxx", args[0].name)
	require.Equal(t, "''", args[0].value)
	require.True(t, args[0].hasValue)
	require.Equal(t, "yyy", args[1].name)
	require.Equal(t, "", args[1].value)
	require.False(t, args[1].hasValue)
}

func TestPropertyValidator_AddAllConstraints(t *testing.T) {
	constraints := defaultConstraints()
	pv := &PropertyValidator{}
	for cn := range constraints {
		t.Run(fmt.Sprintf("&%s{}", cn), func(t *testing.T) {
			err := pv.addTagItem("", "", fmt.Sprintf("&%s", cn))
			require.Nil(t, err)
		})
	}
	require.Equal(t, len(constraints), len(pv.Constraints))
}

func TestPropertyValidator_ProcessTagItems(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.processTagItems("", "", []string{})
	require.Nil(t, err)

	require.False(t, pv.Mandatory)
	require.False(t, pv.NotNull)
	err = pv.processTagItems("", "", []string{tagTokenNotNull, tagTokenMandatory})
	require.Nil(t, err)
	require.True(t, pv.Mandatory)
	require.True(t, pv.NotNull)

	require.Equal(t, 0, len(pv.Constraints))
	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "StringNotEmpty{}"})
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "&StringNotEmpty{}"})
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "[StringNotEmpty{}, &StringNotEmpty{}]"})
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.Constraints))

	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "[)(]"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnopened, 0), err.Error())

	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "[,]"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, ""), err.Error())

	err = pv.processTagItems("", "", []string{"&Bad{}X"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintsFormat, "&Bad{}X"), err.Error())

	err = pv.processTagItems("", "", []string{"UNKNOWN:x"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "UNKNOWN"), err.Error())
}

func TestPropertyValidator_AddTagItem(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.NotNull)
	err := pv.addTagItem("", "", tagTokenNotNull)
	require.Nil(t, err)
	require.True(t, pv.NotNull)
	err = pv.addTagItem("", "", tagTokenNullable)
	require.Nil(t, err)
	require.False(t, pv.NotNull)

	require.False(t, pv.Mandatory)
	err = pv.addTagItem("", "", tagTokenMandatory)
	require.Nil(t, err)
	require.True(t, pv.Mandatory)

	pv.Mandatory = false
	require.Equal(t, 0, len(pv.MandatoryWhen))
	err = pv.addTagItem("", "", tagTokenMandatory+":FOO")
	require.Nil(t, err)
	require.True(t, pv.Mandatory)
	require.Equal(t, 1, len(pv.MandatoryWhen))

	err = pv.addTagItem("", "", tagTokenMandatory+":[BAR,BAZ]")
	require.Nil(t, err)
	require.True(t, pv.Mandatory)
	require.Equal(t, 3, len(pv.MandatoryWhen))

	err = pv.addTagItem("", "", tagTokenOptional)
	require.Nil(t, err)
	require.False(t, pv.Mandatory)
	err = pv.addTagItem("", "", tagTokenRequired)
	require.Nil(t, err)
	require.True(t, pv.Mandatory)

	require.Equal(t, JsonAny, pv.Type)
	types := []string{
		jsonTypeTokenString,
		jsonTypeTokenNumber,
		jsonTypeTokenInteger,
		jsonTypeTokenBoolean,
		jsonTypeTokenObject,
		jsonTypeTokenArray,
		jsonTypeTokenAny,
	}
	for _, ty := range types {
		t.Run(fmt.Sprintf("AddTagItem-Type-%s", ty), func(t *testing.T) {
			err = pv.addTagItem("", "", tagTokenType+":"+ty)
			require.Nil(t, err)
			require.Equal(t, pv.Type.String(), ty)
		})
	}

	err = pv.addTagItem("", "", tagTokenType+": BAD")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownPropertyType, "BAD"), err.Error())
	err = pv.addTagItem("", "", tagTokenType+": "+jsonTypeTokenString+"  ") //extra spaces
	require.Nil(t, err)
	require.Equal(t, JsonString, pv.Type)

	err = pv.addTagItem("", "", "UNKNOWN:xxx")
	require.NotNil(t, err)
}

func TestPropertyValidator_AddTagItemConstraintWithNoCurly(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addTagItem("", "", "&StringNotEmpty")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
}

func TestPropertyValidator_AddTagItemConstraintWithSingleValue(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addTagItem("", "", "&StringMinLength{1}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c := pv.Constraints[0].(*StringMinLength)
	require.Equal(t, 1, c.Value)

	// constraints with tagged default field...
	err = pv.addTagItem("", "", "&StringNotBlank{'test message'}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))
	c2 := pv.Constraints[1].(*StringNotBlank)
	require.Equal(t, "test message", c2.Message)

	err = pv.addTagItem("", "", "&StringPattern{'^([A-Z]{3})$'}")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.Constraints))
	c3 := pv.Constraints[2].(*StringPattern)
	require.Equal(t, "^([A-Z]{3})$", c3.Regexp.String())

	err = pv.addTagItem("", "", "&StringValidToken{['FOO','BAR']}")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.Constraints))
	c4 := pv.Constraints[3].(*StringValidToken)
	require.Equal(t, 2, len(c4.Tokens))
	require.Equal(t, "FOO", c4.Tokens[0])

	// constraint with one field...
	constraintsRegistry.register(true, &constraintWithOneField{})
	defer func() {
		constraintsRegistry.reset()
	}()
	err = pv.addTagItem("", "", "&constraintWithOneField{'FOO'}")
	require.Nil(t, err)
	require.Equal(t, 5, len(pv.Constraints))
	c5 := pv.Constraints[4].(*constraintWithOneField)
	require.Equal(t, "FOO", c5.Message)

	// and a constraint that doesn't have a default or 'Value' field...
	err = pv.addTagItem("", "", "&StringLength{1}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldUnknown, "StringLength", "1"), err.Error())
}

type constraintWithOneField struct {
	Message string
}

func (c *constraintWithOneField) Check(v interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (c *constraintWithOneField) GetMessage(tcx I18nContext) string {
	return ""
}

func TestPropertyValidator_AddTagItemsIgnoresEmpties(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.processTagItems("", "", []string{"", "optional"})
	require.Nil(t, err)
	err = pv.processTagItems("", "", []string{"", "optional", ""})
	require.Nil(t, err)
}

func TestPropertyValidator_AddTagItemConstraint(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem("", "", tagTokenConstraint+":StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))

	err = pv.addTagItem("", "", tagTokenConstraint+":UNKNOWN{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "UNKNOWN"), err.Error())
	require.Equal(t, 1, len(pv.Constraints))

	err = pv.addTagItem("", "", "&StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.addTagItem("", "", "&UNKNOWN{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "UNKNOWN"), err.Error())
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.addTagItem("", "", "&StringValidToken{Tokens:['XXX',\"YYY\",'ZZZ']}")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.Constraints))
	c := pv.Constraints[2].(*StringValidToken)
	require.Equal(t, 3, len(c.Tokens))
	require.Equal(t, "XXX", c.Tokens[0])
	require.Equal(t, "YYY", c.Tokens[1])
	require.Equal(t, "ZZZ", c.Tokens[2])
}

func TestPropertyValidator_AddTagItemWhen(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.WhenConditions))

	err := pv.addTagItem("", "", tagTokenWhen+":TEST1")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.WhenConditions))
	require.Equal(t, "TEST1", pv.WhenConditions[0])

	err = pv.addTagItem("", "", tagTokenWhen+":[TEST2,TEST3]")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.WhenConditions))
	require.Equal(t, "TEST2", pv.WhenConditions[1])
	require.Equal(t, "TEST3", pv.WhenConditions[2])

	err = pv.addTagItem("", "", tagTokenWhen)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgExpectedColon, tagTokenWhen), err.Error())

	err = pv.addTagItem("", "", tagTokenWhen+":[{]")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 0), err.Error())

	err = pv.addTagItem("", "", tagTokenWhen+":\"TEST4\"")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.WhenConditions))
	require.Equal(t, "TEST4", pv.WhenConditions[3])

	err = pv.addTagItem("", "", tagTokenWhen+":['TEST5','TEST6']")
	require.Nil(t, err)
	require.Equal(t, 6, len(pv.WhenConditions))
	require.Equal(t, "TEST5", pv.WhenConditions[4])
	require.Equal(t, "TEST6", pv.WhenConditions[5])
}

func TestPropertyValidator_AddTagItemUnwanted(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.UnwantedConditions))

	err := pv.addTagItem("", "", tagTokenUnwanted+":TEST1")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.UnwantedConditions))
	require.Equal(t, "TEST1", pv.UnwantedConditions[0])

	err = pv.addTagItem("", "", tagTokenUnwanted+":[TEST2,TEST3]")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.UnwantedConditions))
	require.Equal(t, "TEST2", pv.UnwantedConditions[1])
	require.Equal(t, "TEST3", pv.UnwantedConditions[2])

	err = pv.addTagItem("", "", tagTokenUnwanted)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgExpectedColon, tagTokenUnwanted), err.Error())

	err = pv.addTagItem("", "", tagTokenUnwanted+":[{]")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 0), err.Error())

	err = pv.addTagItem("", "", tagTokenUnwanted+":\"TEST4\"")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.UnwantedConditions))
	require.Equal(t, "TEST4", pv.UnwantedConditions[3])

	err = pv.addTagItem("", "", tagTokenUnwanted+":['TEST5','TEST6']")
	require.Nil(t, err)
	require.Equal(t, 6, len(pv.UnwantedConditions))
	require.Equal(t, "TEST5", pv.UnwantedConditions[4])
	require.Equal(t, "TEST6", pv.UnwantedConditions[5])
}

func TestPropertyValidator_AddTagItemMandatoryWhen(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.MandatoryWhen))

	err := pv.addTagItem("", "", tagTokenMandatory+":TEST1")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.MandatoryWhen))
	require.Equal(t, "TEST1", pv.MandatoryWhen[0])

	err = pv.addTagItem("", "", tagTokenMandatory+":[TEST2,TEST3]")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.MandatoryWhen))
	require.Equal(t, "TEST2", pv.MandatoryWhen[1])
	require.Equal(t, "TEST3", pv.MandatoryWhen[2])

	err = pv.addTagItem("", "", tagTokenMandatory+":[{]")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 0), err.Error())

	err = pv.addTagItem("", "", tagTokenMandatory+":\"TEST4\"")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.MandatoryWhen))
	require.Equal(t, "TEST4", pv.MandatoryWhen[3])

	err = pv.addTagItem("", "", tagTokenRequired+":['TEST5','TEST6']")
	require.Nil(t, err)
	require.Equal(t, 6, len(pv.MandatoryWhen))
	require.Equal(t, "TEST5", pv.MandatoryWhen[4])
	require.Equal(t, "TEST6", pv.MandatoryWhen[5])
}

func TestPropertyValidator_AddTagItemOnly(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.Only)
	require.Equal(t, 0, len(pv.OnlyConditions))

	err := pv.addTagItem("", "", tagTokenOnly)
	require.Nil(t, err)
	require.True(t, pv.Only)
	require.Equal(t, 0, len(pv.OnlyConditions))

	err = pv.addTagItem("", "", tagTokenOnly+":TEST1")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.OnlyConditions))
	require.Equal(t, "TEST1", pv.OnlyConditions[0])

	err = pv.addTagItem("", "", tagTokenOnly+":[TEST2,TEST3]")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.OnlyConditions))
	require.Equal(t, "TEST2", pv.OnlyConditions[1])
	require.Equal(t, "TEST3", pv.OnlyConditions[2])

	err = pv.addTagItem("", "", tagTokenOnly+":[{]")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 0), err.Error())

	err = pv.addTagItem("", "", tagTokenOnly+":\"TEST4\"")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.OnlyConditions))
	require.Equal(t, "TEST4", pv.OnlyConditions[3])

	err = pv.addTagItem("", "", tagTokenOnly+":['TEST5','TEST6']")
	require.Nil(t, err)
	require.Equal(t, 6, len(pv.OnlyConditions))
	require.Equal(t, "TEST5", pv.OnlyConditions[4])
	require.Equal(t, "TEST6", pv.OnlyConditions[5])

	err = pv.addTagItem("", "", tagTokenOnlyMsg)
	require.NotNil(t, err)
	err = pv.addTagItem("", "", tagTokenOnlyMsg+":foo bar")
	require.Nil(t, err)
	require.Equal(t, "foo bar", pv.OnlyMessage)
	err = pv.addTagItem("", "", tagTokenOnlyMsg+":'foo bar'")
	require.Nil(t, err)
	require.Equal(t, "foo bar", pv.OnlyMessage)
}

func TestPropertyValidator_AddTagItemStopOnFirst(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.StopOnFirst)

	err := pv.addTagItem("", "", tagTokenStopOnFirst)
	require.Nil(t, err)
	require.True(t, pv.StopOnFirst)

	pv.StopOnFirst = false
	err = pv.addTagItem("", "", tagTokenStopOnFirstAlt)
	require.Nil(t, err)
	require.True(t, pv.StopOnFirst)

	err = pv.addTagItem("", "", tagTokenStopOnFirst+":foo")
	require.NotNil(t, err)
}

func TestPropertyValidator_AddConditionalConstraint(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem("", "", "&[FOO]StringNotEmpty{'Foo Message'}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	constraint := pv.Constraints[0].(*ConditionalConstraint)
	require.Equal(t, 1, len(constraint.When))
	require.Equal(t, "FOO", constraint.When[0])
	innerConstraint, ok := constraint.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "Foo Message", innerConstraint.Message)

	err = pv.addTagItem("", "", "&[BAR,BAZ]StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))
	constraint = pv.Constraints[1].(*ConditionalConstraint)
	require.Equal(t, 2, len(constraint.When))
	require.Equal(t, "BAR", constraint.When[0])
	require.Equal(t, "BAZ", constraint.When[1])
	innerConstraint, ok = constraint.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "", innerConstraint.Message)

	err = pv.addTagItem("", "", "&[}]StringNotEmpty{}")
	require.NotNil(t, err)

	err = pv.addTagItem("", "", "&[StringNotEmpty{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConditionalConstraintsFormat, "&[StringNotEmpty{}"), err.Error())

	err = pv.processTagItems("", "", []string{tagTokenConstraintsPrefix + "[[foo2]StringNotEmpty,[bar2,'baz2']StringNotBlank{'Not Blank'}]"})
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.Constraints))
	constraint = pv.Constraints[2].(*ConditionalConstraint)
	require.Equal(t, 1, len(constraint.When))
	require.Equal(t, "foo2", constraint.When[0])
	innerConstraint, ok = constraint.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "", innerConstraint.Message)
	constraint = pv.Constraints[3].(*ConditionalConstraint)
	require.Equal(t, 2, len(constraint.When))
	require.Equal(t, "bar2", constraint.When[0])
	require.Equal(t, "baz2", constraint.When[1])
	innerConstraint2, ok := constraint.Constraint.(*StringNotBlank)
	require.True(t, ok)
	require.Equal(t, "Not Blank", innerConstraint2.Message)
}

func TestPropertyValidator_AddConditionalConstraintWithExpr(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem("", "", "&<bar && baz>StringNotEmpty{'Foo Message'}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	constraint := pv.Constraints[0].(*ConditionalConstraint)
	require.Equal(t, 0, len(constraint.When))
	require.NotNil(t, constraint.Others)
	require.Equal(t, "bar && baz", constraint.Others.String())
	innerConstraint, ok := constraint.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.Equal(t, "Foo Message", innerConstraint.Message)

	err = pv.addTagItem("", "", "&<not a valid expr>StringNotEmpty{'Foo Message'}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConditionalExpr, "not a valid expr", "unexpected property name start (at position 4)"), err.Error())

	err = pv.addTagItem("", "", "&<no closing StringNotEmpty{'Foo Message'}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConditionalConstraintsFormat, "&<no closing StringNotEmpty{'Foo Message'}"), err.Error())
}

func TestPropertyValidator_AddSetConditionIf(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem("", "", "&SetConditionIf{Constraint:&StringNotEmpty, SetOk:'OK', SetFail:'FAIL'}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c0, ok := pv.Constraints[0].(*SetConditionIf)
	require.True(t, ok)
	require.Equal(t, "OK", c0.SetOk)
	require.Equal(t, "FAIL", c0.SetFail)
	_, ok = c0.Constraint.(*StringNotEmpty)
	require.True(t, ok)
}

func TestPropertyValidator_AddObjectTagItem_RequiredWith(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.RequiredWith))

	err := pv.addTagItem("", "", tagTokenRequiredWithAlt+":!foo")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.RequiredWith))
	otherPty := pv.RequiredWith[0].(*OtherProperty)
	require.Equal(t, "foo", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenRequiredWithAlt+":!bar")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.RequiredWith))
	otherPty = pv.RequiredWith[1].(*OtherProperty)
	require.Equal(t, "bar", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenRequiredWithAlt+":(!bar || !baz)")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.RequiredWith))
	otherGrp := pv.RequiredWith[2].(*OtherGrouping)
	require.Equal(t, 2, len(otherGrp.Of))
	otherPty = otherGrp.Of[0].(*OtherProperty)
	require.Equal(t, "bar", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)
	otherPty = otherGrp.Of[1].(*OtherProperty)
	require.Equal(t, "baz", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, Or, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenRequiredWithAlt+":bad expression")
	require.NotNil(t, err)
}

func TestPropertyValidator_AddObjectTagItem_RequiredWithMsg(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, "", pv.RequiredWithMessage)

	err := pv.addTagItem("", "", tagTokenRequiredWithAltMsg+":fooey")
	require.Nil(t, err)
	require.Equal(t, "fooey", pv.RequiredWithMessage)

	err = pv.addTagItem("", "", tagTokenRequiredWithAltMsg+":'fooey'")
	require.Nil(t, err)
	require.Equal(t, "fooey", pv.RequiredWithMessage)
}

func TestPropertyValidator_AddObjectTagItem_UnwantedWith(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.UnwantedWith))

	err := pv.addTagItem("", "", tagTokenUnwantedWithAlt+":!foo")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.UnwantedWith))
	otherPty := pv.UnwantedWith[0].(*OtherProperty)
	require.Equal(t, "foo", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenUnwantedWithAlt+":!bar")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.UnwantedWith))
	otherPty = pv.UnwantedWith[1].(*OtherProperty)
	require.Equal(t, "bar", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenUnwantedWithAlt+":(!bar || !baz)")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.UnwantedWith))
	otherGrp := pv.UnwantedWith[2].(*OtherGrouping)
	require.Equal(t, 2, len(otherGrp.Of))
	otherPty = otherGrp.Of[0].(*OtherProperty)
	require.Equal(t, "bar", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, And, otherPty.Op)
	otherPty = otherGrp.Of[1].(*OtherProperty)
	require.Equal(t, "baz", otherPty.Name)
	require.True(t, otherPty.Not)
	require.Equal(t, Or, otherPty.Op)

	err = pv.addTagItem("", "", tagTokenUnwantedWithAlt+":bad expression")
	require.NotNil(t, err)
}

func TestPropertyValidator_AddObjectTagItem_UnwantedWithMsg(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, "", pv.UnwantedWithMessage)

	err := pv.addTagItem("", "", tagTokenUnwantedWithAltMsg+":fooey")
	require.Nil(t, err)
	require.Equal(t, "fooey", pv.UnwantedWithMessage)

	err = pv.addTagItem("", "", tagTokenUnwantedWithAltMsg+":'fooey'")
	require.Nil(t, err)
	require.Equal(t, "fooey", pv.UnwantedWithMessage)
}

func TestPropertyValidator_AddTagItemObjWhen(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.Equal(t, 0, len(pv.ObjectValidator.WhenConditions))

	err := pv.addTagItem("", "", tagTokenObjWhen+":TEST1")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.ObjectValidator.WhenConditions))
	require.Equal(t, "TEST1", pv.ObjectValidator.WhenConditions[0])

	err = pv.addTagItem("", "", tagTokenObjWhen+":[TEST2,TEST3]")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.ObjectValidator.WhenConditions))
	require.Equal(t, "TEST2", pv.ObjectValidator.WhenConditions[1])
	require.Equal(t, "TEST3", pv.ObjectValidator.WhenConditions[2])

	err = pv.addTagItem("", "", tagTokenObjWhen)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgExpectedColon, tagTokenObjWhen), err.Error())

	err = pv.addTagItem("", "", tagTokenObjWhen+":[{]")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnclosed, 0), err.Error())

	err = pv.addTagItem("", "", tagTokenObjWhen+":\"TEST4\"")
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.ObjectValidator.WhenConditions))
	require.Equal(t, "TEST4", pv.ObjectValidator.WhenConditions[3])

	err = pv.addTagItem("", "", tagTokenObjWhen+":['TEST5','TEST6']")
	require.Nil(t, err)
	require.Equal(t, 6, len(pv.ObjectValidator.WhenConditions))
	require.Equal(t, "TEST5", pv.ObjectValidator.WhenConditions[4])
	require.Equal(t, "TEST6", pv.ObjectValidator.WhenConditions[5])

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenObjWhen+":TEST1")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenObjWhen), err.Error())
}

func TestPropertyValidator_AddObjectTagItem_IgnoreUnknownProperties(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err := pv.addTagItem("", "", tagTokenObjIgnoreUnknownProperties)
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.IgnoreUnknownProperties)

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenObjIgnoreUnknownProperties)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenObjIgnoreUnknownProperties), err.Error())
}

func TestPropertyValidator_AddObjectTagItem_UnknownProperties(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err := pv.addTagItem("", "", tagTokenObjUnknownProperties+": true")
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err = pv.addTagItem("", "", tagTokenObjUnknownProperties+": FALSE")
	require.Nil(t, err)
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err = pv.addTagItem("", "", tagTokenObjUnknownProperties+": INVALID_BOOL_VALUE")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTagValue, tagTokenObjUnknownProperties, "boolean", "INVALID_BOOL_VALUE"), err.Error())

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenObjUnknownProperties+": FALSE")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenObjUnknownProperties), err.Error())
}

func TestPropertyValidator_AddObjectTagItem_Ordered(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.OrderedPropertyChecks)

	err := pv.addTagItem("", "", tagTokenObjOrdered)
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.OrderedPropertyChecks)

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenObjOrdered)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenObjOrdered), err.Error())
}

func TestPropertyValidator_AddObjectTagItem_Constraint(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{
		Constraints: Constraints{},
	}}
	require.Equal(t, 0, len(pv.ObjectValidator.Constraints))

	err := pv.addTagItem("", "", tagTokenObjConstraint+":&StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.ObjectValidator.Constraints))

	err = pv.addTagItem("", "", tagTokenObjConstraint+":StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.ObjectValidator.Constraints))

	err = pv.addTagItem("", "", tagTokenObjConstraint+":&Unknown{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "Unknown"), err.Error())
	require.Equal(t, 2, len(pv.ObjectValidator.Constraints))

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenObjConstraint+":&WontGetUsed{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenObjConstraint), err.Error())
}

func TestPropertyValidator_AddArrayTagItem_AllowNullItems(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.AllowNullItems)

	err := pv.addTagItem("", "", tagTokenArrAllowNullItems)
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.AllowNullItems)

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", tagTokenArrAllowNullItems)
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgPropertyNotObject, tagTokenArrAllowNullItems), err.Error())
}

func TestParseConstraintSet(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{
		Constraints: Constraints{},
	}}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem("", "", "&ConstraintSet{Message:'Foo',Constraints:[&StringNotEmpty{Message:'msg1'},&StringNotBlank{Message:'msg2'}],Stop:true,OneOf:true}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	constraintSet, ok := pv.Constraints[0].(*ConstraintSet)
	require.True(t, ok)
	require.Equal(t, 2, len(constraintSet.Constraints))
	require.True(t, constraintSet.Stop)
	require.True(t, constraintSet.OneOf)
}

func TestRebuildConstraintWithArgs(t *testing.T) {
	const constraintName = "TEST_CONSTRAINT"
	var orgConstraint = &StringNotEmpty{}
	c, err := rebuildConstraintWithArgs(constraintName, orgConstraint, "")
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, msgNotEmptyString, c.GetMessage(nil))

	c, err = rebuildConstraintWithArgs(constraintName, orgConstraint, ")(")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnopened, 0)), err.Error())
}

func TestRebuildConstraintWithOneArg(t *testing.T) {
	const constraintName = "TEST_CONSTRAINT"
	var orgConstraint = &StringExactLength{Message: "foo"}

	c, err := rebuildConstraintWithArgs(constraintName, orgConstraint, "16")
	require.Nil(t, err)
	ac := c.(*StringExactLength)
	require.Equal(t, 16, ac.Value)

	_, err = rebuildConstraintWithArgs(constraintName, orgConstraint, "'foo'")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldInvalidValue, constraintName, "Value"), err.Error())

	_, err = rebuildConstraintWithArgs(constraintName, orgConstraint, "msg")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldInvalidValue, constraintName, "msg"), err.Error())

	// bools don't need a value...
	c, err = rebuildConstraintWithArgs(constraintName, orgConstraint, "stp")
	require.Nil(t, err)
	ac = c.(*StringExactLength)
	require.True(t, ac.Stop)
}

type dummyNonStructConstraint map[string]struct{}

func (d *dummyNonStructConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (d *dummyNonStructConstraint) GetMessage(tcx I18nContext) string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithNonStructConstraint(t *testing.T) {
	constraintsRegistry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyNonStructConstraint{}
	constraintsRegistry.register(true, testConstraint)
	defer func() {
		constraintsRegistry.reset()
	}()

	c, err := rebuildConstraintWithArgs(testConstraintName, testConstraint, "")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgCannotCreateConstraint, testConstraintName), err.Error())
	require.Nil(t, c)
}

type dummyStructConstraintWithUnexportedFields struct {
	field1 string
	Field2 int
}

func (d *dummyStructConstraintWithUnexportedFields) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (d *dummyStructConstraintWithUnexportedFields) GetMessage(tcx I18nContext) string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithNonPublicField(t *testing.T) {
	constraintsRegistry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyStructConstraintWithUnexportedFields{}
	constraintsRegistry.register(true, testConstraint)
	defer func() {
		constraintsRegistry.reset()
	}()

	c, err := rebuildConstraintWithArgs(testConstraintName, testConstraint, "field1:\"foo\"")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldNotExported, testConstraintName, "field1"), err.Error())
	require.Nil(t, c)
}

type dummyStructConstraintWithFields struct {
	Field1 string
	Field2 int
}

func (d *dummyStructConstraintWithFields) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (d *dummyStructConstraintWithFields) GetMessage(tcx I18nContext) string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithInvalidArgType(t *testing.T) {
	constraintsRegistry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyStructConstraintWithFields{}
	constraintsRegistry.register(true, testConstraint)
	defer func() {
		constraintsRegistry.reset()
	}()

	_, err := rebuildConstraintWithArgs(testConstraintName, testConstraint, "Field2:\"foo\"")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldInvalidValue, testConstraintName, "Field2"), err.Error())
}

func TestBuildConstraintFromTagValueFailsWithBadArgs(t *testing.T) {
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	testConstraint := &dummyStructConstraintWithFields{}
	constraintsRegistry.register(true, testConstraint)

	c, err := buildConstraintFromTagValue("&dummyStructConstraintWithFields{UnknownField:\"\"}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldUnknown, "dummyStructConstraintWithFields", "UnknownField"), err.Error())
	require.Nil(t, c)
}

func TestUnexpectedColonAfterTagToken(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addTagItem("", "", tagTokenNotNull+":")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnexpectedColon, tagTokenNotNull), err.Error())
}

type structWithTypeFields struct {
	AString string
	AInt    int
	AUInt   uint
	AFloat  float64
	ABool   bool
	ARegexp regexp.Regexp
	ASlice  []string
	AOther  Constraint
}

func getTestFieldValue(fieldName string) (reflect.Value, *structWithTypeFields) {
	testStruct := structWithTypeFields{}
	ty := reflect.TypeOf(testStruct)
	newC := reflect.New(ty)
	newS := newC.Interface().(*structWithTypeFields)
	return newC.Elem().FieldByName(fieldName), newS
}

func TestSafeSet_String(t *testing.T) {
	fld, newStruct := getTestFieldValue("AString")
	require.Equal(t, "", newStruct.AString)

	ok := safeSet(fld, "\"baz\"", true)
	require.True(t, ok)
	require.Equal(t, "baz", fld.String())
	require.Equal(t, "baz", newStruct.AString)

	ok = safeSet(fld, "0", true)
	require.False(t, ok)
}

func TestSafeSet_Int(t *testing.T) {
	fld, newStruct := getTestFieldValue("AInt")
	require.Equal(t, 0, newStruct.AInt)

	ok := safeSet(fld, "1", true)
	require.True(t, ok)
	require.Equal(t, int64(1), fld.Int())
	require.Equal(t, 1, newStruct.AInt)

	ok = safeSet(fld, "\"1\"", true)
	require.False(t, ok)
	ok = safeSet(fld, "true", true)
	require.False(t, ok)
	ok = safeSet(fld, "1.1", true)
	require.False(t, ok)
}

func TestSafeSet_UInt(t *testing.T) {
	fld, newStruct := getTestFieldValue("AUInt")
	require.Equal(t, uint(0), newStruct.AUInt)

	ok := safeSet(fld, "1", true)
	require.True(t, ok)
	require.Equal(t, uint64(1), fld.Uint())
	require.Equal(t, uint(1), newStruct.AUInt)

	ok = safeSet(fld, "\"1\"", true)
	require.False(t, ok)
	ok = safeSet(fld, "true", true)
	require.False(t, ok)
	ok = safeSet(fld, "1.1", true)
	require.False(t, ok)
}

func TestSafeSet_Float(t *testing.T) {
	fld, newStruct := getTestFieldValue("AFloat")
	require.Equal(t, float64(0), newStruct.AFloat)

	ok := safeSet(fld, "1", true)
	require.True(t, ok)
	require.Equal(t, float64(1), fld.Float())
	require.Equal(t, float64(1), newStruct.AFloat)

	ok = safeSet(fld, "\"1\"", true)
	require.False(t, ok)
	ok = safeSet(fld, "true", true)
	require.False(t, ok)
}

func TestSafeSet_Bool(t *testing.T) {
	fld, newStruct := getTestFieldValue("ABool")
	require.Equal(t, false, newStruct.ABool)

	ok := safeSet(fld, "true", true)
	require.True(t, ok)
	require.Equal(t, true, fld.Bool())
	require.Equal(t, true, newStruct.ABool)

	fld, newStruct = getTestFieldValue("ABool")
	require.Equal(t, false, newStruct.ABool)
	ok = safeSet(fld, "", false)
	require.True(t, ok)
	require.Equal(t, true, fld.Bool())
	require.Equal(t, true, newStruct.ABool)

	ok = safeSet(fld, "\"1\"", true)
	require.False(t, ok)
	ok = safeSet(fld, "1.1", true)
	require.False(t, ok)
}

func TestSafeSet_Slice(t *testing.T) {
	fld, newStruct := getTestFieldValue("ASlice")
	require.Equal(t, 0, len(newStruct.ASlice))

	ok := safeSet(fld, "[\"foo\", \"bar\"]", true)
	require.True(t, ok)
	arr := fld.Interface().([]string)
	require.Equal(t, 2, len(arr))
	require.Equal(t, 2, len(newStruct.ASlice))
	require.Equal(t, "foo", arr[0])
	require.Equal(t, "bar", arr[1])

	ok = safeSet(fld, "\"foo\"", true)
	require.False(t, ok)
	ok = safeSet(fld, "0", true)
	require.False(t, ok)
}

func TestSafeSet_Regexp(t *testing.T) {
	fld, newStruct := getTestFieldValue("ARegexp")
	require.NotNil(t, newStruct.ARegexp)

	const pattern = "^([a-fA-F0-9]{8})$"
	ok := safeSet(fld, "\""+pattern+"\"", true)
	require.True(t, ok)
	rx := fld.Interface().(regexp.Regexp)
	require.Equal(t, pattern, rx.String())
	require.Equal(t, pattern, newStruct.ARegexp.String())

	ok = safeSet(fld, "^^^", true)
	require.False(t, ok)
	ok = safeSet(fld, "1.1", true)
	require.False(t, ok)
	ok = safeSet(fld, "true", true)
	require.False(t, ok)
}

func TestSafeSet_OtherUnknown(t *testing.T) {
	fld, _ := getTestFieldValue("AOther")

	ok := safeSet(fld, "", true)
	require.False(t, ok)
}

func TestItemsToSlice_String(t *testing.T) {
	sampleArr := make([]string, 0)
	itemType := reflect.TypeOf(sampleArr)

	v, ok := itemsToSlice(itemType, "['foo','bar', 'baz']")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant := v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 3, len(resultant.([]string)))
	require.Equal(t, "foo", (resultant.([]string))[0])
	require.Equal(t, "bar", (resultant.([]string))[1])
	require.Equal(t, "baz", (resultant.([]string))[2])

	v, ok = itemsToSlice(itemType, "[ foo, bar, baz ]")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant = v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 3, len(resultant.([]string)))
	require.Equal(t, "foo", (resultant.([]string))[0])
	require.Equal(t, "bar", (resultant.([]string))[1])
	require.Equal(t, "baz", (resultant.([]string))[2])
}

func TestItemsToSlice_Int(t *testing.T) {
	sampleArr := make([]int, 0)
	itemType := reflect.TypeOf(sampleArr)

	v, ok := itemsToSlice(itemType, "[1,2, 3]")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant := v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 3, len(resultant.([]int)))
	require.Equal(t, 1, (resultant.([]int))[0])
	require.Equal(t, 2, (resultant.([]int))[1])
	require.Equal(t, 3, (resultant.([]int))[2])

	v, ok = itemsToSlice(itemType, "['foo',2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[1.1,2,3]")
	require.False(t, ok)
}

func TestItemsToSlice_UInt(t *testing.T) {
	sampleArr := make([]uint, 0)
	itemType := reflect.TypeOf(sampleArr)

	v, ok := itemsToSlice(itemType, "[1,2, 3]")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant := v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 3, len(resultant.([]uint)))
	require.Equal(t, uint(1), (resultant.([]uint))[0])
	require.Equal(t, uint(2), (resultant.([]uint))[1])
	require.Equal(t, uint(3), (resultant.([]uint))[2])

	v, ok = itemsToSlice(itemType, "['foo',2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[-1,2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[1.1,2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[x,2,3]")
	require.False(t, ok)
}

func TestItemsToSlice_Float(t *testing.T) {
	sampleArr := make([]float64, 0)
	itemType := reflect.TypeOf(sampleArr)

	v, ok := itemsToSlice(itemType, "[1.1,2.2, 3.3]")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant := v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 3, len(resultant.([]float64)))
	require.Equal(t, 1.1, (resultant.([]float64))[0])
	require.Equal(t, 2.2, (resultant.([]float64))[1])
	require.Equal(t, 3.3, (resultant.([]float64))[2])

	v, ok = itemsToSlice(itemType, "['foo',2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[x,2,3]")
	require.False(t, ok)
}

func TestItemsToSlice_Bool(t *testing.T) {
	sampleArr := make([]bool, 0)
	itemType := reflect.TypeOf(sampleArr)

	v, ok := itemsToSlice(itemType, "[1,t,T,TRUE,true,True, 0,f,F,FALSE,false,False]")
	require.True(t, ok)
	require.NotNil(t, v)
	resultant := v.Interface()
	require.NotNil(t, resultant)
	require.Equal(t, 12, len(resultant.([]bool)))
	require.True(t, (resultant.([]bool))[0])
	require.True(t, (resultant.([]bool))[1])
	require.True(t, (resultant.([]bool))[2])
	require.True(t, (resultant.([]bool))[3])
	require.True(t, (resultant.([]bool))[4])
	require.True(t, (resultant.([]bool))[5])
	require.False(t, (resultant.([]bool))[6])
	require.False(t, (resultant.([]bool))[7])
	require.False(t, (resultant.([]bool))[8])
	require.False(t, (resultant.([]bool))[9])
	require.False(t, (resultant.([]bool))[10])
	require.False(t, (resultant.([]bool))[11])

	v, ok = itemsToSlice(itemType, "['foo',2,3]")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[x,2,3]")
	require.False(t, ok)
}

func TestDelimStack(t *testing.T) {
	stk := &delimiterStack{current: nil, stack: []*delimiter{}}
	stk.push('"', 1)
	stk.push('[', 2)
	stk.push('(', 3)
	require.Equal(t, 2, len(stk.stack))
	require.Equal(t, '(', stk.current.open)
	require.Equal(t, 3, stk.current.pos)
	stk.pop()
	require.Equal(t, 1, len(stk.stack))
	require.Equal(t, '[', stk.current.open)
	require.Equal(t, 2, stk.current.pos)
	stk.pop()
	require.Equal(t, 0, len(stk.stack))
	require.Equal(t, '"', stk.current.open)
	require.Equal(t, 1, stk.current.pos)
	stk.pop()
	require.Nil(t, stk.current)
}

func TestWrappedConstraint(t *testing.T) {
	type myStruct struct {
		Foo []string `json:"foo" v8n:"&ArrayConditionalConstraint{When:'!first', Constraint:StringNotEmpty{}}"`
	}
	v, err := ValidatorFor(myStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	// and with bad constraint...
	type myStruct2 struct {
		Foo []string `json:"foo" v8n:"&ArrayConditionalConstraint{When:'!first', Constraint:unknownConstraint{}}"`
	}
	_, err = ValidatorFor(myStruct2{}, nil)
	require.NotNil(t, err)
}

func TestWrappedConstraints(t *testing.T) {
	type myStruct struct {
		Foo []string `json:"foo" v8n:"&ArrayOf{Type:'string', Constraints:[&StringNotEmpty{}]}"`
	}
	v, err := ValidatorFor(myStruct{}, nil)
	require.Nil(t, err)
	require.NotNil(t, v)

	// and with bad constraint...
	type myStruct2 struct {
		Foo []string `json:"foo" v8n:"&ArrayOf{Type:'string', Constraints:[&unknownConstraint{}]}"`
	}
	_, err = ValidatorFor(myStruct2{}, nil)
	require.NotNil(t, err)
}

func TestConditionalConstraintParse(t *testing.T) {
	pv := &PropertyValidator{}

	err := pv.addTagItem("", "", "&ConditionalConstraint{When:['FOO','BAR'],Constraint:&StringNotEmpty{Stop:true}}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	cc, ok := pv.Constraints[0].(*ConditionalConstraint)
	require.True(t, ok)
	require.Equal(t, 2, len(cc.When))
	require.Equal(t, "FOO", cc.When[0])
	require.Equal(t, "BAR", cc.When[1])
	ccc, ok := cc.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.True(t, ccc.Stop)
}

func TestConditionalConstraintParseOthers(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.addTagItem("", "", "&ConditionalConstraint{Others:foo && bar,Constraint:&StringNotEmpty{Stop:true}}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	cc, ok := pv.Constraints[0].(*ConditionalConstraint)
	require.True(t, ok)
	require.NotNil(t, cc.Others)
	require.Equal(t, "foo && bar", cc.Others.String())

	ccc, ok := cc.Constraint.(*StringNotEmpty)
	require.True(t, ok)
	require.True(t, ccc.Stop)

	pv = &PropertyValidator{}
	err = pv.addTagItem("", "", "&ConditionalConstraint{Others:'foo && bar',Constraint:&StringNotEmpty{Stop:true}}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	cc, ok = pv.Constraints[0].(*ConditionalConstraint)
	require.True(t, ok)
	require.NotNil(t, cc.Others)
	require.Equal(t, "foo && bar", cc.Others.String())

	err = pv.addTagItem("", "", "&ConditionalConstraint{Others:not a valid expr,Constraint:&StringNotEmpty{Stop:true}}")
	require.NotNil(t, err)
}

func TestConstraintFieldAbbreviation(t *testing.T) {
	pv := PropertyValidator{}
	err := pv.addTagItem("", "", "&StringValidEmail{}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))

	// 'allow' matches more than one field...
	err = pv.addTagItem("", "", "&StringValidEmail{allow:true}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldUnknown, "StringValidEmail", "allow"), err.Error())

	// case-insensitive...
	pv = PropertyValidator{}
	err = pv.addTagItem("", "", "&StringValidEmail{allowIpv6:true}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c := pv.Constraints[0].(*StringValidEmail)
	require.True(t, c.AllowIPV6)

	// only matches one...
	pv = PropertyValidator{}
	err = pv.addTagItem("", "", "&StringValidEmail{v6:true}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c = pv.Constraints[0].(*StringValidEmail)
	require.True(t, c.AllowIPV6)

	// boolean no value..
	pv = PropertyValidator{}
	err = pv.addTagItem("", "", "&StringValidEmail{v6}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c = pv.Constraints[0].(*StringValidEmail)
	require.True(t, c.AllowIPV6)

	// abbreviated...
	pv = PropertyValidator{}
	err = pv.addTagItem("", "", "&StringValidEmail{Msg:'fooey'}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	c = pv.Constraints[0].(*StringValidEmail)
	require.Equal(t, "fooey", c.Message)
}

func TestReadmeShortening(t *testing.T) {
	pv, err := NewPropertyValidator("&strnocc,&strupper{'Upper only'},&strlen{min:10,max:20,excMin}")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.Constraints))
	_, ok := pv.Constraints[0].(*StringNoControlCharacters)
	require.True(t, ok)
	c1, ok := pv.Constraints[1].(*StringUppercase)
	require.True(t, ok)
	require.Equal(t, "Upper only", c1.Message)
	c2, ok := pv.Constraints[2].(*StringLength)
	require.True(t, ok)
	require.Equal(t, 10, c2.Minimum)
	require.Equal(t, 20, c2.Maximum)
	require.True(t, c2.ExclusiveMin)
}

func TestAbbreviateName(t *testing.T) {
	abbr := abbreviateName("messsssssssage")
	require.Equal(t, "msg", abbr)

	abbr = abbreviateName("message")
	require.Equal(t, "msg", abbr)

	abbr = abbreviateName("stop")
	require.Equal(t, "stp", abbr)

	abbr = abbreviateName("allow")
	require.Equal(t, "alw", abbr)
}

func TestCamelToWords(t *testing.T) {
	w := camelToWords("eMi")
	require.Equal(t, 2, len(w))
	require.Equal(t, "e", w[0])
	require.Equal(t, "mi", w[1])

	w = camelToWords("EMI")
	require.Equal(t, 1, len(w))
	require.Equal(t, "emi", w[0])

	w = camelToWords("Emi")
	require.Equal(t, 1, len(w))
	require.Equal(t, "emi", w[0])

	w = camelToWords("EmI")
	require.Equal(t, 2, len(w))
	require.Equal(t, "em", w[0])
	require.Equal(t, "i", w[1])

	w = camelToWords("emI")
	require.Equal(t, 2, len(w))
	require.Equal(t, "em", w[0])
	require.Equal(t, "i", w[1])

	w = camelToWords("eMI")
	require.Equal(t, 2, len(w))
	require.Equal(t, "e", w[0])
	require.Equal(t, "mi", w[1])

	w = camelToWords("EmiEmi")
	require.Equal(t, 2, len(w))
	require.Equal(t, "emi", w[0])
	require.Equal(t, "emi", w[1])

	w = camelToWords("E1aE2bE3cE")
	require.Equal(t, 4, len(w))
	require.Equal(t, "e1a", w[0])
	require.Equal(t, "e2b", w[1])
	require.Equal(t, "e3c", w[2])
	require.Equal(t, "e", w[3])

	w = camelToWords("AllowIPV6")
	require.Equal(t, 2, len(w))
	require.Equal(t, "allow", w[0])
	require.Equal(t, "ipv6", w[1])

	w = camelToWords("AddCountryCodeTlds")
	require.Equal(t, 4, len(w))
	require.Equal(t, "add", w[0])
	require.Equal(t, "country", w[1])
	require.Equal(t, "code", w[2])
	require.Equal(t, "tlds", w[3])
}

func TestTagParsingComplex(t *testing.T) {
	pv, err := NewPropertyValidator("&cond{w:[FOO,BAR], c:<foo && bar>strnb}, &cond{o:'foo && bar', c:[FOO,BAR]strne}")
	require.Nil(t, err)
	require.NotNil(t, pv)
	require.Equal(t, 2, len(pv.Constraints))
	c0, ok := pv.Constraints[0].(*ConditionalConstraint)
	require.True(t, ok)
	require.Equal(t, 2, len(c0.When))
	require.Equal(t, "FOO", c0.When[0])
	require.Equal(t, "BAR", c0.When[1])
	require.NotNil(t, c0.Constraint)
	c1, ok := c0.Constraint.(*ConditionalConstraint)
	require.True(t, ok)
	require.Equal(t, "foo && bar", c1.Others.String())

	c1, ok = pv.Constraints[1].(*ConditionalConstraint)
	require.True(t, ok)
	require.Nil(t, c1.When)
	require.Equal(t, "foo && bar", c1.Others.String())
	c0, ok = c1.Constraint.(*ConditionalConstraint)
	require.True(t, ok)
	require.Equal(t, 2, len(c0.When))
	require.Equal(t, "FOO", c0.When[0])
	require.Equal(t, "BAR", c0.When[1])
}

func TestTagParsingStructFields(t *testing.T) {
	pv, err := NewPropertyValidator("&strchars{allow:[{\"R16\": [{\"Lo\":32,\"Hi\":512}]}]}")
	require.Nil(t, err)
	require.NotNil(t, pv)
	require.Equal(t, 1, len(pv.Constraints))
	c := pv.Constraints[0].(*StringCharacters)
	require.Equal(t, 1, len(c.AllowRanges))
	ar0 := c.AllowRanges[0]
	require.Equal(t, uint16(32), ar0.R16[0].Lo)
	require.Equal(t, uint16(512), ar0.R16[0].Hi)

	_, err = NewPropertyValidator("&strchars{allow:[{\"R16\": [{\"Lo\":true,\"Hi\":false}]}]}")
	require.NotNil(t, err)
}

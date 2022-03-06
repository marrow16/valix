package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"regexp"
	"testing"
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
	require.Equal(t, "\"(foo)\"", args["str"])
	require.Equal(t, "true", args["bool"])
	require.Equal(t, "0", args["int"])
	require.Equal(t, "1.1", args["float"])

	_, err = argsStringToArgs(constraintName, "str: (,bool:true")
	//                                             012345
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnclosed, 5)), err.Error())

	_, err = argsStringToArgs(constraintName, "str: ),bool:true")
	//                                             012345
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnopened, 5)), err.Error())

	_, err = argsStringToArgs(constraintName, "xxx")
	require.NotNil(t, err)
	_, err = argsStringToArgs(constraintName, "xxx\":\"")
	require.NotNil(t, err)
}

func TestPropertyValidator_ProcessTagItems(t *testing.T) {
	pv := &PropertyValidator{}
	err := pv.processTagItems([]string{})
	require.Nil(t, err)

	require.False(t, pv.Mandatory)
	require.False(t, pv.NotNull)
	err = pv.processTagItems([]string{tagItemNotNull, tagItemMandatory})
	require.True(t, pv.Mandatory)
	require.True(t, pv.NotNull)

	require.Equal(t, 0, len(pv.Constraints))
	err = pv.processTagItems([]string{tagItemConstraintsPrefix + "StringNotEmpty{}"})
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))
	err = pv.processTagItems([]string{tagItemConstraintsPrefix + "&StringNotEmpty{}"})
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.processTagItems([]string{tagItemConstraintsPrefix + "[StringNotEmpty{}, &StringNotEmpty{}]"})
	require.Nil(t, err)
	require.Equal(t, 4, len(pv.Constraints))

	err = pv.processTagItems([]string{tagItemConstraintsPrefix + "[)(]"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnopened, 0), err.Error())

	err = pv.processTagItems([]string{tagItemConstraintsPrefix + "[,]"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintsFormat, ""), err.Error())

	err = pv.processTagItems([]string{"UNKNOWN:x"})
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTokenInTag, "UNKNOWN"), err.Error())
}

func TestPropertyValidator_AddTagItem(t *testing.T) {
	pv := &PropertyValidator{}
	require.False(t, pv.NotNull)
	err := pv.addTagItem(tagItemNotNull)
	require.Nil(t, err)
	require.True(t, pv.NotNull)

	require.False(t, pv.Mandatory)
	err = pv.addTagItem(tagItemMandatory)
	require.Nil(t, err)
	require.True(t, pv.Mandatory)
	err = pv.addTagItem(tagItemOptional)
	require.Nil(t, err)
	require.False(t, pv.Mandatory)

	require.Equal(t, JsonTypeUndefined, pv.Type)
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
			err = pv.addTagItem(tagItemType + ":" + ty)
			require.Nil(t, err)
			require.Equal(t, pv.Type.String(), ty)
		})
	}

	err = pv.addTagItem(tagItemType + ": BAD")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownPropertyType, "BAD"), err.Error())
	err = pv.addTagItem(tagItemType + ": " + jsonTypeTokenString + "  ") //extra spaces
	require.Nil(t, err)
	require.Equal(t, JsonString, pv.Type)

	err = pv.addTagItem("UNKNOWN:xxx")
	require.NotNil(t, err)
}

func TestPropertyValidator_AddTagItemConstraint(t *testing.T) {
	pv := &PropertyValidator{}
	require.Equal(t, 0, len(pv.Constraints))

	err := pv.addTagItem(tagItemConstraint + ":StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.Constraints))

	err = pv.addTagItem(tagItemConstraint + ":UNKNOWN{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "UNKNOWN"), err.Error())
	require.Equal(t, 1, len(pv.Constraints))

	err = pv.addTagItem("&StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.addTagItem("&UNKNOWN{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "UNKNOWN"), err.Error())
	require.Equal(t, 2, len(pv.Constraints))

	err = pv.addTagItem("&StringValidToken{Tokens:['XXX',\"YYY\",'ZZZ']}")
	require.Nil(t, err)
	require.Equal(t, 3, len(pv.Constraints))
	c := pv.Constraints[2].(*StringValidToken)
	require.Equal(t, 3, len(c.Tokens))
	require.Equal(t, "XXX", c.Tokens[0])
	require.Equal(t, "YYY", c.Tokens[1])
	require.Equal(t, "ZZZ", c.Tokens[2])
}

func TestPropertyValidator_AddObjectTagItem_IgnoreUnknownProperties(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err := pv.addTagItem(tagItemObjIgnoreUnknownProperties)
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.IgnoreUnknownProperties)
}

func TestPropertyValidator_AddObjectTagItem_UnknownProperties(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{}}
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err := pv.addTagItem(tagItemObjUnknownProperties + ": true")
	require.Nil(t, err)
	require.True(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err = pv.addTagItem(tagItemObjUnknownProperties + ": FALSE")
	require.Nil(t, err)
	require.False(t, pv.ObjectValidator.IgnoreUnknownProperties)

	err = pv.addTagItem(tagItemObjUnknownProperties + ": INVALID_BOOL_VALUE")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownTagValue, tagItemObjUnknownProperties, "boolean", "INVALID_BOOL_VALUE"), err.Error())
}

func TestPropertyValidator_AddObjectTagItem_Constraint(t *testing.T) {
	pv := &PropertyValidator{ObjectValidator: &Validator{
		Constraints: Constraints{},
	}}
	require.Equal(t, 0, len(pv.ObjectValidator.Constraints))

	err := pv.addTagItem(tagItemObjConstraint + ":&StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 1, len(pv.ObjectValidator.Constraints))

	err = pv.addTagItem(tagItemObjConstraint + ":StringNotEmpty{}")
	require.Nil(t, err)
	require.Equal(t, 2, len(pv.ObjectValidator.Constraints))

	err = pv.addTagItem(tagItemObjConstraint + ":&Unknown{}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgUnknownConstraint, "Unknown"), err.Error())
	require.Equal(t, 2, len(pv.ObjectValidator.Constraints))
}

func TestRebuildConstraintWithArgs(t *testing.T) {
	const constraintName = "TEST_CONSTRAINT"
	var orgConstraint = &StringNotEmpty{}
	c, err := rebuildConstraintWithArgs(constraintName, orgConstraint, "")
	require.Nil(t, err)
	require.NotNil(t, c)
	require.Equal(t, messageNotEmptyString, c.GetMessage())

	c, err = rebuildConstraintWithArgs(constraintName, orgConstraint, ")(")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintArgsParseError, constraintName, fmt.Sprintf(msgUnopened, 0)), err.Error())
}

type dummyNonStructConstraint map[string]struct{}

func (d *dummyNonStructConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (d *dummyNonStructConstraint) GetMessage() string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithNonStructConstraint(t *testing.T) {
	registry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyNonStructConstraint{}
	registry.register(true, testConstraint)
	defer func() {
		registry.reset()
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
func (d *dummyStructConstraintWithUnexportedFields) GetMessage() string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithNonPublicField(t *testing.T) {
	registry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyStructConstraintWithUnexportedFields{}
	registry.register(true, testConstraint)
	defer func() {
		registry.reset()
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
func (d *dummyStructConstraintWithFields) GetMessage() string {
	return ""
}

func TestRebuildConstraintWithArgsFailsWithInvalidArgType(t *testing.T) {
	registry.reset()
	const testConstraintName = "TEST_CONSTRAINT"

	testConstraint := &dummyStructConstraintWithFields{}
	registry.register(true, testConstraint)
	defer func() {
		registry.reset()
	}()

	c, err := rebuildConstraintWithArgs(testConstraintName, testConstraint, "Field2:\"foo\"")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldInvalidValue, testConstraintName, "Field2"), err.Error())
	require.Nil(t, c)
}

func TestBuildConstraintFromTagValueFailsWithBadArgs(t *testing.T) {
	registry.reset()
	testConstraint := &dummyStructConstraintWithFields{}
	registry.register(true, testConstraint)
	defer func() {
		registry.reset()
	}()

	c, err := buildConstraintFromTagValue("&dummyStructConstraintWithFields{UnknownField:\"\"}")
	require.NotNil(t, err)
	require.Equal(t, fmt.Sprintf(msgConstraintFieldUnknown, "dummyStructConstraintWithFields", "UnknownField"), err.Error())
	require.Nil(t, c)
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

	ok := safeSet(fld, "\"baz\"")
	require.True(t, ok)
	require.Equal(t, "baz", fld.String())
	require.Equal(t, "baz", newStruct.AString)

	ok = safeSet(fld, "0")
	require.False(t, ok)
}

func TestSafeSet_Int(t *testing.T) {
	fld, newStruct := getTestFieldValue("AInt")
	require.Equal(t, 0, newStruct.AInt)

	ok := safeSet(fld, "1")
	require.True(t, ok)
	require.Equal(t, int64(1), fld.Int())
	require.Equal(t, 1, newStruct.AInt)

	ok = safeSet(fld, "\"1\"")
	require.False(t, ok)
	ok = safeSet(fld, "true")
	require.False(t, ok)
	ok = safeSet(fld, "1.1")
	require.False(t, ok)
}

func TestSafeSet_UInt(t *testing.T) {
	fld, newStruct := getTestFieldValue("AUInt")
	require.Equal(t, uint(0), newStruct.AUInt)

	ok := safeSet(fld, "1")
	require.True(t, ok)
	require.Equal(t, uint64(1), fld.Uint())
	require.Equal(t, uint(1), newStruct.AUInt)

	ok = safeSet(fld, "\"1\"")
	require.False(t, ok)
	ok = safeSet(fld, "true")
	require.False(t, ok)
	ok = safeSet(fld, "1.1")
	require.False(t, ok)
}

func TestSafeSet_Float(t *testing.T) {
	fld, newStruct := getTestFieldValue("AFloat")
	require.Equal(t, float64(0), newStruct.AFloat)

	ok := safeSet(fld, "1")
	require.True(t, ok)
	require.Equal(t, float64(1), fld.Float())
	require.Equal(t, float64(1), newStruct.AFloat)

	ok = safeSet(fld, "\"1\"")
	require.False(t, ok)
	ok = safeSet(fld, "true")
	require.False(t, ok)
}

func TestSafeSet_Bool(t *testing.T) {
	fld, newStruct := getTestFieldValue("ABool")
	require.Equal(t, false, newStruct.ABool)

	ok := safeSet(fld, "true")
	require.True(t, ok)
	require.Equal(t, true, fld.Bool())
	require.Equal(t, true, newStruct.ABool)

	ok = safeSet(fld, "\"1\"")
	require.False(t, ok)
	ok = safeSet(fld, "1.1")
	require.False(t, ok)
}

func TestSafeSet_Slice(t *testing.T) {
	fld, newStruct := getTestFieldValue("ASlice")
	require.Equal(t, 0, len(newStruct.ASlice))

	ok := safeSet(fld, "[\"foo\", \"bar\"]")
	require.True(t, ok)
	arr := fld.Interface().([]string)
	require.Equal(t, 2, len(arr))
	require.Equal(t, 2, len(newStruct.ASlice))
	require.Equal(t, "foo", arr[0])
	require.Equal(t, "bar", arr[1])

	ok = safeSet(fld, "\"foo\"")
	require.False(t, ok)
	ok = safeSet(fld, "0")
	require.False(t, ok)
}

func TestSafeSet_Regexp(t *testing.T) {
	fld, newStruct := getTestFieldValue("ARegexp")
	require.NotNil(t, newStruct.ARegexp)

	const pattern = "^([a-fA-F0-9]{8})$"
	ok := safeSet(fld, "\""+pattern+"\"")
	require.True(t, ok)
	rx := fld.Interface().(regexp.Regexp)
	require.Equal(t, pattern, rx.String())
	require.Equal(t, pattern, newStruct.ARegexp.String())

	ok = safeSet(fld, "^^^")
	require.False(t, ok)
	ok = safeSet(fld, "1.1")
	require.False(t, ok)
	ok = safeSet(fld, "true")
	require.False(t, ok)
}

func TestSafeSet_OtherUnknown(t *testing.T) {
	fld, _ := getTestFieldValue("AOther")

	ok := safeSet(fld, "")
	require.False(t, ok)
}

func TestItemsToSlice_String(t *testing.T) {
	sampleArr := []string{}
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

	v, ok = itemsToSlice(itemType, "[1,'bar','baz']")
	require.False(t, ok)
	v, ok = itemsToSlice(itemType, "[x,'bar','baz']")
	require.False(t, ok)
}

func TestItemsToSlice_Int(t *testing.T) {
	sampleArr := []int{}
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
	sampleArr := []uint{}
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
	sampleArr := []float64{}
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
	sampleArr := []bool{}
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

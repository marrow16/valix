package valix

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"reflect"
	"testing"
)

func TestTagAliasesRepoInitialised(t *testing.T) {
	require.NotNil(t, tagAliasesRepo)
	require.Equal(t, 0, len(tagAliasesRepo.aliases))
}

func TestCanRegisterTagAliases(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	require.Equal(t, 0, len(tagAliasesRepo.aliases))

	RegisterTagTokenAliases(TagAliases{
		"nnm": "notNull,mandatory",
		"s":   "type:string",
	})
	require.Equal(t, 2, len(tagAliasesRepo.aliases))
	require.Equal(t, "notNull,mandatory", tagAliasesRepo.aliases["nnm"])
	require.Equal(t, "type:string", tagAliasesRepo.aliases["s"])

	RegisterTagTokenAlias("s", "type:string,&StringNotEmpty{}")
	require.Equal(t, 2, len(tagAliasesRepo.aliases))
	require.Equal(t, "type:string,&StringNotEmpty{}", tagAliasesRepo.aliases["s"])
}

func TestTagAliasesRepoClears(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	require.Equal(t, 0, len(tagAliasesRepo.aliases))

	RegisterTagTokenAliases(TagAliases{
		"nnm": "notNull,mandatory",
		"s":   "type:string",
	})
	require.Equal(t, 2, len(tagAliasesRepo.aliases))

	ClearTagTokenAliases()
	require.Equal(t, 0, len(tagAliasesRepo.aliases))
}

func TestTagAliasesRepoErrorsOnNotFound(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	_, err := tagAliasesRepo.resolve([]string{"$non_existent_alias"})
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(errMsgUnknownTagAlias, "non_existent_alias"), err.Error())
}

func TestTagAliasesRepoErrorsOnCyclic(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	tagAliasesRepo.registerSingle("cyclic", "$cyclic")
	_, err := tagAliasesRepo.resolve([]string{"$cyclic"})
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(errMsgCyclicTagAlias, "cyclic"), err.Error())

	tagAliasesRepo.reset()
	tagAliasesRepo.registerSingle("first", "$second")
	tagAliasesRepo.registerSingle("second", "$third")
	tagAliasesRepo.registerSingle("third", "foo,$first,bar")
	_, err = tagAliasesRepo.resolve([]string{"$first"})
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(errMsgCyclicTagAlias, "first"), err.Error())
}

func TestTagAliasesRepoErrorsWithBadReplacement(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	tagAliasesRepo.registerSingle("bad", "'")
	_, err := tagAliasesRepo.resolve([]string{"$bad"})
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(errMsgAliasParse, "bad", "unclosed ''' at position 0"), err.Error())
}

func TestTagAliasesResolvesSingleLevel(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	tagAliasesRepo.registerSingle("test", "FOO")
	result, err := tagAliasesRepo.resolve([]string{"$test"})
	require.NoError(t, err)
	require.Equal(t, 1, len(result))
	require.Equal(t, "FOO", result[0])

	result, err = tagAliasesRepo.resolve([]string{"$test", "FOO2"})
	require.NoError(t, err)
	require.Equal(t, 2, len(result))
	require.Equal(t, "FOO", result[0])
	require.Equal(t, "FOO2", result[1])

	tagAliasesRepo.registerSingle("test", "FOO,BAR,BAZ")
	result, err = tagAliasesRepo.resolve([]string{"$test"})
	require.NoError(t, err)
	require.Equal(t, 3, len(result))
	require.Equal(t, "FOO", result[0])
	require.Equal(t, "BAR", result[1])
	require.Equal(t, "BAZ", result[2])
}

func TestTagAliasesResolvesDeep(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	tagAliasesRepo.registerSingle("first", "FOO1, $second, BAR1")
	tagAliasesRepo.registerSingle("second", "FOO2, $third, BAR2")
	tagAliasesRepo.registerSingle("third", "FOO3, BAR3")

	result, err := tagAliasesRepo.resolve([]string{"$first"})
	require.NoError(t, err)
	require.Equal(t, 6, len(result))
	require.Equal(t, "FOO1", result[0])
	require.Equal(t, "FOO2", result[1])
	require.Equal(t, "FOO3", result[2])
	require.Equal(t, "BAR3", result[3])
	require.Equal(t, "BAR2", result[4])
	require.Equal(t, "BAR1", result[5])
}

func TestTagAliasErrorsCauseTagParsingErrors(t *testing.T) {
	tagAliasesRepo.reset()
	defer tagAliasesRepo.reset()

	tagAliasesRepo.registerSingle("bad", "'")
	tagAliasesRepo.registerSingle("cyclic", "$cyclic")

	pv := &PropertyValidator{}
	err := pv.processV8nTagValue("Foo", "foo", "$unknown")
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgWrapped, "Foo", "foo", fmt.Sprintf(errMsgUnknownTagAlias, "unknown")), err.Error())

	err = pv.processV8nTagValue("Foo", "foo", "$bad")
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgWrapped, "Foo", "foo", fmt.Sprintf(errMsgAliasParse, "bad", "unclosed ''' at position 0")), err.Error())

	err = pv.processV8nTagValue("Foo", "foo", "$cyclic")
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgWrapped, "Foo", "foo", fmt.Sprintf(errMsgCyclicTagAlias, "cyclic")), err.Error())
}

func TestCustomTagTokensInitialised(t *testing.T) {
	require.NotNil(t, customTagTokenRegistry)
	require.Equal(t, 0, len(customTagTokenRegistry.handlers))
}

func TestCanRegisterCustomTagToken(t *testing.T) {
	customTagTokenRegistry.reset()
	defer customTagTokenRegistry.reset()

	custom := &testCustomTagToken{}
	RegisterCustomTagToken("test", custom)

	require.Equal(t, 1, len(customTagTokenRegistry.handlers))
}

func TestCustomTagTokenRegistryClears(t *testing.T) {
	customTagTokenRegistry.reset()
	defer customTagTokenRegistry.reset()

	custom := &testCustomTagToken{}
	RegisterCustomTagToken("test", custom)
	require.Equal(t, 1, len(customTagTokenRegistry.handlers))

	ClearCustomTagTokens()
	require.Equal(t, 0, len(customTagTokenRegistry.handlers))
}

func TestCustomTagTokenUsed(t *testing.T) {
	customTagTokenRegistry.reset()
	defer customTagTokenRegistry.reset()

	type testStruct struct {
		Foo string `json:"foo" v8n:"test:test_value"`
	}
	_, err := ValidatorFor(testStruct{}, nil)
	require.Error(t, err)
	require.Equal(t, fmt.Sprintf(msgWrapped, "Foo", "foo", fmt.Sprintf(msgUnknownTokenInTag, "test")), err.Error())

	custom := &testCustomTagToken{}
	RegisterCustomTagToken("test", custom)
	require.Equal(t, 1, len(customTagTokenRegistry.handlers))

	v, err := ValidatorFor(testStruct{}, nil)
	require.NoError(t, err)
	require.NotNil(t, v)
	require.Equal(t, 1, custom.hits)
	require.Equal(t, "test", custom.token)
	require.True(t, custom.hasValue)
	require.Equal(t, "test_value", custom.value)
	require.Equal(t, "foo", custom.propertyName)
	require.Equal(t, "Foo", custom.fieldName)

	custom.errors = true
	_, err = ValidatorFor(testStruct{}, nil)
	require.Error(t, err)
	require.Equal(t, 2, custom.hits)
	require.Equal(t, fmt.Sprintf(msgWrapped, "Foo", "foo", customTagTokenError), err.Error())
}

const customTagTokenError = "custom tag token errored"

type testCustomTagToken struct {
	hits         int
	token        string
	hasValue     bool
	value        string
	propertyName string
	fieldName    string
	errors       bool
}

func (c *testCustomTagToken) Handle(token string, hasValue bool, tokenValue string, pv *PropertyValidator, propertyName string, fieldName string) error {
	c.hits++
	c.token = token
	c.hasValue = hasValue
	c.value = tokenValue
	c.propertyName = propertyName
	c.fieldName = fieldName
	if c.errors {
		return errors.New(customTagTokenError)
	}
	return nil
}

func TestCustomTagsInitialised(t *testing.T) {
	require.NotNil(t, customTagsRegistry)
	require.Equal(t, 0, len(customTagsRegistry.handlers))
}

func TestCanRegisterCustomTag(t *testing.T) {
	customTagsRegistry.reset()
	defer customTagsRegistry.reset()

	require.Equal(t, 0, len(customTagsRegistry.handlers))

	custom := &testCustomTag{}
	RegisterCustomTag("test", custom)
	require.Equal(t, 1, len(customTagsRegistry.handlers))
}

func TestCustomTagRegistryClears(t *testing.T) {
	customTagsRegistry.reset()
	defer customTagsRegistry.reset()

	require.Equal(t, 0, len(customTagsRegistry.handlers))

	custom := &testCustomTag{}
	RegisterCustomTag("test", custom)
	require.Equal(t, 1, len(customTagsRegistry.handlers))

	ClearCustomTags()
	require.Equal(t, 0, len(customTagsRegistry.handlers))
}

func TestCustomTagUsed(t *testing.T) {
	customTagsRegistry.reset()
	defer customTagsRegistry.reset()

	require.Equal(t, 0, len(customTagsRegistry.handlers))

	custom := &testCustomTag{}
	RegisterCustomTag("test", custom)
	require.Equal(t, 1, len(customTagsRegistry.handlers))

	type testStruct struct {
		Foo string `json:"foo" test:"test_value"`
	}
	v, err := ValidatorFor(testStruct{}, nil)
	require.NoError(t, err)
	require.True(t, v.Properties["foo"].Mandatory)
	require.True(t, v.Properties["foo"].NotNull)
	require.Equal(t, "test_value", custom.value)

	custom.errors = true
	_, err = ValidatorFor(testStruct{}, nil)
	require.Error(t, err)
	require.Equal(t, customTagError, err.Error())
}

const customTagError = "custom tag errored"

type testCustomTag struct {
	value  string
	errors bool
}

func (t *testCustomTag) Handle(tag string, tagValue string, commaParsed []string, pv *PropertyValidator, fld reflect.StructField) error {
	if t.errors {
		return errors.New(customTagError)
	}
	pv.NotNull = true
	pv.Mandatory = true
	t.value = tagValue
	return nil
}

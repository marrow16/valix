package valix

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// TagHandler is the interface for handling custom tags
//
// Custom tags are used when building validators from structs using valix.ValidatorFor
type TagHandler interface {
	Handle(tag string, tagValue string, commaParsed []string, pv *PropertyValidator, fld reflect.StructField) error
}

// RegisterCustomTag registers a custom tag - for use when building validators from structs
// using valix.ValidatorFor, for example...
//	RegisterCustomTag("custom", myTagHandler)
// and then use the custom tag...
//	type MyStruct struct {
//		Foo string `custom:"<custom_value>"`
//	}
// and then build validator for struct...
//	validator, err := valix.ValidatorFor(MyStruct{}, nil)
// will call the custom tag handler when building the validator
func RegisterCustomTag(tag string, handler TagHandler) {
	customTagsRegistry.register(tag, handler)
}

// ClearCustomTags clears any registered custom tags that were registered using RegisterCustomTag
func ClearCustomTags() {
	customTagsRegistry.reset()
}

// CustomTagTokenHandler is the interface for handling custom tag tokens (i.e. custom tokens in the `v8n` tag)
type CustomTagTokenHandler interface {
	Handle(token string, hasValue bool, tokenValue string, pv *PropertyValidator, propertyName string, fieldName string) error
}

// RegisterCustomTagToken registers a custom tag token handler - registered custom tag tokens can be used
// within the `v8n` tag
//
// Example:
//	RegisterCustomTagToken("my_token", myCustomTokenHandler)
// and then use the custom tag token...
//	type MyStruct struct {
//		Foo string `json:"foo" v8n:"my_token: hello"`
//	}
func RegisterCustomTagToken(token string, handler CustomTagTokenHandler) {
	customTagTokenRegistry.register(token, handler)
}

// ClearCustomTagTokens clears any custom tag tokens registered using RegisterCustomTagToken
func ClearCustomTagTokens() {
	customTagTokenRegistry.reset()
}

// RegisterTagTokenAlias registers a tag alias - a tag alias can be used in the `v8n` tag to specify
// multiple tokens using a single alias
//
// Example:
//	RegisterTagTokenAlias("mnnne", "mandatory,notNull,&StringNotEmpty{}"
// and the use the alias with the `v8n` tag...
//	type MyStruct struct {
//		Foo string `json:"foo" v8n:"$mnnne, &StringNotBlank{}"`
//	}
// would be the equivalent of...
//	type MyStruct struct {
//		Foo string `json:"foo" v8n:"mandatory,notNull,&StringNotEmpty{}, &StringNotBlank{}"`
//	}
func RegisterTagTokenAlias(alias string, val string) {
	tagAliasesRepo.registerSingle(alias, val)
}

// RegisterTagTokenAliases register multiple tag aliases - see RegisterTagTokenAlias
func RegisterTagTokenAliases(aliases TagAliases) {
	tagAliasesRepo.register(aliases)
}

// ClearTagTokenAliases clears any tag aliases registered using RegisterTagTokenAlias / RegisterTagTokenAliases
func ClearTagTokenAliases() {
	tagAliasesRepo.clear()
}

type customTagTokenHolder struct {
	handler       CustomTagTokenHandler
	requiresValue bool
}

type customTagTokens struct {
	handlers map[string]*customTagTokenHolder
	sync     *sync.Mutex
}

var customTagTokenRegistry = &customTagTokens{
	handlers: map[string]*customTagTokenHolder{},
	sync:     &sync.Mutex{},
}

// reset for testing
func (r *customTagTokens) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.handlers = map[string]*customTagTokenHolder{}
}

func (r *customTagTokens) register(token string, handler CustomTagTokenHandler) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.handlers[token] = &customTagTokenHolder{
		handler: handler,
	}
}

func (r *customTagTokens) handle(token string, hasValue bool, tokenValue string, pv *PropertyValidator, propertyName string, fieldName string) (bool, error) {
	defer r.sync.Unlock()
	r.sync.Lock()
	if h, ok := r.handlers[token]; ok && h.handler != nil {
		return true, h.handler.Handle(token, hasValue, tokenValue, pv, propertyName, fieldName)
	}
	return false, nil
}

type TagAliases map[string]string

type tagAliasesRepository struct {
	aliases TagAliases
	sync    *sync.Mutex
}

var tagAliasesRepo = &tagAliasesRepository{
	aliases: TagAliases{},
	sync:    &sync.Mutex{},
}

// reset for testing
func (r *tagAliasesRepository) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.aliases = TagAliases{}
}

func (r *tagAliasesRepository) register(aliases TagAliases) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for k, v := range aliases {
		r.aliases[k] = v
	}
}

func (r *tagAliasesRepository) registerSingle(alias string, val string) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.aliases[alias] = val
}

func (r *tagAliasesRepository) clear() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.aliases = TagAliases{}
}

func (r *tagAliasesRepository) resolve(tagItems []string) ([]string, error) {
	foundAliasesAt := map[int]string{}
	for i, tagItem := range tagItems {
		if strings.HasPrefix(tagItem, "$") {
			foundAliasesAt[i] = tagItem[1:]
		}
	}
	if len(foundAliasesAt) == 0 {
		return tagItems, nil
	}
	defer r.sync.Unlock()
	r.sync.Lock()

	result := make([]string, 0, len(tagItems)+(len(foundAliasesAt)*4))
	for i, item := range tagItems {
		if alias, ok := foundAliasesAt[i]; ok {
			resolved, err := r.resolveItem(alias, map[string]bool{})
			if err != nil {
				return nil, err
			}
			result = append(result, resolved...)
		} else {
			result = append(result, item)
		}
	}
	return result, nil
}

const (
	errMsgCyclicTagAlias  = "cyclic tag alias reference '$%s'"
	errMsgUnknownTagAlias = "unknown tag alias reference '$%s'"
	errMsgAliasParse      = "error parsing resolved tag alias '$%s' - %s"
)

func (r *tagAliasesRepository) resolveItem(alias string, seen map[string]bool) ([]string, error) {
	if seen[alias] {
		return nil, fmt.Errorf(errMsgCyclicTagAlias, alias)
	}
	seen[alias] = true
	aliased, ok := r.aliases[alias]
	if !ok {
		return nil, fmt.Errorf(errMsgUnknownTagAlias, alias)
	}
	result, err := parseCommas(aliased)
	if err != nil {
		return nil, fmt.Errorf(errMsgAliasParse, alias, err.Error())
	}
	return r.resolveItems(result, seen)
}

func (r *tagAliasesRepository) resolveItems(items []string, seen map[string]bool) ([]string, error) {
	result := make([]string, 0, len(items))
	for _, item := range items {
		if strings.HasPrefix(item, "$") {
			resolved, err := r.resolveItem(item[1:], seen)
			if err != nil {
				return nil, err
			}
			result = append(result, resolved...)
		} else {
			result = append(result, item)
		}
	}
	return result, nil
}

type customTags struct {
	handlers map[string]TagHandler
	sync     *sync.Mutex
}

var customTagsRegistry = &customTags{
	handlers: map[string]TagHandler{},
	sync:     &sync.Mutex{},
}

func (r *customTags) register(tag string, handler TagHandler) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.handlers[tag] = handler
}

func (r *customTags) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.handlers = map[string]TagHandler{}
}

func (r *customTags) processField(fld reflect.StructField, pv *PropertyValidator) error {
	defer r.sync.Unlock()
	r.sync.Lock()
	var result error = nil
	for k, h := range r.handlers {
		if h != nil {
			if tag, ok := fld.Tag.Lookup(k); ok {
				commaParsed, _ := parseCommas(tag)
				result = h.Handle(k, tag, commaParsed, pv, fld)
				if result != nil {
					break
				}
			}
		}
	}
	return result
}

# Valix
[![GoDoc](https://godoc.org/github.com/marrow16/valix?status.svg)](https://pkg.go.dev/github.com/marrow16/valix)
[![Latest Version](https://img.shields.io/github/v/tag/marrow16/valix.svg?sort=semver&style=flat&label=version&color=blue)](https://github.com/marrow16/valix/releases)
[![codecov](https://codecov.io/gh/marrow16/valix/branch/master/graph/badge.svg)](https://codecov.io/gh/marrow16/valix)
[![Go Report Card](https://goreportcard.com/badge/github.com/marrow16/valix)](https://goreportcard.com/report/github.com/marrow16/valix)
[![Maintainability](https://api.codeclimate.com/v1/badges/1d64bc6c8474c2074f2b/maintainability)](https://codeclimate.com/github/marrow16/valix/maintainability)


Valix - Go package for validating requests

## Contents
* [Overview](#overview)
* [Installation](#installation)
* [Features](#features)
* [Concepts](#concepts)
* [Examples](#examples)
  * [Creating Validators](#creating-validators)
  * [Using Validators](#using-validators)
* [Constraints](#constraints)
  * [Common Constraints](#common-constraints)
  * [Constraint Sets](#constraint-sets)
  * [Custom Constraints](#custom-constraints)
  * [Constraints Registry](#constraints-registry)
  * [Conditional Constraints](#conditional-constraints)
* [Polymorphic Validation](#polymorphic-validation)
* [Validation Tags](#validation-tags)
* [I18n Support](#internationalisation-support)

## Overview

Validate JSON requests in the form of `*http.Request`, `map[string]interface{}` or `[]interface{}`

## Installation
To install Valix, use go get:

    go get github.com/marrow16/valix

To update Valix to the latest version, run:

    go get -u github.com/marrow16/valix

## Features

* Deep validation (define validation for properties where those properties are objects or arrays of objects that also need validating)
* Create validators from structs or define them as code (see [Creating Validators](#creating-validators))
* Validate `http.Request` (into struct or `map[string]interface{}`)
* Finds all validation violations - not just the first one! (see [Using Validators](#using-validators)) and provides information for each violation (property name, path and message)
* Rich set of pre-defined common constraints (see [Common Constraints](#common-constraints))
* Customisable constraints (see [Custom Constraints](#custom-constraints))
* Conditional constraints to support partial polymorphic request models (see [Conditional Constraints](#conditional-constraints))
* Full [Polymorphic Validation](#polymorphic-validation)
* [Support for i18n](#internationalisation-support) - enabling translation of validation messages (inc. `http.Request` language and region detection)
* Validators fully marshalable and unmarshalable *(save/share validators as JSON)*
* Highly extensible *(add your own constraints, messages, presets etc. without a PR to this repository)*
* 100% tested (see [Codecov.io](https://codecov.io/gh/marrow16/valix))
* *Coming soon - generate validators from [OpenApi/Swagger](https://swagger.io/docs/specification/about/) JSON*

## Concepts

Valix is based on the concept that incoming API requests (such as `POST`, `PUT` etc.) should be
validated early - against a definition of what the request body should look like.

At validation (of a JSON object as example) the following steps are performed:
* Optionally check any constraints on the overall object (specified in `Validator.Constraints`)
* Check if there are any unknown properties (unless `Validator.IgnoreUnknownProperties` is set to true)
* For each defined property in `Validator.Properties`:
  * Check the property is present (if `PropertyValidator.Mandatory` is set to true)
  * Check the property value is non-null (if `PropertyValidator.NotNull` is set to true)
  * Check the property value is of the correct type (if `PropertyValidator.Type` is set)
  * Check the property value against constraints (specified in `PropertyValidator.Constraints`)
  * If the property value is an object or array (and `PropertyValidator.ObjectValidator` is specified) - check the value using the validator (see top of this process)

The validator does **_not_** stop on the first problem it finds - it finds all problems and returns them as a list (slice) of 'violations'.  Each violation has a message along with the name and path of the property that failed.
However, custom constraints can be defined that will either stop the entire validation or cease further validation constraints on the current property


## Examples

### Creating Validators

Validators can be created from existing structs - adding `v8n` tags (in conjunction with existing `json` tags), for example:
```go
package main

import (
    "github.com/marrow16/valix"
)

type AddPersonRequest struct {
    Name string `json:"name" v8n:"notNull,mandatory,constraints:[StringNoControlCharacters{},StringLength{Minimum: 1, Maximum: 255}]"`
    Age int `json:"age" v8n:"type:Integer,notNull,mandatory,constraint:PositiveOrZero{}"`
}
var AddPersonRequestValidator = valix.MustCompileValidatorFor(AddPersonRequest{}, nil)
```
(see [Validation Tags](#validation-tags) for documentation on `v8n` tags)

Or in slightly more abbreviated form (using `&` to denote constraint tokens):
```go
package main

import (
    "github.com/marrow16/valix"
)

type AddPersonRequest struct {
    Name string `json:"name" v8n:"notNull,mandatory,&StringNoControlCharacters{},&StringLength{Minimum: 1, Maximum: 255}"`
    Age int `json:"age" v8n:"type:Integer,notNull,mandatory,&PositiveOrZero{}"`
}
var AddPersonRequestValidator = valix.MustCompileValidatorFor(AddPersonRequest{}, nil)
```

The `valix.MustCompileValidatorFor()` function panics if the validator cannot be compiled.  If you do not want a panic but would rather see the compilation error instead then use the `valix.ValidatorFor()` function instead.  


Alternatively, Validators can be expressed effectively without a struct, for example:
```go
package main

import (
    "github.com/marrow16/valix"
)

var CreatePersonRequestValidator = &valix.Validator{
    IgnoreUnknownProperties: false,
    Properties: valix.Properties{
        "name": {
            Type:         valix.JsonString,
            NotNull:      true,
            Mandatory:    true,
            Constraints:  valix.Constraints{
                &valix.StringNoControlCharacters{},
                &valix.StringLength{Minimum: 1, Maximum: 255},
            },
        },
        "age": {
            Type:         valix.JsonInteger,
            NotNull:      true,
            Mandatory:    true,
            Constraints:  valix.Constraints{
                &valix.PositiveOrZero{},
            },
        },
    },
}
```

####  Additional validator options

Validators can have additional properties that control the overall validation behaviour.  These properties are described as follows:

| Property                  | Description                                                                                                                                                                                                                                                                                                                                        |
|---------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `AllowArray`              | (default `false`) Allows the validator to accept JSON arrays - and validate each item in the array<br/>Setting this option to `true` whilst leaving the `DisallowObject` option as `false` means that the validator will accept either a JSON array or object                                                                                      |
| `AllowNullJson`           | Normally, a validator sees Null JSON (i.e. JSON string just containing the word `null`) as a violation - as it represents neither an object nor an array.<br/>Setting this option to `true` disables this behaviour (and results of successful validation may return a `nil` map/slice)<br/>*NB. This option is only used by top-level validators* |
| `DisallowObject`          | (default `false`) Prevents the validator from accepting JSON objects<br/>Should only be set to `true` when `AllowArray` is also set to `true`                                                                                                                                                                                                      |
| `IgnoreUnknownProperties` | Normally, a validator will report as a violation any properties not defined within the validator<br/>Setting this option to `true` means the validator will not check for unknown properties                                                                                                                                                       |
| `OrderedPropertyChecks`   | Normally, a validator checks specified properties in an unpredictable order (as they are stored in a map).<br/>Setting this option to `true` means that the validator will check properties in order - by their `Order` field (or `order` tag) and then by name                                                                                    |
| `StopOnFirst`             | Normally, a validator will find all constraint violations<br/>Setting this option to `true` causes the validator to stop when it finds the first violation<br/>*NB. This option is only used by top-level validators*                                                                                                                              |
| `UseNumber`               | Validators use `json.NewDecoder()` to decode JSON<br/>Setting this option to `true` instructs the validator to call `Decoder.UseNumber()` prior to decoding<br/>*NB. This option is only used by top-level validators*                                                                                                                             |


### Using Validators

Once a validator has been created (using previous examples in [Creating Validators](#creating-validators)), they can be used in several ways:

#### Validating a request into a struct

A request `*http.Request` can be validated into a struct: 
```go
package main

import (
    "encoding/json"
    "net/http"

    "github.com/marrow16/valix"
)

func AddPersonHandler(w http.ResponseWriter, r *http.Request) {
    addPersonReq := &AddPersonRequest{}
    ok, violations, _ := CreatePersonRequestValidator.RequestValidateInto(r, addPersonReq)
    if !ok {
        // write an error response with full info - using violations information
        valix.SortViolationsByPathAndProperty(violations)
        errResponse := map[string]interface{}{
            "$error": "Request invalid",
            "$details": violations,
        }
        w.WriteHeader(http.StatusUnprocessableEntity)
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(errResponse)
        return
    }
    // the addPersonReq will now be a validated struct
    w.Header().Set("Content-Type", "application/json")
    _ = json.NewEncoder(w).Encode(addPersonReq)
}
```

#### Validating a string or reader into a struct

A string, representing JSON, can be validated into a struct:
```go
package main

import (
    "testing"

    "github.com/marrow16/valix"
    "github.com/stretchr/testify/require"
)

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
```
Also, a reader `io.Reader` can be validated into a struct using the `.ValidateReaderIntoStruct()` method of the validator:
```go
package main

import (
    "strings"
    "testing"

    "github.com/marrow16/valix"
    "github.com/stretchr/testify/require"
)

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
```

#### Validating a map
Validators can also validate a `map[string]interface{}` representation of a JSON object:
```go
package main

import (
    "testing"

    "github.com/marrow16/valix"
    "github.com/stretchr/testify/require"
)

func TestValidateMap(t *testing.T) {
    req := map[string]interface{}{
        "name": "",
        "age": -1,
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
        "age": 25,
    }
    ok, _ = AddPersonRequestValidator.Validate(req)
    require.True(t, ok)
}
```

#### Validating a slice
Validators can also validate a slice `[]interface{}` representation of a JSON object, where each object element in the slice is validated:
```go
package main

import (
    "testing"

    "github.com/marrow16/valix"
    "github.com/stretchr/testify/require"
)

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
```

## Constraints

In Valix, a constraint is a particular validation rule that must be satisfied. For a constraint to be used
by the validator it must implement the `valix.Constraint` interface.

### Common Constraints

Valix provides a rich set of over 100 pre-defined common constraints (plus many common regex patterns) - 
see [Constraints Reference](https://github.com/marrow16/valix/wiki/Constraints-Reference) wiki documentation for full reference with examples.  

### Constraint Sets

It is not uncommon in APIs for many properties in different requests to share a common set of constraints.
For this reason, Valix provides a `ConstraintSet` - which is itself a `Constraint` but contains a list of sub-constraints.
```go
type ConstraintSet struct {
    Constraints Constraints
    Message string
}
```
When checking a `ConstraintSet`, the contained constraints are checked sequentially but the overall
set stops on the first failing constraint.   

If a `Message` is provided (non-empty string) then that message is used for any of the failing constraints -
otherwise the individual constraint fail messages are used.

The following is an example of a constraint set which imposes a complex constraint
(although one that could probably be more easily achieved using `valix.StringPattern`)

```go
package main

import (
    "unicode"
    "github.com/marrow16/valix"
)

var MySet = &valix.ConstraintSet{
    Constraints: valix.Constraints{
        &valix.StringTrim{},
        &valix.StringNotEmpty{},
        &valix.StringLength{Minimum: 16, Maximum: 64},
        valix.NewCustomConstraint(func(value interface{}, vcx *valix.ValidatorContext, this *valix.CustomConstraint) (bool, string) {
            if str, ok := value.(string); ok {
                if len(str) == 0 || str[0] < 'A' || str[0] > 'Z' {
                    return false, this.GetMessage(vcx)
                }
            }
            return true, ""
        }, ""),
        &valix.StringCharacters{
            AllowRanges: []unicode.RangeTable{
                {R16: []unicode.Range16{{'0', 'z', 1}}},
            },
            DisallowRanges: []unicode.RangeTable{
                {R16: []unicode.Range16{{0x003a, 0x0040, 1}}},
                {R16: []unicode.Range16{{0x005b, 0x005e, 1}}},
                {R16: []unicode.Range16{{0x0060, 0x0060, 1}}},
            },
        },
    },
    Message: "String value length must be between 16 and 64 chars; must be letters (upper or lower), digits or underscores; must start with an uppercase letter",
}
```

Constraint sets can also be registered, making them available in `v8n` struct tags.


### Custom Constraints

If you need a constraint for a specific domain validation, there are two ways to do this...

Create a re-usable constraint (which implements the `valix.Constraint` interface), example:
```go
package main

import (
    "strings"
    "github.com/marrow16/valix"
)

type NoFoo struct {
}

func (c *NoFoo) Check(value interface{}, vcx *valix.ValidatorContext) (bool, string) {
    if str, ok := value.(string); ok {
        return !strings.Contains(str, "foo"), c.GetMessage(vcx)
    }
    return true, ""
}

func (c *NoFoo) GetMessage(tcx I18nContext) string {
    return "Value must not contain \"foo\""
}
```

Or create a custom constraint on the fly with check function, example:<br>
_Note: Custom constraints using functions means that the validator cannot be marshalled/unmarshalled_
```go
package main

import (
    "strings"
    "github.com/marrow16/valix"
)

var myValidator = &valix.Validator{
    IgnoreUnknownProperties: false,
    Properties: valix.Properties{
        "foo": {
            Type:         valix.JsonString,
            NotNull:      true,
            Mandatory:    true,
            Constraints:  valix.Constraints{
                valix.NewCustomConstraint(func(value interface{}, vcx *valix.ValidatorContext, cc *valix.CustomConstraint) (bool, string) {
                    if str, ok := value.(string); ok {
                        return !strings.Contains(str, "foo"), cc.GetMessage(vcx)
                    }
                    return true, ""
                }, "Value must not contain \"foo\""),
            },
        },
    },
}
```

### Constraints Registry

All of the Valix common constraints are loaded into a registry - the registry enables the `v8n` tags to reference these.

If you want to make your own custom constraint available for use in `v8n` tags, it must also be registered.  For example:

```go
package main

import (
    "strings"
    "github.com/marrow16/valix"
)

func init() {
    valix.RegisterConstraint(&NoFoo{})
}

// and the constraint can now be used in `v8n` tag...
type MyRequest struct {
    Name string `json:"name" v8n:"&NoFoo{}"`
}
```

### Required/Unwanted Properties

When properties are required, or unwanted, according to the presence of other properties - use the `v8n` tag tokens `required_with:<expr>` or `unwanted_with:<expr>` (abbreviated forms `+:` and `-:` respectively).

The simplest expression is just the name of another property - the following example shows making two properties `foo` and `bar` mutually inclusive (i.e. if one is present then the other is required):
```go
type ExampleMutuallyInclusive struct {
    Foo string `json:"foo" v8n:"+:bar, +msg:'foo required when bar present'"`
    Bar string `json:"bar" v8n:"+:foo, +msg:'bar required when foo present'"`
}
```
Or another example, using the `unwanted_with:` token, to make two properties mutually exclusive (i.e. if one is present then the other must not):
```go
type ExampleMutuallyExclusive struct {
    Foo string `json:"foo" v8n:"-:bar, -msg:'foo and bar are mutually exclusive'"`
    Bar string `json:"bar" v8n:"-:foo, -msg:'foo and bar are mutually exclusive'"`
}
```

The `required_with:` and `unwanted_with` can also use more complex boolean expressions - the following example demonstrates making two out of three properties mutually inclusive but not all three:  
```go
type ExampleTwoOfThreeMutuallyInclusive struct {
    Foo string `json:"foo" v8n:"+:(bar || baz) && !(bar && baz), -:bar && baz"`
    Bar string `json:"bar" v8n:"+:(foo || baz) && !(foo && baz), -:foo && baz"`
    Baz string `json:"baz" v8n:"+:(foo || bar) && !(foo && bar), -:foo && bar"`
}
```

The boolean property expressions can also traverse up and down the object tree (using JSON `.` path style notation).  The following demonstrates:
```go
type ExampleUpAndDownRequired struct {
    Foo string `json:"foo" v8n:"+:sub.foo"`
    Bar string `json:"bar" v8n:"+:sub.bar"`
    Sub struct {
        SubFoo string `json:"foo" v8n:"+:..foo"`
        SubBar string `json:"bar" v8n:"+:..bar"`
    } `json:"sub"`
}
```

Additional expression functionality notes:
* Boolean operator `^^` (XOr) is also supported
* Path traversal also supports going up to the root object - e.g. `/.foo.bar` will go up to the root object and then descend down path `foo.bar`
* Prefixing property name with `~` (tilde) means check a context condition rather than property existence - e.g. `~METHOD_POST` checks whether the `METHOD_POST` condition token has been set in the context
* The `+:` and `-:` validation tag tokens correspond to the `valix.PropertyValidator.RequiredWith` and `valix.PropertyValidator.UnwantedWith` fields respectively
* Use the `valix.ParseExpression` or `valix.MustParseExpression` functions to programmatically parse boolean property expressions
* Unfortunately, array index notation (e.g. `sub[0]`) is not *currently* supported (and neither is traversing array values)

### Conditional Constraints

Sometimes, the model of JSON requests needs to vary according to some property condition.  For example, the following is a beverage order for a tea and coffee shop:
```go
type BeverageOrder struct {
    Type     string `json:"type" v8n:"notNull,required,&StringValidToken{Tokens:['tea','coffee']}"`
    Quantity int    `json:"quantity" v8n:"notNull,required,&Positive{}"`
    // only relevant to type="tea"...
    Blend    string `json:"blend" v8n:"notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
    // only relevant to type="coffee"...
    Roast    string `json:"roast" v8n:"notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}
```
The above validation will always expect the `blend` and `roast` properties to be present and their value to be valid.  However, this is not the requirement of the model - we only want:
* the `blend` property to be present and valid when `type="tea"`
* the 'roast' property to be present and valid when `type="coffee"`

These validation requirements can be incorporated by using the _'when conditions'_:
```go
type BeverageOrder struct {
    Type     string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{Parent:true}"`
    Quantity int    `json:"quantity" v8n:"notNull,required,&Positive{}"`
    // only relevant to type="tea"...
    Blend    string `json:"blend" v8n:"when:tea,notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
    // only relevant to type="coffee"...
    Roast    string `json:"roast" v8n:"when:coffee,notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}
```
Note in the above:
* on the `Type` field:
  * `order:-1` means that this property is checked
  * `&SetConditionFrom{}` sets a validator context condition token from the incoming value of property `type`<br>
     therefore, either a validator condition token of `tea` or `coffee` will be set<br>
     Note: the `Parent:true` means properties at the same level as this will see the condition token 
* on the `Blend` field the `when:tea` tag has been added - which means the `blend` property is only checked when there is a validator condition token of `tea` set  
* on the `Roast` field the `when:coffee` tag has been added - which means the `roast` property is only checked when there is a validator condition token of `coffee` set

However, this second example may still not be strict enough - because it allows the `blend` property to be present when the type is `"coffee"` and the `roast` property to be present when the type is `"tea"`<br>
This can be overcome by using the _'unwanted conditions'_:
```go
type BeverageOrderStrict struct {
    Type     string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{Parent:true}"`
    Quantity int    `json:"quantity" v8n:"notNull,required,&Positive{}"`
    // Blend is only relevant when type="tea"...
    Blend    string `json:"blend" v8n:"when:tea,unwanted:!tea,notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
    // Roast is only relevant when type="coffee"...
    Roast    string `json:"roast" v8n:"when:coffee,unwanted:!coffee,notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}
```
Note in the above:
* the `unwanted:!tea` tag has been added to the `Blend` field - 
  which means... _"if condition token of `tea` has __not__ been set then we do not want the `blend` property to be present"_
* the `unwanted:!coffee` tag has been added to the `Roast` field - 
  which means... _"if condition token of `coffee` has __not__ been set then we do not want the `roast` property to be present"_

## Polymorphic Validation

Sometimes, the model of JSON requests needs to vary completely according to some condition.
The previous [conditional constraints](#conditional-constraints) example can solve some of the conditional variance - but when different properties are required/not-required or same properties have different constraints under different conditions then polymorphic validation becomes necessary.

As an example, if the following requests all need to be validated using a single validator:
```json
{
    "type": "tea",
    "quantity": 1,
    "blend": "Earl Grey|English Breakfast|Masala Chai"
}
```
```json
{
    "type": "coffee",
    "quantity": 1,
    "roast": "light|medium|dark"
}
```
```json
{
    "type": "soft",
    "quantity": 1,
    "brand": "Coca Cola",
    "flavor": "Regular|Diet|Zero|Cherry"
}
```
```json
{
    "type": "soft",
    "quantity": 1,
    "brand": "Tango",
    "flavor": "Orange|Apple|Strawberry|Watermelon|Tropical"
}
```
A demonstrated solution to this can be see in the [Polymorphic example](https://github.com/marrow16/valix/blob/master/examples/polymorphic_test.go) code

*Note: Polymorphic validators cannot be derived from struct tags (see [Validation Tags](#validation-tags)) - because structs themselves cannot be polymorphic!*

## Validation Tags
Valix can read tags from struct fields when building validators.  These are the `v8n` tags, in the format:
```go
type example struct {
    Field string `v8n:"token[,token, ...]"`
}
```
Where the tokens correspond to various property validation options - as listed here:

<table>
  <thead>
    <tr>
      <th width="33%">Token</th>
      <th>Purpose & Example</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td><code>constraint:constraint-name{fields...}</code></td>
      <td>
        Adds a constraint to the property (this token can be specified multiple times within the <code>v8n</code> tag.
        The <code>constraint-name</code> must be a Valix common constraint or a previously registered constraint.
        The constraint `fields` can optionally be set.
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"constraint:StringMaxLength{Value:255}"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>constraints:[constraint-name{},...]</code></td>
      <td>
        Adds multiple constraints to the property
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"constraints:[StringNotEmpty{},StringNoControlCharacters{}]"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>&amp;constraint-name{fields...}</code></td>
      <td>
        Adds a constraint to the property (shorthand way of specifying constraint without <code>constraint:</code> or <code>constraints:[]</code> prefix)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"&amp;StringMaxLength{Value:255}"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>&amp;[condition,...]constraint-name{fields...}</code></td>
      <td>
        Adds a conditional constraint to the property - the constraint is only checked when the condition(s) are met
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"&amp;[METHOD_POST]StringNotEmpty{}"`
}</pre>
          The <code>StringNotEmpty</code> constraint is only checked when the <code>METHOD_POST</code> condition token has been set
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>&amp;&lt;expr&gt;constraint-name{fields...}</code></td>
      <td>
        Adds a conditional constraint to the property - the constraint is only checked when the <code>expr</code> evaluates to true
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `json:"foo" v8n:"&amp;&lt;(bar && !baz) || (!bar && baz)&gt;StringNotEmpty{}"`
  Bar string `json:"bar" v8n:"optional"`
  Baz string `json:"baz" v8n:"optional"`
}</pre>
          The <code>StringNotEmpty</code> constraint is only checked when only one of the <code>bar</code> or <code>baz</code> properties are present
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>mandatory</code></td>
      <td>
        Specifies the JSON property must be present
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"mandatory"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>mandatory:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>mandatory:[condition,...]</code>
      </td>
      <td>
        Specifies the JSON property must be present under specified conditions
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"mandatory:[METHOD_POST,METHOD_PATCH]"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>required_with:expr</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>+:expr</code>
      </td>
      <td>
        Specifies the JSON property is required according to the presence/non-presence of other properties (as determined by the <code>expr</code>)<br>
        You can also control the violation message used when the property is required but missing using a <code>required_with_msg:</code> or <code>+msg:</code> tag token
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `json:"foo" v8n:"required_with:bar && baz,+msg:'Sometimes foo is required'"`
  Bar string `json:"bar"`
  Baz string `json:"baz"`
}</pre>
          <em>Means the property <code>foo</code> is required when both <code>bar</code> and <code>baz</code> properties are present</em><br>
          <em>Use <code>+msg</code> or <code>required_with_msg</code> to alter the message used when this constraint fails</em>
        </details>
        <br><em>(see also <a href="#requiredunwanted_properties">Required/Unwanted Properties</a> for further notes and examples on expressions)</em>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>unwanted_with:expr</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>-:expr</code>
      </td>
      <td>
        Specifies the JSON property is unwanted according to the presence/non-presence of other properties (as determined by the <code>expr</code>)<br>
        You can also control the violation message used when the property is present but unwanted using a <code>unwanted_with_msg:</code> or <code>-msg:</code> tag token
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `json:"foo" v8n:"unwanted_with:bar || baz,-msg:'Sometimes foo is unwanted'"`
  Bar string `json:"bar"`
  Baz string `json:"baz"`
}</pre>
          <em>Means the property <code>foo</code> is unwanted when either the <code>bar</code> or <code>baz</code> properties are present</em><br>
          <em>Use <code>-msg</code> or <code>unwanted_with_msg</code> to alter the message used when this constraint fails</em>
        </details>
        <br><em>(see also <a href="#requiredunwanted_properties">Required/Unwanted Properties</a> for further notes and examples on expressions)</em>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>notNull</code></td>
      <td>
        Specifies the JSON value for the property cannot be null
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"notNull"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>nullable</code></td>
      <td>
        Specifies the JSON value for the property can be null (<em>opposite of <code>notNull</code></em>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"nullable"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>optional</code></td>
      <td>
        Specifies the JSON property does not have to be present (<em>opposite of <code>mandatory</code></em>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"optional"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>order:n</code></td>
      <td>
        Specifies the order in which the property should be validated (only respected if parent object is tagged as <code>obj.ordered</code> or parent validator is set to <code>OrderedPropertyChecks</code>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"order:0"`
  Bar string `v8n:"order:1"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>only</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>only:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>only:[condition,...]</code>
      </td>
      <td>
        Specifies that the property must not be present with other properties 
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"only,mandatory"`
  Bar string `v8n:"only,mandatory"`
  Baz string `v8n:"only,mandatory"`
}</pre>
          In the above example, a request with only one of the <code>Foo</code>, <code>Bar</code> or <code>Baz</code> will be valid - specifying more than one of those properties would cause a violation.<br>
          <em>Note that even though all three properties are mandatory - if only one of the properties is present, then the other mandatories are ignored.</em><br>
          <em>Use <code>only_msg</code> to alter the message used when this constraint fails</em>   
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>required</code></td>
      <td>
        <em>same as <code>mandatory</code></em>
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"required"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>required:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>required:[condition,...]</code>
      </td>
      <td>
        <em>same as <code>mandatory:condition</code></em>
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"required:[METHOD_POST,METHOD_PATCH]"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>stop_on_first</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>stop1st</code>
      </td>
      <td>
        Specifies that property validation to stop at the first constraint violation found<br>
        <em>Note: This would be the equivalent of setting <code>Stop</code> on each constraint</em>
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"stop_on_first,&StringNotBlank{},&StringNotEmpty{}"`
}</pre>
          <em>In the above example, only one of the specified constraints would fail</em>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>type:type</code>
      </td>
      <td>
        Specifies (overrides) the type expected for the JSON property value<br>
        Where <code>type</code> must be one of (case-insensitive):<br>
        &nbsp;&nbsp;&nbsp;<code>string</code>, <code>number</code>, <code>integer</code>, <code>boolean</code>, <code>object</code>, <code>array</code> or <code>any</code>
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo json.Number `v8n:"type:integer"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>when:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>when:[condition,...]</code>
      </td>
      <td>
        Adds when condition(s) for the property - where <code>condition</code> is a condition token (that may have been set during validation)<br>
        The property is only validated when these conditions are met (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"when:YES_FOO"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>unwanted:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>unwanted:[condition,...]</code>
      </td>
      <td>
        Adds unwanted condition(s) for the property - where <code>condition</code> is a condition token (that may have been set during validation)<br>
        If the unwanted condition(s) is met but the property is present then this is a validation violation (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"unwanted:NO_FOO"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>obj.ignoreUnknownProperties</code></td>
      <td>
        Sets an object (or array of objects) to ignore unknown properties (ignoring unknown properties means that the validator will not fail if an unknown property is found)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.ignoreUnknownProperties"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>obj.unknownProperties:true|false</code></td>
      <td>
        Sets whether an object is to allow/ignore (<code>true</code>) or disallow (<code>false</code>) unknown properties
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.unknownProperties:false"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>obj.constraint:constraint-name{}</code></td>
      <td>
        Sets a constraint on an entire object or array
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.constraint:Length{Minimum:1,Maximum:16}"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>obj.ordered</code></td>
      <td>
        Sets the object validator to check properties in order<br/>
        (same as <code>Validator.OrderedPropertyChecks</code> in <a href="#additional-validator-options">Additional validator options</a>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubObj struct{
    Foo string `v8n:"order:0"`
    Bar string `v8n:"order:1"`
  } `v8n:"obj.ordered"`
}</pre>
the above will check the properties in order specified by their <code>order:</code> - whereas the following will check the properties in alphabetical order of name...
         <pre>type Example struct {
  SubObj struct{
    Foo string `json:"foo"`
    Bar string `json:"bar"`
  } `v8n:"obj.ordered"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>obj.when:condition</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>obj.when:[condition,...]</code>
      </td>
      <td>
        Adds when condition(s) for the object or array - where <code>condition</code> is a condition token (that may have been set during validation)<br>
        The object/array is only validated when these conditions are met (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.when:YES_SUB"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td>
        <code>arr.allowNulls</code><br>
      </td>
      <td>
        For array (slice) fields, specifies that array elements can be null
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  SubSlice []*struct {
    Foo string
  } `json:"subSlice" v8n:"arr.allowNulls"`
}</pre>
        </details>
      </td>
    </tr>
  </tbody>
</table>

#### Registering your own tag tokens

The `v8n` tag can also support custom tokens which can be registered using `valix.RegisterCustomTagToken`.
Any registered custom tag tokens can be used in the `v8n` tag and will be processed when building a validator for a struct using `valix.ValidatorFor`

An example of how this is used can be found in [_examples/custom_tag_tokens_test.go](https://github.com/marrow16/valix/blob/master/examples/custom_tag_tokens_test.go)

#### Tag token aliases

If you find that you're using the same `v8n` tag tokens repeatedly - you can create aliases for these and then just reference the alias using a `$` prefix.

An example of how this is used can be found in [_examples/tag_aliases_test.go](https://github.com/marrow16/valix/blob/master/examples/tag_aliases_test.go)

#### Abbreviating constraints in tags

When specifying constraints in tags, especially with constraint args, the struct tags can become a little verbose.  For example:
```go
type MyStruct struct {
    Foo string `json:"foo" v8n:"&StringNoControlCharacters{},&StringUppercase{Message:'Upper only'},&StringLength{Minimum:10,Maximum:20,ExclusiveMin:true}"`
}
```

To overcome this, there are several things you can do:
1. Where there are no args for the constraint, the `{}` at the end can be dropped
2. Where the constraint struct has only one field or has a default field (tagged with <code>&#96;v8n:"default"&#96;</code>) then the arg name can be dropped
3. The constraints registry has pre-defined abbreviated forms
4. Constraint arg names can be abbreviated or shortened to closest matching name, e.g.
   1. `Message` can be abbreviated to `Msg` or `msg` (case-insensitive, remove vowels, replace double-characters with single)
   2. `Minimum` can be shortened to `Min` or `min` (or any other variation that matches only one target field name)
   3. `ExclusiveMin` can be shortened to `excMin`, `exMin`, `eMin` etc.
5. Where a constraint field is a `bool` and setting it to true - the value _(`:true`)_ can be omitted 

After those steps, the constraint tags would be:
```go
type MyStruct struct {
    Foo string `json:"foo" v8n:"&strnocc,&strupper{'Upper only'},&strlen{stp,min:10,max:20,excMin}"`
}
```

## Internationalisation Support

Valix has full **I18n** support for translating violation messages - which is both extensible and/or replaceable...

#### I18n support features:

* Support for both language and region (e.g. `en`, `en-GB`, `en-US`, `fr` and `fr-CA` etc.)
* Detection of request `Accept-Language` header<br>*(when using `Validator.RequestValidate` or `Validator.RequestValidateInto`)*
* Fallback language and region support
  * e.g. if `fr-CA` was requested but no Canadian specific translation then `fr` is used
  * e.g. if `mt` *(Maltese)* is an unsupported language but you want the fallback language to be `it` *(Italian)* then set this in the `valix.DefaultFallbackLanguages` variable, e.g. `valix.DefaultFallbackLanguages["mt"] = "it"` 
* Default runtime language and region changeable<br>*(set vars `valix.DefaultLanguage` and/or `valix.DefaultRegion`)*
* Built-in `valix.DefaultTranslator` supports English, French, German, Italian and Spanish
  * more languages and regional variants can be added at runtime
  * replace translator with your own (implementing `valix.Translator` interface)
* Completely replaceable I18n support (replace variable `valix.DefaultI18nProvider` with your own) 

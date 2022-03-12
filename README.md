# Valix
[![GoDoc](https://godoc.org/github.com/marrow16/valix?status.svg)](https://pkg.go.dev/github.com/marrow16/valix)
[![codecov](https://codecov.io/gh/marrow16/valix/branch/master/graph/badge.svg)](https://codecov.io/gh/marrow16/valix)
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
* [Validation Tags](#validation-tags) 

## Overview

Validate JSON requests in the form of `*http.Request`, `map[string]interface{}` or `[]interface{}`

## Installation
To install Valix, use go get:

    go get github.com/marrow16/valix

To update Valix to the latest version, run:

    go get -u github.com/marrow16/valix

## Features

* Deep validation
* Create validators from structs or define them as code (see [Creating Validators](#creating-validators))
* Finds all validation violations - not just the first one! (see [Using Validators](#using-validators))
* Rich set of pre-defined common constraints (see [Common Constraints](#common-constraints))
* Customisable constraints (see [Custom Constraints](#custom-constraints))
* Conditional constraints to support polymorphic request models (see [Conditional Constraints](#conditional-constraints))
* 100% tested (see [Codecov.io](https://codecov.io/gh/marrow16/valix))

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

The validator does **_not_** stop on the first problem it finds - it finds all problems and returns them as a list (slice) of 'violations'.  Each violation has a message plus the name and path of the property that failed.
However, custom constraints can be defined that will either stop the entire validation or cease further validation constraints on the current property

Valix comes with a rich set of common constraints (see [Common Constraints](#common-constraints))

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
(Common constraints are defined in `common_contraints.go` and `changing_constraints.go`)

Valix provides a rich set of pre-defined common constraints - listed here for reference:

| Constraint Name                   | Description                                                                                                                                                                                                                               |
|-----------------------------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `valix.ArrayOf`                   | Check each element in an array value is of the correct type                                                                                                                                                                               |
| `valix.DatetimeFuture`            | Check that a datetime/date (represented as string or time.Time) is in the future                                                                                                                                                          |
| `valix.DatetimeFutureOrPresent`   | Check that a datetime/date (represented as string or time.Time) is in the future or present                                                                                                                                               |
| `valix.DatetimePast`              | Check that a datetime/date (represented as string or time.Time) is in the past                                                                                                                                                            |
| `valix.DatetimePastOrPresent`     | Check that a datetime/date (represented as string or time.Time) is in the past or present                                                                                                                                                 |
| `valix.Length`                    | Check that a value (object, array, string) has minimum and maximum length                                                                                                                                                                 |
| `valix.Maximum`                   | Check that a numeric value is less than or equal to a specified maximum                                                                                                                                                                   |
| `valix.Minimum`                   | Check that a numeric value is greater than or equal to a specified minimum                                                                                                                                                                |
| `valix.Negative`                  | Check that a numeric value is negative                                                                                                                                                                                                    |
| `valix.NegativeOrZero`            | Check that a numeric value is negative or zero                                                                                                                                                                                            |
| `valix.Positive`                  | Check that a numeric value is positive (exc. zero)                                                                                                                                                                                        |
| `valix.PositiveOrZero`            | Check that a numeric value is positive or zero                                                                                                                                                                                            |
| `valix.Range`                     | Check that a numeric value is within a specified minimum and maximum range                                                                                                                                                                |
| `valix.SetConditionFrom`          | Is a utility constraint that can be used to set a condition in the `ValidatorContext` from string value of the property to which this constraint is added.<br/>(see example usage in [Conditional Constraints](#conditional-constraints)) |
| `valix.StringCharacters`          | Check that a string contains only allowable characters  (and does not contain any disallowed characters)                                                                                                                                  |
| `valix.StringLength`              | Check that a string has a minimum and maximum length                                                                                                                                                                                      |
| `valix.StringMaxLength`           | Check that a string has a maximum length                                                                                                                                                                                                  |
| `valix.StringMinLength`           | Check that a string has a minimum length                                                                                                                                                                                                  |
| `valix.StringNormalizeUnicode`    | Sets the Unicode normalization of a string value (to the specified form: NFC, NFKC, NFD, NFKD)                                                                                                                                            |
| `valix.StringNotBlank`            | Check that string value is not blank (i.e. that after removing leading and  trailing whitespace the value is not an empty string)                                                                                                         |
| `valix.StringNotEmpty`            | Check that string value is not empty (i.e. not "")                                                                                                                                                                                        |
| `valix.StringNoControlCharacters` | Check that a string does not contain any control characters (i.e. chars < 32)                                                                                                                                                             |
| `valix.StringPattern`             | Check that a string matches a given regexp pattern                                                                                                                                                                                        |
| `valix.StringTrim`                | Trims a string value                                                                                                                                                                                                                      |
| `valix.StringValidCardNumber`     | Check that a string contains a valid card number (according to Luhn Algorithm)                                                                                                                                                            |
| `valix.StringValidEmail`          | Check that a string contains a valid email address                                                                                                                                                                                        |
| `valix.StringValidISODate`        | Check that a string value is a valid ISO8601 Date format (excluding time)                                                                                                                                                                 |
| `valix.StringValidISODatetime`    | Check that a string value is a valid ISO8601 Date/time format                                                                                                                                                                             |
| `valix.StringValidISODate`        | Check that a string value is a valid ISO8601 Date format (excluding time)                                                                                                                                                                 |
| `valix.StringValidToken`          | Check that a string matches one of a pre-defined list of tokens                                                                                                                                                                           |
| `valix.StringValidUuid`           | Check that a string value is a valid UUID (and optionally of a minimum or specified version)                                                                                                                                              |

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
```

### Custom Constraints

If you need a constraint for a specific domain validation, there are two ways to do this...

Create a custom constraint on the fly, example:
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
```

Or create a re-usable constraint (which implements the `valix.Constraint` interface), example:
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
        return !strings.Contains(str, "foo"), c.GetMessage()
    }
    return true, ""
}

func (c *NoFoo) GetMessage() string {
    return "Value must not contain \"foo\""
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
    Type     string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{}"`
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
  * `&SetConditionFrom{}` sets a validator context condition token from the incoming value of property `type` - therefore, either a validator context token of `tea` or `coffee` will be set
* on the `Blend` field the `when:tea` tag has been added - which means the `blend` property is only checked when there is a validator context token of `tea` set  
* on the `Roast` field the `when:coffee` tag has been added - which means the `roast` property is only checked when there is a validator context token of `coffee` set

However, this second example may still not be strict enough - because it allows the `blend` property to be present when the type is `"coffee"` and the `roast` property to be present when the type is `"tea"`<br>
This can be overcome by using the _'unwanted conditions'_:
```go
type BeverageOrderStrict struct {
    Type     string `json:"type" v8n:"notNull,required,order:-1,&StringValidToken{Tokens:['tea','coffee']},&SetConditionFrom{}"`
    Quantity int    `json:"quantity" v8n:"notNull,required,&Positive{}"`
    // Blend is only relevant when type="tea"...
    Blend    string `json:"blend" v8n:"when:tea,unwanted:!tea,notNull,required,&StringValidToken{Tokens:['Earl Grey','English Breakfast','Masala Chai']}"`
    // Roast is only relevant when type="coffee"...
    Roast    string `json:"roast" v8n:"when:coffee,unwanted:!coffee,notNull,required,&StringValidToken{Tokens:['light','medium','dark']}"`
}
```
Note in the above:
* the `unwanted:!tea` tag has been added to the `Blend` field - 
  which means... _"if context token of `tea` has __not__ been set then we do not want the `blend` property to be present"_
* the `unwanted:!coffee` tag has been added to the `Roast` field - 
  which means... _"if context token of `coffee` has __not__ been set then we do not want the `roast` property to be present"_

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
      <td><code>mandatory</code></td>
      <td>
        Specifies the JSON property must be present
        <pre>type Example struct {
  Foo string `v8n:"mandatory"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>notNull</code></td>
      <td>
        Specifies the JSON value for the property cannot be null
        <pre>type Example struct {
  Foo string `v8n:"notNull"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>nullable</code></td>
      <td>
        Specifies the JSON value for the property can be null (<em>opposite of <code>notNull</code></em>)
        <pre>type Example struct {
  Foo string `v8n:"nullable"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>optional</code></td>
      <td>
        Specifies the JSON property does not have to be present (<em>opposite of <code>mandatory</code></em>)
        <pre>type Example struct {
  Foo string `v8n:"optional"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>order:&lt;n&gt;</code></td>
      <td>
        Specifies the order in which the property should be validated (only respected if parent object is tagged as <code>obj.ordered</code> or parent validator is set to <code>OrderedPropertyChecks</code>)
        <pre>type Example struct {
  Foo string `v8n:"order:0"`
  Bar string `v8n:"order:1"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>required</code></td>
      <td>
        <em>same as <code>mandatory</code></em>
        <pre>type Example struct {
  Foo string `v8n:"required"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>type:&lt;type&gt;</code>
      </td>
      <td>
        Specifies (overrides) the type expected for the JSON property value<br>
        Where <code>&lt;type&gt;</code> must be one of (case-insensitive):<br>
        &nbsp;&nbsp;&nbsp;<code>string</code>, <code>number</code>, <code>integer</code>, <code>boolean</code>, <code>object</code>, <code>array</code> or <code>any</code>
        <pre>type Example struct {
  Foo json.Number `v8n:"type:integer"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>constraint:&lt;name&gt;{fields}</code></td>
      <td>
        Adds a constraint to the property (this token can be specified multiple times within the <code>v8n</code> tag.
        The <code>&lt;name&gt;</code> must be a Valix common constraint or a previously registered constraint.
        The constraint `fields` can optionally be set.
        <pre>type Example struct {
  Foo string `v8n:"constraint:StringMaxLength{Value:255}"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>&amp;&lt;constraint-name&gt;{fields}</code></td>
      <td>
        Adds a constraint to the property (shorthand way of specifying constraint without <code>constraint:</code> or <code>constraints:[]</code> prefix)
        <pre>type Example struct {
  Foo string `v8n:"&amp;StringMaxLength{Value:255}"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>constraints:[&lt;name&gt;{},...]</code></td>
      <td>
        Adds multiple constraints to the property
        <pre>type Example struct {
  Foo string `v8n:"constraints:[StringNotEmpty{},StringNoControlCharacters{}]"`
}</pre>
      </td>
    </tr>
    <tr>
      <td>
        <code>when:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>when:[&lt;condition&gt;,...]</code>
      </td>
      <td>
        Adds when condition(s) for the property - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
        The property is only validated when these conditions are met (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <pre>type Example struct {
  Foo string `v8n:"when:YES_FOO"`
}</pre>
      </td>
    </tr>
    <tr>
      <td>
        <code>unwanted:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>unwanted:[&lt;condition&gt;,...]</code>
      </td>
      <td>
        Adds unwanted condition(s) for the property - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
        If the unwanted condition(s) is met but the property is present then this is a validation violation (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <pre>type Example struct {
  Foo string `v8n:"unwanted:NO_FOO"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>obj.ignoreUnknownProperties</code></td>
      <td>
        Sets an object (or array of objects) to ignore unknown properties (ignoring unknown properties means that the validator will not fail if an unknown property is found)
        <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.ignoreUnknownProperties"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>obj.unknownProperties:true|false</code></td>
      <td>
        Sets whether an object is to allow/ignore (<code>true</code>) or disallow (<code>false</code>) unknown properties
        <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.unknownProperties:false"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>obj.constraint:&lt;name&gt;{}</code></td>
      <td>
        Sets a constraint on an entire object or array
        <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.constraint:Length{Minimum:1,Maximum:16}"`
}</pre>
      </td>
    </tr>
    <tr>
      <td><code>obj.ordered</code></td>
      <td>
        Sets the object validator to check properties in order<br/>
        (same as <code>Validator.OrderedPropertyChecks</code> in <a href="#additional-validator-options">Additional validator options</a>)
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
      </td>
    </tr>
    <tr>
      <td>
        <code>obj.when:&lt;token&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>obj.when:[&lt;token&gt;,...]</code>
      </td>
      <td>
        Adds when condition(s) for the object or array - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
        The object/array is only validated when these conditions are met (see <a href="#conditional-constraints">Conditional Constraints</a>)
        <pre>type Example struct {
  SubObj struct{
    Foo string
  } `json:"subObj" v8n:"obj.when:YES_SUB"`
}</pre>
      </td>
    </tr>
  </tbody>
</table>

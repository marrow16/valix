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

Valix provides a rich set of pre-defined common constraints - listed here for reference:

<table>
    <tr>
        <td>
            <code>ArrayConditionalConstraint</code><br>&nbsp;&nbsp;<code>acond</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a special constraint that wraps another constraint - but the wrapped
            constraint is only checked when the specified array condition is met
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>When</code> <em>string</em>
                        </td>
                        <td>
                            is the special token denoting the array condition on which the wrapped constraint is to be checked<br>
                            One of:
                            <ul>
                              <li><code>"first"</code> array item is the first</li>
                              <li><code>"!first"</code> array item is not the first</li>
                              <li><code>"last"</code> array item is the last</li>
                              <li><code>"!last"</code> array item is not the last</li>
                              <li><code>"%n"</code> modulus <em>n</em> of the array index is zero</li>
                              <li><code>"&gt;n"</code> array index is greater than <em>n</em></li>
                              <li><code>"&lt;n"</code> array index is less than <em>n</em></li>
                              <li><code>"n"</code> array index is <em>n</em></li>
                            </ul>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Constraint</code> <em>Constraint</em>
                        </td>
                        <td>
                            is the wrapped constraint
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Ancestry</code> <em>uint</em>
                        </td>
                        <td>
                            is ancestry depth at which to obtain the current array index information<br><br>
                            <em>Note: the ancestry level is only for arrays in the object tree (and does not need to include other levels).
	                        Therefore, by default the value is 0 (zero) - which means the last encountered array</em>
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>ArrayOf</code><br>&nbsp;&nbsp;<code>aof</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check each element in an array value is of the correct type
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Type</code> <em>string</em>
                        </td>
                        <td>
                            the type to check for each item
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowNullElement</code> <em>bool</em>
                        </td>
                        <td>
                            whether to allow null items in the array
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Constraints</code> <em>[]Constraint</em>
                        </td>
                        <td>
                            is an optional slice of constraints tha each array element must satisfy
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>ArrayUnique</code><br>&nbsp;&nbsp;<code>aunique</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check each element in an array value is unique
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>IgnoreNulls</code> <em>bool</em>
                        </td>
                        <td>
                            whether to ignore null items in the array
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>IgnoreCase</code> <em>bool</em>
                        </td>
                        <td>
                            whether uniqueness is case in-insensitive (for string elements)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>ConditionalConstraint</code>
        </td>
        <td>
            Is a special constraint that wraps another constraint - but the wrapped constraint is only checked when the specified when conditions are met
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>When</code> <em>[]string</em>
                        </td>
                        <td>
                            is the condition tokens that determine when the wrapped constraint is checked
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Constraint</code> <em>Constraint</em>
                        </td>
                        <td>
                            is the wrapped constraint
                        </td>
                    </tr>
                </table>
            </details>
            <em>This constraint does not have a <code>v8n</code> tag abbreviation and cannot be used directly in a <code>v8n</code>.  However, all other constraints can be made conditional in <code>v8n</code> tags by prefixing them with <code>[condition,...]</code></em>
            (see <a href="#validation-tags">Validation Tags</a> )
            <details>
              <summary>Example</summary>
              <pre>type Example struct {
  Foo string `v8n:"&[bar,baz]StringNotEmpty{}"`
}</pre>
              <em>Makes the <code>&StringNotEmpty{}</code> constraint only checked when either <code>bar</code> or <code>baz</code> condition tokens have been set</em> 
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeGreaterThan</code><br>&nbsp;&nbsp;<code>dtgt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is greater than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against (a string representation of date or datetime in ISO format)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeGreaterThanOther</code><br>&nbsp;&nbsp;<code>dtgto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is greater than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeGreaterThanOrEqual</code><br>&nbsp;&nbsp;<code>dtgte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is greater than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against (a string representation of date or datetime in ISO format)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeGreaterThanOrEqualOther</code><br>&nbsp;&nbsp;<code>dtgteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is greater than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeLessThan</code><br>&nbsp;&nbsp;<code>dtlt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is less than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against (a string representation of date or datetime in ISO format)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeLessThanOther</code><br>&nbsp;&nbsp;<code>dtlto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is less than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeLessThanOrEqual</code><br>&nbsp;&nbsp;<code>dtlte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is less than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against (a string representation of date or datetime in ISO format)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeLessThanOrEqualOther</code><br>&nbsp;&nbsp;<code>dtlteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value is less than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeFuture</code><br>&nbsp;&nbsp;<code>dtfuture</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a datetime/date (represented as string or time.Time) is in the future
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeFutureOrPresent</code><br>&nbsp;&nbsp;<code>dtfuturep</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a datetime/date (represented as string or time.Time) is in the future or present
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimePast</code><br>&nbsp;&nbsp;<code>dtpast</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a datetime/date (represented as string or time.Time) is in the past
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimePastOrPresent</code><br>&nbsp;&nbsp;<code>dtpastp</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a datetime/date (represented as string or time.Time) is in the past or present
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeTolerance</code><br>&nbsp;&nbsp;<code>dttol</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value meets a tolerance against a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against (a string representation of date or datetime in ISO format)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Duration</code> <em>int64</em>
                        </td>
                        <td>
                            the tolerance duration amount - which can be positive, negative or zero
                            <ul>
                              <li>
                                For negative values, this is the maximum duration into the past
                              </li>
                              <li>
                                For positive values, this is the maximum duration into the future
                              </li>
                              <li>
                                If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	                            specified.  For example, if the <code>Duration</code> is zero and the <code>Unit</code> is specified as "year" then this constraint
	                            will check the same year
                              </li>
                            </ul>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Unit</code> <em>string</em>
                        </td>
                        <td>
                            is the string token specifying the unit in which the <code>Duration</code> is measured<br>
                            The value can - <code>"millennium"</code>, <code>"century"</code>, <code>"decade"</code>, <code>"year"</code>, <code>"month"</code>, <code>"week"</code>, <code>"day"</code>,
	                        <code>"hour"</code>, <code>"min"|"minute"</code>, <code>"sec"|"second"</code>, <code>"milli"|"millisecond"</code>, <code>"micro"|"microsecond"</code> or <code>"nano"|"nanosecond"</code><br>
                            <em>if this is empty, then "day" is assumed.  If the token is invalid - this constraint will fail any comparisons</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>MinCheck</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>IgnoreNull</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, IgnoreNull makes the constraint less strict by ignoring null values
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeToleranceToNow</code><br>&nbsp;&nbsp;<code>dttolnow</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value meets a tolerance against the current time
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Duration</code> <em>int64</em>
                        </td>
                        <td>
                            the tolerance duration amount - which can be positive, negative or zero
                            <ul>
                              <li>
                                For negative values, this is the maximum duration into the past
                              </li>
                              <li>
                                For positive values, this is the maximum duration into the future
                              </li>
                              <li>
                                If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	                            specified.  For example, if the <code>Duration</code> is zero and the <code>Unit</code> is specified as "year" then this constraint
	                            will check the same year
                              </li>
                            </ul>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Unit</code> <em>string</em>
                        </td>
                        <td>
                            is the string token specifying the unit in which the <code>Duration</code> is measured<br>
                            The value can - <code>"millennium"</code>, <code>"century"</code>, <code>"decade"</code>, <code>"year"</code>, <code>"month"</code>, <code>"week"</code>, <code>"day"</code>,
	                        <code>"hour"</code>, <code>"min"|"minute"</code>, <code>"sec"|"second"</code>, <code>"milli"|"millisecond"</code>, <code>"micro"|"microsecond"</code> or <code>"nano"|"nanosecond"</code><br>
                            <em>if this is empty, then "day" is assumed.  If the token is invalid - this constraint will fail any comparisons</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>MinCheck</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>IgnoreNull</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, IgnoreNull makes the constraint less strict by ignoring null values
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>DatetimeToleranceToOther</code><br>&nbsp;&nbsp;<code>dttolother</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a date/time (as an ISO string) value meets a tolerance against the value of another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Duration</code> <em>int64</em>
                        </td>
                        <td>
                            the tolerance duration amount - which can be positive, negative or zero
                            <ul>
                              <li>
                                For negative values, this is the maximum duration into the past
                              </li>
                              <li>
                                For positive values, this is the maximum duration into the future
                              </li>
                              <li>
                                If the value is zero then the behaviour is assumed to be "same" - but is then dependent on the unit
	                            specified.  For example, if the <code>Duration</code> is zero and the <code>Unit</code> is specified as "year" then this constraint
	                            will check the same year
                              </li>
                            </ul>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Unit</code> <em>string</em>
                        </td>
                        <td>
                            is the string token specifying the unit in which the <code>Duration</code> is measured<br>
                            The value can - <code>"millennium"</code>, <code>"century"</code>, <code>"decade"</code>, <code>"year"</code>, <code>"month"</code>, <code>"week"</code>, <code>"day"</code>,
	                        <code>"hour"</code>, <code>"min"|"minute"</code>, <code>"sec"|"second"</code>, <code>"milli"|"millisecond"</code>, <code>"micro"|"microsecond"</code> or <code>"nano"|"nanosecond"</code><br>
                            <em>if this is empty, then "day" is assumed.  If the token is invalid - this constraint will fail any comparisons</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>MinCheck</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, specifies that the tolerance is a minimum check (rather than the default maximum check)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcTime</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, excludes the time when comparing<br>
                            <em>Note: This also excludes the effect of any timezone offsets specified in either of the compared values</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>IgnoreNull</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, IgnoreNull makes the constraint less strict by ignoring null values
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>EqualsOther</code><br>&nbsp;&nbsp;<code>eqo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a property value equals the value of another named property
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>FailingConstraint</code><br>&nbsp;&nbsp;<code>fail</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a utility constraint that always fails
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>StopAll</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, StopAll stops the entire validation
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>FailWhen</code><br>&nbsp;&nbsp;<code>failw</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a utility constraint that fails when specified conditions are met
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Conditions</code> <em>[]string</em>
                        </td>
                        <td>
                            the conditions under which to fail
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>StopAll</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, StopAll stops the entire validation
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>GreaterThan</code><br>&nbsp;&nbsp;<code>gt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is greater than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>GreaterThanOrEqual</code><br>&nbsp;&nbsp;<code>gte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is greater than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>GreaterThanOther</code><br>&nbsp;&nbsp;<code>gto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is greater than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>GreaterThanOrEqualOther</code><br>&nbsp;&nbsp;<code>gteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            check that a numeric value is greater than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Length</code><br>&nbsp;&nbsp;<code>len</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a value (object, array, string) has minimum and maximum length<br>
            <em>(Although this constraint can be used on string properties - it is advised to use <code>StringLength</code> instead)</em>
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Minimum</code> <em>int</em>
                        </td>
                        <td>
                            the minimum length
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Maximum</code> <em>int</em>
                        </td>
                        <td>
                            the maximum length (only checked if this value is > 0)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>string</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>LengthExact</code><br>&nbsp;&nbsp;<code>lenx</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a value (object, array, string) has a specific length<br>
            <em>(Although this constraint can be used on string properties - it is advised to use <code>StringExactLength</code> instead)</em>
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int</em>
                        </td>
                        <td>
                            the length to check
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>LessThan</code><br>&nbsp;&nbsp;<code>lt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is less than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>LessThanOrEqual</code><br>&nbsp;&nbsp;<code>lte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is less than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>LessThanOther</code><br>&nbsp;&nbsp;<code>lto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is less than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>LessThanOrEqualOther</code><br>&nbsp;&nbsp;<code>lteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is less than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Maximum</code><br>&nbsp;&nbsp;<code>max</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is less than or equal to a specified maximum
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the maximum value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>MaximumInt</code><br>&nbsp;&nbsp;<code>maxi</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that an integer value is less than or equal to a specified maximum
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int64</em>
                        </td>
                        <td>
                            the maximum value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Minimum</code><br>&nbsp;&nbsp;<code>min</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is greater than or equal to a specified minimum
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>float64</em>
                        </td>
                        <td>
                            the minimum value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>MinimumInt</code><br>&nbsp;&nbsp;<code>mini</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that an integer value is greater than or equal to a specified minimum
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int64</em>
                        </td>
                        <td>
                            the minimum value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>MultipleOf</code><br>&nbsp;&nbsp;<code>xof</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that an integer value is a multiple of a specific number
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int64</em>
                        </td>
                        <td>
                            the multiple of value to check
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Negative</code><br>&nbsp;&nbsp;<code>neg</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is negative
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NegativeOrZero</code><br>&nbsp;&nbsp;<code>negz</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is negative or zero
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsCIDR</code><br>&nbsp;&nbsp;<code>isCIDR</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid CIDR (v4 or v6) address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>V4Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only CIDR v4
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>V6Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only CIDR v6
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowLoopback</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows loopback addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowPrivate</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows private addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsHostname</code><br>&nbsp;&nbsp;<code>isHostname</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid hostname
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>AllowIPAddress</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowIPV6</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP v6 address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowLocal</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTldOnly</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with only Tld specified (e.g. "audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGeographicTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGenericTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with generic Tlds (e.g. "some.academy")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowBrandTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with brand Tlds (e.g. "my.audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowInfraTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTestTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional country (and geographic) Tlds to allow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of country (and geographic) Tlds to disallow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional generic Tlds to allow (only checked if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of generic Tlds to disallow (only relevant if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional brand Tlds to allow (only checked if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of brand Tlds to disallow (only relevant if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional local Tlds to allow (only checked if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of local Tlds to disallow (only relevant if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsIP</code><br>&nbsp;&nbsp;<code>isIP</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid IP (v4 or v6) address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>V4Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only IP v4
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>V6Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only IP v6
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Resolvable</code> <em>bool</em>
                        </td>
                        <td>
                            if set, checks that the address is resolvable
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowLoopback</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows loopback addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowPrivate</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows private addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowLocalhost</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows value of "localhost" to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsMac</code><br>&nbsp;&nbsp;<code>isMAC</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid MAC address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsTCP</code><br>&nbsp;&nbsp;<code>isTCP</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid resolvable TCP (v4 or v6) address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>V4Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only TCP v4
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>V6Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only TCP v6
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowLoopback</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows loopback addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowPrivate</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows private addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsTld</code><br>&nbsp;&nbsp;<code>isTld</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid Tld (top level domain)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>AllowGeographicTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows geogrpahic tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGenericTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows generic tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowBrandTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows brand tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            additional country (or geographic) tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            excludes specific country (or geographic) tlds being seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            additional generic tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            excludes specific generic tlds being seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            additional brand tlds to be seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            excludes specific brand tlds being seen as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsUDP</code><br>&nbsp;&nbsp;<code>isUDP</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid resolvable UDP (v4 or v6) address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>V4Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only UDP v4
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>V6Only</code> <em>bool</em>
                        </td>
                        <td>
                            if set, allows only UDP v6
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowLoopback</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows loopback addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowPrivate</code> <em>bool</em>
                        </td>
                        <td>
                            if set, disallows private addresses
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsURI</code><br>&nbsp;&nbsp;<code>isURI</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid URI
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>CheckHost</code> <em>bool</em>
                        </td>
                        <td>
                            if set, the host is also checked (see also AllowIPAddress and others)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowIPAddress</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowIPV6</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP v6 address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowLocal</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTldOnly</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with only Tld specified (e.g. "audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGeographicTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGenericTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with generic Tlds (e.g. "some.academy")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowBrandTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with brand Tlds (e.g. "my.audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowInfraTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTestTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional country (and geographic) Tlds to allow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of country (and geographic) Tlds to disallow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional generic Tlds to allow (only checked if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of generic Tlds to disallow (only relevant if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional brand Tlds to allow (only checked if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of brand Tlds to disallow (only relevant if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional local Tlds to allow (only checked if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of local Tlds to disallow (only relevant if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NetIsURL</code><br>&nbsp;&nbsp;<code>isURL</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is a valid URL
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>CheckHost</code> <em>bool</em>
                        </td>
                        <td>
                            if set, the host is also checked (see also AllowIPAddress and others)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowIPAddress</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowIPV6</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows IP v6 address hostnames
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowLocal</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames 'local' (e.g. "localhost", "local", "localdomain", "127.0.0.1", "::1")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTldOnly</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with only Tld specified (e.g. "audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGeographicTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with geographic Tlds (e.g. "some-company.africa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowGenericTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with generic Tlds (e.g. "some.academy")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowBrandTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with brand Tlds (e.g. "my.audi")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowInfraTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with infrastructure Tlds (e.g. "arpa")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTestTlds</code> <em>bool</em>
                        </td>
                        <td>
                            when set, allows hostnames with test Tlds and test domains (e.g. "example.com", "test.com")
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional country (and geographic) Tlds to allow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcCountryCodeTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of country (and geographic) Tlds to disallow
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional generic Tlds to allow (only checked if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcGenericTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of generic Tlds to disallow (only relevant if AllowGenericTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional brand Tlds to allow (only checked if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcBrandTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of brand Tlds to disallow (only relevant if AllowBrandTlds is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AddLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of additional local Tlds to allow (only checked if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExcLocalTlds</code> <em>[]string</em>
                        </td>
                        <td>
                            is an optional slice of local Tlds to disallow (only relevant if AllowLocal is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NotEmpty</code><br>&nbsp;&nbsp;<code>notempty</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a map or slice property value is not empty (has properties or array elements)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>NotEqualsOther</code><br>&nbsp;&nbsp;<code>neqo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a property value not equals the value of another named property
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Positive</code><br>&nbsp;&nbsp;<code>pos</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is positive (exc. zero)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>PositiveOrZero</code><br>&nbsp;&nbsp;<code>posz</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is positive or zero
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>Range</code><br>&nbsp;&nbsp;<code>range</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a numeric value is within a specified minimum and maximum range
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Minimum</code> <em>float64</em>
                        </td>
                        <td>
                            the minimum value of the range
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Maximum</code> <em>float64</em>
                        </td>
                        <td>
                            the maximum value of the range
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>string</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>RangeInt</code><br>&nbsp;&nbsp;<code>rangei</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that an integer value is within a specified minimum and maximum range
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Minimum</code> <em>int64</em>
                        </td>
                        <td>
                            the minimum value of the range
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Maximum</code> <em>int64</em>
                        </td>
                        <td>
                            the maximum value of the range
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>string</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>SetConditionFrom</code><br>&nbsp;&nbsp;<code>cfrom</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a utility constraint that can be used to set a condition in the <code>ValidatorContext</code> from the value of
            the property to which this constraint is added.<br/>
            <em>(see example usage in <a href="#conditional-constraints">Conditional Constraints</a>)</em>
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Parent</code> <em>bool</em>
                        </td>
                        <td>
                            by default, conditions are set on the current property or object - but specifying true for this field means the condition is set on the parent object too
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Global</code> <em>bool</em>
                        </td>
                        <td>
                            setting this field to true means the condition is set for the entire validator context
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Prefix</code> <em>string</em>
                        </td>
                        <td>
                            is any prefix to be appended to the condition token
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Mapping</code> <em>map[string]string</em>
                        </td>
                        <td>
                            converts the string value to alternate values (if the value is not found in the map then the original value is used
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>NullToken</code> <em>string</em>
                        </td>
                        <td>
                            is the condition token used if the value of the property is null/nil.  If this field is not set
                            and the property value is null at validation - then a condition token of "null" is used
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Format</code> <em>string</em>
                        </td>
                        <td>
                            is an optional format string for dealing with non-string property values
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>SetConditionOnType</code><br>&nbsp;&nbsp;<code>ctype</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a utility constraint that can be used to set a condition in the
            <code>ValidatorContext</code> indicating the type of the property value to which this constraint is added.
            <details>
                <summary>Fields</summary>
                <em>None</em>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>SetConditionProperty</code><br>&nbsp;&nbsp;<code>cpty</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Is a utility constraint that can be used to set a condition in the <code>ValidatorContext</code> from the value of a specified property 
            within the object to which this constraint is attached<br/>
            <em>(see example usage in <a href="#polymorphic-validation">Polymorphic Validation</a>)</em>
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the name of the property to extract the condition value from
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Prefix</code> <em>string</em>
                        </td>
                        <td>
                            is any prefix to be appended to the condition token
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Mapping</code> <em>map[string]string</em>
                        </td>
                        <td>
                            converts the token value to alternate values (if the value is not found in the map then the original value is used
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>NullToken</code> <em>string</em>
                        </td>
                        <td>
                            is the condition token used if the value of the property specified is null/nil.  If this field is not set
                            and the property value is null at validation - then a condition token of "null" is used
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>MissingToken</code> <em>string</em>
                        </td>
                        <td>
                            is the condition token used if the property specified is missing.  If this field is not set
                            and the property is missing at validation - then a condition token of "missing" is used
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Format</code> <em>string</em>
                        </td>
                        <td>
                            is an optional format string for dealing with non-string property values
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringCharacters</code><br>&nbsp;&nbsp;<code>strchars</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string contains only allowable characters (and does not contain any disallowed characters)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>AllowRanges</code> <em>[]*unicode.RangeTable</em>
                        </td>
                        <td>
                            the ranges of characters (runes) that are allowed - each character must be in at least one of these
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowRanges</code> <em>[]*unicode.RangeTable</em>
                        </td>
                        <td>
                            the ranges of characters (runes) that are not allowed - if any character is in any of these ranges then the constraint is violated
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringContains</code><br>&nbsp;&nbsp;<code>contains</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string contains with a given value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to check that the string contains
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Values</code> <em>[]string</em>
                        </td>
                        <td>
                            multiple additional values that the string may contain
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is case-insensitive (by default, the check is case-sensitive)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Not</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is NOT-ed (i.e. checks that the string does not contain)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringEndsWith</code><br>&nbsp;&nbsp;<code>ends</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string ends with a given suffix
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to check that the string ends with
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Values</code> <em>[]string</em>
                        </td>
                        <td>
                            multiple additional values that the string may end with
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is case-insensitive (by default, the check is case-sensitive)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Not</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is NOT-ed (i.e. checks that the string does not end with)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLength</code><br>&nbsp;&nbsp;<code>strlen</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has a minimum and maximum length
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Minimum</code> <em>int</em>
                        </td>
                        <td>
                            the minimum length
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Maximum</code> <em>int</em>
                        </td>
                        <td>
                            the maximum length (only checked if this value is > 0)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>string</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>UseRuneLen</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, uses the rune length (true Unicode length) to check length of string
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringExactLength</code><br>&nbsp;&nbsp;<code>strxlen</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has an exact length
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int</em>
                        </td>
                        <td>
                            the exact length expected
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>UseRuneLen</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, uses the rune length (true Unicode length) to check length of string
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringGreaterThan</code><br>&nbsp;&nbsp;<code>strgt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is greater than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringGreaterThanOrEqual</code><br>&nbsp;&nbsp;<code>strgte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is greater than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringGreaterThanOther</code><br>&nbsp;&nbsp;<code>strgto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is greater than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringGreaterThanOrEqualOther</code><br>&nbsp;&nbsp;<code>strgteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is greater than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLessThan</code><br>&nbsp;&nbsp;<code>strlt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is less than a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLessThanOrEqual</code><br>&nbsp;&nbsp;<code>strlte</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is less than or equal to a specified value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to compare against
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLessThanOther</code><br>&nbsp;&nbsp;<code>strlto</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is less than another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLessThanOrEqualOther</code><br>&nbsp;&nbsp;<code>strlteo</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is less than or equal to another named property value
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>PropertyName</code> <em>string</em>
                        </td>
                        <td>
                            the property name of the other value to compare against<br><br>
                            Note: the <code>PropertyName</code> can also be JSON dot notation path - where leading dots allow traversal up
                            the object tree and names, separated by dots, allow traversal down the object tree.<br>
                            A single dot at start is equivalent to no starting dot (i.e. a property name at the same level)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            when set, the comparison is case-insensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringLowercase</code><br>&nbsp;&nbsp;<code>strlower</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has only lowercase letters
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringMaxLength</code><br>&nbsp;&nbsp;<code>strmax</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has a maximum length
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int</em>
                        </td>
                        <td>
                            the maximum length value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMax</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMax specifies the maximum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>UseRuneLen</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, uses the rune length (true Unicode length) to check length of string
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringMinLength</code><br>&nbsp;&nbsp;<code>strmin</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has a minimum length
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>int</em>
                        </td>
                        <td>
                            the minimum length value
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>ExclusiveMin</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, ExclusiveMin specifies the minimum value is exclusive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>UseRuneLen</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, uses the rune length (true Unicode length) to check length of string
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringNotBlank</code><br>&nbsp;&nbsp;<code>strnb</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is not blank (i.e. that after removing leading and trailing whitespace the value is
            not an empty string)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringNotEmpty</code><br>&nbsp;&nbsp;<code>strne</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that string value is not empty (i.e. not "")
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringNoControlCharacters</code><br>&nbsp;&nbsp;<code>strnocc</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string does not contain any control characters (i.e. chars < 32)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringPattern</code><br>&nbsp;&nbsp;<code>strpatt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string matches a given regexp pattern
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Regexp</code> <em>regexp.Regexp</em>
                        </td>
                        <td>
                            the regexp pattern that the string value must match
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringPresetPattern</code><br>&nbsp;&nbsp;<code>strpreset</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string matches a given preset pattern<br>
            There are Preset patterns are defined in the <code>PatternPresets</code> variable (add your own where required) and
            messages for the preset patterns are defined in the <code>PatternPresetMessages</code><br>
            If the preset pattern requires some extra validation beyond the regexp match, then add a checker to the <code>PatternPresetPostPatternChecks</code> variable
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Preset</code> <em>string</em>
                        </td>
                        <td>
                            the preset token (which must exist in the <code>PatternPresets</code> map)<br>
                            <em>If the specified preset token does not exist - the constraint fails!</em>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
            There are 40+ built-in preset patterns (listed below) - and you can add your own using the <code>valix.RegisterPresetPattern()</code> function.
            <details>
              <summary>List of built-in presets</summary>
              <em>Note: Post check indicates whether value is further checked after regexp match - for example, check digits on card numbers & barcodes</em>
              <table>
                  <tr>
                      <th>token</th>
                      <th>message / description</th>
                      <th>post check?</th>
                  </tr>
                  <tr>
                      <td><code>alpha</code></td>
                      <td>Value must be only alphabet characters (A-Z, a-z)</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>alphaNumeric</code></td>
                      <td>Value must be only alphanumeric characters (A-Z, a-z, 0-9)</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>barcode</code></td>
                      <td>
                        Value must be a valid barcode<br>
                        <em>Checks against EAN, ISBN, ISSN, UPC regexps and verifies using check digit</em>
                      </td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>base64</code></td>
                      <td>Value must be a valid base64 encoded string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>base64URL</code></td>
                      <td>Value must be a valid base64 URL encoded string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>card</code></td>
                      <td>Value must be a valid card number</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>cmyk</code></td>
                      <td>
                        Value must be a valid cmyk() color string<br>
                        <em>(components as <code>0.162</code> or <code>16.2%</code>)</em>
                      </td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>cmyk300</code></td>
                      <td>
                        Value must be a valid cmyk() color string (maximum 300%)<br>
                        <em>Post-check ensures components do not exceed 300%</em><br>
                        <em>(components as <code>0.162</code> or <code>16.2%</code>)</em>
                      </td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN</code></td>
                      <td>
                        Value must be a valid EAN code<br>
                        <em>(EAN-8, 13, 14, 18 or 99)</em>
                      </td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN8</code></td>
                      <td>Value must be a valid EAN-8 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN13</code></td>
                      <td>Value must be a valid EAN-13 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>DUN14</code></td>
                      <td>Value must be a valid DUN-14 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN14</code></td>
                      <td>Value must be a valid EAN-14 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN18</code></td>
                      <td>Value must be a valid EAN-18 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>EAN99</code></td>
                      <td>Value must be a valid EAN-99 code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>e164</code></td>
                      <td>Value must be a valid E.164 code</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>hexadecimal</code></td>
                      <td>Value must be a valid hexadecimal string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>hsl</code></td>
                      <td>Value must be a valid hsl() color string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>hsla</code></td>
                      <td>Value must be a valid hsla() color string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>htmlColor</code></td>
                      <td>Value must be a valid HTML color string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>integer</code></td>
                      <td>Value must be a valid integer string (characters 0-9)</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>ISBN</code></td>
                      <td>Value must be a valid ISBN <em>(10 or 13 digit)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ISBN10</code></td>
                      <td>Value must be a valid ISBN <em>(10 digit only)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ISBN13</code></td>
                      <td>Value must be a valid ISBN <em>(13 digit only)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ISSN</code></td>
                      <td>Value must be a valid ISSN <em>(8 or 13 digit)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ISSN13</code></td>
                      <td>Value must be a valid ISSN <em>(13 digit only)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ISSN8</code></td>
                      <td>Value must be a valid ISSN <em>(8 digit only)</em></td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>numeric</code></td>
                      <td>Value must be a valid number string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>numeric+e</code></td>
                      <td>Value must be a valid number string<br><em>also allows scientific notation - e.g. <code>"1.2e34"</code></em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>numeric+x</code></td>
                      <td>Value must be a valid number string<br><em>also allows scientific notation - e.g. <code>"1.2e34"</code>, or <code>"Inf"</code> or <code>"NaN"</code></em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>publication</code></td>
                      <td>Value must be a valid ISBN or ISSN</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>rgb</code></td>
                      <td>Value must be a valid rgb() color string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>rgba</code></td>
                      <td>Value must be a valid rgba() color string</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>rgb-icc</code></td>
                      <td>
                        Value must be a valid rgb-icc() color string<br>
                        <em>Post-check ensures components are correct</em>
                      </td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>ULID</code></td>
                      <td>Value must be a valid ULID</td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UPC</code></td>
                      <td>Value must be a valid UPC code (UPC-A or UPC-E)</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>UPC-A</code></td>
                      <td>Value must be a valid UPC-A code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>UPC-E</code></td>
                      <td>Value must be a valid UPC-E code</td>
                      <td>&#9989;</td>
                  </tr>
                  <tr>
                      <td><code>uuid</code></td>
                      <td>Value must be a valid UUID <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID</code></td>
                      <td>Value must be a valid UUID <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>uuid1</code></td>
                      <td>Value must be a valid UUID (Version 1) <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID1</code></td>
                      <td>Value must be a valid UUID (Version 1) <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>uuid2</code></td>
                      <td>Value must be a valid UUID (Version 2) <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID2</code></td>
                      <td>Value must be a valid UUID (Version 2) <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>uuid3</code></td>
                      <td>Value must be a valid UUID (Version 3) <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID3</code></td>
                      <td>Value must be a valid UUID (Version 3) <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>uuid4</code></td>
                      <td>Value must be a valid UUID (Version 4) <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID4</code></td>
                      <td>Value must be a valid UUID (Version 4) <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>uuid5</code></td>
                      <td>Value must be a valid UUID (Version 5) <em>(lowercase hex chars only)</em></td>
                      <td>&#10005;</td>
                  </tr>
                  <tr>
                      <td><code>UUID5</code></td>
                      <td>Value must be a valid UUID (Version 5) <em>(upper or lower hex chars)</em></td>
                      <td>&#10005;</td>
                  </tr>
              </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringStartsWith</code><br>&nbsp;&nbsp;<code>starts</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string starts with a given prefix
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Value</code> <em>string</em>
                        </td>
                        <td>
                            the value to check that the string starts with
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Values</code> <em>[]string</em>
                        </td>
                        <td>
                            multiple additional values that the string may start with
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>CaseInsensitive</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is case-insensitive (by default, the check is case-sensitive)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Not</code> <em>bool</em>
                        </td>
                        <td>
                            whether the check is NOT-ed (i.e. checks that the string does not start with)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringUppercase</code><br>&nbsp;&nbsp;<code>strupper</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has only uppercase letters
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidCardNumber</code><br>&nbsp;&nbsp;<code>strvcn</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string contains a valid card number (according to Luhn Algorithm)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>AllowSpaces</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, accepts space separators in the card number (but must appear between each 4 digits)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidCountryCode</code><br>&nbsp;&nbsp;<code>strcountry</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string is a valid ISO-3166 (3166-1 / 3166-2) country code
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Allow3166_2</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-2 country and region codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_2_Obsoletes</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true (along with <code>Allow3166_2</code>), allows ISO-3166-2 obsolete region codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowUserAssigned</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 user assigned country codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_1_ExceptionallyReserved</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 exceptionally reserved country codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_1_IndeterminatelyReserved</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 indeterminately reserved country codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_1_TransitionallyReserved</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 transitionally reserved country codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_1_Deleted</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 deleted country codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Allow3166_1_Numeric</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-3166-1 numeric country codes<br>
                            if <code>AllowUserAssigned</code> is also set to true, also allows user assigned contry code (i.e 900 - 999)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>NumericOnly</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, overrides all other flags (with the exception of <code>AllowUserAssigned</code>) and allows only ISO-3166-1 numeric codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidCurrencyCode</code><br>&nbsp;&nbsp;<code>strccy</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string is a valid ISO-4217 currency code
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>AllowNumeric</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-4217 numeric codes
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowHistorical</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows historical ISO-4217 codes (and historical numeric codes - if <code>AllowNumeric</code> is also set)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowTestCode</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-4217 test code (<code>XTS</code> / <code>963</code>) 
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowNoCode</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows ISO-4217 no currency code (<code>XXX</code> / <code>999</code>) 
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowCrypto</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows commonly used crypto currency codes 
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowUnofficial</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows commonly used unofficial currency codes 
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>NumericOnly</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows only ISO-4217 numeric currency codes 
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidEmail</code><br>&nbsp;&nbsp;<code>stremail</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string contains a valid email address
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidISODate</code><br>&nbsp;&nbsp;<code>strisod</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is a valid ISO8601 Date format (excluding time)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidISODatetime</code><br>&nbsp;&nbsp;<code>strisodt</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is a valid ISO8601 Date/time format
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>NoOffset</code> <em>bool</em>
                        </td>
                        <td>
                            specifies, if set to true, that time offsets are not permitted
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>NoMillis</code> <em>bool</em>
                        </td>
                        <td>
                            specifies, if set to true, that seconds cannot have decimal places
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidJson</code><br>&nbsp;&nbsp;<code>strjson</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Checks that a string is valid json
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>DisallowNullJson</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, disallows <code>"null"</code> as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowValue</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, disallows single value (e.g. <code>"true"</code>) as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowArray</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, disallows JSON arrays as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>DisallowObject</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, disallows JSON objects as valid
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidLanguageCode</code><br>&nbsp;&nbsp;<code>strlang</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string is a valid BCP-47 language code
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidToken</code><br>&nbsp;&nbsp;<code>strtoken</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string matches one of a pre-defined list of tokens
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Tokens</code> <em>[]string</em>
                        </td>
                        <td>
                            the set of allowed tokens for the string
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>IgnoreCase</code> <em>bool</em>
                        </td>
                        <td>
                            set to true to make the token check case in-sensitive
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidTimezone</code><br>&nbsp;&nbsp;<code>strtz</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is a valid timezone
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>LocationOnly</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows location only
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>OffsetOnly</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows offset only
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>AllowNumeric</code> <em>bool</em>
                        </td>
                        <td>
                            if set to true, allows offset to be a numeric value (aswell as a string)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidUnicodeNormalization</code><br>&nbsp;&nbsp;<code>struninorm</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string has the correct Unicode normalization form
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>Form</code> <em>norm.Form</em>
                        </td>
                        <td>
                            the normalization form required - i.e. <code>norm.NFC</code>, <code>norm.NFKC</code>, <code>norm.NFD</code> or <code>norm.NFKD</code>
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
    <tr></tr>
    <tr>
        <td>
            <code>StringValidUuid</code><br>&nbsp;&nbsp;<code>struuid</code>&nbsp;<em>(i18n tag abbr.)</em>
        </td>
        <td>
            Check that a string value is a valid UUID (and optionally of a minimum or specified version)
            <details>
                <summary>Fields</summary>
                <table>
                    <tr>
                        <td>
                            <code>MinVersion</code> <em>uint8</em>
                        </td>
                        <td>
                            the minimum UUID version (optional - if zero this is not checked)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>SpecificVersion</code> <em>uint8</em>
                        </td>
                        <td>
                            the specific UUID version (optional - if zero this is not checked)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Message</code> <em>string</em>
                        </td>
                        <td>
                            the violation message to be used if the constraint fails (if empty, the default violation message is used)
                        </td>
                    </tr>
                    <tr>
                        <td>
                            <code>Stop</code> <em>bool</em>
                        </td>
                        <td>
                            when set to true, prevents further validation checks on the property if this constraint fails
                        </td>
                    </tr>
                </table>
            </details>
        </td>
    </tr>
</table>

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
        <code>mandatory:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>mandatory:[&lt;condition&gt;,...]</code>
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
        <code>required_with:&lt;expr&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>+:&lt;expr&gt;</code>
      </td>
      <td>
        Specifies the JSON property is required according to the presence/non-presence of other properties (as determined by the <code>&lt;expr&gt;</code>)<br>
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
        <code>unwanted_with:&lt;expr&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>-:&lt;expr&gt;</code>
      </td>
      <td>
        Specifies the JSON property is required according to the presence/non-presence of other properties (as determined by the <code>&lt;expr&gt;</code>)<br>
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
      <td><code>order:&lt;n&gt;</code></td>
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
        <code>only:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>only:[&lt;condition&gt;,...]</code>
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
        <code>required:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>required:[&lt;condition&gt;,...]</code>
      </td>
      <td>
        <em>same as <code>mandatory:&lt;condition&gt;</code></em>
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
      <td><code>type:&lt;type&gt;</code>
      </td>
      <td>
        Specifies (overrides) the type expected for the JSON property value<br>
        Where <code>&lt;type&gt;</code> must be one of (case-insensitive):<br>
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
      <td><code>constraint:&lt;name&gt;{fields}</code></td>
      <td>
        Adds a constraint to the property (this token can be specified multiple times within the <code>v8n</code> tag.
        The <code>&lt;name&gt;</code> must be a Valix common constraint or a previously registered constraint.
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
      <td><code>&amp;&lt;constraint-name&gt;{fields}</code></td>
      <td>
        Adds a constraint to the property (shorthand way of specifying constraint without <code>constraint:</code> or <code>constraints:[]</code> prefix)
        <details>
          <summary>Example</summary>
          <pre>type Example struct {
  Foo string `v8n:"&amp;StringMaxLength{Value:255}"`
}</pre>
        The constraint can also be made conditional by prefixing the name with the condition tokens, e.g.
        <pre>type Example struct {
  Foo string `v8n:"&amp;[METHOD_POST]StringNotEmpty{}"`
}</pre>
        </details>
      </td>
    </tr>
    <tr></tr>
    <tr>
      <td><code>constraints:[&lt;name&gt;{},...]</code></td>
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
      <td>
        <code>when:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>when:[&lt;condition&gt;,...]</code>
      </td>
      <td>
        Adds when condition(s) for the property - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
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
        <code>unwanted:&lt;condition&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>unwanted:[&lt;condition&gt;,...]</code>
      </td>
      <td>
        Adds unwanted condition(s) for the property - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
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
      <td><code>obj.constraint:&lt;name&gt;{}</code></td>
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
        <code>obj.when:&lt;token&gt;</code><br>
        &nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;<em>or</em><br>
        <code>obj.when:[&lt;token&gt;,...]</code>
      </td>
      <td>
        Adds when condition(s) for the object or array - where <code>&lt;condition&gt;</code> is a condition token (that may have been set during validation)<br>
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

An example of how this is used can be found in [examples/custom_tag_tokens_test.go](https://github.com/marrow16/valix/blob/master/examples/custom_tag_tokens_test.go)

#### Tag token aliases

If you find that you're using the same `v8n` tag tokens repeatedly - you can create aliases for these and then just reference the alias using a `$` prefix.

An example of how this is used can be found in [examples/tag_aliases_test.go](https://github.com/marrow16/valix/blob/master/examples/tag_aliases_test.go)

#### Abbreviating constraints in tags

When specifying constraints in tags, especially with constraint args, the struct tags can become a little verbose.  For example:
```go
type MyStruct struct {
    Foo string `json:"foo" v8n:"&StringNoControlCharacters{},&StringMinLength{Value:10}"`
}
```

To overcome this, there are several things you can do:
1. Where there are no args for the constraint, the `{}` at the end can be dropped
2. Where the constraint struct has only one field or has a default field (tagged with <code>&#96;v8n:"default"&#96;</code>) then the arg name can be dropped
3. The constraints registry has pre-defined abbreviated forms

After those steps, the constraint tags would be:
```go
type MyStruct struct {
    Foo string `json:"foo" v8n:"&strnocc,&strmin{10}"`
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

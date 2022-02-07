# Valix

Valix - Go package for validating requests

## Overview

Validate requests in the form of `*http.Request`, `map[string]interface{}` or `[]interface{}`

### Installation
To install Valix, use go get:

    go get github.com/marrow16/valix

### Staying up to date
To update Valix to the latest version, run:

    go get -u github.com/marrow16/valix

## Examples
Validators can be expressed effectively as, for example:
```go
package main

import (
	"github.com/marrow16/valix"
)

var personValidator = &valix.Validator{
    IgnoreUnknownProperties: false,
    Properties: valix.Properties{
        "name": {
            PropertyType: valix.PropertyType.String,
            NotNull:      true,
            Mandatory:    true,
            Constraints:  valix.Constraints{
                &valix.StringLength{Minimum: 1, Maximum: 255},
            },
        },
        "age": {
            PropertyType: valix.PropertyType.Int,
            NotNull:      true,
            Mandatory:    true,
            Constraints:  valix.Constraints{
                &valix.PositiveOrZero{},
            },
        },
    },
}
```
and then, given a `*http.Request` containing a body of:
```json
{
  "name": "",
  "age": -1
}
```
which is invalid and can easily be validated by, for example:
```go
package main

import "net/http"

func AddPerson(writer http.ResponseWriter, request *http.Request) error {
    ok, violations, obj := personValidator.RequestValidate(request)
    if !ok {
        // write an error response with full info - using violations
        // information 
    }
    // the 'obj' will now be a validated map[string]interface{} 
    return nil
}
```
The above violations would be:

| Index | .Property | .Path | .Message                                               |
|------:|-----------|-------|--------------------------------------------------------|
|     0 | `"name"`  | `""`  | `"Value length must be between 1 and 255 (inclusive)"` |
|     1 | `"age"`   | `""`  | `"Value must be positive or zero"`                     |

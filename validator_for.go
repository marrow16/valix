package valix

import (
	"errors"
	"reflect"
	"strings"
)

const (
	tagNameV8n  = "v8n"
	tagNameJson = "json"
)

type ValidatorForOptions struct {
	IgnoreUnknownProperties bool
	Constraints             Constraints
	AllowNullJson           bool
	UseNumber               bool
}

const (
	errMsgValidatorForStructOnly = "ValidatorFor can only be used with struct arg"
)

// ValidatorFor creates a Validator for a specified struct
//
// If a Validator cannot be compiled for the supplied struct an error is returned
//
// When evaluating the supplied struct to build a Validator, tags on the struct fields
// are used to further clarify the validation constraints.  These are specified using the
// `v8n` tag (see full documentation for details of this tag)
//
// The `json` tag is also used, if specified, to determine the JSON property name to be used.
func ValidatorFor(vstruct interface{}, options *ValidatorForOptions) (*Validator, error) {
	ty := reflect.TypeOf(vstruct)
	if ty.Kind() != reflect.Struct {
		return nil, errors.New(errMsgValidatorForStructOnly)
	}
	result := emptyValidatorFromOptions(options)

	properties, err := buildPropertyValidators(ty)
	if err != nil {
		return nil, err
	}
	result.Properties = properties
	return result, nil
}

// MustCompileValidatorFor creates a Validator for a specified struct
//
// Similar to ValidatorFor but rather than returning an error, it panics if a
// Validator cannot be compiled for the struct
func MustCompileValidatorFor(vstruct interface{}, options *ValidatorForOptions) *Validator {
	v, err := ValidatorFor(vstruct, options)
	if err != nil {
		panic(err)
	}
	return v
}

func emptyValidatorFromOptions(options *ValidatorForOptions) *Validator {
	result := &Validator{
		IgnoreUnknownProperties: false,
		Properties:              Properties{},
		Constraints:             nil,
		AllowArray:              false,
		DisallowObject:          false,
		AllowNullJson:           false,
		UseNumber:               false,
	}
	if options != nil {
		result.IgnoreUnknownProperties = options.IgnoreUnknownProperties
		result.Constraints = options.Constraints
		result.AllowNullJson = options.AllowNullJson
		result.UseNumber = options.UseNumber
	}
	return result
}

func buildPropertyValidators(ty reflect.Type) (Properties, error) {
	result := Properties{}
	cnt := ty.NumField()
	newTy := reflect.New(ty)
	for i := 0; i < cnt; i++ {
		fld := ty.Field(i)
		actualFld := newTy.Elem().FieldByName(fld.Name)
		if actualFld.CanSet() {
			pv, fn, err := propertyValidatorFromField(fld)
			if err != nil {
				return nil, err
			}
			result[fn] = pv
		}
	}
	return result, nil
}

func propertyValidatorFromField(fld reflect.StructField) (*PropertyValidator, string, error) {
	name := getFieldName(fld)
	result := &PropertyValidator{
		Type:        detectFieldType(fld),
		NotNull:     false,
		Mandatory:   false,
		Constraints: Constraints{},
		ObjectValidator: &Validator{
			IgnoreUnknownProperties: false,
			Properties:              Properties{},
			Constraints:             Constraints{},
			AllowArray:              false,
			DisallowObject:          false,
		},
	}
	if err := result.processV8nTag(fld); err != nil {
		return nil, name, err
	}
	if fld.Type.Kind() == reflect.Struct && result.Type == JsonObject {
		ptys, err := buildPropertyValidators(fld.Type)
		if err != nil {
			return nil, name, err
		}
		result.ObjectValidator.DisallowObject = false
		result.ObjectValidator.AllowArray = false
		result.ObjectValidator.Properties = ptys
	} else if fld.Type.Kind() == reflect.Slice && result.Type == JsonArray && fld.Type.Elem().Kind() == reflect.Struct {
		ptys, err := buildPropertyValidators(fld.Type.Elem())
		if err != nil {
			return nil, name, err
		}
		result.ObjectValidator.DisallowObject = true
		result.ObjectValidator.AllowArray = true
		result.ObjectValidator.Properties = ptys
	} else {
		result.ObjectValidator = nil
	}
	return result, name, nil
}

func getFieldName(fld reflect.StructField) string {
	result := fld.Name
	if tag, ok := fld.Tag.Lookup(tagNameJson); ok {
		if cAt := strings.Index(tag, ","); cAt != -1 {
			result = tag[0:cAt]
		} else {
			result = tag
		}
	}
	return result
}

func detectFieldType(fld reflect.StructField) JsonType {
	switch fld.Type.Kind() {
	case reflect.String:
		return JsonString
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return JsonNumber
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return JsonNumber
	case reflect.Float64, reflect.Float32:
		return JsonNumber
	case reflect.Bool:
		return JsonBoolean
	case reflect.Struct:
		return JsonObject
	case reflect.Slice:
		return JsonArray
	case reflect.Map:
		if fld.Type.Key().Kind() == reflect.String && fld.Type.Elem().Kind() == reflect.Interface {
			// only if it's a 'map[string]interface{}` (because that's a representation of a JSON object)...
			//  .Type.Key().Kind() --^       ^-- .Type.Elem().Kind()
			return JsonObject
		}
	}
	return JsonTypeUndefined
}
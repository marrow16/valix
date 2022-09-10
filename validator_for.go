package valix

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

const (
	tagNameV8n   = "v8n"
	tagNameV8nAs = "v8n-as"
	tagNameJson  = "json"
)

// ValidatorForOptions is used by ValidatorFor and MustCompileValidatorFor to set
// the initial overall validator for the struct
type ValidatorForOptions struct {
	// IgnoreUnknownProperties is whether to ignore unknown properties (default false)
	//
	// Set this to `true` if you want to allow unknown properties
	IgnoreUnknownProperties bool
	// Constraints is an optional slice of Constraint items to be checked on the object/array
	//
	// * These are checked in the order specified and prior to property validator & unknown property checks
	Constraints Constraints
	// AllowNullJson forces validator to accept a request body that is null JSON (i.e. a body containing just `null`)
	AllowNullJson bool
	// AllowArray denotes, when true (default is false), that this validator will allow a JSON array - where each
	// item in the array can be validated as an object
	AllowArray bool
	// DisallowObject denotes, when set to true, that this validator will disallow JSON objects - i.e. that it
	// expects JSON arrays (in which case the AllowArray should also be set to true)
	DisallowObject bool
	// StopOnFirst if set, instructs the validator to stop at the first violation found
	StopOnFirst bool
	// UseNumber forces RequestValidate method to use json.Number when decoding request body
	UseNumber bool
	// OrderedPropertyChecks determines whether properties should be checked in order - when set to true, properties
	// are sorted by PropertyValidator.Order and property name
	//
	// When this is set to false (default) properties are checked in the order in which they appear in the properties map -
	// which is unpredictable
	OrderedPropertyChecks bool
	// OasInfo is additional information (for OpenAPI Specification)
	OasInfo *OasInfo
}

const (
	errMsgValidatorForStructOnly   = "ValidatorFor can only be used with struct arg"
	errMsgCannotFindPropertyInRepo = tagNameV8nAs + " tag cannot find property name '%s' in properties repository"
	errMsgIncompatiblePropertyType = tagNameV8nAs + " tag has incompatible field type for property name '%s'"
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
		StopOnFirst:             false,
	}
	if options != nil {
		result.IgnoreUnknownProperties = options.IgnoreUnknownProperties
		result.Constraints = options.Constraints
		result.AllowNullJson = options.AllowNullJson
		result.UseNumber = options.UseNumber
		result.AllowArray = options.AllowArray
		result.DisallowObject = options.DisallowObject
		result.StopOnFirst = options.StopOnFirst
		result.OrderedPropertyChecks = options.OrderedPropertyChecks
		result.OasInfo = options.OasInfo
	}
	return result
}

func buildPropertyValidators(ty reflect.Type) (Properties, error) {
	cnt := ty.NumField()
	result := make(Properties, cnt)
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
	result, err := initialPropertyValidator(fld, name)
	if err != nil {
		return nil, name, err
	}
	if err := result.processV8nTag(fld.Name, name, fld); err != nil {
		return nil, name, err
	}
	if err := customTagsRegistry.processField(fld, result); err != nil {
		return nil, name, err
	}
	if err := result.processOasTag(fld); err != nil {
		return nil, name, err
	}

	fKind := fld.Type.Kind()
	objValidatorUsed := false
	if result.Type == JsonObject {
		if fKind == reflect.Struct {
			if ptys, err := buildPropertyValidators(fld.Type); err != nil {
				return nil, name, err
			} else if len(ptys) > 0 && result.ObjectValidator != nil {
				result.ObjectValidator.DisallowObject = false
				result.ObjectValidator.AllowArray = false
				result.ObjectValidator.Properties = ptys
				objValidatorUsed = true
			}
		} else if fKind == reflect.Ptr && fld.Type.Elem().Kind() == reflect.Struct {
			if ptys, err := buildPropertyValidators(fld.Type.Elem()); err != nil {
				return nil, name, err
			} else if len(ptys) > 0 && result.ObjectValidator != nil {
				result.ObjectValidator.DisallowObject = false
				result.ObjectValidator.AllowArray = false
				result.ObjectValidator.Properties = ptys
				objValidatorUsed = true
			}
		}
	} else if result.Type == JsonArray && fKind == reflect.Slice {
		if fld.Type.Elem().Kind() == reflect.Struct {
			if ptys, err := buildPropertyValidators(fld.Type.Elem()); err != nil {
				return nil, name, err
			} else if len(ptys) > 0 && result.ObjectValidator != nil {
				result.ObjectValidator.DisallowObject = true
				result.ObjectValidator.AllowArray = true
				result.ObjectValidator.Properties = ptys
				objValidatorUsed = true
			}
		} else if fld.Type.Elem().Kind() == reflect.Ptr && fld.Type.Elem().Elem().Kind() == reflect.Struct {
			if ptys, err := buildPropertyValidators(fld.Type.Elem().Elem()); err != nil {
				return nil, name, err
			} else if len(ptys) > 0 && result.ObjectValidator != nil {
				result.ObjectValidator.DisallowObject = true
				result.ObjectValidator.AllowArray = true
				result.ObjectValidator.Properties = ptys
				objValidatorUsed = true
			}
		}
	}
	if !objValidatorUsed {
		result.ObjectValidator = nil
	}
	return result, name, nil
}

func initialPropertyValidator(fld reflect.StructField, name string) (*PropertyValidator, error) {
	jsonFldType := detectFieldType(fld)
	if tag, ok := fld.Tag.Lookup(tagNameV8nAs); ok {
		asName := ternary(tag == "").string(name, tag)
		result := propertiesRepo.getNamed(asName)
		if result == nil {
			return nil, fmt.Errorf(errMsgCannotFindPropertyInRepo, asName)
		}
		result = result.Clone()
		// check types are compatible...
		if result.Type != JsonAny && result.Type != jsonFldType {
			return nil, fmt.Errorf(errMsgIncompatiblePropertyType, asName)
		}
		result.Type = jsonFldType
		// make sure constraints and object validator are initialised...
		if result.Constraints == nil {
			result.Constraints = Constraints{}
		}
		if result.ObjectValidator == nil {
			result.ObjectValidator = &Validator{
				IgnoreUnknownProperties: false,
				Properties:              Properties{},
				Constraints:             Constraints{},
				AllowArray:              false,
				DisallowObject:          false,
			}
		}
		return result, nil
	}
	return &PropertyValidator{
		Type:        jsonFldType,
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
	}, nil
}

func getFieldName(fld reflect.StructField) string {
	result := fld.Name
	if rn, ok := getDefaultPropertyNameProvider().NameFor(fld); ok {
		result = rn
	}
	return result
}

var timeType = reflect.TypeOf(time.Time{})
var valixTimeType = reflect.TypeOf(Time{})

func detectFieldType(fld reflect.StructField) JsonType {
	k := fld.Type.Kind()
	isPtr := k == reflect.Ptr
	if isPtr {
		k = fld.Type.Elem().Kind()
	}
	switch k {
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
		if (!isPtr && fld.Type.AssignableTo(timeType)) || (isPtr && fld.Type.Elem().AssignableTo(timeType)) ||
			(!isPtr && fld.Type.AssignableTo(valixTimeType)) || (isPtr && fld.Type.Elem().AssignableTo(valixTimeType)) {
			return JsonDatetime
		}
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
	return JsonAny
}

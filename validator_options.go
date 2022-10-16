package valix

type Option interface {
	Apply(on *Validator) error
}

// ValidatorForOptions is used by ValidatorFor and MustCompileValidatorFor to set
// the initial overall validator for the struct
//
// Note: This struct is not deprecated, but is retained for backward compatibility - it is
// now easier to use individual Option's when using ValidatorFor()
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

func (o *ValidatorForOptions) Apply(on *Validator) error {
	on.IgnoreUnknownProperties = o.IgnoreUnknownProperties
	on.Constraints = append(on.Constraints, o.Constraints...)
	on.AllowNullJson = o.AllowNullJson
	on.UseNumber = o.UseNumber
	on.AllowArray = o.AllowArray
	on.DisallowObject = o.DisallowObject
	on.StopOnFirst = o.StopOnFirst
	on.OrderedPropertyChecks = o.OrderedPropertyChecks
	if o.OasInfo != nil {
		on.OasInfo = o.OasInfo
	}
	return nil
}

// OptionIgnoreUnknownProperties option for ValidatorFor - sets Validator to ignore unknown properties
var OptionIgnoreUnknownProperties Option = &optionIgnoreUnknownProperties{true}

// OptionDisallowUnknownProperties option for ValidatorFor - sets Validator to not ignore unknown properties
var OptionDisallowUnknownProperties Option = &optionIgnoreUnknownProperties{false}

type optionIgnoreUnknownProperties struct {
	setting bool
}

func (o *optionIgnoreUnknownProperties) Apply(on *Validator) error {
	on.IgnoreUnknownProperties = o.setting
	return nil
}

// OptionConstraints option for ValidatorFor - adds constraints to the Validator
var OptionConstraints = func(constraints ...Constraint) Option {
	return &optionConstraints{
		constraints: constraints,
	}
}

type optionConstraints struct {
	constraints Constraints
}

func (o *optionConstraints) Apply(on *Validator) error {
	on.Constraints = append(on.Constraints, o.constraints...)
	return nil
}

// OptionAllowNullJson option for ValidatorFor - sets Validator to allow JSON null
var OptionAllowNullJson Option = &optionAllowNullJson{true}

// OptionDisallowNullJson option for ValidatorFor - sets Validator to not allow JSON null
var OptionDisallowNullJson Option = &optionAllowNullJson{false}

type optionAllowNullJson struct {
	setting bool
}

func (o *optionAllowNullJson) Apply(on *Validator) error {
	on.AllowNullJson = o.setting
	return nil
}

// OptionAllowArray option for ValidatorFor - sets Validator to allow JSON array
var OptionAllowArray Option = &optionAllowArray{true}

// OptionDisallowArray option for ValidatorFor - sets Validator to not allow JSON array
var OptionDisallowArray Option = &optionAllowArray{false}

type optionAllowArray struct {
	setting bool
}

func (o *optionAllowArray) Apply(on *Validator) error {
	on.AllowArray = o.setting
	return nil
}

// OptionAllowObject option for ValidatorFor - sets Validator to allow JSON object
var OptionAllowObject Option = &optionDisallowObject{false}

// OptionDisallowObject option for ValidatorFor - sets Validator to not allow JSON object
var OptionDisallowObject Option = &optionDisallowObject{true}

type optionDisallowObject struct {
	setting bool
}

func (o *optionDisallowObject) Apply(on *Validator) error {
	on.DisallowObject = o.setting
	return nil
}

// OptionStopOnFirst option for ValidatorFor - sets Validator to stop on first violation
var OptionStopOnFirst Option = &optionStopOnFirst{true}

// OptionDontStopOnFirst option for ValidatorFor - sets Validator to not stop on first violation
var OptionDontStopOnFirst Option = &optionStopOnFirst{false}

type optionStopOnFirst struct {
	setting bool
}

func (o *optionStopOnFirst) Apply(on *Validator) error {
	on.StopOnFirst = o.setting
	return nil
}

// OptionUseNumber option for ValidatorFor - sets Validator to use json.Number when decoding requests
var OptionUseNumber Option = &optionUseNumber{true}

type optionUseNumber struct {
	setting bool
}

func (o *optionUseNumber) Apply(on *Validator) error {
	on.UseNumber = o.setting
	return nil
}

// OptionOrderedPropertyChecks option for ValidatorFor - sets Validator to do ordered property checks
var OptionOrderedPropertyChecks Option = &optionOrderedPropertyChecks{true}

// OptionUnOrderedPropertyChecks option for ValidatorFor - sets Validator to not do ordered property checks
var OptionUnOrderedPropertyChecks Option = &optionOrderedPropertyChecks{false}

type optionOrderedPropertyChecks struct {
	setting bool
}

func (o *optionOrderedPropertyChecks) Apply(on *Validator) error {
	on.OrderedPropertyChecks = o.setting
	return nil
}
package valix

import (
	"fmt"
	"reflect"
	"sync"
)

const (
	panicMsgConstraintExists = "constraint \"%s\" already exists in registry"
)

type constraintRegistry struct {
	namedConstraints map[string]Constraint
	sync             *sync.Mutex
}

var constraintsRegistry constraintRegistry

func init() {
	constraintsRegistry = constraintRegistry{
		namedConstraints: defaultConstraints(),
		sync:             &sync.Mutex{},
	}
}

func defaultConstraints() map[string]Constraint {
	return map[string]Constraint{
		"ArrayConditionalConstraint":      &ArrayConditionalConstraint{},
		"ArrayOf":                         &ArrayOf{},
		"ArrayUnique":                     &ArrayUnique{},
		"DatetimeGreaterThan":             &DatetimeGreaterThan{},
		"DatetimeGreaterThanOther":        &DatetimeGreaterThanOther{},
		"DatetimeGreaterThanOrEqual":      &DatetimeGreaterThanOrEqual{},
		"DatetimeGreaterThanOrEqualOther": &DatetimeGreaterThanOrEqualOther{},
		"DatetimeFuture":                  &DatetimeFuture{},
		"DatetimeFutureOrPresent":         &DatetimeFutureOrPresent{},
		"DatetimeLessThan":                &DatetimeLessThan{},
		"DatetimeLessThanOther":           &DatetimeLessThanOther{},
		"DatetimeLessThanOrEqual":         &DatetimeLessThanOrEqual{},
		"DatetimeLessThanOrEqualOther":    &DatetimeLessThanOrEqualOther{},
		"DatetimePast":                    &DatetimePast{},
		"DatetimePastOrPresent":           &DatetimePastOrPresent{},
		"DatetimeTolerance":               &DatetimeTolerance{},
		"DatetimeToleranceToNow":          &DatetimeToleranceToNow{},
		"DatetimeToleranceToOther":        &DatetimeToleranceToOther{},
		"EqualsOther":                     &EqualsOther{},
		"FailingConstraint":               &FailingConstraint{},
		"FailWhen":                        &FailWhen{},
		"GreaterThan":                     &GreaterThan{},
		"GreaterThanOrEqual":              &GreaterThanOrEqual{},
		"GreaterThanOther":                &GreaterThanOther{},
		"GreaterThanOrEqualOther":         &GreaterThanOrEqualOther{},
		"Length":                          &Length{},
		"LengthExact":                     &LengthExact{},
		"LessThan":                        &LessThan{},
		"LessThanOrEqual":                 &LessThanOrEqual{},
		"LessThanOther":                   &LessThanOther{},
		"LessThanOrEqualOther":            &LessThanOrEqualOther{},
		"Maximum":                         &Maximum{},
		"MaximumInt":                      &MaximumInt{},
		"Minimum":                         &Minimum{},
		"MinimumInt":                      &MinimumInt{},
		"MultipleOf":                      &MultipleOf{},
		"Negative":                        &Negative{},
		"NegativeOrZero":                  &NegativeOrZero{},
		"NetIsCIDR":                       &NetIsCIDR{},
		"NetIsHostname":                   &NetIsHostname{},
		"NetIsIP":                         &NetIsIP{},
		"NetIsMac":                        &NetIsMac{},
		"NetIsTCP":                        &NetIsTCP{},
		"NetIsTld":                        &NetIsTld{},
		"NetIsUDP":                        &NetIsUDP{},
		"NetIsURI":                        &NetIsURI{},
		"NetIsURL":                        &NetIsURL{},
		"NotEmpty":                        &NotEmpty{},
		"NotEqualsOther":                  &NotEqualsOther{},
		"Positive":                        &Positive{},
		"PositiveOrZero":                  &PositiveOrZero{},
		"Range":                           &Range{},
		"RangeInt":                        &RangeInt{},
		"SetConditionFrom":                &SetConditionFrom{},
		"SetConditionOnType":              &SetConditionOnType{},
		"SetConditionProperty":            &SetConditionProperty{},
		"StringCharacters":                &StringCharacters{},
		"StringContains":                  &StringContains{},
		"StringEndsWith":                  &StringEndsWith{},
		"StringExactLength":               &StringExactLength{},
		"StringGreaterThan":               &StringGreaterThan{},
		"StringGreaterThanOrEqual":        &StringGreaterThanOrEqual{},
		"StringGreaterThanOther":          &StringGreaterThanOther{},
		"StringGreaterThanOrEqualOther":   &StringGreaterThanOrEqualOther{},
		"StringLessThan":                  &StringLessThan{},
		"StringLessThanOrEqual":           &StringLessThanOrEqual{},
		"StringLessThanOther":             &StringLessThanOther{},
		"StringLessThanOrEqualOther":      &StringLessThanOrEqualOther{},
		"StringNoControlCharacters":       &StringNoControlCharacters{},
		"StringNotBlank":                  &StringNotBlank{},
		"StringNotEmpty":                  &StringNotEmpty{},
		"StringLength":                    &StringLength{},
		"StringLowercase":                 &StringLowercase{},
		"StringMaxLength":                 &StringMaxLength{},
		"StringMinLength":                 &StringMinLength{},
		"StringPattern":                   &StringPattern{},
		"StringPresetPattern":             &StringPresetPattern{},
		"StringStartsWith":                &StringStartsWith{},
		"StringUppercase":                 &StringUppercase{},
		"StringValidCardNumber":           &StringValidCardNumber{},
		"StringValidCountryCode":          &StringValidCountryCode{},
		"StringValidCurrencyCode":         &StringValidCurrencyCode{},
		"StringValidEmail":                &StringValidEmail{},
		"StringValidLanguageCode":         &StringValidLanguageCode{},
		"StringValidISODatetime":          &StringValidISODatetime{},
		"StringValidISODate":              &StringValidISODate{},
		"StringValidJson":                 &StringValidJson{},
		"StringValidTimezone":             &StringValidTimezone{},
		"StringValidToken":                &StringValidToken{},
		"StringValidUnicodeNormalization": &StringValidUnicodeNormalization{},
		"StringValidUuid":                 &StringValidUuid{},
		// abbreviations...
		"aof":        &ArrayOf{},
		"acond":      &ArrayConditionalConstraint{},
		"aunique":    &ArrayUnique{},
		"contains":   &StringContains{},
		"dtgt":       &DatetimeGreaterThan{},
		"dtgto":      &DatetimeGreaterThanOther{},
		"dtgte":      &DatetimeGreaterThanOrEqual{},
		"dtgteo":     &DatetimeGreaterThanOrEqualOther{},
		"dtfuture":   &DatetimeFuture{},
		"dtfuturep":  &DatetimeFutureOrPresent{},
		"dtlt":       &DatetimeLessThan{},
		"dtlto":      &DatetimeLessThanOther{},
		"dtlte":      &DatetimeLessThanOrEqual{},
		"dtlteo":     &DatetimeLessThanOrEqualOther{},
		"dtpast":     &DatetimePast{},
		"dtpastp":    &DatetimePastOrPresent{},
		"dttol":      &DatetimeTolerance{},
		"dttolnow":   &DatetimeToleranceToNow{},
		"dttolother": &DatetimeToleranceToOther{},
		"ends":       &StringEndsWith{},
		"eqo":        &EqualsOther{},
		"fail":       &FailingConstraint{},
		"failw":      &FailWhen{},
		"gt":         &GreaterThan{},
		"gte":        &GreaterThanOrEqual{},
		"gto":        &GreaterThanOther{},
		"gteo":       &GreaterThanOrEqualOther{},
		"len":        &Length{},
		"lenx":       &LengthExact{},
		"lt":         &LessThan{},
		"lte":        &LessThanOrEqual{},
		"lto":        &LessThanOther{},
		"lteo":       &LessThanOrEqualOther{},
		"max":        &Maximum{},
		"maxi":       &MaximumInt{},
		"min":        &Minimum{},
		"mini":       &MinimumInt{},
		"notempty":   &NotEmpty{},
		"xof":        &MultipleOf{},
		"neg":        &Negative{},
		"negz":       &NegativeOrZero{},
		"neqo":       &NotEqualsOther{},
		"pos":        &Positive{},
		"posz":       &PositiveOrZero{},
		"range":      &Range{},
		"rangei":     &RangeInt{},
		"cfrom":      &SetConditionFrom{},
		"ctype":      &SetConditionOnType{},
		"cpty":       &SetConditionProperty{},
		"starts":     &StringStartsWith{},
		"strchars":   &StringCharacters{},
		"strxlen":    &StringExactLength{},
		"strgt":      &StringGreaterThan{},
		"strgte":     &StringGreaterThanOrEqual{},
		"strgto":     &StringGreaterThanOther{},
		"strgteo":    &StringGreaterThanOrEqualOther{},
		"strlt":      &StringLessThan{},
		"strlte":     &StringLessThanOrEqual{},
		"strlto":     &StringLessThanOther{},
		"strlteo":    &StringLessThanOrEqualOther{},
		"strnocc":    &StringNoControlCharacters{},
		"strnb":      &StringNotBlank{},
		"strne":      &StringNotEmpty{},
		"strlen":     &StringLength{},
		"strlower":   &StringLowercase{},
		"strmax":     &StringMaxLength{},
		"strmin":     &StringMinLength{},
		"strpatt":    &StringPattern{},
		"strpreset":  &StringPresetPattern{},
		"strupper":   &StringUppercase{},
		"stremail":   &StringValidEmail{},
		"strcountry": &StringValidCountryCode{},
		"strccy":     &StringValidCurrencyCode{},
		"strlang":    &StringValidLanguageCode{},
		"strvcn":     &StringValidCardNumber{},
		"strisodt":   &StringValidISODatetime{},
		"strisod":    &StringValidISODate{},
		"strjson":    &StringValidJson{},
		"strtz":      &StringValidTimezone{},
		"strtoken":   &StringValidToken{},
		"struninorm": &StringValidUnicodeNormalization{},
		"struuid":    &StringValidUuid{},
		"isCIDR":     &NetIsCIDR{},
		"isHostname": &NetIsHostname{},
		"isIP":       &NetIsIP{},
		"isMac":      &NetIsMac{},
		"isTCP":      &NetIsTCP{},
		"isTld":      &NetIsTld{},
		"isUDP":      &NetIsUDP{},
		"isURI":      &NetIsURI{},
		"isURL":      &NetIsURL{},
		// special abbreviations...
		"iso3166":           &StringValidCountryCode{Allow3166_1_Numeric: true, Allow3166_2: true},
		"iso3166-1":         &StringValidCountryCode{},
		"iso3166-1-numeric": &StringValidCountryCode{NumericOnly: true},
		"iso3166-2":         &StringValidCountryCode{Allow3166_2: true},
		"iso4217":           &StringValidCurrencyCode{AllowNumeric: true},
		"iso4217-alpha":     &StringValidCurrencyCode{},
		"iso4217-numeric":   &StringValidCurrencyCode{NumericOnly: true},
		"lang":              &StringValidLanguageCode{},
		// preset patterns...
		presetTokenAlpha:        &StringPresetPattern{Preset: presetTokenAlpha},
		presetTokenAlphaNumeric: &StringPresetPattern{Preset: presetTokenAlphaNumeric},
		presetTokenBarcode:      &StringPresetPattern{Preset: presetTokenBarcode},
		presetTokenBase64:       &StringPresetPattern{Preset: presetTokenBase64},
		presetTokenBase64URL:    &StringPresetPattern{Preset: presetTokenBase64URL},
		presetTokenCard:         &StringPresetPattern{Preset: presetTokenCard},
		presetTokenCMYK:         &StringPresetPattern{Preset: presetTokenCMYK},
		presetTokenCMYK300:      &StringPresetPattern{Preset: presetTokenCMYK300},
		presetTokenE164:         &StringPresetPattern{Preset: presetTokenE164},
		presetTokenEAN:          &StringPresetPattern{Preset: presetTokenEAN},
		presetTokenEAN8:         &StringPresetPattern{Preset: presetTokenEAN8},
		presetTokenEAN13:        &StringPresetPattern{Preset: presetTokenEAN13},
		presetTokenDUN14:        &StringPresetPattern{Preset: presetTokenDUN14},
		presetTokenEAN14:        &StringPresetPattern{Preset: presetTokenEAN14},
		presetTokenEAN18:        &StringPresetPattern{Preset: presetTokenEAN18},
		presetTokenEAN99:        &StringPresetPattern{Preset: presetTokenEAN99},
		presetTokenHexadecimal:  &StringPresetPattern{Preset: presetTokenHexadecimal},
		presetTokenHsl:          &StringPresetPattern{Preset: presetTokenHsl},
		presetTokenHsla:         &StringPresetPattern{Preset: presetTokenHsla},
		presetTokenHtmlColor:    &StringPresetPattern{Preset: presetTokenHtmlColor},
		presetTokenInteger:      &StringPresetPattern{Preset: presetTokenInteger},
		presetTokenISBN:         &StringPresetPattern{Preset: presetTokenISBN},
		presetTokenISBN10:       &StringPresetPattern{Preset: presetTokenISBN10},
		presetTokenISBN13:       &StringPresetPattern{Preset: presetTokenISBN13},
		presetTokenISSN:         &StringPresetPattern{Preset: presetTokenISSN},
		presetTokenISSN8:        &StringPresetPattern{Preset: presetTokenISSN8},
		presetTokenISSN13:       &StringPresetPattern{Preset: presetTokenISSN13},
		presetTokenNumeric:      &StringPresetPattern{Preset: presetTokenNumeric},
		presetTokenNumericE:     &StringPresetPattern{Preset: presetTokenNumericE},
		presetTokenNumericX:     &StringPresetPattern{Preset: presetTokenNumericX},
		presetTokenPublication:  &StringPresetPattern{Preset: presetTokenPublication},
		presetTokenRgb:          &StringPresetPattern{Preset: presetTokenRgb},
		presetTokenRgba:         &StringPresetPattern{Preset: presetTokenRgba},
		presetTokenRgbIcc:       &StringPresetPattern{Preset: presetTokenRgbIcc},
		presetTokenULID:         &StringPresetPattern{Preset: presetTokenULID},
		presetTokenUPC:          &StringPresetPattern{Preset: presetTokenUPC},
		presetTokenUPCA:         &StringPresetPattern{Preset: presetTokenUPCA},
		presetTokenUPCE:         &StringPresetPattern{Preset: presetTokenUPCE},
		presetTokenUuid:         &StringPresetPattern{Preset: presetTokenUuid},
		presetTokenUUID:         &StringPresetPattern{Preset: presetTokenUUID},
		presetTokenUuid1:        &StringPresetPattern{Preset: presetTokenUuid1},
		presetTokenUUID1:        &StringPresetPattern{Preset: presetTokenUUID1},
		presetTokenUuid2:        &StringPresetPattern{Preset: presetTokenUuid2},
		presetTokenUUID2:        &StringPresetPattern{Preset: presetTokenUUID2},
		presetTokenUuid3:        &StringPresetPattern{Preset: presetTokenUuid3},
		presetTokenUUID3:        &StringPresetPattern{Preset: presetTokenUUID3},
		presetTokenUuid4:        &StringPresetPattern{Preset: presetTokenUuid4},
		presetTokenUUID4:        &StringPresetPattern{Preset: presetTokenUUID4},
		presetTokenUuid5:        &StringPresetPattern{Preset: presetTokenUuid5},
		presetTokenUUID5:        &StringPresetPattern{Preset: presetTokenUUID5},
	}
}

func (r *constraintRegistry) checkOverwriteAllowed(overwrite bool, name string) {
	if !overwrite {
		if _, ok := r.namedConstraints[name]; ok {
			panic(fmt.Errorf(panicMsgConstraintExists, name))
		}
	}
}

func (r *constraintRegistry) register(overwrite bool, constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	name := reflect.TypeOf(constraint).Elem().Name()
	r.checkOverwriteAllowed(overwrite, name)
	r.namedConstraints[name] = constraint
}

func (r *constraintRegistry) registerMany(overwrite bool, constraints ...Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for _, constraint := range constraints {
		name := reflect.TypeOf(constraint).Elem().Name()
		r.checkOverwriteAllowed(overwrite, name)
		r.namedConstraints[name] = constraint
	}
}

func (r *constraintRegistry) registerNamed(overwrite bool, name string, constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.checkOverwriteAllowed(overwrite, name)
	r.namedConstraints[name] = constraint
}

func (r *constraintRegistry) registerManyNamed(overwrite bool, constraints map[string]Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for name, constraint := range constraints {
		r.checkOverwriteAllowed(overwrite, name)
		r.namedConstraints[name] = constraint
	}
}

func (r *constraintRegistry) get(name string) (Constraint, bool) {
	defer r.sync.Unlock()
	r.sync.Lock()
	c, ok := r.namedConstraints[name]
	return c, ok
}

// reset for testing
func (r *constraintRegistry) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.namedConstraints = defaultConstraints()
}

// has for testing
func (r *constraintRegistry) has(name string) bool {
	defer r.sync.Unlock()
	r.sync.Lock()
	_, ok := r.namedConstraints[name]
	return ok
}

// GetRegisteredConstraint returns a previously registered constraint
func GetRegisteredConstraint(name string) (Constraint, bool) {
	return constraintsRegistry.get(name)
}

// RegisterConstraint registers a Constraint for use by ValidatorFor
//
// For example:
//   RegisterConstraint(&MyConstraint{})
// will register the named constraint `MyConstraint` (using reflect to determine the name) which can then be used in a tag, e.g.
//   type MyRequest struct {
//      MyProperty string `json:"my_pty" v8n:"constraints:MyConstraint"`
//   }
// Use RegisterNamedConstraint to register a specific name (without using reflect to determine name
//
// Note: this function will panic if the constraint is already registered - use ReRegisterConstraint for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterConstraint(constraint Constraint) {
	constraintsRegistry.register(false, constraint)
}

// ReRegisterConstraint registers a Constraint for use by ValidatorFor
//
// For example:
//   ReRegisterConstraint(&MyConstraint{})
// will register the named constraint `MyConstraint` (using reflect to determine the name) which can then be used in a tag, e.g.
//   type MyRequest struct {
//      MyProperty string `json:"my_pty" v8n:"constraints:MyConstraint"`
//   }
// Use ReRegisterNamedConstraint to register a specific name (without using reflect to determine name
//
// If the constraint is already registered it is overwritten (this function will never panic)
func ReRegisterConstraint(constraint Constraint) {
	constraintsRegistry.register(true, constraint)
}

// RegisterNamedConstraint registers a Constraint for use by ValidatorFor with a specific name (or alias)
//
// For example:
//   RegisterNamedConstraint("myYes", &MyConstraint{SomeFlag: true})
//   RegisterNamedConstraint("myNo", &MyConstraint{SomeFlag: false})
// will register the two named constraints (with different settings) which can then be used in a tag, e.g.
//   type MyRequest struct {
//      MyProperty1 string `json:"my_pty_1" v8n:"constraints:myYes"`
//      MyProperty2 string `json:"my_pty_2" v8n:"constraints:myNo"`
//   }
//
// Note: this function will panic if the constraint is already registered - use ReRegisterNamedConstraint for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterNamedConstraint(name string, constraint Constraint) {
	constraintsRegistry.registerNamed(false, name, constraint)
}

// ReRegisterNamedConstraint registers a Constraint for use by ValidatorFor with a specific name (or alias)
//
// For example:
//   ReRegisterNamedConstraint("myYes", &MyConstraint{SomeFlag: true})
//   ReRegisterNamedConstraint("myNo", &MyConstraint{SomeFlag: false})
// will register the two named constraints (with different default settings) which can then be used in a tag, e.g.
//   type MyRequest struct {
//      MyProperty1 string `json:"my_pty_1" v8n:"constraints:myYes"`
//      MyProperty2 string `json:"my_pty_2" v8n:"constraints:myNo"`
//   }
//
// If the constraint is already registered it is overwritten (this function will never panic)
func ReRegisterNamedConstraint(name string, constraint Constraint) {
	constraintsRegistry.registerNamed(true, name, constraint)
}

// RegisterConstraints registers multiple constraints
//
// Note: this function will panic if the constraint is already registered - use ReRegisterConstraints for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterConstraints(constraints ...Constraint) {
	constraintsRegistry.registerMany(false, constraints...)
}

// ReRegisterConstraints registers multiple constraints
//
// If any of the constraints are already registered they are overwritten (this function will never panic)
func ReRegisterConstraints(constraints ...Constraint) {
	constraintsRegistry.registerMany(true, constraints...)
}

// RegisterNamedConstraints registers multiple named constraints
//
// Note: this function will panic if the constraint is already registered - use ReRegisterNamedConstraints for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterNamedConstraints(constraints map[string]Constraint) {
	constraintsRegistry.registerManyNamed(false, constraints)
}

// ReRegisterNamedConstraints registers multiple named constraints
//
// If any of the constraints are already registered they are overwritten (this function will never panic)
func ReRegisterNamedConstraints(constraints map[string]Constraint) {
	constraintsRegistry.registerManyNamed(true, constraints)
}

// ConstraintsRegistryReset is provided for test purposes - so the constraints registry can
// be cleared of all registered constraints (and reset to just the Valix common constraints)
func ConstraintsRegistryReset() {
	constraintsRegistry.reset()
}

// ConstraintsRegistryHas is provided for test purposes - so the constraints registry can
// be checked to see if a specific constraint name has been registered
func ConstraintsRegistryHas(name string) bool {
	return constraintsRegistry.has(name)
}

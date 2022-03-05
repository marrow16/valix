package valix

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

const (
	panicMsgConstraintExists = "Constraint \"%s\" already exists in registry"
)

type constraintsRegistry struct {
	namedConstraints map[string]Constraint
	sync             *sync.Mutex
}

var registry constraintsRegistry

func init() {
	registry = constraintsRegistry{
		namedConstraints: defaultConstraints(),
		sync:             &sync.Mutex{},
	}
}

func defaultConstraints() map[string]Constraint {
	return map[string]Constraint{
		"ArrayOf":                         &ArrayOf{},
		"DatetimeFuture":                  &DatetimeFuture{},
		"DatetimeFutureOrPresent":         &DatetimeFutureOrPresent{},
		"DatetimePast":                    &DatetimePast{},
		"DatetimePastOrPresent":           &DatetimePastOrPresent{},
		"Length":                          &Length{},
		"Maximum":                         &Maximum{},
		"Minimum":                         &Minimum{},
		"Negative":                        &Negative{},
		"NegativeOrZero":                  &NegativeOrZero{},
		"Positive":                        &Positive{},
		"PositiveOrZero":                  &PositiveOrZero{},
		"Range":                           &Range{},
		"StringCharacters":                &StringCharacters{},
		"StringNoControlCharacters":       &StringNoControlCharacters{},
		"StringNotBlank":                  &StringNotBlank{},
		"StringNotEmpty":                  &StringNotEmpty{},
		"StringLength":                    &StringLength{},
		"StringMaxLength":                 &StringMaxLength{},
		"StringMinLength":                 &StringMinLength{},
		"StringNormalizeUnicode":          &StringNormalizeUnicode{},
		"StringPattern":                   &StringPattern{},
		"StringTrim":                      &StringTrim{},
		"StringValidCardNumber":           &StringValidCardNumber{},
		"StringValidISODatetime":          &StringValidISODatetime{},
		"StringValidISODate":              &StringValidISODate{},
		"StringValidToken":                &StringValidToken{},
		"StringValidUnicodeNormalization": &StringValidUnicodeNormalization{},
		"StringValidUuid":                 &StringValidUuid{},
	}
}

func (r *constraintsRegistry) checkOverwriteAllowed(overwrite bool, name string) {
	if !overwrite {
		if _, ok := r.namedConstraints[name]; ok {
			panic(errors.New(fmt.Sprintf(panicMsgConstraintExists, name)))
		}
	}
}

func (r *constraintsRegistry) register(overwrite bool, constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	name := reflect.TypeOf(constraint).Elem().Name()
	r.checkOverwriteAllowed(overwrite, name)
	r.namedConstraints[name] = constraint
}

func (r *constraintsRegistry) registerMany(overwrite bool, constraints ...Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for _, constraint := range constraints {
		name := reflect.TypeOf(constraint).Elem().Name()
		r.checkOverwriteAllowed(overwrite, name)
		r.namedConstraints[name] = constraint
	}
}

func (r *constraintsRegistry) registerNamed(overwrite bool, name string, constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.checkOverwriteAllowed(overwrite, name)
	r.namedConstraints[name] = constraint
}

func (r *constraintsRegistry) registerManyNamed(overwrite bool, constraints map[string]Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for name, constraint := range constraints {
		r.checkOverwriteAllowed(overwrite, name)
		r.namedConstraints[name] = constraint
	}
}

func (r *constraintsRegistry) get(name string) (Constraint, bool) {
	defer r.sync.Unlock()
	r.sync.Lock()
	c, ok := r.namedConstraints[name]
	return c, ok
}

// reset for testing
func (r *constraintsRegistry) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.namedConstraints = defaultConstraints()
}

// has for testing
func (r *constraintsRegistry) has(name string) bool {
	defer r.sync.Unlock()
	r.sync.Lock()
	_, ok := r.namedConstraints[name]
	return ok
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
	registry.register(false, constraint)
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
	registry.register(true, constraint)
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
	registry.registerNamed(false, name, constraint)
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
	registry.registerNamed(true, name, constraint)
}

// RegisterConstraints registers multiple constraints
//
// Note: this function will panic if the constraint is already registered - use ReRegisterConstraints for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterConstraints(constraints ...Constraint) {
	registry.registerMany(false, constraints...)
}

// ReRegisterConstraints registers multiple constraints
//
// If any of the constraints are already registered they are overwritten (this function will never panic)
func ReRegisterConstraints(constraints ...Constraint) {
	registry.registerMany(true, constraints...)
}

// RegisterNamedConstraints registers multiple named constraints
//
// Note: this function will panic if the constraint is already registered - use ReRegisterNamedConstraints for
// non-panic behaviour where you don't mind the constraint registration being overwritten
func RegisterNamedConstraints(constraints map[string]Constraint) {
	registry.registerManyNamed(false, constraints)
}

// ReRegisterNamedConstraints registers multiple named constraints
//
// If any of the constraints are already registered they are overwritten (this function will never panic)
func ReRegisterNamedConstraints(constraints map[string]Constraint) {
	registry.registerManyNamed(true, constraints)
}

// ConstraintsRegistryReset is provided for test purposes - so the constraints registry can
// be cleared of all registered constraints (and reset to just the Valix common constraints)
func ConstraintsRegistryReset() {
	registry.reset()
}

// ConstraintsRegistryHas is provided for test purposes - so the constraints registry can
// be checked to see if a specific constraint name has been registered
func ConstraintsRegistryHas(name string) bool {
	return registry.has(name)
}

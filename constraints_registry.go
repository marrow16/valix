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

func (r *constraintsRegistry) register(constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	name := reflect.TypeOf(constraint).Elem().Name()
	if _, ok := r.namedConstraints[name]; ok {
		panic(errors.New(fmt.Sprintf(panicMsgConstraintExists, name)))
	}
	r.namedConstraints[name] = constraint
}

func (r *constraintsRegistry) registerMany(constraints ...Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for _, constraint := range constraints {
		name := reflect.TypeOf(constraint).Elem().Name()
		if _, ok := r.namedConstraints[name]; ok {
			panic(errors.New(fmt.Sprintf(panicMsgConstraintExists, name)))
		}
		r.namedConstraints[name] = constraint
	}
}

func (r *constraintsRegistry) registerNamed(name string, constraint Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	if _, ok := r.namedConstraints[name]; ok {
		panic(errors.New(fmt.Sprintf(panicMsgConstraintExists, name)))
	}
	r.namedConstraints[name] = constraint
}

func (r *constraintsRegistry) registerManyNamed(constraints map[string]Constraint) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for name, constraint := range constraints {
		if _, ok := r.namedConstraints[name]; ok {
			panic(errors.New(fmt.Sprintf(panicMsgConstraintExists, name)))
		}
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
func RegisterConstraint(constraint Constraint) {
	registry.register(constraint)
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
func RegisterNamedConstraint(name string, constraint Constraint) {
	registry.registerNamed(name, constraint)
}

// RegisterConstraints registers multiple constraints - see RegisterConstraint
func RegisterConstraints(constraints ...Constraint) {
	registry.registerMany(constraints...)
}

// RegisterNamedConstraints registers multiple named constraints - see RegisterNamedConstraint
func RegisterNamedConstraints(constraints map[string]Constraint) {
	registry.registerManyNamed(constraints)
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

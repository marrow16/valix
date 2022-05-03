package valix

import (
	"fmt"
	"sync"
)

type propertiesRepository struct {
	properties Properties
	panics     bool
	sync       *sync.Mutex
}

const (
	propertiesRepoPanicMsg = "property '%s' not found in properties repository"
)

var propertiesRepo = newPropertiesRepo()

func newPropertiesRepo() *propertiesRepository {
	return &propertiesRepository{
		properties: Properties{},
		panics:     true,
		sync:       &sync.Mutex{},
	}
}

// RegisterProperties registers common properties in the properties repository
//
// The properties repository is used by validators to lookup common properties when:
//
// * a property validator is nil, e.g.
//	v := valix.Validator{
//			Properties: valix.Properties{
//				"foo": nil,
//			},
//		}
//
// * a struct tag of `v8n-as` is used, e.g.
//	type MyStruct struct {
//		Foo string `json:"foo" v8n-as:""`
//	}
// or...
//	type MyStruct struct {
//		FooId string `json:"fooId" v8n-as:"id"`
//	}
func RegisterProperties(properties Properties) {
	propertiesRepo.register(properties)
}

// PropertiesRepoClear clears the properties repository
func PropertiesRepoClear() {
	propertiesRepo.clear()
}

// PropertiesRepoPanics sets whether the properties repository panics when asked for
// a common property that does not exist
func PropertiesRepoPanics(panics bool) {
	propertiesRepo.setPanics(panics)
}

// reset for testing
func (r *propertiesRepository) reset() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.panics = true
	r.properties = Properties{}
}

func (r *propertiesRepository) register(properties Properties) {
	defer r.sync.Unlock()
	r.sync.Lock()
	for k, v := range properties {
		if v == nil {
			r.properties[k] = &PropertyValidator{}
		} else {
			r.properties[k] = v
		}
	}
}

func (r *propertiesRepository) clear() {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.properties = Properties{}
}

func (r *propertiesRepository) setPanics(panics bool) {
	defer r.sync.Unlock()
	r.sync.Lock()
	r.panics = panics
}

func (r *propertiesRepository) fetch(properties Properties) Properties {
	defer r.sync.Unlock()
	r.sync.Lock()
	result := make(Properties, len(properties))
	for propertyName, pv := range properties {
		if pv == nil {
			if rpv, ok := r.properties[propertyName]; ok {
				result[propertyName] = rpv
			} else if r.panics {
				panic(fmt.Errorf(propertiesRepoPanicMsg, propertyName))
			} else {
				result[propertyName] = &PropertyValidator{}
			}
		} else {
			result[propertyName] = pv
		}
	}
	return result
}

// used by ValidatorFor / MustCompileValidatorFor
//
// NB. Does not panic, just returns nil if property not found in repo
func (r *propertiesRepository) getNamed(propertyName string) *PropertyValidator {
	defer r.sync.Unlock()
	r.sync.Lock()
	if rpv, ok := r.properties[propertyName]; ok {
		return rpv
	}
	return nil
}

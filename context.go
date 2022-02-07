package valix

import (
	"fmt"
)

type ValidatorContext struct {
	// ok is the final result of the validator
	ok bool
	// continueAll is whether entire validation is ok to continue
	continueAll bool
	// root is the original starting object (map or slice) for the validator
	root interface{}
	// violations is the collected violations
	violations []*Violation
	// pathStack is the path as each property or array index is checked
	pathStack []*pathStackItem
}

func newValidatorContext(root interface{}) *ValidatorContext {
	return &ValidatorContext{
		ok:          true,
		continueAll: true,
		root:        root,
		violations:  []*Violation{},
		pathStack: []*pathStackItem{
			{
				property: nil,
				path:     "",
				value:    root,
			},
		},
	}
}

// AddViolation adds a Violation to the validation context
//
// Note: Also causes the validator to fail (i.e. return false on CheckFunc)
func (vc *ValidatorContext) AddViolation(v *Violation) {
	vc.violations = append(vc.violations, v)
	vc.ok = false
}

// AddViolationForCurrent adds a Violation to the validation context for
// the current property and path
//
// Note: Also causes the validator to fail (i.e. return false on CheckFunc)
func (vc *ValidatorContext) AddViolationForCurrent(msg string) {
	curr := vc.pathStack[len(vc.pathStack)-1]
	vc.AddViolation(NewViolation(curr.propertyAsString(), curr.path, msg))
}

// Stop causes the entire validation to stop - i.e. not further constraints or
// property value validations are performed
//
// Note: This does not affect whether the validator succeeds or fails
func (vc *ValidatorContext) Stop() {
	vc.continueAll = false
}

// CeaseFurther causes further constraints and property validators on the current
// property to be ceased (i.e. not performed)
//
// Note: This does not affect whether the validator succeeds or fails
func (vc *ValidatorContext) CeaseFurther() {
	vc.pathStack[len(vc.pathStack)-1].stopped = true
}

// CurrentProperty returns the current property - which may be a string (for property name)
// or an int (for array index)
//
// Alternatively, to obtain the current property name use...
//    CurrentPropertyName
// or the current array index use...
//    CurrentArrayIndex
func (vc *ValidatorContext) CurrentProperty() interface{} {
	return vc.pathStack[len(vc.pathStack)-1].property
}

// CurrentPropertyName returns the current property name (or nil if current is an array index)
func (vc *ValidatorContext) CurrentPropertyName() *string {
	pty := vc.CurrentProperty()
	if s, ok := pty.(string); ok {
		return &s
	}
	return nil
}

// CurrentArrayIndex returns the current array index (or nil if current is a property name)
func (vc *ValidatorContext) CurrentArrayIndex() *int {
	pty := vc.CurrentProperty()
	if i, ok := pty.(int); ok {
		return &i
	}
	return nil
}

// CurrentPath returns the current property path
func (vc *ValidatorContext) CurrentPath() string {
	return vc.pathStack[len(vc.pathStack)-1].path
}

func (vc *ValidatorContext) CurrentValue() interface{} {
	return vc.pathStack[len(vc.pathStack)-1].value
	/*
		if pv, ok := vc.AncestorValue(0); ok {
			if apv, aok := pv.([]interface{}); aok {
				if idx := vc.CurrentArrayIndex(); idx != nil {
					return true, apv[*idx]
				}
			} else if opv, ook := pv.(map[string]interface{}); ook {
				if pty := vc.CurrentPropertyName(); pty != nil {
					return true, opv[*pty]
				}
			}
		}
		return false, nil
	*/
}

// CurrentDepth returns the current depth of the context - i.e. how many properties deep in the tree
func (vc *ValidatorContext) CurrentDepth() int {
	return len(vc.pathStack) - 1
}

// AncestorProperty returns an ancestor property - which may be a string (for property name)
// or an int (for array index)
//
// Alternatively, to obtain an ancestor property name use...
//    AncestorPropertyName
// or an ancestor array index use...
//    AncestorArrayIndex
func (vc *ValidatorContext) AncestorProperty(level uint) (interface{}, bool) {
	if itm, ok := vc.ancestorStackItem(level); ok {
		return itm.property, true
	}
	return nil, false
}

// AncestorPropertyName returns an ancestor property name (or nil if ancestor property is an array index)
func (vc *ValidatorContext) AncestorPropertyName(level uint) (*string, bool) {
	if itm, ok := vc.ancestorStackItem(level); ok {
		if s, oks := itm.property.(string); oks {
			return &s, true
		}
	}
	return nil, false
}

// AncestorArrayIndex returns an ancestor array index (or nil if ancestor is a property name)
func (vc *ValidatorContext) AncestorArrayIndex(level uint) (*int, bool) {
	if itm, ok := vc.ancestorStackItem(level); ok {
		if i, oki := itm.property.(int); oki {
			return &i, true
		}
	}
	return nil, false
}

// AncestorPath returns an ancestor property path
//
// The level determines how far up the ancestry - 0 is parent, 1 is grandparent, etc.
func (vc *ValidatorContext) AncestorPath(level uint) (*string, bool) {
	if itm, ok := vc.ancestorStackItem(level); ok {
		path := itm.path
		return &path, true
	}
	return nil, false
}

// AncestorValue returns an ancestor value
func (vc *ValidatorContext) AncestorValue(level uint) (interface{}, bool) {
	if itm, ok := vc.ancestorStackItem(level); ok {
		return itm.value, true
	}
	return nil, false
}

// SetCurrentValue set the value of the current property during validation
//
// Note: Use with extreme caution - altering values that are maps or slices (objects or arrays)
// during validation may cause severe problems where there are descendant validators on them
func (vc *ValidatorContext) SetCurrentValue(v interface{}) bool {
	if pv, ok := vc.AncestorValue(0); ok {
		currStackItem := vc.pathStack[len(vc.pathStack)-1]
		if apv, aok := pv.([]interface{}); aok {
			if idx := vc.CurrentArrayIndex(); idx != nil {
				apv[*idx] = v
				currStackItem.value = v
				return true
			}
		} else if opv, ook := pv.(map[string]interface{}); ook {
			if pty := vc.CurrentPropertyName(); pty != nil {
				opv[*pty] = v
				currStackItem.value = v
				return true
			}
		}
	}
	return false
}

func (vc *ValidatorContext) ancestorStackItem(level uint) (*pathStackItem, bool) {
	// note: -2 because... it's -1 for len and -1 for up one
	idx := len(vc.pathStack) - 2 - int(level)
	if idx < 0 {
		return nil, false
	}
	return vc.pathStack[idx], true
}

func (vc *ValidatorContext) pushPathProperty(property string, value interface{}) {
	curr := vc.pathStack[len(vc.pathStack)-1]
	vc.pathStack = append(vc.pathStack, &pathStackItem{
		property: property,
		path:     curr.asPath(),
		value:    value,
	})
}

func (vc *ValidatorContext) pushPathIndex(idx int, value interface{}) {
	curr := vc.pathStack[len(vc.pathStack)-1]
	vc.pathStack = append(vc.pathStack, &pathStackItem{
		property: idx,
		path:     curr.asPath(),
		value:    value,
	})
}

func (vc *ValidatorContext) popPath() {
	if len(vc.pathStack) > 1 {
		vc.pathStack = vc.pathStack[:len(vc.pathStack)-1]
	}
}

func (vc *ValidatorContext) continuePty() bool {
	return !vc.pathStack[len(vc.pathStack)-1].stopped
}

type pathStackItem struct {
	property interface{}
	path     string
	value    interface{}
	stopped  bool
}

func (p *pathStackItem) asPath() string {
	if p.property == nil {
		return ""
	}
	if pty, ok := p.property.(string); ok {
		if p.path != "" {
			return p.path + "." + pty
		}
		return pty
	}
	return p.path + p.propertyAsString()
}

func (p *pathStackItem) propertyAsString() string {
	if p.property == nil {
		return ""
	}
	if s, ok := p.property.(string); ok {
		return s
	}
	i, _ := p.property.(int)
	return fmt.Sprintf("[%d]", i)
}

package valix

import (
	"fmt"
)

type Context struct {
	ok           bool
	continueAll  bool
	top          interface{}
	topValidator *interface{}
	violations   []*Violation
	pathStack    []*pathStackItem
}

func newContext(top interface{}, topValidator interface{}) *Context {
	return &Context{
		ok:           true,
		continueAll:  true,
		top:          top,
		topValidator: &topValidator,
		violations:   []*Violation{},
		pathStack: []*pathStackItem{
			{
				property: nil,
				path:     "",
				value:    top,
			},
		},
	}
}

// AddViolation adds a Violation to the validation context
//
// Note: Also causes the validator to fail (i.e. return false on validate)
func (c *Context) AddViolation(v *Violation) {
	c.violations = append(c.violations, v)
	c.ok = false
}

// AddViolationForCurrent adds a Violation to the validation context for
// the current property and path
//
// Note: Also causes the validator to fail (i.e. return false on validate)
func (c *Context) AddViolationForCurrent(msg string) {
	curr := c.pathStack[len(c.pathStack)-1]
	c.AddViolation(NewViolation(curr.propertyAsString(), curr.path, msg))
}

// Stop causes the entire validation to stop - i.e. not further constraints or
// property value validations are performed
//
// Note: This does not affect whether the validator succeeds or fails
func (c *Context) Stop() {
	c.continueAll = false
}

// CeaseFurther causes further constraints and property validators on the current
// property to be ceased (i.e. not performed)
//
// Note: This does not affect whether the validator succeeds or fails
func (c *Context) CeaseFurther() {
	c.pathStack[len(c.pathStack)-1].stopped = true
}

// CurrentProperty returns the current property - which may be a string (for property name)
// or an int (for array index)
//
// Alternatively, to obtain the current property name use...
//    CurrentPropertyName
// or the current array index use...
//    CurrentArrayIndex
func (c *Context) CurrentProperty() interface{} {
	return c.pathStack[len(c.pathStack)-1].property
}

// CurrentPropertyName returns the current property name (or nil if current is an array index)
func (c *Context) CurrentPropertyName() *string {
	pty := c.CurrentProperty()
	if s, ok := pty.(string); ok {
		return &s
	}
	return nil
}

// CurrentArrayIndex returns the current array index (or nil if current is a property name)
func (c *Context) CurrentArrayIndex() *int {
	pty := c.CurrentProperty()
	if i, ok := pty.(int); ok {
		return &i
	}
	return nil
}

// CurrentPath returns the current property path
func (c *Context) CurrentPath() string {
	return c.pathStack[len(c.pathStack)-1].path
}

// CurrentDepth returns the current depth of the context - i.e. how many properties deep in the tree
func (c *Context) CurrentDepth() int {
	return len(c.pathStack) - 1
}

// AncestorProperty returns an ancestor property - which may be a string (for property name)
// or an int (for array index)
//
// Alternatively, to obtain an ancestor property name use...
//    AncestorPropertyName
// or an ancestor array index use...
//    AncestorArrayIndex
func (c *Context) AncestorProperty(level uint) (interface{}, bool) {
	if itm, ok := c.ancestorStackItem(level); ok {
		return itm.property, true
	}
	return nil, false
}

// AncestorPropertyName returns an ancestor property name (or nil if ancestor property is an array index)
func (c *Context) AncestorPropertyName(level uint) (*string, bool) {
	if itm, ok := c.ancestorStackItem(level); ok {
		if s, oks := itm.property.(string); oks {
			return &s, true
		}
	}
	return nil, false
}

// AncestorArrayIndex returns an ancestor array index (or nil if ancestor is a property name)
func (c *Context) AncestorArrayIndex(level uint) (*int, bool) {
	if itm, ok := c.ancestorStackItem(level); ok {
		if i, oki := itm.property.(int); oki {
			return &i, true
		}
	}
	return nil, false
}

// AncestorPath returns an ancestor property path
//
// The level determines how far up the ancestry - 0 is parent, 1 is grandparent, etc.
func (c *Context) AncestorPath(level uint) (*string, bool) {
	if itm, ok := c.ancestorStackItem(level); ok {
		path := itm.path
		return &path, true
	}
	return nil, false
}

// AncestorValue returns an ancestor value
func (c *Context) AncestorValue(level uint) (interface{}, bool) {
	if itm, ok := c.ancestorStackItem(level); ok {
		return itm.value, true
	}
	return nil, false
}

// SetCurrentValue set the value of the current property during validation
//
// Note: Use with extreme caution - altering values during validation may cause problems
// for other constraints or validators in the chain
func (c *Context) SetCurrentValue(v interface{}) bool {
	if pv, ok := c.AncestorValue(0); ok {
		if apv, aok := pv.([]interface{}); aok {
			if idx := c.CurrentArrayIndex(); idx != nil {
				apv[*idx] = v
				return true
			}
		} else if opv, ook := pv.(map[string]interface{}); ook {
			if pty := c.CurrentPropertyName(); pty != nil {
				opv[*pty] = v
				return true
			}
		}
	}
	return false
}

func (c *Context) ancestorStackItem(level uint) (*pathStackItem, bool) {
	// note: -2 because... it's -1 for len and -1 for up one
	idx := len(c.pathStack) - 2 - int(level)
	if idx < 0 {
		return nil, false
	}
	return c.pathStack[idx], true
}

func (c *Context) pushPathProperty(property string, value interface{}) {
	curr := c.pathStack[len(c.pathStack)-1]
	c.pathStack = append(c.pathStack, &pathStackItem{
		property: property,
		path:     curr.asPath(),
		value:    value,
	})
}

func (c *Context) pushPathIndex(idx int, value interface{}) {
	curr := c.pathStack[len(c.pathStack)-1]
	c.pathStack = append(c.pathStack, &pathStackItem{
		property: idx,
		path:     curr.asPath(),
		value:    value,
	})
}

func (c *Context) popPath() {
	if len(c.pathStack) > 1 {
		c.pathStack = c.pathStack[:len(c.pathStack)-1]
	}
}

func (c *Context) continuePty() bool {
	return !c.pathStack[len(c.pathStack)-1].stopped
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

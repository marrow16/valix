package valix

import "fmt"

type Context struct {
	ok           bool
	continues    bool
	top          interface{}
	topValidator *interface{}
	violations   []*Violation
	pathStack    []*pathStackItem
}

type pathStackItem struct {
	propertyName string
	path         string
	value        interface{}
}

func newContext(top interface{}, topValidator interface{}) *Context {
	return &Context{
		ok:           true,
		continues:    true,
		top:          top,
		topValidator: &topValidator,
		violations:   []*Violation{},
		pathStack: []*pathStackItem{
			{
				propertyName: "",
				path:         "",
				value:        top,
			},
		},
	}
}

func (c *Context) AddViolation(v *Violation) {
	c.violations = append(c.violations, v)
	c.ok = false
}

func (c *Context) AddViolationForCurrent(msg string) {
	curr := c.pathStack[len(c.pathStack)-1]
	c.AddViolation(NewViolation(curr.propertyName, curr.path, msg))
}

func (c *Context) Stop() {
	c.continues = false
}

func (c *Context) StopWith(v *Violation) {
	c.AddViolation(v)
	c.Stop()
}

func (c *Context) StopWithCurrent(msg string) {
	c.AddViolationForCurrent(msg)
	c.Stop()
}

func (c *Context) PropertyName() string {
	return c.pathStack[len(c.pathStack)-1].propertyName
}

func (c *Context) Path() string {
	return c.pathStack[len(c.pathStack)-1].path
}

func (c *Context) pushPathProperty(propertyName string, value interface{}) {
	curr := c.pathStack[len(c.pathStack)-1]
	newPath := curr.path
	if newPath != "" && curr.propertyName != "" {
		newPath = newPath + "." + curr.propertyName
	} else if curr.propertyName != "" {
		newPath = curr.propertyName
	}
	c.pathStack = append(c.pathStack, &pathStackItem{
		path:         newPath,
		propertyName: propertyName,
		value:        value,
	})
}

func (c *Context) pushPathIndex(idx int, value interface{}) {
	curr := c.pathStack[len(c.pathStack)-1]
	newPath := curr.path
	if newPath != "" && curr.propertyName != "" {
		newPath = newPath + "." + curr.propertyName
	} else if curr.propertyName != "" {
		newPath = curr.propertyName
	}
	c.pathStack = append(c.pathStack, &pathStackItem{
		path:         newPath + fmt.Sprintf("[%d]", idx),
		propertyName: "",
		value:        value,
	})
}

func (c *Context) popPath() {
	if len(c.pathStack) > 1 {
		c.pathStack = c.pathStack[:len(c.pathStack)-1]
	}
}

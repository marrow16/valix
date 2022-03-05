package valix

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestConstraintsRegistryInitialized(t *testing.T) {
	registry.reset()
	require.Equal(t, 29, len(registry.namedConstraints))
}

func TestRegisterConstraint(t *testing.T) {
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))
}

func TestRegisterConstraintPanicsAddingDuplicate(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, myConstraintName), r.(error).Error())
		}
	}()

	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
}

func TestReRegisterConstraintNotPanicsAddingDuplicate(t *testing.T) {
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	ReRegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))

	ReRegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))
}

func TestRegisterNamedConstraint(t *testing.T) {
	const testName = "TestName"
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, registry.has(testName))
}

func TestRegisterNamedConstraintPanicsAddingDuplicate(t *testing.T) {
	const testName = "TestName"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, testName), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, registry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint2{})
}

func TestReRegisterNamedConstraintNotPanicsAddingDuplicate(t *testing.T) {
	const testName = "TestName"
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName))

	ReRegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, registry.has(testName))

	ReRegisterNamedConstraint(testName, &myConstraint2{})
	require.True(t, registry.has(testName))
}

func TestRegisterConstraints(t *testing.T) {
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))
	require.False(t, registry.has(myConstraint2Name))

	RegisterConstraints(&myConstraint{}, &myConstraint2{})
	require.True(t, registry.has(myConstraintName))
	require.True(t, registry.has(myConstraint2Name))
}

func TestRegisterConstraintsPanicsAddingDuplicates(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, myConstraintName), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	RegisterConstraints(&myConstraint{}, &myConstraint{})
}

func TestReRegisterConstraintsNotPanicsAddingDuplicates(t *testing.T) {
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	// register two of the same...
	ReRegisterConstraints(&myConstraint{}, &myConstraint{})
	require.True(t, registry.has(myConstraintName))
}

func TestRegisterNamedConstraints(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName1))
	require.False(t, registry.has(testName2))

	RegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, registry.has(testName1))
	require.True(t, registry.has(testName2))
}

func TestRegisterNamedConstraintsPanicsWithDuplicate(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, testName1), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName1))
	require.False(t, registry.has(testName2))

	// register one first...
	RegisterNamedConstraint(testName1, &myConstraint{})

	RegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, registry.has(testName1))
	require.True(t, registry.has(testName2))
}

func TestReRegisterNamedConstraintsNotPanicsWithDuplicate(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(testName1))
	require.False(t, registry.has(testName2))

	// register them first...
	RegisterNamedConstraint(testName1, &myConstraint{})
	RegisterNamedConstraint(testName2, &myConstraint{})
	c, ok := registry.get(testName1)
	require.True(t, ok)
	require.False(t, c.(*myConstraint).SomeFlag)
	c, ok = registry.get(testName2)
	require.True(t, ok)
	require.False(t, c.(*myConstraint).SomeFlag)

	ReRegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, registry.has(testName1))
	require.True(t, registry.has(testName2))

	c, ok = registry.get(testName1)
	require.True(t, ok)
	require.True(t, c.(*myConstraint).SomeFlag) // default flag should be set now
	c, ok = registry.get(testName2)
	require.True(t, ok)
	require.False(t, c.(*myConstraint2).SomeFlag)
}

func TestConstraintsRegistryReset(t *testing.T) {
	defer registry.reset()
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))

	ConstraintsRegistryReset()
	require.False(t, registry.has(myConstraintName))
}

func TestConstraintsRegistryHas(t *testing.T) {
	defer registry.reset()
	// make sure it isn't there first...
	registry.reset()
	require.False(t, registry.has(myConstraintName))
	require.False(t, ConstraintsRegistryHas(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, registry.has(myConstraintName))
	require.True(t, ConstraintsRegistryHas(myConstraintName))
}

const myConstraintName = "myConstraint"
const myConstraint2Name = "myConstraint2"

type myConstraint struct {
	SomeFlag bool
}

func (my *myConstraint) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (my *myConstraint) GetMessage() string {
	return "My test message"
}

type myConstraint2 struct {
	SomeFlag bool
}

func (my *myConstraint2) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (my *myConstraint2) GetMessage() string {
	return "My test message"
}

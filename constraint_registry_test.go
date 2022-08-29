package valix

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const commonConstraintsCount = 89 // excludes abbreviations (every constraint has an abbreviation)
const commonSpecialAbbrsCount = 8 // special abbreviations

func TestConstraintsRegistryInitialized(t *testing.T) {
	constraintsRegistry.reset()
	require.Equal(t, (commonConstraintsCount*2)+commonSpecialAbbrsCount+len(getBuiltInPresets()), len(constraintsRegistry.namedConstraints))
}

func TestRegisterConstraint(t *testing.T) {
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))
}

func TestRegisterConstraintPanicsAddingDuplicate(t *testing.T) {
	defer func() {
		constraintsRegistry.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, myConstraintName), r.(error).Error())
		}
	}()

	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
}

func TestReRegisterConstraintNotPanicsAddingDuplicate(t *testing.T) {
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(myConstraintName))

	ReRegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))

	ReRegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))
}

func TestRegisterNamedConstraint(t *testing.T) {
	const testName = "TestName"
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, constraintsRegistry.has(testName))
}

func TestRegisterNamedConstraintPanicsAddingDuplicate(t *testing.T) {
	const testName = "TestName"
	defer func() {
		constraintsRegistry.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, testName), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, constraintsRegistry.has(testName))

	RegisterNamedConstraint(testName, &myConstraint2{})
}

func TestReRegisterNamedConstraintNotPanicsAddingDuplicate(t *testing.T) {
	const testName = "TestName"
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(testName))

	ReRegisterNamedConstraint(testName, &myConstraint{})
	require.True(t, constraintsRegistry.has(testName))

	ReRegisterNamedConstraint(testName, &myConstraint2{})
	require.True(t, constraintsRegistry.has(testName))
}

func TestRegisterConstraints(t *testing.T) {
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(myConstraintName))
	require.False(t, constraintsRegistry.has(myConstraint2Name))

	RegisterConstraints(&myConstraint{}, &myConstraint2{})
	require.True(t, constraintsRegistry.has(myConstraintName))
	require.True(t, constraintsRegistry.has(myConstraint2Name))
}

func TestRegisterConstraintsPanicsAddingDuplicates(t *testing.T) {
	defer func() {
		constraintsRegistry.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, myConstraintName), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(myConstraintName))

	RegisterConstraints(&myConstraint{}, &myConstraint{})
}

func TestReRegisterConstraintsNotPanicsAddingDuplicates(t *testing.T) {
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(myConstraintName))

	// register two of the same...
	ReRegisterConstraints(&myConstraint{}, &myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))
}

func TestRegisterNamedConstraints(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(testName1))
	require.False(t, constraintsRegistry.has(testName2))

	RegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, constraintsRegistry.has(testName1))
	require.True(t, constraintsRegistry.has(testName2))
}

func TestRegisterNamedConstraintsPanicsWithDuplicate(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	defer func() {
		constraintsRegistry.reset()
		if r := recover(); r == nil {
			t.Errorf("Expected panic")
		} else {
			require.Equal(t, fmt.Sprintf(panicMsgConstraintExists, testName1), r.(error).Error())
		}
	}()
	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(testName1))
	require.False(t, constraintsRegistry.has(testName2))

	// register one first...
	RegisterNamedConstraint(testName1, &myConstraint{})

	RegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, constraintsRegistry.has(testName1))
	require.True(t, constraintsRegistry.has(testName2))
}

func TestReRegisterNamedConstraintsNotPanicsWithDuplicate(t *testing.T) {
	const testName1 = "TestName1"
	const testName2 = "TestName2"
	// make sure it isn't there first...
	constraintsRegistry.reset()
	defer func() {
		constraintsRegistry.reset()
	}()
	require.False(t, constraintsRegistry.has(testName1))
	require.False(t, constraintsRegistry.has(testName2))

	// register them first...
	RegisterNamedConstraint(testName1, &myConstraint{})
	RegisterNamedConstraint(testName2, &myConstraint{})
	c, ok := constraintsRegistry.get(testName1)
	require.True(t, ok)
	require.False(t, c.(*myConstraint).SomeFlag)
	c, ok = constraintsRegistry.get(testName2)
	require.True(t, ok)
	require.False(t, c.(*myConstraint).SomeFlag)

	ReRegisterNamedConstraints(map[string]Constraint{
		testName1: &myConstraint{SomeFlag: true},
		testName2: &myConstraint2{SomeFlag: false},
	})
	require.True(t, constraintsRegistry.has(testName1))
	require.True(t, constraintsRegistry.has(testName2))

	c, ok = constraintsRegistry.get(testName1)
	require.True(t, ok)
	require.True(t, c.(*myConstraint).SomeFlag) // default flag should be set now
	c, ok = constraintsRegistry.get(testName2)
	require.True(t, ok)
	require.False(t, c.(*myConstraint2).SomeFlag)
}

func TestConstraintsRegistryReset(t *testing.T) {
	defer constraintsRegistry.reset()
	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))

	ConstraintsRegistryReset()
	require.False(t, constraintsRegistry.has(myConstraintName))
}

func TestConstraintsRegistryHas(t *testing.T) {
	defer constraintsRegistry.reset()
	// make sure it isn't there first...
	constraintsRegistry.reset()
	require.False(t, constraintsRegistry.has(myConstraintName))
	require.False(t, ConstraintsRegistryHas(myConstraintName))

	RegisterConstraint(&myConstraint{})
	require.True(t, constraintsRegistry.has(myConstraintName))
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
func (my *myConstraint) GetMessage(tcx I18nContext) string {
	return obtainI18nContext(tcx).TranslateMessage("My test message")
}

type myConstraint2 struct {
	SomeFlag bool
}

func (my *myConstraint2) Check(value interface{}, vcx *ValidatorContext) (bool, string) {
	return true, ""
}
func (my *myConstraint2) GetMessage(tcx I18nContext) string {
	return obtainI18nContext(tcx).TranslateMessage("My test message")
}

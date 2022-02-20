package examples

import (
	"github.com/marrow16/valix"
	"github.com/stretchr/testify/require"
	"testing"
)

func init() {
	valix.RegisterConstraint(&NoFoo{})
}

type myConstraint1 struct {
	Message string
}

func (m *myConstraint1) Check(value interface{}, vcx *valix.ValidatorContext) (bool, string) {
	return true, ""
}

func (m *myConstraint1) GetMessage() string {
	return "This is constraint 2 message"
}

type myConstraint2 struct {
	Message string
}

func (m *myConstraint2) Check(value interface{}, vcx *valix.ValidatorContext) (bool, string) {
	return true, ""
}

func (m *myConstraint2) GetMessage() string {
	return "This is constraint 1 message"
}

var myConstraintSet = &valix.ConstraintSet{
	Constraints: valix.Constraints{
		&myConstraint1{},
		&myConstraint2{},
	},
	Message: "This is constraint set message",
}

func TestCanRegisterConstraints(t *testing.T) {
	defer valix.ConstraintsRegistryReset()

	valix.RegisterConstraint(&myConstraint1{})
	require.True(t, valix.ConstraintsRegistryHas("myConstraint1"))

	valix.RegisterConstraints(&myConstraint2{}, myConstraintSet)
	require.True(t, valix.ConstraintsRegistryHas("myConstraint2"))
	require.True(t, valix.ConstraintsRegistryHas("ConstraintSet")) // uses the name of struct rather than name of variable!

	valix.ConstraintsRegistryReset()
	require.False(t, valix.ConstraintsRegistryHas("myConstraint1"))
	require.False(t, valix.ConstraintsRegistryHas("myConstraint2"))
	require.False(t, valix.ConstraintsRegistryHas("ConstraintSet"))

	valix.RegisterConstraints(&myConstraint1{}, &myConstraint2{})
	valix.RegisterNamedConstraint("AliasedConstraintSetName", myConstraintSet)
	require.True(t, valix.ConstraintsRegistryHas("myConstraint1"))
	require.True(t, valix.ConstraintsRegistryHas("myConstraint2"))
	require.False(t, valix.ConstraintsRegistryHas("ConstraintSet"))
	require.True(t, valix.ConstraintsRegistryHas("AliasedConstraintSetName"))
}

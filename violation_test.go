package valix

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateNewViolationStruct(t *testing.T) {
	v := &Violation{}

	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.Equal(t, "", v.Message)
	require.False(t, v.BadRequest)
}

func TestNewEmptyViolation(t *testing.T) {
	const msg = "MESSAGE"
	v := NewEmptyViolation(msg)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.False(t, v.BadRequest)
}

func TestNewViolation(t *testing.T) {
	const pty = "PTY"
	const path = "PATH"
	const msg = "MESSAGE"
	v := NewViolation(pty, path, msg)

	require.Equal(t, pty, v.Property)
	require.Equal(t, path, v.Path)
	require.Equal(t, msg, v.Message)
	require.False(t, v.BadRequest)
}

func TestNewBadRequestViolation(t *testing.T) {
	const msg = "MESSAGE"
	v := NewBadRequestViolation(msg)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.True(t, v.BadRequest)
}

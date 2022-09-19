package valix

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateNewViolationStruct(t *testing.T) {
	v := &Violation{}

	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.Equal(t, "", v.Message)
	require.False(t, v.BadRequest)
	require.Equal(t, 0, len(v.Codes))
}

func TestNewEmptyViolation(t *testing.T) {
	const msg = "MESSAGE"
	v := NewEmptyViolation(msg)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.False(t, v.BadRequest)
	require.Equal(t, 0, len(v.Codes))
}

func TestNewEmptyViolationWithCodes(t *testing.T) {
	const msg = "MESSAGE"
	v := NewEmptyViolation(msg, "123", 345, true)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.False(t, v.BadRequest)
	require.Equal(t, 3, len(v.Codes))
	require.Equal(t, "123", v.Codes[0])
	require.Equal(t, 345, v.Codes[1])
	require.True(t, v.Codes[2].(bool))
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
	require.Equal(t, 0, len(v.Codes))
}

func TestNewViolationWithCodes(t *testing.T) {
	const pty = "PTY"
	const path = "PATH"
	const msg = "MESSAGE"
	v := NewViolation(pty, path, msg, "123", 345)

	require.Equal(t, pty, v.Property)
	require.Equal(t, path, v.Path)
	require.Equal(t, msg, v.Message)
	require.False(t, v.BadRequest)
	require.Equal(t, 2, len(v.Codes))
	require.Equal(t, "123", v.Codes[0])
	require.Equal(t, 345, v.Codes[1])
}

func TestNewBadRequestViolation(t *testing.T) {
	const msg = "MESSAGE"
	v := NewBadRequestViolation(msg)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.True(t, v.BadRequest)
	require.Equal(t, 0, len(v.Codes))
}

func TestNewBadRequestViolationWithCodes(t *testing.T) {
	const msg = "MESSAGE"
	v := NewBadRequestViolation(msg, "123", 345)

	require.Equal(t, msg, v.Message)
	require.Equal(t, "", v.Property)
	require.Equal(t, "", v.Path)
	require.True(t, v.BadRequest)
	require.Equal(t, 2, len(v.Codes))
	require.Equal(t, "123", v.Codes[0])
	require.Equal(t, 345, v.Codes[1])
}

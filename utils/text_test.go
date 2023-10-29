package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReverseString(t *testing.T) {
	refreshSecretKey := "1234567890123456"
	refreshSecretKey = "rftk" + refreshSecretKey[3:]
	require.Equal(t, "rftk567890123456", refreshSecretKey)


	text := "Hello"
	reversed := ReverseString(text)
	require.NotEmpty(t, reversed)
	require.Equal(t, "olleH", reversed)

	text = "12345"
	reversed = ReverseString(text)
	require.NotEmpty(t, reversed)
	require.Equal(t, "54321", reversed)

	text = "a1s2d3"
	reversed = ReverseString(text)
	require.NotEmpty(t, reversed)
	require.Equal(t, "3d2s1a", reversed)
}

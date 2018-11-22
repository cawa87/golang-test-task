package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestReadLines(t *testing.T) {
	lines, err := readLines("testdata/lines")
	require.NoError(t, err)
	require.Equal(t, []string{
		"one",
		"two",
		"three",
	}, lines)
}

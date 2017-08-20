package main

import "github.com/stretchr/testify/assert"
import "runtime"
import "testing"

func TestShouldLoadConfig(t *testing.T) {
	var config Config;
	assert.NoError(t, config.Load())
}

func TestConfigDefauls(t *testing.T) {
	var config Config;
	assert.NoError(t, config.Load())
	assert.Equal(t, runtime.NumCPU(), config.scanBodyConcurrency)
}

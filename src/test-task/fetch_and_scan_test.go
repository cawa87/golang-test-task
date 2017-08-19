package main

import "github.com/stretchr/testify/assert"
import "testing"

func TestFetchAndScanEmptyUrls(t *testing.T) {
	data, e := fetchAndScan([]string{}, nil)
	assert.NoError(t, e)
	assert.Equal(t, []*FetchAndScanData{}, data)
}

package main

import "github.com/stretchr/testify/assert"
import "net/http"
import "testing"

func TestHeaderGetMediaType(t *testing.T) {
	h := http.Header{}
	h.Set("one",   "text/html; charset=windows-1251")
	h.Set("two",   "text/html")
	h.Set("three", "")
	assert.Equal(t, "text/html", headerGetMediaType(h, "one"))
	assert.Equal(t, "text/html", headerGetMediaType(h, "two"))
	assert.Equal(t, "",          headerGetMediaType(h, "three"))
	assert.Equal(t, "",          headerGetMediaType(h, "four"))
}

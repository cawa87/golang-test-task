package http_plus

import "github.com/stretchr/testify/assert"
import "net/http"
import "testing"

func TestHeaderGetMediaType(t *testing.T) {
	h := http.Header{}
	h.Set("one",   "text/html; charset=windows-1251")
	h.Set("two",   "text/html")
	h.Set("three", "")
	assert.Equal(t, "text/html", HeaderGetMediaType(h, "one"))
	assert.Equal(t, "text/html", HeaderGetMediaType(h, "two"))
	assert.Equal(t, "",          HeaderGetMediaType(h, "three"))
	assert.Equal(t, "",          HeaderGetMediaType(h, "four"))
}

func TestHeaderGetContentLength(t *testing.T) {
	h := http.Header{}
	assert.Equal(t, int64(-1), HeaderGetContentLength(h, -1))
	assert.Equal(t, int64(-1), HeaderGetContentLength(h,  0))
	assert.Equal(t, int64( 1), HeaderGetContentLength(h,  1))
	h.Set("Content-Length", "")
	assert.Equal(t, int64(-1), HeaderGetContentLength(h, -1))
	assert.Equal(t, int64(-1), HeaderGetContentLength(h,  0))
	assert.Equal(t, int64( 1), HeaderGetContentLength(h,  1))
	h.Set("Content-Length", "2")
	assert.Equal(t, int64( 2), HeaderGetContentLength(h, -1))
	assert.Equal(t, int64( 2), HeaderGetContentLength(h,  0))
	assert.Equal(t, int64( 2), HeaderGetContentLength(h,  1))
}

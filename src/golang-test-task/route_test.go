package main

import "bytes"
import "github.com/stretchr/testify/assert"
import "net/http/httptest"
import "testing"

func TestExtractUrls(t *testing.T) {
	body := bytes.NewBufferString(`["one", "two", "three"]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("content-type", "application/json")
	urls, _ := extractUrls(req)
	assert.Equal(t, []string{"one", "two", "three"}, urls)
}

func TestDefaultRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	(&Route{}).ServeHTTP(res, req)
	assert.Equal(t, 404, res.Code)
}

func TestResponseShouldhaveProperContentType(t *testing.T) {
	body := bytes.NewBufferString(`[]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	(&Route{}).ServeHTTP(res, req)
	assert.Equal(t, "application/json", res.Header().Get("content-type"))
}

func TestSendEmptyContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	res := httptest.NewRecorder()
	(&Route{}).ServeHTTP(res, req)
	assert.Equal(t, 400, res.Code)
}

func TestSendEmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("content-type", "application/json")
	_, e := extractUrls(req)
	assert.Error(t, e)
}

func TestSendInvalidJSONStruct(t *testing.T) {
	body := bytes.NewBufferString(`[{"one":1}, {"two":2}]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("content-type", "application/json")
	_, e := extractUrls(req)
	assert.Error(t, e)
}

func TestSendEmptyArray(t *testing.T) {
	body := bytes.NewBufferString(`[]`)
	req := httptest.NewRequest("POST", "/", body)
	req.Header.Set("content-type", "application/json")
	res := httptest.NewRecorder()
	(&Route{}).ServeHTTP(res, req)
	assert.Equal(t, "[]", res.Body.String())
}


package main

import "net/http/httptest"
import "testing"

func TestDefaultRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	(&Route{}).ServeHTTP(res, req)
	if res.Code != 404 { t.Fatal("404 expected") }
}

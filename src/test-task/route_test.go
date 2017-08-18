package main

import "net/http"
import "net/http/httptest"
import "testing"

func TestDefaultRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()

	defer func() {
		if recover() != http.ErrAbortHandler {
			t.Fatal("ErrAbortHandler expected")
	}}()

	(&Route{}).ServeHTTP(res, req)
}

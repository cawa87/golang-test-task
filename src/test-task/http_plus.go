package main

import "mime"
import "net/http"

func headerGetMediaType(header http.Header, key string) string {
	if a := header.Get(key); a != "" {
		t, _, e := mime.ParseMediaType(a)
		if e == nil { return t }}
	return ""
}

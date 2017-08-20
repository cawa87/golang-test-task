package http_plus

import "mime"
import "net/http"
import "strconv"

func HeaderGetMediaType(header http.Header, key string) string {
	if a := header.Get(key); a != "" {
		t, _, e := mime.ParseMediaType(a)
		if e == nil { return t }}
	return ""
}

func HeaderGetContentLength(header http.Header, hint int64) int64 {
	if cl := header.Get("Content-Length"); cl != "" {
		if l, e := strconv.ParseInt(cl, 10, 64); e == nil && l >=0 {
			return l }}
	if hint > 0 { return hint }
	return -1
}

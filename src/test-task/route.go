package main

import "encoding/json"
import "er"
import "net/http"

type Route struct {}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusNotFound); return }

	_, e := extractUrls(req)
	if e != nil { res.WriteHeader(http.StatusBadRequest); return }
}

func extractUrls(req *http.Request) ([]string, error) {
	if req.Header.Get("content-type") != "application/json" {
		return nil, er.Er(nil, "Expected Content-Type: application/json") }

	d := json.NewDecoder(req.Body)
	var urls []string
	if e := d.Decode(&urls); e != nil {
		return nil, er.Er(e, "Invalid JSON struct or parsing error") }

	return urls, nil
}

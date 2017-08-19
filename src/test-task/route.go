package main

import "encoding/json"
import "er"
import "net/http"
import "printlog"

type Route struct {
	fetchPipe chan<- FetchAndScanTask
}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		res.WriteHeader(http.StatusNotFound); return }

	urls, e := extractUrls(req)
	if e != nil {
		res.WriteHeader(http.StatusBadRequest); return }

	data, e := fetchAndScan(urls, r.fetchPipe)
	if e != nil {
		printlog.Error(er.Er(e, "fetchAndScan"))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	if e = writeJson(res, &data); e != nil {
		printlog.Error(e)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func extractUrls(req *http.Request) ([]string, error) {
	if req.Header.Get("content-type") != "application/json" {
		return nil, er.Er(nil, "Expected Content-Type: application/json") }

	jd := json.NewDecoder(req.Body)
	var urls []string
	if e := jd.Decode(&urls); e != nil {
		return nil, er.Er(e, "Invalid JSON struct or parsing error") }

	return urls, nil
}

func writeJson(res http.ResponseWriter, data interface{}) error {
	res.Header().Set("content-type", "application/json")

	bytes, e := json.Marshal(data)
	if e != nil { return er.Er(e, "json.Marhsal") }

	_, e = res.Write(bytes)
	if e != nil { return er.Er(e, "http.ResponseWriter.Write") }

	return nil
}

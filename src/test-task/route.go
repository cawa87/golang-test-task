package main

import "net/http"

type Route struct {}

func (r *Route) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	panic(http.ErrAbortHandler)
}

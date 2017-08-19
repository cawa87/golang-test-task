package main

import "er"
import "net/http"
import "os"
import "printlog"

func main() {
	printlog.Verbosity = printlog.DEBUG

	var config Config
	if e := config.Load(); e != nil {
		printlog.Error(er.Er(e, "Fail to load config"))
		os.Exit(1)
	}

	fetchPipe := make(chan FetchAndScanTask)
	startFetchWorkers(config.fetchUrlsConcurrency, fetchPipe)

	route := &Route{fetchPipe}
	if e := http.ListenAndServe(config.listen, route); e != nil {
		printlog.Error(er.Er(e, "Fail to start server"))
		os.Exit(2)
	}
}

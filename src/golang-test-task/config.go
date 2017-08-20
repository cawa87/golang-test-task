package main

import "er"
import "github.com/docopt/docopt-go"
import "runtime"
import "strconv"

const doc = `
Usage:
	test-task [options]

Options:
	-h, --help  This help
	-l ADDR     Listen on [host]:port [default: :80]
	-f COUNT    Fetch max COUNT urls in parallel [default: 10]
	-s COUNT    Scan max COUNT bodies in parallel [default: 0]
`

type Config struct {
	listen              string
	fetchUrlConcurrency int
	scanBodyConcurrency int
}

func (c* Config) Load() error {
	opt, e := docopt.Parse(doc, nil, true, "", true, true)
	if e != nil { return er.Er(e, "docopt.Parse") }

	c.listen                 = opt["-l"].(string)
	c.fetchUrlConcurrency, _ = strconv.Atoi(opt["-f"].(string))
	c.scanBodyConcurrency, _ = strconv.Atoi(opt["-s"].(string))
	if c.scanBodyConcurrency <= 0 {
		c.scanBodyConcurrency = runtime.NumCPU() }
	return nil
}

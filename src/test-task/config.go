package main

import "er"
import "github.com/docopt/docopt-go"
import "strconv"

const doc = `
Usage:
	test-task [options]

Options:
	-h, --help  This help
	-l ADDR     Listen on [host]:port [default: :8080]
	-f COUNT    Fetch max COUNT urls in parallel [default: 10]
`

type Config struct {
	listen string
	fetchUrlsConcurrency int
}

func (c* Config) Load() error {
	opt, e := docopt.Parse(doc, nil, true, "", true, true)
	if e != nil { return er.Er(e, "docopt.Parse") }

	c.listen                  = opt["-l"].(string)
	c.fetchUrlsConcurrency, _ = strconv.Atoi(opt["-f"].(string))
	return nil
}

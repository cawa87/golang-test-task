package main

import "er"
import "github.com/docopt/docopt-go"

const doc = `
Usage:
	test-task [options]

Options:
	-h, --help  This help
	-l ADDR     Listen on [host]:port [default: :8080]
`

type Config struct {
	listen string
}

func (c* Config) Load() error {
	opt, e := docopt.Parse(doc, nil, true, "", true, true)
	if e != nil { return er.Er(e, "docopt.Parse") }

	c.listen = opt["-l"].(string)
	return nil
}

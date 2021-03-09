package model

import (
	"net/url"
)

// Host represents data our users want to get
type Host struct {
	Domain         *url.URL
	RequestCount   uint64
	Latency        int64
	HTTPStatusCode int
}

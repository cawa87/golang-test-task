package repository

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"github.com/kkucherenkov/golang-test-task/model"
)

//UnknownDomain error message constant
const UnknownDomain = "domain name is not found (please check sites.txt)"

type datasource struct {
	hosts              []string
	storage            map[string]model.Host
	mutex              sync.Mutex
	hostWithMaxLatency model.Host
	hostWithMinLatency model.Host
}

// New create new datasource
func New(hosts []string) Repository {
	ds := &datasource{}

	ds.hosts = hosts
	// Make store
	ds.storage = make(map[string]model.Host)

	for _, host := range hosts {
		url, err := url.Parse(host)
		if err != nil {
			continue
		}
		ds.storage[host] = model.Host{
			Domain:         url,
			RequestCount:   0,
			Latency:        0,
			HTTPStatusCode: 0,
		}
	}

	ds.mutex = sync.Mutex{}

	return ds
}

func (ds *datasource) Store(domain string, latency int64, code int) (err error) {

	ds.mutex.Lock()
	d, err := url.Parse(domain)

	if err != nil {
		return fmt.Errorf("parse domain name (%s) error - %v", domain, err)
	}

	if host, ok := ds.storage[domain]; ok {
		host.Latency = latency
		host.HTTPStatusCode = code
		host.Domain = d
		atomic.AddUint64(&host.RequestCount, 1)

		ds.storage[domain] = host

		go func() {
			_ = ds.onStore(host)
		}()

	} else {
		err = errors.New(UnknownDomain)
	}

	ds.mutex.Unlock()

	return err
}

func (ds *datasource) GetMaxLatency() (l model.Host, err error) {
	return ds.hostWithMaxLatency, nil
}

func (ds *datasource) GetMinLatency() (l model.Host, err error) {
	return ds.hostWithMinLatency, nil
}

func (ds *datasource) GetByDomain(domain string) (site model.Host, err error) {
	host, ok := ds.storage[domain]
	if ok {
		return host, nil
	}
	return host, errors.New(UnknownDomain)

}

func (ds *datasource) GetAllHosts() []string {
	return ds.hosts
}

func (ds *datasource) onStore(host model.Host) error {
	return ds.updateMetrics(host)
}

func (ds *datasource) updateMetrics(host model.Host) error {

	if ds.hostWithMaxLatency.Latency < host.Latency || ds.hostWithMaxLatency.Latency == 0 {
		ds.hostWithMaxLatency = host
	}

	if ds.hostWithMinLatency.Latency > host.Latency || ds.hostWithMinLatency.Latency == 0 {
		ds.hostWithMinLatency = host
	}

	return nil

}

func (ds *datasource) GetCount() int {
	return len(ds.hosts)
}

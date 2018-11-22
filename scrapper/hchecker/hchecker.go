package hchecker

import (
	"context"
	"errors"
	"math"
	"net/http"
	"sync"
	"time"
)

var (
	checkScheme  = "https://"
	checkTimeout = 10 * time.Second

	errNon2XX = errors.New("non-2xx status code")
)

var client = http.Client{
	Transport: &http.Transport{
		DisableKeepAlives: true,
	},
}

type HealthChecker struct {
	mu     sync.Mutex
	status map[string]Status
}

type Status struct {
	Site    string
	Live    bool
	Respond float64 // seconds
}

func New() *HealthChecker {
	return &HealthChecker{
		status: map[string]Status{},
	}
}

func (c *HealthChecker) Watch(sites []string, interval time.Duration) {
	for _, site := range sites {
		go c.watchSite(site, interval)
	}
}

func (c *HealthChecker) Status() (ss []Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, s := range c.status {
		ss = append(ss, s)
	}
	return ss
}

func (c *HealthChecker) StatusOf(site string) Status {
	c.mu.Lock()
	defer c.mu.Unlock()

	s, _ := c.status[site]
	return s
}

func (c *HealthChecker) FindMin() (min Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	min.Respond = math.MaxInt64
	for _, s := range c.status {
		if s.Live && s.Respond < min.Respond {
			min = s
		}
	}
	return min
}

func (c *HealthChecker) FindMax() (max Status) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, s := range c.status {
		if s.Live && s.Respond > max.Respond {
			max = s
		}
	}
	return max
}

func (c *HealthChecker) watchSite(site string, interval time.Duration) {
	for {
		respond, err := checkSite(site, checkTimeout)

		c.mu.Lock()
		c.status[site] = Status{
			Site:    site,
			Live:    err == nil,
			Respond: respond.Seconds(),
		}
		c.mu.Unlock()

		time.Sleep(interval)
	}
}

func checkSite(site string, timeout time.Duration) (time.Duration, error) {
	req, err := http.NewRequest(http.MethodHead, checkScheme+site, nil)
	if err != nil {
		return 0, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	start := time.Now()
	req = req.WithContext(ctx)
	res, err := client.Do(req)
	elapsed := time.Since(start)

	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	if !(200 <= res.StatusCode && res.StatusCode < 300) {
		return elapsed, errNon2XX
	}
	return elapsed, nil
}

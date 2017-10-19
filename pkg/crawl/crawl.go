package crawl

import (
	"bytes"
	"errors"
	"io"
	"mime"

	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

const maxRedirectsCount = 20

var (
	// ErrTooManyRedirects too many redirects
	ErrTooManyRedirects = errors.New("too many redirects")
	// ErrMissingLocation missing location
	ErrMissingLocation = errors.New("missing location")
)

var (
	requestMethod = []byte("GET")
	strLocation   = []byte("Location")
)

var emptyResult = Result{}

// Result parsing result.
type Result struct {
	URL      string    `json:"url"`
	Meta     Meta      `json:"meta"`
	Elements []Element `json:"elements"`
	Err      error     `json:"-"`
}

// Meta meta information of the requested page.
type Meta struct {
	Status        int    `json:"status"`
	ContentType   string `json:"content_type,omitempty"`
	ContentLength int    `json:"content_length,omitempty"`
}

// Element statistics on the tags of the loaded page.
type Element struct {
	TagName string `json:"tag_name"`
	Count   int    `json:"count"`
}

// Crawl scrape page interface.
type Crawl interface {
	Do(url string) (Result, error)
}

type basicCrawl struct {
	client *fasthttp.Client
}

func (c *basicCrawl) getPageTagCounts(b io.Reader) map[string]int {
	z := html.NewTokenizer(b)

	tagCounter := map[string]int{}

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			break
		}
		switch {
		case tt == html.StartTagToken:
			t := z.Token()
			if _, ok := tagCounter[t.Data]; !ok {
				tagCounter[t.Data] = 1
			} else {
				tagCounter[t.Data]++
			}
		}
	}
	return tagCounter
}

func (c *basicCrawl) getRedirectURL(baseURL string, location []byte) string {
	u := fasthttp.AcquireURI()
	u.Update(baseURL)
	u.UpdateBytes(location)
	redirectURL := u.String()
	fasthttp.ReleaseURI(u)
	return redirectURL
}

func (c *basicCrawl) doRequestFollowRedirects(req *fasthttp.Request, url string) (*fasthttp.Response, error) {
	resp := fasthttp.AcquireResponse()

	redirectsCount := 0
	for {
		req.SetRequestURI(url)
		if err := c.client.Do(req, resp); err != nil {
			return nil, err

		}
		if statusCode := resp.Header.StatusCode(); statusCode != fasthttp.StatusMovedPermanently &&
			statusCode != fasthttp.StatusFound &&
			statusCode != fasthttp.StatusSeeOther {
			break
		}

		redirectsCount++
		if redirectsCount > maxRedirectsCount {
			return nil, ErrTooManyRedirects
		}
		location := resp.Header.PeekBytes(strLocation)
		if len(location) == 0 {
			return nil, ErrMissingLocation
		}
		url = c.getRedirectURL(url, location)
	}
	return resp, nil
}

// Do start scrape.
func (c *basicCrawl) Do(url string) (Result, error) {
	req := fasthttp.AcquireRequest()

	defer fasthttp.ReleaseRequest(req)

	req.Header.SetMethodBytes(requestMethod)
	resp, err := c.doRequestFollowRedirects(req, url)

	if resp != nil {
		defer fasthttp.ReleaseResponse(resp)
	}

	if err != nil {
		return emptyResult, err
	}

	contentType, _, err := mime.ParseMediaType(string(resp.Header.ContentType()))
	if err != nil {
		return emptyResult, err
	}

	result := Result{
		URL: url,
		Meta: Meta{
			Status:        resp.StatusCode(),
			ContentLength: resp.Header.ContentLength(),
			ContentType:   contentType,
		},
		Elements: []Element{},
	}

	if resp.StatusCode() == fasthttp.StatusOK {
		tagCounts := c.getPageTagCounts(
			bytes.NewBuffer(resp.Body()),
		)

		for tagName, count := range tagCounts {
			result.Elements = append(result.Elements, Element{
				TagName: tagName,
				Count:   count,
			})
		}
	}
	return result, nil
}

// New create crawl instance.
func New() Crawl {
	return &basicCrawl{
		client: &fasthttp.Client{},
	}
}

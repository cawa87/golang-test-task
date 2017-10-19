package crawl_test

import (
	"testing"

	"github.com/l-vitaly/golang-test-task/pkg/crawl"
)

func TestCrawl(t *testing.T) {
	c := crawl.New()

	result, err := c.Do("http://ya.ru")

	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if want, have := "http://ya.ru", result.URL; want != have {
		t.Errorf("want %s, have %s", want, have)
	}

	if want, have := "text/html", result.Meta.ContentType; want != have {
		t.Errorf("want %s, have %s", want, have)
	}

	if len(result.Elements) == 0 {
		t.Errorf("want > 0, have 0")
	}
}

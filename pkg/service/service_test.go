package service_test

import (
	"testing"

	"github.com/l-vitaly/golang-test-task/pkg/crawl"

	"github.com/l-vitaly/golang-test-task/pkg/service"
)

func TestServicePostURLs(t *testing.T) {

	c := crawl.New()

	s := service.NewBasicService(c, 100)

	results := s.PostURLs([]string{"http://lenta.ru", "http://ya.ru"})

	if want, have := 2, len(results); want != have {
		t.Errorf("want %d, have %d", want, have)
	}
}

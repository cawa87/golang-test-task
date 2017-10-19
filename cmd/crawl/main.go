package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/l-vitaly/golang-test-task/pkg/crawl"
	"github.com/valyala/fasthttp"

	"github.com/go-kit/kit/log"
	"github.com/l-vitaly/golang-test-task/pkg/config"
	"github.com/l-vitaly/golang-test-task/pkg/endpoint"
	"github.com/l-vitaly/golang-test-task/pkg/service"
	"github.com/l-vitaly/golang-test-task/pkg/transport"
)

func main() {
	cfg := config.Parse()

	var logger log.Logger
	logger = log.NewJSONLogger(os.Stdout)
	defer logger.Log("msg", "goodbye")

	logger = log.With(logger, "@timestamp", log.DefaultTimestampUTC)
	logger = log.With(logger, "@message", "info")
	logger = log.With(logger, "caller", log.DefaultCaller)

	logger.Log("msg", "hello")

	errCh := make(chan error)

	c := crawl.New()
	svc := service.New(c, cfg.MaxWorkers, log.With(logger, "component", "Service"))
	e := endpoint.New(svc)
	r := transport.NewHTTPHandler(e, log.With(logger, "component", "HTTP"))

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errCh <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		// HTTP handler.
		logger.Log("transport", "HTTP", "bind-addr", cfg.BindAddr)
		errCh <- fasthttp.ListenAndServe(cfg.BindAddr, r.HandleRequest)
	}()

	logger.Log("info", "started")
	logger.Log("err", <-errCh)
}

package transport

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
	httptransport "github.com/l-vitaly/go-kit/transport/fasthttp"
	"github.com/l-vitaly/golang-test-task/pkg/endpoint"
	"github.com/l-vitaly/golang-test-task/pkg/service"
	"github.com/pquerna/ffjson/ffjson"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var errBadRequest = errors.New("bad request")

type failer interface {
	Fail() error
}

// NewHTTPHandler new server http handler.
func NewHTTPHandler(endpoints endpoint.Set, logger log.Logger) *routing.Router {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(func(ctx context.Context, err error, r *fasthttp.Response) {
			encodeError(err, r)
		}),
	}

	r := routing.New()

	var v1 *routing.RouteGroup
	{
		v1 = r.Group("/v1")

		var urls *routing.RouteGroup
		{
			urls = v1.Group("")
			urls.Post("", httptransport.NewServer(
				endpoints.PostURLsEndpoint,
				decodePostURLsRequest,
				encodePostURLsResponse,
				options...,
			).RouterHandle())
		}
	}
	return r
}

func decodePostURLsRequest(ctx context.Context, r *fasthttp.Request) (interface{}, error) {
	var urls []string
	if len(r.Body()) == 0 {
		return nil, errBadRequest
	}
	if err := ffjson.Unmarshal(r.Body(), &urls); err != nil {
		return nil, err
	}
	return urls, nil
}

func encodePostURLsResponse(ctx context.Context, r *fasthttp.Response, response interface{}) error {
	if e, ok := response.(failer); ok && e.Fail() != nil {
		encodeError(e.Fail(), r)
		return nil
	}
	r.Header.Set("Content-Type", "application/json; charset=utf-8")

	resp := response.(endpoint.PostURLsResponse)

	b, err := ffjson.Marshal(resp.Data)
	if err != nil {
		return err
	}
	r.SetBody(b)
	return nil
}

func encodeError(err error, r *fasthttp.Response) {
	if err == nil {
		panic("encodeError with nil error")
	}
	r.Header.Set("Content-Type", "application/json; charset=utf-8")
	r.SetStatusCode(codeFrom(err))
	b, err := ffjson.Marshal(map[string]string{
		"error": err.Error(),
	})
	if err != nil {
		r.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	r.SetBody(b)
}

func codeFrom(err error) int {
	switch err {
	case service.ErrEmptyURLs, errBadRequest:
		return fasthttp.StatusBadRequest
	default:
		return fasthttp.StatusInternalServerError
	}
}

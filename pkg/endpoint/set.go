package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/l-vitaly/golang-test-task/pkg/crawl"
	"github.com/l-vitaly/golang-test-task/pkg/service"
)

// Set collects all of the endpoints that compose an crawl service.
type Set struct {
	PostURLsEndpoint endpoint.Endpoint
}

// New returns a Endpoints.
func New(svc service.Service) Set {
	var postURLsEndpoint endpoint.Endpoint
	{
		postURLsEndpoint = makePostURLsEndpoint(svc)
	}
	return Set{
		PostURLsEndpoint: postURLsEndpoint,
	}
}

func makePostURLsEndpoint(s service.Service) endpoint.Endpoint {
	return func(_ context.Context, request interface{}) (interface{}, error) {
		results, err := s.PostURLs(request.([]string))
		return PostURLsResponse{Data: results, Err: err}, nil
	}
}

type PostURLsResponse struct {
	Data []crawl.Result `json:"data"`
	Err  error          `json:"error,omitempty"`
}

func (r PostURLsResponse) Fail() error {
	return r.Err
}

package loadbalancer

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-openapi/errors"
)

type Picker interface {
	Pick(backends []*Backend) (*Backend, error)
}

type loadBalancer struct {
	log      *slog.Logger
	picker   Picker
	backends []*Backend
}

// New creates and returns http server.
func New(
	log *slog.Logger,
	p Picker,
	backends []*Backend,
	host string,
	port int,
) *http.Server {
	lb := loadBalancer{ //nolint:exhaustruct
		log:      log,
		picker:   p,
		backends: backends,
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: &lb,
	}

	for _, b := range backends {
		b.ReverseProxy.ErrorHandler = lb.MarkUnavailable
	}

	return srv
}

var ErrBackendsUnavailable = errors.New(
	503, //nolint:mnd
	"none of backends is available",
)

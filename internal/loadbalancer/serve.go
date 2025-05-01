package loadbalancer

import (
	"errors"
	"net/http"
	"net/http/httputil"
	"net/url"

	api "github.com/go-openapi/errors"
	"go.uber.org/atomic"
)

type Backend struct {
	URL          *url.URL
	Alive        *atomic.Bool
	ReverseProxy *httputil.ReverseProxy
	Conns        *atomic.Uint64
}

func (b *loadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	backend, err := b.picker.Pick(b.backends)
	switch {
	case errors.Is(err, ErrBackendsUnavailable):
		b.log.Error("unable to pick:", "error", err)
		api.ServeError(w, r, ErrBackendsUnavailable)
		return
	}

	backend.Conns.Inc()
	defer backend.Conns.Dec()

	b.log.Info("requesting backend:", "url", backend.URL.String())
	backend.ReverseProxy.ServeHTTP(w, r)
}

package loadbalancer

import (
	"net/http"
	"net/url"
	"slices"
)

func (b *loadBalancer) MarkUnavailable(
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	b.log.Error("backend doesn't respond:", "error", err)

	backend := searchBackend(b.backends, r.URL)
	backend.Alive.Store(false)

	// Call balancer again,
	// without problematic server
	b.ServeHTTP(w, r)
}

// searchBackend is a wrapper for convenience.
func searchBackend(in []*Backend, toFind *url.URL) *Backend {
	idx := slices.IndexFunc(in, func(backend *Backend) bool { return backend.URL.Host == toFind.Host })
	return in[idx]
}

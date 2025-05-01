package algorithms

import (
	"math/rand"
	"sync/atomic"

	"github.com/tikhonfedulov/silly-balancer/internal/loadbalancer"
)

type leastConnections struct{}

func (l *leastConnections) Pick(backends []*loadbalancer.Backend) (*loadbalancer.Backend, error) {
	available := availableBackends(backends)
	if len(available) == 0 {
		return nil, loadbalancer.ErrBackendsUnavailable
	}

	lowest := available[0]
	for _, b := range available {
		if b.Conns.Load() < lowest.Conns.Load() {
			lowest = b
		}
	}

	return lowest, nil
}

type roundRobin struct {
	current uint64
}

func (r *roundRobin) Pick(backends []*loadbalancer.Backend) (*loadbalancer.Backend, error) {
	available := availableBackends(backends)
	if len(available) == 0 {
		return nil, loadbalancer.ErrBackendsUnavailable
	}

	pos := atomic.AddUint64(&r.current, 1) - 1
	return available[pos%uint64(len(available))], nil
}

type random struct{}

func (r *random) Pick(backends []*loadbalancer.Backend) (*loadbalancer.Backend, error) {
	available := availableBackends(backends)
	if len(available) == 0 {
		return nil, loadbalancer.ErrBackendsUnavailable
	}

	i := rand.Intn(len(available))
	return available[i], nil
}

// availableBackends is a wrapper for convenience.
func availableBackends(backends []*loadbalancer.Backend) []*loadbalancer.Backend {
	out := make([]*loadbalancer.Backend, 0, len(backends))
	for _, b := range backends {
		if b.Alive.Load() {
			out = append(out, b)
		}
	}

	return out
}

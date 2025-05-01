package health

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/tikhonfedulov/silly-balancer/internal/loadbalancer"
)

type health struct {
	log      *slog.Logger
	backends []*loadbalancer.Backend
	path     string
}

func New(
	log *slog.Logger,
	b []*loadbalancer.Backend,
	path string,
) *health {
	h := &health{ //nolint:exhaustruct
		log:      log,
		backends: b,
		path:     path,
	}

	return h
}

func (h *health) Start() {
	go func() {
		for {
			for _, b := range h.backends {
				u := *b.URL
				u.Path = h.path

				resp, err := http.Get(u.String())
				switch {
				case err != nil:
					h.log.Error("health error", "url", u.String(), "error", err)
					b.Alive.Store(false)
					continue
				case resp.StatusCode != http.StatusOK:
					h.log.Warn("unhealthy server", "url", u.String(), "status", resp.Status)
					b.Alive.Store(false)
					continue
				}

				b.Alive.Store(true)
				resp.Body.Close() //nolint:errcheck
			}

			time.Sleep(time.Minute)
		}
	}()
}

package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http/httputil"
	"net/url"
	"os"

	"go.uber.org/atomic"

	"github.com/tikhonfedulov/silly-balancer/cmd/server/internal/config"
	"github.com/tikhonfedulov/silly-balancer/internal/health"
	"github.com/tikhonfedulov/silly-balancer/internal/loadbalancer"
	"github.com/tikhonfedulov/silly-balancer/internal/loadbalancer/algorithms"
	"github.com/tikhonfedulov/silly-balancer/internal/logger"
)

//nolint:gochecknoglobals
var (
	cfgFile = flag.String("cfg", "./config.yml", "path to Config file")
	debug   = flag.Bool("debug", false, "enable debug logging")
)

func main() {
	flag.Parse()

	var (
		cfg      = must(config.ReadYaml[config.Config](*cfgFile))
		log      = logger.New(*debug)
		picker   = must(algorithms.New(cfg.LoadBalancer.Algorithm))
		backends = must(build(cfg.Backends))
		restorer = health.New(log, backends, cfg.Health.Path)
		balancer = loadbalancer.New(
			log,
			picker,
			backends,
			cfg.LoadBalancer.Host,
			cfg.LoadBalancer.Port,
		)
	)

	log.Debug("started", slog.Any("config", cfg))
	restorer.Start()

	if err := balancer.ListenAndServe(); err != nil {
		log.Error("shutdown", slog.String("error", err.Error()))
		os.Exit(1)
	}
}

func must[T any](v T, err error) T { //nolint:ireturn
	if err != nil {
		panic(err)
	}

	return v
}

func build(loaded []config.Backend) ([]*loadbalancer.Backend, error) {
	out := make([]*loadbalancer.Backend, len(loaded))
	for i := range loaded {
		url, err := url.Parse(loaded[i].URL)
		if err != nil {
			return nil, fmt.Errorf("invalid url: %w", err)
		}

		out[i] = &loadbalancer.Backend{
			URL:          url,
			Alive:        atomic.NewBool(true),
			ReverseProxy: httputil.NewSingleHostReverseProxy(url),
			Conns:        atomic.NewUint64(0),
		}
	}

	return out, nil
}

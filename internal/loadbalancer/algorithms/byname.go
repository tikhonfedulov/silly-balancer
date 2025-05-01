//nolint:all
package algorithms

import (
	"fmt"

	"github.com/tikhonfedulov/silly-balancer/internal/loadbalancer"
)

func New(name string) (loadbalancer.Picker, error) { //nolint:ireturn
	switch name {
	case "random":
		return &random{}, nil
	case "least-connections":
		return &leastConnections{}, nil
	case "round-robin":
		return &roundRobin{}, nil
	default:
		return nil, fmt.Errorf("unexpected algorithm: %s", name)
	}
}

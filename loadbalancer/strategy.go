package loadbalancer

import (
	"errors"

	"github.com/mkauppila/load-balancer/types"
)

var errNoHealthyServers = errors.New("no healthy servers")

type Strategy interface {
	nextHealthyServer() (*types.Server, error)
}

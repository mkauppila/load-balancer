package loadbalancer

import "github.com/mkauppila/load-balancer/types"

type Strategy interface {
	nextHealthyServer() (*types.Server, error)
}

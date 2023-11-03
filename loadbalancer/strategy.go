package loadbalancer

import "github.com/mkauppila/load-balancer/types"

type Strategy interface {
	getNextServer() (*types.Server, error)
}

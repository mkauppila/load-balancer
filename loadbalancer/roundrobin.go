package loadbalancer

import (
	"container/ring"
	"fmt"
	"sync"

	"github.com/mkauppila/load-balancer/types"
)

type RoundRobin struct {
	mu      sync.Mutex
	servers *ring.Ring
}

func CreateRoundRobin(servers []*types.Server) *RoundRobin {
	r := ring.New(len(servers))
	for _, s := range servers {
		r.Value = s
		r = r.Next()
	}
	rr := RoundRobin{}
	rr.servers = r
	return &rr
}

func (r *RoundRobin) nextHealthyServer() (*types.Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// would be nice to iterate the ring, pick the first alive server
	// and move the head accordingly
	retryCounter := r.servers.Len()
	for {
		r.servers = r.servers.Move(1)
		server := r.servers.Value.(*types.Server)
		if server.IsHealthy {
			return server, nil
		}

		retryCounter--
		if retryCounter <= 0 {
			return nil, fmt.Errorf("all servers are dead")
		}
	}
}

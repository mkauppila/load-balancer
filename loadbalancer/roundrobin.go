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

	retryCounter := r.servers.Len()
	for {
		server := r.servers.Value.(*types.Server)
		r.servers = r.servers.Move(1)
		if server.IsHealthy {
			return server, nil
		}

		retryCounter--
		if retryCounter <= 0 {
			return nil, fmt.Errorf("all servers are dead")
		}
	}
}

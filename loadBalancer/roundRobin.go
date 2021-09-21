package loadBalancer

import (
	"container/ring"
	"fmt"
	"sync"
)

type RoundRobin struct {
	mu      sync.Mutex
	servers *ring.Ring
}

func CreateRoundRobin(servers []*Server) *RoundRobin {
	r := ring.New(len(servers))
	for _, s := range servers {
		r.Value = s
		r = r.Next()
	}
	rr := RoundRobin{}
	rr.servers = r
	return &rr
}

func (r *RoundRobin) getNextServer() (*Server, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// would be nice to iterate the ring, pick the first alive server
	// and move the head accordingly
	retryCounter := r.servers.Len()
	for {
		r.servers = r.servers.Move(1)
		server := r.servers.Value.(*Server)
		if server.isHealthy {
			return server, nil
		}

		retryCounter--
		if retryCounter <= 0 {
			return nil, fmt.Errorf("all servers are dead")
		}
	}
}

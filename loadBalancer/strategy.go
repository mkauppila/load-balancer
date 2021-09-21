package loadBalancer

import (
	"container/ring"
	"fmt"
	"math/rand"
	"sync"
)

type Strategy interface {
	getNextServer() (*Server, error)
}

type RoundRobin struct {
	mu      sync.Mutex
	servers *ring.Ring
}

func CreateRoundRobin(servers []*Server) *RoundRobin {
	rr := RoundRobin{}
	r := ring.New(len(servers))
	for _, s := range servers {
		r.Value = s
		r = r.Next()
	}
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

type Random struct {
	servers []*Server
}

func CreateRandom(servers []*Server) *Random {
	return &Random{servers}
}

func (r *Random) getNextServer() (*Server, error) {
	var aliveServers []*Server
	for _, server := range r.servers {
		if server.isHealthy {
			aliveServers = append(aliveServers, server)
		}
	}

	if len(aliveServers) == 0 {
		return nil, fmt.Errorf("all servers are dead")
	}

	index := rand.Intn(len(aliveServers))
	return aliveServers[index], nil
}

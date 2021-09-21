package loadBalancer

import (
	"container/ring"
	"fmt"
)

type Strategy interface {
	getNextServer() (*Server, error)
}

type RoundRobin struct {
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
	retryCounter := r.servers.Len()
	for {
		r.servers = r.servers.Move(1)
		server := r.servers.Value.(*Server)
		if server.isHealthy {
			fmt.Println("this was chosen: ", server)
			return server, nil
		}

		retryCounter--
		if retryCounter <= 0 {
			return nil, fmt.Errorf("all servers are dead")
		}
	}
}

package loadBalancer

import (
	"fmt"
	"math/rand"
)

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

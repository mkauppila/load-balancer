package loadbalancer

import (
	"fmt"
	"math/rand"

	"github.com/mkauppila/load-balancer/types"
)

type Random struct {
	servers []*types.Server
}

func CreateRandom(servers []*types.Server) *Random {
	return &Random{servers}
}

func (r *Random) nextHealthyServer() (*types.Server, error) {
	var aliveServers []*types.Server
	for _, server := range r.servers {
		if server.IsHealthy {
			aliveServers = append(aliveServers, server)
		}
	}

	if len(aliveServers) == 0 {
		return nil, fmt.Errorf("all servers are dead")
	}

	index := rand.Intn(len(aliveServers))
	return aliveServers[index], nil
}

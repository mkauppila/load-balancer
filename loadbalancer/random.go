package loadbalancer

import (
	"math/rand"

	"github.com/mkauppila/load-balancer/types"
)

type Random struct {
	servers []*types.Server
	rand    *rand.Rand
}

func CreateRandom(servers []*types.Server, rand *rand.Rand) *Random {
	return &Random{servers: servers, rand: rand}
}

func (r *Random) nextHealthyServer() (*types.Server, error) {
	var aliveServers []*types.Server
	for _, server := range r.servers {
		if server.IsHealthy {
			aliveServers = append(aliveServers, server)
		}
	}

	if len(aliveServers) == 0 {
		return nil, errNoHealthyServers
	}

	index := r.rand.Intn(len(aliveServers))
	return aliveServers[index], nil
}

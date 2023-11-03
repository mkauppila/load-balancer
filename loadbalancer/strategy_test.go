package loadbalancer

import (
	"fmt"
	"testing"

	"github.com/mkauppila/load-balancer/types"
)

func TestGettingServerRoundRobin(t *testing.T) {
	servers := []*types.Server{
		{Url: "url", IsHealthy: true},
		{Url: "url2", IsHealthy: true},
	}
	rr := CreateRoundRobin(servers)

	servers[0].IsHealthy = false
	s, _ := rr.getNextServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}

	s, _ = rr.getNextServer()
	fmt.Println("1: ", s)
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}
}

package loadbalancer

import (
	"errors"
	"testing"

	"github.com/mkauppila/load-balancer/types"
)

func TestRoundRobinLoopThroughServers(t *testing.T) {
	servers := []*types.Server{
		{Url: "url1", IsHealthy: true},
		{Url: "url2", IsHealthy: true},
		{Url: "url3", IsHealthy: true},
	}
	rr := CreateRoundRobin(servers)

	s, _ := rr.nextHealthyServer()
	if s.Url != "url1" {
		t.Errorf("Wrong server. Expected url1, actual %s", s.Url)
	}

	s, _ = rr.nextHealthyServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}

	s, _ = rr.nextHealthyServer()
	if s.Url != "url3" {
		t.Errorf("Wrong server. Expected url3, actual %s", s.Url)
	}
}

func TestRoundRobinGiveErrorIfAllServersAreUnhealthy(t *testing.T) {
	servers := []*types.Server{
		{Url: "url1", IsHealthy: false},
		{Url: "url2", IsHealthy: false},
		{Url: "url3", IsHealthy: false},
	}
	rr := CreateRoundRobin(servers)

	_, err := rr.nextHealthyServer()
	if !errors.Is(err, errNoHealthyServers) {
		t.Errorf("Expected an error")
	}
}

func TestRoundRobinSkipOverUnhealthyServer(t *testing.T) {
	servers := []*types.Server{
		{Url: "url", IsHealthy: true},
		{Url: "url2", IsHealthy: true},
	}
	rr := CreateRoundRobin(servers)

	servers[0].IsHealthy = false
	s, _ := rr.nextHealthyServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}

	s, _ = rr.nextHealthyServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}
}

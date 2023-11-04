package loadbalancer

import (
	"errors"
	"math/rand"
	"testing"

	"github.com/mkauppila/load-balancer/types"
)

func TestRandomLoopThroughServers(t *testing.T) {
	servers := []*types.Server{
		{Url: "url1", IsHealthy: true},
		{Url: "url2", IsHealthy: true},
		{Url: "url3", IsHealthy: true},
	}
	println(servers)

	r := rand.New(rand.NewSource(100))
	rr := CreateRandom(servers, r)

	// Based on the given rand seed this is the order of
	s, _ := rr.nextHealthyServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url1, actual %s", s.Url)
	}

	s, _ = rr.nextHealthyServer()
	if s.Url != "url3" {
		t.Errorf("Wrong server. Expected url2, actual %s", s.Url)
	}

	s, _ = rr.nextHealthyServer()
	if s.Url != "url2" {
		t.Errorf("Wrong server. Expected url3, actual %s", s.Url)
	}
}
func TestRandomGiveErrorIfAllServersAreUnhealthy(t *testing.T) {
	servers := []*types.Server{
		{Url: "url1", IsHealthy: false},
		{Url: "url2", IsHealthy: false},
		{Url: "url3", IsHealthy: false},
	}
	rr := CreateRandom(servers, rand.New(rand.NewSource(100)))

	_, err := rr.nextHealthyServer()
	if !errors.Is(err, errNoHealthyServers) {
		t.Errorf("Expected an error")
	}
}

func TestRandomSkipOverUnhealthyServer(t *testing.T) {
	servers := []*types.Server{
		{Url: "url", IsHealthy: true},
		{Url: "url2", IsHealthy: true},
	}
	rr := CreateRandom(servers, rand.New(rand.NewSource(100)))

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

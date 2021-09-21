package loadBalancer

import (
	"fmt"
	"testing"
)

func TestGettingServerRoundRobin(t *testing.T) {
	servers := []*Server{
		{Url: "url", isHealthy: true},
		{Url: "url2", isHealthy: true},
	}
	rr := CreateRoundRobin(servers)

	servers[0].isHealthy = false
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

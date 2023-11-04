package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/loadbalancer"
	"github.com/mkauppila/load-balancer/test/httpserver"
	"github.com/mkauppila/load-balancer/types"
)

func TestLoadBalancerRoundRobin(t *testing.T) {
	ctx := context.Background()

	httpServerBasePort := 50_000

	cfg := config.Configuration{}
	for c := 0; c < 3; c++ {
		cfg.Servers = append(cfg.Servers, types.Server{
			Url: fmt.Sprintf("http://localhost:%d", httpServerBasePort+c),
		})
	}
	cfg.HealthCheck = types.HealthCheck{
		Enabled:    false,
		IntervalMs: 10,
		Path:       "/health",
	}
	cfg.Strategy = types.RoundRobin
	cfg.Port = 40_000

	fmt.Println("Setting up the target HTTP servers...")
	var wg sync.WaitGroup
	for i, server := range cfg.Servers {
		wg.Add(1)
		response := fmt.Sprintf("response %d", i)
		t.Log("gen response ", response)
		go func(server types.Server, response string) {
			defer wg.Done()
			readyCtx, cancel := context.WithCancel(context.Background())
			httpserver.RunServer(cancel, ctx, server.Url, response)
			<-readyCtx.Done()
		}(server, response)
	}
	wg.Wait()
	fmt.Println("Target HTTP servers are up and ready")

	srv := loadbalancer.NewLoadBalancer(cfg)
	cancel := srv.Start(ctx)
	// TODO Why does this need to start with http://
	lbAddr := fmt.Sprintf("http://localhost:%d", cfg.Port)
	request, _ := http.NewRequest(http.MethodGet, lbAddr, nil)
	// TODO: check for no error
	response := httptest.NewRecorder()

	srv.ForwardRequest(response, request)
	got := response.Code
	if got != 200 {
		t.Errorf("The request was not success. Got %d", response.Code)
	}

	cancel()
}

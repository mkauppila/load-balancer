package loadbalancer

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/httpserver"
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
	cfg.Strategy = types.StrategyRoundRobin
	cfg.Port = 40_000

	fmt.Println("Setting up the target HTTP servers...")
	var wg sync.WaitGroup
	for i, server := range cfg.Servers {
		wg.Add(1)
		response := fmt.Sprintf("response %d", i)
		go func(server types.Server, response string) {
			defer wg.Done()
			readyCtx, cancel := context.WithCancel(context.Background())
			httpserver.RunServer(cancel, ctx, server.Url, response)
			<-readyCtx.Done()
		}(server, response)
	}
	wg.Wait()
	fmt.Println("Target HTTP servers are up and ready")

	expectedResponses := []string{
		// First loop
		"response 0",
		"response 1",
		"response 2",
		// Second loop
		"response 0",
		"response 1",
		"response 2",
	}

	srv := NewLoadBalancer(cfg)
	cancel := srv.Start(ctx)
	lbAddr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	for _, expectedResponse := range expectedResponses {
		// TODO Why does this need to start with http://
		request, err := http.NewRequest(http.MethodGet, lbAddr, nil)
		if err != nil {
			t.Errorf("Failed to create new HTTP request")
		}
		response := httptest.NewRecorder()

		srv.ForwardRequest(response, request)
		if response.Code != 200 {
			t.Errorf("The request was not success. Got %d", response.Code)
		}

		actual := response.Body.String()
		if expectedResponse != actual {
			t.Errorf("Expected '%s', got '%s'", expectedResponse, actual)
		}
	}

	cancel()
}

func TestLoadBalancerRandom(t *testing.T) {
	ctx := context.Background()

	httpServerBasePort := 55_000

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
	cfg.Strategy = types.StrategyRandom
	cfg.Port = 40_000

	fmt.Println("Setting up the target HTTP servers...")
	var wg sync.WaitGroup
	for i, server := range cfg.Servers {
		wg.Add(1)
		response := fmt.Sprintf("response %d", i)
		go func(server types.Server, response string) {
			defer wg.Done()
			readyCtx, cancel := context.WithCancel(context.Background())
			httpserver.RunServer(cancel, ctx, server.Url, response)
			<-readyCtx.Done()
		}(server, response)
	}
	wg.Wait()
	fmt.Println("Target HTTP servers are up and ready")

	expectedResponses := []string{
		// First loop
		"response 1",
		"response 2",
		"response 2",
		// Second loop
		"response 1",
		"response 2",
		"response 1",
	}

	srv := NewLoadBalancer(cfg)
	cancel := srv.Start(ctx)
	lbAddr := fmt.Sprintf("http://localhost:%d", cfg.Port)

	for _, expectedResponse := range expectedResponses {
		// TODO Why does this need to start with http://
		request, err := http.NewRequest(http.MethodGet, lbAddr, nil)
		if err != nil {
			t.Errorf("Failed to create new HTTP request")
		}
		response := httptest.NewRecorder()

		srv.ForwardRequest(response, request)
		if response.Code != 200 {
			t.Errorf("The request was not success. Got %d", response.Code)
		}

		actual := response.Body.String()
		if expectedResponse != actual {
			t.Errorf("Expected '%s', got '%s'", expectedResponse, actual)
		}
	}

	cancel()
}

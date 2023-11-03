package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/lb"
	"github.com/mkauppila/load-balancer/test/httpserver"
)

func TestSomething(t *testing.T) {
	ctx := context.Background()

	url := "http://localhost:50000"

	cfg := config.Configuration{
		Servers: []config.Server{
			{
				Url: url,
			},
		},
		HealthCheck: config.HealthCheck{
			Enabled:    false,
			IntervalMs: 10,
			Path:       "/health",
		},
		Strategy: "round-robin",
	}

	fmt.Println("Setting up the target HTTP servers...")
	var wg sync.WaitGroup
	for _, server := range cfg.Servers {
		wg.Add(1)
		go func(server config.Server) {
			defer wg.Done()
			readyCtx, cancel := context.WithCancel(context.Background())
			httpserver.RunServer(cancel, ctx, server.Url)
			<-readyCtx.Done()
		}(server)
	}
	wg.Wait()
	fmt.Println("Target HTTP servers are up and ready")

	srv := lb.NewLoadBalancer(cfg)
	cancel := srv.Start(ctx)
	request, _ := http.NewRequest(http.MethodGet, url+"/", nil)
	// TODO: check for no error
	response := httptest.NewRecorder()

	srv.ForwardRequest(response, request)
	got := response.Code
	if got != 200 {
		t.Errorf("The request was not success. Got %d", response.Code)
	}

	cancel()
}

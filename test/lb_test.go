package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

	httpserver.RunServer(ctx, cfg.Servers[0].Url)

	time.Sleep(1 * time.Second)

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

	time.Sleep(1 * time.Second)

	cancel()
}

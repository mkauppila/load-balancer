package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/lb"
)

func TestSomething(t *testing.T) {
	url := "http://localhost:50000"
	cfg := config.Configuration{
		Servers: []config.Server{
			{
				Url: url,
			},
		},
		HealthCheck: config.HealthCheck{
			Enabled:    true,
			IntervalMs: 10,
			Path:       "/health",
		},
		Strategy: "round-robin",
	}

	srv := lb.NewLoadBalancer(cfg)
	cancel := srv.Start(context.Background())

	request, _ := http.NewRequest(http.MethodGet, url, nil)
	// TODO: check for no error
	response := httptest.NewRecorder()

	srv.ForwardRequest(response, request)
	got := response.Code
	if got != 200 {
		t.Errorf("The request was not success. Got %d", response.Code)
	}

	cancel()
}

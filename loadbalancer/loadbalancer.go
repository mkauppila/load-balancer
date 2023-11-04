package loadbalancer

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/types"
)

type LoadBalancer struct {
	allServers  []*types.Server
	port        int
	healthCheck types.HealthCheck
	strategy    Strategy
}

func NewLoadBalancer(conf config.Configuration) LoadBalancer {
	var servers []*types.Server
	for i := 0; i < len(conf.Servers); i++ {
		server := types.Server{Url: conf.Servers[i].Url, IsHealthy: true}
		servers = append(servers, &server)
	}

	loadBalancer := LoadBalancer{
		allServers: servers,
		healthCheck: types.HealthCheck{
			Enabled:    conf.HealthCheck.Enabled,
			Path:       conf.HealthCheck.Path,
			IntervalMs: conf.HealthCheck.IntervalMs,
		},
		port: conf.Port,
	}

	switch conf.Strategy {
	case "random":
		loadBalancer.strategy = CreateRandom(servers)
	case "round-robin":
		loadBalancer.strategy = CreateRoundRobin(servers)
	default:
		fmt.Println("Invalid load balancing strategy. Defaulting to round robin.")
		loadBalancer.strategy = CreateRoundRobin(servers)
	}

	if loadBalancer.healthCheck.Enabled {
		for _, server := range loadBalancer.allServers {
			fmt.Println("Kick up health check for ", server.Url)
			go loadBalancer.doHealthCheck(server)
		}
	}

	return loadBalancer
}

func (b *LoadBalancer) Start(ctx context.Context) context.CancelFunc {
	mux := http.NewServeMux()
	mux.HandleFunc("/", b.ForwardRequest)

	_, cancel := context.WithCancel(ctx)
	go func() {
		// TODO: Allow changing the addr
		addr := fmt.Sprintf("localhost:%d", b.port)
		fmt.Println("Start listening on addr:", addr)
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			log.Println(err)
			cancel()
			return
		}
	}()

	return cancel
}

func (b *LoadBalancer) doHealthCheck(server *types.Server) {
	for {
		client := http.DefaultClient
		response, err := client.Get(server.Url + b.healthCheck.Path)
		if err != nil {
			fmt.Println("Health check failed for", server.Url)
			server.IsHealthy = false
		} else {
			if response.StatusCode == http.StatusOK {
				fmt.Println("Health check OK for ", server.Url)
				server.IsHealthy = true
			} else {
				server.IsHealthy = false
				fmt.Println("Health check failed for", server.Url,
					" wrong response status ", response.StatusCode)
			}
		}

		time.Sleep(time.Millisecond * time.Duration(b.healthCheck.IntervalMs))
	}
}

func (b *LoadBalancer) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	server, err := b.strategy.getNextServer()
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	req, err := http.NewRequestWithContext(r.Context(), r.Method, server.Url, r.Body)
	if err != nil {
		fmt.Println("Invalid HTTP request: ", err)
		return
	}

	req.Header = r.Header
	req.Header.Add("X-Forwarded-For", r.RemoteAddr)
	// TODO: Add rest of the custom headers

	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error: ", err)
		w.WriteHeader(http.StatusBadGateway)
		w.Header().Add("Content-Type", "application/text")
		data := []byte(err.Error())
		b, err := w.Write(data)
		if err != nil || b != len(data) {
			return
		}
		return
	}

	_, err = io.Copy(w, response.Body)
	if err != nil {
		err := response.Body.Close()
		if err != nil {
			return
		}
		return
	}

	err = response.Body.Close()
	if err != nil {
		return
	}
}

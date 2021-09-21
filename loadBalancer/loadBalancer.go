package loadBalancer

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/mkauppila/load-balancer/configuration"
)

type Server struct {
	Url       string
	isHealthy bool
}

type HealthCheck struct {
	Enabled    bool
	Path       string
	IntervalMs int
}

type LoadBalancer struct {
	allServers  []*Server
	healthCheck HealthCheck
	strategy    Strategy
}

func NewLoadBalancer(conf configuration.Configuration) LoadBalancer {
	var servers []*Server
	for i := 0; i < len(conf.Servers); i++ {
		server := Server{Url: conf.Servers[i].Url, isHealthy: true}
		servers = append(servers, &server)
	}

	loadBalancer := LoadBalancer{
		// NextServer: make(chan *Server),
		allServers: servers,
		healthCheck: HealthCheck{
			Enabled:    conf.HealthCheck.Enabled,
			Path:       conf.HealthCheck.Path,
			IntervalMs: conf.HealthCheck.IntervalMs,
		},
	}

	switch conf.Strategy {
	case "random":
		loadBalancer.strategy = CreateRandom(servers)
	case "round-robin":
		loadBalancer.strategy = CreateRoundRobin(servers)
	default:
		panic("Unknown load balancing strategy")
	}

	if loadBalancer.healthCheck.Enabled {
		for _, server := range loadBalancer.allServers {
			fmt.Println("Kick up health check for ", server.Url)
			go loadBalancer.doHealthCheck(server)
		}
	}

	return loadBalancer
}

func (b *LoadBalancer) doHealthCheck(server *Server) {
	for {
		client := http.DefaultClient
		response, err := client.Get(server.Url + b.healthCheck.Path)
		if err != nil {
			fmt.Println("Health check failed for", server.Url)
			server.isHealthy = false
		} else {
			if response.StatusCode == http.StatusOK {
				fmt.Println("Health check OK for ", server.Url)
				server.isHealthy = true
			} else {
				server.isHealthy = false
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
		// TODO: close the request
		return
	}
	fmt.Println("server: ", server)

	req, _ := http.NewRequest(r.Method, server.Url, nil)
	req.Header = r.Header
	req.Header.Add("X-Forwarded-For", r.RemoteAddr)
	// Add rest of the custom headers
	req.Body = r.Body

	// TODO: handle connection refused error gracefully
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error with request")
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Got response back \nres: %s\n", body)
	w.Write(body)
}

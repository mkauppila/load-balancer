package loadBalancer

import (
	"container/ring"
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

type LoadBalancer struct {
	allServers []*Server
	// The healthyServers is access from 2 goroutines atm. Unsafe?
	healthyServers *ring.Ring // rename to healthyServers
	// Actually this should rather be a RWMutex
	NextServer chan *Server
}

func NewLoadBalancer(conf configuration.Configuration) LoadBalancer {
	fmt.Println("conf: ", conf)

	var servers []*Server
	r := ring.New(len(conf.Servers))
	for i := 0; i < len(conf.Servers); i++ {
		server := Server{Url: conf.Servers[i].Url, isHealthy: true}
		r.Value = &server
		servers = append(servers, &server)
		r = r.Next()
	}
	fmt.Println(servers)

	loadBalancer := LoadBalancer{healthyServers: r, NextServer: make(chan *Server), allServers: servers}

	// wont this goroutine dangle if the context is deleted?
	go loadBalancer.nextServerStream()
	for _, server := range loadBalancer.allServers {
		fmt.Println("Kick up health check for ", server.Url)
		go loadBalancer.doHealthCheck(server)
	}
	return loadBalancer
}

func (b *LoadBalancer) nextServerStream() {
	for {
		// will fail with zero servers, ie servers == []
		b.NextServer <- b.getNextServer()
	}
}

func (b *LoadBalancer) getNextServer() *Server {
	b.healthyServers = b.healthyServers.Move(1)
	return b.healthyServers.Value.(*Server)
}

func (b *LoadBalancer) doHealthCheck(server *Server) {
	removeUnhealthy := func(server *Server) {
		ring := b.healthyServers
		for len := 0; len < b.healthyServers.Len(); len++ {
			if ring.Value.(*Server).Url == server.Url {
				ring.Unlink(1)

				break
			}

			ring = ring.Move(1)
		}
	}

	for {
		// Run the check once in 2 seconds
		time.Sleep(time.Second * 2)

		client := http.DefaultClient
		response, err := client.Get(server.Url + "/health")
		if err != nil {
			fmt.Println("Health check failed for", server.Url)
			server.isHealthy = false
			removeUnhealthy(server)
		} else {
			if response.StatusCode == http.StatusOK {
				fmt.Println("Health check OK for ", server.Url)
				server.isHealthy = true
			} else {
				server.isHealthy = false
				fmt.Println("Health check failed for", server.Url,
					" wrong response status ", response.StatusCode)
				removeUnhealthy(server)
			}
		}
	}
}

func (b *LoadBalancer) ForwardRequest(w http.ResponseWriter, r *http.Request) {
	server := <-b.NextServer
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

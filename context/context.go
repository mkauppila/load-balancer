package context

import (
	"container/ring"
	"fmt"
	"net/http"
	"time"
)

type Server struct {
	Url       string
	isHealthy bool
}

type Context struct {
	allServers []*Server
	// The healthyServers is access from 2 goroutines atm. Unsafe?
	healthyServers *ring.Ring // rename to healthyServers
	NextServer     chan *Server
}

func (c *Context) nextServerStream() {
	for {
		c.NextServer <- c.getNextServer()
	}
}

func (c *Context) Close() {
	// TODO: this make the channel panic!
	close(c.NextServer)
}

func (c *Context) getNextServer() *Server {
	c.healthyServers = c.healthyServers.Move(1)
	return c.healthyServers.Value.(*Server)
}

func (c *Context) doHealthCheck(server *Server) {
	removeUnhealthy := func(server *Server) {
		ring := c.healthyServers
		for len := 0; len < c.healthyServers.Len(); len++ {
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

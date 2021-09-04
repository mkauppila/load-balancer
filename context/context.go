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
	servers    *ring.Ring
	NextServer chan Server
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

func (c *Context) getNextServer() Server {
	c.servers = c.servers.Move(1)
	return c.servers.Value.(Server)
}

func (c *Context) doHealthCheck(server Server) {
	fmt.Println("run a health check", server)
	for {
		// Run the check once in 2 seconds
		time.Sleep(time.Second * 2)

		client := http.DefaultClient
		response, err := client.Get(server.Url + "/health")
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
	}
}

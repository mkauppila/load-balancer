package context

import (
	"container/ring"
)

type Server struct {
	Url string
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

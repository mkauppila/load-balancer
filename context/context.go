package context

import "container/ring"

type Server struct {
	Url string
}

type Context struct {
	servers *ring.Ring
	// now I could hide a mutex here so move this to another file?
}

func (c *Context) GetNextServer() Server {
	servers := c.servers.Move(1)
	c.servers = servers // this is no obsolete?
	return servers.Value.(Server)
}

package configuration

import (
	"container/ring"
	"errors"
	"io/ioutil"
	"strings"
)

type Server struct {
	Url string
}

type Configuration struct {
	servers *ring.Ring
	// now I could hide a mutex here so move this to another file?
}

func (c *Configuration) GetNextServer() Server {
	servers := c.servers.Move(1)
	c.servers = servers // this is no obsolete?
	return servers.Value.(Server)
}

// needs to update some values for that server (like connection count)

func ParseConfiguration() (Configuration, error) {
	data, err := ioutil.ReadFile("../lb.conf")
	if err != nil {
		return Configuration{}, errors.New("No configuration file exists.")
	}

	var servers []Server
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			break
		}
		d := strings.Split(line, " ")
		url := d[1]
		server := Server{Url: strings.Trim(url, " \n\t")}

		servers = append(servers, server)
	}

	r := ring.New(len(servers))
	for i := 0; i < len(servers); i++ {
		r.Value = servers[i]
		r = r.Next()
	}

	return Configuration{servers: r}, nil
}

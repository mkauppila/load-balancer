package context

import (
	"container/ring"
	"errors"
	"io/ioutil"
	"strings"
)

func ParseConfiguration() (Context, error) {
	data, err := ioutil.ReadFile("lb.conf")
	if err != nil {
		return Context{}, errors.New("no configuration file exists")
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

	context := Context{servers: r, NextServer: make(chan Server)}
	// wont this goroutine dangle if the context is deleted?
	go context.nextServerStream()
	return context, nil
}

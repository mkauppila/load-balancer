package configuration

import (
	"errors"
	"strconv"
	"strings"
)

// parse configuration separately and
// and create LoadBalancer out of it
// Rename Context to LoadBalancer
type HealthCheck struct {
	Enabled    bool
	IntervalMs int
	Path       string
}

type Server struct {
	Url string
}

type Strategy string

const (
	random     Strategy = "random"
	roundRobin Strategy = "round-robin"
)

type Configuration struct {
	Servers []Server
	HealthCheck
	Strategy
}

func ParseConfiguration(contents []byte) (conf Configuration, err error) {
	for _, line := range strings.Split(string(contents), "\n") {
		if len(line) == 0 {
			continue
		}
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "health_check") {
			healthCheck, err := parseHealthCheck(line)
			if err != nil {
				return conf, err
			}
			conf.HealthCheck = healthCheck
		} else if strings.HasPrefix(line, "server") {
			server := parseServer(line)
			conf.Servers = append(conf.Servers, server)
		} else if strings.HasPrefix(line, "strategy") {
			conf.Strategy = parseStrategy(line)
		}
		// else {
		// 	// unknown, skip or fail?
		// }
	}
	return conf, nil
}

func parseHealthCheck(line string) (hc HealthCheck, err error) {
	splittedLine := strings.Split(line, " ")
	for index, item := range splittedLine {
		switch index {
		case 1:
			if item == "on" {
				hc.Enabled = true
			} else if item == "off" {
				hc.Enabled = false
			} else {
				return hc, errors.New("Unknown option: " + item)
			}
		case 2:
			interval, err := strconv.ParseInt(item, 10, 32)
			if err != nil {
				return hc, errors.New("interval needs to be an integer")
			}
			hc.IntervalMs = int(interval)
		case 3:
			hc.Path = item
		}
	}
	return hc, nil
}

func parseServer(line string) Server {
	splittedLine := strings.Split(line, " ")
	return Server{Url: splittedLine[1]}
}

func parseStrategy(line string) Strategy {
	splittedLine := strings.Split(line, " ")
	switch Strategy(splittedLine[1]) {
	case random:
		return random
	case roundRobin:
		return roundRobin
	default:
		// TODO handle this properly and return a parse error
		return "error"
	}
}

package config

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mkauppila/load-balancer/types"
)

type Configuration struct {
	Servers []types.Server
	types.HealthCheck
	types.Strategy
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

func parseHealthCheck(line string) (hc types.HealthCheck, err error) {
	parts := strings.Split(line, " ")
	for index, item := range parts {
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

func parseServer(line string) types.Server {
	parts := strings.Split(line, " ")
	return types.Server{Url: parts[1]}
}

func parseStrategy(line string) types.Strategy {
	parts := strings.Split(line, " ")
	switch types.Strategy(parts[1]) {
	case types.Random:
		return types.Random
	case types.RoundRobin:
		return types.RoundRobin
	default:
		// TODO handle this properly and return a parse error
		return "error"
	}
}

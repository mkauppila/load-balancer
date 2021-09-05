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

type Configuration struct {
	HealthCheck HealthCheck
	Servers     []Server
}

		if len(line) == 0 {
			break
		}

		if strings.HasPrefix(line, "health_check") {
			healthCheck := parseHealthCheck(line)
			conf.HealthCheck = healthCheck
		} else if strings.HasPrefix(line, "server") {
			server := parseServer(line)
			conf.Servers = append(conf.Servers, server)
		} else {
			// unknown, skip or fail?
		}
	}

	return conf, nil
}

func parseHealthCheck(line string) (hc HealthCheck) {
	splittedLine := strings.Split(line, " ")
	for index, item := range splittedLine {
		switch index {
		case 1:
			if item == "on" {
				hc.Enabled = true
			} else if item == "off" {
				hc.Enabled = false
			} else {
				/// do erorr or warning etc!
			}
		case 2: // interval
			interval, err := strconv.ParseInt(item, 10, 32)
			if err != nil {
				// do some error1
			}
			hc.IntervalMs = int(interval)
		case 3: // path
			hc.Path = item // should check for / or something?
		default:
			// fail!
		}
	}
	return
}

func parseServer(line string) (server Server) {
	splittedLine := strings.Split(line, " ")
	server.Url = splittedLine[1]
	return
}

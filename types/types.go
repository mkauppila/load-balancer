package types

type HealthCheck struct {
	Enabled    bool
	IntervalMs int
	Path       string
}

type Server struct {
	Url       string
	IsHealthy bool
}

type Strategy string

const (
	Random     Strategy = "random"
	RoundRobin Strategy = "round-robin"
)

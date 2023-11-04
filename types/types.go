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
	StrategyRandom     Strategy = "random"
	StrategyRoundRobin Strategy = "round-robin"
)

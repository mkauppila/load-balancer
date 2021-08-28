module load-balancer/lb

go 1.17

replace internal/configuration => ./internal/configuration

require internal/configuration v0.0.0-00010101000000-000000000000 // indirect

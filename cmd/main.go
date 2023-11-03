package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target
// - add different load balancing methods, weighted, least connections, ip_hash, some other hash etc..

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/lb"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	contents, err := os.ReadFile("lb.conf")
	if err != nil {
		panic("no config file exists")
	}

	conf, err := config.ParseConfiguration(contents)
	if err != nil {
		log.Fatalln("Failed to parse config. Error: ", err)
	}

	srv := lb.NewLoadBalancer(conf)

	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
		syscall.SIGTERM,
	)

	_ = srv.Start(ctx)

	<-ctx.Done()
}

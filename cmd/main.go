package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target
// - add different load balancing methods, weighted, least connections, ip_hash, some other hash etc..

import (
	"context"
	"log"
	"net/http"
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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		srv.ForwardRequest(w, r)
	})

	ctx := context.Background()
	ctx, stopFn := signal.NotifyContext(
		ctx,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
		syscall.SIGTERM,
	)

	go func() {
		err := http.ListenAndServe("localhost:4000", nil)
		if err != nil {
			log.Println(err)
			stopFn()
			return
		}
	}()

	<-ctx.Done()
}

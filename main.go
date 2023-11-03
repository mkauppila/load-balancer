package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target
// - add different load balancing methods, weighted, least connections, ip_hash, some other hash etc..

import (
	"log"
	"net/http"
	"os"

	"github.com/mkauppila/load-balancer/config"
	"github.com/mkauppila/load-balancer/lb"
)

func main() {
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

	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}

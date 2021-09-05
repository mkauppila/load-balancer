package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target
// - add config options health checks to enable and set the interval
// - add different load balancing methods, weighted, least connections, ip_hash, some other hash etc..

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mkauppila/load-balancer/configuration"
	"github.com/mkauppila/load-balancer/loadBalancer"
)

func main() {
	contents, err := ioutil.ReadFile("lb.conf")
	if err != nil {
		panic("no configuration file exists")
	}

	conf, err := configuration.ParseConfiguration(contents)
	if err != nil {
		log.Fatalln("Failed to parse configuration. Error: ", err)
	}

	loadBalancer := loadBalancer.NewLoadBalancer(conf)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		loadBalancer.ForwardRequest(w, r)
	})
	// AP: what is DefaultServeMux?
	log.Fatal(http.ListenAndServe("localhost:4000", nil))
}

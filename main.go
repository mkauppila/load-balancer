package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

// load balance a get request in round robin
//  -- add support for round robin (which one is the next?, counter? a circular list?) /OK

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target
// - add health checks if enabled
// - add different load balancing methods, weighted, least connections, ip_hash, some other hash etc..
//

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/mkauppila/load-balancer/context"
)

func forwardRequest(context *context.Context, w http.ResponseWriter, r *http.Request) {
	server := <-context.NextServer

	req, _ := http.NewRequest(r.Method, server.Url, nil)
	req.Header = r.Header
	req.Header.Add("X-Forwarded-For", r.RemoteAddr)
	// Add rest of the custom headers
	req.Body = r.Body

	// TODO: handle connection refused error gracefully
	client := http.DefaultClient
	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error with request")
	}
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Got response back \nres: %s\n", body)
	w.Write(body)
}

func start(context *context.Context) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { forwardRequest(context, w, r) })
	// AP: what is DefaultServeMux?
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	context, err := context.ParseConfiguration()
	if err != nil {
		log.Fatalln("Failed to parse configuration with error: ", err)
	}
	defer context.Close()

	fmt.Println("Starting up!")
	start(&context)
}

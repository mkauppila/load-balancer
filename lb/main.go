package main

// hello world golang http server, returning the process id

// load balance a get request in round robin
//  -- add support for round robin (which one is the next?, counter? a circular list?)
// take in a request, read the HTTP verb (GET, PUT, ...)
//   -- test that the POST, PUT requests work
// create new request to the actual server with existing headers and body /OK
//  - verify that the body works!
// (optional) add a load balancer header to it
// take the response from actual server and relay it back to the client /OK

// TODO maybe?
// - this could support HTTPS and break it before sending the request as plain HTTP to the target

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const servers = "http://localhost:4001"

func forwardRequest(w http.ResponseWriter, r *http.Request) {
	client := http.DefaultClient
	r.Host = servers

	req, _ := http.NewRequest(r.Method, servers, nil)
	req.Header = r.Header
	req.Header.Add("X-Forwarded-For", r.RemoteAddr)
	req.Body = r.Body

	response, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		log.Fatal("Error with request")
	}
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Got response back \nres: %s\n", body)
	w.Write(body)

	defer response.Body.Close()
}

func start() {
	http.HandleFunc("/", forwardRequest)
	// AP: what is DefaultServeMux?
	log.Fatal(http.ListenAndServe(":4000", nil))
}

func main() {
	fmt.Println("Starting up!")
	start()
}

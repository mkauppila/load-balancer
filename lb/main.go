package main

// read: https://docs.aws.amazon.com/elasticloadbalancing/latest/application/x-forwarded-headers.html
// read: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-load-balancer/

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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	defer response.Body.Close()

	body, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("Got response back \nres: %s\n", body)
	w.Write(body)
}

func start() {
	http.HandleFunc("/", forwardRequest)
	// AP: what is DefaultServeMux?
	log.Fatal(http.ListenAndServe(":4000", nil))
}

type Server struct {
	url string
}
type Configuration struct {
	servers []Server
}

func ParseConfiguration() (Configuration, error) {
	data, err := ioutil.ReadFile("../lb.conf")
	if err != nil {
		return Configuration{}, errors.New("No configuration file exists.")
	}

	var servers []Server
	for _, line := range strings.Split(string(data), "\n") {
		if len(line) == 0 {
			break
		}
		d := strings.Split(line, " ")
		url := d[1]
		server := Server{url: url}

		servers = append(servers, server)
	}

	return Configuration{servers: servers}, nil
}

var configuration Configuration

func main() {
	conf, err := ParseConfiguration()
	if err != nil {
		log.Fatalln("Failed to parse configuration...", err)
	}
	fmt.Println(conf)
	configuration = conf

	fmt.Println("Starting up!")
	start()
}

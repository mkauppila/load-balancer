package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func logRequestDetails(r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.Host)
	for name, values := range r.Header {
		fmt.Println(name, values)
	}
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Println("body: ", string(body))
}

func page(w http.ResponseWriter, r *http.Request) {
	logRequestDetails(r)

	fmt.Fprintf(w, fmt.Sprintf("Hello from %d", os.Getpid()))
}

func start() {
	http.HandleFunc("/", page)
	port := os.Getenv("PORT")
	fmt.Printf("Running on port: %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func main() {
	start()
}

package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func logRequestDetails(r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.Host)
	for name, values := range r.Header {
		fmt.Println(name, values)
	}

	body, _ := io.ReadAll(r.Body)
	fmt.Println("body: ", string(body))
}

func page(w http.ResponseWriter, r *http.Request) {
	logRequestDetails(r)

	fmt.Fprintf(w, "Hello from %d", os.Getpid())
}

func main() {
	http.HandleFunc("/", page)
	port := os.Getenv("PORT")
	if port == "" {
		panic("no PORT env defined")
	}
	fmt.Printf("Running on port: %s\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%s", port), nil))
	fmt.Println("Shutting down the test API")
}

package httpserver

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

func page(w http.ResponseWriter, r *http.Request) {
	logRequestDetails(r)
	fmt.Fprintf(w, "Hello from %d", os.Getpid())
}

func RunServer(ctx context.Context, confUrl string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", page)

	addr, _ := url.Parse(confUrl)
	fmt.Printf("Starting on addr: %s:%s\n", addr.Hostname(), addr.Port())
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port()), mux)
		if err != nil {
			cancel()
			return
		}
	}()
}

func logRequestDetails(r *http.Request) {
	fmt.Printf("%s %s\n", r.Method, r.Host)
	for name, values := range r.Header {
		fmt.Println(name, values)
	}

	body, _ := io.ReadAll(r.Body)
	fmt.Println("body: ", string(body))
}

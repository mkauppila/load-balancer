package httpserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"
)

func page(w http.ResponseWriter, r *http.Request) {
	logRequestDetails(r)
	fmt.Fprintf(w, "Hello from %d", os.Getpid())
}

func RunServer(started context.CancelFunc, serverCtx context.Context, confUrl string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", page)

	addr, _ := url.Parse(confUrl)
	fmt.Printf("Starting on addr: %s:%s\n", addr.Hostname(), addr.Port())
	serverCtx, cancel := context.WithCancel(serverCtx)
	go func() {
		l, err := net.Listen("tcp", fmt.Sprintf("%s:%s", addr.Hostname(), addr.Port()))
		if err != nil {
			panic(err)
		}

		started()

		go func() {
			err = http.Serve(l, mux)
			if err != nil {
				cancel()
			}
		}()
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

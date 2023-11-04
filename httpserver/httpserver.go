package httpserver

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
)

func page(w http.ResponseWriter, r *http.Request, response string) {
	//logRequestDetails(r)
	_, _ = fmt.Fprintf(w, response)
}

func RunServer(started context.CancelFunc, ctx context.Context, confUrl string, response string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		page(writer, request, response)
	})

	_, cancel := context.WithCancel(ctx)
	go func() {
		l, err := net.Listen("tcp", confUrl)
		if err != nil {
			panic(err)
		}

		started()

		err = http.Serve(l, mux)
		if err != nil {
			cancel()
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

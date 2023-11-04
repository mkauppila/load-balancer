package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/mkauppila/load-balancer/httpserver"
)

func main() {
	run(os.Args)
}

func run(args []string) {
	cancelFn := func() {}

	ctx, _ := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGABRT,
		syscall.SIGTERM,
	)

	port := os.Getenv("HTTP_PORT")
	addr := fmt.Sprintf("http://localhost:%s", port)

	httpserver.RunServer(cancelFn, ctx, addr, "response")

	<-ctx.Done()
}

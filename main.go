package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"
)

func run() error {
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		return err
	}

	log.Printf("Listening on http://%v", listener.Addr())

	socketServer := newSocketServer()

	server := &http.Server{
		Handler:      socketServer,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}

	err = server.Serve(listener)
	if err != nil {
		log.Printf("Server Initialization Error: %v\n", err)
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	return server.Shutdown(ctx)
}

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

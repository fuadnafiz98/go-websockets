package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server running on 0.0.0.0:8888")

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO"))
	})

	server := &http.Server{
		Addr:    "0.0.0.0:8888",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}

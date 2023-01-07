package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Data struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

func main() {
	fmt.Println("Server running on 0.0.0.0:8888")

	handler := http.NewServeMux()

	handler.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id")
		w.Write([]byte("id is => " + id))
	})

	handler.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		data := &Data{
			Id:       "98",
			Username: "fuadnafiz98",
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err)
		}
	})

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

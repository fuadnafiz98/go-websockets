package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"google.golang.org/protobuf/proto"

	book "github.com/fuadnafiz98/go-websockets/book"
)

type Data struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

var DB_PATH string = "database.pb"

func _main() {
	fmt.Println("Server running on 0.0.0.0:8888")

	handler := http.NewServeMux()

	handler.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			book := &book.Book{
				Id:     1,
				Title:  "Make Time",
				Author: "Un authored",
			}
			fmt.Println(proto.Marshal(book))

		case http.MethodPost:
		case http.MethodPatch:
		case http.MethodDelete:
			id := r.URL.Query().Get("id")
			w.Write([]byte("id is => " + id))
		default:
			http.Error(w, "Method not allowd", http.StatusMethodNotAllowed)
		}
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

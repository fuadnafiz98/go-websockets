# Golang WebServer

It's a repo supposed to be about websockets. But before starting I have to learn how the regular web server works on golang :/

Let's create a good old web server in golang

Project setup:

```bash
go mod init github.com/fuadnafiz98/go-websockets
```

create a file `main.go`

```bash
touch main.go
# or
# > main.go
```

Start from here:

```golang
package main

func main() {
	//
}
```

Let's start from basic.

`http.ListenAndServe(addr, handler)`

this takes an address and a handler. Handler is like the router.
for now keep the handler to `nil` and give it an address. This function starts the server or returns an error.

We have to check the error and thats it!
we have a running server!

```golang
package main

import (
	"fmt"
	"net/http"
)

func main() {
  fmt.Println("Server running on 0.0.0.0:8888")
	err := http.ListenAndServe("0.0.0.0:8888", nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

Now run the file `go run main.go`

Now curl the server `curl localhost:8888`

which will show `404 page not found`

which is good as we didn't setup any handler!

Now we have to write a handler function.

But before that we can declare the http server in a bit different way which will help us a lot in future.

```golang
package main

import (
	"fmt"
	"net/http"
)

func main() {
	fmt.Println("Server running on 0.0.0.0:8888")

	server := &http.Server{
		Addr: "0.0.0.0:8000",
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

This is the same thing as before, we are declaring a `server` instance which is a pointer to a http server. And we are modifiying the Server struct of the http. (I am not sure how to explain this, will modify this later).

Now we can put the address, along with some other configuration of the server.

We can also see there is a function in the `http.Server` struct called `Handler` which is `http.DefaultServerMux` by default but we can implement our own ServerMux.

For implementing that we have to create a new ServerMux `http.NewServeMux()`

```golang
func main() {
	fmt.Println("Server running on 0.0.0.0:8888")

	handler := http.NewServeMux()

	handler.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("HELLO"))
	})

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: handler,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
		return
	}
}
```

Now in the `HandleFunc` we will declare the pattern `/` and write a `(req, res)` function like in node.js. here it just writes the byte 'HELLO'

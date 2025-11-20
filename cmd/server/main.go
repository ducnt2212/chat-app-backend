package main

import (
	"flag"
	"fmt"
	"net/http"
)

func main() {

	addr := flag.String("addr", "", "HTTP network address")

	flag.Parse()

	app := NewApplication(*addr)

	server := http.Server{
		Addr:    app.Addr,
		Handler: app.routes(),
	}

	fmt.Printf("Server is starting on %s\n", app.Addr)

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

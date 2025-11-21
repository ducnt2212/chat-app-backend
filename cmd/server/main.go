//go:debug x509negativeserial=1
package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/ducnt2212/chat-app-backend/internal/repository"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	addr := flag.String("addr", "", "HTTP network address")
	flag.Parse()

	connString := fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_SERVER"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))
	repo, err := repository.NewSQLServerDB(connString)
	if err != nil {
		panic(err)
	}

	app, err := NewApplication(*addr, repo)
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Addr:    app.Addr,
		Handler: app.routes(),
	}

	fmt.Printf("Server is starting on %s\n", app.Addr)

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}

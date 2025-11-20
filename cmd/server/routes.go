package main

import "net/http"

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", app.health)

	return mux

}

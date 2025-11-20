package main

import "net/http"

func (app *Application) health(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("OK"))
}

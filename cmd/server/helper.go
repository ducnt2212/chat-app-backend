package main

import (
	"encoding/json"
	"net/http"
)

func (app *Application) replyJson(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(payload)
}

func (app *Application) replyError(writer http.ResponseWriter, status int, payload string) {
	http.Error(writer, payload, status)
}

package helper

import (
	"encoding/json"
	"net/http"
)

func ReplyJSON(writer http.ResponseWriter, status int, payload any) {
	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(payload)
}

func ReplyJSONError(writer http.ResponseWriter, status int, errorMsg string) {
	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(map[string]string{"error": errorMsg})
}

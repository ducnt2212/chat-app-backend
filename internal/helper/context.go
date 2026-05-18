package helper

import "net/http"

type ContextKey string

const (
	UserIDContextKey ContextKey = "user_id"
)

func GetUserID(request *http.Request) int {
	userID, _ := request.Context().Value(UserIDContextKey).(int)
	return userID
}

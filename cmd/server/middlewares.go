package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/ducnt2212/chat-app-backend/internal/helper"
)

func (app *Application) requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

		msg := fmt.Sprintf("%s - %s %s %s", request.RemoteAddr, request.Proto, request.Method, request.URL.RequestURI())
		app.logger.Info(msg)
		next.ServeHTTP(writer, request)
	})
}

func (app *Application) panicRecover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("Panic recovered: %s - Source: %s", err, debug.Stack())
				app.logger.Warning(msg)
				http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(writer, request)
	})
}

func (app *Application) authChecker(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			helper.ReplyJSONError(writer, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			helper.ReplyJSONError(writer, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenStr := parts[1]
		claims, err := helper.VerifyJwt(tokenStr, os.Getenv("JWT_SECRET"))
		if err != nil {
			switch err {
			case helper.ErrInvalidToken:
				helper.ReplyJSONError(writer, http.StatusUnauthorized, "Invalid token")
			case helper.ErrInvalidTokenClaims:
				helper.ReplyJSONError(writer, http.StatusUnauthorized, "Invalid Token claims")
			default:
				app.logger.Error(err.Error())
				helper.ReplyJSONError(writer, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		userID := int(claims["user_id"].(float64))
		ctx := context.WithValue(request.Context(), helper.UserIDContextKey, userID)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

package main

import (
	"context"
	"fmt"
	"net/http"
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
				msg := fmt.Sprintf("panic recovered: %s - Source: %s", err, debug.Stack())
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
			app.replyJsonError(writer, http.StatusUnauthorized, "Missing Authorization header")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			app.replyJsonError(writer, http.StatusUnauthorized, "Invalid Authorization header format")
			return
		}

		tokenStr := parts[1]
		claims, err := helper.VerifyJwt(tokenStr)
		if err != nil {
			switch err {
			case helper.ErrInvalidToken:
				app.replyJsonError(writer, http.StatusUnauthorized, "Invalid token")
			case helper.ErrInvalidTokenClaims:
				app.replyJsonError(writer, http.StatusUnauthorized, "Invalid Token claims")
			default:
				app.logger.Error(err.Error())
				app.replyJsonError(writer, http.StatusInternalServerError, "Internal Server Error")
			}
			return
		}

		userID := int(claims["user_id"].(float64))
		ctx := context.WithValue(request.Context(), "user_id", userID)

		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

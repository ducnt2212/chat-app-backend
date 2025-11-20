package main

import (
	"net/http"

	"github.com/ducnt2212/chat-app-backend/internal/middleware"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	middlewares := middleware.NewChain(app.RecoverPanic, app.RequestLogger)

	mux.Handle("GET /", middlewares.ThenFunc(app.health))
	mux.Handle("POST /auth/register", middlewares.ThenFunc(app.register))
	mux.Handle("POST /auth/login", middlewares.ThenFunc(app.login))

	return mux
}

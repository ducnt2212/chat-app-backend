package main

import (
	"net/http"

	"github.com/ducnt2212/chat-app-backend/internal/middleware"
)

func (app *Application) routes() http.Handler {
	mux := http.NewServeMux()

	middlewares := middleware.NewChain(app.panicRecover, app.requestLogger)

	mux.Handle("/", middlewares.ThenFunc(app.health))
	mux.Handle("POST /auth/register", middlewares.ThenFunc(app.register))
	mux.Handle("POST /auth/login", middlewares.ThenFunc(app.login))

	authMiddlewares := middlewares.Append(app.authChecker)

	mux.Handle("GET /rooms", authMiddlewares.ThenFunc(app.health))

	return mux
}

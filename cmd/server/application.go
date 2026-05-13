package main

import (
	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/repository"
	"github.com/ducnt2212/chat-app-backend/internal/websocket"
)

type Application struct {
	Addr      string
	logger    *logger.Logger
	repo      repository.IRepository
	wsHandler *websocket.Handler
}

func NewApplication(Addr string, repo repository.IRepository, hub *websocket.Hub) (*Application, error) {
	go hub.Run()
	lg := logger.GetLogger()

	return &Application{
		Addr:      Addr,
		logger:    lg,
		repo:      repo,
		wsHandler: websocket.NewHandler(hub, lg, repo),
	}, nil
}

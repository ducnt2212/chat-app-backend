package main

import (
	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/repository"
	"github.com/ducnt2212/chat-app-backend/internal/service"
	"github.com/ducnt2212/chat-app-backend/internal/websocket"
)

type Application struct {
	Addr           string
	logger         *logger.Logger
	wsHandler      *websocket.Handler
	userService    *service.UserService
	roomService    *service.RoomService
	messageService *service.MessageService
}

func NewApplication(Addr string, repo repository.IRepository, hub *websocket.Hub) (*Application, error) {
	go hub.Run()
	lg := logger.GetLogger()

	userService := service.NewUserService(repo)
	roomService := service.NewRoomService(repo, repo)
	msgService := service.NewMessageService(repo, repo)

	return &Application{
		Addr:           Addr,
		logger:         lg,
		wsHandler:      websocket.NewHandler(hub, lg, roomService, msgService, userService),
		userService:    userService,
		roomService:    roomService,
		messageService: msgService,
	}, nil
}

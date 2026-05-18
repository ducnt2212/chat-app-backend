package websocket

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ducnt2212/chat-app-backend/internal/helper"
	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/service"
	gorilla "github.com/gorilla/websocket"
)

type Handler struct {
	Hub            *Hub
	logger         *logger.Logger
	roomService    *service.RoomService
	messageService *service.MessageService
	userService    *service.UserService
}

func NewHandler(hub *Hub, logger *logger.Logger, roomService *service.RoomService, messageService *service.MessageService, userService *service.UserService) *Handler {
	return &Handler{
		Hub:            hub,
		logger:         logger,
		roomService:    roomService,
		messageService: messageService,
		userService:    userService,
	}
}

var upgrader = gorilla.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (handler *Handler) ServeWS(writer http.ResponseWriter, request *http.Request, userID int) {
	roomID, err := getRoomID(request.URL.Path)
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Room id")
		return
	}

	canAccess, err := handler.roomService.CanAccessRoom(roomID, userID)
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Internal Server Error")
		handler.logger.Error(err.Error())
		return
	}

	if !canAccess {
		helper.ReplyJSONError(writer, http.StatusForbidden, "Forbidden")
		return
	}

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	client := &Client{
		UserID:         userID,
		RoomID:         roomID,
		SendChan:       make(chan Event, 256),
		Conn:           conn,
		Hub:            handler.Hub,
		logger:         handler.logger,
		messageService: handler.messageService,
		userService:    handler.userService,
	}

	handler.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func getRoomID(path string) (int, error) {
	roomIDText := strings.TrimPrefix(path, "/ws/rooms/")
	return strconv.Atoi(roomIDText)
}

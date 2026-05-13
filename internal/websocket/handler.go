package websocket

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ducnt2212/chat-app-backend/internal/logger"
	gorilla "github.com/gorilla/websocket"
)

type Handler struct {
	Hub            *Hub
	logger         *logger.Logger
	messageService MessageService
}

func NewHandler(hub *Hub, logger *logger.Logger, messageService MessageService) *Handler {
	return &Handler{
		Hub:            hub,
		logger:         logger,
		messageService: messageService,
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
		http.Error(writer, "invalid room id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(writer, request, nil)
	if err != nil {
		return
	}

	client := &Client{
		UserID:         userID,
		RoomID:         roomID,
		Send:           make(chan Event, 256),
		Conn:           conn,
		Hub:            handler.Hub,
		logger:         handler.logger,
		messageService: handler.messageService,
	}

	handler.Hub.Register(client)

	go client.WritePump()
	go client.ReadPump()
}

func getRoomID(path string) (int, error) {
	roomIDText := strings.TrimPrefix(path, "/ws/rooms/")
	return strconv.Atoi(roomIDText)
}

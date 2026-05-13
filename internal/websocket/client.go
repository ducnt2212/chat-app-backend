package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/models"
	gorilla "github.com/gorilla/websocket"
)

type Client struct {
	UserID         int
	RoomID         int
	Send           chan Event
	Conn           *gorilla.Conn
	Hub            *Hub
	logger         *logger.Logger
	messageService MessageService
}

type MessageService interface {
	CreateMessage(message models.Message) (models.Message, error)
}

func (client *Client) ReadPump() {
	defer func() {
		client.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		var event Event
		err := client.Conn.ReadJSON(&event)
		if err != nil {
			client.logger.Error(fmt.Sprintf("websocket read error: %v", err))
			break
		}

		switch event.Type {
		case EventSendMessage:
			client.handleSendMessage(event)
		case EventTyping:
			event.RoomID = client.RoomID
			client.Hub.Broadcast(event)
		}
	}
}

func (client *Client) WritePump() {
	defer client.Conn.Close()

	for event := range client.Send {
		err := client.Conn.WriteJSON(event)
		if err != nil {
			client.logger.Error(fmt.Sprintf("Websocket write error: %v", err))
			break
		}
	}
}

func (client *Client) handleSendMessage(event Event) {
	data, err := json.Marshal(event.Payload)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Marshal Send message payload error: %v", err))
		return
	}

	var payload SendMessagePayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Unmarshal Send message payload error: %v", err))
		return
	}

	if payload.Content == "" {
		return
	}

	message, err := client.messageService.CreateMessage(models.Message{
		RoomID:   client.RoomID,
		SenderID: client.UserID,
		Content:  payload.Content,
	})
	if err != nil {
		client.logger.Error(fmt.Sprintf("Send message error: %v", err))
		return
	}

	client.Hub.Broadcast(Event{
		Type:   EventNewMessage,
		RoomID: client.RoomID,
		Payload: NewMessagePayload{
			ID:        message.ID,
			RoomID:    message.RoomID,
			SenderID:  message.SenderID,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format(time.RFC3339Nano),
		},
	})
}

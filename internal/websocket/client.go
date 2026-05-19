package websocket

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/service"
	gorilla "github.com/gorilla/websocket"
)

type Client struct {
	UserID         int
	RoomID         int
	SendChan       chan Event
	Conn           *gorilla.Conn
	Hub            *Hub
	logger         *logger.Logger
	messageService *service.MessageService
	userService    *service.UserService
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
			client.logger.Error(fmt.Sprintf("Websocket read error: %v", err))
			break
		}

		switch event.Type {
		case EventSendMessage:
			client.handleSendMessage(event)
		case EventIsTyping:
			client.handleIsTyping(event)
		}
	}
}

func (client *Client) WritePump() {
	defer func() {
		client.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for event := range client.SendChan {
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
		client.logger.Error(fmt.Sprintf("Marshal Send message event error: %v", err))
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

	message, err := client.messageService.SendMessage(client.RoomID, client.UserID, payload.Content)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Create message error: %v", err))
		return
	}

	user, err := client.userService.GetUserByID(message.SenderID)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Get user error: %v", err))
		return
	}

	client.Hub.Broadcast(
		Event{
			Type:   EventNewMessage,
			RoomID: client.RoomID,
			Payload: NewMessagePayload{
				ID:             message.ID,
				RoomID:         message.RoomID,
				SenderID:       message.SenderID,
				SenderUsername: user.Username,
				Content:        message.Content,
				CreatedAt:      message.CreatedAt.Format(time.RFC3339Nano),
			},
		},
	)
}

func (client *Client) handleIsTyping(event Event) {
	data, err := json.Marshal(event.Payload)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Marshal Is typing event error: %v", err))
		return
	}

	var payload IsTypingPayload
	err = json.Unmarshal(data, &payload)
	if err != nil {
		client.logger.Error(fmt.Sprintf("Unmarshal Is typing payload error: %v", err))
		return
	}

	client.Hub.Broadcast(
		Event{
			Type:   EventIsTyping,
			RoomID: client.RoomID,
			Payload: IsTypingPayload{
				UserID:   client.UserID,
				IsTyping: payload.IsTyping,
			},
		},
	)
}

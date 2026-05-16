package websocket

type SendMessagePayload struct {
	Content string `json:"content"`
}

type NewMessagePayload struct {
	ID             int    `json:"id"`
	RoomID         int    `json:"room_id"`
	SenderID       int    `json:"sender_id"`
	SenderUsername string `json:"sender_username"`
	Content        string `json:"content"`
	CreatedAt      string `json:"created_at"`
}

type IsTypingPayload struct {
	UserID   int  `json:"user_id"`
	IsTyping bool `json:"is_typing"`
}

type UserPresencePayload struct {
	UserID int `json:"user_id"`
}

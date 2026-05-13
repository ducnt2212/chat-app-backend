package websocket

type SendMessagePayload struct {
	Content string `json:"content"`
}

type NewMessagePayload struct {
	ID        int    `json:"id"`
	RoomID    int    `json:"room_id"`
	SenderID  int    `json:"sender_id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

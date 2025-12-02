package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	RoomID    int       `json:"room_id"`
	SenderID  int       `json:"sender_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

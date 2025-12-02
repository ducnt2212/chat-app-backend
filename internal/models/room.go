package models

import "time"

type Room struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	IsPrivate bool      `json:"is_private"`
	CreatedBy int       `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}

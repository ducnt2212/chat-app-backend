package repository

import (
	"database/sql"

	"github.com/ducnt2212/chat-app-backend/internal/models"
)

type UserRepository interface {
	CreateUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(user_id int) (models.User, error)
}

type MessageRepository interface {
	CreateMessage(message models.Message) (models.Message, error)
	ListMessagesByRoom(roomID, limit int, cursor string) ([]models.Message, string, error)
}

type RoomRepository interface {
	CreateRoom(room models.Room) (int, error)
	ListRooms() ([]models.Room, error)
}

type RoomMembersRepository interface {
	AddUserToRoom(roomID int, userID int) error
	IsUserInRoom(roomID int, userID int) (bool, error)
}

type IRepository interface {
	UserRepository
	MessageRepository
	RoomRepository
	RoomMembersRepository
}

type Repository struct {
	DB *sql.DB
}

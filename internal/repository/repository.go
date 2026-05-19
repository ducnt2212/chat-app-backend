package repository

import (
	"database/sql"

	"github.com/ducnt2212/chat-app-backend/internal/models"
)

type UserRepository interface {
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(userID int) (models.User, error)
}

type MessageRepository interface {
	CreateMessage(message models.Message) (models.Message, error)
	ListMessagesByRoom(roomID, limit int, cursor string) ([]models.Message, string, error)
}

type RoomRepository interface {
	CreateRoom(room models.Room) (models.Room, error)
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

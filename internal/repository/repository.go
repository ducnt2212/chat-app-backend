package repository

import (
	"database/sql"

	"github.com/ducnt2212/chat-app-backend/internal/models"
)

type IRepository interface {
	CreateUser(user models.User) (int, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(user_id int) (models.User, error)
	CreateMessage(message models.Message) (models.Message, error)
	ListMessagesByRoom(roomID, limit int, cursor string) ([]models.Message, string, error)
	CreateRoom(room models.Room) (int, error)
	ListRooms() ([]models.Room, error)
	AddUserToRoom(roomID int, userID int) error
	IsUserInRoom(roomID int, userID int) (bool, error)
}

type Repository struct {
	DB *sql.DB
}

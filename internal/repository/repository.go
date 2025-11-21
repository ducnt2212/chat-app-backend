package repository

import (
	"database/sql"

	"github.com/ducnt2212/chat-app-backend/internal/models"
)

type IRepository interface {
	CreateUser(username, email, hashedPassword string) (int, error)
	GetUserByUsername(username string) (models.User, error)
	CreateMessage(senderId int, receiverId int, content string) error
	GetMessagesBetweenUsers(userAId, userBId int) ([]models.Message, error)
}

type Repository struct {
	DB *sql.DB
}

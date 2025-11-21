package repository

import (
	"database/sql"

	"github.com/ducnt2212/chat-app-backend/internal/models"
	_ "github.com/microsoft/go-mssqldb"
)

func NewSQLServerDB(connString string) (*Repository, error) {
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Repository{DB: db}, nil
}

func (repo *Repository) CreateUser(username, email, hashedPassword string) (int, error) {
	stmt, err := repo.DB.Prepare(`
	INSERT INTO users (username, email, hashed_password)
	VALUES (@username, @email, @hashed_password);
	SELECT SCOPE_IDENTITY()`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result := stmt.QueryRow(sql.Named("username", username), sql.Named("email", email), sql.Named("hashed_password", hashedPassword))

	var id int
	err = result.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (repo *Repository) GetUserByUsername(username string) (models.User, error) {
	stmt, err := repo.DB.Prepare(`SELECT id, username, email, hashed_password, created_at FROM users WHERE
	username = @username`)
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRow(sql.Named("username", username))

	user := models.User{}
	err = result.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *Repository) CreateMessage(senderId int, receiverId int, content string) error {
	stmt, err := repo.DB.Prepare(`INSERT INTO messages (sender_id, receiver_id, content) VALUES (@sender_id, @receiver_id, @content)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("@sender_id", senderId), sql.Named("@receiver_id", receiverId), sql.Named("@content", content))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) GetMessagesBetweenUsers(userAId, userBId int) ([]models.Message, error) {
	stmt, err := repo.DB.Prepare(`SELECT id, sender_id, receiver_id, content, created_at
	FROM messages
	WHERE (sender_id = @userAId AND receiver_id = @userBId)
	OR (sender_id = @userBId AND receiver_id = @userAId)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query(sql.Named("@userAId", userAId), sql.Named("@userBId", userBId))
	if err != nil {
		return nil, err
	}
	defer result.Close()

	messages := []models.Message{}
	for result.Next() {
		message := models.Message{}
		err := result.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, err
		}

		messages = append(messages, message)
	}

	return messages, nil
}

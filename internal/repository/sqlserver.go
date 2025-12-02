package repository

import (
	"database/sql"
	"time"

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

func (repo *Repository) CreateUser(user models.User) (int, error) {
	stmt, err := repo.DB.Prepare(`
	INSERT INTO users (username, email, hashed_password)
	VALUES (@username, @email, @hashed_password);
	SELECT SCOPE_IDENTITY()`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	result := stmt.QueryRow(sql.Named("username", user.Username), sql.Named("email", user.Email), sql.Named("hashed_password", user.HashedPassword))

	var id int
	err = result.Scan(&id)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (repo *Repository) GetUserByEmail(email string) (models.User, error) {
	stmt, err := repo.DB.Prepare(`SELECT id, username, email, hashed_password
	FROM users
	WHERE email = @email`)
	if err != nil {
		return models.User{}, err
	}
	defer stmt.Close()

	result := stmt.QueryRow(sql.Named("email", email))

	user := models.User{}
	err = result.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repo *Repository) CreateMessage(message models.Message) error {
	stmt, err := repo.DB.Prepare(`INSERT INTO messages (room_id, sender_id, content) VALUES (@room_id, @sender_id, @content)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sql.Named("room_id", message.RoomID), sql.Named("sender_id", message.SenderID), sql.Named("content", message.Content))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ListMessagesByRoom(roomID, limit int, cursor string) ([]models.Message, string, error) {
	var stmt *sql.Stmt
	var err error
	var nextCursor string = ""

	if cursor == "" {
		stmt, err = repo.DB.Prepare(`SELECT TOP (@limit) *
		FROM messages
		WHERE room_id = @room_id
		ORDER BY created_at DESC`)
	} else {
		stmt, err = repo.DB.Prepare(`SELECT TOP (@limit) *
		FROM messages
		WHERE room_id = @room_id AND created_at < @cursor
		ORDER BY created_at DESC`)
	}

	if err != nil {
		return nil, nextCursor, err
	}
	defer stmt.Close()

	result, err := stmt.Query(sql.Named("room_id", roomID), sql.Named("limit", limit), sql.Named("cursor", cursor))
	if err != nil {
		return nil, nextCursor, err
	}
	defer result.Close()

	messages := []models.Message{}
	for result.Next() {
		message := models.Message{}
		err := result.Scan(&message.ID, &message.RoomID, &message.SenderID, &message.Content, &message.CreatedAt)
		if err != nil {
			return nil, nextCursor, err
		}

		messages = append(messages, message)
	}

	if len(messages) > 0 {
		nextCursor = messages[len(messages)-1].CreatedAt.Format(time.RFC3339Nano)
	}

	return messages, nextCursor, nil
}

func (repo *Repository) CreateRoom(room models.Room) error {
	stmt, err := repo.DB.Prepare(`INSERT INTO rooms (name, is_private, created_by) VALUES (@name, @is_private, @created_by)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Query(sql.Named("name", room.Name), sql.Named("is_private", room.IsPrivate), sql.Named("created_by", room.CreatedBy))
	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) ListRooms() ([]models.Room, error) {
	stmt, err := repo.DB.Prepare(`SELECT id, name, is_private, created_by, created_at
	FROM rooms
	WHERE is_private = 0
	ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var rooms []models.Room
	for result.Next() {
		room := models.Room{}
		err = result.Scan(&room.ID, &room.Name, &room.IsPrivate, &room.CreatedBy, &room.CreatedAt)
		if err != nil {
			return nil, err
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

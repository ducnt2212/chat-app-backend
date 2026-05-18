package service

import (
	"github.com/ducnt2212/chat-app-backend/internal/models"
	"github.com/ducnt2212/chat-app-backend/internal/repository"
)

type MessageService struct {
	messageRepo     repository.MessageRepository
	roomMembersRepo repository.RoomMembersRepository
}

func NewMessageService(messageRepo repository.MessageRepository, roomMembersRepo repository.RoomMembersRepository) *MessageService {
	return &MessageService{
		messageRepo:     messageRepo,
		roomMembersRepo: roomMembersRepo,
	}
}

func (s *MessageService) SendMessage(roomID int, senderID int, content string) (models.Message, error) {
	inRoom, err := s.roomMembersRepo.IsUserInRoom(roomID, senderID)
	if err != nil {
		return models.Message{}, err
	}
	if !inRoom {
		return models.Message{}, ErrForbidden
	}

	msg := models.Message{
		RoomID:   roomID,
		SenderID: senderID,
		Content:  content,
	}

	return s.messageRepo.CreateMessage(msg)
}

func (s *MessageService) ListMessages(roomID int, userID int, limit int, cursor string) ([]models.Message, string, error) {
	inRoom, err := s.roomMembersRepo.IsUserInRoom(roomID, userID)
	if err != nil {
		return nil, "", err
	}
	if !inRoom {
		return nil, "", ErrForbidden
	}

	return s.messageRepo.ListMessagesByRoom(roomID, limit, cursor)
}

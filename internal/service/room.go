package service

import (
	"github.com/ducnt2212/chat-app-backend/internal/models"
	"github.com/ducnt2212/chat-app-backend/internal/repository"
)

type RoomService struct {
	roomRepo        repository.RoomRepository
	roomMembersRepo repository.RoomMembersRepository
}

func NewRoomService(roomRepo repository.RoomRepository, roomMembersRepo repository.RoomMembersRepository) *RoomService {
	return &RoomService{
		roomRepo:        roomRepo,
		roomMembersRepo: roomMembersRepo,
	}
}

func (s *RoomService) CreateRoom(room models.Room, creatorID int) (models.Room, error) {
	id, err := s.roomRepo.CreateRoom(room)
	if err != nil {
		return models.Room{}, err
	}

	if err := s.roomMembersRepo.AddUserToRoom(id, creatorID); err != nil {
		return models.Room{}, err
	}

	room.ID = id
	return room, nil
}

func (s *RoomService) JoinRoom(roomID int, userID int) error {
	return s.roomMembersRepo.AddUserToRoom(roomID, userID)
}

func (s *RoomService) CanAccessRoom(roomID int, userID int) (bool, error) {
	return s.roomMembersRepo.IsUserInRoom(roomID, userID)
}

func (s *RoomService) ListRooms() ([]models.Room, error) {
	return s.roomRepo.ListRooms()
}

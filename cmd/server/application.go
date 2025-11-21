package main

import (
	"github.com/ducnt2212/chat-app-backend/internal/logger"
	"github.com/ducnt2212/chat-app-backend/internal/repository"
)

type Application struct {
	Addr   string
	logger *logger.Logger
	repo   repository.IRepository
}

func NewApplication(Addr string, repo repository.IRepository) (*Application, error) {
	return &Application{
		Addr:   Addr,
		logger: logger.GetLogger(),
		repo:   repo,
	}, nil
}

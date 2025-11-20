package main

import (
	"github.com/ducnt2212/chat-app-backend/internal/logger"
)

type Application struct {
	Addr   string
	logger *logger.Logger
}

func NewApplication(Addr string) *Application {
	return &Application{
		Addr:   Addr,
		logger: logger.GetLogger(),
	}
}

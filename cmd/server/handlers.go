package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ducnt2212/chat-app-backend/internal/helper"
	"github.com/ducnt2212/chat-app-backend/internal/models"
	"github.com/ducnt2212/chat-app-backend/internal/service"
)

func (app *Application) health(writer http.ResponseWriter, request *http.Request) {
	response := map[string]string{
		"response": "OK",
	}
	helper.ReplyJSON(writer, http.StatusOK, response)
}

func (app *Application) register(writer http.ResponseWriter, request *http.Request) {
	var registerForm struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&registerForm); err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	hashedPassword, err := helper.HashPassword(registerForm.Password)
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user := models.User{
		Username:       registerForm.Username,
		Email:          registerForm.Email,
		HashedPassword: hashedPassword,
	}

	err = app.userService.CreateUser(user)
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in creating user")
		return
	}

	response := map[string]string{
		"msg": "User created successfully",
	}
	helper.ReplyJSON(writer, http.StatusCreated, response)
}

func (app *Application) login(writer http.ResponseWriter, request *http.Request) {
	var loginForm struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&loginForm); err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	user, err := app.userService.GetUserByEmail(loginForm.Email)
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !helper.IsCorrectPassword(string(user.HashedPassword), loginForm.Password) {
		helper.ReplyJSONError(writer, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := helper.GenerateJwt(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response := map[string]string{
		"token": token,
	}
	helper.ReplyJSON(writer, http.StatusOK, response)
}

func (app *Application) createRoom(writer http.ResponseWriter, request *http.Request) {
	var createRoomForm struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(request.Body).Decode(&createRoomForm); err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	if createRoomForm.Name == "" {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Room name is required")
		return
	}

	userID := helper.GetUserID(request)

	room := models.Room{
		Name:      createRoomForm.Name,
		IsPrivate: false,
		CreatedBy: userID,
	}

	room, err := app.roomService.CreateRoom(room)
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in creating room")
		return
	}

	response := map[string]string{
		"msg": "Room created successfully",
	}
	helper.ReplyJSON(writer, http.StatusCreated, response)
}

func (app *Application) listRooms(writer http.ResponseWriter, request *http.Request) {
	rooms, err := app.roomService.ListRooms()
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in listing rooms")
		return
	}

	response := map[string][]models.Room{
		"rooms": rooms,
	}
	helper.ReplyJSON(writer, http.StatusOK, response)
}

func (app *Application) joinRoom(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.PathValue("roomID"))
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Room id")
		return
	}

	userID := helper.GetUserID(request)

	err = app.roomService.JoinRoom(roomID, userID)
	if err != nil {
		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in joining room")
		return
	}

	response := map[string]string{
		"msg": "Joined room successfully",
	}
	helper.ReplyJSON(writer, http.StatusOK, response)
}

func (app *Application) sendMessage(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.PathValue("roomID"))
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Not Found")
		return
	}

	userID := helper.GetUserID(request)

	var sendMessageForm struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(request.Body).Decode(&sendMessageForm); err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	if sendMessageForm.Content == "" {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Content is required")
		return
	}

	_, err = app.messageService.SendMessage(roomID, userID, sendMessageForm.Content)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			helper.ReplyJSONError(writer, http.StatusForbidden, "Forbidden")
			return
		}

		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in sending message")
		return
	}

	resposne := map[string]string{"msg": "Sent message successfully"}
	helper.ReplyJSON(writer, http.StatusCreated, resposne)
}

func (app *Application) getMessages(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.PathValue("roomID"))
	if err != nil {
		helper.ReplyJSONError(writer, http.StatusBadRequest, "Not Found")
		return
	}

	userID := helper.GetUserID(request)

	var limit int
	params := request.URL.Query()

	limitParam := params.Get("limit")
	if limitParam == "" {
		limit = 30
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Limit parameter")
			return
		}

		if limit < 1 {
			limit = 1
		} else if limit > 30 {
			limit = 30
		}
	}

	cursor := params.Get("cursor")
	if cursor != "" {
		_, err = time.Parse(time.RFC3339Nano, cursor)
		if err != nil {
			helper.ReplyJSONError(writer, http.StatusBadRequest, "Invalid Cursor parameter")
			return
		}
	}

	messages, nextCursor, err := app.messageService.ListMessages(roomID, userID, limit, cursor)
	if err != nil {
		if errors.Is(err, service.ErrForbidden) {
			helper.ReplyJSONError(writer, http.StatusForbidden, "Forbidden")
			return
		}

		app.logger.Error(err.Error())
		helper.ReplyJSONError(writer, http.StatusInternalServerError, "Server error in getting messages")
		return
	}

	response := map[string]any{
		"messages":    messages,
		"next_cursor": nextCursor,
	}
	helper.ReplyJSON(writer, http.StatusOK, response)
}

func (app *Application) serveWS(writer http.ResponseWriter, request *http.Request) {
	userID := helper.GetUserID(request)

	app.wsHandler.ServeWS(writer, request, userID)
}

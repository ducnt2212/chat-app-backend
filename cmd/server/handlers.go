package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/ducnt2212/chat-app-backend/internal/helper"
	"github.com/ducnt2212/chat-app-backend/internal/models"
)

func (app *Application) health(writer http.ResponseWriter, request *http.Request) {
	response := map[string]string{
		"resposne": "OK",
	}
	app.replyJson(writer, http.StatusOK, response)
}

func (app *Application) register(writer http.ResponseWriter, request *http.Request) {
	var registerForm struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&registerForm); err != nil {
		app.replyJsonError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	hashedPassword, err := helper.HashPassword(registerForm.Password)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	user := models.User{
		Username:       registerForm.Username,
		Email:          registerForm.Email,
		HashedPassword: hashedPassword,
	}

	id, err := app.repo.CreateUser(user)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Server error in creating user")
		return
	}

	app.logger.Info(fmt.Sprintf("Created username: %s with id: %d", registerForm.Username, id))

	response := map[string]string{
		"msg": "User created successfully",
	}
	app.replyJson(writer, http.StatusCreated, response)
}

func (app *Application) login(writer http.ResponseWriter, request *http.Request) {
	var loginForm struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(request.Body).Decode(&loginForm); err != nil {
		app.replyJsonError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	user, err := app.repo.GetUserByEmail(loginForm.Email)
	if err != nil {
		app.replyJsonError(writer, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	if !helper.IsCorrectPassword(string(user.HashedPassword), loginForm.Password) {
		app.replyJsonError(writer, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	token, err := helper.GenerateJwt(user.ID, os.Getenv("JWT_SECRET"))
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	response := map[string]string{
		"token": token,
	}
	app.replyJson(writer, http.StatusOK, response)
}

func (app *Application) createRoom(writer http.ResponseWriter, request *http.Request) {
	var createRoomForm struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(request.Body).Decode(&createRoomForm); err != nil {
		app.replyJsonError(writer, http.StatusBadRequest, "Inavalid Request")
		return
	}

	if createRoomForm.Name == "" {
		app.replyJsonError(writer, http.StatusBadRequest, "Room name is required")
		return
	}

	userID := request.Context().Value("user_id").(int)

	room := models.Room{
		Name:      createRoomForm.Name,
		IsPrivate: false,
		CreatedBy: userID,
	}

	err := app.repo.CreateRoom(room)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Server error in creating room")
		return
	}

	response := map[string]string{
		"msg": "Room created successfully",
	}
	app.replyJson(writer, http.StatusCreated, response)
}

func (app *Application) listRooms(writer http.ResponseWriter, request *http.Request) {
	rooms, err := app.repo.ListRooms()
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Server error in listing rooms")
		return
	}

	response := map[string]any{
		"rooms": rooms,
	}
	app.replyJson(writer, http.StatusOK, response)
}

func (app *Application) sendMessage(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.PathValue("roomID"))
	if err != nil {
		app.replyJsonError(writer, http.StatusNotFound, "Not Found")
		return
	}
	userID := request.Context().Value("user_id").(int)

	var sendMessageForm struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(request.Body).Decode(&sendMessageForm); err != nil {
		app.replyJsonError(writer, http.StatusBadRequest, "Invalid Request")
		return
	}

	if sendMessageForm.Content == "" {
		app.replyJsonError(writer, http.StatusBadRequest, "Content is required")
		return
	}

	msg := models.Message{
		RoomID:   roomID,
		SenderID: userID,
		Content:  sendMessageForm.Content,
	}

	err = app.repo.CreateMessage(msg)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Server error in sending message")
		return
	}

	resposne := map[string]string{"msg": "Sent message successfully"}
	app.replyJson(writer, http.StatusCreated, resposne)
}

func (app *Application) getMessages(writer http.ResponseWriter, request *http.Request) {
	roomID, err := strconv.Atoi(request.PathValue("roomID"))
	if err != nil {
		app.replyJsonError(writer, http.StatusBadRequest, "Not Found")
		return
	}

	var limit int
	params := request.URL.Query()

	limitParam := params.Get("limit")
	if limitParam == "" {
		limit = 30
	} else {
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			app.replyJsonError(writer, http.StatusBadRequest, "Invalid Limit parameter")
			return
		}
	}

	cursorParam := params.Get("cursor")
	if cursorParam != "" {

		_, err = time.Parse(time.RFC3339Nano, cursorParam)
		if err != nil {
			app.replyJsonError(writer, http.StatusBadRequest, "Invalid Cursor parameter")
			return
		}
	}

	messages, nextCursor, err := app.repo.ListMessagesByRoom(roomID, limit, cursorParam)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJsonError(writer, http.StatusInternalServerError, "Server error in getting messages")
		return
	}

	response := map[string]any{
		"messages":    messages,
		"next_cursor": nextCursor,
	}
	app.replyJson(writer, http.StatusOK, response)
}

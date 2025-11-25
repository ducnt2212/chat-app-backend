package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/ducnt2212/chat-app-backend/internal/helper"
	"github.com/ducnt2212/chat-app-backend/internal/models"
)

func (app *Application) health(writer http.ResponseWriter, request *http.Request) {
	app.replyJson(writer, http.StatusOK, "{\"message\": \"OK\"}")
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

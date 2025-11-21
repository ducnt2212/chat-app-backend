package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (app *Application) health(writer http.ResponseWriter, request *http.Request) {
	app.replyJson(writer, http.StatusOK, "{\"message\": \"OK\"}")
}

func (app *Application) register(writer http.ResponseWriter, request *http.Request) {
	var userRegisterForm struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.NewDecoder(request.Body).Decode(&userRegisterForm)

	// Register logic
	// TODO: Be careful with empty string
	id, err := app.repo.CreateUser(userRegisterForm.Username, userRegisterForm.Email, userRegisterForm.Password)
	if err != nil {
		app.logger.Error(err.Error())
		app.replyJson(writer, http.StatusInternalServerError, `{"error":"server error in creating user"}`)
		return
	}

	app.logger.Info(fmt.Sprintf("Created username: %s with id: %d", userRegisterForm.Username, id))

	app.replyJson(writer, http.StatusOK, userRegisterForm)
}

func (app *Application) login(writer http.ResponseWriter, request *http.Request) {
	var UserLoginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	json.NewDecoder(request.Body).Decode(&UserLoginForm)

	// Login logic

	response := fmt.Sprintf("Logged in with Username: %s, Password: %s", UserLoginForm.Username, UserLoginForm.Password)
	app.replyJson(writer, http.StatusOK, response)
}

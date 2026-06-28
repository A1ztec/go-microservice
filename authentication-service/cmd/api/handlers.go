package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	//Todo authenticate the user
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		log.Printf("this error in the readJson: %v", err)
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}
	// validate the user agains the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		log.Printf("this error in the GetByEmail: %v", err)
		app.errorJson(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		log.Printf("this error in the PasswordMatches: %v", err)
		app.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusOK, payload)
}

package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
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

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	go func() {
		defer cancel()
		err := app.logRequest(ctx, "authentication", fmt.Sprintf("Logged in user %s", user.Email))
		if err != nil {
			log.Printf("Error logging request: %v", err)
		}
	}()
	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data:    user,
	}
	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) logRequest(ctx context.Context, name, data string) error {

	var entry struct {
		Name string `json:"name"`
		Data string `json:"data"`
	}
	entry.Name = name
	entry.Data = data

	jsonData, err := json.MarshalIndent(entry, "", "\t")
	if err != nil {
		log.Printf("Error marshalling log entry: %v", err)
		return err
	}
	logServiceUrl := "http://logger-service/log"
	request, err := http.NewRequestWithContext(ctx, "POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		return err
	}
	_, err = app.Client.Do(request)
	if err != nil {
		log.Println("error in sending the request to the logger microservice", err)
		return err
	}
	return nil
}

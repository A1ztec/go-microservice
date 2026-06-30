package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type requestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth,omitempty"`
	Log    LogPayload  `json:"log,omitempty"`
}

type AuthPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}
func (app *Config) HandleSubmission(w http.ResponseWriter, r *http.Request) {

	var requestPayload requestPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err, http.StatusBadRequest)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)
	case "log":
		app.logItem(w, requestPayload.Log)
		log.Printf("this is the log data: %v", requestPayload.Log)
	default:
		app.errorJson(w, errors.New("Invalid Action"), http.StatusBadRequest)
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create somme json and send to auth microservice

	jsonData, _ := json.MarshalIndent(a, "", "\t")
	//call the service
	request, err := http.NewRequest("POST", "http://authentication-service/authenticate", bytes.NewBuffer(jsonData))
	if err != nil {
		app.errorJson(w, err)
		return
	}
	response, err := app.Client.Do(request)
	if err != nil {
		app.errorJson(w, err)
		return
	}
	defer response.Body.Close()

	log.Printf("DEBUG: Auth Service responded with Status Code: %d", response.StatusCode)

	//make sure we get back the correct status code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("Invalid Credentials"), http.StatusUnauthorized)
		return
	} else if response.StatusCode != http.StatusOK {
		app.errorJson(w, errors.New(response.Status), http.StatusInternalServerError)
		return
	}
	// create a variable that we can read the response body into
	var jsonFromService jsonResponse

	// decode the json from the auth service
	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	// if we get an error in the response, send it back
	if jsonFromService.Error {
		app.errorJson(w, errors.New(jsonFromService.Message), http.StatusUnauthorized)
		return
	}
	// create a new payload and send it back
	payload := jsonResponse{
		Error:   false,
		Message: "Authenticated",
		Data:    jsonFromService.Data,
	}
	app.writeJson(w, http.StatusOK, payload)
}

func (app *Config) logItem(w http.ResponseWriter, l LogPayload) {
	// get the body and encode it
	jsonData, err := json.MarshalIndent(l, "", "\t")
	if err != nil {
		log.Printf("Error marshalling request body: %v", err)
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	logServiceUrl := "http://logger-service/log"
	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Printf("Error creating new request: %v", err)
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	request.Header.Set("Content-Type", "application/json")
	res, err := app.Client.Do(request)
	if err != nil {
		log.Println("error in sending the request to the logger microservice", err)
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Printf("Logger service responded with status code: %d", res.StatusCode)
		app.errorJson(w, errors.New("Error calling logger service"), http.StatusInternalServerError)
		return
	}

	var jsonFromService jsonResponse
	err = json.NewDecoder(res.Body).Decode(&jsonFromService)
	if err != nil {
		log.Printf("Error decoding response from logger service: %v", err)
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}

	if jsonFromService.Error {
		log.Printf("Logger service returned an error: %s", jsonFromService.Message)
		app.errorJson(w, errors.New(jsonFromService.Message), http.StatusInternalServerError)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: "Logged",
	}
	err = app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
		app.errorJson(w, err, http.StatusInternalServerError)
		return
	}
}

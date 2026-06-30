package main

import (
	"log"
	"log-service/data"
	"net/http"
)

type jsonPayload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func (app *Config) Writelog(w http.ResponseWriter, r *http.Request) {
	// read json from request body
	var requestPayload jsonPayload
	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		log.Printf("Error reading JSON: %v", err)
		app.errorJson(w, err)
		return
	}
	// create a log entry
	event := data.LogEntry{
		Name: requestPayload.Name,
		Data: requestPayload.Data,
	}
	err = app.Models.LogEntry.Insert(event)
	if err != nil {
		log.Printf("Error inserting log entry: %v", err)
		app.errorJson(w, err)
		return
	}
	payload := jsonResponse{
		Error:   false,
		Message: "logged",
	}
	err = app.writeJson(w, http.StatusOK, payload)
	if err != nil {
		log.Printf("Error writing JSON response: %v", err)
		app.errorJson(w, err)
		return
	}
}

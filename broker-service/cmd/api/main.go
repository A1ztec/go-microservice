package main

import (
	"fmt"
	"log"
	"net/http"
)

const portNumber = "80"

type Config struct {
}

func main() {
	app := Config{}
	log.Printf("starting broker service on port %s\n", portNumber)

	// define a http server and set the handler
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", portNumber),
		Handler: app.routes(),
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

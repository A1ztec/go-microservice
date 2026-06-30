package main

import (
	"context"
	"fmt"
	"log"
	"log-service/data"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	// connect to mongo
	c, err := connectToMongo()
	log.Println("Connected to MongoDB...", c)
	if err != nil {
		log.Panic(err)
	}
	//create a context in order to disconnect from mongo when the application is stopped
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	// close connection to mongo
	defer func() {
		if err = c.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
	app := Config{
		Models: data.New(c),
	}
	log.Println("Starting logger service... on port", os.Getenv("PORT_NUMBER"))
	// go app.serve()
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", os.Getenv("PORT_NUMBER")),
		Handler: app.routes(),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func (app *Config) serve() {
	// src := &http.Server{
	// 	Addr:    fmt.Sprintf(":%s", os.Getenv("PORT_NUMBER")),
	// 	Handler: app.routes(),
	// }
	// err := src.ListenAndServe()
	// if err != nil {
	// 	log.Panic(err)
	// }
}

func connectToMongo() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(os.Getenv("MONGOURL"))
	log.Println("Connecting to MongoDB... , ", os.Getenv("MONGOURL"))
	clientOptions.SetAuth(options.Credential{
		Username: os.Getenv("MONGOUSER"),
		Password: os.Getenv("MONGOPASSWORD"),
	})
	c, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Println("error in connect to mongo")
		return nil, err
	}
	log.Println("Successfully connected to MongoDB")
	return c, nil
}

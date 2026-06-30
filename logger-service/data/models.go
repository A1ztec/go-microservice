package data

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func New(mongo *mongo.Client) Models {
	client = mongo
	return Models{
		LogEntry: LogEntry{},
	}
}

type Models struct {
	LogEntry LogEntry
}

type LogEntry struct {
	ID        string    `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string    `bson:"name" json:"name"`
	Data      string    `bson:"data" json:"data"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

func (l *LogEntry) Insert(entry LogEntry) error {
	collection := client.Database("logs").Collection("logs")
	_, err := collection.InsertOne(context.TODO(), LogEntry{
		Name:      entry.Name,
		Data:      entry.Data,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		log.Println("Error inserting log entry", err)
		return err
	}
	return nil
}
func (l *LogEntry) All() ([]*LogEntry, error) {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := collection.Find(ctx, bson.D{}, opts)
	if err != nil {
		log.Println("Error getting all log entries", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var logs []*LogEntry
	for cursor.Next(ctx) {
		var logEntry LogEntry
		err := cursor.Decode(&logEntry)
		if err != nil {
			log.Println("Error decoding log entry into slice", err)
			return nil, err
		}
		logs = append(logs, &logEntry)
	}
	return logs, nil
}
func (l *LogEntry) GetByID(id string) (*LogEntry, error) {
	// create collection handle
	collection := client.Database("logs").Collection("logs")
	// create context
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// convert id string to ObjectID
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println("Error converting ID to ObjectID", err)
		return nil, err
	}

	// find the document by ID
	var logEntry LogEntry
	err = collection.FindOne(ctx, bson.M{"_id": docId}).Decode(&logEntry)
	if err != nil {
		log.Println("Error finding log entry by ID", err)
		return nil, err
	}
	return &logEntry, nil
}
func (l *LogEntry) DropCollection() error {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := collection.Drop(ctx)
	if err != nil {
		log.Println("Error dropping collection", err)
		return err
	}
	return nil
}

func (l *LogEntry) Update(entry LogEntry) (*mongo.UpdateResult, error) {
	collection := client.Database("logs").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// convert id string to ObjectID
	docId, err := primitive.ObjectIDFromHex(l.ID)
	if err != nil {
		log.Println("Error converting ID to ObjectID", err)
		return nil, err
	}

	// create update document
	update := bson.M{
		"$set": bson.M{
			"name":       entry.Name,
			"data":       entry.Data,
			"updated_at": time.Now(),
		},
	}

	// update the document by ID
	result, err := collection.UpdateOne(ctx, bson.M{"_id": docId}, update)
	if err != nil {
		log.Println("Error updating log entry", err)
		return nil, err
	}
	return result, nil
}

package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBEventData struct {
	Agent           string    `bson:"agent"`
	ChainIDs        []string  `bson:"chainIds"`
	Timestamp       time.Time `bson:"timestamp"`
	TransactionHash string    `bson:"transactionHash"`
	LogIndex        int       `bson:"logIndex"`
}

type EventGetFilter struct {
	TransactionHash string
	LogIndex        int
}

type EventCollection struct {
	Collection *mongo.Collection
}

var eventsCollection *EventCollection

func EventServiceInstance() *EventCollection {
	if eventsCollection == nil {
		collection := GetCollection("events")
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "transactionHash", Value: 1}, {Key: "logIndex", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			panic(err)
		}
		fmt.Println("Index Created: " + name)
		eventsCollection = &EventCollection{
			Collection: collection,
		}
	}

	return eventsCollection
}

func (collection *EventCollection) Add(event DBEventData) error {
	_, err := collection.Collection.InsertOne(context.TODO(), event)
	if err != nil {
		log.Printf("Failed to insert event into MongoDB: %v", err)
		return err
	}
	return nil
}

func (collection *EventCollection) Get(eventFilter EventGetFilter) (DBEventData, error) {
	var event DBEventData

	err := collection.Collection.FindOne(context.TODO(), bson.M{"transactionHash": eventFilter.TransactionHash, "logIndex": eventFilter.LogIndex}).Decode(&event)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return DBEventData{}, err
		}
		return DBEventData{}, err
	}

	return event, nil
}

func (collection *EventCollection) Update() error {
	return fmt.Errorf("missing implementation")
}

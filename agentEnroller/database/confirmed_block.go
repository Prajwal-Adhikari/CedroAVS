package db

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConfirmedBlock struct {
	Chain       string `bson:"chain"`
	BlockNumber int64  `bson:"blockNumber"`
}

type ConfirmedBlockCollection struct {
	Collection *mongo.Collection
}

var confirmedBlock *ConfirmedBlockCollection

func ConfirmedBlockInstance() *ConfirmedBlockCollection {
	if confirmedBlock == nil {
		collection := GetCollection("confirmedBlocks")
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "chain", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			panic(err)
		}
		fmt.Println("Index Created: " + name)
		confirmedBlock = &ConfirmedBlockCollection{
			Collection: collection,
		}
	}

	return confirmedBlock
}

func (collection *ConfirmedBlockCollection) Add(block ConfirmedBlock) error {
	_, err := collection.Collection.InsertOne(context.TODO(), block)
	if err != nil {
		log.Printf("Failed to insert confirmed block into MongoDB: %v", err)
		return err
	}
	return nil
}

func (collection *ConfirmedBlockCollection) Get(chain string) (ConfirmedBlock, error) {
	var confirmedBlock ConfirmedBlock

	err := collection.Collection.FindOne(context.TODO(), bson.M{"chain": chain}).Decode(&confirmedBlock)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return ConfirmedBlock{}, err
		}
		return ConfirmedBlock{}, err
	}

	return confirmedBlock, nil
}

func (collection *ConfirmedBlockCollection) Update(chainName string, blockNumber int64) error {
	_, err := collection.Collection.UpdateOne(context.TODO(), bson.D{{Key: "chain", Value: chainName}}, bson.D{{Key: "$set", Value: bson.M{"blockNumber": blockNumber}}}, options.Update().SetUpsert(true))

	if err != nil {
		log.Printf("Failed to insert confirmed block into MongoDB: %v", err)

		return err
	}

	return nil
}

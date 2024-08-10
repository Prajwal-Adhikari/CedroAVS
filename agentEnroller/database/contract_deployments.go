package db

import (
	"context"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DeploymentNotification struct {
	Chain           string `bson:"chain"`
	TransactionHash string `bson:"transactionHash"`
	ManagerAddress  string `bson:"managerAddress"`
	Status          string `bson:"status"`
	Error           string `bson:"error"`
}

type ChainDeployment struct {
	Chain           string `bson:"chain"`
	ContractAddress string `bson:"contractAddress"`
	TransactionHash string `bson:"transactionHash"`
	ManagerAddress  string `bson:"managerAddress"`
	Status          string `bson:"status"`
	Error           string `bson:"error"`
}

type DBContractDeploymentData struct {
	Owner                  string                   `bson:"owner"`
	Chains                 []ChainDeployment        `bson:"chains"`
	DeploymentNotification []DeploymentNotification `bson:"deploymentNotifications"`
}

type ContractDeploymentGetFilter struct {
	owner common.Address
}

type ContractDeploymentCollection struct {
	Collection *mongo.Collection
}

var contractDeploymentCollection *ContractDeploymentCollection

func ContractDeploymentInstance() *ContractDeploymentCollection {
	if contractDeploymentCollection == nil {
		collection := GetCollection("deployments")
		indexModel := mongo.IndexModel{
			Keys:    bson.D{{Key: "owner", Value: 1}},
			Options: options.Index().SetUnique(true),
		}
		name, err := collection.Indexes().CreateOne(context.TODO(), indexModel)
		if err != nil {
			panic(err)
		}
		fmt.Println("Index Created: " + name)

		contractDeploymentCollection = &ContractDeploymentCollection{
			Collection: collection,
		}
	}

	return contractDeploymentCollection
}

func (collection *ContractDeploymentCollection) Add(contractDeployment DBContractDeploymentData) error {
	_, err := collection.Collection.InsertOne(context.TODO(), contractDeployment)
	if err != nil {
		log.Printf("Failed to insert contract deployment into MongoDB: %v", err)
		return err
	}
	return nil
}

func (collection *ContractDeploymentCollection) Get(contractDeloymentFilter ContractDeploymentGetFilter) (DBContractDeploymentData, error) {
	var contractDeployment DBContractDeploymentData

	err := collection.Collection.FindOne(context.TODO(), bson.M{"owner": contractDeloymentFilter.owner}).Decode(&contractDeployment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return DBContractDeploymentData{}, err
		}
		return DBContractDeploymentData{}, err
	}

	return contractDeployment, nil
}

func (collection *ContractDeploymentCollection) Update() error {
	return fmt.Errorf("missing implementation")
}

func (collection *ContractDeploymentCollection) UpdateDeploymentNotification(owner common.Address, deploymentNotification DeploymentNotification) error {
	filter := bson.M{
		"owner": owner.Hex(),
	}

	update := bson.M{
		"$push": bson.M{
			"deploymentNotifications": deploymentNotification,
		},
	}

	// Define the options to update a single document
	opts := options.Update().SetUpsert(false)

	// Perform the update
	_, err := collection.Collection.UpdateOne(context.Background(), filter, update, opts)
	if err != nil {
		log.Printf("Failed to update deployment notification: %v", err)
		return err
	}

	return nil
}

func (cdc *ContractDeploymentCollection) BulkUpdateDeploymentNotification(owner common.Address, deploymentNotifications []DeploymentNotification) error {
	var models []mongo.WriteModel

	for _, notification := range deploymentNotifications {
		filter := bson.M{
			"owner": owner.String(),
		}

		update := bson.M{
			"$push": bson.M{
				"deploymentNotifications": notification,
			},
		}

		updateModel := mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(true)

		models = append(models, updateModel)
	}

	bulkOption := options.BulkWrite().SetOrdered(true)
	_, err := cdc.Collection.BulkWrite(context.Background(), models, bulkOption)
	if err != nil {
		log.Printf("Failed to perform bulk update: %v", err)
		return err
	}

	// Debug: Fetch the document after update
	var afterUpdate DBContractDeploymentData
	err = cdc.Collection.FindOne(context.Background(), bson.M{"owner": owner}).Decode(&afterUpdate)
	if err != nil {
		return fmt.Errorf("failed to fetch document after update: %v", err)
	}
	log.Printf("After Update: %+v\n", afterUpdate)

	return nil
}

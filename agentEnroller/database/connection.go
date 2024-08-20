package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/config"
)

type DBConnection struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var connection *DBConnection

func InitDBConnection() {
	appConfig := config.LoadConfig()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := appConfig.DBHost
	dbPort := appConfig.DBPort
	dbUser := appConfig.DBUser
	dbPass := appConfig.DBPass
	dbName := appConfig.DBName

	mongodbURI := fmt.Sprintf("mongodb+srv://%s:%s@%s", dbUser, dbPass, dbHost)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))

	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	database := client.Database(dbName)

	connection = &DBConnection{
		Client:   client,
		Database: database,
	}
}

func GetDBConnection() *DBConnection {
	if connection == nil {
		log.Fatal("Database connection is not initialized")
	}
	return connection
}

func GetCollection(collectionName string) *mongo.Collection {
	dbConn := GetDBConnection()
	return dbConn.Database.Collection(collectionName)
}

func CloseDBConnection() {
	if err := connection.Client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Failed to disconnect MongoDB client: %v", err)
	}
}

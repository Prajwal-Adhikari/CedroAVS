package db

import (
	"context"
	"fmt"
	"log"
	"time"

	env "github.com/Cedro-Finance/CedroAVS/aggregator/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var DbClient *mongo.Database

func InitDb() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbHost := env.AppEnv.DBHost
	dbPort := env.AppEnv.DBPort
	dbUser := env.AppEnv.DBUser
	dbPass := env.AppEnv.DBPass
	dbName := env.AppEnv.DBName

	mongodbURI := fmt.Sprintf("mongodb://%s:%s@%s:%s", dbUser, dbPass, dbHost, dbPort)

	if dbUser == "" || dbPass == "" {
		mongodbURI = fmt.Sprintf("mongodb://%s:%s", dbHost, dbPort)
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongodbURI))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	fmt.Println("database connected succesfully")

	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	DbClient = client.Database(dbName)
}

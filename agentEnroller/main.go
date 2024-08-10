package main

import (
	"log"

	db "github.com/Cedro-Finance/CedroAVS/agentEnroller/database"
	"github.com/Cedro-Finance/CedroAVS/agentEnroller/eventlistener"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	// Initialize the DB connection
	db.InitDBConnection()
	defer db.CloseDBConnection()

	// Intitalize eventlisteners
	eventlistener.SyncToLatest()
	eventlistener.Start()
}

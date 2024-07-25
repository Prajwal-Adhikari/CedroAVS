package main

import (
	"log"

	"github.com/Cedro-Finance/CedroAVS/agentEnroller/eventlistener"
	"github.com/lpernett/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	eventlistener.Start()
}

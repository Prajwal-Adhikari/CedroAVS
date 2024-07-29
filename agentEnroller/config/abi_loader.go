package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type ContractData struct {
	ABI      json.RawMessage `json:"abi"`
	Bytecode string          `json:"bytecode"`
}

func LoadContractData(filePath string) (*ContractData, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var contractData ContractData
	err = json.Unmarshal(byteValue, &contractData)
	if err != nil {
		return nil, err
	}

	return &contractData, nil
}

func GetAgentPodAbi() (string, string) {
	contractData, err := LoadContractData("data/abi/agent_pod.json")
	if err != nil {
		log.Fatalf("Failed to load contract data: %v", err)
	}

	parsedABI := string(contractData.ABI)
	bytecode := string(contractData.Bytecode)

	return parsedABI, bytecode
}

func GetAgentManagerAbi() string {
	contractData, err := LoadContractData("data/abi/agent_manage.json")
	if err != nil {
		log.Fatalf("Failed to load contract data: %v", err)
	}

	parsedABI := string(contractData.ABI)

	return parsedABI
}

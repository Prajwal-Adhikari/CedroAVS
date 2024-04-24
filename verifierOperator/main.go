package main

import (
	"fmt"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/keys"
)

func main() {
	configFile := config.LoadConfig("config.yaml")
	_ = configFile
	keystore := keys.LoadKeys(configFile.KeyStore)
	fmt.Println(keystore)
}

package actions

import (
	"fmt"

	"github.com/Cedro-Finance/CedroAVS/verifier/config"
	"github.com/Cedro-Finance/CedroAVS/verifier/operator"
	"github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/urfave/cli"
)

func RegisterOperatorWithAvs(ctx *cli.Context) error {

	configPath := ctx.GlobalString("config")

	configFile, err := config.LoadConfig(configPath)

	if err != nil {
		fmt.Print(err)
		panic("Failed loading config")
	}

	_ = configFile

	op, err := operator.NewOperator(*configFile)

	if err != nil {
		panic("Counldn't start operator")
	}

	operatorEcdsaPrivate, err := ecdsa.ReadKey(configFile.KeyStore.Path, configFile.KeyStore.Password)

	if err != nil {
		panic("No keystore is found")
	}

	err = op.RegisterOperatorWithAVS(operatorEcdsaPrivate)

	if err != nil {
		return err
	}

	// configPath := ctx.GlobalString("")
	// nodeConfig := config.Node

	return nil
}

func DeregisterOperatorWithAvs(ctx *cli.Context) error {

	configPath := ctx.GlobalString("config")

	configFile, err := config.LoadConfig(configPath)

	if err != nil {
		fmt.Print(err)
		panic("Failed loading config")
	}

	_ = configFile

	op, err := operator.NewOperator(*configFile)

	if err != nil {
		panic("Counldn't start operator")
	}

	operatorEcdsaPrivate, err := ecdsa.ReadKey(configFile.KeyStore.Path, configFile.KeyStore.Password)

	if err != nil {
		panic("No keystore is found")
	}

	err = op.DeregisterOperatorWithAVS(operatorEcdsaPrivate)

	if err != nil {
		return err
	}

	// configPath := ctx.GlobalString("")
	// nodeConfig := config.Node

	return nil
}



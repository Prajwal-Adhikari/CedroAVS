package operator

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"time"

	"github.com/Cedro-Finance/CedroAVS/core"
	contractRegistryCoordinator "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	"github.com/ethereum/go-ethereum/crypto"
)

func (o *Operator) RegisterOperatorWithAVS(operatorEcdsaKeyPair *ecdsa.PrivateKey) error {

	// hardcode these things for now
	quorumNumbers := []byte{0,1}
	socket := ""

	fmt.Println(core.ConvertToBN254G1Point(o.blsKeyPair.GetPubKeyG1()))
	fmt.Println(core.ConvertToBN254G2Point(o.blsKeyPair.GetPubKeyG2()))
	// Generate salt and expiry
	privateKeyBytes := []byte(o.blsKeyPair.PrivKey.String())
	salt := [32]byte{}
	copy(salt[:], crypto.Keccak256([]byte("churn"), []byte(time.Now().String()), quorumNumbers[:], privateKeyBytes))

	operatorToAvsRegistrationSigSalt := salt

	curBlockNum, err := o.avs.Eth.BlockNumber(context.Background())
	if err != nil {
		o.logger.Errorf("Unable to get current block number")
		return err
	}
	curBlock, err := o.avs.Eth.BlockByNumber(context.Background(), big.NewInt(int64(curBlockNum)))
	if err != nil {
		o.logger.Errorf("Unable to get current block")
		return err
	}
	sigValidForSeconds := int64(1_000_000)
	operatorToAvsRegistrationSigExpiry := big.NewInt(int64(curBlock.Time()) + sigValidForSeconds)

	o.logger.Info(
		"RegisterOperatorInQuorumWithAVSRegistryCoordinator",
		"quorumNumbers", quorumNumbers[0],
		"socket", socket,
		"operatorToAvsRegistrationSigExpiry", operatorToAvsRegistrationSigExpiry,
	)

	fmt.Println(operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry)

	quorumNumbersToSDK := core.ConvertQuorumNumbersFromBytes(quorumNumbers)

	_, err = o.avs.AvsWriter.RegisterOperatorInQuorumWithAVSRegistryCoordinator(
		context.Background(),
		operatorEcdsaKeyPair, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry,
		o.blsKeyPair, quorumNumbersToSDK, socket,
	)
	if err != nil {
		o.logger.Error("Unable to register operator with avs registry coordinator", err)
		return err
	}
	o.logger.Infof("Registered operator with avs registry coordinator.")

	return nil
}

func (o *Operator) DeregisterOperatorWithAVS(operatorEcdsaKeypair *ecdsa.PrivateKey) error {
	quorumNumbers := []byte{1}

	o.logger.Info(
		"DeregisterOperatorInQuorumWithAVSRegistryCoordinator",
		"quorumNumbers", quorumNumbers[0],
	)

	quorumNumbersToSDK := core.ConvertQuorumNumbersFromBytes(quorumNumbers)

	_, err := o.avs.AvsWriter.DeregisterOperator(
		context.Background(), quorumNumbersToSDK, contractRegistryCoordinator.BN254G1Point{},
	)
	if err != nil {
		o.logger.Error("Unable to register operator with avs registry coordinator", err)
		return err
	}
	o.logger.Infof("Deregistered operator with avs registry coordinator.")

	return nil
}

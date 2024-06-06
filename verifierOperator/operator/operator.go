package operator

import (
	"context"
	"fmt"
	"math/big"
	"os"

	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/agents"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/common"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/keys"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/logger"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/wallet"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	sdkEcdsa "github.com/Layr-Labs/eigensdk-go/crypto/ecdsa"
	"github.com/Layr-Labs/eigensdk-go/signerv2"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	bls "github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdklogging "github.com/Layr-Labs/eigensdk-go/logging"
)

var RegisteredClients map[string]func(config.ChainConfig) common.IClient

type Operator struct {
	avs        *avs.AVSClient
	keystore   *keys.KeyStore
	blsKeyPair *bls.KeyPair

	operatorId                sdktypes.OperatorId
	txMgr                     txmgr.TxManager
	logger                    sdklogging.Logger
	gRPCAggreagatorIpPortAddr string
}

func (r *Operator) mockAgentsForAgents(_ string) map[string]common.IClient {
	return map[string]common.IClient{}
}

func (r *Operator) GetChainBalances(ctx context.Context, eventCh <-chan *avs.Request) {

	for {

		balances := make(map[*big.Int][]agents.BalanceResponse)
		reqEvent := <-eventCh

		balances[reqEvent.ReqId] = []agents.BalanceResponse{}

		agentAddress := reqEvent.AgentAddress

		//TODO: Find the agents for the agentAddress
		reqAgets := r.mockAgentsForAgents(agentAddress)

		balCh := make(chan *agents.BalanceResponse)

		for _, ag := range reqAgets {
			go ag.GetBalance(ctx, reqEvent, balCh)
		}

		go r.AggregateBalanceResponse(ctx, balCh, reqAgets)
	}

}

func (r *Operator) AggregateBalanceResponse(ctx context.Context, balCh <-chan *agents.BalanceResponse, reqAgents map[string]common.IClient) {
	reqIdToResp := make(map[*big.Int][]*agents.BalanceResponse)

	for {
		balResponse := <-balCh

		value, ok := reqIdToResp[balResponse.ReqId]

		if !ok {
			reqIdToResp[balResponse.ReqId] = []*agents.BalanceResponse{balResponse}
		} else {
			value = append(value, balResponse)
			reqIdToResp[balResponse.ReqId] = value
		}

		if len(reqAgents) == len(reqIdToResp) {
			//Aggregated for one reqId

			var totalBalance = big.NewInt(0)

			for _, bal := range reqIdToResp[balResponse.ReqId] {
				totalBalance = totalBalance.Add(totalBalance, bal.TotalBalance)
			}

			//send totalBalance to one channel & then close

			delete(reqIdToResp, balResponse.ReqId)
			return
		}
	}
}

func (r *Operator) SignAndSendBalanceResponse(ctx context.Context, aggr <-chan []*agents.BalanceResponse) {
	for {
		// aggregated := <-aggr

		aggregatedBlnc := agents.AggregateBalanceResponse{
			ReqId: big.NewInt(1),
		}

		hash, err := aggregatedBlnc.SignHash()

		if err != nil {
			logger.GetLogger().Error("Error whole hasing the message")

			return
		}

		signature := r.blsKeyPair.SignMessage(hash)

		fmt.Print(signature)

		//SIgn THE MESSAGE HERE IN WHAT FORMAT?

		//MAKE HTTP CALL TO THE AGGREGATOR SERVICE

	}
}

func NewOperator(cfg config.Config) (*Operator, error) {
	var logLevel sdklogging.LogLevel

	if cfg.Production {
		logLevel = sdklogging.Production
	} else {
		logLevel = sdklogging.Development
	}

	logger, err := logger.NewZapLogger(logLevel)

	if err != nil {
		return nil, err
	}

	avsClient := avs.NewAVSClient(&cfg.AVS, logger)

	if err != nil {
		return nil, err
	}

	ethRpcClient, err := eth.NewClient(cfg.AVS.EthRpcUrl)

	if err != nil {
		panic(err)
	}

	//Operator Key Management

	operatorAddress, err := sdkEcdsa.GetAddressFromKeyStoreFile("")

	if err != nil {
		panic(err)
	}

	ecdsaKeypassword, ok := os.LookupEnv("OPERATOR_ECDSA_KEY_PASSWORD")

	if !ok {
		logger.Warnf("OPERATOR_ECDSA_KEY_PASSWORD env var not set")
	}

	signerConfig := signerv2.Config{
		KeystorePath: "",
		Password:     ecdsaKeypassword,
	}

	signerV2, _, err := signerv2.SignerFromConfig(signerConfig, big.NewInt(1))

	if err != nil {
		panic(err)
	}

	privateKey, err := sdkEcdsa.ReadKey(signerConfig.KeystorePath, signerConfig.Password)

	if err != nil {
		return nil, err
	}

	txSender, err := wallet.NewPrivateKeyWallet(ethRpcClient, signerV2, operatorAddress, logger)

	if err != nil {
		return nil, err
	}

	txMgr := txmgr.NewSimpleTxManager(txSender, ethRpcClient, logger, operatorAddress)

	chainioConfig := clients.BuildAllConfig{
		EthHttpUrl:                 cfg.AVS.EthRpcUrl,
		RegistryCoordinatorAddr:    cfg.AVS.RegistryCoordinatorAddr.String(),
		OperatorStateRetrieverAddr: cfg.AVS.OperatorStateRetrieverAddr.String(),
		AvsName:                    "CedroAVS",
	}

	sdkClients, err := clients.BuildAll(chainioConfig, privateKey, logger)

	operatorId, err := sdkClients.AvsRegistryChainReader.GetOperatorId(&bind.CallOpts{}, operatorAddress)

	if err != nil {
		logger.Error("Cannot get operator id", "err", err)
		return nil, err
	}

	return &Operator{
		avs:                       avsClient,
		blsKeyPair:                cfg.BlsPrivateKey,
		keystore:                  &keys.KeyStore{},
		operatorId:                operatorId,
		txMgr:                     txMgr,
		logger:                    logger,
		gRPCAggreagatorIpPortAddr: cfg.GRpcAggregatorIpAddr,
	}, nil

}

func StartRelay(ctx context.Context, cfg *config.Config) {

	op, err := NewOperator(*cfg)

	if err != nil {
		panic("Couldn't start operator")
	}

	op.logger.Infof("Starting Operator")

	// balanceReqCh := make(chan *avs.Request) //Chanel to accept request from AVS contract
	// // deviationReqCh := make(chan *avs.Request) //Chanel to accept request from offline trigger

	// regAgents := make(map[string]common.IClient)
	// for _, cfg := range cfg.Agents {
	// 	regAgents[cfg.Name] = RegisteredClients[cfg.Name](cfg)
	// }

	// go op.avs.PollEvents(ctx, balanceReqCh)
	// go op.GetChainBalances(ctx, balanceReqCh)
}

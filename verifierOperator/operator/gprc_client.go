package operator

import (
	// "fmt"
	"context"
	"fmt"
	"net/rpc"
	"time"

	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Cedro-Finance/CedroAVS/verifierOperator/config"
	"github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	"github.com/ethereum/go-ethereum/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	// "google.golang.org/grpc"
	// "google.golang.org/grpc/credentials/insecure"
)

type AggregatorGRpcClient struct {
	rpcClient *rpc.Client

	logger     logging.Logger
	operatorId sdktypes.OperatorId

	operatorAddr               common.Address
	OperatorStateRetrieverAddr common.Address
	RegistryCoodrinatorAddr    common.Address

	gRPCAggreagatorIpPortAddr string
	timeout                   time.Duration

	config config.Config
}

func NewAggregatorGRpcClient(config config.Config, operatorId sdktypes.OperatorId, opearatorAddr common.Address, logger logging.Logger) (*AggregatorGRpcClient, error) {
	return &AggregatorGRpcClient{
		rpcClient:                  nil,
		logger:                     logger,
		config:                     config,
		operatorId:                 operatorId,
		operatorAddr:               opearatorAddr,
		OperatorStateRetrieverAddr: config.AVS.OperatorStateRetrieverAddr,
		RegistryCoodrinatorAddr:    config.AVS.RegistryCoordinatorAddr,
		gRPCAggreagatorIpPortAddr:  "",
		timeout:                    1 * time.Second,
	}, nil
}


func (cl *AggregatorGRpcClient) InitOperatorToAggregator() error {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)



	nodeCtx, cancel := context.WithTimeout(context.Background(), cl.timeout)

	defer cancel()

	request := &aggregator.InitOperatorRequest{
		ChainId:                    2,
		OperatorId:                 cl.operatorId[:],
		OperatorAddress:            cl.operatorAddr.Hex(),
		OperatorStateRetrieverAddr: cl.config.AVS.OperatorStateRetrieverAddr.Hex(),
		RegistryCoordinatorAddr:    cl.config.AVS.RegistryCoordinatorAddr.Hex(),
	}

	cl.logger.Info("Init operator to aggregator", "req", fmt.Sprintf("%#v", request))

	reply, err := n.InitOperator(nodeCtx, request)
	if err != nil {
		return fmt.Errorf("call initOperatorToAggregator failed: %v", err.Error())
	}

	if !reply.GetOk() {
		return fmt.Errorf("initOperatorToAggregator failed: %v", reply.GetReason())
	}

	return nil
}



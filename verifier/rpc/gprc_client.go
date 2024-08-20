package rpc

import (
	// "fmt"
	"context"
	"fmt"
	"net/rpc"
	"time"

	aggregator "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/agents"
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
		gRPCAggreagatorIpPortAddr:  config.GRpcAggregatorIpAddr,
		timeout:                    1 * time.Minute,
	}, nil
}

func (cl *AggregatorGRpcClient) InitOperatorToAggregator() error {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	fmt.Println(conn)

	if err != nil {
		return fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)

	nodeCtx, cancel := context.WithTimeout(context.Background(), cl.timeout)

	defer cancel()

	request := &aggregator.InitOperatorRequest{
		OperatorId:                    cl.operatorId[:],
		OperatorAddress:               cl.operatorAddr.Hex(),
		OperatorStateRetrieverAddress: cl.config.AVS.OperatorStateRetrieverAddr.Hex(),
		RegistryCoordinatorAddress:    cl.config.AVS.RegistryCoordinatorAddr.Hex(),
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

func (cl *AggregatorGRpcClient) CreateBalanceTaskToAggregator(req *agents.DeviationRequest) (*tasks.BalanceTaskInfo, error) {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)

	nodeCtx, cancel := context.WithTimeout(context.Background(), cl.timeout)

	defer cancel()

	request := &aggregator.CreateDeviationTaskRequest{
		ReqId:                req.ReqId.Bytes(),
		ReferenceBlockNumber: req.ReferenceBlockNumber,
		AgentAddress:         req.AgentAddress.Hex(),
	}

	cl.logger.Info("CreateAlertTask to aggregator", "req", fmt.Sprintf("%#v", request))

	reply, err := n.CreateDeviationTask(nodeCtx, request)
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed: %v", err.Error())
	}

	info, err := tasks.NewBalanceTaskInfo(reply.GetInfo())
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed by decode alert info: %v", err.Error())
	}

	return info, nil
}

func (cl *AggregatorGRpcClient) SendSignedTaskResponse(ctx context.Context, signedTaskResponse *tasks.SignedDeviationTaskRespRequest) (*aggregator.SignedTaskRespResponse, error) {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)

	nodeCtx, cancel := context.WithTimeout(ctx, cl.timeout)
	defer cancel()

	cl.logger.Debug("Sending Signed Task", "taskResponse", signedTaskResponse.Balance)

	request := &aggregator.SignedDeviationTaskRespRequest{
		OperatorId:               cl.operatorId[:],
		OperatorRequestSignature: signedTaskResponse.BlsSignature.Serialize(),
		BalanceHeader:            signedTaskResponse.Balance.ToPbType(),
	}

	response, err := n.ProcessSignedDeviationTaskResponse(nodeCtx, request)

	if err != nil {
		cl.logger.Error("Error while sending signed response to aggregator", "err", err)
		return nil, err
	}

	return response, nil
}

func (cl *AggregatorGRpcClient) CreateAgentSlashTaskToAggregator(req *agents.AgentSlashRequest) (*tasks.AgentSlashTaskInfo, error) {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)

	nodeCtx, cancel := context.WithTimeout(context.Background(), cl.timeout)

	defer cancel()

	request := &aggregator.CreateAgentSlashTaskRequest{
		// ReqId:                req.ReqId.Bytes(),
		ReferenceBlockNumber: req.ReferenceBlockNumber,
		AgentAddress:         req.AgentAddress.Hex(),
		TxnId:                req.TxnId.Bytes(),
	}

	cl.logger.Info("CreateAlertTask to aggregator", "req", req)

	reply, err := n.CreateAgentSlashTask(nodeCtx, request)
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed: %v", err.Error())
	}

	info, err := tasks.NewAgentSlashTaskInfo(reply.GetInfo())
	if err != nil {
		return nil, fmt.Errorf("call CreateAlertTask failed by decode alert info: %v", err.Error())
	}

	return info, nil
}

func (cl *AggregatorGRpcClient) SendSignedAgentSlashTaskResponse(signedTaskResponse *tasks.SignedAgentSlashTaskRespRequest) (*aggregator.SignedTaskRespResponse, error) {
	conn, err := grpc.NewClient(cl.gRPCAggreagatorIpPortAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("dial initOperatorToAggregator connection failed: %v", err.Error())
	}

	n := aggregator.NewAggregatorClient(conn)

	nodeCtx, cancel := context.WithTimeout(context.Background(), cl.timeout)
	defer cancel()

	cl.logger.Debug("Sending Signed Task", "taskResponse", signedTaskResponse.Balance)

	request := &aggregator.SignedAgentSlashTaskRespRequest{
		OperatorId:               cl.operatorId[:],
		OperatorRequestSignature: signedTaskResponse.BlsSignature.Serialize(),
		FraudHeader:              signedTaskResponse.Balance.ToPbType(),
	}

	response, err := n.ProcessSignedAgentSlashTaskResponse(nodeCtx, request)

	if err != nil {
		cl.logger.Error("Error while sending signed response to aggregator", "err", err)
		return nil, err
	}

	return response, nil
}

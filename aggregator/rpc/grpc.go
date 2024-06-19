package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Cedro-Finance/CedroAVS/core"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
)

type AggregatorRpcHandler interface {
	InitOperator(req *core.InitOperatorRequest) (*core.InitOperatorResponse, error)
	CreateTask(req *core.CreateTaskRequest) (*core.CreateTaskResponse, error)
	ProcessSignedTaskResponse(signedTaskResponse *core.SignedTaskRespRequest) (*core.SignedTaskRespResponse, error)
}

type GRpcHandler struct {
	pb.UnimplementedAggregatorServer
	logger     logging.Logger
	wg         *sync.WaitGroup
	aggreagtor AggregatorRpcHandler
}

func NewGRpcHandler(logger logging.Logger, aggreagator AggregatorRpcHandler) *GRpcHandler {
	return &GRpcHandler{
		logger:     logger,
		aggreagtor: aggreagator,
		wg:         &sync.WaitGroup{},
	}
}

func (g *GRpcHandler) Start(ctx context.Context, addr string) {

	if addr == "" {
		addr = ":50051"
	}

	g.logger.Infof("Start GRpc Server", "addr", addr)

	lis, err := net.Listen("tcp", addr)

	if err != nil {
		g.logger.Fatalf("Failed creating Grpc server", "err", err)
	}

	s := grpc.NewServer()

	pb.RegisterAggregatorServer(s, g)

	g.logger.Infof("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		g.logger.Fatalf("failed to serve: %v", err)
	}

	
}

// Send Init operator to aggregator from operator, will check if the config is matching
func (s *GRpcHandler) InitOperator(ctx context.Context, req *pb.InitOperatorRequest) (*pb.InitOperatorResponse, error) {

	reqs, err := core.NewInitOperatorRequest(req)

	if err != nil {
		return nil, err
	}

	resp, err := s.aggreagtor.InitOperator(reqs)
	if err != nil {
		return nil, fmt.Errorf("initOperator handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

// Create a alert task
func (s *GRpcHandler) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	msg, err := core.NewCreateTaskRequest(req)
	if err != nil {
		return nil, fmt.Errorf("createTask message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.CreateTask(msg)
	if err != nil {
		return nil, fmt.Errorf("createTask handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

// Send signed task for alert
func (s *GRpcHandler) ProcessSignedTaskResponse(ctx context.Context, req *pb.SignedTaskRespRequest) (*pb.SignedTaskRespResponse, error) {
	msg, err := core.NewSignedTaskRespRequest(req)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.ProcessSignedTaskResponse(msg)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

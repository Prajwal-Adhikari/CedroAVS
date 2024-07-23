package rpc

import (
	"context"
	"fmt"
	"net"
	"sync"

	pb "github.com/Cedro-Finance/CedroAVS/api/grpc"
	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Layr-Labs/eigensdk-go/logging"
	"google.golang.org/grpc"
)

type IAggregatorRpcHandler interface {
	InitOperator(req *tasks.InitOperatorRequest) (*tasks.InitOperatorResponse, error)

	CreateTask(task string, req tasks.ICreateTaskRequest) (tasks.ICreateTaskResponse, error)

	ProcessSignedTask(task string, signed *tasks.CreateSignedTaskResRequest) (*tasks.SignedTaskRespResponse, error)
}

type GRpcHandler struct {
	pb.UnimplementedAggregatorServer
	logger     logging.Logger
	wg         *sync.WaitGroup
	aggreagtor IAggregatorRpcHandler
}

func NewGRpcHandler(logger logging.Logger, aggreagator IAggregatorRpcHandler) *GRpcHandler {
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

	reqs, err := tasks.NewInitOperatorRequest(req)

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
func (s *GRpcHandler) CreateDeviationTask(ctx context.Context, req *pb.CreateDeviationTaskRequest) (*pb.CreateDeviationTaskResponse, error) {
	msg, err := tasks.NewCreateDeviationTaskRequest(req)
	if err != nil {
		return nil, fmt.Errorf("createTask message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.CreateTask(tasks.AgentDeviationTask, msg)
	if err != nil {
		return nil, fmt.Errorf("createTask handler error: %v", err.Error())
	}

	taskInfo := resp.Info.(*tasks.BalanceTaskInfo)

	response := &tasks.CreateDeviationTaskResponse{
		Info: *taskInfo,
	}

	return response.ToPbType(), nil
}

// Create a alert task
func (s *GRpcHandler) CreateAgentSlashTask(ctx context.Context, req *pb.CreateAgentSlashTaskRequest) (*pb.CreateAgentSlashTaskResponse, error) {
	msg, err := tasks.NewCreateAgentSlashTaskRequest(req)
	if err != nil {
		return nil, fmt.Errorf("createTask message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.CreateTask(tasks.AgentSlashTask, msg)
	if err != nil {
		return nil, fmt.Errorf("createTask handler error: %v", err.Error())
	}

	taskInfo := resp.Info.(*tasks.AgentSlashTaskInfo)

	response := &tasks.CreateAgentSlashTaskResponse{
		Info: *taskInfo,
	}

	return response.ToPbType(), nil
}

// Send signed task for alert
func (s *GRpcHandler) ProcessSignedDeviationTaskResponse(ctx context.Context, req *pb.SignedDeviationTaskRespRequest) (*pb.SignedTaskRespResponse, error) {
	msg, err := tasks.NewSignedDeviationTaskRespRequest(req)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.ProcessSignedTask(tasks.AgentDeviationTask, msg)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

func (s *GRpcHandler) ProcessSignedAgentSlashTaskResponse(ctx context.Context, req *pb.SignedAgentSlashTaskRespRequest) (*pb.SignedTaskRespResponse, error) {
	msg, err := tasks.NewSignedAgentSlashTaskRespRequest(req)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse message convert error: %v", err.Error())
	}

	resp, err := s.aggreagtor.ProcessSignedTask(tasks.AgentSlashTask, msg)
	if err != nil {
		return nil, fmt.Errorf("processSignedTaskResponse handler error: %v", err.Error())
	}

	return resp.ToPbType(), nil
}

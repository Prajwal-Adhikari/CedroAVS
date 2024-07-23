package tasks

import (
	"context"
	"math/big"
	"sync"

	"github.com/Cedro-Finance/CedroAVS/core/tasks"
	"github.com/Cedro-Finance/CedroAVS/verifier/client/avs"
	"github.com/Cedro-Finance/CedroAVS/verifier/opio"
	"github.com/Layr-Labs/eigensdk-go/logging"
)

type OperatorAgentSlashTask struct {
	req     *avs.AgentSlashRequest
	retries uint64
}

type SlashTaskProcessor struct {
	tasks map[*big.Int](*OperatorAgentSlashTask)

	taskMu sync.RWMutex

	logger logging.Logger

	signCh chan *tasks.SignedAgentSlashTaskRespRequest

	reqCh chan avs.TaskRequest

	opio *opio.OperatorIO
}

var _ Taskprocessor = (*SlashTaskProcessor)(nil)

func NewSlashTask(opio *opio.OperatorIO, logger logging.Logger) *SlashTaskProcessor {
	return &SlashTaskProcessor{
		tasks:  make(map[*big.Int]*OperatorAgentSlashTask),
		signCh: make(chan *tasks.SignedAgentSlashTaskRespRequest, 10),
		opio:   opio,
		logger: logger,
		reqCh:  make(chan avs.TaskRequest),
	}
}

func (r *SlashTaskProcessor) Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	r.logger.Info("Slash Task Processor Started")

	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.reqCh:
			r.logger.Info("Slash Task Request received", "req", req)

			reqData := req.(*avs.AgentSlashRequest)

			newTask, err := r.CreateSlashTask(reqData)

			if err != nil {
				r.logger.Error("Error while creating new task")
			}

			go r.processSlashAgentTask(newTask)

		case sign := <-r.signCh: 
			r.logger.Info("Signature received", "sign", sign)

			go r.sendSignedTaskToAggregator(sign)
		}
	}
}

func (r *SlashTaskProcessor) processSlashAgentTask(task *OperatorAgentSlashTask) {

	reqInfo, err := r.opio.GRpcClient.CreateAgentSlashTaskToAggregator(task.req)

	if err != nil {
		r.logger.Error("Error while creating new task in the aggregator", "err", err)
		return
	}

	reqInfo.InsufficientBalance = task.req.InsufficientBalance

	go r.signSlashResponse(reqInfo)

}

func (r *SlashTaskProcessor) GetSlashTask(txId *big.Int) *OperatorAgentSlashTask {
	r.taskMu.RLock()
	defer r.taskMu.RUnlock()

	res := r.tasks[txId]

	return res
}

func (r *SlashTaskProcessor) CreateSlashTask(req *avs.AgentSlashRequest) (*OperatorAgentSlashTask, error) {
	task := r.GetSlashTask(req.TxnId)

	if task == nil {
		r.taskMu.Lock()

		task = &OperatorAgentSlashTask{
			req:     req,
			retries: 0,
		}

		r.tasks[req.TxnId] = task
		r.taskMu.Unlock()
	}

	return task, nil
}

func (r *SlashTaskProcessor) GetRequestCh() chan avs.TaskRequest {
	return r.reqCh
}

func (r *SlashTaskProcessor) GetSignatureCh() chan *tasks.SignedAgentSlashTaskRespRequest {
	return r.signCh
}

func (s *SlashTaskProcessor) signSlashResponse(agentHeader *tasks.AgentSlashTaskInfo) {
	s.logger.Info("Sending Agent Slash", "header", agentHeader)

	hash, err := agentHeader.GetHeaderHash()

	if err != nil {
		s.logger.Error("Error while hashing the message")
	}

	signature := s.opio.Signer.SignMessage(hash)

	s.logger.Info("Signed the message", "sign", signature)

	s.signCh <- &tasks.SignedAgentSlashTaskRespRequest{
		Balance:      *agentHeader,
		BlsSignature: *signature,
	}
}

func (s *SlashTaskProcessor) sendSignedTaskToAggregator(req *tasks.SignedAgentSlashTaskRespRequest) {
	maxRetry := 10
	retries := 0

RETRY_LOOP:

	for {
		response, err := s.opio.GRpcClient.SendSignedAgentSlashTaskResponse(req)

		if err != nil {
			s.logger.Error("Error while snding Signed Task Response", "err", err)

			if retries < maxRetry {
				retries++
				break RETRY_LOOP
			}

			return
		} else {
			s.taskMu.Lock()
			delete(s.tasks, req.Balance.InsufficientBalance)
			s.taskMu.Unlock()

			s.logger.Info("Successfully posted the signatures to the aggregator",
				"task", req.Balance.TaskIndex,
				"agent", req.Balance.AgentAddress.Hex(),
				"response", response)

			return
		}
	}
}

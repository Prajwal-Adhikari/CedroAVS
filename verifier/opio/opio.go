package opio

import (
	"github.com/Cedro-Finance/CedroAVS/api/graphql"
	"github.com/Cedro-Finance/CedroAVS/verifier/rpc"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
)

type OperatorIO struct {
	TxMgr txmgr.TxManager

	Signer *bls.KeyPair

	GRpcClient *rpc.AggregatorGRpcClient

	GraphClient *graphql.GraphQLClient
}

func NewOpeartorIO(txMg txmgr.TxManager, signer *bls.KeyPair, gRpcClient *rpc.AggregatorGRpcClient, graphClient *graphql.GraphQLClient) *OperatorIO {
	return &OperatorIO{
		TxMgr:       txMg,
		Signer:      signer,
		GRpcClient:  gRpcClient,
		GraphClient: graphClient,
	}
}

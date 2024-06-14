package agents

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"golang.org/x/crypto/sha3"
)

type AggregateBalanceResponse struct {
	//ReqId is the identifier for each request from AVS to get the balance
	ReqId *big.Int

	//Total Balance accumulated from each agent contract for given agent
	ActualBalance *big.Int

	//Total Balance that a agent should have in given timestamp
	IdealBalance *big.Int

	//A blocknumber on which check event is emitted
	ReferenceBlockNumber uint64
}

func (a *AggregateBalanceResponse) EncodeSignHash() ([]byte, error) {

	BalanceType, err := abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{
			Name: "reqId",
			Type: "uint256",
		},
		{
			Name: "actualBalance",
			Type: "uint256",
		},
		{
			Name: "idealBalance",
			Type: "uint256",
		},

		{
			Name: "referenceBlockNumber",
			Type: "uint64",
		},
	})

	if err != nil {
		return nil, err
	}

	arguments := abi.Arguments{
		{
			Type: BalanceType,
		},
	}

	bytes, err := arguments.Pack(a)

	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (a AggregateBalanceResponse) SignHash() ([32]byte, error) {
	balanceBytes, err := a.EncodeSignHash()

	if err != nil {
		return [32]byte{}, err
	}

	var hash [32]byte

	hasher := sha3.NewLegacyKeccak256()

	hasher.Write(balanceBytes)

	copy(hash[:], hasher.Sum(nil)[:32])

	return hash, nil
}

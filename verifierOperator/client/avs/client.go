package avs

import (
	"fmt"

	ethCommon "github.com/ethereum/go-ethereum/common"
)

func Client() {
	fmt.Println("this is eigen client")
}

type AVSClient struct {
	contractAddress ethCommon.Address
}

func NewAVSClient() *AVSClient {
	return &AVSClient{contractAddress: ethCommon.Address{}}
}

type Request struct{}

// polls the balance request events from the AVS contract on eth
func (c *AVSClient) PollEvents() *Request { return &Request{} }

package aggregator

import (
	"testing"

	"github.com/Cedro-Finance/CedroAVS/core/config"
	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	"github.com/stretchr/testify/assert"
)



func TestMain(t *testing.T) {

	cfg := &config.Config{
		BlsPrivateKey:            &bls.KeyPair{},
		Production:               false,
		AVS:                      config.AVSConfig{},
		KeyStore:                 &config.KeyStore{},
		GRpcAggregatorIpAddr:     "",
		Agents:                   []config.ChainConfig{},
		AgentManagerAddr:         [20]byte{},
		EthConfig:                config.EthConfig{},
		AgentManagerStartHeight:  0,
		AgentManagerPollDuration: 0,
		GraphQLUrl:               "",
	}

	actualAggregator, _ := NewAggregator(cfg)

	// assert.Panics(t, func() {
	// 	NewAggregator(cfg)
	// })

	assert.NotNil(t, actualAggregator.grpcServerIpPortAddr)
}

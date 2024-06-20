package avs

import (
	cedroservicemanager "github.com/Cedro-Finance/CedroAVS/core/bindings/CedroServiceManager"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type AVSManagerBinding struct {
	ServiceManager      *cedroservicemanager.CedroServiceManager
	RegistryCoordinator *regcoord.ContractRegistryCoordinator
	ethClient           eth.Client
}

func NewAvsManagerBindings(registryAddr, operatorStateAddr common.Address, ethClient eth.Client) (*AVSManagerBinding, error) {
	contractRegistryCordinator, err := regcoord.NewContractRegistryCoordinator(registryAddr, ethClient)

	if err != nil {
		return nil, err
	}

	serviceManagerAddr, err := contractRegistryCordinator.ServiceManager(&bind.CallOpts{})

	if err != nil {
		return nil, err
	}

	contractServiceManager, err := cedroservicemanager.NewCedroServiceManager(serviceManagerAddr, ethClient)

	if err != nil {
		return nil, err
	}

	return &AVSManagerBinding{
		ServiceManager:      contractServiceManager,
		RegistryCoordinator: contractRegistryCordinator,
		ethClient:           ethClient,
	}, nil
}

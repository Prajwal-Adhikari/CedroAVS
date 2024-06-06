// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package abi

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// AgentManagerMetaData contains all meta data concerning the AgentManager contract.
var AgentManagerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"AIAgentAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"}],\"name\":\"NotAnAgent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"SlashingTimeLimitExceeded\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timestamp\",\"type\":\"uint256\"}],\"name\":\"DeviationCheck\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint8\",\"name\":\"version\",\"type\":\"uint8\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"adapter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"addressToBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"agentContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"agents\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"delistAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"chainIDs\",\"type\":\"bytes32[]\"}],\"name\":\"enrollAgent\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"invokeDeviationCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_agentContract\",\"type\":\"address\"}],\"name\":\"mapAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"pools\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"requestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"route\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_adapter\",\"type\":\"address\"}],\"name\":\"setAdapter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainID\",\"type\":\"bytes32\"}],\"name\":\"setChainID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_id\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"_enable\",\"type\":\"bool\"}],\"name\":\"setPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_route\",\"type\":\"address\"}],\"name\":\"setRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"slashAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slashingTimeThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AgentManagerABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentManagerMetaData.ABI instead.
var AgentManagerABI = AgentManagerMetaData.ABI

// AgentManager is an auto generated Go binding around an Ethereum contract.
type AgentManager struct {
	AgentManagerCaller     // Read-only binding to the contract
	AgentManagerTransactor // Write-only binding to the contract
	AgentManagerFilterer   // Log filterer for contract events
}

// AgentManagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentManagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentManagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentManagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentManagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentManagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentManagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentManagerSession struct {
	Contract     *AgentManager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentManagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentManagerCallerSession struct {
	Contract *AgentManagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// AgentManagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentManagerTransactorSession struct {
	Contract     *AgentManagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// AgentManagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentManagerRaw struct {
	Contract *AgentManager // Generic contract binding to access the raw methods on
}

// AgentManagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentManagerCallerRaw struct {
	Contract *AgentManagerCaller // Generic read-only contract binding to access the raw methods on
}

// AgentManagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentManagerTransactorRaw struct {
	Contract *AgentManagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentManager creates a new instance of AgentManager, bound to a specific deployed contract.
func NewAgentManager(address common.Address, backend bind.ContractBackend) (*AgentManager, error) {
	contract, err := bindAgentManager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentManager{AgentManagerCaller: AgentManagerCaller{contract: contract}, AgentManagerTransactor: AgentManagerTransactor{contract: contract}, AgentManagerFilterer: AgentManagerFilterer{contract: contract}}, nil
}

// NewAgentManagerCaller creates a new read-only instance of AgentManager, bound to a specific deployed contract.
func NewAgentManagerCaller(address common.Address, caller bind.ContractCaller) (*AgentManagerCaller, error) {
	contract, err := bindAgentManager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentManagerCaller{contract: contract}, nil
}

// NewAgentManagerTransactor creates a new write-only instance of AgentManager, bound to a specific deployed contract.
func NewAgentManagerTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentManagerTransactor, error) {
	contract, err := bindAgentManager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentManagerTransactor{contract: contract}, nil
}

// NewAgentManagerFilterer creates a new log filterer instance of AgentManager, bound to a specific deployed contract.
func NewAgentManagerFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentManagerFilterer, error) {
	contract, err := bindAgentManager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentManagerFilterer{contract: contract}, nil
}

// bindAgentManager binds a generic wrapper to an already deployed contract.
func bindAgentManager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentManagerMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentManager *AgentManagerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentManager.Contract.AgentManagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentManager *AgentManagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentManager.Contract.AgentManagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentManager *AgentManagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentManager.Contract.AgentManagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentManager *AgentManagerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentManager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentManager *AgentManagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentManager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentManager *AgentManagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentManager.Contract.contract.Transact(opts, method, params...)
}

// Adapter is a free data retrieval call binding the contract method 0x03eadcfc.
//
// Solidity: function adapter() view returns(address)
func (_AgentManager *AgentManagerCaller) Adapter(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "adapter")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Adapter is a free data retrieval call binding the contract method 0x03eadcfc.
//
// Solidity: function adapter() view returns(address)
func (_AgentManager *AgentManagerSession) Adapter() (common.Address, error) {
	return _AgentManager.Contract.Adapter(&_AgentManager.CallOpts)
}

// Adapter is a free data retrieval call binding the contract method 0x03eadcfc.
//
// Solidity: function adapter() view returns(address)
func (_AgentManager *AgentManagerCallerSession) Adapter() (common.Address, error) {
	return _AgentManager.Contract.Adapter(&_AgentManager.CallOpts)
}

// AddressToBytes32 is a free data retrieval call binding the contract method 0x82c947b7.
//
// Solidity: function addressToBytes32(address addr) pure returns(bytes32)
func (_AgentManager *AgentManagerCaller) AddressToBytes32(opts *bind.CallOpts, addr common.Address) ([32]byte, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "addressToBytes32", addr)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AddressToBytes32 is a free data retrieval call binding the contract method 0x82c947b7.
//
// Solidity: function addressToBytes32(address addr) pure returns(bytes32)
func (_AgentManager *AgentManagerSession) AddressToBytes32(addr common.Address) ([32]byte, error) {
	return _AgentManager.Contract.AddressToBytes32(&_AgentManager.CallOpts, addr)
}

// AddressToBytes32 is a free data retrieval call binding the contract method 0x82c947b7.
//
// Solidity: function addressToBytes32(address addr) pure returns(bytes32)
func (_AgentManager *AgentManagerCallerSession) AddressToBytes32(addr common.Address) ([32]byte, error) {
	return _AgentManager.Contract.AddressToBytes32(&_AgentManager.CallOpts, addr)
}

// AgentContract is a free data retrieval call binding the contract method 0x4760f2e0.
//
// Solidity: function agentContract(address ) view returns(address)
func (_AgentManager *AgentManagerCaller) AgentContract(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "agentContract", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentContract is a free data retrieval call binding the contract method 0x4760f2e0.
//
// Solidity: function agentContract(address ) view returns(address)
func (_AgentManager *AgentManagerSession) AgentContract(arg0 common.Address) (common.Address, error) {
	return _AgentManager.Contract.AgentContract(&_AgentManager.CallOpts, arg0)
}

// AgentContract is a free data retrieval call binding the contract method 0x4760f2e0.
//
// Solidity: function agentContract(address ) view returns(address)
func (_AgentManager *AgentManagerCallerSession) AgentContract(arg0 common.Address) (common.Address, error) {
	return _AgentManager.Contract.AgentContract(&_AgentManager.CallOpts, arg0)
}

// Agents is a free data retrieval call binding the contract method 0x513856c8.
//
// Solidity: function agents(uint256 ) view returns(address)
func (_AgentManager *AgentManagerCaller) Agents(opts *bind.CallOpts, arg0 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "agents", arg0)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Agents is a free data retrieval call binding the contract method 0x513856c8.
//
// Solidity: function agents(uint256 ) view returns(address)
func (_AgentManager *AgentManagerSession) Agents(arg0 *big.Int) (common.Address, error) {
	return _AgentManager.Contract.Agents(&_AgentManager.CallOpts, arg0)
}

// Agents is a free data retrieval call binding the contract method 0x513856c8.
//
// Solidity: function agents(uint256 ) view returns(address)
func (_AgentManager *AgentManagerCallerSession) Agents(arg0 *big.Int) (common.Address, error) {
	return _AgentManager.Contract.Agents(&_AgentManager.CallOpts, arg0)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentManager *AgentManagerCaller) ChainID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "chainID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentManager *AgentManagerSession) ChainID() ([32]byte, error) {
	return _AgentManager.Contract.ChainID(&_AgentManager.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentManager *AgentManagerCallerSession) ChainID() ([32]byte, error) {
	return _AgentManager.Contract.ChainID(&_AgentManager.CallOpts)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address ) view returns(bool)
func (_AgentManager *AgentManagerCaller) IsAgent(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "isAgent", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address ) view returns(bool)
func (_AgentManager *AgentManagerSession) IsAgent(arg0 common.Address) (bool, error) {
	return _AgentManager.Contract.IsAgent(&_AgentManager.CallOpts, arg0)
}

// IsAgent is a free data retrieval call binding the contract method 0x1ffbb064.
//
// Solidity: function isAgent(address ) view returns(bool)
func (_AgentManager *AgentManagerCallerSession) IsAgent(arg0 common.Address) (bool, error) {
	return _AgentManager.Contract.IsAgent(&_AgentManager.CallOpts, arg0)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentManager *AgentManagerCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentManager *AgentManagerSession) Oracle() (common.Address, error) {
	return _AgentManager.Contract.Oracle(&_AgentManager.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentManager *AgentManagerCallerSession) Oracle() (common.Address, error) {
	return _AgentManager.Contract.Oracle(&_AgentManager.CallOpts)
}

// Pools is a free data retrieval call binding the contract method 0xb5217bb4.
//
// Solidity: function pools(bytes32 ) view returns(bool enable, address token)
func (_AgentManager *AgentManagerCaller) Pools(opts *bind.CallOpts, arg0 [32]byte) (struct {
	Enable bool
	Token  common.Address
}, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "pools", arg0)

	outstruct := new(struct {
		Enable bool
		Token  common.Address
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Enable = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Token = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)

	return *outstruct, err

}

// Pools is a free data retrieval call binding the contract method 0xb5217bb4.
//
// Solidity: function pools(bytes32 ) view returns(bool enable, address token)
func (_AgentManager *AgentManagerSession) Pools(arg0 [32]byte) (struct {
	Enable bool
	Token  common.Address
}, error) {
	return _AgentManager.Contract.Pools(&_AgentManager.CallOpts, arg0)
}

// Pools is a free data retrieval call binding the contract method 0xb5217bb4.
//
// Solidity: function pools(bytes32 ) view returns(bool enable, address token)
func (_AgentManager *AgentManagerCallerSession) Pools(arg0 [32]byte) (struct {
	Enable bool
	Token  common.Address
}, error) {
	return _AgentManager.Contract.Pools(&_AgentManager.CallOpts, arg0)
}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID(address ) view returns(uint256)
func (_AgentManager *AgentManagerCaller) RequestID(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "requestID", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID(address ) view returns(uint256)
func (_AgentManager *AgentManagerSession) RequestID(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.RequestID(&_AgentManager.CallOpts, arg0)
}

// RequestID is a free data retrieval call binding the contract method 0x6ee1c916.
//
// Solidity: function requestID(address ) view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) RequestID(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.RequestID(&_AgentManager.CallOpts, arg0)
}

// RequestTimestamp is a free data retrieval call binding the contract method 0x3dfac681.
//
// Solidity: function requestTimestamp(address , uint256 ) view returns(uint256)
func (_AgentManager *AgentManagerCaller) RequestTimestamp(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "requestTimestamp", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RequestTimestamp is a free data retrieval call binding the contract method 0x3dfac681.
//
// Solidity: function requestTimestamp(address , uint256 ) view returns(uint256)
func (_AgentManager *AgentManagerSession) RequestTimestamp(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _AgentManager.Contract.RequestTimestamp(&_AgentManager.CallOpts, arg0, arg1)
}

// RequestTimestamp is a free data retrieval call binding the contract method 0x3dfac681.
//
// Solidity: function requestTimestamp(address , uint256 ) view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) RequestTimestamp(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _AgentManager.Contract.RequestTimestamp(&_AgentManager.CallOpts, arg0, arg1)
}

// Route is a free data retrieval call binding the contract method 0xd96661ef.
//
// Solidity: function route() view returns(address)
func (_AgentManager *AgentManagerCaller) Route(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "route")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Route is a free data retrieval call binding the contract method 0xd96661ef.
//
// Solidity: function route() view returns(address)
func (_AgentManager *AgentManagerSession) Route() (common.Address, error) {
	return _AgentManager.Contract.Route(&_AgentManager.CallOpts)
}

// Route is a free data retrieval call binding the contract method 0xd96661ef.
//
// Solidity: function route() view returns(address)
func (_AgentManager *AgentManagerCallerSession) Route() (common.Address, error) {
	return _AgentManager.Contract.Route(&_AgentManager.CallOpts)
}

// SlashingTimeThreshold is a free data retrieval call binding the contract method 0x104e93eb.
//
// Solidity: function slashingTimeThreshold() view returns(uint256)
func (_AgentManager *AgentManagerCaller) SlashingTimeThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "slashingTimeThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SlashingTimeThreshold is a free data retrieval call binding the contract method 0x104e93eb.
//
// Solidity: function slashingTimeThreshold() view returns(uint256)
func (_AgentManager *AgentManagerSession) SlashingTimeThreshold() (*big.Int, error) {
	return _AgentManager.Contract.SlashingTimeThreshold(&_AgentManager.CallOpts)
}

// SlashingTimeThreshold is a free data retrieval call binding the contract method 0x104e93eb.
//
// Solidity: function slashingTimeThreshold() view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) SlashingTimeThreshold() (*big.Int, error) {
	return _AgentManager.Contract.SlashingTimeThreshold(&_AgentManager.CallOpts)
}

// DelistAgent is a paid mutator transaction binding the contract method 0x8b24d351.
//
// Solidity: function delistAgent(address _agent) returns()
func (_AgentManager *AgentManagerTransactor) DelistAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "delistAgent", _agent)
}

// DelistAgent is a paid mutator transaction binding the contract method 0x8b24d351.
//
// Solidity: function delistAgent(address _agent) returns()
func (_AgentManager *AgentManagerSession) DelistAgent(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.DelistAgent(&_AgentManager.TransactOpts, _agent)
}

// DelistAgent is a paid mutator transaction binding the contract method 0x8b24d351.
//
// Solidity: function delistAgent(address _agent) returns()
func (_AgentManager *AgentManagerTransactorSession) DelistAgent(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.DelistAgent(&_AgentManager.TransactOpts, _agent)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x52cd29f5.
//
// Solidity: function enrollAgent(bytes32[] chainIDs) payable returns()
func (_AgentManager *AgentManagerTransactor) EnrollAgent(opts *bind.TransactOpts, chainIDs [][32]byte) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "enrollAgent", chainIDs)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x52cd29f5.
//
// Solidity: function enrollAgent(bytes32[] chainIDs) payable returns()
func (_AgentManager *AgentManagerSession) EnrollAgent(chainIDs [][32]byte) (*types.Transaction, error) {
	return _AgentManager.Contract.EnrollAgent(&_AgentManager.TransactOpts, chainIDs)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x52cd29f5.
//
// Solidity: function enrollAgent(bytes32[] chainIDs) payable returns()
func (_AgentManager *AgentManagerTransactorSession) EnrollAgent(chainIDs [][32]byte) (*types.Transaction, error) {
	return _AgentManager.Contract.EnrollAgent(&_AgentManager.TransactOpts, chainIDs)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AgentManager *AgentManagerTransactor) Initialize(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "initialize")
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AgentManager *AgentManagerSession) Initialize() (*types.Transaction, error) {
	return _AgentManager.Contract.Initialize(&_AgentManager.TransactOpts)
}

// Initialize is a paid mutator transaction binding the contract method 0x8129fc1c.
//
// Solidity: function initialize() returns()
func (_AgentManager *AgentManagerTransactorSession) Initialize() (*types.Transaction, error) {
	return _AgentManager.Contract.Initialize(&_AgentManager.TransactOpts)
}

// InvokeDeviationCheck is a paid mutator transaction binding the contract method 0xb6c7657a.
//
// Solidity: function invokeDeviationCheck(address _agent) returns()
func (_AgentManager *AgentManagerTransactor) InvokeDeviationCheck(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "invokeDeviationCheck", _agent)
}

// InvokeDeviationCheck is a paid mutator transaction binding the contract method 0xb6c7657a.
//
// Solidity: function invokeDeviationCheck(address _agent) returns()
func (_AgentManager *AgentManagerSession) InvokeDeviationCheck(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.InvokeDeviationCheck(&_AgentManager.TransactOpts, _agent)
}

// InvokeDeviationCheck is a paid mutator transaction binding the contract method 0xb6c7657a.
//
// Solidity: function invokeDeviationCheck(address _agent) returns()
func (_AgentManager *AgentManagerTransactorSession) InvokeDeviationCheck(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.InvokeDeviationCheck(&_AgentManager.TransactOpts, _agent)
}

// MapAgent is a paid mutator transaction binding the contract method 0x7b1269d0.
//
// Solidity: function mapAgent(address _agent, address _agentContract) returns()
func (_AgentManager *AgentManagerTransactor) MapAgent(opts *bind.TransactOpts, _agent common.Address, _agentContract common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "mapAgent", _agent, _agentContract)
}

// MapAgent is a paid mutator transaction binding the contract method 0x7b1269d0.
//
// Solidity: function mapAgent(address _agent, address _agentContract) returns()
func (_AgentManager *AgentManagerSession) MapAgent(_agent common.Address, _agentContract common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.MapAgent(&_AgentManager.TransactOpts, _agent, _agentContract)
}

// MapAgent is a paid mutator transaction binding the contract method 0x7b1269d0.
//
// Solidity: function mapAgent(address _agent, address _agentContract) returns()
func (_AgentManager *AgentManagerTransactorSession) MapAgent(_agent common.Address, _agentContract common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.MapAgent(&_AgentManager.TransactOpts, _agent, _agentContract)
}

// SetAdapter is a paid mutator transaction binding the contract method 0xab1da79c.
//
// Solidity: function setAdapter(address _adapter) returns()
func (_AgentManager *AgentManagerTransactor) SetAdapter(opts *bind.TransactOpts, _adapter common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setAdapter", _adapter)
}

// SetAdapter is a paid mutator transaction binding the contract method 0xab1da79c.
//
// Solidity: function setAdapter(address _adapter) returns()
func (_AgentManager *AgentManagerSession) SetAdapter(_adapter common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetAdapter(&_AgentManager.TransactOpts, _adapter)
}

// SetAdapter is a paid mutator transaction binding the contract method 0xab1da79c.
//
// Solidity: function setAdapter(address _adapter) returns()
func (_AgentManager *AgentManagerTransactorSession) SetAdapter(_adapter common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetAdapter(&_AgentManager.TransactOpts, _adapter)
}

// SetChainID is a paid mutator transaction binding the contract method 0x3af370b5.
//
// Solidity: function setChainID(bytes32 _chainID) returns()
func (_AgentManager *AgentManagerTransactor) SetChainID(opts *bind.TransactOpts, _chainID [32]byte) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setChainID", _chainID)
}

// SetChainID is a paid mutator transaction binding the contract method 0x3af370b5.
//
// Solidity: function setChainID(bytes32 _chainID) returns()
func (_AgentManager *AgentManagerSession) SetChainID(_chainID [32]byte) (*types.Transaction, error) {
	return _AgentManager.Contract.SetChainID(&_AgentManager.TransactOpts, _chainID)
}

// SetChainID is a paid mutator transaction binding the contract method 0x3af370b5.
//
// Solidity: function setChainID(bytes32 _chainID) returns()
func (_AgentManager *AgentManagerTransactorSession) SetChainID(_chainID [32]byte) (*types.Transaction, error) {
	return _AgentManager.Contract.SetChainID(&_AgentManager.TransactOpts, _chainID)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_AgentManager *AgentManagerTransactor) SetOracle(opts *bind.TransactOpts, _oracle common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setOracle", _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_AgentManager *AgentManagerSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetOracle(&_AgentManager.TransactOpts, _oracle)
}

// SetOracle is a paid mutator transaction binding the contract method 0x7adbf973.
//
// Solidity: function setOracle(address _oracle) returns()
func (_AgentManager *AgentManagerTransactorSession) SetOracle(_oracle common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetOracle(&_AgentManager.TransactOpts, _oracle)
}

// SetPool is a paid mutator transaction binding the contract method 0xf90ab8fd.
//
// Solidity: function setPool(bytes32 _id, address _token, bool _enable) returns()
func (_AgentManager *AgentManagerTransactor) SetPool(opts *bind.TransactOpts, _id [32]byte, _token common.Address, _enable bool) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setPool", _id, _token, _enable)
}

// SetPool is a paid mutator transaction binding the contract method 0xf90ab8fd.
//
// Solidity: function setPool(bytes32 _id, address _token, bool _enable) returns()
func (_AgentManager *AgentManagerSession) SetPool(_id [32]byte, _token common.Address, _enable bool) (*types.Transaction, error) {
	return _AgentManager.Contract.SetPool(&_AgentManager.TransactOpts, _id, _token, _enable)
}

// SetPool is a paid mutator transaction binding the contract method 0xf90ab8fd.
//
// Solidity: function setPool(bytes32 _id, address _token, bool _enable) returns()
func (_AgentManager *AgentManagerTransactorSession) SetPool(_id [32]byte, _token common.Address, _enable bool) (*types.Transaction, error) {
	return _AgentManager.Contract.SetPool(&_AgentManager.TransactOpts, _id, _token, _enable)
}

// SetRoute is a paid mutator transaction binding the contract method 0x48472de6.
//
// Solidity: function setRoute(address _route) returns()
func (_AgentManager *AgentManagerTransactor) SetRoute(opts *bind.TransactOpts, _route common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setRoute", _route)
}

// SetRoute is a paid mutator transaction binding the contract method 0x48472de6.
//
// Solidity: function setRoute(address _route) returns()
func (_AgentManager *AgentManagerSession) SetRoute(_route common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetRoute(&_AgentManager.TransactOpts, _route)
}

// SetRoute is a paid mutator transaction binding the contract method 0x48472de6.
//
// Solidity: function setRoute(address _route) returns()
func (_AgentManager *AgentManagerTransactorSession) SetRoute(_route common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SetRoute(&_AgentManager.TransactOpts, _route)
}

// SlashAgent is a paid mutator transaction binding the contract method 0xb5a91e01.
//
// Solidity: function slashAgent(address _agent) returns()
func (_AgentManager *AgentManagerTransactor) SlashAgent(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "slashAgent", _agent)
}

// SlashAgent is a paid mutator transaction binding the contract method 0xb5a91e01.
//
// Solidity: function slashAgent(address _agent) returns()
func (_AgentManager *AgentManagerSession) SlashAgent(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SlashAgent(&_AgentManager.TransactOpts, _agent)
}

// SlashAgent is a paid mutator transaction binding the contract method 0xb5a91e01.
//
// Solidity: function slashAgent(address _agent) returns()
func (_AgentManager *AgentManagerTransactorSession) SlashAgent(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.SlashAgent(&_AgentManager.TransactOpts, _agent)
}

// AgentManagerDeviationCheckIterator is returned from FilterDeviationCheck and is used to iterate over the raw logs and unpacked data for DeviationCheck events raised by the AgentManager contract.
type AgentManagerDeviationCheckIterator struct {
	Event *AgentManagerDeviationCheck // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentManagerDeviationCheckIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerDeviationCheck)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentManagerDeviationCheck)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentManagerDeviationCheckIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerDeviationCheckIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerDeviationCheck represents a DeviationCheck event raised by the AgentManager contract.
type AgentManagerDeviationCheck struct {
	Agent     common.Address
	RequestID *big.Int
	Timestamp *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeviationCheck is a free log retrieval operation binding the contract event 0xa08d0dfd9bd009969096b75b05bd74056e062358eac1eb43af41815db2bd1b77.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID, uint256 timestamp)
func (_AgentManager *AgentManagerFilterer) FilterDeviationCheck(opts *bind.FilterOpts) (*AgentManagerDeviationCheckIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "DeviationCheck")
	if err != nil {
		return nil, err
	}
	return &AgentManagerDeviationCheckIterator{contract: _AgentManager.contract, event: "DeviationCheck", logs: logs, sub: sub}, nil
}

// WatchDeviationCheck is a free log subscription operation binding the contract event 0xa08d0dfd9bd009969096b75b05bd74056e062358eac1eb43af41815db2bd1b77.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID, uint256 timestamp)
func (_AgentManager *AgentManagerFilterer) WatchDeviationCheck(opts *bind.WatchOpts, sink chan<- *AgentManagerDeviationCheck) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "DeviationCheck")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerDeviationCheck)
				if err := _AgentManager.contract.UnpackLog(event, "DeviationCheck", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeviationCheck is a log parse operation binding the contract event 0xa08d0dfd9bd009969096b75b05bd74056e062358eac1eb43af41815db2bd1b77.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID, uint256 timestamp)
func (_AgentManager *AgentManagerFilterer) ParseDeviationCheck(log types.Log) (*AgentManagerDeviationCheck, error) {
	event := new(AgentManagerDeviationCheck)
	if err := _AgentManager.contract.UnpackLog(event, "DeviationCheck", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerInitializedIterator is returned from FilterInitialized and is used to iterate over the raw logs and unpacked data for Initialized events raised by the AgentManager contract.
type AgentManagerInitializedIterator struct {
	Event *AgentManagerInitialized // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *AgentManagerInitializedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerInitialized)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(AgentManagerInitialized)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *AgentManagerInitializedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerInitializedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerInitialized represents a Initialized event raised by the AgentManager contract.
type AgentManagerInitialized struct {
	Version uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AgentManager *AgentManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*AgentManagerInitializedIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AgentManagerInitializedIterator{contract: _AgentManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AgentManager *AgentManagerFilterer) WatchInitialized(opts *bind.WatchOpts, sink chan<- *AgentManagerInitialized) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerInitialized)
				if err := _AgentManager.contract.UnpackLog(event, "Initialized", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseInitialized is a log parse operation binding the contract event 0x7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb3847402498.
//
// Solidity: event Initialized(uint8 version)
func (_AgentManager *AgentManagerFilterer) ParseInitialized(log types.Log) (*AgentManagerInitialized, error) {
	event := new(AgentManagerInitialized)
	if err := _AgentManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package agentmanager

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
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"}],\"name\":\"AIAgentAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"ArgumentsMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"bondAmount\",\"type\":\"uint256\"}],\"name\":\"BondAmountLessThanThreshold\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"deviationPercentage\",\"type\":\"uint256\"}],\"name\":\"DeviationAlreadySet\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InsufficientFees\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"InvalidInitialization\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"}],\"name\":\"NotAnAgent\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"NotInitializing\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"PreviousDeviationInProgress\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"ReentrancyGuardReentrantCall\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"SlashingTimeLimitExceeded\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"TransferEthFailed\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"chainIDs\",\"type\":\"bytes32[]\"}],\"name\":\"AgentPodChainList\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"chainID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agentContract\",\"type\":\"address\"}],\"name\":\"AgentPodCreated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"bondAmount\",\"type\":\"uint256\"}],\"name\":\"BondAmount\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txnId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32[]\",\"name\":\"assets\",\"type\":\"bytes32[]\"},{\"indexed\":false,\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"},{\"indexed\":false,\"internalType\":\"enumData.STATUS\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"Claim\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txnID\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"depositor\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"asset\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principalAmount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"timeInterval\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"requestID\",\"type\":\"uint256\"}],\"name\":\"DeviationCheck\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"version\",\"type\":\"uint64\"}],\"name\":\"Initialized\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountToSlash\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"chainID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"txnId\",\"type\":\"uint256\"}],\"name\":\"SlashAgent\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"acceptedAssetList\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"adapter\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_addr\",\"type\":\"address\"}],\"name\":\"addressToBytes32\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"agentContract\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"agents\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"bondAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"bondThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"delistAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"deviationPercentage\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_txnId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_depositor\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"_assets\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amounts\",\"type\":\"uint256[]\"},{\"internalType\":\"enumData.STATUS\",\"name\":\"_status\",\"type\":\"uint8\"}],\"name\":\"emitClaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_txnID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_depositor\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_asset\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_principalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timeInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interestRate\",\"type\":\"uint256\"}],\"name\":\"emitDeposit\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amountToSlash\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"_chainID\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_txnId\",\"type\":\"uint256\"}],\"name\":\"emitSlash\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_chainIDs\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_fees\",\"type\":\"uint256[]\"}],\"name\":\"enrollAgent\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"initialize\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"invokeDeviationCheck\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"isAgent\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_agentContract\",\"type\":\"address\"}],\"name\":\"mapAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"pools\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"enable\",\"type\":\"bool\"},{\"internalType\":\"address\",\"name\":\"token\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"}],\"name\":\"receivedResponse\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requestFulfilled\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"requestID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"requestTimestamp\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"route\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_adapter\",\"type\":\"address\"}],\"name\":\"setAdapter\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_bondThreshold\",\"type\":\"uint256\"}],\"name\":\"setBondThreshold\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_chainID\",\"type\":\"bytes32\"}],\"name\":\"setChainID\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_deviationPercentage\",\"type\":\"uint256\"}],\"name\":\"setDeviationPercentage\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_id\",\"type\":\"bytes32[]\"},{\"internalType\":\"address[]\",\"name\":\"_token\",\"type\":\"address[]\"},{\"internalType\":\"bool[]\",\"name\":\"_enable\",\"type\":\"bool[]\"}],\"name\":\"setPools\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_route\",\"type\":\"address\"}],\"name\":\"setRoute\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_slashAmount\",\"type\":\"uint256\"}],\"name\":\"slashAgent\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"slashingTimeThreshold\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"updateBondAmount\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"}]",
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

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentManager *AgentManagerCaller) AcceptedAssetList(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "acceptedAssetList", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentManager *AgentManagerSession) AcceptedAssetList(arg0 *big.Int) ([32]byte, error) {
	return _AgentManager.Contract.AcceptedAssetList(&_AgentManager.CallOpts, arg0)
}

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentManager *AgentManagerCallerSession) AcceptedAssetList(arg0 *big.Int) ([32]byte, error) {
	return _AgentManager.Contract.AcceptedAssetList(&_AgentManager.CallOpts, arg0)
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
// Solidity: function addressToBytes32(address _addr) pure returns(bytes32)
func (_AgentManager *AgentManagerCaller) AddressToBytes32(opts *bind.CallOpts, _addr common.Address) ([32]byte, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "addressToBytes32", _addr)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AddressToBytes32 is a free data retrieval call binding the contract method 0x82c947b7.
//
// Solidity: function addressToBytes32(address _addr) pure returns(bytes32)
func (_AgentManager *AgentManagerSession) AddressToBytes32(_addr common.Address) ([32]byte, error) {
	return _AgentManager.Contract.AddressToBytes32(&_AgentManager.CallOpts, _addr)
}

// AddressToBytes32 is a free data retrieval call binding the contract method 0x82c947b7.
//
// Solidity: function addressToBytes32(address _addr) pure returns(bytes32)
func (_AgentManager *AgentManagerCallerSession) AddressToBytes32(_addr common.Address) ([32]byte, error) {
	return _AgentManager.Contract.AddressToBytes32(&_AgentManager.CallOpts, _addr)
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

// BondAmount is a free data retrieval call binding the contract method 0xa862c058.
//
// Solidity: function bondAmount(address ) view returns(uint256)
func (_AgentManager *AgentManagerCaller) BondAmount(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "bondAmount", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BondAmount is a free data retrieval call binding the contract method 0xa862c058.
//
// Solidity: function bondAmount(address ) view returns(uint256)
func (_AgentManager *AgentManagerSession) BondAmount(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.BondAmount(&_AgentManager.CallOpts, arg0)
}

// BondAmount is a free data retrieval call binding the contract method 0xa862c058.
//
// Solidity: function bondAmount(address ) view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) BondAmount(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.BondAmount(&_AgentManager.CallOpts, arg0)
}

// BondThreshold is a free data retrieval call binding the contract method 0x63fd7f49.
//
// Solidity: function bondThreshold() view returns(uint256)
func (_AgentManager *AgentManagerCaller) BondThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "bondThreshold")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BondThreshold is a free data retrieval call binding the contract method 0x63fd7f49.
//
// Solidity: function bondThreshold() view returns(uint256)
func (_AgentManager *AgentManagerSession) BondThreshold() (*big.Int, error) {
	return _AgentManager.Contract.BondThreshold(&_AgentManager.CallOpts)
}

// BondThreshold is a free data retrieval call binding the contract method 0x63fd7f49.
//
// Solidity: function bondThreshold() view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) BondThreshold() (*big.Int, error) {
	return _AgentManager.Contract.BondThreshold(&_AgentManager.CallOpts)
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

// DeviationPercentage is a free data retrieval call binding the contract method 0x4877ba33.
//
// Solidity: function deviationPercentage(address ) view returns(uint256)
func (_AgentManager *AgentManagerCaller) DeviationPercentage(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "deviationPercentage", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DeviationPercentage is a free data retrieval call binding the contract method 0x4877ba33.
//
// Solidity: function deviationPercentage(address ) view returns(uint256)
func (_AgentManager *AgentManagerSession) DeviationPercentage(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.DeviationPercentage(&_AgentManager.CallOpts, arg0)
}

// DeviationPercentage is a free data retrieval call binding the contract method 0x4877ba33.
//
// Solidity: function deviationPercentage(address ) view returns(uint256)
func (_AgentManager *AgentManagerCallerSession) DeviationPercentage(arg0 common.Address) (*big.Int, error) {
	return _AgentManager.Contract.DeviationPercentage(&_AgentManager.CallOpts, arg0)
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

// RequestFulfilled is a free data retrieval call binding the contract method 0xec2089e1.
//
// Solidity: function requestFulfilled(address , uint256 ) view returns(bool)
func (_AgentManager *AgentManagerCaller) RequestFulfilled(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (bool, error) {
	var out []interface{}
	err := _AgentManager.contract.Call(opts, &out, "requestFulfilled", arg0, arg1)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RequestFulfilled is a free data retrieval call binding the contract method 0xec2089e1.
//
// Solidity: function requestFulfilled(address , uint256 ) view returns(bool)
func (_AgentManager *AgentManagerSession) RequestFulfilled(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _AgentManager.Contract.RequestFulfilled(&_AgentManager.CallOpts, arg0, arg1)
}

// RequestFulfilled is a free data retrieval call binding the contract method 0xec2089e1.
//
// Solidity: function requestFulfilled(address , uint256 ) view returns(bool)
func (_AgentManager *AgentManagerCallerSession) RequestFulfilled(arg0 common.Address, arg1 *big.Int) (bool, error) {
	return _AgentManager.Contract.RequestFulfilled(&_AgentManager.CallOpts, arg0, arg1)
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

// EmitClaim is a paid mutator transaction binding the contract method 0xcf10bff3.
//
// Solidity: function emitClaim(address _agent, uint256 _txnId, address _depositor, bytes32[] _assets, uint256[] _amounts, uint8 _status) returns()
func (_AgentManager *AgentManagerTransactor) EmitClaim(opts *bind.TransactOpts, _agent common.Address, _txnId *big.Int, _depositor common.Address, _assets [][32]byte, _amounts []*big.Int, _status uint8) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "emitClaim", _agent, _txnId, _depositor, _assets, _amounts, _status)
}

// EmitClaim is a paid mutator transaction binding the contract method 0xcf10bff3.
//
// Solidity: function emitClaim(address _agent, uint256 _txnId, address _depositor, bytes32[] _assets, uint256[] _amounts, uint8 _status) returns()
func (_AgentManager *AgentManagerSession) EmitClaim(_agent common.Address, _txnId *big.Int, _depositor common.Address, _assets [][32]byte, _amounts []*big.Int, _status uint8) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitClaim(&_AgentManager.TransactOpts, _agent, _txnId, _depositor, _assets, _amounts, _status)
}

// EmitClaim is a paid mutator transaction binding the contract method 0xcf10bff3.
//
// Solidity: function emitClaim(address _agent, uint256 _txnId, address _depositor, bytes32[] _assets, uint256[] _amounts, uint8 _status) returns()
func (_AgentManager *AgentManagerTransactorSession) EmitClaim(_agent common.Address, _txnId *big.Int, _depositor common.Address, _assets [][32]byte, _amounts []*big.Int, _status uint8) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitClaim(&_AgentManager.TransactOpts, _agent, _txnId, _depositor, _assets, _amounts, _status)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x44f11d9e.
//
// Solidity: function emitDeposit(address _agent, uint256 _txnID, address _depositor, bytes32 _asset, uint256 _amount, uint256 _principalAmount, uint256 _timeInterval, uint256 _interestRate) returns()
func (_AgentManager *AgentManagerTransactor) EmitDeposit(opts *bind.TransactOpts, _agent common.Address, _txnID *big.Int, _depositor common.Address, _asset [32]byte, _amount *big.Int, _principalAmount *big.Int, _timeInterval *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "emitDeposit", _agent, _txnID, _depositor, _asset, _amount, _principalAmount, _timeInterval, _interestRate)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x44f11d9e.
//
// Solidity: function emitDeposit(address _agent, uint256 _txnID, address _depositor, bytes32 _asset, uint256 _amount, uint256 _principalAmount, uint256 _timeInterval, uint256 _interestRate) returns()
func (_AgentManager *AgentManagerSession) EmitDeposit(_agent common.Address, _txnID *big.Int, _depositor common.Address, _asset [32]byte, _amount *big.Int, _principalAmount *big.Int, _timeInterval *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitDeposit(&_AgentManager.TransactOpts, _agent, _txnID, _depositor, _asset, _amount, _principalAmount, _timeInterval, _interestRate)
}

// EmitDeposit is a paid mutator transaction binding the contract method 0x44f11d9e.
//
// Solidity: function emitDeposit(address _agent, uint256 _txnID, address _depositor, bytes32 _asset, uint256 _amount, uint256 _principalAmount, uint256 _timeInterval, uint256 _interestRate) returns()
func (_AgentManager *AgentManagerTransactorSession) EmitDeposit(_agent common.Address, _txnID *big.Int, _depositor common.Address, _asset [32]byte, _amount *big.Int, _principalAmount *big.Int, _timeInterval *big.Int, _interestRate *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitDeposit(&_AgentManager.TransactOpts, _agent, _txnID, _depositor, _asset, _amount, _principalAmount, _timeInterval, _interestRate)
}

// EmitSlash is a paid mutator transaction binding the contract method 0x689b94be.
//
// Solidity: function emitSlash(address _agent, uint256 _amountToSlash, bytes32 _chainID, uint256 _txnId) returns()
func (_AgentManager *AgentManagerTransactor) EmitSlash(opts *bind.TransactOpts, _agent common.Address, _amountToSlash *big.Int, _chainID [32]byte, _txnId *big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "emitSlash", _agent, _amountToSlash, _chainID, _txnId)
}

// EmitSlash is a paid mutator transaction binding the contract method 0x689b94be.
//
// Solidity: function emitSlash(address _agent, uint256 _amountToSlash, bytes32 _chainID, uint256 _txnId) returns()
func (_AgentManager *AgentManagerSession) EmitSlash(_agent common.Address, _amountToSlash *big.Int, _chainID [32]byte, _txnId *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitSlash(&_AgentManager.TransactOpts, _agent, _amountToSlash, _chainID, _txnId)
}

// EmitSlash is a paid mutator transaction binding the contract method 0x689b94be.
//
// Solidity: function emitSlash(address _agent, uint256 _amountToSlash, bytes32 _chainID, uint256 _txnId) returns()
func (_AgentManager *AgentManagerTransactorSession) EmitSlash(_agent common.Address, _amountToSlash *big.Int, _chainID [32]byte, _txnId *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EmitSlash(&_AgentManager.TransactOpts, _agent, _amountToSlash, _chainID, _txnId)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x5aa0b067.
//
// Solidity: function enrollAgent(bytes32[] _chainIDs, uint256[] _fees) payable returns()
func (_AgentManager *AgentManagerTransactor) EnrollAgent(opts *bind.TransactOpts, _chainIDs [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "enrollAgent", _chainIDs, _fees)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x5aa0b067.
//
// Solidity: function enrollAgent(bytes32[] _chainIDs, uint256[] _fees) payable returns()
func (_AgentManager *AgentManagerSession) EnrollAgent(_chainIDs [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EnrollAgent(&_AgentManager.TransactOpts, _chainIDs, _fees)
}

// EnrollAgent is a paid mutator transaction binding the contract method 0x5aa0b067.
//
// Solidity: function enrollAgent(bytes32[] _chainIDs, uint256[] _fees) payable returns()
func (_AgentManager *AgentManagerTransactorSession) EnrollAgent(_chainIDs [][32]byte, _fees []*big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.EnrollAgent(&_AgentManager.TransactOpts, _chainIDs, _fees)
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

// ReceivedResponse is a paid mutator transaction binding the contract method 0x6dc408a1.
//
// Solidity: function receivedResponse(address _agent) returns()
func (_AgentManager *AgentManagerTransactor) ReceivedResponse(opts *bind.TransactOpts, _agent common.Address) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "receivedResponse", _agent)
}

// ReceivedResponse is a paid mutator transaction binding the contract method 0x6dc408a1.
//
// Solidity: function receivedResponse(address _agent) returns()
func (_AgentManager *AgentManagerSession) ReceivedResponse(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.ReceivedResponse(&_AgentManager.TransactOpts, _agent)
}

// ReceivedResponse is a paid mutator transaction binding the contract method 0x6dc408a1.
//
// Solidity: function receivedResponse(address _agent) returns()
func (_AgentManager *AgentManagerTransactorSession) ReceivedResponse(_agent common.Address) (*types.Transaction, error) {
	return _AgentManager.Contract.ReceivedResponse(&_AgentManager.TransactOpts, _agent)
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

// SetBondThreshold is a paid mutator transaction binding the contract method 0xe1d8a8b9.
//
// Solidity: function setBondThreshold(uint256 _bondThreshold) returns()
func (_AgentManager *AgentManagerTransactor) SetBondThreshold(opts *bind.TransactOpts, _bondThreshold *big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setBondThreshold", _bondThreshold)
}

// SetBondThreshold is a paid mutator transaction binding the contract method 0xe1d8a8b9.
//
// Solidity: function setBondThreshold(uint256 _bondThreshold) returns()
func (_AgentManager *AgentManagerSession) SetBondThreshold(_bondThreshold *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SetBondThreshold(&_AgentManager.TransactOpts, _bondThreshold)
}

// SetBondThreshold is a paid mutator transaction binding the contract method 0xe1d8a8b9.
//
// Solidity: function setBondThreshold(uint256 _bondThreshold) returns()
func (_AgentManager *AgentManagerTransactorSession) SetBondThreshold(_bondThreshold *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SetBondThreshold(&_AgentManager.TransactOpts, _bondThreshold)
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

// SetDeviationPercentage is a paid mutator transaction binding the contract method 0xc50ec79a.
//
// Solidity: function setDeviationPercentage(uint256 _deviationPercentage) returns()
func (_AgentManager *AgentManagerTransactor) SetDeviationPercentage(opts *bind.TransactOpts, _deviationPercentage *big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setDeviationPercentage", _deviationPercentage)
}

// SetDeviationPercentage is a paid mutator transaction binding the contract method 0xc50ec79a.
//
// Solidity: function setDeviationPercentage(uint256 _deviationPercentage) returns()
func (_AgentManager *AgentManagerSession) SetDeviationPercentage(_deviationPercentage *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SetDeviationPercentage(&_AgentManager.TransactOpts, _deviationPercentage)
}

// SetDeviationPercentage is a paid mutator transaction binding the contract method 0xc50ec79a.
//
// Solidity: function setDeviationPercentage(uint256 _deviationPercentage) returns()
func (_AgentManager *AgentManagerTransactorSession) SetDeviationPercentage(_deviationPercentage *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SetDeviationPercentage(&_AgentManager.TransactOpts, _deviationPercentage)
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

// SetPools is a paid mutator transaction binding the contract method 0x11d7d912.
//
// Solidity: function setPools(bytes32[] _id, address[] _token, bool[] _enable) returns()
func (_AgentManager *AgentManagerTransactor) SetPools(opts *bind.TransactOpts, _id [][32]byte, _token []common.Address, _enable []bool) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "setPools", _id, _token, _enable)
}

// SetPools is a paid mutator transaction binding the contract method 0x11d7d912.
//
// Solidity: function setPools(bytes32[] _id, address[] _token, bool[] _enable) returns()
func (_AgentManager *AgentManagerSession) SetPools(_id [][32]byte, _token []common.Address, _enable []bool) (*types.Transaction, error) {
	return _AgentManager.Contract.SetPools(&_AgentManager.TransactOpts, _id, _token, _enable)
}

// SetPools is a paid mutator transaction binding the contract method 0x11d7d912.
//
// Solidity: function setPools(bytes32[] _id, address[] _token, bool[] _enable) returns()
func (_AgentManager *AgentManagerTransactorSession) SetPools(_id [][32]byte, _token []common.Address, _enable []bool) (*types.Transaction, error) {
	return _AgentManager.Contract.SetPools(&_AgentManager.TransactOpts, _id, _token, _enable)
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

// SlashAgent is a paid mutator transaction binding the contract method 0x690aa3e6.
//
// Solidity: function slashAgent(address _agent, uint256 _slashAmount) returns()
func (_AgentManager *AgentManagerTransactor) SlashAgent(opts *bind.TransactOpts, _agent common.Address, _slashAmount *big.Int) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "slashAgent", _agent, _slashAmount)
}

// SlashAgent is a paid mutator transaction binding the contract method 0x690aa3e6.
//
// Solidity: function slashAgent(address _agent, uint256 _slashAmount) returns()
func (_AgentManager *AgentManagerSession) SlashAgent(_agent common.Address, _slashAmount *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SlashAgent(&_AgentManager.TransactOpts, _agent, _slashAmount)
}

// SlashAgent is a paid mutator transaction binding the contract method 0x690aa3e6.
//
// Solidity: function slashAgent(address _agent, uint256 _slashAmount) returns()
func (_AgentManager *AgentManagerTransactorSession) SlashAgent(_agent common.Address, _slashAmount *big.Int) (*types.Transaction, error) {
	return _AgentManager.Contract.SlashAgent(&_AgentManager.TransactOpts, _agent, _slashAmount)
}

// UpdateBondAmount is a paid mutator transaction binding the contract method 0x03d53713.
//
// Solidity: function updateBondAmount() payable returns()
func (_AgentManager *AgentManagerTransactor) UpdateBondAmount(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentManager.contract.Transact(opts, "updateBondAmount")
}

// UpdateBondAmount is a paid mutator transaction binding the contract method 0x03d53713.
//
// Solidity: function updateBondAmount() payable returns()
func (_AgentManager *AgentManagerSession) UpdateBondAmount() (*types.Transaction, error) {
	return _AgentManager.Contract.UpdateBondAmount(&_AgentManager.TransactOpts)
}

// UpdateBondAmount is a paid mutator transaction binding the contract method 0x03d53713.
//
// Solidity: function updateBondAmount() payable returns()
func (_AgentManager *AgentManagerTransactorSession) UpdateBondAmount() (*types.Transaction, error) {
	return _AgentManager.Contract.UpdateBondAmount(&_AgentManager.TransactOpts)
}

// AgentManagerAgentPodChainListIterator is returned from FilterAgentPodChainList and is used to iterate over the raw logs and unpacked data for AgentPodChainList events raised by the AgentManager contract.
type AgentManagerAgentPodChainListIterator struct {
	Event *AgentManagerAgentPodChainList // Event containing the contract specifics and raw log

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
func (it *AgentManagerAgentPodChainListIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerAgentPodChainList)
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
		it.Event = new(AgentManagerAgentPodChainList)
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
func (it *AgentManagerAgentPodChainListIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerAgentPodChainListIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerAgentPodChainList represents a AgentPodChainList event raised by the AgentManager contract.
type AgentManagerAgentPodChainList struct {
	Agent    common.Address
	ChainIDs [][32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterAgentPodChainList is a free log retrieval operation binding the contract event 0x2f10946b561f38d0050eb99e7eaa47b08ff6387d5ee19942bcfc754a7b530906.
//
// Solidity: event AgentPodChainList(address agent, bytes32[] chainIDs)
func (_AgentManager *AgentManagerFilterer) FilterAgentPodChainList(opts *bind.FilterOpts) (*AgentManagerAgentPodChainListIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "AgentPodChainList")
	if err != nil {
		return nil, err
	}
	return &AgentManagerAgentPodChainListIterator{contract: _AgentManager.contract, event: "AgentPodChainList", logs: logs, sub: sub}, nil
}

// WatchAgentPodChainList is a free log subscription operation binding the contract event 0x2f10946b561f38d0050eb99e7eaa47b08ff6387d5ee19942bcfc754a7b530906.
//
// Solidity: event AgentPodChainList(address agent, bytes32[] chainIDs)
func (_AgentManager *AgentManagerFilterer) WatchAgentPodChainList(opts *bind.WatchOpts, sink chan<- *AgentManagerAgentPodChainList) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "AgentPodChainList")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerAgentPodChainList)
				if err := _AgentManager.contract.UnpackLog(event, "AgentPodChainList", log); err != nil {
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

// ParseAgentPodChainList is a log parse operation binding the contract event 0x2f10946b561f38d0050eb99e7eaa47b08ff6387d5ee19942bcfc754a7b530906.
//
// Solidity: event AgentPodChainList(address agent, bytes32[] chainIDs)
func (_AgentManager *AgentManagerFilterer) ParseAgentPodChainList(log types.Log) (*AgentManagerAgentPodChainList, error) {
	event := new(AgentManagerAgentPodChainList)
	if err := _AgentManager.contract.UnpackLog(event, "AgentPodChainList", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerAgentPodCreatedIterator is returned from FilterAgentPodCreated and is used to iterate over the raw logs and unpacked data for AgentPodCreated events raised by the AgentManager contract.
type AgentManagerAgentPodCreatedIterator struct {
	Event *AgentManagerAgentPodCreated // Event containing the contract specifics and raw log

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
func (it *AgentManagerAgentPodCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerAgentPodCreated)
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
		it.Event = new(AgentManagerAgentPodCreated)
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
func (it *AgentManagerAgentPodCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerAgentPodCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerAgentPodCreated represents a AgentPodCreated event raised by the AgentManager contract.
type AgentManagerAgentPodCreated struct {
	ChainID       [32]byte
	Agent         common.Address
	AgentContract common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterAgentPodCreated is a free log retrieval operation binding the contract event 0xe22697a2849d88697d0ddd171106dc1f6b9698e1f4e370bf51b21fda24009ebd.
//
// Solidity: event AgentPodCreated(bytes32 chainID, address agent, address agentContract)
func (_AgentManager *AgentManagerFilterer) FilterAgentPodCreated(opts *bind.FilterOpts) (*AgentManagerAgentPodCreatedIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "AgentPodCreated")
	if err != nil {
		return nil, err
	}
	return &AgentManagerAgentPodCreatedIterator{contract: _AgentManager.contract, event: "AgentPodCreated", logs: logs, sub: sub}, nil
}

// WatchAgentPodCreated is a free log subscription operation binding the contract event 0xe22697a2849d88697d0ddd171106dc1f6b9698e1f4e370bf51b21fda24009ebd.
//
// Solidity: event AgentPodCreated(bytes32 chainID, address agent, address agentContract)
func (_AgentManager *AgentManagerFilterer) WatchAgentPodCreated(opts *bind.WatchOpts, sink chan<- *AgentManagerAgentPodCreated) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "AgentPodCreated")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerAgentPodCreated)
				if err := _AgentManager.contract.UnpackLog(event, "AgentPodCreated", log); err != nil {
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

// ParseAgentPodCreated is a log parse operation binding the contract event 0xe22697a2849d88697d0ddd171106dc1f6b9698e1f4e370bf51b21fda24009ebd.
//
// Solidity: event AgentPodCreated(bytes32 chainID, address agent, address agentContract)
func (_AgentManager *AgentManagerFilterer) ParseAgentPodCreated(log types.Log) (*AgentManagerAgentPodCreated, error) {
	event := new(AgentManagerAgentPodCreated)
	if err := _AgentManager.contract.UnpackLog(event, "AgentPodCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerBondAmountIterator is returned from FilterBondAmount and is used to iterate over the raw logs and unpacked data for BondAmount events raised by the AgentManager contract.
type AgentManagerBondAmountIterator struct {
	Event *AgentManagerBondAmount // Event containing the contract specifics and raw log

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
func (it *AgentManagerBondAmountIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerBondAmount)
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
		it.Event = new(AgentManagerBondAmount)
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
func (it *AgentManagerBondAmountIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerBondAmountIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerBondAmount represents a BondAmount event raised by the AgentManager contract.
type AgentManagerBondAmount struct {
	Agent      common.Address
	BondAmount *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterBondAmount is a free log retrieval operation binding the contract event 0xe659df645661c396f8a34ce9f3c9763decd7ab9cdd01e745b2a4e60a3ba1413f.
//
// Solidity: event BondAmount(address agent, uint256 bondAmount)
func (_AgentManager *AgentManagerFilterer) FilterBondAmount(opts *bind.FilterOpts) (*AgentManagerBondAmountIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "BondAmount")
	if err != nil {
		return nil, err
	}
	return &AgentManagerBondAmountIterator{contract: _AgentManager.contract, event: "BondAmount", logs: logs, sub: sub}, nil
}

// WatchBondAmount is a free log subscription operation binding the contract event 0xe659df645661c396f8a34ce9f3c9763decd7ab9cdd01e745b2a4e60a3ba1413f.
//
// Solidity: event BondAmount(address agent, uint256 bondAmount)
func (_AgentManager *AgentManagerFilterer) WatchBondAmount(opts *bind.WatchOpts, sink chan<- *AgentManagerBondAmount) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "BondAmount")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerBondAmount)
				if err := _AgentManager.contract.UnpackLog(event, "BondAmount", log); err != nil {
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

// ParseBondAmount is a log parse operation binding the contract event 0xe659df645661c396f8a34ce9f3c9763decd7ab9cdd01e745b2a4e60a3ba1413f.
//
// Solidity: event BondAmount(address agent, uint256 bondAmount)
func (_AgentManager *AgentManagerFilterer) ParseBondAmount(log types.Log) (*AgentManagerBondAmount, error) {
	event := new(AgentManagerBondAmount)
	if err := _AgentManager.contract.UnpackLog(event, "BondAmount", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerClaimIterator is returned from FilterClaim and is used to iterate over the raw logs and unpacked data for Claim events raised by the AgentManager contract.
type AgentManagerClaimIterator struct {
	Event *AgentManagerClaim // Event containing the contract specifics and raw log

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
func (it *AgentManagerClaimIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerClaim)
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
		it.Event = new(AgentManagerClaim)
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
func (it *AgentManagerClaimIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerClaimIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerClaim represents a Claim event raised by the AgentManager contract.
type AgentManagerClaim struct {
	Agent     common.Address
	TxnId     *big.Int
	Depositor common.Address
	Assets    [][32]byte
	Amounts   []*big.Int
	Status    uint8
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterClaim is a free log retrieval operation binding the contract event 0x94b0c7d9c701cfb2c134cf1ab079f0909bab7d2c63bcae89962eeda04db30766.
//
// Solidity: event Claim(address agent, uint256 txnId, address depositor, bytes32[] assets, uint256[] amounts, uint8 status)
func (_AgentManager *AgentManagerFilterer) FilterClaim(opts *bind.FilterOpts) (*AgentManagerClaimIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "Claim")
	if err != nil {
		return nil, err
	}
	return &AgentManagerClaimIterator{contract: _AgentManager.contract, event: "Claim", logs: logs, sub: sub}, nil
}

// WatchClaim is a free log subscription operation binding the contract event 0x94b0c7d9c701cfb2c134cf1ab079f0909bab7d2c63bcae89962eeda04db30766.
//
// Solidity: event Claim(address agent, uint256 txnId, address depositor, bytes32[] assets, uint256[] amounts, uint8 status)
func (_AgentManager *AgentManagerFilterer) WatchClaim(opts *bind.WatchOpts, sink chan<- *AgentManagerClaim) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "Claim")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerClaim)
				if err := _AgentManager.contract.UnpackLog(event, "Claim", log); err != nil {
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

// ParseClaim is a log parse operation binding the contract event 0x94b0c7d9c701cfb2c134cf1ab079f0909bab7d2c63bcae89962eeda04db30766.
//
// Solidity: event Claim(address agent, uint256 txnId, address depositor, bytes32[] assets, uint256[] amounts, uint8 status)
func (_AgentManager *AgentManagerFilterer) ParseClaim(log types.Log) (*AgentManagerClaim, error) {
	event := new(AgentManagerClaim)
	if err := _AgentManager.contract.UnpackLog(event, "Claim", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the AgentManager contract.
type AgentManagerDepositIterator struct {
	Event *AgentManagerDeposit // Event containing the contract specifics and raw log

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
func (it *AgentManagerDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerDeposit)
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
		it.Event = new(AgentManagerDeposit)
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
func (it *AgentManagerDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerDeposit represents a Deposit event raised by the AgentManager contract.
type AgentManagerDeposit struct {
	Agent           common.Address
	TxnID           *big.Int
	Depositor       common.Address
	Asset           [32]byte
	Amount          *big.Int
	PrincipalAmount *big.Int
	TimeInterval    *big.Int
	InterestRate    *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xe068c2883f642d3a15746539679861ccbcfff0bf292b987ff8ef90ff9ce628e2.
//
// Solidity: event Deposit(address agent, uint256 txnID, address depositor, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate)
func (_AgentManager *AgentManagerFilterer) FilterDeposit(opts *bind.FilterOpts) (*AgentManagerDepositIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return &AgentManagerDepositIterator{contract: _AgentManager.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xe068c2883f642d3a15746539679861ccbcfff0bf292b987ff8ef90ff9ce628e2.
//
// Solidity: event Deposit(address agent, uint256 txnID, address depositor, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate)
func (_AgentManager *AgentManagerFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *AgentManagerDeposit) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "Deposit")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerDeposit)
				if err := _AgentManager.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// ParseDeposit is a log parse operation binding the contract event 0xe068c2883f642d3a15746539679861ccbcfff0bf292b987ff8ef90ff9ce628e2.
//
// Solidity: event Deposit(address agent, uint256 txnID, address depositor, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate)
func (_AgentManager *AgentManagerFilterer) ParseDeposit(log types.Log) (*AgentManagerDeposit, error) {
	event := new(AgentManagerDeposit)
	if err := _AgentManager.contract.UnpackLog(event, "Deposit", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterDeviationCheck is a free log retrieval operation binding the contract event 0xf1f6ae32e492ed94fefcd957830c206af14cae19c30c54be0814d27e8cdb5ac2.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID)
func (_AgentManager *AgentManagerFilterer) FilterDeviationCheck(opts *bind.FilterOpts) (*AgentManagerDeviationCheckIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "DeviationCheck")
	if err != nil {
		return nil, err
	}
	return &AgentManagerDeviationCheckIterator{contract: _AgentManager.contract, event: "DeviationCheck", logs: logs, sub: sub}, nil
}

// WatchDeviationCheck is a free log subscription operation binding the contract event 0xf1f6ae32e492ed94fefcd957830c206af14cae19c30c54be0814d27e8cdb5ac2.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID)
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

// ParseDeviationCheck is a log parse operation binding the contract event 0xf1f6ae32e492ed94fefcd957830c206af14cae19c30c54be0814d27e8cdb5ac2.
//
// Solidity: event DeviationCheck(address agent, uint256 requestID)
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
	Version uint64
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterInitialized is a free log retrieval operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgentManager *AgentManagerFilterer) FilterInitialized(opts *bind.FilterOpts) (*AgentManagerInitializedIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "Initialized")
	if err != nil {
		return nil, err
	}
	return &AgentManagerInitializedIterator{contract: _AgentManager.contract, event: "Initialized", logs: logs, sub: sub}, nil
}

// WatchInitialized is a free log subscription operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
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

// ParseInitialized is a log parse operation binding the contract event 0xc7f505b2f371ae2175ee4913f4499e1f2633a7b5936321eed1cdaeb6115181d2.
//
// Solidity: event Initialized(uint64 version)
func (_AgentManager *AgentManagerFilterer) ParseInitialized(log types.Log) (*AgentManagerInitialized, error) {
	event := new(AgentManagerInitialized)
	if err := _AgentManager.contract.UnpackLog(event, "Initialized", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// AgentManagerSlashAgentIterator is returned from FilterSlashAgent and is used to iterate over the raw logs and unpacked data for SlashAgent events raised by the AgentManager contract.
type AgentManagerSlashAgentIterator struct {
	Event *AgentManagerSlashAgent // Event containing the contract specifics and raw log

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
func (it *AgentManagerSlashAgentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentManagerSlashAgent)
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
		it.Event = new(AgentManagerSlashAgent)
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
func (it *AgentManagerSlashAgentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentManagerSlashAgentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentManagerSlashAgent represents a SlashAgent event raised by the AgentManager contract.
type AgentManagerSlashAgent struct {
	Agent         common.Address
	AmountToSlash *big.Int
	ChainID       [32]byte
	TxnId         *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSlashAgent is a free log retrieval operation binding the contract event 0xc01e2e0defd55e48b29c6b1c0abaa6e53cf62dd50601aadb27e0b2706d3da2ea.
//
// Solidity: event SlashAgent(address agent, uint256 amountToSlash, bytes32 chainID, uint256 txnId)
func (_AgentManager *AgentManagerFilterer) FilterSlashAgent(opts *bind.FilterOpts) (*AgentManagerSlashAgentIterator, error) {

	logs, sub, err := _AgentManager.contract.FilterLogs(opts, "SlashAgent")
	if err != nil {
		return nil, err
	}
	return &AgentManagerSlashAgentIterator{contract: _AgentManager.contract, event: "SlashAgent", logs: logs, sub: sub}, nil
}

// WatchSlashAgent is a free log subscription operation binding the contract event 0xc01e2e0defd55e48b29c6b1c0abaa6e53cf62dd50601aadb27e0b2706d3da2ea.
//
// Solidity: event SlashAgent(address agent, uint256 amountToSlash, bytes32 chainID, uint256 txnId)
func (_AgentManager *AgentManagerFilterer) WatchSlashAgent(opts *bind.WatchOpts, sink chan<- *AgentManagerSlashAgent) (event.Subscription, error) {

	logs, sub, err := _AgentManager.contract.WatchLogs(opts, "SlashAgent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentManagerSlashAgent)
				if err := _AgentManager.contract.UnpackLog(event, "SlashAgent", log); err != nil {
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

// ParseSlashAgent is a log parse operation binding the contract event 0xc01e2e0defd55e48b29c6b1c0abaa6e53cf62dd50601aadb27e0b2706d3da2ea.
//
// Solidity: event SlashAgent(address agent, uint256 amountToSlash, bytes32 chainID, uint256 txnId)
func (_AgentManager *AgentManagerFilterer) ParseSlashAgent(log types.Log) (*AgentManagerSlashAgent, error) {
	event := new(AgentManagerSlashAgent)
	if err := _AgentManager.contract.UnpackLog(event, "SlashAgent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

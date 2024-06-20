// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package agentpod

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

// AgentPodMetaData contains all meta data concerning the AgentPod contract.
var AgentPodMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_agentManager\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_oracle\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"ArgumentsMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_txnID\",\"type\":\"uint256\"}],\"name\":\"AssetAlreadyClaimed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"asset\",\"type\":\"bytes32\"}],\"name\":\"AssetAlreadyExists\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"asset\",\"type\":\"bytes32\"}],\"name\":\"AssetNotEnabledByAgent\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"asset\",\"type\":\"bytes32\"}],\"name\":\"AssetNotEnabledByAgentManager\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_txnID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"DepositNotMatured\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"}],\"name\":\"IncorrectRatio\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"sentId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"requiredId\",\"type\":\"uint256\"}],\"name\":\"IncorrectRequestId\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"contractTAG\",\"type\":\"string\"},{\"internalType\":\"uint256\",\"name\":\"valueSent\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"valueToBeGreater\",\"type\":\"uint256\"}],\"name\":\"InsufficientValue\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"SignatureAlreadyUsed\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"SignatureCorrupted\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"StaleMessage\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"address\",\"name\":\"_agent\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_amountToSlash\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"_chainID\",\"type\":\"bytes32\"}],\"name\":\"SlashAgent\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"ROOT_DECIMALS\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SECONDS_IN_A_YEAR\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"acceptedAssetList\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"agentManager\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_assets\",\"type\":\"bytes32[]\"}],\"name\":\"approveAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_txnId\",\"type\":\"uint256\"},{\"internalType\":\"bytes32[]\",\"name\":\"_assets\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_ratio\",\"type\":\"uint256[]\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_asset\",\"type\":\"bytes32\"}],\"name\":\"delistAsset\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"_asset\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timeInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_timeThreshold\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_signature\",\"type\":\"bytes\"}],\"name\":\"depositAsset\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"getAssetBalance\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"assetBalance\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getBalanceOfAgent\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_user\",\"type\":\"address\"}],\"name\":\"getUserTxnIDs\",\"outputs\":[{\"internalType\":\"uint256[]\",\"name\":\"\",\"type\":\"uint256[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"isAssetAccepted\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"operator\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oppurtunity\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"oracle\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"pool\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"poolBalanceOfAgentPod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountInUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"principalBalanceOfAgentPod\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amountInUSD\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"reward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"rewardAmount\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_oppurtunity\",\"type\":\"address\"}],\"name\":\"setOppurtunity\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"setOracle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_pool\",\"type\":\"address\"}],\"name\":\"setPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"signatureUsed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"name\":\"takeAction\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"totalReward\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32[]\",\"name\":\"_token\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"_amount\",\"type\":\"uint256[]\"}],\"name\":\"transferAssetToPool\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"txnIdCounter\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userDeposits\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"txnID\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"asset\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"principalAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"timeInterval\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"interestRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"createdAt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimedAt\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"claimed\",\"type\":\"bool\"},{\"internalType\":\"enumData.STATUS\",\"name\":\"status\",\"type\":\"uint8\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"userTxnIDs\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AgentPodABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentPodMetaData.ABI instead.
var AgentPodABI = AgentPodMetaData.ABI

// AgentPod is an auto generated Go binding around an Ethereum contract.
type AgentPod struct {
	AgentPodCaller     // Read-only binding to the contract
	AgentPodTransactor // Write-only binding to the contract
	AgentPodFilterer   // Log filterer for contract events
}

// AgentPodCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentPodCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentPodTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentPodTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentPodFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentPodFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentPodSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentPodSession struct {
	Contract     *AgentPod         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentPodCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentPodCallerSession struct {
	Contract *AgentPodCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// AgentPodTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentPodTransactorSession struct {
	Contract     *AgentPodTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// AgentPodRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentPodRaw struct {
	Contract *AgentPod // Generic contract binding to access the raw methods on
}

// AgentPodCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentPodCallerRaw struct {
	Contract *AgentPodCaller // Generic read-only contract binding to access the raw methods on
}

// AgentPodTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentPodTransactorRaw struct {
	Contract *AgentPodTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgentPod creates a new instance of AgentPod, bound to a specific deployed contract.
func NewAgentPod(address common.Address, backend bind.ContractBackend) (*AgentPod, error) {
	contract, err := bindAgentPod(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AgentPod{AgentPodCaller: AgentPodCaller{contract: contract}, AgentPodTransactor: AgentPodTransactor{contract: contract}, AgentPodFilterer: AgentPodFilterer{contract: contract}}, nil
}

// NewAgentPodCaller creates a new read-only instance of AgentPod, bound to a specific deployed contract.
func NewAgentPodCaller(address common.Address, caller bind.ContractCaller) (*AgentPodCaller, error) {
	contract, err := bindAgentPod(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentPodCaller{contract: contract}, nil
}

// NewAgentPodTransactor creates a new write-only instance of AgentPod, bound to a specific deployed contract.
func NewAgentPodTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentPodTransactor, error) {
	contract, err := bindAgentPod(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentPodTransactor{contract: contract}, nil
}

// NewAgentPodFilterer creates a new log filterer instance of AgentPod, bound to a specific deployed contract.
func NewAgentPodFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentPodFilterer, error) {
	contract, err := bindAgentPod(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentPodFilterer{contract: contract}, nil
}

// bindAgentPod binds a generic wrapper to an already deployed contract.
func bindAgentPod(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentPodMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentPod *AgentPodRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentPod.Contract.AgentPodCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentPod *AgentPodRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentPod.Contract.AgentPodTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentPod *AgentPodRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentPod.Contract.AgentPodTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AgentPod *AgentPodCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AgentPod.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AgentPod *AgentPodTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentPod.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AgentPod *AgentPodTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AgentPod.Contract.contract.Transact(opts, method, params...)
}

// ROOTDECIMALS is a free data retrieval call binding the contract method 0x09633be0.
//
// Solidity: function ROOT_DECIMALS() view returns(uint256)
func (_AgentPod *AgentPodCaller) ROOTDECIMALS(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "ROOT_DECIMALS")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ROOTDECIMALS is a free data retrieval call binding the contract method 0x09633be0.
//
// Solidity: function ROOT_DECIMALS() view returns(uint256)
func (_AgentPod *AgentPodSession) ROOTDECIMALS() (*big.Int, error) {
	return _AgentPod.Contract.ROOTDECIMALS(&_AgentPod.CallOpts)
}

// ROOTDECIMALS is a free data retrieval call binding the contract method 0x09633be0.
//
// Solidity: function ROOT_DECIMALS() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) ROOTDECIMALS() (*big.Int, error) {
	return _AgentPod.Contract.ROOTDECIMALS(&_AgentPod.CallOpts)
}

// SECONDSINAYEAR is a free data retrieval call binding the contract method 0xc4fa17a4.
//
// Solidity: function SECONDS_IN_A_YEAR() view returns(uint256)
func (_AgentPod *AgentPodCaller) SECONDSINAYEAR(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "SECONDS_IN_A_YEAR")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// SECONDSINAYEAR is a free data retrieval call binding the contract method 0xc4fa17a4.
//
// Solidity: function SECONDS_IN_A_YEAR() view returns(uint256)
func (_AgentPod *AgentPodSession) SECONDSINAYEAR() (*big.Int, error) {
	return _AgentPod.Contract.SECONDSINAYEAR(&_AgentPod.CallOpts)
}

// SECONDSINAYEAR is a free data retrieval call binding the contract method 0xc4fa17a4.
//
// Solidity: function SECONDS_IN_A_YEAR() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) SECONDSINAYEAR() (*big.Int, error) {
	return _AgentPod.Contract.SECONDSINAYEAR(&_AgentPod.CallOpts)
}

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentPod *AgentPodCaller) AcceptedAssetList(opts *bind.CallOpts, arg0 *big.Int) ([32]byte, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "acceptedAssetList", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentPod *AgentPodSession) AcceptedAssetList(arg0 *big.Int) ([32]byte, error) {
	return _AgentPod.Contract.AcceptedAssetList(&_AgentPod.CallOpts, arg0)
}

// AcceptedAssetList is a free data retrieval call binding the contract method 0x05afee7a.
//
// Solidity: function acceptedAssetList(uint256 ) view returns(bytes32)
func (_AgentPod *AgentPodCallerSession) AcceptedAssetList(arg0 *big.Int) ([32]byte, error) {
	return _AgentPod.Contract.AcceptedAssetList(&_AgentPod.CallOpts, arg0)
}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_AgentPod *AgentPodCaller) AgentManager(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "agentManager")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_AgentPod *AgentPodSession) AgentManager() (common.Address, error) {
	return _AgentPod.Contract.AgentManager(&_AgentPod.CallOpts)
}

// AgentManager is a free data retrieval call binding the contract method 0x7622f78d.
//
// Solidity: function agentManager() view returns(address)
func (_AgentPod *AgentPodCallerSession) AgentManager() (common.Address, error) {
	return _AgentPod.Contract.AgentManager(&_AgentPod.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentPod *AgentPodCaller) ChainID(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "chainID")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentPod *AgentPodSession) ChainID() ([32]byte, error) {
	return _AgentPod.Contract.ChainID(&_AgentPod.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() view returns(bytes32)
func (_AgentPod *AgentPodCallerSession) ChainID() ([32]byte, error) {
	return _AgentPod.Contract.ChainID(&_AgentPod.CallOpts)
}

// GetAssetBalance is a free data retrieval call binding the contract method 0x5373433f.
//
// Solidity: function getAssetBalance(address _token) view returns(uint256 assetBalance)
func (_AgentPod *AgentPodCaller) GetAssetBalance(opts *bind.CallOpts, _token common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "getAssetBalance", _token)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetAssetBalance is a free data retrieval call binding the contract method 0x5373433f.
//
// Solidity: function getAssetBalance(address _token) view returns(uint256 assetBalance)
func (_AgentPod *AgentPodSession) GetAssetBalance(_token common.Address) (*big.Int, error) {
	return _AgentPod.Contract.GetAssetBalance(&_AgentPod.CallOpts, _token)
}

// GetAssetBalance is a free data retrieval call binding the contract method 0x5373433f.
//
// Solidity: function getAssetBalance(address _token) view returns(uint256 assetBalance)
func (_AgentPod *AgentPodCallerSession) GetAssetBalance(_token common.Address) (*big.Int, error) {
	return _AgentPod.Contract.GetAssetBalance(&_AgentPod.CallOpts, _token)
}

// GetBalanceOfAgent is a free data retrieval call binding the contract method 0x459c71ea.
//
// Solidity: function getBalanceOfAgent() view returns(uint256)
func (_AgentPod *AgentPodCaller) GetBalanceOfAgent(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "getBalanceOfAgent")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetBalanceOfAgent is a free data retrieval call binding the contract method 0x459c71ea.
//
// Solidity: function getBalanceOfAgent() view returns(uint256)
func (_AgentPod *AgentPodSession) GetBalanceOfAgent() (*big.Int, error) {
	return _AgentPod.Contract.GetBalanceOfAgent(&_AgentPod.CallOpts)
}

// GetBalanceOfAgent is a free data retrieval call binding the contract method 0x459c71ea.
//
// Solidity: function getBalanceOfAgent() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) GetBalanceOfAgent() (*big.Int, error) {
	return _AgentPod.Contract.GetBalanceOfAgent(&_AgentPod.CallOpts)
}

// GetUserTxnIDs is a free data retrieval call binding the contract method 0x06a00bb7.
//
// Solidity: function getUserTxnIDs(address _user) view returns(uint256[])
func (_AgentPod *AgentPodCaller) GetUserTxnIDs(opts *bind.CallOpts, _user common.Address) ([]*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "getUserTxnIDs", _user)

	if err != nil {
		return *new([]*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new([]*big.Int)).(*[]*big.Int)

	return out0, err

}

// GetUserTxnIDs is a free data retrieval call binding the contract method 0x06a00bb7.
//
// Solidity: function getUserTxnIDs(address _user) view returns(uint256[])
func (_AgentPod *AgentPodSession) GetUserTxnIDs(_user common.Address) ([]*big.Int, error) {
	return _AgentPod.Contract.GetUserTxnIDs(&_AgentPod.CallOpts, _user)
}

// GetUserTxnIDs is a free data retrieval call binding the contract method 0x06a00bb7.
//
// Solidity: function getUserTxnIDs(address _user) view returns(uint256[])
func (_AgentPod *AgentPodCallerSession) GetUserTxnIDs(_user common.Address) ([]*big.Int, error) {
	return _AgentPod.Contract.GetUserTxnIDs(&_AgentPod.CallOpts, _user)
}

// IsAssetAccepted is a free data retrieval call binding the contract method 0x97b863eb.
//
// Solidity: function isAssetAccepted(bytes32 ) view returns(bool)
func (_AgentPod *AgentPodCaller) IsAssetAccepted(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "isAssetAccepted", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsAssetAccepted is a free data retrieval call binding the contract method 0x97b863eb.
//
// Solidity: function isAssetAccepted(bytes32 ) view returns(bool)
func (_AgentPod *AgentPodSession) IsAssetAccepted(arg0 [32]byte) (bool, error) {
	return _AgentPod.Contract.IsAssetAccepted(&_AgentPod.CallOpts, arg0)
}

// IsAssetAccepted is a free data retrieval call binding the contract method 0x97b863eb.
//
// Solidity: function isAssetAccepted(bytes32 ) view returns(bool)
func (_AgentPod *AgentPodCallerSession) IsAssetAccepted(arg0 [32]byte) (bool, error) {
	return _AgentPod.Contract.IsAssetAccepted(&_AgentPod.CallOpts, arg0)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_AgentPod *AgentPodCaller) Operator(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "operator")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_AgentPod *AgentPodSession) Operator() (common.Address, error) {
	return _AgentPod.Contract.Operator(&_AgentPod.CallOpts)
}

// Operator is a free data retrieval call binding the contract method 0x570ca735.
//
// Solidity: function operator() view returns(address)
func (_AgentPod *AgentPodCallerSession) Operator() (common.Address, error) {
	return _AgentPod.Contract.Operator(&_AgentPod.CallOpts)
}

// Oppurtunity is a free data retrieval call binding the contract method 0x4e8505ac.
//
// Solidity: function oppurtunity() view returns(address)
func (_AgentPod *AgentPodCaller) Oppurtunity(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "oppurtunity")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oppurtunity is a free data retrieval call binding the contract method 0x4e8505ac.
//
// Solidity: function oppurtunity() view returns(address)
func (_AgentPod *AgentPodSession) Oppurtunity() (common.Address, error) {
	return _AgentPod.Contract.Oppurtunity(&_AgentPod.CallOpts)
}

// Oppurtunity is a free data retrieval call binding the contract method 0x4e8505ac.
//
// Solidity: function oppurtunity() view returns(address)
func (_AgentPod *AgentPodCallerSession) Oppurtunity() (common.Address, error) {
	return _AgentPod.Contract.Oppurtunity(&_AgentPod.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentPod *AgentPodCaller) Oracle(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "oracle")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentPod *AgentPodSession) Oracle() (common.Address, error) {
	return _AgentPod.Contract.Oracle(&_AgentPod.CallOpts)
}

// Oracle is a free data retrieval call binding the contract method 0x7dc0d1d0.
//
// Solidity: function oracle() view returns(address)
func (_AgentPod *AgentPodCallerSession) Oracle() (common.Address, error) {
	return _AgentPod.Contract.Oracle(&_AgentPod.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentPod *AgentPodCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentPod *AgentPodSession) Owner() (common.Address, error) {
	return _AgentPod.Contract.Owner(&_AgentPod.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_AgentPod *AgentPodCallerSession) Owner() (common.Address, error) {
	return _AgentPod.Contract.Owner(&_AgentPod.CallOpts)
}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(address)
func (_AgentPod *AgentPodCaller) Pool(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "pool")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(address)
func (_AgentPod *AgentPodSession) Pool() (common.Address, error) {
	return _AgentPod.Contract.Pool(&_AgentPod.CallOpts)
}

// Pool is a free data retrieval call binding the contract method 0x16f0115b.
//
// Solidity: function pool() view returns(address)
func (_AgentPod *AgentPodCallerSession) Pool() (common.Address, error) {
	return _AgentPod.Contract.Pool(&_AgentPod.CallOpts)
}

// PoolBalanceOfAgentPod is a free data retrieval call binding the contract method 0xa6a43549.
//
// Solidity: function poolBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodCaller) PoolBalanceOfAgentPod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "poolBalanceOfAgentPod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PoolBalanceOfAgentPod is a free data retrieval call binding the contract method 0xa6a43549.
//
// Solidity: function poolBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodSession) PoolBalanceOfAgentPod() (*big.Int, error) {
	return _AgentPod.Contract.PoolBalanceOfAgentPod(&_AgentPod.CallOpts)
}

// PoolBalanceOfAgentPod is a free data retrieval call binding the contract method 0xa6a43549.
//
// Solidity: function poolBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodCallerSession) PoolBalanceOfAgentPod() (*big.Int, error) {
	return _AgentPod.Contract.PoolBalanceOfAgentPod(&_AgentPod.CallOpts)
}

// PrincipalBalanceOfAgentPod is a free data retrieval call binding the contract method 0xd9e529f2.
//
// Solidity: function principalBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodCaller) PrincipalBalanceOfAgentPod(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "principalBalanceOfAgentPod")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// PrincipalBalanceOfAgentPod is a free data retrieval call binding the contract method 0xd9e529f2.
//
// Solidity: function principalBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodSession) PrincipalBalanceOfAgentPod() (*big.Int, error) {
	return _AgentPod.Contract.PrincipalBalanceOfAgentPod(&_AgentPod.CallOpts)
}

// PrincipalBalanceOfAgentPod is a free data retrieval call binding the contract method 0xd9e529f2.
//
// Solidity: function principalBalanceOfAgentPod() view returns(uint256 amountInUSD)
func (_AgentPod *AgentPodCallerSession) PrincipalBalanceOfAgentPod() (*big.Int, error) {
	return _AgentPod.Contract.PrincipalBalanceOfAgentPod(&_AgentPod.CallOpts)
}

// Reward is a free data retrieval call binding the contract method 0x6353586b.
//
// Solidity: function reward(address ) view returns(uint256)
func (_AgentPod *AgentPodCaller) Reward(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "reward", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Reward is a free data retrieval call binding the contract method 0x6353586b.
//
// Solidity: function reward(address ) view returns(uint256)
func (_AgentPod *AgentPodSession) Reward(arg0 common.Address) (*big.Int, error) {
	return _AgentPod.Contract.Reward(&_AgentPod.CallOpts, arg0)
}

// Reward is a free data retrieval call binding the contract method 0x6353586b.
//
// Solidity: function reward(address ) view returns(uint256)
func (_AgentPod *AgentPodCallerSession) Reward(arg0 common.Address) (*big.Int, error) {
	return _AgentPod.Contract.Reward(&_AgentPod.CallOpts, arg0)
}

// RewardAmount is a free data retrieval call binding the contract method 0xf7b2a7be.
//
// Solidity: function rewardAmount() view returns(uint256)
func (_AgentPod *AgentPodCaller) RewardAmount(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "rewardAmount")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// RewardAmount is a free data retrieval call binding the contract method 0xf7b2a7be.
//
// Solidity: function rewardAmount() view returns(uint256)
func (_AgentPod *AgentPodSession) RewardAmount() (*big.Int, error) {
	return _AgentPod.Contract.RewardAmount(&_AgentPod.CallOpts)
}

// RewardAmount is a free data retrieval call binding the contract method 0xf7b2a7be.
//
// Solidity: function rewardAmount() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) RewardAmount() (*big.Int, error) {
	return _AgentPod.Contract.RewardAmount(&_AgentPod.CallOpts)
}

// SignatureUsed is a free data retrieval call binding the contract method 0xbb10c829.
//
// Solidity: function signatureUsed(bytes ) view returns(bool)
func (_AgentPod *AgentPodCaller) SignatureUsed(opts *bind.CallOpts, arg0 []byte) (bool, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "signatureUsed", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SignatureUsed is a free data retrieval call binding the contract method 0xbb10c829.
//
// Solidity: function signatureUsed(bytes ) view returns(bool)
func (_AgentPod *AgentPodSession) SignatureUsed(arg0 []byte) (bool, error) {
	return _AgentPod.Contract.SignatureUsed(&_AgentPod.CallOpts, arg0)
}

// SignatureUsed is a free data retrieval call binding the contract method 0xbb10c829.
//
// Solidity: function signatureUsed(bytes ) view returns(bool)
func (_AgentPod *AgentPodCallerSession) SignatureUsed(arg0 []byte) (bool, error) {
	return _AgentPod.Contract.SignatureUsed(&_AgentPod.CallOpts, arg0)
}

// TotalReward is a free data retrieval call binding the contract method 0x750142e6.
//
// Solidity: function totalReward() view returns(uint256)
func (_AgentPod *AgentPodCaller) TotalReward(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "totalReward")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalReward is a free data retrieval call binding the contract method 0x750142e6.
//
// Solidity: function totalReward() view returns(uint256)
func (_AgentPod *AgentPodSession) TotalReward() (*big.Int, error) {
	return _AgentPod.Contract.TotalReward(&_AgentPod.CallOpts)
}

// TotalReward is a free data retrieval call binding the contract method 0x750142e6.
//
// Solidity: function totalReward() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) TotalReward() (*big.Int, error) {
	return _AgentPod.Contract.TotalReward(&_AgentPod.CallOpts)
}

// TxnIdCounter is a free data retrieval call binding the contract method 0xe4231962.
//
// Solidity: function txnIdCounter() view returns(uint256)
func (_AgentPod *AgentPodCaller) TxnIdCounter(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "txnIdCounter")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TxnIdCounter is a free data retrieval call binding the contract method 0xe4231962.
//
// Solidity: function txnIdCounter() view returns(uint256)
func (_AgentPod *AgentPodSession) TxnIdCounter() (*big.Int, error) {
	return _AgentPod.Contract.TxnIdCounter(&_AgentPod.CallOpts)
}

// TxnIdCounter is a free data retrieval call binding the contract method 0xe4231962.
//
// Solidity: function txnIdCounter() view returns(uint256)
func (_AgentPod *AgentPodCallerSession) TxnIdCounter() (*big.Int, error) {
	return _AgentPod.Contract.TxnIdCounter(&_AgentPod.CallOpts)
}

// UserDeposits is a free data retrieval call binding the contract method 0xf16b41f4.
//
// Solidity: function userDeposits(uint256 ) view returns(uint256 txnID, address user, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate, uint256 createdAt, uint256 claimedAt, bool claimed, uint8 status)
func (_AgentPod *AgentPodCaller) UserDeposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	TxnID           *big.Int
	User            common.Address
	Asset           [32]byte
	Amount          *big.Int
	PrincipalAmount *big.Int
	TimeInterval    *big.Int
	InterestRate    *big.Int
	CreatedAt       *big.Int
	ClaimedAt       *big.Int
	Claimed         bool
	Status          uint8
}, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "userDeposits", arg0)

	outstruct := new(struct {
		TxnID           *big.Int
		User            common.Address
		Asset           [32]byte
		Amount          *big.Int
		PrincipalAmount *big.Int
		TimeInterval    *big.Int
		InterestRate    *big.Int
		CreatedAt       *big.Int
		ClaimedAt       *big.Int
		Claimed         bool
		Status          uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TxnID = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.User = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Asset = *abi.ConvertType(out[2], new([32]byte)).(*[32]byte)
	outstruct.Amount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.PrincipalAmount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.TimeInterval = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.InterestRate = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.CreatedAt = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.ClaimedAt = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.Claimed = *abi.ConvertType(out[9], new(bool)).(*bool)
	outstruct.Status = *abi.ConvertType(out[10], new(uint8)).(*uint8)

	return *outstruct, err

}

// UserDeposits is a free data retrieval call binding the contract method 0xf16b41f4.
//
// Solidity: function userDeposits(uint256 ) view returns(uint256 txnID, address user, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate, uint256 createdAt, uint256 claimedAt, bool claimed, uint8 status)
func (_AgentPod *AgentPodSession) UserDeposits(arg0 *big.Int) (struct {
	TxnID           *big.Int
	User            common.Address
	Asset           [32]byte
	Amount          *big.Int
	PrincipalAmount *big.Int
	TimeInterval    *big.Int
	InterestRate    *big.Int
	CreatedAt       *big.Int
	ClaimedAt       *big.Int
	Claimed         bool
	Status          uint8
}, error) {
	return _AgentPod.Contract.UserDeposits(&_AgentPod.CallOpts, arg0)
}

// UserDeposits is a free data retrieval call binding the contract method 0xf16b41f4.
//
// Solidity: function userDeposits(uint256 ) view returns(uint256 txnID, address user, bytes32 asset, uint256 amount, uint256 principalAmount, uint256 timeInterval, uint256 interestRate, uint256 createdAt, uint256 claimedAt, bool claimed, uint8 status)
func (_AgentPod *AgentPodCallerSession) UserDeposits(arg0 *big.Int) (struct {
	TxnID           *big.Int
	User            common.Address
	Asset           [32]byte
	Amount          *big.Int
	PrincipalAmount *big.Int
	TimeInterval    *big.Int
	InterestRate    *big.Int
	CreatedAt       *big.Int
	ClaimedAt       *big.Int
	Claimed         bool
	Status          uint8
}, error) {
	return _AgentPod.Contract.UserDeposits(&_AgentPod.CallOpts, arg0)
}

// UserTxnIDs is a free data retrieval call binding the contract method 0xcb21992e.
//
// Solidity: function userTxnIDs(address , uint256 ) view returns(uint256)
func (_AgentPod *AgentPodCaller) UserTxnIDs(opts *bind.CallOpts, arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _AgentPod.contract.Call(opts, &out, "userTxnIDs", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// UserTxnIDs is a free data retrieval call binding the contract method 0xcb21992e.
//
// Solidity: function userTxnIDs(address , uint256 ) view returns(uint256)
func (_AgentPod *AgentPodSession) UserTxnIDs(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _AgentPod.Contract.UserTxnIDs(&_AgentPod.CallOpts, arg0, arg1)
}

// UserTxnIDs is a free data retrieval call binding the contract method 0xcb21992e.
//
// Solidity: function userTxnIDs(address , uint256 ) view returns(uint256)
func (_AgentPod *AgentPodCallerSession) UserTxnIDs(arg0 common.Address, arg1 *big.Int) (*big.Int, error) {
	return _AgentPod.Contract.UserTxnIDs(&_AgentPod.CallOpts, arg0, arg1)
}

// ApproveAsset is a paid mutator transaction binding the contract method 0x76b47114.
//
// Solidity: function approveAsset(bytes32[] _assets) returns()
func (_AgentPod *AgentPodTransactor) ApproveAsset(opts *bind.TransactOpts, _assets [][32]byte) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "approveAsset", _assets)
}

// ApproveAsset is a paid mutator transaction binding the contract method 0x76b47114.
//
// Solidity: function approveAsset(bytes32[] _assets) returns()
func (_AgentPod *AgentPodSession) ApproveAsset(_assets [][32]byte) (*types.Transaction, error) {
	return _AgentPod.Contract.ApproveAsset(&_AgentPod.TransactOpts, _assets)
}

// ApproveAsset is a paid mutator transaction binding the contract method 0x76b47114.
//
// Solidity: function approveAsset(bytes32[] _assets) returns()
func (_AgentPod *AgentPodTransactorSession) ApproveAsset(_assets [][32]byte) (*types.Transaction, error) {
	return _AgentPod.Contract.ApproveAsset(&_AgentPod.TransactOpts, _assets)
}

// Claim is a paid mutator transaction binding the contract method 0x8b048ed2.
//
// Solidity: function claim(uint256 _txnId, bytes32[] _assets, uint256[] _ratio) returns()
func (_AgentPod *AgentPodTransactor) Claim(opts *bind.TransactOpts, _txnId *big.Int, _assets [][32]byte, _ratio []*big.Int) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "claim", _txnId, _assets, _ratio)
}

// Claim is a paid mutator transaction binding the contract method 0x8b048ed2.
//
// Solidity: function claim(uint256 _txnId, bytes32[] _assets, uint256[] _ratio) returns()
func (_AgentPod *AgentPodSession) Claim(_txnId *big.Int, _assets [][32]byte, _ratio []*big.Int) (*types.Transaction, error) {
	return _AgentPod.Contract.Claim(&_AgentPod.TransactOpts, _txnId, _assets, _ratio)
}

// Claim is a paid mutator transaction binding the contract method 0x8b048ed2.
//
// Solidity: function claim(uint256 _txnId, bytes32[] _assets, uint256[] _ratio) returns()
func (_AgentPod *AgentPodTransactorSession) Claim(_txnId *big.Int, _assets [][32]byte, _ratio []*big.Int) (*types.Transaction, error) {
	return _AgentPod.Contract.Claim(&_AgentPod.TransactOpts, _txnId, _assets, _ratio)
}

// DelistAsset is a paid mutator transaction binding the contract method 0xe06e60a5.
//
// Solidity: function delistAsset(bytes32 _asset) returns()
func (_AgentPod *AgentPodTransactor) DelistAsset(opts *bind.TransactOpts, _asset [32]byte) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "delistAsset", _asset)
}

// DelistAsset is a paid mutator transaction binding the contract method 0xe06e60a5.
//
// Solidity: function delistAsset(bytes32 _asset) returns()
func (_AgentPod *AgentPodSession) DelistAsset(_asset [32]byte) (*types.Transaction, error) {
	return _AgentPod.Contract.DelistAsset(&_AgentPod.TransactOpts, _asset)
}

// DelistAsset is a paid mutator transaction binding the contract method 0xe06e60a5.
//
// Solidity: function delistAsset(bytes32 _asset) returns()
func (_AgentPod *AgentPodTransactorSession) DelistAsset(_asset [32]byte) (*types.Transaction, error) {
	return _AgentPod.Contract.DelistAsset(&_AgentPod.TransactOpts, _asset)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xff556fbd.
//
// Solidity: function depositAsset(bytes32 _asset, uint256 _amount, uint256 _timeInterval, uint256 _interestRate, uint256 _timeThreshold, bytes _signature) payable returns()
func (_AgentPod *AgentPodTransactor) DepositAsset(opts *bind.TransactOpts, _asset [32]byte, _amount *big.Int, _timeInterval *big.Int, _interestRate *big.Int, _timeThreshold *big.Int, _signature []byte) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "depositAsset", _asset, _amount, _timeInterval, _interestRate, _timeThreshold, _signature)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xff556fbd.
//
// Solidity: function depositAsset(bytes32 _asset, uint256 _amount, uint256 _timeInterval, uint256 _interestRate, uint256 _timeThreshold, bytes _signature) payable returns()
func (_AgentPod *AgentPodSession) DepositAsset(_asset [32]byte, _amount *big.Int, _timeInterval *big.Int, _interestRate *big.Int, _timeThreshold *big.Int, _signature []byte) (*types.Transaction, error) {
	return _AgentPod.Contract.DepositAsset(&_AgentPod.TransactOpts, _asset, _amount, _timeInterval, _interestRate, _timeThreshold, _signature)
}

// DepositAsset is a paid mutator transaction binding the contract method 0xff556fbd.
//
// Solidity: function depositAsset(bytes32 _asset, uint256 _amount, uint256 _timeInterval, uint256 _interestRate, uint256 _timeThreshold, bytes _signature) payable returns()
func (_AgentPod *AgentPodTransactorSession) DepositAsset(_asset [32]byte, _amount *big.Int, _timeInterval *big.Int, _interestRate *big.Int, _timeThreshold *big.Int, _signature []byte) (*types.Transaction, error) {
	return _AgentPod.Contract.DepositAsset(&_AgentPod.TransactOpts, _asset, _amount, _timeInterval, _interestRate, _timeThreshold, _signature)
}

// SetOppurtunity is a paid mutator transaction binding the contract method 0x63c41080.
//
// Solidity: function setOppurtunity(address _oppurtunity) returns()
func (_AgentPod *AgentPodTransactor) SetOppurtunity(opts *bind.TransactOpts, _oppurtunity common.Address) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "setOppurtunity", _oppurtunity)
}

// SetOppurtunity is a paid mutator transaction binding the contract method 0x63c41080.
//
// Solidity: function setOppurtunity(address _oppurtunity) returns()
func (_AgentPod *AgentPodSession) SetOppurtunity(_oppurtunity common.Address) (*types.Transaction, error) {
	return _AgentPod.Contract.SetOppurtunity(&_AgentPod.TransactOpts, _oppurtunity)
}

// SetOppurtunity is a paid mutator transaction binding the contract method 0x63c41080.
//
// Solidity: function setOppurtunity(address _oppurtunity) returns()
func (_AgentPod *AgentPodTransactorSession) SetOppurtunity(_oppurtunity common.Address) (*types.Transaction, error) {
	return _AgentPod.Contract.SetOppurtunity(&_AgentPod.TransactOpts, _oppurtunity)
}

// SetOracle is a paid mutator transaction binding the contract method 0x5b89c489.
//
// Solidity: function setOracle() returns()
func (_AgentPod *AgentPodTransactor) SetOracle(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "setOracle")
}

// SetOracle is a paid mutator transaction binding the contract method 0x5b89c489.
//
// Solidity: function setOracle() returns()
func (_AgentPod *AgentPodSession) SetOracle() (*types.Transaction, error) {
	return _AgentPod.Contract.SetOracle(&_AgentPod.TransactOpts)
}

// SetOracle is a paid mutator transaction binding the contract method 0x5b89c489.
//
// Solidity: function setOracle() returns()
func (_AgentPod *AgentPodTransactorSession) SetOracle() (*types.Transaction, error) {
	return _AgentPod.Contract.SetOracle(&_AgentPod.TransactOpts)
}

// SetPool is a paid mutator transaction binding the contract method 0x4437152a.
//
// Solidity: function setPool(address _pool) returns()
func (_AgentPod *AgentPodTransactor) SetPool(opts *bind.TransactOpts, _pool common.Address) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "setPool", _pool)
}

// SetPool is a paid mutator transaction binding the contract method 0x4437152a.
//
// Solidity: function setPool(address _pool) returns()
func (_AgentPod *AgentPodSession) SetPool(_pool common.Address) (*types.Transaction, error) {
	return _AgentPod.Contract.SetPool(&_AgentPod.TransactOpts, _pool)
}

// SetPool is a paid mutator transaction binding the contract method 0x4437152a.
//
// Solidity: function setPool(address _pool) returns()
func (_AgentPod *AgentPodTransactorSession) SetPool(_pool common.Address) (*types.Transaction, error) {
	return _AgentPod.Contract.SetPool(&_AgentPod.TransactOpts, _pool)
}

// TakeAction is a paid mutator transaction binding the contract method 0x8bdda81e.
//
// Solidity: function takeAction(bytes payload) returns()
func (_AgentPod *AgentPodTransactor) TakeAction(opts *bind.TransactOpts, payload []byte) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "takeAction", payload)
}

// TakeAction is a paid mutator transaction binding the contract method 0x8bdda81e.
//
// Solidity: function takeAction(bytes payload) returns()
func (_AgentPod *AgentPodSession) TakeAction(payload []byte) (*types.Transaction, error) {
	return _AgentPod.Contract.TakeAction(&_AgentPod.TransactOpts, payload)
}

// TakeAction is a paid mutator transaction binding the contract method 0x8bdda81e.
//
// Solidity: function takeAction(bytes payload) returns()
func (_AgentPod *AgentPodTransactorSession) TakeAction(payload []byte) (*types.Transaction, error) {
	return _AgentPod.Contract.TakeAction(&_AgentPod.TransactOpts, payload)
}

// TransferAssetToPool is a paid mutator transaction binding the contract method 0x0f2a731a.
//
// Solidity: function transferAssetToPool(bytes32[] _token, uint256[] _amount) returns()
func (_AgentPod *AgentPodTransactor) TransferAssetToPool(opts *bind.TransactOpts, _token [][32]byte, _amount []*big.Int) (*types.Transaction, error) {
	return _AgentPod.contract.Transact(opts, "transferAssetToPool", _token, _amount)
}

// TransferAssetToPool is a paid mutator transaction binding the contract method 0x0f2a731a.
//
// Solidity: function transferAssetToPool(bytes32[] _token, uint256[] _amount) returns()
func (_AgentPod *AgentPodSession) TransferAssetToPool(_token [][32]byte, _amount []*big.Int) (*types.Transaction, error) {
	return _AgentPod.Contract.TransferAssetToPool(&_AgentPod.TransactOpts, _token, _amount)
}

// TransferAssetToPool is a paid mutator transaction binding the contract method 0x0f2a731a.
//
// Solidity: function transferAssetToPool(bytes32[] _token, uint256[] _amount) returns()
func (_AgentPod *AgentPodTransactorSession) TransferAssetToPool(_token [][32]byte, _amount []*big.Int) (*types.Transaction, error) {
	return _AgentPod.Contract.TransferAssetToPool(&_AgentPod.TransactOpts, _token, _amount)
}

// AgentPodSlashAgentIterator is returned from FilterSlashAgent and is used to iterate over the raw logs and unpacked data for SlashAgent events raised by the AgentPod contract.
type AgentPodSlashAgentIterator struct {
	Event *AgentPodSlashAgent // Event containing the contract specifics and raw log

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
func (it *AgentPodSlashAgentIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentPodSlashAgent)
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
		it.Event = new(AgentPodSlashAgent)
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
func (it *AgentPodSlashAgentIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentPodSlashAgentIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentPodSlashAgent represents a SlashAgent event raised by the AgentPod contract.
type AgentPodSlashAgent struct {
	Agent         common.Address
	AmountToSlash *big.Int
	ChainID       [32]byte
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSlashAgent is a free log retrieval operation binding the contract event 0xcc645481284a350e914981d4131b0f7e9594a28dbaff2a6541e6e903a28ee4f2.
//
// Solidity: event SlashAgent(address _agent, uint256 _amountToSlash, bytes32 _chainID)
func (_AgentPod *AgentPodFilterer) FilterSlashAgent(opts *bind.FilterOpts) (*AgentPodSlashAgentIterator, error) {

	logs, sub, err := _AgentPod.contract.FilterLogs(opts, "SlashAgent")
	if err != nil {
		return nil, err
	}
	return &AgentPodSlashAgentIterator{contract: _AgentPod.contract, event: "SlashAgent", logs: logs, sub: sub}, nil
}

// WatchSlashAgent is a free log subscription operation binding the contract event 0xcc645481284a350e914981d4131b0f7e9594a28dbaff2a6541e6e903a28ee4f2.
//
// Solidity: event SlashAgent(address _agent, uint256 _amountToSlash, bytes32 _chainID)
func (_AgentPod *AgentPodFilterer) WatchSlashAgent(opts *bind.WatchOpts, sink chan<- *AgentPodSlashAgent) (event.Subscription, error) {

	logs, sub, err := _AgentPod.contract.WatchLogs(opts, "SlashAgent")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentPodSlashAgent)
				if err := _AgentPod.contract.UnpackLog(event, "SlashAgent", log); err != nil {
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

// ParseSlashAgent is a log parse operation binding the contract event 0xcc645481284a350e914981d4131b0f7e9594a28dbaff2a6541e6e903a28ee4f2.
//
// Solidity: event SlashAgent(address _agent, uint256 _amountToSlash, bytes32 _chainID)
func (_AgentPod *AgentPodFilterer) ParseSlashAgent(log types.Log) (*AgentPodSlashAgent, error) {
	event := new(AgentPodSlashAgent)
	if err := _AgentPod.contract.UnpackLog(event, "SlashAgent", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

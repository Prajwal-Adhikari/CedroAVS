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

// AgentMetaData contains all meta data concerning the Agent contract.
var AgentMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"balance\",\"type\":\"uint256\"}],\"name\":\"Balance\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"height\",\"type\":\"uint256\"}],\"name\":\"updateBalance\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// AgentABI is the input ABI used to generate the binding from.
// Deprecated: Use AgentMetaData.ABI instead.
var AgentABI = AgentMetaData.ABI

// Agent is an auto generated Go binding around an Ethereum contract.
type Agent struct {
	AgentCaller     // Read-only binding to the contract
	AgentTransactor // Write-only binding to the contract
	AgentFilterer   // Log filterer for contract events
}

// AgentCaller is an auto generated read-only Go binding around an Ethereum contract.
type AgentCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AgentTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AgentFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AgentSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AgentSession struct {
	Contract     *Agent            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AgentCallerSession struct {
	Contract *AgentCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AgentTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AgentTransactorSession struct {
	Contract     *AgentTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AgentRaw is an auto generated low-level Go binding around an Ethereum contract.
type AgentRaw struct {
	Contract *Agent // Generic contract binding to access the raw methods on
}

// AgentCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AgentCallerRaw struct {
	Contract *AgentCaller // Generic read-only contract binding to access the raw methods on
}

// AgentTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AgentTransactorRaw struct {
	Contract *AgentTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAgent creates a new instance of Agent, bound to a specific deployed contract.
func NewAgent(address common.Address, backend bind.ContractBackend) (*Agent, error) {
	contract, err := bindAgent(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Agent{AgentCaller: AgentCaller{contract: contract}, AgentTransactor: AgentTransactor{contract: contract}, AgentFilterer: AgentFilterer{contract: contract}}, nil
}

// NewAgentCaller creates a new read-only instance of Agent, bound to a specific deployed contract.
func NewAgentCaller(address common.Address, caller bind.ContractCaller) (*AgentCaller, error) {
	contract, err := bindAgent(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AgentCaller{contract: contract}, nil
}

// NewAgentTransactor creates a new write-only instance of Agent, bound to a specific deployed contract.
func NewAgentTransactor(address common.Address, transactor bind.ContractTransactor) (*AgentTransactor, error) {
	contract, err := bindAgent(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AgentTransactor{contract: contract}, nil
}

// NewAgentFilterer creates a new log filterer instance of Agent, bound to a specific deployed contract.
func NewAgentFilterer(address common.Address, filterer bind.ContractFilterer) (*AgentFilterer, error) {
	contract, err := bindAgent(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AgentFilterer{contract: contract}, nil
}

// bindAgent binds a generic wrapper to an already deployed contract.
func bindAgent(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AgentMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Agent *AgentRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Agent.Contract.AgentCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Agent *AgentRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Agent.Contract.AgentTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Agent *AgentRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Agent.Contract.AgentTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Agent *AgentCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Agent.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Agent *AgentTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Agent.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Agent *AgentTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Agent.Contract.contract.Transact(opts, method, params...)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x58453fef.
//
// Solidity: function updateBalance(uint256 reqId, uint256 height) returns()
func (_Agent *AgentTransactor) UpdateBalance(opts *bind.TransactOpts, reqId *big.Int, height *big.Int) (*types.Transaction, error) {
	return _Agent.contract.Transact(opts, "updateBalance", reqId, height)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x58453fef.
//
// Solidity: function updateBalance(uint256 reqId, uint256 height) returns()
func (_Agent *AgentSession) UpdateBalance(reqId *big.Int, height *big.Int) (*types.Transaction, error) {
	return _Agent.Contract.UpdateBalance(&_Agent.TransactOpts, reqId, height)
}

// UpdateBalance is a paid mutator transaction binding the contract method 0x58453fef.
//
// Solidity: function updateBalance(uint256 reqId, uint256 height) returns()
func (_Agent *AgentTransactorSession) UpdateBalance(reqId *big.Int, height *big.Int) (*types.Transaction, error) {
	return _Agent.Contract.UpdateBalance(&_Agent.TransactOpts, reqId, height)
}

// AgentBalanceIterator is returned from FilterBalance and is used to iterate over the raw logs and unpacked data for Balance events raised by the Agent contract.
type AgentBalanceIterator struct {
	Event *AgentBalance // Event containing the contract specifics and raw log

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
func (it *AgentBalanceIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AgentBalance)
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
		it.Event = new(AgentBalance)
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
func (it *AgentBalanceIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AgentBalanceIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AgentBalance represents a Balance event raised by the Agent contract.
type AgentBalance struct {
	ReqId   *big.Int
	Height  *big.Int
	Balance *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterBalance is a free log retrieval operation binding the contract event 0xce422aac32fa3a960eff96496b10fd05b1688cfaa1f8d2732407f0c4ac52dbc1.
//
// Solidity: event Balance(uint256 reqId, uint256 height, uint256 balance)
func (_Agent *AgentFilterer) FilterBalance(opts *bind.FilterOpts) (*AgentBalanceIterator, error) {

	logs, sub, err := _Agent.contract.FilterLogs(opts, "Balance")
	if err != nil {
		return nil, err
	}
	return &AgentBalanceIterator{contract: _Agent.contract, event: "Balance", logs: logs, sub: sub}, nil
}

// WatchBalance is a free log subscription operation binding the contract event 0xce422aac32fa3a960eff96496b10fd05b1688cfaa1f8d2732407f0c4ac52dbc1.
//
// Solidity: event Balance(uint256 reqId, uint256 height, uint256 balance)
func (_Agent *AgentFilterer) WatchBalance(opts *bind.WatchOpts, sink chan<- *AgentBalance) (event.Subscription, error) {

	logs, sub, err := _Agent.contract.WatchLogs(opts, "Balance")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AgentBalance)
				if err := _Agent.contract.UnpackLog(event, "Balance", log); err != nil {
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

// ParseBalance is a log parse operation binding the contract event 0xce422aac32fa3a960eff96496b10fd05b1688cfaa1f8d2732407f0c4ac52dbc1.
//
// Solidity: event Balance(uint256 reqId, uint256 height, uint256 balance)
func (_Agent *AgentFilterer) ParseBalance(log types.Log) (*AgentBalance, error) {
	event := new(AgentBalance)
	if err := _Agent.contract.UnpackLog(event, "Balance", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

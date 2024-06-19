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

// AVSMetaData contains all meta data concerning the AVS contract.
var AVSMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"reqId\",\"type\":\"uint256\"}],\"name\":\"GetChainTVL\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"reqId\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AVSABI is the input ABI used to generate the binding from.
// Deprecated: Use AVSMetaData.ABI instead.
var AVSABI = AVSMetaData.ABI

// AVS is an auto generated Go binding around an Ethereum contract.
type AVS struct {
	AVSCaller     // Read-only binding to the contract
	AVSTransactor // Write-only binding to the contract
	AVSFilterer   // Log filterer for contract events
}

// AVSCaller is an auto generated read-only Go binding around an Ethereum contract.
type AVSCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AVSTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AVSTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AVSFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AVSFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AVSSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AVSSession struct {
	Contract     *AVS              // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AVSCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AVSCallerSession struct {
	Contract *AVSCaller    // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// AVSTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AVSTransactorSession struct {
	Contract     *AVSTransactor    // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AVSRaw is an auto generated low-level Go binding around an Ethereum contract.
type AVSRaw struct {
	Contract *AVS // Generic contract binding to access the raw methods on
}

// AVSCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AVSCallerRaw struct {
	Contract *AVSCaller // Generic read-only contract binding to access the raw methods on
}

// AVSTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AVSTransactorRaw struct {
	Contract *AVSTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAVS creates a new instance of AVS, bound to a specific deployed contract.
func NewAVS(address common.Address, backend bind.ContractBackend) (*AVS, error) {
	contract, err := bindAVS(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AVS{AVSCaller: AVSCaller{contract: contract}, AVSTransactor: AVSTransactor{contract: contract}, AVSFilterer: AVSFilterer{contract: contract}}, nil
}

// NewAVSCaller creates a new read-only instance of AVS, bound to a specific deployed contract.
func NewAVSCaller(address common.Address, caller bind.ContractCaller) (*AVSCaller, error) {
	contract, err := bindAVS(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AVSCaller{contract: contract}, nil
}

// NewAVSTransactor creates a new write-only instance of AVS, bound to a specific deployed contract.
func NewAVSTransactor(address common.Address, transactor bind.ContractTransactor) (*AVSTransactor, error) {
	contract, err := bindAVS(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AVSTransactor{contract: contract}, nil
}

// NewAVSFilterer creates a new log filterer instance of AVS, bound to a specific deployed contract.
func NewAVSFilterer(address common.Address, filterer bind.ContractFilterer) (*AVSFilterer, error) {
	contract, err := bindAVS(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AVSFilterer{contract: contract}, nil
}

// bindAVS binds a generic wrapper to an already deployed contract.
func bindAVS(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := AVSMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AVS *AVSRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AVS.Contract.AVSCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AVS *AVSRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AVS.Contract.AVSTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AVS *AVSRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AVS.Contract.AVSTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AVS *AVSCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AVS.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AVS *AVSTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AVS.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AVS *AVSTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AVS.Contract.contract.Transact(opts, method, params...)
}

// ReqId is a free data retrieval call binding the contract method 0x8997ce4e.
//
// Solidity: function reqId() view returns(uint256)
func (_AVS *AVSCaller) ReqId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _AVS.contract.Call(opts, &out, "reqId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// ReqId is a free data retrieval call binding the contract method 0x8997ce4e.
//
// Solidity: function reqId() view returns(uint256)
func (_AVS *AVSSession) ReqId() (*big.Int, error) {
	return _AVS.Contract.ReqId(&_AVS.CallOpts)
}

// ReqId is a free data retrieval call binding the contract method 0x8997ce4e.
//
// Solidity: function reqId() view returns(uint256)
func (_AVS *AVSCallerSession) ReqId() (*big.Int, error) {
	return _AVS.Contract.ReqId(&_AVS.CallOpts)
}

// AVSGetChainTVLIterator is returned from FilterGetChainTVL and is used to iterate over the raw logs and unpacked data for GetChainTVL events raised by the AVS contract.
type AVSGetChainTVLIterator struct {
	Event *AVSGetChainTVL // Event containing the contract specifics and raw log

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
func (it *AVSGetChainTVLIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(AVSGetChainTVL)
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
		it.Event = new(AVSGetChainTVL)
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
func (it *AVSGetChainTVLIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *AVSGetChainTVLIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// AVSGetChainTVL represents a GetChainTVL event raised by the AVS contract.
type AVSGetChainTVL struct {
	ReqId *big.Int
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterGetChainTVL is a free log retrieval operation binding the contract event 0x2d5e6ecde4cb62765ba53ac5731ad9860654adce9c147465c685927ae4a37287.
//
// Solidity: event GetChainTVL(uint256 reqId)
func (_AVS *AVSFilterer) FilterGetChainTVL(opts *bind.FilterOpts) (*AVSGetChainTVLIterator, error) {

	logs, sub, err := _AVS.contract.FilterLogs(opts, "GetChainTVL")
	if err != nil {
		return nil, err
	}
	return &AVSGetChainTVLIterator{contract: _AVS.contract, event: "GetChainTVL", logs: logs, sub: sub}, nil
}

// WatchGetChainTVL is a free log subscription operation binding the contract event 0x2d5e6ecde4cb62765ba53ac5731ad9860654adce9c147465c685927ae4a37287.
//
// Solidity: event GetChainTVL(uint256 reqId)
func (_AVS *AVSFilterer) WatchGetChainTVL(opts *bind.WatchOpts, sink chan<- *AVSGetChainTVL) (event.Subscription, error) {

	logs, sub, err := _AVS.contract.WatchLogs(opts, "GetChainTVL")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(AVSGetChainTVL)
				if err := _AVS.contract.UnpackLog(event, "GetChainTVL", log); err != nil {
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

// ParseGetChainTVL is a log parse operation binding the contract event 0x2d5e6ecde4cb62765ba53ac5731ad9860654adce9c147465c685927ae4a37287.
//
// Solidity: event GetChainTVL(uint256 reqId)
func (_AVS *AVSFilterer) ParseGetChainTVL(log types.Log) (*AVSGetChainTVL, error) {
	event := new(AVSGetChainTVL)
	if err := _AVS.contract.UnpackLog(event, "GetChainTVL", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

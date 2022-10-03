// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package Create2Deployer

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
)

// Create2DeployerMetaData contains all meta data concerning the Create2Deployer contract.
var Create2DeployerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"}],\"name\":\"computeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"computeAddressWithDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"code\",\"type\":\"bytes\"}],\"name\":\"deploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"payoutAddress\",\"type\":\"address\"}],\"name\":\"killCreate2Deployer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
	Bin: "0x6080806040523461005b5760008054336001600160a01b0319821681178355916001600160a01b03909116907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e09080a361053d90816100618239f35b600080fdfe60406080815260048036101561001f575b5050361561001d57600080fd5b005b600091823560e01c8063481286e61461040157806356299481146103be578063644704541461036457806366cfa057146101ce578063715018a6146101715780638da5cb5b146101455763f2fde38b146100795750610010565b34610141576020366003190112610141576001600160a01b0382358181169391929084900361013d576100aa610437565b83156100eb57505082546001600160a01b0319811683178455167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08380a380f35b906020608492519162461bcd60e51b8352820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b6064820152fd5b8480fd5b8280fd5b83823461016d578160031936011261016d57905490516001600160a01b039091168152602090f35b5080fd5b83346101cb57806003193601126101cb5761018a610437565b80546001600160a01b03198116825581906001600160a01b03167f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e08280a380f35b80fd5b50346101415760603660031901126101415781356044359167ffffffffffffffff92838111610360573660238201121561036057808501359380851161034d57825190601f8601601f19908116603f011682019081118282101761033a5783528481526020938785830193602497368982840101116101415780898993018737840101528047106102f8578151156102b85790516001600160a01b0392863592f51615610279578480f35b5162461bcd60e51b8152928301526019908201527f437265617465323a204661696c6564206f6e206465706c6f79000000000000006044820152606490fd5b60648786888188519362461bcd60e51b85528401528201527f437265617465323a2062797465636f6465206c656e677468206973207a65726f6044820152fd5b835162461bcd60e51b8152808801869052601d818801527f437265617465323a20696e73756666696369656e742062616c616e63650000006044820152606490fd5b634e487b7160e01b885260418752602488fd5b634e487b7160e01b875260418652602487fd5b8580fd5b5091903461016d57602036600319011261016d5735916001600160a01b038316808403610141578280808093610398610437565b47908282156103b5575bf1156103ab5782ff5b51903d90823e3d90fd5b506108fc6103a2565b5082346101cb5760603660031901126101cb57604435926001600160a01b039182851685036101cb57506020936103f991602435903561048f565b915191168152f35b5082346101cb57816003193601126101cb57506104266020923090602435903561048f565b90516001600160a01b039091168152f35b6000546001600160a01b0316330361044b57565b606460405162461bcd60e51b815260206004820152602060248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e65726044820152fd5b919060405192602084019260ff60f81b84526bffffffffffffffffffffffff199060601b16602185015260358401526055830152605582526080820182811067ffffffffffffffff8211176104f157604052905190206001600160a01b031690565b634e487b7160e01b600052604160045260246000fdfea26469706673582212202096dd28ab0777bcb4bc81a3935edff6bf19e3ff281fa69b1ff2557428d03a6f64736f6c63430008110033",
}

// Create2DeployerABI is the input ABI used to generate the binding from.
// Deprecated: Use Create2DeployerMetaData.ABI instead.
var Create2DeployerABI = Create2DeployerMetaData.ABI

// Create2DeployerBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use Create2DeployerMetaData.Bin instead.
var Create2DeployerBin = Create2DeployerMetaData.Bin

// DeployCreate2Deployer deploys a new Ethereum contract, binding an instance of Create2Deployer to it.
func DeployCreate2Deployer(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *Create2Deployer, error) {
	parsed, err := Create2DeployerMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(Create2DeployerBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Create2Deployer{Create2DeployerCaller: Create2DeployerCaller{contract: contract}, Create2DeployerTransactor: Create2DeployerTransactor{contract: contract}, Create2DeployerFilterer: Create2DeployerFilterer{contract: contract}}, nil
}

// Create2Deployer is an auto generated Go binding around an Ethereum contract.
type Create2Deployer struct {
	Create2DeployerCaller     // Read-only binding to the contract
	Create2DeployerTransactor // Write-only binding to the contract
	Create2DeployerFilterer   // Log filterer for contract events
}

// Create2DeployerCaller is an auto generated read-only Go binding around an Ethereum contract.
type Create2DeployerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type Create2DeployerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type Create2DeployerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// Create2DeployerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type Create2DeployerSession struct {
	Contract     *Create2Deployer  // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// Create2DeployerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type Create2DeployerCallerSession struct {
	Contract *Create2DeployerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts          // Call options to use throughout this session
}

// Create2DeployerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type Create2DeployerTransactorSession struct {
	Contract     *Create2DeployerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts          // Transaction auth options to use throughout this session
}

// Create2DeployerRaw is an auto generated low-level Go binding around an Ethereum contract.
type Create2DeployerRaw struct {
	Contract *Create2Deployer // Generic contract binding to access the raw methods on
}

// Create2DeployerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type Create2DeployerCallerRaw struct {
	Contract *Create2DeployerCaller // Generic read-only contract binding to access the raw methods on
}

// Create2DeployerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type Create2DeployerTransactorRaw struct {
	Contract *Create2DeployerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewCreate2Deployer creates a new instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2Deployer(address common.Address, backend bind.ContractBackend) (*Create2Deployer, error) {
	contract, err := bindCreate2Deployer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Create2Deployer{Create2DeployerCaller: Create2DeployerCaller{contract: contract}, Create2DeployerTransactor: Create2DeployerTransactor{contract: contract}, Create2DeployerFilterer: Create2DeployerFilterer{contract: contract}}, nil
}

// NewCreate2DeployerCaller creates a new read-only instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerCaller(address common.Address, caller bind.ContractCaller) (*Create2DeployerCaller, error) {
	contract, err := bindCreate2Deployer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerCaller{contract: contract}, nil
}

// NewCreate2DeployerTransactor creates a new write-only instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerTransactor(address common.Address, transactor bind.ContractTransactor) (*Create2DeployerTransactor, error) {
	contract, err := bindCreate2Deployer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerTransactor{contract: contract}, nil
}

// NewCreate2DeployerFilterer creates a new log filterer instance of Create2Deployer, bound to a specific deployed contract.
func NewCreate2DeployerFilterer(address common.Address, filterer bind.ContractFilterer) (*Create2DeployerFilterer, error) {
	contract, err := bindCreate2Deployer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerFilterer{contract: contract}, nil
}

// bindCreate2Deployer binds a generic wrapper to an already deployed contract.
func bindCreate2Deployer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(Create2DeployerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create2Deployer *Create2DeployerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create2Deployer.Contract.Create2DeployerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create2Deployer *Create2DeployerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Create2DeployerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create2Deployer *Create2DeployerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Create2DeployerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Create2Deployer *Create2DeployerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Create2Deployer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Create2Deployer *Create2DeployerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Create2Deployer *Create2DeployerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Create2Deployer.Contract.contract.Transact(opts, method, params...)
}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerCaller) ComputeAddress(opts *bind.CallOpts, salt [32]byte, codeHash [32]byte) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "computeAddress", salt, codeHash)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerSession) ComputeAddress(salt [32]byte, codeHash [32]byte) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddress(&_Create2Deployer.CallOpts, salt, codeHash)
}

// ComputeAddress is a free data retrieval call binding the contract method 0x481286e6.
//
// Solidity: function computeAddress(bytes32 salt, bytes32 codeHash) view returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) ComputeAddress(salt [32]byte, codeHash [32]byte) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddress(&_Create2Deployer.CallOpts, salt, codeHash)
}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerCaller) ComputeAddressWithDeployer(opts *bind.CallOpts, salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "computeAddressWithDeployer", salt, codeHash, deployer)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerSession) ComputeAddressWithDeployer(salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddressWithDeployer(&_Create2Deployer.CallOpts, salt, codeHash, deployer)
}

// ComputeAddressWithDeployer is a free data retrieval call binding the contract method 0x56299481.
//
// Solidity: function computeAddressWithDeployer(bytes32 salt, bytes32 codeHash, address deployer) pure returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) ComputeAddressWithDeployer(salt [32]byte, codeHash [32]byte, deployer common.Address) (common.Address, error) {
	return _Create2Deployer.Contract.ComputeAddressWithDeployer(&_Create2Deployer.CallOpts, salt, codeHash, deployer)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Create2Deployer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerSession) Owner() (common.Address, error) {
	return _Create2Deployer.Contract.Owner(&_Create2Deployer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_Create2Deployer *Create2DeployerCallerSession) Owner() (common.Address, error) {
	return _Create2Deployer.Contract.Owner(&_Create2Deployer.CallOpts)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerTransactor) Deploy(opts *bind.TransactOpts, value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "deploy", value, salt, code)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerSession) Deploy(value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Deploy(&_Create2Deployer.TransactOpts, value, salt, code)
}

// Deploy is a paid mutator transaction binding the contract method 0x66cfa057.
//
// Solidity: function deploy(uint256 value, bytes32 salt, bytes code) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Deploy(value *big.Int, salt [32]byte, code []byte) (*types.Transaction, error) {
	return _Create2Deployer.Contract.Deploy(&_Create2Deployer.TransactOpts, value, salt, code)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerTransactor) KillCreate2Deployer(opts *bind.TransactOpts, payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "killCreate2Deployer", payoutAddress)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerSession) KillCreate2Deployer(payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.KillCreate2Deployer(&_Create2Deployer.TransactOpts, payoutAddress)
}

// KillCreate2Deployer is a paid mutator transaction binding the contract method 0x64470454.
//
// Solidity: function killCreate2Deployer(address payoutAddress) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) KillCreate2Deployer(payoutAddress common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.KillCreate2Deployer(&_Create2Deployer.TransactOpts, payoutAddress)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Create2Deployer.Contract.RenounceOwnership(&_Create2Deployer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Create2Deployer *Create2DeployerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Create2Deployer.Contract.RenounceOwnership(&_Create2Deployer.TransactOpts)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.TransferOwnership(&_Create2Deployer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Create2Deployer *Create2DeployerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Create2Deployer.Contract.TransferOwnership(&_Create2Deployer.TransactOpts, newOwner)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Create2Deployer.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerSession) Receive() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Receive(&_Create2Deployer.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_Create2Deployer *Create2DeployerTransactorSession) Receive() (*types.Transaction, error) {
	return _Create2Deployer.Contract.Receive(&_Create2Deployer.TransactOpts)
}

// Create2DeployerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Create2Deployer contract.
type Create2DeployerOwnershipTransferredIterator struct {
	Event *Create2DeployerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *Create2DeployerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(Create2DeployerOwnershipTransferred)
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
		it.Event = new(Create2DeployerOwnershipTransferred)
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
func (it *Create2DeployerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *Create2DeployerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// Create2DeployerOwnershipTransferred represents a OwnershipTransferred event raised by the Create2Deployer contract.
type Create2DeployerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Create2Deployer *Create2DeployerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*Create2DeployerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Create2Deployer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &Create2DeployerOwnershipTransferredIterator{contract: _Create2Deployer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Create2Deployer *Create2DeployerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *Create2DeployerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Create2Deployer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(Create2DeployerOwnershipTransferred)
				if err := _Create2Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Create2Deployer *Create2DeployerFilterer) ParseOwnershipTransferred(log types.Log) (*Create2DeployerOwnershipTransferred, error) {
	event := new(Create2DeployerOwnershipTransferred)
	if err := _Create2Deployer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

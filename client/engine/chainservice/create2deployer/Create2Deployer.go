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
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"}],\"name\":\"computeAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"codeHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"deployer\",\"type\":\"address\"}],\"name\":\"computeAddressWithDeployer\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"salt\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"code\",\"type\":\"bytes\"}],\"name\":\"deploy\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657610304908161001c8239f35b600080fdfe60406080815260048036101561001457600080fd5b600091823560e01c8063481286e61461022057806356299481146101da576366cfa0571461004157600080fd5b346101ac5760603660031901126101ac5781356044359167ffffffffffffffff928381116101d657366023820112156101d65780850135938085116101c357825190601f8601601f19908116603f01168201908111828210176101b05783528481526020938785830193602497368982840101116101ac57808989930187378401015280471061016a5781511561012a5790516001600160a01b0392863592f516156100eb578480f35b5162461bcd60e51b8152928301526019908201527f437265617465323a204661696c6564206f6e206465706c6f79000000000000006044820152606490fd5b60648786888188519362461bcd60e51b85528401528201527f437265617465323a2062797465636f6465206c656e677468206973207a65726f6044820152fd5b835162461bcd60e51b8152808801869052601d818801527f437265617465323a20696e73756666696369656e742062616c616e63650000006044820152606490fd5b8280fd5b634e487b7160e01b885260418752602488fd5b634e487b7160e01b875260418652602487fd5b8580fd5b50823461021d57606036600319011261021d57604435926001600160a01b0391828516850361021d5750602093610215916024359035610256565b915191168152f35b80fd5b50823461021d578160031936011261021d575061024560209230906024359035610256565b90516001600160a01b039091168152f35b919060405192602084019260ff60f81b84526bffffffffffffffffffffffff199060601b16602185015260358401526055830152605582526080820182811067ffffffffffffffff8211176102b857604052905190206001600160a01b031690565b634e487b7160e01b600052604160045260246000fdfea2646970667358221220567a69f8ed14284d8665ec0c4b2debd2808d8e3bc865809fe000b3b94f3cdbb864736f6c63430008110033",
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

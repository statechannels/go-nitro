// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConsensusApp

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

// ExitFormatAllocation is an auto generated low-level Go binding around an user-defined struct.
type ExitFormatAllocation struct {
	Destination    [32]byte
	Amount         *big.Int
	AllocationType uint8
	Metadata       []byte
}

// ExitFormatAssetMetadata is an auto generated low-level Go binding around an user-defined struct.
type ExitFormatAssetMetadata struct {
	AssetType uint8
	Metadata  []byte
}

// ExitFormatSingleAssetExit is an auto generated low-level Go binding around an user-defined struct.
type ExitFormatSingleAssetExit struct {
	Asset         common.Address
	AssetMetadata ExitFormatAssetMetadata
	Allocations   []ExitFormatAllocation
}

// INitroTypesFixedPart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesFixedPart struct {
	Participants      []common.Address
	ChannelNonce      uint64
	AppDefinition     common.Address
	ChallengeDuration *big.Int
}

// INitroTypesRecoveredVariablePart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesRecoveredVariablePart struct {
	VariablePart INitroTypesVariablePart
	SignedBy     *big.Int
}

// INitroTypesVariablePart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesVariablePart struct {
	Outcome []ExitFormatSingleAssetExit
	AppData []byte
	TurnNum *big.Int
	IsFinal bool
}

// ConsensusAppMetaData contains all meta data concerning the ConsensusApp contract.
var ConsensusAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"requireStateSupported\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657610818908161001c8239f35b600080fdfe61010060048036101561001157600080fd5b600091823560e01c63a0c596351461002857600080fd5b346102b457600319906060368301126102b0578235906001600160401b038083116102ac5760808484360301126102ac5760249384358281116102a857366023820112156102a85780870135928084116102a457868460051b830101913683116102a0576044359382851161028a5760409085360301126102a057608086018681108382111761028e576040528689013582811161028a578701963660238901121561028a57898801356100e36100de82610331565b61030c565b988991808b528b6020809c019160051b83010191368311610286578c8c9101915b83831061026e5750505050875288810135838116810361026a57610140916064918a8a015261013560448201610348565b60408a015201610361565b6060870152866101526100de87610331565b809681520191888201925b84841061024057505050505061017690369087016103c6565b905161021057820151908186925b6101c7575060ff905151911603610199578380f35b60405162461bcd60e51b815292830152600a908201526921756e616e696d6f757360b01b6044820152606490fd5b6000198101908082116101fe57169160ff8091169081146101ec576001019180610184565b634e487b7160e01b8752601186528487fd5b634e487b7160e01b8852601187528588fd5b60405162461bcd60e51b8152808601849052600a818601526907c70726f6f667c213d360b41b6044820152606490fd5b833582811161026657899161025b83928d36918801016103c6565b81520193019261015d565b8c80fd5b8b80fd5b819061027984610348565b8152019101908b90610104565b8e80fd5b8a80fd5b634e487b7160e01b8b5260418a52888bfd5b8980fd5b8880fd5b8780fd5b8580fd5b8380fd5b8280fd5b60405190604082018281106001600160401b038211176102d757604052565b634e487b7160e01b600052604160045260246000fd5b60405190608082018281106001600160401b038211176102d757604052565b6040519190601f01601f191682016001600160401b038111838210176102d757604052565b6001600160401b0381116102d75760051b60200190565b35906001600160a01b038216820361035c57565b600080fd5b359065ffffffffffff8216820361035c57565b81601f8201121561035c578035906001600160401b0382116102d7576103a3601f8301601f191660200161030c565b928284526020838301011161035c57816000926020809301838601378301015290565b604091816080528060a052031261035c576103df6102b8565b60805180358060c0526001600160401b03811161035c576080910160a051031261035c5761040b6102ed565b9060c05160805101926001600160401b0384351161035c5760a0518435850190601f8201121561035c5735916104436100de84610331565b91602083858152019460a051873560c05160805101019060208760051b8301011161035c57939591949360200192905b873560c05160805101019160208660051b840101851015610769576001600160401b0385351161035c576060601f198635850160a05103011261035c57604051928360608101106001600160401b036060860111176107545760206104e2916060860160405287350101610348565b83528435893560c0516080510101016040810135906001600160401b03821161035c57604091601f19910160a05103011261035c5761051f6102b8565b85358a3560c05160805101010160046020604083013583010135101561035c578060406020920135010135815285358a3560c0516080510101016001600160401b03604080830135830101351161035c5760a05161058d916040808201359091019081013501602001610374565b602082015260208401528435893560c05160805101010160e0526001600160401b03606060e05101351161035c5760a05160e0516060810135019690603f8801121561035c576105e36100de6020890135610331565b9860208a818a01358152019260a051606060e05101358d8a35903560c0516080510101010190604060208c013560051b8301011161035c57604001935b606060e05101358d8a35903560c0516080510101010190604060208c013560051b830101861015610732578535916001600160401b03831161035c5760a05160809184019003603f19011261035c578160808f938c61067d6102ed565b95604083606060e05101358435843560c05189510101010101013587526060838160e05101358435843560c0518951010101010101356020880152606060e05101359135903560c05185510101010101013560ff8116810361035c578f906040850152606060e0510135908c35903560c05160805101010101019060a0820135926001600160401b03841161035c57610722602094936040869560a051920101610374565b6060820152815201940193610620565b5050604086019a909a5293855292979295505060209384019390920191610473565b60246000634e487b7160e01b81526041600452fd5b975094959350505050815260c051608051016020810135906001600160401b03821161035c5761079d9160a0519101610374565b60208201526107b4604060c0516080510101610361565b6040820152606060c051608051010135801515810361035c576060820152825260206080510135602083015256fea2646970667358221220ce7b46eb01523edc83e3cc095d5dab327aa47454dbdf0696068996967da4133f64736f6c63430008110033",
}

// ConsensusAppABI is the input ABI used to generate the binding from.
// Deprecated: Use ConsensusAppMetaData.ABI instead.
var ConsensusAppABI = ConsensusAppMetaData.ABI

// ConsensusAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use ConsensusAppMetaData.Bin instead.
var ConsensusAppBin = ConsensusAppMetaData.Bin

// DeployConsensusApp deploys a new Ethereum contract, binding an instance of ConsensusApp to it.
func DeployConsensusApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *ConsensusApp, error) {
	parsed, err := ConsensusAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(ConsensusAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &ConsensusApp{ConsensusAppCaller: ConsensusAppCaller{contract: contract}, ConsensusAppTransactor: ConsensusAppTransactor{contract: contract}, ConsensusAppFilterer: ConsensusAppFilterer{contract: contract}}, nil
}

// ConsensusApp is an auto generated Go binding around an Ethereum contract.
type ConsensusApp struct {
	ConsensusAppCaller     // Read-only binding to the contract
	ConsensusAppTransactor // Write-only binding to the contract
	ConsensusAppFilterer   // Log filterer for contract events
}

// ConsensusAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type ConsensusAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ConsensusAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ConsensusAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ConsensusAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ConsensusAppSession struct {
	Contract     *ConsensusApp     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ConsensusAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ConsensusAppCallerSession struct {
	Contract *ConsensusAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ConsensusAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ConsensusAppTransactorSession struct {
	Contract     *ConsensusAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ConsensusAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type ConsensusAppRaw struct {
	Contract *ConsensusApp // Generic contract binding to access the raw methods on
}

// ConsensusAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ConsensusAppCallerRaw struct {
	Contract *ConsensusAppCaller // Generic read-only contract binding to access the raw methods on
}

// ConsensusAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ConsensusAppTransactorRaw struct {
	Contract *ConsensusAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewConsensusApp creates a new instance of ConsensusApp, bound to a specific deployed contract.
func NewConsensusApp(address common.Address, backend bind.ContractBackend) (*ConsensusApp, error) {
	contract, err := bindConsensusApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ConsensusApp{ConsensusAppCaller: ConsensusAppCaller{contract: contract}, ConsensusAppTransactor: ConsensusAppTransactor{contract: contract}, ConsensusAppFilterer: ConsensusAppFilterer{contract: contract}}, nil
}

// NewConsensusAppCaller creates a new read-only instance of ConsensusApp, bound to a specific deployed contract.
func NewConsensusAppCaller(address common.Address, caller bind.ContractCaller) (*ConsensusAppCaller, error) {
	contract, err := bindConsensusApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusAppCaller{contract: contract}, nil
}

// NewConsensusAppTransactor creates a new write-only instance of ConsensusApp, bound to a specific deployed contract.
func NewConsensusAppTransactor(address common.Address, transactor bind.ContractTransactor) (*ConsensusAppTransactor, error) {
	contract, err := bindConsensusApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ConsensusAppTransactor{contract: contract}, nil
}

// NewConsensusAppFilterer creates a new log filterer instance of ConsensusApp, bound to a specific deployed contract.
func NewConsensusAppFilterer(address common.Address, filterer bind.ContractFilterer) (*ConsensusAppFilterer, error) {
	contract, err := bindConsensusApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ConsensusAppFilterer{contract: contract}, nil
}

// bindConsensusApp binds a generic wrapper to an already deployed contract.
func bindConsensusApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ConsensusAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusApp *ConsensusAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusApp.Contract.ConsensusAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusApp *ConsensusAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusApp.Contract.ConsensusAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusApp *ConsensusAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusApp.Contract.ConsensusAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ConsensusApp *ConsensusAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ConsensusApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ConsensusApp *ConsensusAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ConsensusApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ConsensusApp *ConsensusAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ConsensusApp.Contract.contract.Transact(opts, method, params...)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0xa0c59635.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppCaller) RequireStateSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	var out []interface{}
	err := _ConsensusApp.contract.Call(opts, &out, "requireStateSupported", fixedPart, proof, candidate)

	if err != nil {
		return err
	}

	return err

}

// RequireStateSupported is a free data retrieval call binding the contract method 0xa0c59635.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _ConsensusApp.Contract.RequireStateSupported(&_ConsensusApp.CallOpts, fixedPart, proof, candidate)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0xa0c59635.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppCallerSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _ConsensusApp.Contract.RequireStateSupported(&_ConsensusApp.CallOpts, fixedPart, proof, candidate)
}

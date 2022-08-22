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

// ExitFormatSingleAssetExit is an auto generated low-level Go binding around an user-defined struct.
type ExitFormatSingleAssetExit struct {
	Asset       common.Address
	Metadata    []byte
	Allocations []ExitFormatAllocation
}

// INitroTypesFixedPart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesFixedPart struct {
	ChainId           *big.Int
	Participants      []common.Address
	ChannelNonce      *big.Int
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"requireStateSupported\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50610757806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c806350d8c32a14610030575b600080fd5b61004361003e366004610485565b610045565b005b6100696100518561061b565b61005b84866105cd565b6100648461070f565b61006f565b50505050565b8151156100975760405162461bcd60e51b815260040161008e90610567565b60405180910390fd5b8260200151516100aa82602001516100cf565b60ff16146100ca5760405162461bcd60e51b815260040161008e90610543565b505050565b6000805b82156100ea576000198301909216916001016100d3565b90505b919050565b80356001600160a01b03811681146100ed57600080fd5b600082601f830112610119578081fd5b61012b61012683356105af565b61058b565b82358152602080820191908401835b85358110156102e557813586016060818903601f1901121561015a578586fd5b60405180606082011067ffffffffffffffff6060830111171561017957fe5b6060810160405261018c602083016100f2565b815267ffffffffffffffff604083013511156101a6578687fd5b6101b98960206040850135850101610300565b602082015267ffffffffffffffff606083013511156101d6578687fd5b60608201358201915088603f8301126101ed578687fd5b6101fd61012660208401356105af565b602083810135825281019060408401895b60208601358110156102c7576080823587018e03603f19011215610230578a8bfd5b60405180608082011067ffffffffffffffff6080830111171561024f57fe5b6080818101604090815284358901908101358352606081013560208401526102779101610474565b604082015267ffffffffffffffff60a08435890101351115610297578b8cfd5b6102ad8e8435890160a081013501604001610300565b60608201528452602093840193919091019060010161020e565b5050604083015250855250602093840193919091019060010161013a565b509095945050505050565b803580151581146100ed57600080fd5b600082601f830112610310578081fd5b813567ffffffffffffffff81111561032457fe5b610337601f8201601f191660200161058b565b81815284602083860101111561034b578283fd5b816020850160208301379081016020019190915292915050565b600060408284031215610376578081fd5b50919050565b6000604080838503121561038e578182fd5b805181810167ffffffffffffffff82821081831117156103aa57fe5b8184528294508535818111156103bf57600080fd5b8601608081890312156103d157600080fd5b60c0840183811083821117156103e357fe5b85528035828111156103f457600080fd5b61040089828401610109565b84525060208101358281111561041557600080fd5b61042189828401610300565b60608601525061043285820161045e565b6080850152610443606082016102f0565b60a08501525050815260209384013593019290925292915050565b803565ffffffffffff811681146100ed57600080fd5b803560ff811681146100ed57600080fd5b6000806000806060858703121561049a578384fd5b843567ffffffffffffffff808211156104b1578586fd5b9086019060a082890312156104c4578586fd5b909450602086013590808211156104d9578485fd5b818701915087601f8301126104ec578485fd5b8135818111156104fa578586fd5b886020808302850101111561050d578586fd5b60208301955080945050604087013591508082111561052a578283fd5b5061053787828801610365565b91505092959194509250565b6020808252600a908201526921756e616e696d6f757360b01b604082015260600190565b6020808252600a908201526907c70726f6f667c213d360b41b604082015260600190565b60405181810167ffffffffffffffff811182821017156105a757fe5b604052919050565b600067ffffffffffffffff8211156105c357fe5b5060209081020190565b60006105db610126846105af565b8381526020808201919084845b8781101561060f576105fd368335890161037c565b855293820193908201906001016105e8565b50919695505050505050565b600060a0823603121561062c578081fd5b60405160a0810167ffffffffffffffff828210818311171561064a57fe5b8160405284358352602091508185013581811115610666578485fd5b8501905036601f820112610678578384fd5b8035610686610126826105af565b81815283810190838501368685028601870111156106a2578788fd5b8794505b838510156106cb576106b7816100f2565b8352600194909401939185019185016106a6565b50808587015250505050506106e26040840161045e565b60408201526106f3606084016100f2565b60608201526107046080840161045e565b608082015292915050565b600061071b368361037c565b9291505056fea2646970667358221220ab269368a4be0e85e4e0531e69026144868b71453284a71ace29c26d774d8a6c64736f6c63430007060033",
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

// RequireStateSupported is a free data retrieval call binding the contract method 0x50d8c32a.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppCaller) RequireStateSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	var out []interface{}
	err := _ConsensusApp.contract.Call(opts, &out, "requireStateSupported", fixedPart, proof, candidate)

	if err != nil {
		return err
	}

	return err

}

// RequireStateSupported is a free data retrieval call binding the contract method 0x50d8c32a.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _ConsensusApp.Contract.RequireStateSupported(&_ConsensusApp.CallOpts, fixedPart, proof, candidate)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0x50d8c32a.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_ConsensusApp *ConsensusAppCallerSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _ConsensusApp.Contract.RequireStateSupported(&_ConsensusApp.CallOpts, fixedPart, proof, candidate)
}

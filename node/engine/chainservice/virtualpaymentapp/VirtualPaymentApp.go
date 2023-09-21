// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package VirtualPaymentApp

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

// VirtualPaymentAppMetaData contains all meta data concerning the VirtualPaymentApp contract.
var VirtualPaymentAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"stateIsSupported\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657611203908161001c8239f35b600080fdfe60c0604052600436101561001257600080fd5b6000803560e01c639936d8121461002857600080fd5b346101545760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc81813601126101505760043567ffffffffffffffff9182821161014c57608081833603011261014c5760243591838311610144573660238401121561014457826004013591848311610148573660248460051b86010111610148576044359485116101485760409085360301126101445760246100d594600401930190600401610630565b604051938492151583526020604081850152825192836040860152825b84811061012e57505050828201840152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168101030190f35b81810183015188820188015287955082016100f2565b8580fd5b8680fd5b8480fd5b8280fd5b80fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101ab570180359067ffffffffffffffff82116101ab57602001918160051b360383136101ab57565b600080fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81813603018212156101ab570190565b3565ffffffffffff811681036101ab5790565b6080810190811067ffffffffffffffff82111761021257604052565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6040810190811067ffffffffffffffff82111761021257604052565b6020810190811067ffffffffffffffff82111761021257604052565b6060810190811067ffffffffffffffff82111761021257604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761021257604052565b901561030f578035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1813603018212156101ab570190565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b67ffffffffffffffff81116102125760051b60200190565b359073ffffffffffffffffffffffffffffffffffffffff821682036101ab57565b92919267ffffffffffffffff821161021257604051916103bf60207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f8401160184610295565b8294818452818301116101ab578281602093846000960137010152565b9080601f830112156101ab578160206103f793359101610377565b90565b90929192608093845261040c8161033e565b9060409061041c82519384610295565b8295818452602060a0948186520192600592831b860195825187116101ab5780945b87861061044f575050505050505050565b67ffffffffffffffff86358181116101ab5783019060608083885103126101ab57855192818401848110848211176105ef57875261048c81610356565b84528a518101358381116101ab57810187818a5103126101ab57875190888201828110868211176105ef578952803560048110156101ab5782528c51810135908582116101ab576104e0918b5191016103dc565b8c518201528b5185015286810135908382116101ab5701908751601f830112156101ab578135906105108261033e565b9361051d89519586610295565b8285528c519384808701948d1b820101948b5186116101ab5781015b85811061055957505050505050848201528152865101948651019461043e565b80358381116101ab578201948c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08782510301126101ab578f958c519061059f826101f6565b875181013582528d81013588518301528681013560ff811681036101ab57828f0152808f01358681116101ab578f916105de9251918a519101016103dc565b868201528152855101945101610539565b602460007f4e487b710000000000000000000000000000000000000000000000000000000081526041600452fd5b359065ffffffffffff821682036101ab57565b9092918015611033576001811461069f5760646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6261642070726f6f66206c656e677468000000000000000000000000000000006044820152fd5b6106a981856102d6565b6106b38380610157565b9190506040813603126101ab576040516106cc81610241565b813567ffffffffffffffff81116101ab578201906080823603126101ab57604051916106f7836101f6565b803567ffffffffffffffff81116101ab57810136601f820112156101ab576107269036906020813591016103fa565b8352602081013567ffffffffffffffff81116101ab5761074990369083016103dc565b6020840152606061075c6040830161061d565b916040850192835201359283151584036101ab57602065ffffffffffff93826001966060849501528152019401358452511603610fd55761079f60ff915161117d565b1603610f7757600265ffffffffffff6107c360406107bd87806101b0565b016101e3565b1603610f1957600180602085013560021c1603610ebb576107e483806101b0565b6020810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101ab570180359067ffffffffffffffff82116101ab5760200181360381136101ab57610840913691610377565b916080813603126101ab5760405190610858826101f6565b803567ffffffffffffffff81116101ab57810136601f820112156101ab578035906108828261033e565b916108906040519384610295565b80835260208084019160051b830101913683116101ab57602001905b828210610ea3575050508252602081013567ffffffffffffffff811681036101ab576108f19160609160208501526108e660408201610356565b60408501520161061d565b606082015282805181010392608084126101ab5760607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe06040519561093587610241565b6020840151875201126101ab576040519081606081011067ffffffffffffffff606084011117610212576060820160405260408101519060ff821682036101ab57608091835260608101516020840152015160408201526020840152805167ffffffffffffffff6020830151169073ffffffffffffffffffffffffffffffffffffffff60408401511665ffffffffffff606085015116906040519360a08501608060208701528451809152602060c0870195019060005b818110610e7757505050846000946080948294610a3f94604060209a0152606084015285830152037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08101835282610295565b838151910120865160405190858201928352604082015260408152610a6381610279565b5190208387015190604051858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c8152610aa781610279565b5190209060ff8151169060408682015191015191604051938452868401526040830152606082015282805260015afa15610e6b5773ffffffffffffffffffffffffffffffffffffffff60005116908115610e0d57610b1a73ffffffffffffffffffffffffffffffffffffffff9151611170565b511603610daf57610b5e610b54610b4b610b45610b3f610b4595610b669751996102d6565b806101b0565b80610157565b939095806101b0565b94909236916103fa565b9236916103fa565b90600181511480610da4575b80610d7c575b80610d54575b15610cf6576020610bac6040610b9b83610ba383610b9b89611170565b510151611170565b51015194611170565b51015190838203918211610cc75703610c6957610bca604091611170565b51015180516001101561030f57604001516020015103610c0b576001906040516020810181811067ffffffffffffffff821117610212576040526000815290565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f426f62206e6f742061646a757374656420636f72726563746c790000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f416c696365206e6f742061646a757374656420636f72726563746c79000000006044820152fd5b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6f6e6c79206e617469766520617373657420616c6c6f776564000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff610d7383611170565b51511615610b7e565b5073ffffffffffffffffffffffffffffffffffffffff610d9b82611170565b51511615610b78565b506001825114610b72565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f696e76616c6964207369676e617475726520666f7220766f75636865720000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964207369676e61747572650000000000000000000000000000006044820152fd5b6040513d6000823e3d90fd5b825173ffffffffffffffffffffffffffffffffffffffff168752602096870196909201916001016109ec565b60208091610eb084610356565b8152019101906108ac565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f726564656d7074696f6e206e6f74207369676e656420627920426f62000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d316044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f737466756e642021756e616e696d6f75733b207c70726f6f667c3d3100006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6261642070726f6f665b305d2e7475726e4e756d3b207c70726f6f667c3d31006044820152fd5b5090915060ff611050611049602085013561117d565b9280610157565b92905016036111125765ffffffffffff8061107060406107bd85806101b0565b16156111015761108760406107bd846001956101b0565b16146110eb5760646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d306044820152fd5b6001906040516110fa8161025d565b6000815290565b50506001906040516110fa8161025d565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f21756e616e696d6f75733b207c70726f6f667c3d3000000000000000000000006044820152fd5b80511561030f5760200190565b806000915b61118a575090565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810190808211610cc757169060ff809116908114610cc757600101908061118256fea2646970667358221220a01b040d2b2acba03515a643384930166441976faa0d06d3dd650cdf21475bbc64736f6c63430008110033",
}

// VirtualPaymentAppABI is the input ABI used to generate the binding from.
// Deprecated: Use VirtualPaymentAppMetaData.ABI instead.
var VirtualPaymentAppABI = VirtualPaymentAppMetaData.ABI

// VirtualPaymentAppBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VirtualPaymentAppMetaData.Bin instead.
var VirtualPaymentAppBin = VirtualPaymentAppMetaData.Bin

// DeployVirtualPaymentApp deploys a new Ethereum contract, binding an instance of VirtualPaymentApp to it.
func DeployVirtualPaymentApp(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *VirtualPaymentApp, error) {
	parsed, err := VirtualPaymentAppMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VirtualPaymentAppBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &VirtualPaymentApp{VirtualPaymentAppCaller: VirtualPaymentAppCaller{contract: contract}, VirtualPaymentAppTransactor: VirtualPaymentAppTransactor{contract: contract}, VirtualPaymentAppFilterer: VirtualPaymentAppFilterer{contract: contract}}, nil
}

// VirtualPaymentApp is an auto generated Go binding around an Ethereum contract.
type VirtualPaymentApp struct {
	VirtualPaymentAppCaller     // Read-only binding to the contract
	VirtualPaymentAppTransactor // Write-only binding to the contract
	VirtualPaymentAppFilterer   // Log filterer for contract events
}

// VirtualPaymentAppCaller is an auto generated read-only Go binding around an Ethereum contract.
type VirtualPaymentAppCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtualPaymentAppTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VirtualPaymentAppTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtualPaymentAppFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VirtualPaymentAppFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VirtualPaymentAppSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VirtualPaymentAppSession struct {
	Contract     *VirtualPaymentApp // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// VirtualPaymentAppCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VirtualPaymentAppCallerSession struct {
	Contract *VirtualPaymentAppCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// VirtualPaymentAppTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VirtualPaymentAppTransactorSession struct {
	Contract     *VirtualPaymentAppTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// VirtualPaymentAppRaw is an auto generated low-level Go binding around an Ethereum contract.
type VirtualPaymentAppRaw struct {
	Contract *VirtualPaymentApp // Generic contract binding to access the raw methods on
}

// VirtualPaymentAppCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VirtualPaymentAppCallerRaw struct {
	Contract *VirtualPaymentAppCaller // Generic read-only contract binding to access the raw methods on
}

// VirtualPaymentAppTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VirtualPaymentAppTransactorRaw struct {
	Contract *VirtualPaymentAppTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVirtualPaymentApp creates a new instance of VirtualPaymentApp, bound to a specific deployed contract.
func NewVirtualPaymentApp(address common.Address, backend bind.ContractBackend) (*VirtualPaymentApp, error) {
	contract, err := bindVirtualPaymentApp(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &VirtualPaymentApp{VirtualPaymentAppCaller: VirtualPaymentAppCaller{contract: contract}, VirtualPaymentAppTransactor: VirtualPaymentAppTransactor{contract: contract}, VirtualPaymentAppFilterer: VirtualPaymentAppFilterer{contract: contract}}, nil
}

// NewVirtualPaymentAppCaller creates a new read-only instance of VirtualPaymentApp, bound to a specific deployed contract.
func NewVirtualPaymentAppCaller(address common.Address, caller bind.ContractCaller) (*VirtualPaymentAppCaller, error) {
	contract, err := bindVirtualPaymentApp(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VirtualPaymentAppCaller{contract: contract}, nil
}

// NewVirtualPaymentAppTransactor creates a new write-only instance of VirtualPaymentApp, bound to a specific deployed contract.
func NewVirtualPaymentAppTransactor(address common.Address, transactor bind.ContractTransactor) (*VirtualPaymentAppTransactor, error) {
	contract, err := bindVirtualPaymentApp(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VirtualPaymentAppTransactor{contract: contract}, nil
}

// NewVirtualPaymentAppFilterer creates a new log filterer instance of VirtualPaymentApp, bound to a specific deployed contract.
func NewVirtualPaymentAppFilterer(address common.Address, filterer bind.ContractFilterer) (*VirtualPaymentAppFilterer, error) {
	contract, err := bindVirtualPaymentApp(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VirtualPaymentAppFilterer{contract: contract}, nil
}

// bindVirtualPaymentApp binds a generic wrapper to an already deployed contract.
func bindVirtualPaymentApp(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VirtualPaymentAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VirtualPaymentApp *VirtualPaymentAppRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VirtualPaymentApp.Contract.VirtualPaymentAppCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VirtualPaymentApp *VirtualPaymentAppRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VirtualPaymentApp.Contract.VirtualPaymentAppTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VirtualPaymentApp *VirtualPaymentAppRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VirtualPaymentApp.Contract.VirtualPaymentAppTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_VirtualPaymentApp *VirtualPaymentAppCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _VirtualPaymentApp.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_VirtualPaymentApp *VirtualPaymentAppTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _VirtualPaymentApp.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_VirtualPaymentApp *VirtualPaymentAppTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _VirtualPaymentApp.Contract.contract.Transact(opts, method, params...)
}

// StateIsSupported is a free data retrieval call binding the contract method 0x9936d812.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns(bool, string)
func (_VirtualPaymentApp *VirtualPaymentAppCaller) StateIsSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) (bool, string, error) {
	var out []interface{}
	err := _VirtualPaymentApp.contract.Call(opts, &out, "stateIsSupported", fixedPart, proof, candidate)

	if err != nil {
		return *new(bool), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// StateIsSupported is a free data retrieval call binding the contract method 0x9936d812.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns(bool, string)
func (_VirtualPaymentApp *VirtualPaymentAppSession) StateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) (bool, string, error) {
	return _VirtualPaymentApp.Contract.StateIsSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

// StateIsSupported is a free data retrieval call binding the contract method 0x9936d812.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns(bool, string)
func (_VirtualPaymentApp *VirtualPaymentAppCallerSession) StateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) (bool, string, error) {
	return _VirtualPaymentApp.Contract.StateIsSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

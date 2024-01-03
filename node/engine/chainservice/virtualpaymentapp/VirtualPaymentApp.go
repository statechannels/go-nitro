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
	Bin: "0x60808060405234610016576111bf908161001b8239f35b5f80fdfe60c06040526004361015610011575f80fd5b5f803560e01c639936d81214610025575f80fd5b346101515760607ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc818136011261014d5760043567ffffffffffffffff918282116101495760808183360301126101495760243591838311610141573660238401121561014157826004013591848311610145573660248460051b86010111610145576044359485116101455760409085360301126101415760246100d2946004019301906004016105eb565b604051938492151583526020604081850152825192836040860152825b84811061012b57505050828201840152601f017fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0168101030190f35b81810183015188820188015287955082016100ef565b8580fd5b8680fd5b8480fd5b8280fd5b80fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101a8570180359067ffffffffffffffff82116101a857602001918160051b360383136101a857565b5f80fd5b9035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff81813603018212156101a8570190565b3565ffffffffffff811681036101a85790565b6060810190811067ffffffffffffffff82111761020e57604052565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52604160045260245ffd5b6040810190811067ffffffffffffffff82111761020e57604052565b6080810190811067ffffffffffffffff82111761020e57604052565b6020810190811067ffffffffffffffff82111761020e57604052565b90601f7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0910116810190811067ffffffffffffffff82111761020e57604052565b9015610309578035907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffc1813603018212156101a8570190565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52603260045260245ffd5b67ffffffffffffffff811161020e5760051b60200190565b359073ffffffffffffffffffffffffffffffffffffffff821682036101a857565b92919267ffffffffffffffff821161020e57604051916103b760207fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0601f840116018461028f565b8294818452818301116101a8578281602093845f960137010152565b9080601f830112156101a8578160206103ee9335910161036f565b90565b90929192608093845261040381610336565b906040906104138251938461028f565b8295818452602060a0948186520192600592831b860195825187116101a85780945b878610610446575050505050505050565b67ffffffffffffffff86358181116101a85783019060608083885103126101a857855192610473846101f2565b61047c8161034e565b84528a518101358381116101a857810187818a5103126101a8578751906104a28261023b565b803560048110156101a85782528c51810135908582116101a8576104c9918b5191016103d3565b8c518201528b5185015286810135908382116101a85701908751601f830112156101a8578135906104f982610336565b936105068951958661028f565b8285528c519384808701948d1b820101948b5186116101a85781015b858110610542575050505050508482015281528651019486510194610435565b80358381116101a8578201948c7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe08782510301126101a8578f958c519061058882610257565b875181013582528d81013588518301528681013560ff811681036101a857828f0152808f01358681116101a8578f916105c79251918a519101016103d3565b868201528152855101945101610522565b359065ffffffffffff821682036101a857565b9092918015610ff7576001811461065a5760646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601060248201527f6261642070726f6f66206c656e677468000000000000000000000000000000006044820152fd5b61066481856102d0565b61066e8380610154565b9190506040813603126101a8576040516106878161023b565b813567ffffffffffffffff81116101a8578201906080823603126101a857604051916106b283610257565b803567ffffffffffffffff81116101a857810136601f820112156101a8576106e19036906020813591016103f1565b8352602081013567ffffffffffffffff81116101a85761070490369083016103d3565b60208401526060610717604083016105d8565b916040850192835201359283151584036101a857602065ffffffffffff93826001966060849501528152019401358452511603610f995761075a60ff915161113a565b1603610f3b57600265ffffffffffff61077e604061077887806101ac565b016101df565b1603610edd5761078e8280610154565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810191508111610c8e576001809160ff602087013591161c1603610e7f576107d783806101ac565b6020810135907fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe1813603018212156101a8570180359067ffffffffffffffff82116101a85760200181360381136101a85761083391369161036f565b916080813603126101a8576040519061084b82610257565b803567ffffffffffffffff81116101a857810136601f820112156101a85780359061087582610336565b91610883604051938461028f565b80835260208084019160051b830101913683116101a857602001905b828210610e67575050508252602081013567ffffffffffffffff811681036101a8576108e49160609160208501526108d96040820161034e565b6040850152016105d8565b606082015282805181010392608084126101a85760607fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0604051956109288761023b565b6020840151875201126101a85760405190610942826101f2565b60408101519060ff821682036101a857608091835260608101516020840152015160408201526020840152805167ffffffffffffffff6020830151169073ffffffffffffffffffffffffffffffffffffffff60408401511665ffffffffffff606085015116906040519360a08501608060208701528451809152602060c087019501905f5b818110610e3b57505050845f946080948294610a1994604060209a0152606084015285830152037fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe0810183528261028f565b838151910120865160405190858201928352604082015260408152610a3d816101f2565b5190208387015190604051858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c8152610a81816101f2565b5190209060ff8151169060408682015191015191604051938452868401526040830152606082015282805260015afa15610e305773ffffffffffffffffffffffffffffffffffffffff5f5116908115610dd257610af373ffffffffffffffffffffffffffffffffffffffff915161112d565b511603610d7457610b37610b2d610b24610b1e610b18610b1e95610b3f9751996102d0565b806101ac565b80610154565b939095806101ac565b94909236916103f1565b9236916103f1565b90600181511480610d69575b80610d41575b80610d19575b15610cbb576020610b856040610b7483610b7c83610b748961112d565b51015161112d565b5101519461112d565b51015190838203918211610c8e5703610c3057610ba360409161112d565b51015180516001101561030957604001516020015103610bd257600190604051610bcc81610273565b5f815290565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601a60248201527f426f62206e6f742061646a757374656420636f72726563746c790000000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f416c696365206e6f742061646a757374656420636f72726563746c79000000006044820152fd5b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601960248201527f6f6e6c79206e617469766520617373657420616c6c6f776564000000000000006044820152fd5b5073ffffffffffffffffffffffffffffffffffffffff610d388361112d565b51511615610b57565b5073ffffffffffffffffffffffffffffffffffffffff610d608261112d565b51511615610b51565b506001825114610b4b565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601d60248201527f696e76616c6964207369676e617475726520666f7220766f75636865720000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601160248201527f496e76616c6964207369676e61747572650000000000000000000000000000006044820152fd5b6040513d5f823e3d90fd5b825173ffffffffffffffffffffffffffffffffffffffff168752602096870196909201916001016109c7565b60208091610e748461034e565b81520191019061089f565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601c60248201527f726564656d7074696f6e206e6f74207369676e656420627920426f62000000006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d316044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601e60248201527f706f737466756e642021756e616e696d6f75733b207c70726f6f667c3d3100006044820152fd5b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601f60248201527f6261642070726f6f665b305d2e7475726e4e756d3b207c70726f6f667c3d31006044820152fd5b5090915060ff61101461100d602085013561113a565b9280610154565b92905016036110cf5765ffffffffffff80611034604061077885806101ac565b16156110be5761104b6040610778846001956101ac565b16146110af5760646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d306044820152fd5b600190604051610bcc81610273565b5050600190604051610bcc81610273565b60646040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601560248201527f21756e616e696d6f75733b207c70726f6f667c3d3000000000000000000000006044820152fd5b8051156103095760200190565b805f915b611146575090565b7fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff810190808211610c8e57169060ff809116908114610c8e57600101908061113e56fea264697066735822122027350ec4f707020d05bf4f616324db9ebd25f1041382401342b89ea575c7f7ee64736f6c63430008140033",
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

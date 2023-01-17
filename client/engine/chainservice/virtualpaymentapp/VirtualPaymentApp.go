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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"requireStateSupported\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x6080806040523461001657610e0d908161001c8239f35b600080fdfe60c0604052600436101561001257600080fd5b6000803560e01c630442daa41461002857600080fd5b346100c8576003196060368201126100c4576004356001600160401b03918282116100c05760808183360301126100c057602435918383116100b857366023840112156100b8578260040135918483116100bc573660248460051b860101116100bc576044359485116100bc5760409085360301126100b85760246100b594600401930190600401610435565b80f35b8480fd5b8580fd5b8380fd5b5080fd5b80fd5b903590601e198136030182121561010057018035906001600160401b03821161010057602001918160051b3603831361010057565b600080fd5b903590607e1981360301821215610100570190565b3565ffffffffffff811681036101005790565b901561014857803590603e1981360301821215610100570190565b634e487b7160e01b600052603260045260246000fd5b608081019081106001600160401b0382111761017957604052565b634e487b7160e01b600052604160045260246000fd5b604081019081106001600160401b0382111761017957604052565b606081019081106001600160401b0382111761017957604052565b90601f801991011681019081106001600160401b0382111761017957604052565b6001600160401b0381116101795760051b60200190565b35906001600160a01b038216820361010057565b9291926001600160401b038211610179576040519161023a601f8201601f1916602001846101c5565b829481845281830111610100578281602093846000960137010152565b9080601f830112156101005781602061027293359101610211565b90565b909291926080938452610287816101e6565b90604090610297825193846101c5565b8295818452602060a0948186520192600592831b860195825187116101005780945b8786106102ca575050505050505050565b6001600160401b038635818111610100578301906060808388510312610100578551928184018481108482111761040d578752610306816101fd565b84528a51810135838111610100578851610321918301610257565b8b5185015286810135908382116101005701908751601f830112156101005781359061034c826101e6565b93610359895195866101c5565b8285528c519384808701948d1b820101948b5186116101005781015b8581106103955750505050505084820152815286510194865101946102b9565b8035838111610100578201948c601f19878251030112610100578f958c51906103bd8261015e565b875181013582528d81013588518301528681013560ff8116810361010057828f0152808f0135868111610100578f916103fc9251918a51910101610257565b868201528152855101945101610375565b60246000634e487b7160e01b81526041600452fd5b359065ffffffffffff8216820361010057565b9092918015610c3d576001811461047e5760405162461bcd60e51b815260206004820152601060248201526f0c4c2c840e0e4dedecc40d8cadccee8d60831b6044820152606490fd5b610488818561012d565b61049283806100cb565b919050604081360312610100576040516104ab8161018f565b81356001600160401b0381116101005782019060808236031261010057604051916104d58361015e565b80356001600160401b03811161010057810136601f8201121561010057610503903690602081359101610275565b835260208101356001600160401b038111610100576105259036908301610257565b6020840152606061053860408301610422565b9160408501928352013592831515840361010057602065ffffffffffff93826001966060849501528152019401358452511603610bf85761057b60ff9151610da5565b1603610bb357600265ffffffffffff61059f60406105998780610105565b0161011a565b1603610b6f57600180602085013560021c1603610b2a576105c08380610105565b602081013590601e198136030182121561010057018035906001600160401b038211610100576020018136038113610100576105fd913691610211565b9160808136031261010057604051906106158261015e565b80356001600160401b03811161010057810136601f8201121561010057803561063d816101e6565b9161064b60405193846101c5565b81835260208301903660208460051b83010111610100579060208201915b60208460051b8201018310610b0f5750505050825260208101356001600160401b0381168103610100576106b69160609160208501526106ab604082016101fd565b604085015201610422565b60608201528280518101039260808412610100576060604051946106d98661018f565b60208301518652601f1901126101005760405190606082018281106001600160401b038211176101795760405260408101519060ff821682036101005760809183526060810151602084015201516040820152602084015280516001600160401b036020830151169060018060a01b0360408401511665ffffffffffff606085015116906040519360a08501608060208701528451809152602060c0870195019060005b818110610af0575050508460009460809482946107b294604060209a015260608401528583015203601f1981018352826101c5565b8381519101208651604051908582019283526040820152604081526107d6816101aa565b5190208387015190604051858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c815261081a816101aa565b5190209060ff8151169060408682015191015191604051938452868401526040830152606082015282805260015afa15610ae4576000516001600160a01b0316908115610aab57516001600160a01b039061087490610d98565b511603610a66576108b86108ae6108a561089f61089961089f956108c097519961012d565b80610105565b806100cb565b93909580610105565b9490923691610275565b923691610275565b90600181511480610a5b575b80610a40575b80610a25575b156109e057602061090660406108f5836108fd836108f589610d98565b510151610d98565b51015194610d98565b510151908382039182116109ca570361098557610924604091610d98565b5101518051600110156101485760400151602001510361094057565b60405162461bcd60e51b815260206004820152601a60248201527f426f62206e6f742061646a757374656420636f72726563746c790000000000006044820152606490fd5b60405162461bcd60e51b815260206004820152601c60248201527f416c696365206e6f742061646a757374656420636f72726563746c79000000006044820152606490fd5b634e487b7160e01b600052601160045260246000fd5b60405162461bcd60e51b815260206004820152601960248201527f6f6e6c79206e617469766520617373657420616c6c6f776564000000000000006044820152606490fd5b506001600160a01b03610a3783610d98565b515116156108d8565b506001600160a01b03610a5282610d98565b515116156108d2565b5060018251146108cc565b60405162461bcd60e51b815260206004820152601d60248201527f696e76616c6964207369676e617475726520666f7220766f75636865720000006044820152606490fd5b60405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b6044820152606490fd5b6040513d6000823e3d90fd5b82516001600160a01b031687526020968701969092019160010161077d565b6020808093610b1d866101fd565b8152019301929150610669565b60405162461bcd60e51b815260206004820152601c60248201527f726564656d7074696f6e206e6f74207369676e656420627920426f62000000006044820152606490fd5b606460405162461bcd60e51b815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d316044820152fd5b60405162461bcd60e51b815260206004820152601e60248201527f706f737466756e642021756e616e696d6f75733b207c70726f6f667c3d3100006044820152606490fd5b60405162461bcd60e51b815260206004820152601f60248201527f6261642070726f6f665b305d2e7475726e4e756d3b207c70726f6f667c3d31006044820152606490fd5b5090915060ff610c5a610c536020850135610da5565b92806100cb565b9290501603610d5b5765ffffffffffff80610c7a60406105998580610105565b1615610d5757600181610c9260406105998680610105565b1614610d5757600390610caa60406105998580610105565b1614610cf457606460405162461bcd60e51b815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d306044820152fd5b610d0081606092610105565b013580159081150361010057610d1257565b60405162461bcd60e51b815260206004820152601e60248201527f2166696e616c3b207475726e4e756d3d33202626207c70726f6f667c3d3000006044820152606490fd5b5050565b60405162461bcd60e51b8152602060048201526015602482015274021756e616e696d6f75733b207c70726f6f667c3d3605c1b6044820152606490fd5b8051156101485760200190565b806000915b610db2575090565b6000198101908082116109ca57169060ff8091169081146109ca576001019080610daa56fea2646970667358221220337ebae2885c3ca326247c615fab77e5a564aae3ab02ef65a39d46ac4a85ddb864736f6c63430008110033",
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
	parsed, err := abi.JSON(strings.NewReader(VirtualPaymentAppABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
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

// RequireStateSupported is a free data retrieval call binding the contract method 0x0442daa4.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppCaller) RequireStateSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	var out []interface{}
	err := _VirtualPaymentApp.contract.Call(opts, &out, "requireStateSupported", fixedPart, proof, candidate)

	if err != nil {
		return err
	}

	return err

}

// RequireStateSupported is a free data retrieval call binding the contract method 0x0442daa4.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _VirtualPaymentApp.Contract.RequireStateSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0x0442daa4.
//
// Solidity: function requireStateSupported((address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppCallerSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _VirtualPaymentApp.Contract.RequireStateSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

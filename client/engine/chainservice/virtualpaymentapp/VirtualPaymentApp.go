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
	ChainId           *big.Int
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"requireStateSupported\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b506112c5806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063438b017214610030575b600080fd5b61004361003e366004610a43565b610045565b005b8161016a576100576020850185610f85565b90506100668260200135610309565b60ff161461008f5760405162461bcd60e51b815260040161008690610df6565b60405180910390fd5b610099818061104b565b6100aa906060810190604001610b96565b65ffffffffffff166100bb57610303565b6100c5818061104b565b6100d6906060810190604001610b96565b65ffffffffffff16600114156100eb57610303565b6100f5818061104b565b610106906060810190604001610b96565b65ffffffffffff166003141561015257610120818061104b565b610131906080810190606001610a29565b61014d5760405162461bcd60e51b815260040161008690610cfe565b610303565b60405162461bcd60e51b815260040161008690610d97565b60018214156102eb576101b28383600081811061018357fe5b9050602002810190610195919061102c565b61019e906111a0565b6101ab6020870187610f85565b905061032c565b6101bc818061104b565b6101cd906060810190604001610b96565b65ffffffffffff166002146101f45760405162461bcd60e51b815260040161008690610c0d565b6102038160200135600261038e565b61021f5760405162461bcd60e51b815260040161008690610e5c565b600061028161022e838061104b565b61023c906020810190610fe8565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201919091525061027c92508991506110ad9050565b61039b565b90506102e58484600081811061029357fe5b90506020028101906102a5919061102c565b6102af908061104b565b6102b99080610fd2565b6102c2916110a0565b6102cc848061104b565b6102d69080610fd2565b6102df916110a0565b83610449565b50610303565b60405162461bcd60e51b815260040161008690610dcc565b50505050565b6000805b82156103245760001983019092169160010161030d565b90505b919050565b81600001516040015165ffffffffffff1660011461035c5760405162461bcd60e51b815260040161008690610cc7565b8061036a8360200151610309565b60ff161461038a5760405162461bcd60e51b815260040161008690610e25565b5050565b60ff161c60019081161490565b600080838060200190518101906103b29190610b00565b905060006103f56103c285610633565b83516040516103d5929190602001610be1565b6040516020818303038152906040528051906020012083602001516106a2565b9050836020015160008151811061040857fe5b60200260200101516001600160a01b0316816001600160a01b0316146104405760405162461bcd60e51b815260040161008690610eca565b50519392505050565b8251600114801561045b575081516001145b8015610491575060006001600160a01b03168360008151811061047a57fe5b6020026020010151600001516001600160a01b0316145b80156104c7575060006001600160a01b0316826000815181106104b057fe5b6020026020010151600001516001600160a01b0316145b6104e35760405162461bcd60e51b815260040161008690610e93565b826000815181106104f057fe5b6020026020010151604001516000600181111561050957fe5b8151811061051357fe5b60200260200101516020015181111561053e5760405162461bcd60e51b815260040161008690610c42565b808360008151811061054c57fe5b6020026020010151604001516000600181111561056557fe5b8151811061056f57fe5b602002602001015160200151038260008151811061058957fe5b602002602001015160400151600060018111156105a257fe5b815181106105ac57fe5b602002602001015160200151146105d55760405162461bcd60e51b815260040161008690610d35565b80826000815181106105e357fe5b6020026020010151604001516001808111156105fb57fe5b8151811061060557fe5b6020026020010151602001511461062e5760405162461bcd60e51b815260040161008690610c90565b505050565b600061063d61075c565b82511461065c5760405162461bcd60e51b815260040161008690610c65565b815160208084015160408086015160608701516080880151925161068596959293919201610f01565b604051602081830303815290604052805190602001209050919050565b600080836040516020016106b69190610bb0565b6040516020818303038152906040528051906020012090506000600182856000015186602001518760400151604051600081526020016040526040516106ff9493929190610bef565b6020604051602081039080840390855afa158015610721573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166107545760405162461bcd60e51b815260040161008690610d6c565b949350505050565b4690565b600061077361076e84611083565b611060565b8381529050602081018260005b8581101561091f57813585016060818903121561079c57600080fd5b6040518060608201106001600160401b03606083011117156107ba57fe5b606081016040526107ca82610929565b81526001600160401b03602083013511156107e457600080fd5b6107f48960208401358401610976565b60208201526001600160401b036040830135111561081157600080fd5b60408201358201915088601f83011261082957600080fd5b61083661076e8335611083565b8235815260208082019190840160005b8535811015610901576080823587018e03601f1901121561086657600080fd5b6040518060808201106001600160401b036080830111171561088457fe5b6080810160405260208335880101358152604083358801013560208201526108b160608435890101610a1e565b60408201526001600160401b03608084358901013511156108d157600080fd5b6108e78e84358901608081013501602001610976565b606082015284526020938401939190910190600101610846565b50506040830152508452506020928301929190910190600101610780565b5050509392505050565b80356001600160a01b038116811461032757600080fd5b600082601f830112610950578081fd5b61095f83833560208501610760565b9392505050565b8035801515811461032757600080fd5b600082601f830112610986578081fd5b81356001600160401b0381111561099957fe5b6109ac601f8201601f1916602001611060565b8181528460208386010111156109c0578283fd5b816020850160208301379081016020019190915292915050565b6000604082840312156109eb578081fd5b50919050565b803565ffffffffffff8116811461032757600080fd5b80356001600160401b038116811461032757600080fd5b80356103278161127d565b600060208284031215610a3a578081fd5b61095f82610966565b60008060008060608587031215610a58578283fd5b84356001600160401b0380821115610a6e578485fd5b9086019060a08289031215610a81578485fd5b90945060208601359080821115610a96578485fd5b818701915087601f830112610aa9578485fd5b813581811115610ab7578586fd5b8860208083028501011115610aca578586fd5b602083019550809450506040870135915080821115610ae7578283fd5b50610af4878288016109da565b91505092959194509250565b60008183036080811215610b12578182fd5b604080518181016001600160401b038282108183111715610b2f57fe5b818452865183526060601f1986011215610b47578586fd5b835194506060850191508482108183111715610b5f57fe5b5082526020850151610b708161127d565b835284820151602080850191909152606090950151918301919091529283015250919050565b600060208284031215610ba7578081fd5b61095f826109f1565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b918252602082015260400190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252818101527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d31604082015260600190565b602080825260099082015268756e646572666c6f7760b81b604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b6020808252601a908201527f426f62206e6f742061646a757374656420636f72726563746c79000000000000604082015260600190565b6020808252601f908201527f6261642070726f6f665b305d2e7475726e4e756d3b207c70726f6f667c3d3100604082015260600190565b6020808252601e908201527f2166696e616c3b207475726e4e756d3d33202626207c70726f6f667c3d300000604082015260600190565b6020808252601c908201527f416c696365206e6f742061646a757374656420636f72726563746c7900000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252818101527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d30604082015260600190565b60208082526010908201526f0c4c2c840e0e4dedecc40d8cadccee8d60831b604082015260600190565b602080825260159082015274021756e616e696d6f75733b207c70726f6f667c3d3605c1b604082015260600190565b6020808252601e908201527f706f737466756e642021756e616e696d6f75733b207c70726f6f667c3d310000604082015260600190565b6020808252601c908201527f726564656d7074696f6e206e6f74207369676e656420627920426f6200000000604082015260600190565b60208082526019908201527f6f6e6c79206e617469766520617373657420616c6c6f77656400000000000000604082015260600190565b6020808252601d908201527f696e76616c6964207369676e617475726520666f7220766f7563686572000000604082015260600190565b600060a08201878352602060a08185015281885180845260c086019150828a019350845b81811015610f4a5784516001600160a01b031683529383019391830191600101610f25565b50506001600160401b039790971660408501525050506001600160a01b0392909216606083015265ffffffffffff1660809091015292915050565b6000808335601e19843603018112610f9b578283fd5b8301803591506001600160401b03821115610fb4578283fd5b6020908101925081023603821315610fcb57600080fd5b9250929050565b6000808335601e19843603018112610f9b578182fd5b6000808335601e19843603018112610ffe578182fd5b8301803591506001600160401b03821115611017578283fd5b602001915036819003821315610fcb57600080fd5b60008235603e19833603018112611041578182fd5b9190910192915050565b60008235607e19833603018112611041578182fd5b6040518181016001600160401b038111828210171561107b57fe5b604052919050565b60006001600160401b0382111561109657fe5b5060209081020190565b600061095f368484610760565b600060a082360312156110be578081fd5b60405160a081016001600160401b0382821081831117156110db57fe5b81604052843583526020915081850135818111156110f7578485fd5b8501905036601f820112611109578384fd5b803561111761076e82611083565b8181528381019083850136868502860187011115611133578788fd5b8794505b8385101561115c5761114881610929565b835260019490940193918501918501611137565b508085870152505050505061117360408401610a07565b604082015261118460608401610929565b6060820152611195608084016109f1565b608082015292915050565b600060408083360312156111b2578182fd5b80518181016001600160401b0382821081831117156111cd57fe5b8184528535818111156111de578586fd5b860160803682900312156111f0578586fd5b60c08401838110838211171561120257fe5b8552803582811115611212578687fd5b61121e36828401610940565b845250602081013582811115611232578687fd5b61123e36828401610976565b60608601525061124f8582016109f1565b608085015261126060608201610966565b60a085015250508152602093840135938101939093525090919050565b60ff8116811461128c57600080fd5b5056fea2646970667358221220e8aa2d32787d7d292a143346875fc304d42e51ed798b0e6ac951bd11d68fcc7764736f6c63430007060033",
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

// RequireStateSupported is a free data retrieval call binding the contract method 0x438b0172.
//
// Solidity: function requireStateSupported((uint256,address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppCaller) RequireStateSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	var out []interface{}
	err := _VirtualPaymentApp.contract.Call(opts, &out, "requireStateSupported", fixedPart, proof, candidate)

	if err != nil {
		return err
	}

	return err

}

// RequireStateSupported is a free data retrieval call binding the contract method 0x438b0172.
//
// Solidity: function requireStateSupported((uint256,address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _VirtualPaymentApp.Contract.RequireStateSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0x438b0172.
//
// Solidity: function requireStateSupported((uint256,address[],uint64,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns()
func (_VirtualPaymentApp *VirtualPaymentAppCallerSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) error {
	return _VirtualPaymentApp.Contract.RequireStateSupported(&_VirtualPaymentApp.CallOpts, fixedPart, proof, candidate)
}

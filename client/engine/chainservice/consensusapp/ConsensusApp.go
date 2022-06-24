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

// INitroTypesSignature is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesSignature struct {
	V uint8
	R [32]byte
	S [32]byte
}

// INitroTypesSignedVariablePart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesSignedVariablePart struct {
	VariablePart INitroTypesVariablePart
	Sigs         []INitroTypesSignature
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
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"latestSupportedState\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50611397806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e0ca7d9314610030575b600080fd5b61004361003e366004610ab3565b610059565b6040516100509190610fb2565b60405180910390f35b610061610749565b6001821461008a5760405162461bcd60e51b815260040161008190610f44565b60405180910390fd5b6100a56100968561125c565b6100a08486611109565b6100df565b828260008181106100b257fe5b90506020028101906100c49190611095565b6100ce90806110b4565b6100d79061134f565b949350505050565b6020820151518151600090839060001981019081106100fa57fe5b60200260200101516000015160400151905061011682846101b4565b6000828260010165ffffffffffff168161012c57fe5b0690506000805b85518110156101ab5761015a8787838151811061014c57fe5b60200260200101518561030f565b8015610185576101858287838151811061017057fe5b6020026020010151600001516040015161043c565b85818151811061019157fe5b602090810291909101015151604001519150600101610133565b50505050505050565b80518083108015906101c65750600081115b6101e25760405162461bcd60e51b815260040161008190610df1565b60ff8311156102035760405162461bcd60e51b815260040161008190610e5d565b60008260008151811061021257fe5b602002602001015160000151604001518360018551038151811061023257fe5b602002602001015160000151604001510365ffffffffffff1690508381111561026d5760405162461bcd60e51b815260040161008190610d60565b6000805b84518110156102e257600085828151811061028857fe5b60200260200101516040015183169050806000146102b85760405162461bcd60e51b815260040161008190610e28565b8582815181106102c457fe5b60200260200101516040015183179250508080600101915050610271565b5060018560020a0381146103085760405162461bcd60e51b815260040161008190610eb7565b5050505050565b6000826020015151116103345760405162461bcd60e51b815260040161008190610f7b565b610341826040015161046f565b60ff16826020015151146103675760405162461bcd60e51b815260040161008190610f16565b610384826040015183600001516040015183866020015151610492565b60006103938360400151610511565b905060005b8151811015610308576104186103cd6103b0876105b8565b86516020810151815160408301516060909301519192909161062e565b856020015183815181106103dd57fe5b602002602001015187602001518585815181106103f657fe5b602002602001015160ff168151811061040b57fe5b602002602001015161066a565b6104345760405162461bcd60e51b815260040161008190610eee565b600101610398565b8065ffffffffffff168265ffffffffffff161061046b5760405162461bcd60e51b815260040161008190610dc2565b5050565b6000805b821561048a57600019830190921691600101610473565b90505b919050565b600061049d85610511565b905060005b81518110156105095782848665ffffffffffff1603816104be57fe5b068385858585815181106104ce57fe5b602002602001015160ff160103816104e257fe5b0611156105015760405162461bcd60e51b815260040161008190610d29565b6001016104a2565b505050505050565b6060600061051e8361046f565b60ff166001600160401b038111801561053657600080fd5b50604051908082528060200260200182016040528015610560578160200160208202803683370190505b5090506000805b84156105af57600285066001141561059f5781838260ff168151811061058957fe5b60ff909216602092830291909101909101526001015b600194851c949190910190610567565b50909392505050565b60006105c2610693565b8251146105e15760405162461bcd60e51b815260040161008190610d97565b6105e9610693565b8260200151836040015184606001518560800151604051602001610611959493929190611016565b604051602081830303815290604052805190602001209050919050565b60008585858585604051602001610649959493929190610cbe565b60405160208183030381529060405280519060200120905095945050505050565b60006106768484610697565b6001600160a01b0316826001600160a01b03161490509392505050565b4690565b600080836040516020016106ab9190610c8d565b6040516020818303038152906040528051906020012090506000600182856000015186602001518760400151604051600081526020016040526040516106f49493929190610d0b565b6020604051602081039080840390855afa158015610716573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166100d75760405162461bcd60e51b815260040161008190610e8c565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b80356001600160a01b038116811461048d57600080fd5b600082601f8301126107a2578081fd5b6107b46107af83356110ec565b6110c9565b82358152602080820191908401835b853581101561096957813586016060818903601f190112156107e3578586fd5b6040518060608201106001600160401b036060830111171561080157fe5b606081016040526108146020830161077b565b81526001600160401b036040830135111561082d578687fd5b6108408960206040850135850101610984565b60208201526001600160401b036060830135111561085c578687fd5b60608201358201915088603f830112610873578687fd5b6108836107af60208401356110ec565b602083810135825281019060408401895b602086013581101561094b576080823587018e03603f190112156108b6578a8bfd5b6040518060808201106001600160401b03608083011117156108d457fe5b6080818101604090815284358901908101358352606081013560208401526108fc9101610aa2565b60408201526001600160401b0360a0843589010135111561091b578b8cfd5b6109318e8435890160a081013501604001610984565b606082015284526020938401939190910190600101610894565b505060408301525085525060209384019391909101906001016107c3565b509095945050505050565b8035801515811461048d57600080fd5b600082601f830112610994578081fd5b81356001600160401b038111156109a757fe5b6109ba601f8201601f19166020016110c9565b8181528460208386010111156109ce578283fd5b816020850160208301379081016020019190915292915050565b6000608082840312156109f9578081fd5b604051608081016001600160401b038282108183111715610a1657fe5b816040528293508435915080821115610a2e57600080fd5b610a3a86838701610792565b83526020850135915080821115610a5057600080fd5b50610a5d85828601610984565b602083015250610a6f60408401610a8c565b6040820152610a8060608401610974565b60608201525092915050565b803565ffffffffffff8116811461048d57600080fd5b803560ff8116811461048d57600080fd5b600080600060408486031215610ac7578283fd5b83356001600160401b0380821115610add578485fd5b9085019060a08288031215610af0578485fd5b90935060208501359080821115610b05578384fd5b818601915086601f830112610b18578384fd5b813581811115610b26578485fd5b8760208083028501011115610b39578485fd5b6020830194508093505050509250925092565b6000815180845260208085018081965082840281019150828601855b85811015610c35578284038952815180516001600160a01b0316855285810151606087870181905290610b9d82880182610c42565b604093840151888203898601528051808352908a0194919250898301908a810284018b018d5b82811015610c1857858203601f190184528751805183528d8101518e8401528581015160ff16868401528701516080888401819052610c0481850183610c42565b998f0199958f019593505050600101610bc3565b509e8b019e99505050948801945050506001919091019050610b68565b5091979650505050505050565b60008151808452815b81811015610c6757602081850181015186830182015201610c4b565b81811115610c785782602083870101525b50601f01601f19169290920160200192915050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b600086825260a06020830152610cd760a0830187610c42565b8281036040840152610ce98187610b4c565b65ffffffffffff95909516606084015250509015156080909101529392505050565b93845260ff9290921660208401526040830152606082015260800190565b60208082526017908201527f556e61636365707461626c652073696773206f72646572000000000000000000604082015260600190565b6020808252601c908201527f4f6e6c79206f6e6520726f756e642d726f62696e20616c6c6f77656400000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b6020808252601590820152741d1d5c9b939d5b481b9bdd081a5b98dc99585cd959605a1b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b6020808252818101527f45786365737320736967732066726f6d206f6e65207061727469636970616e74604082015260600190565b602080825260159082015274546f6f206d616e79207061727469636970616e747360581b604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601d908201527f4c61636b696e67207061727469636970616e74207369676e6174757265000000604082015260600190565b6020808252600e908201526d24b73b30b634b21039b4b3b732b960911b604082015260600190565b6020808252601490820152731f1cda59dcdf08084f481f1cda59db9959109e5f60621b604082015260600190565b60208082526018908201527f7c7369676e65645661726961626c6550617274737c213d310000000000000000604082015260600190565b60208082526017908201527f496e73756666696369656e74207369676e617475726573000000000000000000604082015260600190565b600060208252825160806020840152610fce60a0840182610b4c565b90506020840151601f19848303016040850152610feb8282610c42565b91505065ffffffffffff60408501511660608401526060840151151560808401528091505092915050565b600060a08201878352602060a08185015281885180845260c086019150828a019350845b8181101561105f5784516001600160a01b03168352938301939183019160010161103a565b505065ffffffffffff97881660408601526001600160a01b03969096166060850152505050921660809092019190915292915050565b60008235605e198336030181126110aa578182fd5b9190910192915050565b60008235607e198336030181126110aa578182fd5b6040518181016001600160401b03811182821017156110e457fe5b604052919050565b60006001600160401b038211156110ff57fe5b5060209081020190565b60006111176107af846110ec565b8381526020808201919084845b8781101561125057813587016060813603121561113f578687fd5b60408051606081016001600160401b03828210818311171561115d57fe5b9083528335908082111561116f578a8bfd5b61117b368387016109e8565b83528785013591508082111561118f578a8bfd5b9084019036601f8301126111a1578a8bfd5b81356111af6107af826110ec565b81815289810190848b01366060850287018d0111156111cc578e8ffd5b8e95505b8386101561122d57606081360312156111e7578e8ffd5b87516060810181811087821117156111fb57fe5b895261120682610aa2565b8152818d01358d8201528882013589820152835260019590950194918b01916060016111d0565b50858b015250505050918101359082015285529382019390820190600101611124565b50919695505050505050565b600060a0823603121561126d578081fd5b60405160a081016001600160401b03828210818311171561128a57fe5b81604052843583526020915081850135818111156112a6578485fd5b8501905036601f8201126112b8578384fd5b80356112c66107af826110ec565b81815283810190838501368685028601870111156112e2578788fd5b8794505b8385101561130b576112f78161077b565b8352600194909401939185019185016112e6565b508085870152505050505061132260408401610a8c565b60408201526113336060840161077b565b606082015261134460808401610a8c565b608082015292915050565b600061135b36836109e8565b9291505056fea2646970667358221220464a76b60adf074bcc5a72122061c3abb084e9433cd4202cbcd98a6b2370a87964736f6c63430007060033",
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

// LatestSupportedState is a free data retrieval call binding the contract method 0xe0ca7d93.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[],uint256)[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_ConsensusApp *ConsensusAppCaller) LatestSupportedState(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	var out []interface{}
	err := _ConsensusApp.contract.Call(opts, &out, "latestSupportedState", fixedPart, signedVariableParts)

	if err != nil {
		return *new(INitroTypesVariablePart), err
	}

	out0 := *abi.ConvertType(out[0], new(INitroTypesVariablePart)).(*INitroTypesVariablePart)

	return out0, err

}

// LatestSupportedState is a free data retrieval call binding the contract method 0xe0ca7d93.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[],uint256)[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_ConsensusApp *ConsensusAppSession) LatestSupportedState(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	return _ConsensusApp.Contract.LatestSupportedState(&_ConsensusApp.CallOpts, fixedPart, signedVariableParts)
}

// LatestSupportedState is a free data retrieval call binding the contract method 0xe0ca7d93.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[],uint256)[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_ConsensusApp *ConsensusAppCallerSession) LatestSupportedState(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	return _ConsensusApp.Contract.LatestSupportedState(&_ConsensusApp.CallOpts, fixedPart, signedVariableParts)
}

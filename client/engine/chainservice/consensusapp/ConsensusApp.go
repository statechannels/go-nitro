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
	Bin: "0x608060405234801561001057600080fd5b5061141e806100206000396000f3fe608060405234801561001057600080fd5b506004361061002b5760003560e01c8063e0ca7d9314610030575b600080fd5b61004361003e366004610b0a565b610059565b6040516100509190611039565b60405180910390f35b610061610779565b6001821461008a5760405162461bcd60e51b815260040161008190610fcb565b60405180910390fd5b6100a5610096856112e3565b6100a08486611190565b6100df565b828260008181106100b257fe5b90506020028101906100c4919061111c565b6100ce908061113b565b6100d7906113d6565b949350505050565b6020820151518151600090839060001981019081106100fa57fe5b60200260200101516000015160400151905061011682846101b4565b6000828260010165ffffffffffff168161012c57fe5b0690506000805b85518110156101ab5761015a8787838151811061014c57fe5b602002602001015185610340565b8015610185576101858287838151811061017057fe5b60200260200101516000015160400151610474565b85818151811061019157fe5b602090810291909101015151604001519150600101610133565b50505050505050565b80518083108015906101c65750600081115b6101e25760405162461bcd60e51b815260040161008190610e48565b6000826001845103815181106101f457fe5b6020026020010151600001516040015165ffffffffffff169050838160010110156102315760405162461bcd60e51b815260040161008190610f6d565b60ff8411156102525760405162461bcd60e51b815260040161008190610eb4565b60008360008151811061026157fe5b6020026020010151600001516040015165ffffffffffff16820390508481111561029d5760405162461bcd60e51b815260040161008190610db7565b6000805b85518110156103125760008682815181106102b857fe5b60200260200101516040015183169050806000146102e85760405162461bcd60e51b815260040161008190610e7f565b8682815181106102f457fe5b602002602001015160400151831792505080806001019150506102a1565b5060018660020a0381146103385760405162461bcd60e51b815260040161008190610f0e565b505050505050565b6000826020015151116103655760405162461bcd60e51b815260040161008190611002565b61037282604001516104a7565b60ff16826020015151146103985760405162461bcd60e51b815260040161008190610f9d565b6103b58260400151836000015160400151838660200151516104ca565b60006103c48360400151610541565b905060005b815181101561046d576104496103fe6103e1876105e8565b86516020810151815160408301516060909301519192909161065e565b8560200151838151811061040e57fe5b6020026020010151876020015185858151811061042757fe5b602002602001015160ff168151811061043c57fe5b602002602001015161069a565b6104655760405162461bcd60e51b815260040161008190610f45565b6001016103c9565b5050505050565b8065ffffffffffff168265ffffffffffff16106104a35760405162461bcd60e51b815260040161008190610e19565b5050565b6000805b82156104c2576000198301909216916001016104ab565b90505b919050565b60006104d585610541565b905060005b81518110156103385782848665ffffffffffff1603816104f657fe5b0683858585858151811061050657fe5b602002602001015160ff1601038161051a57fe5b0611156105395760405162461bcd60e51b815260040161008190610d80565b6001016104da565b6060600061054e836104a7565b60ff166001600160401b038111801561056657600080fd5b50604051908082528060200260200182016040528015610590578160200160208202803683370190505b5090506000805b84156105df5760028506600114156105cf5781838260ff16815181106105b957fe5b60ff909216602092830291909101909101526001015b600194851c949190910190610597565b50909392505050565b60006105f26106c3565b8251146106115760405162461bcd60e51b815260040161008190610dee565b6106196106c3565b826020015183604001518460600151856080015160405160200161064195949392919061109d565b604051602081830303815290604052805190602001209050919050565b60008585858585604051602001610679959493929190610d15565b60405160208183030381529060405280519060200120905095945050505050565b60006106a684846106c7565b6001600160a01b0316826001600160a01b03161490509392505050565b4690565b600080836040516020016106db9190610ce4565b6040516020818303038152906040528051906020012090506000600182856000015186602001518760400151604051600081526020016040526040516107249493929190610d62565b6020604051602081039080840390855afa158015610746573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166100d75760405162461bcd60e51b815260040161008190610ee3565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b80356001600160a01b03811681146104c557600080fd5b600082601f8301126107d2578081fd5b6107e46107df8335611173565b611150565b82358152602080820191908401835b85358110156109c0576060823587018803601f19011215610812578485fd5b6040518060608201106001600160401b036060830111171561083057fe5b60608101604052610846602084358901016107ab565b81526001600160401b0360408435890101351115610862578586fd5b61087888843589016040810135016020016109db565b60208201526001600160401b0360608435890101351115610897578586fd5b82358701606081013501603f0188136108ae578586fd5b6108c76107df8435890160608101350160200135611173565b83358801606081013501602081810135835282019190604001885b86358b01606081013501602001358110156109a35786358b016060810135018235016080818e03603f19011215610917578a8bfd5b6040518060808201106001600160401b036080830111171561093557fe5b60808101604052604082013581526060820135602082015261095960808301610af9565b60408201526001600160401b0360a08301351115610975578b8cfd5b6109888e604060a08501358501016109db565b606082015285525060209384019391909101906001016108e2565b5050604083015250845260209384019391909101906001016107f3565b509095945050505050565b803580151581146104c557600080fd5b600082601f8301126109eb578081fd5b81356001600160401b038111156109fe57fe5b610a11601f8201601f1916602001611150565b818152846020838601011115610a25578283fd5b816020850160208301379081016020019190915292915050565b600060808284031215610a50578081fd5b604051608081016001600160401b038282108183111715610a6d57fe5b816040528293508435915080821115610a8557600080fd5b610a91868387016107c2565b83526020850135915080821115610aa757600080fd5b50610ab4858286016109db565b602083015250610ac660408401610ae3565b6040820152610ad7606084016109cb565b60608201525092915050565b803565ffffffffffff811681146104c557600080fd5b803560ff811681146104c557600080fd5b600080600060408486031215610b1e578283fd5b83356001600160401b0380821115610b34578485fd5b9085019060a08288031215610b47578485fd5b90935060208501359080821115610b5c578384fd5b818601915086601f830112610b6f578384fd5b813581811115610b7d578485fd5b8760208083028501011115610b90578485fd5b6020830194508093505050509250925092565b6000815180845260208085018081965082840281019150828601855b85811015610c8c578284038952815180516001600160a01b0316855285810151606087870181905290610bf482880182610c99565b604093840151888203898601528051808352908a0194919250898301908a810284018b018d5b82811015610c6f57858203601f190184528751805183528d8101518e8401528581015160ff16868401528701516080888401819052610c5b81850183610c99565b998f0199958f019593505050600101610c1a565b509e8b019e99505050948801945050506001919091019050610bbf565b5091979650505050505050565b60008151808452815b81811015610cbe57602081850181015186830182015201610ca2565b81811115610ccf5782602083870101525b50601f01601f19169290920160200192915050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b600086825260a06020830152610d2e60a0830187610c99565b8281036040840152610d408187610ba3565b65ffffffffffff95909516606084015250509015156080909101529392505050565b93845260ff9290921660208401526040830152606082015260800190565b60208082526017908201527f556e61636365707461626c652073696773206f72646572000000000000000000604082015260600190565b6020808252601c908201527f4f6e6c79206f6e6520726f756e642d726f62696e20616c6c6f77656400000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b6020808252601590820152741d1d5c9b939d5b481b9bdd081a5b98dc99585cd959605a1b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b6020808252818101527f45786365737320736967732066726f6d206f6e65207061727469636970616e74604082015260600190565b602080825260159082015274546f6f206d616e79207061727469636970616e747360581b604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601d908201527f4c61636b696e67207061727469636970616e74207369676e6174757265000000604082015260600190565b6020808252600e908201526d24b73b30b634b21039b4b3b732b960911b604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b6020808252601490820152731f1cda59dcdf08084f481f1cda59db9959109e5f60621b604082015260600190565b60208082526018908201527f7c7369676e65645661726961626c6550617274737c213d310000000000000000604082015260600190565b60208082526017908201527f496e73756666696369656e74207369676e617475726573000000000000000000604082015260600190565b60006020825282516080602084015261105560a0840182610ba3565b90506020840151601f198483030160408501526110728282610c99565b91505065ffffffffffff60408501511660608401526060840151151560808401528091505092915050565b600060a08201878352602060a08185015281885180845260c086019150828a019350845b818110156110e65784516001600160a01b0316835293830193918301916001016110c1565b505065ffffffffffff97881660408601526001600160a01b03969096166060850152505050921660809092019190915292915050565b60008235605e19833603018112611131578182fd5b9190910192915050565b60008235607e19833603018112611131578182fd5b6040518181016001600160401b038111828210171561116b57fe5b604052919050565b60006001600160401b0382111561118657fe5b5060209081020190565b600061119e6107df84611173565b8381526020808201919084845b878110156112d75781358701606081360312156111c6578687fd5b60408051606081016001600160401b0382821081831117156111e457fe5b908352833590808211156111f6578a8bfd5b61120236838701610a3f565b835287850135915080821115611216578a8bfd5b9084019036601f830112611228578a8bfd5b81356112366107df82611173565b81815289810190848b01366060850287018d011115611253578e8ffd5b8e95505b838610156112b4576060813603121561126e578e8ffd5b875160608101818110878211171561128257fe5b895261128d82610af9565b8152818d01358d8201528882013589820152835260019590950194918b0191606001611257565b50858b0152505050509181013590820152855293820193908201906001016111ab565b50919695505050505050565b600060a082360312156112f4578081fd5b60405160a081016001600160401b03828210818311171561131157fe5b816040528435835260209150818501358181111561132d578485fd5b8501905036601f82011261133f578384fd5b803561134d6107df82611173565b8181528381019083850136868502860187011115611369578788fd5b8794505b838510156113925761137e816107ab565b83526001949094019391850191850161136d565b50808587015250505050506113a960408401610ae3565b60408201526113ba606084016107ab565b60608201526113cb60808401610ae3565b608082015292915050565b60006113e23683610a3f565b9291505056fea2646970667358221220a6253a0ee0a8479f1f6e7bae1386686b4cd61147123581d07ed7f35216b425ce64736f6c63430007060033",
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

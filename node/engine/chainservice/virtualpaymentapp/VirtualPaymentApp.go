// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package VirtualPaymentApp

import (
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
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
	Bin: "0x6080806040523461001657610eee908161001c8239f35b600080fdfe60c0604052600436101561001257600080fd5b6000803560e01c639936d8121461002857600080fd5b346101175760606003198181360112610113576004356001600160401b039182821161010f57608081833603011261010f576024359183831161010757366023840112156101075782600401359184831161010b573660248460051b8601011161010b5760443594851161010b5760409085360301126101075760246100b6946004019301906004016104dd565b604051938492151583526020604081850152825192836040860152825b8481106100f157505050828201840152601f01601f19168101030190f35b81810183015188820188015287955082016100d3565b8580fd5b8680fd5b8480fd5b8280fd5b80fd5b903590601e198136030182121561014f57018035906001600160401b03821161014f57602001918160051b3603831361014f57565b600080fd5b903590607e198136030182121561014f570190565b3565ffffffffffff8116810361014f5790565b608081019081106001600160401b0382111761019757604052565b634e487b7160e01b600052604160045260246000fd5b604081019081106001600160401b0382111761019757604052565b602081019081106001600160401b0382111761019757604052565b606081019081106001600160401b0382111761019757604052565b90601f801991011681019081106001600160401b0382111761019757604052565b901561023a57803590603e198136030182121561014f570190565b634e487b7160e01b600052603260045260246000fd5b6001600160401b0381116101975760051b60200190565b35906001600160a01b038216820361014f57565b9291926001600160401b03821161019757604051916102a4601f8201601f1916602001846101fe565b82948184528183011161014f578281602093846000960137010152565b9080601f8301121561014f578160206102dc9335910161027b565b90565b9092919260809384526102f181610250565b90604090610301825193846101fe565b8295818452602060a0948186520192600592831b8601958251871161014f5780945b878610610334575050505050505050565b6001600160401b03863581811161014f57830190606080838851031261014f57855192818401848110848211176104b557875261037081610267565b84528a5181013583811161014f57810187818a51031261014f57875190888201828110868211176104b55789528035600481101561014f5782528c518101359085821161014f576103c4918b5191016102c1565b8c518201528b51850152868101359083821161014f5701908751601f8301121561014f578135906103f482610250565b93610401895195866101fe565b8285528c519384808701948d1b820101948b51861161014f5781015b85811061043d575050505050508482015281528651019486510194610323565b803583811161014f578201948c601f1987825103011261014f578f958c51906104658261017c565b875181013582528d81013588518301528681013560ff8116810361014f57828f0152808f013586811161014f578f916104a49251918a519101016102c1565b86820152815285510194510161041d565b60246000634e487b7160e01b81526041600452fd5b359065ffffffffffff8216820361014f57565b9092918015610d0257600181146105265760405162461bcd60e51b815260206004820152601060248201526f0c4c2c840e0e4dedecc40d8cadccee8d60831b6044820152606490fd5b610530818561021f565b61053a838061011a565b91905060408136031261014f57604051610553816101ad565b81356001600160401b03811161014f5782019060808236031261014f576040519161057d8361017c565b80356001600160401b03811161014f57810136601f8201121561014f576105ab9036906020813591016102df565b835260208101356001600160401b03811161014f576105cd90369083016102c1565b602084015260606105e0604083016104ca565b9160408501928352013592831515840361014f57602065ffffffffffff93826001966060849501528152019401358452511603610cbd5761062360ff9151610e86565b1603610c7857600265ffffffffffff61064760406106418780610154565b01610169565b1603610c3457600180602085013560021c1603610bef576106688380610154565b602081013590601e198136030182121561014f57018035906001600160401b03821161014f57602001813603811361014f576106a591369161027b565b9160808136031261014f57604051906106bd8261017c565b80356001600160401b03811161014f57810136601f8201121561014f578035906106e682610250565b916106f460405193846101fe565b80835260208084019160051b8301019136831161014f57602001905b828210610bd757505050825260208101356001600160401b038116810361014f5761075491606091602085015261074960408201610267565b6040850152016104ca565b6060820152828051810103926080841261014f57606060405194610777866101ad565b60208301518652601f19011261014f57604051908160608101106001600160401b03606084011117610197576060820160405260408101519060ff8216820361014f5760809183526060810151602084015201516040820152602084015280516001600160401b036020830151169060018060a01b0360408401511665ffffffffffff606085015116906040519360a08501608060208701528451809152602060c0870195019060005b818110610bb85750505084600094608094829461085694604060209a015260608401528583015203601f1981018352826101fe565b83815191012086516040519085820192835260408201526040815261087a816101e3565b5190208387015190604051858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c81526108be816101e3565b5190209060ff8151169060408682015191015191604051938452868401526040830152606082015282805260015afa15610bac576000516001600160a01b0316908115610b7357516001600160a01b039061091890610e79565b511603610b2e5761095c61095261094961094361093d6109439561096497519961021f565b80610154565b8061011a565b93909580610154565b94909236916102df565b9236916102df565b90600181511480610b23575b80610b08575b80610aed575b15610aa85760206109aa6040610999836109a18361099989610e79565b510151610e79565b51015194610e79565b51015190838203918211610a925703610a4d576109c8604091610e79565b51015180516001101561023a57604001516020015103610a0857600190604051602081018181106001600160401b03821117610197576040526000815290565b60405162461bcd60e51b815260206004820152601a60248201527f426f62206e6f742061646a757374656420636f72726563746c790000000000006044820152606490fd5b60405162461bcd60e51b815260206004820152601c60248201527f416c696365206e6f742061646a757374656420636f72726563746c79000000006044820152606490fd5b634e487b7160e01b600052601160045260246000fd5b60405162461bcd60e51b815260206004820152601960248201527f6f6e6c79206e617469766520617373657420616c6c6f776564000000000000006044820152606490fd5b506001600160a01b03610aff83610e79565b5151161561097c565b506001600160a01b03610b1a82610e79565b51511615610976565b506001825114610970565b60405162461bcd60e51b815260206004820152601d60248201527f696e76616c6964207369676e617475726520666f7220766f75636865720000006044820152606490fd5b60405162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b6044820152606490fd5b6040513d6000823e3d90fd5b82516001600160a01b0316875260209687019690920191600101610821565b60208091610be484610267565b815201910190610710565b60405162461bcd60e51b815260206004820152601c60248201527f726564656d7074696f6e206e6f74207369676e656420627920426f62000000006044820152606490fd5b606460405162461bcd60e51b815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d316044820152fd5b60405162461bcd60e51b815260206004820152601e60248201527f706f737466756e642021756e616e696d6f75733b207c70726f6f667c3d3100006044820152606490fd5b60405162461bcd60e51b815260206004820152601f60248201527f6261642070726f6f665b305d2e7475726e4e756d3b207c70726f6f667c3d31006044820152606490fd5b509092915060ff610d20610d196020860135610e86565b928061011a565b9290501603610e3c5760409165ffffffffffff80610d42856106418580610154565b1615610e2d57600181610d59866106418680610154565b1614610e2d57600390610d70856106418580610154565b1614610db9576064835162461bcd60e51b815260206004820152602060248201527f6261642063616e646964617465207475726e4e756d3b207c70726f6f667c3d306044820152fd5b610dc581606092610154565b013580159081150361014f57610de95760019151610de2816101c8565b6000815290565b815162461bcd60e51b815260206004820152601e60248201527f2166696e616c3b207475726e4e756d3d33202626207c70726f6f667c3d3000006044820152606490fd5b505060019151610de2816101c8565b60405162461bcd60e51b8152602060048201526015602482015274021756e616e696d6f75733b207c70726f6f667c3d3605c1b6044820152606490fd5b80511561023a5760200190565b806000915b610e93575090565b600019810190808211610a9257169060ff809116908114610a92576001019080610e8b56fea26469706673582212200ac4f9688e843b00b74ec990595d81931e7829618b86008349d127cb7417ed3d64736f6c63430008110033",
}

// VirtualPaymentAppInstance represents a deployed instance of the VirtualPaymentApp contract.
type VirtualPaymentAppInstance struct {
	VirtualPaymentApp
	address common.Address
	backend bind.ContractBackend
}

func NewVirtualPaymentAppInstance(c *VirtualPaymentApp, address common.Address, backend bind.ContractBackend) *VirtualPaymentAppInstance {
	return &VirtualPaymentAppInstance{Db: *c, address: address, backend: backend}
}

func (i *VirtualPaymentAppInstance) Address() common.Address {
	return i.address
}

func (i *VirtualPaymentAppInstance) Backend() bind.ContractBackend {
	return i.backend
}

// VirtualPaymentApp is an auto generated Go binding around an Ethereum contract.
type VirtualPaymentApp struct {
	abi        abi.ABI
	deployCode []byte
}

// NewVirtualPaymentApp creates a new instance of VirtualPaymentApp.
func NewVirtualPaymentApp() (*VirtualPaymentApp, error) {
	parsed, err := VirtualPaymentAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	code := common.Hex2Bytes(VirtualPaymentAppMetaData.Bin)
	return &VirtualPaymentApp{abi: *parsed, deployCode: code}, nil
}

func (_VirtualPaymentApp *VirtualPaymentApp) DeployCode() []byte {
	return _VirtualPaymentApp.deployCode
}

func (_VirtualPaymentApp *VirtualPaymentApp) PackConstructor() ([]byte, error) {
	return _VirtualPaymentApp.abi.Pack("")
}

// StateIsSupported is a free data retrieval call binding the contract method 0x9936d812.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns(bool, string)
func (_VirtualPaymentApp *VirtualPaymentApp) PackStateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) ([]byte, error) {
	return _VirtualPaymentApp.abi.Pack("stateIsSupported", fixedPart, proof, candidate)
}

func (_VirtualPaymentApp *VirtualPaymentApp) UnpackStateIsSupported(data []byte) (struct {
	Arg  bool
	Arg0 string
}, error) {
	out, err := _VirtualPaymentApp.abi.Unpack("stateIsSupported", data)

	outstruct := new(struct {
		Arg  bool
		Arg0 string
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Arg = *abi.ConvertType(out[0], new(bool)).(*bool)
	outstruct.Arg0 = *abi.ConvertType(out[1], new(string)).(*string)

	return *outstruct, err

}

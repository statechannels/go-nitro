// Code generated via abigen V2 - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ConsensusApp

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

// ConsensusAppMetaData contains all meta data concerning the ConsensusApp contract.
var ConsensusAppMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint256\",\"name\":\"signedBy\",\"type\":\"uint256\"}],\"internalType\":\"structINitroTypes.RecoveredVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"stateIsSupported\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608080604052346100165761088b908161001c8239f35b600080fdfe6101006004908136101561001257600080fd5b600090813560e01c639936d8121461002957600080fd5b3461032757606091600319918383360112610327578435906001600160401b039081831161032357608085843603011261032357602490813583811161031f573660238201121561031f57808901359084821161031b57838260051b8201013681116103175760443591868311610313576040809a8436030112610313576080850185811088821117610301578a52878c01358781116102fd578801973660238a0112156102fd578c8901356100e66100e1826103a4565b61037f565b998a91808c52896020809d019160051b830101913683116102f9578a8d9101915b8383106102e1575050505086528681013588811681036102dd57610142916064918b890152610138604482016103bb565b8d890152016103d4565b8b860152876101536100e1866103a4565b809581520190868101915b8383106102b357505050506101769036908b01610439565b905161028457840151908186925b61023b575060ff90515191160361020d57845196838801928311888410176101fc575050839291959352808352815194859360018552838286015280518094860152825b8481106101e657505050828201840152601f01601f19168101030190f35b81810183015188820188015287955082016101c8565b634e487b7160e01b85526041905283fd5b86600a6064928588519362461bcd60e51b85528401528201526921756e616e696d6f757360b01b6044820152fd5b60001981019080821161027257169160ff809116908114610260576001019180610184565b634e487b7160e01b875260118a528387fd5b634e487b7160e01b885260118b528488fd5b865162461bcd60e51b8152808a01869052600a818501526907c70726f6f667c213d360b41b6044820152606490fd5b82358981116102d9578a916102ce83928b3691870101610439565b81520192019161015e565b8b80fd5b8a80fd5b81906102ec846103bb565b8152019101908c90610107565b8d80fd5b8980fd5b634e487b7160e01b8a5260418d52868afd5b8880fd5b8780fd5b8680fd5b8580fd5b8380fd5b5080fd5b60405190604082018281106001600160401b0382111761034a57604052565b634e487b7160e01b600052604160045260246000fd5b60405190608082018281106001600160401b0382111761034a57604052565b6040519190601f01601f191682016001600160401b0381118382101761034a57604052565b6001600160401b03811161034a5760051b60200190565b35906001600160a01b03821682036103cf57565b600080fd5b359065ffffffffffff821682036103cf57565b81601f820112156103cf578035906001600160401b03821161034a57610416601f8301601f191660200161037f565b92828452602083830101116103cf57816000926020809301838601378301015290565b604091816080528060a05203126103cf5761045261032b565b60805180358060c0526001600160401b0381116103cf576080910160a05103126103cf5761047e610360565b9060c05160805101926001600160401b038435116103cf5760a0518435850190601f820112156103cf5735916104b66100e1846103a4565b91602083858152019460a051873560c05160805101019060208760051b830101116103cf57939591949360200192905b873560c05160805101019160208660051b8401018510156107dc576001600160401b038535116103cf576060601f198635850160a0510301126103cf57604051928360608101106001600160401b036060860111176107c75760206105559160608601604052873501016103bb565b83528435893560c0516080510101016040810135906001600160401b0382116103cf57604091601f19910160a0510301126103cf5761059261032b565b85358a3560c0516080510101016004602060408301358301013510156103cf578060406020920135010135815285358a3560c0516080510101016001600160401b0360408083013583010135116103cf5760a0516106009160408082013590910190810135016020016103e7565b602082015260208401528435893560c05160805101010160e0526001600160401b03606060e0510135116103cf5760a05160e0516060810135019690603f880112156103cf576106566100e160208901356103a4565b9860208a818a01358152019260a051606060e05101358d8a35903560c0516080510101010190604060208c013560051b830101116103cf57604001935b606060e05101358d8a35903560c0516080510101010190604060208c013560051b8301018610156107a5578535916001600160401b0383116103cf5760a05160809184019003603f1901126103cf578160808f938c6106f0610360565b95604083606060e05101358435843560c05189510101010101013587526060838160e05101358435843560c0518951010101010101356020880152606060e05101359135903560c05185510101010101013560ff811681036103cf578f906040850152606060e0510135908c35903560c05160805101010101019060a0820135926001600160401b0384116103cf57610795602094936040869560a0519201016103e7565b6060820152815201940193610693565b5050604086019a909a52938552929792955050602093840193909201916104e6565b60246000634e487b7160e01b81526041600452fd5b975094959350505050815260c051608051016020810135906001600160401b0382116103cf576108109160a05191016103e7565b6020820152610827604060c05160805101016103d4565b6040820152606060c05160805101013580151581036103cf576060820152825260206080510135602083015256fea26469706673582212202ebfc8d244add2e41388e3e2284dde760225eee0448710e6aa2574fec676211c64736f6c63430008110033",
}

// ConsensusAppInstance represents a deployed instance of the ConsensusApp contract.
type ConsensusAppInstance struct {
	ConsensusApp
	address common.Address
	backend bind.ContractBackend
}

func NewConsensusAppInstance(c *ConsensusApp, address common.Address, backend bind.ContractBackend) *ConsensusAppInstance {
	return &ConsensusAppInstance{Db: *c, address: address, backend: backend}
}

func (i *ConsensusAppInstance) Address() common.Address {
	return i.address
}

func (i *ConsensusAppInstance) Backend() bind.ContractBackend {
	return i.backend
}

// ConsensusApp is an auto generated Go binding around an Ethereum contract.
type ConsensusApp struct {
	abi        abi.ABI
	deployCode []byte
}

// NewConsensusApp creates a new instance of ConsensusApp.
func NewConsensusApp() (*ConsensusApp, error) {
	parsed, err := ConsensusAppMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	code := common.Hex2Bytes(ConsensusAppMetaData.Bin)
	return &ConsensusApp{abi: *parsed, deployCode: code}, nil
}

func (_ConsensusApp *ConsensusApp) DeployCode() []byte {
	return _ConsensusApp.deployCode
}

func (_ConsensusApp *ConsensusApp) PackConstructor() ([]byte, error) {
	return _ConsensusApp.abi.Pack("")
}

// StateIsSupported is a free data retrieval call binding the contract method 0x9936d812.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256)[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),uint256) candidate) pure returns(bool, string)
func (_ConsensusApp *ConsensusApp) PackStateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesRecoveredVariablePart, candidate INitroTypesRecoveredVariablePart) ([]byte, error) {
	return _ConsensusApp.abi.Pack("stateIsSupported", fixedPart, proof, candidate)
}

func (_ConsensusApp *ConsensusApp) UnpackStateIsSupported(data []byte) (struct {
	Arg  bool
	Arg0 string
}, error) {
	out, err := _ConsensusApp.abi.Unpack("stateIsSupported", data)

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

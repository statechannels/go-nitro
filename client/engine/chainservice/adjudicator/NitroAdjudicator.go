// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.
// generated like this:
// cd nitro-protocol
// solc @statechannels/exit-format/=$(pwd)/node_modules/@statechannels/exit-format/ @openzeppelin/contracts/=$(pwd)/node_modules/@openzeppelin/contracts/ ./contracts/NitroAdjudicator.sol --optimize --bin -o build
// abigen --abi=./build/NitroAdjudicator.abi --bin=./build/NitroAdjudicator.bin --pkg=NitroAdjudicator --out=NitroAdjudicator.go

package NitroAdjudicator

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

// IForceMoveAppVariablePart is an auto generated low-level Go binding around an user-defined struct.
type IForceMoveAppVariablePart struct {
	Outcome []byte
	AppData []byte
}

// IForceMoveFixedPart is an auto generated low-level Go binding around an user-defined struct.
type IForceMoveFixedPart struct {
	ChainId           *big.Int
	Participants      []common.Address
	ChannelNonce      *big.Int
	AppDefinition     common.Address
	ChallengeDuration *big.Int
}

// IForceMoveSignature is an auto generated low-level Go binding around an user-defined struct.
type IForceMoveSignature struct {
	V uint8
	R [32]byte
	S [32]byte
}

// IMultiAssetHolderClaimArgs is an auto generated low-level Go binding around an user-defined struct.
type IMultiAssetHolderClaimArgs struct {
	SourceChannelId                 [32]byte
	SourceStateHash                 [32]byte
	SourceOutcomeBytes              []byte
	SourceAssetIndex                *big.Int
	IndexOfTargetInSource           *big.Int
	TargetStateHash                 [32]byte
	TargetOutcomeBytes              []byte
	TargetAssetIndex                *big.Int
	TargetAllocationIndicesToPayout []*big.Int
}

// NitroAdjudicatorMetaData contains all meta data concerning the NitroAdjudicator contract.
var NitroAdjudicatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"name\":\"compute_claim_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newSourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newTargetAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numParticipants\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numSigs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numWhoSignedWhats\",\"type\":\"uint256\"}],\"name\":\"requireValidInput\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"variablePartAB\",\"type\":\"tuple[2]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"sig\",\"type\":\"tuple\"}],\"name\":\"respond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nParticipants\",\"type\":\"uint256\"},{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"ab\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint48\",\"name\":\"turnNumB\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"}],\"name\":\"validTransition\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614a46806100206000396000f3fe6080604052600436106100fe5760003560e01c80636775b17311610095578063c7df14e211610064578063c7df14e2146102af578063cfaa7978146102cf578063da4cdf73146102ef578063e29cffe01461030f578063f198bea91461033f576100fe565b80636775b173146102225780636c3655221461024f578063732f92081461026f578063be5c2a311461028f576100fe565b80633033730e116100d15780633033730e1461019e57806330776841146101be578063552cfa50146101de578063564b81ef1461020d576100fe565b80630149b7621461010357806311e9f17814610125578063166e56cd1461015e5780632fb1d2701461018b575b600080fd5b34801561010f57600080fd5b5061012361011e3660046137ec565b61035f565b005b34801561013157600080fd5b50610145610140366004613aa0565b6103b7565b6040516101559493929190613fa6565b60405180910390f35b34801561016a57600080fd5b5061017e6101793660046133f9565b610705565b6040516101559190613ff1565b610123610199366004613424565b610722565b3480156101aa57600080fd5b506101236101b9366004613b09565b6109b0565b3480156101ca57600080fd5b506101236101d93660046136f8565b610a30565b3480156101ea57600080fd5b506101fe6101f9366004613693565b610b10565b60405161015593929190614906565b34801561021957600080fd5b5061017e610b2b565b34801561022e57600080fd5b5061024261023d36600461398e565b610b2f565b6040516101559190613fe6565b34801561025b57600080fd5b5061012361026a3660046136ab565b610b4a565b34801561027b57600080fd5b5061012361028a366004613b81565b610f0f565b34801561029b57600080fd5b506102426102aa366004613b50565b610f2e565b3480156102bb57600080fd5b5061017e6102ca366004613693565b610fb3565b3480156102db57600080fd5b506101236102ea366004613b81565b610fc5565b3480156102fb57600080fd5b5061012361030a36600461345e565b610fd4565b34801561031b57600080fd5b5061032f61032a366004613a09565b611115565b6040516101559493929190613f5b565b34801561034b57600080fd5b5061012361035a3660046138b3565b61168c565b610373866020015151855184518451610f2e565b50600061037f876117f8565b905061038a8161186e565b61039481876118a5565b6103a3868686848b88886118e8565b506103ae818761194e565b50505050505050565b6060600060606000808551116103ce5785516103d1565b84515b6001600160401b03811180156103e657600080fd5b5060405190808252806020026020018201604052801561042057816020015b61040d612ceb565b8152602001906001900390816104055790505b5091506000905085516001600160401b038111801561043e57600080fd5b5060405190808252806020026020018201604052801561047857816020015b610465612ceb565b81526020019060019003908161045d5790505b50935060019250866000805b88518110156106f95788818151811061049957fe5b6020026020010151600001518782815181106104b157fe5b602002602001015160000181815250508881815181106104cd57fe5b6020026020010151604001518782815181106104e557fe5b60200260200101516040019060ff16908160ff168152505088818151811061050957fe5b60200260200101516060015187828151811061052157fe5b60200260200101516060018190525060006105538a838151811061054157fe5b602002602001015160200151856119d5565b9050885160001480610582575088518310801561058257508189848151811061057857fe5b6020026020010151145b1561069457600260ff168a848151811061059857fe5b60200260200101516040015160ff1614156105ce5760405162461bcd60e51b81526004016105c590614347565b60405180910390fd5b808a83815181106105db57fe5b602002602001015160200151038883815181106105f457fe5b6020026020010151602001818152505060405180608001604052808b848151811061061b57fe5b60200260200101516000015181526020018281526020018b848151811061063e57fe5b60200260200101516040015160ff1681526020018b848151811061065e57fe5b60200260200101516060015181525086848151811061067957fe5b602002602001018190525080850194508260010192506106c9565b8982815181106106a057fe5b6020026020010151602001518883815181106106b857fe5b602002602001015160200181815250505b8782815181106106d557fe5b6020026020010151602001516000146106ed57600096505b90920391600101610484565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b61072b836119ef565b156107485760405162461bcd60e51b81526004016105c59061458c565b6001600160a01b03841660009081526001602090815260408083208684529091528120548381101561078c5760405162461bcd60e51b81526004016105c5906140bb565b61079684846119fb565b81106107b45760405162461bcd60e51b81526004016105c590614691565b6107c8816107c286866119fb565b90611a55565b91506001600160a01b0386166107fc578234146107f75760405162461bcd60e51b81526004016105c5906146ff565b61089a565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd9061082c90339030908790600401613efd565b602060405180830381600087803b15801561084657600080fd5b505af115801561085a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061087e9190613677565b61089a5760405162461bcd60e51b81526004016105c59061440f565b60006108a682846119fb565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a71590610905908a9087908690613f3a565b60405180910390a26001600160a01b0387166103ae5760006109278585611a55565b90506000336001600160a01b03168260405161094290612b60565b60006040518083038185875af1925050503d806000811461097f576040519150601f19603f3d011682016040523d82523d6000602084013e610984565b606091505b50509050806109a55760405162461bcd60e51b81526004016105c590614762565b505050505050505050565b60606000806109c2888589888a611ab2565b92509250925060608060006109ef84878d815181106109dd57fe5b602002602001015160400151896103b7565b93509350509250610a068b868c8b8a888a88611b30565b610a23868c81518110610a1557fe5b602002602001015183611bdc565b5050505050505050505050565b606080600080610a3f85611c14565b9350935093509350606080606060006060888a6060015181518110610a6057fe5b60200260200101516040015190506060888b60e0015181518110610a8057fe5b6020026020010151604001519050610aa48783838e608001518f6101000151611115565b809650819750829850839950505050505050610afd89878a878c8e6060015181518110610acd57fe5b6020026020010151604001518e6080015181518110610ae857fe5b6020026020010151600001518c898c89611dbd565b6109a5878a60e0015181518110610a1557fe5b6000806000610b1e84611ec0565b9196909550909350915050565b4690565b6000610b3e8686868686611ede565b90505b95945050505050565b610b5383611fb1565b610b6581838051906020012085611fe4565b60016060610b7284612033565b9050606081516001600160401b0381118015610b8d57600080fd5b50604051908082528060200260200182016040528015610bc757816020015b610bb4612d11565b815260200190600190039081610bac5790505b509050606082516001600160401b0381118015610be357600080fd5b50604051908082528060200260200182016040528015610c0d578160200160208202803683370190505b509050606083516001600160401b0381118015610c2957600080fd5b50604051908082528060200260200182016040528015610c53578160200160208202803683370190505b50905060005b8451811015610df557610c6a612d11565b858281518110610c7657fe5b602002602001015190506060816040015190506000878481518110610c9757fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008d815260200190815260200160002054868581518110610ce857fe5b6020026020010181815250506060600060606000610d5b8a8981518110610d0b57fe5b60200260200101518760006001600160401b0381118015610d2b57600080fd5b50604051908082528060200260200182016040528015610d55578160200160208202803683370190505b506103b7565b935093509350935082610d6d5760009c505b80898981518110610d7a57fe5b602002602001018181525050838c8981518110610d9357fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b8981518110610dd657fe5b6020026020010181905250505050505050508080600101915050610c59565b5060005b8451811015610ea9576000858281518110610e1057fe5b6020026020010151600001519050828281518110610e2a57fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208e835290935291909120805491909103905583518a906000805160206149f1833981519152908490879082908110610e8257fe5b6020026020010151604051610e98929190614825565b60405180910390a250600101610df9565b508415610ec457600088815260208190526040812055610efc565b600087604051602001610ed79190614071565b604051602081830303815290604052805190602001209050610efa898883612049565b505b610f05836120af565b5050505050505050565b6000610f20888888888888886120df565b9050610f0581866000610b4a565b6000838510158015610f405750600084115b610f5c5760405162461bcd60e51b81526004016105c5906142a9565b8483148015610f6a57508482145b610f865760405162461bcd60e51b81526004016105c590614446565b60ff851115610fa75760405162461bcd60e51b81526004016105c5906142e0565b5060015b949350505050565b60006020819052908152604090205481565b610f05878787878787876120df565b6000610fdf846117f8565b9050600080610fed83611ec0565b50865160208101519051929450909250600091611014918691868c865b60200201516122a9565b6020878101519081015190519192506000916110399187916001888101908e9061100a565b87515180516020918201206040805160808101825265ffffffffffff808a16825288169381019390935282018590526060820181905291925061107c90876122ed565b6020890151805165ffffffffffff60018801168161109657fe5b06815181106110a157fe5b60200260200101516001600160a01b03166110bc8389612300565b6001600160a01b0316146110e25760405162461bcd60e51b81526004016105c5906143a9565b6110fb8960200151518b8a886001018d60600151611ede565b50611109868660010161194e565b50505050505050505050565b606080606060008088516001600160401b038111801561113457600080fd5b5060405190808252806020026020018201604052801561116e57816020015b61115b612ceb565b8152602001906001900390816111535790505b50945087516001600160401b038111801561118857600080fd5b506040519080825280602002602001820160405280156111c257816020015b6111af612ceb565b8152602001906001900390816111a75790505b50935087516001600160401b03811180156111dc57600080fd5b5060405190808252806020026020018201604052801561121657816020015b611203612ceb565b8152602001906001900390816111fb5790505b50925060005b89518110156113025789818151811061123157fe5b60200260200101516000015186828151811061124957fe5b6020026020010151600001818152505089818151811061126557fe5b60200260200101516020015186828151811061127d57fe5b6020026020010151602001818152505089818151811061129957fe5b6020026020010151606001518682815181106112b157fe5b6020026020010151606001819052508981815181106112cc57fe5b6020026020010151604001518682815181106112e457fe5b602090810291909101015160ff90911660409091015260010161121c565b5060005b88518110156114ad5788818151811061131b57fe5b60200260200101516000015185828151811061133357fe5b6020026020010151600001818152505088818151811061134f57fe5b60200260200101516020015185828151811061136757fe5b6020026020010151602001818152505088818151811061138357fe5b60200260200101516060015185828151811061139b57fe5b6020026020010151606001819052508881815181106113b657fe5b6020026020010151604001518582815181106113ce57fe5b60200260200101516040019060ff16908160ff16815250508881815181106113f257fe5b60200260200101516000015184828151811061140a57fe5b60200260200101516000018181525050600084828151811061142857fe5b6020026020010151602001818152505088818151811061144457fe5b60200260200101516060015184828151811061145c57fe5b60200260200101516060018190525088818151811061147757fe5b60200260200101516040015184828151811061148f57fe5b602090810291909101015160ff909116604090910152600101611306565b508960005b888110156114f457816114c4576114f4565b60006114e78c83815181106114d557fe5b602002602001015160200151846119d5565b90920391506001016114b2565b506000611518828c8b8151811061150757fe5b6020026020010151602001516119d5565b9050606061153c8c8b8151811061152b57fe5b60200260200101516060015161236e565b905060005b815181101561167b57826115545761167b565b60005b8851811015611672578361156a57611672565b88818151811061157657fe5b60200260200101516000015183838151811061158e57fe5b6020026020010151141561166a5760006115bf8e83815181106115ad57fe5b602002602001015160200151866119d5565b905080850394508b51600014806115f357508b51871080156115f35750818c88815181106115e957fe5b6020026020010151145b1561166457808a838151811061160557fe5b60200260200101516020018181510391508181525050808b8e8151811061162857fe5b602002602001015160200181815103915081815250508089838151811061164b57fe5b6020908102919091018101510152968701966001909601955b50611672565b600101611557565b50600101611541565b505050505095509550955095915050565b6116a0876020015151865185518551610f2e565b5060006116ac886117f8565b905060006116b982612384565b60028111156116c457fe5b14156116d9576116d481886123ce565b611708565b60016116e482612384565b60028111156116ef57fe5b14156116ff576116d481886118a5565b6117088161186e565b6000611719888888858d8a8a6118e8565b905061172a818a602001518561240d565b817ff6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee898b60800151420160008a60ff16118d8c8b8b6040516117729796959493929190614846565b60405180910390a26117db60405180608001604052808a65ffffffffffff1681526020018b60800151420165ffffffffffff1681526020018381526020018960018b5103815181106117c057fe5b60200260200101516000015180519060200120815250612467565b600092835260208390526040909220919091555050505050505050565b6000611802610b2b565b8251146118215760405162461bcd60e51b81526004016105c5906140f2565b611829610b2b565b82602001518360400151846060015185608001516040516020016118519594939291906147da565b604051602081830303815290604052805190602001209050919050565b600261187982612384565b600281111561188457fe5b14156118a25760405162461bcd60e51b81526004016105c59061411d565b50565b60006118b083611ec0565b505090508065ffffffffffff168265ffffffffffff16116118e35760405162461bcd60e51b81526004016105c590614084565b505050565b600060606118f989898989896124b8565b905061190c898660200151838787612681565b6119285760405162461bcd60e51b81526004016105c59061427d565b8060018251038151811061193857fe5b6020026020010151915050979650505050505050565b6040805160808101825265ffffffffffff83168152600060208201819052918101829052606081019190915261198390612467565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0826040516119c99190614833565b60405180910390a25050565b60008183116119e457826119e6565b815b90505b92915050565b60a081901c155b919050565b6000828201838110156119e6576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611aac576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6060600080611ac087612757565b611ac986611fb1565b611adb85858051906020012088611fe4565b611ae484612033565b9250828881518110611af257fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a908110611b6a57fe5b602002602001015160400181905250611baa868686604051602001611b8f9190613fd3565b60405160208183030381529060405280519060200120612049565b856000805160206149f18339815191528984604051611bca929190614825565b60405180910390a25050505050505050565b611c10604051806060016040528084600001516001600160a01b0316815260200184602001518152602001838152506127b6565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611c4890612757565b611c5185611fb1565b611c678a60200151858051906020012087611fe4565b611c7084612033565b9850611c7b82612033565b9750888381518110611c8957fe5b6020908102919091010151519650600260ff16898481518110611ca857fe5b6020026020010151604001518281518110611cbf57fe5b60200260200101516040015160ff1614611ceb5760405162461bcd60e51b81526004016105c5906146c8565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a9085908110611d2057fe5b6020026020010151604001518b6080015181518110611d3b57fe5b6020026020010151600001519050876001600160a01b0316898381518110611d5f57fe5b6020026020010151600001516001600160a01b031614611d915760405162461bcd60e51b81526004016105c590614310565b611d9a81611fb1565b611db08b60a00151848051906020012083611fe4565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b9084908110611e0357fe5b602002602001015160400181905250611e2c838d602001518c604051602001611b8f9190613fd3565b85878281518110611e3957fe5b602002602001015160400181905250611e62888d60a0015189604051602001611b8f9190613fd3565b826000805160206149f18339815191528387604051611e82929190614825565b60405180910390a2876000805160206149f18339815191528287604051611eaa929190614825565b60405180910390a2505050505050505050505050565b60009081526020819052604090205460d081901c9160a082901c9190565b600080611eed87878787612862565b90506001816001811115611efd57fe5b1415611fa4578451602086015160405163fd7a2f6560e01b81526001600160a01b0386169263fd7a2f6592611f389289908d90600401614799565b60206040518083038186803b158015611f5057600080fd5b505afa158015611f64573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f889190613677565b611fa45760405162461bcd60e51b81526004016105c590614180565b5060019695505050505050565b6002611fbc82612384565b6002811115611fc757fe5b146118a25760405162461bcd60e51b81526004016105c5906145f3565b6000611fef82611ec0565b92505050611ffd8484612946565b6001600160a01b0316816001600160a01b03161461202d5760405162461bcd60e51b81526004016105c5906143e0565b50505050565b6060818060200190518101906119e99190613570565b60008061205585611ec0565b5091509150600061209560405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250612467565b600096875260208790526040909620959095555050505050565b60005b8151811015611c10576120d78282815181106120ca57fe5b60200260200101516127b6565b6001016120b2565b60006120ea876117f8565b90506120f58161186e565b61210b8760200151518560ff1684518651610f2e565b508360ff168860010165ffffffffffff16101561213a5760405162461bcd60e51b81526004016105c5906145c3565b60608460ff166001600160401b038111801561215557600080fd5b5060405190808252806020026020018201604052801561217f578160200160208202803683370190505b50905060005b8560ff168165ffffffffffff1610156121d7576121b08389898960ff16856001018f010360016122a9565b828265ffffffffffff16815181106121c457fe5b6020908102919091010152600101612185565b506121e9898960200151838688612681565b6122055760405162461bcd60e51b81526004016105c590614555565b60008680519060200120905061224e6040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b815260200183815250612467565b60008085815260200190815260200160002081905550827f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516122949190614833565b60405180910390a25050979650505050505050565b600085856122b686612033565b85856040516020016122cc959493929190613ffa565b60405160208183030381529060405280519060200120905095945050505050565b6122f78282612972565b611c10816129a8565b600080836040516020016123149190613ecc565b6040516020818303038152906040528051906020012090506000612346828560000151866020015187604001516129db565b90506001600160a01b038116610fab5760405162461bcd60e51b81526004016105c59061437e565b6060818060200190518101906119e991906134e1565b60008061239083611ec0565b5091505065ffffffffffff81166123ab5760009150506119f6565b428165ffffffffffff16116123c45760029150506119f6565b60019150506119f6565b60006123d983611ec0565b505090508065ffffffffffff168265ffffffffffff1610156118e35760405162461bcd60e51b81526004016105c590614149565b600061243f846040516020016124239190614047565b6040516020818303038152906040528051906020012083612300565b905061244b8184612a80565b61202d5760405162461bcd60e51b81526004016105c590614216565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b161792916124a491612946565b6001600160a01b0316919091179392505050565b60608085516001600160401b03811180156124d257600080fd5b506040519080825280602002602001820160405280156124fc578160200160208202803683370190505b509050600160ff86168803016000805b88518165ffffffffffff161015612673578089518b03600101019150612586878a8365ffffffffffff168151811061254057fe5b6020026020010151602001518b8465ffffffffffff168151811061256057fe5b602002602001015160000151858765ffffffffffff168765ffffffffffff1610156122a9565b848265ffffffffffff168151811061259a57fe5b6020026020010181815250508965ffffffffffff168265ffffffffffff16101561266b5761266986602001515160405180604001604052808665ffffffffffff168665ffffffffffff1610151515151581526020018665ffffffffffff168660010165ffffffffffff1610151515151581525060405180604001604052808d8665ffffffffffff168151811061262c57fe5b602002602001015181526020018d8660010165ffffffffffff168151811061265057fe5b6020026020010151815250856001018a60600151611ede565b505b60010161250c565b509198975050505050505050565b835183516000919061269584898484612ad6565b6126b15760405162461bcd60e51b81526004016105c59061447b565b60005b82811015612748576000612704888784815181106126ce57fe5b602002602001015160ff16815181106126e357fe5b60200260200101518884815181106126f757fe5b6020026020010151612300565b905088828151811061271257fe5b60200260200101516001600160a01b0316816001600160a01b03161461273f576000945050505050610b41565b506001016126b4565b50600198975050505050505050565b60005b8151816001011015611c105781816001018151811061277557fe5b602002602001015182828151811061278957fe5b6020026020010151106127ae5760405162461bcd60e51b81526004016105c5906141e6565b60010161275a565b805160005b8260400151518110156118e3576000836040015182815181106127da57fe5b60200260200101516000015190506000846040015183815181106127fa57fe5b6020026020010151602001519050612811826119ef565b1561282e576128298461282384612b60565b83612b63565b612858565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b50506001016127bb565b6020830151600090156128a5576128848360015b602002015151845151612c73565b6128a05760405162461bcd60e51b81526004016105c590614623565b61293b565b8351156128c45760405162461bcd60e51b81526004016105c590614736565b846002028265ffffffffffff161015612933576128e2836001612876565b6128fe5760405162461bcd60e51b81526004016105c5906144b0565b6020838101518101518451909101516129179190612c73565b6128a05760405162461bcd60e51b81526004016105c59061451e565b506001610fab565b506000949350505050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60008181526020819052604090205461298c908390612cd7565b611c105760405162461bcd60e51b81526004016105c5906144e7565b60016129b382612384565b60028111156129be57fe5b146118a25760405162461bcd60e51b81526004016105c5906141b7565b6000601b8460ff1610156129f057601b840193505b8360ff16601b14158015612a0857508360ff16601c14155b15612a1557506000610fab565b60018585858560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa158015612a6f573d6000803e3d6000fd5b505050602060405103519050610fab565b6000805b8251811015612acc57828181518110612a9957fe5b60200260200101516001600160a01b0316846001600160a01b03161415612ac45760019150506119e9565b600101612a84565b5060009392505050565b600082855114612af85760405162461bcd60e51b81526004016105c59061465a565b60005b83811015612b5457600084828765ffffffffffff1687010381612b1a57fe5b0690508381888481518110612b2b57fe5b602002602001015160ff16016001011015612b4b57600092505050610fab565b50600101612afb565b50600195945050505050565b90565b6001600160a01b038316612bf3576000826001600160a01b031682604051612b8a90612b60565b60006040518083038185875af1925050503d8060008114612bc7576040519150601f19603f3d011682016040523d82523d6000602084013e612bcc565b606091505b5050905080612bed5760405162461bcd60e51b81526004016105c59061424d565b506118e3565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb90612c219085908590600401613f21565b602060405180830381600087803b158015612c3b57600080fd5b505af1158015612c4f573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061202d9190613677565b815181516000916001918114808314612c8f5760009250612ccd565b600160208701838101602088015b600284838510011415612cc8578051835114612cbc5760009650600093505b60209283019201612c9d565b505050505b5090949350505050565b600081612ce384612467565b149392505050565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b80356119f6816149be565b600082601f830112612d56578081fd5b604051604081018181106001600160401b0382111715612d7257fe5b80604052508091508284604085011115612d8b57600080fd5b60005b6002811015612db7578135612da2816149d3565b83526020928301929190910190600101612d8e565b50505092915050565b600082601f830112612dd0578081fd5b8135612de3612dde82614954565b614931565b818152915060208083019084810160005b84811015612e8f5781358701608080601f19838c03011215612e1557600080fd5b604080518281016001600160401b038282108183111715612e3257fe5b9083528488013582528483013582890152606090612e518287016133e3565b83850152938501359380851115612e6757600080fd5b50612e768d8986880101613156565b9082015287525050509282019290820190600101612df4565b505050505092915050565b600082601f830112612eaa578081fd5b8151612eb8612dde82614954565b818152915060208083019084810160005b84811015612e8f5781518701608080601f19838c03011215612eea57600080fd5b604080518281016001600160401b038282108183111715612f0757fe5b9083528488015182528483015182890152606090612f268287016133ee565b83850152938501519380851115612f3c57600080fd5b50612f4b8d89868801016131a4565b9082015287525050509282019290820190600101612ec9565b600082601f830112612f74578081fd5b8135612f82612dde82614954565b8181529150602080830190848101606080850287018301881015612fa557600080fd5b60005b85811015612fcc57612fba89846132ee565b85529383019391810191600101612fa8565b50505050505092915050565b600082601f830112612fe8578081fd5b604051604081018181106001600160401b038211171561300457fe5b6040529050808260005b6002811015612db757613024868335870161334b565b8352602092830192919091019060010161300e565b600082601f830112613049578081fd5b8135613057612dde82614954565b818152915060208083019084810160005b84811015612e8f5761307f888484358a010161334b565b84529282019290820190600101613068565b600082601f8301126130a1578081fd5b81356130af612dde82614954565b8181529150602080830190848101818402860182018710156130d057600080fd5b60005b84811015612e8f578135845292820192908201906001016130d3565b600082601f8301126130ff578081fd5b813561310d612dde82614954565b81815291506020808301908481018184028601820187101561312e57600080fd5b60005b84811015612e8f578135613144816149e1565b84529282019290820190600101613131565b600082601f830112613166578081fd5b8135613174612dde82614971565b915080825283602082850101111561318b57600080fd5b8060208401602084013760009082016020015292915050565b600082601f8301126131b4578081fd5b81516131c2612dde82614971565b91508082528360208285010111156131d957600080fd5b6131ea816020840160208601614992565b5092915050565b600060a08284031215613202578081fd5b60405160a081016001600160401b03828210818311171561321f57fe5b816040528293508435835260209150818501358181111561323f57600080fd5b85019050601f8101861361325257600080fd5b8035613260612dde82614954565b81815283810190838501858402850186018a101561327d57600080fd5b600094505b838510156132a9578035613295816149be565b835260019490940193918501918501613282565b50808587015250505050506132c0604084016133cd565b60408201526132d160608401612d3b565b60608201526132e2608084016133cd565b60808201525092915050565b6000606082840312156132ff578081fd5b604051606081018181106001600160401b038211171561331b57fe5b604052905080823561332c816149e1565b8082525060208301356020820152604083013560408201525092915050565b60006040828403121561335c578081fd5b604051604081016001600160401b03828210818311171561337957fe5b81604052829350843591508082111561339157600080fd5b61339d86838701613156565b835260208501359150808211156133b357600080fd5b506133c085828601613156565b6020830152505092915050565b803565ffffffffffff811681146119f657600080fd5b80356119f6816149e1565b80516119f6816149e1565b6000806040838503121561340b578182fd5b8235613416816149be565b946020939093013593505050565b60008060008060808587031215613439578182fd5b8435613444816149be565b966020860135965060408601359560600135945092505050565b60008060008060e08587031215613473578182fd5b61347d8686612d46565b935060408501356001600160401b0380821115613498578384fd5b6134a4888389016131f1565b945060608701359150808211156134b9578384fd5b506134c687828801612fd8565b9250506134d686608087016132ee565b905092959194509250565b600060208083850312156134f3578182fd5b82516001600160401b03811115613508578283fd5b8301601f81018513613518578283fd5b8051613526612dde82614954565b8181528381019083850185840285018601891015613542578687fd5b8694505b83851015613564578051835260019490940193918501918501613546565b50979650505050505050565b60006020808385031215613582578182fd5b82516001600160401b0380821115613598578384fd5b818501915085601f8301126135ab578384fd5b81516135b9612dde82614954565b81815284810190848601875b8481101561366857815187016060818d03601f190112156135e457898afd5b60408051606081018181108a821117156135fa57fe5b8252828b0151613609816149be565b8152828201518981111561361b578c8dfd5b6136298f8d838701016131a4565b8c8301525060608301518981111561363f578c8dfd5b61364d8f8d83870101612e9a565b928201929092528652505092870192908701906001016135c5565b50909998505050505050505050565b600060208284031215613688578081fd5b81516119e6816149d3565b6000602082840312156136a4578081fd5b5035919050565b6000806000606084860312156136bf578081fd5b8335925060208401356001600160401b038111156136db578182fd5b6136e786828701613156565b925050604084013590509250925092565b600060208284031215613709578081fd5b81356001600160401b038082111561371f578283fd5b8184019150610120808387031215613735578384fd5b61373e81614931565b9050823581526020830135602082015260408301358281111561375f578485fd5b61376b87828601613156565b604083015250606083013560608201526080830135608082015260a083013560a082015260c0830135828111156137a0578485fd5b6137ac87828601613156565b60c08301525060e083013560e082015261010080840135838111156137cf578586fd5b6137db88828701613091565b918301919091525095945050505050565b60008060008060008060c08789031215613804578384fd5b86356001600160401b038082111561381a578586fd5b6138268a838b016131f1565b975061383460208a016133cd565b96506040890135915080821115613849578586fd5b6138558a838b01613039565b955061386360608a016133e3565b94506080890135915080821115613878578384fd5b6138848a838b01612f64565b935060a0890135915080821115613899578283fd5b506138a689828a016130ef565b9150509295509295509295565b6000806000806000806000610120888a0312156138ce578485fd5b87356001600160401b03808211156138e4578687fd5b6138f08b838c016131f1565b98506138fe60208b016133cd565b975060408a0135915080821115613913578687fd5b61391f8b838c01613039565b965061392d60608b016133e3565b955060808a0135915080821115613942578283fd5b61394e8b838c01612f64565b945060a08a0135915080821115613963578283fd5b506139708a828b016130ef565b9250506139808960c08a016132ee565b905092959891949750929550565b600080600080600060c086880312156139a5578283fd5b853594506139b68760208801612d46565b935060608601356001600160401b038111156139d0578384fd5b6139dc88828901612fd8565b9350506139eb608087016133cd565b915060a08601356139fb816149be565b809150509295509295909350565b600080600080600060a08688031215613a20578283fd5b8535945060208601356001600160401b0380821115613a3d578485fd5b613a4989838a01612dc0565b95506040880135915080821115613a5e578485fd5b613a6a89838a01612dc0565b9450606088013593506080880135915080821115613a86578283fd5b50613a9388828901613091565b9150509295509295909350565b600080600060608486031215613ab4578081fd5b8335925060208401356001600160401b0380821115613ad1578283fd5b613add87838801612dc0565b93506040860135915080821115613af2578283fd5b50613aff86828701613091565b9150509250925092565b600080600080600060a08688031215613b20578283fd5b853594506020860135935060408601356001600160401b0380821115613b44578485fd5b613a6a89838a01613156565b60008060008060808587031215613b65578182fd5b5050823594602084013594506040840135936060013592509050565b600080600080600080600060e0888a031215613b9b578081fd5b613ba4886133cd565b965060208801356001600160401b0380821115613bbf578283fd5b613bcb8b838c016131f1565b975060408a0135915080821115613be0578283fd5b613bec8b838c01613156565b965060608a0135915080821115613c01578283fd5b613c0d8b838c01613156565b9550613c1b60808b016133e3565b945060a08a0135915080821115613c30578283fd5b613c3c8b838c016130ef565b935060c08a0135915080821115613c51578283fd5b50613c5e8a828b01612f64565b91505092959891949750929550565b6000815180845260208085019450808401835b83811015613ca55781516001600160a01b031687529582019590820190600101613c80565b509495945050505050565b6000815180845260208085018081965082840281019150828601855b85811015613d23578284038952815180518552858101518686015260408082015160ff1690860152606090810151608091860182905290613d0f81870183613e72565b9a87019a9550505090840190600101613ccc565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613ca5578151805160ff16885283810151848901526040908101519088015260609096019590820190600101613d43565b6000815180845260208085018081965082840281019150828601855b85811015613d23578284038952815180516001600160a01b0316855285810151606087870181905290613dc982880182613e72565b91505060408083015192508682038188015250613de68183613cb0565b9a87019a9550505090840190600101613d94565b6000815180845260208085018081965082840281019150828601855b85811015613d23578284038952613e2e848351613e9e565b98850198935090840190600101613e16565b6000815180845260208085019450808401835b83811015613ca557815160ff1687529582019590820190600101613e53565b60008151808452613e8a816020860160208601614992565b601f01601f19169290920160200192915050565b6000815160408452613eb36040850182613e72565b905060208301518482036020860152610b418282613e72565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b600060808252613f6e6080830187613cb0565b8281036020840152613f808187613cb0565b90508281036040840152613f948186613cb0565b91505082606083015295945050505050565b600060808252613fb96080830187613cb0565b85151560208401528281036040840152613f948186613cb0565b6000602082526119e66020830184613d78565b901515815260200190565b90815260200190565b600086825260a0602083015261401360a0830187613e72565b82810360408401526140258187613d78565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b6000602082526119e66020830184613e72565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b6020808252601f908201527f496e76616c696420466f7263654d6f7665417070205472616e736974696f6e00604082015260600190565b60208082526015908201527427379037b733b7b4b7339031b430b63632b733b29760591b604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b602080825260129082015271496e76616c6964207369676e61747572657360701b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b602080825260169082015275546f6f206d616e79207061727469636970616e74732160501b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f5369676e6572206e6f7420617574686f72697a6564206d6f7665720000000000604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252818101527f426164207c7369676e6174757265737c767c77686f5369676e6564576861747c604082015260600190565b6020808252818101527f556e61636365707461626c652077686f5369676e656457686174206172726179604082015260600190565b60208082526018908201527f4f7574636f6d65206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601c908201527f737461747573284368616e6e656c4461746129213d73746f7261676500000000604082015260600190565b60208082526018908201527f61707044617461206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601d908201527f496e76616c6964207369676e617475726573202f2021697346696e616c000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b60208082526017908201527f4f7574636f6d65206368616e676520766572626f74656e000000000000000000604082015260600190565b6020808252601e908201527f7c77686f5369676e6564576861747c213d6e5061727469636970616e74730000604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b602080825260129082015271697346696e616c20726574726f677261646560701b604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b6000608082526147ac6080830187613e9e565b82810360208401526147be8187613e9e565b65ffffffffffff95909516604084015250506060015292915050565b600086825260a060208301526147f360a0830187613c6d565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff808a1683528089166020840152871515604084015260e06060840152865160e0840152602087015160a061010085015261488d610180850182613c6d565b6040890151831661012086015260608901516001600160a01b03166101408601526080808a015184166101608701528582039086015290506148cf8188613dfa565b91505082810360a08401526148e48186613d30565b905082810360c08401526148f88185613e40565b9a9950505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b038111828210171561494c57fe5b604052919050565b60006001600160401b0382111561496757fe5b5060209081020190565b60006001600160401b0382111561498457fe5b50601f01601f191660200190565b60005b838110156149ad578181015183820152602001614995565b8381111561202d5750506000910152565b6001600160a01b03811681146118a257600080fd5b80151581146118a257600080fd5b60ff811681146118a257600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a2646970667358221220994750233fad66a0b8514403ce8e565b8455917bbcaa7983f00fb7ed697abed564736f6c63430007040033",
}

// NitroAdjudicatorABI is the input ABI used to generate the binding from.
// Deprecated: Use NitroAdjudicatorMetaData.ABI instead.
var NitroAdjudicatorABI = NitroAdjudicatorMetaData.ABI

// NitroAdjudicatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NitroAdjudicatorMetaData.Bin instead.
var NitroAdjudicatorBin = NitroAdjudicatorMetaData.Bin

// DeployNitroAdjudicator deploys a new Ethereum contract, binding an instance of NitroAdjudicator to it.
func DeployNitroAdjudicator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NitroAdjudicator, error) {
	parsed, err := NitroAdjudicatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NitroAdjudicatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NitroAdjudicator{NitroAdjudicatorCaller: NitroAdjudicatorCaller{contract: contract}, NitroAdjudicatorTransactor: NitroAdjudicatorTransactor{contract: contract}, NitroAdjudicatorFilterer: NitroAdjudicatorFilterer{contract: contract}}, nil
}

// NitroAdjudicator is an auto generated Go binding around an Ethereum contract.
type NitroAdjudicator struct {
	NitroAdjudicatorCaller     // Read-only binding to the contract
	NitroAdjudicatorTransactor // Write-only binding to the contract
	NitroAdjudicatorFilterer   // Log filterer for contract events
}

// NitroAdjudicatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type NitroAdjudicatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjudicatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NitroAdjudicatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjudicatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NitroAdjudicatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjudicatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NitroAdjudicatorSession struct {
	Contract     *NitroAdjudicator // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NitroAdjudicatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NitroAdjudicatorCallerSession struct {
	Contract *NitroAdjudicatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts           // Call options to use throughout this session
}

// NitroAdjudicatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NitroAdjudicatorTransactorSession struct {
	Contract     *NitroAdjudicatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// NitroAdjudicatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type NitroAdjudicatorRaw struct {
	Contract *NitroAdjudicator // Generic contract binding to access the raw methods on
}

// NitroAdjudicatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NitroAdjudicatorCallerRaw struct {
	Contract *NitroAdjudicatorCaller // Generic read-only contract binding to access the raw methods on
}

// NitroAdjudicatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NitroAdjudicatorTransactorRaw struct {
	Contract *NitroAdjudicatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNitroAdjudicator creates a new instance of NitroAdjudicator, bound to a specific deployed contract.
func NewNitroAdjudicator(address common.Address, backend bind.ContractBackend) (*NitroAdjudicator, error) {
	contract, err := bindNitroAdjudicator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicator{NitroAdjudicatorCaller: NitroAdjudicatorCaller{contract: contract}, NitroAdjudicatorTransactor: NitroAdjudicatorTransactor{contract: contract}, NitroAdjudicatorFilterer: NitroAdjudicatorFilterer{contract: contract}}, nil
}

// NewNitroAdjudicatorCaller creates a new read-only instance of NitroAdjudicator, bound to a specific deployed contract.
func NewNitroAdjudicatorCaller(address common.Address, caller bind.ContractCaller) (*NitroAdjudicatorCaller, error) {
	contract, err := bindNitroAdjudicator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorCaller{contract: contract}, nil
}

// NewNitroAdjudicatorTransactor creates a new write-only instance of NitroAdjudicator, bound to a specific deployed contract.
func NewNitroAdjudicatorTransactor(address common.Address, transactor bind.ContractTransactor) (*NitroAdjudicatorTransactor, error) {
	contract, err := bindNitroAdjudicator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorTransactor{contract: contract}, nil
}

// NewNitroAdjudicatorFilterer creates a new log filterer instance of NitroAdjudicator, bound to a specific deployed contract.
func NewNitroAdjudicatorFilterer(address common.Address, filterer bind.ContractFilterer) (*NitroAdjudicatorFilterer, error) {
	contract, err := bindNitroAdjudicator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorFilterer{contract: contract}, nil
}

// bindNitroAdjudicator binds a generic wrapper to an already deployed contract.
func bindNitroAdjudicator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NitroAdjudicatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NitroAdjudicator *NitroAdjudicatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NitroAdjudicator.Contract.NitroAdjudicatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NitroAdjudicator *NitroAdjudicatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.NitroAdjudicatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NitroAdjudicator *NitroAdjudicatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.NitroAdjudicatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NitroAdjudicator *NitroAdjudicatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NitroAdjudicator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NitroAdjudicator *NitroAdjudicatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NitroAdjudicator *NitroAdjudicatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.contract.Transact(opts, method, params...)
}

// ComputeClaimEffectsAndInteractions is a free data retrieval call binding the contract method 0xe29cffe0.
//
// Solidity: function compute_claim_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource, uint256[] targetAllocationIndicesToPayout) pure returns((bytes32,uint256,uint8,bytes)[] newSourceAllocations, (bytes32,uint256,uint8,bytes)[] newTargetAllocations, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorCaller) ComputeClaimEffectsAndInteractions(opts *bind.CallOpts, initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "compute_claim_effects_and_interactions", initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)

	outstruct := new(struct {
		NewSourceAllocations []ExitFormatAllocation
		NewTargetAllocations []ExitFormatAllocation
		ExitAllocations      []ExitFormatAllocation
		TotalPayouts         *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewSourceAllocations = *abi.ConvertType(out[0], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)
	outstruct.NewTargetAllocations = *abi.ConvertType(out[1], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)
	outstruct.ExitAllocations = *abi.ConvertType(out[2], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)
	outstruct.TotalPayouts = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ComputeClaimEffectsAndInteractions is a free data retrieval call binding the contract method 0xe29cffe0.
//
// Solidity: function compute_claim_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource, uint256[] targetAllocationIndicesToPayout) pure returns((bytes32,uint256,uint8,bytes)[] newSourceAllocations, (bytes32,uint256,uint8,bytes)[] newTargetAllocations, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorSession) ComputeClaimEffectsAndInteractions(initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	return _NitroAdjudicator.Contract.ComputeClaimEffectsAndInteractions(&_NitroAdjudicator.CallOpts, initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)
}

// ComputeClaimEffectsAndInteractions is a free data retrieval call binding the contract method 0xe29cffe0.
//
// Solidity: function compute_claim_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource, uint256[] targetAllocationIndicesToPayout) pure returns((bytes32,uint256,uint8,bytes)[] newSourceAllocations, (bytes32,uint256,uint8,bytes)[] newTargetAllocations, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ComputeClaimEffectsAndInteractions(initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	return _NitroAdjudicator.Contract.ComputeClaimEffectsAndInteractions(&_NitroAdjudicator.CallOpts, initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)
}

// ComputeTransferEffectsAndInteractions is a free data retrieval call binding the contract method 0x11e9f178.
//
// Solidity: function compute_transfer_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] allocations, uint256[] indices) pure returns((bytes32,uint256,uint8,bytes)[] newAllocations, bool allocatesOnlyZeros, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorCaller) ComputeTransferEffectsAndInteractions(opts *bind.CallOpts, initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "compute_transfer_effects_and_interactions", initialHoldings, allocations, indices)

	outstruct := new(struct {
		NewAllocations     []ExitFormatAllocation
		AllocatesOnlyZeros bool
		ExitAllocations    []ExitFormatAllocation
		TotalPayouts       *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NewAllocations = *abi.ConvertType(out[0], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)
	outstruct.AllocatesOnlyZeros = *abi.ConvertType(out[1], new(bool)).(*bool)
	outstruct.ExitAllocations = *abi.ConvertType(out[2], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)
	outstruct.TotalPayouts = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// ComputeTransferEffectsAndInteractions is a free data retrieval call binding the contract method 0x11e9f178.
//
// Solidity: function compute_transfer_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] allocations, uint256[] indices) pure returns((bytes32,uint256,uint8,bytes)[] newAllocations, bool allocatesOnlyZeros, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorSession) ComputeTransferEffectsAndInteractions(initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	return _NitroAdjudicator.Contract.ComputeTransferEffectsAndInteractions(&_NitroAdjudicator.CallOpts, initialHoldings, allocations, indices)
}

// ComputeTransferEffectsAndInteractions is a free data retrieval call binding the contract method 0x11e9f178.
//
// Solidity: function compute_transfer_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] allocations, uint256[] indices) pure returns((bytes32,uint256,uint8,bytes)[] newAllocations, bool allocatesOnlyZeros, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ComputeTransferEffectsAndInteractions(initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	return _NitroAdjudicator.Contract.ComputeTransferEffectsAndInteractions(&_NitroAdjudicator.CallOpts, initialHoldings, allocations, indices)
}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorCaller) GetChainID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "getChainID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorSession) GetChainID() (*big.Int, error) {
	return _NitroAdjudicator.Contract.GetChainID(&_NitroAdjudicator.CallOpts)
}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) GetChainID() (*big.Int, error) {
	return _NitroAdjudicator.Contract.GetChainID(&_NitroAdjudicator.CallOpts)
}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorCaller) Holdings(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "holdings", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorSession) Holdings(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _NitroAdjudicator.Contract.Holdings(&_NitroAdjudicator.CallOpts, arg0, arg1)
}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) Holdings(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _NitroAdjudicator.Contract.Holdings(&_NitroAdjudicator.CallOpts, arg0, arg1)
}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCaller) RequireValidInput(opts *bind.CallOpts, numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "requireValidInput", numParticipants, numStates, numSigs, numWhoSignedWhats)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorSession) RequireValidInput(numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	return _NitroAdjudicator.Contract.RequireValidInput(&_NitroAdjudicator.CallOpts, numParticipants, numStates, numSigs, numWhoSignedWhats)
}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) RequireValidInput(numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	return _NitroAdjudicator.Contract.RequireValidInput(&_NitroAdjudicator.CallOpts, numParticipants, numStates, numSigs, numWhoSignedWhats)
}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjudicator *NitroAdjudicatorCaller) StatusOf(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "statusOf", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjudicator *NitroAdjudicatorSession) StatusOf(arg0 [32]byte) ([32]byte, error) {
	return _NitroAdjudicator.Contract.StatusOf(&_NitroAdjudicator.CallOpts, arg0)
}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) StatusOf(arg0 [32]byte) ([32]byte, error) {
	return _NitroAdjudicator.Contract.StatusOf(&_NitroAdjudicator.CallOpts, arg0)
}

// UnpackStatus is a free data retrieval call binding the contract method 0x552cfa50.
//
// Solidity: function unpackStatus(bytes32 channelId) view returns(uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
func (_NitroAdjudicator *NitroAdjudicatorCaller) UnpackStatus(opts *bind.CallOpts, channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "unpackStatus", channelId)

	outstruct := new(struct {
		TurnNumRecord *big.Int
		FinalizesAt   *big.Int
		Fingerprint   *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TurnNumRecord = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.FinalizesAt = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.Fingerprint = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// UnpackStatus is a free data retrieval call binding the contract method 0x552cfa50.
//
// Solidity: function unpackStatus(bytes32 channelId) view returns(uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
func (_NitroAdjudicator *NitroAdjudicatorSession) UnpackStatus(channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	return _NitroAdjudicator.Contract.UnpackStatus(&_NitroAdjudicator.CallOpts, channelId)
}

// UnpackStatus is a free data retrieval call binding the contract method 0x552cfa50.
//
// Solidity: function unpackStatus(bytes32 channelId) view returns(uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) UnpackStatus(channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	return _NitroAdjudicator.Contract.UnpackStatus(&_NitroAdjudicator.CallOpts, channelId)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCaller) ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "validTransition", nParticipants, isFinalAB, ab, turnNumB, appDefinition)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Checkpoint(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Claim(opts *bind.TransactOpts, claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "claim", claimArgs)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Claim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Claim(&_NitroAdjudicator.TransactOpts, claimArgs)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Claim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Claim(&_NitroAdjudicator.TransactOpts, claimArgs)
}

// Conclude is a paid mutator transaction binding the contract method 0xcfaa7978.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0xcfaa7978.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0xcfaa7978.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x732f9208.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", largestTurnNum, fixedPart, appData, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x732f9208.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x732f9208.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "deposit", asset, channelId, expectedHeld, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Deposit(asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Deposit(&_NitroAdjudicator.TransactOpts, asset, channelId, expectedHeld, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Deposit(asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Deposit(&_NitroAdjudicator.TransactOpts, asset, channelId, expectedHeld, amount)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Respond(opts *bind.TransactOpts, isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "respond", isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Respond(isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Respond(&_NitroAdjudicator.TransactOpts, isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Respond(isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Respond(&_NitroAdjudicator.TransactOpts, isFinalAB, fixedPart, variablePartAB, sig)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Transfer(opts *bind.TransactOpts, assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "transfer", assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Transfer(assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Transfer(&_NitroAdjudicator.TransactOpts, assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Transfer(assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Transfer(&_NitroAdjudicator.TransactOpts, assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "transferAllAssets", channelId, outcomeBytes, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) TransferAllAssets(channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.TransferAllAssets(&_NitroAdjudicator.TransactOpts, channelId, outcomeBytes, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) TransferAllAssets(channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.TransferAllAssets(&_NitroAdjudicator.TransactOpts, channelId, outcomeBytes, stateHash)
}

// NitroAdjudicatorAllocationUpdatedIterator is returned from FilterAllocationUpdated and is used to iterate over the raw logs and unpacked data for AllocationUpdated events raised by the NitroAdjudicator contract.
type NitroAdjudicatorAllocationUpdatedIterator struct {
	Event *NitroAdjudicatorAllocationUpdated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NitroAdjudicatorAllocationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorAllocationUpdated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NitroAdjudicatorAllocationUpdated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NitroAdjudicatorAllocationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorAllocationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorAllocationUpdated represents a AllocationUpdated event raised by the NitroAdjudicator contract.
type NitroAdjudicatorAllocationUpdated struct {
	ChannelId       [32]byte
	AssetIndex      *big.Int
	InitialHoldings *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAllocationUpdated is a free log retrieval operation binding the contract event 0xb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterAllocationUpdated(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjudicatorAllocationUpdatedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "AllocationUpdated", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorAllocationUpdatedIterator{contract: _NitroAdjudicator.contract, event: "AllocationUpdated", logs: logs, sub: sub}, nil
}

// WatchAllocationUpdated is a free log subscription operation binding the contract event 0xb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchAllocationUpdated(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorAllocationUpdated, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "AllocationUpdated", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorAllocationUpdated)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "AllocationUpdated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseAllocationUpdated is a log parse operation binding the contract event 0xb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseAllocationUpdated(log types.Log) (*NitroAdjudicatorAllocationUpdated, error) {
	event := new(NitroAdjudicatorAllocationUpdated)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "AllocationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjudicatorChallengeClearedIterator is returned from FilterChallengeCleared and is used to iterate over the raw logs and unpacked data for ChallengeCleared events raised by the NitroAdjudicator contract.
type NitroAdjudicatorChallengeClearedIterator struct {
	Event *NitroAdjudicatorChallengeCleared // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NitroAdjudicatorChallengeClearedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorChallengeCleared)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NitroAdjudicatorChallengeCleared)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NitroAdjudicatorChallengeClearedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorChallengeClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorChallengeCleared represents a ChallengeCleared event raised by the NitroAdjudicator contract.
type NitroAdjudicatorChallengeCleared struct {
	ChannelId        [32]byte
	NewTurnNumRecord *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChallengeCleared is a free log retrieval operation binding the contract event 0x07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0.
//
// Solidity: event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterChallengeCleared(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjudicatorChallengeClearedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "ChallengeCleared", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorChallengeClearedIterator{contract: _NitroAdjudicator.contract, event: "ChallengeCleared", logs: logs, sub: sub}, nil
}

// WatchChallengeCleared is a free log subscription operation binding the contract event 0x07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0.
//
// Solidity: event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchChallengeCleared(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorChallengeCleared, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "ChallengeCleared", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorChallengeCleared)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "ChallengeCleared", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChallengeCleared is a log parse operation binding the contract event 0x07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0.
//
// Solidity: event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseChallengeCleared(log types.Log) (*NitroAdjudicatorChallengeCleared, error) {
	event := new(NitroAdjudicatorChallengeCleared)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "ChallengeCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjudicatorChallengeRegisteredIterator is returned from FilterChallengeRegistered and is used to iterate over the raw logs and unpacked data for ChallengeRegistered events raised by the NitroAdjudicator contract.
type NitroAdjudicatorChallengeRegisteredIterator struct {
	Event *NitroAdjudicatorChallengeRegistered // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NitroAdjudicatorChallengeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorChallengeRegistered)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NitroAdjudicatorChallengeRegistered)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NitroAdjudicatorChallengeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorChallengeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorChallengeRegistered represents a ChallengeRegistered event raised by the NitroAdjudicator contract.
type NitroAdjudicatorChallengeRegistered struct {
	ChannelId     [32]byte
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	IsFinal       bool
	FixedPart     IForceMoveFixedPart
	VariableParts []IForceMoveAppVariablePart
	Sigs          []IForceMoveSignature
	WhoSignedWhat []uint8
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChallengeRegistered is a free log retrieval operation binding the contract event 0xf6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterChallengeRegistered(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjudicatorChallengeRegisteredIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "ChallengeRegistered", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorChallengeRegisteredIterator{contract: _NitroAdjudicator.contract, event: "ChallengeRegistered", logs: logs, sub: sub}, nil
}

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0xf6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchChallengeRegistered(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorChallengeRegistered, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "ChallengeRegistered", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorChallengeRegistered)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "ChallengeRegistered", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseChallengeRegistered is a log parse operation binding the contract event 0xf6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseChallengeRegistered(log types.Log) (*NitroAdjudicatorChallengeRegistered, error) {
	event := new(NitroAdjudicatorChallengeRegistered)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "ChallengeRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjudicatorConcludedIterator is returned from FilterConcluded and is used to iterate over the raw logs and unpacked data for Concluded events raised by the NitroAdjudicator contract.
type NitroAdjudicatorConcludedIterator struct {
	Event *NitroAdjudicatorConcluded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NitroAdjudicatorConcludedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorConcluded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NitroAdjudicatorConcluded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NitroAdjudicatorConcludedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorConcludedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorConcluded represents a Concluded event raised by the NitroAdjudicator contract.
type NitroAdjudicatorConcluded struct {
	ChannelId   [32]byte
	FinalizesAt *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterConcluded is a free log retrieval operation binding the contract event 0x4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901.
//
// Solidity: event Concluded(bytes32 indexed channelId, uint48 finalizesAt)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterConcluded(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjudicatorConcludedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "Concluded", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorConcludedIterator{contract: _NitroAdjudicator.contract, event: "Concluded", logs: logs, sub: sub}, nil
}

// WatchConcluded is a free log subscription operation binding the contract event 0x4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901.
//
// Solidity: event Concluded(bytes32 indexed channelId, uint48 finalizesAt)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchConcluded(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorConcluded, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "Concluded", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorConcluded)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "Concluded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseConcluded is a log parse operation binding the contract event 0x4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901.
//
// Solidity: event Concluded(bytes32 indexed channelId, uint48 finalizesAt)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseConcluded(log types.Log) (*NitroAdjudicatorConcluded, error) {
	event := new(NitroAdjudicatorConcluded)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "Concluded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjudicatorDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the NitroAdjudicator contract.
type NitroAdjudicatorDepositedIterator struct {
	Event *NitroAdjudicatorDeposited // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *NitroAdjudicatorDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorDeposited)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(NitroAdjudicatorDeposited)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *NitroAdjudicatorDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorDeposited represents a Deposited event raised by the NitroAdjudicator contract.
type NitroAdjudicatorDeposited struct {
	Destination         [32]byte
	Asset               common.Address
	AmountDeposited     *big.Int
	DestinationHoldings *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterDeposited(opts *bind.FilterOpts, destination [][32]byte) (*NitroAdjudicatorDepositedIterator, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "Deposited", destinationRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorDepositedIterator{contract: _NitroAdjudicator.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorDeposited, destination [][32]byte) (event.Subscription, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "Deposited", destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorDeposited)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "Deposited", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseDeposited is a log parse operation binding the contract event 0x2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseDeposited(log types.Log) (*NitroAdjudicatorDeposited, error) {
	event := new(NitroAdjudicatorDeposited)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

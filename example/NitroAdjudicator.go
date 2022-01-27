// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package example

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

// NitroAdjucatorMetaData contains all meta data concerning the NitroAdjucator contract.
var NitroAdjucatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"indexed\":false,\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"name\":\"compute_claim_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newSourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newTargetAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"appPartHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"outcomeHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes32\",\"name\":\"appPartHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numParticipants\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numSigs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numWhoSignedWhats\",\"type\":\"uint256\"}],\"name\":\"requireValidInput\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"variablePartAB\",\"type\":\"tuple[2]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"sig\",\"type\":\"tuple\"}],\"name\":\"respond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nParticipants\",\"type\":\"uint256\"},{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"internalType\":\"bytes\",\"name\":\"outcome\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"ab\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint48\",\"name\":\"turnNumB\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"}],\"name\":\"validTransition\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614acc806100206000396000f3fe6080604052600436106100fe5760003560e01c80636775b17311610095578063c7df14e211610064578063c7df14e2146102af578063da4cdf73146102cf578063e03cdb0c146102ef578063e29cffe01461030f578063f198bea91461033f576100fe565b80636775b173146102225780636c3655221461024f5780637ff6a9821461026f578063be5c2a311461028f576100fe565b80633033730e116100d15780633033730e1461019e57806330776841146101be578063552cfa50146101de578063564b81ef1461020d576100fe565b80630149b7621461010357806311e9f17814610125578063166e56cd1461015e5780632fb1d2701461018b575b600080fd5b34801561010f57600080fd5b5061012361011e3660046137dc565b61035f565b005b34801561013157600080fd5b50610145610140366004613a90565b6103b7565b6040516101559493929190613f79565b60405180910390f35b34801561016a57600080fd5b5061017e6101793660046133e9565b610705565b604051610155919061404a565b610123610199366004613414565b610722565b3480156101aa57600080fd5b506101236101b9366004613af9565b6109b0565b3480156101ca57600080fd5b506101236101d93660046136e8565b610a30565b3480156101ea57600080fd5b506101fe6101f9366004613683565b610b0f565b6040516101559392919061498c565b34801561021957600080fd5b5061017e610b2a565b34801561022e57600080fd5b5061024261023d36600461397e565b610b2e565b604051610155919061403f565b34801561025b57600080fd5b5061012361026a36600461369b565b610b49565b34801561027b57600080fd5b5061012361028a366004613b71565b610f06565b34801561029b57600080fd5b506102426102aa366004613b40565b610f15565b3480156102bb57600080fd5b5061017e6102ca366004613683565b610f9a565b3480156102db57600080fd5b506101236102ea36600461344e565b610fac565b3480156102fb57600080fd5b5061012361030a366004613c2d565b6110eb565b34801561031b57600080fd5b5061032f61032a3660046139f9565b611111565b6040516101559493929190613f2e565b34801561034b57600080fd5b5061012361035a3660046138a3565b611688565b610373866020015151855184518451610f15565b50600061037f876117f4565b905061038a8161185b565b6103948187611892565b6103a3868686848b88886118d5565b506103ae818761193a565b50505050505050565b6060600060606000808551116103ce5785516103d1565b84515b6001600160401b03811180156103e657600080fd5b5060405190808252806020026020018201604052801561042057816020015b61040d612d0f565b8152602001906001900390816104055790505b5091506000905085516001600160401b038111801561043e57600080fd5b5060405190808252806020026020018201604052801561047857816020015b610465612d0f565b81526020019060019003908161045d5790505b50935060019250866000805b88518110156106f95788818151811061049957fe5b6020026020010151600001518782815181106104b157fe5b602002602001015160000181815250508881815181106104cd57fe5b6020026020010151604001518782815181106104e557fe5b60200260200101516040019060ff16908160ff168152505088818151811061050957fe5b60200260200101516060015187828151811061052157fe5b60200260200101516060018190525060006105538a838151811061054157fe5b602002602001015160200151856119c1565b9050885160001480610582575088518310801561058257508189848151811061057857fe5b6020026020010151145b1561069457600260ff168a848151811061059857fe5b60200260200101516040015160ff1614156105ce5760405162461bcd60e51b81526004016105c590614371565b60405180910390fd5b808a83815181106105db57fe5b602002602001015160200151038883815181106105f457fe5b6020026020010151602001818152505060405180608001604052808b848151811061061b57fe5b60200260200101516000015181526020018281526020018b848151811061063e57fe5b60200260200101516040015160ff1681526020018b848151811061065e57fe5b60200260200101516060015181525086848151811061067957fe5b602002602001018190525080850194508260010192506106c9565b8982815181106106a057fe5b6020026020010151602001518883815181106106b857fe5b602002602001015160200181815250505b8782815181106106d557fe5b6020026020010151602001516000146106ed57600096505b90920391600101610484565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b61072b836119db565b156107485760405162461bcd60e51b81526004016105c5906145b6565b6001600160a01b03841660009081526001602090815260408083208684529091528120548381101561078c5760405162461bcd60e51b81526004016105c5906140e5565b61079684846119e7565b81106107b45760405162461bcd60e51b81526004016105c5906146bb565b6107c8816107c286866119e7565b90611a41565b91506001600160a01b0386166107fc578234146107f75760405162461bcd60e51b81526004016105c590614729565b61089a565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd9061082c90339030908790600401613ed0565b602060405180830381600087803b15801561084657600080fd5b505af115801561085a573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061087e9190613667565b61089a5760405162461bcd60e51b81526004016105c590614439565b60006108a682846119e7565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a71590610905908a9087908690613f0d565b60405180910390a26001600160a01b0387166103ae5760006109278585611a41565b90506000336001600160a01b03168260405161094290612b84565b60006040518083038185875af1925050503d806000811461097f576040519150601f19603f3d011682016040523d82523d6000602084013e610984565b606091505b50509050806109a55760405162461bcd60e51b81526004016105c59061478c565b505050505050505050565b60008060006109c2888589888a611a9e565b92509250925060008060006109ef84878d815181106109dd57fe5b602002602001015160400151896103b7565b93509350509250610a068b868c8b8a888a88611b1c565b610a23868c81518110610a1557fe5b602002602001015183611bc8565b5050505050505050505050565b600080600080610a3f85611c00565b93509350935093506060806060600080888a6060015181518110610a5f57fe5b60200260200101516040015190506000888b60e0015181518110610a7f57fe5b6020026020010151604001519050610aa38783838e608001518f6101000151611111565b809650819750829850839950505050505050610afc89878a878c8e6060015181518110610acc57fe5b6020026020010151604001518e6080015181518110610ae757fe5b6020026020010151600001518c898c89611da9565b6109a5878a60e0015181518110610a1557fe5b6000806000610b1d84611eac565b9196909550909350915050565b4690565b6000610b3d8686868686611eca565b90505b95945050505050565b610b5283611f9d565b610b6481838051906020012085611fd0565b60016000610b718461201f565b9050600081516001600160401b0381118015610b8c57600080fd5b50604051908082528060200260200182016040528015610bc657816020015b610bb3612d35565b815260200190600190039081610bab5790505b509050600082516001600160401b0381118015610be257600080fd5b50604051908082528060200260200182016040528015610c0c578160200160208202803683370190505b509050600083516001600160401b0381118015610c2857600080fd5b50604051908082528060200260200182016040528015610c52578160200160208202803683370190505b50905060005b8451811015610dec576000858281518110610c6f57fe5b602002602001015190506000816040015190506000878481518110610c9057fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008d815260200190815260200160002054868581518110610ce157fe5b602002602001018181525050600080600080610d528a8981518110610d0257fe5b60200260200101518760006001600160401b0381118015610d2257600080fd5b50604051908082528060200260200182016040528015610d4c578160200160208202803683370190505b506103b7565b935093509350935082610d645760009c505b80898981518110610d7157fe5b602002602001018181525050838c8981518110610d8a57fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b8981518110610dcd57fe5b6020026020010181905250505050505050508080600101915050610c58565b5060005b8451811015610ea0576000858281518110610e0757fe5b6020026020010151600001519050828281518110610e2157fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208e835290935291909120805491909103905583518a90600080516020614a77833981519152908490879082908110610e7957fe5b6020026020010151604051610e8f929190614879565b60405180910390a250600101610df0565b508415610ebb57600088815260208190526040812055610ef3565b600087604051602001610ece919061409b565b604051602081830303815290604052805190602001209050610ef1898883612035565b505b610efc8361209b565b5050505050505050565b610efc878787878787876120cb565b6000838510158015610f275750600084115b610f435760405162461bcd60e51b81526004016105c5906142d3565b8483148015610f5157508482145b610f6d5760405162461bcd60e51b81526004016105c590614470565b60ff851115610f8e5760405162461bcd60e51b81526004016105c59061430a565b5060015b949350505050565b60006020819052908152604090205481565b6000610fb7846117f4565b9050600080610fc583611eac565b508651518051602091820120818901515180519201919091208a5193955091935091600090611004908690888c8c865b602002015160200151886122d2565b905060006110226001808801908d906020020151898d8d6001610ff5565b905061105e60405180608001604052808865ffffffffffff1681526020018765ffffffffffff1681526020018481526020018681525088612369565b60208a0151805165ffffffffffff60018901168161107857fe5b068151811061108357fe5b60200260200101516001600160a01b031661109e828a61237c565b6001600160a01b0316146110c45760405162461bcd60e51b81526004016105c5906143d3565b6110dd8a60200151518c8b896001018e60600151611eca565b50610a23878760010161193a565b835160208501206000611103898989858989896120cb565b90506109a581876000610b49565b606080606060008088516001600160401b038111801561113057600080fd5b5060405190808252806020026020018201604052801561116a57816020015b611157612d0f565b81526020019060019003908161114f5790505b50945087516001600160401b038111801561118457600080fd5b506040519080825280602002602001820160405280156111be57816020015b6111ab612d0f565b8152602001906001900390816111a35790505b50935087516001600160401b03811180156111d857600080fd5b5060405190808252806020026020018201604052801561121257816020015b6111ff612d0f565b8152602001906001900390816111f75790505b50925060005b89518110156112fe5789818151811061122d57fe5b60200260200101516000015186828151811061124557fe5b6020026020010151600001818152505089818151811061126157fe5b60200260200101516020015186828151811061127957fe5b6020026020010151602001818152505089818151811061129557fe5b6020026020010151606001518682815181106112ad57fe5b6020026020010151606001819052508981815181106112c857fe5b6020026020010151604001518682815181106112e057fe5b602090810291909101015160ff909116604090910152600101611218565b5060005b88518110156114a95788818151811061131757fe5b60200260200101516000015185828151811061132f57fe5b6020026020010151600001818152505088818151811061134b57fe5b60200260200101516020015185828151811061136357fe5b6020026020010151602001818152505088818151811061137f57fe5b60200260200101516060015185828151811061139757fe5b6020026020010151606001819052508881815181106113b257fe5b6020026020010151604001518582815181106113ca57fe5b60200260200101516040019060ff16908160ff16815250508881815181106113ee57fe5b60200260200101516000015184828151811061140657fe5b60200260200101516000018181525050600084828151811061142457fe5b6020026020010151602001818152505088818151811061144057fe5b60200260200101516060015184828151811061145857fe5b60200260200101516060018190525088818151811061147357fe5b60200260200101516040015184828151811061148b57fe5b602090810291909101015160ff909116604090910152600101611302565b508960005b888110156114f057816114c0576114f0565b60006114e38c83815181106114d157fe5b602002602001015160200151846119c1565b90920391506001016114ae565b506000611514828c8b8151811061150357fe5b6020026020010151602001516119c1565b905060006115388c8b8151811061152757fe5b60200260200101516060015161242e565b905060005b8151811015611677578261155057611677565b60005b885181101561166e57836115665761166e565b88818151811061157257fe5b60200260200101516000015183838151811061158a57fe5b602002602001015114156116665760006115bb8e83815181106115a957fe5b602002602001015160200151866119c1565b905080850394508b51600014806115ef57508b51871080156115ef5750818c88815181106115e557fe5b6020026020010151145b1561166057808a838151811061160157fe5b60200260200101516020018181510391508181525050808b8e8151811061162457fe5b602002602001015160200181815103915081815250508089838151811061164757fe5b6020908102919091018101510152968701966001909601955b5061166e565b600101611553565b5060010161153d565b505050505095509550955095915050565b61169c876020015151865185518551610f15565b5060006116a8886117f4565b905060006116b582612444565b60028111156116c057fe5b14156116d5576116d0818861248e565b611704565b60016116e082612444565b60028111156116eb57fe5b14156116fb576116d08188611892565b6117048161185b565b6000611715888888858d8a8a6118d5565b9050611726818a60200151856124cd565b817ff6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee898b60800151420160008a60ff16118d8c8b8b60405161176e97969594939291906148cc565b60405180910390a26117d760405180608001604052808a65ffffffffffff1681526020018b60800151420165ffffffffffff1681526020018381526020018960018b5103815181106117bc57fe5b60200260200101516000015180519060200120815250612527565b600092835260208390526040909220919091555050505050505050565b60006117fe610b2a565b82511461181d5760405162461bcd60e51b81526004016105c59061411c565b611825610b2a565b602080840151604080860151905161183e949301614848565b604051602081830303815290604052805190602001209050919050565b600261186682612444565b600281111561187157fe5b141561188f5760405162461bcd60e51b81526004016105c590614147565b50565b600061189d83611eac565b505090508065ffffffffffff168265ffffffffffff16116118d05760405162461bcd60e51b81526004016105c5906140ae565b505050565b6000806118e58989898989612578565b90506118f889866020015183878761274a565b6119145760405162461bcd60e51b81526004016105c5906142a7565b8060018251038151811061192457fe5b6020026020010151915050979650505050505050565b6040805160808101825265ffffffffffff83168152600060208201819052918101829052606081019190915261196f90612527565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0826040516119b59190614887565b60405180910390a25050565b60008183116119d057826119d2565b815b90505b92915050565b60a081901c155b919050565b6000828201838110156119d2576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611a98576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6060600080611aac87612820565b611ab586611f9d565b611ac785858051906020012088611fd0565b611ad08461201f565b9250828881518110611ade57fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a908110611b5657fe5b602002602001015160400181905250611b96868686604051602001611b7b9190613fa6565b60405160208183030381529060405280519060200120612035565b85600080516020614a778339815191528984604051611bb6929190614879565b60405180910390a25050505050505050565b611bfc604051806060016040528084600001516001600160a01b03168152602001846020015181526020018381525061287f565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611c3490612820565b611c3d85611f9d565b611c538a60200151858051906020012087611fd0565b611c5c8461201f565b9850611c678261201f565b9750888381518110611c7557fe5b6020908102919091010151519650600260ff16898481518110611c9457fe5b6020026020010151604001518281518110611cab57fe5b60200260200101516040015160ff1614611cd75760405162461bcd60e51b81526004016105c5906146f2565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a9085908110611d0c57fe5b6020026020010151604001518b6080015181518110611d2757fe5b6020026020010151600001519050876001600160a01b0316898381518110611d4b57fe5b6020026020010151600001516001600160a01b031614611d7d5760405162461bcd60e51b81526004016105c59061433a565b611d8681611f9d565b611d9c8b60a00151848051906020012083611fd0565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b9084908110611def57fe5b602002602001015160400181905250611e18838d602001518c604051602001611b7b9190613fa6565b85878281518110611e2557fe5b602002602001015160400181905250611e4e888d60a0015189604051602001611b7b9190613fa6565b82600080516020614a778339815191528387604051611e6e929190614879565b60405180910390a287600080516020614a778339815191528287604051611e96929190614879565b60405180910390a2505050505050505050505050565b60009081526020819052604090205460d081901c9160a082901c9190565b600080611ed98787878761292b565b90506001816001811115611ee957fe5b1415611f90578451602086015160405163fd7a2f6560e01b81526001600160a01b0386169263fd7a2f6592611f249289908d90600401614807565b60206040518083038186803b158015611f3c57600080fd5b505afa158015611f50573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611f749190613667565b611f905760405162461bcd60e51b81526004016105c5906141aa565b5060019695505050505050565b6002611fa882612444565b6002811115611fb357fe5b1461188f5760405162461bcd60e51b81526004016105c59061461d565b6000611fdb82611eac565b92505050611fe98484612a0f565b6001600160a01b0316816001600160a01b0316146120195760405162461bcd60e51b81526004016105c59061440a565b50505050565b6060818060200190518101906119d59190613560565b60008061204185611eac565b5091509150600061208160405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250612527565b600096875260208790526040909620959095555050505050565b60005b8151811015611bfc576120c38282815181106120b657fe5b602002602001015161287f565b60010161209e565b60006120d6876117f4565b90506120e18161185b565b6120f78760200151518560ff1684518651610f15565b508360ff168860010165ffffffffffff1610156121265760405162461bcd60e51b81526004016105c5906145ed565b60008460ff166001600160401b038111801561214157600080fd5b5060405190808252806020026020018201604052801561216b578160200160208202803683370190505b50905060005b8560ff168165ffffffffffff16101561220d576040518060a001604052808760ff16836001018d010365ffffffffffff168152602001600115158152602001848152602001898152602001888152506040516020016121d091906147c3565b60405160208183030381529060405280519060200120828265ffffffffffff16815181106121fa57fe5b6020908102919091010152600101612171565b5061221f89896020015183868861274a565b61223b5760405162461bcd60e51b81526004016105c59061457f565b6122786040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b815260200188815250612527565b60008084815260200190815260200160002081905550817f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516122be9190614887565b60405180910390a250979650505050505050565b60006040518060a001604052808865ffffffffffff168152602001871515815260200186815260200185608001518660600151866040516020016123189392919061489a565b6040516020818303038152906040528051906020012081526020018381525060405160200161234791906147c3565b6040516020818303038152906040528051906020012090509695505050505050565b6123738282612a3b565b611bfc81612a71565b600080836040516020016123909190613e9f565b6040516020818303038152906040528051906020012090506000600182856000015186602001518760400151604051600081526020016040526040516123d9949392919061407d565b6020604051602081039080840390855afa1580156123fb573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610f925760405162461bcd60e51b81526004016105c5906143a8565b6060818060200190518101906119d591906134d1565b60008061245083611eac565b5091505065ffffffffffff811661246b5760009150506119e2565b428165ffffffffffff16116124845760029150506119e2565b60019150506119e2565b600061249983611eac565b505090508065ffffffffffff168265ffffffffffff1610156118d05760405162461bcd60e51b81526004016105c590614173565b60006124ff846040516020016124e39190614053565b604051602081830303815290604052805190602001208361237c565b905061250b8184612aa4565b6120195760405162461bcd60e51b81526004016105c590614240565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161256491612a0f565b6001600160a01b0316919091179392505050565b6060600085516001600160401b038111801561259357600080fd5b506040519080825280602002602001820160405280156125bd578160200160208202803683370190505b509050600160ff86168803016000805b88518165ffffffffffff16101561273c578089518b0360010101915061264f828465ffffffffffff168465ffffffffffff16101589898d8665ffffffffffff168151811061261757fe5b6020026020010151602001518e8765ffffffffffff168151811061263757fe5b602002602001015160000151805190602001206122d2565b848265ffffffffffff168151811061266357fe5b6020026020010181815250508965ffffffffffff168265ffffffffffff1610156127345761273286602001515160405180604001604052808665ffffffffffff168665ffffffffffff1610151515151581526020018665ffffffffffff168660010165ffffffffffff1610151515151581525060405180604001604052808d8665ffffffffffff16815181106126f557fe5b602002602001015181526020018d8660010165ffffffffffff168151811061271957fe5b6020026020010151815250856001018a60600151611eca565b505b6001016125cd565b509198975050505050505050565b835183516000919061275e84898484612afa565b61277a5760405162461bcd60e51b81526004016105c5906144a5565b60005b828110156128115760006127cd8887848151811061279757fe5b602002602001015160ff16815181106127ac57fe5b60200260200101518884815181106127c057fe5b602002602001015161237c565b90508882815181106127db57fe5b60200260200101516001600160a01b0316816001600160a01b031614612808576000945050505050610b40565b5060010161277d565b50600198975050505050505050565b60005b8151816001011015611bfc5781816001018151811061283e57fe5b602002602001015182828151811061285257fe5b6020026020010151106128775760405162461bcd60e51b81526004016105c590614210565b600101612823565b805160005b8260400151518110156118d0576000836040015182815181106128a357fe5b60200260200101516000015190506000846040015183815181106128c357fe5b60200260200101516020015190506128da826119db565b156128f7576128f2846128ec84612b84565b83612b87565b612921565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b5050600101612884565b60208301516000901561296e5761294d8360015b602002015151845151612c97565b6129695760405162461bcd60e51b81526004016105c59061464d565b612a04565b83511561298d5760405162461bcd60e51b81526004016105c590614760565b846002028265ffffffffffff1610156129fc576129ab83600161293f565b6129c75760405162461bcd60e51b81526004016105c5906144da565b6020838101518101518451909101516129e09190612c97565b6129695760405162461bcd60e51b81526004016105c590614548565b506001610f92565b506000949350505050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b600081815260208190526040902054612a55908390612cfb565b611bfc5760405162461bcd60e51b81526004016105c590614511565b6001612a7c82612444565b6002811115612a8757fe5b1461188f5760405162461bcd60e51b81526004016105c5906141e1565b6000805b8251811015612af057828181518110612abd57fe5b60200260200101516001600160a01b0316846001600160a01b03161415612ae85760019150506119d5565b600101612aa8565b5060009392505050565b600082855114612b1c5760405162461bcd60e51b81526004016105c590614684565b60005b83811015612b7857600084828765ffffffffffff1687010381612b3e57fe5b0690508381888481518110612b4f57fe5b602002602001015160ff16016001011015612b6f57600092505050610f92565b50600101612b1f565b50600195945050505050565b90565b6001600160a01b038316612c17576000826001600160a01b031682604051612bae90612b84565b60006040518083038185875af1925050503d8060008114612beb576040519150601f19603f3d011682016040523d82523d6000602084013e612bf0565b606091505b5050905080612c115760405162461bcd60e51b81526004016105c590614277565b506118d0565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb90612c459085908590600401613ef4565b602060405180830381600087803b158015612c5f57600080fd5b505af1158015612c73573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906120199190613667565b815181516000916001918114808314612cb35760009250612cf1565b600160208701838101602088015b600284838510011415612cec578051835114612ce05760009650600093505b60209283019201612cc1565b505050505b5090949350505050565b600081612d0784612527565b149392505050565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b80356119e281614a44565b600082601f830112612d7a578081fd5b604051604081018181106001600160401b0382111715612d9657fe5b8060405250808385604086011115612dac578384fd5b835b6002811015612dd7578135612dc281614a59565b83526020928301929190910190600101612dae565b509195945050505050565b600082601f830112612df2578081fd5b81356020612e07612e02836149da565b6149b7565b82815281810190858301855b85811015612eac5781358801608080601f19838d03011215612e33578889fd5b604080518281016001600160401b038282108183111715612e5057fe5b908352848a0135825284830135828b0152606090612e6f8287016133d3565b83850152938501359380851115612e84578c8dfd5b50612e938e8b86880101613152565b9082015287525050509284019290840190600101612e13565b5090979650505050505050565b600082601f830112612ec9578081fd5b81516020612ed9612e02836149da565b82815281810190858301855b85811015612eac5781518801608080601f19838d03011215612f05578889fd5b604080518281016001600160401b038282108183111715612f2257fe5b908352848a0151825284830151828b0152606090612f418287016133de565b83850152938501519380851115612f56578c8dfd5b50612f658e8b8688010161319e565b9082015287525050509284019290840190600101612ee5565b600082601f830112612f8e578081fd5b81356020612f9e612e02836149da565b82815281810190858301606080860288018501891015612fbc578687fd5b865b8681101561273c57612fd08a846132de565b85529385019391810191600101612fbe565b600082601f830112612ff2578081fd5b604051604081018181106001600160401b038211171561300e57fe5b6040528083835b6002811015612dd75761302b878335880161333b565b83526020928301929190910190600101613015565b600082601f830112613050578081fd5b81356020613060612e02836149da565b82815281810190858301855b85811015612eac57613083898684358b010161333b565b8452928401929084019060010161306c565b600082601f8301126130a5578081fd5b813560206130b5612e02836149da565b82815281810190858301838502870184018810156130d1578586fd5b855b85811015612eac578135845292840192908401906001016130d3565b600082601f8301126130ff578081fd5b8135602061310f612e02836149da565b828152818101908583018385028701840188101561312b578586fd5b855b85811015612eac57813561314081614a67565b8452928401929084019060010161312d565b600082601f830112613162578081fd5b8135613170612e02826149f7565b818152846020838601011115613184578283fd5b816020850160208301379081016020019190915292915050565b600082601f8301126131ae578081fd5b81516131bc612e02826149f7565b8181528460208386010111156131d0578283fd5b610f92826020830160208701614a18565b600060a082840312156131f2578081fd5b60405160a081016001600160401b03828210818311171561320f57fe5b816040528293508435835260209150818501358181111561322f57600080fd5b85019050601f8101861361324257600080fd5b8035613250612e02826149da565b81815283810190838501858402850186018a101561326d57600080fd5b600094505b8385101561329957803561328581614a44565b835260019490940193918501918501613272565b50808587015250505050506132b0604084016133bd565b60408201526132c160608401612d5f565b60608201526132d2608084016133bd565b60808201525092915050565b6000606082840312156132ef578081fd5b604051606081018181106001600160401b038211171561330b57fe5b604052905080823561331c81614a67565b8082525060208301356020820152604083013560408201525092915050565b60006040828403121561334c578081fd5b604051604081016001600160401b03828210818311171561336957fe5b81604052829350843591508082111561338157600080fd5b61338d86838701613152565b835260208501359150808211156133a357600080fd5b506133b085828601613152565b6020830152505092915050565b803565ffffffffffff811681146119e257600080fd5b80356119e281614a67565b80516119e281614a67565b600080604083850312156133fb578182fd5b823561340681614a44565b946020939093013593505050565b60008060008060808587031215613429578182fd5b843561343481614a44565b966020860135965060408601359560600135945092505050565b60008060008060e08587031215613463578182fd5b61346d8686612d6a565b935060408501356001600160401b0380821115613488578384fd5b613494888389016131e1565b945060608701359150808211156134a9578384fd5b506134b687828801612fe2565b9250506134c686608087016132de565b905092959194509250565b600060208083850312156134e3578182fd5b82516001600160401b038111156134f8578283fd5b8301601f81018513613508578283fd5b8051613516612e02826149da565b8181528381019083850185840285018601891015613532578687fd5b8694505b83851015613554578051835260019490940193918501918501613536565b50979650505050505050565b60006020808385031215613572578182fd5b82516001600160401b0380821115613588578384fd5b818501915085601f83011261359b578384fd5b81516135a9612e02826149da565b81815284810190848601875b8481101561365857815187016060818d03601f190112156135d457898afd5b60408051606081018181108a821117156135ea57fe5b8252828b01516135f981614a44565b8152828201518981111561360b578c8dfd5b6136198f8d8387010161319e565b8c8301525060608301518981111561362f578c8dfd5b61363d8f8d83870101612eb9565b928201929092528652505092870192908701906001016135b5565b50909998505050505050505050565b600060208284031215613678578081fd5b81516119d281614a59565b600060208284031215613694578081fd5b5035919050565b6000806000606084860312156136af578081fd5b8335925060208401356001600160401b038111156136cb578182fd5b6136d786828701613152565b925050604084013590509250925092565b6000602082840312156136f9578081fd5b81356001600160401b038082111561370f578283fd5b8184019150610120808387031215613725578384fd5b61372e816149b7565b9050823581526020830135602082015260408301358281111561374f578485fd5b61375b87828601613152565b604083015250606083013560608201526080830135608082015260a083013560a082015260c083013582811115613790578485fd5b61379c87828601613152565b60c08301525060e083013560e082015261010080840135838111156137bf578586fd5b6137cb88828701613095565b918301919091525095945050505050565b60008060008060008060c087890312156137f4578384fd5b86356001600160401b038082111561380a578586fd5b6138168a838b016131e1565b975061382460208a016133bd565b96506040890135915080821115613839578586fd5b6138458a838b01613040565b955061385360608a016133d3565b94506080890135915080821115613868578384fd5b6138748a838b01612f7e565b935060a0890135915080821115613889578283fd5b5061389689828a016130ef565b9150509295509295509295565b6000806000806000806000610120888a0312156138be578485fd5b87356001600160401b03808211156138d4578687fd5b6138e08b838c016131e1565b98506138ee60208b016133bd565b975060408a0135915080821115613903578687fd5b61390f8b838c01613040565b965061391d60608b016133d3565b955060808a0135915080821115613932578283fd5b61393e8b838c01612f7e565b945060a08a0135915080821115613953578283fd5b506139608a828b016130ef565b9250506139708960c08a016132de565b905092959891949750929550565b600080600080600060c08688031215613995578283fd5b853594506139a68760208801612d6a565b935060608601356001600160401b038111156139c0578384fd5b6139cc88828901612fe2565b9350506139db608087016133bd565b915060a08601356139eb81614a44565b809150509295509295909350565b600080600080600060a08688031215613a10578283fd5b8535945060208601356001600160401b0380821115613a2d578485fd5b613a3989838a01612de2565b95506040880135915080821115613a4e578485fd5b613a5a89838a01612de2565b9450606088013593506080880135915080821115613a76578283fd5b50613a8388828901613095565b9150509295509295909350565b600080600060608486031215613aa4578081fd5b8335925060208401356001600160401b0380821115613ac1578283fd5b613acd87838801612de2565b93506040860135915080821115613ae2578283fd5b50613aef86828701613095565b9150509250925092565b600080600080600060a08688031215613b10578283fd5b853594506020860135935060408601356001600160401b0380821115613b34578485fd5b613a5a89838a01613152565b60008060008060808587031215613b55578182fd5b5050823594602084013594506040840135936060013592509050565b600080600080600080600060e0888a031215613b8b578081fd5b613b94886133bd565b965060208801356001600160401b0380821115613baf578283fd5b613bbb8b838c016131e1565b975060408a0135965060608a0135955060808a01359150613bdb82614a67565b90935060a08901359080821115613bf0578283fd5b613bfc8b838c016130ef565b935060c08a0135915080821115613c11578283fd5b50613c1e8a828b01612f7e565b91505092959891949750929550565b600080600080600080600060e0888a031215613c47578081fd5b613c50886133bd565b965060208801356001600160401b0380821115613c6b578283fd5b613c778b838c016131e1565b975060408a0135965060608a0135915080821115613c93578283fd5b613c9f8b838c01613152565b9550613cad60808b016133d3565b945060a08a0135915080821115613bf0578283fd5b6000815180845260208085019450808401835b83811015613cfa5781516001600160a01b031687529582019590820190600101613cd5565b509495945050505050565b6000815180845260208085018081965082840281019150828601855b85811015613d78578284038952815180518552858101518686015260408082015160ff1690860152606090810151608091860182905290613d6481870183613e45565b9a87019a9550505090840190600101613d21565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613cfa578151805160ff16885283810151848901526040908101519088015260609096019590820190600101613d98565b6000815180845260208085018081965082840281019150828601855b85811015613d78578284038952613e01848351613e71565b98850198935090840190600101613de9565b6000815180845260208085019450808401835b83811015613cfa57815160ff1687529582019590820190600101613e26565b60008151808452613e5d816020860160208601614a18565b601f01601f19169290920160200192915050565b6000815160408452613e866040850182613e45565b905060208301518482036020860152610b408282613e45565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b600060808252613f416080830187613d05565b8281036020840152613f538187613d05565b90508281036040840152613f678186613d05565b91505082606083015295945050505050565b600060808252613f8c6080830187613d05565b85151560208401528281036040840152613f678186613d05565b60208082528251828201819052600091906040908185019080840286018301878501865b8381101561403157888303603f19018552815180516001600160a01b031684528781015160608986018190529061400382870182613e45565b9150508782015191508481038886015261401d8183613d05565b968901969450505090860190600101613fca565b509098975050505050505050565b901515815260200190565b90815260200190565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6000602082526119d26020830184613e45565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b6020808252601f908201527f496e76616c696420466f7263654d6f7665417070205472616e736974696f6e00604082015260600190565b60208082526015908201527427379037b733b7b4b7339031b430b63632b733b29760591b604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b602080825260129082015271496e76616c6964207369676e61747572657360701b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b602080825260169082015275546f6f206d616e79207061727469636970616e74732160501b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f5369676e6572206e6f7420617574686f72697a6564206d6f7665720000000000604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252818101527f426164207c7369676e6174757265737c767c77686f5369676e6564576861747c604082015260600190565b6020808252818101527f556e61636365707461626c652077686f5369676e656457686174206172726179604082015260600190565b60208082526018908201527f4f7574636f6d65206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601c908201527f737461747573284368616e6e656c4461746129213d73746f7261676500000000604082015260600190565b60208082526018908201527f61707044617461206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601d908201527f496e76616c6964207369676e617475726573202f2021697346696e616c000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b60208082526017908201527f4f7574636f6d65206368616e676520766572626f74656e000000000000000000604082015260600190565b6020808252601e908201527f7c77686f5369676e6564576861747c213d6e5061727469636970616e74730000604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b602080825260129082015271697346696e616c20726574726f677261646560701b604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b600060a08201905065ffffffffffff835116825260208301511515602083015260408301516040830152606083015160608301526080830151608083015292915050565b60006080825261481a6080830187613e71565b828103602084015261482c8187613e71565b65ffffffffffff95909516604084015250506060015292915050565b6000848252606060208301526148616060830185613cc2565b905065ffffffffffff83166040830152949350505050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b65ffffffffffff841681526001600160a01b0383166020820152606060408201819052600090610b4090830184613e45565b600065ffffffffffff808a1683528089166020840152871515604084015260e06060840152865160e0840152602087015160a0610100850152614913610180850182613cc2565b6040890151831661012086015260608901516001600160a01b03166101408601526080808a015184166101608701528582039086015290506149558188613dcd565b91505082810360a084015261496a8186613d85565b905082810360c084015261497e8185613e13565b9a9950505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b03811182821017156149d257fe5b604052919050565b60006001600160401b038211156149ed57fe5b5060209081020190565b60006001600160401b03821115614a0a57fe5b50601f01601f191660200190565b60005b83811015614a33578181015183820152602001614a1b565b838111156120195750506000910152565b6001600160a01b038116811461188f57600080fd5b801515811461188f57600080fd5b60ff8116811461188f57600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a2646970667358221220e69671a37cad4aac4bb3d4c6444b16b5a668830800bb8d890708120f35d7947664736f6c63430007060033",
}

// NitroAdjucatorABI is the input ABI used to generate the binding from.
// Deprecated: Use NitroAdjucatorMetaData.ABI instead.
var NitroAdjucatorABI = NitroAdjucatorMetaData.ABI

// NitroAdjucatorBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use NitroAdjucatorMetaData.Bin instead.
var NitroAdjucatorBin = NitroAdjucatorMetaData.Bin

// DeployNitroAdjucator deploys a new Ethereum contract, binding an instance of NitroAdjucator to it.
func DeployNitroAdjucator(auth *bind.TransactOpts, backend bind.ContractBackend) (common.Address, *types.Transaction, *NitroAdjucator, error) {
	parsed, err := NitroAdjucatorMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(NitroAdjucatorBin), backend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &NitroAdjucator{NitroAdjucatorCaller: NitroAdjucatorCaller{contract: contract}, NitroAdjucatorTransactor: NitroAdjucatorTransactor{contract: contract}, NitroAdjucatorFilterer: NitroAdjucatorFilterer{contract: contract}}, nil
}

// NitroAdjucator is an auto generated Go binding around an Ethereum contract.
type NitroAdjucator struct {
	NitroAdjucatorCaller     // Read-only binding to the contract
	NitroAdjucatorTransactor // Write-only binding to the contract
	NitroAdjucatorFilterer   // Log filterer for contract events
}

// NitroAdjucatorCaller is an auto generated read-only Go binding around an Ethereum contract.
type NitroAdjucatorCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjucatorTransactor is an auto generated write-only Go binding around an Ethereum contract.
type NitroAdjucatorTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjucatorFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type NitroAdjucatorFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// NitroAdjucatorSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type NitroAdjucatorSession struct {
	Contract     *NitroAdjucator   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// NitroAdjucatorCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type NitroAdjucatorCallerSession struct {
	Contract *NitroAdjucatorCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// NitroAdjucatorTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type NitroAdjucatorTransactorSession struct {
	Contract     *NitroAdjucatorTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// NitroAdjucatorRaw is an auto generated low-level Go binding around an Ethereum contract.
type NitroAdjucatorRaw struct {
	Contract *NitroAdjucator // Generic contract binding to access the raw methods on
}

// NitroAdjucatorCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type NitroAdjucatorCallerRaw struct {
	Contract *NitroAdjucatorCaller // Generic read-only contract binding to access the raw methods on
}

// NitroAdjucatorTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type NitroAdjucatorTransactorRaw struct {
	Contract *NitroAdjucatorTransactor // Generic write-only contract binding to access the raw methods on
}

// NewNitroAdjucator creates a new instance of NitroAdjucator, bound to a specific deployed contract.
func NewNitroAdjucator(address common.Address, backend bind.ContractBackend) (*NitroAdjucator, error) {
	contract, err := bindNitroAdjucator(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucator{NitroAdjucatorCaller: NitroAdjucatorCaller{contract: contract}, NitroAdjucatorTransactor: NitroAdjucatorTransactor{contract: contract}, NitroAdjucatorFilterer: NitroAdjucatorFilterer{contract: contract}}, nil
}

// NewNitroAdjucatorCaller creates a new read-only instance of NitroAdjucator, bound to a specific deployed contract.
func NewNitroAdjucatorCaller(address common.Address, caller bind.ContractCaller) (*NitroAdjucatorCaller, error) {
	contract, err := bindNitroAdjucator(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorCaller{contract: contract}, nil
}

// NewNitroAdjucatorTransactor creates a new write-only instance of NitroAdjucator, bound to a specific deployed contract.
func NewNitroAdjucatorTransactor(address common.Address, transactor bind.ContractTransactor) (*NitroAdjucatorTransactor, error) {
	contract, err := bindNitroAdjucator(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorTransactor{contract: contract}, nil
}

// NewNitroAdjucatorFilterer creates a new log filterer instance of NitroAdjucator, bound to a specific deployed contract.
func NewNitroAdjucatorFilterer(address common.Address, filterer bind.ContractFilterer) (*NitroAdjucatorFilterer, error) {
	contract, err := bindNitroAdjucator(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorFilterer{contract: contract}, nil
}

// bindNitroAdjucator binds a generic wrapper to an already deployed contract.
func bindNitroAdjucator(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(NitroAdjucatorABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NitroAdjucator *NitroAdjucatorRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NitroAdjucator.Contract.NitroAdjucatorCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NitroAdjucator *NitroAdjucatorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.NitroAdjucatorTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NitroAdjucator *NitroAdjucatorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.NitroAdjucatorTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_NitroAdjucator *NitroAdjucatorCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _NitroAdjucator.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_NitroAdjucator *NitroAdjucatorTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_NitroAdjucator *NitroAdjucatorTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.contract.Transact(opts, method, params...)
}

// ComputeClaimEffectsAndInteractions is a free data retrieval call binding the contract method 0xe29cffe0.
//
// Solidity: function compute_claim_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource, uint256[] targetAllocationIndicesToPayout) pure returns((bytes32,uint256,uint8,bytes)[] newSourceAllocations, (bytes32,uint256,uint8,bytes)[] newTargetAllocations, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjucator *NitroAdjucatorCaller) ComputeClaimEffectsAndInteractions(opts *bind.CallOpts, initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "compute_claim_effects_and_interactions", initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)

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
func (_NitroAdjucator *NitroAdjucatorSession) ComputeClaimEffectsAndInteractions(initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	return _NitroAdjucator.Contract.ComputeClaimEffectsAndInteractions(&_NitroAdjucator.CallOpts, initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)
}

// ComputeClaimEffectsAndInteractions is a free data retrieval call binding the contract method 0xe29cffe0.
//
// Solidity: function compute_claim_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource, uint256[] targetAllocationIndicesToPayout) pure returns((bytes32,uint256,uint8,bytes)[] newSourceAllocations, (bytes32,uint256,uint8,bytes)[] newTargetAllocations, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjucator *NitroAdjucatorCallerSession) ComputeClaimEffectsAndInteractions(initialHoldings *big.Int, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int, targetAllocationIndicesToPayout []*big.Int) (struct {
	NewSourceAllocations []ExitFormatAllocation
	NewTargetAllocations []ExitFormatAllocation
	ExitAllocations      []ExitFormatAllocation
	TotalPayouts         *big.Int
}, error) {
	return _NitroAdjucator.Contract.ComputeClaimEffectsAndInteractions(&_NitroAdjucator.CallOpts, initialHoldings, sourceAllocations, targetAllocations, indexOfTargetInSource, targetAllocationIndicesToPayout)
}

// ComputeTransferEffectsAndInteractions is a free data retrieval call binding the contract method 0x11e9f178.
//
// Solidity: function compute_transfer_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] allocations, uint256[] indices) pure returns((bytes32,uint256,uint8,bytes)[] newAllocations, bool allocatesOnlyZeros, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjucator *NitroAdjucatorCaller) ComputeTransferEffectsAndInteractions(opts *bind.CallOpts, initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "compute_transfer_effects_and_interactions", initialHoldings, allocations, indices)

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
func (_NitroAdjucator *NitroAdjucatorSession) ComputeTransferEffectsAndInteractions(initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	return _NitroAdjucator.Contract.ComputeTransferEffectsAndInteractions(&_NitroAdjucator.CallOpts, initialHoldings, allocations, indices)
}

// ComputeTransferEffectsAndInteractions is a free data retrieval call binding the contract method 0x11e9f178.
//
// Solidity: function compute_transfer_effects_and_interactions(uint256 initialHoldings, (bytes32,uint256,uint8,bytes)[] allocations, uint256[] indices) pure returns((bytes32,uint256,uint8,bytes)[] newAllocations, bool allocatesOnlyZeros, (bytes32,uint256,uint8,bytes)[] exitAllocations, uint256 totalPayouts)
func (_NitroAdjucator *NitroAdjucatorCallerSession) ComputeTransferEffectsAndInteractions(initialHoldings *big.Int, allocations []ExitFormatAllocation, indices []*big.Int) (struct {
	NewAllocations     []ExitFormatAllocation
	AllocatesOnlyZeros bool
	ExitAllocations    []ExitFormatAllocation
	TotalPayouts       *big.Int
}, error) {
	return _NitroAdjucator.Contract.ComputeTransferEffectsAndInteractions(&_NitroAdjucator.CallOpts, initialHoldings, allocations, indices)
}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjucator *NitroAdjucatorCaller) GetChainID(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "getChainID")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjucator *NitroAdjucatorSession) GetChainID() (*big.Int, error) {
	return _NitroAdjucator.Contract.GetChainID(&_NitroAdjucator.CallOpts)
}

// GetChainID is a free data retrieval call binding the contract method 0x564b81ef.
//
// Solidity: function getChainID() pure returns(uint256)
func (_NitroAdjucator *NitroAdjucatorCallerSession) GetChainID() (*big.Int, error) {
	return _NitroAdjucator.Contract.GetChainID(&_NitroAdjucator.CallOpts)
}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjucator *NitroAdjucatorCaller) Holdings(opts *bind.CallOpts, arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "holdings", arg0, arg1)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjucator *NitroAdjucatorSession) Holdings(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _NitroAdjucator.Contract.Holdings(&_NitroAdjucator.CallOpts, arg0, arg1)
}

// Holdings is a free data retrieval call binding the contract method 0x166e56cd.
//
// Solidity: function holdings(address , bytes32 ) view returns(uint256)
func (_NitroAdjucator *NitroAdjucatorCallerSession) Holdings(arg0 common.Address, arg1 [32]byte) (*big.Int, error) {
	return _NitroAdjucator.Contract.Holdings(&_NitroAdjucator.CallOpts, arg0, arg1)
}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorCaller) RequireValidInput(opts *bind.CallOpts, numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "requireValidInput", numParticipants, numStates, numSigs, numWhoSignedWhats)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorSession) RequireValidInput(numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	return _NitroAdjucator.Contract.RequireValidInput(&_NitroAdjucator.CallOpts, numParticipants, numStates, numSigs, numWhoSignedWhats)
}

// RequireValidInput is a free data retrieval call binding the contract method 0xbe5c2a31.
//
// Solidity: function requireValidInput(uint256 numParticipants, uint256 numStates, uint256 numSigs, uint256 numWhoSignedWhats) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorCallerSession) RequireValidInput(numParticipants *big.Int, numStates *big.Int, numSigs *big.Int, numWhoSignedWhats *big.Int) (bool, error) {
	return _NitroAdjucator.Contract.RequireValidInput(&_NitroAdjucator.CallOpts, numParticipants, numStates, numSigs, numWhoSignedWhats)
}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjucator *NitroAdjucatorCaller) StatusOf(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "statusOf", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjucator *NitroAdjucatorSession) StatusOf(arg0 [32]byte) ([32]byte, error) {
	return _NitroAdjucator.Contract.StatusOf(&_NitroAdjucator.CallOpts, arg0)
}

// StatusOf is a free data retrieval call binding the contract method 0xc7df14e2.
//
// Solidity: function statusOf(bytes32 ) view returns(bytes32)
func (_NitroAdjucator *NitroAdjucatorCallerSession) StatusOf(arg0 [32]byte) ([32]byte, error) {
	return _NitroAdjucator.Contract.StatusOf(&_NitroAdjucator.CallOpts, arg0)
}

// UnpackStatus is a free data retrieval call binding the contract method 0x552cfa50.
//
// Solidity: function unpackStatus(bytes32 channelId) view returns(uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
func (_NitroAdjucator *NitroAdjucatorCaller) UnpackStatus(opts *bind.CallOpts, channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "unpackStatus", channelId)

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
func (_NitroAdjucator *NitroAdjucatorSession) UnpackStatus(channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	return _NitroAdjucator.Contract.UnpackStatus(&_NitroAdjucator.CallOpts, channelId)
}

// UnpackStatus is a free data retrieval call binding the contract method 0x552cfa50.
//
// Solidity: function unpackStatus(bytes32 channelId) view returns(uint48 turnNumRecord, uint48 finalizesAt, uint160 fingerprint)
func (_NitroAdjucator *NitroAdjucatorCallerSession) UnpackStatus(channelId [32]byte) (struct {
	TurnNumRecord *big.Int
	FinalizesAt   *big.Int
	Fingerprint   *big.Int
}, error) {
	return _NitroAdjucator.Contract.UnpackStatus(&_NitroAdjucator.CallOpts, channelId)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorCaller) ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	var out []interface{}
	err := _NitroAdjucator.contract.Call(opts, &out, "validTransition", nParticipants, isFinalAB, ab, turnNumB, appDefinition)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjucator.Contract.ValidTransition(&_NitroAdjucator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6775b173.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, (bytes,bytes)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjucator *NitroAdjucatorCallerSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjucator.Contract.ValidTransition(&_NitroAdjucator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "challenge", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Challenge(&_NitroAdjucator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xf198bea9.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Challenge(&_NitroAdjucator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "checkpoint", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Checkpoint(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Checkpoint(&_NitroAdjucator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x0149b762.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, (bytes,bytes)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Checkpoint(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Checkpoint(&_NitroAdjucator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Claim(opts *bind.TransactOpts, claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "claim", claimArgs)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Claim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Claim(&_NitroAdjucator.TransactOpts, claimArgs)
}

// Claim is a paid mutator transaction binding the contract method 0x30776841.
//
// Solidity: function claim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256,uint256[]) claimArgs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Claim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Claim(&_NitroAdjucator.TransactOpts, claimArgs)
}

// Conclude is a paid mutator transaction binding the contract method 0x7ff6a982.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes32 outcomeHash, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Conclude(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeHash [32]byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "conclude", largestTurnNum, fixedPart, appPartHash, outcomeHash, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x7ff6a982.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes32 outcomeHash, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeHash [32]byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Conclude(&_NitroAdjucator.TransactOpts, largestTurnNum, fixedPart, appPartHash, outcomeHash, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x7ff6a982.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes32 outcomeHash, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeHash [32]byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Conclude(&_NitroAdjucator.TransactOpts, largestTurnNum, fixedPart, appPartHash, outcomeHash, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xe03cdb0c.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "concludeAndTransferAllAssets", largestTurnNum, fixedPart, appPartHash, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xe03cdb0c.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjucator.TransactOpts, largestTurnNum, fixedPart, appPartHash, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xe03cdb0c.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes32 appPartHash, bytes outcomeBytes, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appPartHash [32]byte, outcomeBytes []byte, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjucator.TransactOpts, largestTurnNum, fixedPart, appPartHash, outcomeBytes, numStates, whoSignedWhat, sigs)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Deposit(opts *bind.TransactOpts, asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "deposit", asset, channelId, expectedHeld, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjucator *NitroAdjucatorSession) Deposit(asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Deposit(&_NitroAdjucator.TransactOpts, asset, channelId, expectedHeld, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x2fb1d270.
//
// Solidity: function deposit(address asset, bytes32 channelId, uint256 expectedHeld, uint256 amount) payable returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Deposit(asset common.Address, channelId [32]byte, expectedHeld *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Deposit(&_NitroAdjucator.TransactOpts, asset, channelId, expectedHeld, amount)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Respond(opts *bind.TransactOpts, isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "respond", isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Respond(isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Respond(&_NitroAdjucator.TransactOpts, isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xda4cdf73.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Respond(isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Respond(&_NitroAdjucator.TransactOpts, isFinalAB, fixedPart, variablePartAB, sig)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) Transfer(opts *bind.TransactOpts, assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "transfer", assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjucator *NitroAdjucatorSession) Transfer(assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Transfer(&_NitroAdjucator.TransactOpts, assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// Transfer is a paid mutator transaction binding the contract method 0x3033730e.
//
// Solidity: function transfer(uint256 assetIndex, bytes32 fromChannelId, bytes outcomeBytes, bytes32 stateHash, uint256[] indices) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) Transfer(assetIndex *big.Int, fromChannelId [32]byte, outcomeBytes []byte, stateHash [32]byte, indices []*big.Int) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.Transfer(&_NitroAdjucator.TransactOpts, assetIndex, fromChannelId, outcomeBytes, stateHash, indices)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjucator *NitroAdjucatorTransactor) TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjucator.contract.Transact(opts, "transferAllAssets", channelId, outcomeBytes, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjucator *NitroAdjucatorSession) TransferAllAssets(channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.TransferAllAssets(&_NitroAdjucator.TransactOpts, channelId, outcomeBytes, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x6c365522.
//
// Solidity: function transferAllAssets(bytes32 channelId, bytes outcomeBytes, bytes32 stateHash) returns()
func (_NitroAdjucator *NitroAdjucatorTransactorSession) TransferAllAssets(channelId [32]byte, outcomeBytes []byte, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjucator.Contract.TransferAllAssets(&_NitroAdjucator.TransactOpts, channelId, outcomeBytes, stateHash)
}

// NitroAdjucatorAllocationUpdatedIterator is returned from FilterAllocationUpdated and is used to iterate over the raw logs and unpacked data for AllocationUpdated events raised by the NitroAdjucator contract.
type NitroAdjucatorAllocationUpdatedIterator struct {
	Event *NitroAdjucatorAllocationUpdated // Event containing the contract specifics and raw log

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
func (it *NitroAdjucatorAllocationUpdatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjucatorAllocationUpdated)
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
		it.Event = new(NitroAdjucatorAllocationUpdated)
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
func (it *NitroAdjucatorAllocationUpdatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjucatorAllocationUpdatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjucatorAllocationUpdated represents a AllocationUpdated event raised by the NitroAdjucator contract.
type NitroAdjucatorAllocationUpdated struct {
	ChannelId       [32]byte
	AssetIndex      *big.Int
	InitialHoldings *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAllocationUpdated is a free log retrieval operation binding the contract event 0xb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings)
func (_NitroAdjucator *NitroAdjucatorFilterer) FilterAllocationUpdated(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjucatorAllocationUpdatedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.FilterLogs(opts, "AllocationUpdated", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorAllocationUpdatedIterator{contract: _NitroAdjucator.contract, event: "AllocationUpdated", logs: logs, sub: sub}, nil
}

// WatchAllocationUpdated is a free log subscription operation binding the contract event 0xb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings)
func (_NitroAdjucator *NitroAdjucatorFilterer) WatchAllocationUpdated(opts *bind.WatchOpts, sink chan<- *NitroAdjucatorAllocationUpdated, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.WatchLogs(opts, "AllocationUpdated", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjucatorAllocationUpdated)
				if err := _NitroAdjucator.contract.UnpackLog(event, "AllocationUpdated", log); err != nil {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) ParseAllocationUpdated(log types.Log) (*NitroAdjucatorAllocationUpdated, error) {
	event := new(NitroAdjucatorAllocationUpdated)
	if err := _NitroAdjucator.contract.UnpackLog(event, "AllocationUpdated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjucatorChallengeClearedIterator is returned from FilterChallengeCleared and is used to iterate over the raw logs and unpacked data for ChallengeCleared events raised by the NitroAdjucator contract.
type NitroAdjucatorChallengeClearedIterator struct {
	Event *NitroAdjucatorChallengeCleared // Event containing the contract specifics and raw log

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
func (it *NitroAdjucatorChallengeClearedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjucatorChallengeCleared)
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
		it.Event = new(NitroAdjucatorChallengeCleared)
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
func (it *NitroAdjucatorChallengeClearedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjucatorChallengeClearedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjucatorChallengeCleared represents a ChallengeCleared event raised by the NitroAdjucator contract.
type NitroAdjucatorChallengeCleared struct {
	ChannelId        [32]byte
	NewTurnNumRecord *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterChallengeCleared is a free log retrieval operation binding the contract event 0x07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0.
//
// Solidity: event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)
func (_NitroAdjucator *NitroAdjucatorFilterer) FilterChallengeCleared(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjucatorChallengeClearedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.FilterLogs(opts, "ChallengeCleared", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorChallengeClearedIterator{contract: _NitroAdjucator.contract, event: "ChallengeCleared", logs: logs, sub: sub}, nil
}

// WatchChallengeCleared is a free log subscription operation binding the contract event 0x07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0.
//
// Solidity: event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord)
func (_NitroAdjucator *NitroAdjucatorFilterer) WatchChallengeCleared(opts *bind.WatchOpts, sink chan<- *NitroAdjucatorChallengeCleared, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.WatchLogs(opts, "ChallengeCleared", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjucatorChallengeCleared)
				if err := _NitroAdjucator.contract.UnpackLog(event, "ChallengeCleared", log); err != nil {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) ParseChallengeCleared(log types.Log) (*NitroAdjucatorChallengeCleared, error) {
	event := new(NitroAdjucatorChallengeCleared)
	if err := _NitroAdjucator.contract.UnpackLog(event, "ChallengeCleared", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjucatorChallengeRegisteredIterator is returned from FilterChallengeRegistered and is used to iterate over the raw logs and unpacked data for ChallengeRegistered events raised by the NitroAdjucator contract.
type NitroAdjucatorChallengeRegisteredIterator struct {
	Event *NitroAdjucatorChallengeRegistered // Event containing the contract specifics and raw log

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
func (it *NitroAdjucatorChallengeRegisteredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjucatorChallengeRegistered)
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
		it.Event = new(NitroAdjucatorChallengeRegistered)
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
func (it *NitroAdjucatorChallengeRegisteredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjucatorChallengeRegisteredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjucatorChallengeRegistered represents a ChallengeRegistered event raised by the NitroAdjucator contract.
type NitroAdjucatorChallengeRegistered struct {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) FilterChallengeRegistered(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjucatorChallengeRegisteredIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.FilterLogs(opts, "ChallengeRegistered", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorChallengeRegisteredIterator{contract: _NitroAdjucator.contract, event: "ChallengeRegistered", logs: logs, sub: sub}, nil
}

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0xf6c285d62578fdf94d2e5c698650728f1d64a497add9bba112b4ac4d5c489cee.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (bytes,bytes)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
func (_NitroAdjucator *NitroAdjucatorFilterer) WatchChallengeRegistered(opts *bind.WatchOpts, sink chan<- *NitroAdjucatorChallengeRegistered, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.WatchLogs(opts, "ChallengeRegistered", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjucatorChallengeRegistered)
				if err := _NitroAdjucator.contract.UnpackLog(event, "ChallengeRegistered", log); err != nil {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) ParseChallengeRegistered(log types.Log) (*NitroAdjucatorChallengeRegistered, error) {
	event := new(NitroAdjucatorChallengeRegistered)
	if err := _NitroAdjucator.contract.UnpackLog(event, "ChallengeRegistered", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjucatorConcludedIterator is returned from FilterConcluded and is used to iterate over the raw logs and unpacked data for Concluded events raised by the NitroAdjucator contract.
type NitroAdjucatorConcludedIterator struct {
	Event *NitroAdjucatorConcluded // Event containing the contract specifics and raw log

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
func (it *NitroAdjucatorConcludedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjucatorConcluded)
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
		it.Event = new(NitroAdjucatorConcluded)
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
func (it *NitroAdjucatorConcludedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjucatorConcludedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjucatorConcluded represents a Concluded event raised by the NitroAdjucator contract.
type NitroAdjucatorConcluded struct {
	ChannelId   [32]byte
	FinalizesAt *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterConcluded is a free log retrieval operation binding the contract event 0x4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901.
//
// Solidity: event Concluded(bytes32 indexed channelId, uint48 finalizesAt)
func (_NitroAdjucator *NitroAdjucatorFilterer) FilterConcluded(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjucatorConcludedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.FilterLogs(opts, "Concluded", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorConcludedIterator{contract: _NitroAdjucator.contract, event: "Concluded", logs: logs, sub: sub}, nil
}

// WatchConcluded is a free log subscription operation binding the contract event 0x4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901.
//
// Solidity: event Concluded(bytes32 indexed channelId, uint48 finalizesAt)
func (_NitroAdjucator *NitroAdjucatorFilterer) WatchConcluded(opts *bind.WatchOpts, sink chan<- *NitroAdjucatorConcluded, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjucator.contract.WatchLogs(opts, "Concluded", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjucatorConcluded)
				if err := _NitroAdjucator.contract.UnpackLog(event, "Concluded", log); err != nil {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) ParseConcluded(log types.Log) (*NitroAdjucatorConcluded, error) {
	event := new(NitroAdjucatorConcluded)
	if err := _NitroAdjucator.contract.UnpackLog(event, "Concluded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// NitroAdjucatorDepositedIterator is returned from FilterDeposited and is used to iterate over the raw logs and unpacked data for Deposited events raised by the NitroAdjucator contract.
type NitroAdjucatorDepositedIterator struct {
	Event *NitroAdjucatorDeposited // Event containing the contract specifics and raw log

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
func (it *NitroAdjucatorDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjucatorDeposited)
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
		it.Event = new(NitroAdjucatorDeposited)
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
func (it *NitroAdjucatorDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjucatorDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjucatorDeposited represents a Deposited event raised by the NitroAdjucator contract.
type NitroAdjucatorDeposited struct {
	Destination         [32]byte
	Asset               common.Address
	AmountDeposited     *big.Int
	DestinationHoldings *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
func (_NitroAdjucator *NitroAdjucatorFilterer) FilterDeposited(opts *bind.FilterOpts, destination [][32]byte) (*NitroAdjucatorDepositedIterator, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _NitroAdjucator.contract.FilterLogs(opts, "Deposited", destinationRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjucatorDepositedIterator{contract: _NitroAdjucator.contract, event: "Deposited", logs: logs, sub: sub}, nil
}

// WatchDeposited is a free log subscription operation binding the contract event 0x2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 amountDeposited, uint256 destinationHoldings)
func (_NitroAdjucator *NitroAdjucatorFilterer) WatchDeposited(opts *bind.WatchOpts, sink chan<- *NitroAdjucatorDeposited, destination [][32]byte) (event.Subscription, error) {

	var destinationRule []interface{}
	for _, destinationItem := range destination {
		destinationRule = append(destinationRule, destinationItem)
	}

	logs, sub, err := _NitroAdjucator.contract.WatchLogs(opts, "Deposited", destinationRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjucatorDeposited)
				if err := _NitroAdjucator.contract.UnpackLog(event, "Deposited", log); err != nil {
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
func (_NitroAdjucator *NitroAdjucatorFilterer) ParseDeposited(log types.Log) (*NitroAdjucatorDeposited, error) {
	event := new(NitroAdjucatorDeposited)
	if err := _NitroAdjucator.contract.UnpackLog(event, "Deposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

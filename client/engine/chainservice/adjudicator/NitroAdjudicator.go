// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

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

// ExitFormatSingleAssetExit is an auto generated low-level Go binding around an user-defined struct.
type ExitFormatSingleAssetExit struct {
	Asset       common.Address
	Metadata    []byte
	Allocations []ExitFormatAllocation
}

// IForceMoveAppVariablePart is an auto generated low-level Go binding around an user-defined struct.
type IForceMoveAppVariablePart struct {
	Outcome []ExitFormatSingleAssetExit
	AppData []byte
	TurnNum *big.Int
	IsFinal bool
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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"isFinalCount\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"name\":\"compute_claim_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newSourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newTargetAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint48\",\"name\":\"largestTurnNum\",\"type\":\"uint48\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numParticipants\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numSigs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numWhoSignedWhats\",\"type\":\"uint256\"}],\"name\":\"requireValidInput\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"variablePartAB\",\"type\":\"tuple[2]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"sig\",\"type\":\"tuple\"}],\"name\":\"respond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nParticipants\",\"type\":\"uint256\"},{\"internalType\":\"bool[2]\",\"name\":\"isFinalAB\",\"type\":\"bool[2]\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"ab\",\"type\":\"tuple[2]\"},{\"internalType\":\"uint48\",\"name\":\"turnNumB\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"}],\"name\":\"validTransition\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614ba5806100206000396000f3fe6080604052600436106100fe5760003560e01c80636dfa8b7411610095578063be5c2a3111610064578063be5c2a31146102af578063c7df14e2146102cf578063dd5b9932146102ef578063e29cffe01461030f578063fb6d51391461033f576100fe565b80636dfa8b74146102225780639ae6008b1461024f578063a26d793a1461026f578063af69c9d71461028f576100fe565b806330776841116100d1578063307768411461019e578063552cfa50146101be578063564b81ef146101ed5780635973edd914610202576100fe565b806311e9f17814610103578063166e56cd1461013c5780632fb1d270146101695780633033730e1461017e575b600080fd5b34801561010f57600080fd5b5061012361011e366004613be7565b61035f565b6040516101339493929190614118565b60405180910390f35b34801561014857600080fd5b5061015c610157366004613540565b6106ad565b6040516101339190614163565b61017c61017736600461356b565b6106ca565b005b34801561018a57600080fd5b5061017c610199366004613c50565b610959565b3480156101aa57600080fd5b5061017c6101b936600461383f565b6109d9565b3480156101ca57600080fd5b506101de6101d93660046137da565b610ab9565b60405161013393929190614a65565b3480156101f957600080fd5b5061015c610ad4565b34801561020e57600080fd5b5061017c61021d3660046135a5565b610ad8565b34801561022e57600080fd5b5061024261023d366004613ad5565b610c2b565b6040516101339190614158565b34801561025b57600080fd5b5061017c61026a366004613cc8565b610c46565b34801561027b57600080fd5b5061017c61028a366004613933565b610c5f565b34801561029b57600080fd5b5061017c6102aa3660046137f2565b610cae565b3480156102bb57600080fd5b506102426102ca366004613c97565b61104f565b3480156102db57600080fd5b5061015c6102ea3660046137da565b6110d4565b3480156102fb57600080fd5b5061017c61030a3660046139fa565b6110e6565b34801561031b57600080fd5b5061032f61032a366004613b50565b61125a565b60405161013394939291906140cd565b34801561034b57600080fd5b5061017c61035a366004613cc8565b6117d1565b606060006060600080855111610376578551610379565b84515b6001600160401b038111801561038e57600080fd5b506040519080825280602002602001820160405280156103c857816020015b6103b5612d23565b8152602001906001900390816103ad5790505b5091506000905085516001600160401b03811180156103e657600080fd5b5060405190808252806020026020018201604052801561042057816020015b61040d612d23565b8152602001906001900390816104055790505b50935060019250866000805b88518110156106a15788818151811061044157fe5b60200260200101516000015187828151811061045957fe5b6020026020010151600001818152505088818151811061047557fe5b60200260200101516040015187828151811061048d57fe5b60200260200101516040019060ff16908160ff16815250508881815181106104b157fe5b6020026020010151606001518782815181106104c957fe5b60200260200101516060018190525060006104fb8a83815181106104e957fe5b602002602001015160200151856117f0565b905088516000148061052a575088518310801561052a57508189848151811061052057fe5b6020026020010151145b1561063c57600260ff168a848151811061054057fe5b60200260200101516040015160ff1614156105765760405162461bcd60e51b815260040161056d906144a6565b60405180910390fd5b808a838151811061058357fe5b6020026020010151602001510388838151811061059c57fe5b6020026020010151602001818152505060405180608001604052808b84815181106105c357fe5b60200260200101516000015181526020018281526020018b84815181106105e657fe5b60200260200101516040015160ff1681526020018b848151811061060657fe5b60200260200101516060015181525086848151811061062157fe5b60200260200101819052508085019450826001019250610671565b89828151811061064857fe5b60200260200101516020015188838151811061066057fe5b602002602001015160200181815250505b87828151811061067d57fe5b60200260200101516020015160001461069557600096505b9092039160010161042c565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b6106d38361180a565b156106f05760405162461bcd60e51b815260040161056d906146eb565b6001600160a01b0384166000908152600160209081526040808320868452909152812054838110156107345760405162461bcd60e51b815260040161056d9061421a565b61073e8484611816565b811061075c5760405162461bcd60e51b815260040161056d906147f0565b6107708161076a8686611816565b90611870565b91506001600160a01b0386166107a45782341461079f5760405162461bcd60e51b815260040161056d9061485e565b610842565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906107d49033903090879060040161406f565b602060405180830381600087803b1580156107ee57600080fd5b505af1158015610802573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061082691906137be565b6108425760405162461bcd60e51b815260040161056d9061456e565b600061084e8284611816565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715906108ad908a90879086906140ac565b60405180910390a26001600160a01b0387166109505760006108cf8585611870565b90506000336001600160a01b0316826040516108ea90612af7565b60006040518083038185875af1925050503d8060008114610927576040519150601f19603f3d011682016040523d82523d6000602084013e61092c565b606091505b505090508061094d5760405162461bcd60e51b815260040161056d906148c1565b50505b50505050505050565b606060008061096b888589888a6118cd565b925092509250606080600061099884878d8151811061098657fe5b6020026020010151604001518961035f565b935093505092506109af8b868c8b8a888a8861194b565b6109cc868c815181106109be57fe5b6020026020010151836119f7565b5050505050505050505050565b6060806000806109e885611a2f565b9350935093509350606080606060006060888a6060015181518110610a0957fe5b60200260200101516040015190506060888b60e0015181518110610a2957fe5b6020026020010151604001519050610a4d8783838e608001518f610100015161125a565b809650819750829850839950505050505050610aa689878a878c8e6060015181518110610a7657fe5b6020026020010151604001518e6080015181518110610a9157fe5b6020026020010151600001518c898c89611bd8565b61094d878a60e00151815181106109be57fe5b6000806000610ac784611cdb565b9196909550909350915050565b4690565b6000610ae384611cf9565b9050600080610af183611cdb565b50865160208101519051929450909250600091610b18918691868c865b6020020151611d6f565b602087810151908101519051919250600091610b3d9187916001888101908e90610b0e565b875151909150600090610b4f90611dab565b805190602001209050610b9260405180608001604052808765ffffffffffff1681526020018665ffffffffffff1681526020018581526020018381525087611dd4565b6020890151805165ffffffffffff600188011681610bac57fe5b0681518110610bb757fe5b60200260200101516001600160a01b0316610bd28389611de7565b6001600160a01b031614610bf85760405162461bcd60e51b815260040161056d90614508565b610c118960200151518b8a886001018d60600151611e55565b50610c1f8686600101611f28565b50505050505050505050565b6000610c3a8686868686611e55565b90505b95945050505050565b610c5587878787878787611faf565b5050505050505050565b610c7386602001515185518451845161104f565b506000610c7f87611cf9565b9050610c8a81612181565b610c9481876121b8565b610ca3868686848b88886121fb565b506109508187611f28565b610cb783612261565b610cd181610cc484611dab565b8051906020012085612294565b81516001906060906001600160401b0381118015610cee57600080fd5b50604051908082528060200260200182016040528015610d2857816020015b610d15612d49565b815260200190600190039081610d0d5790505b509050606084516001600160401b0381118015610d4457600080fd5b50604051908082528060200260200182016040528015610d6e578160200160208202803683370190505b509050606085516001600160401b0381118015610d8a57600080fd5b50604051908082528060200260200182016040528015610db4578160200160208202803683370190505b50905060005b8651811015610f5657610dcb612d49565b878281518110610dd757fe5b602002602001015190506060816040015190506000898481518110610df857fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c815260200190815260200160002054868581518110610e4957fe5b6020026020010181815250506060600060606000610ebc8a8981518110610e6c57fe5b60200260200101518760006001600160401b0381118015610e8c57600080fd5b50604051908082528060200260200182016040528015610eb6578160200160208202803683370190505b5061035f565b935093509350935082610ece5760009b505b80898981518110610edb57fe5b602002602001018181525050838e8981518110610ef457fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b8981518110610f3757fe5b6020026020010181905250505050505050508080600101915050610dba565b5060005b865181101561100a576000878281518110610f7157fe5b6020026020010151600001519050828281518110610f8b57fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d835290935291909120805491909103905583518990600080516020614b50833981519152908490879082908110610fe357fe5b6020026020010151604051610ff9929190614984565b60405180910390a250600101610f5a565b50831561102557600087815260208190526040812055611046565b600061103087611dab565b8051906020012090506110448887836122e3565b505b61095083612349565b60008385101580156110615750600084115b61107d5760405162461bcd60e51b815260040161056d90614408565b848314801561108b57508482145b6110a75760405162461bcd60e51b815260040161056d906145a5565b60ff8511156110c85760405162461bcd60e51b815260040161056d9061443f565b5060015b949350505050565b60006020819052908152604090205481565b6110fa87602001515186518551855161104f565b50600061110688611cf9565b9050600061111382612379565b600281111561111e57fe5b14156111335761112e81886123c3565b611162565b600161113e82612379565b600281111561114957fe5b14156111595761112e81886121b8565b61116281612181565b6000611173888888858d8a8a6121fb565b9050611184818a6020015185612402565b817fbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f38898b60800151420160008a60ff16118d8c8b8b6040516111cc97969594939291906149a5565b60405180910390a261123d60405180608001604052808a65ffffffffffff1681526020018b60800151420165ffffffffffff16815260200183815260200161122e8a60018c51038151811061121d57fe5b602002602001015160000151611dab565b8051906020012081525061245c565b600092835260208390526040909220919091555050505050505050565b606080606060008088516001600160401b038111801561127957600080fd5b506040519080825280602002602001820160405280156112b357816020015b6112a0612d23565b8152602001906001900390816112985790505b50945087516001600160401b03811180156112cd57600080fd5b5060405190808252806020026020018201604052801561130757816020015b6112f4612d23565b8152602001906001900390816112ec5790505b50935087516001600160401b038111801561132157600080fd5b5060405190808252806020026020018201604052801561135b57816020015b611348612d23565b8152602001906001900390816113405790505b50925060005b89518110156114475789818151811061137657fe5b60200260200101516000015186828151811061138e57fe5b602002602001015160000181815250508981815181106113aa57fe5b6020026020010151602001518682815181106113c257fe5b602002602001015160200181815250508981815181106113de57fe5b6020026020010151606001518682815181106113f657fe5b60200260200101516060018190525089818151811061141157fe5b60200260200101516040015186828151811061142957fe5b602090810291909101015160ff909116604090910152600101611361565b5060005b88518110156115f25788818151811061146057fe5b60200260200101516000015185828151811061147857fe5b6020026020010151600001818152505088818151811061149457fe5b6020026020010151602001518582815181106114ac57fe5b602002602001015160200181815250508881815181106114c857fe5b6020026020010151606001518582815181106114e057fe5b6020026020010151606001819052508881815181106114fb57fe5b60200260200101516040015185828151811061151357fe5b60200260200101516040019060ff16908160ff168152505088818151811061153757fe5b60200260200101516000015184828151811061154f57fe5b60200260200101516000018181525050600084828151811061156d57fe5b6020026020010151602001818152505088818151811061158957fe5b6020026020010151606001518482815181106115a157fe5b6020026020010151606001819052508881815181106115bc57fe5b6020026020010151604001518482815181106115d457fe5b602090810291909101015160ff90911660409091015260010161144b565b508960005b88811015611639578161160957611639565b600061162c8c838151811061161a57fe5b602002602001015160200151846117f0565b90920391506001016115f7565b50600061165d828c8b8151811061164c57fe5b6020026020010151602001516117f0565b905060606116818c8b8151811061167057fe5b6020026020010151606001516124ad565b905060005b81518110156117c05782611699576117c0565b60005b88518110156117b757836116af576117b7565b8881815181106116bb57fe5b6020026020010151600001518383815181106116d357fe5b602002602001015114156117af5760006117048e83815181106116f257fe5b602002602001015160200151866117f0565b905080850394508b516000148061173857508b51871080156117385750818c888151811061172e57fe5b6020026020010151145b156117a957808a838151811061174a57fe5b60200260200101516020018181510391508181525050808b8e8151811061176d57fe5b602002602001015160200181815103915081815250508089838151811061179057fe5b6020908102919091018101510152968701966001909601955b506117b7565b60010161169c565b50600101611686565b505050505095509550955095915050565b60006117e288888888888888611faf565b9050610c5581866000610cae565b60008183116117ff5782611801565b815b90505b92915050565b60a081901c155b919050565b600082820183811015611801576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b6000828211156118c7576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60606000806118db876124c3565b6118e486612261565b6118f685858051906020012088612294565b6118ff84612522565b925082888151811061190d57fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a90811061198557fe5b6020026020010151604001819052506119c58686866040516020016119aa9190614145565b604051602081830303815290604052805190602001206122e3565b85600080516020614b5083398151915289846040516119e5929190614984565b60405180910390a25050505050505050565b611a2b604051806060016040528084600001516001600160a01b031681526020018460200151815260200183815250612538565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611a63906124c3565b611a6c85612261565b611a828a60200151858051906020012087612294565b611a8b84612522565b9850611a9682612522565b9750888381518110611aa457fe5b6020908102919091010151519650600260ff16898481518110611ac357fe5b6020026020010151604001518281518110611ada57fe5b60200260200101516040015160ff1614611b065760405162461bcd60e51b815260040161056d90614827565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a9085908110611b3b57fe5b6020026020010151604001518b6080015181518110611b5657fe5b6020026020010151600001519050876001600160a01b0316898381518110611b7a57fe5b6020026020010151600001516001600160a01b031614611bac5760405162461bcd60e51b815260040161056d9061446f565b611bb581612261565b611bcb8b60a00151848051906020012083612294565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b9084908110611c1e57fe5b602002602001015160400181905250611c47838d602001518c6040516020016119aa9190614145565b85878281518110611c5457fe5b602002602001015160400181905250611c7d888d60a00151896040516020016119aa9190614145565b82600080516020614b508339815191528387604051611c9d929190614984565b60405180910390a287600080516020614b508339815191528287604051611cc5929190614984565b60405180910390a2505050505050505050505050565b60009081526020819052604090205460d081901c9160a082901c9190565b6000611d03610ad4565b825114611d225760405162461bcd60e51b815260040161056d90614251565b611d2a610ad4565b8260200151836040015184606001518560800151604051602001611d52959493929190614939565b604051602081830303815290604052805190602001209050919050565b60008585858585604051602001611d8a95949392919061416c565b60405160208183030381529060405280519060200120905095945050505050565b606081604051602001611dbe9190614145565b6040516020818303038152906040529050919050565b611dde82826125e4565b611a2b8161261a565b60008083604051602001611dfb919061403e565b6040516020818303038152906040528051906020012090506000611e2d8285600001518660200151876040015161264d565b90506001600160a01b0381166110cc5760405162461bcd60e51b815260040161056d906144dd565b600080611e64878787876126f2565b90506001816001811115611e7457fe5b1415611f1b57845160208601516040516345cb2b6960e01b81526001600160a01b038616926345cb2b6992611eaf9289908d906004016148f8565b60206040518083038186803b158015611ec757600080fd5b505afa158015611edb573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190611eff91906137be565b611f1b5760405162461bcd60e51b815260040161056d906142df565b5060019695505050505050565b6040805160808101825265ffffffffffff831681526000602082018190529181018290526060810191909152611f5d9061245c565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e082604051611fa39190614992565b60405180910390a25050565b6000611fba87611cf9565b9050611fc581612181565b611fdb8760200151518560ff168451865161104f565b508360ff168860010165ffffffffffff16101561200a5760405162461bcd60e51b815260040161056d90614722565b60608460ff166001600160401b038111801561202557600080fd5b5060405190808252806020026020018201604052801561204f578160200160208202803683370190505b50905060005b8560ff168165ffffffffffff1610156120a7576120808389898960ff16856001018f01036001611d6f565b828265ffffffffffff168151811061209457fe5b6020908102919091010152600101612055565b506120b98989602001518386886127d6565b6120d55760405162461bcd60e51b815260040161056d906146b4565b60006120e087611dab565b8051906020012090506121266040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b81526020018381525061245c565b60008085815260200190815260200160002081905550827f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc9014260405161216c9190614992565b60405180910390a25050979650505050505050565b600261218c82612379565b600281111561219757fe5b14156121b55760405162461bcd60e51b815260040161056d9061427c565b50565b60006121c383611cdb565b505090508065ffffffffffff168265ffffffffffff16116121f65760405162461bcd60e51b815260040161056d906141e3565b505050565b6000606061220c89898989896128ac565b905061221f8986602001518387876127d6565b61223b5760405162461bcd60e51b815260040161056d906143dc565b8060018251038151811061224b57fe5b6020026020010151915050979650505050505050565b600261226c82612379565b600281111561227757fe5b146121b55760405162461bcd60e51b815260040161056d90614752565b600061229f82611cdb565b925050506122ad8484612a75565b6001600160a01b0316816001600160a01b0316146122dd5760405162461bcd60e51b815260040161056d9061453f565b50505050565b6000806122ef85611cdb565b5091509150600061232f60405180608001604052808565ffffffffffff1681526020018465ffffffffffff1681526020018781526020018681525061245c565b600096875260208790526040909620959095555050505050565b60005b8151811015611a2b5761237182828151811061236457fe5b6020026020010151612538565b60010161234c565b60008061238583611cdb565b5091505065ffffffffffff81166123a0576000915050611811565b428165ffffffffffff16116123b9576002915050611811565b6001915050611811565b60006123ce83611cdb565b505090508065ffffffffffff168265ffffffffffff1610156121f65760405162461bcd60e51b815260040161056d906142a8565b60006124348460405160200161241891906141b9565b6040516020818303038152906040528051906020012083611de7565b90506124408184612aa1565b6122dd5760405162461bcd60e51b815260040161056d90614375565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161249991612a75565b6001600160a01b0316919091179392505050565b6060818060200190518101906118049190613628565b60005b8151816001011015611a2b578181600101815181106124e157fe5b60200260200101518282815181106124f557fe5b60200260200101511061251a5760405162461bcd60e51b815260040161056d90614345565b6001016124c6565b60608180602001905181019061180491906136b7565b805160005b8260400151518110156121f65760008360400151828151811061255c57fe5b602002602001015160000151905060008460400151838151811061257c57fe5b60200260200101516020015190506125938261180a565b156125b0576125ab846125a584612af7565b83612afa565b6125da565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b505060010161253d565b6000818152602081905260409020546125fe908390612c0a565b611a2b5760405162461bcd60e51b815260040161056d90614646565b600161262582612379565b600281111561263057fe5b146121b55760405162461bcd60e51b815260040161056d90614316565b6000601b8460ff16101561266257601b840193505b8360ff16601b1415801561267a57508360ff16601c14155b15612687575060006110cc565b60018585858560405160008152602001604052604051808581526020018460ff1681526020018381526020018281526020019450505050506020604051602081039080840390855afa1580156126e1573d6000803e3d6000fd5b5050506020604051035190506110cc565b602083015160009015612735576127148360015b602002015151845151612c1e565b6127305760405162461bcd60e51b815260040161056d90614782565b6127cb565b8351156127545760405162461bcd60e51b815260040161056d90614895565b846002028265ffffffffffff1610156127c357612772836001612706565b61278e5760405162461bcd60e51b815260040161056d9061460f565b6020838101518101518451909101516127a79190612c35565b6127305760405162461bcd60e51b815260040161056d9061467d565b5060016110cc565b506000949350505050565b83518351600091906127ea84898484612c99565b6128065760405162461bcd60e51b815260040161056d906145da565b60005b8281101561289d5760006128598887848151811061282357fe5b602002602001015160ff168151811061283857fe5b602002602001015188848151811061284c57fe5b6020026020010151611de7565b905088828151811061286757fe5b60200260200101516001600160a01b0316816001600160a01b031614612894576000945050505050610c3d565b50600101612809565b50600198975050505050505050565b60608085516001600160401b03811180156128c657600080fd5b506040519080825280602002602001820160405280156128f0578160200160208202803683370190505b509050600160ff86168803016000805b88518165ffffffffffff161015612a67578089518b0360010101915061297a878a8365ffffffffffff168151811061293457fe5b6020026020010151602001518b8465ffffffffffff168151811061295457fe5b602002602001015160000151858765ffffffffffff168765ffffffffffff161015611d6f565b848265ffffffffffff168151811061298e57fe5b6020026020010181815250508965ffffffffffff168265ffffffffffff161015612a5f57612a5d86602001515160405180604001604052808665ffffffffffff168665ffffffffffff1610151515151581526020018665ffffffffffff168660010165ffffffffffff1610151515151581525060405180604001604052808d8665ffffffffffff1681518110612a2057fe5b602002602001015181526020018d8660010165ffffffffffff1681518110612a4457fe5b6020026020010151815250856001018a60600151611e55565b505b600101612900565b509198975050505050505050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b6000805b8251811015612aed57828181518110612aba57fe5b60200260200101516001600160a01b0316846001600160a01b03161415612ae5576001915050611804565b600101612aa5565b5060009392505050565b90565b6001600160a01b038316612b8a576000826001600160a01b031682604051612b2190612af7565b60006040518083038185875af1925050503d8060008114612b5e576040519150601f19603f3d011682016040523d82523d6000602084013e612b63565b606091505b5050905080612b845760405162461bcd60e51b815260040161056d906143ac565b506121f6565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb90612bb89085908590600401614093565b602060405180830381600087803b158015612bd257600080fd5b505af1158015612be6573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122dd91906137be565b600081612c168461245c565b149392505050565b6000611801612c2c84611dab565b612c3584611dab565b815181516000916001918114808314612c515760009250612c8f565b600160208701838101602088015b600284838510011415612c8a578051835114612c7e5760009650600093505b60209283019201612c5f565b505050505b5090949350505050565b600082855114612cbb5760405162461bcd60e51b815260040161056d906147b9565b60005b83811015612d1757600084828765ffffffffffff1687010381612cdd57fe5b0690508381888481518110612cee57fe5b602002602001015160ff16016001011015612d0e576000925050506110cc565b50600101612cbe565b50600195945050505050565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b803561181181614b1d565b600082601f830112612d8e578081fd5b604051604081018181106001600160401b0382111715612daa57fe5b80604052508091508284604085011115612dc357600080fd5b60005b6002811015612def578135612dda81614b32565b83526020928301929190910190600101612dc6565b50505092915050565b600082601f830112612e08578081fd5b8135612e1b612e1682614ab3565b614a90565b818152915060208083019084810160005b84811015612ec75781358701608080601f19838c03011215612e4d57600080fd5b604080518281016001600160401b038282108183111715612e6a57fe5b9083528488013582528483013582890152606090612e8982870161352a565b83850152938501359380851115612e9f57600080fd5b50612eae8d898688010161327b565b9082015287525050509282019290820190600101612e2c565b505050505092915050565b600082601f830112612ee2578081fd5b8151612ef0612e1682614ab3565b818152915060208083019084810160005b84811015612ec75781518701608080601f19838c03011215612f2257600080fd5b604080518281016001600160401b038282108183111715612f3f57fe5b9083528488015182528483015182890152606090612f5e828701613535565b83850152938501519380851115612f7457600080fd5b50612f838d89868801016132c9565b9082015287525050509282019290820190600101612f01565b600082601f830112612fac578081fd5b8135612fba612e1682614ab3565b8181529150602080830190848101606080850287018301881015612fdd57600080fd5b60005b8581101561300457612ff28984613413565b85529383019391810191600101612fe0565b50505050505092915050565b600082601f830112613020578081fd5b813561302e612e1682614ab3565b818152915060208083019084810160005b84811015612ec75781358701606080601f19838c0301121561306057600080fd5b604080518281016001600160401b03828210818311171561307d57fe5b908352848801359061308e82614b1d565b9082528483013590808211156130a357600080fd5b6130b18e8a8489010161327b565b838a01529385013593808511156130c757600080fd5b50506130d78c8885870101612df8565b9181019190915286525050928201929082019060010161303f565b600082601f830112613102578081fd5b604051604081018181106001600160401b038211171561311e57fe5b6040529050808260005b6002811015612def5761313e8683358701613470565b83526020928301929190910190600101613128565b600082601f830112613163578081fd5b8135613171612e1682614ab3565b818152915060208083019084810160005b84811015612ec757613199888484358a0101613470565b84529282019290820190600101613182565b600082601f8301126131bb578081fd5b81356131c9612e1682614ab3565b8181529150602080830190848101818402860182018710156131ea57600080fd5b60005b84811015612ec7578135845292820192908201906001016131ed565b600082601f830112613219578081fd5b8135613227612e1682614ab3565b81815291506020808301908481018184028601820187101561324857600080fd5b60005b84811015612ec757813561325e81614b40565b8452928201929082019060010161324b565b803561181181614b32565b600082601f83011261328b578081fd5b8135613299612e1682614ad0565b91508082528360208285010111156132b057600080fd5b8060208401602084013760009082016020015292915050565b600082601f8301126132d9578081fd5b81516132e7612e1682614ad0565b91508082528360208285010111156132fe57600080fd5b61330f816020840160208601614af1565b5092915050565b600060a08284031215613327578081fd5b60405160a081016001600160401b03828210818311171561334457fe5b816040528293508435835260209150818501358181111561336457600080fd5b85019050601f8101861361337757600080fd5b8035613385612e1682614ab3565b81815283810190838501858402850186018a10156133a257600080fd5b600094505b838510156133ce5780356133ba81614b1d565b8352600194909401939185019185016133a7565b50808587015250505050506133e560408401613514565b60408201526133f660608401612d73565b606082015261340760808401613514565b60808201525092915050565b600060608284031215613424578081fd5b604051606081018181106001600160401b038211171561344057fe5b604052905080823561345181614b40565b8082525060208301356020820152604083013560408201525092915050565b600060808284031215613481578081fd5b604051608081016001600160401b03828210818311171561349e57fe5b8160405282935084359150808211156134b657600080fd5b6134c286838701613010565b835260208501359150808211156134d857600080fd5b506134e58582860161327b565b6020830152506134f760408401613514565b604082015261350860608401613270565b60608201525092915050565b803565ffffffffffff8116811461181157600080fd5b803561181181614b40565b805161181181614b40565b60008060408385031215613552578182fd5b823561355d81614b1d565b946020939093013593505050565b60008060008060808587031215613580578182fd5b843561358b81614b1d565b966020860135965060408601359560600135945092505050565b60008060008060e085870312156135ba578182fd5b6135c48686612d7e565b935060408501356001600160401b03808211156135df578384fd5b6135eb88838901613316565b94506060870135915080821115613600578384fd5b5061360d878288016130f2565b92505061361d8660808701613413565b905092959194509250565b6000602080838503121561363a578182fd5b82516001600160401b0381111561364f578283fd5b8301601f8101851361365f578283fd5b805161366d612e1682614ab3565b8181528381019083850185840285018601891015613689578687fd5b8694505b838510156136ab57805183526001949094019391850191850161368d565b50979650505050505050565b600060208083850312156136c9578182fd5b82516001600160401b03808211156136df578384fd5b818501915085601f8301126136f2578384fd5b8151613700612e1682614ab3565b81815284810190848601875b848110156137af57815187016060818d03601f1901121561372b57898afd5b60408051606081018181108a8211171561374157fe5b8252828b015161375081614b1d565b81528282015189811115613762578c8dfd5b6137708f8d838701016132c9565b8c83015250606083015189811115613786578c8dfd5b6137948f8d83870101612ed2565b9282019290925286525050928701929087019060010161370c565b50909998505050505050505050565b6000602082840312156137cf578081fd5b815161180181614b32565b6000602082840312156137eb578081fd5b5035919050565b600080600060608486031215613806578081fd5b8335925060208401356001600160401b03811115613822578182fd5b61382e86828701613010565b925050604084013590509250925092565b600060208284031215613850578081fd5b81356001600160401b0380821115613866578283fd5b818401915061012080838703121561387c578384fd5b61388581614a90565b905082358152602083013560208201526040830135828111156138a6578485fd5b6138b28782860161327b565b604083015250606083013560608201526080830135608082015260a083013560a082015260c0830135828111156138e7578485fd5b6138f38782860161327b565b60c08301525060e083013560e08201526101008084013583811115613916578586fd5b613922888287016131ab565b918301919091525095945050505050565b60008060008060008060c0878903121561394b578384fd5b86356001600160401b0380821115613961578586fd5b61396d8a838b01613316565b975061397b60208a01613514565b96506040890135915080821115613990578586fd5b61399c8a838b01613153565b95506139aa60608a0161352a565b945060808901359150808211156139bf578384fd5b6139cb8a838b01612f9c565b935060a08901359150808211156139e0578283fd5b506139ed89828a01613209565b9150509295509295509295565b6000806000806000806000610120888a031215613a15578485fd5b87356001600160401b0380821115613a2b578687fd5b613a378b838c01613316565b9850613a4560208b01613514565b975060408a0135915080821115613a5a578687fd5b613a668b838c01613153565b9650613a7460608b0161352a565b955060808a0135915080821115613a89578283fd5b613a958b838c01612f9c565b945060a08a0135915080821115613aaa578283fd5b50613ab78a828b01613209565b925050613ac78960c08a01613413565b905092959891949750929550565b600080600080600060c08688031215613aec578283fd5b85359450613afd8760208801612d7e565b935060608601356001600160401b03811115613b17578384fd5b613b23888289016130f2565b935050613b3260808701613514565b915060a0860135613b4281614b1d565b809150509295509295909350565b600080600080600060a08688031215613b67578283fd5b8535945060208601356001600160401b0380821115613b84578485fd5b613b9089838a01612df8565b95506040880135915080821115613ba5578485fd5b613bb189838a01612df8565b9450606088013593506080880135915080821115613bcd578283fd5b50613bda888289016131ab565b9150509295509295909350565b600080600060608486031215613bfb578081fd5b8335925060208401356001600160401b0380821115613c18578283fd5b613c2487838801612df8565b93506040860135915080821115613c39578283fd5b50613c46868287016131ab565b9150509250925092565b600080600080600060a08688031215613c67578283fd5b853594506020860135935060408601356001600160401b0380821115613c8b578485fd5b613bb189838a0161327b565b60008060008060808587031215613cac578182fd5b5050823594602084013594506040840135936060013592509050565b600080600080600080600060e0888a031215613ce2578081fd5b613ceb88613514565b965060208801356001600160401b0380821115613d06578283fd5b613d128b838c01613316565b975060408a0135915080821115613d27578283fd5b613d338b838c0161327b565b965060608a0135915080821115613d48578283fd5b613d548b838c01613010565b9550613d6260808b0161352a565b945060a08a0135915080821115613d77578283fd5b613d838b838c01613209565b935060c08a0135915080821115613d98578283fd5b50613da58a828b01612f9c565b91505092959891949750929550565b6000815180845260208085019450808401835b83811015613dec5781516001600160a01b031687529582019590820190600101613dc7565b509495945050505050565b6000815180845260208085018081965082840281019150828601855b85811015613e6a578284038952815180518552858101518686015260408082015160ff1690860152606090810151608091860182905290613e5681870183613fb9565b9a87019a9550505090840190600101613e13565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613dec578151805160ff16885283810151848901526040908101519088015260609096019590820190600101613e8a565b6000815180845260208085018081965082840281019150828601855b85811015613e6a578284038952815180516001600160a01b0316855285810151606087870181905290613f1082880182613fb9565b91505060408083015192508682038188015250613f2d8183613df7565b9a87019a9550505090840190600101613edb565b6000815180845260208085018081965082840281019150828601855b85811015613e6a578284038952613f75848351613fe5565b98850198935090840190600101613f5d565b6000815180845260208085019450808401835b83811015613dec57815160ff1687529582019590820190600101613f9a565b60008151808452613fd1816020860160208601614af1565b601f01601f19169290920160200192915050565b6000815160808452613ffa6080850182613ebf565b9050602083015184820360208601526140138282613fb9565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b6000608082526140e06080830187613df7565b82810360208401526140f28187613df7565b905082810360408401526141068186613df7565b91505082606083015295945050505050565b60006080825261412b6080830187613df7565b851515602084015282810360408401526141068186613df7565b6000602082526118016020830184613ebf565b901515815260200190565b90815260200190565b600086825260a0602083015261418560a0830187613fb9565b82810360408401526141978187613ebf565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b6020808252601f908201527f496e76616c696420466f7263654d6f7665417070205472616e736974696f6e00604082015260600190565b60208082526015908201527427379037b733b7b4b7339031b430b63632b733b29760591b604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b602080825260129082015271496e76616c6964207369676e61747572657360701b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b602080825260169082015275546f6f206d616e79207061727469636970616e74732160501b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f5369676e6572206e6f7420617574686f72697a6564206d6f7665720000000000604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252818101527f426164207c7369676e6174757265737c767c77686f5369676e6564576861747c604082015260600190565b6020808252818101527f556e61636365707461626c652077686f5369676e656457686174206172726179604082015260600190565b60208082526018908201527f4f7574636f6d65206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601c908201527f737461747573284368616e6e656c4461746129213d73746f7261676500000000604082015260600190565b60208082526018908201527f61707044617461206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601d908201527f496e76616c6964207369676e617475726573202f2021697346696e616c000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b60208082526017908201527f4f7574636f6d65206368616e676520766572626f74656e000000000000000000604082015260600190565b6020808252601e908201527f7c77686f5369676e6564576861747c213d6e5061727469636970616e74730000604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b602080825260129082015271697346696e616c20726574726f677261646560701b604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b60006080825261490b6080830187613fe5565b828103602084015261491d8187613fe5565b65ffffffffffff95909516604084015250506060015292915050565b600086825260a0602083015261495260a0830187613db4565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff808a1683528089166020840152871515604084015260e06060840152865160e0840152602087015160a06101008501526149ec610180850182613db4565b6040890151831661012086015260608901516001600160a01b03166101408601526080808a01518416610160870152858203908601529050614a2e8188613f41565b91505082810360a0840152614a438186613e77565b905082810360c0840152614a578185613f87565b9a9950505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b0381118282101715614aab57fe5b604052919050565b60006001600160401b03821115614ac657fe5b5060209081020190565b60006001600160401b03821115614ae357fe5b50601f01601f191660200190565b60005b83811015614b0c578181015183820152602001614af4565b838111156122dd5750506000910152565b6001600160a01b03811681146121b557600080fd5b80151581146121b557600080fd5b60ff811681146121b557600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a26469706673582212201aaf00142654e5629bc7503e0275b25e1a7d4f3b32ef3ada18857456264495cd64736f6c63430007040033",
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

// ValidTransition is a free data retrieval call binding the contract method 0x6dfa8b74.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCaller) ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "validTransition", nParticipants, isFinalAB, ab, turnNumB, appDefinition)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidTransition is a free data retrieval call binding the contract method 0x6dfa8b74.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// ValidTransition is a free data retrieval call binding the contract method 0x6dfa8b74.
//
// Solidity: function validTransition(uint256 nParticipants, bool[2] isFinalAB, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, uint48 turnNumB, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ValidTransition(nParticipants *big.Int, isFinalAB [2]bool, ab [2]IForceMoveAppVariablePart, turnNumB *big.Int, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, isFinalAB, ab, turnNumB, appDefinition)
}

// Challenge is a paid mutator transaction binding the contract method 0xdd5b9932.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xdd5b9932.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0xdd5b9932.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xa26d793a.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xa26d793a.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart IForceMoveFixedPart, largestTurnNum *big.Int, variableParts []IForceMoveAppVariablePart, isFinalCount uint8, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, largestTurnNum, variableParts, isFinalCount, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0xa26d793a.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, uint48 largestTurnNum, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, uint8 isFinalCount, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
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

// Conclude is a paid mutator transaction binding the contract method 0x9ae6008b.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x9ae6008b.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x9ae6008b.
//
// Solidity: function conclude(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xfb6d5139.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xfb6d5139.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xfb6d5139.
//
// Solidity: function concludeAndTransferAllAssets(uint48 largestTurnNum, (uint256,address[],uint48,address,uint48) fixedPart, bytes appData, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(largestTurnNum *big.Int, fixedPart IForceMoveFixedPart, appData []byte, outcome []ExitFormatSingleAssetExit, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, largestTurnNum, fixedPart, appData, outcome, numStates, whoSignedWhat, sigs)
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

// Respond is a paid mutator transaction binding the contract method 0x5973edd9.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Respond(opts *bind.TransactOpts, isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "respond", isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0x5973edd9.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Respond(isFinalAB [2]bool, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Respond(&_NitroAdjudicator.TransactOpts, isFinalAB, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0x5973edd9.
//
// Solidity: function respond(bool[2] isFinalAB, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
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

// TransferAllAssets is a paid mutator transaction binding the contract method 0xaf69c9d7.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcome []ExitFormatSingleAssetExit, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "transferAllAssets", channelId, outcome, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0xaf69c9d7.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) TransferAllAssets(channelId [32]byte, outcome []ExitFormatSingleAssetExit, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.TransferAllAssets(&_NitroAdjudicator.TransactOpts, channelId, outcome, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0xaf69c9d7.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,bytes,(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) TransferAllAssets(channelId [32]byte, outcome []ExitFormatSingleAssetExit, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.TransferAllAssets(&_NitroAdjudicator.TransactOpts, channelId, outcome, stateHash)
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

// FilterChallengeRegistered is a free log retrieval operation binding the contract event 0xbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f38.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
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

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0xbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f38.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
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

// ParseChallengeRegistered is a log parse operation binding the contract event 0xbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f38.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat)
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

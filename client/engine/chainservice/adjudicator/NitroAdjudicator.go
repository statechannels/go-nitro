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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"indexed\":false,\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"indexed\":false,\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"indexed\":false,\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[]\",\"name\":\"variableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"name\":\"compute_claim_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newSourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newTargetAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart\",\"name\":\"latestVariablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart\",\"name\":\"latestVariablePart\",\"type\":\"tuple\"},{\"internalType\":\"uint8\",\"name\":\"numStates\",\"type\":\"uint8\"},{\"internalType\":\"uint8[]\",\"name\":\"whoSignedWhat\",\"type\":\"uint8[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"numParticipants\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numStates\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numSigs\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"numWhoSignedWhats\",\"type\":\"uint256\"}],\"name\":\"requireValidInput\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structIForceMove.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"variablePartAB\",\"type\":\"tuple[2]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structIForceMove.Signature\",\"name\":\"sig\",\"type\":\"tuple\"}],\"name\":\"respond\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"nParticipants\",\"type\":\"uint256\"},{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structIForceMoveApp.VariablePart[2]\",\"name\":\"ab\",\"type\":\"tuple[2]\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"}],\"name\":\"validTransition\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"pure\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50614a82806100206000396000f3fe6080604052600436106100fe5760003560e01c806366d1b8a211610095578063af69c9d711610064578063af69c9d7146102af578063be5c2a31146102cf578063c36b7e4e146102ef578063c7df14e21461030f578063e29cffe01461032f576100fe565b806366d1b8a21461022257806372c7f16d1461024257806380039982146102625780638bf6ed3914610282576100fe565b80633033730e116100d15780633033730e1461019e57806330776841146101be578063552cfa50146101de578063564b81ef1461020d576100fe565b806311e9f17814610103578063166e56cd1461013c578063180f6ff0146101695780632fb1d2701461018b575b600080fd5b34801561010f57600080fd5b5061012361011e366004613b45565b61035f565b6040516101339493929190613fe2565b60405180910390f35b34801561014857600080fd5b5061015c6101573660046134be565b6106ad565b604051610133919061402d565b34801561017557600080fd5b506101896101843660046139ff565b6106ca565b005b6101896101993660046134e9565b6106f3565b3480156101aa57600080fd5b506101896101b9366004613c06565b610982565b3480156101ca57600080fd5b506101896101d936600461373a565b610a02565b3480156101ea57600080fd5b506101fe6101f93660046136d5565b610ae1565b60405161013393929190614942565b34801561021957600080fd5b5061015c610afc565b34801561022e57600080fd5b5061018961023d366004613946565b610b00565b34801561024e57600080fd5b5061018961025d36600461389f565b610c6e565b34801561026e57600080fd5b5061018961027d3660046139ff565b610ccc565b34801561028e57600080fd5b506102a261029d366004613bae565b610cd9565b6040516101339190614022565b3480156102bb57600080fd5b506101896102ca3660046136ed565b610cee565b3480156102db57600080fd5b506102a26102ea366004613c4d565b611072565b3480156102fb57600080fd5b5061018961030a36600461382e565b6110f6565b34801561031b57600080fd5b5061015c61032a3660046136d5565b611246565b34801561033b57600080fd5b5061034f61034a366004613abb565b611258565b6040516101339493929190613f97565b606060006060600080855111610376578551610379565b84515b6001600160401b038111801561038e57600080fd5b506040519080825280602002602001820160405280156103c857816020015b6103b5612d08565b8152602001906001900390816103ad5790505b5091506000905085516001600160401b03811180156103e657600080fd5b5060405190808252806020026020018201604052801561042057816020015b61040d612d08565b8152602001906001900390816104055790505b50935060019250866000805b88518110156106a15788818151811061044157fe5b60200260200101516000015187828151811061045957fe5b6020026020010151600001818152505088818151811061047557fe5b60200260200101516040015187828151811061048d57fe5b60200260200101516040019060ff16908160ff16815250508881815181106104b157fe5b6020026020010151606001518782815181106104c957fe5b60200260200101516060018190525060006104fb8a83815181106104e957fe5b602002602001015160200151856117cf565b905088516000148061052a575088518310801561052a57508189848151811061052057fe5b6020026020010151145b1561063c57600260ff168a848151811061054057fe5b60200260200101516040015160ff1614156105765760405162461bcd60e51b815260040161056d9061438e565b60405180910390fd5b808a838151811061058357fe5b6020026020010151602001510388838151811061059c57fe5b6020026020010151602001818152505060405180608001604052808b84815181106105c357fe5b60200260200101516000015181526020018281526020018b84815181106105e657fe5b60200260200101516040015160ff1681526020018b848151811061060657fe5b60200260200101516060015181525086848151811061062157fe5b60200260200101819052508085019450826001019250610671565b89828151811061064857fe5b60200260200101516020015188838151811061066057fe5b602002602001015160200181815250505b87828151811061067d57fe5b60200260200101516020015160001461069557600096505b9092039160010161042c565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b60006106d986868686866117e9565b85519091506106eb9082906000610cee565b505050505050565b6106fc836119bc565b156107195760405162461bcd60e51b815260040161056d906145d3565b6001600160a01b03841660009081526001602090815260408083208684529091528120548381101561075d5760405162461bcd60e51b815260040161056d90614102565b61076784846119c8565b81106107855760405162461bcd60e51b815260040161056d906146d8565b6107998161079386866119c8565b90611a22565b91506001600160a01b0386166107cd578234146107c85760405162461bcd60e51b815260040161056d90614746565b61086b565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906107fd90339030908790600401613f39565b602060405180830381600087803b15801561081757600080fd5b505af115801561082b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061084f91906136b9565b61086b5760405162461bcd60e51b815260040161056d90614456565b600061087782846119c8565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715906108d6908a9087908690613f76565b60405180910390a26001600160a01b0387166109795760006108f88585611a22565b90506000336001600160a01b03168260405161091390612b0c565b60006040518083038185875af1925050503d8060008114610950576040519150601f19603f3d011682016040523d82523d6000602084013e610955565b606091505b50509050806109765760405162461bcd60e51b815260040161056d906147a9565b50505b50505050505050565b6000806000610994888589888a611a7f565b92509250925060008060006109c184878d815181106109af57fe5b6020026020010151604001518961035f565b935093505092506109d88b868c8b8a888a88611afd565b6109f5868c815181106109e757fe5b602002602001015183611ba9565b5050505050505050505050565b600080600080610a1185611be1565b93509350935093506060806060600080888a6060015181518110610a3157fe5b60200260200101516040015190506000888b60e0015181518110610a5157fe5b6020026020010151604001519050610a758783838e608001518f6101000151611258565b809650819750829850839950505050505050610ace89878a878c8e6060015181518110610a9e57fe5b6020026020010151604001518e6080015181518110610ab957fe5b6020026020010151600001518c898c89611d8a565b610976878a60e00151815181106109e757fe5b6000806000610aef84611e8d565b9196909550909350915050565b4690565b610b14856020015151855185518551611072565b506000610b2086611eab565b90506000610b2d86611f21565b6040015190506000610b3e83611f48565b6002811115610b4957fe5b1415610b5e57610b598282611f92565b610b8d565b6001610b6983611f48565b6002811115610b7457fe5b1415610b8457610b598282611fd6565b610b8d82612014565b6000610b9c87848a898961204b565b9050610bad818960200151866120b9565b827fbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f38838a608001514201610be08b611f21565b606001518c8c8c8c604051610bfb9796959493929190614882565b60405180910390a2610c5260405180608001604052808465ffffffffffff1681526020018a60800151420165ffffffffffff168152602001838152602001610c4b610c458b611f21565b51612119565b9052612132565b6000938452602084905260409093209290925550505050505050565b610c82846020015151845184518451611072565b506000610c8e85611eab565b90506000610c9b85611f21565b604001519050610caa82612014565b610cb48282611fd6565b610cc1858388878761204b565b506106eb8282612183565b6106eb85858585856117e9565b6000610ce684848461220a565b949350505050565b610cf7836122cc565b610d0a81610d0484612119565b856122ff565b81516001906000906001600160401b0381118015610d2757600080fd5b50604051908082528060200260200182016040528015610d6157816020015b610d4e612d2e565b815260200190600190039081610d465790505b509050600084516001600160401b0381118015610d7d57600080fd5b50604051908082528060200260200182016040528015610da7578160200160208202803683370190505b509050600085516001600160401b0381118015610dc357600080fd5b50604051908082528060200260200182016040528015610ded578160200160208202803683370190505b50905060005b8651811015610f87576000878281518110610e0a57fe5b602002602001015190506000816040015190506000898481518110610e2b57fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c815260200190815260200160002054868581518110610e7c57fe5b602002602001018181525050600080600080610eed8a8981518110610e9d57fe5b60200260200101518760006001600160401b0381118015610ebd57600080fd5b50604051908082528060200260200182016040528015610ee7578160200160208202803683370190505b5061035f565b935093509350935082610eff5760009b505b80898981518110610f0c57fe5b602002602001018181525050838e8981518110610f2557fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b8981518110610f6857fe5b6020026020010181905250505050505050508080600101915050610df3565b5060005b865181101561103b576000878281518110610fa257fe5b6020026020010151600001519050828281518110610fbc57fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d835290935291909120805491909103905583518990600080516020614a2d83398151915290849087908290811061101457fe5b602002602001015160405161102a929190614861565b60405180910390a250600101610f8b565b50831561105657600087815260208190526040812055611069565b611069878661106489612119565b612348565b610979836123ae565b60008385101580156110845750600084115b6110a05760405162461bcd60e51b815260040161056d906142f0565b84831480156110ae57508482145b6110ca5760405162461bcd60e51b815260040161056d9061448d565b60ff8511156110eb5760405162461bcd60e51b815260040161056d90614327565b506001949350505050565b600061110184611eab565b905060008061110f83611e8d565b5086516020810151905192945090925060009161113a918691868a865b6020020151606001516123de565b60208781015190810151905191925060009161115f9187916001888101908c9061112c565b90506111b460405180608001604052808665ffffffffffff1681526020018565ffffffffffff1681526020018481526020016111ac8a6000600281106111a157fe5b602002015151612119565b90528661241a565b6020880151805165ffffffffffff6001870116816111ce57fe5b06815181106111d957fe5b60200260200101516001600160a01b03166111f4828861242d565b6001600160a01b03161461121a5760405162461bcd60e51b815260040161056d906143f0565b61122e886020015151888a6060015161220a565b5061123c8585600101612183565b5050505050505050565b60006020819052908152604090205481565b606080606060008088516001600160401b038111801561127757600080fd5b506040519080825280602002602001820160405280156112b157816020015b61129e612d08565b8152602001906001900390816112965790505b50945087516001600160401b03811180156112cb57600080fd5b5060405190808252806020026020018201604052801561130557816020015b6112f2612d08565b8152602001906001900390816112ea5790505b50935087516001600160401b038111801561131f57600080fd5b5060405190808252806020026020018201604052801561135957816020015b611346612d08565b81526020019060019003908161133e5790505b50925060005b89518110156114455789818151811061137457fe5b60200260200101516000015186828151811061138c57fe5b602002602001015160000181815250508981815181106113a857fe5b6020026020010151602001518682815181106113c057fe5b602002602001015160200181815250508981815181106113dc57fe5b6020026020010151606001518682815181106113f457fe5b60200260200101516060018190525089818151811061140f57fe5b60200260200101516040015186828151811061142757fe5b602090810291909101015160ff90911660409091015260010161135f565b5060005b88518110156115f05788818151811061145e57fe5b60200260200101516000015185828151811061147657fe5b6020026020010151600001818152505088818151811061149257fe5b6020026020010151602001518582815181106114aa57fe5b602002602001015160200181815250508881815181106114c657fe5b6020026020010151606001518582815181106114de57fe5b6020026020010151606001819052508881815181106114f957fe5b60200260200101516040015185828151811061151157fe5b60200260200101516040019060ff16908160ff168152505088818151811061153557fe5b60200260200101516000015184828151811061154d57fe5b60200260200101516000018181525050600084828151811061156b57fe5b6020026020010151602001818152505088818151811061158757fe5b60200260200101516060015184828151811061159f57fe5b6020026020010151606001819052508881815181106115ba57fe5b6020026020010151604001518482815181106115d257fe5b602090810291909101015160ff909116604090910152600101611449565b508960005b88811015611637578161160757611637565b600061162a8c838151811061161857fe5b602002602001015160200151846117cf565b90920391506001016115f5565b50600061165b828c8b8151811061164a57fe5b6020026020010151602001516117cf565b9050600061167f8c8b8151811061166e57fe5b6020026020010151606001516124df565b905060005b81518110156117be5782611697576117be565b60005b88518110156117b557836116ad576117b5565b8881815181106116b957fe5b6020026020010151600001518383815181106116d157fe5b602002602001015114156117ad5760006117028e83815181106116f057fe5b602002602001015160200151866117cf565b905080850394508b516000148061173657508b51871080156117365750818c888151811061172c57fe5b6020026020010151145b156117a757808a838151811061174857fe5b60200260200101516020018181510391508181525050808b8e8151811061176b57fe5b602002602001015160200181815103915081815250508089838151811061178e57fe5b6020908102919091018101510152968701966001909601955b506117b5565b60010161169a565b50600101611684565b505050505095509550955095915050565b60008183116117de57826117e0565b815b90505b92915050565b60006117f486611eab565b90506117ff81612014565b6118158660200151518560ff1684518651611072565b508360ff16856040015160010165ffffffffffff1610156118485760405162461bcd60e51b815260040161056d9061460a565b60008460ff166001600160401b038111801561186357600080fd5b5060405190808252806020026020018201604052801561188d578160200160208202803683370190505b50905060005b8560ff168165ffffffffffff1610156118f1576118ca83886020015189600001518960ff16856001018c60400151010360016123de565b828265ffffffffffff16815181106118de57fe5b6020908102919091010152600101611893565b50611907866040015188602001518386886124f5565b6119235760405162461bcd60e51b815260040161056d9061459c565b6119646040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b8152602001610c4b8960000151612119565b60008084815260200190815260200160002081905550817f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516119aa919061486f565b60405180910390a25095945050505050565b60a081901c155b919050565b6000828201838110156117e0576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611a79576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6060600080611a8d876125cb565b611a96866122cc565b611aa8858580519060200120886122ff565b611ab18461262a565b9250828881518110611abf57fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a908110611b3757fe5b602002602001015160400181905250611b77868686604051602001611b5c919061400f565b60405160208183030381529060405280519060200120612348565b85600080516020614a2d8339815191528984604051611b97929190614861565b60405180910390a25050505050505050565b611bdd604051806060016040528084600001516001600160a01b031681526020018460200151815260200183815250612640565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611c15906125cb565b611c1e856122cc565b611c348a602001518580519060200120876122ff565b611c3d8461262a565b9850611c488261262a565b9750888381518110611c5657fe5b6020908102919091010151519650600260ff16898481518110611c7557fe5b6020026020010151604001518281518110611c8c57fe5b60200260200101516040015160ff1614611cb85760405162461bcd60e51b815260040161056d9061470f565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a9085908110611ced57fe5b6020026020010151604001518b6080015181518110611d0857fe5b6020026020010151600001519050876001600160a01b0316898381518110611d2c57fe5b6020026020010151600001516001600160a01b031614611d5e5760405162461bcd60e51b815260040161056d90614357565b611d67816122cc565b611d7d8b60a001518480519060200120836122ff565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b9084908110611dd057fe5b602002602001015160400181905250611df9838d602001518c604051602001611b5c919061400f565b85878281518110611e0657fe5b602002602001015160400181905250611e2f888d60a0015189604051602001611b5c919061400f565b82600080516020614a2d8339815191528387604051611e4f929190614861565b60405180910390a287600080516020614a2d8339815191528287604051611e77929190614861565b60405180910390a2505050505050505050505050565b60009081526020819052604090205460d081901c9160a082901c9190565b6000611eb5610afc565b825114611ed45760405162461bcd60e51b815260040161056d90614139565b611edc610afc565b8260200151836040015184606001518560800151604051602001611f04959493929190614816565b604051602081830303815290604052805190602001209050919050565b611f29612d58565b81600183510381518110611f3957fe5b60200260200101519050919050565b600080611f5483611e8d565b5091505065ffffffffffff8116611f6f5760009150506119c3565b428165ffffffffffff1611611f885760029150506119c3565b60019150506119c3565b6000611f9d83611e8d565b505090508065ffffffffffff168265ffffffffffff161015611fd15760405162461bcd60e51b815260040161056d90614190565b505050565b6000611fe183611e8d565b505090508065ffffffffffff168265ffffffffffff1611611fd15760405162461bcd60e51b815260040161056d906140cb565b600261201f82611f48565b600281111561202a57fe5b14156120485760405162461bcd60e51b815260040161056d90614164565b50565b6000806120598787876126ec565b905061207861206788611f21565b6040015186602001518387876124f5565b6120945760405162461bcd60e51b815260040161056d906142c4565b806001825103815181106120a457fe5b60200260200101519150505b95945050505050565b60006120eb846040516020016120cf9190614083565b604051602081830303815290604052805190602001208361242d565b90506120f78184612879565b6121135760405162461bcd60e51b815260040161056d9061425d565b50505050565b6000612124826128cf565b805190602001209050919050565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161216f916128f8565b6001600160a01b0316919091179392505050565b6040805160808101825265ffffffffffff8316815260006020820181905291810182905260608101919091526121b890612132565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0826040516121fe919061486f565b60405180910390a25050565b6000806122178585612924565b9050600181600181111561222757fe5b14156110eb5783516020850151604051630a4d2cb960e21b81526001600160a01b03861692632934b2e492612260928a906004016147e0565b60206040518083038186803b15801561227857600080fd5b505afa15801561228c573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906122b091906136b9565b6110eb5760405162461bcd60e51b815260040161056d906141c7565b60026122d782611f48565b60028111156122e257fe5b146120485760405162461bcd60e51b815260040161056d9061463a565b600061230a82611e8d565b9250505061231884846128f8565b6001600160a01b0316816001600160a01b0316146121135760405162461bcd60e51b815260040161056d90614427565b60008061235485611e8d565b5091509150600061239460405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250612132565b600096875260208790526040909620959095555050505050565b60005b8151811015611bdd576123d68282815181106123c957fe5b6020026020010151612640565b6001016123b1565b600085858585856040516020016123f9959493929190614036565b60405160208183030381529060405280519060200120905095945050505050565b6124248282612a19565b611bdd81612a4f565b600080836040516020016124419190613f08565b60405160208183030381529060405280519060200120905060006001828560000151866020015187604001516040516000815260200160405260405161248a94939291906140ad565b6020604051602081039080840390855afa1580156124ac573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610ce65760405162461bcd60e51b815260040161056d906143c5565b6060818060200190518101906117e39190613523565b835183516000919061250984898484612a82565b6125255760405162461bcd60e51b815260040161056d906144c2565b60005b828110156125bc5760006125788887848151811061254257fe5b602002602001015160ff168151811061255757fe5b602002602001015188848151811061256b57fe5b602002602001015161242d565b905088828151811061258657fe5b60200260200101516001600160a01b0316816001600160a01b0316146125b35760009450505050506120b0565b50600101612528565b50600198975050505050505050565b60005b8151816001011015611bdd578181600101815181106125e957fe5b60200260200101518282815181106125fd57fe5b6020026020010151106126225760405162461bcd60e51b815260040161056d9061422d565b6001016125ce565b6060818060200190518101906117e391906135b2565b805160005b826040015151811015611fd15760008360400151828151811061266457fe5b602002602001015160000151905060008460400151838151811061268457fe5b602002602001015160200151905061269b826119bc565b156126b8576126b3846126ad84612b0c565b83612b0f565b6126e2565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b5050600101612645565b6060600084516001600160401b038111801561270757600080fd5b50604051908082528060200260200182016040528015612731578160200160208202803683370190505b50905060005b85518165ffffffffffff161015612870576127d185878365ffffffffffff168151811061276057fe5b602002602001015160200151888465ffffffffffff168151811061278057fe5b602002602001015160000151898565ffffffffffff16815181106127a057fe5b6020026020010151604001518a8665ffffffffffff16815181106127c057fe5b6020026020010151606001516123de565b828265ffffffffffff16815181106127e557fe5b60200260200101818152505060018651038165ffffffffffff161015612868576128668460200151516040518060400160405280898565ffffffffffff168151811061282d57fe5b60200260200101518152602001898560010165ffffffffffff168151811061285157fe5b6020026020010151815250866060015161220a565b505b600101612737565b50949350505050565b6000805b82518110156128c55782818151811061289257fe5b60200260200101516001600160a01b0316846001600160a01b031614156128bd5760019150506117e3565b60010161287d565b5060009392505050565b6060816040516020016128e2919061400f565b6040516020818303038152906040529050919050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b6020810151606001516000901561296b5761294a8260015b602002015151835151612c1f565b6129665760405162461bcd60e51b815260040161056d9061466a565b612a10565b8151606001511561298e5760405162461bcd60e51b815260040161056d9061477d565b6002830282600160200201516040015165ffffffffffff161015612a08576129b782600161293c565b6129d35760405162461bcd60e51b815260040161056d906144f7565b6020828101518101518351909101516129ec9190612c3b565b6129665760405162461bcd60e51b815260040161056d90614565565b5060016117e3565b50600092915050565b600081815260208190526040902054612a33908390612c9f565b611bdd5760405162461bcd60e51b815260040161056d9061452e565b6001612a5a82611f48565b6002811115612a6557fe5b146120485760405162461bcd60e51b815260040161056d906141fe565b600082855114612aa45760405162461bcd60e51b815260040161056d906146a1565b60005b83811015612b0057600084828765ffffffffffff1687010381612ac657fe5b0690508381888481518110612ad757fe5b602002602001015160ff16016001011015612af757600092505050610ce6565b50600101612aa7565b50600195945050505050565b90565b6001600160a01b038316612b9f576000826001600160a01b031682604051612b3690612b0c565b60006040518083038185875af1925050503d8060008114612b73576040519150601f19603f3d011682016040523d82523d6000602084013e612b78565b606091505b5050905080612b995760405162461bcd60e51b815260040161056d90614294565b50611fd1565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb90612bcd9085908590600401613f5d565b602060405180830381600087803b158015612be757600080fd5b505af1158015612bfb573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061211391906136b9565b60006117e0612c2d846128cf565b612c36846128cf565b612cb3565b815181516000916001918114808314612c575760009250612c95565b600160208701838101602088015b600284838510011415612c90578051835114612c845760009650600093505b60209283019201612c65565b505050505b5090949350505050565b600081612cab84612132565b149392505050565b815181516000916001918114808314612ccf5760009250612c95565b600160208701838101602088015b600284838510011415612c90578051835114612cfc5760009650600093505b60209283019201612cdd565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b80356119c3816149fa565b600082601f830112612da5578081fd5b81356020612dba612db583614990565b61496d565b82815281810190858301855b85811015612e5f5781358801608080601f19838d03011215612de6578889fd5b604080518281016001600160401b038282108183111715612e0357fe5b908352848a0135825284830135828b0152606090612e228287016134a8565b83850152938501359380851115612e37578c8dfd5b50612e468e8b86880101613205565b9082015287525050509284019290840190600101612dc6565b5090979650505050505050565b600082601f830112612e7c578081fd5b81516020612e8c612db583614990565b82815281810190858301855b85811015612e5f5781518801608080601f19838d03011215612eb8578889fd5b604080518281016001600160401b038282108183111715612ed557fe5b908352848a0151825284830151828b0152606090612ef48287016134b3565b83850152938501519380851115612f09578c8dfd5b50612f188e8b86880101613251565b9082015287525050509284019290840190600101612e98565b600082601f830112612f41578081fd5b81356020612f51612db583614990565b82815281810190858301606080860288018501891015612f6f578687fd5b865b86811015612f9557612f838a84613391565b85529385019391810191600101612f71565b509198975050505050505050565b600082601f830112612fb3578081fd5b81356020612fc3612db583614990565b82815281810190858301855b85811015612e5f5781358801606080601f19838d03011215612fef578889fd5b604080518281016001600160401b03828210818311171561300c57fe5b908352848a01359061301d826149fa565b908252848301359080821115613031578c8dfd5b61303f8f8c84890101613205565b838c0152938501359380851115613054578c8dfd5b50506130648d8a85870101612d95565b91810191909152865250509284019290840190600101612fcf565b600082601f83011261308f578081fd5b604051604081018181106001600160401b03821117156130ab57fe5b6040528083835b60028110156130dd576130c887833588016133ee565b835260209283019291909101906001016130b2565b509195945050505050565b600082601f8301126130f8578081fd5b81356020613108612db583614990565b82815281810190858301855b85811015612e5f5761312b898684358b01016133ee565b84529284019290840190600101613114565b600082601f83011261314d578081fd5b8135602061315d612db583614990565b8281528181019085830183850287018401881015613179578586fd5b855b85811015612e5f5781358452928401929084019060010161317b565b600082601f8301126131a7578081fd5b813560206131b7612db583614990565b82815281810190858301838502870184018810156131d3578586fd5b855b85811015612e5f5781356131e881614a1d565b845292840192908401906001016131d5565b80356119c381614a0f565b600082601f830112613215578081fd5b8135613223612db5826149ad565b818152846020838601011115613237578283fd5b816020850160208301379081016020019190915292915050565b600082601f830112613261578081fd5b815161326f612db5826149ad565b818152846020838601011115613283578283fd5b610ce68260208301602087016149ce565b600060a082840312156132a5578081fd5b60405160a081016001600160401b0382821081831117156132c257fe5b81604052829350843583526020915081850135818111156132e257600080fd5b85019050601f810186136132f557600080fd5b8035613303612db582614990565b81815283810190838501858402850186018a101561332057600080fd5b600094505b8385101561334c578035613338816149fa565b835260019490940193918501918501613325565b508085870152505050505061336360408401613492565b604082015261337460608401612d8a565b606082015261338560808401613492565b60808201525092915050565b6000606082840312156133a2578081fd5b604051606081018181106001600160401b03821117156133be57fe5b60405290508082356133cf81614a1d565b8082525060208301356020820152604083013560408201525092915050565b6000608082840312156133ff578081fd5b604051608081016001600160401b03828210818311171561341c57fe5b81604052829350843591508082111561343457600080fd5b61344086838701612fa3565b8352602085013591508082111561345657600080fd5b5061346385828601613205565b60208301525061347560408401613492565b6040820152613486606084016131fa565b60608201525092915050565b803565ffffffffffff811681146119c357600080fd5b80356119c381614a1d565b80516119c381614a1d565b600080604083850312156134d0578182fd5b82356134db816149fa565b946020939093013593505050565b600080600080608085870312156134fe578182fd5b8435613509816149fa565b966020860135965060408601359560600135945092505050565b60006020808385031215613535578182fd5b82516001600160401b0381111561354a578283fd5b8301601f8101851361355a578283fd5b8051613568612db582614990565b8181528381019083850185840285018601891015613584578687fd5b8694505b838510156135a6578051835260019490940193918501918501613588565b50979650505050505050565b600060208083850312156135c4578182fd5b82516001600160401b03808211156135da578384fd5b818501915085601f8301126135ed578384fd5b81516135fb612db582614990565b81815284810190848601875b848110156136aa57815187016060818d03601f1901121561362657898afd5b60408051606081018181108a8211171561363c57fe5b8252828b015161364b816149fa565b8152828201518981111561365d578c8dfd5b61366b8f8d83870101613251565b8c83015250606083015189811115613681578c8dfd5b61368f8f8d83870101612e6c565b92820192909252865250509287019290870190600101613607565b50909998505050505050505050565b6000602082840312156136ca578081fd5b81516117e081614a0f565b6000602082840312156136e6578081fd5b5035919050565b600080600060608486031215613701578081fd5b8335925060208401356001600160401b0381111561371d578182fd5b61372986828701612fa3565b925050604084013590509250925092565b60006020828403121561374b578081fd5b81356001600160401b0380821115613761578283fd5b8184019150610120808387031215613777578384fd5b6137808161496d565b905082358152602083013560208201526040830135828111156137a1578485fd5b6137ad87828601613205565b604083015250606083013560608201526080830135608082015260a083013560a082015260c0830135828111156137e2578485fd5b6137ee87828601613205565b60c08301525060e083013560e08201526101008084013583811115613811578586fd5b61381d8882870161313d565b918301919091525095945050505050565b600080600060a08486031215613842578081fd5b83356001600160401b0380821115613858578283fd5b61386487838801613294565b94506020860135915080821115613879578283fd5b506138868682870161307f565b9250506138968560408601613391565b90509250925092565b600080600080608085870312156138b4578182fd5b84356001600160401b03808211156138ca578384fd5b6138d688838901613294565b955060208701359150808211156138eb578384fd5b6138f7888389016130e8565b9450604087013591508082111561390c578384fd5b61391888838901612f31565b9350606087013591508082111561392d578283fd5b5061393a87828801613197565b91505092959194509250565b600080600080600060e0868803121561395d578283fd5b85356001600160401b0380821115613973578485fd5b61397f89838a01613294565b96506020880135915080821115613994578485fd5b6139a089838a016130e8565b955060408801359150808211156139b5578485fd5b6139c189838a01612f31565b945060608801359150808211156139d6578283fd5b506139e388828901613197565b9250506139f38760808801613391565b90509295509295909350565b600080600080600060a08688031215613a16578283fd5b85356001600160401b0380821115613a2c578485fd5b613a3889838a01613294565b96506020880135915080821115613a4d578485fd5b613a5989838a016133ee565b955060408801359150613a6b82614a1d565b90935060608701359080821115613a80578283fd5b613a8c89838a01613197565b93506080880135915080821115613aa1578283fd5b50613aae88828901612f31565b9150509295509295909350565b600080600080600060a08688031215613ad2578283fd5b8535945060208601356001600160401b0380821115613aef578485fd5b613afb89838a01612d95565b95506040880135915080821115613b10578485fd5b613b1c89838a01612d95565b9450606088013593506080880135915080821115613b38578283fd5b50613aae8882890161313d565b600080600060608486031215613b59578081fd5b8335925060208401356001600160401b0380821115613b76578283fd5b613b8287838801612d95565b93506040860135915080821115613b97578283fd5b50613ba48682870161313d565b9150509250925092565b600080600060608486031215613bc2578081fd5b8335925060208401356001600160401b03811115613bde578182fd5b613bea8682870161307f565b9250506040840135613bfb816149fa565b809150509250925092565b600080600080600060a08688031215613c1d578283fd5b853594506020860135935060408601356001600160401b0380821115613c41578485fd5b613b1c89838a01613205565b60008060008060808587031215613c62578182fd5b5050823594602084013594506040840135936060013592509050565b6000815180845260208085019450808401835b83811015613cb65781516001600160a01b031687529582019590820190600101613c91565b509495945050505050565b6000815180845260208085018081965082840281019150828601855b85811015613d34578284038952815180518552858101518686015260408082015160ff1690860152606090810151608091860182905290613d2081870183613e83565b9a87019a9550505090840190600101613cdd565b5091979650505050505050565b6000815180845260208085019450808401835b83811015613cb6578151805160ff16885283810151848901526040908101519088015260609096019590820190600101613d54565b6000815180845260208085018081965082840281019150828601855b85811015613d34578284038952815180516001600160a01b0316855285810151606087870181905290613dda82880182613e83565b91505060408083015192508682038188015250613df78183613cc1565b9a87019a9550505090840190600101613da5565b6000815180845260208085018081965082840281019150828601855b85811015613d34578284038952613e3f848351613eaf565b98850198935090840190600101613e27565b6000815180845260208085019450808401835b83811015613cb657815160ff1687529582019590820190600101613e64565b60008151808452613e9b8160208601602086016149ce565b601f01601f19169290920160200192915050565b6000815160808452613ec46080850182613d89565b905060208301518482036020860152613edd8282613e83565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b600060808252613faa6080830187613cc1565b8281036020840152613fbc8187613cc1565b90508281036040840152613fd08186613cc1565b91505082606083015295945050505050565b600060808252613ff56080830187613cc1565b85151560208401528281036040840152613fd08186613cc1565b6000602082526117e06020830184613d89565b901515815260200190565b90815260200190565b600086825260a0602083015261404f60a0830187613e83565b82810360408401526140618187613d89565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b6020808252601f908201527f496e76616c696420466f7263654d6f7665417070205472616e736974696f6e00604082015260600190565b60208082526015908201527427379037b733b7b4b7339031b430b63632b733b29760591b604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b602080825260129082015271496e76616c6964207369676e61747572657360701b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b602080825260169082015275546f6f206d616e79207061727469636970616e74732160501b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f5369676e6572206e6f7420617574686f72697a6564206d6f7665720000000000604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252818101527f426164207c7369676e6174757265737c767c77686f5369676e6564576861747c604082015260600190565b6020808252818101527f556e61636365707461626c652077686f5369676e656457686174206172726179604082015260600190565b60208082526018908201527f4f7574636f6d65206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601c908201527f737461747573284368616e6e656c4461746129213d73746f7261676500000000604082015260600190565b60208082526018908201527f61707044617461206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601d908201527f496e76616c6964207369676e617475726573202f2021697346696e616c000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b60208082526017908201527f4f7574636f6d65206368616e676520766572626f74656e000000000000000000604082015260600190565b6020808252601e908201527f7c77686f5369676e6564576861747c213d6e5061727469636970616e74730000604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b602080825260129082015271697346696e616c20726574726f677261646560701b604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b6000606082526147f36060830186613eaf565b82810360208401526148058186613eaf565b915050826040830152949350505050565b600086825260a0602083015261482f60a0830187613c7e565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff808a1683528089166020840152871515604084015260e06060840152865160e0840152602087015160a06101008501526148c9610180850182613c7e565b6040890151831661012086015260608901516001600160a01b03166101408601526080808a0151841661016087015285820390860152905061490b8188613e0b565b91505082810360a08401526149208186613d41565b905082810360c08401526149348185613e51565b9a9950505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b038111828210171561498857fe5b604052919050565b60006001600160401b038211156149a357fe5b5060209081020190565b60006001600160401b038211156149c057fe5b50601f01601f191660200190565b60005b838110156149e95781810151838201526020016149d1565b838111156121135750506000910152565b6001600160a01b038116811461204857600080fd5b801515811461204857600080fd5b60ff8116811461204857600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a2646970667358221220ccb75aeb41e937bd931d86e2c3179ff812e54cb5ea4c868cb2234536b6eb5ca464736f6c63430007060033",
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

// ValidTransition is a free data retrieval call binding the contract method 0x8bf6ed39.
//
// Solidity: function validTransition(uint256 nParticipants, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCaller) ValidTransition(opts *bind.CallOpts, nParticipants *big.Int, ab [2]IForceMoveAppVariablePart, appDefinition common.Address) (bool, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "validTransition", nParticipants, ab, appDefinition)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// ValidTransition is a free data retrieval call binding the contract method 0x8bf6ed39.
//
// Solidity: function validTransition(uint256 nParticipants, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorSession) ValidTransition(nParticipants *big.Int, ab [2]IForceMoveAppVariablePart, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, ab, appDefinition)
}

// ValidTransition is a free data retrieval call binding the contract method 0x8bf6ed39.
//
// Solidity: function validTransition(uint256 nParticipants, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] ab, address appDefinition) pure returns(bool)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ValidTransition(nParticipants *big.Int, ab [2]IForceMoveAppVariablePart, appDefinition common.Address) (bool, error) {
	return _NitroAdjudicator.Contract.ValidTransition(&_NitroAdjudicator.CallOpts, nParticipants, ab, appDefinition)
}

// Challenge is a paid mutator transaction binding the contract method 0x66d1b8a2.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, variableParts, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x66d1b8a2.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, variableParts, sigs, whoSignedWhat, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x66d1b8a2.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8, challengerSig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, variableParts, sigs, whoSignedWhat, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x72c7f16d.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, variableParts, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x72c7f16d.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, variableParts, sigs, whoSignedWhat)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x72c7f16d.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[] variableParts, (uint8,bytes32,bytes32)[] sigs, uint8[] whoSignedWhat) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Checkpoint(fixedPart IForceMoveFixedPart, variableParts []IForceMoveAppVariablePart, sigs []IForceMoveSignature, whoSignedWhat []uint8) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, variableParts, sigs, whoSignedWhat)
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

// Conclude is a paid mutator transaction binding the contract method 0x80039982.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x80039982.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
}

// Conclude is a paid mutator transaction binding the contract method 0x80039982.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x180f6ff0.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x180f6ff0.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x180f6ff0.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool) latestVariablePart, uint8 numStates, uint8[] whoSignedWhat, (uint8,bytes32,bytes32)[] sigs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(fixedPart IForceMoveFixedPart, latestVariablePart IForceMoveAppVariablePart, numStates uint8, whoSignedWhat []uint8, sigs []IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, latestVariablePart, numStates, whoSignedWhat, sigs)
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

// Respond is a paid mutator transaction binding the contract method 0xc36b7e4e.
//
// Solidity: function respond((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Respond(opts *bind.TransactOpts, fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "respond", fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xc36b7e4e.
//
// Solidity: function respond((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Respond(fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Respond(&_NitroAdjudicator.TransactOpts, fixedPart, variablePartAB, sig)
}

// Respond is a paid mutator transaction binding the contract method 0xc36b7e4e.
//
// Solidity: function respond((uint256,address[],uint48,address,uint48) fixedPart, ((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool)[2] variablePartAB, (uint8,bytes32,bytes32) sig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Respond(fixedPart IForceMoveFixedPart, variablePartAB [2]IForceMoveAppVariablePart, sig IForceMoveSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Respond(&_NitroAdjudicator.TransactOpts, fixedPart, variablePartAB, sig)
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

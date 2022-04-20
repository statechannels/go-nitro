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
	Bin: "0x608060405234801561001057600080fd5b506151d1806100206000396000f3fe6080604052600436106100fe5760003560e01c806366d1b8a211610095578063af69c9d711610064578063af69c9d7146102af578063be5c2a31146102cf578063c36b7e4e146102ef578063c7df14e21461030f578063e29cffe01461032f576100fe565b806366d1b8a21461022257806372c7f16d1461024257806380039982146102625780638bf6ed3914610282576100fe565b80633033730e116100d15780633033730e1461019e57806330776841146101be578063552cfa50146101de578063564b81ef1461020d576100fe565b806311e9f17814610103578063166e56cd1461013c578063180f6ff0146101695780632fb1d2701461018b575b600080fd5b34801561010f57600080fd5b5061012361011e3660046140ac565b61035f565b6040516101339493929190614599565b60405180910390f35b34801561014857600080fd5b5061015c6101573660046138a2565b6106ad565b60405161013391906145e4565b34801561017557600080fd5b50610189610184366004613ef0565b6106ca565b005b6101896101993660046138cd565b6106f3565b3480156101aa57600080fd5b506101896101b936600461417d565b610982565b3480156101ca57600080fd5b506101896101d9366004613b3a565b610a02565b3480156101ea57600080fd5b506101fe6101f9366004613ad5565b610ae1565b60405161013393929190614f92565b34801561021957600080fd5b5061015c610afc565b34801561022e57600080fd5b5061018961023d366004613d5c565b610b00565b34801561024e57600080fd5b5061018961025d366004613ca3565b610d72565b34801561026e57600080fd5b5061018961027d366004613e37565b610e78565b34801561028e57600080fd5b506102a261029d366004614115565b610f1e565b60405161013391906145d9565b3480156102bb57600080fd5b506101896102ca366004613aed565b610f33565b3480156102db57600080fd5b506102a26102ea3660046141c4565b6112b7565b3480156102fb57600080fd5b5061018961030a366004613c2e565b61133b565b34801561031b57600080fd5b5061015c61032a366004613ad5565b6115dc565b34801561033b57600080fd5b5061034f61034a366004614022565b6115ee565b604051610133949392919061454e565b606060006060600080855111610376578551610379565b84515b6001600160401b038111801561038e57600080fd5b506040519080825280602002602001820160405280156103c857816020015b6103b561309e565b8152602001906001900390816103ad5790505b5091506000905085516001600160401b03811180156103e657600080fd5b5060405190808252806020026020018201604052801561042057816020015b61040d61309e565b8152602001906001900390816104055790505b50935060019250866000805b88518110156106a15788818151811061044157fe5b60200260200101516000015187828151811061045957fe5b6020026020010151600001818152505088818151811061047557fe5b60200260200101516040015187828151811061048d57fe5b60200260200101516040019060ff16908160ff16815250508881815181106104b157fe5b6020026020010151606001518782815181106104c957fe5b60200260200101516060018190525060006104fb8a83815181106104e957fe5b60200260200101516020015185611b65565b905088516000148061052a575088518310801561052a57508189848151811061052057fe5b6020026020010151145b1561063c57600260ff168a848151811061054057fe5b60200260200101516040015160ff1614156105765760405162461bcd60e51b815260040161056d90614945565b60405180910390fd5b808a838151811061058357fe5b6020026020010151602001510388838151811061059c57fe5b6020026020010151602001818152505060405180608001604052808b84815181106105c357fe5b60200260200101516000015181526020018281526020018b84815181106105e657fe5b60200260200101516040015160ff1681526020018b848151811061060657fe5b60200260200101516060015181525086848151811061062157fe5b60200260200101819052508085019450826001019250610671565b89828151811061064857fe5b60200260200101516020015188838151811061066057fe5b602002602001015160200181815250505b87828151811061067d57fe5b60200260200101516020015160001461069557600096505b9092039160010161042c565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b60006106d98686868686611b7f565b85519091506106eb9082906000610f33565b505050505050565b6106fc83611d52565b156107195760405162461bcd60e51b815260040161056d90614b8a565b6001600160a01b03841660009081526001602090815260408083208684529091528120548381101561075d5760405162461bcd60e51b815260040161056d906146b9565b6107678484611d5e565b81106107855760405162461bcd60e51b815260040161056d90614c8f565b610799816107938686611d5e565b90611db8565b91506001600160a01b0386166107cd578234146107c85760405162461bcd60e51b815260040161056d90614cfd565b61086b565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906107fd903390309087906004016144f0565b602060405180830381600087803b15801561081757600080fd5b505af115801561082b573d6000803e3d6000fd5b505050506040513d601f19601f8201168201806040525081019061084f9190613ab9565b61086b5760405162461bcd60e51b815260040161056d90614a0d565b60006108778284611d5e565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a715906108d6908a908790869061452d565b60405180910390a26001600160a01b0387166109795760006108f88585611db8565b90506000336001600160a01b03168260405161091390612ea2565b60006040518083038185875af1925050503d8060008114610950576040519150601f19603f3d011682016040523d82523d6000602084013e610955565b606091505b50509050806109765760405162461bcd60e51b815260040161056d90614d60565b50505b50505050505050565b6000806000610994888589888a611e15565b92509250925060008060006109c184878d815181106109af57fe5b6020026020010151604001518961035f565b935093505092506109d88b868c8b8a888a88611e93565b6109f5868c815181106109e757fe5b602002602001015183611f3f565b5050505050505050505050565b600080600080610a1185611f77565b93509350935093506060806060600080888a6060015181518110610a3157fe5b60200260200101516040015190506000888b60e0015181518110610a5157fe5b6020026020010151604001519050610a758783838e608001518f61010001516115ee565b809650819750829850839950505050505050610ace89878a878c8e6060015181518110610a9e57fe5b6020026020010151604001518e6080015181518110610ab957fe5b6020026020010151600001518c898c89612120565b610976878a60e00151815181106109e757fe5b6000806000610aef84612223565b9196909550909350915050565b4690565b610b1c610b106020890189614fbd565b885190915086856112b7565b506000610b30610b2b89615105565b612241565b90506000610b3d886122b7565b6040015190506000610b4e836122de565b6002811115610b5957fe5b1415610b6e57610b698282612328565b610b9d565b6001610b79836122de565b6002811115610b8457fe5b1415610b9457610b69828261236c565b610b9d826123aa565b6000610c398984610bad8d615105565b8b8b808060200260200160405190810160405280939291908181526020016000905b82821015610bfb57610bec60608302860136819003810190614007565b81526020019060010190610bcf565b50505050508a8a808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506123e192505050565b9050610c9281610c4c60208d018d614fbd565b80806020026020016040519081016040528093929190818152602001838360200280828437600092019190915250610c8d9250505036889003880188614007565b61244f565b827fbc93eb35d1e2104744f8837e55265ad2f79d1c3f0843e40e757505a1c9007f3883610cc560a08e0160808f016141f5565b4201610cd08d6122b7565b606001518e8e8e8e8e8e604051610cef99989796959493929190614e6d565b60405180910390a2610d5460405180608001604052808465ffffffffffff1681526020018c6080016020810190610d2691906141f5565b420165ffffffffffff168152602001838152602001610d4d610d478d6122b7565b516124af565b90526124c8565b60009384526020849052604090932092909255505050505050505050565b610d8c610d826020890189614fbd565b90508685846112b7565b506000610d9b610b2b89615105565b90506000610db1610dac888a6150f8565b6122b7565b604001519050610dc0826123aa565b610dca828261236c565b610e6d610dd7888a6150f8565b83610de18c615105565b8989808060200260200160405190810160405280939291908181526020016000905b82821015610e2f57610e2060608302860136819003810190614007565b81526020019060010190610e03565b50505050508888808060200260200160405190810160405280939291908181526020018383602002808284376000920191909152506123e192505050565b506109768282612519565b610f14610e8488615105565b610e8d88615111565b878787808060200260200160405190810160405280939291908181526020018383602002808284376000920182905250604080516020808d02820181019092528b815294508b93508a925082919085015b82821015610f0a57610efb60608302860136819003810190614007565b81526020019060010190610ede565b5050505050611b7f565b5050505050505050565b6000610f2b8484846125a0565b949350505050565b610f3c83612662565b610f4f81610f49846124af565b85612695565b81516001906000906001600160401b0381118015610f6c57600080fd5b50604051908082528060200260200182016040528015610fa657816020015b610f936130c4565b815260200190600190039081610f8b5790505b509050600084516001600160401b0381118015610fc257600080fd5b50604051908082528060200260200182016040528015610fec578160200160208202803683370190505b509050600085516001600160401b038111801561100857600080fd5b50604051908082528060200260200182016040528015611032578160200160208202803683370190505b50905060005b86518110156111cc57600087828151811061104f57fe5b60200260200101519050600081604001519050600089848151811061107057fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c8152602001908152602001600020548685815181106110c157fe5b6020026020010181815250506000806000806111328a89815181106110e257fe5b60200260200101518760006001600160401b038111801561110257600080fd5b5060405190808252806020026020018201604052801561112c578160200160208202803683370190505b5061035f565b9350935093509350826111445760009b505b8089898151811061115157fe5b602002602001018181525050838e898151811061116a57fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b89815181106111ad57fe5b6020026020010181905250505050505050508080600101915050611038565b5060005b86518110156112805760008782815181106111e757fe5b602002602001015160000151905082828151811061120157fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d83529093529190912080549190910390558351899060008051602061517c83398151915290849087908290811061125957fe5b602002602001015160405161126f929190614e4c565b60405180910390a2506001016111d0565b50831561129b576000878152602081905260408120556112ae565b6112ae87866112a9896124af565b6126de565b61097983612744565b60008385101580156112c95750600084115b6112e55760405162461bcd60e51b815260040161056d906148a7565b84831480156112f357508482145b61130f5760405162461bcd60e51b815260040161056d90614a44565b60ff8511156113305760405162461bcd60e51b815260040161056d906148de565b506001949350505050565b6000611349610b2b85615105565b905060008061135783612223565b50909250905060006114058461136d888061505d565b61137b906020810190615019565b8080601f01602080910402602001604051908101604052809392919081815260200183838082843760009201829052508c935091506113b79050565b6020028101906113c7919061505d565b6113d19080615003565b6113da916150dd565b868a60005b6020028101906113ef919061505d565b611400906080810190606001613a9d565b612774565b905060006114978561141a60208a018a61505d565b611428906020810190615019565b8080601f0160208091040260200160405190810160405280939291908181526020018383808284376000920191909152508c9250600191506114679050565b602002810190611477919061505d565b6114819080615003565b61148a916150dd565b6001808901908c906113df565b905061150960405180608001604052808665ffffffffffff1681526020018565ffffffffffff1681526020018481526020016115018a6000600281106114d957fe5b6020028101906114e9919061505d565b6114f39080615003565b6114fc916150dd565b6124af565b9052866127b0565b6115166020890189614fbd565b61152360208b018b614fbd565b90508660010165ffffffffffff168161153857fe5b0681811061154257fe5b90506020020160208101906115579190613886565b6001600160a01b031661157882611573368a90038a018a614007565b6127c3565b6001600160a01b03161461159e5760405162461bcd60e51b815260040161056d906149a7565b6115ce6115ae60208a018a614fbd565b90506115b9896150ea565b6115c960808c0160608d01613886565b6125a0565b50610f148585600101612519565b60006020819052908152604090205481565b606080606060008088516001600160401b038111801561160d57600080fd5b5060405190808252806020026020018201604052801561164757816020015b61163461309e565b81526020019060019003908161162c5790505b50945087516001600160401b038111801561166157600080fd5b5060405190808252806020026020018201604052801561169b57816020015b61168861309e565b8152602001906001900390816116805790505b50935087516001600160401b03811180156116b557600080fd5b506040519080825280602002602001820160405280156116ef57816020015b6116dc61309e565b8152602001906001900390816116d45790505b50925060005b89518110156117db5789818151811061170a57fe5b60200260200101516000015186828151811061172257fe5b6020026020010151600001818152505089818151811061173e57fe5b60200260200101516020015186828151811061175657fe5b6020026020010151602001818152505089818151811061177257fe5b60200260200101516060015186828151811061178a57fe5b6020026020010151606001819052508981815181106117a557fe5b6020026020010151604001518682815181106117bd57fe5b602090810291909101015160ff9091166040909101526001016116f5565b5060005b8851811015611986578881815181106117f457fe5b60200260200101516000015185828151811061180c57fe5b6020026020010151600001818152505088818151811061182857fe5b60200260200101516020015185828151811061184057fe5b6020026020010151602001818152505088818151811061185c57fe5b60200260200101516060015185828151811061187457fe5b60200260200101516060018190525088818151811061188f57fe5b6020026020010151604001518582815181106118a757fe5b60200260200101516040019060ff16908160ff16815250508881815181106118cb57fe5b6020026020010151600001518482815181106118e357fe5b60200260200101516000018181525050600084828151811061190157fe5b6020026020010151602001818152505088818151811061191d57fe5b60200260200101516060015184828151811061193557fe5b60200260200101516060018190525088818151811061195057fe5b60200260200101516040015184828151811061196857fe5b602090810291909101015160ff9091166040909101526001016117df565b508960005b888110156119cd578161199d576119cd565b60006119c08c83815181106119ae57fe5b60200260200101516020015184611b65565b909203915060010161198b565b5060006119f1828c8b815181106119e057fe5b602002602001015160200151611b65565b90506000611a158c8b81518110611a0457fe5b602002602001015160600151612875565b905060005b8151811015611b545782611a2d57611b54565b60005b8851811015611b4b5783611a4357611b4b565b888181518110611a4f57fe5b602002602001015160000151838381518110611a6757fe5b60200260200101511415611b43576000611a988e8381518110611a8657fe5b60200260200101516020015186611b65565b905080850394508b5160001480611acc57508b5187108015611acc5750818c8881518110611ac257fe5b6020026020010151145b15611b3d57808a8381518110611ade57fe5b60200260200101516020018181510391508181525050808b8e81518110611b0157fe5b6020026020010151602001818151039150818152505080898381518110611b2457fe5b6020908102919091018101510152968701966001909601955b50611b4b565b600101611a30565b50600101611a1a565b505050505095509550955095915050565b6000818311611b745782611b76565b815b90505b92915050565b6000611b8a86612241565b9050611b95816123aa565b611bab8660200151518560ff16845186516112b7565b508360ff16856040015160010165ffffffffffff161015611bde5760405162461bcd60e51b815260040161056d90614bc1565b60008460ff166001600160401b0381118015611bf957600080fd5b50604051908082528060200260200182016040528015611c23578160200160208202803683370190505b50905060005b8560ff168165ffffffffffff161015611c8757611c6083886020015189600001518960ff16856001018c6040015101036001612774565b828265ffffffffffff1681518110611c7457fe5b6020908102919091010152600101611c29565b50611c9d8660400151886020015183868861288b565b611cb95760405162461bcd60e51b815260040161056d90614b53565b611cfa6040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b8152602001610d4d89600001516124af565b60008084815260200190815260200160002081905550817f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc90142604051611d409190614e5a565b60405180910390a25095945050505050565b60a081901c155b919050565b600082820183811015611b76576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611e0f576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6060600080611e2387612961565b611e2c86612662565b611e3e85858051906020012088612695565b611e47846129c0565b9250828881518110611e5557fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a908110611ecd57fe5b602002602001015160400181905250611f0d868686604051602001611ef291906145c6565b604051602081830303815290604052805190602001206126de565b8560008051602061517c8339815191528984604051611f2d929190614e4c565b60405180910390a25050505050505050565b611f73604051806060016040528084600001516001600160a01b0316815260200184602001518152602001838152506129d6565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611fab90612961565b611fb485612662565b611fca8a60200151858051906020012087612695565b611fd3846129c0565b9850611fde826129c0565b9750888381518110611fec57fe5b6020908102919091010151519650600260ff1689848151811061200b57fe5b602002602001015160400151828151811061202257fe5b60200260200101516040015160ff161461204e5760405162461bcd60e51b815260040161056d90614cc6565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a908590811061208357fe5b6020026020010151604001518b608001518151811061209e57fe5b6020026020010151600001519050876001600160a01b03168983815181106120c257fe5b6020026020010151600001516001600160a01b0316146120f45760405162461bcd60e51b815260040161056d9061490e565b6120fd81612662565b6121138b60a00151848051906020012083612695565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b908490811061216657fe5b60200260200101516040018190525061218f838d602001518c604051602001611ef291906145c6565b8587828151811061219c57fe5b6020026020010151604001819052506121c5888d60a0015189604051602001611ef291906145c6565b8260008051602061517c83398151915283876040516121e5929190614e4c565b60405180910390a28760008051602061517c833981519152828760405161220d929190614e4c565b60405180910390a2505050505050505050505050565b60009081526020819052604090205460d081901c9160a082901c9190565b600061224b610afc565b82511461226a5760405162461bcd60e51b815260040161056d906146f0565b612272610afc565b826020015183604001518460600151856080015160405160200161229a959493929190614dcd565b604051602081830303815290604052805190602001209050919050565b6122bf6130ee565b816001835103815181106122cf57fe5b60200260200101519050919050565b6000806122ea83612223565b5091505065ffffffffffff8116612305576000915050611d59565b428165ffffffffffff161161231e576002915050611d59565b6001915050611d59565b600061233383612223565b505090508065ffffffffffff168265ffffffffffff1610156123675760405162461bcd60e51b815260040161056d90614747565b505050565b600061237783612223565b505090508065ffffffffffff168265ffffffffffff16116123675760405162461bcd60e51b815260040161056d90614682565b60026123b5826122de565b60028111156123c057fe5b14156123de5760405162461bcd60e51b815260040161056d9061471b565b50565b6000806123ef878787612a82565b905061240e6123fd886122b7565b60400151866020015183878761288b565b61242a5760405162461bcd60e51b815260040161056d9061487b565b8060018251038151811061243a57fe5b60200260200101519150505b95945050505050565b600061248184604051602001612465919061463a565b60405160208183030381529060405280519060200120836127c3565b905061248d8184612c0f565b6124a95760405162461bcd60e51b815260040161056d90614814565b50505050565b60006124ba82612c65565b805190602001209050919050565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161250591612c8e565b6001600160a01b0316919091179392505050565b6040805160808101825265ffffffffffff83168152600060208201819052918101829052606081019190915261254e906124c8565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0826040516125949190614e5a565b60405180910390a25050565b6000806125ad8585612cba565b905060018160018111156125bd57fe5b14156113305783516020850151604051630a4d2cb960e21b81526001600160a01b03861692632934b2e4926125f6928a90600401614d97565b60206040518083038186803b15801561260e57600080fd5b505afa158015612622573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906126469190613ab9565b6113305760405162461bcd60e51b815260040161056d9061477e565b600261266d826122de565b600281111561267857fe5b146123de5760405162461bcd60e51b815260040161056d90614bf1565b60006126a082612223565b925050506126ae8484612c8e565b6001600160a01b0316816001600160a01b0316146124a95760405162461bcd60e51b815260040161056d906149de565b6000806126ea85612223565b5091509150600061272a60405180608001604052808565ffffffffffff1681526020018465ffffffffffff168152602001878152602001868152506124c8565b600096875260208790526040909620959095555050505050565b60005b8151811015611f735761276c82828151811061275f57fe5b60200260200101516129d6565b600101612747565b6000858585858560405160200161278f9594939291906145ed565b60405160208183030381529060405280519060200120905095945050505050565b6127ba8282612daf565b611f7381612de5565b600080836040516020016127d791906144bf565b6040516020818303038152906040528051906020012090506000600182856000015186602001518760400151604051600081526020016040526040516128209493929190614664565b6020604051602081039080840390855afa158015612842573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116610f2b5760405162461bcd60e51b815260040161056d9061497c565b606081806020019051810190611b799190613907565b835183516000919061289f84898484612e18565b6128bb5760405162461bcd60e51b815260040161056d90614a79565b60005b8281101561295257600061290e888784815181106128d857fe5b602002602001015160ff16815181106128ed57fe5b602002602001015188848151811061290157fe5b60200260200101516127c3565b905088828151811061291c57fe5b60200260200101516001600160a01b0316816001600160a01b031614612949576000945050505050612446565b506001016128be565b50600198975050505050505050565b60005b8151816001011015611f735781816001018151811061297f57fe5b602002602001015182828151811061299357fe5b6020026020010151106129b85760405162461bcd60e51b815260040161056d906147e4565b600101612964565b606081806020019051810190611b799190613996565b805160005b826040015151811015612367576000836040015182815181106129fa57fe5b6020026020010151600001519050600084604001518381518110612a1a57fe5b6020026020010151602001519050612a3182611d52565b15612a4e57612a4984612a4384612ea2565b83612ea5565b612a78565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b50506001016129db565b6060600084516001600160401b0381118015612a9d57600080fd5b50604051908082528060200260200182016040528015612ac7578160200160208202803683370190505b50905060005b85518165ffffffffffff161015612c0657612b6785878365ffffffffffff1681518110612af657fe5b602002602001015160200151888465ffffffffffff1681518110612b1657fe5b602002602001015160000151898565ffffffffffff1681518110612b3657fe5b6020026020010151604001518a8665ffffffffffff1681518110612b5657fe5b602002602001015160600151612774565b828265ffffffffffff1681518110612b7b57fe5b60200260200101818152505060018651038165ffffffffffff161015612bfe57612bfc8460200151516040518060400160405280898565ffffffffffff1681518110612bc357fe5b60200260200101518152602001898560010165ffffffffffff1681518110612be757fe5b602002602001015181525086606001516125a0565b505b600101612acd565b50949350505050565b6000805b8251811015612c5b57828181518110612c2857fe5b60200260200101516001600160a01b0316846001600160a01b03161415612c53576001915050611b79565b600101612c13565b5060009392505050565b606081604051602001612c7891906145c6565b6040516020818303038152906040529050919050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60208101516060015160009015612d0157612ce08260015b602002015151835151612fb5565b612cfc5760405162461bcd60e51b815260040161056d90614c21565b612da6565b81516060015115612d245760405162461bcd60e51b815260040161056d90614d34565b6002830282600160200201516040015165ffffffffffff161015612d9e57612d4d826001612cd2565b612d695760405162461bcd60e51b815260040161056d90614aae565b602082810151810151835190910151612d829190612fd1565b612cfc5760405162461bcd60e51b815260040161056d90614b1c565b506001611b79565b50600092915050565b600081815260208190526040902054612dc9908390613035565b611f735760405162461bcd60e51b815260040161056d90614ae5565b6001612df0826122de565b6002811115612dfb57fe5b146123de5760405162461bcd60e51b815260040161056d906147b5565b600082855114612e3a5760405162461bcd60e51b815260040161056d90614c58565b60005b83811015612e9657600084828765ffffffffffff1687010381612e5c57fe5b0690508381888481518110612e6d57fe5b602002602001015160ff16016001011015612e8d57600092505050610f2b565b50600101612e3d565b50600195945050505050565b90565b6001600160a01b038316612f35576000826001600160a01b031682604051612ecc90612ea2565b60006040518083038185875af1925050503d8060008114612f09576040519150601f19603f3d011682016040523d82523d6000602084013e612f0e565b606091505b5050905080612f2f5760405162461bcd60e51b815260040161056d9061484b565b50612367565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb90612f639085908590600401614514565b602060405180830381600087803b158015612f7d57600080fd5b505af1158015612f91573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906124a99190613ab9565b6000611b76612fc384612c65565b612fcc84612c65565b613049565b815181516000916001918114808314612fed576000925061302b565b600160208701838101602088015b60028483851001141561302657805183511461301a5760009650600093505b60209283019201612ffb565b505050505b5090949350505050565b600081613041846124c8565b149392505050565b815181516000916001918114808314613065576000925061302b565b600160208701838101602088015b6002848385100114156130265780518351146130925760009650600093505b60209283019201613073565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b600061313361312e8461509f565b61507c565b83815290506020808201908360005b868110156131ec5781358601606080828b03121561315f57600080fd5b60408051918201916001600160401b03808411828510171561317d57fe5b92825283359261318c84615149565b9281528387013592808411156131a157600080fd5b6131ad8d8587016135ab565b88830152828501359350808411156131c457600080fd5b506131d18c848601613294565b91810191909152865250509282019290820190600101613142565b505050509392505050565b60006001600160401b0383111561320a57fe5b602061321781850261507c565b9150818360005b868110156131ec5761323388833588016137b6565b8352918301919083019060010161321e565b600061325361312e8461509f565b83815290506020808201908360005b868110156131ec5761327788833588016137b6565b84529282019290820190600101613262565b8035611d5981615149565b600082601f8301126132a4578081fd5b813560206132b461312e8361509f565b82815281810190858301855b858110156133595781358801608080601f19838d030112156132e0578889fd5b604080518281016001600160401b0382821081831117156132fd57fe5b908352848a0135825284830135828b015260609061331c828701613870565b83850152938501359380851115613331578c8dfd5b506133408e8b868801016135ab565b90820152875250505092840192908401906001016132c0565b5090979650505050505050565b600082601f830112613376578081fd5b8151602061338661312e8361509f565b82815281810190858301855b858110156133595781518801608080601f19838d030112156133b2578889fd5b604080518281016001600160401b0382821081831117156133cf57fe5b908352848a0151825284830151828b01526060906133ee82870161387b565b83850152938501519380851115613403578c8dfd5b506134128e8b868801016135f7565b9082015287525050509284019290840190600101613392565b60008083601f84011261343c578182fd5b5081356001600160401b03811115613452578182fd5b60208301915083602060608302850101111561346d57600080fd5b9250929050565b600082601f830112613484578081fd5b8135602061349461312e8361509f565b828152818101908583016060808602880185018910156134b2578687fd5b865b868110156134d8576134c68a84613759565b855293850193918101916001016134b4565b509198975050505050505050565b600082601f8301126134f6578081fd5b611b7683833560208501613120565b60008083601f840112613516578182fd5b5081356001600160401b0381111561352c578182fd5b602083019150836020808302850101111561346d57600080fd5b600082601f830112613556578081fd5b8135602061356661312e8361509f565b8281528181019085830183850287018401881015613582578586fd5b855b8581101561335957813584529284019290840190600101613584565b8035611d598161515e565b600082601f8301126135bb578081fd5b81356135c961312e826150bc565b8181528460208386010111156135dd578283fd5b816020850160208301379081016020019190915292915050565b600082601f830112613607578081fd5b815161361561312e826150bc565b818152846020838601011115613629578283fd5b610f2b82602083016020870161511d565b600060a08284031215611b74578081fd5b600060a0828403121561365c578081fd5b60405160a081016001600160401b03828210818311171561367957fe5b816040528293508435835260209150818501358181111561369957600080fd5b85019050601f810186136136ac57600080fd5b80356136ba61312e8261509f565b81815283810190838501858402850186018a10156136d757600080fd5b600094505b838510156137035780356136ef81615149565b8352600194909401939185019185016136dc565b508085870152505050505061371a6040840161385a565b604082015261372b60608401613289565b606082015261373c6080840161385a565b60808201525092915050565b600060608284031215611b74578081fd5b60006060828403121561376a578081fd5b604051606081018181106001600160401b038211171561378657fe5b60405290508082356137978161516c565b8082525060208301356020820152604083013560408201525092915050565b6000608082840312156137c7578081fd5b604051608081016001600160401b0382821081831117156137e457fe5b8160405282935084359150808211156137fc57600080fd5b613808868387016134e6565b8352602085013591508082111561381e57600080fd5b5061382b858286016135ab565b60208301525061383d6040840161385a565b604082015261384e606084016135a0565b60608201525092915050565b803565ffffffffffff81168114611d5957600080fd5b8035611d598161516c565b8051611d598161516c565b600060208284031215613897578081fd5b8135611b7681615149565b600080604083850312156138b4578081fd5b82356138bf81615149565b946020939093013593505050565b600080600080608085870312156138e2578182fd5b84356138ed81615149565b966020860135965060408601359560600135945092505050565b60006020808385031215613919578182fd5b82516001600160401b0381111561392e578283fd5b8301601f8101851361393e578283fd5b805161394c61312e8261509f565b8181528381019083850185840285018601891015613968578687fd5b8694505b8385101561398a57805183526001949094019391850191850161396c565b50979650505050505050565b600060208083850312156139a8578182fd5b82516001600160401b03808211156139be578384fd5b818501915085601f8301126139d1578384fd5b81516139df61312e8261509f565b81815284810190848601875b84811015613a8e57815187016060818d03601f19011215613a0a57898afd5b60408051606081018181108a82111715613a2057fe5b8252828b0151613a2f81615149565b81528282015189811115613a41578c8dfd5b613a4f8f8d838701016135f7565b8c83015250606083015189811115613a65578c8dfd5b613a738f8d83870101613366565b928201929092528652505092870192908701906001016139eb565b50909998505050505050505050565b600060208284031215613aae578081fd5b8135611b768161515e565b600060208284031215613aca578081fd5b8151611b768161515e565b600060208284031215613ae6578081fd5b5035919050565b600080600060608486031215613b01578081fd5b8335925060208401356001600160401b03811115613b1d578182fd5b613b29868287016134e6565b925050604084013590509250925092565b600060208284031215613b4b578081fd5b81356001600160401b0380821115613b61578283fd5b8184019150610120808387031215613b77578384fd5b613b808161507c565b90508235815260208301356020820152604083013582811115613ba1578485fd5b613bad878286016135ab565b604083015250606083013560608201526080830135608082015260a083013560a082015260c083013582811115613be2578485fd5b613bee878286016135ab565b60c08301525060e083013560e08201526101008084013583811115613c11578586fd5b613c1d88828701613546565b918301919091525095945050505050565b600080600060a08486031215613c42578081fd5b83356001600160401b0380821115613c58578283fd5b613c648783880161363a565b94506020860135915080821115613c79578283fd5b50840160408101861015613c8b578182fd5b9150613c9a8560408601613748565b90509250925092565b60008060008060008060006080888a031215613cbd578485fd5b87356001600160401b0380821115613cd3578687fd5b613cdf8b838c0161363a565b985060208a0135915080821115613cf4578687fd5b613d008b838c01613505565b909850965060408a0135915080821115613d18578485fd5b613d248b838c0161342b565b909650945060608a0135915080821115613d3c578384fd5b50613d498a828b01613505565b989b979a50959850939692959293505050565b600080600080600080600060e0888a031215613d76578081fd5b87356001600160401b0380821115613d8c578283fd5b613d988b838c0161363a565b985060208a0135915080821115613dad578283fd5b818a0191508a601f830112613dc0578283fd5b613dcf8b833560208501613245565b975060408a0135915080821115613de4578283fd5b613df08b838c0161342b565b909750955060608a0135915080821115613e08578283fd5b50613e158a828b01613505565b9094509250613e2990508960808a01613748565b905092959891949750929550565b600080600080600080600060a0888a031215613e51578081fd5b87356001600160401b0380821115613e67578283fd5b613e738b838c0161363a565b985060208a0135915080821115613e88578283fd5b908901906080828c031215613e9b578283fd5b819750613eaa60408b01613870565b965060608a0135915080821115613ebf578283fd5b613ecb8b838c01613505565b909650945060808a0135915080821115613ee3578283fd5b50613d498a828b0161342b565b600080600080600060a08688031215613f07578283fd5b85356001600160401b0380821115613f1d578485fd5b613f2989838a0161364b565b9650602091508188013581811115613f3f578586fd5b613f4b8a828b016137b6565b9650506040880135613f5c8161516c565b9450606088013581811115613f6f578384fd5b8801601f81018a13613f7f578384fd5b8035613f8d61312e8261509f565b81815284810190838601868402850187018e1015613fa9578788fd5b8794505b83851015613fd4578035613fc08161516c565b835260019490940193918601918601613fad565b5096505050506080880135915080821115613fed578283fd5b50613ffa88828901613474565b9150509295509295909350565b600060608284031215614018578081fd5b611b768383613759565b600080600080600060a08688031215614039578283fd5b8535945060208601356001600160401b0380821115614056578485fd5b61406289838a01613294565b95506040880135915080821115614077578485fd5b61408389838a01613294565b945060608801359350608088013591508082111561409f578283fd5b50613ffa88828901613546565b6000806000606084860312156140c0578081fd5b8335925060208401356001600160401b03808211156140dd578283fd5b6140e987838801613294565b935060408601359150808211156140fe578283fd5b5061410b86828701613546565b9150509250925092565b600080600060608486031215614129578081fd5b8335925060208401356001600160401b03811115614145578182fd5b8401601f81018613614155578182fd5b614161866002836131f7565b925050604084013561417281615149565b809150509250925092565b600080600080600060a08688031215614194578283fd5b853594506020860135935060408601356001600160401b03808211156141b8578485fd5b61408389838a016135ab565b600080600080608085870312156141d9578182fd5b5050823594602084013594506040840135936060013592509050565b600060208284031215614206578081fd5b611b768261385a565b6001600160a01b03169052565b60008284526020808501945082825b8581101561425957813561423e81615149565b6001600160a01b03168752958201959082019060010161422b565b509495945050505050565b6000815180845260208085018081965082840281019150828601855b858110156142d7578284038952815180518552858101518686015260408082015160ff16908601526060908101516080918601829052906142c38187018361442e565b9a87019a9550505090840190600101614280565b5091979650505050505050565b60008284526020808501945082825b858110156142595781356143068161516c565b60ff168752818301358388015260408083013590880152606096870196909101906001016142f3565b6000815180845260208085018081965082840281019150828601855b858110156142d7578284038952815180516001600160a01b03168552858101516060878701819052906143808288018261442e565b9150506040808301519250868203818801525061439d8183614264565b9a87019a955050509084019060010161434b565b6000815180845260208085018081965082840281019150828601855b858110156142d75782840389526143e584835161445a565b988501989350908401906001016143cd565b60008284526020808501945082825b858110156142595781356144198161516c565b60ff1687529582019590820190600101614406565b6000815180845261444681602086016020860161511d565b601f01601f19169290920160200192915050565b600081516080845261446f608085018261432f565b905060208301518482036020860152614488828261442e565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b65ffffffffffff169052565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b6000608082526145616080830187614264565b82810360208401526145738187614264565b905082810360408401526145878186614264565b91505082606083015295945050505050565b6000608082526145ac6080830187614264565b851515602084015282810360408401526145878186614264565b600060208252611b76602083018461432f565b901515815260200190565b90815260200190565b600086825260a0602083015261460660a083018761442e565b8281036040840152614618818761432f565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b6020808252601f908201527f496e76616c696420466f7263654d6f7665417070205472616e736974696f6e00604082015260600190565b60208082526015908201527427379037b733b7b4b7339031b430b63632b733b29760591b604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b602080825260129082015271496e76616c6964207369676e61747572657360701b604082015260600190565b6020808252601d908201527f496e73756666696369656e74206f722065786365737320737461746573000000604082015260600190565b602080825260169082015275546f6f206d616e79207061727469636970616e74732160501b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601b908201527f5369676e6572206e6f7420617574686f72697a6564206d6f7665720000000000604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252818101527f426164207c7369676e6174757265737c767c77686f5369676e6564576861747c604082015260600190565b6020808252818101527f556e61636365707461626c652077686f5369676e656457686174206172726179604082015260600190565b60208082526018908201527f4f7574636f6d65206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601c908201527f737461747573284368616e6e656c4461746129213d73746f7261676500000000604082015260600190565b60208082526018908201527f61707044617461206368616e676520666f7262696464656e0000000000000000604082015260600190565b6020808252601d908201527f496e76616c6964207369676e617475726573202f2021697346696e616c000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b6020808252601690820152756c6172676573745475726e4e756d20746f6f206c6f7760501b604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b60208082526017908201527f4f7574636f6d65206368616e676520766572626f74656e000000000000000000604082015260600190565b6020808252601e908201527f7c77686f5369676e6564576861747c213d6e5061727469636970616e74730000604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b602080825260129082015271697346696e616c20726574726f677261646560701b604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b600060608252614daa606083018661445a565b8281036020840152614dbc818661445a565b915050826040830152949350505050565b600060a08201878352602060a08185015281885180845260c086019150828a019350845b81811015614e165784516001600160a01b031683529383019391830191600101614df1565b505065ffffffffffff97881660408601526001600160a01b03969096166060850152505050921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff808c168352808b16602084015250881515604083015260e06060830152873560e08301526020880135601e19893603018112614eb0578182fd5b880180356001600160401b03811115614ec7578283fd5b6020810236038a1315614ed8578283fd5b60a0610100850152614ef26101808501826020850161421c565b915050614f0160408a0161385a565b614f0f6101208501826144b3565b50614f1c60608a01613289565b614f2a61014085018261420f565b50614f3760808a0161385a565b614f456101608501826144b3565b508281036080840152614f5881896143b1565b905082810360a0840152614f6d8187896142e4565b905082810360c0840152614f828185876143f7565b9c9b505050505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6000808335601e19843603018112614fd3578283fd5b8301803591506001600160401b03821115614fec578283fd5b602090810192508102360382131561346d57600080fd5b6000808335601e19843603018112614fd3578182fd5b6000808335601e1984360301811261502f578182fd5b8301803591506001600160401b03821115615048578283fd5b60200191503681900382131561346d57600080fd5b60008235607e19833603018112615072578182fd5b9190910192915050565b6040518181016001600160401b038111828210171561509757fe5b604052919050565b60006001600160401b038211156150b257fe5b5060209081020190565b60006001600160401b038211156150cf57fe5b50601f01601f191660200190565b6000611b76368484613120565b6000611b79366002846131f7565b6000611b76368484613245565b6000611b79368361364b565b6000611b7936836137b6565b60005b83811015615138578181015183820152602001615120565b838111156124a95750506000910152565b6001600160a01b03811681146123de57600080fd5b80151581146123de57600080fd5b60ff811681146123de57600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a264697066735822122024a4c14a74381044a18cda67ae963dde693ac36e2514a6aa936136898ac8ea2464736f6c63430007060033",
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

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
}

// INitroTypesVariablePart is an auto generated low-level Go binding around an user-defined struct.
type INitroTypesVariablePart struct {
	Outcome []ExitFormatSingleAssetExit
	AppData []byte
	TurnNum *big.Int
	IsFinal bool
}

// NitroAdjudicatorMetaData contains all meta data concerning the NitroAdjudicator contract.
var NitroAdjudicatorMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"claim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"uint256[]\",\"name\":\"targetAllocationIndicesToPayout\",\"type\":\"uint256[]\"}],\"name\":\"compute_claim_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newSourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newTargetAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"latestSupportedState\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b5061420c806100206000396000f3fe6080604052600436106100e85760003560e01c8063307768411161008a578063af69c9d711610059578063af69c9d714610279578063c7df14e214610299578063c927c2d3146102b9578063e29cffe0146102d9576100e8565b806330776841146101f55780633503f8f014610215578063552cfa5014610235578063564b81ef14610264576100e8565b806312c52592116100c657806312c5259214610175578063166e56cd146101955780632fb1d270146101c25780633033730e146101d5576100e8565b80630d36d531146100ed5780630d96cc221461010f57806311e9f17814610145575b600080fd5b3480156100f957600080fd5b5061010d6101083660046132d1565b610309565b005b34801561011b57600080fd5b5061012f61012a366004613238565b610331565b60405161013c9190613f30565b60405180910390f35b34801561015157600080fd5b506101656101603660046134f6565b6103f2565b60405161013c949392919061395b565b34801561018157600080fd5b5061010d610190366004613331565b610740565b3480156101a157600080fd5b506101b56101b0366004612f9d565b610891565b60405161013c919061399b565b61010d6101d0366004612fc8565b6108ae565b3480156101e157600080fd5b5061010d6101f036600461355f565b610b3d565b34801561020157600080fd5b5061010d610210366004613144565b610bbd565b34801561022157600080fd5b5061010d6102303660046132d1565b610c9c565b34801561024157600080fd5b506102556102503660046130df565b610ce9565b60405161013c9392919061409f565b34801561027057600080fd5b506101b5610d04565b34801561028557600080fd5b5061010d6102943660046130f7565b610d13565b3480156102a557600080fd5b506101b56102b43660046130df565b611097565b3480156102c557600080fd5b5061010d6102d43660046132d1565b6110a9565b3480156102e557600080fd5b506102f96102f436600461345f565b6110b3565b60405161013c9493929190613910565b6000610315838361162a565b905061032c81610324846116e0565b516000610d13565b505050565b610339612733565b6000610345838561412b565b90506103576080860160608701612f81565b6001600160a01b031663f409586c8661037861037282614138565b8561170c565b6040518363ffffffff1660e01b8152600401610395929190613e45565b60006040518083038186803b1580156103ad57600080fd5b505afa1580156103c1573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103e991908101906133a2565b95945050505050565b60606000606060008085511161040957855161040c565b84515b6001600160401b038111801561042157600080fd5b5060405190808252806020026020018201604052801561045b57816020015b610448612765565b8152602001906001900390816104405790505b5091506000905085516001600160401b038111801561047957600080fd5b506040519080825280602002602001820160405280156104b357816020015b6104a0612765565b8152602001906001900390816104985790505b50935060019250866000805b8851811015610734578881815181106104d457fe5b6020026020010151600001518782815181106104ec57fe5b6020026020010151600001818152505088818151811061050857fe5b60200260200101516040015187828151811061052057fe5b60200260200101516040019060ff16908160ff168152505088818151811061054457fe5b60200260200101516060015187828151811061055c57fe5b602002602001015160600181905250600061058e8a838151811061057c57fe5b602002602001015160200151856117b6565b90508851600014806105bd57508851831080156105bd5750818984815181106105b357fe5b6020026020010151145b156106cf57600260ff168a84815181106105d357fe5b60200260200101516040015160ff1614156106095760405162461bcd60e51b815260040161060090613c03565b60405180910390fd5b808a838151811061061657fe5b6020026020010151602001510388838151811061062f57fe5b6020026020010151602001818152505060405180608001604052808b848151811061065657fe5b60200260200101516000015181526020018281526020018b848151811061067957fe5b60200260200101516040015160ff1681526020018b848151811061069957fe5b6020026020010151606001518152508684815181106106b457fe5b60200260200101819052508085019450826001019250610704565b8982815181106106db57fe5b6020026020010151602001518883815181106106f357fe5b602002602001015160200181815250505b87828151811061071057fe5b60200260200101516020015160001461072857600096505b909203916001016104bf565b50505093509350935093565b600061074b846117ce565b90506000610758846116e0565b604001519050600061076983611844565b600281111561077457fe5b141561078957610784828261188e565b6107b8565b600161079483611844565b600281111561079f57fe5b14156107af5761078482826118cd565b6107b88261190b565b60006107c5868685611942565b90506107d6818760200151866119fd565b827f5bc2767c561afcb4ffbef365ef0757d1d4389267288cbedb1014507e356ff1ec8388608001514201610809896116e0565b606001518a8a604051610820959493929190613faf565b60405180910390a261087760405180608001604052808465ffffffffffff1681526020018860800151420165ffffffffffff16815260200183815260200161087061086a896116e0565b51611a57565b9052611a70565b600093845260208490526040909320929092555050505050565b600160209081526000928352604080842090915290825290205481565b6108b783611ac1565b156108d45760405162461bcd60e51b815260040161060090613d02565b6001600160a01b0384166000908152600160209081526040808320868452909152812054838110156109185760405162461bcd60e51b815260040161060090613a70565b6109228484611ac8565b81106109405760405162461bcd60e51b815260040161060090613d69565b6109548161094e8686611ac8565b90611b22565b91506001600160a01b038616610988578234146109835760405162461bcd60e51b815260040161060090613dd7565b610a26565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906109b8903390309087906004016138b2565b602060405180830381600087803b1580156109d257600080fd5b505af11580156109e6573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a0a91906130c3565b610a265760405162461bcd60e51b815260040161060090613ccb565b6000610a328284611ac8565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a71590610a91908a90879086906138ef565b60405180910390a26001600160a01b038716610b34576000610ab38585611b22565b90506000336001600160a01b031682604051610ace90612620565b60006040518083038185875af1925050503d8060008114610b0b576040519150601f19603f3d011682016040523d82523d6000602084013e610b10565b606091505b5050905080610b315760405162461bcd60e51b815260040161060090613e0e565b50505b50505050505050565b6000806000610b4f888589888a611b7f565b9250925092506000806000610b7c84878d81518110610b6a57fe5b602002602001015160400151896103f2565b93509350509250610b938b868c8b8a888a88611bfd565b610bb0868c81518110610ba257fe5b602002602001015183611ca9565b5050505050505050505050565b600080600080610bcc85611ce1565b93509350935093506060806060600080888a6060015181518110610bec57fe5b60200260200101516040015190506000888b60e0015181518110610c0c57fe5b6020026020010151604001519050610c308783838e608001518f61010001516110b3565b809650819750829850839950505050505050610c8989878a878c8e6060015181518110610c5957fe5b6020026020010151604001518e6080015181518110610c7457fe5b6020026020010151600001518c898c89611e8a565b610b31878a60e0015181518110610ba257fe5b6000610ca7836117ce565b90506000610cb4836116e0565b604001519050610cc38261190b565b610ccd82826118cd565b610cd8848484611942565b50610ce38282611f8d565b50505050565b6000806000610cf784612014565b9196909550909350915050565b6000610d0e612032565b905090565b610d1c83612036565b610d2f81610d2984611a57565b85612069565b81516001906000906001600160401b0381118015610d4c57600080fd5b50604051908082528060200260200182016040528015610d8657816020015b610d7361278b565b815260200190600190039081610d6b5790505b509050600084516001600160401b0381118015610da257600080fd5b50604051908082528060200260200182016040528015610dcc578160200160208202803683370190505b509050600085516001600160401b0381118015610de857600080fd5b50604051908082528060200260200182016040528015610e12578160200160208202803683370190505b50905060005b8651811015610fac576000878281518110610e2f57fe5b602002602001015190506000816040015190506000898481518110610e5057fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c815260200190815260200160002054868581518110610ea157fe5b602002602001018181525050600080600080610f128a8981518110610ec257fe5b60200260200101518760006001600160401b0381118015610ee257600080fd5b50604051908082528060200260200182016040528015610f0c578160200160208202803683370190505b506103f2565b935093509350935082610f245760009b505b80898981518110610f3157fe5b602002602001018181525050838e8981518110610f4a57fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b8981518110610f8d57fe5b6020026020010181905250505050505050508080600101915050610e18565b5060005b8651811015611060576000878281518110610fc757fe5b6020026020010151600001519050828281518110610fe157fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d8352909352919091208054919091039055835189906000805160206141b783398151915290849087908290811061103957fe5b602002602001015160405161104f929190613f8e565b60405180910390a250600101610fb0565b50831561107b5760008781526020819052604081205561108e565b61108e878661108989611a57565b6120b2565b610b3483612118565b60006020819052908152604090205481565b61032c828261162a565b606080606060008088516001600160401b03811180156110d257600080fd5b5060405190808252806020026020018201604052801561110c57816020015b6110f9612765565b8152602001906001900390816110f15790505b50945087516001600160401b038111801561112657600080fd5b5060405190808252806020026020018201604052801561116057816020015b61114d612765565b8152602001906001900390816111455790505b50935087516001600160401b038111801561117a57600080fd5b506040519080825280602002602001820160405280156111b457816020015b6111a1612765565b8152602001906001900390816111995790505b50925060005b89518110156112a0578981815181106111cf57fe5b6020026020010151600001518682815181106111e757fe5b6020026020010151600001818152505089818151811061120357fe5b60200260200101516020015186828151811061121b57fe5b6020026020010151602001818152505089818151811061123757fe5b60200260200101516060015186828151811061124f57fe5b60200260200101516060018190525089818151811061126a57fe5b60200260200101516040015186828151811061128257fe5b602090810291909101015160ff9091166040909101526001016111ba565b5060005b885181101561144b578881815181106112b957fe5b6020026020010151600001518582815181106112d157fe5b602002602001015160000181815250508881815181106112ed57fe5b60200260200101516020015185828151811061130557fe5b6020026020010151602001818152505088818151811061132157fe5b60200260200101516060015185828151811061133957fe5b60200260200101516060018190525088818151811061135457fe5b60200260200101516040015185828151811061136c57fe5b60200260200101516040019060ff16908160ff168152505088818151811061139057fe5b6020026020010151600001518482815181106113a857fe5b6020026020010151600001818152505060008482815181106113c657fe5b602002602001015160200181815250508881815181106113e257fe5b6020026020010151606001518482815181106113fa57fe5b60200260200101516060018190525088818151811061141557fe5b60200260200101516040015184828151811061142d57fe5b602090810291909101015160ff9091166040909101526001016112a4565b508960005b88811015611492578161146257611492565b60006114858c838151811061147357fe5b602002602001015160200151846117b6565b9092039150600101611450565b5060006114b6828c8b815181106114a557fe5b6020026020010151602001516117b6565b905060006114da8c8b815181106114c957fe5b602002602001015160600151612148565b905060005b815181101561161957826114f257611619565b60005b8851811015611610578361150857611610565b88818151811061151457fe5b60200260200101516000015183838151811061152c57fe5b6020026020010151141561160857600061155d8e838151811061154b57fe5b602002602001015160200151866117b6565b905080850394508b516000148061159157508b51871080156115915750818c888151811061158757fe5b6020026020010151145b1561160257808a83815181106115a357fe5b60200260200101516020018181510391508181525050808b8e815181106115c657fe5b60200260200101516020018181510391508181525050808983815181106115e957fe5b6020908102919091018101510152968701966001909601955b50611610565b6001016114f5565b506001016114df565b505050505095509550955095915050565b6000611635836117ce565b90506116408161190b565b61164b838383611942565b5061168c6040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b815260200161087061086a866116e0565b60008083815260200190815260200160002081905550807f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516116d29190613f9c565b60405180910390a292915050565b6116e8612733565b816001835103815181106116f857fe5b60200260200101516000015190505b919050565b6060600082516001600160401b038111801561172757600080fd5b5060405190808252806020026020018201604052801561176157816020015b61174e6127b5565b8152602001906001900390816117465790505b50905060005b83518110156117ac5761178d8585838151811061178057fe5b602002602001015161215e565b82828151811061179957fe5b6020908102919091010152600101611767565b5090505b92915050565b60008183116117c557826117c7565b815b9392505050565b60006117d8612032565b8251146117f75760405162461bcd60e51b815260040161060090613aa7565b6117ff612032565b8260200151836040015184606001518560800151604051602001611827959493929190613f43565b604051602081830303815290604052805190602001209050919050565b60008061185083612014565b5091505065ffffffffffff811661186b576000915050611707565b428165ffffffffffff1611611884576002915050611707565b6001915050611707565b600061189983612014565b505090508065ffffffffffff168265ffffffffffff16101561032c5760405162461bcd60e51b815260040161060090613afe565b60006118d883612014565b505090508065ffffffffffff168265ffffffffffff161161032c5760405162461bcd60e51b815260040161060090613a39565b600261191682611844565b600281111561192157fe5b141561193f5760405162461bcd60e51b815260040161060090613ad2565b50565b60008084606001516001600160a01b031663f409586c86611963888861170c565b6040518363ffffffff1660e01b8152600401611980929190613f0b565b60006040518083038186803b15801561199857600080fd5b505afa1580156119ac573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526119d491908101906133a2565b90506119e08185612206565b6103e9838260200151836000015184604001518560600151612285565b6000611a2f84604051602001611a1391906139f1565b60405160208183030381529060405280519060200120836122c1565b9050611a3b818461237b565b610ce35760405162461bcd60e51b815260040161060090613b65565b6000611a62826123d1565b805190602001209050919050565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b16179291611aad916123fa565b6001600160a01b0316919091179392505050565b60a01c1590565b6000828201838110156117c7576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611b79576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b6060600080611b8d87612426565b611b9686612036565b611ba885858051906020012088612069565b611bb184612485565b9250828881518110611bbf57fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a908110611c3757fe5b602002602001015160400181905250611c77868686604051602001611c5c9190613988565b604051602081830303815290604052805190602001206120b2565b856000805160206141b78339815191528984604051611c97929190613f8e565b60405180910390a25050505050505050565b611cdd604051806060016040528084600001516001600160a01b03168152602001846020015181526020018381525061249b565b5050565b8051604082015160608381015160c085015160e08601516101008701519395869560009586959294919390611d1590612426565b611d1e85612036565b611d348a60200151858051906020012087612069565b611d3d84612485565b9850611d4882612485565b9750888381518110611d5657fe5b6020908102919091010151519650600260ff16898481518110611d7557fe5b6020026020010151604001518281518110611d8c57fe5b60200260200101516040015160ff1614611db85760405162461bcd60e51b815260040161060090613da0565b6001600160a01b03871660009081526001602090815260408083208884529091528120548a519097508a9085908110611ded57fe5b6020026020010151604001518b6080015181518110611e0857fe5b6020026020010151600001519050876001600160a01b0316898381518110611e2c57fe5b6020026020010151600001516001600160a01b031614611e5e5760405162461bcd60e51b815260040161060090613bcc565b611e6781612036565b611e7d8b60a00151848051906020012083612069565b5050505050509193509193565b885160608a015160e08b01516001600160a01b038b166000908152600160209081526040808320868452909152902080548590039055895189908b9084908110611ed057fe5b602002602001015160400181905250611ef9838d602001518c604051602001611c5c9190613988565b85878281518110611f0657fe5b602002602001015160400181905250611f2f888d60a0015189604051602001611c5c9190613988565b826000805160206141b78339815191528387604051611f4f929190613f8e565b60405180910390a2876000805160206141b78339815191528287604051611f77929190613f8e565b60405180910390a2505050505050505050505050565b6040805160808101825265ffffffffffff831681526000602082018190529181018290526060810191909152611fc290611a70565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e0826040516120089190613f9c565b60405180910390a25050565b60009081526020819052604090205460d081901c9160a082901c9190565b4690565b600261204182611844565b600281111561204c57fe5b1461193f5760405162461bcd60e51b815260040161060090613d39565b600061207482612014565b9250505061208284846123fa565b6001600160a01b0316816001600160a01b031614610ce35760405162461bcd60e51b815260040161060090613c9c565b6000806120be85612014565b509150915060006120fe60405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250611a70565b600096875260208790526040909620959095555050505050565b60005b8151811015611cdd5761214082828151811061213357fe5b602002602001015161249b565b60010161211b565b6060818060200190518101906117b09190613002565b6121666127b5565b60408051808201909152825181526000602082018190525b8460200151518110156117ac5760005b8460200151518110156121fd576121e16121ac878760000151612547565b866020015183815181106121bc57fe5b6020026020010151886020015185815181106121d457fe5b6020026020010151612593565b156121f557602083018051600283900a0190525b60010161218e565b5060010161217e565b6122698160018351038151811061221957fe5b6020026020010151600001516040516020016122359190613f30565b604051602081830303815290604052836040516020016122559190613f30565b6040516020818303038152906040526125bc565b611cdd5760405162461bcd60e51b815260040161060090613c3a565b600085858585856040516020016122a09594939291906139a4565b60405160208183030381529060405280519060200120905095945050505050565b600080836040516020016122d59190613881565b60405160208183030381529060405280519060200120905060006001828560000151866020015187604001516040516000815260200160405260405161231e9493929190613a1b565b6020604051602081039080840390855afa158015612340573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b0381166123735760405162461bcd60e51b815260040161060090613c71565b949350505050565b6000805b82518110156123c75782818151811061239457fe5b60200260200101516001600160a01b0316846001600160a01b031614156123bf5760019150506117b0565b60010161237f565b5060009392505050565b6060816040516020016123e49190613988565b6040516020818303038152906040529050919050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60005b8151816001011015611cdd5781816001018151811061244457fe5b602002602001015182828151811061245857fe5b60200260200101511061247d5760405162461bcd60e51b815260040161060090613b35565b600101612429565b6060818060200190518101906117b09190613091565b805160005b82604001515181101561032c576000836040015182815181106124bf57fe5b60200260200101516000015190506000846040015183815181106124df57fe5b60200260200101516020015190506124f682611ac1565b156125135761250e8461250884612620565b83612623565b61253d565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b50506001016124a0565b6000612552836117ce565b6020808401518451604080870151606088015191516125759695919291016139a4565b60405160208183030381529060405280519060200120905092915050565b600061259f84846122c1565b6001600160a01b0316826001600160a01b03161490509392505050565b8151815160009160019181148083146125d85760009250612616565b600160208701838101602088015b6002848385100114156126115780518351146126055760009650600093505b602092830192016125e6565b505050505b5090949350505050565b90565b6001600160a01b0383166126b3576000826001600160a01b03168260405161264a90612620565b60006040518083038185875af1925050503d8060008114612687576040519150601f19603f3d011682016040523d82523d6000602084013e61268c565b606091505b50509050806126ad5760405162461bcd60e51b815260040161060090613b9c565b5061032c565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb906126e190859085906004016138d6565b602060405180830381600087803b1580156126fb57600080fd5b505af115801561270f573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610ce391906130c3565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b60405180604001604052806127c8612733565b8152602001600081525090565b60006127e86127e3846140ed565b6140ca565b83815290506020808201908360005b868110156129065781358601604080828b03121561281457600080fd5b80518181016001600160401b03828210818311171561282f57fe5b81845284358181111561284157600080fd5b85016080818f03121561285357600080fd5b60c08401838110838211171561286557fe5b855280358281111561287657600080fd5b6128828f828401612b44565b845250888101358281111561289657600080fd5b6128a28f828401612d6c565b6060860152506128b3818601612f55565b60808501526128c460608201612d56565b60a085015250818352878501359350808411156128e057600080fd5b50506128ee8b838501612ab3565b818701528652505092820192908201906001016127f7565b505050509392505050565b803561170781614170565b600082601f83011261292c578081fd5b8135602061293c6127e3836140ed565b82815281810190858301855b858110156129e15781358801608080601f19838d03011215612968578889fd5b604080518281016001600160401b03828210818311171561298557fe5b908352848a0135825284830135828b01526060906129a4828701612f6b565b838501529385013593808511156129b9578c8dfd5b506129c88e8b86880101612d6c565b9082015287525050509284019290840190600101612948565b5090979650505050505050565b600082601f8301126129fe578081fd5b81516020612a0e6127e3836140ed565b82815281810190858301855b858110156129e15781518801608080601f19838d03011215612a3a578889fd5b604080518281016001600160401b038282108183111715612a5757fe5b908352848a0151825284830151828b0152606090612a76828701612f76565b83850152938501519380851115612a8b578c8dfd5b50612a9a8e8b86880101612db8565b9082015287525050509284019290840190600101612a1a565b600082601f830112612ac3578081fd5b81356020612ad36127e3836140ed565b82815281810190858301606080860288018501891015612af1578687fd5b865b86811015612b1757612b058a84612ef8565b85529385019391810191600101612af3565b509198975050505050505050565b600082601f830112612b35578081fd5b6117c7838335602085016127d5565b600082601f830112612b54578081fd5b81356020612b646127e3836140ed565b82815281810190858301855b858110156129e15781358801606080601f19838d03011215612b90578889fd5b604080518281016001600160401b038282108183111715612bad57fe5b908352848a013590612bbe82614170565b908252848301359080821115612bd2578c8dfd5b612be08f8c84890101612d6c565b838c0152938501359380851115612bf5578c8dfd5b5050612c058d8a8587010161291c565b91810191909152865250509284019290840190600101612b70565b600082601f830112612c30578081fd5b81516020612c406127e3836140ed565b82815281810190858301855b858110156129e15781518801606080601f19838d03011215612c6c578889fd5b604080518281016001600160401b038282108183111715612c8957fe5b908352848a015190612c9a82614170565b908252848301519080821115612cae578c8dfd5b612cbc8f8c84890101612db8565b838c0152938501519380851115612cd1578c8dfd5b5050612ce18d8a858701016129ee565b91810191909152865250509284019290840190600101612c4c565b600082601f830112612d0c578081fd5b81356020612d1c6127e3836140ed565b8281528181019085830183850287018401881015612d38578586fd5b855b858110156129e157813584529284019290840190600101612d3a565b803561170781614185565b805161170781614185565b600082601f830112612d7c578081fd5b8135612d8a6127e38261410a565b818152846020838601011115612d9e578283fd5b816020850160208301379081016020019190915292915050565b600082601f830112612dc8578081fd5b8151612dd66127e38261410a565b818152846020838601011115612dea578283fd5b612373826020830160208701614144565b600060a08284031215612e0c578081fd5b60405160a081016001600160401b038282108183111715612e2957fe5b8160405282935084358352602091508185013581811115612e4957600080fd5b85019050601f81018613612e5c57600080fd5b8035612e6a6127e3826140ed565b81815283810190838501858402850186018a1015612e8757600080fd5b600094505b83851015612eb3578035612e9f81614170565b835260019490940193918501918501612e8c565b5080858701525050505050612eca60408401612f55565b6040820152612edb60608401612911565b6060820152612eec60808401612f55565b60808201525092915050565b600060608284031215612f09578081fd5b604051606081018181106001600160401b0382111715612f2557fe5b6040529050808235612f36816141a7565b8082525060208301356020820152604083013560408201525092915050565b803561170781614193565b805161170781614193565b8035611707816141a7565b8051611707816141a7565b600060208284031215612f92578081fd5b81356117c781614170565b60008060408385031215612faf578081fd5b8235612fba81614170565b946020939093013593505050565b60008060008060808587031215612fdd578182fd5b8435612fe881614170565b966020860135965060408601359560600135945092505050565b60006020808385031215613014578182fd5b82516001600160401b03811115613029578283fd5b8301601f81018513613039578283fd5b80516130476127e3826140ed565b8181528381019083850185840285018601891015613063578687fd5b8694505b83851015613085578051835260019490940193918501918501613067565b50979650505050505050565b6000602082840312156130a2578081fd5b81516001600160401b038111156130b7578182fd5b61237384828501612c20565b6000602082840312156130d4578081fd5b81516117c781614185565b6000602082840312156130f0578081fd5b5035919050565b60008060006060848603121561310b578081fd5b8335925060208401356001600160401b03811115613127578182fd5b61313386828701612b44565b925050604084013590509250925092565b600060208284031215613155578081fd5b81356001600160401b038082111561316b578283fd5b8184019150610120808387031215613181578384fd5b61318a816140ca565b905082358152602083013560208201526040830135828111156131ab578485fd5b6131b787828601612d6c565b604083015250606083013560608201526080830135608082015260a083013560a082015260c0830135828111156131ec578485fd5b6131f887828601612d6c565b60c08301525060e083013560e0820152610100808401358381111561321b578586fd5b61322788828701612cfc565b918301919091525095945050505050565b60008060006040848603121561324c578081fd5b83356001600160401b0380821115613262578283fd5b9085019060a08288031215613275578283fd5b9093506020850135908082111561328a578283fd5b818601915086601f83011261329d578283fd5b8135818111156132ab578384fd5b87602080830285010111156132be578384fd5b6020830194508093505050509250925092565b600080604083850312156132e3578182fd5b82356001600160401b03808211156132f9578384fd5b61330586838701612dfb565b9350602085013591508082111561331a578283fd5b5061332785828601612b25565b9150509250929050565b600080600060a08486031215613345578081fd5b83356001600160401b038082111561335b578283fd5b61336787838801612dfb565b9450602086013591508082111561337c578283fd5b5061338986828701612b25565b9250506133998560408601612ef8565b90509250925092565b6000602082840312156133b3578081fd5b81516001600160401b03808211156133c9578283fd5b90830190608082860312156133dc578283fd5b6040516080810181811083821117156133f157fe5b604052825182811115613402578485fd5b61340e87828601612c20565b825250602083015182811115613422578485fd5b61342e87828601612db8565b60208301525061344060408401612f60565b604082015261345160608401612d61565b606082015295945050505050565b600080600080600060a08688031215613476578283fd5b8535945060208601356001600160401b0380821115613493578485fd5b61349f89838a0161291c565b955060408801359150808211156134b4578485fd5b6134c089838a0161291c565b94506060880135935060808801359150808211156134dc578283fd5b506134e988828901612cfc565b9150509295509295909350565b60008060006060848603121561350a578081fd5b8335925060208401356001600160401b0380821115613527578283fd5b6135338783880161291c565b93506040860135915080821115613548578283fd5b5061355586828701612cfc565b9150509250925092565b600080600080600060a08688031215613576578283fd5b853594506020860135935060408601356001600160401b038082111561359a578485fd5b6134c089838a01612d6c565b6001600160a01b03169052565b60008284526020808501945082825b858110156135f05781356135d581614170565b6001600160a01b0316875295820195908201906001016135c2565b509495945050505050565b6000815180845260208085019450808401835b838110156135f05781516001600160a01b03168752958201959082019060010161360e565b6000815180845260208085018081965082840281019150828601855b858110156136a6578284038952815180518552858101518686015260408082015160ff169086015260609081015160809186018290529061369281870183613791565b9a87019a955050509084019060010161364f565b5091979650505050505050565b6000815180845260208085018081965082840281019150828601855b858110156136a65782840389528151604081518187526136f18288018261381c565b928801519688019690965250988501989350908401906001016136cf565b6000815180845260208085018081965082840281019150828601855b858110156136a6578284038952815180516001600160a01b031685528581015160608787018190529061376082880182613791565b9150506040808301519250868203818801525061377d8183613633565b9a87019a955050509084019060010161372b565b600081518084526137a9816020860160208601614144565b601f01601f19169290920160200192915050565b600081518352602082015160a060208501526137dc60a08501826135fb565b60408481015165ffffffffffff908116918701919091526060808601516001600160a01b0316908701526080948501511693909401929092525090919050565b6000815160808452613831608085018261370f565b90506020830151848203602086015261384a8282613791565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b65ffffffffffff169052565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b6000608082526139236080830187613633565b82810360208401526139358187613633565b905082810360408401526139498186613633565b91505082606083015295945050505050565b60006080825261396e6080830187613633565b851515602084015282810360408401526139498186613633565b6000602082526117c7602083018461370f565b90815260200190565b600086825260a060208301526139bd60a0830187613791565b82810360408401526139cf818761370f565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b6020808252601a908201527f7661726961626c6550617274206e6f7420746865206c6173742e000000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b600060408252833560408301526020840135601e19853603018112613e68578182fd5b840180356001600160401b03811115613e7f578283fd5b602081023603861315613e90578283fd5b60a06060850152613ea860e0850182602085016135b3565b915050613eb760408601612f55565b613ec46080850182613875565b50613ed160608601612911565b613ede60a08501826135a6565b50613eeb60808601612f55565b613ef860c0850182613875565b5082810360208401526103e981856136b3565b600060408252613f1e60408301856137bd565b82810360208401526103e981856136b3565b6000602082526117c7602083018461381c565b600086825260a06020830152613f5c60a08301876135fb565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff80881683526020818816818501526040915086151582850152606060a081860152613fe660a08601886137bd565b858103608087015286518082528382019084810283018501858a01885b8381101561408b57858303601f19018552815180518a85526140278b86018261381c565b918a0151858303868c01528051808452908b01928d92508b01905b80831015614076578351805160ff1683528c8101518d8401528d01518d830152928b019260019290920191908a0190614042565b50968a01969450505090870190600101614003565b50909e9d5050505050505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b03811182821017156140e557fe5b604052919050565b60006001600160401b0382111561410057fe5b5060209081020190565b60006001600160401b0382111561411d57fe5b50601f01601f191660200190565b60006117c73684846127d5565b60006117b03683612dfb565b60005b8381101561415f578181015183820152602001614147565b83811115610ce35750506000910152565b6001600160a01b038116811461193f57600080fd5b801515811461193f57600080fd5b65ffffffffffff8116811461193f57600080fd5b60ff8116811461193f57600080fdfeb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679a26469706673582212205eeab93441f63f7cace90181b86061e591c1d61fd1e6c7a0858d434ada313bd064736f6c63430007060033",
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

// LatestSupportedState is a free data retrieval call binding the contract method 0x0d96cc22.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_NitroAdjudicator *NitroAdjudicatorCaller) LatestSupportedState(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "latestSupportedState", fixedPart, signedVariableParts)

	if err != nil {
		return *new(INitroTypesVariablePart), err
	}

	out0 := *abi.ConvertType(out[0], new(INitroTypesVariablePart)).(*INitroTypesVariablePart)

	return out0, err

}

// LatestSupportedState is a free data retrieval call binding the contract method 0x0d96cc22.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_NitroAdjudicator *NitroAdjudicatorSession) LatestSupportedState(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	return _NitroAdjudicator.Contract.LatestSupportedState(&_NitroAdjudicator.CallOpts, fixedPart, signedVariableParts)
}

// LatestSupportedState is a free data retrieval call binding the contract method 0x0d96cc22.
//
// Solidity: function latestSupportedState((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) pure returns(((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool))
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) LatestSupportedState(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (INitroTypesVariablePart, error) {
	return _NitroAdjudicator.Contract.LatestSupportedState(&_NitroAdjudicator.CallOpts, fixedPart, signedVariableParts)
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

// Challenge is a paid mutator transaction binding the contract method 0x12c52592.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, signedVariableParts, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x12c52592.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x12c52592.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x3503f8f0.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, signedVariableParts)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x3503f8f0.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x3503f8f0.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Checkpoint(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
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

// Conclude is a paid mutator transaction binding the contract method 0xc927c2d3.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", fixedPart, signedVariableParts)
}

// Conclude is a paid mutator transaction binding the contract method 0xc927c2d3.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
}

// Conclude is a paid mutator transaction binding the contract method 0xc927c2d3.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x0d36d531.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", fixedPart, signedVariableParts)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x0d36d531.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x0d36d531.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, signedVariableParts []INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, signedVariableParts)
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
	ChannelId           [32]byte
	TurnNumRecord       *big.Int
	FinalizesAt         *big.Int
	IsFinal             bool
	FixedPart           INitroTypesFixedPart
	SignedVariableParts []INitroTypesSignedVariablePart
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterChallengeRegistered is a free log retrieval operation binding the contract event 0x5bc2767c561afcb4ffbef365ef0757d1d4389267288cbedb1014507e356ff1ec.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts)
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

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0x5bc2767c561afcb4ffbef365ef0757d1d4389267288cbedb1014507e356ff1ec.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts)
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

// ParseChallengeRegistered is a log parse operation binding the contract event 0x5bc2767c561afcb4ffbef365ef0757d1d4389267288cbedb1014507e356ff1ec.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] signedVariableParts)
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

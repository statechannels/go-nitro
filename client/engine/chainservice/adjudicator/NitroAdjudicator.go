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
	SourceChannelId       [32]byte
	SourceStateHash       [32]byte
	SourceOutcomeBytes    []byte
	SourceAssetIndex      *big.Int
	IndexOfTargetInSource *big.Int
	TargetStateHash       [32]byte
	TargetOutcomeBytes    []byte
	TargetAssetIndex      *big.Int
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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"}],\"name\":\"Reclaimed\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"}],\"name\":\"compute_reclaim_effects\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"requireStateSupported\",\"outputs\":[],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50613bf5806100206000396000f3fe6080604052600436106100e85760003560e01c8063566d54c61161008a578063998240af11610059578063998240af14610279578063af69c9d714610299578063c7df14e2146102b9578063d3c4e738146102d9576100e8565b8063566d54c6146101ec57806359e168f3146102195780637d3f9f281461023957806385554cc814610259576100e8565b80632fb1d270116100c65780632fb1d270146101755780633033730e14610188578063552cfa50146101a8578063564b81ef146101d7576100e8565b8063103a40c4146100ed57806311e9f1781461010f578063166e56cd14610148575b600080fd5b3480156100f957600080fd5b5061010d610108366004612d86565b6102f9565b005b34801561011b57600080fd5b5061012f61012a366004612e62565b61044c565b60405161013f9493929190613345565b60405180910390f35b34801561015457600080fd5b50610168610163366004612931565b61079a565b60405161013f9190613397565b61010d61018336600461295c565b6107b7565b34801561019457600080fd5b5061010d6101a3366004612ec1565b610a46565b3480156101b457600080fd5b506101c86101c3366004612b22565b610ac6565b60405161013f93929190613ab0565b3480156101e357600080fd5b50610168610ae1565b3480156101f857600080fd5b5061020c610207366004612996565b610af0565b60405161013f9190613332565b34801561022557600080fd5b5061010d610234366004612d03565b610d6a565b34801561024557600080fd5b5061010d610254366004612c46565b610d90565b34801561026557600080fd5b5061010d610274366004612d03565b610e33565b34801561028557600080fd5b5061010d610294366004612d03565b610e3e565b3480156102a557600080fd5b5061010d6102b4366004612b3a565b610e81565b3480156102c557600080fd5b506101686102d4366004612b22565b611217565b3480156102e557600080fd5b5061010d6102f4366004612b76565b611229565b600061030485611299565b83516040015190915060006103188361130f565b600281111561032357fe5b141561033857610333828261135a565b610367565b60016103438361130f565b600281111561034e57fe5b141561035e57610333828261139e565b610367826113dc565b610372868686611413565b6000610382878660000151611485565b9050610393818860200151866114d2565b827f84b316f9a919207471680bd56ff706bd0dbbfce1c925b9253027d9e09274560a83896080015142018860000151606001518b8b8b6040516103db96959493929190613a0a565b60405180910390a261043160405180608001604052808465ffffffffffff1681526020018960800151420165ffffffffffff16815260200183815260200161042a88600001516000015161152c565b9052611545565b60009384526020849052604090932092909255505050505050565b606060006060600080855111610463578551610466565b84515b6001600160401b038111801561047b57600080fd5b506040519080825280602002602001820160405280156104b557816020015b6104a2612170565b81526020019060019003908161049a5790505b5091506000905085516001600160401b03811180156104d357600080fd5b5060405190808252806020026020018201604052801561050d57816020015b6104fa612170565b8152602001906001900390816104f25790505b50935060019250866000805b885181101561078e5788818151811061052e57fe5b60200260200101516000015187828151811061054657fe5b6020026020010151600001818152505088818151811061056257fe5b60200260200101516040015187828151811061057a57fe5b60200260200101516040019060ff16908160ff168152505088818151811061059e57fe5b6020026020010151606001518782815181106105b657fe5b60200260200101516060018190525060006105e88a83815181106105d657fe5b60200260200101516020015185611596565b9050885160001480610617575088518310801561061757508189848151811061060d57fe5b6020026020010151145b1561072957600260ff168a848151811061062d57fe5b60200260200101516040015160ff1614156106635760405162461bcd60e51b815260040161065a90613689565b60405180910390fd5b808a838151811061067057fe5b6020026020010151602001510388838151811061068957fe5b6020026020010151602001818152505060405180608001604052808b84815181106106b057fe5b60200260200101516000015181526020018281526020018b84815181106106d357fe5b60200260200101516040015160ff1681526020018b84815181106106f357fe5b60200260200101516060015181525086848151811061070e57fe5b6020026020010181905250808501945082600101925061075e565b89828151811061073557fe5b60200260200101516020015188838151811061074d57fe5b602002602001015160200181815250505b87828151811061076a57fe5b60200260200101516020015160001461078257600096505b90920391600101610519565b50505093509350935093565b600160209081526000928352604080842090915290825290205481565b6107c0836115ae565b156107dd5760405162461bcd60e51b815260040161065a90613751565b6001600160a01b0384166000908152600160209081526040808320868452909152812054838110156108215760405162461bcd60e51b815260040161065a9061346c565b61082b84846115b5565b81106108495760405162461bcd60e51b815260040161065a906137b8565b61085d8161085786866115b5565b9061160f565b91506001600160a01b0386166108915782341461088c5760405162461bcd60e51b815260040161065a90613826565b61092f565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906108c1903390309087906004016132d4565b602060405180830381600087803b1580156108db57600080fd5b505af11580156108ef573d6000803e3d6000fd5b505050506040513d601f19601f820116820180604052508101906109139190612b06565b61092f5760405162461bcd60e51b815260040161065a9061371a565b600061093b82846115b5565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a7159061099a908a9087908690613311565b60405180910390a26001600160a01b038716610a3d5760006109bc858561160f565b90506000336001600160a01b0316826040516109d79061205d565b60006040518083038185875af1925050503d8060008114610a14576040519150601f19603f3d011682016040523d82523d6000602084013e610a19565b606091505b5050905080610a3a5760405162461bcd60e51b815260040161065a9061385d565b50505b50505050505050565b6000806000610a58888589888a61166c565b9250925092506000806000610a8584878d81518110610a7357fe5b6020026020010151604001518961044c565b93509350509250610a9c8b868c8b8a888a886116ea565b610ab9868c81518110610aab57fe5b6020026020010151836117a8565b5050505050505050505050565b6000806000610ad4846117e0565b9196909550909350915050565b6000610aeb6117fe565b905090565b6060600060018551036001600160401b0381118015610b0e57600080fd5b50604051908082528060200260200182016040528015610b4857816020015b610b35612170565b815260200190600190039081610b2d5790505b5090506000858481518110610b5957fe5b602002602001015190506000610b728260600151611802565b905060008060008060005b8b51811015610d035789811415610b975760019450610cfb565b60405180608001604052808d8381518110610bae57fe5b60200260200101516000015181526020018d8381518110610bcb57fe5b60200260200101516020015181526020018d8381518110610be857fe5b60200260200101516040015160ff1681526020018d8381518110610c0857fe5b602002602001015160600151815250888381518110610c2357fe5b602002602001018190525085600001518c8281518110610c3f57fe5b6020026020010151600001511415610c91578a600081518110610c5e57fe5b602002602001015160200151888381518110610c7657fe5b60200260200101516020018181510191508181525050600193505b85602001518c8281518110610ca257fe5b6020026020010151600001511415610cf4578a600181518110610cc157fe5b602002602001015160200151888381518110610cd957fe5b60200260200101516020018181510191508181525050600192505b6001909101905b600101610b7d565b5083610d215760405162461bcd60e51b815260040161065a90613623565b82610d3e5760405162461bcd60e51b815260040161065a906134a3565b81610d5b5760405162461bcd60e51b815260040161065a90613527565b50949998505050505050505050565b6000610d7784848461181e565b825151909150610d8a9082906000610e81565b50505050565b6000610d9c8385613b3c565b9050610dae6080860160608701612915565b6001600160a01b03166350d8c32a86610dcf610dc982613b49565b856118d9565b610de9610ddb8a613b49565b610de488613b55565b611981565b6040518463ffffffff1660e01b8152600401610e0793929190613894565b60006040518083038186803b158015610e1f57600080fd5b505afa158015610a3a573d6000803e3d6000fd5b610d8a83838361181e565b6000610e4984611299565b825160400151909150610e5b826113dc565b610e65828261139e565b610e70858585611413565b610e7a8282611a46565b5050505050565b610e8a83611acd565b610e9d81610e978461152c565b85611b00565b81516001906000906001600160401b0381118015610eba57600080fd5b50604051908082528060200260200182016040528015610ef457816020015b610ee1612196565b815260200190600190039081610ed95790505b509050600084516001600160401b0381118015610f1057600080fd5b50604051908082528060200260200182016040528015610f3a578160200160208202803683370190505b509050600085516001600160401b0381118015610f5657600080fd5b50604051908082528060200260200182016040528015610f80578160200160208202803683370190505b50905060005b865181101561111a576000878281518110610f9d57fe5b602002602001015190506000816040015190506000898481518110610fbe57fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c81526020019081526020016000205486858151811061100f57fe5b6020026020010181815250506000806000806110808a898151811061103057fe5b60200260200101518760006001600160401b038111801561105057600080fd5b5060405190808252806020026020018201604052801561107a578160200160208202803683370190505b5061044c565b9350935093509350826110925760009b505b8089898151811061109f57fe5b602002602001018181525050838e89815181106110b857fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b89815181106110fb57fe5b6020026020010181905250505050505050508080600101915050610f86565b5060005b86518110156111e057600087828151811061113557fe5b602002602001015160000151905082828151811061114f57fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d8352909352919091208054919091039055835189907fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c1306376799084908790829081106111b957fe5b60200260200101516040516111cf9291906139e9565b60405180910390a25060010161111e565b5083156111fb5760008781526020819052604081205561120e565b61120e87866112098961152c565b611b49565b610a3d83611baf565b60006020819052908152604090205481565b60008061123583611bdf565b91509150606060008385606001518151811061124d57fe5b60200260200101516040015190506000838660e001518151811061126d57fe5b602002602001015160400151905061128a82828860800151610af0565b92505050610d8a848483611d53565b60006112a36117fe565b8251146112c25760405162461bcd60e51b815260040161065a906134d0565b6112ca6117fe565b82602001518360400151846060015185608001516040516020016112f295949392919061399e565b604051602081830303815290604052805190602001209050919050565b60008061131b836117e0565b5091505065ffffffffffff8116611336576000915050611355565b428165ffffffffffff161161134f576002915050611355565b60019150505b919050565b6000611365836117e0565b505090508065ffffffffffff168265ffffffffffff1610156113995760405162461bcd60e51b815260040161065a90613555565b505050565b60006113a9836117e0565b505090508065ffffffffffff168265ffffffffffff16116113995760405162461bcd60e51b815260040161065a90613435565b60026113e78261130f565b60028111156113f257fe5b14156114105760405162461bcd60e51b815260040161065a906134fb565b50565b82606001516001600160a01b03166350d8c32a8461143186866118d9565b61143b8786611981565b6040518463ffffffff1660e01b815260040161145993929190613979565b60006040518083038186803b15801561147157600080fd5b505afa158015610a3d573d6000803e3d6000fd5b600061149083611299565b6020808401518451604080870151606088015191516114b39695919291016133a0565b6040516020818303038152906040528051906020012090505b92915050565b6000611504846040516020016114e891906133ed565b6040516020818303038152906040528051906020012083611dd7565b90506115108184611e91565b610d8a5760405162461bcd60e51b815260040161065a906135bc565b600061153782611ee7565b805190602001209050919050565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161158291611f10565b6001600160a01b0316919091179392505050565b60008183116115a557826115a7565b815b9392505050565b60a01c1590565b6000828201838110156115a7576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611666576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b606060008061167a87611f3c565b61168386611acd565b61169585858051906020012088611b00565b61169e84611f9b565b92508288815181106116ac57fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a90811061172457fe5b6020026020010151604001819052506117648686866040516020016117499190613384565b60405160208183030381529060405280519060200120611b49565b857fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c13063767989846040516117969291906139e9565b60405180910390a25050505050505050565b6117dc604051806060016040528084600001516001600160a01b031681526020018460200151815260200183815250611fb1565b5050565b60009081526020819052604090205460d081901c9160a082901c9190565b4690565b61180a6121c0565b818060200190518101906114cc9190612e1b565b600061182984611299565b9050611834816113dc565b61183f848484611413565b6118846040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b815260200161042a85600001516000015161152c565b60008083815260200190815260200160002081905550807f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516118ca91906139f7565b60405180910390a29392505050565b6060600082516001600160401b03811180156118f457600080fd5b5060405190808252806020026020018201604052801561192e57816020015b61191b6121d7565b8152602001906001900390816119135790505b50905060005b83518110156119795761195a8585838151811061194d57fe5b6020026020010151611981565b82828151811061196657fe5b6020908102919091010152600101611934565b509392505050565b6119896121d7565b60408051808201909152825181526000602082018190525b8360200151518110156119795760006119de6119c1878760000151611485565b866020015184815181106119d157fe5b6020026020010151611dd7565b905060005b866020015151811015611a3c5786602001518181518110611a0057fe5b60200260200101516001600160a01b0316826001600160a01b03161415611a3457602084018051600283900a019052611a3c565b6001016119e3565b50506001016119a1565b6040805160808101825265ffffffffffff831681526000602082018190529181018290526060810191909152611a7b90611545565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e082604051611ac191906139f7565b60405180910390a25050565b6002611ad88261130f565b6002811115611ae357fe5b146114105760405162461bcd60e51b815260040161065a90613788565b6000611b0b826117e0565b92505050611b198484611f10565b6001600160a01b0316816001600160a01b031614610d8a5760405162461bcd60e51b815260040161065a906136eb565b600080611b55856117e0565b50915091506000611b9560405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250611545565b600096875260208790526040909620959095555050505050565b60005b81518110156117dc57611bd7828281518110611bca57fe5b6020026020010151611fb1565b600101611bb2565b8051604082015160608381015160c085015160e086015192948594909390929190611c0985611acd565b611c1f8860200151858051906020012087611b00565b611c2884611f9b565b9650611c3382611f9b565b95506000878481518110611c4357fe5b6020908102919091010151519050600260ff16888581518110611c6257fe5b6020026020010151604001518a6080015181518110611c7d57fe5b60200260200101516040015160ff1614611ca95760405162461bcd60e51b815260040161065a906137ef565b6000888581518110611cb757fe5b6020026020010151604001518a6080015181518110611cd257fe5b6020026020010151600001519050816001600160a01b0316888481518110611cf657fe5b6020026020010151600001516001600160a01b031614611d285760405162461bcd60e51b815260040161065a90613652565b611d3181611acd565b611d478a60a00151858051906020012083611b00565b50505050505050915091565b8251606084015183518390859083908110611d6a57fe5b602002602001015160400181905250611d93828660200151866040516020016117499190613384565b845160608601516040517f4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b91611dc891613397565b60405180910390a25050505050565b60008083604051602001611deb91906132a3565b604051602081830303815290604052805190602001209050600060018285600001518660200151876040015160405160008152602001604052604051611e349493929190613417565b6020604051602081039080840390855afa158015611e56573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116611e895760405162461bcd60e51b815260040161065a906136c0565b949350505050565b6000805b8251811015611edd57828181518110611eaa57fe5b60200260200101516001600160a01b0316846001600160a01b03161415611ed55760019150506114cc565b600101611e95565b5060009392505050565b606081604051602001611efa9190613384565b6040516020818303038152906040529050919050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60005b81518160010110156117dc57818160010181518110611f5a57fe5b6020026020010151828281518110611f6e57fe5b602002602001015110611f935760405162461bcd60e51b815260040161065a9061358c565b600101611f3f565b6060818060200190518101906114cc91906129ff565b805160005b82604001515181101561139957600083604001518281518110611fd557fe5b6020026020010151600001519050600084604001518381518110611ff557fe5b602002602001015160200151905061200c826115ae565b15612029576120248461201e8461205d565b83612060565b612053565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b5050600101611fb6565b90565b6001600160a01b0383166120f0576000826001600160a01b0316826040516120879061205d565b60006040518083038185875af1925050503d80600081146120c4576040519150601f19603f3d011682016040523d82523d6000602084013e6120c9565b606091505b50509050806120ea5760405162461bcd60e51b815260040161065a906135f3565b50611399565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb9061211e90859085906004016132f8565b602060405180830381600087803b15801561213857600080fd5b505af115801561214c573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610d8a9190612b06565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b604080518082019091526000808252602082015290565b60405180604001604052806121ea6121f7565b8152602001600081525090565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b600061223c61223784613afe565b613adb565b83815290506020808201908360005b868110156122725761226088833588016127eb565b8452928201929082019060010161224b565b505050509392505050565b803561135581613b8d565b600082601f830112612298578081fd5b813560206122a861223783613afe565b82815281810190858301855b8581101561234d5781358801608080601f19838d030112156122d4578889fd5b604080518281016001600160401b0382821081831117156122f157fe5b908352848a0135825284830135828b01526060906123108287016128ff565b83850152938501359380851115612325578c8dfd5b506123348e8b868801016125f1565b90820152875250505092840192908401906001016122b4565b5090979650505050505050565b600082601f83011261236a578081fd5b8151602061237a61223783613afe565b82815281810190858301855b8581101561234d5781518801608080601f19838d030112156123a6578889fd5b604080518281016001600160401b0382821081831117156123c357fe5b908352848a0151825284830151828b01526060906123e282870161290a565b838501529385015193808511156123f7578c8dfd5b506124068e8b8688010161263d565b9082015287525050509284019290840190600101612386565b600082601f83011261242f578081fd5b8135602061243f61223783613afe565b8281528181019085830160608086028801850189101561245d578687fd5b865b86811015612483576124718a8461277d565b8552938501939181019160010161245f565b509198975050505050505050565b600082601f8301126124a1578081fd5b6115a783833560208501612229565b600082601f8301126124c0578081fd5b813560206124d061223783613afe565b82815281810190858301855b8581101561234d5781358801606080601f19838d030112156124fc578889fd5b604080518281016001600160401b03828210818311171561251957fe5b908352848a01359061252a82613b8d565b90825284830135908082111561253e578c8dfd5b61254c8f8c848901016125f1565b838c0152938501359380851115612561578c8dfd5b50506125718d8a85870101612288565b918101919091528652505092840192908401906001016124dc565b600082601f83011261259c578081fd5b813560206125ac61223783613afe565b82815281810190858301838502870184018810156125c8578586fd5b855b8581101561234d578135845292840192908401906001016125ca565b803561135581613ba2565b600082601f830112612601578081fd5b813561260f61223782613b1b565b818152846020838601011115612623578283fd5b816020850160208301379081016020019190915292915050565b600082601f83011261264d578081fd5b815161265b61223782613b1b565b81815284602083860101111561266f578283fd5b611e89826020830160208701613b61565b600060a08284031215612691578081fd5b60405160a081016001600160401b0382821081831117156126ae57fe5b81604052829350843583526020915081850135818111156126ce57600080fd5b85019050601f810186136126e157600080fd5b80356126ef61223782613afe565b81815283810190838501858402850186018a101561270c57600080fd5b600094505b8385101561273857803561272481613b8d565b835260019490940193918501918501612711565b508085870152505050505061274f604084016128e9565b60408201526127606060840161227d565b6060820152612771608084016128e9565b60808201525092915050565b60006060828403121561278e578081fd5b604051606081018181106001600160401b03821117156127aa57fe5b60405290508082356127bb81613bb0565b8082525060208301356020820152604083013560408201525092915050565b6000604082840312156115a5578081fd5b600060408083850312156127fd578182fd5b80518181016001600160401b03828210818311171561281857fe5b81845282945085358181111561282d57600080fd5b86016080818903121561283f57600080fd5b60c08401838110838211171561285157fe5b855280358281111561286257600080fd5b61286e898284016124b0565b84525060208101358281111561288357600080fd5b61288f898284016125f1565b6060860152506128a08582016128e9565b60808501526128b1606082016125e6565b60a08501525081835260208601359350808411156128ce57600080fd5b50506128dc8583860161241f565b6020820152505092915050565b803565ffffffffffff8116811461135557600080fd5b803561135581613bb0565b805161135581613bb0565b600060208284031215612926578081fd5b81356115a781613b8d565b60008060408385031215612943578081fd5b823561294e81613b8d565b946020939093013593505050565b60008060008060808587031215612971578182fd5b843561297c81613b8d565b966020860135965060408601359560600135945092505050565b6000806000606084860312156129aa578081fd5b83356001600160401b03808211156129c0578283fd5b6129cc87838801612288565b945060208601359150808211156129e1578283fd5b506129ee86828701612288565b925050604084013590509250925092565b60006020808385031215612a11578182fd5b82516001600160401b0380821115612a27578384fd5b818501915085601f830112612a3a578384fd5b8151612a4861223782613afe565b81815284810190848601875b84811015612af757815187016060818d03601f19011215612a7357898afd5b60408051606081018181108a82111715612a8957fe5b8252828b0151612a9881613b8d565b81528282015189811115612aaa578c8dfd5b612ab88f8d8387010161263d565b8c83015250606083015189811115612ace578c8dfd5b612adc8f8d8387010161235a565b92820192909252865250509287019290870190600101612a54565b50909998505050505050505050565b600060208284031215612b17578081fd5b81516115a781613ba2565b600060208284031215612b33578081fd5b5035919050565b600080600060608486031215612b4e578081fd5b8335925060208401356001600160401b03811115612b6a578182fd5b6129ee868287016124b0565b600060208284031215612b87578081fd5b81356001600160401b0380821115612b9d578283fd5b8184019150610100808387031215612bb3578384fd5b612bbc81613adb565b90508235815260208301356020820152604083013582811115612bdd578485fd5b612be9878286016125f1565b604083015250606083013560608201526080830135608082015260a083013560a082015260c083013582811115612c1e578485fd5b612c2a878286016125f1565b60c08301525060e083013560e082015280935050505092915050565b60008060008060608587031215612c5b578182fd5b84356001600160401b0380821115612c71578384fd5b9086019060a08289031215612c84578384fd5b90945060208601359080821115612c99578384fd5b818701915087601f830112612cac578384fd5b813581811115612cba578485fd5b8860208083028501011115612ccd578485fd5b602083019550809450506040870135915080821115612cea578283fd5b50612cf7878288016127da565b91505092959194509250565b600080600060608486031215612d17578081fd5b83356001600160401b0380821115612d2d578283fd5b612d3987838801612680565b94506020860135915080821115612d4e578283fd5b612d5a87838801612491565b93506040860135915080821115612d6f578283fd5b50612d7c868287016127eb565b9150509250925092565b60008060008060c08587031215612d9b578182fd5b84356001600160401b0380821115612db1578384fd5b612dbd88838901612680565b95506020870135915080821115612dd2578384fd5b612dde88838901612491565b94506040870135915080821115612df3578384fd5b50612e00878288016127eb565b925050612e10866060870161277d565b905092959194509250565b600060408284031215612e2c578081fd5b604051604081018181106001600160401b0382111715612e4857fe5b604052825181526020928301519281019290925250919050565b600080600060608486031215612e76578081fd5b8335925060208401356001600160401b0380821115612e93578283fd5b612e9f87838801612288565b93506040860135915080821115612eb4578283fd5b50612d7c8682870161258c565b600080600080600060a08688031215612ed8578283fd5b853594506020860135935060408601356001600160401b0380821115612efc578485fd5b612f0889838a016125f1565b9450606088013593506080880135915080821115612f24578283fd5b50612f318882890161258c565b9150509295509295909350565b6001600160a01b03169052565b60008284526020808501945082825b85811015612f88578135612f6d81613b8d565b6001600160a01b031687529582019590820190600101612f5a565b509495945050505050565b6000815180845260208085019450808401835b83811015612f885781516001600160a01b031687529582019590820190600101612fa6565b6000815180845260208085018081965082840281019150828601855b8581101561303e578284038952815180518552858101518686015260408082015160ff169086015260609081015160809186018290529061302a81870183613113565b9a87019a9550505090840190600101612fe7565b5091979650505050505050565b6000815180845260208085018081965082840281019150828601855b8581101561303e57828403895261307f84835161319e565b98850198935090840190600101613067565b6000815180845260208085018081965082840281019150828601855b8581101561303e578284038952815180516001600160a01b03168552858101516060878701819052906130e282880182613113565b915050604080830151925086820381880152506130ff8183612fcb565b9a87019a95505050908401906001016130ad565b6000815180845261312b816020860160208601613b61565b601f01601f19169290920160200192915050565b600081518352602082015160a0602085015261315e60a0850182612f93565b60408481015165ffffffffffff908116918701919091526060808601516001600160a01b0316908701526080948501511693909401929092525090919050565b60008151604084526131b3604085018261323e565b602093840151949093019390935250919050565b6000604082518185526131dc8286018261323e565b602085810151878303888301528051808452908201935090918591908301905b80831015613232578451805160ff16835284810151858401528601518683015293830193600192909201916060909101906131fc565b50979650505050505050565b60008151608084526132536080850182613091565b90506020830151848203602086015261326c8282613113565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b65ffffffffffff169052565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b6000602082526115a76020830184612fcb565b6000608082526133586080830187612fcb565b851515602084015282810360408401526133728186612fcb565b91505082606083015295945050505050565b6000602082526115a76020830184613091565b90815260200190565b600086825260a060208301526133b960a0830187613113565b82810360408401526133cb8187613091565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b60208082526013908201527218dbdd5b19081b9bdd08199a5b99081b19599d606a1b604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526014908201527318dbdd5b19081b9bdd08199a5b99081c9a59da1d60621b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b60208082526015908201527418dbdd5b19081b9bdd08199a5b99081d185c99d95d605a1b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b600060608252843560608301526020850135601e198636030181126138b7578182fd5b850180356001600160401b038111156138ce578283fd5b6020810236038713156138df578283fd5b60a060808501526138f861010085018260208501612f4b565b915050613907604087016128e9565b61391460a0850182613297565b506139216060870161227d565b61392e60c0850182612f3e565b5061393b608087016128e9565b61394860e0850182613297565b50828103602084015261395b818661304b565b9050828103604084015261396f818561319e565b9695505050505050565b60006060825261398c606083018661313f565b828103602084015261395b818661304b565b600086825260a060208301526139b760a0830187612f93565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff8089168352602081891681850152871515604085015260c06060850152613a3d60c085018861313f565b848103608086015286518082529092508183019082810284018301838901865b83811015613a8b57601f19878403018552613a798383516131c7565b94860194925090850190600101613a5d565b505086810360a0880152613a9f81896131c7565b9d9c50505050505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b0381118282101715613af657fe5b604052919050565b60006001600160401b03821115613b1157fe5b5060209081020190565b60006001600160401b03821115613b2e57fe5b50601f01601f191660200190565b60006115a7368484612229565b60006114cc3683612680565b60006114cc36836127eb565b60005b83811015613b7c578181015183820152602001613b64565b83811115610d8a5750506000910152565b6001600160a01b038116811461141057600080fd5b801515811461141057600080fd5b60ff8116811461141057600080fdfea26469706673582212204bb9bd13dd3e47e0d5682bfecef686babf658e4bc1c6c88463f4242212b5c29f64736f6c63430007060033",
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

// ComputeReclaimEffects is a free data retrieval call binding the contract method 0x566d54c6.
//
// Solidity: function compute_reclaim_effects((bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource) pure returns((bytes32,uint256,uint8,bytes)[])
func (_NitroAdjudicator *NitroAdjudicatorCaller) ComputeReclaimEffects(opts *bind.CallOpts, sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int) ([]ExitFormatAllocation, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "compute_reclaim_effects", sourceAllocations, targetAllocations, indexOfTargetInSource)

	if err != nil {
		return *new([]ExitFormatAllocation), err
	}

	out0 := *abi.ConvertType(out[0], new([]ExitFormatAllocation)).(*[]ExitFormatAllocation)

	return out0, err

}

// ComputeReclaimEffects is a free data retrieval call binding the contract method 0x566d54c6.
//
// Solidity: function compute_reclaim_effects((bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource) pure returns((bytes32,uint256,uint8,bytes)[])
func (_NitroAdjudicator *NitroAdjudicatorSession) ComputeReclaimEffects(sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int) ([]ExitFormatAllocation, error) {
	return _NitroAdjudicator.Contract.ComputeReclaimEffects(&_NitroAdjudicator.CallOpts, sourceAllocations, targetAllocations, indexOfTargetInSource)
}

// ComputeReclaimEffects is a free data retrieval call binding the contract method 0x566d54c6.
//
// Solidity: function compute_reclaim_effects((bytes32,uint256,uint8,bytes)[] sourceAllocations, (bytes32,uint256,uint8,bytes)[] targetAllocations, uint256 indexOfTargetInSource) pure returns((bytes32,uint256,uint8,bytes)[])
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) ComputeReclaimEffects(sourceAllocations []ExitFormatAllocation, targetAllocations []ExitFormatAllocation, indexOfTargetInSource *big.Int) ([]ExitFormatAllocation, error) {
	return _NitroAdjudicator.Contract.ComputeReclaimEffects(&_NitroAdjudicator.CallOpts, sourceAllocations, targetAllocations, indexOfTargetInSource)
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

// RequireStateSupported is a free data retrieval call binding the contract method 0x7d3f9f28.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) pure returns()
func (_NitroAdjudicator *NitroAdjudicatorCaller) RequireStateSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) error {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "requireStateSupported", fixedPart, proof, candidate)

	if err != nil {
		return err
	}

	return err

}

// RequireStateSupported is a free data retrieval call binding the contract method 0x7d3f9f28.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) pure returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) error {
	return _NitroAdjudicator.Contract.RequireStateSupported(&_NitroAdjudicator.CallOpts, fixedPart, proof, candidate)
}

// RequireStateSupported is a free data retrieval call binding the contract method 0x7d3f9f28.
//
// Solidity: function requireStateSupported((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) pure returns()
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) RequireStateSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) error {
	return _NitroAdjudicator.Contract.RequireStateSupported(&_NitroAdjudicator.CallOpts, fixedPart, proof, candidate)
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

// Challenge is a paid mutator transaction binding the contract method 0x103a40c4.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, proof, candidate, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x103a40c4.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x103a40c4.
//
// Solidity: function challenge((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x998240af.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, proof, candidate)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x998240af.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x998240af.
//
// Solidity: function checkpoint((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Checkpoint(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0x85554cc8.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", fixedPart, proof, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0x85554cc8.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0x85554cc8.
//
// Solidity: function conclude((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x59e168f3.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", fixedPart, proof, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x59e168f3.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0x59e168f3.
//
// Solidity: function concludeAndTransferAllAssets((uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
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

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Reclaim(opts *bind.TransactOpts, claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "reclaim", claimArgs)
}

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Reclaim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Reclaim(&_NitroAdjudicator.TransactOpts, claimArgs)
}

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) claimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Reclaim(claimArgs IMultiAssetHolderClaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Reclaim(&_NitroAdjudicator.TransactOpts, claimArgs)
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
	FixedPart     INitroTypesFixedPart
	Proof         []INitroTypesSignedVariablePart
	Candidate     INitroTypesSignedVariablePart
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterChallengeRegistered is a free log retrieval operation binding the contract event 0x84b316f9a919207471680bd56ff706bd0dbbfce1c925b9253027d9e09274560a.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0x84b316f9a919207471680bd56ff706bd0dbbfce1c925b9253027d9e09274560a.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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

// ParseChallengeRegistered is a log parse operation binding the contract event 0x84b316f9a919207471680bd56ff706bd0dbbfce1c925b9253027d9e09274560a.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (uint256,address[],uint48,address,uint48) fixedPart, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,bytes,(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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

// NitroAdjudicatorReclaimedIterator is returned from FilterReclaimed and is used to iterate over the raw logs and unpacked data for Reclaimed events raised by the NitroAdjudicator contract.
type NitroAdjudicatorReclaimedIterator struct {
	Event *NitroAdjudicatorReclaimed // Event containing the contract specifics and raw log

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
func (it *NitroAdjudicatorReclaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(NitroAdjudicatorReclaimed)
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
		it.Event = new(NitroAdjudicatorReclaimed)
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
func (it *NitroAdjudicatorReclaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *NitroAdjudicatorReclaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// NitroAdjudicatorReclaimed represents a Reclaimed event raised by the NitroAdjudicator contract.
type NitroAdjudicatorReclaimed struct {
	ChannelId  [32]byte
	AssetIndex *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterReclaimed is a free log retrieval operation binding the contract event 0x4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b.
//
// Solidity: event Reclaimed(bytes32 indexed channelId, uint256 assetIndex)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) FilterReclaimed(opts *bind.FilterOpts, channelId [][32]byte) (*NitroAdjudicatorReclaimedIterator, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.FilterLogs(opts, "Reclaimed", channelIdRule)
	if err != nil {
		return nil, err
	}
	return &NitroAdjudicatorReclaimedIterator{contract: _NitroAdjudicator.contract, event: "Reclaimed", logs: logs, sub: sub}, nil
}

// WatchReclaimed is a free log subscription operation binding the contract event 0x4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b.
//
// Solidity: event Reclaimed(bytes32 indexed channelId, uint256 assetIndex)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) WatchReclaimed(opts *bind.WatchOpts, sink chan<- *NitroAdjudicatorReclaimed, channelId [][32]byte) (event.Subscription, error) {

	var channelIdRule []interface{}
	for _, channelIdItem := range channelId {
		channelIdRule = append(channelIdRule, channelIdItem)
	}

	logs, sub, err := _NitroAdjudicator.contract.WatchLogs(opts, "Reclaimed", channelIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(NitroAdjudicatorReclaimed)
				if err := _NitroAdjudicator.contract.UnpackLog(event, "Reclaimed", log); err != nil {
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

// ParseReclaimed is a log parse operation binding the contract event 0x4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b.
//
// Solidity: event Reclaimed(bytes32 indexed channelId, uint256 assetIndex)
func (_NitroAdjudicator *NitroAdjudicatorFilterer) ParseReclaimed(log types.Log) (*NitroAdjudicatorReclaimed, error) {
	event := new(NitroAdjudicatorReclaimed)
	if err := _NitroAdjudicator.contract.UnpackLog(event, "Reclaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

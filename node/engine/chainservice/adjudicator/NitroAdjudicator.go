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

// IMultiAssetHolderReclaimArgs is an auto generated low-level Go binding around an user-defined struct.
type IMultiAssetHolderReclaimArgs struct {
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
	Participants      []common.Address
	ChannelNonce      uint64
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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"finalHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"}],\"name\":\"Reclaimed\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"}],\"name\":\"compute_reclaim_effects\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"}],\"internalType\":\"structIMultiAssetHolder.ReclaimArgs\",\"name\":\"reclaimArgs\",\"type\":\"tuple\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"stateIsSupported\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x60808060405234610016576137d2908161001c8239f35b600080fdfe608060405260048036101561001357600080fd5b60003560e01c90816311e9f1781461172e578163166e56cd146116e65781632fb1d270146113d45781633033730e1461119a57816331afa0b414610ebf578163552cfa5014610e7a578163566d54c614610e105781635685b7dc14610bbf5781636d2a9c9214610ac45781638286a0601461076d578163c7df14e214610743578163d3c4e73814610468578163ec346235146100d7575063ee049b50146100b957600080fd5b346100d2576100d06100ca36611ef5565b906123a4565b005b600080fd5b346100d2576100f1906100e936611ef5565b8093916123a4565b915151906100fe83613182565b6101078261340b565b61011084613755565b93915050604051926020938481019360008552604082015260408152610135816117bf565b519092206001600160a01b03929083169083160361042c575060019183519161015d8361187d565b9261016b604051948561185c565b80845261017a601f199161187d565b018260005b8281106103ed575050506101938551611f38565b61019d8651611f38565b9060005b8751811015610285576101b48189611fc2565b5188604082015191610215876101ca8685611fc2565b515116938460005260018a528d6040600020906000528a526040600020546101f28789611fc2565b526101fd8688611fc2565b51906040519161020c83611841565b600083526128cc565b9094911561027c575b91604061023f8880989796948e966102396102779c8f611fc2565b52611fc2565b510152015160405192610251846117bf565b83528883015260408201526102668289611fc2565b526102718188611fc2565b50611f6a565b6101a1565b60009c5061021e565b50919095879460005b825181101561034f5780877fc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a8b6102d684896102cd61034a988b611fc2565b51511692611fc2565b518160005260018b526040600020846000528b526102fa6040600020918254611fd6565b9055610306848a611fc2565b519060005260018a526040600020836000528a526040600020549061034260405192839287846040919493926060820195825260208201520152565b0390a2611f6a565b61028e565b5085858389156103955750600091825252600060408120555b60005b81518110156100d0578061038b6103856103909385611fc2565b51612f7a565b611f6a565b61036b565b906103de6103a460009361340b565b6103ad85613755565b509190604051926103bd846117f0565b65ffffffffffff8092168452168483015284604083015260608201526136e4565b92825252604060002055610368565b6040516103f9816117bf565b600081526040516104098161180b565b60008152606090818582015284830152604082015282828801015201839061017f565b60405162461bcd60e51b815290810183905260156024820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b6044820152606490fd5b346100d257600319602036820181136100d2578235916001600160401b03938484116100d257610100809285360301126100d2576040519182018281108682111761072e57604052838101358252828201906024850135825260448501358681116100d2576104dc908236918801016118bd565b906040840191825260608401956064810135875260808501906084810135825260a086019860a48201358a5260c48201359081116100d25761052560e4918536918501016118bd565b918260c089015201359860e08701938a85526105608851965161055b8c519861054d81613182565b8a518351848f012090613119565b613442565b809661056b85613442565b9c6001600160a01b0391600260ff60406105a281610597868961058e828d611fc2565b5151169a611fc2565b5101518c5190611fc2565b51015116036106ea578e9695949392916105cf60406105c46105d7948e611fc2565b5101518a5190611fc2565b515197611fc2565b515116036106a657506106478a9b604061063c7f4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b9d610635888f6106529b9961069a9f9e9d9b99610628889b613182565b5191815191012090613119565b5189611fc2565b510151945190611fc2565b510151905191612c77565b9184519260406106638a5185611fc2565b510152519060405161069181610683898201948a86526040830190612b56565b03601f19810183528261185c565b519020916131db565b519251604051908152a2005b60405162461bcd60e51b81529081018a9052601d60248201527f746172676574417373657420213d2067756172616e74656541737365740000006044820152606490fd5b60405162461bcd60e51b81528086018e9052601a60248201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e0000000000006044820152606490fd5b604182634e487b7160e01b6000525260246000fd5b346100d25760203660031901126100d2573560005260006020526020604060002054604051908152f35b346100d25760c03660031901126100d2576001600160401b039080358281116100d25761079d9036908301611c61565b906024358381116100d2576107b59036908301611ed7565b926044359081116100d2576107cd9036908301611d28565b60603660631901126100d2576040516107e5816117bf565b60643560ff811681036100d2578152608435602082015260a435604082015261080d846132f1565b9365ffffffffffff6040845101511693610826866136b1565b61082f8161224c565b610a935765ffffffffffff61084387613755565b5050168510610a4f575b61086161085b85898561253c565b9061226c565b6108ba61086f855184613394565b936108b584519160405160208101908882526040808201526009606082015268666f7263654d6f766560b81b6080820152608081526108ad81611826565b51902061322c565b6124f5565b15610a0b57506108de65ffffffffffff60608301511665ffffffffffff421661229c565b606084510151151565ffffffffffff60405192878452166020830152604082015260c0606082015261091360c08201836122b6565b8181036080830152875180825260208201906020808260051b8501019a01916000905b8282106109dc57505050505061099b65ffffffffffff60608194897f0e6d8300485cb09fa95f22b89b46f7b0cc3029f1bbf257a0884414d415546cf8868061098d6109c89e9f6109a39a810360a08401528d612334565b0390a201511683421661229c565b93515161340b565b92604051946109b1866117f0565b8552166020840152604083015260608201526136e4565b906000526000602052604060002055600080f35b909192939a6020806109fb6001938f601f199082030186528851612334565b9d96019493919091019101610936565b60649060206040519162461bcd60e51b8352820152601f60248201527f4368616c6c656e676572206973206e6f742061207061727469636970616e74006044820152fd5b60649060206040519162461bcd60e51b8352820152601860248201527f7475726e4e756d5265636f7264206465637265617365642e00000000000000006044820152fd5b6001610a9e876136b1565b610aa78161224c565b03610abb57610ab6858761279a565b61084d565b610ab6866127fe565b346100d25760603660031901126100d2576001600160401b0381358181116100d257610af39036908401611c61565b916024358281116100d257610b0b9036908301611ed7565b926044359283116100d257610b7d61085b610b4c6020947f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e096369101611d28565b610b55846132f1565b9665ffffffffffff6040835101511694610b6e896127fe565b610b78868a61279a565b61253c565b610ba8604051610b8c816117f0565b82815260008482015260006040820152600060608201526136e4565b8460005260008352604060002055604051908152a2005b346100d2576003196060368201126100d2578135906001600160401b03928383116100d257828101833603946080848701126100d257602435948186116100d257366023870112156100d25785840135958287116100d2573660248860051b830101116100d257604435928084116100d25760408785360301126100d257604483019081359860018060a01b0393848b16809b036100d257610c9b91610c77610c7d92610c6c368c611c61565b926024369201611e68565b90612613565b95610c95610c8b368a611c61565b9136908b01611d28565b906126a4565b966040519a634c9b6c0960e11b8c526060828d015260e48c01973590602219018112156100d25785016024810197910135908282116100d2578160051b360388136100d257608060648d01528190528a9897966101048a0196959493929160005b818110610ddc5750505094610d719465ffffffffffff610d4d6064878d9b9760009f9d99610d62988e6084819f610d386024610d419601611c3a565b16910152611b14565b1660a48c015201611c4e565b1660c4880152848783030160248801526121f7565b918483030160448501526121d7565b03915afa908115610dd057600090600092610daa575b50610da660405192839215158352604060208401526040830190611a54565b0390f35b9050610dc991503d806000833e610dc1818361185c565b810190612027565b9082610d87565b6040513d6000823e3d90fd5b9198999a5091929394959660019086610df48b611b14565b168152602080910199019101918c9a9998979695949392610cfc565b346100d25760603660031901126100d2576001600160401b039080358281116100d257610e409036908301611904565b6024359283116100d257610e5d610e6692610da694369101611904565b60443591612c77565b604051918291602083526020830190611a79565b346100d25760203660031901126100d257610e9760609135613755565b6040805165ffffffffffff94851681529390921660208401526001600160a01b031690820152f35b346100d2576060806003193601126100d2576024356001600160401b0381116100d257610eef9036908401611b28565b90610efa8335613182565b610f11610f068361340b565b843590604435613119565b600191805191610f208361187d565b92610f2e604051948561185c565b808452610f3d601f199161187d565b019060005b82811061115b57505050610f568151611f38565b610f608251611f38565b9160005b815181101561103557610f778183611fc2565b51604081015190610fc86001600160a01b03610f938587611fc2565b5151169283600052600160205260406000208b35600052602052604060002054610fbd8689611fc2565b526101fd8588611fc2565b909392911561102c575b9160209186959493610fe7611027988c611fc2565b526040610ff4878a611fc2565b510152015160405192611006846117bf565b83526020830152604082015261101c8288611fc2565b526102718187611fc2565b610f64565b60009a50610fd2565b5090919260005b8251811015611101576110fc906001600160a01b0361105b8286611fc2565b5151166110688288611fc2565b51816000526020906001825260406000208b3560005282526110906040600020918254611fd6565b905561109c8388611fc2565b5191600052600181526040600020908a35600052527fc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a6103426040600020546040519182918d359587846040919493926060820195825260208201520152565b61103c565b508582861561113e5750356000526000602052600060408120555b60005b81518110156100d0578061038b6103856111399385611fc2565b61111f565b9061114b6111569261340b565b9060443590356131db565b61111c565b60209060405161116a816117bf565b6000815260405161117a8161180b565b600081528390858282015281830152846040830152828801015201610f42565b346100d25760a03660031901126100d2576001600160401b039080359060446024803582358681116100d2576111d390369086016118bd565b946064968735906084359081116100d2576111f190369088016119d3565b9260005b600181018082116113c057855181101561126d5761121e6112168388611fc2565b519187611fc2565b5111156112335761122e90611f6a565b6111f5565b60405162461bcd60e51b81526020818a015260168188015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b818901528a90fd5b50508761139b6100d0959461128184613182565b8361129a84519461055b6020978897888401208a613119565b6001600160a01b03979091907fc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a90896112d38686611fc2565b5151168060005260018852604060002084600052885261136961130b6040600020549360406113028a8a611fc2565b510151856128cc565b909d9291508460005260018c526040600020886000528c526113336040600020918254611fd6565b905560406113418a8a611fc2565b5101526040518a8101908b825261135f81610683604082018c612b56565b51902090866131db565b6000908152600188526040808220858352895290819020548151878152602081019390935290820152606090a2611fc2565b519384511693015190604051936113b1856117bf565b84528301526040820152612f7a565b8660118a634e487b7160e01b600052526000fd5b60803660031901126100d2576113e8611afe565b60249081359160649384358460a01c156116a45760018060a01b038416918260005260209360018552604060002087600052855260449081356040600020540361166c57846114f7578334036114b6575050507f87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b9495505b816000526001835260406000208660005283526114836040600020918254612853565b905560009081526001825260408082208683529092528190205481516001600160a01b03939093168352602083015290a2005b60405162461bcd60e51b8152928301869052601f908301527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400908201528690fd5b604098929398518681016323b872dd60e01b8152338383015230848301528a8583015284825261152682611826565b604051916115338361180b565b8883527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656489840152873b1561162b57600061157f939281925190828b5af16115796130e9565b90613780565b8051806115b5575b5050505050507f87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b9394611460565b818891810103126100d257866115cb9101611fe3565b156115d7578080611587565b60405162461bcd60e51b8152938401869052602a908401527f5361666545524332303a204552433230206f7065726174696f6e20646964206e90830152691bdd081cdd58d8d9595960b21b90820152608490fd5b60405162461bcd60e51b81528088018a9052601d818601527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000818701528690fd5b60405162461bcd60e51b8152928301869052601490830152731a195b1908084f48195e1c1958dd195912195b1960621b908201528690fd5b50601f859260206040519362461bcd60e51b85528401528201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e006044820152fd5b346100d25760403660031901126100d2576001600160a01b03611707611afe565b16600052600160205260406000206024356000526020526020604060002054604051908152f35b346100d25760603660031901126100d2576001600160401b036024358181116100d25761175e9036908401611904565b906044359081116100d2576117a0926117b59261178161178893369084016119d3565b91356128cc565b92939190604051958695608087526080870190611a79565b91151560208601528482036040860152611a79565b9060608301520390f35b606081019081106001600160401b038211176117da57604052565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b038211176117da57604052565b604081019081106001600160401b038211176117da57604052565b60a081019081106001600160401b038211176117da57604052565b602081019081106001600160401b038211176117da57604052565b90601f801991011681019081106001600160401b038211176117da57604052565b6001600160401b0381116117da5760051b60200190565b359060ff821682036100d257565b6001600160401b0381116117da57601f01601f191660200190565b81601f820112156100d2578035906118d4826118a2565b926118e2604051948561185c565b828452602083830101116100d257816000926020809301838601378301015290565b9080601f830112156100d257813561191b8161187d565b9260409161192b8351958661185c565b808552602093848087019260051b840101938185116100d257858401925b85841061195a575050505050505090565b6001600160401b0384358181116100d257860191608080601f1985880301126100d257845190611989826117f0565b8a8501358252858501358b8301526060906119a5828701611894565b878401528501359384116100d2576119c4878c809796819701016118bd565b90820152815201930192611949565b81601f820112156100d2578035916119ea8361187d565b926119f8604051948561185c565b808452602092838086019260051b8201019283116100d2578301905b828210611a22575050505090565b81358152908301908301611a14565b60005b838110611a445750506000910152565b8181015183820152602001611a34565b90602091611a6d81518092818552858086019101611a31565b601f01601f1916010190565b908082519081815260208091019281808460051b8301019501936000915b848310611aa75750505050505090565b9091929394958480611aee600193601f198682030187528a51805182528381015184830152604060ff81830151169083015260608091015191608080928201520190611a54565b9801930193019194939290611a97565b600435906001600160a01b03821682036100d257565b35906001600160a01b03821682036100d257565b9080601f830112156100d2578135611b3f8161187d565b92604091611b4f8351958661185c565b808552602093848087019260051b840101938185116100d257858401925b858410611b7e575050505050505090565b6001600160401b0384358181116100d257860191601f196060848703820181136100d257855191611bae836117bf565b611bb98c8701611b14565b8352868601358581116100d2578790870191828a0301126100d257865190611be08261180b565b8c81013560048110156100d2578252878101358681116100d2578d8a91611c089301016118bd565b8c8201528b8301528401359283116100d257611c2b868b80969581960101611904565b85820152815201930192611b6d565b35906001600160401b03821682036100d257565b359065ffffffffffff821682036100d257565b9190916080818403126100d25760405190611c7b826117f0565b819381356001600160401b0381116100d25782019080601f830112156100d257813590611ca78261187d565b91611cb5604051938461185c565b808352602093848085019260051b8201019283116100d2578401905b828210611d1157505050606092611d0c9284928652611cf1818301611c3a565b90860152611d0160408201611b14565b604086015201611c4e565b910152565b848091611d1d84611b14565b815201910190611cd1565b919060409283818303126100d2578351848101916001600160401b0395828410878511176117da57838152829682358181116100d2578301926080848803126100d257611d74866117f0565b83358281116100d25787611d89918601611b28565b8652602095868501358381116100d25788611da59187016118bd565b606095869182890152611db9868201611c4e565b6080890152013580151581036100d25760a08701528552858101359182116100d257019085601f830112156100d257813590611df48261187d565b96611e018251988961185c565b8288528685818a019402850101938185116100d2578701925b848410611e2b575050505050500152565b85848303126100d2578786918451611e42816117bf565b611e4b87611894565b815282870135838201528587013586820152815201930192611e1a565b92919092611e758461187d565b91611e83604051938461185c565b829480845260208094019060051b8301928284116100d25780915b848310611ead57505050505050565b82356001600160401b0381116100d2578691611ecc8684938601611d28565b815201920191611e9e565b9080601f830112156100d257816020611ef293359101611e68565b90565b9060406003198301126100d2576001600160401b036004358181116100d25783611f2191600401611c61565b926024359182116100d257611ef291600401611d28565b90611f428261187d565b611f4f604051918261185c565b8281528092611f60601f199161187d565b0190602036910137565b6000198114611f795760010190565b634e487b7160e01b600052601160045260246000fd5b805115611f9c5760200190565b634e487b7160e01b600052603260045260246000fd5b805160011015611f9c5760400190565b8051821015611f9c5760209160051b010190565b91908203918211611f7957565b519081151582036100d257565b90929192611ffd816118a2565b9161200b604051938461185c565b8294828452828201116100d2576020612025930190611a31565b565b91906040838203126100d25761203c83611fe3565b926020810151906001600160401b0382116100d257019080601f830112156100d2578151611ef292602001611ff0565b906080810182519060808352815180915260a0830190602060a08260051b8601019301916000905b8282106120d5575050505060606120b8819260208601518582036020870152611a54565b9365ffffffffffff60408201511660408501520151151591015290565b90919293609f19868203018252845160018060a01b03815116825260208101516060602084015280516004808210156121c257506020604092612129926060870152015182608086015260a0850190611a54565b910151916040818303910152815180825260208201906020808260051b85010194019260005b82811061216f575050505050602080600192960192019201909291612094565b90919293946020806121b5600193601f19878203018952608060608b5180518452858101518685015260ff60408201511660408501520151918160608201520190611a54565b970195019392910161214f565b602190634e487b7160e01b6000525260246000fd5b906020806121ee845160408552604085019061206c565b93015191015290565b90815180825260208092019182818360051b85019501936000915b8483106122225750505050505090565b909192939495848061223c83856001950387528a516121d7565b9801930193019194939290612212565b6003111561225657565b634e487b7160e01b600052602160045260246000fd5b156122745750565b60405162461bcd60e51b815260206004820152908190612298906024830190611a54565b0390fd5b91909165ffffffffffff80809416911601918211611f7957565b906080810191805160808352805180945260a083019360208092019060005b81811061231757505050808201516001600160401b0316908301526040808201516001600160a01b03169083015260609081015165ffffffffffff1691015290565b82516001600160a01b0316875295830195918301916001016122d5565b8051906123496040928385528385019061206c565b9060208091015193818184039101528080855193848152019401926000905b83821061237757505050505090565b8451805160ff16875280840151878501528101518682015260609095019493820193600190910190612368565b91906123af836132f1565b926123b9846127fe565b606082510151156124ba5760206123d083836126a4565b015190816000925b612495575060ff905151911603612463577f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc9016020849261244c61242565ffffffffffff421692515161340b565b60405190612432826117f0565b6000825283858301526000604083015260608201526136e4565b8460005260008352604060002055604051908152a2565b60405162461bcd60e51b815260206004820152600a60248201526921756e616e696d6f757360b01b6044820152606490fd5b600019810190808211611f7957169160ff809116908114611f795760010191806123d8565b60405162461bcd60e51b815260206004820152601360248201527214dd185d19481b5d5cdd08189948199a5b985b606a1b6044820152606490fd5b60005b8251811015612534576001600160a01b03806125148386611fc2565b51169083161461252c5761252790611f6a565b6124f8565b505050600190565b505050600090565b600091612593610d62929594956125a661256d61256660018060a01b036040850151169584612613565b98836126a4565b9760405198899687958695634c9b6c0960e11b87526060600488015260648701906122b6565b60031993848783030160248801526121f7565b03915afa918215610dd05760009081936125bf57509190565b906125d59293503d8091833e610dc1818361185c565b9091565b604051906125e68261180b565b60006020836040516125f7816117f0565b6060815260608382015283604082015283606082015281520152565b81519161261f8361187d565b9261262d604051948561185c565b80845261263c601f199161187d565b0160005b81811061268d57505060005b8151811015612687578061266d6126666126829385611fc2565b51856126a4565b6126778287611fc2565b526102718186611fc2565b61264c565b50505090565b6020906126986125d9565b82828801015201612640565b91906126ae6125d9565b50805190604051916126bf8361180b565b82526020928383019260009283855283955b80820180515188101561278e57612701906126fa896126f3869896518d613394565b9251611fc2565b519061322c565b6001600160a01b0390811694909390865b8a51805182101561277e57612728828892611fc2565b5116871461273e5761273990611f6a565b612712565b929891955093509060ff811161276a57906001612760921b8751178752611f6a565b95929190926126d1565b634e487b7160e01b86526011600452602486fd5b5050935093509561276090611f6a565b50505093509350505090565b6127a390613755565b505065ffffffffffff809116911611156127b957565b60405162461bcd60e51b815260206004820152601c60248201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e000000006044820152606490fd5b6128096002916136b1565b6128128161224c565b1461281957565b60405162461bcd60e51b815260206004820152601260248201527121b430b73732b6103334b730b634bd32b21760711b6044820152606490fd5b91908201809211611f7957565b9061286a8261187d565b6040906128798251918261185c565b838152809361288a601f199161187d565b0191600091825b84811061289f575050505050565b60209083516128ad816117f0565b8581528286818301528686830152606080830152828501015201612891565b9192908351801515600014612b4b576128e490612860565b916000916128f28151612860565b95600190818097938960009586935b61290f575b50505050505050565b909192939495978351851015612b42576129298585611fc2565b51516129358685611fc2565b515260409060ff80836129488989611fc2565b51015116836129578988611fc2565b5101526060806129678989611fc2565b510151816129758a89611fc2565b510152602093846129868a8a611fc2565b51015186811115612b3c575085965b8d8b51908b8215928315612b12575b505050600014612ae15750600283828f6129be908c611fc2565b5101511614612a9e578f96959493868f918f612a5b90612a6194612a6d988f988f908f91612a679a898f94612a368f869288612a1183612a0b8884612a03848e611fc2565b510151611fd6565b93611fc2565b510152612a1e8187611fc2565b51519885612a2c8389611fc2565b5101511695611fc2565b51015194825196612a46886117f0565b87528601528401528201526102398383611fc2565b50612853565b9c611f6a565b95611fc2565b510151612a95575b612a8891612a8291611fd6565b93611f6a565b91909493928a9085612901565b60009a50612a75565b5162461bcd60e51b815260048101859052601b60248201527f63616e6e6f74207472616e7366657220612067756172616e74656500000000006044820152606490fd5b9050612a6d925088915084612afc83959e989796958a611fc2565b51015184612b0a8484611fc2565b510152611fc2565b821092509082612b27575b50508e8b386129a4565b612b339192508d611fc2565b51148a8f612b1d565b96612995565b97829150612906565b506128e48151612860565b90815180825260208092019182818360051b82019501936000915b848310612b815750505050505090565b9091929394958181038352865160018060a01b0381511682528581015190606091828885015280516004808210156121c257509188612bd99285948796839801520151604092839182608088015260a0870190611a54565b91015193828183039101528351908181528581019286808460051b8401019601946000915b848310612c21575050505050505090806001929801930193019194939290612b71565b9193959780612c63600193959799601f19898203018b526080878d5180518452858101518685015260ff89820151168985015201519181898201520190611a54565b99019701930190918b979695939492612bfe565b80516000198101908111611f7957612c8e90612860565b91612c998483611fc2565b51606081015192604094855191612caf8361180b565b600095868452866020809501528781805181010312612f765787805191612cd58361180b565b85810151835201519084810191825287998890899c8a988b5b87518d1015612e57578f848e14612e48578c8f8f90612d5b858f8f908f612d158782611fc2565b51519582612d238984611fc2565b5101516060612d398a60ff85612a2c8389611fc2565b51015193825198612d498a6117f0565b89528801528601526060850152611fc2565b52612d66848d611fc2565b5087159081612e32575b50612df8575b501580612de3575b612d95575b612a61612d8f91611f6a565b9b612cee565b9e5098612dd8908f612dc38b612db98f612daf8391611fb2565b510151938d611fc2565b5101918251612853565b905289612dcf8d611fb2565b51015190612853565b60019e909990612d83565b50612dee8d89611fc2565b5151875114612d7e565b829c919650612dcf818c612e218f612db9612e289882612e188199611f8f565b51015194611fc2565b9052611f8f565b996001948c612d76565b612e3d91508b611fc2565b51518851148f612d70565b509b9d50612d8f60019e611f6a565b509899509c969a99505093999250505015612f3a5715612f005715612ec55783015103612e8357505090565b60649250519062461bcd60e51b825280600483015260248201527f746f74616c5265636c61696d6564213d67756172616e7465652e616d6f756e746044820152fd5b825162461bcd60e51b815260048101859052601460248201527318dbdd5b19081b9bdd08199a5b99081c9a59da1d60621b6044820152606490fd5b835162461bcd60e51b815260048101869052601360248201527218dbdd5b19081b9bdd08199a5b99081b19599d606a1b6044820152606490fd5b845162461bcd60e51b815260048101879052601560248201527418dbdd5b19081b9bdd08199a5b99081d185c99d95d605a1b6044820152606490fd5b8680fd5b80516001600160a01b03908116919060005b60408084019081519182518410156129065784612faa858095611fc2565b515191612fbb602095869251611fc2565b510151918060a01c156000146130be57168761303457600080809381935af1612fe26130e9565b5015612ff8575050612ff390611f6a565b612f8c565b60649250519062461bcd60e51b825260048201526016602482015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b6044820152fd5b825163a9059cbb60e01b81526001600160a01b039190911660048201526024810191909152929190818460448160008b5af19081156130b4575061307e575b50612ff39150611f6a565b82813d83116130ad575b613092818361185c565b810103126100d2576130a6612ff392611fe3565b5038613073565b503d613088565b513d6000823e3d90fd5b60008981526001865284812091815294525091208054612ff393926130e291612853565b9055611f6a565b3d15613114573d906130fa826118a2565b91613108604051938461185c565b82523d6000602084013e565b606090565b9161312390613755565b936001600160a01b039384935061313b925090613728565b1691160361314557565b60405162461bcd60e51b81526020600482015260156024820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b6044820152606490fd5b61318d6002916136b1565b6131968161224c565b0361319d57565b60405162461bcd60e51b815260206004820152601660248201527521b430b73732b6103737ba103334b730b634bd32b21760511b6044820152606490fd5b919061321b916131ea84613755565b509290604051936131fa856117f0565b65ffffffffffff8092168552166020840152604083015260608201526136e4565b906000526000602052604060002055565b90600060806020926040948551858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c8152613272816117bf565b5190209060ff8151169086868201519101519187519384528684015286830152606082015282805260015afa156130b457600051906001600160a01b038216156132ba575090565b5162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b6044820152606490fd5b8051906020916001600160401b03838301511660018060a01b039165ffffffffffff606084604087015116950151166040519485938785019760a086019060808a5285518092528060c088019601976000905b838210613377575050505060408501526060840152608083015203601f198101835261337191508261185c565b51902090565b895181168852988201988a98509682019660019190910190613344565b6133716133a36133e1926132f1565b926020810151815191606065ffffffffffff60408301511691015115156133f460405196879460208601998a5260a0604087015260c0860190611a54565b601f199586868303016060870152612b56565b91608084015260a08301520390810183528261185c565b604051613371816106836020820194602086526040830190612b56565b9080601f830112156100d2578151611ef292602001611ff0565b805181016020828203126100d25760208201516001600160401b0381116100d25760208201603f8285010112156100d25760208184010151906134848261187d565b93613492604051958661185c565b82855260208501916020850160408560051b8385010101116100d257604081830101925b60408560051b838501010184106134d05750505050505090565b83516001600160401b0381116100d257601f19908484010160608189038301126100d25760405191613501836117bf565b60408201516001600160a01b03811681036100d257835260608201516001600160401b0381116100d257604090830191828b0301126100d257604051906135478261180b565b604081015160048110156100d25782526060810151906001600160401b0382116100d257604061357d9260208d01920101613428565b6020820152602083015260808101516001600160401b0381116100d25760208901605f8284010112156100d25760408183010151906135bb8261187d565b926135c9604051948561185c565b828452602084019060208c0160608560051b8584010101116100d257606083820101915b60608560051b8584010101831061361657505050505060408201528152602093840193016134b6565b82516001600160401b0381116100d257608083860182018f03603f1901126100d25760405191613645836117f0565b8386018201606081015184526080810151602085015260a0015160ff811681036100d257604084015260c082878601010151926001600160401b0384116100d2578f6020949360608695866136a19401928b8a01010101613428565b60608201528152019201916135ed565b6136c165ffffffffffff91613755565b5090501680156000146136d45750600090565b42106136df57600290565b600190565b65ffffffffffff60d01b815160d01b1665ffffffffffff60a01b602083015160a01b1617906137236040820151606060018060a01b0393015190613728565b161790565b6040519160208301918252604083015260408252613745826117bf565b905190206001600160a01b031690565b60005260006020526040600020548060d01c9165ffffffffffff8260a01c169160018060a01b031690565b9091901561378c575090565b8151156122745750805190602001fdfea264697066735822122060af7968d12c4e7bf679e20a79d3d88eb209f037ccedf8026a107c3f7d941b7064736f6c63430008110033",
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"}],\"name\":\"Reclaimed\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"}],\"name\":\"compute_reclaim_effects\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32[]\",\"name\":\"channelIds\",\"type\":\"bytes32[]\"},{\"internalType\":\"uint256[]\",\"name\":\"expectedHelds\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256[]\",\"name\":\"amounts\",\"type\":\"uint256[]\"}],\"name\":\"deposit_batch\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"}],\"internalType\":\"structIMultiAssetHolder.ReclaimArgs\",\"name\":\"reclaimArgs\",\"type\":\"tuple\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint64\",\"name\":\"channelNonce\",\"type\":\"uint64\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"proof\",\"type\":\"tuple[]\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart\",\"name\":\"candidate\",\"type\":\"tuple\"}],\"name\":\"stateIsSupported\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"},{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"components\":[{\"internalType\":\"enumExitFormat.AssetType\",\"name\":\"assetType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.AssetMetadata\",\"name\":\"assetMetadata\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608080604052346100165761394d908161001c8239f35b600080fdfe608060405260048036101561001357600080fd5b60003560e01c806311e9f17814611868578063166e56cd146118205780632fb1d2701461150d5780633033730e1461130057806331afa0b41461104d578063552cfa5014611007578063566d54c614610f9c5780635685b7dc14610d655780636d2a9c9214610c695780638286a06014610911578063c7df14e2146108e6578063d3c4e7381461060a578063ec346235146102ab578063ee049b50146102925763f7cb41ba146100c257600080fd5b608036600319011261028d576100d6611c39565b6001600160401b039160243583811161028d576100f69036908301611d75565b9360443581811161028d5761010e9036908501611d75565b91606496873582811161028d576101289036908801611d75565b94909181811480610284575b156102415794976001600160a01b03169460005b81811061015157005b87868c8c6101858561017d818a610176828c8f828a61016f926129cb565b359b6129cb565b35976129cb565b35928a6129cb565b35906040519260208401946302fb1d2760e41b86528d6024860152604485015283015260849081830152815260c081018181108982111761022c5760405251600092839290305af16101d56129db565b50156101e9576101e4906120d5565b610148565b60405162461bcd60e51b81526020818b0152601d60248201527f6465706f7369745f62617463683a206465706f736974206661696c656400000060448201528b90fd5b60418d634e487b7160e01b6000525260246000fd5b60405162461bcd60e51b81526020818a0152601860248201527f4172726179206c656e67746873206d757374206d61746368000000000000000060448201528a90fd5b50858214610134565b600080fd5b3461028d576102a96102a336612060565b9061250f565b005b503461028d576102c6906102be36612060565b80939161250f565b915151906102d3836132fd565b6102dc82613586565b6102e5846138d0565b9391505060405192602093848101936000855260408201526040815261030a816118fa565b519092206001600160a01b0392908316908316036105ce5750600191835191610332836119b8565b926103406040519485611997565b80845261034f601f19916119b8565b018260005b82811061058f5750505061036885516120a3565b61037286516120a3565b9060005b875181101561045a57610389818961212d565b51886040820151916103ea8761039f868561212d565b515116938460005260018a528d6040600020906000528a526040600020546103c7878961212d565b526103d2868861212d565b5190604051916103e18361197c565b60008352612a77565b90949115610451575b9160406104148880989796948e9661040e61044c9c8f61212d565b5261212d565b510152015160405192610426846118fa565b835288830152604082015261043b828961212d565b52610446818861212d565b506120d5565b610376565b60009c506103f3565b509186889560005b82518110156104f157808461047a6104ec938661212d565b515116610487828861212d565b51906000526001895260406000208a60005289526104ab6040600020918254612141565b9055887fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c13063767960406104db848b61212d565b518151908582528c820152a26120d5565b610462565b508686838a156105375750600091825252600060408120555b60005b81518110156102a9578061052d610527610532938561212d565b51613125565b6120d5565b61050d565b90610580610546600093613586565b61054f856138d0565b5091906040519261055f8461192b565b65ffffffffffff80921684521684830152846040830152606082015261385f565b9282525260406000205561050a565b60405161059b816118fa565b600081526040516105ab81611946565b600081526060908185820152848301526040820152828288010152018390610354565b60405162461bcd60e51b815290810183905260156024820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b6044820152606490fd5b503461028d576003196020368201811361028d578235916001600160401b039384841161028d576101008092853603011261028d57604051918201828110868211176108d1576040528381013582528282019060248501358252604485013586811161028d5761067f908236918801016119f8565b906040840191825260608401956064810135875260808501906084810135825260a086019860a48201358a5260c482013590811161028d576106c860e4918536918501016119f8565b918260c089015201359860e08701938a8552610703885196516106fe8c51986106f0816132fd565b8a518351848f012090613294565b6135bd565b809661070e856135bd565b9c6001600160a01b0391600260ff60406107458161073a8689610731828d61212d565b5151169a61212d565b5101518c519061212d565b510151160361088d578e969594939291610772604061076761077a948e61212d565b5101518a519061212d565b51519761212d565b5151160361084957506107ea8a9b60406107df7f4d3754632451ebba9812a9305e7bca17b67a17186a5cff93d2e9ae1b01e3d27b9d6107d8888f6107f59b9961083d9f9e9d9b996107cb889b6132fd565b5191815191012090613294565b518961212d565b51015194519061212d565b510151905191612e22565b9184519260406108068a518561212d565b510152519060405161083481610826898201948a86526040830190612d01565b03601f198101835282611997565b51902091613356565b519251604051908152a2005b60405162461bcd60e51b81529081018a9052601d60248201527f746172676574417373657420213d2067756172616e74656541737365740000006044820152606490fd5b60405162461bcd60e51b81528086018e9052601a60248201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e0000000000006044820152606490fd5b604182634e487b7160e01b6000525260246000fd5b503461028d57602036600319011261028d573560005260006020526020604060002054604051908152f35b503461028d5760c036600319011261028d576001600160401b0390803582811161028d576109429036908301611dcc565b9060243583811161028d5761095a9036908301612042565b9260443590811161028d576109729036908301611e93565b606036606319011261028d5760405161098a816118fa565b60643560ff8116810361028d578152608435602082015260a43560408201526109b28461346c565b9365ffffffffffff60408451015116936109cb8661382c565b6109d4816123b7565b610c385765ffffffffffff6109e8876138d0565b5050168510610bf4575b610a06610a008589856126a7565b906123d7565b610a5f610a1485518461350f565b93610a5a84519160405160208101908882526040808201526009606082015268666f7263654d6f766560b81b608082015260808152610a5281611961565b5190206133a7565b612660565b15610bb05750610a8365ffffffffffff60608301511665ffffffffffff4216612407565b606084510151151565ffffffffffff60405192878452166020830152604082015260c06060820152610ab860c0820183612421565b8181036080830152875180825260208201906020808260051b8501019a01916000905b828210610b81575050505050610b4065ffffffffffff60608194897f0e6d8300485cb09fa95f22b89b46f7b0cc3029f1bbf257a0884414d415546cf88680610b32610b6d9e9f610b489a810360a08401528d61249f565b0390a2015116834216612407565b935151613586565b9260405194610b568661192b565b85521660208401526040830152606082015261385f565b906000526000602052604060002055600080f35b909192939a602080610ba06001938f601f19908203018652885161249f565b9d96019493919091019101610adb565b60649060206040519162461bcd60e51b8352820152601f60248201527f4368616c6c656e676572206973206e6f742061207061727469636970616e74006044820152fd5b60649060206040519162461bcd60e51b8352820152601860248201527f7475726e4e756d5265636f7264206465637265617365642e00000000000000006044820152fd5b6001610c438761382c565b610c4c816123b7565b03610c6057610c5b8587612905565b6109f2565b610c5b86612969565b503461028d57606036600319011261028d576001600160401b03813581811161028d57610c999036908401611dcc565b9160243582811161028d57610cb19036908301612042565b9260443592831161028d57610d23610a00610cf26020947f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e096369101611e93565b610cfb8461346c565b9665ffffffffffff6040835101511694610d1489612969565b610d1e868a612905565b6126a7565b610d4e604051610d328161192b565b828152600084820152600060408201526000606082015261385f565b8460005260008352604060002055604051908152a2005b503461028d5760031960603682011261028d578135906001600160401b039283831161028d578281018336039460808487011261028d5760243581811161028d57610db39036908501611d75565b95906044359280841161028d57604087853603011261028d57604483019081359860018060a01b0393848b16809b0361028d57610e2791610e03610e0992610dfb368c611dcc565b923691611fd3565b9061277e565b95610e21610e17368a611dcc565b9136908b01611e93565b9061280f565b966040519a634c9b6c0960e11b8c526060828d015260e48c019735906022190181121561028d57850160248101979101359082821161028d578160051b3603881361028d57608060648d01528190528a9897966101048a0196959493929160005b818110610f685750505094610efd9465ffffffffffff610ed96064878d9b9760009f9d99610eee988e6084819f610ec46024610ecd9601611da5565b16910152611c4f565b1660a48c015201611db9565b1660c488015284878303016024880152612362565b91848303016044850152612342565b03915afa908115610f5c57600090600092610f36575b50610f3260405192839215158352604060208401526040830190611b8f565b0390f35b9050610f5591503d806000833e610f4d8183611997565b810190612192565b9038610f13565b6040513d6000823e3d90fd5b9198999a5091929394959660019086610f808b611c4f565b168152602080910199019101918c9a9998979695949392610e88565b503461028d57606036600319011261028d576001600160401b0390803582811161028d57610fcd9036908301611a3f565b60243592831161028d57610fea610ff392610f3294369101611a3f565b60443591612e22565b604051918291602083526020830190611bb4565b503461028d57602036600319011261028d57611025606091356138d0565b6040805165ffffffffffff94851681529390921660208401526001600160a01b031690820152f35b503461028d5760608060031936011261028d576024356001600160401b03811161028d5761107e9036908401611c63565b9061108983356132fd565b6110a061109583613586565b843590604435613294565b6001918051916110af836119b8565b926110bd6040519485611997565b8084526110cc601f19916119b8565b019060005b8281106112c1575050506110e581516120a3565b6110ef82516120a3565b9160005b81518110156111c457611106818361212d565b516040810151906111576001600160a01b03611122858761212d565b5151169283600052600160205260406000208b3560005260205260406000205461114c868961212d565b526103d2858861212d565b90939291156111bb575b91602091869594936111766111b6988c61212d565b526040611183878a61212d565b510152015160405192611195846118fa565b8352602083015260408201526111ab828861212d565b52610446818761212d565b6110f3565b60009a50611161565b5090919260005b825181101561126757611262906001600160a01b036111ea828661212d565b5151166111f7828861212d565b51906000526020906001825260406000208a35600052825261121f6040600020918254612141565b905561122b828761212d565b51604051918383528201527fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c13063767960408a3592a26120d5565b6111cb565b50858286156112a45750356000526000602052600060408120555b60005b81518110156102a9578061052d61052761129f938561212d565b611285565b906112b16112bc92613586565b906044359035613356565b611282565b6020906040516112d0816118fa565b600081526040516112e081611946565b6000815283908582820152818301528460408301528288010152016110d1565b503461028d5760a036600319011261028d576001600160401b0390803590604460248035823586811161028d5761133a90369086016119f8565b9460649687359060843590811161028d576113589036908801611b0e565b9260005b600181018082116114f95785518110156113d45761138561137d838861212d565b51918761212d565b51111561139a57611395906120d5565b61135c565b60405162461bcd60e51b81526020818a015260168188015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b818901528a90fd5b5050876114d46102a995946113e8846132fd565b836114018451946106fe6020978897888401208a613294565b907fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c130637679604060018060a01b03996114c48b61143b888861212d565b5151169a8b60005260018a5283600020866000528a5261146e846000205493856114658b8b61212d565b51015185612a77565b91509d60005260018c5285600020886000528c5261149186600020918254612141565b90558461149e8a8a61212d565b51015283518a8101908b82526114ba816108268882018c612d01565b5190209086613356565b81519086825288820152a261212d565b519384511693015190604051936114ea856118fa565b84528301526040820152613125565b8660118a634e487b7160e01b600052526000fd5b50608036600319011261028d57611522611c39565b60249081359160649384358460a01c156117de5760018060a01b03841691826000526020936001855260406000208760005285526044908135604060002054036117a65784611631578334036115f0575050507f87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b9495505b816000526001835260406000208660005283526115bd60406000209182546129be565b905560009081526001825260408082208683529092528190205481516001600160a01b03939093168352602083015290a2005b60405162461bcd60e51b8152928301869052601f908301527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400908201528690fd5b604098929398518681016323b872dd60e01b8152338383015230848301528a8583015284825261166082611961565b6040519161166d83611946565b8883527f5361666545524332303a206c6f772d6c6576656c2063616c6c206661696c656489840152873b156117655760006116b9939281925190828b5af16116b36129db565b906138fb565b8051806116ef575b5050505050507f87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b939461159a565b8188918101031261028d5786611705910161214e565b156117115780806116c1565b60405162461bcd60e51b8152938401869052602a908401527f5361666545524332303a204552433230206f7065726174696f6e20646964206e90830152691bdd081cdd58d8d9595960b21b90820152608490fd5b60405162461bcd60e51b81528088018a9052601d818601527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000818701528690fd5b60405162461bcd60e51b8152928301869052601490830152731a195b1908084f48195e1c1958dd195912195b1960621b908201528690fd5b50601f859260206040519362461bcd60e51b85528401528201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e006044820152fd5b3461028d57604036600319011261028d576001600160a01b03611841611c39565b16600052600160205260406000206024356000526020526020604060002054604051908152f35b503461028d57606036600319011261028d576001600160401b0360243581811161028d576118999036908401611a3f565b9060443590811161028d576118db926118f0926118bc6118c39336908401611b0e565b9135612a77565b92939190604051958695608087526080870190611bb4565b91151560208601528482036040860152611bb4565b9060608301520390f35b606081019081106001600160401b0382111761191557604052565b634e487b7160e01b600052604160045260246000fd5b608081019081106001600160401b0382111761191557604052565b604081019081106001600160401b0382111761191557604052565b60a081019081106001600160401b0382111761191557604052565b602081019081106001600160401b0382111761191557604052565b90601f801991011681019081106001600160401b0382111761191557604052565b6001600160401b0381116119155760051b60200190565b359060ff8216820361028d57565b6001600160401b03811161191557601f01601f191660200190565b81601f8201121561028d57803590611a0f826119dd565b92611a1d6040519485611997565b8284526020838301011161028d57816000926020809301838601378301015290565b9080601f8301121561028d578135611a56816119b8565b92604091611a6683519586611997565b808552602093848087019260051b8401019381851161028d57858401925b858410611a95575050505050505090565b6001600160401b03843581811161028d57860191608080601f19858803011261028d57845190611ac48261192b565b8a8501358252858501358b830152606090611ae08287016119cf565b8784015285013593841161028d57611aff878c809796819701016119f8565b90820152815201930192611a84565b81601f8201121561028d57803591611b25836119b8565b92611b336040519485611997565b808452602092838086019260051b82010192831161028d578301905b828210611b5d575050505090565b81358152908301908301611b4f565b60005b838110611b7f5750506000910152565b8181015183820152602001611b6f565b90602091611ba881518092818552858086019101611b6c565b601f01601f1916010190565b908082519081815260208091019281808460051b8301019501936000915b848310611be25750505050505090565b9091929394958480611c29600193601f198682030187528a51805182528381015184830152604060ff81830151169083015260608091015191608080928201520190611b8f565b9801930193019194939290611bd2565b600435906001600160a01b038216820361028d57565b35906001600160a01b038216820361028d57565b9080601f8301121561028d578135611c7a816119b8565b92604091611c8a83519586611997565b808552602093848087019260051b8401019381851161028d57858401925b858410611cb9575050505050505090565b6001600160401b03843581811161028d57860191601f1960608487038201811361028d57855191611ce9836118fa565b611cf48c8701611c4f565b83528686013585811161028d578790870191828a03011261028d57865190611d1b82611946565b8c810135600481101561028d5782528781013586811161028d578d8a91611d439301016119f8565b8c8201528b83015284013592831161028d57611d66868b80969581960101611a3f565b85820152815201930192611ca8565b9181601f8401121561028d578235916001600160401b03831161028d576020808501948460051b01011161028d57565b35906001600160401b038216820361028d57565b359065ffffffffffff8216820361028d57565b91909160808184031261028d5760405190611de68261192b565b819381356001600160401b03811161028d5782019080601f8301121561028d57813590611e12826119b8565b91611e206040519384611997565b808352602093848085019260051b82010192831161028d578401905b828210611e7c57505050606092611e779284928652611e5c818301611da5565b90860152611e6c60408201611c4f565b604086015201611db9565b910152565b848091611e8884611c4f565b815201910190611e3c565b9190604092838183031261028d578351848101916001600160401b039582841087851117611915578381528296823581811161028d5783019260808488031261028d57611edf8661192b565b833582811161028d5787611ef4918601611c63565b86526020958685013583811161028d5788611f109187016119f8565b606095869182890152611f24868201611db9565b60808901520135801515810361028d5760a087015285528581013591821161028d57019085601f8301121561028d57813590611f5f826119b8565b96611f6c82519889611997565b8288528685818a0194028501019381851161028d578701925b848410611f96575050505050500152565b858483031261028d578786918451611fad816118fa565b611fb6876119cf565b815282870135838201528587013586820152815201930192611f85565b92919092611fe0846119b8565b91611fee6040519384611997565b829480845260208094019060051b83019282841161028d5780915b84831061201857505050505050565b82356001600160401b03811161028d5786916120378684938601611e93565b815201920191612009565b9080601f8301121561028d5781602061205d93359101611fd3565b90565b90604060031983011261028d576001600160401b0360043581811161028d578361208c91600401611dcc565b9260243591821161028d5761205d91600401611e93565b906120ad826119b8565b6120ba6040519182611997565b82815280926120cb601f19916119b8565b0190602036910137565b60001981146120e45760010190565b634e487b7160e01b600052601160045260246000fd5b8051156121075760200190565b634e487b7160e01b600052603260045260246000fd5b8051600110156121075760400190565b80518210156121075760209160051b010190565b919082039182116120e457565b5190811515820361028d57565b90929192612168816119dd565b916121766040519384611997565b82948284528282011161028d576020612190930190611b6c565b565b919060408382031261028d576121a78361214e565b926020810151906001600160401b03821161028d57019080601f8301121561028d57815161205d9260200161215b565b906080810182519060808352815180915260a0830190602060a08260051b8601019301916000905b82821061224057505050506060612223819260208601518582036020870152611b8f565b9365ffffffffffff60408201511660408501520151151591015290565b90919293609f19868203018252845160018060a01b038151168252602081015160606020840152805160048082101561232d57506020604092612294926060870152015182608086015260a0850190611b8f565b910151916040818303910152815180825260208201906020808260051b85010194019260005b8281106122da5750505050506020806001929601920192019092916121ff565b9091929394602080612320600193601f19878203018952608060608b5180518452858101518685015260ff60408201511660408501520151918160608201520190611b8f565b97019501939291016122ba565b602190634e487b7160e01b6000525260246000fd5b9060208061235984516040855260408501906121d7565b93015191015290565b90815180825260208092019182818360051b85019501936000915b84831061238d5750505050505090565b90919293949584806123a783856001950387528a51612342565b980193019301919493929061237d565b600311156123c157565b634e487b7160e01b600052602160045260246000fd5b156123df5750565b60405162461bcd60e51b815260206004820152908190612403906024830190611b8f565b0390fd5b91909165ffffffffffff808094169116019182116120e457565b906080810191805160808352805180945260a083019360208092019060005b81811061248257505050808201516001600160401b0316908301526040808201516001600160a01b03169083015260609081015165ffffffffffff1691015290565b82516001600160a01b031687529583019591830191600101612440565b8051906124b4604092838552838501906121d7565b9060208091015193818184039101528080855193848152019401926000905b8382106124e257505050505090565b8451805160ff168752808401518785015281015186820152606090950194938201936001909101906124d3565b919061251a8361346c565b9261252484612969565b6060825101511561262557602061253b838361280f565b015190816000925b612600575060ff9051519116036125ce577f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901602084926125b761259065ffffffffffff4216925151613586565b6040519061259d8261192b565b60008252838583015260006040830152606082015261385f565b8460005260008352604060002055604051908152a2565b60405162461bcd60e51b815260206004820152600a60248201526921756e616e696d6f757360b01b6044820152606490fd5b6000198101908082116120e457169160ff8091169081146120e4576001019180612543565b60405162461bcd60e51b815260206004820152601360248201527214dd185d19481b5d5cdd08189948199a5b985b606a1b6044820152606490fd5b60005b825181101561269f576001600160a01b038061267f838661212d565b51169083161461269757612692906120d5565b612663565b505050600190565b505050600090565b6000916126fe610eee929594956127116126d86126d160018060a01b03604085015116958461277e565b988361280f565b9760405198899687958695634c9b6c0960e11b8752606060048801526064870190612421565b6003199384878303016024880152612362565b03915afa918215610f5c57600090819361272a57509190565b906127409293503d8091833e610f4d8183611997565b9091565b6040519061275182611946565b60006020836040516127628161192b565b6060815260608382015283604082015283606082015281520152565b81519161278a836119b8565b926127986040519485611997565b8084526127a7601f19916119b8565b0160005b8181106127f857505060005b81518110156127f257806127d86127d16127ed938561212d565b518561280f565b6127e2828761212d565b52610446818661212d565b6127b7565b50505090565b602090612803612744565b828288010152016127ab565b9190612819612744565b508051906040519161282a83611946565b82526020928383019260009283855283955b8082018051518810156128f95761286c906128658961285e869896518d61350f565b925161212d565b51906133a7565b6001600160a01b0390811694909390865b8a5180518210156128e95761289382889261212d565b511687146128a9576128a4906120d5565b61287d565b929891955093509060ff81116128d5579060016128cb921b87511787526120d5565b959291909261283c565b634e487b7160e01b86526011600452602486fd5b505093509350956128cb906120d5565b50505093509350505090565b61290e906138d0565b505065ffffffffffff8091169116111561292457565b60405162461bcd60e51b815260206004820152601c60248201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e000000006044820152606490fd5b61297460029161382c565b61297d816123b7565b1461298457565b60405162461bcd60e51b815260206004820152601260248201527121b430b73732b6103334b730b634bd32b21760711b6044820152606490fd5b919082018092116120e457565b91908110156121075760051b0190565b3d15612a06573d906129ec826119dd565b916129fa6040519384611997565b82523d6000602084013e565b606090565b90612a15826119b8565b604090612a2482519182611997565b8381528093612a35601f19916119b8565b0191600091825b848110612a4a575050505050565b6020908351612a588161192b565b8581528286818301528686830152606080830152828501015201612a3c565b9192908351801515600014612cf657612a8f90612a0b565b91600091612a9d8151612a0b565b95600190818097938960009586935b612aba575b50505050505050565b909192939495978351851015612ced57612ad4858561212d565b5151612ae0868561212d565b515260409060ff8083612af3898961212d565b5101511683612b02898861212d565b510152606080612b12898961212d565b51015181612b208a8961212d565b51015260209384612b318a8a61212d565b51015186811115612ce7575085965b8d8b51908b8215928315612cbd575b505050600014612c8c5750600283828f612b69908c61212d565b5101511614612c49578f96959493868f918f612c0690612c0c94612c18988f988f908f91612c129a898f94612be18f869288612bbc83612bb68884612bae848e61212d565b510151612141565b9361212d565b510152612bc9818761212d565b51519885612bd7838961212d565b510151169561212d565b51015194825196612bf18861192b565b875286015284015282015261040e838361212d565b506129be565b9c6120d5565b9561212d565b510151612c40575b612c3391612c2d91612141565b936120d5565b91909493928a9085612aac565b60009a50612c20565b5162461bcd60e51b815260048101859052601b60248201527f63616e6e6f74207472616e7366657220612067756172616e74656500000000006044820152606490fd5b9050612c18925088915084612ca783959e989796958a61212d565b51015184612cb5848461212d565b51015261212d565b821092509082612cd2575b50508e8b38612b4f565b612cde9192508d61212d565b51148a8f612cc8565b96612b40565b97829150612ab1565b50612a8f8151612a0b565b90815180825260208092019182818360051b82019501936000915b848310612d2c5750505050505090565b9091929394958181038352865160018060a01b03815116825285810151906060918288850152805160048082101561232d57509188612d849285948796839801520151604092839182608088015260a0870190611b8f565b91015193828183039101528351908181528581019286808460051b8401019601946000915b848310612dcc575050505050505090806001929801930193019194939290612d1c565b9193959780612e0e600193959799601f19898203018b526080878d5180518452858101518685015260ff89820151168985015201519181898201520190611b8f565b99019701930190918b979695939492612da9565b805160001981019081116120e457612e3990612a0b565b91612e44848361212d565b51606081015192604094855191612e5a83611946565b6000958684528660208095015287818051810103126131215787805191612e8083611946565b85810151835201519084810191825287998890899c8a988b5b87518d1015613002578f848e14612ff3578c8f8f90612f06858f8f908f612ec0878261212d565b51519582612ece898461212d565b5101516060612ee48a60ff85612bd7838961212d565b51015193825198612ef48a61192b565b8952880152860152606085015261212d565b52612f11848d61212d565b5087159081612fdd575b50612fa3575b501580612f8e575b612f40575b612c0c612f3a916120d5565b9b612e99565b9e5098612f83908f612f6e8b612f648f612f5a839161211d565b510151938d61212d565b51019182516129be565b905289612f7a8d61211d565b510151906129be565b60019e909990612f2e565b50612f998d8961212d565b5151875114612f29565b829c919650612f7a818c612fcc8f612f64612fd39882612fc381996120fa565b5101519461212d565b90526120fa565b996001948c612f21565b612fe891508b61212d565b51518851148f612f1b565b509b9d50612f3a60019e6120d5565b509899509c969a995050939992505050156130e557156130ab5715613070578301510361302e57505090565b60649250519062461bcd60e51b825280600483015260248201527f746f74616c5265636c61696d6564213d67756172616e7465652e616d6f756e746044820152fd5b825162461bcd60e51b815260048101859052601460248201527318dbdd5b19081b9bdd08199a5b99081c9a59da1d60621b6044820152606490fd5b835162461bcd60e51b815260048101869052601360248201527218dbdd5b19081b9bdd08199a5b99081b19599d606a1b6044820152606490fd5b845162461bcd60e51b815260048101879052601560248201527418dbdd5b19081b9bdd08199a5b99081d185c99d95d605a1b6044820152606490fd5b8680fd5b80516001600160a01b03908116919060005b6040808401908151918251841015612ab1578461315585809561212d565b51519161316660209586925161212d565b510151918060a01c156000146132695716876131df57600080809381935af161318d6129db565b50156131a357505061319e906120d5565b613137565b60649250519062461bcd60e51b825260048201526016602482015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b6044820152fd5b825163a9059cbb60e01b81526001600160a01b039190911660048201526024810191909152929190818460448160008b5af190811561325f5750613229575b5061319e91506120d5565b82813d8311613258575b61323d8183611997565b8101031261028d5761325161319e9261214e565b503861321e565b503d613233565b513d6000823e3d90fd5b6000898152600186528481209181529452509120805461319e939261328d916129be565b90556120d5565b9161329e906138d0565b936001600160a01b03938493506132b69250906138a3565b169116036132c057565b60405162461bcd60e51b81526020600482015260156024820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b6044820152606490fd5b61330860029161382c565b613311816123b7565b0361331857565b60405162461bcd60e51b815260206004820152601660248201527521b430b73732b6103737ba103334b730b634bd32b21760511b6044820152606490fd5b919061339691613365846138d0565b509290604051936133758561192b565b65ffffffffffff80921685521660208401526040830152606082015261385f565b906000526000602052604060002055565b90600060806020926040948551858101917f19457468657265756d205369676e6564204d6573736167653a0a3332000000008352603c820152603c81526133ed816118fa565b5190209060ff8151169086868201519101519187519384528684015286830152606082015282805260015afa1561325f57600051906001600160a01b03821615613435575090565b5162461bcd60e51b8152602060048201526011602482015270496e76616c6964207369676e617475726560781b6044820152606490fd5b8051906020916001600160401b03838301511660018060a01b039165ffffffffffff606084604087015116950151166040519485938785019760a086019060808a5285518092528060c088019601976000905b8382106134f2575050505060408501526060840152608083015203601f19810183526134ec915082611997565b51902090565b895181168852988201988a985096820196600191909101906134bf565b6134ec61351e61355c9261346c565b926020810151815191606065ffffffffffff604083015116910151151561356f60405196879460208601998a5260a0604087015260c0860190611b8f565b601f199586868303016060870152612d01565b91608084015260a083015203908101835282611997565b6040516134ec816108266020820194602086526040830190612d01565b9080601f8301121561028d57815161205d9260200161215b565b8051810160208282031261028d5760208201516001600160401b03811161028d5760208201603f82850101121561028d5760208184010151906135ff826119b8565b9361360d6040519586611997565b82855260208501916020850160408560051b83850101011161028d57604081830101925b60408560051b8385010101841061364b5750505050505090565b83516001600160401b03811161028d57601f199084840101606081890383011261028d576040519161367c836118fa565b60408201516001600160a01b038116810361028d57835260608201516001600160401b03811161028d57604090830191828b03011261028d57604051906136c282611946565b6040810151600481101561028d5782526060810151906001600160401b03821161028d5760406136f89260208d019201016135a3565b6020820152602083015260808101516001600160401b03811161028d5760208901605f82840101121561028d576040818301015190613736826119b8565b926137446040519485611997565b828452602084019060208c0160608560051b85840101011161028d57606083820101915b60608560051b858401010183106137915750505050506040820152815260209384019301613631565b82516001600160401b03811161028d57608083860182018f03603f19011261028d57604051916137c08361192b565b8386018201606081015184526080810151602085015260a0015160ff8116810361028d57604084015260c082878601010151926001600160401b03841161028d578f60209493606086958661381c9401928b8a010101016135a3565b6060820152815201920191613768565b61383c65ffffffffffff916138d0565b50905016801560001461384f5750600090565b421061385a57600290565b600190565b65ffffffffffff60d01b815160d01b1665ffffffffffff60a01b602083015160a01b16179061389e6040820151606060018060a01b03930151906138a3565b161790565b60405191602083019182526040830152604082526138c0826118fa565b905190206001600160a01b031690565b60005260006020526040600020548060d01c9165ffffffffffff8260a01c169160018060a01b031690565b90919015613907575090565b8151156123df5750805190602001fdfea2646970667358221220c54d1238ff95db42065823e00d05476cfa0413f246256cd8daf89c3380173a1564736f6c63430008110033",
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
	parsed, err := NitroAdjudicatorMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
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

// StateIsSupported is a free data retrieval call binding the contract method 0x5685b7dc.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) view returns(bool, string)
func (_NitroAdjudicator *NitroAdjudicatorCaller) StateIsSupported(opts *bind.CallOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (bool, string, error) {
	var out []interface{}
	err := _NitroAdjudicator.contract.Call(opts, &out, "stateIsSupported", fixedPart, proof, candidate)

	if err != nil {
		return *new(bool), *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)
	out1 := *abi.ConvertType(out[1], new(string)).(*string)

	return out0, out1, err

}

// StateIsSupported is a free data retrieval call binding the contract method 0x5685b7dc.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) view returns(bool, string)
func (_NitroAdjudicator *NitroAdjudicatorSession) StateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (bool, string, error) {
	return _NitroAdjudicator.Contract.StateIsSupported(&_NitroAdjudicator.CallOpts, fixedPart, proof, candidate)
}

// StateIsSupported is a free data retrieval call binding the contract method 0x5685b7dc.
//
// Solidity: function stateIsSupported((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) view returns(bool, string)
func (_NitroAdjudicator *NitroAdjudicatorCallerSession) StateIsSupported(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (bool, string, error) {
	return _NitroAdjudicator.Contract.StateIsSupported(&_NitroAdjudicator.CallOpts, fixedPart, proof, candidate)
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

// Challenge is a paid mutator transaction binding the contract method 0x8286a060.
//
// Solidity: function challenge((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Challenge(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "challenge", fixedPart, proof, candidate, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x8286a060.
//
// Solidity: function challenge((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Challenge(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate, challengerSig)
}

// Challenge is a paid mutator transaction binding the contract method 0x8286a060.
//
// Solidity: function challenge((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate, (uint8,bytes32,bytes32) challengerSig) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Challenge(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart, challengerSig INitroTypesSignature) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Challenge(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate, challengerSig)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x6d2a9c92.
//
// Solidity: function checkpoint((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Checkpoint(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "checkpoint", fixedPart, proof, candidate)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x6d2a9c92.
//
// Solidity: function checkpoint((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Checkpoint(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// Checkpoint is a paid mutator transaction binding the contract method 0x6d2a9c92.
//
// Solidity: function checkpoint((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Checkpoint(fixedPart INitroTypesFixedPart, proof []INitroTypesSignedVariablePart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Checkpoint(&_NitroAdjudicator.TransactOpts, fixedPart, proof, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0xee049b50.
//
// Solidity: function conclude((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Conclude(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "conclude", fixedPart, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0xee049b50.
//
// Solidity: function conclude((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Conclude(fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, candidate)
}

// Conclude is a paid mutator transaction binding the contract method 0xee049b50.
//
// Solidity: function conclude((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Conclude(fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Conclude(&_NitroAdjudicator.TransactOpts, fixedPart, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xec346235.
//
// Solidity: function concludeAndTransferAllAssets((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) ConcludeAndTransferAllAssets(opts *bind.TransactOpts, fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "concludeAndTransferAllAssets", fixedPart, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xec346235.
//
// Solidity: function concludeAndTransferAllAssets((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, candidate)
}

// ConcludeAndTransferAllAssets is a paid mutator transaction binding the contract method 0xec346235.
//
// Solidity: function concludeAndTransferAllAssets((address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) ConcludeAndTransferAllAssets(fixedPart INitroTypesFixedPart, candidate INitroTypesSignedVariablePart) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.ConcludeAndTransferAllAssets(&_NitroAdjudicator.TransactOpts, fixedPart, candidate)
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

// DepositBatch is a paid mutator transaction binding the contract method 0xf7cb41ba.
//
// Solidity: function deposit_batch(address asset, bytes32[] channelIds, uint256[] expectedHelds, uint256[] amounts) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) DepositBatch(opts *bind.TransactOpts, asset common.Address, channelIds [][32]byte, expectedHelds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "deposit_batch", asset, channelIds, expectedHelds, amounts)
}

// DepositBatch is a paid mutator transaction binding the contract method 0xf7cb41ba.
//
// Solidity: function deposit_batch(address asset, bytes32[] channelIds, uint256[] expectedHelds, uint256[] amounts) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) DepositBatch(asset common.Address, channelIds [][32]byte, expectedHelds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.DepositBatch(&_NitroAdjudicator.TransactOpts, asset, channelIds, expectedHelds, amounts)
}

// DepositBatch is a paid mutator transaction binding the contract method 0xf7cb41ba.
//
// Solidity: function deposit_batch(address asset, bytes32[] channelIds, uint256[] expectedHelds, uint256[] amounts) payable returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) DepositBatch(asset common.Address, channelIds [][32]byte, expectedHelds []*big.Int, amounts []*big.Int) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.DepositBatch(&_NitroAdjudicator.TransactOpts, asset, channelIds, expectedHelds, amounts)
}

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) reclaimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) Reclaim(opts *bind.TransactOpts, reclaimArgs IMultiAssetHolderReclaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "reclaim", reclaimArgs)
}

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) reclaimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) Reclaim(reclaimArgs IMultiAssetHolderReclaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Reclaim(&_NitroAdjudicator.TransactOpts, reclaimArgs)
}

// Reclaim is a paid mutator transaction binding the contract method 0xd3c4e738.
//
// Solidity: function reclaim((bytes32,bytes32,bytes,uint256,uint256,bytes32,bytes,uint256) reclaimArgs) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactorSession) Reclaim(reclaimArgs IMultiAssetHolderReclaimArgs) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.Reclaim(&_NitroAdjudicator.TransactOpts, reclaimArgs)
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

// TransferAllAssets is a paid mutator transaction binding the contract method 0x31afa0b4.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorTransactor) TransferAllAssets(opts *bind.TransactOpts, channelId [32]byte, outcome []ExitFormatSingleAssetExit, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.contract.Transact(opts, "transferAllAssets", channelId, outcome, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x31afa0b4.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
func (_NitroAdjudicator *NitroAdjudicatorSession) TransferAllAssets(channelId [32]byte, outcome []ExitFormatSingleAssetExit, stateHash [32]byte) (*types.Transaction, error) {
	return _NitroAdjudicator.Contract.TransferAllAssets(&_NitroAdjudicator.TransactOpts, channelId, outcome, stateHash)
}

// TransferAllAssets is a paid mutator transaction binding the contract method 0x31afa0b4.
//
// Solidity: function transferAllAssets(bytes32 channelId, (address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[] outcome, bytes32 stateHash) returns()
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
	FinalHoldings   *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterAllocationUpdated is a free log retrieval operation binding the contract event 0xc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings, uint256 finalHoldings)
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

// WatchAllocationUpdated is a free log subscription operation binding the contract event 0xc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings, uint256 finalHoldings)
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

// ParseAllocationUpdated is a log parse operation binding the contract event 0xc36da2054c5669d6dac211b7366d59f2d369151c21edf4940468614b449e0b9a.
//
// Solidity: event AllocationUpdated(bytes32 indexed channelId, uint256 assetIndex, uint256 initialHoldings, uint256 finalHoldings)
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

// FilterChallengeRegistered is a free log retrieval operation binding the contract event 0x0e6d8300485cb09fa95f22b89b46f7b0cc3029f1bbf257a0884414d415546cf8.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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

// WatchChallengeRegistered is a free log subscription operation binding the contract event 0x0e6d8300485cb09fa95f22b89b46f7b0cc3029f1bbf257a0884414d415546cf8.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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

// ParseChallengeRegistered is a log parse operation binding the contract event 0x0e6d8300485cb09fa95f22b89b46f7b0cc3029f1bbf257a0884414d415546cf8.
//
// Solidity: event ChallengeRegistered(bytes32 indexed channelId, uint48 turnNumRecord, uint48 finalizesAt, bool isFinal, (address[],uint64,address,uint48) fixedPart, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[])[] proof, (((address,(uint8,bytes),(bytes32,uint256,uint8,bytes)[])[],bytes,uint48,bool),(uint8,bytes32,bytes32)[]) candidate)
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
	DestinationHoldings *big.Int
	Raw                 types.Log // Blockchain specific contextual infos
}

// FilterDeposited is a free log retrieval operation binding the contract event 0x87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 destinationHoldings)
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

// WatchDeposited is a free log subscription operation binding the contract event 0x87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 destinationHoldings)
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

// ParseDeposited is a log parse operation binding the contract event 0x87d4c0b5e30d6808bc8a94ba1c4d839b29d664151551a31753387ee9ef48429b.
//
// Solidity: event Deposited(bytes32 indexed destination, address asset, uint256 destinationHoldings)
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

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
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"}],\"name\":\"AllocationUpdated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"newTurnNumRecord\",\"type\":\"uint48\"}],\"name\":\"ChallengeCleared\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"indexed\":false,\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"ChallengeRegistered\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"}],\"name\":\"Concluded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amountDeposited\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"destinationHoldings\",\"type\":\"uint256\"}],\"name\":\"Deposited\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature\",\"name\":\"challengerSig\",\"type\":\"tuple\"}],\"name\":\"challenge\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"checkpoint\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"sourceAllocations\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"targetAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"}],\"name\":\"compute_reclaim_effects\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"initialHoldings\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"compute_transfer_effects_and_interactions\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"newAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"bool\",\"name\":\"allocatesOnlyZeros\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"exitAllocations\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"totalPayouts\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"conclude\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"concludeAndTransferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"expectedHeld\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getChainID\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"holdings\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"chainId\",\"type\":\"uint256\"},{\"internalType\":\"address[]\",\"name\":\"participants\",\"type\":\"address[]\"},{\"internalType\":\"uint48\",\"name\":\"channelNonce\",\"type\":\"uint48\"},{\"internalType\":\"address\",\"name\":\"appDefinition\",\"type\":\"address\"},{\"internalType\":\"uint48\",\"name\":\"challengeDuration\",\"type\":\"uint48\"}],\"internalType\":\"structINitroTypes.FixedPart\",\"name\":\"fixedPart\",\"type\":\"tuple\"},{\"components\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"variablePart\",\"type\":\"tuple\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"v\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"r\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"s\",\"type\":\"bytes32\"}],\"internalType\":\"structINitroTypes.Signature[]\",\"name\":\"sigs\",\"type\":\"tuple[]\"}],\"internalType\":\"structINitroTypes.SignedVariablePart[]\",\"name\":\"signedVariableParts\",\"type\":\"tuple[]\"}],\"name\":\"latestSupportedState\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes\",\"name\":\"appData\",\"type\":\"bytes\"},{\"internalType\":\"uint48\",\"name\":\"turnNum\",\"type\":\"uint48\"},{\"internalType\":\"bool\",\"name\":\"isFinal\",\"type\":\"bool\"}],\"internalType\":\"structINitroTypes.VariablePart\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"pure\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"sourceChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"sourceStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sourceOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"sourceAssetIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"indexOfTargetInSource\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"targetStateHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"targetOutcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"targetAssetIndex\",\"type\":\"uint256\"}],\"internalType\":\"structIMultiAssetHolder.ClaimArgs\",\"name\":\"claimArgs\",\"type\":\"tuple\"}],\"name\":\"reclaim\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"statusOf\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"assetIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"fromChannelId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"outcomeBytes\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"},{\"internalType\":\"uint256[]\",\"name\":\"indices\",\"type\":\"uint256[]\"}],\"name\":\"transfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"destination\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint8\",\"name\":\"allocationType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"metadata\",\"type\":\"bytes\"}],\"internalType\":\"structExitFormat.Allocation[]\",\"name\":\"allocations\",\"type\":\"tuple[]\"}],\"internalType\":\"structExitFormat.SingleAssetExit[]\",\"name\":\"outcome\",\"type\":\"tuple[]\"},{\"internalType\":\"bytes32\",\"name\":\"stateHash\",\"type\":\"bytes32\"}],\"name\":\"transferAllAssets\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"channelId\",\"type\":\"bytes32\"}],\"name\":\"unpackStatus\",\"outputs\":[{\"internalType\":\"uint48\",\"name\":\"turnNumRecord\",\"type\":\"uint48\"},{\"internalType\":\"uint48\",\"name\":\"finalizesAt\",\"type\":\"uint48\"},{\"internalType\":\"uint160\",\"name\":\"fingerprint\",\"type\":\"uint160\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561001057600080fd5b50613e0d806100206000396000f3fe6080604052600436106100e85760003560e01c80633503f8f01161008a578063af69c9d711610059578063af69c9d714610286578063c7df14e2146102a6578063c927c2d3146102c6578063d3c4e738146102e6576100e8565b80633503f8f0146101f5578063552cfa5014610215578063564b81ef14610244578063566d54c614610259576100e8565b806312c52592116100c657806312c5259214610175578063166e56cd146101955780632fb1d270146101c25780633033730e146101d5576100e8565b80630d36d531146100ed5780630d96cc221461010f57806311e9f17814610145575b600080fd5b3480156100f957600080fd5b5061010d610108366004612e7f565b610306565b005b34801561011b57600080fd5b5061012f61012a366004612de6565b61032e565b60405161013c9190613b51565b60405180910390f35b34801561015157600080fd5b50610165610160366004613054565b6103ef565b60405161013c94939291906134b7565b34801561018157600080fd5b5061010d610190366004612edf565b61073d565b3480156101a157600080fd5b506101b56101b0366004612ba6565b61088e565b60405161013c9190613509565b61010d6101d0366004612bd1565b6108ab565b3480156101e157600080fd5b5061010d6101f03660046130bd565b610b3a565b34801561020157600080fd5b5061010d610210366004612e7f565b610bba565b34801561022157600080fd5b50610235610230366004612cc2565b610c07565b60405161013c93929190613cc0565b34801561025057600080fd5b506101b5610c22565b34801561026557600080fd5b50610279610274366004612c0b565b610c31565b60405161013c91906134a4565b34801561029257600080fd5b5061010d6102a1366004612cda565b610edf565b3480156102b257600080fd5b506101b56102c1366004612cc2565b611275565b3480156102d257600080fd5b5061010d6102e1366004612e7f565b611287565b3480156102f257600080fd5b5061010d610301366004612d16565b611291565b60006103128383611301565b905061032981610321846113b7565b516000610edf565b505050565b610336612321565b60006103428385613d4c565b90506103546080860160608701612b8a565b6001600160a01b031663f409586c8661037561036f82613d59565b856113e3565b6040518363ffffffff1660e01b8152600401610392929190613a66565b60006040518083038186803b1580156103aa57600080fd5b505afa1580156103be573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526103e69190810190612f97565b95945050505050565b606060006060600080855111610406578551610409565b84515b6001600160401b038111801561041e57600080fd5b5060405190808252806020026020018201604052801561045857816020015b610445612353565b81526020019060019003908161043d5790505b5091506000905085516001600160401b038111801561047657600080fd5b506040519080825280602002602001820160405280156104b057816020015b61049d612353565b8152602001906001900390816104955790505b50935060019250866000805b8851811015610731578881815181106104d157fe5b6020026020010151600001518782815181106104e957fe5b6020026020010151600001818152505088818151811061050557fe5b60200260200101516040015187828151811061051d57fe5b60200260200101516040019060ff16908160ff168152505088818151811061054157fe5b60200260200101516060015187828151811061055957fe5b602002602001015160600181905250600061058b8a838151811061057957fe5b6020026020010151602001518561148d565b90508851600014806105ba57508851831080156105ba5750818984815181106105b057fe5b6020026020010151145b156106cc57600260ff168a84815181106105d057fe5b60200260200101516040015160ff1614156106065760405162461bcd60e51b81526004016105fd90613824565b60405180910390fd5b808a838151811061061357fe5b6020026020010151602001510388838151811061062c57fe5b6020026020010151602001818152505060405180608001604052808b848151811061065357fe5b60200260200101516000015181526020018281526020018b848151811061067657fe5b60200260200101516040015160ff1681526020018b848151811061069657fe5b6020026020010151606001518152508684815181106106b157fe5b60200260200101819052508085019450826001019250610701565b8982815181106106d857fe5b6020026020010151602001518883815181106106f057fe5b602002602001015160200181815250505b87828151811061070d57fe5b60200260200101516020015160001461072557600096505b909203916001016104bc565b50505093509350935093565b6000610748846114a5565b90506000610755846113b7565b60400151905060006107668361151b565b600281111561077157fe5b1415610786576107818282611565565b6107b5565b60016107918361151b565b600281111561079c57fe5b14156107ac5761078182826115a4565b6107b5826115e2565b60006107c2868685611619565b90506107d3818760200151866116d4565b827f5bc2767c561afcb4ffbef365ef0757d1d4389267288cbedb1014507e356ff1ec8388608001514201610806896113b7565b606001518a8a60405161081d959493929190613bd0565b60405180910390a261087460405180608001604052808465ffffffffffff1681526020018860800151420165ffffffffffff16815260200183815260200161086d610867896113b7565b5161172e565b9052611747565b600093845260208490526040909320929092555050505050565b600160209081526000928352604080842090915290825290205481565b6108b483611798565b156108d15760405162461bcd60e51b81526004016105fd90613923565b6001600160a01b0384166000908152600160209081526040808320868452909152812054838110156109155760405162461bcd60e51b81526004016105fd906135de565b61091f848461179f565b811061093d5760405162461bcd60e51b81526004016105fd9061398a565b6109518161094b868661179f565b906117f9565b91506001600160a01b038616610985578234146109805760405162461bcd60e51b81526004016105fd906139f8565b610a23565b6040516323b872dd60e01b81526001600160a01b038716906323b872dd906109b590339030908790600401613446565b602060405180830381600087803b1580156109cf57600080fd5b505af11580156109e3573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610a079190612ca6565b610a235760405162461bcd60e51b81526004016105fd906138ec565b6000610a2f828461179f565b6001600160a01b03881660009081526001602090815260408083208a8452909152908190208290555190915086907f2dcdaad87b561ba5a69835009b4c53ef9d3c41ca6cc9574049187659d6c6a71590610a8e908a9087908690613483565b60405180910390a26001600160a01b038716610b31576000610ab085856117f9565b90506000336001600160a01b031682604051610acb9061220e565b60006040518083038185875af1925050503d8060008114610b08576040519150601f19603f3d011682016040523d82523d6000602084013e610b0d565b606091505b5050905080610b2e5760405162461bcd60e51b81526004016105fd90613a2f565b50505b50505050505050565b6000806000610b4c888589888a611856565b9250925092506000806000610b7984878d81518110610b6757fe5b602002602001015160400151896103ef565b93509350509250610b908b868c8b8a888a886118d4565b610bad868c81518110610b9f57fe5b602002602001015183611992565b5050505050505050505050565b6000610bc5836114a5565b90506000610bd2836113b7565b604001519050610be1826115e2565b610beb82826115a4565b610bf6848484611619565b50610c0182826119ca565b50505050565b6000806000610c1584611a51565b9196909550909350915050565b6000610c2c611a6f565b905090565b6060600060018551036001600160401b0381118015610c4f57600080fd5b50604051908082528060200260200182016040528015610c8957816020015b610c76612353565b815260200190600190039081610c6e5790505b5090506000858481518110610c9a57fe5b60200260200101519050600280811115610cb057fe5b60ff16816040015160ff1614610cd85760405162461bcd60e51b81526004016105fd906137fb565b6000610ce78260600151611a73565b905060008060008060005b8b51811015610e785789811415610d0c5760019450610e70565b60405180608001604052808d8381518110610d2357fe5b60200260200101516000015181526020018d8381518110610d4057fe5b60200260200101516020015181526020018d8381518110610d5d57fe5b60200260200101516040015160ff1681526020018d8381518110610d7d57fe5b602002602001015160600151815250888381518110610d9857fe5b602002602001018190525085600001518c8281518110610db457fe5b6020026020010151600001511415610e06578a600081518110610dd357fe5b602002602001015160200151888381518110610deb57fe5b60200260200101516020018181510191508181525050600193505b85602001518c8281518110610e1757fe5b6020026020010151600001511415610e69578a600181518110610e3657fe5b602002602001015160200151888381518110610e4e57fe5b60200260200101516020018181510191508181525050600192505b6001909101905b600101610cf2565b5083610e965760405162461bcd60e51b81526004016105fd90613795565b82610eb35760405162461bcd60e51b81526004016105fd90613615565b81610ed05760405162461bcd60e51b81526004016105fd90613699565b50949998505050505050505050565b610ee883611a8f565b610efb81610ef58461172e565b85611ac2565b81516001906000906001600160401b0381118015610f1857600080fd5b50604051908082528060200260200182016040528015610f5257816020015b610f3f612379565b815260200190600190039081610f375790505b509050600084516001600160401b0381118015610f6e57600080fd5b50604051908082528060200260200182016040528015610f98578160200160208202803683370190505b509050600085516001600160401b0381118015610fb457600080fd5b50604051908082528060200260200182016040528015610fde578160200160208202803683370190505b50905060005b8651811015611178576000878281518110610ffb57fe5b60200260200101519050600081604001519050600089848151811061101c57fe5b602002602001015160000151905060016000826001600160a01b03166001600160a01b0316815260200190815260200160002060008c81526020019081526020016000205486858151811061106d57fe5b6020026020010181815250506000806000806110de8a898151811061108e57fe5b60200260200101518760006001600160401b03811180156110ae57600080fd5b506040519080825280602002602001820160405280156110d8578160200160208202803683370190505b506103ef565b9350935093509350826110f05760009b505b808989815181106110fd57fe5b602002602001018181525050838e898151811061111657fe5b6020026020010151604001819052506040518060600160405280866001600160a01b0316815260200188602001518152602001838152508b898151811061115957fe5b6020026020010181905250505050505050508080600101915050610fe4565b5060005b865181101561123e57600087828151811061119357fe5b60200260200101516000015190508282815181106111ad57fe5b6020908102919091018101516001600160a01b03831660009081526001835260408082208d8352909352919091208054919091039055835189907fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c13063767990849087908290811061121757fe5b602002602001015160405161122d929190613baf565b60405180910390a25060010161117c565b5083156112595760008781526020819052604081205561126c565b61126c87866112678961172e565b611b0b565b610b3183611b71565b60006020819052908152604090205481565b6103298282611301565b60008061129d83611ba1565b9150915060606000838560600151815181106112b557fe5b60200260200101516040015190506000838660e00151815181106112d557fe5b60200260200101516040015190506112f282828860800151610c31565b92505050610c01848483611d11565b600061130c836114a5565b9050611317816115e2565b611322838383611619565b506113636040518060800160405280600065ffffffffffff1681526020014265ffffffffffff1681526020016000801b815260200161086d610867866113b7565b60008083815260200190815260200160002081905550807f4f465027a3d06ea73dd12be0f5c5fc0a34e21f19d6eaed4834a7a944edabc901426040516113a99190613bbd565b60405180910390a292915050565b6113bf612321565b816001835103815181106113cf57fe5b60200260200101516000015190505b919050565b6060600082516001600160401b03811180156113fe57600080fd5b5060405190808252806020026020018201604052801561143857816020015b6114256123a3565b81526020019060019003908161141d5790505b50905060005b8351811015611483576114648585838151811061145757fe5b6020026020010151611d58565b82828151811061147057fe5b602090810291909101015260010161143e565b5090505b92915050565b600081831161149c578261149e565b815b9392505050565b60006114af611a6f565b8251146114ce5760405162461bcd60e51b81526004016105fd90613642565b6114d6611a6f565b82602001518360400151846060015185608001516040516020016114fe959493929190613b64565b604051602081830303815290604052805190602001209050919050565b60008061152783611a51565b5091505065ffffffffffff81166115425760009150506113de565b428165ffffffffffff161161155b5760029150506113de565b60019150506113de565b600061157083611a51565b505090508065ffffffffffff168265ffffffffffff1610156103295760405162461bcd60e51b81526004016105fd906136c7565b60006115af83611a51565b505090508065ffffffffffff168265ffffffffffff16116103295760405162461bcd60e51b81526004016105fd906135a7565b60026115ed8261151b565b60028111156115f857fe5b14156116165760405162461bcd60e51b81526004016105fd9061366d565b50565b60008084606001516001600160a01b031663f409586c8661163a88886113e3565b6040518363ffffffff1660e01b8152600401611657929190613b2c565b60006040518083038186803b15801561166f57600080fd5b505afa158015611683573d6000803e3d6000fd5b505050506040513d6000823e601f3d908101601f191682016040526116ab9190810190612f97565b90506116b78185611e1d565b6103e6838260200151836000015184604001518560600151611e9c565b6000611706846040516020016116ea919061355f565b6040516020818303038152906040528051906020012083611ed8565b90506117128184611f92565b610c015760405162461bcd60e51b81526004016105fd9061372e565b600061173982611fe8565b805190602001209050919050565b805160208201516040830151606084015160009360d01b6001600160d01b03191660a093841b65ffffffffffff60a01b1617929161178491612011565b6001600160a01b0316919091179392505050565b60a01c1590565b60008282018381101561149e576040805162461bcd60e51b815260206004820152601b60248201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604482015290519081900360640190fd5b600082821115611850576040805162461bcd60e51b815260206004820152601e60248201527f536166654d6174683a207375627472616374696f6e206f766572666c6f770000604482015290519081900360640190fd5b50900390565b60606000806118648761203d565b61186d86611a8f565b61187f85858051906020012088611ac2565b6118888461209c565b925082888151811061189657fe5b602090810291909101810151516001600160a01b03811660009081526001835260408082209982529890925296902054929895975091955050505050565b6001600160a01b03871660009081526001602090815260408083208984529091529020805482900390558351839085908a90811061190e57fe5b60200260200101516040018190525061194e86868660405160200161193391906134f6565b60405160208183030381529060405280519060200120611b0b565b857fb3917fd12b23b8d48703d554ab284c5b1912bb5c67e710c7534a56c1306376798984604051611980929190613baf565b60405180910390a25050505050505050565b6119c6604051806060016040528084600001516001600160a01b0316815260200184602001518152602001838152506120b2565b5050565b6040805160808101825265ffffffffffff8316815260006020820181905291810182905260608101919091526119ff90611747565b60008084815260200190815260200160002081905550817f07da0a0674fb921e484018c8b81d80e292745e5d8ed134b580c8b9c631c5e9e082604051611a459190613bbd565b60405180910390a25050565b60009081526020819052604090205460d081901c9160a082901c9190565b4690565b611a7b6123c3565b818060200190518101906114879190612f50565b6002611a9a8261151b565b6002811115611aa557fe5b146116165760405162461bcd60e51b81526004016105fd9061395a565b6000611acd82611a51565b92505050611adb8484612011565b6001600160a01b0316816001600160a01b031614610c015760405162461bcd60e51b81526004016105fd906138bd565b600080611b1785611a51565b50915091506000611b5760405180608001604052808565ffffffffffff1681526020018465ffffffffffff16815260200187815260200186815250611747565b600096875260208790526040909620959095555050505050565b60005b81518110156119c657611b99828281518110611b8c57fe5b60200260200101516120b2565b600101611b74565b8051604082015160608381015160c085015160e086015192948594909390929190611bcb85611a8f565b611be18860200151858051906020012087611ac2565b611bea8461209c565b9650611bf58261209c565b95506000878481518110611c0557fe5b6020908102919091010151519050600260ff16888581518110611c2457fe5b6020026020010151604001518381518110611c3b57fe5b60200260200101516040015160ff1614611c675760405162461bcd60e51b81526004016105fd906139c1565b6000888581518110611c7557fe5b6020026020010151604001518a6080015181518110611c9057fe5b6020026020010151600001519050816001600160a01b0316888481518110611cb457fe5b6020026020010151600001516001600160a01b031614611ce65760405162461bcd60e51b81526004016105fd906137c4565b611cef81611a8f565b611d058a60a00151858051906020012083611ac2565b50505050505050915091565b8251606084015183518390859083908110611d2857fe5b602002602001015160400181905250611d518286602001518660405160200161193391906134f6565b5050505050565b611d606123a3565b60408051808201909152825181526000602082018190525b836020015151811015611483576000611db5611d9887876000015161215e565b86602001518481518110611da857fe5b6020026020010151611ed8565b905060005b866020015151811015611e135786602001518181518110611dd757fe5b60200260200101516001600160a01b0316826001600160a01b03161415611e0b57602084018051600283900a019052611e13565b600101611dba565b5050600101611d78565b611e8081600183510381518110611e3057fe5b602002602001015160000151604051602001611e4c9190613b51565b60405160208183030381529060405283604051602001611e6c9190613b51565b6040516020818303038152906040526121aa565b6119c65760405162461bcd60e51b81526004016105fd9061385b565b60008585858585604051602001611eb7959493929190613512565b60405160208183030381529060405280519060200120905095945050505050565b60008083604051602001611eec9190613415565b604051602081830303815290604052805190602001209050600060018285600001518660200151876040015160405160008152602001604052604051611f359493929190613589565b6020604051602081039080840390855afa158015611f57573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b038116611f8a5760405162461bcd60e51b81526004016105fd90613892565b949350505050565b6000805b8251811015611fde57828181518110611fab57fe5b60200260200101516001600160a01b0316846001600160a01b03161415611fd6576001915050611487565b600101611f96565b5060009392505050565b606081604051602001611ffb91906134f6565b6040516020818303038152906040529050919050565b604080516020808201949094528082019290925280518083038201815260609092019052805191012090565b60005b81518160010110156119c65781816001018151811061205b57fe5b602002602001015182828151811061206f57fe5b6020026020010151106120945760405162461bcd60e51b81526004016105fd906136fe565b600101612040565b6060818060200190518101906114879190612c74565b805160005b826040015151811015610329576000836040015182815181106120d657fe5b60200260200101516000015190506000846040015183815181106120f657fe5b602002602001015160200151905061210d82611798565b1561212a576121258461211f8461220e565b83612211565b612154565b6001600160a01b038416600090815260016020908152604080832085845290915290208054820190555b50506001016120b7565b6000612169836114a5565b60208084015184516040808701516060880151915161218c969591929101613512565b60405160208183030381529060405280519060200120905092915050565b8151815160009160019181148083146121c65760009250612204565b600160208701838101602088015b6002848385100114156121ff5780518351146121f35760009650600093505b602092830192016121d4565b505050505b5090949350505050565b90565b6001600160a01b0383166122a1576000826001600160a01b0316826040516122389061220e565b60006040518083038185875af1925050503d8060008114612275576040519150601f19603f3d011682016040523d82523d6000602084013e61227a565b606091505b505090508061229b5760405162461bcd60e51b81526004016105fd90613765565b50610329565b60405163a9059cbb60e01b81526001600160a01b0384169063a9059cbb906122cf908590859060040161346a565b602060405180830381600087803b1580156122e957600080fd5b505af11580156122fd573d6000803e3d6000fd5b505050506040513d601f19601f82011682018060405250810190610c019190612ca6565b60405180608001604052806060815260200160608152602001600065ffffffffffff1681526020016000151581525090565b604080516080810182526000808252602082018190529181019190915260608082015290565b604051806060016040528060006001600160a01b0316815260200160608152602001606081525090565b60405180604001604052806123b6612321565b8152602001600081525090565b604080518082019091526000808252602082015290565b60006123ed6123e884613d0e565b613ceb565b83815290506020808201908360005b8681101561250f5781358601604080828b03121561241957600080fd5b80518181016001600160401b03828210818311171561243457fe5b81845284358181111561244657600080fd5b85016080818f03121561245857600080fd5b60c08401838110838211171561246a57fe5b855280358281111561247b57600080fd5b6124878f82840161274d565b845250888101358281111561249b57600080fd5b6124a78f828401612975565b9050606081818701526124bb878401612b5e565b60808701526124cb81840161295f565b60a0870152505050818352878501359350808411156124e957600080fd5b50506124f78b8385016126bc565b818701528652505092820192908201906001016123fc565b505050509392505050565b80356113de81613d91565b600082601f830112612535578081fd5b813560206125456123e883613d0e565b82815281810190858301855b858110156125ea5781358801608080601f19838d03011215612571578889fd5b604080518281016001600160401b03828210818311171561258e57fe5b908352848a0135825284830135828b01526060906125ad828701612b74565b838501529385013593808511156125c2578c8dfd5b506125d18e8b86880101612975565b9082015287525050509284019290840190600101612551565b5090979650505050505050565b600082601f830112612607578081fd5b815160206126176123e883613d0e565b82815281810190858301855b858110156125ea5781518801608080601f19838d03011215612643578889fd5b604080518281016001600160401b03828210818311171561266057fe5b908352848a0151825284830151828b015260609061267f828701612b7f565b83850152938501519380851115612694578c8dfd5b506126a38e8b868801016129c1565b9082015287525050509284019290840190600101612623565b600082601f8301126126cc578081fd5b813560206126dc6123e883613d0e565b828152818101908583016060808602880185018910156126fa578687fd5b865b868110156127205761270e8a84612b01565b855293850193918101916001016126fc565b509198975050505050505050565b600082601f83011261273e578081fd5b61149e838335602085016123da565b600082601f83011261275d578081fd5b8135602061276d6123e883613d0e565b82815281810190858301855b858110156125ea5781358801606080601f19838d03011215612799578889fd5b604080518281016001600160401b0382821081831117156127b657fe5b908352848a0135906127c782613d91565b9082528483013590808211156127db578c8dfd5b6127e98f8c84890101612975565b838c01529385013593808511156127fe578c8dfd5b505061280e8d8a85870101612525565b91810191909152865250509284019290840190600101612779565b600082601f830112612839578081fd5b815160206128496123e883613d0e565b82815281810190858301855b858110156125ea5781518801606080601f19838d03011215612875578889fd5b604080518281016001600160401b03828210818311171561289257fe5b908352848a0151906128a382613d91565b9082528483015190808211156128b7578c8dfd5b6128c58f8c848901016129c1565b838c01529385015193808511156128da578c8dfd5b50506128ea8d8a858701016125f7565b91810191909152865250509284019290840190600101612855565b600082601f830112612915578081fd5b813560206129256123e883613d0e565b8281528181019085830183850287018401881015612941578586fd5b855b858110156125ea57813584529284019290840190600101612943565b80356113de81613da6565b80516113de81613da6565b600082601f830112612985578081fd5b81356129936123e882613d2b565b8181528460208386010111156129a7578283fd5b816020850160208301379081016020019190915292915050565b600082601f8301126129d1578081fd5b81516129df6123e882613d2b565b8181528460208386010111156129f3578283fd5b611f8a826020830160208701613d65565b600060a08284031215612a15578081fd5b60405160a081016001600160401b038282108183111715612a3257fe5b8160405282935084358352602091508185013581811115612a5257600080fd5b85019050601f81018613612a6557600080fd5b8035612a736123e882613d0e565b81815283810190838501858402850186018a1015612a9057600080fd5b600094505b83851015612abc578035612aa881613d91565b835260019490940193918501918501612a95565b5080858701525050505050612ad360408401612b5e565b6040820152612ae46060840161251a565b6060820152612af560808401612b5e565b60808201525092915050565b600060608284031215612b12578081fd5b604051606081018181106001600160401b0382111715612b2e57fe5b6040529050808235612b3f81613dc8565b8082525060208301356020820152604083013560408201525092915050565b80356113de81613db4565b80516113de81613db4565b80356113de81613dc8565b80516113de81613dc8565b600060208284031215612b9b578081fd5b813561149e81613d91565b60008060408385031215612bb8578081fd5b8235612bc381613d91565b946020939093013593505050565b60008060008060808587031215612be6578182fd5b8435612bf181613d91565b966020860135965060408601359560600135945092505050565b600080600060608486031215612c1f578081fd5b83356001600160401b0380821115612c35578283fd5b612c4187838801612525565b94506020860135915080821115612c56578283fd5b50612c6386828701612525565b925050604084013590509250925092565b600060208284031215612c85578081fd5b81516001600160401b03811115612c9a578182fd5b611f8a84828501612829565b600060208284031215612cb7578081fd5b815161149e81613da6565b600060208284031215612cd3578081fd5b5035919050565b600080600060608486031215612cee578081fd5b8335925060208401356001600160401b03811115612d0a578182fd5b612c638682870161274d565b600060208284031215612d27578081fd5b81356001600160401b0380821115612d3d578283fd5b8184019150610100808387031215612d53578384fd5b612d5c81613ceb565b90508235815260208301356020820152604083013582811115612d7d578485fd5b612d8987828601612975565b604083015250606083013560608201526080830135608082015260a083013560a082015260c083013582811115612dbe578485fd5b612dca87828601612975565b60c08301525060e083013560e082015280935050505092915050565b600080600060408486031215612dfa578081fd5b83356001600160401b0380821115612e10578283fd5b9085019060a08288031215612e23578283fd5b90935060208501359080821115612e38578283fd5b818601915086601f830112612e4b578283fd5b813581811115612e59578384fd5b8760208083028501011115612e6c578384fd5b6020830194508093505050509250925092565b60008060408385031215612e91578182fd5b82356001600160401b0380821115612ea7578384fd5b612eb386838701612a04565b93506020850135915080821115612ec8578283fd5b50612ed58582860161272e565b9150509250929050565b600080600060a08486031215612ef3578081fd5b83356001600160401b0380821115612f09578283fd5b612f1587838801612a04565b94506020860135915080821115612f2a578283fd5b50612f378682870161272e565b925050612f478560408601612b01565b90509250925092565b600060408284031215612f61578081fd5b604051604081018181106001600160401b0382111715612f7d57fe5b604052825181526020928301519281019290925250919050565b600060208284031215612fa8578081fd5b81516001600160401b0380821115612fbe578283fd5b9083019060808286031215612fd1578283fd5b604051608081018181108382111715612fe657fe5b604052825182811115612ff7578485fd5b61300387828601612829565b825250602083015182811115613017578485fd5b613023878286016129c1565b60208301525061303560408401612b69565b60408201526130466060840161296a565b606082015295945050505050565b600080600060608486031215613068578081fd5b8335925060208401356001600160401b0380821115613085578283fd5b61309187838801612525565b935060408601359150808211156130a6578283fd5b506130b386828701612905565b9150509250925092565b600080600080600060a086880312156130d4578283fd5b853594506020860135935060408601356001600160401b03808211156130f8578485fd5b61310489838a01612975565b9450606088013593506080880135915080821115613120578283fd5b5061312d88828901612905565b9150509295509295909350565b6001600160a01b03169052565b60008284526020808501945082825b8581101561318457813561316981613d91565b6001600160a01b031687529582019590820190600101613156565b509495945050505050565b6000815180845260208085019450808401835b838110156131845781516001600160a01b0316875295820195908201906001016131a2565b6000815180845260208085018081965082840281019150828601855b8581101561323a578284038952815180518552858101518686015260408082015160ff169086015260609081015160809186018290529061322681870183613325565b9a87019a95505050908401906001016131e3565b5091979650505050505050565b6000815180845260208085018081965082840281019150828601855b8581101561323a578284038952815160408151818752613285828801826133b0565b92880151968801969096525098850198935090840190600101613263565b6000815180845260208085018081965082840281019150828601855b8581101561323a578284038952815180516001600160a01b03168552858101516060878701819052906132f482880182613325565b9150506040808301519250868203818801525061331181836131c7565b9a87019a95505050908401906001016132bf565b6000815180845261333d816020860160208601613d65565b601f01601f19169290920160200192915050565b600081518352602082015160a0602085015261337060a085018261318f565b60408481015165ffffffffffff908116918701919091526060808601516001600160a01b0316908701526080948501511693909401929092525090919050565b60008151608084526133c560808501826132a3565b9050602083015184820360208601526133de8282613325565b91505065ffffffffffff60408401511660408501526060830151151560608501528091505092915050565b65ffffffffffff169052565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b039384168152919092166020820152604081019190915260600190565b6001600160a01b03929092168252602082015260400190565b6001600160a01b039390931683526020830191909152604082015260600190565b60006020825261149e60208301846131c7565b6000608082526134ca60808301876131c7565b851515602084015282810360408401526134e481866131c7565b91505082606083015295945050505050565b60006020825261149e60208301846132a3565b90815260200190565b600086825260a0602083015261352b60a0830187613325565b828103604084015261353d81876132a3565b65ffffffffffff95909516606084015250509015156080909101529392505050565b90815260406020820181905260099082015268666f7263654d6f766560b81b606082015260800190565b93845260ff9290921660208401526040830152606082015260800190565b6020808252601c908201527f7475726e4e756d5265636f7264206e6f7420696e637265617365642e00000000604082015260600190565b60208082526017908201527f686f6c64696e6773203c20657870656374656448656c64000000000000000000604082015260600190565b60208082526013908201527218dbdd5b19081b9bdd08199a5b99081b19599d606a1b604082015260600190565b602080825260119082015270125b98dbdc9c9958dd0818da185a5b9259607a1b604082015260600190565b60208082526012908201527121b430b73732b6103334b730b634bd32b21760711b604082015260600190565b60208082526014908201527318dbdd5b19081b9bdd08199a5b99081c9a59da1d60621b604082015260600190565b60208082526018908201527f7475726e4e756d5265636f7264206465637265617365642e0000000000000000604082015260600190565b602080825260169082015275125b991a58d95cc81b5d5cdd081899481cdbdc9d195960521b604082015260600190565b6020808252601f908201527f4368616c6c656e676572206973206e6f742061207061727469636970616e7400604082015260600190565b602080825260169082015275086deead8c840dcdee840e8e4c2dce6cccae4408aa8960531b604082015260600190565b60208082526015908201527418dbdd5b19081b9bdd08199a5b99081d185c99d95d605a1b604082015260600190565b6020808252601d908201527f746172676574417373657420213d2067756172616e7465654173736574000000604082015260600190565b6020808252600f908201526e6e6f7420612067756172616e74656560881b604082015260600190565b6020808252601b908201527f63616e6e6f74207472616e7366657220612067756172616e7465650000000000604082015260600190565b6020808252601a908201527f7661726961626c6550617274206e6f7420746865206c6173742e000000000000604082015260600190565b602080825260119082015270496e76616c6964207369676e617475726560781b604082015260600190565b6020808252601590820152741a5b98dbdc9c9958dd08199a5b99d95c9c1c9a5b9d605a1b604082015260600190565b60208082526018908201527f436f756c64206e6f74206465706f736974204552433230730000000000000000604082015260600190565b6020808252601f908201527f4465706f73697420746f2065787465726e616c2064657374696e6174696f6e00604082015260600190565b60208082526016908201527521b430b73732b6103737ba103334b730b634bd32b21760511b604082015260600190565b6020808252601b908201527f686f6c64696e677320616c72656164792073756666696369656e740000000000604082015260600190565b6020808252601a908201527f6e6f7420612067756172616e74656520616c6c6f636174696f6e000000000000604082015260600190565b6020808252601f908201527f496e636f7272656374206d73672e76616c756520666f72206465706f73697400604082015260600190565b6020808252601d908201527f436f756c64206e6f7420726566756e64206578636573732066756e6473000000604082015260600190565b600060408252833560408301526020840135601e19853603018112613a89578182fd5b840180356001600160401b03811115613aa0578283fd5b602081023603861315613ab1578283fd5b60a06060850152613ac960e085018260208501613147565b915050613ad860408601612b5e565b613ae56080850182613409565b50613af26060860161251a565b613aff60a085018261313a565b50613b0c60808601612b5e565b613b1960c0850182613409565b5082810360208401526103e68185613247565b600060408252613b3f6040830185613351565b82810360208401526103e68185613247565b60006020825261149e60208301846133b0565b600086825260a06020830152613b7d60a083018761318f565b65ffffffffffff95861660408401526001600160a01b0394909416606083015250921660809092019190915292915050565b918252602082015260400190565b65ffffffffffff91909116815260200190565b600065ffffffffffff80881683526020818816818501526040915086151582850152606060a081860152613c0760a0860188613351565b858103608087015286518082528382019084810283018501858a01885b83811015613cac57858303601f19018552815180518a8552613c488b8601826133b0565b918a0151858303868c01528051808452908b01928d92508b01905b80831015613c97578351805160ff1683528c8101518d8401528d01518d830152928b019260019290920191908a0190613c63565b50968a01969450505090870190600101613c24565b50909e9d5050505050505050505050505050565b65ffffffffffff93841681529190921660208201526001600160a01b03909116604082015260600190565b6040518181016001600160401b0381118282101715613d0657fe5b604052919050565b60006001600160401b03821115613d2157fe5b5060209081020190565b60006001600160401b03821115613d3e57fe5b50601f01601f191660200190565b600061149e3684846123da565b60006114873683612a04565b60005b83811015613d80578181015183820152602001613d68565b83811115610c015750506000910152565b6001600160a01b038116811461161657600080fd5b801515811461161657600080fd5b65ffffffffffff8116811461161657600080fd5b60ff8116811461161657600080fdfea26469706673582212201444df6db57af9539080f29723636e2490b821ef439a5ba6774a7c4d9449b29164736f6c63430007060033",
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

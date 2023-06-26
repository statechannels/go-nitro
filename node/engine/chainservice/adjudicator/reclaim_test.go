package NitroAdjudicator

import (
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/go-cmp/cmp"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/types"
)

func TestComputeReclaimEffects(t *testing.T) {
	Alice := types.Destination(common.HexToHash("0xa"))
	Bob := types.Destination(common.HexToHash("0xb"))

	type TestCaseInputs struct {
		sourceAllocations     []outcome.Allocation
		targetAllocations     []outcome.Allocation
		indexOfTargetInSource uint
	}

	type TestCaseOutputs struct {
		newSourceAllocations []outcome.Allocation
	}

	type TestCase struct {
		inputs  TestCaseInputs
		outputs TestCaseOutputs
	}

	metadata, err := outcome.GuaranteeMetadata{Left: Alice, Right: Bob}.Encode()
	if err != nil {
		t.Fatal(err)
	}

	testCase1 := TestCase{
		inputs: TestCaseInputs{
			indexOfTargetInSource: 2,
			sourceAllocations: []outcome.Allocation{
				{
					Destination:    Alice,
					Amount:         big.NewInt(2),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Bob,
					Amount:         big.NewInt(2),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    [32]byte{},
					Amount:         big.NewInt(6),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       metadata,
				},
			},
			targetAllocations: []outcome.Allocation{
				{
					Destination:    Alice,
					Amount:         big.NewInt(1),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Bob,
					Amount:         big.NewInt(5),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
			},
		},
		outputs: TestCaseOutputs{
			newSourceAllocations: []outcome.Allocation{
				{
					Destination:    Alice,
					Amount:         big.NewInt(3),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Bob,
					Amount:         big.NewInt(7),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
			},
		},
	}

	offChainNewSourceAllocations, err := computeReclaimEffects(
		testCase1.inputs.sourceAllocations,
		testCase1.inputs.targetAllocations,
		testCase1.inputs.indexOfTargetInSource,
	)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(offChainNewSourceAllocations, testCase1.outputs.newSourceAllocations); diff != "" {
		t.Fatalf("newSourceAllocations does not match expectation :\n%s", diff)
	}
}

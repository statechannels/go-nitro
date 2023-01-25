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
	Irene := types.Destination(common.HexToHash("0x1"))
	Bob := types.Destination(common.HexToHash("0xb"))

	type TestCaseInputs struct {
		sourceAllocations     []outcome.Allocation
		targetAllocations     []outcome.Allocation
		indexOfTargetInSource uint
		releaseFees           bool
	}

	type TestCaseOutputs struct {
		newSourceAllocations []outcome.Allocation
	}

	type TestCase struct {
		inputs  TestCaseInputs
		outputs TestCaseOutputs
	}

	metadata, err := outcome.GuaranteeMetadata{Left: Alice, Right: Irene}.Encode()

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
					Destination:    Irene,
					Amount:         big.NewInt(2),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    [32]byte{},
					Amount:         big.NewInt(7),
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
					Destination:    Irene,
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
			releaseFees: true,
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
					Destination:    Irene,
					Amount:         big.NewInt(8),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
			},
		},
	}

	testCase2 := TestCase{
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
					Destination:    Irene,
					Amount:         big.NewInt(2),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    [32]byte{},
					Amount:         big.NewInt(7),
					AllocationType: outcome.GuaranteeAllocationType,
					Metadata:       metadata,
				},
			},
			targetAllocations: []outcome.Allocation{
				{
					Destination:    Alice,
					Amount:         big.NewInt(6),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Irene,
					Amount:         big.NewInt(1),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Bob,
					Amount:         big.NewInt(0),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
			},
			releaseFees: false,
		},
		outputs: TestCaseOutputs{
			newSourceAllocations: []outcome.Allocation{
				{
					Destination:    Alice,
					Amount:         big.NewInt(9),
					AllocationType: outcome.NormalAllocationType,
					Metadata:       []byte{},
				},
				{
					Destination:    Irene,
					Amount:         big.NewInt(2),
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
		testCase1.inputs.releaseFees,
	)
	if err != nil {
		t.Fatal(err)
	}

	for i, tc := range []TestCase{testCase1, testCase2} {
		if diff := cmp.Diff(offChainNewSourceAllocations, tc.outputs.newSourceAllocations); diff != "" {
			t.Fatalf("test case %v; newSourceAllocations does not match expectation :\n%s", i, diff)
		}

	}

}

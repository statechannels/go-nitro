package virtual_funding

import (
	"encoding/json"

	"github.com/statechannels/go-nitro/channel/state"
	"github.com/statechannels/go-nitro/channel/state/outcome"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

type Objective struct {
	status          string // TODO protocols.ObjectiveStatus
	jointPreFund    state.State
	preFundHash     types.Bytes32 // computed from jointPreFund (cached for performance)
	jointPostFund   state.State
	signedPreFund   []bool        // indexed by participant
	postFundHash    types.Bytes32 // computed from jointPreFund (cached for performance)
	signedPostFund  []bool        // indexed by participant
	jointChannelId  types.Bytes32
	ledgerChannelId types.Bytes32
}

type WaitingFor uint // Enumerable states of the objective

const (
	CompletePrefund WaitingFor = iota
	LedgerChannelCounterSignature
)

// Generic methods
func (o *Objective) Clone() Objective {
	j, _ := json.Marshal(o)
	var newObjective Objective
	_ = json.Unmarshal(j, &newObjective)
	return newObjective
}

func (o *Objective) Id() protocols.ObjectiveId {
	return protocols.ObjectiveId(`virtual-funding:` + o.jointChannelId.String())
}

func (o *Objective) Approve() Objective {
	newObjective := o.Clone()
	newObjective.status = `approved`
	return newObjective
}
func (o *Objective) Reject() Objective {
	newObjective := o.Clone()
	newObjective.status = `rejected`
	return newObjective
}

// Specific methods

func (o *Objective) Update(event protocols.ObjectiveEvent) Objective {
	newObjective := o.Clone()

	// process new signatures
	for hash, sig := range event.Sigs { // TODO extract into a helper function?
		signer, _ := state.RecoverEthereumMessageSigner(hash[:], sig)
		signerIndex := state.IndexOf(o.jointPreFund.Participants, signer)
		switch hash {
		case newObjective.preFundHash:
			newObjective.signedPreFund[signerIndex] = true
		case newObjective.postFundHash:
			newObjective.signedPostFund[signerIndex] = true
		}
	}

	// process holdings updates -- not relevant to this protocol

	// process adjudication updates -- not relevant to this protocol

	return newObjective
}

func (*Objective) Crank() (Objective, SideEffects, WaitingFor, error) {

}

type ConstructorParams struct {
	chainId              *types.Uint256
	appDefinition        types.Address
	appData              types.Bytes
	outcome              outcome.Exit
	myExitDestination    types.Bytes32
	theirExitDestination types.Bytes32
	mySigningAddress     types.Address
	theirSigningAddress  types.Address
	hubSigningAddress    types.Address // TODO generalize to >1 hub
	ChallengeDuration    *types.Uint256
}

// New makes a new virtual_funding Objective
func New(params ConstructorParams) (Objective, error) {
	o := Objective{}
	o.jointPreFund = state.State{
		ChainId:           params.chainId,
		Participants:      []types.Address{params.mySigningAddress, params.theirSigningAddress, params.hubSigningAddress},
		AppDefinition:     params.appDefinition,
		AppData:           params.appData,
		Outcome:           params.outcome,
		ChallengeDuration: params.ChallengeDuration,
		TurnNum:           state.PrefundTurnNum,
		IsFinal:           false,
	}

	var err error
	o.jointChannelId, err = o.jointPreFund.ChannelId()
	if err != nil {
		return Objective{}, err
	}
	return o, nil
}

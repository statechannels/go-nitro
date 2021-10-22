// Package Store contains the interface for a go-nitro store.
package store // import "github.com/statechannels/go-nitro/client/engine/store"

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/types"
)

// Store is responsible for persisting objectives, objective metadata, states, signatures, private keys and blockchain data
type Store interface {
	GetChannelSecretKey() *[]byte // Get a pointer to a secret key for signing channel updates

	GetObjectiveById(protocols.ObjectiveId) protocols.Objective // Read an existing objective
	GetObjectiveByChannelId(types.Bytes32) protocols.Objective  // Get the objective that currently owns the channel with the supplied ChannelId
	SetObjective(protocols.Objective) error                     // Write an objective

	UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) // updates progressLastMadeAt information for an objective
}

type TestStore struct{}

func (TestStore) GetChannelSecretKey() *[]byte {
	k := common.Hex2Bytes(`187bb12e927c1652377405f81d93ce948a593f7d66cfba383ee761858b05921a`)
	return &k
}

func (TestStore) GetObjectiveById(protocols.ObjectiveId) protocols.Objective {
	return protocols.TestObjective{}
}
func (TestStore) GetObjectiveByChannelId(types.Bytes32) protocols.Objective {
	return protocols.TestObjective{}
}

func (TestStore) SetObjective(protocols.Objective) error {
	return nil
}
func (TestStore) UpdateProgressLastMadeAt(protocols.ObjectiveId, protocols.WaitingFor) {}

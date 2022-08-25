package engine

import (
	"fmt"

	"github.com/statechannels/go-nitro/protocols"
	"github.com/statechannels/go-nitro/protocols/directdefund"
	"github.com/statechannels/go-nitro/protocols/directfund"
	"github.com/statechannels/go-nitro/protocols/virtualdefund"
	"github.com/statechannels/go-nitro/protocols/virtualfund"
)

// handleObjectiveRequest handles an ObjectiveRequest (triggered by a client API call).
// It will attempt to spawn a new, approved objective.
func (e *Engine) handleObjectiveRequest(or protocols.ObjectiveRequest) (EngineEvent, error) {
	defer e.metrics.RecordFunctionDuration()()

	if or == nil {
		panic("tried to handle nil objective request")
	}

	switch request := or.(type) {

	case virtualfund.ObjectiveRequest:
		e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
		vfo, err := virtualfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetConsensusChannel)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
		}

		err = e.registerPaymentChannel(vfo)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("could not register channel with payment/receipt manager: %w", err)
		}
		return e.attemptProgress(&vfo)

	case virtualdefund.ObjectiveRequest:
		e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
		vdfo, err := virtualdefund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetChannelById, e.store.GetConsensusChannel)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
		}
		return e.attemptProgress(&vdfo)

	case directfund.ObjectiveRequest:
		e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
		dfo, err := directfund.NewObjective(request, true, *e.store.GetAddress(), e.store.GetChannelsByParticipant, e.store.GetConsensusChannel)
		if err != nil {
			return EngineEvent{}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
		}
		return e.attemptProgress(&dfo)

	case directdefund.ObjectiveRequest:
		e.metrics.RecordObjectiveStarted(request.Id(*e.store.GetAddress()))
		ddfo, err := directdefund.NewObjective(request, true, e.store.GetConsensusChannelById)
		if err != nil {
			return EngineEvent{FailedObjectives: []protocols.ObjectiveId{request.Id(*e.store.GetAddress())}}, fmt.Errorf("handleAPIEvent: Could not create objective for %+v: %w", request, err)
		}
		// If ddfo creation was successful, destroy the consensus channel to prevent it being used (a Channel will now take over governance)
		e.store.DestroyConsensusChannel(request.ChannelId)
		return e.attemptProgress(&ddfo)

	default:
		return EngineEvent{}, fmt.Errorf("handleAPIEvent: Unknown objective type %T", request)
	}

}

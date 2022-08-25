package engine

import "github.com/statechannels/go-nitro/protocols"

// executeSideEffects executes the SideEffects declared by cranking an Objective or handling a payment request.
func (e *Engine) executeSideEffects(sideEffects protocols.SideEffects) error {
	defer e.metrics.RecordFunctionDuration()()

	for _, message := range sideEffects.MessagesToSend {

		e.logger.Printf("Sending message %+v", protocols.SummarizeMessage(message))
		e.msg.Send(message)
		e.metrics.RecordOutgoingMessage(message)
	}
	for _, tx := range sideEffects.TransactionsToSubmit {
		e.logger.Printf("Sending chain transaction for channel %s", tx.ChannelId())
		err := e.chain.SendTransaction(tx)
		if err != nil {
			return err
		}
	}
	for _, proposal := range sideEffects.ProposalsToProcess {
		e.fromLedger <- proposal
	}
	return nil
}

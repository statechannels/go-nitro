## 0001 -- Off-chain protocols as Flowcharts

## Status

Accepted

## Context

We require a consistent approach to encoding off-chain state channel protocols which:

- tames complexity
- reduces data duplication and synchronization issues
- allows for "restartability" of the protocol if it is interrupted before completing
- allows for easy testing

## Decision

An objective encodes a state chart (or state machine), but in a slightly non-standard manner:

- The **enumerable state** of the objective consists of several "pause points" where the current peer cannot progress without external input `X`. Hence we call theses states "Waiting For `X`".
- The **enumerable state** of the objective is not stored anywhere: it is computed from the **extended state**.
- The **extended state** of the objective (sometimes known as the "context") contains potentially infinite data such as off-chain messages, signatures and so on.
- State transitions are triggered by **events** in one of two ways:
  - new information comes to light and the extended state is **updated**
  - the objective is **cranked** -- driven towards completion -- and side effects are generated and returned. This may also update the extended state.
- The objective progress can be thought of as a flow chart (rather than a conventional state chart). When a blocking condition is identified, the state transition method call returns "early". As more progress is made, the execution reaches deeper into the method until eventually it reaches the end and the objective is complete.

Updating and cranking MUST be strictly pure functions, returning a fresh copy of the objective (not mutating its inputs).

## Consequences

After careful consideration, we discard the more conventional explicit finite state machine or ("state chart') model, although many of its essential elements remain.

A side-by-side comparison can be found in [this document](./flowchart-vs-statechart.md)

Complexity is tamed by the guiding principle of a flowchart "trying to make as much progress as possible before being blocked by events outside of its control".

Restartability is achieved "for free", since the flowchart actually restarts every time the objective is cranked. There is no "progress data" or finite state stored, other than some metadata derived from the extended state. This metadata may be useful for logging, debugging, statistics and so on -- but is not part of the core logic of the protocol.

For this reason it is imposisble for the finite state to get out of sync with the extended state -- for example you are `WaitingForCompletePrefund` when in fact you have a complete pre fund.

Unit testing is straightforward -- since `Update` and `Crank` are pure functions there is no need to spin up or mock a database, messaging service or chain service in order to get started unit testing.

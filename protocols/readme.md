This package currently contains two prototype implementations of the direct funding protocol.

They have been written to be as close as possible to show of the differences in approach.

`common.go` contains common data structures and methods. In particular:

- the _enumerable_ state of the protocol.
- the _extended_ state of the protocol (plus methods)
- the side effects
- the errors

The first approach is a state chart approach:

- the enumerable state and extended state are combined into (simply) a "state"
- a number of discrete event types are defined, as well as a rich event structure (roughly speaking the intersection of all the data for all of the event types)
- a `NexState` method contains the business logic. It is a "reducer" that switches on the current enumerable state and defers to one of a number of per-state reducers. It returns a new state, some side effects and an error
- it is accompanied by a mermaid diagram (manually constructed) in a comment

The second approach is a flow chart approach:

- there is an implicit mechanism for handling events and updating the extended state (not yet shown)
- a single `Crank` method is defined on the extended state (only)
- it does not take an event (or any other parameter)
- the enumerable state is computed by inspecting the extended state, and returned along with side effects and an error
- the flowchart is implemented by returning early from the execution when conditions are not (yet) met
- it is accompanied by a mermaid diagram (manually constructed) in a comment

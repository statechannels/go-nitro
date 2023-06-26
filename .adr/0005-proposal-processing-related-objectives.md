# 0005 -- Proposal Processing

## Status

Deprecated - see [ADR-0006](./0006-proposal-processing-ledger-effects.md)

## Definitions

- references to **virtualfundging** objectives are intended as generic references to either `virtualfund` or `virtualdefund`.

## Context

The go-nitro engine receives peer messages labelled with specific ObjectiveIDs.

Because:

- virtualfunding objectives alter the state of a node's ledger channel(s) by adding and removing guarantees
- alterations to the ledger channel state must occur in the correct order (see [ADR-0003](0003-consensus-ledger-channels.md))

The progress of virtualfunding objectives `A` and `B` can interfere with one another.

The engine needs a mechanism for re-cranking objectives `B`, `C`, `D`, ... which may or may not have become unblocked after progress on objective `A` alters the state of one or more ledger channels.

## Decision

The engine uses a `Related Objectives` strategy to restart progess on objectives which were blocked.

The workflow is:

- receive a peer message, labelled for a speficic ObjectiveID
- `Update()` this objective with any included data
- `Crank()` this objective
- pass the updated objective into `attemptProgressForRelatedObjectives()`

`attemptProgressForRelatedObjectives()` queries the objective's ledger channels for other objectives which are now potentially unblocked, and the process recurses on those objectives.

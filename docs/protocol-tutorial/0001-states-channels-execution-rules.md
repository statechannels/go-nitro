# States, Channels and Execution Rules

A state channel can be thought of as a set of data structures (called "states") committed to and exchanged between a fixed set of actors (which we call participants), together with some execution rules.

!!! info

    "Committing to" a state is typically done by digitially signing it.

## States

In Nitro protocol, a state has the following type (on chain in Solidity, off-chain in Typescript and Go):
=== "Solidity"

    ```solidity
    struct State {
        // participants sign the hash of this
        bytes32 channelId; // keccack(FixedPart)
        bytes appData;
        bytes outcome;
        uint48 turnNum;
        bool isFinal;
    }
    ```

=== "TypeScript"

    ``` ts
    import {Channel, Outcome, State} from '@statechannels/nitro-protocol';

    const state: State = {
        turnNum: 0,
        isFinal: false,
        channel,
        challengeDuration,
        outcome,
        appDefinition,
        appData
    };

    ```

=== "Go"

    ``` Go
    import (
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/statechannels/go-nitro/channel/state"
    "github.com/statechannels/go-nitro/channel/state/outcome"
    "github.com/statechannels/go-nitro/internal/testactors"
    "github.com/statechannels/go-nitro/types"
    )

    var testState = state.State{
        ChainId: chainId,
        Participants: []types.Address{
            testactors.Alice.Address(),
            testactors.Bob.Address(),
            },
        ChannelNonce: big.NewInt(37140676580),
        AppDefinition: someAppDefinition,
        ChallengeDuration: big.NewInt(60),
        AppData: []byte{},
        Outcome: testOutcome,
        TurnNum: 5,
        IsFinal: false,
    }

    ```

## Execution Rules

The rules dictate the conditions under which a state may be considered **supported** by the underlying blockchain, and also dictate how one supported state may supercede another. In this manner, state channels may be "updated" as participants follow the rules to support state after state.

If a state is supported by the underylying blockchain, it has a chance to be the **final state** for the channel. The final state influences how any assets locked into the channel will be dispersed.

Unlike the rules of the underlying blockchain -- which dictate which state history is canoncial via Proof of Work, Proof of Stake, Proof of Authority (or some other hardcoded mechanism) -- Nitro protocol allows for the state channel update rules to vary from one application to another. One state channel might proceed only by unanimous consensus (all parties must digitially sign a state to make it supported), and another might proceed in a round-robin fashion (each party has the chance to support the next state unilaterally).

The rules for how one supported state may supercede another are very simple. Each state has a version number, with greater version numbers superceding lesser ones.

The state channel rules are enshrined in two places on the blockchain: firstly, in the **core protocol**, and secondly in the **application rules**.

### Core protocol rules

Nitro is a very open protocol, and has become more open over time (in particular, v2 is more open compared with v1). This means that very little is stipulated at the core protocol level. Each application gets full control over when a state can be considered supported. The only things enforced by the core protocol are:

- the rule that higher turn numbers take precedence over lower ones
- an escape hatch for an "instant checkout" of the channel, which bypasses the application rules altogether (TODO link)

Otherwise, the core protocol defers to the application rules.

### Application rules

Each channel is required to specify application rules in a contract adhering to the following on chain interface:

```solidity

/**
 * @dev The IForceMoveApp interface calls for its children to implement an application-specific requireStateSupported function, defining the state machine of a ForceMove state channel DApp.
 */
interface IForceMoveApp is INitroTypes {
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied.
     * @dev Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert when invalid support proof and a candidate are supplied.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate. May be omitted when `candidate` constitutes a support proof itself.
     * @param candidate Recovered variable part the proof was supplied for. Also may constitute a support proof itself.
     */
    function requireStateSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external pure;
}

```

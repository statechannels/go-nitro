---
description: How to program a state channel to fit your application.
---

# Execution Rules

A channel's execution rules dictate the conditions under which a state may be considered **supported** by the underlying blockchain, and also dictate how one supported state may supercede another. In this manner, state channels may be "updated" as participants follow the rules to support state after state.

If a state is supported by the underylying blockchain, it has a chance to be the **final state** for the channel. The final state influences how any assets locked into the channel will be dispersed.

Unlike the rules of the underlying blockchain -- which dictate which state history is canoncial via Proof of Work, Proof of Stake, Proof of Authority (or some other hardcoded mechanism) -- Nitro protocol allows for the state channel update rules to vary from one application to another. One state channel might proceed only by unanimous consensus (all parties must digitially sign a state to make it supported), and another might proceed in a round-robin fashion (each party has the chance to support the next state unilaterally).

The rules for how one supported state may supercede another are very simple. Each state has a version number, with greater version numbers superceding lesser ones.

The state channel rules are enshrined in two places on the blockchain: firstly, in the **core protocol**, and secondly in the **application rules**.

Participants _may_ provide "support proofs" to the blockchain in order to record the execution of the state channel. They will typically prefer to not do that, since it costs gas -- but they will keep such "support proofs" in hand in case they need to use them.

### Core protocol rules

Nitro is a very open protocol: This means that very little is stipulated at the core protocol level. Each application gets full control over when a state can be considered supported. The only things enforced by the core protocol are:

- the rule that higher turn numbers take precedence over lower ones
- an escape hatch for an ["instant checkout"](./0070-finalizing-a-channel.md#happy-path) of the channel, which bypasses the application rules altogether

Otherwise, the core protocol defers to the application rules.

### Application rules

Each channel is required to specify application rules in a contract adhering to the following on chain interface:

```solidity

/**
 * @dev The IForceMoveApp interface calls for its children to implement an application-specific stateIsSupported function, defining the state machine of a ForceMove state channel DApp.
 */
interface IForceMoveApp is INitroTypes {
    /**
     * @notice Encodes application-specific rules for a particular ForceMove-compliant state channel. Must revert or return false when invalid support proof and a candidate are supplied.
     * @dev Depending on the application, it might be desirable to narrow the state mutability of an implementation to 'pure' to make security analysis easier.
     * @param fixedPart Fixed part of the state channel.
     * @param proof Array of recovered variable parts which constitutes a support proof for the candidate. May be omitted when `candidate` constitutes a support proof itself.
     * @param candidate Recovered variable part the proof was supplied for. Also may constitute a support proof itself.
     */
    function stateIsSupported(
        FixedPart calldata fixedPart,
        RecoveredVariablePart[] calldata proof,
        RecoveredVariablePart calldata candidate
    ) external view returns (bool, string memory);
}

```

!!! info

    Although the above interface allows for a `view` function, we recommend that you use a `pure` function wherever possible. Doing so makes the execution rules easier to reason about and verify.

### Auxiliary application rules

Some of the more advanced features in Nitro are actually expressed themselves as Nitro Applications which we call _auxiliary applications_. There are a couple of important ones:

#### `ConsensusApp`

The consensus app encodes a very simple rule for execution -- in order for a state to be supported, it must be _unanimously countersigned_ -- that is, signed by _all_ of the channel participants. See the [source code](https://github.com/statechannels/go-nitro/blob/main/nitro-protocol/contracts/ConsensusApp.sol).

Ledger channels are a special type of channel used to fund other channels -- they are an example of a channel which run the `ConsensusApp`.

#### `VirtualPaymentApp`

The virtual payment app allows a _payer_ to pay a _payee_ via their inirection connection through `n` intermediaries. Payments are simply signed "vouchers" sent from the _payer_ to the _payee_. This app is in effect a mini state channel adjudicator, which requires unanimous consensu for most state execution, but parses vouchers and allows for other transitions via _forced transtiions_ (or unilateral consensus). See the [source code](https://github.com/statechannels/go-nitro/blob/main/nitro-protocol/contracts/VirtualPaymentApp.sol).

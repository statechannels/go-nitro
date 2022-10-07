# States, Channels and Execution Rules

A state channel can be thought of as a set of data structures (called "states") committed to and exchanged between a fixed set of actors (which we call participants), together with some execution rules.

!!! info

    In Nitro, participants "commit to" a state by digitially signing it.

## States

In Nitro protocol, a state is broken up into fixed and variable parts:

=== "Solidity"

    ```solidity
    import {ExitFormat as Outcome} from '@statechannels/exit-format/contracts/ExitFormat.sol';

    struct FixedPart {
        uint256 chainId;
        address[] participants;
        uint48 channelNonce;
        address appDefinition;
        uint48 challengeDuration;
    }

    struct VariablePart {
        Outcome.SingleAssetExit[] outcome; // (1)
        bytes appData;
        uint48 turnNum;
        bool isFinal;
    }
    ```

    1. This composite type is explained in the section on [outcomes](./0002-outcomes.md).

=== "TypeScript"

    ```typescript
        import * as ExitFormat from '@statechannels/exit-format';
        import {Address, Bytes, Bytes32, Uint256, Uint48, Uint64} from '@statechannels/nitro-protocol'; // (1)

        export interface FixedPart {
            chainId: Uint256;
            participants: Address[];
            channelNonce: Uint64;
            appDefinition: Address;
            challengeDuration: Uint48;
        }

        export interface VariablePart {
            outcome: ExitFormat.Exit; // (2)
            appData: Bytes;
            turnNum: Uint48;
            isFinal: boolean;
        }
    ```

    1. `Bytes32`, `Bytes`, `Address`, `Uint256`, `Uint64` are aliases to the Javascript `string` type. They are respresented as hex strings. `Uint48` is aliased to a `number`.
    2. This composite type is explained in the section on [outcomes](./0002-outcomes.md).

=== "Go"

    ```Go
    import (
        "github.com/statechannels/go-nitro/channel/state/outcome"
        "github.com/statechannels/go-nitro/types" // (1)
    )

    type (
        FixedPart struct {
            ChainId           *types.Uint256
            Participants      []types.Address
            ChannelNonce      uint64
            AppDefinition     types.Address
            ChallengeDuration uint32
        }

        VariablePart struct {
            AppData types.Bytes
            Outcome outcome.Exit // (2)
            TurnNum uint64
            IsFinal bool
        }
    )
    ```

    1. `types.Address` is an alias to go-ethereum's [`common.Address`](https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.8/common#Address) type. `types.Bytes32` is an alias to go-ethereum's [`common.Hash`](https://pkg.go.dev/github.com/ethereum/go-ethereum@v1.10.8/common#Hash) type.
    2. This composite type is explained in the section on [outcomes](./0002-outcomes.md).

!!! info

    States are usually submitted to the blockchain as a single fixed part and multiple variable parts.

Let's take each property in turn:

### Chain id

This needs to match the id of the chain where assets are to be locked (i.e. the 'root' of the funding graph for this channel). In the event of a mismatch, the channel cannot be concluded and funds cannot be unlocked.

### Participants

This is a dynamic array of Ethereum addresses, each derived from an ECDSA private key in the usual manner. Each address represents a participant in the state channel who is able to commit to state updates and thereby cause the channel to finalize on chain.

!!! warning

    Before joining a state channel, you (or your off-chain software) should check that it has length at least 2, but no more than 255, and include a public key (account) that you control. Each entry should be a nonzero ethereum address.

### ChannelNonce

This is a unique number used to differentiate channels with an otherwise identical `FixedPart`. For example, if the same participants want to run the same kind of channel on the same chain as a previous channel, they can choose a new `ChannelNonce` to prevent state updates from from the existing channel being replayed.

!!! warning

    You should never join a channel which re-uses a channel nonce.

### AppDefinition

This is an Ethereum address where a Nitro application has been deployed. This is a contract conforming to the `ForceMoveApp` and defining [application rules](#application-rules).

!!! warning

    You should have confidence that the application is not malicious or suffering from security flaws. You should inspect the source code (which should be publically available and verifiable) or appeal to a trusted authority to do this.

### ChallengeDuration

This is duration (in seconds) of the challenge-response window. If a challenge is raised on chain at time `t`, the channel will finalize at `t + ChallengeDuration` unless cleared by a subqsequent on-chain transaction.

!!! warning

    This should be at least 1 block time (~15 seconds on mainnet) and less than `2^48-1` seconds. Whatever it is set to, the channel should be closed long before `2^48 - 1 - challengeDuration`. In practice we recommend somewhere between 5 minutes and 5 months.

### AppData

The AppData is optional data which may be interpreted by the Nitro application and affect the execution rules of the channel -- see the section on [application rules](#application-rules). For example, it could describe the state of a chess board or include the hash of a secret.

### Outcome

This describes how funds will be disbursed if the channel were to finalize in the current state. See the section on [Outcomes](0002-outcomes.md).

### TurnNum

The turn number is the mechanism by which newer states take precedence over older ones. The turn number usually increments as the channel progresses.

!!! warning

    The turn number must not exceed 281,474,976,710,655 because then it will overflow on chain. It should not exceed 4,294,967,295 because it may then overflow off-chain. It is very unlikely a channel would ever have this many updates.

### IsFinal

This is a boolean flag which allows the [channel execution rules](#execution-rules) to be bypassed and for the channel to be finalized "instantly" without waiting for the challenge-response window to lapse.

!!! warning

    As soon as an `isFinal=true` state is _enabled_ (that is to say, you cannot prevent it from becoming supported) it is not safe to continue executing the state channel. It should be finalized immediately.

## Channel IDs

Channels are identified by the hash of the `FixedPart`` of the state (those parts that may _not_ vary):

```solidity

  bytes32 channelId = keccak256(
      abi.encode(
          fixedPart.chainId,
          fixedPart.participants,
          fixedPart.channelNonce,
          fixedPart.appDefinition,
          fixedPart.challengeDuration
      )
  );

```

## State commitments

To commit to a state, a hash is formed as follows:

```solidity
 bytes32 stateHash = keccak256(abi.encode(
        channelId,
        variablePart.appData,
        variablePart.outcome,
        variablePart.turnNum,
        variablePart.isFinal
    ));
```

and this hash is signed using an _ephemeral_ Ethereum private key. _Ephemeral_ in this context means a dedicated private key, generated solely for the purpose of executing the state channel.

## Execution Rules

The rules dictate the conditions under which a state may be considered **supported** by the underlying blockchain, and also dictate how one supported state may supercede another. In this manner, state channels may be "updated" as participants follow the rules to support state after state.

If a state is supported by the underylying blockchain, it has a chance to be the **final state** for the channel. The final state influences how any assets locked into the channel will be dispersed.

Unlike the rules of the underlying blockchain -- which dictate which state history is canoncial via Proof of Work, Proof of Stake, Proof of Authority (or some other hardcoded mechanism) -- Nitro protocol allows for the state channel update rules to vary from one application to another. One state channel might proceed only by unanimous consensus (all parties must digitially sign a state to make it supported), and another might proceed in a round-robin fashion (each party has the chance to support the next state unilaterally).

The rules for how one supported state may supercede another are very simple. Each state has a version number, with greater version numbers superceding lesser ones.

The state channel rules are enshrined in two places on the blockchain: firstly, in the **core protocol**, and secondly in the **application rules**.

### Core protocol rules

Nitro is a very open protocol, and has become more open over time (in particular, v2 is more open compared with v1). This means that very little is stipulated at the core protocol level. Each application gets full control over when a state can be considered supported. The only things enforced by the core protocol are:

- the rule that higher turn numbers take precedence over lower ones
- an escape hatch for an ["instant checkout"](./0006-finalizing-a-channel.md#happy-path) of the channel, which bypasses the application rules altogether

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

### Auxiliary application rules

Some of the more advanced features in Nitro are actually expressed themselves as Nitro Applications which we call _auxiliary applications_. There are a couple of important ones:

#### `ConsensusApp`

The consensus app encodes a very simple rule for execution -- in order for a state to be supported, it must be _unanimously countersigned_ -- that is, signed by _all_ of the channel participants. See the [source code](https://github.com/statechannels/go-nitro/blob/main/nitro-protocol/contracts/ConsensusApp.sol).

Ledger channels are a special type of channel used to fund other channels -- they are an example of a channel which run the `ConsensusApp`.

#### `VirtualPaymentApp`

The virtual payment app allows a _payer_ to pay a _payee_ via their inirection connection through `n` intermediaries. Payments are simply signed "vouchers" sent from the _payer_ to the _payee_. This app is in effect a mini state channel adjudicator, which requires unanimous consensu for most state execution, but parses vouchers and allows for other transitions via _forced transtiions_ (or unilateral consensus). See the [source code](https://github.com/statechannels/go-nitro/blob/main/nitro-protocol/contracts/VirtualPaymentApp.sol).

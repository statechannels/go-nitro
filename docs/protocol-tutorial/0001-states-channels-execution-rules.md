# States, Channels and Execution Rules

A state channel can be thought of as a set of data structures (called "states") committed to and exchanged between a fixed set of actors (which we call participants), together with some execution rules.

!!! info

    In Nitro, "committing to" a state menas digitially signing it.

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

## Channel IDs

Channels are identified by the hash of the `FixedPart`` of the state (those parts that may _not_ vary):

```solidity

  struct FixedPart {
      uint256 chainId;
      address[] participants;
      uint48 channelNonce;
      address appDefinition;
      uint48 challengeDuration;
  }

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

The remainding fields of the state may vary, and are known as the `VariablePart`:

```solidity
       struct VariablePart {
        Outcome.SingleAssetExit[] outcome;
        bytes appData;
        uint48 turnNum;
        bool isFinal;
    }
```

## State commitments

To commit to a state, a hash is formed as follows:

```solidity
 bytes32 stateHash = keccak256(abi.encode(
        channelId,
        vp.appData,
        vp.outcome,
        vp.turnNum,
        vp.isFinal
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

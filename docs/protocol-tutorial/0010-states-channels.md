# States & Channels

A state channel can be thought of as a set of data structures (called "states") committed to and exchanged between a fixed set of actors (which we call participants), together with some execution rules.

!!! info

    In Nitro, participants "commit to" a state by digitially signing it.

A state channel controls [funds](./0050-funding-a-channel.md) which are locked up -- either on an L1 blockchain or on some other ledger such as another state channel.

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

    1. This composite type is explained in the section on [outcomes](./0030-outcomes.md).

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
    2. This composite type is explained in the section on [outcomes](./0030-outcomes.md).

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
    2. This composite type is explained in the section on [outcomes](./0030-outcomes.md).

!!! info

    States are usually submitted to the blockchain as a single fixed part and multiple variable parts. This is known as a "support proof".

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

This describes how funds will be disbursed if the channel were to finalize in the current state. See the section on [Outcomes](./0030-outcomes.md).

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

# Finalizing a channel

## Happy path

Finalization of a state channel is a necessary step before defunding it. In the so-called 'happy' case, all participants cooperate to achieve this.

A participant wishing to end the state channel will sign a state with `isFinal = true`. Then, the other participants may support that state. Once a full set of `n` such signatures exists (this set is known as a **finalization proof**) the channel is said to be 'closed' or 'finalized'.

### Off chain

In most cases, the channel would be finalized and defunded [off chain](./off-chain-funding), and no contract calls are necessary.

### On chain -- calling `conclude`

In the case where assets were deposited against the channel on chain (the case of direct funding), anyone in possession of a finalization proof may use it to finalize the `outcome` on-chain. They would do this by calling `conclude` on the adjudicator. This enables [assets to be released](./release-assets).

The conclude method allows anyone with sufficient off-chain state to immediately finalize an outcome for a channel without having to wait for a challenge to expire (more on that later).

The off-chain state(s) is submitted (in an optimized format), and once relevant checks have passed, an expired challenge is stored against the `channelId`. (This is an implementation detail -- the important point is that the chain shows that the channel has been finalized.)

In the following example the participants support the state by countersigning it, without increasing the turn number:

```typescript
// In lesson6.test.ts

/* Import ethereum wallet utilities  */
import { BigNumber, ethers } from "ethers";
const { AddressZero, HashZero } = ethers.constants;

/* Import statechannels wallet utilities  */
import {
  Channel,
  State,
  getFixedPart,
  hashOutcome,
  signStates,
  hashAppPart,
} from "@statechannels/nitro-protocol";

/* Construct a final state */
const participants = [];
const wallets: ethers.Wallet[] = [];
for (let i = 0; i < 3; i++) {
  wallets[i] = ethers.Wallet.createRandom();
  participants[i] = wallets[i].address;
}
const chainId = "0x1234";
const channelNonce = BigNumber.from(0).toHexString();
const channel: Channel = { chainId, channelNonce, participants };
const largestTurnNum = 4;
const state: State = {
  isFinal: true,
  channel,
  outcome: [],
  appDefinition: AddressZero,
  appData: HashZero,
  challengeDuration: 86400, // 1 day
  turnNum: largestTurnNum,
};

/* Generate a finalization proof */
const whoSignedWhat = [0, 0, 0];
const sigs = await signStates([state], wallets, whoSignedWhat);

/*
  Call conclude
*/
const numStates = 1;
const fixedPart = getFixedPart(state);
const appPartHash = hashAppPart(state);
const outcomeHash = hashOutcome(state.outcome);
const tx = NitroAdjudicator.conclude(
  largestTurnNum,
  fixedPart,
  appPartHash,
  outcomeHash,
  numStates,
  whoSignedWhat,
  sigs
);
```

Notice we imported `hashOutcome` and `hashAppPart` in order to provide the `conclude` method with the correct calldata.

## Sad path

When cooperation breaks down, it is possible to finalize a state channel without requiring on-demand cooperation of all `n` participants. This is the so-called 'sad' path to finalizing a channel, and it requires a supported (but not necessarily `isFinal`) state(s) being submitted to the chain.

The `challenge` function allows anyone holding the appropriate off-chain state(s) to register a _challenge state_ on chain. It is designed to ensure that a state channel can progress or be finalized in the event of inactivity on behalf of a participant (e.g. the current mover).

The required data for this method consists of a single state, along with `n` signatures. Once these are submitted (in an optimized format), and once relevant checks have passed, an `outcome` is registered against the `channelId`, with a finalization time set at some delay after the transaction is processed.

This delay allows the challenge to be cleared by a timely and well-formed [respond](./clear-a-challenge#call-respond) or [checkpoint](./clear-a-challenge#call-checkpoint) transaction. We'll get to those shortly. If no such transaction is forthcoming, the challenge will time out, allowing the `outcome` registered to be finalized. A finalized outcome can then be used to extract funds from the channel (more on that below, too).

!!! tip

    The `challengeDuration` is a [fixed parameter](./execute-state-transitions#construct-a-state-with-the-correct-format) expressed in seconds, that is set when a channel is proposed. It should be set to a value low enough that participants may close a channel in a reasonable amount of time in the case that a counterparty becomes unresponsive; but high enough that malicious challenges can be detected and responded to without placing unreasonable liveness requirements on responders. A `challengeDuration` of 1 day is a reasonable starting point, but the "right" value will likely depend on your application.

### Call `challenge`

!!!note

    The challenger needs to sign this data:

    ```
    keccak256(abi.encode(challengeStateHash, 'forceMove'))
    ```

    in order to form `challengerSig`. This signals their intent to challenge this channel with this particular state. This mechanism allows the challenge to be authorized only by a channel participant.


    We provide a handy utility function `signChallengeMessage` to form this signature.

```typescript
// In lesson7.test.ts

import { signChallengeMessage } from "@statechannels/nitro-protocol";

const participants = [];
const wallets: ethers.Wallet[] = [];
for (let i = 0; i < 3; i++) {
  wallets[i] = ethers.Wallet.createRandom();
  participants[i] = wallets[i].address;
}
const chainId = "0x1234";
const channelNonce = 0;
const channel: Channel = { chainId, channelNonce, participants };

/* Choose a challenger */
const challenger = wallets[0];

/* Construct a progression of states */
const largestTurnNum = 8;
const isFinalCount = 0;
const appDatas = [0, 1, 2];
const states: State[] = appDatas.map((data, idx) => ({
  turnNum: largestTurnNum - appDatas.length + 1 + idx,
  isFinal: idx > appDatas.length - isFinalCount,
  channel,
  challengeDuration: 86400, // 1 day
  outcome: [],
  appDefinition: process.env.TRIVIAL_APP_ADDRESS,
  appData: HashZero,
}));

/* Construct a support proof */
const whoSignedWhat = [0, 1, 2];
const signatures = await signStates(states, wallets, whoSignedWhat);

/* Form the challengeSignature */
const challengeSignedState: SignedState = signState(
  states[states.length - 1],
  challenger.privateKey
);
const challengeSignature = signChallengeMessage(
  [challengeSignedState],
  challenger.privateKey
);

/* Submit the challenge transaction */
const variableParts = states.map((state) => getVariablePart(state));
const fixedPart = getFixedPart(states[0]);

const tx = NitroAdjudicator.challenge(
  fixedPart,
  largestTurnNum,
  variableParts,
  isFinalCount,
  signatures,
  whoSignedWhat,
  challengeSignature
);
```

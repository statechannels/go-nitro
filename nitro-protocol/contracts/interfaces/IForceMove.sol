// SPDX-License-Identifier: MIT
pragma solidity 0.7.6;
pragma experimental ABIEncoderV2;

import './IForceMoveApp.sol';

/**
 * @dev The IForceMove interface defines the interface that an implementation of ForceMove should implement. ForceMove protocol allows state channels to be adjudicated and finalized.
 */
interface IForceMove {
    struct Signature {
        uint8 v;
        bytes32 r;
        bytes32 s;
    }

    struct FixedPart {
        uint256 chainId;
        address[] participants;
        uint48 channelNonce;
        address appDefinition;
        uint48 challengeDuration;
    }

    struct State {
        // participants sign the hash of this
        bytes32 channelId; // keccack(chainId,participants,channelNonce,appDefinition,challengeDuration)
        bytes appData;
        bytes outcome;
        uint48 turnNum;
        bool isFinal;
    }

    /**
     * @notice Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @dev Registers a challenge against a state channel. A challenge will either prompt another participant into clearing the challenge (via one of the other methods), or cause the channel to finalize at a specific time.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param variableParts An ordered array of structs, each decribing the properties of the state channel that may change with each state update. Length is from 1 to the number of participants (inclusive).
     * @param sigs An array of signatures that support the state with the `largestTurnNum`. There must be one for each participant, e.g.: [sig-from-p0, sig-from-p1, ...]
     * @param whoSignedWhat An array denoting which participant has signed which state: `participant[i]` signed the state with index `whoSignedWhat[i]`.
     * @param challengerSig The signature of a participant on the keccak256 of the abi.encode of (supportedStateHash, 'forceMove').
     */
    function challenge(
        FixedPart memory fixedPart,
        IForceMoveApp.VariablePart[] memory variableParts,
        Signature[] memory sigs,
        uint8[] memory whoSignedWhat,
        Signature memory challengerSig
    ) external;

    /**
     * @notice Repsonds to an ongoing challenge registered against a state channel.
     * @dev Repsonds to an ongoing challenge registered against a state channel.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param variablePartAB An pair of structs, each decribing the properties of the state channel that may change with each state update (for the challenge state and for the response state).
     * @param sig The responder's signature on the `responseStateHash`.
     */
    function respond(
        FixedPart memory fixedPart,
        IForceMoveApp.VariablePart[2] memory variablePartAB,
        // variablePartAB[0] = challengeVariablePart
        // variablePartAB[1] = responseVariablePart
        Signature memory sig
    ) external;

    /**
     * @notice Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @dev Overwrites the `turnNumRecord` stored against a channel by providing a state with higher turn number, supported by a signature from each participant.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param variableParts An ordered array of structs, each decribing the properties of the state channel that may change with each state update.
     * @param sigs An array of signatures that support the state with the `largestTurnNum`: one for each participant, in participant order (e.g. [sig of participant[0], sig of participant[1], ...]).
     * @param whoSignedWhat An array denoting which participant has signed which state: `participant[i]` signed the state with index `whoSignedWhat[i]`.
     */
    function checkpoint(
        FixedPart memory fixedPart,
        IForceMoveApp.VariablePart[] memory variableParts,
        Signature[] memory sigs,
        uint8[] memory whoSignedWhat
    ) external;

    /**
     * @notice Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @dev Finalizes a channel by providing a finalization proof. External wrapper for _conclude.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param latestVariablePart Latest variable part in finalization proof. Must have the largest turnNum and the same appData and outcome as all other variable parts in finalization proof.
     * @param numStates The number of states in the finalization proof.
     * @param whoSignedWhat An array denoting which participant has signed which state: `participant[i]` signed the state with index `whoSignedWhat[i]`.
     * @param sigs An array of signatures that support the state with the `largestTurnNum`: one for each participant, in participant order (e.g. [sig of participant[0], sig of participant[1], ...]).
     */
    function conclude(
        FixedPart memory fixedPart,
        IForceMoveApp.VariablePart memory latestVariablePart,
        uint8 numStates,
        uint8[] memory whoSignedWhat,
        Signature[] memory sigs
    ) external;

    // events

    /**
     * @dev Indicates that a challenge has been registered against `channelId`.
     * @param channelId Unique identifier for a state channel.
     * @param turnNumRecord A turnNum that (the adjudicator knows) is supported by a signature from each participant.
     * @param finalizesAt The unix timestamp when `channelId` will finalize.
     * @param isFinal Boolean denoting whether the challenge state is final.
     * @param fixedPart Data describing properties of the state channel that do not change with state updates.
     * @param variableParts An ordered array of structs, each decribing the properties of the state channel that may change with each state update.
     * @param sigs A list of Signatures that supported the challenge: one for each participant, in participant order (e.g. [sig of participant[0], sig of participant[1], ...]).
     * @param whoSignedWhat Indexing information to identify which signature was by which participant
     */
    event ChallengeRegistered(
        bytes32 indexed channelId,
        uint48 turnNumRecord,
        uint48 finalizesAt,
        bool isFinal,
        FixedPart fixedPart,
        IForceMoveApp.VariablePart[] variableParts,
        Signature[] sigs,
        uint8[] whoSignedWhat
    );

    /**
     * @dev Indicates that a challenge, previously registered against `channelId`, has been cleared.
     * @param channelId Unique identifier for a state channel.
     * @param newTurnNumRecord A turnNum that (the adjudicator knows) is supported by a signature from each participant.
     */
    event ChallengeCleared(bytes32 indexed channelId, uint48 newTurnNumRecord);

    /**
     * @dev Indicates that a challenge has been registered against `channelId`.
     * @param channelId Unique identifier for a state channel.
     * @param finalizesAt The unix timestamp when `channelId` finalized.
     */
    event Concluded(bytes32 indexed channelId, uint48 finalizesAt);
}
